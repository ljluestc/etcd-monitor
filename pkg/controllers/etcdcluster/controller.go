package etcdcluster

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	etcdv1alpha1 "github.com/etcd-monitor/taskmaster/api/etcd/v1alpha1"
	"github.com/etcd-monitor/taskmaster/pkg/controllers/util"
	clientset "github.com/etcd-monitor/taskmaster/pkg/generated/clientset/versioned"
	etcdinformers "github.com/etcd-monitor/taskmaster/pkg/generated/informers/externalversions"
)

// Controller is the controller implementation for EtcdCluster resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// etcdclientset is a clientset for our own API group
	etcdclientset clientset.Interface

	etcdClustersLister etcdv1alpha1.EtcdClusterLister
	etcdClustersSynced cache.InformerSynced

	workqueue workqueue.RateLimitingInterface
	recorder  util.EventRecorder
	clientBuilder util.ClientBuilder
}

// NewController returns a new etcdcluster controller
func NewController(
	kubeclientset kubernetes.Interface,
	etcdclientset clientset.Interface,
	kubeInformerFactory informers.SharedInformerFactory,
	etcdInformerFactory etcdinformers.SharedInformerFactory,
	clientBuilder util.ClientBuilder,
) *Controller {
	klog.V(4).Info("Creating EtcdCluster controller")

	// Obtain references to shared informers
	etcdClusterInformer := etcdInformerFactory.Etcd().EtcdClusters()

	controller := &Controller{
		kubeclientset:      kubeclientset,
		etcdclientset:      etcdclientset,
		etcdClustersLister: etcdClusterInformer.Lister(),
		etcdClustersSynced: etcdClusterInformer.Informer().HasSynced,
		workqueue:          workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "EtcdClusters"),
		recorder:           util.NewEventRecorder(kubeclientset),
		clientBuilder:      clientBuilder,
	}

	klog.Info("Setting up event handlers for EtcdCluster")

	etcdClusterInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueEtcdCluster,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueEtcdCluster(new)
		},
		DeleteFunc: controller.enqueueEtcdClusterForDelete,
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	klog.Info("Starting EtcdCluster controller")

	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.etcdClustersSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message off the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here to avoid this item ever being worked again by
			// processing it.
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// EtcdCluster resource to be synced.
		if err := c.syncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// If no error occurs we Forget this item so it will not be retried again.
		c.workqueue.Forget(obj)
		klog.V(4).Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the EtcdCluster resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the EtcdCluster resource with this namespace/name
	etcdCluster, err := c.etcdClustersLister.EtcdClusters(namespace).Get(name)
	if err != nil {
		// The EtcdCluster resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			klog.V(4).Infof("etcdCluster '%s' in work queue no longer exists", key)
			return nil
		}
		return err
	}

	// Here you would implement the actual reconciliation logic for your EtcdCluster
	// For now, we'll just log that we're processing it.
	klog.V(2).Infof("Processing EtcdCluster: %s/%s", etcdCluster.Namespace, etcdCluster.Name)

	// TODO: Implement actual cluster management logic (create, update, delete, status)

	return nil
}

// enqueueEtcdCluster takes an EtcdCluster resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed objects of invalid type.
func (c *Controller) enqueueEtcdCluster(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

// enqueueEtcdClusterForDelete takes an EtcdCluster resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed objects of invalid type.
func (c *Controller) enqueueEtcdClusterForDelete(obj interface{}) {
	var key string
	var err error
	tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
	if !ok {
		runtime.HandleError(fmt.Errorf("error decoding object, invalid type"))
		return
	}
	key, err = cache.MetaNamespaceKeyFunc(tombstone.Obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}