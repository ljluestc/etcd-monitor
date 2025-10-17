// +build integration

package integration_tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestEtcdIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	config := DefaultTestConfig()
	client := SetupEtcdClient(t, config)
	defer CleanupEtcdClient(t, client, "/test")

	t.Run("Put and Get", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		key := GenerateTestKey("/test/integration")
		value := "test-value"

		// Put
		_, err := client.Put(ctx, key, value)
		require.NoError(t, err)

		// Get
		resp, err := client.Get(ctx, key)
		require.NoError(t, err)
		require.Equal(t, int64(1), resp.Count)
		assert.Equal(t, value, string(resp.Kvs[0].Value))

		// Cleanup
		_, err = client.Delete(ctx, key)
		require.NoError(t, err)
	})

	t.Run("Put with Lease", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		key := GenerateTestKey("/test/lease")
		value := "lease-value"

		// Create lease
		leaseResp, err := client.Grant(ctx, 5)
		require.NoError(t, err)

		// Put with lease
		_, err = client.Put(ctx, key, value, clientv3.WithLease(leaseResp.ID))
		require.NoError(t, err)

		// Verify key exists
		VerifyKeyExists(t, client, key)

		// Wait for lease to expire
		time.Sleep(6 * time.Second)

		// Verify key no longer exists
		VerifyKeyNotExists(t, client, key)
	})

	t.Run("Watch", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		key := GenerateTestKey("/test/watch")
		value := "watch-value"

		// Start watching
		watchChan := client.Watch(ctx, key)

		// Put value in goroutine
		go func() {
			time.Sleep(100 * time.Millisecond)
			putCtx, putCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer putCancel()
			client.Put(putCtx, key, value)
		}()

		// Wait for watch event
		select {
		case resp := <-watchChan:
			require.NotNil(t, resp)
			require.Greater(t, len(resp.Events), 0)
			assert.Equal(t, value, string(resp.Events[0].Kv.Value))
		case <-ctx.Done():
			t.Fatal("Timeout waiting for watch event")
		}

		// Cleanup
		delCtx, delCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer delCancel()
		client.Delete(delCtx, key)
	})

	t.Run("Transaction", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		key := GenerateTestKey("/test/txn")
		oldValue := "old-value"
		newValue := "new-value"

		// Put initial value
		_, err := client.Put(ctx, key, oldValue)
		require.NoError(t, err)

		// Transaction with compare-and-swap
		txn := client.Txn(ctx)
		cmp := clientv3.Compare(clientv3.Value(key), "=", oldValue)
		putOp := clientv3.OpPut(key, newValue)
		getOp := clientv3.OpGet(key)

		resp, err := txn.If(cmp).Then(putOp).Else(getOp).Commit()
		require.NoError(t, err)
		assert.True(t, resp.Succeeded, "Transaction should succeed")

		// Verify new value
		getResp, err := client.Get(ctx, key)
		require.NoError(t, err)
		assert.Equal(t, newValue, string(getResp.Kvs[0].Value))

		// Cleanup
		_, err = client.Delete(ctx, key)
		require.NoError(t, err)
	})

	t.Run("Prefix Operations", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		prefix := "/test/prefix"
		keys := []string{
			prefix + "/key1",
			prefix + "/key2",
			prefix + "/key3",
		}

		// Put multiple keys
		for i, key := range keys {
			_, err := client.Put(ctx, key, string(rune('a'+i)))
			require.NoError(t, err)
		}

		// Get with prefix
		resp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
		require.NoError(t, err)
		assert.Equal(t, int64(len(keys)), resp.Count)

		// Delete with prefix
		delResp, err := client.Delete(ctx, prefix, clientv3.WithPrefix())
		require.NoError(t, err)
		assert.Equal(t, int64(len(keys)), delResp.Deleted)

		// Verify all keys deleted
		resp, err = client.Get(ctx, prefix, clientv3.WithPrefix())
		require.NoError(t, err)
		assert.Equal(t, int64(0), resp.Count)
	})

	t.Run("Cluster Health", func(t *testing.T) {
		VerifyClusterHealth(t, client)
	})

	t.Run("Member List", func(t *testing.T) {
		members := GetClusterMembers(t, client)
		assert.Greater(t, len(members), 0, "Should have at least one member")

		for _, member := range members {
			assert.NotEmpty(t, member.Name)
			assert.Greater(t, len(member.ClientURLs), 0)
		}
	})

	t.Run("Endpoint Status", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		endpoints := client.Endpoints()
		require.Greater(t, len(endpoints), 0)

		for _, endpoint := range endpoints {
			status, err := client.Status(ctx, endpoint)
			require.NoError(t, err)
			assert.NotEmpty(t, status.Version)
			assert.Greater(t, status.DbSize, int64(0))
		}
	})

	t.Run("Compact and Defragment", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Put and delete some data to create history
		key := GenerateTestKey("/test/compact")
		for i := 0; i < 10; i++ {
			_, err := client.Put(ctx, key, string(rune('a'+i)))
			require.NoError(t, err)
		}

		// Get current revision
		resp, err := client.Get(ctx, key)
		require.NoError(t, err)
		revision := resp.Header.Revision

		// Compact
		_, err = client.Compact(ctx, revision)
		require.NoError(t, err)

		// Cleanup
		_, err = client.Delete(ctx, key)
		require.NoError(t, err)
	})

	t.Run("Keep Alive Lease", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Create lease
		leaseResp, err := client.Grant(ctx, 5)
		require.NoError(t, err)

		// Start keep alive
		keepAliveChan, err := client.KeepAlive(ctx, leaseResp.ID)
		require.NoError(t, err)

		// Wait for a few keep alive responses
		count := 0
		timeout := time.After(10 * time.Second)
		for count < 3 {
			select {
			case resp := <-keepAliveChan:
				require.NotNil(t, resp)
				count++
			case <-timeout:
				t.Fatal("Timeout waiting for keep alive responses")
			}
		}

		// Revoke lease
		_, err = client.Revoke(ctx, leaseResp.ID)
		require.NoError(t, err)
	})

	t.Run("Time To Live", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Create lease
		leaseResp, err := client.Grant(ctx, 10)
		require.NoError(t, err)

		// Get TTL
		ttlResp, err := client.TimeToLive(ctx, leaseResp.ID)
		require.NoError(t, err)
		assert.Greater(t, ttlResp.TTL, int64(0))
		assert.LessOrEqual(t, ttlResp.TTL, int64(10))

		// Revoke lease
		_, err = client.Revoke(ctx, leaseResp.ID)
		require.NoError(t, err)
	})
}

func TestEtcdConcurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	config := DefaultTestConfig()
	client := SetupEtcdClient(t, config)
	defer CleanupEtcdClient(t, client, "/test/concurrent")

	t.Run("Concurrent Writes", func(t *testing.T) {
		prefix := "/test/concurrent"
		numGoroutines := 10
		numOperations := 10

		done := make(chan bool, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer func() { done <- true }()

				for j := 0; j < numOperations; j++ {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					key := GenerateTestKey(prefix)
					_, err := client.Put(ctx, key, "value")
					cancel()
					require.NoError(t, err)
				}
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < numGoroutines; i++ {
			<-done
		}

		// Verify all keys written
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		resp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
		require.NoError(t, err)
		assert.Equal(t, int64(numGoroutines*numOperations), resp.Count)
	})
}
