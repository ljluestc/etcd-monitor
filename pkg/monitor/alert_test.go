package monitor

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestAlertManager(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	thresholds := AlertThresholds{
		MaxLatencyMs:            100,
		MaxDatabaseSizeMB:       8192,
		MinAvailableNodes:       2,
		MaxLeaderChangesPerHour: 3,
	}

	am := NewAlertManager(thresholds, logger)

	t.Run("TriggerAlert", func(t *testing.T) {
		alert := Alert{
			Level:     AlertLevelWarning,
			Type:      AlertTypeHighLatency,
			Message:   "High latency detected",
			Timestamp: time.Now(),
		}

		// Trigger alert
		am.TriggerAlert(alert)

		// Check alert history
		history := am.GetAlertHistory()
		if len(history) != 1 {
			t.Errorf("Expected 1 alert in history, got %d", len(history))
		}

		if history[0].Message != alert.Message {
			t.Errorf("Expected message '%s', got '%s'", alert.Message, history[0].Message)
		}
	})

	t.Run("Deduplication", func(t *testing.T) {
		alert := Alert{
			Level:     AlertLevelWarning,
			Type:      AlertTypeHighLatency,
			Message:   "Duplicate alert",
			Timestamp: time.Now(),
		}

		// Trigger same alert twice
		am.TriggerAlert(alert)
		am.TriggerAlert(alert)

		// Should only have one additional alert (due to deduplication)
		history := am.GetAlertHistory()
		if len(history) != 2 { // 1 from previous test + 1 new
			t.Errorf("Expected 2 alerts in history, got %d", len(history))
		}
	})

	t.Run("GetActiveAlerts", func(t *testing.T) {
		active := am.GetActiveAlerts()
		if len(active) == 0 {
			t.Log("No active alerts (expected if dedupe window expired)")
		}
	})

	t.Run("ClearAlert", func(t *testing.T) {
		am.ClearAlert(AlertTypeHighLatency, "High latency detected")
		// Alert should be cleared from active alerts
	})
}

func TestEmailChannel(t *testing.T) {
	t.Skip("Email channel requires SMTP server")

	ec := &EmailChannel{
		SMTPServer: "localhost",
		SMTPPort:   1025, // mailhog test server
		From:       "monitor@example.com",
		To:         []string{"admin@example.com"},
	}

	alert := Alert{
		Level:     AlertLevelCritical,
		Type:      AlertTypeClusterHealth,
		Message:   "Test alert",
		Timestamp: time.Now(),
		Details: map[string]interface{}{
			"cluster": "test-cluster",
			"reason":  "Test notification",
		},
	}

	err := ec.Send(alert)
	if err != nil {
		t.Errorf("Failed to send email: %v", err)
	}
}

func TestSlackChannel(t *testing.T) {
	t.Skip("Slack channel requires webhook URL")

	sc := &SlackChannel{
		WebhookURL: "https://hooks.slack.com/services/test",
		Channel:    "#alerts",
		Username:   "etcd-monitor",
	}

	alert := Alert{
		Level:     AlertLevelWarning,
		Type:      AlertTypeHighDiskUsage,
		Message:   "Disk usage above threshold",
		Timestamp: time.Now(),
	}

	err := sc.Send(alert)
	if err != nil {
		t.Logf("Expected failure without real webhook: %v", err)
	}
}

func TestConsoleChannel(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cc := NewConsoleChannel(logger)

	alert := Alert{
		Level:     AlertLevelInfo,
		Type:      AlertTypeClusterHealth,
		Message:   "Cluster is healthy",
		Timestamp: time.Now(),
	}

	err := cc.Send(alert)
	if err != nil {
		t.Errorf("Console channel should not fail: %v", err)
	}
}
