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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:defaulter-gen=TypeMeta

// EtcdInspection is a specification for a EtcdInspection resource
type EtcdInspection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   EtcdInspectionSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status EtcdInspectionStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// EtcdInspectionSpec is the spec for a EtcdInspectionSpec resource
type EtcdInspectionSpec struct {
	ClusterName        string `json:"clusterName" protobuf:"bytes,1,opt,name=clusterName"`
	InspectionType     string `json:"inspectionType" protobuf:"bytes,2,opt,name=inspectionType"`
	InspectionProvider string `json:"inspectionProvider,omitempty" protobuf:"bytes,3,opt,name=inspectionProvider"`
}

// EtcdInspectionStatus is the status for a EtcdInspection resource
type EtcdInspectionStatus struct {
	Phase string `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EtcdInspectionList is a list of EtcdInspection resources
type EtcdInspectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []EtcdInspection `json:"items" protobuf:"bytes,2,rep,name=items"`
}
