package internalinterfaces

import (
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	cache "k8s.io/client-go/tools/cache"
)

type NewInformerFunc func(runtime.Object, time.Duration) cache.SharedIndexInformer

type SharedInformerFactory interface {
	Start(stopCh <-chan struct{})
	InformerFor(obj runtime.Object, newFunc NewInformerFunc) cache.SharedIndexInformer
}

type TweakListOptionsFunc func(*v1.ListOptions)
