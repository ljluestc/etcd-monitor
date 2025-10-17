package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
)

type Interface interface {
	EtcdClusters() EtcdClusterInformer
	EtcdInspections() EtcdInspectionInformer
}

type EtcdClusterInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() EtcdClusterLister
}

type EtcdInspectionInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() EtcdInspectionLister
}

type EtcdClusterLister interface {
	List(selector labels.Selector) (ret []*EtcdCluster, err error)
	EtcdClusters(namespace string) EtcdClusterNamespaceLister
}

type EtcdClusterNamespaceLister interface {
	List(selector labels.Selector) (ret []*EtcdCluster, err error)
	Get(name string) (*EtcdCluster, error)
}

type EtcdInspectionLister interface {
	List(selector labels.Selector) (ret []*EtcdInspection, err error)
	EtcdInspections(namespace string) EtcdInspectionNamespaceLister
}

type EtcdInspectionNamespaceLister interface {
	List(selector labels.Selector) (ret []*EtcdInspection, err error)
	Get(name string) (*EtcdInspection, error)
}

func NewInformer(factory cache.SharedIndexInformer, namespace string) Interface {
	return &version{factory: factory, namespace: namespace}
}

type version struct {
	factory   cache.SharedIndexInformer
	namespace string
}

func (v *version) EtcdClusters() EtcdClusterInformer {
	return &etcdClusterInformer{factory: v.factory, namespace: v.namespace}
}

func (v *version) EtcdInspections() EtcdInspectionInformer {
	return &etcdInspectionInformer{factory: v.factory, namespace: v.namespace}
}

type etcdClusterInformer struct {
	factory   cache.SharedIndexInformer
	namespace string
}

func (f *etcdClusterInformer) Informer() cache.SharedIndexInformer {
	return f.factory
}

func (f *etcdClusterInformer) Lister() EtcdClusterLister {
	return &etcdClusterLister{factory: f.factory, namespace: f.namespace}
}

type etcdClusterLister struct {
	factory   cache.SharedIndexInformer
	namespace string
}

func (l *etcdClusterLister) List(selector labels.Selector) (ret []*EtcdCluster, err error) {
	err = cache.ListAll(l.factory.GetIndexer(), selector, func(m interface{}) {
		ret = append(ret, m.(*EtcdCluster))
	})
	return ret, err
}

func (l *etcdClusterLister) EtcdClusters(namespace string) EtcdClusterNamespaceLister {
	return &etcdClusterNamespaceLister{factory: l.factory, namespace: namespace}
}

type etcdClusterNamespaceLister struct {
	factory   cache.SharedIndexInformer
	namespace string
}

func (l *etcdClusterNamespaceLister) List(selector labels.Selector) (ret []*EtcdCluster, err error) {
	err = cache.ListAllByNamespace(l.factory.GetIndexer(), l.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*EtcdCluster))
	})
	return ret, err
}

func (l *etcdClusterNamespaceLister) Get(name string) (*EtcdCluster, error) {
	obj, exists, err := l.factory.GetIndexer().GetByKey(l.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(schema.GroupResource{Resource: "etcdclusters"}, name)
	}
	return obj.(*EtcdCluster), nil
}

type etcdInspectionInformer struct {
	factory   cache.SharedIndexInformer
	namespace string
}

func (f *etcdInspectionInformer) Informer() cache.SharedIndexInformer {
	return f.factory
}

func (f *etcdInspectionInformer) Lister() EtcdInspectionLister {
	return &etcdInspectionLister{factory: f.factory, namespace: f.namespace}
}

type etcdInspectionLister struct {
	factory   cache.SharedIndexInformer
	namespace string
}

func (l *etcdInspectionLister) List(selector labels.Selector) (ret []*EtcdInspection, err error) {
	err = cache.ListAll(l.factory.GetIndexer(), selector, func(m interface{}) {
		ret = append(ret, m.(*EtcdInspection))
	})
	return ret, err
}

func (l *etcdInspectionLister) EtcdInspections(namespace string) EtcdInspectionNamespaceLister {
	return &etcdInspectionNamespaceLister{factory: l.factory, namespace: namespace}
}

type etcdInspectionNamespaceLister struct {
	factory   cache.SharedIndexInformer
	namespace string
}

func (l *etcdInspectionNamespaceLister) List(selector labels.Selector) (ret []*EtcdInspection, err error) {
	err = cache.ListAllByNamespace(l.factory.GetIndexer(), l.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*EtcdInspection))
	})
	return ret, err
}

func (l *etcdInspectionNamespaceLister) Get(name string) (*EtcdInspection, error) {
	obj, exists, err := l.factory.GetIndexer().GetByKey(l.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(schema.GroupResource{Resource: "etcdinspections"}, name)
	}
	return obj.(*EtcdInspection), nil
}
