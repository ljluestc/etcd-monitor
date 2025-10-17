package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// MockEtcdClientForService is a comprehensive mock for etcd client
type MockEtcdClientForService struct {
	statusResponse *clientv3.StatusResponse
	statusError    error
	memberResponse *clientv3.MemberListResponse
	memberError    error
	alarmResponse  *clientv3.AlarmResponse
	alarmError     error
	watchChan      chan clientv3.WatchResponse
	closeError     error
}

func (m *MockEtcdClientForService) Status(ctx context.Context, endpoint string) (*clientv3.StatusResponse, error) {
	return m.statusResponse, m.statusError
}

func (m *MockEtcdClientForService) MemberList(ctx context.Context) (*clientv3.MemberListResponse, error) {
	return m.memberResponse, m.memberError
}

func (m *MockEtcdClientForService) AlarmList(ctx context.Context) (*clientv3.AlarmResponse, error) {
	return m.alarmResponse, m.alarmError
}

func (m *MockEtcdClientForService) Watch(ctx context.Context, key string, opts ...clientv3.OpOption) clientv3.WatchChan {
	return m.watchChan
}

func (m *MockEtcdClientForService) Close() error {
	return m.closeError
}

func TestMonitorService_Comprehensive(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("NewMonitorService", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)
		assert.NotNil(t, ms)
		assert.Equal(t, config, ms.config)
		assert.Equal(t, logger, ms.logger)
		assert.NotNil(t, ms.ctx)
		assert.NotNil(t, ms.cancel)
		assert.False(t, ms.isRunning)
	})

	t.Run("NewMonitorService with nil logger", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, nil)
		require.NoError(t, err)
		assert.NotNil(t, ms)
		assert.NotNil(t, ms.logger) // Should create a production logger
	})

	t.Run("Start - Success", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		// Mock the etcd client creation
		ms.client = &MockEtcdClientForService{
			statusResponse: &clientv3.StatusResponse{
				Header: &clientv3.ResponseHeader{},
				Version: "3.5.0",
				DbSize: 1024,
				Leader: 1,
				RaftIndex: 100,
				RaftTerm: 1,
				RaftAppliedIndex: 100,
			},
			memberResponse: &clientv3.MemberListResponse{
				Header: &clientv3.ResponseHeader{},
				Members: []*clientv3.Member{
					{ID: 1, Name: "member1", PeerURLs: []string{"http://localhost:2380"}, ClientURLs: []string{"http://localhost:2379"}},
				},
			},
			alarmResponse: &clientv3.AlarmResponse{
				Header: &clientv3.ResponseHeader{},
				Alarms: []*clientv3.AlarmMember{},
			},
			watchChan: make(chan clientv3.WatchResponse, 1),
		}

		err = ms.Start()
		require.NoError(t, err)
		assert.True(t, ms.IsRunning())

		// Clean up
		ms.Stop()
	})

	t.Run("Start - Already Running", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		// Mock the etcd client
		ms.client = &MockEtcdClientForService{}
		ms.isRunning = true // Simulate already running

		err = ms.Start()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already running")
	})

	t.Run("Start - etcd Connection Error", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"invalid-endpoint"},
			DialTimeout:         1 * time.Millisecond, // Very short timeout
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		err = ms.Start()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to connect to etcd")
	})

	t.Run("Stop - Not Running", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		err = ms.Stop()
		require.NoError(t, err)
		assert.False(t, ms.IsRunning())
	})

	t.Run("Stop - Running", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		// Mock the etcd client
		ms.client = &MockEtcdClientForService{
			statusResponse: &clientv3.StatusResponse{
				Header: &clientv3.ResponseHeader{},
				Version: "3.5.0",
				DbSize: 1024,
				Leader: 1,
				RaftIndex: 100,
				RaftTerm: 1,
				RaftAppliedIndex: 100,
			},
			memberResponse: &clientv3.MemberListResponse{
				Header: &clientv3.ResponseHeader{},
				Members: []*clientv3.Member{
					{ID: 1, Name: "member1", PeerURLs: []string{"http://localhost:2380"}, ClientURLs: []string{"http://localhost:2379"}},
				},
			},
			alarmResponse: &clientv3.AlarmResponse{
				Header: &clientv3.ResponseHeader{},
				Alarms: []*clientv3.AlarmMember{},
			},
			watchChan: make(chan clientv3.WatchResponse, 1),
		}

		err = ms.Start()
		require.NoError(t, err)

		err = ms.Stop()
		require.NoError(t, err)
		assert.False(t, ms.IsRunning())
	})

	t.Run("GetClusterStatus - Success", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		// Mock health checker
		mockHealthChecker := &MockHealthChecker{
			clusterStatus: &ClusterStatus{
				Healthy:     true,
				LeaderID:    1,
				MemberCount: 3,
				QuorumSize:  2,
				HasLeader:   true,
				LastCheck:   time.Now(),
			},
		}
		ms.healthChecker = mockHealthChecker

		status, err := ms.GetClusterStatus()
		require.NoError(t, err)
		assert.NotNil(t, status)
		assert.True(t, status.Healthy)
		assert.Equal(t, uint64(1), status.LeaderID)
	})

	t.Run("GetClusterStatus - No Health Checker", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		status, err := ms.GetClusterStatus()
		assert.Error(t, err)
		assert.Nil(t, status)
		assert.Contains(t, err.Error(), "health checker not initialized")
	})

	t.Run("GetCurrentMetrics - Success", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		// Mock metrics collector
		mockMetricsCollector := &MockMetricsCollector{
			metrics: &MetricsSnapshot{
				Timestamp:       time.Now(),
				RequestRate:     100.0,
				ReadLatencyP50:  1.0,
				ReadLatencyP95:  5.0,
				ReadLatencyP99:  10.0,
				DBSize:          1024 * 1024,
				DBSizeInUse:     512 * 1024,
			},
		}
		ms.metricsCollector = mockMetricsCollector

		metrics, err := ms.GetCurrentMetrics()
		require.NoError(t, err)
		assert.NotNil(t, metrics)
		assert.Equal(t, 100.0, metrics.RequestRate)
	})

	t.Run("GetCurrentMetrics - No Metrics Collector", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		metrics, err := ms.GetCurrentMetrics()
		assert.Error(t, err)
		assert.Nil(t, metrics)
		assert.Contains(t, err.Error(), "metrics collector not initialized")
	})

	t.Run("GetHealthChecker", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		// Initially nil
		assert.Nil(t, ms.GetHealthChecker())

		// Set health checker
		mockHealthChecker := &MockHealthChecker{}
		ms.healthChecker = mockHealthChecker
		assert.Equal(t, mockHealthChecker, ms.GetHealthChecker())
	})

	t.Run("GetMetricsCollector", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		// Initially nil
		assert.Nil(t, ms.GetMetricsCollector())

		// Set metrics collector
		mockMetricsCollector := &MockMetricsCollector{}
		ms.metricsCollector = mockMetricsCollector
		assert.Equal(t, mockMetricsCollector, ms.GetMetricsCollector())
	})

	t.Run("GetAlertManager", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		// Initially nil
		assert.Nil(t, ms.GetAlertManager())

		// Set alert manager
		mockAlertManager := &MockAlertManager{}
		ms.alertManager = mockAlertManager
		assert.Equal(t, mockAlertManager, ms.GetAlertManager())
	})

	t.Run("IsRunning", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		// Initially not running
		assert.False(t, ms.IsRunning())

		// Set running state
		ms.isRunning = true
		assert.True(t, ms.IsRunning())
	})
}

func TestMonitorService_HealthChecks(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("CheckHealthAlerts - Unhealthy Cluster", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		// Mock alert manager
		mockAlertManager := &MockAlertManager{}
		ms.alertManager = mockAlertManager

		// Test unhealthy cluster
		status := &ClusterStatus{
			Healthy: false,
		}

		ms.checkHealthAlerts(status)

		// Should trigger alert
		alerts := mockAlertManager.GetActiveAlerts()
		assert.Len(t, alerts, 1)
		assert.Equal(t, AlertLevelCritical, alerts[0].Level)
		assert.Equal(t, AlertTypeClusterHealth, alerts[0].Type)
	})

	t.Run("CheckHealthAlerts - No Leader", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		// Mock alert manager
		mockAlertManager := &MockAlertManager{}
		ms.alertManager = mockAlertManager

		// Test no leader
		status := &ClusterStatus{
			Healthy:   true,
			HasLeader: false,
		}

		ms.checkHealthAlerts(status)

		// Should trigger alert
		alerts := mockAlertManager.GetActiveAlerts()
		assert.Len(t, alerts, 1)
		assert.Equal(t, AlertLevelCritical, alerts[0].Level)
		assert.Equal(t, AlertTypeLeaderElection, alerts[0].Type)
	})

	t.Run("CheckHealthAlerts - Network Partition", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		// Mock alert manager
		mockAlertManager := &MockAlertManager{}
		ms.alertManager = mockAlertManager

		// Test network partition
		status := &ClusterStatus{
			Healthy:         true,
			HasLeader:       true,
			NetworkPartition: true,
		}

		ms.checkHealthAlerts(status)

		// Should trigger alert
		alerts := mockAlertManager.GetActiveAlerts()
		assert.Len(t, alerts, 1)
		assert.Equal(t, AlertLevelCritical, alerts[0].Level)
		assert.Equal(t, AlertTypeNetworkPartition, alerts[0].Type)
	})

	t.Run("CheckHealthAlerts - Alarms", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		// Mock alert manager
		mockAlertManager := &MockAlertManager{}
		ms.alertManager = mockAlertManager

		// Test with alarms
		status := &ClusterStatus{
			Healthy: true,
			HasLeader: true,
			Alarms: []AlarmInfo{
				{Type: "CORRUPT", MemberID: 1, Triggered: time.Now()},
			},
		}

		ms.checkHealthAlerts(status)

		// Should trigger alert
		alerts := mockAlertManager.GetActiveAlerts()
		assert.Len(t, alerts, 1)
		assert.Equal(t, AlertLevelWarning, alerts[0].Level)
		assert.Equal(t, AlertTypeEtcdAlarm, alerts[0].Type)
	})
}

func TestMonitorService_MetricAlerts(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("CheckMetricAlerts - High Latency", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
			AlertThresholds: AlertThresholds{
				MaxLatencyMs: 50,
			},
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		// Mock alert manager
		mockAlertManager := &MockAlertManager{}
		ms.alertManager = mockAlertManager

		// Test high latency
		metrics := &MetricsSnapshot{
			WriteLatencyP99: 100.0, // Above threshold
		}

		ms.checkMetricAlerts(metrics)

		// Should trigger alert
		alerts := mockAlertManager.GetActiveAlerts()
		assert.Len(t, alerts, 1)
		assert.Equal(t, AlertLevelWarning, alerts[0].Level)
		assert.Equal(t, AlertTypeHighLatency, alerts[0].Type)
	})

	t.Run("CheckMetricAlerts - High Disk Usage", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
			AlertThresholds: AlertThresholds{
				MaxDatabaseSizeMB: 100,
			},
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		// Mock alert manager
		mockAlertManager := &MockAlertManager{}
		ms.alertManager = mockAlertManager

		// Test high disk usage
		metrics := &MetricsSnapshot{
			DBSize: 200 * 1024 * 1024, // 200MB, above threshold
		}

		ms.checkMetricAlerts(metrics)

		// Should trigger alert
		alerts := mockAlertManager.GetActiveAlerts()
		assert.Len(t, alerts, 1)
		assert.Equal(t, AlertLevelWarning, alerts[0].Level)
		assert.Equal(t, AlertTypeHighDiskUsage, alerts[0].Type)
	})

	t.Run("CheckMetricAlerts - High Proposal Queue", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		// Mock alert manager
		mockAlertManager := &MockAlertManager{}
		ms.alertManager = mockAlertManager

		// Test high proposal queue
		metrics := &MetricsSnapshot{
			ProposalPending: 150, // Above threshold of 100
		}

		ms.checkMetricAlerts(metrics)

		// Should trigger alert
		alerts := mockAlertManager.GetActiveAlerts()
		assert.Len(t, alerts, 1)
		assert.Equal(t, AlertLevelWarning, alerts[0].Level)
		assert.Equal(t, AlertTypeHighProposalQueue, alerts[0].Type)
	})
}

func TestMonitorService_Concurrency(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Concurrent Operations", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)

		// Mock components
		mockHealthChecker := &MockHealthChecker{
			clusterStatus: &ClusterStatus{
				Healthy:     true,
				LeaderID:    1,
				MemberCount: 3,
				QuorumSize:  2,
				HasLeader:   true,
				LastCheck:   time.Now(),
			},
		}
		ms.healthChecker = mockHealthChecker

		mockMetricsCollector := &MockMetricsCollector{
			metrics: &MetricsSnapshot{
				Timestamp:   time.Now(),
				RequestRate: 100.0,
			},
		}
		ms.metricsCollector = mockMetricsCollector

		// Test concurrent access
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				defer func() { done <- true }()
				
				// Test IsRunning
				_ = ms.IsRunning()
				
				// Test GetClusterStatus
				status, err := ms.GetClusterStatus()
				if err == nil {
					assert.NotNil(t, status)
				}
				
				// Test GetCurrentMetrics
				metrics, err := ms.GetCurrentMetrics()
				if err == nil {
					assert.NotNil(t, metrics)
				}
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

func TestMonitorService_EdgeCases(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Nil Config", func(t *testing.T) {
		ms, err := NewMonitorService(nil, logger)
		require.NoError(t, err)
		assert.NotNil(t, ms)
	})

	t.Run("Empty Endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{},
			DialTimeout:         5 * time.Second,
			HealthCheckInterval: 10 * time.Second,
			MetricsInterval:     5 * time.Second,
			WatchInterval:       1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)
		assert.NotNil(t, ms)
	})

	t.Run("Zero Timeouts", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         0,
			HealthCheckInterval: 0,
			MetricsInterval:     0,
			WatchInterval:       0,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)
		assert.NotNil(t, ms)
	})

	t.Run("Negative Timeouts", func(t *testing.T) {
		config := &Config{
			Endpoints:           []string{"http://localhost:2379"},
			DialTimeout:         -1 * time.Second,
			HealthCheckInterval: -1 * time.Second,
			MetricsInterval:     -1 * time.Second,
			WatchInterval:       -1 * time.Second,
		}

		ms, err := NewMonitorService(config, logger)
		require.NoError(t, err)
		assert.NotNil(t, ms)
	})
}

// Mock implementations for testing

type MockHealthChecker struct {
	clusterStatus    *ClusterStatus
	clusterStatusErr error
	memberList       []*MemberInfo
	memberListErr    error
}

func (m *MockHealthChecker) CheckClusterHealth(ctx context.Context) (*ClusterStatus, error) {
	return m.clusterStatus, m.clusterStatusErr
}

func (m *MockHealthChecker) GetMemberList(ctx context.Context) ([]*MemberInfo, error) {
	return m.memberList, m.memberListErr
}

func (m *MockHealthChecker) CheckNodeHealth(ctx context.Context, endpoint string) (bool, error) {
	return true, nil
}

func (m *MockHealthChecker) DetectNetworkPartition(ctx context.Context) (bool, error) {
	return false, nil
}

type MockMetricsCollector struct {
	metrics    *MetricsSnapshot
	metricsErr error
}

func (m *MockMetricsCollector) CollectMetrics(ctx context.Context) (*MetricsSnapshot, error) {
	return m.metrics, m.metricsErr
}

type MockAlertManager struct {
	activeAlerts []Alert
	alertHistory []Alert
}

func (m *MockAlertManager) TriggerAlert(alert Alert) {
	m.activeAlerts = append(m.activeAlerts, alert)
	m.alertHistory = append(m.alertHistory, alert)
}

func (m *MockAlertManager) GetActiveAlerts() []Alert {
	return m.activeAlerts
}

func (m *MockAlertManager) GetAlertHistory() []Alert {
	return m.alertHistory
}

func (m *MockAlertManager) ClearAlert(alertType AlertType, message string) {
	// Simple implementation for testing
	for i, alert := range m.activeAlerts {
		if alert.Type == alertType && alert.Message == message {
			m.activeAlerts = append(m.activeAlerts[:i], m.activeAlerts[i+1:]...)
			break
		}
	}
}

func (m *MockAlertManager) ClearAllAlerts() {
	m.activeAlerts = []Alert{}
}
