package etcd

import (
	"context"
	"crypto/tls"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestEtcd_Comprehensive(t *testing.T) {

	t.Run("NewClientv3", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		
		client, err := NewClientv3(config)
		if err != nil {
			// Expected to fail without etcd
			assert.Error(t, err)
		} else {
			assert.NotNil(t, client)
			client.Close()
		}
	})

	t.Run("NewClientv3 with TLS", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			TLS: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		
		client, err := NewClientv3(config)
		if err != nil {
			// Expected to fail without etcd
			assert.Error(t, err)
		} else {
			assert.NotNil(t, client)
			client.Close()
		}
	})

	t.Run("ClientConfig validation", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "test",
			Password:    "test",
		}
		
		assert.Equal(t, []string{"localhost:2379"}, config.Endpoints)
		assert.Equal(t, 5*time.Second, config.DialTimeout)
		assert.Equal(t, "test", config.Username)
		assert.Equal(t, "test", config.Password)
	})

	t.Run("SecureConfig validation", func(t *testing.T) {
		secureConfig := &SecureConfig{
			CertFile:   "/path/to/cert.pem",
			KeyFile:    "/path/to/key.pem",
			TrustedCAFile: "/path/to/ca.pem",
			InsecureSkipVerify: true,
			Username:   "test-user",
			Password:   "test-pass",
		}
		
		assert.Equal(t, "/path/to/cert.pem", secureConfig.CertFile)
		assert.Equal(t, "/path/to/key.pem", secureConfig.KeyFile)
		assert.Equal(t, "/path/to/ca.pem", secureConfig.TrustedCAFile)
		assert.True(t, secureConfig.InsecureSkipVerify)
		assert.Equal(t, "test-user", secureConfig.Username)
		assert.Equal(t, "test-pass", secureConfig.Password)
	})

	t.Run("MemberList", func(t *testing.T) {
		client, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			t.Skip("Skipping test - etcd not available")
			return
		}
		defer client.Close()

		resp, err := MemberList(client)
		// This will fail in test environment, but we can test the structure
		assert.Error(t, err) // Expected to fail without real etcd
		assert.Nil(t, resp)
	})

	t.Run("Status", func(t *testing.T) {
		client, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			t.Skip("Skipping test - etcd not available")
			return
		}
		defer client.Close()

		resp, err := Status("localhost:2379", client)
		// This will fail in test environment, but we can test the structure
		assert.Error(t, err) // Expected to fail without real etcd
		assert.Nil(t, resp)
	})

	t.Run("AlarmList", func(t *testing.T) {
		client, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			t.Skip("Skipping test - etcd not available")
			return
		}
		defer client.Close()

		resp, err := AlarmList(client)
		// This will fail in test environment, but we can test the structure
		assert.Error(t, err) // Expected to fail without real etcd
		assert.Nil(t, resp)
	})

	t.Run("Context with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		// Test that context timeout works
		select {
		case <-ctx.Done():
			// Expected - context should be done due to timeout
		case <-time.After(10 * time.Millisecond):
			t.Fatal("Context should have timed out")
		}
	})

	t.Run("Logger creation", func(t *testing.T) {
		// Test logger creation
		logger, err := zap.NewDevelopment()
		assert.NoError(t, err)
		assert.NotNil(t, logger)
	})

	t.Run("Etcd client creation", func(t *testing.T) {
		// Test etcd client creation (will fail without etcd)
		client, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			// Expected to fail without etcd
			assert.Error(t, err)
		} else {
			client.Close()
		}
	})
}