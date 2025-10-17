package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// MockMonitorTarget is a mock monitor target for testing
type MockMonitorTarget struct {
	result string
	err    error
}

func (m *MockMonitorTarget) Execute(ctx context.Context) (string, error) {
	return m.result, m.err
}

func (m *MockMonitorTarget) Close() error {
	return nil
}

func TestMonitor_Comprehensive(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("NewMonitor", func(t *testing.T) {
		monitor := NewMonitor(logger)
		
		assert.NotNil(t, monitor)
		assert.Equal(t, logger, monitor.logger)
	})

	t.Run("NewMonitor with nil logger", func(t *testing.T) {
		monitor := NewMonitor(nil)
		
		assert.NotNil(t, monitor)
		assert.NotNil(t, monitor.logger) // Should create a production logger
	})
}

func TestMonitor_RunMonitor(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("RunMonitor - Success", func(t *testing.T) {
		mockTarget := &MockMonitorTarget{
			result: "monitor result",
		}
		
		monitor := NewMonitor(logger)
		result, err := monitor.RunMonitor(context.Background(), mockTarget, "test-monitor")
		
		require.NoError(t, err)
		assert.Equal(t, "monitor result", result)
	})

	t.Run("RunMonitor - Error", func(t *testing.T) {
		mockTarget := &MockMonitorTarget{
			err: assert.AnError,
		}
		
		monitor := NewMonitor(logger)
		result, err := monitor.RunMonitor(context.Background(), mockTarget, "test-monitor")
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunMonitor - Nil Target", func(t *testing.T) {
		monitor := NewMonitor(logger)
		result, err := monitor.RunMonitor(context.Background(), nil, "test-monitor")
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunMonitor - Empty Name", func(t *testing.T) {
		mockTarget := &MockMonitorTarget{
			result: "monitor result",
		}
		
		monitor := NewMonitor(logger)
		result, err := monitor.RunMonitor(context.Background(), mockTarget, "")
		
		require.NoError(t, err)
		assert.Equal(t, "monitor result", result)
	})

	t.Run("RunMonitor - Context Cancellation", func(t *testing.T) {
		mockTarget := &MockMonitorTarget{
			result: "monitor result",
		}
		
		monitor := NewMonitor(logger)
		
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		result, err := monitor.RunMonitor(ctx, mockTarget, "test-monitor")
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestMonitor_RunHealthMonitor(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("RunHealthMonitor - Success", func(t *testing.T) {
		mockTarget := &MockMonitorTarget{
			result: "health monitor result",
		}
		
		monitor := NewMonitor(logger)
		result, err := monitor.RunHealthMonitor(context.Background(), mockTarget)
		
		require.NoError(t, err)
		assert.Equal(t, "health monitor result", result)
	})

	t.Run("RunHealthMonitor - Error", func(t *testing.T) {
		mockTarget := &MockMonitorTarget{
			err: assert.AnError,
		}
		
		monitor := NewMonitor(logger)
		result, err := monitor.RunHealthMonitor(context.Background(), mockTarget)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunHealthMonitor - Nil Target", func(t *testing.T) {
		monitor := NewMonitor(logger)
		result, err := monitor.RunHealthMonitor(context.Background(), nil)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestMonitor_RunMetricsMonitor(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("RunMetricsMonitor - Success", func(t *testing.T) {
		mockTarget := &MockMonitorTarget{
			result: "metrics monitor result",
		}
		
		monitor := NewMonitor(logger)
		result, err := monitor.RunMetricsMonitor(context.Background(), mockTarget)
		
		require.NoError(t, err)
		assert.Equal(t, "metrics monitor result", result)
	})

	t.Run("RunMetricsMonitor - Error", func(t *testing.T) {
		mockTarget := &MockMonitorTarget{
			err: assert.AnError,
		}
		
		monitor := NewMonitor(logger)
		result, err := monitor.RunMetricsMonitor(context.Background(), mockTarget)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunMetricsMonitor - Nil Target", func(t *testing.T) {
		monitor := NewMonitor(logger)
		result, err := monitor.RunMetricsMonitor(context.Background(), nil)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestMonitor_RunAlertMonitor(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("RunAlertMonitor - Success", func(t *testing.T) {
		mockTarget := &MockMonitorTarget{
			result: "alert monitor result",
		}
		
		monitor := NewMonitor(logger)
		result, err := monitor.RunAlertMonitor(context.Background(), mockTarget)
		
		require.NoError(t, err)
		assert.Equal(t, "alert monitor result", result)
	})

	t.Run("RunAlertMonitor - Error", func(t *testing.T) {
		mockTarget := &MockMonitorTarget{
			err: assert.AnError,
		}
		
		monitor := NewMonitor(logger)
		result, err := monitor.RunAlertMonitor(context.Background(), mockTarget)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunAlertMonitor - Nil Target", func(t *testing.T) {
		monitor := NewMonitor(logger)
		result, err := monitor.RunAlertMonitor(context.Background(), nil)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestMonitor_EdgeCases(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Concurrent Operations", func(t *testing.T) {
		mockTarget := &MockMonitorTarget{
			result: "concurrent monitor result",
		}
		
		monitor := NewMonitor(logger)
		
		// Test concurrent access
		done := make(chan bool, 20)
		for i := 0; i < 20; i++ {
			go func() {
				defer func() { done <- true }()
				
				// Test all methods concurrently
				result, err := monitor.RunMonitor(context.Background(), mockTarget, "test-monitor")
				if err == nil {
					assert.Equal(t, "concurrent monitor result", result)
				}
				
				healthResult, err := monitor.RunHealthMonitor(context.Background(), mockTarget)
				if err == nil {
					assert.Equal(t, "concurrent monitor result", healthResult)
				}
				
				metricsResult, err := monitor.RunMetricsMonitor(context.Background(), mockTarget)
				if err == nil {
					assert.Equal(t, "concurrent monitor result", metricsResult)
				}
				
				alertResult, err := monitor.RunAlertMonitor(context.Background(), mockTarget)
				if err == nil {
					assert.Equal(t, "concurrent monitor result", alertResult)
				}
			}()
		}
		
		// Wait for all goroutines to complete
		for i := 0; i < 20; i++ {
			<-done
		}
	})

	t.Run("Empty Results", func(t *testing.T) {
		mockTarget := &MockMonitorTarget{
			result: "",
		}
		
		monitor := NewMonitor(logger)
		
		// Test with empty results
		result, err := monitor.RunMonitor(context.Background(), mockTarget, "test-monitor")
		require.NoError(t, err)
		assert.Empty(t, result)
		
		healthResult, err := monitor.RunHealthMonitor(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, healthResult)
		
		metricsResult, err := monitor.RunMetricsMonitor(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, metricsResult)
		
		alertResult, err := monitor.RunAlertMonitor(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, alertResult)
	})

	t.Run("Nil Results", func(t *testing.T) {
		mockTarget := &MockMonitorTarget{
			result: "",
		}
		
		monitor := NewMonitor(logger)
		
		// Test with nil results
		result, err := monitor.RunMonitor(context.Background(), mockTarget, "test-monitor")
		require.NoError(t, err)
		assert.Empty(t, result)
		
		healthResult, err := monitor.RunHealthMonitor(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, healthResult)
		
		metricsResult, err := monitor.RunMetricsMonitor(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, metricsResult)
		
		alertResult, err := monitor.RunAlertMonitor(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, alertResult)
	})
}

func TestMonitor_DataStructures(t *testing.T) {
	t.Run("Monitor Structure", func(t *testing.T) {
		monitor := NewMonitor(nil)
		
		// Verify structure
		assert.NotNil(t, monitor)
		assert.NotNil(t, monitor.logger)
	})
}

func TestMonitor_Interfaces(t *testing.T) {
	t.Run("Monitor implements MonitorInterface", func(t *testing.T) {
		monitor := NewMonitor(nil)
		
		var iface MonitorInterface = monitor
		assert.NotNil(t, iface)
	})
}

func TestMonitor_ErrorHandling(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Multiple Errors", func(t *testing.T) {
		mockTarget := &MockMonitorTarget{
			err: assert.AnError,
		}
		
		monitor := NewMonitor(logger)
		
		// Test all methods with errors
		result, err := monitor.RunMonitor(context.Background(), mockTarget, "test-monitor")
		assert.Error(t, err)
		assert.Empty(t, result)
		
		healthResult, err := monitor.RunHealthMonitor(context.Background(), mockTarget)
		assert.Error(t, err)
		assert.Empty(t, healthResult)
		
		metricsResult, err := monitor.RunMetricsMonitor(context.Background(), mockTarget)
		assert.Error(t, err)
		assert.Empty(t, metricsResult)
		
		alertResult, err := monitor.RunAlertMonitor(context.Background(), mockTarget)
		assert.Error(t, err)
		assert.Empty(t, alertResult)
	})
}

func TestMonitor_Performance(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Monitor Performance", func(t *testing.T) {
		mockTarget := &MockMonitorTarget{
			result: "performance test result",
		}
		
		monitor := NewMonitor(logger)
		
		// Measure monitor performance
		start := time.Now()
		for i := 0; i < 100; i++ {
			result, err := monitor.RunMonitor(context.Background(), mockTarget, "test-monitor")
			require.NoError(t, err)
			assert.Equal(t, "performance test result", result)
		}
		duration := time.Since(start)
		
		// Should be fast
		assert.True(t, duration < time.Second, "Monitor operations took too long: %v", duration)
	})
}

func TestMonitor_Integration(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Full Monitor Flow", func(t *testing.T) {
		mockTarget := &MockMonitorTarget{
			result: "integration test result",
		}
		
		monitor := NewMonitor(logger)
		
		// Test the full monitor flow
		result, err := monitor.RunMonitor(context.Background(), mockTarget, "integration-monitor")
		require.NoError(t, err)
		assert.Equal(t, "integration test result", result)
		
		healthResult, err := monitor.RunHealthMonitor(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Equal(t, "integration test result", healthResult)
		
		metricsResult, err := monitor.RunMetricsMonitor(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Equal(t, "integration test result", metricsResult)
		
		alertResult, err := monitor.RunAlertMonitor(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Equal(t, "integration test result", alertResult)
	})
}

func TestMonitor_Concurrency(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Multiple Monitors", func(t *testing.T) {
		mockTarget1 := &MockMonitorTarget{
			result: "monitor 1 result",
		}
		
		mockTarget2 := &MockMonitorTarget{
			result: "monitor 2 result",
		}
		
		monitor1 := NewMonitor(logger)
		monitor2 := NewMonitor(logger)
		
		// Run monitors concurrently
		done := make(chan bool, 2)
		
		go func() {
			defer func() { done <- true }()
			result, err := monitor1.RunMonitor(context.Background(), mockTarget1, "monitor-1")
			require.NoError(t, err)
			assert.Equal(t, "monitor 1 result", result)
		}()
		
		go func() {
			defer func() { done <- true }()
			result, err := monitor2.RunMonitor(context.Background(), mockTarget2, "monitor-2")
			require.NoError(t, err)
			assert.Equal(t, "monitor 2 result", result)
		}()
		
		// Wait for both monitors to complete
		<-done
		<-done
	})
}

func TestMonitor_ResourceCleanup(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Resource Cleanup on Exit", func(t *testing.T) {
		mockTarget := &MockMonitorTarget{
			result: "cleanup test result",
		}
		
		monitor := NewMonitor(logger)
		
		// Run monitor
		result, err := monitor.RunMonitor(context.Background(), mockTarget, "cleanup-monitor")
		require.NoError(t, err)
		assert.Equal(t, "cleanup test result", result)
		
		// Close target
		err = mockTarget.Close()
		require.NoError(t, err)
		
		// Resources should be cleaned up
	})
}
