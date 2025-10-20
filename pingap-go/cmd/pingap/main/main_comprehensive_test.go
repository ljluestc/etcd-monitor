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
}

func TestMain_EdgeCases(t *testing.T) {
	t.Run("Main with empty config", func(t *testing.T) {
		config := &Config{}
		logger, _ := zap.NewDevelopment()

		main := NewMain(config, logger)
		assert.NotNil(t, main)
		assert.NotNil(t, main.config)
		assert.NotNil(t, main.logger)
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
}

func TestMain_Integration(t *testing.T) {
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
}
