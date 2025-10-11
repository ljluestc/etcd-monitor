package etcd

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
	"k8s.io/klog/v2"
)

// MemberList returns the member list from etcd cluster
func MemberList(client *clientv3.Client) (*clientv3.MemberListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), client.Endpoints()[0])
	defer cancel()

	resp, err := client.MemberList(ctx)
	if err != nil {
		klog.Errorf("failed to get member list: %v", err)
		return nil, err
	}

	return resp, nil
}

// Status returns the status of an etcd member
func Status(endpoint string, client *clientv3.Client) (*clientv3.StatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), client.Endpoints()[0])
	defer cancel()

	resp, err := client.Status(ctx, endpoint)
	if err != nil {
		klog.Errorf("failed to get status for %s: %v", endpoint, err)
		return nil, err
	}

	return resp, nil
}

// AlarmList returns the list of alarms from etcd cluster
func AlarmList(client *clientv3.Client) (*clientv3.AlarmResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), client.Endpoints()[0])
	defer cancel()

	resp, err := client.AlarmList(ctx)
	if err != nil {
		klog.Errorf("failed to get alarm list: %v", err)
		return nil, err
	}

	return resp, nil
}

// Alarm represents an etcd alarm
type Alarm struct {
	MemberID  uint64
	AlarmType string
}

// GetEtcdAlarms gets all alarms from the etcd cluster
func GetEtcdAlarms(endpoints []string, config *ClientConfig) ([]*Alarm, error) {
	client, err := NewClientv3(config)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	alarmResp, err := AlarmList(client)
	if err != nil {
		return nil, err
	}

	alarms := make([]*Alarm, 0)
	for _, alarm := range alarmResp.Alarms {
		alarms = append(alarms, &Alarm{
			MemberID:  alarm.MemberID,
			AlarmType: alarm.Alarm.String(),
		})
	}

	return alarms, nil
}

// GetMemberEndpoint returns the client endpoint for a member ID
func GetMemberEndpoint(client *clientv3.Client, memberID uint64) (string, error) {
	memberList, err := MemberList(client)
	if err != nil {
		return "", err
	}

	for _, member := range memberList.Members {
		if member.ID == memberID {
			if len(member.ClientURLs) > 0 {
				return member.ClientURLs[0], nil
			}
		}
	}

	return "", fmt.Errorf("member %d not found", memberID)
}
