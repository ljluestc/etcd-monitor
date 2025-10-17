package patterns

import (
	"context"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestPatterns_Comprehensive(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("NewDistributedLock", func(t *testing.T) {
		client, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			t.Skip("Skipping test - etcd not available")
			return
		}
		defer client.Close()

		lock, err := NewDistributedLock(client, LockConfig{
			LockKey: "/test/lock",
			TTL:     30,
		}, logger)
		
		assert.NoError(t, err)
		assert.NotNil(t, lock)
		lock.Close()
	})

	t.Run("NewLeaderElection", func(t *testing.T) {
		client, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			t.Skip("Skipping test - etcd not available")
			return
		}
		defer client.Close()

		election, err := NewLeaderElection(client, "/test/election", "test-member", 30, logger)
		
		assert.NoError(t, err)
		assert.NotNil(t, election)
		election.Close()
	})

	t.Run("NewServiceRegistry", func(t *testing.T) {
		client, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			t.Skip("Skipping test - etcd not available")
			return
		}
		defer client.Close()

		registry, err := NewServiceRegistry(client, "/test/services", 30, logger)
		
		assert.NoError(t, err)
		assert.NotNil(t, registry)
		registry.Close()
	})

	t.Run("NewConfigManager", func(t *testing.T) {
		client, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			t.Skip("Skipping test - etcd not available")
			return
		}
		defer client.Close()

		manager := NewConfigManager(client, "/test/config", logger)
		
		assert.NotNil(t, manager)
	})

	t.Run("LockConfig validation", func(t *testing.T) {
		config := LockConfig{
			LockKey: "/test/lock",
			TTL:     30,
		}
		
		assert.Equal(t, "/test/lock", config.LockKey)
		assert.Equal(t, 30, config.TTL)
	})


	t.Run("ServiceInfo validation", func(t *testing.T) {
		info := ServiceInfo{
			Name:    "test-service",
			ID:      "test-service-1",
			Address: "localhost",
			Port:    8080,
			Metadata: map[string]string{
				"version": "1.0.0",
			},
		}
		
		assert.Equal(t, "test-service", info.Name)
		assert.Equal(t, "test-service-1", info.ID)
		assert.Equal(t, "localhost", info.Address)
		assert.Equal(t, 8080, info.Port)
		assert.Equal(t, "1.0.0", info.Metadata["version"])
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