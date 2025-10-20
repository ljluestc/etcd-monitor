package plugin

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestPlugin_Comprehensive(t *testing.T) {
	t.Run("NewPlugin with valid config", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		plugin := NewPlugin(config, logger)
		assert.NotNil(t, plugin)
		assert.NotNil(t, plugin.config)
		assert.NotNil(t, plugin.logger)
	})

	t.Run("NewPlugin with nil config", func(t *testing.T) {
		logger, _ := zap.NewDevelopment()

		plugin := NewPlugin(nil, logger)
		assert.NotNil(t, plugin)
		assert.Nil(t, plugin.config)
		assert.NotNil(t, plugin.logger)
	})

	t.Run("NewPlugin with nil logger", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}

		plugin := NewPlugin(config, nil)
		assert.NotNil(t, plugin)
		assert.NotNil(t, plugin.config)
		assert.NotNil(t, plugin.logger) // Should create a production logger
	})

	t.Run("Start plugin", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		plugin := NewPlugin(config, logger)
		err := plugin.Start()
		assert.NoError(t, err)
	})

	t.Run("Stop plugin", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		plugin := NewPlugin(config, logger)
		err := plugin.Start()
		require.NoError(t, err)

		plugin.Stop()
		// Should not panic or error
	})

	t.Run("Run plugin", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		plugin := NewPlugin(config, logger)
		ctx := context.Background()
		
		err := plugin.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("GetPluginResults", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		plugin := NewPlugin(config, logger)
		results := plugin.GetPluginResults()
		assert.NotNil(t, results)
	})

	t.Run("GetPluginStatus", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		plugin := NewPlugin(config, logger)
		status := plugin.GetPluginStatus()
		assert.NotNil(t, status)
	})

	t.Run("GetPluginMetrics", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		plugin := NewPlugin(config, logger)
		metrics := plugin.GetPluginMetrics()
		assert.NotNil(t, metrics)
	})
}

func TestPlugin_EdgeCases(t *testing.T) {
	t.Run("Plugin with empty config", func(t *testing.T) {
		config := &Config{}
		logger, _ := zap.NewDevelopment()

		plugin := NewPlugin(config, logger)
		assert.NotNil(t, plugin)
		assert.NotNil(t, plugin.config)
		assert.NotNil(t, plugin.logger)
	})

	t.Run("Plugin with nil endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		plugin := NewPlugin(config, logger)
		assert.NotNil(t, plugin)
		assert.Nil(t, plugin.config.Endpoints)
	})

	t.Run("Plugin with empty endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		plugin := NewPlugin(config, logger)
		assert.NotNil(t, plugin)
		assert.Len(t, plugin.config.Endpoints, 0)
	})

	t.Run("Plugin with zero timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 0,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		plugin := NewPlugin(config, logger)
		assert.NotNil(t, plugin)
		assert.Equal(t, time.Duration(0), plugin.config.DialTimeout)
	})

	t.Run("Plugin with negative timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: -1 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		plugin := NewPlugin(config, logger)
		assert.NotNil(t, plugin)
		assert.Equal(t, -1*time.Second, plugin.config.DialTimeout)
	})

	t.Run("Plugin with empty username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		plugin := NewPlugin(config, logger)
		assert.NotNil(t, plugin)
		assert.Empty(t, plugin.config.Username)
	})

	t.Run("Plugin with empty password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		plugin := NewPlugin(config, logger)
		assert.NotNil(t, plugin)
		assert.Empty(t, plugin.config.Password)
	})
}

func TestPlugin_Performance(t *testing.T) {
	t.Run("Plugin creation performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Create 1000 Plugin instances
		for i := 0; i < 1000; i++ {
			plugin := NewPlugin(config, logger)
			assert.NotNil(t, plugin)
		}
		
		duration := time.Since(start)
		
		// Should create 1000 instances in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Creating 1000 Plugin instances took too long: %v", duration)
	})

	t.Run("Plugin start/stop performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Start and stop 100 plugins
		for i := 0; i < 100; i++ {
			plugin := NewPlugin(config, logger)
			err := plugin.Start()
			require.NoError(t, err)
			plugin.Stop()
		}
		
		duration := time.Since(start)
		
		// Should start and stop 100 plugins in less than 1 second
		assert.True(t, duration < 1*time.Second, "Starting and stopping 100 plugins took too long: %v", duration)
	})
}

func TestPlugin_Integration(t *testing.T) {
	t.Run("Plugin full lifecycle", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		plugin := NewPlugin(config, logger)
		assert.NotNil(t, plugin)

		// Start plugin
		err := plugin.Start()
		assert.NoError(t, err)

		// Run plugin
		ctx := context.Background()
		err = plugin.Run(ctx)
		assert.NoError(t, err)

		// Get results
		results := plugin.GetPluginResults()
		assert.NotNil(t, results)

		// Get status
		status := plugin.GetPluginStatus()
		assert.NotNil(t, status)

		// Get metrics
		metrics := plugin.GetPluginMetrics()
		assert.NotNil(t, metrics)

		// Stop plugin
		plugin.Stop()
	})

	t.Run("Plugin with multiple concurrent operations", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		plugin := NewPlugin(config, logger)
		assert.NotNil(t, plugin)

		// Start plugin
		err := plugin.Start()
		require.NoError(t, err)

		// Perform multiple concurrent operations
		done := make(chan bool, 10)
		
		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				// Run plugin
				ctx := context.Background()
				err := plugin.Run(ctx)
				assert.NoError(t, err)
				
				// Get results
				results := plugin.GetPluginResults()
				assert.NotNil(t, results)
				
				// Get status
				status := plugin.GetPluginStatus()
				assert.NotNil(t, status)
				
				// Get metrics
				metrics := plugin.GetPluginMetrics()
				assert.NotNil(t, metrics)
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Stop plugin
		plugin.Stop()
	})
}
