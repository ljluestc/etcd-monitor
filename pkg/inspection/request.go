package inspection

import (
	"context"
	"strings"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"k8s.io/klog/v2"

	etcdv1alpha1 "github.com/etcd-monitor/taskmaster/api/etcd/v1alpha1"
	"github.com/etcd-monitor/taskmaster/pkg/etcd"
	featureutil "github.com/etcd-monitor/taskmaster/pkg/featureprovider/util"
	"github.com/etcd-monitor/taskmaster/pkg/inspection/metrics"
)

// watcherInfo holds information about an active watcher
type watcherInfo struct {
	prefix       string
	resourceName string
	cancel       context.CancelFunc
}

var (
	// Global map to track active watchers per cluster
	activeWatchers = make(map[string][]watcherInfo)
	watcherMutex   sync.RWMutex
)

// CollectEtcdClusterRequest monitors etcd requests by watching key prefixes
func (c *Server) CollectEtcdClusterRequest(inspection *etcdv1alpha1.EtcdInspection) error {
	namespace, name := inspection.Namespace, inspection.Spec.ClusterName
	cluster, clientConfig, err := c.GetEtcdClusterInfo(namespace, name)
	defer func() {
		if err != nil {
			featureutil.IncrFailedInspectionCounter(name, etcdv1alpha1.KStoneFeatureRequest)
		}
	}()
	if err != nil {
		klog.Errorf("load tlsConfig failed, namespace is %s, name is %s, err is %v", namespace, name, err)
		return err
	}

	// Create etcd client
	cli, err := etcd.NewClientv3(clientConfig)
	if err != nil {
		klog.Errorf("failed to create etcd client for cluster %s, err is %v", name, err)
		return err
	}

	// Define key prefixes to watch (common Kubernetes resource prefixes)
	watchPrefixes := []struct {
		prefix       string
		resourceName string
	}{
		{"/registry/pods", "pods"},
		{"/registry/services", "services"},
		{"/registry/deployments", "deployments"},
		{"/registry/configmaps", "configmaps"},
		{"/registry/secrets", "secrets"},
		{"/registry/", "all"},
	}

	// Clean up existing watchers for this cluster
	c.stopWatchersForCluster(cluster.Name)

	// Start watchers for each prefix
	for _, wp := range watchPrefixes {
		ctx, cancel := context.WithCancel(context.Background())

		// Store watcher info
		watcherMutex.Lock()
		activeWatchers[cluster.Name] = append(activeWatchers[cluster.Name], watcherInfo{
			prefix:       wp.prefix,
			resourceName: wp.resourceName,
			cancel:       cancel,
		})
		watcherMutex.Unlock()

		// Start watching in a goroutine
		go c.watchKeyPrefix(ctx, cli, cluster.Name, wp.prefix, wp.resourceName)
	}

	// Also collect current key count for each prefix
	for _, wp := range watchPrefixes {
		count, err := c.countKeysWithPrefix(cli, wp.prefix)
		if err != nil {
			klog.Warningf("failed to count keys for prefix %s: %v", wp.prefix, err)
			continue
		}

		labels := map[string]string{
			"clusterName":  cluster.Name,
			"etcdPrefix":   wp.prefix,
			"resourceName": wp.resourceName,
		}
		metrics.EtcdKeyTotal.With(labels).Set(float64(count))
	}

	return nil
}

// watchKeyPrefix watches a specific key prefix and tracks operations
func (c *Server) watchKeyPrefix(ctx context.Context, cli *clientv3.Client, clusterName, prefix, resourceName string) {
	klog.V(2).Infof("starting watcher for cluster %s, prefix %s", clusterName, prefix)

	watcher := clientv3.NewWatcher(cli)
	defer watcher.Close()

	watchChan := watcher.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())

	for {
		select {
		case <-ctx.Done():
			klog.V(2).Infof("stopping watcher for cluster %s, prefix %s", clusterName, prefix)
			return
		case watchResp := <-watchChan:
			if watchResp.Err() != nil {
				klog.Errorf("watch error for cluster %s, prefix %s: %v", clusterName, prefix, watchResp.Err())
				// Attempt to restart watch after a delay
				time.Sleep(5 * time.Second)
				watchChan = watcher.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())
				continue
			}

			for _, event := range watchResp.Events {
				c.trackRequestEvent(clusterName, prefix, resourceName, event)
			}
		}
	}
}

// trackRequestEvent tracks a single watch event and updates metrics
func (c *Server) trackRequestEvent(clusterName, prefix, resourceName string, event *clientv3.Event) {
	var grpcMethod string

	switch event.Type {
	case clientv3.EventTypePut:
		if event.IsCreate() {
			grpcMethod = "CREATE"
		} else {
			grpcMethod = "PUT"
		}
	case clientv3.EventTypeDelete:
		grpcMethod = "DELETE"
	default:
		grpcMethod = "UNKNOWN"
	}

	labels := map[string]string{
		"clusterName":  clusterName,
		"grpcMethod":   grpcMethod,
		"etcdPrefix":   prefix,
		"resourceName": resourceName,
	}

	metrics.EtcdRequestTotal.With(labels).Inc()

	klog.V(4).Infof("tracked request: cluster=%s, method=%s, prefix=%s, key=%s",
		clusterName, grpcMethod, prefix, string(event.Kv.Key))
}

// countKeysWithPrefix counts the number of keys with a given prefix
func (c *Server) countKeysWithPrefix(cli *clientv3.Client, prefix string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCountOnly())
	if err != nil {
		return 0, err
	}

	return resp.Count, nil
}

// stopWatchersForCluster stops all watchers for a given cluster
func (c *Server) stopWatchersForCluster(clusterName string) {
	watcherMutex.Lock()
	defer watcherMutex.Unlock()

	if watchers, exists := activeWatchers[clusterName]; exists {
		for _, watcher := range watchers {
			watcher.cancel()
		}
		delete(activeWatchers, clusterName)
		klog.V(2).Infof("stopped all watchers for cluster %s", clusterName)
	}
}

// StopAllWatchers stops all active watchers (useful for cleanup)
func (c *Server) StopAllWatchers() {
	watcherMutex.Lock()
	defer watcherMutex.Unlock()

	for clusterName, watchers := range activeWatchers {
		for _, watcher := range watchers {
			watcher.cancel()
		}
		klog.V(2).Infof("stopped all watchers for cluster %s", clusterName)
	}

	activeWatchers = make(map[string][]watcherInfo)
}

// parseResourceNameFromKey extracts resource name from etcd key
func parseResourceNameFromKey(key string) string {
	parts := strings.Split(key, "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	return "unknown"
}
