package inspection

import (
	"strconv"

	"k8s.io/klog/v2"

	etcdv1alpha1 "github.com/etcd-monitor/taskmaster/api/etcd/v1alpha1"
	"github.com/etcd-monitor/taskmaster/pkg/clusterprovider"
	featureutil "github.com/etcd-monitor/taskmaster/pkg/featureprovider/util"
	"github.com/etcd-monitor/taskmaster/pkg/inspection/metrics"
)

var alarmTypeList = []string{"NOSPACE", "CORRUPT"}

// CollectAlarmList collects the alarms of etcd, and
// transfer them to prometheus metrics
func (c *Server) CollectAlarmList(inspection *etcdv1alpha1.EtcdInspection) error {
	namespace, name := inspection.Namespace, inspection.Spec.ClusterName
	cluster, clientConfig, err := c.GetEtcdClusterInfo(namespace, name)
	defer func() {
		if err != nil {
			featureutil.IncrFailedInspectionCounter(name, etcdv1alpha1.KStoneFeatureAlarm)
		}
	}()
	if err != nil {
		klog.Errorf("load tlsConfig failed, namespace is %s, name is %s, err is %v", namespace, name, err)
		return err
	}

	alarms, err := clusterprovider.GetEtcdAlarms([]string{cluster.Status.ServiceName}, clientConfig)
	if err != nil {
		return err
	}

	for _, m := range cluster.Status.Members {
		if len(alarms) == 0 {
			cleanAllAlarmMetrics(cluster.Name, m.Endpoint)
		}
		for _, a := range alarms {
			if m.MemberId == strconv.FormatUint(a.MemberID, 10) {
				labels := map[string]string{
					"clusterName": cluster.Name,
					"endpoint":    m.Endpoint,
					"alarmType":   a.AlarmType,
				}
				metrics.EtcdEndpointAlarm.With(labels).Set(1)
			}
		}
	}
	return nil
}

func cleanAllAlarmMetrics(clusterName, endpoint string) {
	for _, alarmType := range alarmTypeList {
		labels := map[string]string{
			"clusterName": clusterName,
			"endpoint":    endpoint,
			"alarmType":   alarmType,
		}
		metrics.EtcdEndpointAlarm.With(labels).Set(0)
	}
}
