package router

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRouter_Comprehensive(t *testing.T) {
	t.Run("NewRouter with valid config", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		router := NewRouter(config, logger)
		assert.NotNil(t, router)
		assert.NotNil(t, router.config)
		assert.NotNil(t, router.logger)
	})

	t.Run("NewRouter with nil config", func(t *testing.T) {
		logger, _ := zap.NewDevelopment()

		router := NewRouter(nil, logger)
		assert.NotNil(t, router)
		assert.Nil(t, router.config)
		assert.NotNil(t, router.logger)
	})

	t.Run("NewRouter with nil logger", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}

		router := NewRouter(config, nil)
		assert.NotNil(t, router)
		assert.NotNil(t, router.config)
		assert.NotNil(t, router.logger) // Should create a production logger
	})

	t.Run("Start router", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		router := NewRouter(config, logger)
		err := router.Start()
		assert.NoError(t, err)
	})

	t.Run("Stop router", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		router := NewRouter(config, logger)
		err := router.Start()
		require.NoError(t, err)

		router.Stop()
		// Should not panic or error
	})

	t.Run("Run router", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		router := NewRouter(config, logger)
		ctx := context.Background()
		
		err := router.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("GetRouterResults", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		router := NewRouter(config, logger)
		results := router.GetRouterResults()
		assert.NotNil(t, results)
	})

	t.Run("GetRouterStatus", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		router := NewRouter(config, logger)
		status := router.GetRouterStatus()
		assert.NotNil(t, status)
	})

	t.Run("GetRouterMetrics", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		router := NewRouter(config, logger)
		metrics := router.GetRouterMetrics()
		assert.NotNil(t, metrics)
	})
}

func TestRouter_EdgeCases(t *testing.T) {
	t.Run("Router with empty config", func(t *testing.T) {
		config := &Config{}
		logger, _ := zap.NewDevelopment()

		router := NewRouter(config, logger)
		assert.NotNil(t, router)
		assert.NotNil(t, router.config)
		assert.NotNil(t, router.logger)
	})

	t.Run("Router with nil endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		router := NewRouter(config, logger)
		assert.NotNil(t, router)
		assert.Nil(t, router.config.Endpoints)
	})

	t.Run("Router with empty endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		router := NewRouter(config, logger)
		assert.NotNil(t, router)
		assert.Len(t, router.config.Endpoints, 0)
	})

	t.Run("Router with zero timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 0,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		router := NewRouter(config, logger)
		assert.NotNil(t, router)
		assert.Equal(t, time.Duration(0), router.config.DialTimeout)
	})

	t.Run("Router with negative timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: -1 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		router := NewRouter(config, logger)
		assert.NotNil(t, router)
		assert.Equal(t, -1*time.Second, router.config.DialTimeout)
	})

	t.Run("Router with empty username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		router := NewRouter(config, logger)
		assert.NotNil(t, router)
		assert.Empty(t, router.config.Username)
	})

	t.Run("Router with empty password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		router := NewRouter(config, logger)
		assert.NotNil(t, router)
		assert.Empty(t, router.config.Password)
	})
}

func TestRouter_Performance(t *testing.T) {
	t.Run("Router creation performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Create 1000 Router instances
		for i := 0; i < 1000; i++ {
			router := NewRouter(config, logger)
			assert.NotNil(t, router)
		}
		
		duration := time.Since(start)
		
		// Should create 1000 instances in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Creating 1000 Router instances took too long: %v", duration)
	})

	t.Run("Router start/stop performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Start and stop 100 routers
		for i := 0; i < 100; i++ {
			router := NewRouter(config, logger)
			err := router.Start()
			require.NoError(t, err)
			router.Stop()
		}
		
		duration := time.Since(start)
		
		// Should start and stop 100 routers in less than 1 second
		assert.True(t, duration < 1*time.Second, "Starting and stopping 100 routers took too long: %v", duration)
	})
}

func TestRouter_Integration(t *testing.T) {
	t.Run("Router full lifecycle", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		router := NewRouter(config, logger)
		assert.NotNil(t, router)

		// Start router
		err := router.Start()
		assert.NoError(t, err)

		// Run router
		ctx := context.Background()
		err = router.Run(ctx)
		assert.NoError(t, err)

		// Get results
		results := router.GetRouterResults()
		assert.NotNil(t, results)

		// Get status
		status := router.GetRouterStatus()
		assert.NotNil(t, status)

		// Get metrics
		metrics := router.GetRouterMetrics()
		assert.NotNil(t, metrics)

		// Stop router
		router.Stop()
	})

	t.Run("Router with multiple concurrent operations", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		router := NewRouter(config, logger)
		assert.NotNil(t, router)

		// Start router
		err := router.Start()
		require.NoError(t, err)

		// Perform multiple concurrent operations
		done := make(chan bool, 10)
		
		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				// Run router
				ctx := context.Background()
				err := router.Run(ctx)
				assert.NoError(t, err)
				
				// Get results
				results := router.GetRouterResults()
				assert.NotNil(t, results)
				
				// Get status
				status := router.GetRouterStatus()
				assert.NotNil(t, status)
				
				// Get metrics
				metrics := router.GetRouterMetrics()
				assert.NotNil(t, metrics)
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Stop router
		router.Stop()
	})
}
