package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// MockEtcdClient is a mock implementation of etcd client for testing
type MockEtcdClient struct {
	members     []*clientv3.Member
	status      *clientv3.StatusResponse
	alarms      []*clientv3.AlarmMember
	shouldError bool
}

func (m *MockEtcdClient) MemberList(ctx context.Context) (*clientv3.MemberListResponse, error) {
	if m.shouldError {
		return nil, assert.AnError
	}
	return &clientv3.MemberListResponse{Members: m.members}, nil
}

func (m *MockEtcdClient) Status(ctx context.Context, endpoint string) (*clientv3.StatusResponse, error) {
	if m.shouldError {
		return nil, assert.AnError
	}
	return m.status, nil
}

func (m *MockEtcdClient) AlarmList(ctx context.Context) (*clientv3.AlarmResponse, error) {
	if m.shouldError {
		return nil, assert.AnError
	}
	return &clientv3.AlarmResponse{Alarms: m.alarms}, nil
}

func (m *MockEtcdClient) Close() error {
	return nil
}

func TestMonitorService_Comprehensive(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("NewMonitorService with valid config", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			HealthCheckInterval: 30 * time.Second,
			MetricsInterval:     10 * time.Second,
			AlertThresholds: AlertThresholds{
				MaxLatencyMs:           100,
				MaxDatabaseSizeMB:      8192,
				MinAvailableNodes:      2,
				MaxLeaderChangesPerHour: 3,
				MaxErrorRate:           0.05,
				MinDiskSpacePercent:    10.0,
			},
			BenchmarkEnabled:  false,
			BenchmarkInterval: 1 * time.Hour,
		}

		service, err := NewMonitorService(config, logger)
		assert.NoError(t, err)
		assert.NotNil(t, service)
		assert.NotNil(t, service.config)
		assert.NotNil(t, service.logger)
		assert.NotNil(t, service.alertManager)
		assert.NotNil(t, service.healthChecker)
		assert.NotNil(t, service.metricsCollector)
	})

	t.Run("NewMonitorService with nil config", func(t *testing.T) {
		service, err := NewMonitorService(nil, logger)
		assert.Error(t, err)
		assert.Nil(t, service)
	})

	t.Run("NewMonitorService with nil logger", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		service, err := NewMonitorService(config, nil)
		assert.NoError(t, err)
		assert.NotNil(t, service)
		assert.NotNil(t, service.logger) // Should create a production logger
	})

	t.Run("Start and Stop service", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		service, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		// Start service
		err = service.Start()
		assert.NoError(t, err)

		// Stop service
		service.Stop()
	})

	t.Run("GetMetrics", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		service, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		metrics := service.GetMetrics()
		assert.NotNil(t, metrics)
	})

	t.Run("GetAlerts", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		service, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		alerts := service.GetAlerts()
		assert.NotNil(t, alerts)
	})

	t.Run("TriggerAlert", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		service, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		alert := Alert{
			Level:     AlertLevelCritical,
			Type:      AlertTypeClusterHealth,
			Message:   "Test alert",
			Timestamp: time.Now(),
		}

		service.TriggerAlert(alert)
		
		alerts := service.GetAlerts()
		assert.Len(t, alerts, 1)
	})

	t.Run("ClearAlert", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		service, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		alert := Alert{
			Level:     AlertLevelWarning,
			Type:      AlertTypeHighLatency,
			Message:   "Test alert",
			Timestamp: time.Now(),
		}

		service.TriggerAlert(alert)
		service.ClearAlert(alert.Type, alert.Message)
		
		alerts := service.GetAlerts()
		assert.Len(t, alerts, 0)
	})
}

func TestHealthChecker_Comprehensive(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("NewHealthChecker with mock client", func(t *testing.T) {
		mockClient := &MockEtcdClient{
			members: []*clientv3.Member{
				{ID: 1, Name: "member1", PeerURLs: []string{"http://localhost:2380"}},
			},
			status: &clientv3.StatusResponse{
				Header: &clientv3.ResponseHeader{},
				Version: "3.5.0",
				DbSize:  1024,
				Leader:  1,
				RaftIndex: 100,
				RaftTerm: 1,
			},
			alarms: []*clientv3.AlarmMember{},
		}

		hc := NewHealthChecker(mockClient, logger)
		assert.NotNil(t, hc)
		assert.NotNil(t, hc.client)
		assert.NotNil(t, hc.logger)
	})

	t.Run("CheckClusterHealth with mock client", func(t *testing.T) {
		mockClient := &MockEtcdClient{
			members: []*clientv3.Member{
				{ID: 1, Name: "member1", PeerURLs: []string{"http://localhost:2380"}},
			},
			status: &clientv3.StatusResponse{
				Header: &clientv3.ResponseHeader{},
				Version: "3.5.0",
				DbSize:  1024,
				Leader:  1,
				RaftIndex: 100,
				RaftTerm: 1,
			},
			alarms: []*clientv3.AlarmMember{},
		}

		hc := NewHealthChecker(mockClient, logger)
		ctx := context.Background()
		
		status, err := hc.CheckClusterHealth(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.True(t, status.Healthy)
	})

	t.Run("CheckClusterHealth with error", func(t *testing.T) {
		mockClient := &MockEtcdClient{
			shouldError: true,
		}

		hc := NewHealthChecker(mockClient, logger)
		ctx := context.Background()
		
		status, err := hc.CheckClusterHealth(ctx)
		assert.Error(t, err)
		assert.NotNil(t, status)
		assert.False(t, status.Healthy)
	})

	t.Run("CheckMemberHealth with mock client", func(t *testing.T) {
		mockClient := &MockEtcdClient{
			status: &clientv3.StatusResponse{
				Header: &clientv3.ResponseHeader{},
				Version: "3.5.0",
				DbSize:  1024,
				Leader:  1,
				RaftIndex: 100,
				RaftTerm: 1,
			},
		}

		hc := NewHealthChecker(mockClient, logger)
		ctx := context.Background()
		
		health, err := hc.CheckMemberHealth(ctx, "localhost:2379")
		assert.NoError(t, err)
		assert.NotNil(t, health)
	})

	t.Run("CheckAlarmStatus with mock client", func(t *testing.T) {
		mockClient := &MockEtcdClient{
			alarms: []*clientv3.AlarmMember{
				{MemberID: 1, Alarm: clientv3.AlarmTypeNONE},
			},
		}

		hc := NewHealthChecker(mockClient, logger)
		ctx := context.Background()
		
		alarms, err := hc.CheckAlarmStatus(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, alarms)
	})

	t.Run("CheckDiskUsage with mock client", func(t *testing.T) {
		mockClient := &MockEtcdClient{
			status: &clientv3.StatusResponse{
				Header: &clientv3.ResponseHeader{},
				DbSize: 1024,
			},
		}

		hc := NewHealthChecker(mockClient, logger)
		ctx := context.Background()
		
		usage, err := hc.CheckDiskUsage(ctx, "localhost:2379")
		assert.NoError(t, err)
		assert.NotNil(t, usage)
	})

	t.Run("CheckNetworkLatency with mock client", func(t *testing.T) {
		mockClient := &MockEtcdClient{}

		hc := NewHealthChecker(mockClient, logger)
		ctx := context.Background()
		
		latency, err := hc.CheckNetworkLatency(ctx, "localhost:2379")
		assert.NoError(t, err)
		assert.NotNil(t, latency)
	})

	t.Run("CheckLeaderElection with mock client", func(t *testing.T) {
		mockClient := &MockEtcdClient{
			status: &clientv3.StatusResponse{
				Header: &clientv3.ResponseHeader{},
				Leader: 1,
			},
		}

		hc := NewHealthChecker(mockClient, logger)
		ctx := context.Background()
		
		election, err := hc.CheckLeaderElection(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, election)
	})

	t.Run("CheckProposalQueue with mock client", func(t *testing.T) {
		mockClient := &MockEtcdClient{
			status: &clientv3.StatusResponse{
				Header: &clientv3.ResponseHeader{},
				RaftIndex: 100,
			},
		}

		hc := NewHealthChecker(mockClient, logger)
		ctx := context.Background()
		
		queue, err := hc.CheckProposalQueue(ctx, "localhost:2379")
		assert.NoError(t, err)
		assert.NotNil(t, queue)
	})
}

func TestMetricsCollector_Comprehensive(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("NewMetricsCollector with mock client", func(t *testing.T) {
		mockClient := &MockEtcdClient{}

		mc := NewMetricsCollector(mockClient, logger)
		assert.NotNil(t, mc)
		assert.NotNil(t, mc.client)
		assert.NotNil(t, mc.logger)
	})

	t.Run("CollectMetrics with mock client", func(t *testing.T) {
		mockClient := &MockEtcdClient{
			status: &clientv3.StatusResponse{
				Header: &clientv3.ResponseHeader{},
				Version: "3.5.0",
				DbSize:  1024,
				Leader:  1,
				RaftIndex: 100,
				RaftTerm: 1,
			},
		}

		mc := NewMetricsCollector(mockClient, logger)
		ctx := context.Background()
		
		metrics, err := mc.CollectMetrics(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, metrics)
	})

	t.Run("CollectMetrics with error", func(t *testing.T) {
		mockClient := &MockEtcdClient{
			shouldError: true,
		}

		mc := NewMetricsCollector(mockClient, logger)
		ctx := context.Background()
		
		metrics, err := mc.CollectMetrics(ctx)
		assert.Error(t, err)
		assert.Nil(t, metrics)
	})

	t.Run("GetMetrics", func(t *testing.T) {
		mockClient := &MockEtcdClient{}

		mc := NewMetricsCollector(mockClient, logger)
		
		metrics := mc.GetMetrics()
		assert.NotNil(t, metrics)
	})
}

func TestAlertManager_EdgeCases(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	thresholds := AlertThresholds{}

	t.Run("Alert deduplication with same message", func(t *testing.T) {
		am := NewAlertManager(thresholds, logger)
		
		alert := Alert{
			Level:     AlertLevelWarning,
			Type:      AlertTypeHighLatency,
			Message:   "Same message",
			Timestamp: time.Now(),
		}

		// Trigger same alert multiple times
		am.TriggerAlert(alert)
		am.TriggerAlert(alert)
		am.TriggerAlert(alert)

		active := am.GetActiveAlerts()
		assert.Len(t, active, 1) // Should be deduplicated
	})

	t.Run("Alert history limit", func(t *testing.T) {
		am := NewAlertManager(thresholds, logger)
		
		// Trigger many alerts to test history limit
		for i := 0; i < 1001; i++ {
			alert := Alert{
				Level:     AlertLevelInfo,
				Type:      AlertTypeEtcdAlarm,
				Message:   fmt.Sprintf("Alert %d", i),
				Timestamp: time.Now(),
			}
			am.TriggerAlert(alert)
		}

		history := am.GetAlertHistory()
		assert.Len(t, history, 1000) // Should be limited to maxHistory
	})

	t.Run("Alert with nil details", func(t *testing.T) {
		am := NewAlertManager(thresholds, logger)
		
		alert := Alert{
			Level:     AlertLevelCritical,
			Type:      AlertTypeClusterHealth,
			Message:   "Alert with nil details",
			Details:   nil,
			Timestamp: time.Now(),
		}

		am.TriggerAlert(alert)
		
		active := am.GetActiveAlerts()
		assert.Len(t, active, 1)
		assert.Nil(t, active[0].Details)
	})

	t.Run("Alert with empty details", func(t *testing.T) {
		am := NewAlertManager(thresholds, logger)
		
		alert := Alert{
			Level:     AlertLevelInfo,
			Type:      AlertTypeEtcdAlarm,
			Message:   "Alert with empty details",
			Details:   map[string]interface{}{},
			Timestamp: time.Now(),
		}

		am.TriggerAlert(alert)
		
		active := am.GetActiveAlerts()
		assert.Len(t, active, 1)
		assert.NotNil(t, active[0].Details)
		assert.Len(t, active[0].Details, 0)
	})

	t.Run("Clear non-existent alert", func(t *testing.T) {
		am := NewAlertManager(thresholds, logger)
		
		// Try to clear an alert that doesn't exist
		am.ClearAlert(AlertTypeHighLatency, "Non-existent alert")
		
		active := am.GetActiveAlerts()
		assert.Len(t, active, 0)
	})

	t.Run("Multiple alert types", func(t *testing.T) {
		am := NewAlertManager(thresholds, logger)
		
		alertTypes := []AlertType{
			AlertTypeClusterHealth,
			AlertTypeLeaderElection,
			AlertTypeNetworkPartition,
			AlertTypeHighLatency,
			AlertTypeHighDiskUsage,
			AlertTypeHighProposalQueue,
			AlertTypeEtcdAlarm,
		}

		for _, alertType := range alertTypes {
			alert := Alert{
				Level:     AlertLevelWarning,
				Type:      alertType,
				Message:   fmt.Sprintf("Alert for %s", alertType),
				Timestamp: time.Now(),
			}
			am.TriggerAlert(alert)
		}

		active := am.GetActiveAlerts()
		assert.Len(t, active, len(alertTypes))
	})

	t.Run("Multiple alert levels", func(t *testing.T) {
		am := NewAlertManager(thresholds, logger)
		
		alertLevels := []AlertLevel{
			AlertLevelInfo,
			AlertLevelWarning,
			AlertLevelCritical,
		}

		for _, level := range alertLevels {
			alert := Alert{
				Level:     level,
				Type:      AlertTypeClusterHealth,
				Message:   fmt.Sprintf("Alert with level %s", level),
				Timestamp: time.Now(),
			}
			am.TriggerAlert(alert)
		}

		active := am.GetActiveAlerts()
		assert.Len(t, active, len(alertLevels))
	})
}

func TestAlertManager_Performance(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	thresholds := AlertThresholds{}

	t.Run("High volume alert processing", func(t *testing.T) {
		am := NewAlertManager(thresholds, logger)
		
		start := time.Now()
		
		// Process 1000 alerts
		for i := 0; i < 1000; i++ {
			alert := Alert{
				Level:     AlertLevelInfo,
				Type:      AlertTypeEtcdAlarm,
				Message:   fmt.Sprintf("Alert %d", i),
				Timestamp: time.Now(),
			}
			am.TriggerAlert(alert)
		}
		
		duration := time.Since(start)
		
		// Should process 1000 alerts in less than 1 second
		assert.True(t, duration < time.Second, "Processing 1000 alerts took too long: %v", duration)
		
		history := am.GetAlertHistory()
		assert.Len(t, history, 1000)
	})

	t.Run("Concurrent alert processing", func(t *testing.T) {
		am := NewAlertManager(thresholds, logger)
		
		done := make(chan bool, 10)
		
		// Process alerts concurrently
		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				for j := 0; j < 100; j++ {
					alert := Alert{
						Level:     AlertLevelInfo,
						Type:      AlertTypeEtcdAlarm,
						Message:   fmt.Sprintf("Concurrent alert %d-%d", id, j),
						Timestamp: time.Now(),
					}
					am.TriggerAlert(alert)
				}
			}(i)
		}
		
		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
		
		history := am.GetAlertHistory()
		assert.Len(t, history, 1000)
	})
}

func TestAlertManager_Integration(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	thresholds := AlertThresholds{}

	t.Run("Full alert lifecycle", func(t *testing.T) {
		am := NewAlertManager(thresholds, logger)
		
		// 1. Trigger alert
		alert := Alert{
			Level:     AlertLevelCritical,
			Type:      AlertTypeClusterHealth,
			Message:   "Cluster health issue",
			Details:   map[string]interface{}{"node": "node1", "status": "down"},
			Timestamp: time.Now(),
		}
		
		am.TriggerAlert(alert)
		
		// 2. Verify alert is active
		active := am.GetActiveAlerts()
		assert.Len(t, active, 1)
		assert.Equal(t, alert.Message, active[0].Message)
		
		// 3. Verify alert is in history
		history := am.GetAlertHistory()
		assert.Len(t, history, 1)
		assert.Equal(t, alert.Message, history[0].Message)
		
		// 4. Clear alert
		am.ClearAlert(alert.Type, alert.Message)
		
		// 5. Verify alert is no longer active
		active = am.GetActiveAlerts()
		assert.Len(t, active, 0)
		
		// 6. Verify alert is still in history
		history = am.GetAlertHistory()
		assert.Len(t, history, 1)
	})

	t.Run("Alert threshold validation", func(t *testing.T) {
		thresholds := AlertThresholds{
			MaxLatencyMs:           100,
			MaxDatabaseSizeMB:      8192,
			MinAvailableNodes:      2,
			MaxLeaderChangesPerHour: 3,
			MaxErrorRate:           0.05,
			MinDiskSpacePercent:    10.0,
		}
		
		am := NewAlertManager(thresholds, logger)
		
		// Test that thresholds are properly set
		assert.Equal(t, 100, am.thresholds.MaxLatencyMs)
		assert.Equal(t, 8192, am.thresholds.MaxDatabaseSizeMB)
		assert.Equal(t, 2, am.thresholds.MinAvailableNodes)
		assert.Equal(t, 3, am.thresholds.MaxLeaderChangesPerHour)
		assert.Equal(t, 0.05, am.thresholds.MaxErrorRate)
		assert.Equal(t, 10.0, am.thresholds.MinDiskSpacePercent)
	})
}
