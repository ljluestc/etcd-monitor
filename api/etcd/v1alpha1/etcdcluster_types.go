/*
Copyright 2023 Clay.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EtcdCluster is the Schema for the etcdclusters API
type EtcdCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   EtcdClusterSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status EtcdClusterStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// Defines cluster types, storage backends, and feature gates
type EtcdClusterType string
type EtcdStorageBackend string
type KStoneFeature string

const (
	// Cluster types
	EtcdClusterImported EtcdClusterType = "imported"

	// Storage backends
	EtcdStorageV2 EtcdStorageBackend = "v2"
	EtcdStorageV3 EtcdStorageBackend = "v3"

	// Feature gates
	KStoneFeatureHealthy     KStoneFeature = "healthy"
	KStoneFeatureAlarm       KStoneFeature = "alarm"
	KStoneFeatureRequest     KStoneFeature = "request"
	KStoneFeatureConsistency KStoneFeature = "consistency"
)

const (
	KStoneFeatureAnno = "kstone.tkestack.io/feature"
)

// EtcdClusterSpec defines the desired state of EtcdCluster
type EtcdClusterSpec struct {
	ClusterType    EtcdClusterType    `json:"clusterType" protobuf:"bytes,1,opt,name=clusterType"`
	StorageBackend EtcdStorageBackend `json:"storageBackend,omitempty" protobuf:"bytes,2,opt,name=storageBackend"`
	Version        string             `json:"version,omitempty" protobuf:"bytes,3,opt,name=version"`
	Size           int32              `json:"size,omitempty" protobuf:"varint,4,opt,name=size"`
	SecureConfig   *SecureConfig      `json:"secureConfig,omitempty" protobuf:"bytes,5,opt,name=secureConfig"`
}

// SecureConfig contains TLS configuration
type SecureConfig struct {
	TLSSecret corev1.LocalObjectReference `json:"tlsSecret,omitempty" protobuf:"bytes,1,opt,name=tlsSecret"`
}

// EtcdClusterStatus defines the observed state of EtcdCluster
type EtcdClusterStatus struct {
	Phase       string       `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase"`
	ServiceName string       `json:"serviceName,omitempty" protobuf:"bytes,2,opt,name=serviceName"`
	Members     []EtcdMember `json:"members,omitempty" protobuf:"bytes,3,rep,name=members"`
}

// EtcdMember represents a member of the etcd cluster
type EtcdMember struct {
	MemberId           string `json:"memberId,omitempty" protobuf:"bytes,1,opt,name=memberId"`
	Name               string `json:"name,omitempty" protobuf:"bytes,2,opt,name=name"`
	Endpoint           string `json:"endpoint,omitempty" protobuf:"bytes,3,opt,name=endpoint"`
	ExtensionClientUrl string `json:"extensionClientUrl,omitempty" protobuf:"bytes,4,opt,name=extensionClientUrl"`
	Version            string `json:"version,omitempty" protobuf:"bytes,5,opt,name=version"`
	IsLearner          bool   `json:"isLearner,omitempty" protobuf:"varint,6,opt,name=isLearner"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EtcdClusterList contains a list of EtcdCluster
type EtcdClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []EtcdCluster `json:"items" protobuf:"bytes,2,rep,name=items"`
}
