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
	"k8s.io/apimachinery/pkg/runtime"
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

// DeepCopyObject returns a generically typed copy of an object
func (in *EtcdInspection) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyObject returns a generically typed copy of an object
func (in *EtcdInspectionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopy returns a deep copy of the EtcdInspection
func (in *EtcdInspection) DeepCopy() *EtcdInspection {
	if in == nil {
		return nil
	}
	out := new(EtcdInspection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies all properties of this object into another object of the same type
func (in *EtcdInspection) DeepCopyInto(out *EtcdInspection) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy returns a deep copy of the EtcdInspectionList
func (in *EtcdInspectionList) DeepCopy() *EtcdInspectionList {
	if in == nil {
		return nil
	}
	out := new(EtcdInspectionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies all properties of this object into another object of the same type
func (in *EtcdInspectionList) DeepCopyInto(out *EtcdInspectionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EtcdInspection, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy returns a deep copy of the EtcdInspectionSpec
func (in *EtcdInspectionSpec) DeepCopy() *EtcdInspectionSpec {
	if in == nil {
		return nil
	}
	out := new(EtcdInspectionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies all properties of this object into another object of the same type
func (in *EtcdInspectionSpec) DeepCopyInto(out *EtcdInspectionSpec) {
	*out = *in
}

// DeepCopy returns a deep copy of the EtcdInspectionStatus
func (in *EtcdInspectionStatus) DeepCopy() *EtcdInspectionStatus {
	if in == nil {
		return nil
	}
	out := new(EtcdInspectionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies all properties of this object into another object of the same type
func (in *EtcdInspectionStatus) DeepCopyInto(out *EtcdInspectionStatus) {
	*out = *in
}
