package v1alpha1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

type EtcdV1alpha1Interface interface {
	EtcdClusters(namespace string) EtcdClusterInterface
	EtcdInspections(namespace string) EtcdInspectionInterface
}

type EtcdClusterInterface interface {
	Create(ctx context.Context, etcdCluster *EtcdCluster, opts metav1.CreateOptions) (*EtcdCluster, error)
	Update(ctx context.Context, etcdCluster *EtcdCluster, opts metav1.UpdateOptions) (*EtcdCluster, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*EtcdCluster, error)
	List(ctx context.Context, opts metav1.ListOptions) (*EtcdClusterList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type EtcdInspectionInterface interface {
	Create(ctx context.Context, etcdInspection *EtcdInspection, opts metav1.CreateOptions) (*EtcdInspection, error)
	Update(ctx context.Context, etcdInspection *EtcdInspection, opts metav1.UpdateOptions) (*EtcdInspection, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*EtcdInspection, error)
	List(ctx context.Context, opts metav1.ListOptions) (*EtcdInspectionList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type etcdV1alpha1Client struct {
	restClient rest.Interface
}

func New(c rest.Interface) EtcdV1alpha1Interface {
	return &etcdV1alpha1Client{restClient: c}
}

func (c *etcdV1alpha1Client) EtcdClusters(namespace string) EtcdClusterInterface {
	return &etcdClustersClient{
		client:    c,
		namespace: namespace,
	}
}

func (c *etcdV1alpha1Client) EtcdInspections(namespace string) EtcdInspectionInterface {
	return &etcdInspectionsClient{
		client:    c,
		namespace: namespace,
	}
}

type etcdClustersClient struct {
	client    *etcdV1alpha1Client
	namespace string
}

func (c *etcdClustersClient) Create(ctx context.Context, etcdCluster *EtcdCluster, opts metav1.CreateOptions) (*EtcdCluster, error) {
	result := &EtcdCluster{}
	err := c.client.restClient.
		Post().
		Namespace(c.namespace).
		Resource("etcdclusters").
		VersionedParams(&opts, metav1.ParameterCodec).
		Body(etcdCluster).
		Do(ctx).
		Into(result)
	return result, err
}

func (c *etcdClustersClient) Update(ctx context.Context, etcdCluster *EtcdCluster, opts metav1.UpdateOptions) (*EtcdCluster, error) {
	result := &EtcdCluster{}
	err := c.client.restClient.
		Put().
		Namespace(c.namespace).
		Resource("etcdclusters").
		Name(etcdCluster.Name).
		VersionedParams(&opts, metav1.ParameterCodec).
		Body(etcdCluster).
		Do(ctx).
		Into(result)
	return result, err
}

func (c *etcdClustersClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.restClient.
		Delete().
		Namespace(c.namespace).
		Resource("etcdclusters").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *etcdClustersClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*EtcdCluster, error) {
	result := &EtcdCluster{}
	err := c.client.restClient.
		Get().
		Namespace(c.namespace).
		Resource("etcdclusters").
		Name(name).
		VersionedParams(&opts, metav1.ParameterCodec).
		Do(ctx).
		Into(result)
	return result, err
}

func (c *etcdClustersClient) List(ctx context.Context, opts metav1.ListOptions) (*EtcdClusterList, error) {
	result := &EtcdClusterList{}
	err := c.client.restClient.
		Get().
		Namespace(c.namespace).
		Resource("etcdclusters").
		VersionedParams(&opts, metav1.ParameterCodec).
		Do(ctx).
		Into(result)
	return result, err
}

func (c *etcdClustersClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.restClient.
		Get().
		Namespace(c.namespace).
		Resource("etcdclusters").
		VersionedParams(&opts, metav1.ParameterCodec).
		Watch(ctx)
}

type etcdInspectionsClient struct {
	client    *etcdV1alpha1Client
	namespace string
}

func (c *etcdInspectionsClient) Create(ctx context.Context, etcdInspection *EtcdInspection, opts metav1.CreateOptions) (*EtcdInspection, error) {
	result := &EtcdInspection{}
	err := c.client.restClient.
		Post().
		Namespace(c.namespace).
		Resource("etcdinspections").
		VersionedParams(&opts, metav1.ParameterCodec).
		Body(etcdInspection).
		Do(ctx).
		Into(result)
	return result, err
}

func (c *etcdInspectionsClient) Update(ctx context.Context, etcdInspection *EtcdInspection, opts metav1.UpdateOptions) (*EtcdInspection, error) {
	result := &EtcdInspection{}
	err := c.client.restClient.
		Put().
		Namespace(c.namespace).
		Resource("etcdinspections").
		Name(etcdInspection.Name).
		VersionedParams(&opts, metav1.ParameterCodec).
		Body(etcdInspection).
		Do(ctx).
		Into(result)
	return result, err
}

func (c *etcdInspectionsClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.restClient.
		Delete().
		Namespace(c.namespace).
		Resource("etcdinspections").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *etcdInspectionsClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*EtcdInspection, error) {
	result := &EtcdInspection{}
	err := c.client.restClient.
		Get().
		Namespace(c.namespace).
		Resource("etcdinspections").
		Name(name).
		VersionedParams(&opts, metav1.ParameterCodec).
		Do(ctx).
		Into(result)
	return result, err
}

func (c *etcdInspectionsClient) List(ctx context.Context, opts metav1.ListOptions) (*EtcdInspectionList, error) {
	result := &EtcdInspectionList{}
	err := c.client.restClient.
		Get().
		Namespace(c.namespace).
		Resource("etcdinspections").
		VersionedParams(&opts, metav1.ParameterCodec).
		Do(ctx).
		Into(result)
	return result, err
}

func (c *etcdInspectionsClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.restClient.
		Get().
		Namespace(c.namespace).
		Resource("etcdinspections").
		VersionedParams(&opts, metav1.ParameterCodec).
		Watch(ctx)
}

func NewForConfig(c *rest.Config) (EtcdV1alpha1Interface, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return New(client), nil
}

func NewForConfigAndClient(c *rest.Config, h rest.Interface) (EtcdV1alpha1Interface, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	return New(h), nil
}

func setConfigDefaults(config *rest.Config) error {
	gv := SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = runtime.NewSimpleNegotiatedSerializer(runtime.SerializerInfo{})
	return nil
}
