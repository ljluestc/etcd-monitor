package signals

import (
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupSignalHandler_Comprehensive(t *testing.T) {
	t.Run("SetupSignalHandler", func(t *testing.T) {
		stopCh := SetupSignalHandler()
		
		assert.NotNil(t, stopCh)
		assert.IsType(t, (<-chan struct{})(nil), stopCh)
	})

	// Note: Testing SetupSignalHandler called twice would cause a panic
	// because the channel is closed on first call. This is expected behavior.

	// Note: Testing signal handling is complex without actually sending signals
	// The SetupSignalHandler function is tested above for basic functionality

	t.Run("Shutdown signals are defined", func(t *testing.T) {
		// Test that shutdownSignals contains the expected signals
		expectedSignals := []os.Signal{os.Interrupt, syscall.SIGTERM}
		
		assert.Equal(t, len(expectedSignals), len(shutdownSignals))
		for i, expected := range expectedSignals {
			assert.Equal(t, expected, shutdownSignals[i])
		}
	})

	t.Run("OnlyOneSignalHandler channel", func(t *testing.T) {
		// Test that onlyOneSignalHandler is properly initialized
		assert.NotNil(t, onlyOneSignalHandler)
		assert.IsType(t, make(chan struct{}), onlyOneSignalHandler)
	})
}