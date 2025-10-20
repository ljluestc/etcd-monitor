package upstream

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestUpstream_Comprehensive(t *testing.T) {
	t.Run("NewUpstream with valid config", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		upstream := NewUpstream(config, logger)
		assert.NotNil(t, upstream)
		assert.NotNil(t, upstream.config)
		assert.NotNil(t, upstream.logger)
	})

	t.Run("NewUpstream with nil config", func(t *testing.T) {
		logger, _ := zap.NewDevelopment()

		upstream := NewUpstream(nil, logger)
		assert.NotNil(t, upstream)
		assert.Nil(t, upstream.config)
		assert.NotNil(t, upstream.logger)
	})

	t.Run("NewUpstream with nil logger", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}

		upstream := NewUpstream(config, nil)
		assert.NotNil(t, upstream)
		assert.NotNil(t, upstream.config)
		assert.NotNil(t, upstream.logger) // Should create a production logger
	})

	t.Run("Start upstream", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		upstream := NewUpstream(config, logger)
		err := upstream.Start()
		assert.NoError(t, err)
	})

	t.Run("Stop upstream", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		upstream := NewUpstream(config, logger)
		err := upstream.Start()
		require.NoError(t, err)

		upstream.Stop()
		// Should not panic or error
	})

	t.Run("Run upstream", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		upstream := NewUpstream(config, logger)
		ctx := context.Background()
		
		err := upstream.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("GetUpstreamResults", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		upstream := NewUpstream(config, logger)
		results := upstream.GetUpstreamResults()
		assert.NotNil(t, results)
	})

	t.Run("GetUpstreamStatus", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		upstream := NewUpstream(config, logger)
		status := upstream.GetUpstreamStatus()
		assert.NotNil(t, status)
	})

	t.Run("GetUpstreamMetrics", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		upstream := NewUpstream(config, logger)
		metrics := upstream.GetUpstreamMetrics()
		assert.NotNil(t, metrics)
	})
}

func TestUpstream_EdgeCases(t *testing.T) {
	t.Run("Upstream with empty config", func(t *testing.T) {
		config := &Config{}
		logger, _ := zap.NewDevelopment()

		upstream := NewUpstream(config, logger)
		assert.NotNil(t, upstream)
		assert.NotNil(t, upstream.config)
		assert.NotNil(t, upstream.logger)
	})

	t.Run("Upstream with nil endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		upstream := NewUpstream(config, logger)
		assert.NotNil(t, upstream)
		assert.Nil(t, upstream.config.Endpoints)
	})

	t.Run("Upstream with empty endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		upstream := NewUpstream(config, logger)
		assert.NotNil(t, upstream)
		assert.Len(t, upstream.config.Endpoints, 0)
	})

	t.Run("Upstream with zero timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 0,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		upstream := NewUpstream(config, logger)
		assert.NotNil(t, upstream)
		assert.Equal(t, time.Duration(0), upstream.config.DialTimeout)
	})

	t.Run("Upstream with negative timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: -1 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		upstream := NewUpstream(config, logger)
		assert.NotNil(t, upstream)
		assert.Equal(t, -1*time.Second, upstream.config.DialTimeout)
	})

	t.Run("Upstream with empty username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		upstream := NewUpstream(config, logger)
		assert.NotNil(t, upstream)
		assert.Empty(t, upstream.config.Username)
	})

	t.Run("Upstream with empty password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		upstream := NewUpstream(config, logger)
		assert.NotNil(t, upstream)
		assert.Empty(t, upstream.config.Password)
	})
}

func TestUpstream_Performance(t *testing.T) {
	t.Run("Upstream creation performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Create 1000 Upstream instances
		for i := 0; i < 1000; i++ {
			upstream := NewUpstream(config, logger)
			assert.NotNil(t, upstream)
		}
		
		duration := time.Since(start)
		
		// Should create 1000 instances in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Creating 1000 Upstream instances took too long: %v", duration)
	})

	t.Run("Upstream start/stop performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Start and stop 100 upstreams
		for i := 0; i < 100; i++ {
			upstream := NewUpstream(config, logger)
			err := upstream.Start()
			require.NoError(t, err)
			upstream.Stop()
		}
		
		duration := time.Since(start)
		
		// Should start and stop 100 upstreams in less than 1 second
		assert.True(t, duration < 1*time.Second, "Starting and stopping 100 upstreams took too long: %v", duration)
	})
}

func TestUpstream_Integration(t *testing.T) {
	t.Run("Upstream full lifecycle", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		upstream := NewUpstream(config, logger)
		assert.NotNil(t, upstream)

		// Start upstream
		err := upstream.Start()
		assert.NoError(t, err)

		// Run upstream
		ctx := context.Background()
		err = upstream.Run(ctx)
		assert.NoError(t, err)

		// Get results
		results := upstream.GetUpstreamResults()
		assert.NotNil(t, results)

		// Get status
		status := upstream.GetUpstreamStatus()
		assert.NotNil(t, status)

		// Get metrics
		metrics := upstream.GetUpstreamMetrics()
		assert.NotNil(t, metrics)

		// Stop upstream
		upstream.Stop()
	})

	t.Run("Upstream with multiple concurrent operations", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		upstream := NewUpstream(config, logger)
		assert.NotNil(t, upstream)

		// Start upstream
		err := upstream.Start()
		require.NoError(t, err)

		// Perform multiple concurrent operations
		done := make(chan bool, 10)
		
		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				// Run upstream
				ctx := context.Background()
				err := upstream.Run(ctx)
				assert.NoError(t, err)
				
				// Get results
				results := upstream.GetUpstreamResults()
				assert.NotNil(t, results)
				
				// Get status
				status := upstream.GetUpstreamStatus()
				assert.NotNil(t, status)
				
				// Get metrics
				metrics := upstream.GetUpstreamMetrics()
				assert.NotNil(t, metrics)
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Stop upstream
		upstream.Stop()
	})
}
