// Package benchmark provides performance benchmarking for etcd clusters
package benchmark

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// BenchmarkType defines the type of benchmark
type BenchmarkType string

const (
	BenchmarkTypeWrite      BenchmarkType = "write"
	BenchmarkTypeRead       BenchmarkType = "read"
	BenchmarkTypeMixed      BenchmarkType = "mixed"
	BenchmarkTypeRange      BenchmarkType = "range"
	BenchmarkTypeDelete     BenchmarkType = "delete"
	BenchmarkTypeTxn        BenchmarkType = "txn"
)

// Config holds benchmark configuration
type Config struct {
	Type           BenchmarkType
	Connections    int           // Number of concurrent connections
	Clients        int           // Number of clients per connection
	KeySize        int           // Size of keys in bytes
	ValueSize      int           // Size of values in bytes
	TotalOperations int          // Total number of operations
	Duration       time.Duration // Duration of benchmark (alternative to TotalOperations)
	KeyPrefix      string        // Prefix for test keys
	TargetLeader   bool          // Whether to target leader only
	RateLimit      int           // Max operations per second (0 = unlimited)
}

// Result contains benchmark results
type Result struct {
	Type               BenchmarkType
	StartTime          time.Time
	EndTime            time.Time
	Duration           time.Duration
	TotalOperations    int
	SuccessfulOps      int
	FailedOps          int
	Throughput         float64 // ops/sec
	AvgLatency         float64 // ms
	MinLatency         float64 // ms
	MaxLatency         float64 // ms
	P50Latency         float64 // ms
	P95Latency         float64 // ms
	P99Latency         float64 // ms
	LatencyHistogram   map[string]int
	ErrorTypes         map[string]int
}

// Runner executes benchmarks
type Runner struct {
	client *clientv3.Client
	config *Config
	logger *zap.Logger
}

// NewRunner creates a new benchmark runner
func NewRunner(client *clientv3.Client, config *Config, logger *zap.Logger) *Runner {
	return &Runner{
		client: client,
		config: config,
		logger: logger,
	}
}

// Run executes the benchmark
func (r *Runner) Run(ctx context.Context) (*Result, error) {
	r.logger.Info("Starting benchmark",
		zap.String("type", string(r.config.Type)),
		zap.Int("connections", r.config.Connections),
		zap.Int("clients", r.config.Clients),
		zap.Int("operations", r.config.TotalOperations))

	result := &Result{
		Type:             r.config.Type,
		StartTime:        time.Now(),
		LatencyHistogram: make(map[string]int),
		ErrorTypes:       make(map[string]int),
	}

	// Execute benchmark based on type
	var err error
	switch r.config.Type {
	case BenchmarkTypeWrite:
		err = r.runWriteBenchmark(ctx, result)
	case BenchmarkTypeRead:
		err = r.runReadBenchmark(ctx, result)
	case BenchmarkTypeMixed:
		err = r.runMixedBenchmark(ctx, result)
	default:
		return nil, fmt.Errorf("unsupported benchmark type: %s", r.config.Type)
	}

	if err != nil {
		return nil, err
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	return result, nil
}

// runWriteBenchmark performs write operations benchmark
func (r *Runner) runWriteBenchmark(ctx context.Context, result *Result) error {
	var (
		wg             sync.WaitGroup
		latencies      []float64
		latenciesMutex sync.Mutex
		successOps     int64
		failedOps      int64
		opsPerClient   = r.config.TotalOperations / (r.config.Connections * r.config.Clients)
	)

	// Create workers
	for conn := 0; conn < r.config.Connections; conn++ {
		for client := 0; client < r.config.Clients; client++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				for i := 0; i < opsPerClient; i++ {
					key := fmt.Sprintf("%s/write-%d-%d", r.config.KeyPrefix, workerID, i)
					value := generateRandomString(r.config.ValueSize)

					start := time.Now()
					_, err := r.client.Put(ctx, key, value)
					latency := time.Since(start)

					if err != nil {
						atomic.AddInt64(&failedOps, 1)
						r.logger.Debug("Write operation failed", zap.Error(err))
					} else {
						atomic.AddInt64(&successOps, 1)
						latenciesMutex.Lock()
						latencies = append(latencies, float64(latency.Microseconds())/1000.0)
						latenciesMutex.Unlock()
					}

					// Rate limiting
					if r.config.RateLimit > 0 {
						time.Sleep(time.Second / time.Duration(r.config.RateLimit))
					}
				}
			}(conn*r.config.Clients + client)
		}
	}

	wg.Wait()

	// Calculate statistics
	result.SuccessfulOps = int(successOps)
	result.FailedOps = int(failedOps)
	result.TotalOperations = result.SuccessfulOps + result.FailedOps
	r.calculateLatencyStats(latencies, result)

	return nil
}

// runReadBenchmark performs read operations benchmark
func (r *Runner) runReadBenchmark(ctx context.Context, result *Result) error {
	// First, populate keys to read
	numKeys := 10000
	r.logger.Info("Populating keys for read benchmark", zap.Int("count", numKeys))

	for i := 0; i < numKeys; i++ {
		key := fmt.Sprintf("%s/read-%d", r.config.KeyPrefix, i)
		value := generateRandomString(r.config.ValueSize)
		_, err := r.client.Put(ctx, key, value)
		if err != nil {
			return fmt.Errorf("failed to populate keys: %w", err)
		}
	}

	// Now perform read benchmark
	var (
		wg             sync.WaitGroup
		latencies      []float64
		latenciesMutex sync.Mutex
		successOps     int64
		failedOps      int64
		opsPerClient   = r.config.TotalOperations / (r.config.Connections * r.config.Clients)
	)

	for conn := 0; conn < r.config.Connections; conn++ {
		for client := 0; client < r.config.Clients; client++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for i := 0; i < opsPerClient; i++ {
					// Random key selection
					keyIndex := rand.Intn(numKeys)
					key := fmt.Sprintf("%s/read-%d", r.config.KeyPrefix, keyIndex)

					start := time.Now()
					_, err := r.client.Get(ctx, key)
					latency := time.Since(start)

					if err != nil {
						atomic.AddInt64(&failedOps, 1)
					} else {
						atomic.AddInt64(&successOps, 1)
						latenciesMutex.Lock()
						latencies = append(latencies, float64(latency.Microseconds())/1000.0)
						latenciesMutex.Unlock()
					}

					// Rate limiting
					if r.config.RateLimit > 0 {
						time.Sleep(time.Second / time.Duration(r.config.RateLimit))
					}
				}
			}()
		}
	}

	wg.Wait()

	// Cleanup
	r.logger.Info("Cleaning up test keys")
	_, _ = r.client.Delete(ctx, r.config.KeyPrefix+"/read-", clientv3.WithPrefix())

	// Calculate statistics
	result.SuccessfulOps = int(successOps)
	result.FailedOps = int(failedOps)
	result.TotalOperations = result.SuccessfulOps + result.FailedOps
	r.calculateLatencyStats(latencies, result)

	return nil
}

// runMixedBenchmark performs mixed read/write operations
func (r *Runner) runMixedBenchmark(ctx context.Context, result *Result) error {
	// 70% reads, 30% writes
	readRatio := 0.7

	var (
		wg             sync.WaitGroup
		latencies      []float64
		latenciesMutex sync.Mutex
		successOps     int64
		failedOps      int64
		opsPerClient   = r.config.TotalOperations / (r.config.Connections * r.config.Clients)
	)

	for conn := 0; conn < r.config.Connections; conn++ {
		for client := 0; client < r.config.Clients; client++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				for i := 0; i < opsPerClient; i++ {
					key := fmt.Sprintf("%s/mixed-%d-%d", r.config.KeyPrefix, workerID, i%1000)

					var err error
					start := time.Now()

					// Decide read or write
					if rand.Float64() < readRatio {
						_, err = r.client.Get(ctx, key)
					} else {
						value := generateRandomString(r.config.ValueSize)
						_, err = r.client.Put(ctx, key, value)
					}

					latency := time.Since(start)

					if err != nil {
						atomic.AddInt64(&failedOps, 1)
					} else {
						atomic.AddInt64(&successOps, 1)
						latenciesMutex.Lock()
						latencies = append(latencies, float64(latency.Microseconds())/1000.0)
						latenciesMutex.Unlock()
					}

					// Rate limiting
					if r.config.RateLimit > 0 {
						time.Sleep(time.Second / time.Duration(r.config.RateLimit))
					}
				}
			}(conn*r.config.Clients + client)
		}
	}

	wg.Wait()

	// Cleanup
	_, _ = r.client.Delete(ctx, r.config.KeyPrefix+"/mixed-", clientv3.WithPrefix())

	// Calculate statistics
	result.SuccessfulOps = int(successOps)
	result.FailedOps = int(failedOps)
	result.TotalOperations = result.SuccessfulOps + result.FailedOps
	r.calculateLatencyStats(latencies, result)

	return nil
}

// calculateLatencyStats calculates latency statistics
func (r *Runner) calculateLatencyStats(latencies []float64, result *Result) {
	if len(latencies) == 0 {
		return
	}

	sort.Float64s(latencies)

	// Calculate percentiles
	result.P50Latency = percentile(latencies, 0.50)
	result.P95Latency = percentile(latencies, 0.95)
	result.P99Latency = percentile(latencies, 0.99)
	result.MinLatency = latencies[0]
	result.MaxLatency = latencies[len(latencies)-1]

	// Calculate average
	sum := 0.0
	for _, lat := range latencies {
		sum += lat
	}
	result.AvgLatency = sum / float64(len(latencies))

	// Calculate throughput
	result.Throughput = float64(result.SuccessfulOps) / result.Duration.Seconds()

	// Build histogram
	for _, lat := range latencies {
		bucket := getLatencyBucket(lat)
		result.LatencyHistogram[bucket]++
	}
}

// percentile calculates the percentile value
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

// getLatencyBucket returns the histogram bucket for a latency value
func getLatencyBucket(latency float64) string {
	switch {
	case latency < 1:
		return "<1ms"
	case latency < 5:
		return "1-5ms"
	case latency < 10:
		return "5-10ms"
	case latency < 50:
		return "10-50ms"
	case latency < 100:
		return "50-100ms"
	case latency < 500:
		return "100-500ms"
	default:
		return ">500ms"
	}
}

// generateRandomString generates a random string of specified length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// PrintResult prints the benchmark result
func PrintResult(result *Result, logger *zap.Logger) {
	logger.Info("Benchmark completed",
		zap.String("type", string(result.Type)),
		zap.Duration("duration", result.Duration),
		zap.Int("total_ops", result.TotalOperations),
		zap.Int("successful_ops", result.SuccessfulOps),
		zap.Int("failed_ops", result.FailedOps),
		zap.Float64("throughput", result.Throughput),
		zap.Float64("avg_latency", result.AvgLatency),
		zap.Float64("p50_latency", result.P50Latency),
		zap.Float64("p95_latency", result.P95Latency),
		zap.Float64("p99_latency", result.P99Latency))
}
