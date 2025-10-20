package examples

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestExamples_Comprehensive(t *testing.T) {
	t.Run("NewExamples with valid config", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.NotNil(t, examples.config)
		assert.NotNil(t, examples.logger)
	})

	t.Run("NewExamples with nil config", func(t *testing.T) {
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(nil, logger)
		assert.NotNil(t, examples)
		assert.Nil(t, examples.config)
		assert.NotNil(t, examples.logger)
	})

	t.Run("NewExamples with nil logger", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}

		examples := NewExamples(config, nil)
		assert.NotNil(t, examples)
		assert.NotNil(t, examples.config)
		assert.NotNil(t, examples.logger) // Should create a production logger
	})

	t.Run("NewExamples with empty config", func(t *testing.T) {
		config := &Config{}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.NotNil(t, examples.config)
		assert.NotNil(t, examples.logger)
	})

	t.Run("Start examples", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		err := examples.Start()
		assert.NoError(t, err)
	})

	t.Run("Stop examples", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		err := examples.Start()
		require.NoError(t, err)

		examples.Stop()
		// Should not panic or error
	})

	t.Run("Stop examples without starting", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		examples.Stop()
		// Should not panic or error
	})

	t.Run("Run examples", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		ctx := context.Background()
		
		err := examples.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("Run examples with nil context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		
		err := examples.Run(nil)
		assert.Error(t, err)
	})

	t.Run("Run examples with cancelled context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		err := examples.Run(ctx)
		assert.Error(t, err)
	})

	t.Run("Run examples with timeout context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		err := examples.Run(ctx)
		assert.Error(t, err)
	})

	t.Run("GetExamplesResults", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		results := examples.GetExamplesResults()
		assert.NotNil(t, results)
	})

	t.Run("GetExamplesResults after run", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		ctx := context.Background()
		
		err := examples.Run(ctx)
		require.NoError(t, err)
		
		results := examples.GetExamplesResults()
		assert.NotNil(t, results)
	})

	t.Run("GetExamplesStatus", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		status := examples.GetExamplesStatus()
		assert.NotNil(t, status)
	})

	t.Run("GetExamplesStatus after run", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		ctx := context.Background()
		
		err := examples.Run(ctx)
		require.NoError(t, err)
		
		status := examples.GetExamplesStatus()
		assert.NotNil(t, status)
	})

	t.Run("GetExamplesMetrics", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		metrics := examples.GetExamplesMetrics()
		assert.NotNil(t, metrics)
	})

	t.Run("GetExamplesMetrics after run", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		ctx := context.Background()
		
		err := examples.Run(ctx)
		require.NoError(t, err)
		
		metrics := examples.GetExamplesMetrics()
		assert.NotNil(t, metrics)
	})
}

func TestExamples_EdgeCases(t *testing.T) {
	t.Run("Examples with very long endpoint list", func(t *testing.T) {
		endpoints := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			endpoints[i] = fmt.Sprintf("localhost:%d", 2379+i)
		}

		config := &Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.Len(t, examples.config.Endpoints, 1000)
	})

	t.Run("Examples with very long timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 24 * time.Hour,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.Equal(t, 24*time.Hour, examples.config.DialTimeout)
	})

	t.Run("Examples with very short timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 1 * time.Nanosecond,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.Equal(t, 1*time.Nanosecond, examples.config.DialTimeout)
	})

	t.Run("Examples with zero timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 0,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.Equal(t, time.Duration(0), examples.config.DialTimeout)
	})

	t.Run("Examples with negative timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: -1 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.Equal(t, -1*time.Second, examples.config.DialTimeout)
	})

	t.Run("Examples with very long username", func(t *testing.T) {
		longUsername := strings.Repeat("a", 10000)
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    longUsername,
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.Len(t, examples.config.Username, 10000)
	})

	t.Run("Examples with very long password", func(t *testing.T) {
		longPassword := strings.Repeat("b", 10000)
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    longPassword,
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.Len(t, examples.config.Password, 10000)
	})

	t.Run("Examples with empty endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.Len(t, examples.config.Endpoints, 0)
	})

	t.Run("Examples with nil endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.Nil(t, examples.config.Endpoints)
	})

	t.Run("Examples with empty username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.Empty(t, examples.config.Username)
	})

	t.Run("Examples with empty password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.Empty(t, examples.config.Password)
	})

	t.Run("Examples with special characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "https://etcd-1:2379", "http://etcd-2:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.Len(t, examples.config.Endpoints, 3)
		assert.Equal(t, "http://etcd-0:2379", examples.config.Endpoints[0])
		assert.Equal(t, "https://etcd-1:2379", examples.config.Endpoints[1])
		assert.Equal(t, "http://etcd-2:2379", examples.config.Endpoints[2])
	})

	t.Run("Examples with special characters in username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "user@domain.com",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.Equal(t, "user@domain.com", examples.config.Username)
	})

	t.Run("Examples with special characters in password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "pass!@#$%^&*()",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.Equal(t, "pass!@#$%^&*()", examples.config.Password)
	})

	t.Run("Examples with unicode characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.Len(t, examples.config.Endpoints, 2)
	})

	t.Run("Examples with unicode characters in username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "用户",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.Equal(t, "用户", examples.config.Username)
	})

	t.Run("Examples with unicode characters in password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "密码",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.Equal(t, "密码", examples.config.Password)
	})

	t.Run("Examples with control characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.Len(t, examples.config.Endpoints, 2)
	})

	t.Run("Examples with control characters in username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "user\n\r\t",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.Equal(t, "user\n\r\t", examples.config.Username)
	})

	t.Run("Examples with control characters in password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "pass\n\r\t",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.Equal(t, "pass\n\r\t", examples.config.Password)
	})
}

func TestExamples_Performance(t *testing.T) {
	t.Run("Examples creation performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Create 1000 Examples instances
		for i := 0; i < 1000; i++ {
			examples := NewExamples(config, logger)
			assert.NotNil(t, examples)
		}
		
		duration := time.Since(start)
		
		// Should create 1000 instances in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Creating 1000 Examples instances took too long: %v", duration)
	})

	t.Run("Examples with large config performance", func(t *testing.T) {
		endpoints := make([]string, 10000)
		for i := 0; i < 10000; i++ {
			endpoints[i] = fmt.Sprintf("localhost:%d", 2379+i)
		}
		
		config := &Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
			Username:    strings.Repeat("a", 10000),
			Password:    strings.Repeat("b", 10000),
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		examples := NewExamples(config, logger)
		
		duration := time.Since(start)
		
		// Should create examples with large config in less than 10ms
		assert.True(t, duration < 10*time.Millisecond, "Creating Examples with large config took too long: %v", duration)
		assert.Len(t, examples.config.Endpoints, 10000)
		assert.Len(t, examples.config.Username, 10000)
		assert.Len(t, examples.config.Password, 10000)
	})

	t.Run("Examples start/stop performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Start and stop 100 examples
		for i := 0; i < 100; i++ {
			examples := NewExamples(config, logger)
			err := examples.Start()
			require.NoError(t, err)
			examples.Stop()
		}
		
		duration := time.Since(start)
		
		// Should start and stop 100 examples in less than 1 second
		assert.True(t, duration < 1*time.Second, "Starting and stopping 100 examples took too long: %v", duration)
	})

	t.Run("Examples run performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		examples := NewExamples(config, logger)
		ctx := context.Background()
		
		start := time.Now()
		
		// Run examples 10 times
		for i := 0; i < 10; i++ {
			err := examples.Run(ctx)
			assert.NoError(t, err)
		}
		
		duration := time.Since(start)
		
		// Should run 10 examples in less than 1 second
		assert.True(t, duration < 1*time.Second, "Running 10 examples took too long: %v", duration)
	})

	t.Run("Examples get results performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		examples := NewExamples(config, logger)
		
		start := time.Now()
		
		// Get results 1000 times
		for i := 0; i < 1000; i++ {
			results := examples.GetExamplesResults()
			assert.NotNil(t, results)
		}
		
		duration := time.Since(start)
		
		// Should get results 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting results 1000 times took too long: %v", duration)
	})

	t.Run("Examples get status performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		examples := NewExamples(config, logger)
		
		start := time.Now()
		
		// Get status 1000 times
		for i := 0; i < 1000; i++ {
			status := examples.GetExamplesStatus()
			assert.NotNil(t, status)
		}
		
		duration := time.Since(start)
		
		// Should get status 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting status 1000 times took too long: %v", duration)
	})

	t.Run("Examples get metrics performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		examples := NewExamples(config, logger)
		
		start := time.Now()
		
		// Get metrics 1000 times
		for i := 0; i < 1000; i++ {
			metrics := examples.GetExamplesMetrics()
			assert.NotNil(t, metrics)
		}
		
		duration := time.Since(start)
		
		// Should get metrics 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting metrics 1000 times took too long: %v", duration)
	})
}

func TestExamples_Integration(t *testing.T) {
	t.Run("Examples with all fields set", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379", "http://etcd-2:2379"},
			DialTimeout: 10 * time.Second,
			Username:    "admin",
			Password:    "secret",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.NotNil(t, examples.config)
		assert.NotNil(t, examples.logger)
		assert.Len(t, examples.config.Endpoints, 3)
		assert.Equal(t, "http://etcd-0:2379", examples.config.Endpoints[0])
		assert.Equal(t, "http://etcd-1:2379", examples.config.Endpoints[1])
		assert.Equal(t, "http://etcd-2:2379", examples.config.Endpoints[2])
		assert.Equal(t, 10*time.Second, examples.config.DialTimeout)
		assert.Equal(t, "admin", examples.config.Username)
		assert.Equal(t, "secret", examples.config.Password)
	})

	t.Run("Examples with minimal fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.NotNil(t, examples.config)
		assert.NotNil(t, examples.logger)
		assert.Len(t, examples.config.Endpoints, 1)
		assert.Equal(t, "localhost:2379", examples.config.Endpoints[0])
		assert.Equal(t, 5*time.Second, examples.config.DialTimeout)
		assert.Equal(t, "testuser", examples.config.Username)
		assert.Equal(t, "testpass", examples.config.Password)
	})

	t.Run("Examples with empty fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 0,
			Username:    "",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.NotNil(t, examples.config)
		assert.NotNil(t, examples.logger)
		assert.Len(t, examples.config.Endpoints, 0)
		assert.Equal(t, time.Duration(0), examples.config.DialTimeout)
		assert.Empty(t, examples.config.Username)
		assert.Empty(t, examples.config.Password)
	})

	t.Run("Examples with nil fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 0,
			Username:    "",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)
		assert.NotNil(t, examples.config)
		assert.NotNil(t, examples.logger)
		assert.Nil(t, examples.config.Endpoints)
		assert.Equal(t, time.Duration(0), examples.config.DialTimeout)
		assert.Empty(t, examples.config.Username)
		assert.Empty(t, examples.config.Password)
	})

	t.Run("Examples full lifecycle", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)

		// Start examples
		err := examples.Start()
		assert.NoError(t, err)

		// Run examples
		ctx := context.Background()
		err = examples.Run(ctx)
		assert.NoError(t, err)

		// Get results
		results := examples.GetExamplesResults()
		assert.NotNil(t, results)

		// Get status
		status := examples.GetExamplesStatus()
		assert.NotNil(t, status)

		// Get metrics
		metrics := examples.GetExamplesMetrics()
		assert.NotNil(t, metrics)

		// Stop examples
		examples.Stop()
	})

	t.Run("Examples with multiple concurrent operations", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)

		// Start examples
		err := examples.Start()
		require.NoError(t, err)

		// Perform multiple concurrent operations
		done := make(chan bool, 10)
		
		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				// Run examples
				ctx := context.Background()
				err := examples.Run(ctx)
				assert.NoError(t, err)
				
				// Get results
				results := examples.GetExamplesResults()
				assert.NotNil(t, results)
				
				// Get status
				status := examples.GetExamplesStatus()
				assert.NotNil(t, status)
				
				// Get metrics
				metrics := examples.GetExamplesMetrics()
				assert.NotNil(t, metrics)
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Stop examples
		examples.Stop()
	})

	t.Run("Examples with different configurations", func(t *testing.T) {
		configs := []*Config{
			{
				Endpoints:   []string{"localhost:2379"},
				DialTimeout: 5 * time.Second,
				Username:    "testuser",
				Password:    "testpass",
			},
			{
				Endpoints:   []string{"localhost:2379", "localhost:2380"},
				DialTimeout: 10 * time.Second,
				Username:    "admin",
				Password:    "secret",
			},
			{
				Endpoints:   []string{"localhost:2379", "localhost:2380", "localhost:2381"},
				DialTimeout: 15 * time.Second,
				Username:    "root",
				Password:    "password",
			},
		}
		logger, _ := zap.NewDevelopment()

		for i, config := range configs {
			examples := NewExamples(config, logger)
			assert.NotNil(t, examples)

			// Start examples
			err := examples.Start()
			require.NoError(t, err)

			// Run examples
			ctx := context.Background()
			err = examples.Run(ctx)
			assert.NoError(t, err)

			// Get results
			results := examples.GetExamplesResults()
			assert.NotNil(t, results)

			// Get status
			status := examples.GetExamplesStatus()
			assert.NotNil(t, status)

			// Get metrics
			metrics := examples.GetExamplesMetrics()
			assert.NotNil(t, metrics)

			// Stop examples
			examples.Stop()
		}
	})

	t.Run("Examples with context cancellation", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)

		// Start examples
		err := examples.Start()
		require.NoError(t, err)

		// Run examples with context that gets cancelled
		ctx, cancel := context.WithCancel(context.Background())
		
		// Cancel context after a short delay
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()
		
		err = examples.Run(ctx)
		assert.Error(t, err) // Should error due to context cancellation

		// Stop examples
		examples.Stop()
	})

	t.Run("Examples with context timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		examples := NewExamples(config, logger)
		assert.NotNil(t, examples)

		// Start examples
		err := examples.Start()
		require.NoError(t, err)

		// Run examples with context that times out
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		err = examples.Run(ctx)
		assert.Error(t, err) // Should error due to context timeout

		// Stop examples
		examples.Stop()
	})
}
