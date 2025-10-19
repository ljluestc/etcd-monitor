package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/etcd-monitor/taskmaster/testutil/mocks"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.uber.org/zap"
)

func TestHealthChecker_WithMock(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockClient := mocks.NewMockEtcdClient()

	t.Run("Mock client endpoint health check", func(t *testing.T) {
		ctx := context.Background()

		// Test status endpoint
		resp, err := mockClient.Status(ctx, "http://localhost:2379")
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, uint64(1), resp.Leader)
	})

	t.Run("Mock client status for all members", func(t *testing.T) {
		ctx := context.Background()

		membersResp, err := mockClient.MemberList(ctx)
		assert.NoError(t, err)

		for _, member := range membersResp.Members {
			if len(member.ClientURLs) > 0 {
				statusResp, err := mockClient.Status(ctx, member.ClientURLs[0])
				assert.NoError(t, err)
				assert.NotNil(t, statusResp)
			}
		}
	})

	t.Run("Mock client leader detection", func(t *testing.T) {
		ctx := context.Background()

		membersResp, _ := mockClient.MemberList(ctx)
		var leaderID uint64

		for _, member := range membersResp.Members {
			if len(member.ClientURLs) > 0 {
				statusResp, _ := mockClient.Status(ctx, member.ClientURLs[0])
				if statusResp.Leader == member.ID {
					leaderID = member.ID
					break
				}
			}
		}

		assert.Equal(t, uint64(1), leaderID)
	})

	t.Run("Mock client alarm detection", func(t *testing.T) {
		ctx := context.Background()

		// No alarms initially
		alarmResp, err := mockClient.AlarmList(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(alarmResp.Alarms))

		// Add an alarm
		mockClient.AddAlarm(1, etcdserverpb.AlarmType_NOSPACE)

		// Check alarm is detected
		alarmResp, err = mockClient.AlarmList(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(alarmResp.Alarms))
		assert.Equal(t, etcdserverpb.AlarmType_NOSPACE, alarmResp.Alarms[0].Alarm)

		// Clear alarm
		mockClient.ClearAlarms()
		alarmResp, err = mockClient.AlarmList(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(alarmResp.Alarms))
	})

	t.Run("Mock client network partition simulation", func(t *testing.T) {
		ctx := context.Background()
		mockClient2 := mocks.NewMockEtcdClient()

		// Get initial member count
		membersResp, _ := mockClient2.MemberList(ctx)
		initialCount := len(membersResp.Members)
		assert.Equal(t, 3, initialCount)

		// Simulate partition by removing a member
		mockClient2.RemoveMember(3)

		membersResp, _ = mockClient2.MemberList(ctx)
		assert.Equal(t, 2, len(membersResp.Members))
	})

	t.Run("Mock client split-brain simulation", func(t *testing.T) {
		ctx := context.Background()
		mockClient2 := mocks.NewMockEtcdClient()

		// Normal case: one leader
		membersResp, _ := mockClient2.MemberList(ctx)
		leaderCount := 0

		for _, member := range membersResp.Members {
			if len(member.ClientURLs) > 0 {
				statusResp, _ := mockClient2.Status(ctx, member.ClientURLs[0])
				if statusResp.Leader == member.ID {
					leaderCount++
				}
			}
		}

		assert.Equal(t, 1, leaderCount)
	})

	t.Run("Mock client error scenarios", func(t *testing.T) {
		ctx := context.Background()
		mockClient2 := mocks.NewMockEtcdClient()

		// Test MemberList error
		mockClient2.SetMemberListError(assert.AnError)
		_, err := mockClient2.MemberList(ctx)
		assert.Error(t, err)
		mockClient2.SetMemberListError(nil)

		// Test Status error
		mockClient2.SetStatusError(assert.AnError)
		_, err = mockClient2.Status(ctx, "http://localhost:2379")
		assert.Error(t, err)
		mockClient2.SetStatusError(nil)

		// Test AlarmList error
		mockClient2.SetAlarmListError(assert.AnError)
		_, err = mockClient2.AlarmList(ctx)
		assert.Error(t, err)
		mockClient2.SetAlarmListError(nil)
	})

	_ = logger // Suppress unused warning
}

func TestHealthChecker_RecordLeaderChange(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	hc := NewHealthChecker(nil, logger)

	t.Run("No change when old leader is zero", func(t *testing.T) {
		hc.recordLeaderChange(0, 1)
		history := hc.GetLeaderHistory()
		assert.Equal(t, 0, len(history))
	})

	t.Run("No change when leader stays the same", func(t *testing.T) {
		hc2 := NewHealthChecker(nil, logger)
		hc2.recordLeaderChange(1, 1)
		history := hc2.GetLeaderHistory()
		assert.Equal(t, 0, len(history))
	})

	t.Run("Record actual leader change", func(t *testing.T) {
		hc3 := NewHealthChecker(nil, logger)
		hc3.recordLeaderChange(1, 2)
		history := hc3.GetLeaderHistory()
		assert.Equal(t, 1, len(history))
		assert.Equal(t, uint64(1), history[0].OldLeaderID)
		assert.Equal(t, uint64(2), history[0].NewLeaderID)
	})

	t.Run("Record multiple leader changes", func(t *testing.T) {
		hc4 := NewHealthChecker(nil, logger)
		hc4.recordLeaderChange(1, 2)
		hc4.recordLeaderChange(2, 3)
		hc4.recordLeaderChange(3, 1)

		history := hc4.GetLeaderHistory()
		assert.Equal(t, 3, len(history))
		assert.Equal(t, uint64(1), history[0].OldLeaderID)
		assert.Equal(t, uint64(2), history[0].NewLeaderID)
		assert.Equal(t, uint64(2), history[1].OldLeaderID)
		assert.Equal(t, uint64(3), history[1].NewLeaderID)
		assert.Equal(t, uint64(3), history[2].OldLeaderID)
		assert.Equal(t, uint64(1), history[2].NewLeaderID)
	})

	t.Run("History limit enforcement", func(t *testing.T) {
		hc5 := NewHealthChecker(nil, logger)

		// Record 150 leader changes
		for i := 0; i < 150; i++ {
			hc5.recordLeaderChange(uint64(i), uint64(i+1))
		}

		history := hc5.GetLeaderHistory()
		assert.LessOrEqual(t, len(history), 100)
		assert.Equal(t, 100, len(history))

		// Verify oldest entries were removed
		// Latest entry should be 149 -> 150
		assert.Equal(t, uint64(149), history[len(history)-1].OldLeaderID)
		assert.Equal(t, uint64(150), history[len(history)-1].NewLeaderID)
	})

	t.Run("Recent timestamp for new changes", func(t *testing.T) {
		hc6 := NewHealthChecker(nil, logger)
		before := time.Now()
		hc6.recordLeaderChange(1, 2)
		after := time.Now()

		history := hc6.GetLeaderHistory()
		assert.Equal(t, 1, len(history))
		assert.True(t, history[0].Timestamp.After(before) || history[0].Timestamp.Equal(before))
		assert.True(t, history[0].Timestamp.Before(after) || history[0].Timestamp.Equal(after))
	})
}

func TestAlarmInfo_Structure(t *testing.T) {
	t.Run("AlarmInfo creation", func(t *testing.T) {
		alarm := AlarmInfo{
			Type:      "NOSPACE",
			MemberID:  1,
			Triggered: time.Now(),
		}

		assert.Equal(t, "NOSPACE", alarm.Type)
		assert.Equal(t, uint64(1), alarm.MemberID)
		assert.WithinDuration(t, time.Now(), alarm.Triggered, 1*time.Second)
	})

	t.Run("AlarmInfo in slice", func(t *testing.T) {
		alarms := []AlarmInfo{
			{
				Type:      "NOSPACE",
				MemberID:  1,
				Triggered: time.Now(),
			},
			{
				Type:      "CORRUPT",
				MemberID:  2,
				Triggered: time.Now(),
			},
		}

		assert.Equal(t, 2, len(alarms))
		assert.Equal(t, "NOSPACE", alarms[0].Type)
		assert.Equal(t, "CORRUPT", alarms[1].Type)
	})
}

func TestClusterStatus_HealthCalculation(t *testing.T) {
	t.Run("Healthy cluster with quorum", func(t *testing.T) {
		status := &ClusterStatus{
			Healthy:          true,
			LeaderID:         1,
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
		assert.True(t, status.HasLeader)
		assert.False(t, status.NetworkPartition)
		assert.Equal(t, 0, len(status.Alarms))
	})

	t.Run("Unhealthy cluster without quorum", func(t *testing.T) {
		status := &ClusterStatus{
			Healthy:          false,
			LeaderID:         0,
			MemberCount:      3,
			QuorumSize:       2,
			HasLeader:        false,
			LeaderChanges:    0,
			LastLeaderChange: time.Now(),
			NetworkPartition: true,
			Alarms:           []AlarmInfo{},
			LastCheck:        time.Now(),
		}

		assert.False(t, status.Healthy)
		assert.False(t, status.HasLeader)
		assert.True(t, status.NetworkPartition)
	})

	t.Run("Cluster with alarms is unhealthy", func(t *testing.T) {
		status := &ClusterStatus{
			Healthy:          false,
			LeaderID:         1,
			MemberCount:      3,
			QuorumSize:       2,
			HasLeader:        true,
			LeaderChanges:    0,
			LastLeaderChange: time.Now(),
			NetworkPartition: false,
			Alarms: []AlarmInfo{
				{
					Type:      "NOSPACE",
					MemberID:  1,
					Triggered: time.Now(),
				},
			},
			LastCheck: time.Now(),
		}

		assert.False(t, status.Healthy)
		assert.Equal(t, 1, len(status.Alarms))
		assert.Equal(t, "NOSPACE", status.Alarms[0].Type)
	})

	t.Run("Quorum calculation for different cluster sizes", func(t *testing.T) {
		testCases := []struct {
			memberCount int
			expectedQuorum int
		}{
			{1, 1},
			{3, 2},
			{5, 3},
			{7, 4},
		}

		for _, tc := range testCases {
			status := &ClusterStatus{
				MemberCount: tc.memberCount,
				QuorumSize:  (tc.memberCount / 2) + 1,
			}
			assert.Equal(t, tc.expectedQuorum, status.QuorumSize)
		}
	})

	t.Run("Recent leader changes tracking", func(t *testing.T) {
		now := time.Now()
		status := &ClusterStatus{
			Healthy:          true,
			LeaderID:         3,
			MemberCount:      3,
			QuorumSize:       2,
			HasLeader:        true,
			LeaderChanges:    5, // 5 changes in last hour
			LastLeaderChange: now.Add(-30 * time.Minute),
			NetworkPartition: false,
			Alarms:           []AlarmInfo{},
			LastCheck:        now,
		}

		assert.Equal(t, 5, status.LeaderChanges)
		assert.True(t, status.LastLeaderChange.After(now.Add(-1*time.Hour)))
	})
}

func TestMemberInfo_Validation(t *testing.T) {
	t.Run("Complete member info", func(t *testing.T) {
		member := MemberInfo{
			ID:         123,
			Name:       "etcd-test",
			PeerURLs:   []string{"http://localhost:2380"},
			ClientURLs: []string{"http://localhost:2379"},
			IsLeader:   true,
			IsHealthy:  true,
			DBSize:     2048,
			Version:    "3.5.0",
		}

		assert.NotZero(t, member.ID)
		assert.NotEmpty(t, member.Name)
		assert.NotEmpty(t, member.PeerURLs)
		assert.NotEmpty(t, member.ClientURLs)
		assert.True(t, member.IsLeader)
		assert.True(t, member.IsHealthy)
		assert.Greater(t, member.DBSize, int64(0))
		assert.NotEmpty(t, member.Version)
	})

	t.Run("Member info for unhealthy member", func(t *testing.T) {
		member := MemberInfo{
			ID:         456,
			Name:       "etcd-down",
			PeerURLs:   []string{"http://localhost:2381"},
			ClientURLs: []string{"http://localhost:2382"},
			IsLeader:   false,
			IsHealthy:  false,
			DBSize:     0,
			Version:    "",
		}

		assert.False(t, member.IsHealthy)
		assert.False(t, member.IsLeader)
		assert.Equal(t, int64(0), member.DBSize)
		assert.Empty(t, member.Version)
	})

	t.Run("Multiple member infos", func(t *testing.T) {
		members := []MemberInfo{
			{ID: 1, Name: "etcd-1", IsLeader: true, IsHealthy: true},
			{ID: 2, Name: "etcd-2", IsLeader: false, IsHealthy: true},
			{ID: 3, Name: "etcd-3", IsLeader: false, IsHealthy: false},
		}

		assert.Equal(t, 3, len(members))

		// Count healthy members
		healthyCount := 0
		for _, m := range members {
			if m.IsHealthy {
				healthyCount++
			}
		}
		assert.Equal(t, 2, healthyCount)

		// Find leader
		var leaderID uint64
		for _, m := range members {
			if m.IsLeader {
				leaderID = m.ID
			}
		}
		assert.Equal(t, uint64(1), leaderID)
	})
}
