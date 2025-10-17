package clusterprovider

import (
	"context"
	"sort"

	"k8s.io/klog/v2"

	etcdv1alpha1 "github.com/etcd-monitor/taskmaster/api/etcd/v1alpha1"
	"github.com/etcd-monitor/taskmaster/pkg/etcd"
	"github.com/etcd-monitor/taskmaster/pkg/etcd/client"
)

// GetStorageMemberEndpoints returns the member endpoints from etcd cluster
func GetStorageMemberEndpoints(cluster *etcdv1alpha1.EtcdCluster) []string {
	endpoints := make([]string, 0)
	for _, member := range cluster.Status.Members {
		if member.ExtensionClientUrl != "" {
			endpoints = append(endpoints, member.ExtensionClientUrl)
		} else {
			endpoints = append(endpoints, member.Endpoint)
		}
	}
	return endpoints
}

// PopulateExtensionClientURL populates the extension client URL map
func PopulateExtensionClientURL(cluster *etcdv1alpha1.EtcdCluster) (map[string]string, error) {
	extClientURLMap := make(map[string]string)
	if cluster.Annotations == nil {
		return extClientURLMap, nil
	}

	// Parse annotation for extension client URLs
	// Expected format: "extClientURL=memberID:url,memberID:url"
	// This is a simplified implementation
	return extClientURLMap, nil
}

// GetRuntimeEtcdMembers retrieves the current etcd members from a running cluster
func GetRuntimeEtcdMembers(endpoints []string, config *etcd.ClientConfig, storageBackend etcdv1alpha1.EtcdStorageBackend) ([]etcdv1alpha1.EtcdMember, error) {
	members := make([]etcdv1alpha1.EtcdMember, 0)

	versionClient, err := client.GetEtcdClientProvider(storageBackend, &client.VersionContext{
		Config: config,
	})
	if err != nil {
		klog.Errorf("failed to get etcd client provider, err: %v", err)
		return members, err
	}
	defer versionClient.Close()

	memberList, err := versionClient.MemberList()
	if err != nil {
		klog.Errorf("failed to get member list, err: %v", err)
		return members, err
	}

	// Get status from one endpoint to get version and leader info
	var status *client.Member
	if len(endpoints) > 0 {
		status, err = versionClient.Status(endpoints[0])
		if err != nil {
			klog.Warningf("failed to get status from %s, err: %v", endpoints[0], err)
		}
	}

	for _, m := range memberList {
		endpoint := ""
		if len(m.ClientURLs) > 0 {
			endpoint = m.ClientURLs[0]
		}

		member := etcdv1alpha1.EtcdMember{
			MemberId:  m.ID,
			Name:      m.Name,
			Endpoint:  endpoint,
			IsLearner: m.IsLearner,
		}

		if status != nil {
			member.Version = status.Version
		}

		members = append(members, member)
	}

	return members, nil
}

// GetEtcdAlarms retrieves alarms from etcd cluster
func GetEtcdAlarms(endpoints []string, config *etcd.ClientConfig) ([]*etcd.Alarm, error) {
	alarms := make([]*etcd.Alarm, 0)

	cli, err := etcd.NewClientv3(config)
	if err != nil {
		klog.Errorf("failed to create etcd client, err: %v", err)
		return alarms, err
	}
	defer cli.Close()

	ctx := context.Background()
	alarmList, err := cli.AlarmList(ctx)
	if err != nil {
		klog.Errorf("failed to get alarm list, err: %v", err)
		return alarms, err
	}

	for _, alarm := range alarmList.Alarms {
		alarmType := "UNKNOWN"
		switch alarm.Alarm {
		case 1:
			alarmType = "NOSPACE"
		case 2:
			alarmType = "CORRUPT"
		}

		alarms = append(alarms, &etcd.Alarm{
			MemberID:  alarm.MemberID,
			AlarmType: alarmType,
		})
	}

	return alarms, nil
}

// IsMemberListEqual checks if two member lists are equal
func IsMemberListEqual(list1, list2 []etcdv1alpha1.EtcdMember) bool {
	if len(list1) != len(list2) {
		return false
	}

	// Sort by member ID for comparison
	sort.Slice(list1, func(i, j int) bool {
		return list1[i].MemberId < list1[j].MemberId
	})
	sort.Slice(list2, func(i, j int) bool {
		return list2[i].MemberId < list2[j].MemberId
	})

	for i := range list1 {
		if list1[i].MemberId != list2[i].MemberId ||
			list1[i].Name != list2[i].Name ||
			list1[i].Endpoint != list2[i].Endpoint ||
			list1[i].IsLearner != list2[i].IsLearner {
			return false
		}
	}

	return true
}
