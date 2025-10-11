// Package api provides Prometheus metrics export
package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/etcd-monitor/taskmaster/pkg/monitor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// PrometheusExporter exports metrics in Prometheus format
type PrometheusExporter struct {
	monitorService *monitor.MonitorService
	logger         *zap.Logger

	// Cluster health metrics
	clusterHealthy       prometheus.Gauge
	clusterHasLeader     prometheus.Gauge
	clusterMemberCount   prometheus.Gauge
	clusterQuorumSize    prometheus.Gauge
	clusterLeaderChanges prometheus.Counter

	// Performance metrics
	readLatencyP50     prometheus.Gauge
	readLatencyP95     prometheus.Gauge
	readLatencyP99     prometheus.Gauge
	writeLatencyP50    prometheus.Gauge
	writeLatencyP95    prometheus.Gauge
	writeLatencyP99    prometheus.Gauge
	requestRate        prometheus.Gauge

	// Database metrics
	dbSize            prometheus.Gauge
	dbSizeInUse       prometheus.Gauge

	// Raft metrics
	proposalCommitted prometheus.Gauge
	proposalApplied   prometheus.Gauge
	proposalPending   prometheus.Gauge
	proposalFailed    prometheus.Gauge

	// Resource metrics
	memoryUsage          prometheus.Gauge
	diskUsage            prometheus.Gauge
	activeConnections    prometheus.Gauge
	watcherCount         prometheus.Gauge

	// Performance metrics
	fsyncDurationP95  prometheus.Gauge
	commitDurationP95 prometheus.Gauge
}

// NewPrometheusExporter creates a new Prometheus exporter
func NewPrometheusExporter(monitorService *monitor.MonitorService, logger *zap.Logger) *PrometheusExporter {
	pe := &PrometheusExporter{
		monitorService: monitorService,
		logger:         logger,
	}

	// Register metrics
	pe.registerMetrics()

	// Start metrics collection
	go pe.collectMetrics()

	return pe
}

// registerMetrics registers all Prometheus metrics
func (pe *PrometheusExporter) registerMetrics() {
	// Cluster health metrics
	pe.clusterHealthy = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "cluster",
		Name:      "healthy",
		Help:      "Whether the etcd cluster is healthy (1 = healthy, 0 = unhealthy)",
	})

	pe.clusterHasLeader = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "cluster",
		Name:      "has_leader",
		Help:      "Whether the etcd cluster has a leader (1 = has leader, 0 = no leader)",
	})

	pe.clusterMemberCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "cluster",
		Name:      "member_count",
		Help:      "Number of members in the etcd cluster",
	})

	pe.clusterQuorumSize = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "cluster",
		Name:      "quorum_size",
		Help:      "Required quorum size for the etcd cluster",
	})

	pe.clusterLeaderChanges = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "etcd",
		Subsystem: "cluster",
		Name:      "leader_changes_total",
		Help:      "Total number of leader changes",
	})

	// Performance metrics
	pe.readLatencyP50 = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "request",
		Name:      "read_latency_p50_milliseconds",
		Help:      "Read request latency P50 in milliseconds",
	})

	pe.readLatencyP95 = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "request",
		Name:      "read_latency_p95_milliseconds",
		Help:      "Read request latency P95 in milliseconds",
	})

	pe.readLatencyP99 = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "request",
		Name:      "read_latency_p99_milliseconds",
		Help:      "Read request latency P99 in milliseconds",
	})

	pe.writeLatencyP50 = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "request",
		Name:      "write_latency_p50_milliseconds",
		Help:      "Write request latency P50 in milliseconds",
	})

	pe.writeLatencyP95 = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "request",
		Name:      "write_latency_p95_milliseconds",
		Help:      "Write request latency P95 in milliseconds",
	})

	pe.writeLatencyP99 = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "request",
		Name:      "write_latency_p99_milliseconds",
		Help:      "Write request latency P99 in milliseconds",
	})

	pe.requestRate = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "request",
		Name:      "rate_per_second",
		Help:      "Request rate in operations per second",
	})

	// Database metrics
	pe.dbSize = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "mvcc",
		Name:      "db_total_size_bytes",
		Help:      "Total database size in bytes",
	})

	pe.dbSizeInUse = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "mvcc",
		Name:      "db_total_size_in_use_bytes",
		Help:      "Database size in use in bytes",
	})

	// Raft metrics
	pe.proposalCommitted = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "server",
		Name:      "proposals_committed_total",
		Help:      "Total number of consensus proposals committed",
	})

	pe.proposalApplied = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "server",
		Name:      "proposals_applied_total",
		Help:      "Total number of consensus proposals applied",
	})

	pe.proposalPending = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "server",
		Name:      "proposals_pending",
		Help:      "Current number of pending proposals",
	})

	pe.proposalFailed = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "server",
		Name:      "proposals_failed_total",
		Help:      "Total number of failed proposals",
	})

	// Resource metrics
	pe.memoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "server",
		Name:      "memory_usage_bytes",
		Help:      "Memory usage in bytes",
	})

	pe.diskUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "server",
		Name:      "disk_usage_bytes",
		Help:      "Disk usage in bytes",
	})

	pe.activeConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "server",
		Name:      "active_connections",
		Help:      "Number of active client connections",
	})

	pe.watcherCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "server",
		Name:      "watchers",
		Help:      "Number of active watchers",
	})

	// Performance metrics
	pe.fsyncDurationP95 = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "disk",
		Name:      "wal_fsync_duration_p95_milliseconds",
		Help:      "WAL fsync duration P95 in milliseconds",
	})

	pe.commitDurationP95 = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "etcd",
		Subsystem: "disk",
		Name:      "backend_commit_duration_p95_milliseconds",
		Help:      "Backend commit duration P95 in milliseconds",
	})

	// Register all metrics
	prometheus.MustRegister(
		pe.clusterHealthy,
		pe.clusterHasLeader,
		pe.clusterMemberCount,
		pe.clusterQuorumSize,
		pe.clusterLeaderChanges,
		pe.readLatencyP50,
		pe.readLatencyP95,
		pe.readLatencyP99,
		pe.writeLatencyP50,
		pe.writeLatencyP95,
		pe.writeLatencyP99,
		pe.requestRate,
		pe.dbSize,
		pe.dbSizeInUse,
		pe.proposalCommitted,
		pe.proposalApplied,
		pe.proposalPending,
		pe.proposalFailed,
		pe.memoryUsage,
		pe.diskUsage,
		pe.activeConnections,
		pe.watcherCount,
		pe.fsyncDurationP95,
		pe.commitDurationP95,
	)

	pe.logger.Info("Prometheus metrics registered")
}

// collectMetrics continuously collects metrics
func (pe *PrometheusExporter) collectMetrics() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		pe.updateMetrics()
	}
}

// updateMetrics updates all Prometheus metrics
func (pe *PrometheusExporter) updateMetrics() {
	// Get cluster status
	status, err := pe.monitorService.GetClusterStatus()
	if err != nil {
		pe.logger.Error("Failed to get cluster status for metrics", zap.Error(err))
		return
	}

	// Update cluster metrics
	if status.Healthy {
		pe.clusterHealthy.Set(1)
	} else {
		pe.clusterHealthy.Set(0)
	}

	if status.HasLeader {
		pe.clusterHasLeader.Set(1)
	} else {
		pe.clusterHasLeader.Set(0)
	}

	pe.clusterMemberCount.Set(float64(status.MemberCount))
	pe.clusterQuorumSize.Set(float64(status.QuorumSize))
	pe.clusterLeaderChanges.Add(float64(status.LeaderChanges))

	// Get current metrics
	metrics, err := pe.monitorService.GetCurrentMetrics()
	if err != nil {
		pe.logger.Error("Failed to get current metrics", zap.Error(err))
		return
	}

	// Update performance metrics
	pe.readLatencyP50.Set(metrics.ReadLatencyP50)
	pe.readLatencyP95.Set(metrics.ReadLatencyP95)
	pe.readLatencyP99.Set(metrics.ReadLatencyP99)
	pe.writeLatencyP50.Set(metrics.WriteLatencyP50)
	pe.writeLatencyP95.Set(metrics.WriteLatencyP95)
	pe.writeLatencyP99.Set(metrics.WriteLatencyP99)
	pe.requestRate.Set(metrics.RequestRate)

	// Update database metrics
	pe.dbSize.Set(float64(metrics.DBSize))
	pe.dbSizeInUse.Set(float64(metrics.DBSizeInUse))

	// Update Raft metrics
	pe.proposalCommitted.Set(float64(metrics.ProposalCommitted))
	pe.proposalApplied.Set(float64(metrics.ProposalApplied))
	pe.proposalPending.Set(float64(metrics.ProposalPending))
	pe.proposalFailed.Set(float64(metrics.ProposalFailed))

	// Update resource metrics
	pe.memoryUsage.Set(float64(metrics.MemoryUsage))
	pe.diskUsage.Set(float64(metrics.DiskUsage))
	pe.activeConnections.Set(float64(metrics.ActiveConnections))
	pe.watcherCount.Set(float64(metrics.WatcherCount))

	// Update performance metrics
	pe.fsyncDurationP95.Set(metrics.FSyncDurationP95)
	pe.commitDurationP95.Set(metrics.CommitDurationP95)
}

// Handler returns the HTTP handler for Prometheus metrics
func (pe *PrometheusExporter) Handler() http.Handler {
	return promhttp.Handler()
}

// RegisterWithServer registers the Prometheus handler with the API server
func (pe *PrometheusExporter) RegisterWithServer(server *Server) {
	server.router.Handle("/metrics", pe.Handler()).Methods("GET")
	pe.logger.Info("Prometheus metrics endpoint registered at /metrics")
}

// GetMetricsSummary returns a text summary of current metrics
func (pe *PrometheusExporter) GetMetricsSummary() string {
	status, err := pe.monitorService.GetClusterStatus()
	if err != nil {
		return fmt.Sprintf("Error getting status: %v", err)
	}

	metrics, err := pe.monitorService.GetCurrentMetrics()
	if err != nil {
		return fmt.Sprintf("Error getting metrics: %v", err)
	}

	summary := fmt.Sprintf(`
etcd Metrics Summary:
=====================
Cluster Health: %v
Has Leader: %v
Members: %d
Quorum Size: %d
Leader Changes: %d

Performance:
  Read Latency P99: %.2f ms
  Write Latency P99: %.2f ms
  Request Rate: %.2f ops/sec

Database:
  Size: %d bytes (%.2f MB)
  In Use: %d bytes (%.2f MB)

Raft:
  Proposals Committed: %d
  Proposals Applied: %d
  Proposals Pending: %d
  Proposals Failed: %d
`,
		status.Healthy,
		status.HasLeader,
		status.MemberCount,
		status.QuorumSize,
		status.LeaderChanges,
		metrics.ReadLatencyP99,
		metrics.WriteLatencyP99,
		metrics.RequestRate,
		metrics.DBSize,
		float64(metrics.DBSize)/(1024*1024),
		metrics.DBSizeInUse,
		float64(metrics.DBSizeInUse)/(1024*1024),
		metrics.ProposalCommitted,
		metrics.ProposalApplied,
		metrics.ProposalPending,
		metrics.ProposalFailed,
	)

	return summary
}
