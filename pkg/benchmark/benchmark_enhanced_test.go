package benchmark

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewRunner(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Create runner with nil client", func(t *testing.T) {
		config := &Config{Type: BenchmarkTypeWrite}
		runner := NewRunner(nil, config, logger)
		assert.NotNil(t, runner)
		assert.Nil(t, runner.client)
		assert.Equal(t, config, runner.config)
		assert.Equal(t, logger, runner.logger)
	})
}

func TestRunner_Run_ErrorCases(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Run with nil client returns error", func(t *testing.T) {
		config := &Config{
			Type:            BenchmarkTypeWrite,
			Connections:     1,
			Clients:         1,
			TotalOperations: 10,
			KeyPrefix:       "/bench",
		}

		runner := NewRunner(nil, config, logger)
		result, err := runner.Run(context.Background())
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "etcd client is nil")
	})
}

func TestCalculateLatencyStats(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := &Config{Type: BenchmarkTypeWrite}
	runner := NewRunner(nil, config, logger)

	t.Run("Calculate stats with empty latencies", func(t *testing.T) {
		result := &Result{
			StartTime:        time.Now(),
			EndTime:          time.Now().Add(1 * time.Second),
			Duration:         1 * time.Second,
			SuccessfulOps:    0,
			LatencyHistogram: make(map[string]int),
		}

		runner.calculateLatencyStats([]float64{}, result)

		assert.Equal(t, 0.0, result.P50Latency)
		assert.Equal(t, 0.0, result.P95Latency)
		assert.Equal(t, 0.0, result.P99Latency)
	})

	t.Run("Calculate stats with sample latencies", func(t *testing.T) {
		latencies := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
		result := &Result{
			StartTime:        time.Now(),
			EndTime:          time.Now().Add(1 * time.Second),
			Duration:         1 * time.Second,
			SuccessfulOps:    10,
			LatencyHistogram: make(map[string]int),
		}

		runner.calculateLatencyStats(latencies, result)

		assert.Greater(t, result.P50Latency, 0.0)
		assert.Greater(t, result.P95Latency, 0.0)
		assert.Greater(t, result.P99Latency, 0.0)
		assert.Equal(t, 1.0, result.MinLatency)
		assert.Equal(t, 10.0, result.MaxLatency)
		assert.Equal(t, 5.5, result.AvgLatency)
		assert.Equal(t, 10.0, result.Throughput) // 10 ops / 1 second
	})

	t.Run("Calculate stats builds histogram", func(t *testing.T) {
		latencies := []float64{0.5, 2.0, 7.0, 15.0, 75.0, 250.0, 600.0}
		result := &Result{
			StartTime:        time.Now(),
			EndTime:          time.Now().Add(1 * time.Second),
			Duration:         1 * time.Second,
			SuccessfulOps:    7,
			LatencyHistogram: make(map[string]int),
		}

		runner.calculateLatencyStats(latencies, result)

		assert.Greater(t, len(result.LatencyHistogram), 0)
		assert.Contains(t, result.LatencyHistogram, "<1ms")
		assert.Contains(t, result.LatencyHistogram, "1-5ms")
	})
}

func TestPercentile(t *testing.T) {
	t.Run("Percentile with empty slice", func(t *testing.T) {
		result := percentile([]float64{}, 0.50)
		assert.Equal(t, 0.0, result)
	})

	t.Run("Percentile with single value", func(t *testing.T) {
		result := percentile([]float64{5.5}, 0.50)
		assert.Equal(t, 5.5, result)
	})

	t.Run("Percentile P50", func(t *testing.T) {
		values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		result := percentile(values, 0.50)
		assert.Equal(t, 5.0, result)
	})

	t.Run("Percentile P95", func(t *testing.T) {
		values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		result := percentile(values, 0.95)
		assert.Equal(t, 9.0, result)
	})

	t.Run("Percentile P99", func(t *testing.T) {
		values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		result := percentile(values, 0.99)
		assert.Equal(t, 9.0, result)
	})
}

func TestGetLatencyBucket(t *testing.T) {
	testCases := []struct {
		latency  float64
		expected string
	}{
		{0.5, "<1ms"},
		{2.0, "1-5ms"},
		{7.0, "5-10ms"},
		{15.0, "10-50ms"},
		{75.0, "50-100ms"},
		{250.0, "100-500ms"},
		{600.0, ">500ms"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			bucket := getLatencyBucket(tc.latency)
			assert.Equal(t, tc.expected, bucket)
		})
	}
}

func TestGenerateRandomString(t *testing.T) {
	t.Run("Generate string of specified length", func(t *testing.T) {
		str := generateRandomString(32)
		assert.Equal(t, 32, len(str))
	})

	t.Run("Generate zero-length string", func(t *testing.T) {
		str := generateRandomString(0)
		assert.Equal(t, 0, len(str))
	})

	t.Run("Generate different strings", func(t *testing.T) {
		str1 := generateRandomString(100)
		str2 := generateRandomString(100)
		// Statistically unlikely to be the same
		assert.NotEqual(t, str1, str2)
	})

	t.Run("Generate string with valid characters", func(t *testing.T) {
		str := generateRandomString(50)
		for _, char := range str {
			assert.True(t, (char >= 'a' && char <= 'z') ||
				(char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9'))
		}
	})
}

func TestPrintResult(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Print result without panic", func(t *testing.T) {
		result := &Result{
			Type:          BenchmarkTypeWrite,
			StartTime:     time.Now(),
			EndTime:       time.Now().Add(1 * time.Second),
			Duration:      1 * time.Second,
			TotalOperations: 100,
			SuccessfulOps: 95,
			FailedOps:     5,
			Throughput:    95.0,
			AvgLatency:    10.5,
			P50Latency:    8.0,
			P95Latency:    20.0,
			P99Latency:    35.0,
		}

		assert.NotPanics(t, func() {
			PrintResult(result, logger)
		})
	})
}

func TestConfig_Structure(t *testing.T) {
	t.Run("Config with all fields", func(t *testing.T) {
		config := &Config{
			Type:            BenchmarkTypeWrite,
			Connections:     5,
			Clients:         10,
			KeySize:         32,
			ValueSize:       256,
			TotalOperations: 10000,
			Duration:        60 * time.Second,
			KeyPrefix:       "/benchmark",
			TargetLeader:    true,
			RateLimit:       1000,
		}

		assert.Equal(t, BenchmarkTypeWrite, config.Type)
		assert.Equal(t, 5, config.Connections)
		assert.Equal(t, 10, config.Clients)
		assert.Equal(t, 32, config.KeySize)
		assert.Equal(t, 256, config.ValueSize)
		assert.Equal(t, 10000, config.TotalOperations)
		assert.Equal(t, 60*time.Second, config.Duration)
		assert.Equal(t, "/benchmark", config.KeyPrefix)
		assert.True(t, config.TargetLeader)
		assert.Equal(t, 1000, config.RateLimit)
	})
}

func TestResult_Structure(t *testing.T) {
	t.Run("Result with all fields", func(t *testing.T) {
		start := time.Now()
		end := start.Add(5 * time.Second)

		result := &Result{
			Type:             BenchmarkTypeMixed,
			StartTime:        start,
			EndTime:          end,
			Duration:         end.Sub(start),
			TotalOperations:  500,
			SuccessfulOps:    480,
			FailedOps:        20,
			Throughput:       96.0,
			AvgLatency:       12.5,
			MinLatency:       1.0,
			MaxLatency:       150.0,
			P50Latency:       10.0,
			P95Latency:       45.0,
			P99Latency:       100.0,
			LatencyHistogram: map[string]int{"<1ms": 50, "1-5ms": 200},
			ErrorTypes:       map[string]int{"timeout": 10, "unavailable": 10},
		}

		assert.Equal(t, BenchmarkTypeMixed, result.Type)
		assert.Equal(t, 5*time.Second, result.Duration)
		assert.Equal(t, 500, result.TotalOperations)
		assert.Equal(t, 480, result.SuccessfulOps)
		assert.Equal(t, 20, result.FailedOps)
		assert.Equal(t, 96.0, result.Throughput)
		assert.Equal(t, 12.5, result.AvgLatency)
		assert.Equal(t, 2, len(result.LatencyHistogram))
		assert.Equal(t, 2, len(result.ErrorTypes))
	})
}

func TestBenchmarkType_Constants(t *testing.T) {
	t.Run("All benchmark types defined", func(t *testing.T) {
		types := []BenchmarkType{
			BenchmarkTypeWrite,
			BenchmarkTypeRead,
			BenchmarkTypeMixed,
			BenchmarkTypeRange,
			BenchmarkTypeDelete,
			BenchmarkTypeTxn,
		}

		assert.Equal(t, 6, len(types))
		assert.Equal(t, BenchmarkType("write"), BenchmarkTypeWrite)
		assert.Equal(t, BenchmarkType("read"), BenchmarkTypeRead)
		assert.Equal(t, BenchmarkType("mixed"), BenchmarkTypeMixed)
		assert.Equal(t, BenchmarkType("range"), BenchmarkTypeRange)
		assert.Equal(t, BenchmarkType("delete"), BenchmarkTypeDelete)
		assert.Equal(t, BenchmarkType("txn"), BenchmarkTypeTxn)
	})
}
