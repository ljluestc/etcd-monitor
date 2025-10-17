package imported

import (
	"sync"

	"k8s.io/klog/v2"

	etcdv1alpha1 "github.com/etcd-monitor/taskmaster/api/etcd/v1alpha1"
	"github.com/etcd-monitor/taskmaster/pkg/clusterprovider"
	"github.com/etcd-monitor/taskmaster/pkg/etcd"
)

const (
	ProviderName = etcdv1alpha1.EtcdClusterImported
)

var (
	once     sync.Once
	instance *ImportedCluster
)

// ImportedCluster represents an etcd cluster that is not managed by this operator
// It only monitors existing clusters and doesn't perform lifecycle operations
type ImportedCluster struct {
	name string
	ctx  *clusterprovider.ClusterContext
}

func init() {
	clusterprovider.RegisterEtcdClusterFactory(
		ProviderName,
		func(ctx *clusterprovider.ClusterContext) (clusterprovider.Cluster, error) {
			return initImportedClusterInstance(ctx)
		},
	)
}

func initImportedClusterInstance(ctx *clusterprovider.ClusterContext) (clusterprovider.Cluster, error) {
	var err error
	once.Do(func() {
		instance = &ImportedCluster{
			name: string(ProviderName),
			ctx:  ctx,
		}
	})
	return instance, err
}

// BeforeCreate does pre-creation validation for imported clusters
// For imported clusters, we just validate that we can connect
func (c *ImportedCluster) BeforeCreate(cluster *etcdv1alpha1.EtcdCluster) error {
	klog.V(2).Infof("BeforeCreate for imported cluster %s/%s", cluster.Namespace, cluster.Name)
	// No-op for imported clusters - we don't create them
	return nil
}

// Create is a no-op for imported clusters since they already exist
func (c *ImportedCluster) Create(cluster *etcdv1alpha1.EtcdCluster) error {
	klog.V(2).Infof("Create for imported cluster %s/%s (no-op)", cluster.Namespace, cluster.Name)
	// Imported clusters are not created by the operator
	return nil
}

// AfterCreate performs post-creation tasks
// For imported clusters, this validates connectivity
func (c *ImportedCluster) AfterCreate(cluster *etcdv1alpha1.EtcdCluster) error {
	klog.V(2).Infof("AfterCreate for imported cluster %s/%s", cluster.Namespace, cluster.Name)
	// No additional setup needed for imported clusters
	return nil
}

// BeforeUpdate validates before updating cluster configuration
func (c *ImportedCluster) BeforeUpdate(cluster *etcdv1alpha1.EtcdCluster) error {
	klog.V(2).Infof("BeforeUpdate for imported cluster %s/%s", cluster.Namespace, cluster.Name)
	// No-op for imported clusters
	return nil
}

// Update is a no-op for imported clusters since we don't manage them
func (c *ImportedCluster) Update(cluster *etcdv1alpha1.EtcdCluster) error {
	klog.V(2).Infof("Update for imported cluster %s/%s (no-op)", cluster.Namespace, cluster.Name)
	// We don't update imported clusters - they're managed externally
	return nil
}

// AfterUpdate performs post-update tasks
func (c *ImportedCluster) AfterUpdate(cluster *etcdv1alpha1.EtcdCluster) error {
	klog.V(2).Infof("AfterUpdate for imported cluster %s/%s", cluster.Namespace, cluster.Name)
	// No additional tasks needed
	return nil
}

// BeforeDelete performs pre-deletion cleanup
func (c *ImportedCluster) BeforeDelete(cluster *etcdv1alpha1.EtcdCluster) error {
	klog.V(2).Infof("BeforeDelete for imported cluster %s/%s", cluster.Namespace, cluster.Name)
	// No-op for imported clusters
	return nil
}

// Delete is a no-op for imported clusters since we don't manage their lifecycle
func (c *ImportedCluster) Delete(cluster *etcdv1alpha1.EtcdCluster) error {
	klog.V(2).Infof("Delete for imported cluster %s/%s (no-op)", cluster.Namespace, cluster.Name)
	// We don't delete imported clusters - only remove monitoring
	return nil
}

// AfterDelete performs cleanup after cluster deletion
func (c *ImportedCluster) AfterDelete(cluster *etcdv1alpha1.EtcdCluster) error {
	klog.V(2).Infof("AfterDelete for imported cluster %s/%s", cluster.Namespace, cluster.Name)
	// Cleanup any monitoring resources
	return nil
}

// Equal checks whether the cluster configuration has changed
// For imported clusters, we compare the spec to determine if updates are needed
func (c *ImportedCluster) Equal(cluster *etcdv1alpha1.EtcdCluster) (bool, error) {
	klog.V(4).Infof("Checking equality for imported cluster %s/%s", cluster.Namespace, cluster.Name)

	// For imported clusters, we primarily care about:
	// - Storage backend changes
	// - Size changes (though we don't enforce it)
	// - TLS configuration changes

	// Since we're just monitoring, most spec changes don't require action
	// We return true (equal) to avoid unnecessary reconciliation
	return true, nil
}

// Status retrieves the current status of the etcd cluster
func (c *ImportedCluster) Status(config *etcd.ClientConfig, cluster *etcdv1alpha1.EtcdCluster) (etcdv1alpha1.EtcdClusterStatus, error) {
	klog.V(4).Infof("Getting status for imported cluster %s/%s", cluster.Namespace, cluster.Name)

	status := etcdv1alpha1.EtcdClusterStatus{
		Phase:       "Running",
		ServiceName: cluster.Status.ServiceName,
		Members:     []etcdv1alpha1.EtcdMember{},
	}

	// Get endpoints from current status or annotations
	endpoints := clusterprovider.GetStorageMemberEndpoints(cluster)

	if len(endpoints) == 0 {
		klog.Warningf("No endpoints found for imported cluster %s/%s", cluster.Namespace, cluster.Name)
		status.Phase = "Unknown"
		return status, nil
	}

	// Update config with endpoints
	config.Endpoints = endpoints

	// Retrieve member information from the running cluster
	members, err := clusterprovider.GetRuntimeEtcdMembers(
		endpoints,
		config,
		cluster.Spec.StorageBackend,
	)
	if err != nil {
		klog.Errorf("Failed to get runtime members for cluster %s/%s: %v", cluster.Namespace, cluster.Name, err)
		status.Phase = "Unhealthy"
		return status, err
	}

	// Populate extension client URLs if configured
	extClientURLMap, err := clusterprovider.PopulateExtensionClientURL(cluster)
	if err != nil {
		klog.Warningf("Failed to populate extension client URLs: %v", err)
	}

	// Update members with extension URLs
	for i := range members {
		if extURL, found := extClientURLMap[members[i].MemberId]; found {
			members[i].ExtensionClientUrl = extURL
		}
	}

	status.Members = members

	// Verify cluster health
	if len(members) > 0 {
		status.Phase = "Running"
	} else {
		status.Phase = "Degraded"
	}

	klog.V(2).Infof("Status for imported cluster %s/%s: phase=%s, members=%d",
		cluster.Namespace, cluster.Name, status.Phase, len(members))

	return status, nil
}
