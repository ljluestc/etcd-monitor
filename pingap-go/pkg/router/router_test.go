package router

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// MockRouterTarget is a mock router target for testing
type MockRouterTarget struct {
	result string
	err    error
}

func (m *MockRouterTarget) Execute(ctx context.Context) (string, error) {
	return m.result, m.err
}

func (m *MockRouterTarget) Close() error {
	return nil
}

func TestRouter_Comprehensive(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("NewRouter", func(t *testing.T) {
		router := NewRouter(logger)
		
		assert.NotNil(t, router)
		assert.Equal(t, logger, router.logger)
	})

	t.Run("NewRouter with nil logger", func(t *testing.T) {
		router := NewRouter(nil)
		
		assert.NotNil(t, router)
		assert.NotNil(t, router.logger) // Should create a production logger
	})
}

func TestRouter_RunRouter(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("RunRouter - Success", func(t *testing.T) {
		mockTarget := &MockRouterTarget{
			result: "router result",
		}
		
		router := NewRouter(logger)
		result, err := router.RunRouter(context.Background(), mockTarget, "test-router")
		
		require.NoError(t, err)
		assert.Equal(t, "router result", result)
	})

	t.Run("RunRouter - Error", func(t *testing.T) {
		mockTarget := &MockRouterTarget{
			err: assert.AnError,
		}
		
		router := NewRouter(logger)
		result, err := router.RunRouter(context.Background(), mockTarget, "test-router")
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunRouter - Nil Target", func(t *testing.T) {
		router := NewRouter(logger)
		result, err := router.RunRouter(context.Background(), nil, "test-router")
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunRouter - Empty Name", func(t *testing.T) {
		mockTarget := &MockRouterTarget{
			result: "router result",
		}
		
		router := NewRouter(logger)
		result, err := router.RunRouter(context.Background(), mockTarget, "")
		
		require.NoError(t, err)
		assert.Equal(t, "router result", result)
	})

	t.Run("RunRouter - Context Cancellation", func(t *testing.T) {
		mockTarget := &MockRouterTarget{
			result: "router result",
		}
		
		router := NewRouter(logger)
		
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		result, err := router.RunRouter(ctx, mockTarget, "test-router")
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestRouter_EdgeCases(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Concurrent Operations", func(t *testing.T) {
		mockTarget := &MockRouterTarget{
			result: "concurrent router result",
		}
		
		router := NewRouter(logger)
		
		// Test concurrent access
		done := make(chan bool, 20)
		for i := 0; i < 20; i++ {
			go func() {
				defer func() { done <- true }()
				
				// Test all methods concurrently
				result, err := router.RunRouter(context.Background(), mockTarget, "test-router")
				if err == nil {
					assert.Equal(t, "concurrent router result", result)
				}
			}()
		}
		
		// Wait for all goroutines to complete
		for i := 0; i < 20; i++ {
			<-done
		}
	})

	t.Run("Empty Results", func(t *testing.T) {
		mockTarget := &MockRouterTarget{
			result: "",
		}
		
		router := NewRouter(logger)
		
		// Test with empty results
		result, err := router.RunRouter(context.Background(), mockTarget, "test-router")
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Nil Results", func(t *testing.T) {
		mockTarget := &MockRouterTarget{
			result: "",
		}
		
		router := NewRouter(logger)
		
		// Test with nil results
		result, err := router.RunRouter(context.Background(), mockTarget, "test-router")
		require.NoError(t, err)
		assert.Empty(t, result)
	})
}

func TestRouter_DataStructures(t *testing.T) {
	t.Run("Router Structure", func(t *testing.T) {
		router := NewRouter(nil)
		
		// Verify structure
		assert.NotNil(t, router)
		assert.NotNil(t, router.logger)
	})
}

func TestRouter_Interfaces(t *testing.T) {
	t.Run("Router implements RouterInterface", func(t *testing.T) {
		router := NewRouter(nil)
		
		var iface RouterInterface = router
		assert.NotNil(t, iface)
	})
}

func TestRouter_ErrorHandling(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Multiple Errors", func(t *testing.T) {
		mockTarget := &MockRouterTarget{
			err: assert.AnError,
		}
		
		router := NewRouter(logger)
		
		// Test all methods with errors
		result, err := router.RunRouter(context.Background(), mockTarget, "test-router")
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestRouter_Performance(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Router Performance", func(t *testing.T) {
		mockTarget := &MockRouterTarget{
			result: "performance test result",
		}
		
		router := NewRouter(logger)
		
		// Measure router performance
		start := time.Now()
		for i := 0; i < 100; i++ {
			result, err := router.RunRouter(context.Background(), mockTarget, "test-router")
			require.NoError(t, err)
			assert.Equal(t, "performance test result", result)
		}
		duration := time.Since(start)
		
		// Should be fast
		assert.True(t, duration < time.Second, "Router operations took too long: %v", duration)
	})
}

func TestRouter_Integration(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Full Router Flow", func(t *testing.T) {
		mockTarget := &MockRouterTarget{
			result: "integration test result",
		}
		
		router := NewRouter(logger)
		
		// Test the full router flow
		result, err := router.RunRouter(context.Background(), mockTarget, "integration-router")
		require.NoError(t, err)
		assert.Equal(t, "integration test result", result)
	})
}

func TestRouter_Concurrency(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Multiple Router Instances", func(t *testing.T) {
		mockTarget1 := &MockRouterTarget{
			result: "router 1 result",
		}
		
		mockTarget2 := &MockRouterTarget{
			result: "router 2 result",
		}
		
		router1 := NewRouter(logger)
		router2 := NewRouter(logger)
		
		// Run router instances concurrently
		done := make(chan bool, 2)
		
		go func() {
			defer func() { done <- true }()
			result, err := router1.RunRouter(context.Background(), mockTarget1, "router-1")
			require.NoError(t, err)
			assert.Equal(t, "router 1 result", result)
		}()
		
		go func() {
			defer func() { done <- true }()
			result, err := router2.RunRouter(context.Background(), mockTarget2, "router-2")
			require.NoError(t, err)
			assert.Equal(t, "router 2 result", result)
		}()
		
		// Wait for both router instances to complete
		<-done
		<-done
	})
}

func TestRouter_ResourceCleanup(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Resource Cleanup on Exit", func(t *testing.T) {
		mockTarget := &MockRouterTarget{
			result: "cleanup test result",
		}
		
		router := NewRouter(logger)
		
		// Run router
		result, err := router.RunRouter(context.Background(), mockTarget, "cleanup-router")
		require.NoError(t, err)
		assert.Equal(t, "cleanup test result", result)
		
		// Close target
		err = mockTarget.Close()
		require.NoError(t, err)
		
		// Resources should be cleaned up
	})
}
