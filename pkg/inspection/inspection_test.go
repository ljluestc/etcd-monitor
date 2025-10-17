package inspection

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// MockInspectionTarget is a mock inspection target for testing
type MockInspectionTarget struct {
	result string
	err    error
}

func (m *MockInspectionTarget) Execute(ctx context.Context) (string, error) {
	return m.result, m.err
}

func (m *MockInspectionTarget) Close() error {
	return nil
}

func TestInspection_Comprehensive(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("NewInspection", func(t *testing.T) {
		inspection := NewInspection(logger)
		
		assert.NotNil(t, inspection)
		assert.Equal(t, logger, inspection.logger)
	})

	t.Run("NewInspection with nil logger", func(t *testing.T) {
		inspection := NewInspection(nil)
		
		assert.NotNil(t, inspection)
		assert.NotNil(t, inspection.logger) // Should create a production logger
	})
}

func TestInspection_RunInspection(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("RunInspection - Success", func(t *testing.T) {
		mockTarget := &MockInspectionTarget{
			result: "inspection result",
		}
		
		inspection := NewInspection(logger)
		result, err := inspection.RunInspection(context.Background(), mockTarget, "test-inspection")
		
		require.NoError(t, err)
		assert.Equal(t, "inspection result", result)
	})

	t.Run("RunInspection - Error", func(t *testing.T) {
		mockTarget := &MockInspectionTarget{
			err: assert.AnError,
		}
		
		inspection := NewInspection(logger)
		result, err := inspection.RunInspection(context.Background(), mockTarget, "test-inspection")
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunInspection - Nil Target", func(t *testing.T) {
		inspection := NewInspection(logger)
		result, err := inspection.RunInspection(context.Background(), nil, "test-inspection")
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunInspection - Empty Name", func(t *testing.T) {
		mockTarget := &MockInspectionTarget{
			result: "inspection result",
		}
		
		inspection := NewInspection(logger)
		result, err := inspection.RunInspection(context.Background(), mockTarget, "")
		
		require.NoError(t, err)
		assert.Equal(t, "inspection result", result)
	})

	t.Run("RunInspection - Context Cancellation", func(t *testing.T) {
		mockTarget := &MockInspectionTarget{
			result: "inspection result",
		}
		
		inspection := NewInspection(logger)
		
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		result, err := inspection.RunInspection(ctx, mockTarget, "test-inspection")
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestInspection_RunAlarmInspection(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("RunAlarmInspection - Success", func(t *testing.T) {
		mockTarget := &MockInspectionTarget{
			result: "alarm inspection result",
		}
		
		inspection := NewInspection(logger)
		result, err := inspection.RunAlarmInspection(context.Background(), mockTarget)
		
		require.NoError(t, err)
		assert.Equal(t, "alarm inspection result", result)
	})

	t.Run("RunAlarmInspection - Error", func(t *testing.T) {
		mockTarget := &MockInspectionTarget{
			err: assert.AnError,
		}
		
		inspection := NewInspection(logger)
		result, err := inspection.RunAlarmInspection(context.Background(), mockTarget)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunAlarmInspection - Nil Target", func(t *testing.T) {
		inspection := NewInspection(logger)
		result, err := inspection.RunAlarmInspection(context.Background(), nil)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestInspection_RunConsistencyInspection(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("RunConsistencyInspection - Success", func(t *testing.T) {
		mockTarget := &MockInspectionTarget{
			result: "consistency inspection result",
		}
		
		inspection := NewInspection(logger)
		result, err := inspection.RunConsistencyInspection(context.Background(), mockTarget)
		
		require.NoError(t, err)
		assert.Equal(t, "consistency inspection result", result)
	})

	t.Run("RunConsistencyInspection - Error", func(t *testing.T) {
		mockTarget := &MockInspectionTarget{
			err: assert.AnError,
		}
		
		inspection := NewInspection(logger)
		result, err := inspection.RunConsistencyInspection(context.Background(), mockTarget)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunConsistencyInspection - Nil Target", func(t *testing.T) {
		inspection := NewInspection(logger)
		result, err := inspection.RunConsistencyInspection(context.Background(), nil)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestInspection_RunHealthyInspection(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("RunHealthyInspection - Success", func(t *testing.T) {
		mockTarget := &MockInspectionTarget{
			result: "healthy inspection result",
		}
		
		inspection := NewInspection(logger)
		result, err := inspection.RunHealthyInspection(context.Background(), mockTarget)
		
		require.NoError(t, err)
		assert.Equal(t, "healthy inspection result", result)
	})

	t.Run("RunHealthyInspection - Error", func(t *testing.T) {
		mockTarget := &MockInspectionTarget{
			err: assert.AnError,
		}
		
		inspection := NewInspection(logger)
		result, err := inspection.RunHealthyInspection(context.Background(), mockTarget)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunHealthyInspection - Nil Target", func(t *testing.T) {
		inspection := NewInspection(logger)
		result, err := inspection.RunHealthyInspection(context.Background(), nil)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestInspection_RunRequestInspection(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("RunRequestInspection - Success", func(t *testing.T) {
		mockTarget := &MockInspectionTarget{
			result: "request inspection result",
		}
		
		inspection := NewInspection(logger)
		result, err := inspection.RunRequestInspection(context.Background(), mockTarget)
		
		require.NoError(t, err)
		assert.Equal(t, "request inspection result", result)
	})

	t.Run("RunRequestInspection - Error", func(t *testing.T) {
		mockTarget := &MockInspectionTarget{
			err: assert.AnError,
		}
		
		inspection := NewInspection(logger)
		result, err := inspection.RunRequestInspection(context.Background(), mockTarget)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("RunRequestInspection - Nil Target", func(t *testing.T) {
		inspection := NewInspection(logger)
		result, err := inspection.RunRequestInspection(context.Background(), nil)
		
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestInspection_EdgeCases(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Concurrent Operations", func(t *testing.T) {
		mockTarget := &MockInspectionTarget{
			result: "concurrent inspection result",
		}
		
		inspection := NewInspection(logger)
		
		// Test concurrent access
		done := make(chan bool, 20)
		for i := 0; i < 20; i++ {
			go func() {
				defer func() { done <- true }()
				
				// Test all methods concurrently
				result, err := inspection.RunInspection(context.Background(), mockTarget, "test-inspection")
				if err == nil {
					assert.Equal(t, "concurrent inspection result", result)
				}
				
				alarmResult, err := inspection.RunAlarmInspection(context.Background(), mockTarget)
				if err == nil {
					assert.Equal(t, "concurrent inspection result", alarmResult)
				}
				
				consistencyResult, err := inspection.RunConsistencyInspection(context.Background(), mockTarget)
				if err == nil {
					assert.Equal(t, "concurrent inspection result", consistencyResult)
				}
				
				healthyResult, err := inspection.RunHealthyInspection(context.Background(), mockTarget)
				if err == nil {
					assert.Equal(t, "concurrent inspection result", healthyResult)
				}
				
				requestResult, err := inspection.RunRequestInspection(context.Background(), mockTarget)
				if err == nil {
					assert.Equal(t, "concurrent inspection result", requestResult)
				}
			}()
		}
		
		// Wait for all goroutines to complete
		for i := 0; i < 20; i++ {
			<-done
		}
	})

	t.Run("Empty Results", func(t *testing.T) {
		mockTarget := &MockInspectionTarget{
			result: "",
		}
		
		inspection := NewInspection(logger)
		
		// Test with empty results
		result, err := inspection.RunInspection(context.Background(), mockTarget, "test-inspection")
		require.NoError(t, err)
		assert.Empty(t, result)
		
		alarmResult, err := inspection.RunAlarmInspection(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, alarmResult)
		
		consistencyResult, err := inspection.RunConsistencyInspection(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, consistencyResult)
		
		healthyResult, err := inspection.RunHealthyInspection(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, healthyResult)
		
		requestResult, err := inspection.RunRequestInspection(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, requestResult)
	})

	t.Run("Nil Results", func(t *testing.T) {
		mockTarget := &MockInspectionTarget{
			result: "",
		}
		
		inspection := NewInspection(logger)
		
		// Test with nil results
		result, err := inspection.RunInspection(context.Background(), mockTarget, "test-inspection")
		require.NoError(t, err)
		assert.Empty(t, result)
		
		alarmResult, err := inspection.RunAlarmInspection(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, alarmResult)
		
		consistencyResult, err := inspection.RunConsistencyInspection(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, consistencyResult)
		
		healthyResult, err := inspection.RunHealthyInspection(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, healthyResult)
		
		requestResult, err := inspection.RunRequestInspection(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Empty(t, requestResult)
	})
}

func TestInspection_DataStructures(t *testing.T) {
	t.Run("Inspection Structure", func(t *testing.T) {
		inspection := NewInspection(nil)
		
		// Verify structure
		assert.NotNil(t, inspection)
		assert.NotNil(t, inspection.logger)
	})
}

func TestInspection_Interfaces(t *testing.T) {
	t.Run("Inspection implements InspectionInterface", func(t *testing.T) {
		inspection := NewInspection(nil)
		
		var iface InspectionInterface = inspection
		assert.NotNil(t, iface)
	})
}

func TestInspection_ErrorHandling(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Multiple Errors", func(t *testing.T) {
		mockTarget := &MockInspectionTarget{
			err: assert.AnError,
		}
		
		inspection := NewInspection(logger)
		
		// Test all methods with errors
		result, err := inspection.RunInspection(context.Background(), mockTarget, "test-inspection")
		assert.Error(t, err)
		assert.Empty(t, result)
		
		alarmResult, err := inspection.RunAlarmInspection(context.Background(), mockTarget)
		assert.Error(t, err)
		assert.Empty(t, alarmResult)
		
		consistencyResult, err := inspection.RunConsistencyInspection(context.Background(), mockTarget)
		assert.Error(t, err)
		assert.Empty(t, consistencyResult)
		
		healthyResult, err := inspection.RunHealthyInspection(context.Background(), mockTarget)
		assert.Error(t, err)
		assert.Empty(t, healthyResult)
		
		requestResult, err := inspection.RunRequestInspection(context.Background(), mockTarget)
		assert.Error(t, err)
		assert.Empty(t, requestResult)
	})
}

func TestInspection_Performance(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Inspection Performance", func(t *testing.T) {
		mockTarget := &MockInspectionTarget{
			result: "performance test result",
		}
		
		inspection := NewInspection(logger)
		
		// Measure inspection performance
		start := time.Now()
		for i := 0; i < 100; i++ {
			result, err := inspection.RunInspection(context.Background(), mockTarget, "test-inspection")
			require.NoError(t, err)
			assert.Equal(t, "performance test result", result)
		}
		duration := time.Since(start)
		
		// Should be fast
		assert.True(t, duration < time.Second, "Inspection operations took too long: %v", duration)
	})
}

func TestInspection_Integration(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Full Inspection Flow", func(t *testing.T) {
		mockTarget := &MockInspectionTarget{
			result: "integration test result",
		}
		
		inspection := NewInspection(logger)
		
		// Test the full inspection flow
		result, err := inspection.RunInspection(context.Background(), mockTarget, "integration-inspection")
		require.NoError(t, err)
		assert.Equal(t, "integration test result", result)
		
		alarmResult, err := inspection.RunAlarmInspection(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Equal(t, "integration test result", alarmResult)
		
		consistencyResult, err := inspection.RunConsistencyInspection(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Equal(t, "integration test result", consistencyResult)
		
		healthyResult, err := inspection.RunHealthyInspection(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Equal(t, "integration test result", healthyResult)
		
		requestResult, err := inspection.RunRequestInspection(context.Background(), mockTarget)
		require.NoError(t, err)
		assert.Equal(t, "integration test result", requestResult)
	})
}

func TestInspection_Concurrency(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Multiple Inspections", func(t *testing.T) {
		mockTarget1 := &MockInspectionTarget{
			result: "inspection 1 result",
		}
		
		mockTarget2 := &MockInspectionTarget{
			result: "inspection 2 result",
		}
		
		inspection1 := NewInspection(logger)
		inspection2 := NewInspection(logger)
		
		// Run inspections concurrently
		done := make(chan bool, 2)
		
		go func() {
			defer func() { done <- true }()
			result, err := inspection1.RunInspection(context.Background(), mockTarget1, "inspection-1")
			require.NoError(t, err)
			assert.Equal(t, "inspection 1 result", result)
		}()
		
		go func() {
			defer func() { done <- true }()
			result, err := inspection2.RunInspection(context.Background(), mockTarget2, "inspection-2")
			require.NoError(t, err)
			assert.Equal(t, "inspection 2 result", result)
		}()
		
		// Wait for both inspections to complete
		<-done
		<-done
	})
}

func TestInspection_ResourceCleanup(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Resource Cleanup on Exit", func(t *testing.T) {
		mockTarget := &MockInspectionTarget{
			result: "cleanup test result",
		}
		
		inspection := NewInspection(logger)
		
		// Run inspection
		result, err := inspection.RunInspection(context.Background(), mockTarget, "cleanup-inspection")
		require.NoError(t, err)
		assert.Equal(t, "cleanup test result", result)
		
		// Close target
		err = mockTarget.Close()
		require.NoError(t, err)
		
		// Resources should be cleaned up
	})
}