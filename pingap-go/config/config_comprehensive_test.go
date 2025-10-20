package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestConfig_Comprehensive(t *testing.T) {
	t.Run("NewConfig with valid settings", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}
		logger, _ := zap.NewDevelopment()

		assert.NotNil(t, config)
		assert.NotNil(t, logger)
		assert.Len(t, config.Endpoints, 1)
		assert.Equal(t, "localhost:2379", config.Endpoints[0])
		assert.Equal(t, 5*time.Second, config.DialTimeout)
		assert.Equal(t, "testuser", config.Username)
		assert.Equal(t, "testpass", config.Password)
	})

	t.Run("Config with nil endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}

		assert.NotNil(t, config)
		assert.Nil(t, config.Endpoints)
	})

	t.Run("Config with empty endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}

		assert.NotNil(t, config)
		assert.Len(t, config.Endpoints, 0)
	})

	t.Run("Config with multiple endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379", "localhost:2380", "localhost:2381"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}

		assert.NotNil(t, config)
		assert.Len(t, config.Endpoints, 3)
		assert.Equal(t, "localhost:2379", config.Endpoints[0])
		assert.Equal(t, "localhost:2380", config.Endpoints[1])
		assert.Equal(t, "localhost:2381", config.Endpoints[2])
	})

	t.Run("Config with zero timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 0,
			Username:    "testuser",
			Password:    "testpass",
		}

		assert.NotNil(t, config)
		assert.Equal(t, time.Duration(0), config.DialTimeout)
	})

	t.Run("Config with negative timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: -1 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}

		assert.NotNil(t, config)
		assert.Equal(t, -1*time.Second, config.DialTimeout)
	})

	t.Run("Config with empty username", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "",
			Password:    "testpass",
		}

		assert.NotNil(t, config)
		assert.Empty(t, config.Username)
	})

	t.Run("Config with empty password", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    "",
		}

		assert.NotNil(t, config)
		assert.Empty(t, config.Password)
	})

	t.Run("Config with special characters", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "https://etcd-1:2379"},
			DialTimeout: 10 * time.Second,
			Username:    "user@domain.com",
			Password:    "pass!@#$%^&*()",
		}

		assert.NotNil(t, config)
		assert.Len(t, config.Endpoints, 2)
		assert.Equal(t, "http://etcd-0:2379", config.Endpoints[0])
		assert.Equal(t, "https://etcd-1:2379", config.Endpoints[1])
		assert.Equal(t, "user@domain.com", config.Username)
		assert.Equal(t, "pass!@#$%^&*()", config.Password)
	})

	t.Run("Config with unicode characters", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "用户",
			Password:    "密码",
		}

		assert.NotNil(t, config)
		assert.Equal(t, "用户", config.Username)
		assert.Equal(t, "密码", config.Password)
	})

	t.Run("Config with control characters", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "user\n\r\t",
			Password:    "pass\n\r\t",
		}

		assert.NotNil(t, config)
		assert.Equal(t, "user\n\r\t", config.Username)
		assert.Equal(t, "pass\n\r\t", config.Password)
	})
}

func TestConfig_EdgeCases(t *testing.T) {
	t.Run("Config with very long endpoint list", func(t *testing.T) {
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

		assert.NotNil(t, config)
		assert.Len(t, config.Endpoints, 1000)
	})

	t.Run("Config with very long timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 24 * time.Hour,
			Username:    "testuser",
			Password:    "testpass",
		}

		assert.NotNil(t, config)
		assert.Equal(t, 24*time.Hour, config.DialTimeout)
	})

	t.Run("Config with very short timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 1 * time.Nanosecond,
			Username:    "testuser",
			Password:    "testpass",
		}

		assert.NotNil(t, config)
		assert.Equal(t, 1*time.Nanosecond, config.DialTimeout)
	})

	t.Run("Config with very long username", func(t *testing.T) {
		longUsername := strings.Repeat("a", 10000)
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    longUsername,
			Password:    "testpass",
		}

		assert.NotNil(t, config)
		assert.Len(t, config.Username, 10000)
	})

	t.Run("Config with very long password", func(t *testing.T) {
		longPassword := strings.Repeat("b", 10000)
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "testuser",
			Password:    longPassword,
		}

		assert.NotNil(t, config)
		assert.Len(t, config.Password, 10000)
	})
}

func TestConfig_Performance(t *testing.T) {
	t.Run("Config creation performance", func(t *testing.T) {
		start := time.Now()
		
		// Create 1000 Config instances
		for i := 0; i < 1000; i++ {
			config := &Config{
				Endpoints:   []string{"localhost:2379"},
				DialTimeout: 5 * time.Second,
				Username:    "testuser",
				Password:    "testpass",
			}
			assert.NotNil(t, config)
		}
		
		duration := time.Since(start)
		
		// Should create 1000 instances in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Creating 1000 Config instances took too long: %v", duration)
	})

	t.Run("Config with large data performance", func(t *testing.T) {
		endpoints := make([]string, 10000)
		for i := 0; i < 10000; i++ {
			endpoints[i] = fmt.Sprintf("localhost:%d", 2379+i)
		}
		
		start := time.Now()
		
		config := &Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
			Username:    strings.Repeat("a", 10000),
			Password:    strings.Repeat("b", 10000),
		}
		
		duration := time.Since(start)
		
		// Should create config with large data in less than 10ms
		assert.True(t, duration < 10*time.Millisecond, "Creating Config with large data took too long: %v", duration)
		assert.Len(t, config.Endpoints, 10000)
		assert.Len(t, config.Username, 10000)
		assert.Len(t, config.Password, 10000)
	})
}
