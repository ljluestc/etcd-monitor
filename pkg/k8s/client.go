package k8s

import (
	"time"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	clientset "github.com/etcd-monitor/taskmaster/pkg/generated/clientset/versioned"
	informers "github.com/etcd-monitor/taskmaster/pkg/generated/informers/externalversions"
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
	etcdInformerFactory := informers.NewSharedInformerFactory(clustetClient, time.Second*30)

	return kubeClient, clustetClient, kubeInformerFactory, etcdInformerFactory, nil
}
