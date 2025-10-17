package featureprovider

import (
	etcdv1alpha1 "github.com/etcd-monitor/taskmaster/api/etcd/v1alpha1"
	"github.com/etcd-monitor/taskmaster/pkg/controllers/util"
	"github.com/etcd-monitor/taskmaster/pkg/etcd"
)

// Feature is an abstract, pluggable interface for cluster features.
type Feature interface {
	// Equal checks whether the feature needs to be updated
	Equal(cluster *etcdv1alpha1.EtcdCluster) bool

	// Sync synchronizes the latest feature configuration
	Sync(cluster *etcdv1alpha1.EtcdCluster) error

	// Do executes inspection tasks.
	Do(task *etcdv1alpha1.EtcdInspection) error
}

type FeatureContext struct {
	ClientBuilder      util.ClientBuilder
	ClientConfigGetter etcd.ClientConfigGetter
}
