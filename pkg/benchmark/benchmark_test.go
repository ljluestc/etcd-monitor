package benchmark

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestBenchmark_Comprehensive(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("NewRunner", func(t *testing.T) {
		config := &Config{
			Type:            BenchmarkTypeWrite,
			Connections:     1,
			Clients:         1,
			KeySize:         10,
			ValueSize:       100,
			TotalOperations: 10,
			KeyPrefix:       "test",
		}

		runner := NewRunner(nil, config, logger)
		
		assert.NotNil(t, runner)
		assert.Equal(t, config, runner.config)
		assert.Equal(t, logger, runner.logger)
	})

	t.Run("Config validation", func(t *testing.T) {
		config := &Config{
			Type:            BenchmarkTypeWrite,
			Connections:     1,
			Clients:         1,
			KeySize:         10,
			ValueSize:       100,
			TotalOperations: 10,
			KeyPrefix:       "test",
		}

		assert.Equal(t, BenchmarkTypeWrite, config.Type)
		assert.Equal(t, 1, config.Connections)
		assert.Equal(t, 1, config.Clients)
		assert.Equal(t, 10, config.KeySize)
		assert.Equal(t, 100, config.ValueSize)
		assert.Equal(t, 10, config.TotalOperations)
		assert.Equal(t, "test", config.KeyPrefix)
	})

	t.Run("BenchmarkType constants", func(t *testing.T) {
		assert.Equal(t, BenchmarkType("write"), BenchmarkTypeWrite)
		assert.Equal(t, BenchmarkType("read"), BenchmarkTypeRead)
		assert.Equal(t, BenchmarkType("mixed"), BenchmarkTypeMixed)
		assert.Equal(t, BenchmarkType("range"), BenchmarkTypeRange)
		assert.Equal(t, BenchmarkType("delete"), BenchmarkTypeDelete)
		assert.Equal(t, BenchmarkType("txn"), BenchmarkTypeTxn)
	})

	t.Run("Result structure", func(t *testing.T) {
		result := &Result{
			Type:             BenchmarkTypeWrite,
			StartTime:        time.Now(),
			EndTime:          time.Now().Add(time.Second),
			Duration:         time.Second,
			TotalOperations:  100,
			SuccessfulOps:    95,
			FailedOps:        5,
			Throughput:       100.0,
			AvgLatency:       10.5,
			MinLatency:       1.0,
			MaxLatency:       50.0,
			P50Latency:       8.0,
			P95Latency:       25.0,
			P99Latency:       45.0,
			LatencyHistogram: make(map[string]int),
			ErrorTypes:       make(map[string]int),
		}

		assert.Equal(t, BenchmarkTypeWrite, result.Type)
		assert.Equal(t, 100, result.TotalOperations)
		assert.Equal(t, 95, result.SuccessfulOps)
		assert.Equal(t, 5, result.FailedOps)
		assert.Equal(t, 100.0, result.Throughput)
		assert.Equal(t, 10.5, result.AvgLatency)
	})

	t.Run("generateRandomString", func(t *testing.T) {
		str := generateRandomString(10)
		assert.Len(t, str, 10)
		
		str2 := generateRandomString(10)
		assert.Len(t, str2, 10)
		// Should be different (very high probability)
		assert.NotEqual(t, str, str2)
	})

	t.Run("getLatencyBucket", func(t *testing.T) {
		assert.Equal(t, "<1ms", getLatencyBucket(0.5))
		assert.Equal(t, "1-5ms", getLatencyBucket(2.0))
		assert.Equal(t, "5-10ms", getLatencyBucket(7.0))
		assert.Equal(t, "10-50ms", getLatencyBucket(25.0))
		assert.Equal(t, "50-100ms", getLatencyBucket(75.0))
		assert.Equal(t, "100-500ms", getLatencyBucket(200.0))
		assert.Equal(t, ">500ms", getLatencyBucket(600.0))
	})

	t.Run("percentile calculation", func(t *testing.T) {
		latencies := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
		
		assert.Equal(t, 5.0, percentile(latencies, 0.5))  // P50
		assert.Equal(t, 9.0, percentile(latencies, 0.95)) // P95
		assert.Equal(t, 9.0, percentile(latencies, 0.99)) // P99
	})

	t.Run("percentile with empty slice", func(t *testing.T) {
		latencies := []float64{}
		assert.Equal(t, 0.0, percentile(latencies, 0.5))
	})

	t.Run("PrintResult with valid result", func(t *testing.T) {
		result := &Result{
			Type:            BenchmarkTypeWrite,
			Duration:        time.Second,
			TotalOperations: 100,
			SuccessfulOps:   95,
			FailedOps:       5,
			Throughput:      100.0,
			AvgLatency:      10.5,
			P50Latency:      8.0,
			P95Latency:      25.0,
			P99Latency:      45.0,
		}

		// This should not panic
		PrintResult(result, logger)
	})

	t.Run("Runner with nil client should error", func(t *testing.T) {
		config := &Config{
			Type:            BenchmarkTypeWrite,
			Connections:     1,
			Clients:         1,
			TotalOperations: 1,
		}

		runner := NewRunner(nil, config, logger)
		
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		
		_, err := runner.Run(ctx)
		assert.Error(t, err) // Should error due to nil client
		assert.Contains(t, err.Error(), "etcd client is nil")
	})
}