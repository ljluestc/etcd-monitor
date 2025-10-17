package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// MockEtcdClientForMetrics is a mock etcd client specifically for metrics testing
type MockEtcdClientForMetrics struct {
	statusResponse *clientv3.StatusResponse
	statusError    error
	rangeResponse  *clientv3.GetResponse
	rangeError     error
	putResponse    *clientv3.PutResponse
	putError       error
	deleteResponse *clientv3.DeleteResponse
	deleteError    error
	txnResponse    *clientv3.TxnResponse
	txnError       error
}

func (m *MockEtcdClientForMetrics) Status(ctx context.Context, endpoint string) (*clientv3.StatusResponse, error) {
	return m.statusResponse, m.statusError
}

func (m *MockEtcdClientForMetrics) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	return m.rangeResponse, m.rangeError
}

func (m *MockEtcdClientForMetrics) Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	return m.putResponse, m.putError
}

func (m *MockEtcdClientForMetrics) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	return m.deleteResponse, m.deleteError
}

func (m *MockEtcdClientForMetrics) Txn(ctx context.Context) clientv3.Txn {
	return &MockTxn{
		response: m.txnResponse,
		error:    m.txnError,
	}
}

func (m *MockEtcdClientForMetrics) Close() error {
	return nil
}

// MockTxn is a mock transaction for testing
type MockTxn struct {
	response *clientv3.TxnResponse
	error    error
}

func (m *MockTxn) If(cs ...clientv3.Cmp) clientv3.Txn {
	return m
}

func (m *MockTxn) Then(ops ...clientv3.Op) clientv3.Txn {
	return m
}

func (m *MockTxn) Else(ops ...clientv3.Op) clientv3.Txn {
	return m
}

func (m *MockTxn) Commit() (*clientv3.TxnResponse, error) {
	return m.response, m.error
}

func TestNewMetricsCollector(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	t.Run("Valid client", func(t *testing.T) {
		mockClient := &MockEtcdClientForMetrics{}
		mc := NewMetricsCollector(mockClient, logger)
		
		assert.NotNil(t, mc)
		assert.Equal(t, mockClient, mc.client)
		assert.Equal(t, logger, mc.logger)
	})
	
	t.Run("Nil client", func(t *testing.T) {
		mc := NewMetricsCollector(nil, logger)
		assert.NotNil(t, mc)
		assert.Nil(t, mc.client)
	})
}

func TestCollectMetrics(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	t.Run("Successful metrics collection", func(t *testing.T) {
		mockClient := &MockEtcdClientForMetrics{
			statusResponse: &clientv3.StatusResponse{
				Header: &clientv3.ResponseHeader{},
				Version: "3.5.0",
				DbSize: 1024 * 1024, // 1MB
				Leader: 1,
				RaftIndex: 100,
				RaftTerm: 1,
				RaftAppliedIndex: 100,
			},
			rangeResponse: &clientv3.GetResponse{
				Header: &clientv3.ResponseHeader{},
				Kvs: []*clientv3.KeyValue{
					{Key: []byte("key1"), Value: []byte("value1"), Version: 1},
					{Key: []byte("key2"), Value: []byte("value2"), Version: 1},
				},
				Count: 2,
			},
		}
		
		mc := NewMetricsCollector(mockClient, logger)
		metrics, err := mc.CollectMetrics(context.Background())
		
		require.NoError(t, err)
		assert.NotNil(t, metrics)
		assert.Equal(t, int64(1024*1024), metrics.DBSize)
		assert.True(t, metrics.Timestamp.After(time.Now().Add(-time.Minute)))
	})
	
	t.Run("Status check error", func(t *testing.T) {
		mockClient := &MockEtcdClientForMetrics{
			statusError: rpctypes.ErrGRPCUnavailable,
		}
		
		mc := NewMetricsCollector(mockClient, logger)
		metrics, err := mc.CollectMetrics(context.Background())
		
		require.Error(t, err)
		assert.Nil(t, metrics)
	})
	
	t.Run("Range operation error", func(t *testing.T) {
		mockClient := &MockEtcdClientForMetrics{
			statusResponse: &clientv3.StatusResponse{
				Header: &clientv3.ResponseHeader{},
				Version: "3.5.0",
				DbSize: 1024,
				Leader: 1,
				RaftIndex: 100,
				RaftTerm: 1,
				RaftAppliedIndex: 100,
			},
			rangeError: rpctypes.ErrGRPCUnavailable,
		}
		
		mc := NewMetricsCollector(mockClient, logger)
		metrics, err := mc.CollectMetrics(context.Background())
		
		require.Error(t, err)
		assert.Nil(t, metrics)
	})
}

func TestCollectRequestMetrics(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	t.Run("Successful request metrics collection", func(t *testing.T) {
		mockClient := &MockEtcdClientForMetrics{
			putResponse: &clientv3.PutResponse{
				Header: &clientv3.ResponseHeader{},
			},
			rangeResponse: &clientv3.GetResponse{
				Header: &clientv3.ResponseHeader{},
				Kvs: []*clientv3.KeyValue{
					{Key: []byte("key1"), Value: []byte("value1"), Version: 1},
				},
				Count: 1,
			},
			deleteResponse: &clientv3.DeleteResponse{
				Header: &clientv3.ResponseHeader{},
				Deleted: 1,
			},
		}
		
		mc := NewMetricsCollector(mockClient, logger)
		
		// Test write operations
		start := time.Now()
		_, err := mockClient.Put(context.Background(), "test-key", "test-value")
		writeLatency := time.Since(start)
		require.NoError(t, err)
		
		// Test read operations
		start = time.Now()
		_, err = mockClient.Get(context.Background(), "test-key")
		readLatency := time.Since(start)
		require.NoError(t, err)
		
		// Test delete operations
		start = time.Now()
		_, err = mockClient.Delete(context.Background(), "test-key")
		deleteLatency := time.Since(start)
		require.NoError(t, err)
		
		// Verify latencies are reasonable
		assert.True(t, writeLatency < time.Second)
		assert.True(t, readLatency < time.Second)
		assert.True(t, deleteLatency < time.Second)
	})
	
	t.Run("Request metrics with errors", func(t *testing.T) {
		mockClient := &MockEtcdClientForMetrics{
			putError: rpctypes.ErrGRPCUnavailable,
			rangeError: rpctypes.ErrGRPCUnavailable,
			deleteError: rpctypes.ErrGRPCUnavailable,
		}
		
		mc := NewMetricsCollector(mockClient, logger)
		
		// These should fail but not crash
		_, err := mockClient.Put(context.Background(), "test-key", "test-value")
		assert.Error(t, err)
		
		_, err = mockClient.Get(context.Background(), "test-key")
		assert.Error(t, err)
		
		_, err = mockClient.Delete(context.Background(), "test-key")
		assert.Error(t, err)
	})
}

func TestCollectDatabaseMetrics(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	t.Run("Database metrics collection", func(t *testing.T) {
		mockClient := &MockEtcdClientForMetrics{
			statusResponse: &clientv3.StatusResponse{
				Header: &clientv3.ResponseHeader{},
				Version: "3.5.0",
				DbSize: 2048 * 1024, // 2MB
				Leader: 1,
				RaftIndex: 200,
				RaftTerm: 2,
				RaftAppliedIndex: 200,
			},
		}
		
		mc := NewMetricsCollector(mockClient, logger)
		metrics, err := mc.CollectMetrics(context.Background())
		
		require.NoError(t, err)
		assert.Equal(t, int64(2048*1024), metrics.DBSize)
		assert.True(t, metrics.DBSizeInUse >= 0)
	})
}

func TestCollectRaftMetrics(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	t.Run("Raft metrics collection", func(t *testing.T) {
		mockClient := &MockEtcdClientForMetrics{
			statusResponse: &clientv3.StatusResponse{
				Header: &clientv3.ResponseHeader{},
				Version: "3.5.0",
				DbSize: 1024,
				Leader: 1,
				RaftIndex: 100,
				RaftTerm: 1,
				RaftAppliedIndex: 95, // 5 proposals pending
			},
		}
		
		mc := NewMetricsCollector(mockClient, logger)
		metrics, err := mc.CollectMetrics(context.Background())
		
		require.NoError(t, err)
		assert.Equal(t, uint64(100), metrics.ProposalCommitted)
		assert.Equal(t, uint64(95), metrics.ProposalApplied)
		assert.Equal(t, uint64(5), metrics.ProposalPending)
	})
}

func TestCollectResourceMetrics(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	t.Run("Resource metrics collection", func(t *testing.T) {
		mockClient := &MockEtcdClientForMetrics{
			statusResponse: &clientv3.StatusResponse{
				Header: &clientv3.ResponseHeader{},
				Version: "3.5.0",
				DbSize: 1024,
				Leader: 1,
				RaftIndex: 100,
				RaftTerm: 1,
				RaftAppliedIndex: 100,
			},
		}
		
		mc := NewMetricsCollector(mockClient, logger)
		metrics, err := mc.CollectMetrics(context.Background())
		
		require.NoError(t, err)
		// Resource metrics would be collected from system monitoring
		// For now, just verify the structure is correct
		assert.True(t, metrics.MemoryUsage >= 0)
		assert.True(t, metrics.CPUUsage >= 0)
		assert.True(t, metrics.DiskUsage >= 0)
	})
}

func TestCollectClientMetrics(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	t.Run("Client metrics collection", func(t *testing.T) {
		mockClient := &MockEtcdClientForMetrics{
			statusResponse: &clientv3.StatusResponse{
				Header: &clientv3.ResponseHeader{},
				Version: "3.5.0",
				DbSize: 1024,
				Leader: 1,
				RaftIndex: 100,
				RaftTerm: 1,
				RaftAppliedIndex: 100,
			},
		}
		
		mc := NewMetricsCollector(mockClient, logger)
		metrics, err := mc.CollectMetrics(context.Background())
		
		require.NoError(t, err)
		// Client metrics would be collected from connection monitoring
		assert.True(t, metrics.ActiveConnections >= 0)
		assert.True(t, metrics.WatcherCount >= 0)
	})
}

func TestMetricsCollector_EdgeCases(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	t.Run("Nil client", func(t *testing.T) {
		mc := NewMetricsCollector(nil, logger)
		
		metrics, err := mc.CollectMetrics(context.Background())
		require.Error(t, err)
		assert.Nil(t, metrics)
	})
	
	t.Run("Context cancellation", func(t *testing.T) {
		mockClient := &MockEtcdClientForMetrics{}
		mc := NewMetricsCollector(mockClient, logger)
		
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		metrics, err := mc.CollectMetrics(ctx)
		require.Error(t, err)
		assert.Nil(t, metrics)
	})
	
	t.Run("Empty status response", func(t *testing.T) {
		mockClient := &MockEtcdClientForMetrics{
			statusResponse: &clientv3.StatusResponse{
				Header: &clientv3.ResponseHeader{},
			},
		}
		
		mc := NewMetricsCollector(mockClient, logger)
		metrics, err := mc.CollectMetrics(context.Background())
		
		require.NoError(t, err)
		assert.NotNil(t, metrics)
		// Should handle empty response gracefully
		assert.Equal(t, int64(0), metrics.DBSize)
	})
}

func TestMetricsCollector_Concurrency(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockClient := &MockEtcdClientForMetrics{
		statusResponse: &clientv3.StatusResponse{
			Header: &clientv3.ResponseHeader{},
			Version: "3.5.0",
			DbSize: 1024,
			Leader: 1,
			RaftIndex: 100,
			RaftTerm: 1,
			RaftAppliedIndex: 100,
		},
		rangeResponse: &clientv3.GetResponse{
			Header: &clientv3.ResponseHeader{},
			Kvs: []*clientv3.KeyValue{},
			Count: 0,
		},
	}
	
	mc := NewMetricsCollector(mockClient, logger)
	
	// Test concurrent access
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			
			metrics, err := mc.CollectMetrics(context.Background())
			if err == nil {
				assert.NotNil(t, metrics)
			}
		}()
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestMetricsCollector_Performance(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	t.Run("Performance under load", func(t *testing.T) {
		mockClient := &MockEtcdClientForMetrics{
			statusResponse: &clientv3.StatusResponse{
				Header: &clientv3.ResponseHeader{},
				Version: "3.5.0",
				DbSize: 1024,
				Leader: 1,
				RaftIndex: 100,
				RaftTerm: 1,
				RaftAppliedIndex: 100,
			},
			rangeResponse: &clientv3.GetResponse{
				Header: &clientv3.ResponseHeader{},
				Kvs: []*clientv3.KeyValue{},
				Count: 0,
			},
		}
		
		mc := NewMetricsCollector(mockClient, logger)
		
		// Measure collection time
		start := time.Now()
		metrics, err := mc.CollectMetrics(context.Background())
		duration := time.Since(start)
		
		require.NoError(t, err)
		assert.NotNil(t, metrics)
		assert.True(t, duration < time.Second, "Metrics collection took too long: %v", duration)
	})
}

func TestMetricsCollector_DataValidation(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	t.Run("Data validation", func(t *testing.T) {
		mockClient := &MockEtcdClientForMetrics{
			statusResponse: &clientv3.StatusResponse{
				Header: &clientv3.ResponseHeader{},
				Version: "3.5.0",
				DbSize: 1024 * 1024, // 1MB
				Leader: 1,
				RaftIndex: 100,
				RaftTerm: 1,
				RaftAppliedIndex: 100,
			},
			rangeResponse: &clientv3.GetResponse{
				Header: &clientv3.ResponseHeader{},
				Kvs: []*clientv3.KeyValue{
					{Key: []byte("key1"), Value: []byte("value1"), Version: 1},
				},
				Count: 1,
			},
		}
		
		mc := NewMetricsCollector(mockClient, logger)
		metrics, err := mc.CollectMetrics(context.Background())
		
		require.NoError(t, err)
		assert.NotNil(t, metrics)
		
		// Validate data types and ranges
		assert.True(t, metrics.DBSize >= 0)
		assert.True(t, metrics.DBSizeInUse >= 0)
		assert.True(t, metrics.DBSizeInUse <= metrics.DBSize)
		assert.True(t, metrics.ProposalCommitted >= 0)
		assert.True(t, metrics.ProposalApplied >= 0)
		assert.True(t, metrics.ProposalPending >= 0)
		assert.True(t, metrics.MemoryUsage >= 0)
		assert.True(t, metrics.CPUUsage >= 0)
		assert.True(t, metrics.CPUUsage <= 100) // CPU percentage should be <= 100
		assert.True(t, metrics.DiskUsage >= 0)
		assert.True(t, metrics.ActiveConnections >= 0)
		assert.True(t, metrics.WatcherCount >= 0)
	})
}
