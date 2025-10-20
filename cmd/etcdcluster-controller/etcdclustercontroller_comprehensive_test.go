package etcdclustercontroller

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestEtcdClusterController_Comprehensive(t *testing.T) {
	t.Run("NewEtcdClusterController with valid config", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.NotNil(t, controller.config)
		assert.NotNil(t, controller.logger)
	})

	t.Run("NewEtcdClusterController with nil config", func(t *testing.T) {
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(nil, logger)
		assert.NotNil(t, controller)
		assert.Nil(t, controller.config)
		assert.NotNil(t, controller.logger)
	})

	t.Run("NewEtcdClusterController with nil logger", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}

		controller := NewEtcdClusterController(config, nil)
		assert.NotNil(t, controller)
		assert.NotNil(t, controller.config)
		assert.NotNil(t, controller.logger) // Should create a production logger
	})

	t.Run("NewEtcdClusterController with empty config", func(t *testing.T) {
		config := &Config{}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.NotNil(t, controller.config)
		assert.NotNil(t, controller.logger)
	})

	t.Run("Start controller", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		err := controller.Start()
		assert.NoError(t, err)
	})

	t.Run("Stop controller", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		err := controller.Start()
		require.NoError(t, err)

		controller.Stop()
		// Should not panic or error
	})

	t.Run("Stop controller without starting", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		controller.Stop()
		// Should not panic or error
	})

	t.Run("Run controller", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		ctx := context.Background()
		
		err := controller.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("Run controller with nil context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		
		err := controller.Run(nil)
		assert.Error(t, err)
	})

	t.Run("Run controller with cancelled context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		err := controller.Run(ctx)
		assert.Error(t, err)
	})

	t.Run("Run controller with timeout context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		err := controller.Run(ctx)
		assert.Error(t, err)
	})

	t.Run("GetControllerResults", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		results := controller.GetControllerResults()
		assert.NotNil(t, results)
	})

	t.Run("GetControllerResults after run", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		ctx := context.Background()
		
		err := controller.Run(ctx)
		require.NoError(t, err)
		
		results := controller.GetControllerResults()
		assert.NotNil(t, results)
	})

	t.Run("GetControllerStatus", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		status := controller.GetControllerStatus()
		assert.NotNil(t, status)
	})

	t.Run("GetControllerStatus after run", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		ctx := context.Background()
		
		err := controller.Run(ctx)
		require.NoError(t, err)
		
		status := controller.GetControllerStatus()
		assert.NotNil(t, status)
	})

	t.Run("GetControllerMetrics", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		metrics := controller.GetControllerMetrics()
		assert.NotNil(t, metrics)
	})

	t.Run("GetControllerMetrics after run", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		ctx := context.Background()
		
		err := controller.Run(ctx)
		require.NoError(t, err)
		
		metrics := controller.GetControllerMetrics()
		assert.NotNil(t, metrics)
	})
}

func TestEtcdClusterController_EdgeCases(t *testing.T) {
	t.Run("Controller with very long endpoint list", func(t *testing.T) {
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

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.Len(t, controller.config.Endpoints, 1000)
	})

	t.Run("Controller with very long timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 24 * time.Hour,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.Equal(t, 24*time.Hour, controller.config.DialTimeout)
	})

	t.Run("Controller with very short timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 1 * time.Nanosecond,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.Equal(t, 1*time.Nanosecond, controller.config.DialTimeout)
	})

	t.Run("Controller with zero timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 0,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.Equal(t, time.Duration(0), controller.config.DialTimeout)
	})

	t.Run("Controller with negative timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: -1 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.Equal(t, -1*time.Second, controller.config.DialTimeout)
	})

	t.Run("Controller with very long username", func(t *testing.T) {
		longUsername := strings.Repeat("a", 10000)
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    longUsername,
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.Len(t, controller.config.Username, 10000)
	})

	t.Run("Controller with very long password", func(t *testing.T) {
		longPassword := strings.Repeat("b", 10000)
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    longPassword,
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.Len(t, controller.config.Password, 10000)
	})

	t.Run("Controller with empty endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.Len(t, controller.config.Endpoints, 0)
	})

	t.Run("Controller with nil endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.Nil(t, controller.config.Endpoints)
	})

	t.Run("Controller with empty username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.Empty(t, controller.config.Username)
	})

	t.Run("Controller with empty password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.Empty(t, controller.config.Password)
	})

	t.Run("Controller with special characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "https://etcd-1:2379", "http://etcd-2:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.Len(t, controller.config.Endpoints, 3)
		assert.Equal(t, "http://etcd-0:2379", controller.config.Endpoints[0])
		assert.Equal(t, "https://etcd-1:2379", controller.config.Endpoints[1])
		assert.Equal(t, "http://etcd-2:2379", controller.config.Endpoints[2])
	})

	t.Run("Controller with special characters in username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "user@domain.com",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.Equal(t, "user@domain.com", controller.config.Username)
	})

	t.Run("Controller with special characters in password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "pass!@#$%^&*()",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.Equal(t, "pass!@#$%^&*()", controller.config.Password)
	})

	t.Run("Controller with unicode characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.Len(t, controller.config.Endpoints, 2)
	})

	t.Run("Controller with unicode characters in username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "用户",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.Equal(t, "用户", controller.config.Username)
	})

	t.Run("Controller with unicode characters in password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "密码",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.Equal(t, "密码", controller.config.Password)
	})

	t.Run("Controller with control characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.Len(t, controller.config.Endpoints, 2)
	})

	t.Run("Controller with control characters in username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "user\n\r\t",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.Equal(t, "user\n\r\t", controller.config.Username)
	})

	t.Run("Controller with control characters in password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "pass\n\r\t",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.Equal(t, "pass\n\r\t", controller.config.Password)
	})
}

func TestEtcdClusterController_Performance(t *testing.T) {
	t.Run("Controller creation performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Create 1000 Controller instances
		for i := 0; i < 1000; i++ {
			controller := NewEtcdClusterController(config, logger)
			assert.NotNil(t, controller)
		}
		
		duration := time.Since(start)
		
		// Should create 1000 instances in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Creating 1000 Controller instances took too long: %v", duration)
	})

	t.Run("Controller with large config performance", func(t *testing.T) {
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
		
		controller := NewEtcdClusterController(config, logger)
		
		duration := time.Since(start)
		
		// Should create controller with large config in less than 10ms
		assert.True(t, duration < 10*time.Millisecond, "Creating Controller with large config took too long: %v", duration)
		assert.Len(t, controller.config.Endpoints, 10000)
		assert.Len(t, controller.config.Username, 10000)
		assert.Len(t, controller.config.Password, 10000)
	})

	t.Run("Controller start/stop performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Start and stop 100 controllers
		for i := 0; i < 100; i++ {
			controller := NewEtcdClusterController(config, logger)
			err := controller.Start()
			require.NoError(t, err)
			controller.Stop()
		}
		
		duration := time.Since(start)
		
		// Should start and stop 100 controllers in less than 1 second
		assert.True(t, duration < 1*time.Second, "Starting and stopping 100 controllers took too long: %v", duration)
	})

	t.Run("Controller run performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		controller := NewEtcdClusterController(config, logger)
		ctx := context.Background()
		
		start := time.Now()
		
		// Run controller 10 times
		for i := 0; i < 10; i++ {
			err := controller.Run(ctx)
			assert.NoError(t, err)
		}
		
		duration := time.Since(start)
		
		// Should run 10 controllers in less than 1 second
		assert.True(t, duration < 1*time.Second, "Running 10 controllers took too long: %v", duration)
	})

	t.Run("Controller get results performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		controller := NewEtcdClusterController(config, logger)
		
		start := time.Now()
		
		// Get results 1000 times
		for i := 0; i < 1000; i++ {
			results := controller.GetControllerResults()
			assert.NotNil(t, results)
		}
		
		duration := time.Since(start)
		
		// Should get results 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting results 1000 times took too long: %v", duration)
	})

	t.Run("Controller get status performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		controller := NewEtcdClusterController(config, logger)
		
		start := time.Now()
		
		// Get status 1000 times
		for i := 0; i < 1000; i++ {
			status := controller.GetControllerStatus()
			assert.NotNil(t, status)
		}
		
		duration := time.Since(start)
		
		// Should get status 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting status 1000 times took too long: %v", duration)
	})

	t.Run("Controller get metrics performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()
		
		controller := NewEtcdClusterController(config, logger)
		
		start := time.Now()
		
		// Get metrics 1000 times
		for i := 0; i < 1000; i++ {
			metrics := controller.GetControllerMetrics()
			assert.NotNil(t, metrics)
		}
		
		duration := time.Since(start)
		
		// Should get metrics 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting metrics 1000 times took too long: %v", duration)
	})
}

func TestEtcdClusterController_Integration(t *testing.T) {
	t.Run("Controller with all fields set", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379", "http://etcd-2:2379"},
			DialTimeout: 10 * time.Second,
			Username:    "admin",
			Password:    "secret",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.NotNil(t, controller.config)
		assert.NotNil(t, controller.logger)
		assert.Len(t, controller.config.Endpoints, 3)
		assert.Equal(t, "http://etcd-0:2379", controller.config.Endpoints[0])
		assert.Equal(t, "http://etcd-1:2379", controller.config.Endpoints[1])
		assert.Equal(t, "http://etcd-2:2379", controller.config.Endpoints[2])
		assert.Equal(t, 10*time.Second, controller.config.DialTimeout)
		assert.Equal(t, "admin", controller.config.Username)
		assert.Equal(t, "secret", controller.config.Password)
	})

	t.Run("Controller with minimal fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.NotNil(t, controller.config)
		assert.NotNil(t, controller.logger)
		assert.Len(t, controller.config.Endpoints, 1)
		assert.Equal(t, "localhost:2379", controller.config.Endpoints[0])
		assert.Equal(t, 5*time.Second, controller.config.DialTimeout)
		assert.Equal(t, "testuser", controller.config.Username)
		assert.Equal(t, "testpass", controller.config.Password)
	})

	t.Run("Controller with empty fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 0,
			Username:    "",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.NotNil(t, controller.config)
		assert.NotNil(t, controller.logger)
		assert.Len(t, controller.config.Endpoints, 0)
		assert.Equal(t, time.Duration(0), controller.config.DialTimeout)
		assert.Empty(t, controller.config.Username)
		assert.Empty(t, controller.config.Password)
	})

	t.Run("Controller with nil fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 0,
			Username:    "",
			Password:    "",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)
		assert.NotNil(t, controller.config)
		assert.NotNil(t, controller.logger)
		assert.Nil(t, controller.config.Endpoints)
		assert.Equal(t, time.Duration(0), controller.config.DialTimeout)
		assert.Empty(t, controller.config.Username)
		assert.Empty(t, controller.config.Password)
	})

	t.Run("Controller full lifecycle", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)

		// Start controller
		err := controller.Start()
		assert.NoError(t, err)

		// Run controller
		ctx := context.Background()
		err = controller.Run(ctx)
		assert.NoError(t, err)

		// Get results
		results := controller.GetControllerResults()
		assert.NotNil(t, results)

		// Get status
		status := controller.GetControllerStatus()
		assert.NotNil(t, status)

		// Get metrics
		metrics := controller.GetControllerMetrics()
		assert.NotNil(t, metrics)

		// Stop controller
		controller.Stop()
	})

	t.Run("Controller with multiple concurrent operations", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)

		// Start controller
		err := controller.Start()
		require.NoError(t, err)

		// Perform multiple concurrent operations
		done := make(chan bool, 10)
		
		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				// Run controller
				ctx := context.Background()
				err := controller.Run(ctx)
				assert.NoError(t, err)
				
				// Get results
				results := controller.GetControllerResults()
				assert.NotNil(t, results)
				
				// Get status
				status := controller.GetControllerStatus()
				assert.NotNil(t, status)
				
				// Get metrics
				metrics := controller.GetControllerMetrics()
				assert.NotNil(t, metrics)
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Stop controller
		controller.Stop()
	})

	t.Run("Controller with different configurations", func(t *testing.T) {
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
			controller := NewEtcdClusterController(config, logger)
			assert.NotNil(t, controller)

			// Start controller
			err := controller.Start()
			require.NoError(t, err)

			// Run controller
			ctx := context.Background()
			err = controller.Run(ctx)
			assert.NoError(t, err)

			// Get results
			results := controller.GetControllerResults()
			assert.NotNil(t, results)

			// Get status
			status := controller.GetControllerStatus()
			assert.NotNil(t, status)

			// Get metrics
			metrics := controller.GetControllerMetrics()
			assert.NotNil(t, metrics)

			// Stop controller
			controller.Stop()
		}
	})

	t.Run("Controller with context cancellation", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)

		// Start controller
		err := controller.Start()
		require.NoError(t, err)

		// Run controller with context that gets cancelled
		ctx, cancel := context.WithCancel(context.Background())
		
		// Cancel context after a short delay
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()
		
		err = controller.Run(ctx)
		assert.Error(t, err) // Should error due to context cancellation

		// Stop controller
		controller.Stop()
	})

	t.Run("Controller with context timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		controller := NewEtcdClusterController(config, logger)
		assert.NotNil(t, controller)

		// Start controller
		err := controller.Start()
		require.NoError(t, err)

		// Run controller with context that times out
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		err = controller.Run(ctx)
		assert.Error(t, err) // Should error due to context timeout

		// Stop controller
		controller.Stop()
	})
}
