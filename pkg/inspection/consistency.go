package inspection

import (
	"k8s.io/klog/v2"

	etcdv1alpha1 "github.com/etcd-monitor/taskmaster/api/etcd/v1alpha1"
	"github.com/etcd-monitor/taskmaster/pkg/etcd"
	featureutil "github.com/etcd-monitor/taskmaster/pkg/featureprovider/util"
	"github.com/etcd-monitor/taskmaster/pkg/inspection/metrics"
)

// memberData represents metadata for a cluster member
type memberData struct {
	endpoint          string
	keyTotal          int64
	revision          int64
	index             int64
	raftAppliedIndex  int64
	raftIndex         int64
}

// CollectClusterConsistentData collects consistency data from all etcd members
// and calculates differences between members
func (c *Server) CollectClusterConsistentData(inspection *etcdv1alpha1.EtcdInspection) error {
	namespace, name := inspection.Namespace, inspection.Spec.ClusterName
	cluster, clientConfig, err := c.GetEtcdClusterInfo(namespace, name)
	defer func() {
		if err != nil {
			featureutil.IncrFailedInspectionCounter(name, etcdv1alpha1.KStoneFeatureConsistency)
		}
	}()
	if err != nil {
		klog.Errorf("load tlsConfig failed, namespace is %s, name is %s, err is %v", namespace, name, err)
		return err
	}

	// Collect statistics from all members
	endpoints := make([]string, 0, len(cluster.Status.Members))
	for _, m := range cluster.Status.Members {
		endpoints = append(endpoints, m.ExtensionClientUrl)
	}

	// Create etcd client
	cli, err := etcd.NewClientv3(clientConfig)
	if err != nil {
		klog.Errorf("failed to create etcd client for cluster %s, err is %v", name, err)
		return err
	}
	defer cli.Close()

	// Create stat collector
	config := &etcd.ClientConfig{
		Endpoints: endpoints,
	}
	stat := etcd.NewStat(config, string(cluster.Spec.StorageBackend))

	// Collect metadata from all members

	memberStats := make([]memberData, 0, len(endpoints))

	for _, endpoint := range endpoints {
		// Get key total
		keyTotal, err := stat.GetTotalKeyNum(endpoint)
		if err != nil {
			klog.Warningf("failed to get key total for %s: %v", endpoint, err)
			continue
		}

		// Get index information
		indexInfo, err := stat.GetIndex(endpoint)
		if err != nil {
			klog.Warningf("failed to get index for %s: %v", endpoint, err)
			continue
		}

		memberStats = append(memberStats, memberData{
			endpoint:          endpoint,
			keyTotal:          keyTotal,
			revision:          indexInfo.Revision,
			index:             indexInfo.Index,
			raftAppliedIndex:  int64(indexInfo.RaftAppliedIndex),
			raftIndex:         int64(indexInfo.RaftIndex),
		})
	}

	if len(memberStats) == 0 {
		klog.Warningf("no member stats collected for cluster %s", name)
		return nil
	}

	// Calculate differences
	keyTotalDiff := calculateDifference(memberStats, func(m memberData) int64 { return m.keyTotal })
	revisionDiff := calculateDifference(memberStats, func(m memberData) int64 { return m.revision })
	indexDiff := calculateDifference(memberStats, func(m memberData) int64 { return m.index })
	raftAppliedIndexDiff := calculateDifference(memberStats, func(m memberData) int64 { return m.raftAppliedIndex })
	raftIndexDiff := calculateDifference(memberStats, func(m memberData) int64 { return m.raftIndex })

	// Update metrics
	clusterLabels := map[string]string{"clusterName": cluster.Name}

	metrics.EtcdNodeDiffTotal.With(clusterLabels).Set(float64(keyTotalDiff))
	metrics.EtcdNodeRevisionDiff.With(clusterLabels).Set(float64(revisionDiff))
	metrics.EtcdNodeIndexDiff.With(clusterLabels).Set(float64(indexDiff))
	metrics.EtcdNodeRaftAppliedIndexDiff.With(clusterLabels).Set(float64(raftAppliedIndexDiff))
	metrics.EtcdNodeRaftIndexDiff.With(clusterLabels).Set(float64(raftIndexDiff))

	klog.V(2).Infof("cluster %s consistency check: keyTotal diff=%d, revision diff=%d, index diff=%d, raftAppliedIndex diff=%d, raftIndex diff=%d",
		name, keyTotalDiff, revisionDiff, indexDiff, raftAppliedIndexDiff, raftIndexDiff)

	return nil
}

// calculateDifference calculates the difference between max and min values
func calculateDifference(stats []memberData, valueFunc func(memberData) int64) int64 {
	if len(stats) == 0 {
		return 0
	}

	min := valueFunc(stats[0])
	max := valueFunc(stats[0])

	for _, stat := range stats[1:] {
		value := valueFunc(stat)
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}

	return max - min
}
