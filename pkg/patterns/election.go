package patterns

import (
	"context"
	"fmt"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/zap"
)

// LeaderElection provides distributed leader election
type LeaderElection struct {
	client    *clientv3.Client
	session   *concurrency.Session
	election  *concurrency.Election
	electionKey string
	nodeID    string
	ttl       int
	logger    *zap.Logger
	isLeader  bool
	mu        sync.RWMutex
	onElected func()
	onDefeat  func()
}

// NewLeaderElection creates a new leader election instance
func NewLeaderElection(client *clientv3.Client, electionKey string, nodeID string, ttl int, logger *zap.Logger) (*LeaderElection, error) {
	if ttl == 0 {
		ttl = 30
	}

	// Create session
	session, err := concurrency.NewSession(client, concurrency.WithTTL(ttl))
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Create election
	election := concurrency.NewElection(session, electionKey)

	return &LeaderElection{
		client:      client,
		session:     session,
		election:    election,
		electionKey: electionKey,
		nodeID:      nodeID,
		ttl:         ttl,
		logger:      logger,
		isLeader:    false,
	}, nil
}

// Campaign campaigns to become the leader
func (le *LeaderElection) Campaign(ctx context.Context) error {
	le.logger.Info("Starting election campaign", zap.String("node_id", le.nodeID))

	if err := le.election.Campaign(ctx, le.nodeID); err != nil {
		return fmt.Errorf("campaign failed: %w", err)
	}

	le.mu.Lock()
	le.isLeader = true
	le.mu.Unlock()

	le.logger.Info("Elected as leader", zap.String("node_id", le.nodeID))

	// Call elected callback
	if le.onElected != nil {
		go le.onElected()
	}

	return nil
}

// Resign gives up leadership
func (le *LeaderElection) Resign(ctx context.Context) error {
	le.mu.Lock()
	defer le.mu.Unlock()

	if !le.isLeader {
		return fmt.Errorf("not the leader")
	}

	if err := le.election.Resign(ctx); err != nil {
		return fmt.Errorf("resign failed: %w", err)
	}

	le.isLeader = false
	le.logger.Info("Resigned from leadership", zap.String("node_id", le.nodeID))

	// Call defeat callback
	if le.onDefeat != nil {
		go le.onDefeat()
	}

	return nil
}

// IsLeader returns whether this node is the current leader
func (le *LeaderElection) IsLeader() bool {
	le.mu.RLock()
	defer le.mu.RUnlock()
	return le.isLeader
}

// GetLeader returns the current leader
func (le *LeaderElection) GetLeader(ctx context.Context) (string, error) {
	resp, err := le.election.Leader(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get leader: %w", err)
	}

	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("no leader elected")
	}

	return string(resp.Kvs[0].Value), nil
}

// Observe watches for leader changes
func (le *LeaderElection) Observe(ctx context.Context) <-chan string {
	observeChan := make(chan string)

	go func() {
		defer close(observeChan)

		ch := le.election.Observe(ctx)
		for {
			select {
			case <-ctx.Done():
				return
			case resp, ok := <-ch:
				if !ok {
					return
				}
				if len(resp.Kvs) > 0 {
					leader := string(resp.Kvs[0].Value)
					observeChan <- leader

					// Check if we lost leadership
					le.mu.Lock()
					if le.isLeader && leader != le.nodeID {
						le.isLeader = false
						if le.onDefeat != nil {
							go le.onDefeat()
						}
						le.logger.Info("Lost leadership", zap.String("new_leader", leader))
					}
					le.mu.Unlock()
				}
			}
		}
	}()

	return observeChan
}

// OnElected sets callback for when node becomes leader
func (le *LeaderElection) OnElected(callback func()) {
	le.onElected = callback
}

// OnDefeat sets callback for when node loses leadership
func (le *LeaderElection) OnDefeat(callback func()) {
	le.onDefeat = callback
}

// Close closes the election and session
func (le *LeaderElection) Close() error {
	le.mu.Lock()
	defer le.mu.Unlock()

	if le.isLeader {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := le.election.Resign(ctx); err != nil {
			le.logger.Error("Failed to resign during close", zap.Error(err))
		}
		le.isLeader = false
	}

	return le.session.Close()
}

// LeaderElectionManager manages multiple leader elections
type LeaderElectionManager struct {
	client    *clientv3.Client
	elections map[string]*LeaderElection
	mu        sync.RWMutex
	logger    *zap.Logger
}

// NewLeaderElectionManager creates a new leader election manager
func NewLeaderElectionManager(client *clientv3.Client, logger *zap.Logger) *LeaderElectionManager {
	return &LeaderElectionManager{
		client:    client,
		elections: make(map[string]*LeaderElection),
		logger:    logger,
	}
}

// CreateElection creates a new election
func (lem *LeaderElectionManager) CreateElection(electionKey, nodeID string, ttl int) (*LeaderElection, error) {
	lem.mu.Lock()
	defer lem.mu.Unlock()

	if _, exists := lem.elections[electionKey]; exists {
		return nil, fmt.Errorf("election %s already exists", electionKey)
	}

	election, err := NewLeaderElection(lem.client, electionKey, nodeID, ttl, lem.logger)
	if err != nil {
		return nil, err
	}

	lem.elections[electionKey] = election
	return election, nil
}

// GetElection retrieves an election
func (lem *LeaderElectionManager) GetElection(electionKey string) (*LeaderElection, error) {
	lem.mu.RLock()
	defer lem.mu.RUnlock()

	election, exists := lem.elections[electionKey]
	if !exists {
		return nil, fmt.Errorf("election %s not found", electionKey)
	}

	return election, nil
}

// CloseElection closes and removes an election
func (lem *LeaderElectionManager) CloseElection(electionKey string) error {
	lem.mu.Lock()
	defer lem.mu.Unlock()

	election, exists := lem.elections[electionKey]
	if !exists {
		return fmt.Errorf("election %s not found", electionKey)
	}

	if err := election.Close(); err != nil {
		return err
	}

	delete(lem.elections, electionKey)
	return nil
}

// CloseAll closes all elections
func (lem *LeaderElectionManager) CloseAll() error {
	lem.mu.Lock()
	defer lem.mu.Unlock()

	var lastErr error
	for key, election := range lem.elections {
		if err := election.Close(); err != nil {
			lem.logger.Error("Failed to close election", zap.String("key", key), zap.Error(err))
			lastErr = err
		}
	}

	lem.elections = make(map[string]*LeaderElection)
	return lastErr
}

// LeaderTask represents a task that should only run on the leader
type LeaderTask struct {
	election *LeaderElection
	task     func(ctx context.Context) error
	interval time.Duration
	logger   *zap.Logger
	cancel   context.CancelFunc
}

// NewLeaderTask creates a new leader task
func NewLeaderTask(election *LeaderElection, task func(ctx context.Context) error, interval time.Duration, logger *zap.Logger) *LeaderTask {
	return &LeaderTask{
		election: election,
		task:     task,
		interval: interval,
		logger:   logger,
	}
}

// Start starts executing the task if this node is the leader
func (lt *LeaderTask) Start(ctx context.Context) {
	taskCtx, cancel := context.WithCancel(ctx)
	lt.cancel = cancel

	// Set up callbacks
	lt.election.OnElected(func() {
		lt.logger.Info("Starting leader task")
		go lt.runTask(taskCtx)
	})

	lt.election.OnDefeat(func() {
		lt.logger.Info("Stopping leader task")
		if lt.cancel != nil {
			lt.cancel()
		}
	})

	// If already leader, start immediately
	if lt.election.IsLeader() {
		go lt.runTask(taskCtx)
	}
}

// runTask runs the task periodically
func (lt *LeaderTask) runTask(ctx context.Context) {
	ticker := time.NewTicker(lt.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !lt.election.IsLeader() {
				return
			}

			if err := lt.task(ctx); err != nil {
				lt.logger.Error("Leader task failed", zap.Error(err))
			}
		}
	}
}

// Stop stops the task
func (lt *LeaderTask) Stop() {
	if lt.cancel != nil {
		lt.cancel()
	}
}

// HighAvailabilityCoordinator coordinates HA pairs
type HighAvailabilityCoordinator struct {
	election  *LeaderElection
	primary   func(ctx context.Context) error
	secondary func(ctx context.Context) error
	logger    *zap.Logger
	cancel    context.CancelFunc
}

// NewHighAvailabilityCoordinator creates a new HA coordinator
func NewHighAvailabilityCoordinator(
	election *LeaderElection,
	primary func(ctx context.Context) error,
	secondary func(ctx context.Context) error,
	logger *zap.Logger,
) *HighAvailabilityCoordinator {
	return &HighAvailabilityCoordinator{
		election:  election,
		primary:   primary,
		secondary: secondary,
		logger:    logger,
	}
}

// Start starts the HA coordinator
func (hac *HighAvailabilityCoordinator) Start(ctx context.Context) error {
	coordCtx, cancel := context.WithCancel(ctx)
	hac.cancel = cancel

	// Campaign for leadership
	go func() {
		if err := hac.election.Campaign(coordCtx); err != nil {
			hac.logger.Error("Campaign failed", zap.Error(err))
			return
		}
	}()

	// Observe leadership changes
	observeChan := hac.election.Observe(coordCtx)

	go func() {
		var currentCancel context.CancelFunc

		for {
			select {
			case <-coordCtx.Done():
				if currentCancel != nil {
					currentCancel()
				}
				return
			case leader, ok := <-observeChan:
				if !ok {
					return
				}

				// Cancel current role
				if currentCancel != nil {
					currentCancel()
				}

				// Start new role
				roleCtx, roleCancel := context.WithCancel(coordCtx)
				currentCancel = roleCancel

				if leader == hac.election.nodeID {
					hac.logger.Info("Running as primary")
					go hac.runRole(roleCtx, hac.primary)
				} else {
					hac.logger.Info("Running as secondary")
					go hac.runRole(roleCtx, hac.secondary)
				}
			}
		}
	}()

	return nil
}

// runRole runs the role function
func (hac *HighAvailabilityCoordinator) runRole(ctx context.Context, roleFn func(ctx context.Context) error) {
	if err := roleFn(ctx); err != nil {
		hac.logger.Error("Role function failed", zap.Error(err))
	}
}

// Stop stops the coordinator
func (hac *HighAvailabilityCoordinator) Stop() error {
	if hac.cancel != nil {
		hac.cancel()
	}
	return hac.election.Close()
}
