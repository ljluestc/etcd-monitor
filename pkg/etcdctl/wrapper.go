// Package etcdctl provides a wrapper for etcdctl operations
package etcdctl

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// Wrapper provides a high-level interface for etcdctl-like operations
type Wrapper struct {
	client *clientv3.Client
	logger *zap.Logger
}

// NewWrapper creates a new etcdctl wrapper
func NewWrapper(client *clientv3.Client, logger *zap.Logger) *Wrapper {
	return &Wrapper{
		client: client,
		logger: logger,
	}
}

// Put stores a key-value pair
func (w *Wrapper) Put(ctx context.Context, key, value string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	w.logger.Debug("Put operation", zap.String("key", key))
	return w.client.Put(ctx, key, value, opts...)
}

// Get retrieves a value by key
func (w *Wrapper) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	w.logger.Debug("Get operation", zap.String("key", key))
	return w.client.Get(ctx, key, opts...)
}

// Delete removes a key
func (w *Wrapper) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	w.logger.Debug("Delete operation", zap.String("key", key))
	return w.client.Delete(ctx, key, opts...)
}

// Watch watches for changes on a key
func (w *Wrapper) Watch(ctx context.Context, key string, opts ...clientv3.OpOption) clientv3.WatchChan {
	w.logger.Debug("Watch operation", zap.String("key", key))
	return w.client.Watch(ctx, key, opts...)
}

// MemberList returns the list of cluster members
func (w *Wrapper) MemberList(ctx context.Context) (*clientv3.MemberListResponse, error) {
	w.logger.Debug("Member list operation")
	return w.client.MemberList(ctx)
}

// MemberAdd adds a new member to the cluster
func (w *Wrapper) MemberAdd(ctx context.Context, peerURLs []string) (*clientv3.MemberAddResponse, error) {
	w.logger.Info("Adding member", zap.Strings("peer_urls", peerURLs))
	return w.client.MemberAdd(ctx, peerURLs)
}

// MemberRemove removes a member from the cluster
func (w *Wrapper) MemberRemove(ctx context.Context, memberID uint64) (*clientv3.MemberRemoveResponse, error) {
	w.logger.Info("Removing member", zap.Uint64("member_id", memberID))
	return w.client.MemberRemove(ctx, memberID)
}

// EndpointHealth checks the health of cluster endpoints
func (w *Wrapper) EndpointHealth(ctx context.Context) ([]EndpointHealth, error) {
	w.logger.Debug("Checking endpoint health")

	endpoints := w.client.Endpoints()
	results := make([]EndpointHealth, 0, len(endpoints))

	for _, endpoint := range endpoints {
		start := time.Now()
		_, err := w.client.Status(ctx, endpoint)
		latency := time.Since(start)

		health := EndpointHealth{
			Endpoint: endpoint,
			Healthy:  err == nil,
			Latency:  latency,
		}

		if err != nil {
			health.Error = err.Error()
		}

		results = append(results, health)
	}

	return results, nil
}

// EndpointHealth represents the health of an endpoint
type EndpointHealth struct {
	Endpoint string
	Healthy  bool
	Latency  time.Duration
	Error    string
}

// EndpointStatus returns the status of cluster endpoints
func (w *Wrapper) EndpointStatus(ctx context.Context) ([]EndpointStatus, error) {
	w.logger.Debug("Getting endpoint status")

	endpoints := w.client.Endpoints()
	results := make([]EndpointStatus, 0, len(endpoints))

	for _, endpoint := range endpoints {
		statusResp, err := w.client.Status(ctx, endpoint)
		if err != nil {
			w.logger.Warn("Failed to get endpoint status",
				zap.String("endpoint", endpoint),
				zap.Error(err))
			continue
		}

		status := EndpointStatus{
			Endpoint:        endpoint,
			Version:         statusResp.Version,
			DBSize:          statusResp.DbSize,
			DBSizeInUse:     statusResp.DbSizeInUse,
			Leader:          statusResp.Leader,
			RaftIndex:       statusResp.RaftIndex,
			RaftTerm:        statusResp.RaftTerm,
			RaftAppliedIndex: statusResp.RaftAppliedIndex,
		}

		results = append(results, status)
	}

	return results, nil
}

// EndpointStatus represents the status of an endpoint
type EndpointStatus struct {
	Endpoint         string
	Version          string
	DBSize           int64
	DBSizeInUse      int64
	Leader           uint64
	RaftIndex        uint64
	RaftTerm         uint64
	RaftAppliedIndex uint64
}

// Compact compacts the etcd keyspace
func (w *Wrapper) Compact(ctx context.Context, rev int64, opts ...clientv3.CompactOption) (*clientv3.CompactResponse, error) {
	w.logger.Info("Compacting keyspace", zap.Int64("revision", rev))
	return w.client.Compact(ctx, rev, opts...)
}

// Defragment defragments the etcd database
func (w *Wrapper) Defragment(ctx context.Context, endpoint string) error {
	w.logger.Info("Defragmenting", zap.String("endpoint", endpoint))

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	_, err := w.client.Defragment(timeoutCtx, endpoint)
	return err
}

// SnapshotSave saves a snapshot of the etcd database
func (w *Wrapper) SnapshotSave(ctx context.Context, path string) error {
	w.logger.Info("Saving snapshot", zap.String("path", path))

	// Get snapshot from one of the endpoints
	endpoints := w.client.Endpoints()
	if len(endpoints) == 0 {
		return fmt.Errorf("no endpoints available")
	}

	// Use the first endpoint
	// Note: This is a simplified implementation for testing
	// In a real implementation, you would need to create a proper snapshot
	return fmt.Errorf("snapshot save not implemented in test environment")
}

// AlarmList lists all alarms
func (w *Wrapper) AlarmList(ctx context.Context) (*clientv3.AlarmResponse, error) {
	w.logger.Debug("Listing alarms")
	return w.client.AlarmList(ctx)
}

// AlarmDisarm disarms an alarm
func (w *Wrapper) AlarmDisarm(ctx context.Context, alarm *clientv3.AlarmMember) (*clientv3.AlarmResponse, error) {
	w.logger.Info("Disarming alarm",
		zap.String("type", alarm.Alarm.String()),
		zap.Uint64("member_id", alarm.MemberID))
	return w.client.AlarmDisarm(ctx, alarm)
}

// Grant creates a new lease
func (w *Wrapper) Grant(ctx context.Context, ttl int64) (*clientv3.LeaseGrantResponse, error) {
	w.logger.Debug("Granting lease", zap.Int64("ttl", ttl))
	return w.client.Grant(ctx, ttl)
}

// Revoke revokes a lease
func (w *Wrapper) Revoke(ctx context.Context, leaseID clientv3.LeaseID) (*clientv3.LeaseRevokeResponse, error) {
	w.logger.Debug("Revoking lease", zap.Int64("lease_id", int64(leaseID)))
	return w.client.Revoke(ctx, leaseID)
}

// KeepAlive keeps a lease alive
func (w *Wrapper) KeepAlive(ctx context.Context, leaseID clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	w.logger.Debug("Starting lease keep-alive", zap.Int64("lease_id", int64(leaseID)))
	return w.client.KeepAlive(ctx, leaseID)
}

// TimeToLive returns the remaining TTL of a lease
func (w *Wrapper) TimeToLive(ctx context.Context, leaseID clientv3.LeaseID, opts ...clientv3.LeaseOption) (*clientv3.LeaseTimeToLiveResponse, error) {
	w.logger.Debug("Getting lease TTL", zap.Int64("lease_id", int64(leaseID)))
	return w.client.TimeToLive(ctx, leaseID, opts...)
}

// Transaction executes a transaction
func (w *Wrapper) Transaction(ctx context.Context) clientv3.Txn {
	w.logger.Debug("Creating transaction")
	return w.client.Txn(ctx)
}

// GetWithPrefix retrieves all keys with a given prefix
func (w *Wrapper) GetWithPrefix(ctx context.Context, prefix string) (*clientv3.GetResponse, error) {
	w.logger.Debug("Get with prefix", zap.String("prefix", prefix))
	return w.client.Get(ctx, prefix, clientv3.WithPrefix())
}

// DeleteWithPrefix deletes all keys with a given prefix
func (w *Wrapper) DeleteWithPrefix(ctx context.Context, prefix string) (*clientv3.DeleteResponse, error) {
	w.logger.Debug("Delete with prefix", zap.String("prefix", prefix))
	return w.client.Delete(ctx, prefix, clientv3.WithPrefix())
}

// GetRange retrieves a range of keys
func (w *Wrapper) GetRange(ctx context.Context, start, end string) (*clientv3.GetResponse, error) {
	w.logger.Debug("Get range", zap.String("start", start), zap.String("end", end))
	return w.client.Get(ctx, start, clientv3.WithRange(end))
}

// PutWithLease stores a key-value pair with a lease
func (w *Wrapper) PutWithLease(ctx context.Context, key, value string, leaseID clientv3.LeaseID) (*clientv3.PutResponse, error) {
	w.logger.Debug("Put with lease",
		zap.String("key", key),
		zap.Int64("lease_id", int64(leaseID)))
	return w.client.Put(ctx, key, value, clientv3.WithLease(leaseID))
}

// CompareAndSwap performs a compare-and-swap operation
func (w *Wrapper) CompareAndSwap(ctx context.Context, key, oldValue, newValue string) (bool, error) {
	w.logger.Debug("Compare-and-swap", zap.String("key", key))

	txn := w.client.Txn(ctx)
	cmp := clientv3.Compare(clientv3.Value(key), "=", oldValue)
	putOp := clientv3.OpPut(key, newValue)
	getOp := clientv3.OpGet(key)

	resp, err := txn.If(cmp).Then(putOp).Else(getOp).Commit()
	if err != nil {
		return false, err
	}

	return resp.Succeeded, nil
}

// ExecuteCommand executes an etcdctl-like command
func (w *Wrapper) ExecuteCommand(ctx context.Context, command string, args []string) (interface{}, error) {
	switch command {
	case "put":
		if len(args) < 2 {
			return nil, fmt.Errorf("put requires key and value")
		}
		return w.Put(ctx, args[0], args[1])

	case "get":
		if len(args) < 1 {
			return nil, fmt.Errorf("get requires key")
		}
		return w.Get(ctx, args[0])

	case "del", "delete":
		if len(args) < 1 {
			return nil, fmt.Errorf("delete requires key")
		}
		return w.Delete(ctx, args[0])

	case "member list":
		return w.MemberList(ctx)

	case "endpoint health":
		return w.EndpointHealth(ctx)

	case "endpoint status":
		return w.EndpointStatus(ctx)

	case "alarm list":
		return w.AlarmList(ctx)

	case "compact":
		if len(args) < 1 {
			return nil, fmt.Errorf("compact requires revision")
		}
		var rev int64
		fmt.Sscanf(args[0], "%d", &rev)
		return w.Compact(ctx, rev)

	default:
		return nil, fmt.Errorf("unknown command: %s", command)
	}
}
