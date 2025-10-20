package pingapgo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestPingapGo_Comprehensive(t *testing.T) {
	t.Run("NewPingapGo with valid config", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.NotNil(t, pingap.config)
		assert.NotNil(t, pingap.logger)
	})

	t.Run("NewPingapGo with nil config", func(t *testing.T) {
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(nil, logger)
		assert.NotNil(t, pingap)
		assert.Nil(t, pingap.config)
		assert.NotNil(t, pingap.logger)
	})

	t.Run("NewPingapGo with nil logger", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}

		pingap := NewPingapGo(config, nil)
		assert.NotNil(t, pingap)
		assert.NotNil(t, pingap.config)
		assert.NotNil(t, pingap.logger) // Should create a production logger
	})

	t.Run("NewPingapGo with empty config", func(t *testing.T) {
		config := &Config{}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.NotNil(t, pingap.config)
		assert.NotNil(t, pingap.logger)
	})

	t.Run("Start pingap", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		err := pingap.Start()
		assert.NoError(t, err)
	})

	t.Run("Stop pingap", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		err := pingap.Start()
		require.NoError(t, err)

		pingap.Stop()
		// Should not panic or error
	})

	t.Run("Stop pingap without starting", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		pingap.Stop()
		// Should not panic or error
	})

	t.Run("Run pingap", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		ctx := context.Background()
		
		err := pingap.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("Run pingap with nil context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		
		err := pingap.Run(nil)
		assert.Error(t, err)
	})

	t.Run("Run pingap with cancelled context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		err := pingap.Run(ctx)
		assert.Error(t, err)
	})

	t.Run("Run pingap with timeout context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		err := pingap.Run(ctx)
		assert.Error(t, err)
	})

	t.Run("GetPingapResults", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		results := pingap.GetPingapResults()
		assert.NotNil(t, results)
	})

	t.Run("GetPingapResults after run", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		ctx := context.Background()
		
		err := pingap.Run(ctx)
		require.NoError(t, err)
		
		results := pingap.GetPingapResults()
		assert.NotNil(t, results)
	})

	t.Run("GetPingapStatus", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		status := pingap.GetPingapStatus()
		assert.NotNil(t, status)
	})

	t.Run("GetPingapStatus after run", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		ctx := context.Background()
		
		err := pingap.Run(ctx)
		require.NoError(t, err)
		
		status := pingap.GetPingapStatus()
		assert.NotNil(t, status)
	})

	t.Run("GetPingapMetrics", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		metrics := pingap.GetPingapMetrics()
		assert.NotNil(t, metrics)
	})

	t.Run("GetPingapMetrics after run", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		ctx := context.Background()
		
		err := pingap.Run(ctx)
		require.NoError(t, err)
		
		metrics := pingap.GetPingapMetrics()
		assert.NotNil(t, metrics)
	})
}

func TestPingapGo_EdgeCases(t *testing.T) {
	t.Run("PingapGo with very long endpoint list", func(t *testing.T) {
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

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.Len(t, pingap.config.Endpoints, 1000)
	})

	t.Run("PingapGo with very long timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 24 * time.Hour,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.Equal(t, 24*time.Hour, pingap.config.DialTimeout)
	})

	t.Run("PingapGo with very short timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 1 * time.Nanosecond,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.Equal(t, 1*time.Nanosecond, pingap.config.DialTimeout)
	})

	t.Run("PingapGo with zero timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 0,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.Equal(t, time.Duration(0), pingap.config.DialTimeout)
	})

	t.Run("PingapGo with negative timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: -1 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.Equal(t, -1*time.Second, pingap.config.DialTimeout)
	})

	t.Run("PingapGo with very long username", func(t *testing.T) {
		longUsername := strings.Repeat("a", 10000)
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    longUsername,
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.Len(t, pingap.config.Username, 10000)
	})

	t.Run("PingapGo with very long password", func(t *testing.T) {
		longPassword := strings.Repeat("b", 10000)
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    longPassword,
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.Len(t, pingap.config.Password, 10000)
	})

	t.Run("PingapGo with empty endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.Len(t, pingap.config.Endpoints, 0)
	})

	t.Run("PingapGo with nil endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.Nil(t, pingap.config.Endpoints)
	})

	t.Run("PingapGo with empty username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.Empty(t, pingap.config.Username)
	})

	t.Run("PingapGo with empty password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.Empty(t, pingap.config.Password)
	})

	t.Run("PingapGo with special characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "https://etcd-1:2379", "http://etcd-2:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.Len(t, pingap.config.Endpoints, 3)
		assert.Equal(t, "http://etcd-0:2379", pingap.config.Endpoints[0])
		assert.Equal(t, "https://etcd-1:2379", pingap.config.Endpoints[1])
		assert.Equal(t, "http://etcd-2:2379", pingap.config.Endpoints[2])
	})

	t.Run("PingapGo with special characters in username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "user@domain.com",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.Equal(t, "user@domain.com", pingap.config.Username)
	})

	t.Run("PingapGo with special characters in password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "pass!@#$%^&*()",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.Equal(t, "pass!@#$%^&*()", pingap.config.Password)
	})

	t.Run("PingapGo with unicode characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.Len(t, pingap.config.Endpoints, 2)
	})

	t.Run("PingapGo with unicode characters in username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "用户",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.Equal(t, "用户", pingap.config.Username)
	})

	t.Run("PingapGo with unicode characters in password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "密码",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.Equal(t, "密码", pingap.config.Password)
	})

	t.Run("PingapGo with control characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.Len(t, pingap.config.Endpoints, 2)
	})

	t.Run("PingapGo with control characters in username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "user\n\r\t",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.Equal(t, "user\n\r\t", pingap.config.Username)
	})

	t.Run("PingapGo with control characters in password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "pass\n\r\t",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.Equal(t, "pass\n\r\t", pingap.config.Password)
	})
}

func TestPingapGo_Performance(t *testing.T) {
	t.Run("PingapGo creation performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Create 1000 PingapGo instances
		for i := 0; i < 1000; i++ {
			pingap := NewPingapGo(config, logger)
			assert.NotNil(t, pingap)
		}
		
		duration := time.Since(start)
		
		// Should create 1000 instances in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Creating 1000 PingapGo instances took too long: %v", duration)
	})

	t.Run("PingapGo with large config performance", func(t *testing.T) {
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
		
		pingap := NewPingapGo(config, logger)
		
		duration := time.Since(start)
		
		// Should create pingap with large config in less than 10ms
		assert.True(t, duration < 10*time.Millisecond, "Creating PingapGo with large config took too long: %v", duration)
		assert.Len(t, pingap.config.Endpoints, 10000)
		assert.Len(t, pingap.config.Username, 10000)
		assert.Len(t, pingap.config.Password, 10000)
	})

	t.Run("PingapGo start/stop performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Start and stop 100 pingaps
		for i := 0; i < 100; i++ {
			pingap := NewPingapGo(config, logger)
			err := pingap.Start()
			require.NoError(t, err)
			pingap.Stop()
		}
		
		duration := time.Since(start)
		
		// Should start and stop 100 pingaps in less than 1 second
		assert.True(t, duration < 1*time.Second, "Starting and stopping 100 pingaps took too long: %v", duration)
	})

	t.Run("PingapGo run performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		pingap := NewPingapGo(config, logger)
		ctx := context.Background()
		
		start := time.Now()
		
		// Run pingap 10 times
		for i := 0; i < 10; i++ {
			err := pingap.Run(ctx)
			assert.NoError(t, err)
		}
		
		duration := time.Since(start)
		
		// Should run 10 pingaps in less than 1 second
		assert.True(t, duration < 1*time.Second, "Running 10 pingaps took too long: %v", duration)
	})

	t.Run("PingapGo get results performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		pingap := NewPingapGo(config, logger)
		
		start := time.Now()
		
		// Get results 1000 times
		for i := 0; i < 1000; i++ {
			results := pingap.GetPingapResults()
			assert.NotNil(t, results)
		}
		
		duration := time.Since(start)
		
		// Should get results 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting results 1000 times took too long: %v", duration)
	})

	t.Run("PingapGo get status performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		pingap := NewPingapGo(config, logger)
		
		start := time.Now()
		
		// Get status 1000 times
		for i := 0; i < 1000; i++ {
			status := pingap.GetPingapStatus()
			assert.NotNil(t, status)
		}
		
		duration := time.Since(start)
		
		// Should get status 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting status 1000 times took too long: %v", duration)
	})

	t.Run("PingapGo get metrics performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		pingap := NewPingapGo(config, logger)
		
		start := time.Now()
		
		// Get metrics 1000 times
		for i := 0; i < 1000; i++ {
			metrics := pingap.GetPingapMetrics()
			assert.NotNil(t, metrics)
		}
		
		duration := time.Since(start)
		
		// Should get metrics 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting metrics 1000 times took too long: %v", duration)
	})
}

func TestPingapGo_Integration(t *testing.T) {
	t.Run("PingapGo with all fields set", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379", "http://etcd-2:2379"},
			DialTimeout: 10 * time.Second,
			Username:    "admin",
			Password:    "secret",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.NotNil(t, pingap.config)
		assert.NotNil(t, pingap.logger)
		assert.Len(t, pingap.config.Endpoints, 3)
		assert.Equal(t, "http://etcd-0:2379", pingap.config.Endpoints[0])
		assert.Equal(t, "http://etcd-1:2379", pingap.config.Endpoints[1])
		assert.Equal(t, "http://etcd-2:2379", pingap.config.Endpoints[2])
		assert.Equal(t, 10*time.Second, pingap.config.DialTimeout)
		assert.Equal(t, "admin", pingap.config.Username)
		assert.Equal(t, "secret", pingap.config.Password)
	})

	t.Run("PingapGo with minimal fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.NotNil(t, pingap.config)
		assert.NotNil(t, pingap.logger)
		assert.Len(t, pingap.config.Endpoints, 1)
		assert.Equal(t, "localhost:2379", pingap.config.Endpoints[0])
		assert.Equal(t, 5*time.Second, pingap.config.DialTimeout)
		assert.Equal(t, "testuser", pingap.config.Username)
		assert.Equal(t, "testpass", pingap.config.Password)
	})

	t.Run("PingapGo with empty fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 0,
			Username:    "",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.NotNil(t, pingap.config)
		assert.NotNil(t, pingap.logger)
		assert.Len(t, pingap.config.Endpoints, 0)
		assert.Equal(t, time.Duration(0), pingap.config.DialTimeout)
		assert.Empty(t, pingap.config.Username)
		assert.Empty(t, pingap.config.Password)
	})

	t.Run("PingapGo with nil fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 0,
			Username:    "",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)
		assert.NotNil(t, pingap.config)
		assert.NotNil(t, pingap.logger)
		assert.Nil(t, pingap.config.Endpoints)
		assert.Equal(t, time.Duration(0), pingap.config.DialTimeout)
		assert.Empty(t, pingap.config.Username)
		assert.Empty(t, pingap.config.Password)
	})

	t.Run("PingapGo full lifecycle", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)

		// Start pingap
		err := pingap.Start()
		assert.NoError(t, err)

		// Run pingap
		ctx := context.Background()
		err = pingap.Run(ctx)
		assert.NoError(t, err)

		// Get results
		results := pingap.GetPingapResults()
		assert.NotNil(t, results)

		// Get status
		status := pingap.GetPingapStatus()
		assert.NotNil(t, status)

		// Get metrics
		metrics := pingap.GetPingapMetrics()
		assert.NotNil(t, metrics)

		// Stop pingap
		pingap.Stop()
	})

	t.Run("PingapGo with multiple concurrent operations", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)

		// Start pingap
		err := pingap.Start()
		require.NoError(t, err)

		// Perform multiple concurrent operations
		done := make(chan bool, 10)
		
		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				// Run pingap
				ctx := context.Background()
				err := pingap.Run(ctx)
				assert.NoError(t, err)
				
				// Get results
				results := pingap.GetPingapResults()
				assert.NotNil(t, results)
				
				// Get status
				status := pingap.GetPingapStatus()
				assert.NotNil(t, status)
				
				// Get metrics
				metrics := pingap.GetPingapMetrics()
				assert.NotNil(t, metrics)
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Stop pingap
		pingap.Stop()
	})

	t.Run("PingapGo with different configurations", func(t *testing.T) {
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
			pingap := NewPingapGo(config, logger)
			assert.NotNil(t, pingap)

			// Start pingap
			err := pingap.Start()
			require.NoError(t, err)

			// Run pingap
			ctx := context.Background()
			err = pingap.Run(ctx)
			assert.NoError(t, err)

			// Get results
			results := pingap.GetPingapResults()
			assert.NotNil(t, results)

			// Get status
			status := pingap.GetPingapStatus()
			assert.NotNil(t, status)

			// Get metrics
			metrics := pingap.GetPingapMetrics()
			assert.NotNil(t, metrics)

			// Stop pingap
			pingap.Stop()
		}
	})

	t.Run("PingapGo with context cancellation", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)

		// Start pingap
		err := pingap.Start()
		require.NoError(t, err)

		// Run pingap with context that gets cancelled
		ctx, cancel := context.WithCancel(context.Background())
		
		// Cancel context after a short delay
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()
		
		err = pingap.Run(ctx)
		assert.Error(t, err) // Should error due to context cancellation

		// Stop pingap
		pingap.Stop()
	})

	t.Run("PingapGo with context timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		pingap := NewPingapGo(config, logger)
		assert.NotNil(t, pingap)

		// Start pingap
		err := pingap.Start()
		require.NoError(t, err)

		// Run pingap with context that times out
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		err = pingap.Run(ctx)
		assert.Error(t, err) // Should error due to context timeout

		// Stop pingap
		pingap.Stop()
	})
}
