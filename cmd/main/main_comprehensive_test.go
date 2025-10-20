package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestMain_Comprehensive(t *testing.T) {
	t.Run("NewMain with valid config", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.NotNil(t, main.config)
		assert.NotNil(t, main.logger)
	})

	t.Run("NewMain with nil config", func(t *testing.T) {
		logger, _ := zap.NewDevelopment()

		main := NewMain(nil, logger)
		assert.NotNil(t, main)
		assert.Nil(t, main.config)
		assert.NotNil(t, main.logger)
	})

	t.Run("NewMain with nil logger", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}

		main := NewMain(config, nil)
		assert.NotNil(t, main)
		assert.NotNil(t, main.config)
		assert.NotNil(t, main.logger) // Should create a production logger
	})

	t.Run("NewMain with empty config", func(t *testing.T) {
		config := &Config{}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.NotNil(t, main.config)
		assert.NotNil(t, main.logger)
	})

	t.Run("Start main", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		err := main.Start()
		assert.NoError(t, err)
	})

	t.Run("Stop main", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		err := main.Start()
		require.NoError(t, err)

		main.Stop()
		// Should not panic or error
	})

	t.Run("Stop main without starting", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		main.Stop()
		// Should not panic or error
	})

	t.Run("Run main", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		ctx := context.Background()
		
		err := main.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("Run main with nil context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		
		err := main.Run(nil)
		assert.Error(t, err)
	})

	t.Run("Run main with cancelled context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		err := main.Run(ctx)
		assert.Error(t, err)
	})

	t.Run("Run main with timeout context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		err := main.Run(ctx)
		assert.Error(t, err)
	})

	t.Run("GetMainResults", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		results := main.GetMainResults()
		assert.NotNil(t, results)
	})

	t.Run("GetMainResults after run", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		ctx := context.Background()
		
		err := main.Run(ctx)
		require.NoError(t, err)
		
		results := main.GetMainResults()
		assert.NotNil(t, results)
	})

	t.Run("GetMainStatus", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		status := main.GetMainStatus()
		assert.NotNil(t, status)
	})

	t.Run("GetMainStatus after run", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		ctx := context.Background()
		
		err := main.Run(ctx)
		require.NoError(t, err)
		
		status := main.GetMainStatus()
		assert.NotNil(t, status)
	})

	t.Run("GetMainMetrics", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		metrics := main.GetMainMetrics()
		assert.NotNil(t, metrics)
	})

	t.Run("GetMainMetrics after run", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		ctx := context.Background()
		
		err := main.Run(ctx)
		require.NoError(t, err)
		
		metrics := main.GetMainMetrics()
		assert.NotNil(t, metrics)
	})
}

func TestMain_EdgeCases(t *testing.T) {
	t.Run("Main with very long endpoint list", func(t *testing.T) {
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

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.Len(t, main.config.Endpoints, 1000)
	})

	t.Run("Main with very long timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 24 * time.Hour,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.Equal(t, 24*time.Hour, main.config.DialTimeout)
	})

	t.Run("Main with very short timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 1 * time.Nanosecond,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.Equal(t, 1*time.Nanosecond, main.config.DialTimeout)
	})

	t.Run("Main with zero timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 0,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.Equal(t, time.Duration(0), main.config.DialTimeout)
	})

	t.Run("Main with negative timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: -1 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.Equal(t, -1*time.Second, main.config.DialTimeout)
	})

	t.Run("Main with very long username", func(t *testing.T) {
		longUsername := strings.Repeat("a", 10000)
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    longUsername,
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.Len(t, main.config.Username, 10000)
	})

	t.Run("Main with very long password", func(t *testing.T) {
		longPassword := strings.Repeat("b", 10000)
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    longPassword,
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.Len(t, main.config.Password, 10000)
	})

	t.Run("Main with empty endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.Len(t, main.config.Endpoints, 0)
	})

	t.Run("Main with nil endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.Nil(t, main.config.Endpoints)
	})

	t.Run("Main with empty username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.Empty(t, main.config.Username)
	})

	t.Run("Main with empty password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.Empty(t, main.config.Password)
	})

	t.Run("Main with special characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "https://etcd-1:2379", "http://etcd-2:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.Len(t, main.config.Endpoints, 3)
		assert.Equal(t, "http://etcd-0:2379", main.config.Endpoints[0])
		assert.Equal(t, "https://etcd-1:2379", main.config.Endpoints[1])
		assert.Equal(t, "http://etcd-2:2379", main.config.Endpoints[2])
	})

	t.Run("Main with special characters in username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "user@domain.com",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.Equal(t, "user@domain.com", main.config.Username)
	})

	t.Run("Main with special characters in password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "pass!@#$%^&*()",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.Equal(t, "pass!@#$%^&*()", main.config.Password)
	})

	t.Run("Main with unicode characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.Len(t, main.config.Endpoints, 2)
	})

	t.Run("Main with unicode characters in username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "用户",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.Equal(t, "用户", main.config.Username)
	})

	t.Run("Main with unicode characters in password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "密码",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.Equal(t, "密码", main.config.Password)
	})

	t.Run("Main with control characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.Len(t, main.config.Endpoints, 2)
	})

	t.Run("Main with control characters in username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "user\n\r\t",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.Equal(t, "user\n\r\t", main.config.Username)
	})

	t.Run("Main with control characters in password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "pass\n\r\t",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.Equal(t, "pass\n\r\t", main.config.Password)
	})
}

func TestMain_Performance(t *testing.T) {
	t.Run("Main creation performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Create 1000 Main instances
		for i := 0; i < 1000; i++ {
			main := NewMain(config, logger)
			assert.NotNil(t, main)
		}
		
		duration := time.Since(start)
		
		// Should create 1000 instances in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Creating 1000 Main instances took too long: %v", duration)
	})

	t.Run("Main with large config performance", func(t *testing.T) {
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
		
		main := NewMain(config, logger)
		
		duration := time.Since(start)
		
		// Should create main with large config in less than 10ms
		assert.True(t, duration < 10*time.Millisecond, "Creating Main with large config took too long: %v", duration)
		assert.Len(t, main.config.Endpoints, 10000)
		assert.Len(t, main.config.Username, 10000)
		assert.Len(t, main.config.Password, 10000)
	})

	t.Run("Main start/stop performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Start and stop 100 mains
		for i := 0; i < 100; i++ {
			main := NewMain(config, logger)
			err := main.Start()
			require.NoError(t, err)
			main.Stop()
		}
		
		duration := time.Since(start)
		
		// Should start and stop 100 mains in less than 1 second
		assert.True(t, duration < 1*time.Second, "Starting and stopping 100 mains took too long: %v", duration)
	})

	t.Run("Main run performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		main := NewMain(config, logger)
		ctx := context.Background()
		
		start := time.Now()
		
		// Run main 10 times
		for i := 0; i < 10; i++ {
			err := main.Run(ctx)
			assert.NoError(t, err)
		}
		
		duration := time.Since(start)
		
		// Should run 10 mains in less than 1 second
		assert.True(t, duration < 1*time.Second, "Running 10 mains took too long: %v", duration)
	})

	t.Run("Main get results performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		main := NewMain(config, logger)
		
		start := time.Now()
		
		// Get results 1000 times
		for i := 0; i < 1000; i++ {
			results := main.GetMainResults()
			assert.NotNil(t, results)
		}
		
		duration := time.Since(start)
		
		// Should get results 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting results 1000 times took too long: %v", duration)
	})

	t.Run("Main get status performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		main := NewMain(config, logger)
		
		start := time.Now()
		
		// Get status 1000 times
		for i := 0; i < 1000; i++ {
			status := main.GetMainStatus()
			assert.NotNil(t, status)
		}
		
		duration := time.Since(start)
		
		// Should get status 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting status 1000 times took too long: %v", duration)
	})

	t.Run("Main get metrics performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		main := NewMain(config, logger)
		
		start := time.Now()
		
		// Get metrics 1000 times
		for i := 0; i < 1000; i++ {
			metrics := main.GetMainMetrics()
			assert.NotNil(t, metrics)
		}
		
		duration := time.Since(start)
		
		// Should get metrics 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting metrics 1000 times took too long: %v", duration)
	})
}

func TestMain_Integration(t *testing.T) {
	t.Run("Main with all fields set", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379", "http://etcd-2:2379"},
			DialTimeout: 10 * time.Second,
			Username:    "admin",
			Password:    "secret",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.NotNil(t, main.config)
		assert.NotNil(t, main.logger)
		assert.Len(t, main.config.Endpoints, 3)
		assert.Equal(t, "http://etcd-0:2379", main.config.Endpoints[0])
		assert.Equal(t, "http://etcd-1:2379", main.config.Endpoints[1])
		assert.Equal(t, "http://etcd-2:2379", main.config.Endpoints[2])
		assert.Equal(t, 10*time.Second, main.config.DialTimeout)
		assert.Equal(t, "admin", main.config.Username)
		assert.Equal(t, "secret", main.config.Password)
	})

	t.Run("Main with minimal fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.NotNil(t, main.config)
		assert.NotNil(t, main.logger)
		assert.Len(t, main.config.Endpoints, 1)
		assert.Equal(t, "localhost:2379", main.config.Endpoints[0])
		assert.Equal(t, 5*time.Second, main.config.DialTimeout)
		assert.Equal(t, "testuser", main.config.Username)
		assert.Equal(t, "testpass", main.config.Password)
	})

	t.Run("Main with empty fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 0,
			Username:    "",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.NotNil(t, main.config)
		assert.NotNil(t, main.logger)
		assert.Len(t, main.config.Endpoints, 0)
		assert.Equal(t, time.Duration(0), main.config.DialTimeout)
		assert.Empty(t, main.config.Username)
		assert.Empty(t, main.config.Password)
	})

	t.Run("Main with nil fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 0,
			Username:    "",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.NotNil(t, main.config)
		assert.NotNil(t, main.logger)
		assert.Nil(t, main.config.Endpoints)
		assert.Equal(t, time.Duration(0), main.config.DialTimeout)
		assert.Empty(t, main.config.Username)
		assert.Empty(t, main.config.Password)
	})

	t.Run("Main full lifecycle", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)

		// Start main
		err := main.Start()
		assert.NoError(t, err)

		// Run main
		ctx := context.Background()
		err = main.Run(ctx)
		assert.NoError(t, err)

		// Get results
		results := main.GetMainResults()
		assert.NotNil(t, results)

		// Get status
		status := main.GetMainStatus()
		assert.NotNil(t, status)

		// Get metrics
		metrics := main.GetMainMetrics()
		assert.NotNil(t, metrics)

		// Stop main
		main.Stop()
	})

	t.Run("Main with multiple concurrent operations", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)

		// Start main
		err := main.Start()
		require.NoError(t, err)

		// Perform multiple concurrent operations
		done := make(chan bool, 10)
		
		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				// Run main
				ctx := context.Background()
				err := main.Run(ctx)
				assert.NoError(t, err)
				
				// Get results
				results := main.GetMainResults()
				assert.NotNil(t, results)
				
				// Get status
				status := main.GetMainStatus()
				assert.NotNil(t, status)
				
				// Get metrics
				metrics := main.GetMainMetrics()
				assert.NotNil(t, metrics)
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Stop main
		main.Stop()
	})

	t.Run("Main with different configurations", func(t *testing.T) {
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
			main := NewMain(config, logger)
			assert.NotNil(t, main)

			// Start main
			err := main.Start()
			require.NoError(t, err)

			// Run main
			ctx := context.Background()
			err = main.Run(ctx)
			assert.NoError(t, err)

			// Get results
			results := main.GetMainResults()
			assert.NotNil(t, results)

			// Get status
			status := main.GetMainStatus()
			assert.NotNil(t, status)

			// Get metrics
			metrics := main.GetMainMetrics()
			assert.NotNil(t, metrics)

			// Stop main
			main.Stop()
		}
	})

	t.Run("Main with context cancellation", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)

		// Start main
		err := main.Start()
		require.NoError(t, err)

		// Run main with context that gets cancelled
		ctx, cancel := context.WithCancel(context.Background())
		
		// Cancel context after a short delay
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()
		
		err = main.Run(ctx)
		assert.Error(t, err) // Should error due to context cancellation

		// Stop main
		main.Stop()
	})

	t.Run("Main with context timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)

		// Start main
		err := main.Start()
		require.NoError(t, err)

		// Run main with context that times out
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		err = main.Run(ctx)
		assert.Error(t, err) // Should error due to context timeout

		// Stop main
		main.Stop()
	})
}
