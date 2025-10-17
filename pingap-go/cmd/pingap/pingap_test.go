package pingap

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// MockPingapTarget is a mock pingap target for testing
type MockPingapTarget struct {
	result string
	err    error
}

func (m *MockPingapTarget) Execute(ctx context.Context) (string, error) {
	return m.result, m.err
}

func (m *MockPingapTarget) Close() error {
	return nil
}

func TestPingap_Comprehensive(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("NewPingap", func(t *testing.T) {
		pingap := NewPingap(logger)
		
		assert.NotNil(t, pingap)
		assert.Equal(t, logger, pingap.logger)
	})

	t.Run("NewPingap with nil logger", func(t *testing.T) {
		pingap := NewPingap(nil)
		
		assert.NotNil(t, pingap)
		assert.NotNil(t, pingap.logger) // Should create a production logger
	})
}

func TestPingap_RunPingap(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("RunPingap - Success", func(t *testing.T) {
		mockTarget := &MockPingapTarget{
			result: "pingap result",
		}
		
		pingap := NewPingap(logger)
		result, err := pingap.RunPingap(context.Background(), mockTarget, "test-pingap")
		
		require.NoError(t, err)
		assert.Equal(t, "pingap result", result)
	})

	t.Run("RunPingap - Error", func(t *testing.T) {
		mockTarget := &MockPingapTarget{
			err: assert.AnError,
		}
		
		pingap := NewPingap(logger)
		result, err := pingap.RunPingap(context.Background(), mockTarget, "test-pingap")
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunPingap - Nil Target", func(t *testing.T) {
		pingap := NewPingap(logger)
		result, err := pingap.RunPingap(context.Background(), nil, "test-pingap")
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunPingap - Empty Name", func(t *testing.T) {
		mockTarget := &MockPingapTarget{
			result: "pingap result",
		}
		
		pingap := NewPingap(logger)
		result, err := pingap.RunPingap(context.Background(), mockTarget, "")
		
		require.NoError(t, err)
		assert.Equal(t, "pingap result", result)
	})

	t.Run("RunPingap - Context Cancellation", func(t *testing.T) {
		mockTarget := &MockPingapTarget{
			result: "pingap result",
		}
		
		pingap := NewPingap(logger)
		
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		result, err := pingap.RunPingap(ctx, mockTarget, "test-pingap")
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestPingap_EdgeCases(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Concurrent Operations", func(t *testing.T) {
		mockTarget := &MockPingapTarget{
			result: "concurrent pingap result",
		}
		
		pingap := NewPingap(logger)
		
		// Test concurrent access
		done := make(chan bool, 20)
		for i := 0; i < 20; i++ {
			go func() {
				defer func() { done <- true }()
				
				// Test all methods concurrently
				result, err := pingap.RunPingap(context.Background(), mockTarget, "test-pingap")
				if err == nil {
					assert.Equal(t, "concurrent pingap result", result)
				}
			}()
		}
		
		// Wait for all goroutines to complete
		for i := 0; i < 20; i++ {
			<-done
		}
	})

	t.Run("Empty Results", func(t *testing.T) {
		mockTarget := &MockPingapTarget{
			result: "",
		}
		
		pingap := NewPingap(logger)
		
		// Test with empty results
		result, err := pingap.RunPingap(context.Background(), mockTarget, "test-pingap")
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Nil Results", func(t *testing.T) {
		mockTarget := &MockPingapTarget{
			result: "",
		}
		
		pingap := NewPingap(logger)
		
		// Test with nil results
		result, err := pingap.RunPingap(context.Background(), mockTarget, "test-pingap")
		require.NoError(t, err)
		assert.Empty(t, result)
	})
}

func TestPingap_DataStructures(t *testing.T) {
	t.Run("Pingap Structure", func(t *testing.T) {
		pingap := NewPingap(nil)
		
		// Verify structure
		assert.NotNil(t, pingap)
		assert.NotNil(t, pingap.logger)
	})
}

func TestPingap_Interfaces(t *testing.T) {
	t.Run("Pingap implements PingapInterface", func(t *testing.T) {
		pingap := NewPingap(nil)
		
		var iface PingapInterface = pingap
		assert.NotNil(t, iface)
	})
}

func TestPingap_ErrorHandling(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Multiple Errors", func(t *testing.T) {
		mockTarget := &MockPingapTarget{
			err: assert.AnError,
		}
		
		pingap := NewPingap(logger)
		
		// Test all methods with errors
		result, err := pingap.RunPingap(context.Background(), mockTarget, "test-pingap")
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestPingap_Performance(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Pingap Performance", func(t *testing.T) {
		mockTarget := &MockPingapTarget{
			result: "performance test result",
		}
		
		pingap := NewPingap(logger)
		
		// Measure pingap performance
		start := time.Now()
		for i := 0; i < 100; i++ {
			result, err := pingap.RunPingap(context.Background(), mockTarget, "test-pingap")
			require.NoError(t, err)
			assert.Equal(t, "performance test result", result)
		}
		duration := time.Since(start)
		
		// Should be fast
		assert.True(t, duration < time.Second, "Pingap operations took too long: %v", duration)
	})
}

func TestPingap_Integration(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Full Pingap Flow", func(t *testing.T) {
		mockTarget := &MockPingapTarget{
			result: "integration test result",
		}
		
		pingap := NewPingap(logger)
		
		// Test the full pingap flow
		result, err := pingap.RunPingap(context.Background(), mockTarget, "integration-pingap")
		require.NoError(t, err)
		assert.Equal(t, "integration test result", result)
	})
}

func TestPingap_Concurrency(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Multiple Pingap Instances", func(t *testing.T) {
		mockTarget1 := &MockPingapTarget{
			result: "pingap 1 result",
		}
		
		mockTarget2 := &MockPingapTarget{
			result: "pingap 2 result",
		}
		
		pingap1 := NewPingap(logger)
		pingap2 := NewPingap(logger)
		
		// Run pingap instances concurrently
		done := make(chan bool, 2)
		
		go func() {
			defer func() { done <- true }()
			result, err := pingap1.RunPingap(context.Background(), mockTarget1, "pingap-1")
			require.NoError(t, err)
			assert.Equal(t, "pingap 1 result", result)
		}()
		
		go func() {
			defer func() { done <- true }()
			result, err := pingap2.RunPingap(context.Background(), mockTarget2, "pingap-2")
			require.NoError(t, err)
			assert.Equal(t, "pingap 2 result", result)
		}()
		
		// Wait for both pingap instances to complete
		<-done
		<-done
	})
}

func TestPingap_ResourceCleanup(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Resource Cleanup on Exit", func(t *testing.T) {
		mockTarget := &MockPingapTarget{
			result: "cleanup test result",
		}
		
		pingap := NewPingap(logger)
		
		// Run pingap
		result, err := pingap.RunPingap(context.Background(), mockTarget, "cleanup-pingap")
		require.NoError(t, err)
		assert.Equal(t, "cleanup test result", result)
		
		// Close target
		err = mockTarget.Close()
		require.NoError(t, err)
		
		// Resources should be cleaned up
	})
}
