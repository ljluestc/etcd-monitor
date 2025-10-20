package etcdmonitor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestEtcdMonitor_Comprehensive(t *testing.T) {
	t.Run("NewEtcdMonitor with valid config", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.NotNil(t, monitor.config)
		assert.NotNil(t, monitor.logger)
	})

	t.Run("NewEtcdMonitor with nil config", func(t *testing.T) {
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(nil, logger)
		assert.NotNil(t, monitor)
		assert.Nil(t, monitor.config)
		assert.NotNil(t, monitor.logger)
	})

	t.Run("NewEtcdMonitor with nil logger", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}

		monitor := NewEtcdMonitor(config, nil)
		assert.NotNil(t, monitor)
		assert.NotNil(t, monitor.config)
		assert.NotNil(t, monitor.logger) // Should create a production logger
	})

	t.Run("NewEtcdMonitor with empty config", func(t *testing.T) {
		config := &Config{}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.NotNil(t, monitor.config)
		assert.NotNil(t, monitor.logger)
	})

	t.Run("Start monitor", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		err := monitor.Start()
		assert.NoError(t, err)
	})

	t.Run("Stop monitor", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		err := monitor.Start()
		require.NoError(t, err)

		monitor.Stop()
		// Should not panic or error
	})

	t.Run("Stop monitor without starting", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		monitor.Stop()
		// Should not panic or error
	})

	t.Run("Run monitor", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		ctx := context.Background()
		
		err := monitor.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("Run monitor with nil context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		
		err := monitor.Run(nil)
		assert.Error(t, err)
	})

	t.Run("Run monitor with cancelled context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		err := monitor.Run(ctx)
		assert.Error(t, err)
	})

	t.Run("Run monitor with timeout context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		err := monitor.Run(ctx)
		assert.Error(t, err)
	})

	t.Run("GetMonitorResults", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		results := monitor.GetMonitorResults()
		assert.NotNil(t, results)
	})

	t.Run("GetMonitorResults after run", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		ctx := context.Background()
		
		err := monitor.Run(ctx)
		require.NoError(t, err)
		
		results := monitor.GetMonitorResults()
		assert.NotNil(t, results)
	})

	t.Run("GetMonitorStatus", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		status := monitor.GetMonitorStatus()
		assert.NotNil(t, status)
	})

	t.Run("GetMonitorStatus after run", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		ctx := context.Background()
		
		err := monitor.Run(ctx)
		require.NoError(t, err)
		
		status := monitor.GetMonitorStatus()
		assert.NotNil(t, status)
	})

	t.Run("GetMonitorMetrics", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		metrics := monitor.GetMonitorMetrics()
		assert.NotNil(t, metrics)
	})

	t.Run("GetMonitorMetrics after run", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		ctx := context.Background()
		
		err := monitor.Run(ctx)
		require.NoError(t, err)
		
		metrics := monitor.GetMonitorMetrics()
		assert.NotNil(t, metrics)
	})
}

func TestEtcdMonitor_EdgeCases(t *testing.T) {
	t.Run("Monitor with very long endpoint list", func(t *testing.T) {
		endpoints := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			endpoints[i] = fmt.Sprintf("localhost:%d", 2379+i)
		}

		config := &Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.Len(t, monitor.config.Endpoints, 1000)
	})

	t.Run("Monitor with very long timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 24 * time.Hour,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.Equal(t, 24*time.Hour, monitor.config.DialTimeout)
	})

	t.Run("Monitor with very short timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 1 * time.Nanosecond,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.Equal(t, 1*time.Nanosecond, monitor.config.DialTimeout)
	})

	t.Run("Monitor with zero timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 0,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.Equal(t, time.Duration(0), monitor.config.DialTimeout)
	})

	t.Run("Monitor with negative timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: -1 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.Equal(t, -1*time.Second, monitor.config.DialTimeout)
	})

	t.Run("Monitor with very long username", func(t *testing.T) {
		longUsername := strings.Repeat("a", 10000)
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    longUsername,
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.Len(t, monitor.config.Username, 10000)
	})

	t.Run("Monitor with very long password", func(t *testing.T) {
		longPassword := strings.Repeat("b", 10000)
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    longPassword,
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.Len(t, monitor.config.Password, 10000)
	})

	t.Run("Monitor with empty endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.Len(t, monitor.config.Endpoints, 0)
	})

	t.Run("Monitor with nil endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.Nil(t, monitor.config.Endpoints)
	})

	t.Run("Monitor with empty username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.Empty(t, monitor.config.Username)
	})

	t.Run("Monitor with empty password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.Empty(t, monitor.config.Password)
	})

	t.Run("Monitor with special characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "https://etcd-1:2379", "http://etcd-2:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.Len(t, monitor.config.Endpoints, 3)
		assert.Equal(t, "http://etcd-0:2379", monitor.config.Endpoints[0])
		assert.Equal(t, "https://etcd-1:2379", monitor.config.Endpoints[1])
		assert.Equal(t, "http://etcd-2:2379", monitor.config.Endpoints[2])
	})

	t.Run("Monitor with special characters in username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "user@domain.com",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.Equal(t, "user@domain.com", monitor.config.Username)
	})

	t.Run("Monitor with special characters in password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "pass!@#$%^&*()",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.Equal(t, "pass!@#$%^&*()", monitor.config.Password)
	})

	t.Run("Monitor with unicode characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.Len(t, monitor.config.Endpoints, 2)
	})

	t.Run("Monitor with unicode characters in username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "用户",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.Equal(t, "用户", monitor.config.Username)
	})

	t.Run("Monitor with unicode characters in password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "密码",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.Equal(t, "密码", monitor.config.Password)
	})

	t.Run("Monitor with control characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.Len(t, monitor.config.Endpoints, 2)
	})

	t.Run("Monitor with control characters in username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "user\n\r\t",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.Equal(t, "user\n\r\t", monitor.config.Username)
	})

	t.Run("Monitor with control characters in password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "pass\n\r\t",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.Equal(t, "pass\n\r\t", monitor.config.Password)
	})
}

func TestEtcdMonitor_Performance(t *testing.T) {
	t.Run("Monitor creation performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Create 1000 Monitor instances
		for i := 0; i < 1000; i++ {
			monitor := NewEtcdMonitor(config, logger)
			assert.NotNil(t, monitor)
		}
		
		duration := time.Since(start)
		
		// Should create 1000 instances in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Creating 1000 Monitor instances took too long: %v", duration)
	})

	t.Run("Monitor with large config performance", func(t *testing.T) {
		endpoints := make([]string, 10000)
		for i := 0; i < 10000; i++ {
			endpoints[i] = fmt.Sprintf("localhost:%d", 2379+i)
		}
		
		config := &Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
			Username:    strings.Repeat("a", 10000),
			Password:    strings.Repeat("b", 10000),
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		monitor := NewEtcdMonitor(config, logger)
		
		duration := time.Since(start)
		
		// Should create monitor with large config in less than 10ms
		assert.True(t, duration < 10*time.Millisecond, "Creating Monitor with large config took too long: %v", duration)
		assert.Len(t, monitor.config.Endpoints, 10000)
		assert.Len(t, monitor.config.Username, 10000)
		assert.Len(t, monitor.config.Password, 10000)
	})

	t.Run("Monitor start/stop performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Start and stop 100 monitors
		for i := 0; i < 100; i++ {
			monitor := NewEtcdMonitor(config, logger)
			err := monitor.Start()
			require.NoError(t, err)
			monitor.Stop()
		}
		
		duration := time.Since(start)
		
		// Should start and stop 100 monitors in less than 1 second
		assert.True(t, duration < 1*time.Second, "Starting and stopping 100 monitors took too long: %v", duration)
	})

	t.Run("Monitor run performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		monitor := NewEtcdMonitor(config, logger)
		ctx := context.Background()
		
		start := time.Now()
		
		// Run monitor 10 times
		for i := 0; i < 10; i++ {
			err := monitor.Run(ctx)
			assert.NoError(t, err)
		}
		
		duration := time.Since(start)
		
		// Should run 10 monitors in less than 1 second
		assert.True(t, duration < 1*time.Second, "Running 10 monitors took too long: %v", duration)
	})

	t.Run("Monitor get results performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		monitor := NewEtcdMonitor(config, logger)
		
		start := time.Now()
		
		// Get results 1000 times
		for i := 0; i < 1000; i++ {
			results := monitor.GetMonitorResults()
			assert.NotNil(t, results)
		}
		
		duration := time.Since(start)
		
		// Should get results 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting results 1000 times took too long: %v", duration)
	})

	t.Run("Monitor get status performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		monitor := NewEtcdMonitor(config, logger)
		
		start := time.Now()
		
		// Get status 1000 times
		for i := 0; i < 1000; i++ {
			status := monitor.GetMonitorStatus()
			assert.NotNil(t, status)
		}
		
		duration := time.Since(start)
		
		// Should get status 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting status 1000 times took too long: %v", duration)
	})

	t.Run("Monitor get metrics performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		monitor := NewEtcdMonitor(config, logger)
		
		start := time.Now()
		
		// Get metrics 1000 times
		for i := 0; i < 1000; i++ {
			metrics := monitor.GetMonitorMetrics()
			assert.NotNil(t, metrics)
		}
		
		duration := time.Since(start)
		
		// Should get metrics 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting metrics 1000 times took too long: %v", duration)
	})
}

func TestEtcdMonitor_Integration(t *testing.T) {
	t.Run("Monitor with all fields set", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379", "http://etcd-2:2379"},
			DialTimeout: 10 * time.Second,
			Username:    "admin",
			Password:    "secret",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.NotNil(t, monitor.config)
		assert.NotNil(t, monitor.logger)
		assert.Len(t, monitor.config.Endpoints, 3)
		assert.Equal(t, "http://etcd-0:2379", monitor.config.Endpoints[0])
		assert.Equal(t, "http://etcd-1:2379", monitor.config.Endpoints[1])
		assert.Equal(t, "http://etcd-2:2379", monitor.config.Endpoints[2])
		assert.Equal(t, 10*time.Second, monitor.config.DialTimeout)
		assert.Equal(t, "admin", monitor.config.Username)
		assert.Equal(t, "secret", monitor.config.Password)
	})

	t.Run("Monitor with minimal fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.NotNil(t, monitor.config)
		assert.NotNil(t, monitor.logger)
		assert.Len(t, monitor.config.Endpoints, 1)
		assert.Equal(t, "localhost:2379", monitor.config.Endpoints[0])
		assert.Equal(t, 5*time.Second, monitor.config.DialTimeout)
		assert.Equal(t, "testuser", monitor.config.Username)
		assert.Equal(t, "testpass", monitor.config.Password)
	})

	t.Run("Monitor with empty fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 0,
			Username:    "",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.NotNil(t, monitor.config)
		assert.NotNil(t, monitor.logger)
		assert.Len(t, monitor.config.Endpoints, 0)
		assert.Equal(t, time.Duration(0), monitor.config.DialTimeout)
		assert.Empty(t, monitor.config.Username)
		assert.Empty(t, monitor.config.Password)
	})

	t.Run("Monitor with nil fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 0,
			Username:    "",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)
		assert.NotNil(t, monitor.config)
		assert.NotNil(t, monitor.logger)
		assert.Nil(t, monitor.config.Endpoints)
		assert.Equal(t, time.Duration(0), monitor.config.DialTimeout)
		assert.Empty(t, monitor.config.Username)
		assert.Empty(t, monitor.config.Password)
	})

	t.Run("Monitor full lifecycle", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)

		// Start monitor
		err := monitor.Start()
		assert.NoError(t, err)

		// Run monitor
		ctx := context.Background()
		err = monitor.Run(ctx)
		assert.NoError(t, err)

		// Get results
		results := monitor.GetMonitorResults()
		assert.NotNil(t, results)

		// Get status
		status := monitor.GetMonitorStatus()
		assert.NotNil(t, status)

		// Get metrics
		metrics := monitor.GetMonitorMetrics()
		assert.NotNil(t, metrics)

		// Stop monitor
		monitor.Stop()
	})

	t.Run("Monitor with multiple concurrent operations", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)

		// Start monitor
		err := monitor.Start()
		require.NoError(t, err)

		// Perform multiple concurrent operations
		done := make(chan bool, 10)
		
		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				// Run monitor
				ctx := context.Background()
				err := monitor.Run(ctx)
				assert.NoError(t, err)
				
				// Get results
				results := monitor.GetMonitorResults()
				assert.NotNil(t, results)
				
				// Get status
				status := monitor.GetMonitorStatus()
				assert.NotNil(t, status)
				
				// Get metrics
				metrics := monitor.GetMonitorMetrics()
				assert.NotNil(t, metrics)
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Stop monitor
		monitor.Stop()
	})

	t.Run("Monitor with different configurations", func(t *testing.T) {
		configs := []*Config{
			{
				Endpoints:   []string{"localhost:2379"},
				DialTimeout: 5 * time.Second,
				Username:    "testuser",
				Password:    "testpass",
			},
			{
				Endpoints:   []string{"localhost:2379", "localhost:2380"},
				DialTimeout: 10 * time.Second,
				Username:    "admin",
				Password:    "secret",
			},
			{
				Endpoints:   []string{"localhost:2379", "localhost:2380", "localhost:2381"},
				DialTimeout: 15 * time.Second,
				Username:    "root",
				Password:    "password",
			},
		}
		logger, _ := zap.NewDevelopment()

		for i, config := range configs {
			monitor := NewEtcdMonitor(config, logger)
			assert.NotNil(t, monitor)

			// Start monitor
			err := monitor.Start()
			require.NoError(t, err)

			// Run monitor
			ctx := context.Background()
			err = monitor.Run(ctx)
			assert.NoError(t, err)

			// Get results
			results := monitor.GetMonitorResults()
			assert.NotNil(t, results)

			// Get status
			status := monitor.GetMonitorStatus()
			assert.NotNil(t, status)

			// Get metrics
			metrics := monitor.GetMonitorMetrics()
			assert.NotNil(t, metrics)

			// Stop monitor
			monitor.Stop()
		}
	})

	t.Run("Monitor with context cancellation", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)

		// Start monitor
		err := monitor.Start()
		require.NoError(t, err)

		// Run monitor with context that gets cancelled
		ctx, cancel := context.WithCancel(context.Background())
		
		// Cancel context after a short delay
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()
		
		err = monitor.Run(ctx)
		assert.Error(t, err) // Should error due to context cancellation

		// Stop monitor
		monitor.Stop()
	})

	t.Run("Monitor with context timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		monitor := NewEtcdMonitor(config, logger)
		assert.NotNil(t, monitor)

		// Start monitor
		err := monitor.Start()
		require.NoError(t, err)

		// Run monitor with context that times out
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		err = monitor.Run(ctx)
		assert.Error(t, err) // Should error due to context timeout

		// Stop monitor
		monitor.Stop()
	})
}
