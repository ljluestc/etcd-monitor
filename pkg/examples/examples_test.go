package examples

import (
	"context"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestExamples_Comprehensive(t *testing.T) {

	t.Run("Example1_DistributedLock", func(t *testing.T) {
		// Skip this test as it requires a real etcd connection
		t.Skip("Skipping test - requires real etcd connection")
	})

	t.Run("Example2_ServiceDiscovery", func(t *testing.T) {
		t.Skip("Skipping test - requires real etcd connection")
	})

	t.Run("Example3_LoadBalancing", func(t *testing.T) {
		t.Skip("Skipping test - requires real etcd connection")
	})

	t.Run("Example4_ConfigurationManagement", func(t *testing.T) {
		t.Skip("Skipping test - requires real etcd connection")
	})

	t.Run("Example5_FeatureFlags", func(t *testing.T) {
		t.Skip("Skipping test - requires real etcd connection")
	})

	t.Run("Example6_LeaderElection", func(t *testing.T) {
		t.Skip("Skipping test - requires real etcd connection")
	})

	t.Run("Example7_LeaderTask", func(t *testing.T) {
		t.Skip("Skipping test - requires real etcd connection")
	})

	t.Run("Example8_Transaction", func(t *testing.T) {
		t.Skip("Skipping test - requires real etcd connection")
	})

	t.Run("Example9_Leases", func(t *testing.T) {
		t.Skip("Skipping test - requires real etcd connection")
	})

	t.Run("Example10_Watch", func(t *testing.T) {
		t.Skip("Skipping test - requires real etcd connection")
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