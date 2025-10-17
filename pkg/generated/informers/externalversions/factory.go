package externalversions

import (
	"reflect"
	"sync"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	cache "k8s.io/client-go/tools/cache"

	etcdv1alpha1 "github.com/etcd-monitor/taskmaster/api/etcd/v1alpha1"
	versioned "github.com/etcd-monitor/taskmaster/pkg/generated/clientset/versioned"
)

type SharedInformerFactory interface {
	Start(stopCh <-chan struct{})
	WaitForCacheSync(stopCh <-chan struct{}) map[reflect.Type]bool
	Etcd() etcdv1alpha1.Interface
}

type sharedInformerFactory struct {
	client           versioned.Interface
	namespace        string
	defaultResync    time.Duration
	customResync     map[reflect.Type]time.Duration
	lock             sync.Mutex
	informers        map[reflect.Type]cache.SharedIndexInformer
	startedInformers map[reflect.Type]bool
}

type SharedInformerOption func(*sharedInformerFactory)

func WithNamespace(namespace string) SharedInformerOption {
	return func(factory *sharedInformerFactory) {
		factory.namespace = namespace
	}
}

func NewSharedInformerFactory(client versioned.Interface, defaultResync time.Duration) SharedInformerFactory {
	return NewSharedInformerFactoryWithOptions(client, defaultResync)
}

func NewSharedInformerFactoryWithOptions(client versioned.Interface, defaultResync time.Duration, options ...SharedInformerOption) SharedInformerFactory {
	factory := &sharedInformerFactory{
		client:           client,
		namespace:        v1.NamespaceAll,
		defaultResync:    defaultResync,
		informers:        make(map[reflect.Type]cache.SharedIndexInformer),
		startedInformers: make(map[reflect.Type]bool),
		customResync:     make(map[reflect.Type]time.Duration),
	}

	for _, opt := range options {
		opt(factory)
	}

	return factory
}

func (f *sharedInformerFactory) Start(stopCh <-chan struct{}) {
	f.lock.Lock()
	defer f.lock.Unlock()

	for informerType, informer := range f.informers {
		if !f.startedInformers[informerType] {
			go informer.Run(stopCh)
			f.startedInformers[informerType] = true
		}
	}
}

func (f *sharedInformerFactory) WaitForCacheSync(stopCh <-chan struct{}) map[reflect.Type]bool {
	informers := func() map[reflect.Type]cache.SharedIndexInformer {
		f.lock.Lock()
		defer f.lock.Unlock()

		informers := map[reflect.Type]cache.SharedIndexInformer{}
		for informerType, informer := range f.informers {
			if f.startedInformers[informerType] {
				informers[informerType] = informer
			}
		}
		return informers
	}()

	res := map[reflect.Type]bool{}
	for informType, informer := range informers {
		res[informType] = cache.WaitForCacheSync(stopCh, informer.HasSynced)
	}
	return res
}

func (f *sharedInformerFactory) Etcd() etcdv1alpha1.Interface {
	// Create a simple mock informer for now
	return &mockEtcdInformer{}
}

type mockEtcdInformer struct{}

func (m *mockEtcdInformer) EtcdClusters() etcdv1alpha1.EtcdClusterInformer {
	return &mockEtcdClusterInformer{}
}

func (m *mockEtcdInformer) EtcdInspections() etcdv1alpha1.EtcdInspectionInformer {
	return &mockEtcdInspectionInformer{}
}

type mockEtcdClusterInformer struct{}

func (m *mockEtcdClusterInformer) Informer() cache.SharedIndexInformer {
	return nil
}

func (m *mockEtcdClusterInformer) Lister() etcdv1alpha1.EtcdClusterLister {
	return &mockEtcdClusterLister{}
}

type mockEtcdClusterLister struct{}

func (m *mockEtcdClusterLister) List(selector labels.Selector) (ret []*etcdv1alpha1.EtcdCluster, err error) {
	return nil, nil
}

func (m *mockEtcdClusterLister) EtcdClusters(namespace string) etcdv1alpha1.EtcdClusterNamespaceLister {
	return &mockEtcdClusterNamespaceLister{}
}

type mockEtcdClusterNamespaceLister struct{}

func (m *mockEtcdClusterNamespaceLister) List(selector labels.Selector) (ret []*etcdv1alpha1.EtcdCluster, err error) {
	return nil, nil
}

func (m *mockEtcdClusterNamespaceLister) Get(name string) (*etcdv1alpha1.EtcdCluster, error) {
	return nil, nil
}

type mockEtcdInspectionInformer struct{}

func (m *mockEtcdInspectionInformer) Informer() cache.SharedIndexInformer {
	return nil
}

func (m *mockEtcdInspectionInformer) Lister() etcdv1alpha1.EtcdInspectionLister {
	return &mockEtcdInspectionLister{}
}

type mockEtcdInspectionLister struct{}

func (m *mockEtcdInspectionLister) List(selector labels.Selector) (ret []*etcdv1alpha1.EtcdInspection, err error) {
	return nil, nil
}

func (m *mockEtcdInspectionLister) EtcdInspections(namespace string) etcdv1alpha1.EtcdInspectionNamespaceLister {
	return &mockEtcdInspectionNamespaceLister{}
}

type mockEtcdInspectionNamespaceLister struct{}

func (m *mockEtcdInspectionNamespaceLister) List(selector labels.Selector) (ret []*etcdv1alpha1.EtcdInspection, err error) {
	return nil, nil
}

func (m *mockEtcdInspectionNamespaceLister) Get(name string) (*etcdv1alpha1.EtcdInspection, error) {
	return nil, nil
}
