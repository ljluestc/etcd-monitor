package etcd

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
	"k8s.io/klog/v2"
)

// Stat is an interface for collecting etcd statistics
type Stat interface {
	// GetTotalKeyNum returns the total number of keys in etcd
	GetTotalKeyNum(endpoint string) (int64, error)
	// GetIndex returns metadata indices for consistency checking
	GetIndex(endpoint string) (*IndexInfo, error)
}

// IndexInfo contains etcd metadata indices
type IndexInfo struct {
	Revision         int64
	Index            int64
	RaftAppliedIndex uint64
	RaftIndex        uint64
}

// V3Stat implements Stat interface for etcd v3
type V3Stat struct {
	config *ClientConfig
}

// NewV3Stat creates a new V3Stat
func NewV3Stat(config *ClientConfig) Stat {
	return &V3Stat{
		config: config,
	}
}

// GetTotalKeyNum returns the total number of keys in etcd v3
func (s *V3Stat) GetTotalKeyNum(endpoint string) (int64, error) {
	client, err := NewClientv3(s.config)
	if err != nil {
		return 0, err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), s.config.DialTimeout)
	defer cancel()

	// Get all keys with count only
	resp, err := client.Get(ctx, "\x00", clientv3.WithFromKey(), clientv3.WithCountOnly())
	if err != nil {
		klog.Errorf("failed to get key count: %v", err)
		return 0, err
	}

	return resp.Count, nil
}

// GetIndex returns metadata indices for consistency checking
func (s *V3Stat) GetIndex(endpoint string) (*IndexInfo, error) {
	client, err := NewClientv3(s.config)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	statusResp, err := Status(endpoint, client)
	if err != nil {
		return nil, err
	}

	info := &IndexInfo{
		Revision:         statusResp.Header.Revision,
		RaftAppliedIndex: statusResp.RaftAppliedIndex,
		RaftIndex:        statusResp.RaftIndex,
	}

	return info, nil
}

// V2Stat implements Stat interface for etcd v2
type V2Stat struct {
	config *ClientConfig
}

// NewV2Stat creates a new V2Stat
func NewV2Stat(config *ClientConfig) Stat {
	return &V2Stat{
		config: config,
	}
}

// GetTotalKeyNum returns the total number of keys in etcd v2
func (s *V2Stat) GetTotalKeyNum(endpoint string) (int64, error) {
	// For v2, we need to use the HTTP API to count keys
	// This is a simplified implementation
	klog.Warning("v2 key counting not fully implemented")
	return 0, fmt.Errorf("v2 key counting not supported")
}

// GetIndex returns metadata indices for etcd v2
func (s *V2Stat) GetIndex(endpoint string) (*IndexInfo, error) {
	// For v2, we can get some stats from the stats endpoint
	stats, err := GetStats(endpoint, s.config)
	if err != nil {
		return nil, err
	}

	info := &IndexInfo{
		// V2 doesn't have all the same fields as v3
		Index: 0, // Would need to parse from stats
	}

	klog.V(4).Infof("Got v2 stats for %s: %+v", endpoint, stats)

	return info, nil
}

// NewStat creates a new Stat based on the storage backend version
func NewStat(config *ClientConfig, version string) Stat {
	if version == "v2" {
		return NewV2Stat(config)
	}
	return NewV3Stat(config)
}
