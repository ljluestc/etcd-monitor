package monitor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAlertManager_Comprehensive(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	thresholds := AlertThresholds{}

	t.Run("NewAlertManager", func(t *testing.T) {
		am := NewAlertManager(thresholds, logger)
		assert.NotNil(t, am)
		assert.NotNil(t, am.logger)
		assert.NotNil(t, am.activeAlerts)
		assert.NotNil(t, am.alertHistory)
	})

	t.Run("AlertManager with nil logger", func(t *testing.T) {
		am := NewAlertManager(thresholds, nil)
		assert.NotNil(t, am)
		assert.NotNil(t, am.logger) // Should create a production logger
	})

	t.Run("TriggerAlert", func(t *testing.T) {
		am := NewAlertManager(thresholds, logger)
		
		alert := Alert{
			Level:     AlertLevelCritical,
			Type:      AlertTypeClusterHealth,
			Message:   "Test alert",
			Details:   map[string]interface{}{"key": "value"},
			Timestamp: time.Now(),
		}

		am.TriggerAlert(alert)
		
		active := am.GetActiveAlerts()
		assert.Len(t, active, 1)
		assert.Equal(t, alert.Level, active[0].Level)
		assert.Equal(t, alert.Type, active[0].Type)
		assert.Equal(t, alert.Message, active[0].Message)
	})

	t.Run("ClearAlert", func(t *testing.T) {
		am := NewAlertManager(thresholds, logger)
		
		alert := Alert{
			Level:     AlertLevelWarning,
			Type:      AlertTypeHighLatency,
			Message:   "Test alert",
			Timestamp: time.Now(),
		}

		am.TriggerAlert(alert)
		active := am.GetActiveAlerts()
		assert.Len(t, active, 1)

		am.ClearAlert(alert.Type, alert.Message)
		active = am.GetActiveAlerts()
		assert.Len(t, active, 0)
	})

	t.Run("GetAlertHistory", func(t *testing.T) {
		am := NewAlertManager(thresholds, logger)
		
		alert := Alert{
			Level:     AlertLevelInfo,
			Type:      AlertTypeEtcdAlarm,
			Message:   "Test alert",
			Timestamp: time.Now(),
		}

		am.TriggerAlert(alert)
		
		history := am.GetAlertHistory()
		assert.Len(t, history, 1)
		assert.Equal(t, alert.Level, history[0].Level)
		assert.Equal(t, alert.Type, history[0].Type)
		assert.Equal(t, alert.Message, history[0].Message)
	})

	t.Run("Alert constants", func(t *testing.T) {
		// Test alert level constants
		assert.Equal(t, "critical", string(AlertLevelCritical))
		assert.Equal(t, "warning", string(AlertLevelWarning))
		assert.Equal(t, "info", string(AlertLevelInfo))

		// Test alert type constants
		assert.Equal(t, "cluster_health", string(AlertTypeClusterHealth))
		assert.Equal(t, "high_latency", string(AlertTypeHighLatency))
		assert.Equal(t, "leader_election", string(AlertTypeLeaderElection))
		assert.Equal(t, "network_partition", string(AlertTypeNetworkPartition))
		assert.Equal(t, "high_disk_usage", string(AlertTypeHighDiskUsage))
		assert.Equal(t, "high_proposal_queue", string(AlertTypeHighProposalQueue))
		assert.Equal(t, "etcd_alarm", string(AlertTypeEtcdAlarm))
	})

	t.Run("Alert structure", func(t *testing.T) {
		alert := Alert{
			Level:     AlertLevelCritical,
			Type:      AlertTypeClusterHealth,
			Message:   "Test alert",
			Details:   map[string]interface{}{"key": "value", "number": 42},
			Timestamp: time.Now(),
		}

		assert.Equal(t, AlertLevelCritical, alert.Level)
		assert.Equal(t, AlertTypeClusterHealth, alert.Type)
		assert.Equal(t, "Test alert", alert.Message)
		assert.NotNil(t, alert.Details)
		assert.NotZero(t, alert.Timestamp)
	})

	t.Run("Console channel", func(t *testing.T) {
		consoleChannel := NewConsoleChannel(logger)
		assert.NotNil(t, consoleChannel)
		assert.NotNil(t, consoleChannel.logger)
	})

	// Note: Email and Slack channels are not implemented yet

	t.Run("Logger creation", func(t *testing.T) {
		// Test logger creation
		logger, err := zap.NewDevelopment()
		assert.NoError(t, err)
		assert.NotNil(t, logger)
	})
}