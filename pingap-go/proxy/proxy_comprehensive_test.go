package proxy

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestProxy_Comprehensive(t *testing.T) {
	t.Run("NewProxy with valid config", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		proxy := NewProxy(config, logger)
		assert.NotNil(t, proxy)
		assert.NotNil(t, proxy.config)
		assert.NotNil(t, proxy.logger)
	})

	t.Run("NewProxy with nil config", func(t *testing.T) {
		logger, _ := zap.NewDevelopment()

		proxy := NewProxy(nil, logger)
		assert.NotNil(t, proxy)
		assert.Nil(t, proxy.config)
		assert.NotNil(t, proxy.logger)
	})

	t.Run("NewProxy with nil logger", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}

		proxy := NewProxy(config, nil)
		assert.NotNil(t, proxy)
		assert.NotNil(t, proxy.config)
		assert.NotNil(t, proxy.logger) // Should create a production logger
	})

	t.Run("Start proxy", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		proxy := NewProxy(config, logger)
		err := proxy.Start()
		assert.NoError(t, err)
	})

	t.Run("Stop proxy", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		proxy := NewProxy(config, logger)
		err := proxy.Start()
		require.NoError(t, err)

		proxy.Stop()
		// Should not panic or error
	})

	t.Run("Run proxy", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		proxy := NewProxy(config, logger)
		ctx := context.Background()
		
		err := proxy.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("GetProxyResults", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		proxy := NewProxy(config, logger)
		results := proxy.GetProxyResults()
		assert.NotNil(t, results)
	})

	t.Run("GetProxyStatus", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		proxy := NewProxy(config, logger)
		status := proxy.GetProxyStatus()
		assert.NotNil(t, status)
	})

	t.Run("GetProxyMetrics", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		proxy := NewProxy(config, logger)
		metrics := proxy.GetProxyMetrics()
		assert.NotNil(t, metrics)
	})
}

func TestProxy_EdgeCases(t *testing.T) {
	t.Run("Proxy with empty config", func(t *testing.T) {
		config := &Config{}
		logger, _ := zap.NewDevelopment()

		proxy := NewProxy(config, logger)
		assert.NotNil(t, proxy)
		assert.NotNil(t, proxy.config)
		assert.NotNil(t, proxy.logger)
	})

	t.Run("Proxy with nil endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		proxy := NewProxy(config, logger)
		assert.NotNil(t, proxy)
		assert.Nil(t, proxy.config.Endpoints)
	})

	t.Run("Proxy with empty endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		proxy := NewProxy(config, logger)
		assert.NotNil(t, proxy)
		assert.Len(t, proxy.config.Endpoints, 0)
	})

	t.Run("Proxy with zero timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 0,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		proxy := NewProxy(config, logger)
		assert.NotNil(t, proxy)
		assert.Equal(t, time.Duration(0), proxy.config.DialTimeout)
	})

	t.Run("Proxy with negative timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: -1 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		proxy := NewProxy(config, logger)
		assert.NotNil(t, proxy)
		assert.Equal(t, -1*time.Second, proxy.config.DialTimeout)
	})

	t.Run("Proxy with empty username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		proxy := NewProxy(config, logger)
		assert.NotNil(t, proxy)
		assert.Empty(t, proxy.config.Username)
	})

	t.Run("Proxy with empty password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		proxy := NewProxy(config, logger)
		assert.NotNil(t, proxy)
		assert.Empty(t, proxy.config.Password)
	})
}

func TestProxy_Performance(t *testing.T) {
	t.Run("Proxy creation performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Create 1000 Proxy instances
		for i := 0; i < 1000; i++ {
			proxy := NewProxy(config, logger)
			assert.NotNil(t, proxy)
		}
		
		duration := time.Since(start)
		
		// Should create 1000 instances in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Creating 1000 Proxy instances took too long: %v", duration)
	})

	t.Run("Proxy start/stop performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Start and stop 100 proxies
		for i := 0; i < 100; i++ {
			proxy := NewProxy(config, logger)
			err := proxy.Start()
			require.NoError(t, err)
			proxy.Stop()
		}
		
		duration := time.Since(start)
		
		// Should start and stop 100 proxies in less than 1 second
		assert.True(t, duration < 1*time.Second, "Starting and stopping 100 proxies took too long: %v", duration)
	})
}

func TestProxy_Integration(t *testing.T) {
	t.Run("Proxy full lifecycle", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		proxy := NewProxy(config, logger)
		assert.NotNil(t, proxy)

		// Start proxy
		err := proxy.Start()
		assert.NoError(t, err)

		// Run proxy
		ctx := context.Background()
		err = proxy.Run(ctx)
		assert.NoError(t, err)

		// Get results
		results := proxy.GetProxyResults()
		assert.NotNil(t, results)

		// Get status
		status := proxy.GetProxyStatus()
		assert.NotNil(t, status)

		// Get metrics
		metrics := proxy.GetProxyMetrics()
		assert.NotNil(t, metrics)

		// Stop proxy
		proxy.Stop()
	})

	t.Run("Proxy with multiple concurrent operations", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		proxy := NewProxy(config, logger)
		assert.NotNil(t, proxy)

		// Start proxy
		err := proxy.Start()
		require.NoError(t, err)

		// Perform multiple concurrent operations
		done := make(chan bool, 10)
		
		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				// Run proxy
				ctx := context.Background()
				err := proxy.Run(ctx)
				assert.NoError(t, err)
				
				// Get results
				results := proxy.GetProxyResults()
				assert.NotNil(t, results)
				
				// Get status
				status := proxy.GetProxyStatus()
				assert.NotNil(t, status)
				
				// Get metrics
				metrics := proxy.GetProxyMetrics()
				assert.NotNil(t, metrics)
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Stop proxy
		proxy.Stop()
	})
}
