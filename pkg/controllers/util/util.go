package util

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	"github.com/etcd-monitor/taskmaster/pkg/k8s"
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

// EventRecorder is a simple event recorder interface
type EventRecorder interface {
	Event(object runtime.Object, eventtype, reason, message string)
	Eventf(object runtime.Object, eventtype, reason, messageFmt string, args ...interface{})
	AnnotatedEventf(object runtime.Object, annotations map[string]string, eventtype, reason, messageFmt string, args ...interface{})
}

// eventRecorder implements EventRecorder interface
type eventRecorder struct {
}

// NewEventRecorder creates a new event recorder
func NewEventRecorder(client clientset.Interface) EventRecorder {
	// For now, return a no-op recorder to avoid compilation issues
	return &eventRecorder{}
}

func (e *eventRecorder) Event(object runtime.Object, eventtype, reason, message string) {
	// No-op implementation for now
	klog.V(2).Infof("Event: %s %s %s", eventtype, reason, message)
}

func (e *eventRecorder) Eventf(object runtime.Object, eventtype, reason, messageFmt string, args ...interface{}) {
	// No-op implementation for now
	klog.V(2).Infof("Eventf: %s %s %s", eventtype, reason, fmt.Sprintf(messageFmt, args...))
}

func (e *eventRecorder) AnnotatedEventf(object runtime.Object, annotations map[string]string, eventtype, reason, messageFmt string, args ...interface{}) {
	// No-op implementation for now
	klog.V(2).Infof("AnnotatedEventf: %s %s %s", eventtype, reason, fmt.Sprintf(messageFmt, args...))
}

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
