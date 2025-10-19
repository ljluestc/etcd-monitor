package mocks

import (
	"context"
	"fmt"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

// MockEtcdClient is a mock implementation of etcd client for testing
type MockEtcdClient struct {
	// KV operations
	kvStore map[string]string
	kvMutex sync.RWMutex

	// Cluster state
	members     []*etcdserverpb.Member
	leaderID    uint64
	alarms      []*etcdserverpb.AlarmMember
	memberMutex sync.RWMutex

	// Mock responses
	memberListError error
	statusError     error
	getError        error
	putError        error
	deleteError     error
	alarmListError  error

	// Latency simulation
	simulateLatency bool
	readLatency     time.Duration
	writeLatency    time.Duration

	// Operation counters
	getCount    int64
	putCount    int64
	deleteCount int64
	opMutex     sync.Mutex
}

// NewMockEtcdClient creates a new mock etcd client with default cluster setup
func NewMockEtcdClient() *MockEtcdClient {
	return &MockEtcdClient{
		kvStore: make(map[string]string),
		members: []*etcdserverpb.Member{
			{
				ID:   1,
				Name: "etcd-1",
				PeerURLs: []string{
					"http://localhost:2380",
				},
				ClientURLs: []string{
					"http://localhost:2379",
				},
			},
			{
				ID:   2,
				Name: "etcd-2",
				PeerURLs: []string{
					"http://localhost:2381",
				},
				ClientURLs: []string{
					"http://localhost:2382",
				},
			},
			{
				ID:   3,
				Name: "etcd-3",
				PeerURLs: []string{
					"http://localhost:2383",
				},
				ClientURLs: []string{
					"http://localhost:2384",
				},
			},
		},
		leaderID:        1,
		alarms:          make([]*etcdserverpb.AlarmMember, 0),
		simulateLatency: false,
		readLatency:     1 * time.Millisecond,
		writeLatency:    2 * time.Millisecond,
	}
}

// Get retrieves a value from the mock KV store
func (m *MockEtcdClient) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	if m.simulateLatency {
		time.Sleep(m.readLatency)
	}

	m.opMutex.Lock()
	m.getCount++
	m.opMutex.Unlock()

	if m.getError != nil {
		return nil, m.getError
	}

	m.kvMutex.RLock()
	defer m.kvMutex.RUnlock()

	value, exists := m.kvStore[key]
	if !exists {
		return &clientv3.GetResponse{
			Header: &etcdserverpb.ResponseHeader{},
			Kvs:    []*mvccpb.KeyValue{},
		}, nil
	}

	return &clientv3.GetResponse{
		Header: &etcdserverpb.ResponseHeader{},
		Kvs: []*mvccpb.KeyValue{
			{
				Key:   []byte(key),
				Value: []byte(value),
			},
		},
	}, nil
}

// Put stores a key-value pair in the mock KV store
func (m *MockEtcdClient) Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	if m.simulateLatency {
		time.Sleep(m.writeLatency)
	}

	m.opMutex.Lock()
	m.putCount++
	m.opMutex.Unlock()

	if m.putError != nil {
		return nil, m.putError
	}

	m.kvMutex.Lock()
	defer m.kvMutex.Unlock()

	m.kvStore[key] = val

	return &clientv3.PutResponse{
		Header: &etcdserverpb.ResponseHeader{},
	}, nil
}

// Delete removes a key from the mock KV store
func (m *MockEtcdClient) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	if m.simulateLatency {
		time.Sleep(m.writeLatency)
	}

	m.opMutex.Lock()
	m.deleteCount++
	m.opMutex.Unlock()

	if m.deleteError != nil {
		return nil, m.deleteError
	}

	m.kvMutex.Lock()
	defer m.kvMutex.Unlock()

	delete(m.kvStore, key)

	return &clientv3.DeleteResponse{
		Header: &etcdserverpb.ResponseHeader{},
	}, nil
}

// MemberList returns the list of cluster members
func (m *MockEtcdClient) MemberList(ctx context.Context) (*clientv3.MemberListResponse, error) {
	if m.memberListError != nil {
		return nil, m.memberListError
	}

	m.memberMutex.RLock()
	defer m.memberMutex.RUnlock()

	return &clientv3.MemberListResponse{
		Header:  &etcdserverpb.ResponseHeader{},
		Members: m.members,
	}, nil
}

// Status returns the status of a specific endpoint
func (m *MockEtcdClient) Status(ctx context.Context, endpoint string) (*clientv3.StatusResponse, error) {
	if m.statusError != nil {
		return nil, m.statusError
	}

	m.memberMutex.RLock()
	defer m.memberMutex.RUnlock()

	// Find the member for this endpoint
	var memberID uint64
	for _, member := range m.members {
		for _, url := range member.ClientURLs {
			if url == endpoint {
				memberID = member.ID
				break
			}
		}
		if memberID != 0 {
			break
		}
	}

	if memberID == 0 {
		return nil, fmt.Errorf("endpoint not found: %s", endpoint)
	}

	return &clientv3.StatusResponse{
		Header:           &etcdserverpb.ResponseHeader{},
		Version:          "3.5.0",
		DbSize:           1024 * 1024,     // 1 MB
		DbSizeInUse:      512 * 1024,      // 512 KB
		Leader:           m.leaderID,
		RaftIndex:        1000,
		RaftTerm:         5,
		RaftAppliedIndex: 1000,
	}, nil
}

// AlarmList returns the list of active alarms
func (m *MockEtcdClient) AlarmList(ctx context.Context) (*clientv3.AlarmResponse, error) {
	if m.alarmListError != nil {
		return nil, m.alarmListError
	}

	m.memberMutex.RLock()
	defer m.memberMutex.RUnlock()

	return &clientv3.AlarmResponse{
		Header: &etcdserverpb.ResponseHeader{},
		Alarms: m.alarms,
	}, nil
}

// SetMemberListError sets an error to be returned by MemberList
func (m *MockEtcdClient) SetMemberListError(err error) {
	m.memberListError = err
}

// SetStatusError sets an error to be returned by Status
func (m *MockEtcdClient) SetStatusError(err error) {
	m.statusError = err
}

// SetGetError sets an error to be returned by Get
func (m *MockEtcdClient) SetGetError(err error) {
	m.getError = err
}

// SetPutError sets an error to be returned by Put
func (m *MockEtcdClient) SetPutError(err error) {
	m.putError = err
}

// SetDeleteError sets an error to be returned by Delete
func (m *MockEtcdClient) SetDeleteError(err error) {
	m.deleteError = err
}

// SetAlarmListError sets an error to be returned by AlarmList
func (m *MockEtcdClient) SetAlarmListError(err error) {
	m.alarmListError = err
}

// SetLeader changes the current leader
func (m *MockEtcdClient) SetLeader(leaderID uint64) {
	m.memberMutex.Lock()
	defer m.memberMutex.Unlock()
	m.leaderID = leaderID
}

// AddAlarm adds an alarm to the cluster
func (m *MockEtcdClient) AddAlarm(memberID uint64, alarmType etcdserverpb.AlarmType) {
	m.memberMutex.Lock()
	defer m.memberMutex.Unlock()

	m.alarms = append(m.alarms, &etcdserverpb.AlarmMember{
		MemberID: memberID,
		Alarm:    alarmType,
	})
}

// ClearAlarms removes all alarms
func (m *MockEtcdClient) ClearAlarms() {
	m.memberMutex.Lock()
	defer m.memberMutex.Unlock()
	m.alarms = make([]*etcdserverpb.AlarmMember, 0)
}

// SetLatencySimulation enables/disables latency simulation
func (m *MockEtcdClient) SetLatencySimulation(enabled bool, readLatency, writeLatency time.Duration) {
	m.simulateLatency = enabled
	m.readLatency = readLatency
	m.writeLatency = writeLatency
}

// GetOperationCounts returns the number of operations performed
func (m *MockEtcdClient) GetOperationCounts() (get, put, delete int64) {
	m.opMutex.Lock()
	defer m.opMutex.Unlock()
	return m.getCount, m.putCount, m.deleteCount
}

// ResetOperationCounts resets operation counters to zero
func (m *MockEtcdClient) ResetOperationCounts() {
	m.opMutex.Lock()
	defer m.opMutex.Unlock()
	m.getCount = 0
	m.putCount = 0
	m.deleteCount = 0
}

// AddMember adds a new member to the cluster
func (m *MockEtcdClient) AddMember(id uint64, name string, peerURLs, clientURLs []string) {
	m.memberMutex.Lock()
	defer m.memberMutex.Unlock()

	m.members = append(m.members, &etcdserverpb.Member{
		ID:         id,
		Name:       name,
		PeerURLs:   peerURLs,
		ClientURLs: clientURLs,
	})
}

// RemoveMember removes a member from the cluster
func (m *MockEtcdClient) RemoveMember(id uint64) {
	m.memberMutex.Lock()
	defer m.memberMutex.Unlock()

	for i, member := range m.members {
		if member.ID == id {
			m.members = append(m.members[:i], m.members[i+1:]...)
			break
		}
	}
}

// SetMembers replaces the entire member list
func (m *MockEtcdClient) SetMembers(members []*etcdserverpb.Member) {
	m.memberMutex.Lock()
	defer m.memberMutex.Unlock()
	m.members = members
}

// Close implements the client interface
func (m *MockEtcdClient) Close() error {
	// Mock implementation - nothing to close
	return nil
}

// Compact implements the client interface (stub)
func (m *MockEtcdClient) Compact(ctx context.Context, rev int64, opts ...clientv3.CompactOption) (*clientv3.CompactResponse, error) {
	return &clientv3.CompactResponse{}, nil
}

// Do implements the client interface (stub)
func (m *MockEtcdClient) Do(ctx context.Context, op clientv3.Op) (clientv3.OpResponse, error) {
	return clientv3.OpResponse{}, nil
}

// Txn implements the client interface (stub)
func (m *MockEtcdClient) Txn(ctx context.Context) clientv3.Txn {
	return nil
}
