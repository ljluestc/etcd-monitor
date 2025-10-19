package monitor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAlertManager_AddChannel(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	thresholds := AlertThresholds{
		MaxLatencyMs:            100,
		MaxDatabaseSizeMB:       8192,
		MinAvailableNodes:       2,
		MaxLeaderChangesPerHour: 3,
	}

	am := NewAlertManager(thresholds, logger)

	t.Run("Add console channel", func(t *testing.T) {
		channel := NewConsoleChannel(logger)
		am.AddChannel(channel)

		// Verify channel was added by triggering an alert
		alert := Alert{
			Level:     AlertLevelInfo,
			Type:      AlertTypeClusterHealth,
			Message:   "Test alert for channel",
			Timestamp: time.Now(),
		}

		// This should send to the console channel
		am.TriggerAlert(alert)

		// Check that alert was recorded
		history := am.GetAlertHistory()
		assert.Greater(t, len(history), 0)
	})

	t.Run("Add multiple channels", func(t *testing.T) {
		channel1 := NewConsoleChannel(logger)
		channel2 := NewConsoleChannel(logger)

		am2 := NewAlertManager(thresholds, logger)
		am2.AddChannel(channel1)
		am2.AddChannel(channel2)

		alert := Alert{
			Level:     AlertLevelWarning,
			Type:      AlertTypeHighLatency,
			Message:   "Test alert for multiple channels",
			Timestamp: time.Now(),
		}

		// This should send to both channels
		am2.TriggerAlert(alert)

		history := am2.GetAlertHistory()
		assert.Greater(t, len(history), 0)
	})
}

func TestEmailChannel_Structure(t *testing.T) {
	t.Run("Email channel creation", func(t *testing.T) {
		ec := &EmailChannel{
			SMTPServer: "smtp.example.com",
			SMTPPort:   587,
			From:       "monitor@example.com",
			To:         []string{"admin@example.com"},
			Username:   "user",
			Password:   "pass",
		}

		assert.Equal(t, "smtp.example.com", ec.SMTPServer)
		assert.Equal(t, 587, ec.SMTPPort)
		assert.Equal(t, "monitor@example.com", ec.From)
		assert.Equal(t, []string{"admin@example.com"}, ec.To)
	})

	t.Run("Email channel name", func(t *testing.T) {
		ec := &EmailChannel{}
		assert.Equal(t, "email", ec.Name())
	})

	t.Run("Email channel send with invalid config", func(t *testing.T) {
		// Test with invalid SMTP config (should fail gracefully)
		ec := &EmailChannel{
			SMTPServer: "invalid.smtp.server",
			SMTPPort:   9999,
			From:       "test@example.com",
			To:         []string{"admin@example.com"},
		}

		alert := Alert{
			Level:     AlertLevelInfo,
			Type:      AlertTypeClusterHealth,
			Message:   "Test alert",
			Timestamp: time.Now(),
		}

		err := ec.Send(alert)
		// Should error because SMTP server doesn't exist
		assert.Error(t, err)
	})
}

func TestSlackChannel_Structure(t *testing.T) {
	t.Run("Slack channel creation", func(t *testing.T) {
		sc := &SlackChannel{
			WebhookURL: "https://hooks.slack.com/services/test",
			Channel:    "#alerts",
			Username:   "etcd-monitor",
		}

		assert.Equal(t, "https://hooks.slack.com/services/test", sc.WebhookURL)
		assert.Equal(t, "#alerts", sc.Channel)
		assert.Equal(t, "etcd-monitor", sc.Username)
	})

	t.Run("Slack channel name", func(t *testing.T) {
		sc := &SlackChannel{}
		assert.Equal(t, "slack", sc.Name())
	})

	t.Run("Slack channel send with invalid webhook", func(t *testing.T) {
		// Test with invalid webhook URL (should fail gracefully)
		sc := &SlackChannel{
			WebhookURL: "http://invalid.webhook.url/test",
			Channel:    "#test",
			Username:   "bot",
		}

		alert := Alert{
			Level:     AlertLevelWarning,
			Type:      AlertTypeHighLatency,
			Message:   "Test alert",
			Timestamp: time.Now(),
		}

		err := sc.Send(alert)
		// Should error because webhook doesn't exist
		assert.Error(t, err)
	})
}

func TestPagerDutyChannel_Structure(t *testing.T) {
	t.Run("PagerDuty channel creation", func(t *testing.T) {
		pdc := &PagerDutyChannel{
			IntegrationKey: "test-key",
		}

		assert.Equal(t, "test-key", pdc.IntegrationKey)
	})

	t.Run("PagerDuty channel name", func(t *testing.T) {
		pdc := &PagerDutyChannel{}
		assert.Equal(t, "pagerduty", pdc.Name())
	})

	t.Run("PagerDuty channel send with invalid key", func(t *testing.T) {
		pdc := &PagerDutyChannel{
			IntegrationKey: "invalid-key",
		}

		alert := Alert{
			Level:     AlertLevelCritical,
			Type:      AlertTypeClusterHealth,
			Message:   "Test alert",
			Timestamp: time.Now(),
		}

		err := pdc.Send(alert)
		// Should error because integration key is invalid
		assert.Error(t, err)
	})
}

func TestWebhookChannel_Structure(t *testing.T) {
	t.Run("Webhook channel creation", func(t *testing.T) {
		wc := &WebhookChannel{
			URL: "http://example.com/webhook",
			Headers: map[string]string{
				"Authorization": "Bearer token",
			},
		}

		assert.Equal(t, "http://example.com/webhook", wc.URL)
		assert.Equal(t, "Bearer token", wc.Headers["Authorization"])
	})

	t.Run("Webhook channel name", func(t *testing.T) {
		wc := &WebhookChannel{}
		assert.Equal(t, "webhook", wc.Name())
	})

	t.Run("Webhook channel send with invalid URL", func(t *testing.T) {
		wc := &WebhookChannel{
			URL: "http://invalid.webhook.url/test",
		}

		alert := Alert{
			Level:     AlertLevelInfo,
			Type:      AlertTypeEtcdAlarm,
			Message:   "Test alert",
			Timestamp: time.Now(),
		}

		err := wc.Send(alert)
		// Should error because URL doesn't exist
		assert.Error(t, err)
	})
}

func TestConsoleChannel_Comprehensive(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Console channel with nil logger", func(t *testing.T) {
		cc := NewConsoleChannel(nil)
		assert.NotNil(t, cc)

		alert := Alert{
			Level:     AlertLevelInfo,
			Type:      AlertTypeClusterHealth,
			Message:   "Test with nil logger",
			Timestamp: time.Now(),
		}

		err := cc.Send(alert)
		assert.NoError(t, err)
	})

	t.Run("Console channel name", func(t *testing.T) {
		cc := NewConsoleChannel(logger)
		assert.Equal(t, "console", cc.Name())
	})

	t.Run("Console channel send different alert levels", func(t *testing.T) {
		cc := NewConsoleChannel(logger)

		levels := []AlertLevel{
			AlertLevelInfo,
			AlertLevelWarning,
			AlertLevelCritical,
		}

		for _, level := range levels {
			alert := Alert{
				Level:     level,
				Type:      AlertTypeClusterHealth,
				Message:   "Test alert level " + string(level),
				Timestamp: time.Now(),
			}

			err := cc.Send(alert)
			assert.NoError(t, err)
		}
	})

	t.Run("Console channel send with details", func(t *testing.T) {
		cc := NewConsoleChannel(logger)

		alert := Alert{
			Level:     AlertLevelWarning,
			Type:      AlertTypeHighLatency,
			Message:   "High latency detected",
			Timestamp: time.Now(),
			Details: map[string]interface{}{
				"latency_ms": 150,
				"endpoint":   "localhost:2379",
			},
		}

		err := cc.Send(alert)
		assert.NoError(t, err)
	})
}

func TestAlertManager_ClearAlert_Comprehensive(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	thresholds := AlertThresholds{
		MaxLatencyMs:            100,
		MaxDatabaseSizeMB:       8192,
		MinAvailableNodes:       2,
		MaxLeaderChangesPerHour: 3,
	}

	t.Run("Clear active alert", func(t *testing.T) {
		am := NewAlertManager(thresholds, logger)

		alert := Alert{
			Level:     AlertLevelWarning,
			Type:      AlertTypeHighLatency,
			Message:   "Latency spike",
			Timestamp: time.Now(),
		}

		am.TriggerAlert(alert)

		// Clear the alert
		am.ClearAlert(AlertTypeHighLatency, "Latency spike")

		// Alert should no longer be active
		active := am.GetActiveAlerts()
		found := false
		for _, a := range active {
			if a.Type == AlertTypeHighLatency && a.Message == "Latency spike" {
				found = true
				break
			}
		}
		assert.False(t, found, "Alert should have been cleared from active alerts")
	})

	t.Run("Clear non-existent alert", func(t *testing.T) {
		am := NewAlertManager(thresholds, logger)

		// Try to clear an alert that was never triggered
		am.ClearAlert(AlertTypeHighDiskUsage, "Disk full")

		// Should not panic or error
		assert.NotNil(t, am)
	})
}

func TestAlertTypes(t *testing.T) {
	t.Run("All alert types defined", func(t *testing.T) {
		alertTypes := []AlertType{
			AlertTypeClusterHealth,
			AlertTypeHighLatency,
			AlertTypeHighDiskUsage,
			AlertTypeLeaderElection,
			AlertTypeEtcdAlarm,
		}

		for _, alertType := range alertTypes {
			assert.NotEmpty(t, string(alertType))
		}
	})

	t.Run("All alert levels defined", func(t *testing.T) {
		alertLevels := []AlertLevel{
			AlertLevelInfo,
			AlertLevelWarning,
			AlertLevelCritical,
		}

		for _, level := range alertLevels {
			assert.NotEmpty(t, string(level))
		}
	})
}
