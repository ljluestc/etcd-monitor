package monitor

import (
	"context"
	"fmt"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// HealthChecker performs health checks on etcd clusters
type HealthChecker struct {
	client        *clientv3.Client
	logger        *zap.Logger
	mu            sync.RWMutex
	leaderHistory []LeaderChange
	maxHistory    int
}

// LeaderChange records a leader change event
type LeaderChange struct {
	Timestamp  time.Time
	OldLeaderID uint64
	NewLeaderID uint64
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(client *clientv3.Client, logger *zap.Logger) *HealthChecker {
	return &HealthChecker{
		client:        client,
		logger:        logger,
		leaderHistory: make([]LeaderChange, 0),
		maxHistory:    100,
	}
}

// CheckClusterHealth performs a comprehensive health check
func (hc *HealthChecker) CheckClusterHealth(ctx context.Context) (*ClusterStatus, error) {
	status := &ClusterStatus{
		LastCheck: time.Now(),
		Healthy:   true,
	}

	// Get member list
	membersResp, err := hc.client.MemberList(ctx)
	if err != nil {
		status.Healthy = false
		return status, fmt.Errorf("failed to get member list: %w", err)
	}

	status.MemberCount = len(membersResp.Members)
	status.QuorumSize = (status.MemberCount / 2) + 1

	// Check each member's health
	healthyMembers := 0
	var currentLeaderID uint64

	for _, member := range membersResp.Members {
		memberHealthy, isLeader, err := hc.checkMemberHealth(ctx, member)
		if err != nil {
			hc.logger.Warn("Failed to check member health",
				zap.Uint64("member_id", member.ID),
				zap.Error(err))
			continue
		}

		if memberHealthy {
			healthyMembers++
		}

		if isLeader {
			currentLeaderID = member.ID
			status.HasLeader = true
		}
	}

	// Check quorum
	if healthyMembers < status.QuorumSize {
		status.Healthy = false
		status.NetworkPartition = true
	}

	// Check leader changes
	if currentLeaderID != 0 {
		hc.recordLeaderChange(status.LeaderID, currentLeaderID)
		status.LeaderID = currentLeaderID
	} else {
		status.HasLeader = false
		status.Healthy = false
	}

	// Count recent leader changes
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	recentChanges := 0
	for _, change := range hc.leaderHistory {
		if change.Timestamp.After(oneHourAgo) {
			recentChanges++
		}
	}
	status.LeaderChanges = recentChanges

	// Get the last leader change time
	if len(hc.leaderHistory) > 0 {
		status.LastLeaderChange = hc.leaderHistory[len(hc.leaderHistory)-1].Timestamp
	}

	// Check for alarms
	alarmResp, err := hc.client.AlarmList(ctx)
	if err != nil {
		hc.logger.Warn("Failed to get alarm list", zap.Error(err))
	} else {
		for _, alarm := range alarmResp.Alarms {
			status.Alarms = append(status.Alarms, AlarmInfo{
				Type:      alarm.Alarm.String(),
				MemberID:  alarm.MemberID,
				Triggered: time.Now(),
			})
		}
		if len(status.Alarms) > 0 {
			status.Healthy = false
		}
	}

	return status, nil
}

// checkMemberHealth checks if a specific member is healthy
func (hc *HealthChecker) checkMemberHealth(ctx context.Context, member *clientv3.Member) (healthy bool, isLeader bool, err error) {
	// Use status endpoint to check health
	statusResp, err := hc.client.Status(ctx, member.ClientURLs[0])
	if err != nil {
		return false, false, err
	}

	// Member is healthy if we can get its status
	healthy = true

	// Check if this member is the leader
	isLeader = statusResp.Leader == member.ID

	return healthy, isLeader, nil
}

// recordLeaderChange records a leader change event
func (hc *HealthChecker) recordLeaderChange(oldLeaderID, newLeaderID uint64) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	// Only record if there's an actual change
	if oldLeaderID == 0 || oldLeaderID == newLeaderID {
		return
	}

	change := LeaderChange{
		Timestamp:   time.Now(),
		OldLeaderID: oldLeaderID,
		NewLeaderID: newLeaderID,
	}

	hc.leaderHistory = append(hc.leaderHistory, change)

	// Keep only the last maxHistory changes
	if len(hc.leaderHistory) > hc.maxHistory {
		hc.leaderHistory = hc.leaderHistory[1:]
	}

	hc.logger.Info("Leader change detected",
		zap.Uint64("old_leader", oldLeaderID),
		zap.Uint64("new_leader", newLeaderID))
}

// GetLeaderHistory returns the leader change history
func (hc *HealthChecker) GetLeaderHistory() []LeaderChange {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	// Return a copy
	history := make([]LeaderChange, len(hc.leaderHistory))
	copy(history, hc.leaderHistory)
	return history
}

// CheckEndpointHealth checks the health of a specific endpoint
func (hc *HealthChecker) CheckEndpointHealth(ctx context.Context, endpoint string) (bool, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Try to get status from the endpoint
	_, err := hc.client.Status(timeoutCtx, endpoint)
	if err != nil {
		return false, err
	}

	return true, nil
}

// DetectSplitBrain attempts to detect split-brain scenarios
func (hc *HealthChecker) DetectSplitBrain(ctx context.Context) (bool, error) {
	membersResp, err := hc.client.MemberList(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get member list: %w", err)
	}

	leaderCount := 0
	var leaders []uint64

	// Check each member to see if it thinks it's the leader
	for _, member := range membersResp.Members {
		if len(member.ClientURLs) == 0 {
			continue
		}

		statusResp, err := hc.client.Status(ctx, member.ClientURLs[0])
		if err != nil {
			hc.logger.Warn("Failed to get member status",
				zap.Uint64("member_id", member.ID),
				zap.Error(err))
			continue
		}

		if statusResp.Leader == member.ID {
			leaderCount++
			leaders = append(leaders, member.ID)
		}
	}

	// Split-brain: multiple members think they're the leader
	if leaderCount > 1 {
		hc.logger.Error("Split-brain detected",
			zap.Int("leader_count", leaderCount),
			zap.Uint64s("leaders", leaders))
		return true, nil
	}

	return false, nil
}

// CheckNetworkLatency measures network latency to each member
func (hc *HealthChecker) CheckNetworkLatency(ctx context.Context) (map[uint64]time.Duration, error) {
	latencies := make(map[uint64]time.Duration)

	membersResp, err := hc.client.MemberList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get member list: %w", err)
	}

	for _, member := range membersResp.Members {
		if len(member.ClientURLs) == 0 {
			continue
		}

		start := time.Now()
		_, err := hc.client.Status(ctx, member.ClientURLs[0])
		duration := time.Since(start)

		if err != nil {
			hc.logger.Warn("Failed to measure latency",
				zap.Uint64("member_id", member.ID),
				zap.Error(err))
			continue
		}

		latencies[member.ID] = duration
	}

	return latencies, nil
}

// PerformHealthCheck performs a quick health check and returns a boolean
func (hc *HealthChecker) PerformHealthCheck(ctx context.Context) bool {
	status, err := hc.CheckClusterHealth(ctx)
	if err != nil {
		return false
	}
	return status.Healthy
}
