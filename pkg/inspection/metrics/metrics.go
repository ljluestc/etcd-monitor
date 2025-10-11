package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	EtcdNodeDiffTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "inspection",
		Name:      "etcd_node_diff_total",
		Help:      "total etcd node diff key",
	}, []string{"clusterName"})

	EtcdEndpointHealthy = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "inspection",
		Name:      "etcd_endpoint_healthy",
		Help:      "The healthy of etcd member",
	}, []string{"clusterName", "endpoint"})

	EtcdRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem: "inspection",
		Name:      "etcd_request_total",
		Help:      "The total number of etcd requests",
	}, []string{"clusterName", "grpcMethod", "etcdPrefix", "resourceName"})

	EtcdKeyTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "inspection",
		Name:      "etcd_key_total",
		Help:      "The total number of etcd key",
	}, []string{"clusterName", "etcdPrefix", "resourceName"})

	EtcdEndpointAlarm = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "inspection",
		Name:      "etcd_endpoint_alarm",
		Help:      "The alarm of etcd member",
	}, []string{"clusterName", "endpoint", "alarmType"})

	EtcdNodeRevisionDiff = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "inspection",
		Name:      "etcd_node_revision_diff_total",
		Help:      "The revision difference between all member",
	}, []string{"clusterName"})

	EtcdNodeIndexDiff = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "inspection",
		Name:      "etcd_node_index_diff_total",
		Help:      "The index difference between all member",
	}, []string{"clusterName"})

	EtcdNodeRaftAppliedIndexDiff = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "inspection",
		Name:      "etcd_node_raft_applied_index_diff_total",
		Help:      "The raftAppliedIndex difference between all member",
	}, []string{"clusterName"})

	EtcdNodeRaftIndexDiff = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "inspection",
		Name:      "etcd_node_raft_index_diff_total",
		Help:      "The raftIndex difference between all member",
	}, []string{"clusterName"})

	EtcdInspectionFailedNum = prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem: "inspection",
		Name:      "etcd_inspection_failed_total",
		Help:      "The total number of failed inspections",
	}, []string{"clusterName", "inspectionType"})
)

func init() {
	prometheus.MustRegister(EtcdNodeDiffTotal)
	prometheus.MustRegister(EtcdEndpointHealthy)
	prometheus.MustRegister(EtcdRequestTotal)
	prometheus.MustRegister(EtcdKeyTotal)
	prometheus.MustRegister(EtcdEndpointAlarm)
	prometheus.MustRegister(EtcdNodeRevisionDiff)
	prometheus.MustRegister(EtcdNodeIndexDiff)
	prometheus.MustRegister(EtcdNodeRaftAppliedIndexDiff)
	prometheus.MustRegister(EtcdNodeRaftIndexDiff)
	prometheus.MustRegister(EtcdInspectionFailedNum)
}
