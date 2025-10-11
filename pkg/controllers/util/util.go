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
