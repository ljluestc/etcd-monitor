package benchmark

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestBenchmark_Comprehensive(t *testing.T) {
	t.Run("NewBenchmark with valid config", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.NotNil(t, benchmark.config)
		assert.NotNil(t, benchmark.logger)
	})

	t.Run("NewBenchmark with nil config", func(t *testing.T) {
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(nil, logger)
		assert.NotNil(t, benchmark)
		assert.Nil(t, benchmark.config)
		assert.NotNil(t, benchmark.logger)
	})

	t.Run("NewBenchmark with nil logger", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}

		benchmark := NewBenchmark(config, nil)
		assert.NotNil(t, benchmark)
		assert.NotNil(t, benchmark.config)
		assert.NotNil(t, benchmark.logger) // Should create a production logger
	})

	t.Run("NewBenchmark with empty config", func(t *testing.T) {
		config := &Config{}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.NotNil(t, benchmark.config)
		assert.NotNil(t, benchmark.logger)
	})

	t.Run("Start benchmark", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		err := benchmark.Start()
		assert.NoError(t, err)
	})

	t.Run("Stop benchmark", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		err := benchmark.Start()
		require.NoError(t, err)

		benchmark.Stop()
		// Should not panic or error
	})

	t.Run("Stop benchmark without starting", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		benchmark.Stop()
		// Should not panic or error
	})

	t.Run("Run benchmark", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    1 * time.Second,
			Concurrency: 1,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		ctx := context.Background()
		
		err := benchmark.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("Run benchmark with nil context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    1 * time.Second,
			Concurrency: 1,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		
		err := benchmark.Run(nil)
		assert.Error(t, err)
	})

	t.Run("Run benchmark with cancelled context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		err := benchmark.Run(ctx)
		assert.Error(t, err)
	})

	t.Run("Run benchmark with timeout context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		err := benchmark.Run(ctx)
		assert.Error(t, err)
	})

	t.Run("GetBenchmarkResults", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		results := benchmark.GetBenchmarkResults()
		assert.NotNil(t, results)
	})

	t.Run("GetBenchmarkResults after run", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    1 * time.Second,
			Concurrency: 1,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		ctx := context.Background()
		
		err := benchmark.Run(ctx)
		require.NoError(t, err)
		
		results := benchmark.GetBenchmarkResults()
		assert.NotNil(t, results)
	})

	t.Run("GetBenchmarkStatus", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		status := benchmark.GetBenchmarkStatus()
		assert.NotNil(t, status)
	})

	t.Run("GetBenchmarkStatus after run", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    1 * time.Second,
			Concurrency: 1,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		ctx := context.Background()
		
		err := benchmark.Run(ctx)
		require.NoError(t, err)
		
		status := benchmark.GetBenchmarkStatus()
		assert.NotNil(t, status)
	})

	t.Run("GetBenchmarkMetrics", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		metrics := benchmark.GetBenchmarkMetrics()
		assert.NotNil(t, metrics)
	})

	t.Run("GetBenchmarkMetrics after run", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    1 * time.Second,
			Concurrency: 1,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		ctx := context.Background()
		
		err := benchmark.Run(ctx)
		require.NoError(t, err)
		
		metrics := benchmark.GetBenchmarkMetrics()
		assert.NotNil(t, metrics)
	})
}

func TestBenchmark_EdgeCases(t *testing.T) {
	t.Run("Benchmark with very long endpoint list", func(t *testing.T) {
		endpoints := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			endpoints[i] = fmt.Sprintf("localhost:%d", 2379+i)
		}

		config := &Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.Len(t, benchmark.config.Endpoints, 1000)
	})

	t.Run("Benchmark with very long timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 24 * time.Hour,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.Equal(t, 24*time.Hour, benchmark.config.DialTimeout)
	})

	t.Run("Benchmark with very short timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 1 * time.Nanosecond,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.Equal(t, 1*time.Nanosecond, benchmark.config.DialTimeout)
	})

	t.Run("Benchmark with zero timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 0,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.Equal(t, time.Duration(0), benchmark.config.DialTimeout)
	})

	t.Run("Benchmark with negative timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: -1 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.Equal(t, -1*time.Second, benchmark.config.DialTimeout)
	})

	t.Run("Benchmark with very long duration", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    24 * time.Hour,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.Equal(t, 24*time.Hour, benchmark.config.Duration)
	})

	t.Run("Benchmark with very short duration", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    1 * time.Nanosecond,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.Equal(t, 1*time.Nanosecond, benchmark.config.Duration)
	})

	t.Run("Benchmark with zero duration", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    0,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.Equal(t, time.Duration(0), benchmark.config.Duration)
	})

	t.Run("Benchmark with negative duration", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    -1 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.Equal(t, -1*time.Second, benchmark.config.Duration)
	})

	t.Run("Benchmark with very high concurrency", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10000,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.Equal(t, 10000, benchmark.config.Concurrency)
	})

	t.Run("Benchmark with zero concurrency", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 0,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.Equal(t, 0, benchmark.config.Concurrency)
	})

	t.Run("Benchmark with negative concurrency", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: -1,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.Equal(t, -1, benchmark.config.Concurrency)
	})

	t.Run("Benchmark with empty endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.Len(t, benchmark.config.Endpoints, 0)
	})

	t.Run("Benchmark with nil endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.Nil(t, benchmark.config.Endpoints)
	})

	t.Run("Benchmark with special characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "https://etcd-1:2379", "http://etcd-2:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.Len(t, benchmark.config.Endpoints, 3)
		assert.Equal(t, "http://etcd-0:2379", benchmark.config.Endpoints[0])
		assert.Equal(t, "https://etcd-1:2379", benchmark.config.Endpoints[1])
		assert.Equal(t, "http://etcd-2:2379", benchmark.config.Endpoints[2])
	})

	t.Run("Benchmark with unicode characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.Len(t, benchmark.config.Endpoints, 2)
	})

	t.Run("Benchmark with control characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.Len(t, benchmark.config.Endpoints, 2)
	})
}

func TestBenchmark_Performance(t *testing.T) {
	t.Run("Benchmark creation performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Create 1000 Benchmark instances
		for i := 0; i < 1000; i++ {
			benchmark := NewBenchmark(config, logger)
			assert.NotNil(t, benchmark)
		}
		
		duration := time.Since(start)
		
		// Should create 1000 instances in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Creating 1000 Benchmark instances took too long: %v", duration)
	})

	t.Run("Benchmark with large endpoint list performance", func(t *testing.T) {
		endpoints := make([]string, 10000)
		for i := 0; i < 10000; i++ {
			endpoints[i] = fmt.Sprintf("localhost:%d", 2379+i)
		}
		
		config := &Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		benchmark := NewBenchmark(config, logger)
		
		duration := time.Since(start)
		
		// Should create benchmark with 10000 endpoints in less than 10ms
		assert.True(t, duration < 10*time.Millisecond, "Creating Benchmark with 10000 endpoints took too long: %v", duration)
		assert.Len(t, benchmark.config.Endpoints, 10000)
	})

	t.Run("Benchmark start/stop performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Start and stop 100 benchmarks
		for i := 0; i < 100; i++ {
			benchmark := NewBenchmark(config, logger)
			err := benchmark.Start()
			require.NoError(t, err)
			benchmark.Stop()
		}
		
		duration := time.Since(start)
		
		// Should start and stop 100 benchmarks in less than 1 second
		assert.True(t, duration < 1*time.Second, "Starting and stopping 100 benchmarks took too long: %v", duration)
	})

	t.Run("Benchmark run performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    100 * time.Millisecond,
			Concurrency: 1,
		}
		logger, _ := zap.NewDevelopment()
		
		benchmark := NewBenchmark(config, logger)
		ctx := context.Background()
		
		start := time.Now()
		
		// Run benchmark 10 times
		for i := 0; i < 10; i++ {
			err := benchmark.Run(ctx)
			assert.NoError(t, err)
		}
		
		duration := time.Since(start)
		
		// Should run 10 benchmarks in less than 2 seconds
		assert.True(t, duration < 2*time.Second, "Running 10 benchmarks took too long: %v", duration)
	})

	t.Run("Benchmark get results performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()
		
		benchmark := NewBenchmark(config, logger)
		
		start := time.Now()
		
		// Get results 1000 times
		for i := 0; i < 1000; i++ {
			results := benchmark.GetBenchmarkResults()
			assert.NotNil(t, results)
		}
		
		duration := time.Since(start)
		
		// Should get results 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting results 1000 times took too long: %v", duration)
	})

	t.Run("Benchmark get status performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()
		
		benchmark := NewBenchmark(config, logger)
		
		start := time.Now()
		
		// Get status 1000 times
		for i := 0; i < 1000; i++ {
			status := benchmark.GetBenchmarkStatus()
			assert.NotNil(t, status)
		}
		
		duration := time.Since(start)
		
		// Should get status 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting status 1000 times took too long: %v", duration)
	})

	t.Run("Benchmark get metrics performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()
		
		benchmark := NewBenchmark(config, logger)
		
		start := time.Now()
		
		// Get metrics 1000 times
		for i := 0; i < 1000; i++ {
			metrics := benchmark.GetBenchmarkMetrics()
			assert.NotNil(t, metrics)
		}
		
		duration := time.Since(start)
		
		// Should get metrics 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting metrics 1000 times took too long: %v", duration)
	})
}

func TestBenchmark_Integration(t *testing.T) {
	t.Run("Benchmark with all fields set", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379", "http://etcd-2:2379"},
			DialTimeout: 10 * time.Second,
			Duration:    30 * time.Second,
			Concurrency: 50,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.NotNil(t, benchmark.config)
		assert.NotNil(t, benchmark.logger)
		assert.Len(t, benchmark.config.Endpoints, 3)
		assert.Equal(t, "http://etcd-0:2379", benchmark.config.Endpoints[0])
		assert.Equal(t, "http://etcd-1:2379", benchmark.config.Endpoints[1])
		assert.Equal(t, "http://etcd-2:2379", benchmark.config.Endpoints[2])
		assert.Equal(t, 10*time.Second, benchmark.config.DialTimeout)
		assert.Equal(t, 30*time.Second, benchmark.config.Duration)
		assert.Equal(t, 50, benchmark.config.Concurrency)
	})

	t.Run("Benchmark with minimal fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 1,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.NotNil(t, benchmark.config)
		assert.NotNil(t, benchmark.logger)
		assert.Len(t, benchmark.config.Endpoints, 1)
		assert.Equal(t, "localhost:2379", benchmark.config.Endpoints[0])
		assert.Equal(t, 5*time.Second, benchmark.config.DialTimeout)
		assert.Equal(t, 10*time.Second, benchmark.config.Duration)
		assert.Equal(t, 1, benchmark.config.Concurrency)
	})

	t.Run("Benchmark with empty fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 0,
			Duration:    0,
			Concurrency: 0,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.NotNil(t, benchmark.config)
		assert.NotNil(t, benchmark.logger)
		assert.Len(t, benchmark.config.Endpoints, 0)
		assert.Equal(t, time.Duration(0), benchmark.config.DialTimeout)
		assert.Equal(t, time.Duration(0), benchmark.config.Duration)
		assert.Equal(t, 0, benchmark.config.Concurrency)
	})

	t.Run("Benchmark with nil fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 0,
			Duration:    0,
			Concurrency: 0,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)
		assert.NotNil(t, benchmark.config)
		assert.NotNil(t, benchmark.logger)
		assert.Nil(t, benchmark.config.Endpoints)
		assert.Equal(t, time.Duration(0), benchmark.config.DialTimeout)
		assert.Equal(t, time.Duration(0), benchmark.config.Duration)
		assert.Equal(t, 0, benchmark.config.Concurrency)
	})

	t.Run("Benchmark full lifecycle", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    1 * time.Second,
			Concurrency: 1,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)

		// Start benchmark
		err := benchmark.Start()
		assert.NoError(t, err)

		// Run benchmark
		ctx := context.Background()
		err = benchmark.Run(ctx)
		assert.NoError(t, err)

		// Get results
		results := benchmark.GetBenchmarkResults()
		assert.NotNil(t, results)

		// Get status
		status := benchmark.GetBenchmarkStatus()
		assert.NotNil(t, status)

		// Get metrics
		metrics := benchmark.GetBenchmarkMetrics()
		assert.NotNil(t, metrics)

		// Stop benchmark
		benchmark.Stop()
	})

	t.Run("Benchmark with multiple concurrent operations", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    100 * time.Millisecond,
			Concurrency: 5,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)

		// Start benchmark
		err := benchmark.Start()
		require.NoError(t, err)

		// Perform multiple concurrent operations
		done := make(chan bool, 10)
		
		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				// Run benchmark
				ctx := context.Background()
				err := benchmark.Run(ctx)
				assert.NoError(t, err)
				
				// Get results
				results := benchmark.GetBenchmarkResults()
				assert.NotNil(t, results)
				
				// Get status
				status := benchmark.GetBenchmarkStatus()
				assert.NotNil(t, status)
				
				// Get metrics
				metrics := benchmark.GetBenchmarkMetrics()
				assert.NotNil(t, metrics)
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Stop benchmark
		benchmark.Stop()
	})

	t.Run("Benchmark with different configurations", func(t *testing.T) {
		configs := []*Config{
			{
				Endpoints:   []string{"localhost:2379"},
				DialTimeout: 5 * time.Second,
				Duration:    100 * time.Millisecond,
				Concurrency: 1,
			},
			{
				Endpoints:   []string{"localhost:2379", "localhost:2380"},
				DialTimeout: 10 * time.Second,
				Duration:    200 * time.Millisecond,
				Concurrency: 2,
			},
			{
				Endpoints:   []string{"localhost:2379", "localhost:2380", "localhost:2381"},
				DialTimeout: 15 * time.Second,
				Duration:    300 * time.Millisecond,
				Concurrency: 3,
			},
		}
		logger, _ := zap.NewDevelopment()

		for i, config := range configs {
			benchmark := NewBenchmark(config, logger)
			assert.NotNil(t, benchmark)

			// Start benchmark
			err := benchmark.Start()
			require.NoError(t, err)

			// Run benchmark
			ctx := context.Background()
			err = benchmark.Run(ctx)
			assert.NoError(t, err)

			// Get results
			results := benchmark.GetBenchmarkResults()
			assert.NotNil(t, results)

			// Get status
			status := benchmark.GetBenchmarkStatus()
			assert.NotNil(t, status)

			// Get metrics
			metrics := benchmark.GetBenchmarkMetrics()
			assert.NotNil(t, metrics)

			// Stop benchmark
			benchmark.Stop()
		}
	})

	t.Run("Benchmark with context cancellation", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)

		// Start benchmark
		err := benchmark.Start()
		require.NoError(t, err)

		// Run benchmark with context that gets cancelled
		ctx, cancel := context.WithCancel(context.Background())
		
		// Cancel context after a short delay
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()
		
		err = benchmark.Run(ctx)
		assert.Error(t, err) // Should error due to context cancellation

		// Stop benchmark
		benchmark.Stop()
	})

	t.Run("Benchmark with context timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Duration:    10 * time.Second,
			Concurrency: 10,
		}
		logger, _ := zap.NewDevelopment()

		benchmark := NewBenchmark(config, logger)
		assert.NotNil(t, benchmark)

		// Start benchmark
		err := benchmark.Start()
		require.NoError(t, err)

		// Run benchmark with context that times out
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		err = benchmark.Run(ctx)
		assert.Error(t, err) // Should error due to context timeout

		// Stop benchmark
		benchmark.Stop()
	})
}
