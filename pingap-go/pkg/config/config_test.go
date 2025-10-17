package config

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// MockConfigTarget is a mock config target for testing
type MockConfigTarget struct {
	result string
	err    error
}

func (m *MockConfigTarget) Execute(ctx context.Context) (string, error) {
	return m.result, m.err
}

func (m *MockConfigTarget) Close() error {
	return nil
}

func TestConfig_Comprehensive(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("NewConfig", func(t *testing.T) {
		config := NewConfig(logger)
		
		assert.NotNil(t, config)
		assert.Equal(t, logger, config.logger)
	})

	t.Run("NewConfig with nil logger", func(t *testing.T) {
		config := NewConfig(nil)
		
		assert.NotNil(t, config)
		assert.NotNil(t, config.logger) // Should create a production logger
	})
}

func TestConfig_RunConfig(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("RunConfig - Success", func(t *testing.T) {
		mockTarget := &MockConfigTarget{
			result: "config result",
		}
		
		config := NewConfig(logger)
		result, err := config.RunConfig(context.Background(), mockTarget, "test-config")
		
		require.NoError(t, err)
		assert.Equal(t, "config result", result)
	})

	t.Run("RunConfig - Error", func(t *testing.T) {
		mockTarget := &MockConfigTarget{
			err: assert.AnError,
		}
		
		config := NewConfig(logger)
		result, err := config.RunConfig(context.Background(), mockTarget, "test-config")
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunConfig - Nil Target", func(t *testing.T) {
		config := NewConfig(logger)
		result, err := config.RunConfig(context.Background(), nil, "test-config")
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunConfig - Empty Name", func(t *testing.T) {
		mockTarget := &MockConfigTarget{
			result: "config result",
		}
		
		config := NewConfig(logger)
		result, err := config.RunConfig(context.Background(), mockTarget, "")
		
		require.NoError(t, err)
		assert.Equal(t, "config result", result)
	})

	t.Run("RunConfig - Context Cancellation", func(t *testing.T) {
		mockTarget := &MockConfigTarget{
			result: "config result",
		}
		
		config := NewConfig(logger)
		
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		result, err := config.RunConfig(ctx, mockTarget, "test-config")
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestConfig_RunLoader(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("RunLoader - Success", func(t *testing.T) {
		mockTarget := &MockConfigTarget{
			result: "loader result",
		}
		
		config := NewConfig(logger)
		result, err := config.RunLoader(context.Background(), mockTarget)
		
		require.NoError(t, err)
		assert.Equal(t, "loader result", result)
	})

	t.Run("RunLoader - Error", func(t *testing.T) {
		mockTarget := &MockConfigTarget{
			err: assert.AnError,
		}
		
		config := NewConfig(logger)
		result, err := config.RunLoader(context.Background(), mockTarget)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunLoader - Nil Target", func(t *testing.T) {
		config := NewConfig(logger)
		result, err := config.RunLoader(context.Background(), nil)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestConfig_RunTypes(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("RunTypes - Success", func(t *testing.T) {
		mockTarget := &MockConfigTarget{
			result: "types result",
		}
		
		config := NewConfig(logger)
		result, err := config.RunTypes(context.Background(), mockTarget)
		
		require.NoError(t, err)
		assert.Equal(t, "types result", result)
	})

	t.Run("RunTypes - Error", func(t *testing.T) {
		mockTarget := &MockConfigTarget{
			err: assert.AnError,
		}
		
		config := NewConfig(logger)
		result, err := config.RunTypes(context.Background(), mockTarget)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunTypes - Nil Target", func(t *testing.T) {
		config := NewConfig(logger)
		result, err := config.RunTypes(context.Background(), nil)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestConfig_EdgeCases(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Concurrent Operations", func(t *testing.T) {
		mockTarget := &MockConfigTarget{
			result: "concurrent config result",
		}
		
		config := NewConfig(logger)
		
		// Test concurrent access
		done := make(chan bool, 20)
		for i := 0; i < 20; i++ {
			go func() {
				defer func() { done <- true }()
				
				// Test all methods concurrently
				result, err := config.RunConfig(context.Background(), mockTarget, "test-config")
				if err == nil {
					assert.Equal(t, "concurrent config result", result)
				}
				
				loaderResult, err := config.RunLoader(context.Background(), mockTarget)
				if err == nil {
					assert.Equal(t, "concurrent config result", loaderResult)
				}
				
				typesResult, err := config.RunTypes(context.Background(), mockTarget)
				if err == nil {
					assert.Equal(t, "concurrent config result", typesResult)
				}
			}()
		}
		
		// Wait for all goroutines to complete
		for i := 0; i < 20; i++ {
			<-done
		}
	})

	t.Run("Empty Results", func(t *testing.T) {
		mockTarget := &MockConfigTarget{
			result: "",
		}
		
		config := NewConfig(logger)
		
		// Test with empty results
		result, err := config.RunConfig(context.Background(), mockTarget, "test-config")
		require.NoError(t, err)
		assert.Empty(t, result)
		
		loaderResult, err := config.RunLoader(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, loaderResult)
		
		typesResult, err := config.RunTypes(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, typesResult)
	})

	t.Run("Nil Results", func(t *testing.T) {
		mockTarget := &MockConfigTarget{
			result: "",
		}
		
		config := NewConfig(logger)
		
		// Test with nil results
		result, err := config.RunConfig(context.Background(), mockTarget, "test-config")
		require.NoError(t, err)
		assert.Empty(t, result)
		
		loaderResult, err := config.RunLoader(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, loaderResult)
		
		typesResult, err := config.RunTypes(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, typesResult)
	})
}

func TestConfig_DataStructures(t *testing.T) {
	t.Run("Config Structure", func(t *testing.T) {
		config := NewConfig(nil)
		
		// Verify structure
		assert.NotNil(t, config)
		assert.NotNil(t, config.logger)
	})
}

func TestConfig_Interfaces(t *testing.T) {
	t.Run("Config implements ConfigInterface", func(t *testing.T) {
		config := NewConfig(nil)
		
		var iface ConfigInterface = config
		assert.NotNil(t, iface)
	})
}

func TestConfig_ErrorHandling(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Multiple Errors", func(t *testing.T) {
		mockTarget := &MockConfigTarget{
			err: assert.AnError,
		}
		
		config := NewConfig(logger)
		
		// Test all methods with errors
		result, err := config.RunConfig(context.Background(), mockTarget, "test-config")
		assert.Error(t, err)
		assert.Empty(t, result)
		
		loaderResult, err := config.RunLoader(context.Background(), mockTarget)
		assert.Error(t, err)
		assert.Empty(t, loaderResult)
		
		typesResult, err := config.RunTypes(context.Background(), mockTarget)
		assert.Error(t, err)
		assert.Empty(t, typesResult)
	})
}

func TestConfig_Performance(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Config Performance", func(t *testing.T) {
		mockTarget := &MockConfigTarget{
			result: "performance test result",
		}
		
		config := NewConfig(logger)
		
		// Measure config performance
		start := time.Now()
		for i := 0; i < 100; i++ {
			result, err := config.RunConfig(context.Background(), mockTarget, "test-config")
			require.NoError(t, err)
			assert.Equal(t, "performance test result", result)
		}
		duration := time.Since(start)
		
		// Should be fast
		assert.True(t, duration < time.Second, "Config operations took too long: %v", duration)
	})
}

func TestConfig_Integration(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Full Config Flow", func(t *testing.T) {
		mockTarget := &MockConfigTarget{
			result: "integration test result",
		}
		
		config := NewConfig(logger)
		
		// Test the full config flow
		result, err := config.RunConfig(context.Background(), mockTarget, "integration-config")
		require.NoError(t, err)
		assert.Equal(t, "integration test result", result)
		
		loaderResult, err := config.RunLoader(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Equal(t, "integration test result", loaderResult)
		
		typesResult, err := config.RunTypes(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Equal(t, "integration test result", typesResult)
	})
}

func TestConfig_Concurrency(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Multiple Config Instances", func(t *testing.T) {
		mockTarget1 := &MockConfigTarget{
			result: "config 1 result",
		}
		
		mockTarget2 := &MockConfigTarget{
			result: "config 2 result",
		}
		
		config1 := NewConfig(logger)
		config2 := NewConfig(logger)
		
		// Run config instances concurrently
		done := make(chan bool, 2)
		
		go func() {
			defer func() { done <- true }()
			result, err := config1.RunConfig(context.Background(), mockTarget1, "config-1")
			require.NoError(t, err)
			assert.Equal(t, "config 1 result", result)
		}()
		
		go func() {
			defer func() { done <- true }()
			result, err := config2.RunConfig(context.Background(), mockTarget2, "config-2")
			require.NoError(t, err)
			assert.Equal(t, "config 2 result", result)
		}()
		
		// Wait for both config instances to complete
		<-done
		<-done
	})
}

func TestConfig_ResourceCleanup(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Resource Cleanup on Exit", func(t *testing.T) {
		mockTarget := &MockConfigTarget{
			result: "cleanup test result",
		}
		
		config := NewConfig(logger)
		
		// Run config
		result, err := config.RunConfig(context.Background(), mockTarget, "cleanup-config")
		require.NoError(t, err)
		assert.Equal(t, "cleanup test result", result)
		
		// Close target
		err = mockTarget.Close()
		require.NoError(t, err)
		
		// Resources should be cleaned up
	})
}
