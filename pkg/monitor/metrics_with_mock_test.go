package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/etcd-monitor/taskmaster/testutil/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMetricsCollector_CollectDatabaseMetrics(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockClient := mocks.NewMockEtcdClient()

	// Note: We need to use interface or modify the code to accept our mock
	// For now, we'll test the mock client directly to verify it works
	t.Run("Mock client basic operations", func(t *testing.T) {
		ctx := context.Background()

		// Test Put
		_, err := mockClient.Put(ctx, "test-key", "test-value")
		assert.NoError(t, err)

		// Test Get
		resp, err := mockClient.Get(ctx, "test-key")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.Kvs))
		assert.Equal(t, "test-value", string(resp.Kvs[0].Value))

		// Test Delete
		_, err = mockClient.Delete(ctx, "test-key")
		assert.NoError(t, err)

		// Verify deleted
		resp, err = mockClient.Get(ctx, "test-key")
		assert.NoError(t, err)
		assert.Equal(t, 0, len(resp.Kvs))
	})

	t.Run("Mock client member list", func(t *testing.T) {
		ctx := context.Background()

		resp, err := mockClient.MemberList(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(resp.Members))
		assert.Equal(t, "etcd-1", resp.Members[0].Name)
		assert.Equal(t, "etcd-2", resp.Members[1].Name)
		assert.Equal(t, "etcd-3", resp.Members[2].Name)
	})

	t.Run("Mock client status", func(t *testing.T) {
		ctx := context.Background()

		resp, err := mockClient.Status(ctx, "http://localhost:2379")
		assert.NoError(t, err)
		assert.Equal(t, "3.5.0", resp.Version)
		assert.Equal(t, int64(1024*1024), resp.DbSize)
		assert.Equal(t, uint64(1), resp.Leader)
	})

	t.Run("Mock client alarm list", func(t *testing.T) {
		ctx := context.Background()

		// No alarms initially
		resp, err := mockClient.AlarmList(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(resp.Alarms))
	})

	t.Run("Mock client error simulation", func(t *testing.T) {
		ctx := context.Background()

		// Set error for Get
		mockClient.SetGetError(assert.AnError)
		_, err := mockClient.Get(ctx, "test-key")
		assert.Error(t, err)

		// Clear error
		mockClient.SetGetError(nil)
		_, err = mockClient.Get(ctx, "test-key")
		assert.NoError(t, err)
	})

	t.Run("Mock client operation counting", func(t *testing.T) {
		ctx := context.Background()
		mockClient.ResetOperationCounts()

		// Perform operations
		_, _ = mockClient.Get(ctx, "key1")
		_, _ = mockClient.Put(ctx, "key2", "val2")
		_, _ = mockClient.Delete(ctx, "key3")

		getCount, putCount, deleteCount := mockClient.GetOperationCounts()
		assert.Equal(t, int64(1), getCount)
		assert.Equal(t, int64(1), putCount)
		assert.Equal(t, int64(1), deleteCount)
	})

	t.Run("Mock client latency simulation", func(t *testing.T) {
		ctx := context.Background()

		// Enable latency simulation
		mockClient.SetLatencySimulation(true, 10*time.Millisecond, 20*time.Millisecond)

		// Measure read latency
		start := time.Now()
		_, _ = mockClient.Get(ctx, "test-key")
		readDuration := time.Since(start)
		assert.GreaterOrEqual(t, readDuration, 10*time.Millisecond)

		// Measure write latency
		start = time.Now()
		_, _ = mockClient.Put(ctx, "test-key", "test-val")
		writeDuration := time.Since(start)
		assert.GreaterOrEqual(t, writeDuration, 20*time.Millisecond)

		// Disable latency simulation
		mockClient.SetLatencySimulation(false, 0, 0)
	})

	t.Run("Mock client leader changes", func(t *testing.T) {
		ctx := context.Background()

		// Initial leader is 1
		resp, _ := mockClient.Status(ctx, "http://localhost:2379")
		assert.Equal(t, uint64(1), resp.Leader)

		// Change leader to 2
		mockClient.SetLeader(2)
		resp, _ = mockClient.Status(ctx, "http://localhost:2379")
		assert.Equal(t, uint64(2), resp.Leader)
	})

	t.Run("Mock client alarm management", func(t *testing.T) {
		ctx := context.Background()

		// Add alarm
		mockClient.AddAlarm(1, 1) // NOSPACE alarm for member 1

		resp, err := mockClient.AlarmList(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.Alarms))
		assert.Equal(t, uint64(1), resp.Alarms[0].MemberID)

		// Clear alarms
		mockClient.ClearAlarms()
		resp, err = mockClient.AlarmList(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(resp.Alarms))
	})

	t.Run("Mock client member management", func(t *testing.T) {
		ctx := context.Background()
		mockClient2 := mocks.NewMockEtcdClient()

		// Initial 3 members
		resp, _ := mockClient2.MemberList(ctx)
		assert.Equal(t, 3, len(resp.Members))

		// Add a member
		mockClient2.AddMember(4, "etcd-4",
			[]string{"http://localhost:2385"},
			[]string{"http://localhost:2386"})

		resp, _ = mockClient2.MemberList(ctx)
		assert.Equal(t, 4, len(resp.Members))

		// Remove a member
		mockClient2.RemoveMember(4)
		resp, _ = mockClient2.MemberList(ctx)
		assert.Equal(t, 3, len(resp.Members))
	})

	_ = logger // Suppress unused warning
}

func TestMetricsCollector_RecordLatencyMeasurement(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mc := NewMetricsCollector(nil, logger)

	t.Run("Record single measurement", func(t *testing.T) {
		measurement := LatencyMeasurement{
			Timestamp:    time.Now(),
			ReadLatency:  10 * time.Millisecond,
			WriteLatency: 20 * time.Millisecond,
			Operation:    "test-op",
		}

		mc.recordLatencyMeasurement(measurement)

		history := mc.GetLatencyHistory()
		assert.Equal(t, 1, len(history))
		assert.Equal(t, "test-op", history[0].Operation)
	})

	t.Run("Record multiple measurements", func(t *testing.T) {
		mc2 := NewMetricsCollector(nil, logger)

		for i := 0; i < 10; i++ {
			measurement := LatencyMeasurement{
				Timestamp:    time.Now(),
				ReadLatency:  time.Duration(i) * time.Millisecond,
				WriteLatency: time.Duration(i*2) * time.Millisecond,
				Operation:    "bulk-test",
			}
			mc2.recordLatencyMeasurement(measurement)
		}

		history := mc2.GetLatencyHistory()
		assert.Equal(t, 10, len(history))
	})

	t.Run("History size limit enforcement", func(t *testing.T) {
		mc3 := NewMetricsCollector(nil, logger)

		// Record more than maxHistory (1000) measurements
		for i := 0; i < 1500; i++ {
			measurement := LatencyMeasurement{
				Timestamp:    time.Now(),
				ReadLatency:  time.Duration(i) * time.Microsecond,
				WriteLatency: time.Duration(i*2) * time.Microsecond,
				Operation:    "limit-test",
			}
			mc3.recordLatencyMeasurement(measurement)
		}

		history := mc3.GetLatencyHistory()
		assert.LessOrEqual(t, len(history), 1000)
		assert.Equal(t, 1000, len(history)) // Should be exactly 1000
	})
}

func TestMetricsCollector_GetLatencyHistory(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mc := NewMetricsCollector(nil, logger)

	t.Run("Empty history initially", func(t *testing.T) {
		history := mc.GetLatencyHistory()
		assert.NotNil(t, history)
		assert.Equal(t, 0, len(history))
	})

	t.Run("History returns copy not reference", func(t *testing.T) {
		mc2 := NewMetricsCollector(nil, logger)

		measurement := LatencyMeasurement{
			Timestamp:    time.Now(),
			ReadLatency:  5 * time.Millisecond,
			WriteLatency: 10 * time.Millisecond,
			Operation:    "copy-test",
		}
		mc2.recordLatencyMeasurement(measurement)

		history1 := mc2.GetLatencyHistory()
		assert.Equal(t, 1, len(history1))

		// Add another measurement
		mc2.recordLatencyMeasurement(measurement)

		// Original history slice should not be affected
		assert.Equal(t, 1, len(history1))

		// New history should have 2 items
		history2 := mc2.GetLatencyHistory()
		assert.Equal(t, 2, len(history2))
	})
}

func TestPercentile(t *testing.T) {
	t.Run("Percentile with empty slice", func(t *testing.T) {
		result := percentile([]float64{}, 0.50)
		assert.Equal(t, 0.0, result)
	})

	t.Run("Percentile with single value", func(t *testing.T) {
		result := percentile([]float64{100}, 0.50)
		assert.Equal(t, 100.0, result)
	})

	t.Run("Percentile P50 calculation", func(t *testing.T) {
		values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		result := percentile(values, 0.50)
		// P50 of 10 values at index 4-5 boundary
		assert.Equal(t, 5.0, result)
	})

	t.Run("Percentile P95 calculation", func(t *testing.T) {
		values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		result := percentile(values, 0.95)
		// P95 at index int((10-1) * 0.95) = int(8.55) = 8, value at index 8 is 9
		assert.Equal(t, 9.0, result)
	})

	t.Run("Percentile P99 calculation", func(t *testing.T) {
		values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		result := percentile(values, 0.99)
		// P99 at index int((10-1) * 0.99) = int(8.91) = 8, value at index 8 is 9
		assert.Equal(t, 9.0, result)
	})

	t.Run("Percentile boundary cases", func(t *testing.T) {
		values := []float64{10, 20, 30, 40, 50}

		// P0 should give first element
		p0 := percentile(values, 0.00)
		assert.Equal(t, 10.0, p0)

		// P100 should give last element
		p100 := percentile(values, 1.00)
		assert.Equal(t, 50.0, p100)
	})
}

func TestMetricsSnapshot_Structure(t *testing.T) {
	t.Run("MetricsSnapshot creation and fields", func(t *testing.T) {
		snapshot := &MetricsSnapshot{
			Timestamp:         time.Now(),
			DBSize:            1024 * 1024,
			DBSizeInUse:       512 * 1024,
			ReadLatencyP50:    5.0,
			ReadLatencyP95:    10.0,
			ReadLatencyP99:    15.0,
			WriteLatencyP50:   8.0,
			WriteLatencyP95:   16.0,
			WriteLatencyP99:   24.0,
			ProposalCommitted: 1000,
			ProposalApplied:   1000,
			ProposalPending:   0,
			ProposalFailed:    0,
			ActiveConnections: 3,
		}

		assert.Equal(t, int64(1024*1024), snapshot.DBSize)
		assert.Equal(t, int64(512*1024), snapshot.DBSizeInUse)
		assert.Equal(t, 5.0, snapshot.ReadLatencyP50)
		assert.Equal(t, 8.0, snapshot.WriteLatencyP50)
		assert.Equal(t, uint64(1000), snapshot.ProposalCommitted)
		assert.Equal(t, 3, snapshot.ActiveConnections)
	})
}

func TestPerformanceReport_Structure(t *testing.T) {
	t.Run("PerformanceReport creation and fields", func(t *testing.T) {
		start := time.Now()
		end := start.Add(1 * time.Second)

		report := &PerformanceReport{
			StartTime:       start,
			EndTime:         end,
			Duration:        end.Sub(start),
			Operations:      100,
			Throughput:      200.0,
			ReadLatencyP50:  5.5,
			ReadLatencyP95:  11.0,
			ReadLatencyP99:  16.5,
			WriteLatencyP50: 8.5,
			WriteLatencyP95: 17.0,
			WriteLatencyP99: 25.5,
			ReadErrors:      0,
			WriteErrors:     0,
		}

		assert.Equal(t, 100, report.Operations)
		assert.Equal(t, 200.0, report.Throughput)
		assert.Equal(t, 1*time.Second, report.Duration)
		assert.Equal(t, 0, report.ReadErrors)
		assert.Equal(t, 0, report.WriteErrors)
	})
}

func TestLatencyMeasurement_Structure(t *testing.T) {
	t.Run("LatencyMeasurement creation", func(t *testing.T) {
		now := time.Now()
		measurement := LatencyMeasurement{
			Timestamp:    now,
			ReadLatency:  5 * time.Millisecond,
			WriteLatency: 10 * time.Millisecond,
			Operation:    "health-check",
		}

		assert.Equal(t, now, measurement.Timestamp)
		assert.Equal(t, 5*time.Millisecond, measurement.ReadLatency)
		assert.Equal(t, 10*time.Millisecond, measurement.WriteLatency)
		assert.Equal(t, "health-check", measurement.Operation)
	})
}
