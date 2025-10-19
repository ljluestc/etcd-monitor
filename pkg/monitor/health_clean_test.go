package monitor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHealthChecker_LeaderHistory(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	hc := NewHealthChecker(nil, logger)

	t.Run("Empty leader history initially", func(t *testing.T) {
		history := hc.GetLeaderHistory()
		assert.NotNil(t, history)
		assert.Equal(t, 0, len(history))
	})

	t.Run("Leader history after recording changes", func(t *testing.T) {
		hc.recordLeaderChange(123, 456)

		history := hc.GetLeaderHistory()
		assert.Equal(t, 1, len(history))
		assert.Equal(t, uint64(123), history[0].OldLeaderID)
		assert.Equal(t, uint64(456), history[0].NewLeaderID)
	})

	t.Run("Leader history with multiple changes", func(t *testing.T) {
		hc2 := NewHealthChecker(nil, logger)
		hc2.recordLeaderChange(1, 2)
		hc2.recordLeaderChange(2, 3)
		hc2.recordLeaderChange(3, 1)

		history := hc2.GetLeaderHistory()
		assert.Equal(t, 3, len(history))
	})

	t.Run("Leader history respects max limit", func(t *testing.T) {
		hc3 := NewHealthChecker(nil, logger)

		// Record many leader changes to test history limit
		for i := 0; i < 150; i++ {
			hc3.recordLeaderChange(uint64(i), uint64(i+1))
		}

		history := hc3.GetLeaderHistory()

		// Should be limited to maxHistory (100)
		assert.LessOrEqual(t, len(history), 100)
	})
}

func TestClusterStatus_Basic(t *testing.T) {
	t.Run("Cluster status creation", func(t *testing.T) {
		status := &ClusterStatus{
			Healthy:          true,
			LeaderID:         1234,
			MemberCount:      3,
			QuorumSize:       2,
			HasLeader:        true,
			LeaderChanges:    0,
			LastLeaderChange: time.Now(),
			NetworkPartition: false,
			Alarms:           []AlarmInfo{},
			LastCheck:        time.Now(),
		}

		assert.True(t, status.Healthy)
		assert.Equal(t, 3, status.MemberCount)
		assert.Equal(t, 2, status.QuorumSize)
		assert.Equal(t, uint64(1234), status.LeaderID)
	})
}

func TestMemberInfo_Basic(t *testing.T) {
	t.Run("Member info structure", func(t *testing.T) {
		member := MemberInfo{
			ID:         123,
			Name:       "etcd-1",
			PeerURLs:   []string{"http://localhost:2380"},
			ClientURLs: []string{"http://localhost:2379"},
			IsLeader:   true,
			IsHealthy:  true,
			DBSize:     1024,
			Version:    "3.5.0",
		}

		assert.Equal(t, uint64(123), member.ID)
		assert.Equal(t, "etcd-1", member.Name)
		assert.True(t, member.IsLeader)
		assert.True(t, member.IsHealthy)
		assert.Equal(t, int64(1024), member.DBSize)
		assert.Equal(t, "3.5.0", member.Version)
	})
}

func TestLeaderChange_Basic(t *testing.T) {
	t.Run("Leader change creation", func(t *testing.T) {
		lc := LeaderChange{
			OldLeaderID: 100,
			NewLeaderID: 200,
			Timestamp:   time.Now(),
		}

		assert.Equal(t, uint64(100), lc.OldLeaderID)
		assert.Equal(t, uint64(200), lc.NewLeaderID)
		assert.WithinDuration(t, time.Now(), lc.Timestamp, 1*time.Second)
	})
}
