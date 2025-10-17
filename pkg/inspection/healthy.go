package inspection

import (
	"k8s.io/klog/v2"

	etcdv1alpha1 "github.com/etcd-monitor/taskmaster/api/etcd/v1alpha1"
	"github.com/etcd-monitor/taskmaster/pkg/etcd"
	featureutil "github.com/etcd-monitor/taskmaster/pkg/featureprovider/util"
	"github.com/etcd-monitor/taskmaster/pkg/inspection/metrics"
)

func (c *Server) CollectMemberHealthy(inspection *etcdv1alpha1.EtcdInspection) error {
	namespace, name := inspection.Namespace, inspection.Spec.ClusterName
	cluster, clientConfig, err := c.GetEtcdClusterInfo(namespace, name)
	defer func() {
		if err != nil {
			featureutil.IncrFailedInspectionCounter(name, etcdv1alpha1.KStoneFeatureHealthy)
		}
	}()
	if err != nil {
		klog.Errorf("load tlsConfig failed, namespace is %s, name is %s, err is %v", namespace, name, err)
		return err
	}

	for _, m := range cluster.Status.Members {
		healthy, hErr := etcd.MemberHealthy(m.ExtensionClientUrl, clientConfig)
		labels := map[string]string{
			"clusterName": cluster.Name,
			"endpoint":    m.Endpoint,
		}
		if hErr != nil || !healthy {
			metrics.EtcdEndpointHealthy.With(labels).Set(0)
		} else {
			metrics.EtcdEndpointHealthy.With(labels).Set(1)
		}
	}
	return nil
}
