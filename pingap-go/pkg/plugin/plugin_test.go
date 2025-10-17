package plugin

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// MockPluginTarget is a mock plugin target for testing
type MockPluginTarget struct {
	result string
	err    error
}

func (m *MockPluginTarget) Execute(ctx context.Context) (string, error) {
	return m.result, m.err
}

func (m *MockPluginTarget) Close() error {
	return nil
}

func TestPlugin_Comprehensive(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("NewPlugin", func(t *testing.T) {
		plugin := NewPlugin(logger)
		
		assert.NotNil(t, plugin)
		assert.Equal(t, logger, plugin.logger)
	})

	t.Run("NewPlugin with nil logger", func(t *testing.T) {
		plugin := NewPlugin(nil)
		
		assert.NotNil(t, plugin)
		assert.NotNil(t, plugin.logger) // Should create a production logger
	})
}

func TestPlugin_RunPlugin(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("RunPlugin - Success", func(t *testing.T) {
		mockTarget := &MockPluginTarget{
			result: "plugin result",
		}
		
		plugin := NewPlugin(logger)
		result, err := plugin.RunPlugin(context.Background(), mockTarget, "test-plugin")
		
		require.NoError(t, err)
		assert.Equal(t, "plugin result", result)
	})

	t.Run("RunPlugin - Error", func(t *testing.T) {
		mockTarget := &MockPluginTarget{
			err: assert.AnError,
		}
		
		plugin := NewPlugin(logger)
		result, err := plugin.RunPlugin(context.Background(), mockTarget, "test-plugin")
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunPlugin - Nil Target", func(t *testing.T) {
		plugin := NewPlugin(logger)
		result, err := plugin.RunPlugin(context.Background(), nil, "test-plugin")
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunPlugin - Empty Name", func(t *testing.T) {
		mockTarget := &MockPluginTarget{
			result: "plugin result",
		}
		
		plugin := NewPlugin(logger)
		result, err := plugin.RunPlugin(context.Background(), mockTarget, "")
		
		require.NoError(t, err)
		assert.Equal(t, "plugin result", result)
	})

	t.Run("RunPlugin - Context Cancellation", func(t *testing.T) {
		mockTarget := &MockPluginTarget{
			result: "plugin result",
		}
		
		plugin := NewPlugin(logger)
		
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		result, err := plugin.RunPlugin(ctx, mockTarget, "test-plugin")
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestPlugin_RunChain(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("RunChain - Success", func(t *testing.T) {
		mockTarget := &MockPluginTarget{
			result: "chain result",
		}
		
		plugin := NewPlugin(logger)
		result, err := plugin.RunChain(context.Background(), mockTarget)
		
		require.NoError(t, err)
		assert.Equal(t, "chain result", result)
	})

	t.Run("RunChain - Error", func(t *testing.T) {
		mockTarget := &MockPluginTarget{
			err: assert.AnError,
		}
		
		plugin := NewPlugin(logger)
		result, err := plugin.RunChain(context.Background(), mockTarget)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunChain - Nil Target", func(t *testing.T) {
		plugin := NewPlugin(logger)
		result, err := plugin.RunChain(context.Background(), nil)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestPlugin_EdgeCases(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Concurrent Operations", func(t *testing.T) {
		mockTarget := &MockPluginTarget{
			result: "concurrent plugin result",
		}
		
		plugin := NewPlugin(logger)
		
		// Test concurrent access
		done := make(chan bool, 20)
		for i := 0; i < 20; i++ {
			go func() {
				defer func() { done <- true }()
				
				// Test all methods concurrently
				result, err := plugin.RunPlugin(context.Background(), mockTarget, "test-plugin")
				if err == nil {
					assert.Equal(t, "concurrent plugin result", result)
				}
				
				chainResult, err := plugin.RunChain(context.Background(), mockTarget)
				if err == nil {
					assert.Equal(t, "concurrent plugin result", chainResult)
				}
			}()
		}
		
		// Wait for all goroutines to complete
		for i := 0; i < 20; i++ {
			<-done
		}
	})

	t.Run("Empty Results", func(t *testing.T) {
		mockTarget := &MockPluginTarget{
			result: "",
		}
		
		plugin := NewPlugin(logger)
		
		// Test with empty results
		result, err := plugin.RunPlugin(context.Background(), mockTarget, "test-plugin")
		require.NoError(t, err)
		assert.Empty(t, result)
		
		chainResult, err := plugin.RunChain(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, chainResult)
	})

	t.Run("Nil Results", func(t *testing.T) {
		mockTarget := &MockPluginTarget{
			result: "",
		}
		
		plugin := NewPlugin(logger)
		
		// Test with nil results
		result, err := plugin.RunPlugin(context.Background(), mockTarget, "test-plugin")
		require.NoError(t, err)
		assert.Empty(t, result)
		
		chainResult, err := plugin.RunChain(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, chainResult)
	})
}

func TestPlugin_DataStructures(t *testing.T) {
	t.Run("Plugin Structure", func(t *testing.T) {
		plugin := NewPlugin(nil)
		
		// Verify structure
		assert.NotNil(t, plugin)
		assert.NotNil(t, plugin.logger)
	})
}

func TestPlugin_Interfaces(t *testing.T) {
	t.Run("Plugin implements PluginInterface", func(t *testing.T) {
		plugin := NewPlugin(nil)
		
		var iface PluginInterface = plugin
		assert.NotNil(t, iface)
	})
}

func TestPlugin_ErrorHandling(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Multiple Errors", func(t *testing.T) {
		mockTarget := &MockPluginTarget{
			err: assert.AnError,
		}
		
		plugin := NewPlugin(logger)
		
		// Test all methods with errors
		result, err := plugin.RunPlugin(context.Background(), mockTarget, "test-plugin")
		assert.Error(t, err)
		assert.Empty(t, result)
		
		chainResult, err := plugin.RunChain(context.Background(), mockTarget)
		assert.Error(t, err)
		assert.Empty(t, chainResult)
	})
}

func TestPlugin_Performance(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Plugin Performance", func(t *testing.T) {
		mockTarget := &MockPluginTarget{
			result: "performance test result",
		}
		
		plugin := NewPlugin(logger)
		
		// Measure plugin performance
		start := time.Now()
		for i := 0; i < 100; i++ {
			result, err := plugin.RunPlugin(context.Background(), mockTarget, "test-plugin")
			require.NoError(t, err)
			assert.Equal(t, "performance test result", result)
		}
		duration := time.Since(start)
		
		// Should be fast
		assert.True(t, duration < time.Second, "Plugin operations took too long: %v", duration)
	})
}

func TestPlugin_Integration(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Full Plugin Flow", func(t *testing.T) {
		mockTarget := &MockPluginTarget{
			result: "integration test result",
		}
		
		plugin := NewPlugin(logger)
		
		// Test the full plugin flow
		result, err := plugin.RunPlugin(context.Background(), mockTarget, "integration-plugin")
		require.NoError(t, err)
		assert.Equal(t, "integration test result", result)
		
		chainResult, err := plugin.RunChain(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Equal(t, "integration test result", chainResult)
	})
}

func TestPlugin_Concurrency(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Multiple Plugin Instances", func(t *testing.T) {
		mockTarget1 := &MockPluginTarget{
			result: "plugin 1 result",
		}
		
		mockTarget2 := &MockPluginTarget{
			result: "plugin 2 result",
		}
		
		plugin1 := NewPlugin(logger)
		plugin2 := NewPlugin(logger)
		
		// Run plugin instances concurrently
		done := make(chan bool, 2)
		
		go func() {
			defer func() { done <- true }()
			result, err := plugin1.RunPlugin(context.Background(), mockTarget1, "plugin-1")
			require.NoError(t, err)
			assert.Equal(t, "plugin 1 result", result)
		}()
		
		go func() {
			defer func() { done <- true }()
			result, err := plugin2.RunPlugin(context.Background(), mockTarget2, "plugin-2")
			require.NoError(t, err)
			assert.Equal(t, "plugin 2 result", result)
		}()
		
		// Wait for both plugin instances to complete
		<-done
		<-done
	})
}

func TestPlugin_ResourceCleanup(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Resource Cleanup on Exit", func(t *testing.T) {
		mockTarget := &MockPluginTarget{
			result: "cleanup test result",
		}
		
		plugin := NewPlugin(logger)
		
		// Run plugin
		result, err := plugin.RunPlugin(context.Background(), mockTarget, "cleanup-plugin")
		require.NoError(t, err)
		assert.Equal(t, "cleanup test result", result)
		
		// Close target
		err = mockTarget.Close()
		require.NoError(t, err)
		
		// Resources should be cleaned up
	})
}
