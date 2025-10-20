package signals

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSignalHandler_Comprehensive(t *testing.T) {
	t.Run("NewSignalHandler with valid config", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)
		assert.NotNil(t, handler.config)
		assert.NotNil(t, handler.logger)
	})

	t.Run("NewSignalHandler with nil config", func(t *testing.T) {
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(nil, logger)
		assert.NotNil(t, handler)
		assert.Nil(t, handler.config)
		assert.NotNil(t, handler.logger)
	})

	t.Run("NewSignalHandler with nil logger", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}

		handler := NewSignalHandler(config, nil)
		assert.NotNil(t, handler)
		assert.NotNil(t, handler.config)
		assert.NotNil(t, handler.logger) // Should create a production logger
	})

	t.Run("NewSignalHandler with empty config", func(t *testing.T) {
		config := &Config{}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)
		assert.NotNil(t, handler.config)
		assert.NotNil(t, handler.logger)
	})

	t.Run("Start handler", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		err := handler.Start()
		assert.NoError(t, err)
	})

	t.Run("Stop handler", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		err := handler.Start()
		require.NoError(t, err)

		handler.Stop()
		// Should not panic or error
	})

	t.Run("Stop handler without starting", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		handler.Stop()
		// Should not panic or error
	})

	t.Run("Wait for signal", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		err := handler.Start()
		require.NoError(t, err)

		// Wait for signal with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		sig := handler.WaitForSignal(ctx)
		assert.NotNil(t, sig)

		handler.Stop()
	})

	t.Run("Wait for signal with cancelled context", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		err := handler.Start()
		require.NoError(t, err)

		// Wait for signal with cancelled context
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		sig := handler.WaitForSignal(ctx)
		assert.Nil(t, sig) // Should return nil for cancelled context

		handler.Stop()
	})

	t.Run("Wait for signal with timeout context", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		err := handler.Start()
		require.NoError(t, err)

		// Wait for signal with timeout context
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		sig := handler.WaitForSignal(ctx)
		assert.Nil(t, sig) // Should return nil for timeout

		handler.Stop()
	})

	t.Run("Wait for signal with nil context", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		err := handler.Start()
		require.NoError(t, err)

		// Wait for signal with nil context
		sig := handler.WaitForSignal(nil)
		assert.NotNil(t, sig) // Should still work with nil context

		handler.Stop()
	})

	t.Run("GetSignalChannel", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		channel := handler.GetSignalChannel()
		assert.NotNil(t, channel)
	})

	t.Run("GetSignalChannel after start", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		err := handler.Start()
		require.NoError(t, err)

		channel := handler.GetSignalChannel()
		assert.NotNil(t, channel)

		handler.Stop()
	})

	t.Run("GetSignalChannel after stop", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		err := handler.Start()
		require.NoError(t, err)

		handler.Stop()

		channel := handler.GetSignalChannel()
		assert.NotNil(t, channel)
	})

	t.Run("GetSignalInfo", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		info := handler.GetSignalInfo()
		assert.NotNil(t, info)
	})

	t.Run("GetSignalInfo after start", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		err := handler.Start()
		require.NoError(t, err)

		info := handler.GetSignalInfo()
		assert.NotNil(t, info)

		handler.Stop()
	})

	t.Run("GetSignalInfo after stop", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		err := handler.Start()
		require.NoError(t, err)

		handler.Stop()

		info := handler.GetSignalInfo()
		assert.NotNil(t, info)
	})
}

func TestSignalHandler_EdgeCases(t *testing.T) {
	t.Run("SignalHandler with empty signals list", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)
		assert.Len(t, handler.config.Signals, 0)
	})

	t.Run("SignalHandler with nil signals list", func(t *testing.T) {
		config := &Config{
			Signals: nil,
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)
		assert.Nil(t, handler.config.Signals)
	})

	t.Run("SignalHandler with single signal", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)
		assert.Len(t, handler.config.Signals, 1)
		assert.Equal(t, syscall.SIGTERM, handler.config.Signals[0])
	})

	t.Run("SignalHandler with multiple signals", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)
		assert.Len(t, handler.config.Signals, 3)
		assert.Equal(t, syscall.SIGTERM, handler.config.Signals[0])
		assert.Equal(t, syscall.SIGINT, handler.config.Signals[1])
		assert.Equal(t, syscall.SIGHUP, handler.config.Signals[2])
	})

	t.Run("SignalHandler with duplicate signals", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)
		assert.Len(t, handler.config.Signals, 3)
		assert.Equal(t, syscall.SIGTERM, handler.config.Signals[0])
		assert.Equal(t, syscall.SIGTERM, handler.config.Signals[1])
		assert.Equal(t, syscall.SIGINT, handler.config.Signals[2])
	})

	t.Run("SignalHandler with very long signals list", func(t *testing.T) {
		signals := make([]os.Signal, 1000)
		for i := 0; i < 1000; i++ {
			signals[i] = syscall.SIGTERM
		}

		config := &Config{
			Signals: signals,
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)
		assert.Len(t, handler.config.Signals, 1000)
	})

	t.Run("SignalHandler with all available signals", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{
				syscall.SIGTERM,
				syscall.SIGINT,
				syscall.SIGHUP,
				syscall.SIGQUIT,
				syscall.SIGUSR1,
				syscall.SIGUSR2,
			},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)
		assert.Len(t, handler.config.Signals, 6)
	})

	t.Run("SignalHandler with invalid signals", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.Signal(999)}, // Invalid signal number
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)
		assert.Len(t, handler.config.Signals, 1)
		assert.Equal(t, syscall.Signal(999), handler.config.Signals[0])
	})

	t.Run("SignalHandler with zero signal", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.Signal(0)},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)
		assert.Len(t, handler.config.Signals, 1)
		assert.Equal(t, syscall.Signal(0), handler.config.Signals[0])
	})

	t.Run("SignalHandler with negative signal", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.Signal(-1)},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)
		assert.Len(t, handler.config.Signals, 1)
		assert.Equal(t, syscall.Signal(-1), handler.config.Signals[0])
	})

	t.Run("SignalHandler with maximum signal", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.Signal(255)},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)
		assert.Len(t, handler.config.Signals, 1)
		assert.Equal(t, syscall.Signal(255), handler.config.Signals[0])
	})
}

func TestSignalHandler_Performance(t *testing.T) {
	t.Run("SignalHandler creation performance", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Create 1000 SignalHandler instances
		for i := 0; i < 1000; i++ {
			handler := NewSignalHandler(config, logger)
			assert.NotNil(t, handler)
		}
		
		duration := time.Since(start)
		
		// Should create 1000 instances in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Creating 1000 SignalHandler instances took too long: %v", duration)
	})

	t.Run("SignalHandler with large signals list performance", func(t *testing.T) {
		signals := make([]os.Signal, 10000)
		for i := 0; i < 10000; i++ {
			signals[i] = syscall.SIGTERM
		}
		
		config := &Config{
			Signals: signals,
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		handler := NewSignalHandler(config, logger)
		
		duration := time.Since(start)
		
		// Should create handler with 10000 signals in less than 10ms
		assert.True(t, duration < 10*time.Millisecond, "Creating SignalHandler with 10000 signals took too long: %v", duration)
		assert.Len(t, handler.config.Signals, 10000)
	})

	t.Run("SignalHandler start/stop performance", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Start and stop 100 handlers
		for i := 0; i < 100; i++ {
			handler := NewSignalHandler(config, logger)
			err := handler.Start()
			require.NoError(t, err)
			handler.Stop()
		}
		
		duration := time.Since(start)
		
		// Should start and stop 100 handlers in less than 1 second
		assert.True(t, duration < 1*time.Second, "Starting and stopping 100 handlers took too long: %v", duration)
	})

	t.Run("SignalHandler wait for signal performance", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()
		
		handler := NewSignalHandler(config, logger)
		err := handler.Start()
		require.NoError(t, err)
		
		start := time.Now()
		
		// Wait for signal 100 times with timeout
		for i := 0; i < 100; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
			sig := handler.WaitForSignal(ctx)
			cancel()
			assert.Nil(t, sig) // Should timeout
		}
		
		duration := time.Since(start)
		
		// Should wait for signal 100 times in less than 1 second
		assert.True(t, duration < 1*time.Second, "Waiting for signal 100 times took too long: %v", duration)
		
		handler.Stop()
	})

	t.Run("SignalHandler get signal channel performance", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()
		
		handler := NewSignalHandler(config, logger)
		
		start := time.Now()
		
		// Get signal channel 1000 times
		for i := 0; i < 1000; i++ {
			channel := handler.GetSignalChannel()
			assert.NotNil(t, channel)
		}
		
		duration := time.Since(start)
		
		// Should get signal channel 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting signal channel 1000 times took too long: %v", duration)
	})

	t.Run("SignalHandler get signal info performance", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()
		
		handler := NewSignalHandler(config, logger)
		
		start := time.Now()
		
		// Get signal info 1000 times
		for i := 0; i < 1000; i++ {
			info := handler.GetSignalInfo()
			assert.NotNil(t, info)
		}
		
		duration := time.Since(start)
		
		// Should get signal info 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting signal info 1000 times took too long: %v", duration)
	})
}

func TestSignalHandler_Integration(t *testing.T) {
	t.Run("SignalHandler with all fields set", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)
		assert.NotNil(t, handler.config)
		assert.NotNil(t, handler.logger)
		assert.Len(t, handler.config.Signals, 3)
		assert.Equal(t, syscall.SIGTERM, handler.config.Signals[0])
		assert.Equal(t, syscall.SIGINT, handler.config.Signals[1])
		assert.Equal(t, syscall.SIGHUP, handler.config.Signals[2])
	})

	t.Run("SignalHandler with minimal fields", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)
		assert.NotNil(t, handler.config)
		assert.NotNil(t, handler.logger)
		assert.Len(t, handler.config.Signals, 1)
		assert.Equal(t, syscall.SIGTERM, handler.config.Signals[0])
	})

	t.Run("SignalHandler with empty fields", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)
		assert.NotNil(t, handler.config)
		assert.NotNil(t, handler.logger)
		assert.Len(t, handler.config.Signals, 0)
	})

	t.Run("SignalHandler with nil fields", func(t *testing.T) {
		config := &Config{
			Signals: nil,
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)
		assert.NotNil(t, handler.config)
		assert.NotNil(t, handler.logger)
		assert.Nil(t, handler.config.Signals)
	})

	t.Run("SignalHandler full lifecycle", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)

		// Start handler
		err := handler.Start()
		assert.NoError(t, err)

		// Get signal channel
		channel := handler.GetSignalChannel()
		assert.NotNil(t, channel)

		// Get signal info
		info := handler.GetSignalInfo()
		assert.NotNil(t, info)

		// Wait for signal with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		sig := handler.WaitForSignal(ctx)
		cancel()
		assert.NotNil(t, sig)

		// Stop handler
		handler.Stop()
	})

	t.Run("SignalHandler with multiple concurrent operations", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)

		// Start handler
		err := handler.Start()
		require.NoError(t, err)

		// Perform multiple concurrent operations
		done := make(chan bool, 10)
		
		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				// Get signal channel
				channel := handler.GetSignalChannel()
				assert.NotNil(t, channel)
				
				// Get signal info
				info := handler.GetSignalInfo()
				assert.NotNil(t, info)
				
				// Wait for signal with timeout
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
				sig := handler.WaitForSignal(ctx)
				cancel()
				assert.NotNil(t, sig)
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Stop handler
		handler.Stop()
	})

	t.Run("SignalHandler with different signal types", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{
				syscall.SIGTERM,
				syscall.SIGINT,
				syscall.SIGHUP,
				syscall.SIGQUIT,
				syscall.SIGUSR1,
				syscall.SIGUSR2,
			},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)

		// Start handler
		err := handler.Start()
		require.NoError(t, err)

		// Test with different signal types
		for _, sig := range config.Signals {
			// Get signal channel
			channel := handler.GetSignalChannel()
			assert.NotNil(t, channel)
			
			// Get signal info
			info := handler.GetSignalInfo()
			assert.NotNil(t, info)
			
			// Wait for signal with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			receivedSig := handler.WaitForSignal(ctx)
			cancel()
			assert.NotNil(t, receivedSig)
		}

		// Stop handler
		handler.Stop()
	})

	t.Run("SignalHandler with context cancellation", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)

		// Start handler
		err := handler.Start()
		require.NoError(t, err)

		// Wait for signal with context that gets cancelled
		ctx, cancel := context.WithCancel(context.Background())
		
		// Cancel context after a short delay
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()
		
		sig := handler.WaitForSignal(ctx)
		assert.Nil(t, sig) // Should return nil for cancelled context

		// Stop handler
		handler.Stop()
	})

	t.Run("SignalHandler with context timeout", func(t *testing.T) {
		config := &Config{
			Signals: []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		}
		logger, _ := zap.NewDevelopment()

		handler := NewSignalHandler(config, logger)
		assert.NotNil(t, handler)

		// Start handler
		err := handler.Start()
		require.NoError(t, err)

		// Wait for signal with context that times out
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		sig := handler.WaitForSignal(ctx)
		assert.Nil(t, sig) // Should return nil for timeout

		// Stop handler
		handler.Stop()
	})
}
