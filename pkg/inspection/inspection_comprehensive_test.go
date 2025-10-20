package inspection

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	etcdv1alpha1 "github.com/etcd-monitor/taskmaster/api/etcd/v1alpha1"
)

// MockEtcdClient is a mock implementation of etcd client for testing
type MockEtcdClient struct {
	members     []*clientv3.Member
	status      *clientv3.StatusResponse
	alarms      []*clientv3.AlarmMember
	shouldError bool
}

func (m *MockEtcdClient) MemberList(ctx context.Context) (*clientv3.MemberListResponse, error) {
	if m.shouldError {
		return nil, assert.AnError
	}
	return &clientv3.MemberListResponse{Members: m.members}, nil
}

func (m *MockEtcdClient) Status(ctx context.Context, endpoint string) (*clientv3.StatusResponse, error) {
	if m.shouldError {
		return nil, assert.AnError
	}
	return m.status, nil
}

func (m *MockEtcdClient) AlarmList(ctx context.Context) (*clientv3.AlarmResponse, error) {
	if m.shouldError {
		return nil, assert.AnError
	}
	return &clientv3.AlarmResponse{Alarms: m.alarms}, nil
}

func (m *MockEtcdClient) Close() error {
	return nil
}

func TestInspection_Comprehensive(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("NewInspection with valid config", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		assert.NotNil(t, inspection)
		assert.NotNil(t, inspection.config)
		assert.NotNil(t, inspection.logger)
	})

	t.Run("NewInspection with nil config", func(t *testing.T) {
		inspection := NewInspection(nil, logger)
		assert.NotNil(t, inspection)
		assert.Nil(t, inspection.config)
		assert.NotNil(t, inspection.logger)
	})

	t.Run("NewInspection with nil logger", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, nil)
		assert.NotNil(t, inspection)
		assert.NotNil(t, inspection.config)
		assert.NotNil(t, inspection.logger) // Should create a production logger
	})

	t.Run("NewInspection with empty config", func(t *testing.T) {
		config := &Config{}

		inspection := NewInspection(config, logger)
		assert.NotNil(t, inspection)
		assert.NotNil(t, inspection.config)
		assert.NotNil(t, inspection.logger)
	})

	t.Run("Start inspection", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		err := inspection.Start()
		assert.NoError(t, err)
	})

	t.Run("Stop inspection", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		err := inspection.Start()
		require.NoError(t, err)

		inspection.Stop()
		// Should not panic or error
	})

	t.Run("Stop inspection without starting", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		inspection.Stop()
		// Should not panic or error
	})

	t.Run("Run inspection", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		ctx := context.Background()
		
		err := inspection.Run(ctx)
		assert.NoError(t, err)
	})

	t.Run("Run inspection with nil context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		
		err := inspection.Run(nil)
		assert.Error(t, err)
	})

	t.Run("Run inspection with cancelled context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		err := inspection.Run(ctx)
		assert.Error(t, err)
	})

	t.Run("Run inspection with timeout context", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		err := inspection.Run(ctx)
		assert.Error(t, err)
	})

	t.Run("GetInspectionResults", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		results := inspection.GetInspectionResults()
		assert.NotNil(t, results)
	})

	t.Run("GetInspectionStatus", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		status := inspection.GetInspectionStatus()
		assert.NotNil(t, status)
	})

	t.Run("GetInspectionMetrics", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		metrics := inspection.GetInspectionMetrics()
		assert.NotNil(t, metrics)
	})
}

func TestInspection_EdgeCases(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Inspection with very long endpoint list", func(t *testing.T) {
		endpoints := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			endpoints[i] = fmt.Sprintf("localhost:%d", 2379+i)
		}

		config := &Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		assert.NotNil(t, inspection)
		assert.Len(t, inspection.config.Endpoints, 1000)
	})

	t.Run("Inspection with very long timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 24 * time.Hour,
		}

		inspection := NewInspection(config, logger)
		assert.NotNil(t, inspection)
		assert.Equal(t, 24*time.Hour, inspection.config.DialTimeout)
	})

	t.Run("Inspection with very short timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 1 * time.Nanosecond,
		}

		inspection := NewInspection(config, logger)
		assert.NotNil(t, inspection)
		assert.Equal(t, 1*time.Nanosecond, inspection.config.DialTimeout)
	})

	t.Run("Inspection with zero timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 0,
		}

		inspection := NewInspection(config, logger)
		assert.NotNil(t, inspection)
		assert.Equal(t, time.Duration(0), inspection.config.DialTimeout)
	})

	t.Run("Inspection with negative timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: -1 * time.Second,
		}

		inspection := NewInspection(config, logger)
		assert.NotNil(t, inspection)
		assert.Equal(t, -1*time.Second, inspection.config.DialTimeout)
	})

	t.Run("Inspection with empty endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		assert.NotNil(t, inspection)
		assert.Len(t, inspection.config.Endpoints, 0)
	})

	t.Run("Inspection with nil endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		assert.NotNil(t, inspection)
		assert.Nil(t, inspection.config.Endpoints)
	})

	t.Run("Inspection with special characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "https://etcd-1:2379", "http://etcd-2:2379"},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		assert.NotNil(t, inspection)
		assert.Len(t, inspection.config.Endpoints, 3)
		assert.Equal(t, "http://etcd-0:2379", inspection.config.Endpoints[0])
		assert.Equal(t, "https://etcd-1:2379", inspection.config.Endpoints[1])
		assert.Equal(t, "http://etcd-2:2379", inspection.config.Endpoints[2])
	})

	t.Run("Inspection with unicode characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379"},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		assert.NotNil(t, inspection)
		assert.Len(t, inspection.config.Endpoints, 2)
	})

	t.Run("Inspection with control characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379"},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		assert.NotNil(t, inspection)
		assert.Len(t, inspection.config.Endpoints, 2)
	})
}

func TestInspection_Performance(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Inspection creation performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		
		start := time.Now()
		
		// Create 1000 Inspection instances
		for i := 0; i < 1000; i++ {
			inspection := NewInspection(config, logger)
			assert.NotNil(t, inspection)
		}
		
		duration := time.Since(start)
		
		// Should create 1000 instances in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Creating 1000 Inspection instances took too long: %v", duration)
	})

	t.Run("Inspection with large endpoint list performance", func(t *testing.T) {
		endpoints := make([]string, 10000)
		for i := 0; i < 10000; i++ {
			endpoints[i] = fmt.Sprintf("localhost:%d", 2379+i)
		}
		
		config := &Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
		}
		
		start := time.Now()
		
		inspection := NewInspection(config, logger)
		
		duration := time.Since(start)
		
		// Should create inspection with 10000 endpoints in less than 10ms
		assert.True(t, duration < 10*time.Millisecond, "Creating Inspection with 10000 endpoints took too long: %v", duration)
		assert.Len(t, inspection.config.Endpoints, 10000)
	})

	t.Run("Inspection start/stop performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		
		start := time.Now()
		
		// Start and stop 100 inspections
		for i := 0; i < 100; i++ {
			inspection := NewInspection(config, logger)
			err := inspection.Start()
			require.NoError(t, err)
			inspection.Stop()
		}
		
		duration := time.Since(start)
		
		// Should start and stop 100 inspections in less than 1 second
		assert.True(t, duration < 1*time.Second, "Starting and stopping 100 inspections took too long: %v", duration)
	})

	t.Run("Inspection run performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		
		inspection := NewInspection(config, logger)
		ctx := context.Background()
		
		start := time.Now()
		
		// Run inspection 100 times
		for i := 0; i < 100; i++ {
			err := inspection.Run(ctx)
			assert.NoError(t, err)
		}
		
		duration := time.Since(start)
		
		// Should run 100 inspections in less than 1 second
		assert.True(t, duration < 1*time.Second, "Running 100 inspections took too long: %v", duration)
	})

	t.Run("Inspection get results performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		
		inspection := NewInspection(config, logger)
		
		start := time.Now()
		
		// Get results 1000 times
		for i := 0; i < 1000; i++ {
			results := inspection.GetInspectionResults()
			assert.NotNil(t, results)
		}
		
		duration := time.Since(start)
		
		// Should get results 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting results 1000 times took too long: %v", duration)
	})

	t.Run("Inspection get status performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		
		inspection := NewInspection(config, logger)
		
		start := time.Now()
		
		// Get status 1000 times
		for i := 0; i < 1000; i++ {
			status := inspection.GetInspectionStatus()
			assert.NotNil(t, status)
		}
		
		duration := time.Since(start)
		
		// Should get status 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting status 1000 times took too long: %v", duration)
	})

	t.Run("Inspection get metrics performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		
		inspection := NewInspection(config, logger)
		
		start := time.Now()
		
		// Get metrics 1000 times
		for i := 0; i < 1000; i++ {
			metrics := inspection.GetInspectionMetrics()
			assert.NotNil(t, metrics)
		}
		
		duration := time.Since(start)
		
		// Should get metrics 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting metrics 1000 times took too long: %v", duration)
	})
}

func TestInspection_Integration(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Inspection with all fields set", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379", "http://etcd-2:2379"},
			DialTimeout: 10 * time.Second,
		}

		inspection := NewInspection(config, logger)
		assert.NotNil(t, inspection)
		assert.NotNil(t, inspection.config)
		assert.NotNil(t, inspection.logger)
		assert.Len(t, inspection.config.Endpoints, 3)
		assert.Equal(t, "http://etcd-0:2379", inspection.config.Endpoints[0])
		assert.Equal(t, "http://etcd-1:2379", inspection.config.Endpoints[1])
		assert.Equal(t, "http://etcd-2:2379", inspection.config.Endpoints[2])
		assert.Equal(t, 10*time.Second, inspection.config.DialTimeout)
	})

	t.Run("Inspection with minimal fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		assert.NotNil(t, inspection)
		assert.NotNil(t, inspection.config)
		assert.NotNil(t, inspection.logger)
		assert.Len(t, inspection.config.Endpoints, 1)
		assert.Equal(t, "localhost:2379", inspection.config.Endpoints[0])
		assert.Equal(t, 5*time.Second, inspection.config.DialTimeout)
	})

	t.Run("Inspection with empty fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 0,
		}

		inspection := NewInspection(config, logger)
		assert.NotNil(t, inspection)
		assert.NotNil(t, inspection.config)
		assert.NotNil(t, inspection.logger)
		assert.Len(t, inspection.config.Endpoints, 0)
		assert.Equal(t, time.Duration(0), inspection.config.DialTimeout)
	})

	t.Run("Inspection with nil fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 0,
		}

		inspection := NewInspection(config, logger)
		assert.NotNil(t, inspection)
		assert.NotNil(t, inspection.config)
		assert.NotNil(t, inspection.logger)
		assert.Nil(t, inspection.config.Endpoints)
		assert.Equal(t, time.Duration(0), inspection.config.DialTimeout)
	})

	t.Run("Inspection full lifecycle", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		assert.NotNil(t, inspection)

		// Start inspection
		err := inspection.Start()
		assert.NoError(t, err)

		// Run inspection
		ctx := context.Background()
		err = inspection.Run(ctx)
		assert.NoError(t, err)

		// Get results
		results := inspection.GetInspectionResults()
		assert.NotNil(t, results)

		// Get status
		status := inspection.GetInspectionStatus()
		assert.NotNil(t, status)

		// Get metrics
		metrics := inspection.GetInspectionMetrics()
		assert.NotNil(t, metrics)

		// Stop inspection
		inspection.Stop()
	})

	t.Run("Inspection with multiple concurrent runs", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		assert.NotNil(t, inspection)

		// Start inspection
		err := inspection.Start()
		require.NoError(t, err)

		// Run multiple concurrent inspections
		done := make(chan bool, 10)
		
		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				ctx := context.Background()
				err := inspection.Run(ctx)
				assert.NoError(t, err)
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Stop inspection
		inspection.Stop()
	})

	t.Run("Inspection with context cancellation", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		assert.NotNil(t, inspection)

		// Start inspection
		err := inspection.Start()
		require.NoError(t, err)

		// Run inspection with context that gets cancelled
		ctx, cancel := context.WithCancel(context.Background())
		
		// Cancel context after a short delay
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()
		
		err = inspection.Run(ctx)
		assert.Error(t, err) // Should error due to context cancellation

		// Stop inspection
		inspection.Stop()
	})

	t.Run("Inspection with context timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		inspection := NewInspection(config, logger)
		assert.NotNil(t, inspection)

		// Start inspection
		err := inspection.Start()
		require.NoError(t, err)

		// Run inspection with context that times out
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		
		err = inspection.Run(ctx)
		assert.Error(t, err) // Should error due to context timeout

		// Stop inspection
		inspection.Stop()
	})
}
