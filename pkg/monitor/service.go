// Package monitor provides comprehensive monitoring for etcd clusters
package monitor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// MonitorService is the main service for monitoring etcd clusters
type MonitorService struct {
	client          *clientv3.Client
	config          *Config
	logger          *zap.Logger
	healthChecker   *HealthChecker
	metricsCollector *MetricsCollector
	alertManager    *AlertManager
	ctx             context.Context
	cancel          context.CancelFunc
	wg              sync.WaitGroup
	mu              sync.RWMutex
	isRunning       bool
}

// Config holds the configuration for the monitor service
type Config struct {
	// etcd connection
	Endpoints   []string
	DialTimeout time.Duration
	TLS         *TLSConfig

	// Collection intervals
	HealthCheckInterval  time.Duration
	MetricsInterval      time.Duration
	WatchInterval        time.Duration

	// Alert configuration
	AlertThresholds AlertThresholds

	// Benchmark configuration
	BenchmarkEnabled bool
	BenchmarkInterval time.Duration
}

// TLSConfig holds TLS configuration
type TLSConfig struct {
	CertFile      string
	KeyFile       string
	CAFile        string
	InsecureSkipVerify bool
}

// AlertThresholds defines thresholds for alerts
type AlertThresholds struct {
	MaxLatencyMs           int
	MaxDatabaseSizeMB      int
	MinAvailableNodes      int
	MaxLeaderChangesPerHour int
	MaxErrorRate           float64
	MinDiskSpacePercent    float64
}

// ClusterStatus represents the current status of the cluster
type ClusterStatus struct {
	Healthy           bool
	LeaderID          uint64
	MemberCount       int
	QuorumSize        int
	HasLeader         bool
	LeaderChanges     int
	LastLeaderChange  time.Time
	NetworkPartition  bool
	Alarms            []AlarmInfo
	LastCheck         time.Time
}

// AlarmInfo represents an alarm in the cluster
type AlarmInfo struct {
	Type      string
	MemberID  uint64
	Triggered time.Time
}

// MetricsSnapshot represents a point-in-time metrics snapshot
type MetricsSnapshot struct {
	Timestamp            time.Time

	// Request metrics
	RequestRate          float64  // ops/sec
	ReadLatencyP50       float64  // ms
	ReadLatencyP95       float64  // ms
	ReadLatencyP99       float64  // ms
	WriteLatencyP50      float64  // ms
	WriteLatencyP95      float64  // ms
	WriteLatencyP99      float64  // ms

	// Database metrics
	DBSize               int64    // bytes
	DBSizeInUse          int64    // bytes

	// Raft metrics
	ProposalCommitted    uint64
	ProposalApplied      uint64
	ProposalPending      uint64
	ProposalFailed       uint64

	// Resource metrics
	MemoryUsage          uint64   // bytes
	CPUUsage             float64  // percentage
	DiskUsage            uint64   // bytes
	NetworkIn            uint64   // bytes/sec
	NetworkOut           uint64   // bytes/sec

	// Client metrics
	ActiveConnections    int
	WatcherCount         int

	// Performance metrics
	FSyncDurationP95     float64  // ms
	CommitDurationP95    float64  // ms
}

// NewMonitorService creates a new monitoring service
func NewMonitorService(config *Config, logger *zap.Logger) (*MonitorService, error) {
	if logger == nil {
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("failed to create logger: %w", err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	ms := &MonitorService{
		config:    config,
		logger:    logger,
		ctx:       ctx,
		cancel:    cancel,
		isRunning: false,
	}

	return ms, nil
}

// Start begins monitoring the etcd cluster
func (ms *MonitorService) Start() error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if ms.isRunning {
		return fmt.Errorf("monitor service already running")
	}

	// Connect to etcd
	clientConfig := clientv3.Config{
		Endpoints:   ms.config.Endpoints,
		DialTimeout: ms.config.DialTimeout,
	}

	var err error
	ms.client, err = clientv3.New(clientConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to etcd: %w", err)
	}

	// Initialize components
	ms.healthChecker = NewHealthChecker(ms.client, ms.logger)
	ms.metricsCollector = NewMetricsCollector(ms.client, ms.logger)
	ms.alertManager = NewAlertManager(ms.config.AlertThresholds, ms.logger)

	// Start monitoring goroutines
	ms.wg.Add(3)
	go ms.runHealthChecks()
	go ms.runMetricsCollection()
	go ms.runWatcher()

	ms.isRunning = true
	ms.logger.Info("Monitor service started", zap.Strings("endpoints", ms.config.Endpoints))

	return nil
}

// Stop stops the monitoring service
func (ms *MonitorService) Stop() error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if !ms.isRunning {
		return nil
	}

	ms.cancel()
	ms.wg.Wait()

	if ms.client != nil {
		if err := ms.client.Close(); err != nil {
			ms.logger.Error("Error closing etcd client", zap.Error(err))
		}
	}

	ms.isRunning = false
	ms.logger.Info("Monitor service stopped")

	return nil
}

// runHealthChecks performs periodic health checks
func (ms *MonitorService) runHealthChecks() {
	defer ms.wg.Done()

	ticker := time.NewTicker(ms.config.HealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ms.ctx.Done():
			return
		case <-ticker.C:
			status, err := ms.healthChecker.CheckClusterHealth(ms.ctx)
			if err != nil {
				ms.logger.Error("Health check failed", zap.Error(err))
				continue
			}

			// Check for alerts
			ms.checkHealthAlerts(status)
		}
	}
}

// runMetricsCollection collects metrics periodically
func (ms *MonitorService) runMetricsCollection() {
	defer ms.wg.Done()

	ticker := time.NewTicker(ms.config.MetricsInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ms.ctx.Done():
			return
		case <-ticker.C:
			metrics, err := ms.metricsCollector.CollectMetrics(ms.ctx)
			if err != nil {
				ms.logger.Error("Metrics collection failed", zap.Error(err))
				continue
			}

			// Check for metric-based alerts
			ms.checkMetricAlerts(metrics)
		}
	}
}

// runWatcher watches for key changes and cluster events
func (ms *MonitorService) runWatcher() {
	defer ms.wg.Done()

	// Watch for member changes
	watchChan := ms.client.Watch(ms.ctx, "/", clientv3.WithPrefix(), clientv3.WithPrevKV())

	for {
		select {
		case <-ms.ctx.Done():
			return
		case watchResp := <-watchChan:
			if watchResp.Err() != nil {
				ms.logger.Error("Watch error", zap.Error(watchResp.Err()))
				continue
			}

			for _, event := range watchResp.Events {
				ms.handleWatchEvent(event)
			}
		}
	}
}

// handleWatchEvent processes watch events
func (ms *MonitorService) handleWatchEvent(event *clientv3.Event) {
	switch event.Type {
	case mvccpb.PUT:
		ms.logger.Debug("Key created/updated",
			zap.String("key", string(event.Kv.Key)),
			zap.Int64("version", event.Kv.Version))
	case mvccpb.DELETE:
		ms.logger.Debug("Key deleted",
			zap.String("key", string(event.Kv.Key)))
	}
}

// checkHealthAlerts checks health status and triggers alerts
func (ms *MonitorService) checkHealthAlerts(status *ClusterStatus) {
	if !status.Healthy {
		ms.alertManager.TriggerAlert(Alert{
			Level:    AlertLevelCritical,
			Type:     AlertTypeClusterHealth,
			Message:  "Cluster is unhealthy",
			Details:  map[string]interface{}{"status": status},
			Timestamp: time.Now(),
		})
	}

	if !status.HasLeader {
		ms.alertManager.TriggerAlert(Alert{
			Level:    AlertLevelCritical,
			Type:     AlertTypeLeaderElection,
			Message:  "Cluster has no leader",
			Timestamp: time.Now(),
		})
	}

	if status.NetworkPartition {
		ms.alertManager.TriggerAlert(Alert{
			Level:    AlertLevelCritical,
			Type:     AlertTypeNetworkPartition,
			Message:  "Network partition detected",
			Timestamp: time.Now(),
		})
	}

	if len(status.Alarms) > 0 {
		for _, alarm := range status.Alarms {
			ms.alertManager.TriggerAlert(Alert{
				Level:    AlertLevelWarning,
				Type:     AlertTypeEtcdAlarm,
				Message:  fmt.Sprintf("etcd alarm: %s", alarm.Type),
				Details:  map[string]interface{}{"alarm": alarm},
				Timestamp: time.Now(),
			})
		}
	}
}

// checkMetricAlerts checks metrics and triggers alerts
func (ms *MonitorService) checkMetricAlerts(metrics *MetricsSnapshot) {
	thresholds := ms.config.AlertThresholds

	// Check latency
	if metrics.WriteLatencyP99 > float64(thresholds.MaxLatencyMs) {
		ms.alertManager.TriggerAlert(Alert{
			Level:    AlertLevelWarning,
			Type:     AlertTypeHighLatency,
			Message:  fmt.Sprintf("High write latency: %.2fms (threshold: %dms)", metrics.WriteLatencyP99, thresholds.MaxLatencyMs),
			Details:  map[string]interface{}{"p99_latency": metrics.WriteLatencyP99},
			Timestamp: time.Now(),
		})
	}

	// Check database size
	dbSizeMB := float64(metrics.DBSize) / (1024 * 1024)
	if dbSizeMB > float64(thresholds.MaxDatabaseSizeMB) {
		ms.alertManager.TriggerAlert(Alert{
			Level:    AlertLevelWarning,
			Type:     AlertTypeHighDiskUsage,
			Message:  fmt.Sprintf("Database size exceeds threshold: %.2fMB (threshold: %dMB)", dbSizeMB, thresholds.MaxDatabaseSizeMB),
			Details:  map[string]interface{}{"db_size_mb": dbSizeMB},
			Timestamp: time.Now(),
		})
	}

	// Check pending proposals
	if metrics.ProposalPending > 100 {
		ms.alertManager.TriggerAlert(Alert{
			Level:    AlertLevelWarning,
			Type:     AlertTypeHighProposalQueue,
			Message:  fmt.Sprintf("High pending proposals: %d", metrics.ProposalPending),
			Details:  map[string]interface{}{"pending": metrics.ProposalPending},
			Timestamp: time.Now(),
		})
	}
}

// GetClusterStatus returns the current cluster status
func (ms *MonitorService) GetClusterStatus() (*ClusterStatus, error) {
	if ms.healthChecker == nil {
		return nil, fmt.Errorf("health checker not initialized")
	}
	return ms.healthChecker.CheckClusterHealth(context.Background())
}

// GetCurrentMetrics returns the current metrics snapshot
func (ms *MonitorService) GetCurrentMetrics() (*MetricsSnapshot, error) {
	if ms.metricsCollector == nil {
		return nil, fmt.Errorf("metrics collector not initialized")
	}
	return ms.metricsCollector.CollectMetrics(context.Background())
}

// IsRunning returns whether the service is running
func (ms *MonitorService) IsRunning() bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.isRunning
}

// GetHealthChecker returns the health checker instance
func (ms *MonitorService) GetHealthChecker() *HealthChecker {
	return ms.healthChecker
}

// GetMetricsCollector returns the metrics collector instance
func (ms *MonitorService) GetMetricsCollector() *MetricsCollector {
	return ms.metricsCollector
}

// GetAlertManager returns the alert manager instance
func (ms *MonitorService) GetAlertManager() *AlertManager {
	return ms.alertManager
}
