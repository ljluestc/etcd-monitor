package pingap

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestPingap_Comprehensive(t *testing.T) {
	t.Run("NewPingap with valid config", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingap(config, logger)
		assert.NotNil(t, pingap)
		assert.NotNil(t, pingap.config)
		assert.NotNil(t, pingap.logger)
	})

	t.Run("NewPingap with nil config", func(t *testing.T) {
		logger, _ := zap.NewDevelopment()

		pingap := NewPingap(nil, logger)
		assert.NotNil(t, pingap)
		assert.Nil(t, pingap.config)
		assert.NotNil(t, pingap.logger)
	})

	t.Run("NewPingap with nil logger", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}

		pingap := NewPingap(config, nil)
		assert.NotNil(t, pingap)
		assert.NotNil(t, pingap.config)
		assert.NotNil(t, pingap.logger) // Should create a production logger
	})

	t.Run("Start pingap", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingap(config, logger)
		err := pingap.Start()
		assert.NoError(t, err)
	})

	t.Run("Stop pingap", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingap(config, logger)
		err := pingap.Start()
		require.NoError(t, err)

		pingap.Stop()
		// Should not panic or error
	})

	t.Run("Run pingap", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingap(config, logger)
		ctx := context.Background()
		
		err := pingap.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("GetPingapResults", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingap(config, logger)
		results := pingap.GetPingapResults()
		assert.NotNil(t, results)
	})

	t.Run("GetPingapStatus", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingap(config, logger)
		status := pingap.GetPingapStatus()
		assert.NotNil(t, status)
	})

	t.Run("GetPingapMetrics", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingap(config, logger)
		metrics := pingap.GetPingapMetrics()
		assert.NotNil(t, metrics)
	})
}

func TestPingap_EdgeCases(t *testing.T) {
	t.Run("Pingap with empty config", func(t *testing.T) {
		config := &Config{}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingap(config, logger)
		assert.NotNil(t, pingap)
		assert.NotNil(t, pingap.config)
		assert.NotNil(t, pingap.logger)
	})

	t.Run("Pingap with nil endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingap(config, logger)
		assert.NotNil(t, pingap)
		assert.Nil(t, pingap.config.Endpoints)
	})

	t.Run("Pingap with empty endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingap(config, logger)
		assert.NotNil(t, pingap)
		assert.Len(t, pingap.config.Endpoints, 0)
	})

	t.Run("Pingap with zero timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 0,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingap(config, logger)
		assert.NotNil(t, pingap)
		assert.Equal(t, time.Duration(0), pingap.config.DialTimeout)
	})

	t.Run("Pingap with negative timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: -1 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingap(config, logger)
		assert.NotNil(t, pingap)
		assert.Equal(t, -1*time.Second, pingap.config.DialTimeout)
	})

	t.Run("Pingap with empty username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingap(config, logger)
		assert.NotNil(t, pingap)
		assert.Empty(t, pingap.config.Username)
	})

	t.Run("Pingap with empty password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingap(config, logger)
		assert.NotNil(t, pingap)
		assert.Empty(t, pingap.config.Password)
	})
}

func TestPingap_Performance(t *testing.T) {
	t.Run("Pingap creation performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Create 1000 Pingap instances
		for i := 0; i < 1000; i++ {
			pingap := NewPingap(config, logger)
			assert.NotNil(t, pingap)
		}
		
		duration := time.Since(start)
		
		// Should create 1000 instances in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Creating 1000 Pingap instances took too long: %v", duration)
	})

	t.Run("Pingap start/stop performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Start and stop 100 pingaps
		for i := 0; i < 100; i++ {
			pingap := NewPingap(config, logger)
			err := pingap.Start()
			require.NoError(t, err)
			pingap.Stop()
		}
		
		duration := time.Since(start)
		
		// Should start and stop 100 pingaps in less than 1 second
		assert.True(t, duration < 1*time.Second, "Starting and stopping 100 pingaps took too long: %v", duration)
	})
}

func TestPingap_Integration(t *testing.T) {
	t.Run("Pingap full lifecycle", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingap(config, logger)
		assert.NotNil(t, pingap)

		// Start pingap
		err := pingap.Start()
		assert.NoError(t, err)

		// Run pingap
		ctx := context.Background()
		err = pingap.Run(ctx)
		assert.NoError(t, err)

		// Get results
		results := pingap.GetPingapResults()
		assert.NotNil(t, results)

		// Get status
		status := pingap.GetPingapStatus()
		assert.NotNil(t, status)

		// Get metrics
		metrics := pingap.GetPingapMetrics()
		assert.NotNil(t, metrics)

		// Stop pingap
		pingap.Stop()
	})

	t.Run("Pingap with multiple concurrent operations", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingap(config, logger)
		assert.NotNil(t, pingap)

		// Start pingap
		err := pingap.Start()
		require.NoError(t, err)

		// Perform multiple concurrent operations
		done := make(chan bool, 10)
		
		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				// Run pingap
				ctx := context.Background()
				err := pingap.Run(ctx)
				assert.NoError(t, err)
				
				// Get results
				results := pingap.GetPingapResults()
				assert.NotNil(t, results)
				
				// Get status
				status := pingap.GetPingapStatus()
				assert.NotNil(t, status)
				
				// Get metrics
				metrics := pingap.GetPingapMetrics()
				assert.NotNil(t, metrics)
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Stop pingap
		pingap.Stop()
	})
}
