// +build integration

package integration_tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// TestConfig holds configuration for integration tests
type TestConfig struct {
	EtcdEndpoints []string
	DialTimeout   time.Duration
}

// DefaultTestConfig returns default test configuration
func DefaultTestConfig() *TestConfig {
	return &TestConfig{
		EtcdEndpoints: []string{"localhost:2379"},
		DialTimeout:   5 * time.Second,
	}
}

// SetupEtcdClient creates a new etcd client for testing
func SetupEtcdClient(t *testing.T, config *TestConfig) *clientv3.Client {
	if config == nil {
		config = DefaultTestConfig()
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   config.EtcdEndpoints,
		DialTimeout: config.DialTimeout,
	})
	require.NoError(t, err, "Failed to create etcd client")

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), config.DialTimeout)
	defer cancel()

	_, err = client.Status(ctx, config.EtcdEndpoints[0])
	require.NoError(t, err, "Failed to connect to etcd")

	return client
}

// CleanupEtcdClient closes etcd client and cleans up test data
func CleanupEtcdClient(t *testing.T, client *clientv3.Client, prefix string) {
	if client == nil {
		return
	}

	// Clean up test data
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if prefix != "" {
		_, err := client.Delete(ctx, prefix, clientv3.WithPrefix())
		if err != nil {
			t.Logf("Failed to cleanup test data: %v", err)
		}
	}

	// Close client
	err := client.Close()
	if err != nil {
		t.Logf("Failed to close etcd client: %v", err)
	}
}

// WaitForEtcd waits for etcd to be ready
func WaitForEtcd(endpoints []string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		client, err := clientv3.New(clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: 2 * time.Second,
		})
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			_, err = client.Status(ctx, endpoints[0])
			cancel()
			client.Close()

			if err == nil {
				return nil
			}
		}

		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for etcd to be ready")
}

// GenerateTestKey generates a unique test key
func GenerateTestKey(prefix string) string {
	return fmt.Sprintf("%s/test-%d", prefix, time.Now().UnixNano())
}

// PutTestData puts test data into etcd
func PutTestData(t *testing.T, client *clientv3.Client, key, value string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Put(ctx, key, value)
	require.NoError(t, err, "Failed to put test data")
}

// GetTestData retrieves test data from etcd
func GetTestData(t *testing.T, client *clientv3.Client, key string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, key)
	require.NoError(t, err, "Failed to get test data")
	require.Greater(t, len(resp.Kvs), 0, "No data found for key: %s", key)

	return string(resp.Kvs[0].Value)
}

// DeleteTestData deletes test data from etcd
func DeleteTestData(t *testing.T, client *clientv3.Client, key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Delete(ctx, key)
	require.NoError(t, err, "Failed to delete test data")
}

// CreateTestLease creates a lease for testing
func CreateTestLease(t *testing.T, client *clientv3.Client, ttl int64) clientv3.LeaseID {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Grant(ctx, ttl)
	require.NoError(t, err, "Failed to create lease")

	return resp.ID
}

// RevokeTestLease revokes a test lease
func RevokeTestLease(t *testing.T, client *clientv3.Client, leaseID clientv3.LeaseID) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Revoke(ctx, leaseID)
	require.NoError(t, err, "Failed to revoke lease")
}

// VerifyKeyExists verifies that a key exists in etcd
func VerifyKeyExists(t *testing.T, client *clientv3.Client, key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, key)
	require.NoError(t, err, "Failed to get key")
	require.Greater(t, len(resp.Kvs), 0, "Key does not exist: %s", key)
}

// VerifyKeyNotExists verifies that a key does not exist in etcd
func VerifyKeyNotExists(t *testing.T, client *clientv3.Client, key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, key)
	require.NoError(t, err, "Failed to get key")
	require.Equal(t, 0, len(resp.Kvs), "Key exists but should not: %s", key)
}

// WatchForKey waits for a key to appear with a specific value
func WatchForKey(t *testing.T, client *clientv3.Client, key, expectedValue string, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	watchChan := client.Watch(ctx, key)

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("Timeout waiting for key %s with value %s", key, expectedValue)
		case resp := <-watchChan:
			for _, event := range resp.Events {
				if string(event.Kv.Value) == expectedValue {
					return
				}
			}
		}
	}
}

// GetClusterMembers returns the list of cluster members
func GetClusterMembers(t *testing.T, client *clientv3.Client) []*clientv3.Member {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.MemberList(ctx)
	require.NoError(t, err, "Failed to get cluster members")

	return resp.Members
}

// VerifyClusterHealth verifies that the cluster is healthy
func VerifyClusterHealth(t *testing.T, client *clientv3.Client) {
	endpoints := client.Endpoints()

	for _, endpoint := range endpoints {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		status, err := client.Status(ctx, endpoint)
		cancel()

		require.NoError(t, err, "Failed to get status for endpoint %s", endpoint)
		require.NotNil(t, status, "Status is nil for endpoint %s", endpoint)
	}
}
