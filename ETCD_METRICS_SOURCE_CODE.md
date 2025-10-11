# Complete Source Code Analysis: etcd-metrics Repository

## Repository Information
- **Repository**: https://github.com/clay-wangzhi/etcd-metrics
- **Purpose**: Kubernetes Operator for monitoring etcd clusters
- **License**: Apache License 2.0
- **Author**: Clay Wang

---

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
│   │   │   │   └── alarm.go
│   │   │   ├── consistency/
│   │   │   │   └── consistency.go
│   │   │   ├── healthy/
│   │   │   │   └── healthy.go
│   │   │   ├── request/
│   │   │   │   └── request.go
│   │   │   └── providers.go
│   │   ├── util/
│   │   │   └── util.go
│   │   ├── feature.go
│   │   └── plugins.go
│   ├── generated/
│   │   ├── clientset/versioned/
│   │   ├── informers/externalversions/
│   │   └── listers/etcd/v1alpha1/
│   ├── inspection/
│   │   ├── metrics/
│   │   │   └── metrics.go
│   │   ├── alarm.go
│   │   ├── consistency.go
│   │   ├── healthy.go
│   │   ├── inspection.go
│   │   └── request.go
│   ├── k8s/
│   │   └── client.go
│   └── signals/
│       └── signal.go
├── config/
├── deploy/
├── hack/
├── vendor/
├── Dockerfile
├── Makefile
├── go.mod
├── go.sum
├── PROJECT
└── README.md
```

---

## 2. Core Architecture Patterns

### 2.1 Controller Pattern
The project implements the Kubernetes Operator pattern with two main controllers:
1. **EtcdCluster Controller**: Manages the lifecycle of etcd clusters
2. **EtcdInspection Controller**: Performs health checks and monitoring tasks

### 2.2 Plugin Architecture
The codebase uses a plugin-based factory pattern for extensibility:
- **Cluster Providers**: Different etcd cluster implementations (e.g., imported clusters)
- **Feature Providers**: Pluggable inspection features (alarm, consistency, healthy, request)
- **Client Versions**: Support for different etcd storage backends (v2, v3)

### 2.3 Key Design Patterns
- **Factory Pattern**: For creating providers and clients
- **Singleton Pattern**: For feature providers using `sync.Once`
- **Observer Pattern**: Using Kubernetes informers and watchers
- **Strategy Pattern**: Different inspection strategies via feature providers

---

## 3. Complete Source Code Files

### 3.1 API Definitions

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
	Phase       string        `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase"`
	ServiceName string        `json:"serviceName,omitempty" protobuf:"bytes,2,opt,name=serviceName"`
	Members     []EtcdMember  `json:"members,omitempty" protobuf:"bytes,3,rep,name=members"`
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
```

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

// EtcdInspectionList is a list of EtcdInspection resources
type EtcdInspectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []EtcdInspection `json:"items" protobuf:"bytes,2,rep,name=items"`
}
```

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

### 3.2 Command Entry Points

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
	defer klog.Flush()

	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	err := flags.Parse(os.Args[1:])
	if err != nil {
		klog.Errorf("failed to parse args, err is %v", err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```

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

// NewEtcdClusterControllerCommand creates a *cobra.Command object with default parameters
func NewEtcdClusterControllerCommand(out io.Writer) *cobra.Command {
	cc := &EtcdClusterCommand{out: out}
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "run cluster controller",
		Long: `The cluster controller is a daemon, it will watches the changes of etcdcluster resources
 through the apiserver and makes changes attempting to move the current state towards the desired state.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				klog.V(1).Infof("FLAG: --%s=%q", flag.Name, flag.Value)
			})
			if err := cc.run(); err != nil {
				return err
			}
			return nil
		},
	}
	// Add command flags
	cmd.Flags().StringVar(&cc.kubeconfig, "kubeconfig", "", "Path to a kubeconfig file")
	cmd.Flags().StringVar(&cc.masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig.")
	cmd.Flags().StringVar(&cc.labelSelector, "label-selector", "", "Label selector for filtering resources")
	cmd.Flags().StringVar(&cc.leaseLockName, "lease-lock-name", "", "Name of the lease lock")
	cmd.Flags().StringVar(&cc.leaseLockNamespace, "lease-lock-namespace", "", "Namespace of the lease lock")
	return cmd
}

func (c *EtcdClusterCommand) run() error {
	// Set up signals to gracefully shutdown
	stopCh := signals.SetupSignalHandler()

	cfg, err := k8s.GetClientConfig(c.kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, etcdClient, kubeInformerFactory, etcdInformerFactory, err := k8s.GenerateInformer(cfg, c.labelSelector)
	if err != nil {
		klog.Fatalf("Error building clientset: %s", err.Error())
	}

	controller := etcdcluster.NewController(
		kubeClient,
		etcdClient,
		kubeInformerFactory,
		etcdInformerFactory,
		util.NewSimpleClientBuilder(c.kubeconfig),
	)

	// Start informer factories
	kubeInformerFactory.Start(stopCh)
	etcdInformerFactory.Start(stopCh)

	// Leader election setup
	if c.leaseLockName != "" && c.leaseLockNamespace != "" {
		id := string(uuid.NewUUID())
		lock := &resourcelock.LeaseLock{
			LeaseMeta: metav1.ObjectMeta{
				Name:      c.leaseLockName,
				Namespace: c.leaseLockNamespace,
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
					if err = controller.Run(2, stopCh); err != nil {
						klog.Fatalf("Error running controller: %s", err.Error())
					}
				},
				OnStoppedLeading: func() {
					klog.Infof("leader lost: %s", id)
					os.Exit(0)
				},
				OnNewLeader: func(identity string) {
					if identity == id {
						return
					}
					klog.Infof("new leader elected: %s", identity)
				},
			},
		})
	} else {
		if err = controller.Run(2, stopCh); err != nil {
			klog.Fatalf("Error running controller: %s", err.Error())
		}
	}

	<-stopCh
	return nil
}
```

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
			return nil
		},
	}
	// Add command flags
	cmd.Flags().StringVar(&cc.kubeconfig, "kubeconfig", "", "Path to a kubeconfig file")
	cmd.Flags().StringVar(&cc.masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig.")
	cmd.Flags().StringVar(&cc.labelSelector, "label-selector", "", "Label selector for filtering resources")
	cmd.Flags().StringVar(&cc.leaseLockName, "lease-lock-name", "", "Name of the lease lock")
	cmd.Flags().StringVar(&cc.leaseLockNamespace, "lease-lock-namespace", "", "Namespace of the lease lock")
	return cmd
}

func (c *EtcdInspectionCommand) run() error {
	// Set up signals to gracefully shutdown
	stopCh := signals.SetupSignalHandler()

	cfg, err := k8s.GetClientConfig(c.kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, etcdClient, kubeInformerFactory, etcdInformerFactory, err := k8s.GenerateInformer(cfg, c.labelSelector)
	if err != nil {
		klog.Fatalf("Error building clientset: %s", err.Error())
	}

	controller := etcdinspection.NewController(
		kubeClient,
		etcdClient,
		kubeInformerFactory,
		etcdInformerFactory,
		util.NewSimpleClientBuilder(c.kubeconfig),
	)

	// Start informer factories
	kubeInformerFactory.Start(stopCh)
	etcdInformerFactory.Start(stopCh)

	// Leader election setup
	if c.leaseLockName != "" && c.leaseLockNamespace != "" {
		id := string(uuid.NewUUID())
		lock := &resourcelock.LeaseLock{
			LeaseMeta: metav1.ObjectMeta{
				Name:      c.leaseLockName,
				Namespace: c.leaseLockNamespace,
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
					if err = controller.Run(2, stopCh); err != nil {
						klog.Fatalf("Error running controller: %s", err.Error())
					}
				},
				OnStoppedLeading: func() {
					klog.Infof("leader lost: %s", id)
					os.Exit(0)
				},
				OnNewLeader: func(identity string) {
					if identity == id {
						return
					}
					klog.Infof("new leader elected: %s", identity)
				},
			},
		})
	} else {
		if err = controller.Run(2, stopCh); err != nil {
			klog.Fatalf("Error running controller: %s", err.Error())
		}
	}

	<-stopCh
	return nil
}
```

---

### 3.3 Controllers

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

#### pkg/controllers/etcdinspection/etcdinspection-controller.go
**Note**: Similar structure to the cluster controller but focused on inspection tasks. It:

- Watches EtcdInspection custom resources
- Triggers inspection tasks based on feature providers
- Exposes Prometheus metrics via HTTP endpoint
- Uses Gin web framework for the metrics server

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
	}
	return builder
}

type simpleClientBuilder struct {
	kubeconfig string
}

func (b *simpleClientBuilder) ConfigOrDie() *restclient.Config {
	cfg, err := k8s.GetClientConfig(b.kubeconfig)
	if err != nil {
		panic(err)
	}
	return cfg
}

func (b *simpleClientBuilder) ClientOrDie() clientset.Interface {
	return clientset.NewForConfigOrDie(b.ConfigOrDie())
}

func (b *simpleClientBuilder) DynamicClientOrDie() dynamic.Interface {
	return dynamic.NewForConfigOrDie(b.ConfigOrDie())
}
```

---

### 3.4 Cluster Provider

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

#### pkg/clusterprovider/helper.go
**Note**: Contains helper functions for:
- Getting storage member endpoints
- Populating extension client URL maps
- Getting runtime etcd members
- Getting etcd alarms
- Checking member list equality

#### pkg/clusterprovider/providers/providers.go
```go
package providers

import (
	_ "etcd-operator/pkg/clusterprovider/providers/imported" // import imported provider
)
```

#### pkg/clusterprovider/providers/imported/cluster.go
**Note**: Implements the Cluster interface for imported etcd clusters (clusters not created by the operator). Key methods:
- Uses singleton pattern with `sync.Once`
- Registers itself via `init()` function
- Implements lifecycle hooks (BeforeCreate, Create, AfterCreate, etc.)
- Retrieves cluster status from existing etcd endpoints

---

### 3.5 Etcd Client and Utilities

#### pkg/etcd/client.go
**Note**: Provides client configuration management:
- `ClientConfig` struct with endpoints, timeouts, TLS settings
- `SecureConfig` for authentication credentials
- `ClientConfigSecret` for retrieving credentials from Kubernetes secrets
- Functions to create etcd v2 and v3 clients

#### pkg/etcd/health.go
**Note**: Implements health checking:
- HTTP-based health check client
- TLS configuration support
- `/health` endpoint monitoring
- `/version` endpoint for version detection
- `/v2/stats/self` for etcd statistics

#### pkg/etcd/helper.go
**Note**: Utility functions for:
- Creating etcd v2 and v3 clients
- Member list operations
- Status checking
- Certificate and authentication management

#### pkg/etcd/stats.go
**Note**: Statistics collection interface:
- `Stat` interface for v2 and v3 backends
- `GetTotalKeyNum()`: Counts keys in etcd
- `GetIndex()`: Retrieves metadata indices for consistency checking
- Supports both v2 and v3 storage backends

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
	}
	return f(ctx)
}
```

#### pkg/etcd/client/versions/providers.go
```go
package versions

import (
	_ "etcd-operator/pkg/etcd/client/versions/v3" // import etcd client of v3
)
```

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

func (c *V3) Close() {
	c.cli.Close()
}

func init() {
	client.RegisterEtcdClientFactory(etcdv1alpha1.EtcdStorageV3, NewV3Client)
}

func NewV3Client(ctx *client.VersionContext) (client.VersionClient, error) {
	cli, err := etcd.NewClientv3(ctx.Config)
	if err != nil {
		return nil, err
	}
	return &V3{
		ctx: ctx,
		cli: cli,
	}, nil
}
```

---

### 3.6 Feature Provider

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
		features = append(features, feature)
	}
	return features
}
```

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
	labels := map[string]string{
		"clusterName":    clusterName,
		"inspectionType": string(featureName),
	}
	fmt.Print(labels)
	// metrics.EtcdInspectionFailedNum.With(labels).Inc()
}
```

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
}

func (c *FeatureAlarm) Sync(cluster *etcdv1alpha1.EtcdCluster) error {
	return c.inspection.Sync(cluster, etcdv1alpha1.KStoneFeatureAlarm)
}

func (c *FeatureAlarm) Do(inspection *etcdv1alpha1.EtcdInspection) error {
	return c.inspection.CollectAlarmList(inspection)
}
```

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
}

func (c *FeatureConsistency) Sync(cluster *etcdv1alpha1.EtcdCluster) error {
	return c.inspection.Sync(cluster, etcdv1alpha1.KStoneFeatureConsistency)
}

func (c *FeatureConsistency) Do(inspection *etcdv1alpha1.EtcdInspection) error {
	return c.inspection.CollectClusterConsistentData(inspection)
}
```

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
}

func (c *FeatureHealthy) Sync(cluster *etcdv1alpha1.EtcdCluster) error {
	return c.inspection.Sync(cluster, etcdv1alpha1.KStoneFeatureHealthy)
}

func (c *FeatureHealthy) Do(inspection *etcdv1alpha1.EtcdInspection) error {
	return c.inspection.CollectMemberHealthy(inspection)
}
```

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
}

func (c *FeatureRequest) Sync(cluster *etcdv1alpha1.EtcdCluster) error {
	return c.inspection.Sync(cluster, etcdv1alpha1.KStoneFeatureRequest)
}

func (c *FeatureRequest) Do(inspection *etcdv1alpha1.EtcdInspection) error {
	return c.inspection.CollectEtcdClusterRequest(inspection)
}
```

---

### 3.7 Inspection Module

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
```

#### pkg/inspection/consistency.go
**Note**: Implements consistency checking by:
- Collecting metadata from all etcd members
- Comparing key totals, revisions, and raft indices
- Calculating differences between members
- Exposing metrics for inconsistencies

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

#### pkg/inspection/request.go
**Note**: Monitors etcd requests by:
- Creating watchers on specified key prefixes
- Tracking PUT, DELETE, and other operations
- Counting requests per resource type
- Using etcd watch API with event channels

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
```

---

### 3.8 Kubernetes Utilities

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

	clustetClient, err := clientset.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error building example clientset: %s", err.Error())
		return nil, nil, nil, nil, err
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	etcdInformerFactory := informers.NewSharedInformerFactoryWithOptions(
		clustetClient,
		time.Second*30,
		informers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.LabelSelector = labelSelector
		}),
	)

	return kubeClient, clustetClient, kubeInformerFactory, etcdInformerFactory, nil
}
```

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

## 4. Key Implementation Details

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

## 5. Dependencies

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

### 5.2 Code Generation

The project uses Kubernetes code generation tools:
- **deepcopy-gen**: Generates deep copy methods
- **client-gen**: Generates typed clients
- **informer-gen**: Generates informers
- **lister-gen**: Generates listers

These are reflected in the `pkg/generated/` directory structure.

---

## 6. Deployment Architecture

### 6.1 Kubernetes Resources

The operator defines two Custom Resource Definitions (CRDs):
1. **EtcdCluster**: Represents an etcd cluster to be monitored
2. **EtcdInspection**: Represents an inspection task

### 6.2 Controllers

Two controllers run (typically as separate pods):
1. **EtcdCluster Controller**: Manages cluster lifecycle
2. **EtcdInspection Controller**: Performs monitoring tasks

### 6.3 Metrics Endpoint

The inspection controller exposes metrics on port 8080:
- `/metrics`: Prometheus metrics endpoint
- Metrics include health, alarms, consistency, and request counts

---

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
