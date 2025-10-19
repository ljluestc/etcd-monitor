package monitor

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// MetricsCollector collects performance metrics from etcd
type MetricsCollector struct {
	client         *clientv3.Client
	logger         *zap.Logger
	mu             sync.RWMutex
	latencyHistory []LatencyMeasurement
	maxHistory     int
}

// LatencyMeasurement records a latency measurement
type LatencyMeasurement struct {
	Timestamp     time.Time
	ReadLatency   time.Duration
	WriteLatency  time.Duration
	Operation     string
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(client *clientv3.Client, logger *zap.Logger) *MetricsCollector {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	return &MetricsCollector{
		client:         client,
		logger:         logger,
		latencyHistory: make([]LatencyMeasurement, 0),
		maxHistory:     1000,
	}
}

// CollectMetrics collects all metrics from the cluster
func (mc *MetricsCollector) CollectMetrics(ctx context.Context) (*MetricsSnapshot, error) {
	snapshot := &MetricsSnapshot{
		Timestamp: time.Now(),
	}

	// Collect from multiple sources in parallel
	var wg sync.WaitGroup
	var mu sync.Mutex
	errors := make([]error, 0)

	// Collect database metrics
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := mc.collectDatabaseMetrics(ctx, snapshot); err != nil {
			mu.Lock()
			errors = append(errors, fmt.Errorf("database metrics: %w", err))
			mu.Unlock()
		}
	}()

	// Collect latency metrics
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := mc.collectLatencyMetrics(ctx, snapshot); err != nil {
			mu.Lock()
			errors = append(errors, fmt.Errorf("latency metrics: %w", err))
			mu.Unlock()
		}
	}()

	// Collect Raft metrics
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := mc.collectRaftMetrics(ctx, snapshot); err != nil {
			mu.Lock()
			errors = append(errors, fmt.Errorf("raft metrics: %w", err))
			mu.Unlock()
		}
	}()

	// Collect client metrics
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := mc.collectClientMetrics(ctx, snapshot); err != nil {
			mu.Lock()
			errors = append(errors, fmt.Errorf("client metrics: %w", err))
			mu.Unlock()
		}
	}()

	wg.Wait()

	if len(errors) > 0 {
		mc.logger.Warn("Some metrics collection failed",
			zap.Int("error_count", len(errors)))
	}

	return snapshot, nil
}

// collectDatabaseMetrics collects database size and related metrics
func (mc *MetricsCollector) collectDatabaseMetrics(ctx context.Context, snapshot *MetricsSnapshot) error {
	// Get status from the leader
	membersResp, err := mc.client.MemberList(ctx)
	if err != nil {
		return fmt.Errorf("failed to get member list: %w", err)
	}

	// Find the leader and get its status
	for _, member := range membersResp.Members {
		if len(member.ClientURLs) == 0 {
			continue
		}

		statusResp, err := mc.client.Status(ctx, member.ClientURLs[0])
		if err != nil {
			continue
		}

		if statusResp.Leader == member.ID {
			snapshot.DBSize = statusResp.DbSize
			snapshot.DBSizeInUse = statusResp.DbSizeInUse
			break
		}
	}

	return nil
}

// collectLatencyMetrics measures and collects latency metrics
func (mc *MetricsCollector) collectLatencyMetrics(ctx context.Context, snapshot *MetricsSnapshot) error {
	// Measure read latency
	readLatencies := make([]float64, 0, 10)
	for i := 0; i < 10; i++ {
		start := time.Now()
		_, err := mc.client.Get(ctx, "/etcd-monitor/health-check")
		duration := time.Since(start)

		if err == nil {
			readLatencies = append(readLatencies, float64(duration.Milliseconds()))
		}
	}

	// Calculate percentiles for reads
	if len(readLatencies) > 0 {
		sort.Float64s(readLatencies)
		snapshot.ReadLatencyP50 = percentile(readLatencies, 0.50)
		snapshot.ReadLatencyP95 = percentile(readLatencies, 0.95)
		snapshot.ReadLatencyP99 = percentile(readLatencies, 0.99)
	}

	// Measure write latency
	writeLatencies := make([]float64, 0, 10)
	testKey := fmt.Sprintf("/etcd-monitor/latency-test-%d", time.Now().Unix())
	for i := 0; i < 10; i++ {
		start := time.Now()
		_, err := mc.client.Put(ctx, testKey, "test-value")
		duration := time.Since(start)

		if err == nil {
			writeLatencies = append(writeLatencies, float64(duration.Milliseconds()))
		}
	}

	// Clean up test key
	_, _ = mc.client.Delete(ctx, testKey)

	// Calculate percentiles for writes
	if len(writeLatencies) > 0 {
		sort.Float64s(writeLatencies)
		snapshot.WriteLatencyP50 = percentile(writeLatencies, 0.50)
		snapshot.WriteLatencyP95 = percentile(writeLatencies, 0.95)
		snapshot.WriteLatencyP99 = percentile(writeLatencies, 0.99)
	}

	// Record measurement
	mc.recordLatencyMeasurement(LatencyMeasurement{
		Timestamp:    time.Now(),
		ReadLatency:  time.Duration(snapshot.ReadLatencyP50) * time.Millisecond,
		WriteLatency: time.Duration(snapshot.WriteLatencyP50) * time.Millisecond,
		Operation:    "health-check",
	})

	return nil
}

// collectRaftMetrics collects Raft consensus metrics
func (mc *MetricsCollector) collectRaftMetrics(ctx context.Context, snapshot *MetricsSnapshot) error {
	// Get status from each member to aggregate Raft stats
	membersResp, err := mc.client.MemberList(ctx)
	if err != nil {
		return fmt.Errorf("failed to get member list: %w", err)
	}

	for _, member := range membersResp.Members {
		if len(member.ClientURLs) == 0 {
			continue
		}

		statusResp, err := mc.client.Status(ctx, member.ClientURLs[0])
		if err != nil {
			mc.logger.Warn("Failed to get member status for Raft metrics",
				zap.Uint64("member_id", member.ID),
				zap.Error(err))
			continue
		}

		// Aggregate Raft statistics
		snapshot.ProposalCommitted += statusResp.RaftAppliedIndex
		snapshot.ProposalApplied += statusResp.RaftAppliedIndex
		// Note: Pending and failed proposals would need to be collected from Prometheus metrics
	}

	return nil
}

// collectClientMetrics collects client connection metrics
func (mc *MetricsCollector) collectClientMetrics(ctx context.Context, snapshot *MetricsSnapshot) error {
	// Get active client connections
	// This would typically come from Prometheus metrics endpoint
	// For now, we'll use a placeholder implementation

	// Count active watchers by checking endpoint status
	membersResp, err := mc.client.MemberList(ctx)
	if err != nil {
		return fmt.Errorf("failed to get member list: %w", err)
	}

	snapshot.ActiveConnections = len(membersResp.Members)

	return nil
}

// recordLatencyMeasurement records a latency measurement
func (mc *MetricsCollector) recordLatencyMeasurement(measurement LatencyMeasurement) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.latencyHistory = append(mc.latencyHistory, measurement)

	// Keep only the last maxHistory measurements
	if len(mc.latencyHistory) > mc.maxHistory {
		mc.latencyHistory = mc.latencyHistory[1:]
	}
}

// GetLatencyHistory returns the latency measurement history
func (mc *MetricsCollector) GetLatencyHistory() []LatencyMeasurement {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	// Return a copy
	history := make([]LatencyMeasurement, len(mc.latencyHistory))
	copy(history, mc.latencyHistory)
	return history
}

// CalculateRequestRate calculates the request rate (ops/sec)
func (mc *MetricsCollector) CalculateRequestRate(ctx context.Context, duration time.Duration) (float64, error) {
	startTime := time.Now()

	// Get initial metrics
	initialMetrics, err := mc.CollectMetrics(ctx)
	if err != nil {
		return 0, err
	}

	// Wait for the specified duration
	time.Sleep(duration)

	// Get final metrics
	finalMetrics, err := mc.CollectMetrics(ctx)
	if err != nil {
		return 0, err
	}

	// Calculate rate
	elapsed := time.Since(startTime).Seconds()
	totalOps := float64(finalMetrics.ProposalCommitted - initialMetrics.ProposalCommitted)
	rate := totalOps / elapsed

	return rate, nil
}

// MeasurePerformance performs a comprehensive performance measurement
func (mc *MetricsCollector) MeasurePerformance(ctx context.Context, operations int) (*PerformanceReport, error) {
	report := &PerformanceReport{
		StartTime:  time.Now(),
		Operations: operations,
	}

	readLatencies := make([]float64, 0, operations)
	writeLatencies := make([]float64, 0, operations)

	// Perform read operations
	for i := 0; i < operations; i++ {
		start := time.Now()
		_, err := mc.client.Get(ctx, fmt.Sprintf("/perf-test/key-%d", i))
		duration := time.Since(start)

		if err == nil {
			readLatencies = append(readLatencies, float64(duration.Microseconds()))
		} else {
			report.ReadErrors++
		}
	}

	// Perform write operations
	for i := 0; i < operations; i++ {
		start := time.Now()
		_, err := mc.client.Put(ctx, fmt.Sprintf("/perf-test/key-%d", i), "test-value")
		duration := time.Since(start)

		if err == nil {
			writeLatencies = append(writeLatencies, float64(duration.Microseconds()))
		} else {
			report.WriteErrors++
		}
	}

	// Calculate statistics
	if len(readLatencies) > 0 {
		sort.Float64s(readLatencies)
		report.ReadLatencyP50 = percentile(readLatencies, 0.50) / 1000  // Convert to ms
		report.ReadLatencyP95 = percentile(readLatencies, 0.95) / 1000
		report.ReadLatencyP99 = percentile(readLatencies, 0.99) / 1000
	}

	if len(writeLatencies) > 0 {
		sort.Float64s(writeLatencies)
		report.WriteLatencyP50 = percentile(writeLatencies, 0.50) / 1000
		report.WriteLatencyP95 = percentile(writeLatencies, 0.95) / 1000
		report.WriteLatencyP99 = percentile(writeLatencies, 0.99) / 1000
	}

	report.EndTime = time.Now()
	report.Duration = report.EndTime.Sub(report.StartTime)
	report.Throughput = float64(operations*2) / report.Duration.Seconds() // reads + writes

	return report, nil
}

// PerformanceReport contains performance test results
type PerformanceReport struct {
	StartTime       time.Time
	EndTime         time.Time
	Duration        time.Duration
	Operations      int
	Throughput      float64
	ReadLatencyP50  float64
	ReadLatencyP95  float64
	ReadLatencyP99  float64
	WriteLatencyP50 float64
	WriteLatencyP95 float64
	WriteLatencyP99 float64
	ReadErrors      int
	WriteErrors     int
}

// percentile calculates the percentile value from a sorted slice
func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}

	index := int(float64(len(sorted)-1) * p)
	if index < 0 {
		index = 0
	}
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	return sorted[index]
}
