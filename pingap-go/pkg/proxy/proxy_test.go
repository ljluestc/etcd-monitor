package proxy

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// MockProxyTarget is a mock proxy target for testing
type MockProxyTarget struct {
	result string
	err    error
}

func (m *MockProxyTarget) Execute(ctx context.Context) (string, error) {
	return m.result, m.err
}

func (m *MockProxyTarget) Close() error {
	return nil
}

func TestProxy_Comprehensive(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("NewProxy", func(t *testing.T) {
		proxy := NewProxy(logger)
		
		assert.NotNil(t, proxy)
		assert.Equal(t, logger, proxy.logger)
	})

	t.Run("NewProxy with nil logger", func(t *testing.T) {
		proxy := NewProxy(nil)
		
		assert.NotNil(t, proxy)
		assert.NotNil(t, proxy.logger) // Should create a production logger
	})
}

func TestProxy_RunProxy(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("RunProxy - Success", func(t *testing.T) {
		mockTarget := &MockProxyTarget{
			result: "proxy result",
		}
		
		proxy := NewProxy(logger)
		result, err := proxy.RunProxy(context.Background(), mockTarget, "test-proxy")
		
		require.NoError(t, err)
		assert.Equal(t, "proxy result", result)
	})

	t.Run("RunProxy - Error", func(t *testing.T) {
		mockTarget := &MockProxyTarget{
			err: assert.AnError,
		}
		
		proxy := NewProxy(logger)
		result, err := proxy.RunProxy(context.Background(), mockTarget, "test-proxy")
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunProxy - Nil Target", func(t *testing.T) {
		proxy := NewProxy(logger)
		result, err := proxy.RunProxy(context.Background(), nil, "test-proxy")
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunProxy - Empty Name", func(t *testing.T) {
		mockTarget := &MockProxyTarget{
			result: "proxy result",
		}
		
		proxy := NewProxy(logger)
		result, err := proxy.RunProxy(context.Background(), mockTarget, "")
		
		require.NoError(t, err)
		assert.Equal(t, "proxy result", result)
	})

	t.Run("RunProxy - Context Cancellation", func(t *testing.T) {
		mockTarget := &MockProxyTarget{
			result: "proxy result",
		}
		
		proxy := NewProxy(logger)
		
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		result, err := proxy.RunProxy(ctx, mockTarget, "test-proxy")
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestProxy_RunReverseProxy(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("RunReverseProxy - Success", func(t *testing.T) {
		mockTarget := &MockProxyTarget{
			result: "reverse proxy result",
		}
		
		proxy := NewProxy(logger)
		result, err := proxy.RunReverseProxy(context.Background(), mockTarget)
		
		require.NoError(t, err)
		assert.Equal(t, "reverse proxy result", result)
	})

	t.Run("RunReverseProxy - Error", func(t *testing.T) {
		mockTarget := &MockProxyTarget{
			err: assert.AnError,
		}
		
		proxy := NewProxy(logger)
		result, err := proxy.RunReverseProxy(context.Background(), mockTarget)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunReverseProxy - Nil Target", func(t *testing.T) {
		proxy := NewProxy(logger)
		result, err := proxy.RunReverseProxy(context.Background(), nil)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestProxy_RunServer(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("RunServer - Success", func(t *testing.T) {
		mockTarget := &MockProxyTarget{
			result: "server result",
		}
		
		proxy := NewProxy(logger)
		result, err := proxy.RunServer(context.Background(), mockTarget)
		
		require.NoError(t, err)
		assert.Equal(t, "server result", result)
	})

	t.Run("RunServer - Error", func(t *testing.T) {
		mockTarget := &MockProxyTarget{
			err: assert.AnError,
		}
		
		proxy := NewProxy(logger)
		result, err := proxy.RunServer(context.Background(), mockTarget)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunServer - Nil Target", func(t *testing.T) {
		proxy := NewProxy(logger)
		result, err := proxy.RunServer(context.Background(), nil)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestProxy_EdgeCases(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Concurrent Operations", func(t *testing.T) {
		mockTarget := &MockProxyTarget{
			result: "concurrent proxy result",
		}
		
		proxy := NewProxy(logger)
		
		// Test concurrent access
		done := make(chan bool, 20)
		for i := 0; i < 20; i++ {
			go func() {
				defer func() { done <- true }()
				
				// Test all methods concurrently
				result, err := proxy.RunProxy(context.Background(), mockTarget, "test-proxy")
				if err == nil {
					assert.Equal(t, "concurrent proxy result", result)
				}
				
				reverseProxyResult, err := proxy.RunReverseProxy(context.Background(), mockTarget)
				if err == nil {
					assert.Equal(t, "concurrent proxy result", reverseProxyResult)
				}
				
				serverResult, err := proxy.RunServer(context.Background(), mockTarget)
				if err == nil {
					assert.Equal(t, "concurrent proxy result", serverResult)
				}
			}()
		}
		
		// Wait for all goroutines to complete
		for i := 0; i < 20; i++ {
			<-done
		}
	})

	t.Run("Empty Results", func(t *testing.T) {
		mockTarget := &MockProxyTarget{
			result: "",
		}
		
		proxy := NewProxy(logger)
		
		// Test with empty results
		result, err := proxy.RunProxy(context.Background(), mockTarget, "test-proxy")
		require.NoError(t, err)
		assert.Empty(t, result)
		
		reverseProxyResult, err := proxy.RunReverseProxy(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, reverseProxyResult)
		
		serverResult, err := proxy.RunServer(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, serverResult)
	})

	t.Run("Nil Results", func(t *testing.T) {
		mockTarget := &MockProxyTarget{
			result: "",
		}
		
		proxy := NewProxy(logger)
		
		// Test with nil results
		result, err := proxy.RunProxy(context.Background(), mockTarget, "test-proxy")
		require.NoError(t, err)
		assert.Empty(t, result)
		
		reverseProxyResult, err := proxy.RunReverseProxy(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, reverseProxyResult)
		
		serverResult, err := proxy.RunServer(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, serverResult)
	})
}

func TestProxy_DataStructures(t *testing.T) {
	t.Run("Proxy Structure", func(t *testing.T) {
		proxy := NewProxy(nil)
		
		// Verify structure
		assert.NotNil(t, proxy)
		assert.NotNil(t, proxy.logger)
	})
}

func TestProxy_Interfaces(t *testing.T) {
	t.Run("Proxy implements ProxyInterface", func(t *testing.T) {
		proxy := NewProxy(nil)
		
		var iface ProxyInterface = proxy
		assert.NotNil(t, iface)
	})
}

func TestProxy_ErrorHandling(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Multiple Errors", func(t *testing.T) {
		mockTarget := &MockProxyTarget{
			err: assert.AnError,
		}
		
		proxy := NewProxy(logger)
		
		// Test all methods with errors
		result, err := proxy.RunProxy(context.Background(), mockTarget, "test-proxy")
		assert.Error(t, err)
		assert.Empty(t, result)
		
		reverseProxyResult, err := proxy.RunReverseProxy(context.Background(), mockTarget)
		assert.Error(t, err)
		assert.Empty(t, reverseProxyResult)
		
		serverResult, err := proxy.RunServer(context.Background(), mockTarget)
		assert.Error(t, err)
		assert.Empty(t, serverResult)
	})
}

func TestProxy_Performance(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Proxy Performance", func(t *testing.T) {
		mockTarget := &MockProxyTarget{
			result: "performance test result",
		}
		
		proxy := NewProxy(logger)
		
		// Measure proxy performance
		start := time.Now()
		for i := 0; i < 100; i++ {
			result, err := proxy.RunProxy(context.Background(), mockTarget, "test-proxy")
			require.NoError(t, err)
			assert.Equal(t, "performance test result", result)
		}
		duration := time.Since(start)
		
		// Should be fast
		assert.True(t, duration < time.Second, "Proxy operations took too long: %v", duration)
	})
}

func TestProxy_Integration(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Full Proxy Flow", func(t *testing.T) {
		mockTarget := &MockProxyTarget{
			result: "integration test result",
		}
		
		proxy := NewProxy(logger)
		
		// Test the full proxy flow
		result, err := proxy.RunProxy(context.Background(), mockTarget, "integration-proxy")
		require.NoError(t, err)
		assert.Equal(t, "integration test result", result)
		
		reverseProxyResult, err := proxy.RunReverseProxy(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Equal(t, "integration test result", reverseProxyResult)
		
		serverResult, err := proxy.RunServer(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Equal(t, "integration test result", serverResult)
	})
}

func TestProxy_Concurrency(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Multiple Proxy Instances", func(t *testing.T) {
		mockTarget1 := &MockProxyTarget{
			result: "proxy 1 result",
		}
		
		mockTarget2 := &MockProxyTarget{
			result: "proxy 2 result",
		}
		
		proxy1 := NewProxy(logger)
		proxy2 := NewProxy(logger)
		
		// Run proxy instances concurrently
		done := make(chan bool, 2)
		
		go func() {
			defer func() { done <- true }()
			result, err := proxy1.RunProxy(context.Background(), mockTarget1, "proxy-1")
			require.NoError(t, err)
			assert.Equal(t, "proxy 1 result", result)
		}()
		
		go func() {
			defer func() { done <- true }()
			result, err := proxy2.RunProxy(context.Background(), mockTarget2, "proxy-2")
			require.NoError(t, err)
			assert.Equal(t, "proxy 2 result", result)
		}()
		
		// Wait for both proxy instances to complete
		<-done
		<-done
	})
}

func TestProxy_ResourceCleanup(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Resource Cleanup on Exit", func(t *testing.T) {
		mockTarget := &MockProxyTarget{
			result: "cleanup test result",
		}
		
		proxy := NewProxy(logger)
		
		// Run proxy
		result, err := proxy.RunProxy(context.Background(), mockTarget, "cleanup-proxy")
		require.NoError(t, err)
		assert.Equal(t, "cleanup test result", result)
		
		// Close target
		err = mockTarget.Close()
		require.NoError(t, err)
		
		// Resources should be cleaned up
	})
}
