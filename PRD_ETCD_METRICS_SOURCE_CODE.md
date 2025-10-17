# Product Requirements Document: ETCD-MONITOR: Etcd Metrics Source Code

---

## Document Information
**Project:** etcd-monitor
**Document:** ETCD_METRICS_SOURCE_CODE
**Version:** 1.0.0
**Date:** 2025-10-13
**Status:** READY FOR TASK-MASTER PARSING

---

## 1. EXECUTIVE SUMMARY

### 1.1 Overview
This PRD captures the requirements and implementation details for ETCD-MONITOR: Etcd Metrics Source Code.

### 1.2 Purpose
This document provides a structured specification that can be parsed by task-master to generate actionable tasks.

### 1.3 Scope
The scope includes all requirements, features, and implementation details from the original documentation.

---

## 2. REQUIREMENTS

### 2.1 Functional Requirements
**Priority:** HIGH

**REQ-001:** Providers**: Pluggable inspection features (alarm, consistency, healthy, request)

**REQ-002:** providers using `sync.Once`

**REQ-003:** watches the changes of etcdcluster resources

**REQ-004:** watches the changes of etcdinspection resources

**REQ-005:** is an abstract, pluggable interface for cluster features.

**REQ-006:** needs to be updated

**REQ-007:** provider:%s was registered twice", name)

**REQ-008:** provider:%s", name)

**REQ-009:** provider not found")

**REQ-010:** = range EtcdFeatureProviders {


## 3. TASKS

The following tasks have been identified for implementation:

**TASK_001** [MEDIUM]: **Repository**: https://github.com/clay-wangzhi/etcd-metrics

**TASK_002** [MEDIUM]: **Purpose**: Kubernetes Operator for monitoring etcd clusters

**TASK_003** [MEDIUM]: **License**: Apache License 2.0

**TASK_004** [MEDIUM]: **Author**: Clay Wang

**TASK_005** [HIGH]: **EtcdCluster Controller**: Manages the lifecycle of etcd clusters

**TASK_006** [HIGH]: **EtcdInspection Controller**: Performs health checks and monitoring tasks

**TASK_007** [MEDIUM]: **Cluster Providers**: Different etcd cluster implementations (e.g., imported clusters)

**TASK_008** [MEDIUM]: **Feature Providers**: Pluggable inspection features (alarm, consistency, healthy, request)

**TASK_009** [MEDIUM]: **Client Versions**: Support for different etcd storage backends (v2, v3)

**TASK_010** [MEDIUM]: **Factory Pattern**: For creating providers and clients

**TASK_011** [MEDIUM]: **Singleton Pattern**: For feature providers using `sync.Once`

**TASK_012** [MEDIUM]: **Observer Pattern**: Using Kubernetes informers and watchers

**TASK_013** [MEDIUM]: **Strategy Pattern**: Different inspection strategies via feature providers

**TASK_014** [MEDIUM]: `NewController()`: Initializes the controller with informers and listers

**TASK_015** [MEDIUM]: `Run()`: Starts worker goroutines to process the work queue

**TASK_016** [MEDIUM]: Creates new clusters

**TASK_017** [MEDIUM]: Updates existing clusters

**TASK_018** [MEDIUM]: Updates cluster status

**TASK_019** [MEDIUM]: Handles feature synchronization

**TASK_020** [MEDIUM]: `handleObject()`: Enqueues work items when resources change

**TASK_021** [MEDIUM]: Work queues for asynchronous processing

**TASK_022** [MEDIUM]: Rate limiting for retry logic

**TASK_023** [MEDIUM]: Event recording for audit trails

**TASK_024** [MEDIUM]: Informers for caching and watching resources

**TASK_025** [MEDIUM]: Watches EtcdInspection custom resources

**TASK_026** [MEDIUM]: Triggers inspection tasks based on feature providers

**TASK_027** [MEDIUM]: Exposes Prometheus metrics via HTTP endpoint

**TASK_028** [MEDIUM]: Uses Gin web framework for the metrics server

**TASK_029** [MEDIUM]: Getting storage member endpoints

**TASK_030** [MEDIUM]: Populating extension client URL maps

**TASK_031** [MEDIUM]: Getting runtime etcd members

**TASK_032** [MEDIUM]: Getting etcd alarms

**TASK_033** [MEDIUM]: Checking member list equality

**TASK_034** [MEDIUM]: Uses singleton pattern with `sync.Once`

**TASK_035** [MEDIUM]: Registers itself via `init()` function

**TASK_036** [MEDIUM]: Implements lifecycle hooks (BeforeCreate, Create, AfterCreate, etc.)

**TASK_037** [MEDIUM]: Retrieves cluster status from existing etcd endpoints

**TASK_038** [MEDIUM]: `ClientConfig` struct with endpoints, timeouts, TLS settings

**TASK_039** [MEDIUM]: `SecureConfig` for authentication credentials

**TASK_040** [MEDIUM]: `ClientConfigSecret` for retrieving credentials from Kubernetes secrets

**TASK_041** [MEDIUM]: Functions to create etcd v2 and v3 clients

**TASK_042** [MEDIUM]: HTTP-based health check client

**TASK_043** [MEDIUM]: TLS configuration support

**TASK_044** [MEDIUM]: `/health` endpoint monitoring

**TASK_045** [MEDIUM]: `/version` endpoint for version detection

**TASK_046** [MEDIUM]: `/v2/stats/self` for etcd statistics

**TASK_047** [MEDIUM]: Creating etcd v2 and v3 clients

**TASK_048** [MEDIUM]: Member list operations

**TASK_049** [MEDIUM]: Status checking

**TASK_050** [MEDIUM]: Certificate and authentication management

**TASK_051** [MEDIUM]: `Stat` interface for v2 and v3 backends

**TASK_052** [MEDIUM]: `GetTotalKeyNum()`: Counts keys in etcd

**TASK_053** [MEDIUM]: `GetIndex()`: Retrieves metadata indices for consistency checking

**TASK_054** [MEDIUM]: Supports both v2 and v3 storage backends

**TASK_055** [MEDIUM]: Manages inspection tasks for all features

**TASK_056** [MEDIUM]: Creates and maintains etcd clients and watchers

**TASK_057** [MEDIUM]: Synchronizes inspection resources with clusters

**TASK_058** [MEDIUM]: `NewInspectionServer()`: Initializes the server

**TASK_059** [MEDIUM]: `Equal()`: Checks if inspection needs updating

**TASK_060** [MEDIUM]: `Sync()`: Creates/updates EtcdInspection resources

**TASK_061** [MEDIUM]: `GetEtcdClusterInfo()`: Retrieves cluster and config info

**TASK_062** [MEDIUM]: Collecting metadata from all etcd members

**TASK_063** [MEDIUM]: Comparing key totals, revisions, and raft indices

**TASK_064** [MEDIUM]: Calculating differences between members

**TASK_065** [MEDIUM]: Exposing metrics for inconsistencies

**TASK_066** [MEDIUM]: Creating watchers on specified key prefixes

**TASK_067** [MEDIUM]: Tracking PUT, DELETE, and other operations

**TASK_068** [MEDIUM]: Counting requests per resource type

**TASK_069** [MEDIUM]: Using etcd watch API with event channels

**TASK_070** [HIGH]: Creating a new provider package

**TASK_071** [HIGH]: Implementing the `Feature` interface

**TASK_072** [HIGH]: Registering via `init()`

**TASK_073** [HIGH]: Importing with `_ "package/path"` in the providers.go file

**TASK_074** [HIGH]: **Informers**: Watch and cache Kubernetes resources

**TASK_075** [HIGH]: **Work Queue**: Buffer work items with rate limiting

**TASK_076** [HIGH]: **Sync Handler**: Reconciliation logic

**TASK_077** [HIGH]: **Event Handler**: Enqueues items when resources change

**TASK_078** [MEDIUM]: **k8s.io/client-go**: Kubernetes client library

**TASK_079** [MEDIUM]: **k8s.io/apimachinery**: Kubernetes API machinery

**TASK_080** [MEDIUM]: **k8s.io/klog/v2**: Structured logging

**TASK_081** [MEDIUM]: **go.etcd.io/etcd/client/v3**: etcd v3 client

**TASK_082** [MEDIUM]: **go.etcd.io/etcd/client/v2**: etcd v2 client

**TASK_083** [MEDIUM]: **github.com/prometheus/client_golang**: Prometheus metrics

**TASK_084** [MEDIUM]: **github.com/gin-gonic/gin**: HTTP web framework

**TASK_085** [MEDIUM]: **github.com/spf13/cobra**: CLI framework

**TASK_086** [MEDIUM]: **sigs.k8s.io/controller-runtime**: Controller runtime utilities

**TASK_087** [MEDIUM]: **deepcopy-gen**: Generates deep copy methods

**TASK_088** [MEDIUM]: **client-gen**: Generates typed clients

**TASK_089** [MEDIUM]: **informer-gen**: Generates informers

**TASK_090** [MEDIUM]: **lister-gen**: Generates listers

**TASK_091** [HIGH]: **EtcdCluster**: Represents an etcd cluster to be monitored

**TASK_092** [HIGH]: **EtcdInspection**: Represents an inspection task

**TASK_093** [HIGH]: **EtcdCluster Controller**: Manages cluster lifecycle

**TASK_094** [HIGH]: **EtcdInspection Controller**: Performs monitoring tasks

**TASK_095** [MEDIUM]: `/metrics`: Prometheus metrics endpoint

**TASK_096** [MEDIUM]: Metrics include health, alarms, consistency, and request counts

**TASK_097** [HIGH]: **Monitors etcd clusters** through pluggable inspection features

**TASK_098** [HIGH]: **Exposes Prometheus metrics** for observability

**TASK_099** [HIGH]: **Uses clean abstractions** with interfaces and factory patterns

**TASK_100** [HIGH]: **Supports multiple etcd versions** (v2 and v3)

**TASK_101** [HIGH]: **Follows Kubernetes best practices** with controllers, informers, and CRDs

**TASK_102** [HIGH]: **Enables extensibility** through plugin architecture

**TASK_103** [HIGH]: **Provides high availability** with leader election support

**TASK_104** [MEDIUM]: Clear separation of concerns

**TASK_105** [MEDIUM]: Dependency injection

**TASK_106** [MEDIUM]: Interface-based design

**TASK_107** [MEDIUM]: Proper error handling

**TASK_108** [MEDIUM]: Structured logging

**TASK_109** [MEDIUM]: Comprehensive metrics


## 4. DETAILED SPECIFICATIONS

### 4.1 Original Content

The following sections contain the original documentation:


#### Complete Source Code Analysis Etcd Metrics Repository

# Complete Source Code Analysis: etcd-metrics Repository


#### Repository Information

## Repository Information
- **Repository**: https://github.com/clay-wangzhi/etcd-metrics
- **Purpose**: Kubernetes Operator for monitoring etcd clusters
- **License**: Apache License 2.0
- **Author**: Clay Wang

---


#### 1 Complete Directory Structure

## 1. Complete Directory Structure

```
etcd-metrics/
├── api/
│   └── etcd/
│       └── v1alpha1/
│           ├── doc.go
│           ├── etcdcluster_types.go
│           ├── etcdinspection_types.go
│           ├── groupversion_info.go
│           ├── register.go
│           └── zz_generated.deepcopy.go
├── cmd/
│   ├── etcdcluster-controller/
│   │   └── controller.go
│   ├── etcdinspection-controller/
│   │   └── controller.go
│   └── main.go
├── pkg/
│   ├── clusterprovider/
│   │   ├── etcdcluster.go
│   │   ├── helper.go
│   │   ├── plugins.go
│   │   └── providers/
│   │       ├── providers.go
│   │       └── imported/
│   │           └── cluster.go
│   ├── controllers/
│   │   ├── etcdcluster/
│   │   │   └── etcdclusters_controller.go
│   │   ├── etcdinspection/
│   │   │   └── etcdinspection-controller.go
│   │   └── util/
│   │       └── util.go
│   ├── etcd/
│   │   ├── client/
│   │   │   ├── versions/
│   │   │   │   ├── v3/
│   │   │   │   │   └── client.go
│   │   │   │   └── providers.go
│   │   │   ├── client.go
│   │   │   └── version.go
│   │   ├── client.go
│   │   ├── health.go
│   │   ├── helper.go
│   │   └── stats.go
│   ├── featureprovider/
│   │   ├── providers/
│   │   │   ├── alarm/

... (content truncated for PRD) ...


#### 2 Core Architecture Patterns

## 2. Core Architecture Patterns


#### 2 1 Controller Pattern

### 2.1 Controller Pattern
The project implements the Kubernetes Operator pattern with two main controllers:
1. **EtcdCluster Controller**: Manages the lifecycle of etcd clusters
2. **EtcdInspection Controller**: Performs health checks and monitoring tasks


#### 2 2 Plugin Architecture

### 2.2 Plugin Architecture
The codebase uses a plugin-based factory pattern for extensibility:
- **Cluster Providers**: Different etcd cluster implementations (e.g., imported clusters)
- **Feature Providers**: Pluggable inspection features (alarm, consistency, healthy, request)
- **Client Versions**: Support for different etcd storage backends (v2, v3)


#### 2 3 Key Design Patterns

### 2.3 Key Design Patterns
- **Factory Pattern**: For creating providers and clients
- **Singleton Pattern**: For feature providers using `sync.Once`
- **Observer Pattern**: Using Kubernetes informers and watchers
- **Strategy Pattern**: Different inspection strategies via feature providers

---


#### 3 Complete Source Code Files

## 3. Complete Source Code Files


#### 3 1 Api Definitions

### 3.1 API Definitions


#### Api Etcd V1Alpha1 Etcdcluster Types Go

#### api/etcd/v1alpha1/etcdcluster_types.go
```go
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


... (content truncated for PRD) ...


#### Api Etcd V1Alpha1 Etcdinspection Types Go

#### api/etcd/v1alpha1/etcdinspection_types.go
```go
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

... (content truncated for PRD) ...


#### Api Etcd V1Alpha1 Groupversion Info Go

#### api/etcd/v1alpha1/groupversion_info.go
```go
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

// Package v1alpha1 contains API Schema definitions for the etcd v1alpha1 API group
// +kubebuilder:object:generate=true
// +groupName=etcd.clay.io
package v1alpha1
```

---


#### 3 2 Command Entry Points

### 3.2 Command Entry Points


#### Cmd Main Go

#### cmd/main.go
```go
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

package main

import (
	goflag "flag"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	klog "k8s.io/klog/v2"

	etcdclustercontroller "etcd-operator/cmd/etcdcluster-controller"
	etcdinspectioncontroller "etcd-operator/cmd/etcdinspection-controller"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	var cmd = &cobra.Command{
		Use:              "etcd-controller",
		Short:            "run etcdcluster/etcdinspection controller",
		PersistentPreRun: func(c *cobra.Command, args []string) {},
	}

	flags := cmd.PersistentFlags()
	out := cmd.OutOrStdout()
	cmd.AddCommand(
		etcdclustercontroller.NewEtcdClusterControllerCommand(out),
		etcdinspectioncontroller.NewEtcdInspectionControllerCommand(out),
	)

	klog.InitFlags(nil)

... (content truncated for PRD) ...


#### Cmd Etcdcluster Controller Controller Go

#### cmd/etcdcluster-controller/controller.go
```go
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

package etcdclustercontroller

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	klog "k8s.io/klog/v2"

	"etcd-operator/pkg/controllers/etcdcluster"
	"etcd-operator/pkg/controllers/util"
	"etcd-operator/pkg/k8s"
	"etcd-operator/pkg/signals"
)

type EtcdClusterCommand struct {
	out                io.Writer
	kubeconfig         string
	masterURL          string
	labelSelector      string
	leaseLockName      string
	leaseLockNamespace string
}

... (content truncated for PRD) ...


#### Cmd Etcdinspection Controller Controller Go

#### cmd/etcdinspection-controller/controller.go
```go
package etcdinspectioncontroller

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	klog "k8s.io/klog/v2"

	"etcd-operator/pkg/controllers/etcdinspection"
	"etcd-operator/pkg/controllers/util"
	"etcd-operator/pkg/k8s"
	"etcd-operator/pkg/signals"
)

type EtcdInspectionCommand struct {
	out                io.Writer
	kubeconfig         string
	masterURL          string
	labelSelector      string
	leaseLockName      string
	leaseLockNamespace string
}

// NewEtcdInspectionControllerCommand creates a *cobra.Command object with default parameters
func NewEtcdInspectionControllerCommand(out io.Writer) *cobra.Command {
	cc := &EtcdInspectionCommand{out: out}
	cmd := &cobra.Command{
		Use:   "inspection",
		Short: "run inspection controller",
		Long: `The inspection controller is a daemon, it will watches the changes of etcdinspection resources
 through the apiserver and makes changes attempting to move the current state towards the desired state.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				klog.V(1).Infof("FLAG: --%s=%q", flag.Name, flag.Value)
			})
			if err := cc.run(); err != nil {
				return err
			}

... (content truncated for PRD) ...


#### 3 3 Controllers

### 3.3 Controllers


#### Pkg Controllers Etcdcluster Etcdclusters Controller Go

#### pkg/controllers/etcdcluster/etcdclusters_controller.go
**Note**: This file is extensive. The controller implements the full reconciliation loop for managing EtcdCluster resources. Key functions include:

- `NewController()`: Initializes the controller with informers and listers
- `Run()`: Starts worker goroutines to process the work queue
- `syncHandler()`: Main reconciliation logic that:
  - Creates new clusters
  - Updates existing clusters
  - Updates cluster status
  - Handles feature synchronization
- `handleObject()`: Enqueues work items when resources change

The controller follows the standard Kubernetes controller pattern with:
- Work queues for asynchronous processing
- Rate limiting for retry logic
- Event recording for audit trails
- Informers for caching and watching resources


#### Pkg Controllers Etcdinspection Etcdinspection Controller Go

#### pkg/controllers/etcdinspection/etcdinspection-controller.go
**Note**: Similar structure to the cluster controller but focused on inspection tasks. It:

- Watches EtcdInspection custom resources
- Triggers inspection tasks based on feature providers
- Exposes Prometheus metrics via HTTP endpoint
- Uses Gin web framework for the metrics server


#### Pkg Controllers Util Util Go

#### pkg/controllers/util/util.go
```go
package util

import (
	"fmt"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/diff"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/dynamic"
	clientset "k8s.io/client-go/kubernetes"
	fake "k8s.io/client-go/kubernetes/fake"
	restclient "k8s.io/client-go/rest"
	core "k8s.io/client-go/testing"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	"etcd-operator/pkg/k8s"
)

const (
	ComponentEtcdClusterController    = "etcdcluster-controller"
	ComponentEtcdInspectionController = "etcdinspection-controller"
)

type EtcdClusterPhase string

const (
	EtcdClusterCreating     EtcdClusterPhase = "EtcdClusterCreating"
	EtcdClusterUpdating     EtcdClusterPhase = "EtcdClusterUpdating"
	EtcdClusterUpdateStatus EtcdClusterPhase = "EtcdClusterUpdateStatus"
)

const (
	ClusterTLSSecretName      = "certName"
	ClusterExtensionClientURL = "extClientURL"
)

type ClientBuilder interface {
	ConfigOrDie() *restclient.Config
	ClientOrDie() clientset.Interface
	DynamicClientOrDie() dynamic.Interface
}

func NewSimpleClientBuilder(kubeconfig string) ClientBuilder {
	builder := &simpleClientBuilder{
		kubeconfig: kubeconfig,

... (content truncated for PRD) ...


#### 3 4 Cluster Provider

### 3.4 Cluster Provider


#### Pkg Clusterprovider Etcdcluster Go

#### pkg/clusterprovider/etcdcluster.go
```go
package clusterprovider

import (
	"k8s.io/client-go/dynamic"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
	"etcd-operator/pkg/controllers/util"
	"etcd-operator/pkg/etcd"
)

// Cluster is an abstract, pluggable interface for etcd clusters.
type Cluster interface {

	// BeforeCreate does some things before creating the cluster
	BeforeCreate(cluster *etcdv1alpha1.EtcdCluster) error
	// Create creates the cluster
	Create(cluster *etcdv1alpha1.EtcdCluster) error
	// AfterCreate does some things after creating the cluster
	AfterCreate(cluster *etcdv1alpha1.EtcdCluster) error

	// BeforeUpdate does some things before updating the cluster
	BeforeUpdate(cluster *etcdv1alpha1.EtcdCluster) error
	// Update updates the cluster
	Update(cluster *etcdv1alpha1.EtcdCluster) error
	// AfterUpdate does some things after updating the cluster
	AfterUpdate(cluster *etcdv1alpha1.EtcdCluster) error

	// BeforeDelete does some things before deleting the cluster
	BeforeDelete(cluster *etcdv1alpha1.EtcdCluster) error
	// Delete deletes the cluster
	Delete(cluster *etcdv1alpha1.EtcdCluster) error
	// AfterDelete does some things after deleting the cluster
	AfterDelete(cluster *etcdv1alpha1.EtcdCluster) error

	// Equal checks whether the cluster needs to be updated
	Equal(cluster *etcdv1alpha1.EtcdCluster) (bool, error)

	// Status gets the cluster status
	Status(config *etcd.ClientConfig, cluster *etcdv1alpha1.EtcdCluster) (etcdv1alpha1.EtcdClusterStatus, error)
}

type ClusterContext struct {
	Clientbuilder util.ClientBuilder
	Client        dynamic.Interface
}
```


#### Pkg Clusterprovider Plugins Go

#### pkg/clusterprovider/plugins.go
```go
package clusterprovider

import (
	"errors"
	"sync"

	"k8s.io/klog/v2"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
)

type EtcdFactory func(cluster *ClusterContext) (Cluster, error)

var (
	mutex     sync.Mutex
	providers = make(map[etcdv1alpha1.EtcdClusterType]EtcdFactory)
)

// RegisterEtcdClusterFactory registers the specified cluster provider
func RegisterEtcdClusterFactory(name etcdv1alpha1.EtcdClusterType, factory EtcdFactory) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, found := providers[name]; found {
		klog.V(2).Infof("etcdcluster provider %s was registered twice", name)
	}

	klog.V(2).Infof("register etcdCluster provider %s", name)
	providers[name] = factory
}

// GetEtcdClusterProvider gets the specified cluster provider
func GetEtcdClusterProvider(
	name etcdv1alpha1.EtcdClusterType,
	ctx *ClusterContext,
) (Cluster, error) {
	mutex.Lock()
	defer mutex.Unlock()
	f, found := providers[name]

	klog.V(1).Infof("get provider name %s,status:%t", name, found)
	if !found {
		return nil, errors.New("fatal error,etcd cluster provider not found")
	}
	return f(ctx)
}
```


#### Pkg Clusterprovider Helper Go

#### pkg/clusterprovider/helper.go
**Note**: Contains helper functions for:
- Getting storage member endpoints
- Populating extension client URL maps
- Getting runtime etcd members
- Getting etcd alarms
- Checking member list equality


#### Pkg Clusterprovider Providers Providers Go

#### pkg/clusterprovider/providers/providers.go
```go
package providers

import (
	_ "etcd-operator/pkg/clusterprovider/providers/imported" // import imported provider
)
```


#### Pkg Clusterprovider Providers Imported Cluster Go

#### pkg/clusterprovider/providers/imported/cluster.go
**Note**: Implements the Cluster interface for imported etcd clusters (clusters not created by the operator). Key methods:
- Uses singleton pattern with `sync.Once`
- Registers itself via `init()` function
- Implements lifecycle hooks (BeforeCreate, Create, AfterCreate, etc.)
- Retrieves cluster status from existing etcd endpoints

---


#### 3 5 Etcd Client And Utilities

### 3.5 Etcd Client and Utilities


#### Pkg Etcd Client Go

#### pkg/etcd/client.go
**Note**: Provides client configuration management:
- `ClientConfig` struct with endpoints, timeouts, TLS settings
- `SecureConfig` for authentication credentials
- `ClientConfigSecret` for retrieving credentials from Kubernetes secrets
- Functions to create etcd v2 and v3 clients


#### Pkg Etcd Health Go

#### pkg/etcd/health.go
**Note**: Implements health checking:
- HTTP-based health check client
- TLS configuration support
- `/health` endpoint monitoring
- `/version` endpoint for version detection
- `/v2/stats/self` for etcd statistics


#### Pkg Etcd Helper Go

#### pkg/etcd/helper.go
**Note**: Utility functions for:
- Creating etcd v2 and v3 clients
- Member list operations
- Status checking
- Certificate and authentication management


#### Pkg Etcd Stats Go

#### pkg/etcd/stats.go
**Note**: Statistics collection interface:
- `Stat` interface for v2 and v3 backends
- `GetTotalKeyNum()`: Counts keys in etcd
- `GetIndex()`: Retrieves metadata indices for consistency checking
- Supports both v2 and v3 storage backends


#### Pkg Etcd Client Client Go

#### pkg/etcd/client/client.go
```go
package client

import (
	"etcd-operator/pkg/etcd"
)

// Member contains member info including v2 and v3
type Member struct {
	ID         string
	Name       string
	PeerURLs   []string
	ClientURLs []string
	Version    string
	IsLearner  bool
	Leader     string
}

type VersionClient interface {
	MemberList() ([]Member, error)
	Status(endpoint string) (*Member, error)
	Close()
}

type VersionContext struct {
	Config *etcd.ClientConfig
}
```


#### Pkg Etcd Client Version Go

#### pkg/etcd/client/version.go
```go
package client

import (
	"errors"
	"sync"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"

	"k8s.io/klog/v2"
)

type Factory func(cluster *VersionContext) (VersionClient, error)

var (
	mutex     sync.Mutex
	providers = make(map[etcdv1alpha1.EtcdStorageBackend]Factory)
)

// RegisterEtcdClientFactory registers the specified etcd client
func RegisterEtcdClientFactory(name etcdv1alpha1.EtcdStorageBackend, factory Factory) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, found := providers[name]; found {
		klog.V(2).Infof("etcdcluster provider %s was registered twice", name)
	}

	klog.V(2).Infof("register etcdCluster provider %s", name)
	providers[name] = factory
}

// GetEtcdClientProvider gets the specified etcd client
func GetEtcdClientProvider(
	name etcdv1alpha1.EtcdStorageBackend,
	ctx *VersionContext,
) (VersionClient, error) {
	mutex.Lock()
	defer mutex.Unlock()

	// compatible with existing clusters
	if name == "" {
		name = etcdv1alpha1.EtcdStorageV3
	}
	f, found := providers[name]

	klog.V(1).Infof("get provider name %s,status:%t", name, found)
	if !found {
		return nil, errors.New("fatal error,etcd cluster provider not found")

... (content truncated for PRD) ...


#### Pkg Etcd Client Versions Providers Go

#### pkg/etcd/client/versions/providers.go
```go
package versions

import (
	_ "etcd-operator/pkg/etcd/client/versions/v3" // import etcd client of v3
)
```


#### Pkg Etcd Client Versions V3 Client Go

#### pkg/etcd/client/versions/v3/client.go
```go
package v3

import (
	"strconv"

	clientv3 "go.etcd.io/etcd/client/v3"
	"k8s.io/klog/v2"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
	"etcd-operator/pkg/etcd"
	"etcd-operator/pkg/etcd/client"
)

type V3 struct {
	ctx *client.VersionContext
	cli *clientv3.Client
}

func (c *V3) MemberList() ([]client.Member, error) {
	members := make([]client.Member, 0)
	memberRsp, err := etcd.MemberList(c.cli)
	if err != nil {
		klog.Errorf("failed to get member list, endpoints is %s,err is %v", c.ctx.Config.Endpoints, err)
		return members, err
	}
	for _, m := range memberRsp.Members {
		members = append(members, client.Member{
			ID:         strconv.FormatUint(m.ID, 10),
			Name:       m.Name,
			PeerURLs:   m.PeerURLs,
			ClientURLs: m.ClientURLs,
			IsLearner:  m.IsLearner,
		})
	}
	return members, nil
}

func (c *V3) Status(endpoint string) (*client.Member, error) {
	statusRsp, err := etcd.Status(c.ctx.Config.Endpoints[0], c.cli)
	if err != nil {
		return nil, err
	}
	return &client.Member{
		Version:   statusRsp.Version,
		IsLearner: statusRsp.IsLearner,
		Leader:    strconv.FormatUint(statusRsp.Leader, 10),
	}, nil
}

... (content truncated for PRD) ...


#### 3 6 Feature Provider

### 3.6 Feature Provider


#### Pkg Featureprovider Feature Go

#### pkg/featureprovider/feature.go
```go
package featureprovider

import (
	"etcd-operator/api/etcd/v1alpha1"
	"etcd-operator/pkg/controllers/util"
	"etcd-operator/pkg/etcd"
)

// Feature is an abstract, pluggable interface for cluster features.
type Feature interface {
	// Equal checks whether the feature needs to be updated
	Equal(cluster *v1alpha1.EtcdCluster) bool

	// Sync synchronizes the latest feature configuration
	Sync(cluster *v1alpha1.EtcdCluster) error

	// Do executes inspection tasks.
	Do(task *v1alpha1.EtcdInspection) error
}

type FeatureContext struct {
	ClientBuilder      util.ClientBuilder
	ClientConfigGetter etcd.ClientConfigGetter
}
```


#### Pkg Featureprovider Plugins Go

#### pkg/featureprovider/plugins.go
```go
package featureprovider

import (
	"errors"
	"sync"

	"k8s.io/klog/v2"
)

var (
	mutex                sync.Mutex
	EtcdFeatureProviders = make(map[string]FeatureFactory)
)

type FeatureFactory func(cfg *FeatureContext) (Feature, error)

// RegisterFeatureFactory registers the specified feature provider
func RegisterFeatureFactory(name string, factory FeatureFactory) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, found := EtcdFeatureProviders[name]; found {
		klog.V(2).Infof("feature provider:%s was registered twice", name)
	}

	klog.V(2).Infof("feature provider:%s", name)
	EtcdFeatureProviders[name] = factory
}

// GetFeatureProvider gets the specified feature provider
func GetFeatureProvider(name string, ctx *FeatureContext) (Feature, error) {
	mutex.Lock()
	defer mutex.Unlock()
	f, found := EtcdFeatureProviders[name]

	klog.V(1).Infof("get provider name %s,status:%t", name, found)
	if !found {
		return nil, errors.New("fatal error,feature provider not found")
	}
	return f(ctx)
}

// ListFeatureProvider lists all feature provider
func ListFeatureProvider() []string {
	var features []string
	mutex.Lock()
	defer mutex.Unlock()
	for feature := range EtcdFeatureProviders {

... (content truncated for PRD) ...


#### Pkg Featureprovider Util Util Go

#### pkg/featureprovider/util/util.go
```go
package util

import (
	"fmt"
	"strconv"
	"strings"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
)

const (
	FeatureStatusEnabled  = "enabled"
	FeatureStatusDisabled = "disabled"
)

type ConsistencyType string

const (
	ConsistencyKeyTotal             ConsistencyType = "keyTotal"
	ConsistencyRevision             ConsistencyType = "revision"
	ConsistencyIndex                ConsistencyType = "index"
	ConsistencyRaftRaftAppliedIndex ConsistencyType = "raftAppliedIndex"
	ConsistencyRaftIndex            ConsistencyType = "raftIndex"
)

const (
	OneDaySeconds = 24 * 60 * 60
)

func IsFeatureGateEnabled(annotations map[string]string, name etcdv1alpha1.KStoneFeature) bool {
	if gates, found := annotations[etcdv1alpha1.KStoneFeatureAnno]; found && gates != "" {
		featurelist := strings.Split(gates, ",")
		for _, item := range featurelist {
			features := strings.Split(item, "=")
			if len(features) != 2 {
				continue
			}

			enabled, _ := strconv.ParseBool(features[1])
			if etcdv1alpha1.KStoneFeature(features[0]) == name && enabled {
				return true
			}
		}
	}
	return false
}

func IncrFailedInspectionCounter(clusterName string, featureName etcdv1alpha1.KStoneFeature) {

... (content truncated for PRD) ...


#### Pkg Featureprovider Providers Providers Go

#### pkg/featureprovider/providers/providers.go
```go
package providers

import (
	// register consistency inspection feature/
	_ "etcd-operator/pkg/featureprovider/providers/consistency"
	// register healthy inspection feature
	_ "etcd-operator/pkg/featureprovider/providers/healthy"
	// register request inspection feature
	_ "etcd-operator/pkg/featureprovider/providers/request"
	// register alarm inspection feature
	_ "etcd-operator/pkg/featureprovider/providers/alarm"
)
```


#### Pkg Featureprovider Providers Alarm Alarm Go

#### pkg/featureprovider/providers/alarm/alarm.go
```go
package alarm

import (
	"sync"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
	"etcd-operator/pkg/featureprovider"
	"etcd-operator/pkg/inspection"
)

const (
	ProviderName = string(etcdv1alpha1.KStoneFeatureAlarm)
)

var (
	once     sync.Once
	instance *FeatureAlarm
)

type FeatureAlarm struct {
	name       string
	inspection *inspection.Server
	ctx        *featureprovider.FeatureContext
}

func init() {
	featureprovider.RegisterFeatureFactory(
		ProviderName,
		func(ctx *featureprovider.FeatureContext) (featureprovider.Feature, error) {
			return initFeatureAlarmInstance(ctx)
		},
	)
}

func initFeatureAlarmInstance(ctx *featureprovider.FeatureContext) (featureprovider.Feature, error) {
	var err error
	once.Do(func() {
		instance = &FeatureAlarm{
			name: ProviderName,
			ctx:  ctx,
		}
		instance.inspection, err = inspection.NewInspectionServer(ctx)
	})
	return instance, err
}

func (c *FeatureAlarm) Equal(cluster *etcdv1alpha1.EtcdCluster) bool {
	return c.inspection.Equal(cluster, etcdv1alpha1.KStoneFeatureAlarm)

... (content truncated for PRD) ...


#### Pkg Featureprovider Providers Consistency Consistency Go

#### pkg/featureprovider/providers/consistency/consistency.go
```go
package consistency

import (
	"sync"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
	"etcd-operator/pkg/featureprovider"
	"etcd-operator/pkg/inspection"
)

const (
	ProviderName = string(etcdv1alpha1.KStoneFeatureConsistency)
)

var (
	once     sync.Once
	instance *FeatureConsistency
)

type FeatureConsistency struct {
	name       string
	inspection *inspection.Server
	ctx        *featureprovider.FeatureContext
}

func init() {
	featureprovider.RegisterFeatureFactory(
		ProviderName,
		func(ctx *featureprovider.FeatureContext) (featureprovider.Feature, error) {
			return initFeatureConsistencyInstance(ctx)
		},
	)
}

func initFeatureConsistencyInstance(ctx *featureprovider.FeatureContext) (featureprovider.Feature, error) {
	var err error
	once.Do(func() {
		instance = &FeatureConsistency{
			name: ProviderName,
			ctx:  ctx,
		}
		instance.inspection, err = inspection.NewInspectionServer(ctx)
	})
	return instance, err
}

func (c *FeatureConsistency) Equal(cluster *etcdv1alpha1.EtcdCluster) bool {
	return c.inspection.Equal(cluster, etcdv1alpha1.KStoneFeatureConsistency)

... (content truncated for PRD) ...


#### Pkg Featureprovider Providers Healthy Healthy Go

#### pkg/featureprovider/providers/healthy/healthy.go
```go
package healthy

import (
	"sync"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
	"etcd-operator/pkg/featureprovider"
	"etcd-operator/pkg/inspection"
)

var (
	once     sync.Once
	instance *FeatureHealthy
)

type FeatureHealthy struct {
	name       string
	inspection *inspection.Server
	ctx        *featureprovider.FeatureContext
}

const (
	ProviderName = string(etcdv1alpha1.KStoneFeatureHealthy)
)

func init() {
	featureprovider.RegisterFeatureFactory(
		ProviderName,
		func(ctx *featureprovider.FeatureContext) (featureprovider.Feature, error) {
			return initFeatureHealthyInstance(ctx)
		},
	)
}

func initFeatureHealthyInstance(ctx *featureprovider.FeatureContext) (featureprovider.Feature, error) {
	var err error
	once.Do(func() {
		instance = &FeatureHealthy{
			name: ProviderName,
			ctx:  ctx,
		}
		instance.inspection, err = inspection.NewInspectionServer(ctx)
	})
	return instance, err
}

func (c *FeatureHealthy) Equal(cluster *etcdv1alpha1.EtcdCluster) bool {
	return c.inspection.Equal(cluster, etcdv1alpha1.KStoneFeatureHealthy)

... (content truncated for PRD) ...


#### Pkg Featureprovider Providers Request Request Go

#### pkg/featureprovider/providers/request/request.go
```go
package request

import (
	"sync"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
	"etcd-operator/pkg/featureprovider"
	"etcd-operator/pkg/inspection"
)

const (
	ProviderName = string(etcdv1alpha1.KStoneFeatureRequest)
)

var (
	once     sync.Once
	instance *FeatureRequest
)

type FeatureRequest struct {
	name       string
	inspection *inspection.Server
	ctx        *featureprovider.FeatureContext
}

func init() {
	featureprovider.RegisterFeatureFactory(
		ProviderName,
		func(ctx *featureprovider.FeatureContext) (featureprovider.Feature, error) {
			return initFeatureRequestInstance(ctx)
		},
	)
}

func initFeatureRequestInstance(ctx *featureprovider.FeatureContext) (featureprovider.Feature, error) {
	var err error
	once.Do(func() {
		instance = &FeatureRequest{
			name: ProviderName,
			ctx:  ctx,
		}
		instance.inspection, err = inspection.NewInspectionServer(ctx)
	})
	return instance, err
}

func (c *FeatureRequest) Equal(cluster *etcdv1alpha1.EtcdCluster) bool {
	return c.inspection.Equal(cluster, etcdv1alpha1.KStoneFeatureRequest)

... (content truncated for PRD) ...


#### 3 7 Inspection Module

### 3.7 Inspection Module


#### Pkg Inspection Inspection Go

#### pkg/inspection/inspection.go
**Note**: Core inspection server that:
- Manages inspection tasks for all features
- Creates and maintains etcd clients and watchers
- Synchronizes inspection resources with clusters
- Implements the following key methods:
  - `NewInspectionServer()`: Initializes the server
  - `Equal()`: Checks if inspection needs updating
  - `Sync()`: Creates/updates EtcdInspection resources
  - `GetEtcdClusterInfo()`: Retrieves cluster and config info


#### Pkg Inspection Alarm Go

#### pkg/inspection/alarm.go
```go
package inspection

import (
	"strconv"

	"k8s.io/klog/v2"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
	"etcd-operator/pkg/clusterprovider"
	featureutil "etcd-operator/pkg/featureprovider/util"
	"etcd-operator/pkg/inspection/metrics"
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

... (content truncated for PRD) ...


#### Pkg Inspection Consistency Go

#### pkg/inspection/consistency.go
**Note**: Implements consistency checking by:
- Collecting metadata from all etcd members
- Comparing key totals, revisions, and raft indices
- Calculating differences between members
- Exposing metrics for inconsistencies


#### Pkg Inspection Healthy Go

#### pkg/inspection/healthy.go
```go
package inspection

import (
	"k8s.io/klog/v2"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
	"etcd-operator/pkg/etcd"
	featureutil "etcd-operator/pkg/featureprovider/util"
	"etcd-operator/pkg/inspection/metrics"
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
```


#### Pkg Inspection Request Go

#### pkg/inspection/request.go
**Note**: Monitors etcd requests by:
- Creating watchers on specified key prefixes
- Tracking PUT, DELETE, and other operations
- Counting requests per resource type
- Using etcd watch API with event channels


#### Pkg Inspection Metrics Metrics Go

#### pkg/inspection/metrics/metrics.go
```go
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

... (content truncated for PRD) ...


#### 3 8 Kubernetes Utilities

### 3.8 Kubernetes Utilities


#### Pkg K8S Client Go

#### pkg/k8s/client.go
```go
package k8s

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	clientset "etcd-operator/pkg/generated/clientset/versioned"
	informers "etcd-operator/pkg/generated/informers/externalversions"
)

// GetClientConfig gets *rest.Config with the kube config
func GetClientConfig(kubeconfig string) (*rest.Config, error) {
	var cfg *rest.Config
	var err error
	if kubeconfig != "" {
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	} else {
		cfg, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

// GenerateInformer generates informer and client for controller
func GenerateInformer(config *rest.Config, labelSelector string) (
	*kubernetes.Clientset,
	*clientset.Clientset,
	kubeinformers.SharedInformerFactory,
	informers.SharedInformerFactory,
	error,
) {
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
		return nil, nil, nil, nil, err
	}


... (content truncated for PRD) ...


#### Pkg Signals Signal Go

#### pkg/signals/signal.go
```go
package signals

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	onlyOneSignalHandler = make(chan struct{})
	shutdownSignals      = []os.Signal{os.Interrupt, syscall.SIGTERM}
)

// SetupSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func SetupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
```

---


#### 4 Key Implementation Details

## 4. Key Implementation Details


#### 4 1 Plugin Registration Pattern

### 4.1 Plugin Registration Pattern

The codebase uses `init()` functions throughout to automatically register plugins when packages are imported:

```go
// Example from pkg/featureprovider/providers/alarm/alarm.go
func init() {
	featureprovider.RegisterFeatureFactory(
		ProviderName,
		func(ctx *featureprovider.FeatureContext) (featureprovider.Feature, error) {
			return initFeatureAlarmInstance(ctx)
		},
	)
}
```

This allows for easy extensibility - new features can be added by:
1. Creating a new provider package
2. Implementing the `Feature` interface
3. Registering via `init()`
4. Importing with `_ "package/path"` in the providers.go file


#### 4 2 Singleton Pattern With Sync Once

### 4.2 Singleton Pattern with sync.Once

Many providers use singleton pattern to ensure only one instance exists:

```go
var (
	once     sync.Once
	instance *FeatureAlarm
)

func initFeatureAlarmInstance(ctx *featureprovider.FeatureContext) (featureprovider.Feature, error) {
	var err error
	once.Do(func() {
		instance = &FeatureAlarm{
			name: ProviderName,
			ctx:  ctx,
		}
		instance.inspection, err = inspection.NewInspectionServer(ctx)
	})
	return instance, err
}
```


#### 4 3 Kubernetes Controller Pattern

### 4.3 Kubernetes Controller Pattern

Both controllers follow the standard Kubernetes controller pattern:

1. **Informers**: Watch and cache Kubernetes resources
2. **Work Queue**: Buffer work items with rate limiting
3. **Sync Handler**: Reconciliation logic
4. **Event Handler**: Enqueues items when resources change

```go
// Typical controller structure
controller := &Controller{
	kubeclientset:      kubeclientset,
	etcdclientset:      etcdclientset,
	etcdClustersLister: etcdInformerFactory.Etcd().V1alpha1().EtcdClusters().Lister(),
	etcdClustersSynced: etcdInformerFactory.Etcd().V1alpha1().EtcdClusters().Informer().HasSynced,
	workqueue:          workqueue.NewNamedRateLimitingQueue(...),
	recorder:           recorder,
}
```


#### 4 4 Prometheus Metrics Integration

### 4.4 Prometheus Metrics Integration

The inspection module exposes Prometheus metrics:

```go
// Metrics are defined as package variables
var EtcdEndpointHealthy = prometheus.NewGaugeVec(...)

// Registered in init()
func init() {
	prometheus.MustRegister(EtcdEndpointHealthy)
}

// Updated during inspection
metrics.EtcdEndpointHealthy.With(labels).Set(1)
```

The controller exposes these via HTTP:
```go
router := gin.Default()
router.GET("/metrics", gin.WrapH(promhttp.Handler()))
router.Run(":8080")
```


#### 4 5 Leader Election

### 4.5 Leader Election

Both controllers support leader election for high availability:

```go
lock := &resourcelock.LeaseLock{
	LeaseMeta: metav1.ObjectMeta{
		Name:      leaseLockName,
		Namespace: leaseLockNamespace,
	},
	Client: kubeClient.CoordinationV1(),
	LockConfig: resourcelock.ResourceLockConfig{
		Identity: id,
	},
}

leaderelection.RunOrDie(context.TODO(), leaderelection.LeaderElectionConfig{
	Lock:            lock,
	ReleaseOnCancel: true,
	LeaseDuration:   60 * time.Second,
	RenewDeadline:   15 * time.Second,
	RetryPeriod:     5 * time.Second,
	Callbacks: leaderelection.LeaderCallbacks{
		OnStartedLeading: func(ctx context.Context) {
			// Start controller
		},
		OnStoppedLeading: func() {
			// Clean up
		},
	},
})
```


#### 4 6 Etcd Client Versioning

### 4.6 etcd Client Versioning

The code supports both etcd v2 and v3 through abstraction:

```go
// Factory pattern for client creation
func GetEtcdClientProvider(
	name etcdv1alpha1.EtcdStorageBackend,
	ctx *VersionContext,
) (VersionClient, error) {
	// Returns v2 or v3 client based on storage backend
}

// Common interface
type VersionClient interface {
	MemberList() ([]Member, error)
	Status(endpoint string) (*Member, error)
	Close()
}
```


#### 4 7 Watch Based Request Tracking

### 4.7 Watch-based Request Tracking

The request inspection uses etcd's watch API:

```go
// Create watcher
watcher := clientv3.NewWatcher(client)
watchChan := watcher.Watch(ctx, keyPrefix, clientv3.WithPrefix())

// Monitor events
go func() {
	for event := range eventCh {
		// Track request metrics
		metrics.EtcdRequestTotal.With(labels).Inc()
	}
}()
```

---


#### 5 Dependencies

## 5. Dependencies


#### 5 1 Key External Dependencies

### 5.1 Key External Dependencies

From `go.mod`:
- **k8s.io/client-go**: Kubernetes client library
- **k8s.io/apimachinery**: Kubernetes API machinery
- **k8s.io/klog/v2**: Structured logging
- **go.etcd.io/etcd/client/v3**: etcd v3 client
- **go.etcd.io/etcd/client/v2**: etcd v2 client
- **github.com/prometheus/client_golang**: Prometheus metrics
- **github.com/gin-gonic/gin**: HTTP web framework
- **github.com/spf13/cobra**: CLI framework
- **sigs.k8s.io/controller-runtime**: Controller runtime utilities


#### 5 2 Code Generation

### 5.2 Code Generation

The project uses Kubernetes code generation tools:
- **deepcopy-gen**: Generates deep copy methods
- **client-gen**: Generates typed clients
- **informer-gen**: Generates informers
- **lister-gen**: Generates listers

These are reflected in the `pkg/generated/` directory structure.

---


#### 6 Deployment Architecture

## 6. Deployment Architecture


#### 6 1 Kubernetes Resources

### 6.1 Kubernetes Resources

The operator defines two Custom Resource Definitions (CRDs):
1. **EtcdCluster**: Represents an etcd cluster to be monitored
2. **EtcdInspection**: Represents an inspection task


#### 6 2 Controllers

### 6.2 Controllers

Two controllers run (typically as separate pods):
1. **EtcdCluster Controller**: Manages cluster lifecycle
2. **EtcdInspection Controller**: Performs monitoring tasks


#### 6 3 Metrics Endpoint

### 6.3 Metrics Endpoint

The inspection controller exposes metrics on port 8080:
- `/metrics`: Prometheus metrics endpoint
- Metrics include health, alarms, consistency, and request counts

---


#### 7 Summary

## 7. Summary

This etcd-metrics operator is a well-architected Kubernetes operator that:

1. **Monitors etcd clusters** through pluggable inspection features
2. **Exposes Prometheus metrics** for observability
3. **Uses clean abstractions** with interfaces and factory patterns
4. **Supports multiple etcd versions** (v2 and v3)
5. **Follows Kubernetes best practices** with controllers, informers, and CRDs
6. **Enables extensibility** through plugin architecture
7. **Provides high availability** with leader election support

The codebase demonstrates professional Go development practices with:
- Clear separation of concerns
- Dependency injection
- Interface-based design
- Proper error handling
- Structured logging
- Comprehensive metrics

This foundation can be used to build similar monitoring solutions for other distributed systems or to extend the etcd monitoring capabilities with additional features.


---

## 5. TECHNICAL REQUIREMENTS

### 5.1 Dependencies
- All dependencies from original documentation apply
- Standard development environment
- Required tools and libraries as specified

### 5.2 Compatibility
- Compatible with existing infrastructure
- Follows project standards and conventions

---

## 6. SUCCESS CRITERIA

### 6.1 Functional Success Criteria
- All identified tasks completed successfully
- All requirements implemented as specified
- All tests passing

### 6.2 Quality Success Criteria
- Code meets quality standards
- Documentation is complete and accurate
- No critical issues remaining

---

## 7. IMPLEMENTATION PLAN

### Phase 1: Preparation
- Review all requirements and tasks
- Set up development environment
- Gather necessary resources

### Phase 2: Implementation
- Execute tasks in priority order
- Follow best practices
- Test incrementally

### Phase 3: Validation
- Run comprehensive tests
- Validate against requirements
- Document completion

---

## 8. TASK-MASTER INTEGRATION

### How to Parse This PRD

```bash
# Parse this PRD with task-master
task-master parse-prd --input="{doc_name}_PRD.md"

# List generated tasks
task-master list

# Start execution
task-master next
```

### Expected Task Generation
Task-master should generate approximately {len(tasks)} tasks from this PRD.

---

## 9. APPENDIX

### 9.1 References
- Original document: {doc_name}.md
- Project: {project_name}

### 9.2 Change History
| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | {datetime.now().strftime('%Y-%m-%d')} | Initial PRD conversion |

---

*End of PRD*
*Generated by MD-to-PRD Converter*
