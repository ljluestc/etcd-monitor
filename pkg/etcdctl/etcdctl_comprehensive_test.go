package etcdctl

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestEtcdCtl_Comprehensive(t *testing.T) {
	t.Run("NewEtcdCtl with valid config", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.NotNil(t, etcdctl.config)
		assert.NotNil(t, etcdctl.logger)
	})

	t.Run("NewEtcdCtl with nil config", func(t *testing.T) {
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(nil, logger)
		assert.NotNil(t, etcdctl)
		assert.Nil(t, etcdctl.config)
		assert.NotNil(t, etcdctl.logger)
	})

	t.Run("NewEtcdCtl with nil logger", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}

		etcdctl := NewEtcdCtl(config, nil)
		assert.NotNil(t, etcdctl)
		assert.NotNil(t, etcdctl.config)
		assert.NotNil(t, etcdctl.logger) // Should create a production logger
	})

	t.Run("NewEtcdCtl with empty config", func(t *testing.T) {
		config := &Config{}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.NotNil(t, etcdctl.config)
		assert.NotNil(t, etcdctl.logger)
	})

	t.Run("Start etcdctl", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		err := etcdctl.Start()
		assert.NoError(t, err)
	})

	t.Run("Stop etcdctl", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		err := etcdctl.Start()
		require.NoError(t, err)

		etcdctl.Stop()
		// Should not panic or error
	})

	t.Run("Stop etcdctl without starting", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		etcdctl.Stop()
		// Should not panic or error
	})

	t.Run("ExecuteCommand", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		ctx := context.Background()
		
		result, err := etcdctl.ExecuteCommand(ctx, "version")
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("ExecuteCommand with nil context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		
		result, err := etcdctl.ExecuteCommand(nil, "version")
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("ExecuteCommand with cancelled context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		result, err := etcdctl.ExecuteCommand(ctx, "version")
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("ExecuteCommand with timeout context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		result, err := etcdctl.ExecuteCommand(ctx, "version")
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("ExecuteCommand with empty command", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		ctx := context.Background()
		
		result, err := etcdctl.ExecuteCommand(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("ExecuteCommand with nil command", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		ctx := context.Background()
		
		result, err := etcdctl.ExecuteCommand(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("ExecuteCommand with invalid command", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		ctx := context.Background()
		
		result, err := etcdctl.ExecuteCommand(ctx, "invalid-command")
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("ExecuteCommand with special characters", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		ctx := context.Background()
		
		result, err := etcdctl.ExecuteCommand(ctx, "version --help")
		assert.Error(t, err) // Should error due to invalid command
		assert.Nil(t, result)
	})

	t.Run("ExecuteCommand with unicode characters", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		ctx := context.Background()
		
		result, err := etcdctl.ExecuteCommand(ctx, "version 测试")
		assert.Error(t, err) // Should error due to invalid command
		assert.Nil(t, result)
	})

	t.Run("ExecuteCommand with control characters", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		ctx := context.Background()
		
		result, err := etcdctl.ExecuteCommand(ctx, "version\n\r\t")
		assert.Error(t, err) // Should error due to invalid command
		assert.Nil(t, result)
	})

	t.Run("GetCommandResult", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		result := etcdctl.GetCommandResult()
		assert.NotNil(t, result)
	})

	t.Run("GetCommandResult after execution", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		ctx := context.Background()
		
		_, err := etcdctl.ExecuteCommand(ctx, "version")
		require.NoError(t, err)
		
		result := etcdctl.GetCommandResult()
		assert.NotNil(t, result)
	})

	t.Run("GetCommandStatus", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		status := etcdctl.GetCommandStatus()
		assert.NotNil(t, status)
	})

	t.Run("GetCommandStatus after execution", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		ctx := context.Background()
		
		_, err := etcdctl.ExecuteCommand(ctx, "version")
		require.NoError(t, err)
		
		status := etcdctl.GetCommandStatus()
		assert.NotNil(t, status)
	})

	t.Run("GetCommandMetrics", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		metrics := etcdctl.GetCommandMetrics()
		assert.NotNil(t, metrics)
	})

	t.Run("GetCommandMetrics after execution", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		ctx := context.Background()
		
		_, err := etcdctl.ExecuteCommand(ctx, "version")
		require.NoError(t, err)
		
		metrics := etcdctl.GetCommandMetrics()
		assert.NotNil(t, metrics)
	})
}

func TestEtcdCtl_EdgeCases(t *testing.T) {
	t.Run("EtcdCtl with very long endpoint list", func(t *testing.T) {
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

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.Len(t, etcdctl.config.Endpoints, 1000)
	})

	t.Run("EtcdCtl with very long timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 24 * time.Hour,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.Equal(t, 24*time.Hour, etcdctl.config.DialTimeout)
	})

	t.Run("EtcdCtl with very short timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 1 * time.Nanosecond,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.Equal(t, 1*time.Nanosecond, etcdctl.config.DialTimeout)
	})

	t.Run("EtcdCtl with zero timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 0,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.Equal(t, time.Duration(0), etcdctl.config.DialTimeout)
	})

	t.Run("EtcdCtl with negative timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: -1 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.Equal(t, -1*time.Second, etcdctl.config.DialTimeout)
	})

	t.Run("EtcdCtl with very long username", func(t *testing.T) {
		longUsername := strings.Repeat("a", 10000)
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    longUsername,
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.Len(t, etcdctl.config.Username, 10000)
	})

	t.Run("EtcdCtl with very long password", func(t *testing.T) {
		longPassword := strings.Repeat("b", 10000)
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    longPassword,
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.Len(t, etcdctl.config.Password, 10000)
	})

	t.Run("EtcdCtl with empty endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.Len(t, etcdctl.config.Endpoints, 0)
	})

	t.Run("EtcdCtl with nil endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.Nil(t, etcdctl.config.Endpoints)
	})

	t.Run("EtcdCtl with empty username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.Empty(t, etcdctl.config.Username)
	})

	t.Run("EtcdCtl with empty password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.Empty(t, etcdctl.config.Password)
	})

	t.Run("EtcdCtl with special characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "https://etcd-1:2379", "http://etcd-2:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.Len(t, etcdctl.config.Endpoints, 3)
		assert.Equal(t, "http://etcd-0:2379", etcdctl.config.Endpoints[0])
		assert.Equal(t, "https://etcd-1:2379", etcdctl.config.Endpoints[1])
		assert.Equal(t, "http://etcd-2:2379", etcdctl.config.Endpoints[2])
	})

	t.Run("EtcdCtl with special characters in username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "user@domain.com",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.Equal(t, "user@domain.com", etcdctl.config.Username)
	})

	t.Run("EtcdCtl with special characters in password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "pass!@#$%^&*()",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.Equal(t, "pass!@#$%^&*()", etcdctl.config.Password)
	})

	t.Run("EtcdCtl with unicode characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.Len(t, etcdctl.config.Endpoints, 2)
	})

	t.Run("EtcdCtl with unicode characters in username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "用户",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.Equal(t, "用户", etcdctl.config.Username)
	})

	t.Run("EtcdCtl with unicode characters in password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "密码",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.Equal(t, "密码", etcdctl.config.Password)
	})

	t.Run("EtcdCtl with control characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.Len(t, etcdctl.config.Endpoints, 2)
	})

	t.Run("EtcdCtl with control characters in username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "user\n\r\t",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.Equal(t, "user\n\r\t", etcdctl.config.Username)
	})

	t.Run("EtcdCtl with control characters in password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "pass\n\r\t",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.Equal(t, "pass\n\r\t", etcdctl.config.Password)
	})

	t.Run("EtcdCtl with very long command", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		longCommand := strings.Repeat("version ", 1000)
		
		ctx := context.Background()
		result, err := etcdctl.ExecuteCommand(ctx, longCommand)
		assert.Error(t, err) // Should error due to invalid command
		assert.Nil(t, result)
	})

	t.Run("EtcdCtl with special characters in command", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		specialCommand := "version --help@domain.com"
		
		ctx := context.Background()
		result, err := etcdctl.ExecuteCommand(ctx, specialCommand)
		assert.Error(t, err) // Should error due to invalid command
		assert.Nil(t, result)
	})

	t.Run("EtcdCtl with unicode characters in command", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		unicodeCommand := "version 测试"
		
		ctx := context.Background()
		result, err := etcdctl.ExecuteCommand(ctx, unicodeCommand)
		assert.Error(t, err) // Should error due to invalid command
		assert.Nil(t, result)
	})

	t.Run("EtcdCtl with control characters in command", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		controlCommand := "version\n\r\t"
		
		ctx := context.Background()
		result, err := etcdctl.ExecuteCommand(ctx, controlCommand)
		assert.Error(t, err) // Should error due to invalid command
		assert.Nil(t, result)
	})
}

func TestEtcdCtl_Performance(t *testing.T) {
	t.Run("EtcdCtl creation performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Create 1000 EtcdCtl instances
		for i := 0; i < 1000; i++ {
			etcdctl := NewEtcdCtl(config, logger)
			assert.NotNil(t, etcdctl)
		}
		
		duration := time.Since(start)
		
		// Should create 1000 instances in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Creating 1000 EtcdCtl instances took too long: %v", duration)
	})

	t.Run("EtcdCtl with large config performance", func(t *testing.T) {
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
		
		etcdctl := NewEtcdCtl(config, logger)
		
		duration := time.Since(start)
		
		// Should create etcdctl with large config in less than 10ms
		assert.True(t, duration < 10*time.Millisecond, "Creating EtcdCtl with large config took too long: %v", duration)
		assert.Len(t, etcdctl.config.Endpoints, 10000)
		assert.Len(t, etcdctl.config.Username, 10000)
		assert.Len(t, etcdctl.config.Password, 10000)
	})

	t.Run("EtcdCtl start/stop performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Start and stop 100 etcdctl instances
		for i := 0; i < 100; i++ {
			etcdctl := NewEtcdCtl(config, logger)
			err := etcdctl.Start()
			require.NoError(t, err)
			etcdctl.Stop()
		}
		
		duration := time.Since(start)
		
		// Should start and stop 100 instances in less than 1 second
		assert.True(t, duration < 1*time.Second, "Starting and stopping 100 etcdctl instances took too long: %v", duration)
	})

	t.Run("EtcdCtl execute command performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		etcdctl := NewEtcdCtl(config, logger)
		ctx := context.Background()
		
		start := time.Now()
		
		// Execute command 100 times
		for i := 0; i < 100; i++ {
			_, err := etcdctl.ExecuteCommand(ctx, "version")
			assert.NoError(t, err)
		}
		
		duration := time.Since(start)
		
		// Should execute command 100 times in less than 2 seconds
		assert.True(t, duration < 2*time.Second, "Executing command 100 times took too long: %v", duration)
	})

	t.Run("EtcdCtl get command result performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		etcdctl := NewEtcdCtl(config, logger)
		
		start := time.Now()
		
		// Get command result 1000 times
		for i := 0; i < 1000; i++ {
			result := etcdctl.GetCommandResult()
			assert.NotNil(t, result)
		}
		
		duration := time.Since(start)
		
		// Should get command result 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting command result 1000 times took too long: %v", duration)
	})

	t.Run("EtcdCtl get command status performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		etcdctl := NewEtcdCtl(config, logger)
		
		start := time.Now()
		
		// Get command status 1000 times
		for i := 0; i < 1000; i++ {
			status := etcdctl.GetCommandStatus()
			assert.NotNil(t, status)
		}
		
		duration := time.Since(start)
		
		// Should get command status 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting command status 1000 times took too long: %v", duration)
	})

	t.Run("EtcdCtl get command metrics performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		etcdctl := NewEtcdCtl(config, logger)
		
		start := time.Now()
		
		// Get command metrics 1000 times
		for i := 0; i < 1000; i++ {
			metrics := etcdctl.GetCommandMetrics()
			assert.NotNil(t, metrics)
		}
		
		duration := time.Since(start)
		
		// Should get command metrics 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting command metrics 1000 times took too long: %v", duration)
	})
}

func TestEtcdCtl_Integration(t *testing.T) {
	t.Run("EtcdCtl with all fields set", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379", "http://etcd-2:2379"},
			DialTimeout: 10 * time.Second,
			Username:    "admin",
			Password:    "secret",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.NotNil(t, etcdctl.config)
		assert.NotNil(t, etcdctl.logger)
		assert.Len(t, etcdctl.config.Endpoints, 3)
		assert.Equal(t, "http://etcd-0:2379", etcdctl.config.Endpoints[0])
		assert.Equal(t, "http://etcd-1:2379", etcdctl.config.Endpoints[1])
		assert.Equal(t, "http://etcd-2:2379", etcdctl.config.Endpoints[2])
		assert.Equal(t, 10*time.Second, etcdctl.config.DialTimeout)
		assert.Equal(t, "admin", etcdctl.config.Username)
		assert.Equal(t, "secret", etcdctl.config.Password)
	})

	t.Run("EtcdCtl with minimal fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.NotNil(t, etcdctl.config)
		assert.NotNil(t, etcdctl.logger)
		assert.Len(t, etcdctl.config.Endpoints, 1)
		assert.Equal(t, "localhost:2379", etcdctl.config.Endpoints[0])
		assert.Equal(t, 5*time.Second, etcdctl.config.DialTimeout)
		assert.Equal(t, "testuser", etcdctl.config.Username)
		assert.Equal(t, "testpass", etcdctl.config.Password)
	})

	t.Run("EtcdCtl with empty fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 0,
			Username:    "",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.NotNil(t, etcdctl.config)
		assert.NotNil(t, etcdctl.logger)
		assert.Len(t, etcdctl.config.Endpoints, 0)
		assert.Equal(t, time.Duration(0), etcdctl.config.DialTimeout)
		assert.Empty(t, etcdctl.config.Username)
		assert.Empty(t, etcdctl.config.Password)
	})

	t.Run("EtcdCtl with nil fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 0,
			Username:    "",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)
		assert.NotNil(t, etcdctl.config)
		assert.NotNil(t, etcdctl.logger)
		assert.Nil(t, etcdctl.config.Endpoints)
		assert.Equal(t, time.Duration(0), etcdctl.config.DialTimeout)
		assert.Empty(t, etcdctl.config.Username)
		assert.Empty(t, etcdctl.config.Password)
	})

	t.Run("EtcdCtl full lifecycle", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)

		// Start etcdctl
		err := etcdctl.Start()
		assert.NoError(t, err)

		// Execute command
		ctx := context.Background()
		result, err := etcdctl.ExecuteCommand(ctx, "version")
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// Get command result
		commandResult := etcdctl.GetCommandResult()
		assert.NotNil(t, commandResult)

		// Get command status
		status := etcdctl.GetCommandStatus()
		assert.NotNil(t, status)

		// Get command metrics
		metrics := etcdctl.GetCommandMetrics()
		assert.NotNil(t, metrics)

		// Stop etcdctl
		etcdctl.Stop()
	})

	t.Run("EtcdCtl with multiple concurrent operations", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)

		// Start etcdctl
		err := etcdctl.Start()
		require.NoError(t, err)

		// Perform multiple concurrent operations
		done := make(chan bool, 10)
		
		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				// Execute command
				ctx := context.Background()
				result, err := etcdctl.ExecuteCommand(ctx, "version")
				assert.NoError(t, err)
				assert.NotNil(t, result)
				
				// Get command result
				commandResult := etcdctl.GetCommandResult()
				assert.NotNil(t, commandResult)
				
				// Get command status
				status := etcdctl.GetCommandStatus()
				assert.NotNil(t, status)
				
				// Get command metrics
				metrics := etcdctl.GetCommandMetrics()
				assert.NotNil(t, metrics)
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Stop etcdctl
		etcdctl.Stop()
	})

	t.Run("EtcdCtl with different commands", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)

		// Start etcdctl
		err := etcdctl.Start()
		require.NoError(t, err)

		// Test with different commands
		commands := []string{
			"version",
			"endpoint health",
			"endpoint status",
			"member list",
			"snapshot status",
		}

		for _, command := range commands {
			// Execute command
			ctx := context.Background()
			result, err := etcdctl.ExecuteCommand(ctx, command)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			
			// Get command result
			commandResult := etcdctl.GetCommandResult()
			assert.NotNil(t, commandResult)
			
			// Get command status
			status := etcdctl.GetCommandStatus()
			assert.NotNil(t, status)
			
			// Get command metrics
			metrics := etcdctl.GetCommandMetrics()
			assert.NotNil(t, metrics)
		}

		// Stop etcdctl
		etcdctl.Stop()
	})

	t.Run("EtcdCtl with context cancellation", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)

		// Start etcdctl
		err := etcdctl.Start()
		require.NoError(t, err)

		// Execute command with context that gets cancelled
		ctx, cancel := context.WithCancel(context.Background())
		
		// Cancel context after a short delay
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()
		
		result, err := etcdctl.ExecuteCommand(ctx, "version")
		assert.Error(t, err) // Should error due to context cancellation
		assert.Nil(t, result)

		// Stop etcdctl
		etcdctl.Stop()
	})

	t.Run("EtcdCtl with context timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		etcdctl := NewEtcdCtl(config, logger)
		assert.NotNil(t, etcdctl)

		// Start etcdctl
		err := etcdctl.Start()
		require.NoError(t, err)

		// Execute command with context that times out
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		result, err := etcdctl.ExecuteCommand(ctx, "version")
		assert.Error(t, err) // Should error due to context timeout
		assert.Nil(t, result)

		// Stop etcdctl
		etcdctl.Stop()
	})
}
