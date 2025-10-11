package inspection

import (
	"context"
	"fmt"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
	"etcd-operator/pkg/controllers/util"
	"etcd-operator/pkg/etcd"
	"etcd-operator/pkg/featureprovider"
	featureutil "etcd-operator/pkg/featureprovider/util"
	clientset "etcd-operator/pkg/generated/clientset/versioned"
)

// Server manages inspection tasks for etcd clusters
type Server struct {
	ctx                *featureprovider.FeatureContext
	kubeClient         kubernetes.Interface
	etcdClient         clientset.Interface
	inspectionCache    map[string]*inspectionCacheEntry
	inspectionCacheMux sync.RWMutex
}

// inspectionCacheEntry caches inspection resources per cluster
type inspectionCacheEntry struct {
	cluster          *etcdv1alpha1.EtcdCluster
	clientConfig     *etcd.ClientConfig
	inspectionStatus map[etcdv1alpha1.KStoneFeature]bool
}

// NewInspectionServer creates a new inspection server
func NewInspectionServer(ctx *featureprovider.FeatureContext) (*Server, error) {
	kubeClient := ctx.ClientBuilder.ClientOrDie()
	cfg := ctx.ClientBuilder.ConfigOrDie()

	etcdClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd clientset: %v", err)
	}

	server := &Server{
		ctx:             ctx,
		kubeClient:      kubeClient,
		etcdClient:      etcdClient,
		inspectionCache: make(map[string]*inspectionCacheEntry),
	}

	klog.V(2).Info("inspection server initialized")
	return server, nil
}

// Equal checks whether the inspection needs to be updated
func (s *Server) Equal(cluster *etcdv1alpha1.EtcdCluster, feature etcdv1alpha1.KStoneFeature) bool {
	s.inspectionCacheMux.RLock()
	defer s.inspectionCacheMux.RUnlock()

	cacheKey := s.getCacheKey(cluster.Namespace, cluster.Name)
	entry, exists := s.inspectionCache[cacheKey]

	if !exists {
		klog.V(2).Infof("cluster %s/%s not in cache", cluster.Namespace, cluster.Name)
		return false
	}

	// Check if feature is enabled in annotations
	featureEnabled := featureutil.IsFeatureGateEnabled(cluster.Annotations, feature)

	// Check if inspection status matches
	inspectionExists, ok := entry.inspectionStatus[feature]
	if !ok {
		klog.V(2).Infof("feature %s not tracked for cluster %s/%s", feature, cluster.Namespace, cluster.Name)
		return false
	}

	// Both should be enabled or both should be disabled
	equal := featureEnabled == inspectionExists

	klog.V(3).Infof("equal check for cluster %s/%s, feature %s: enabled=%v, exists=%v, equal=%v",
		cluster.Namespace, cluster.Name, feature, featureEnabled, inspectionExists, equal)

	return equal
}

// Sync synchronizes the latest feature configuration
func (s *Server) Sync(cluster *etcdv1alpha1.EtcdCluster, feature etcdv1alpha1.KStoneFeature) error {
	namespace, name := cluster.Namespace, cluster.Name
	featureEnabled := featureutil.IsFeatureGateEnabled(cluster.Annotations, feature)

	klog.V(2).Infof("syncing cluster %s/%s, feature %s, enabled=%v", namespace, name, feature, featureEnabled)

	inspectionName := fmt.Sprintf("%s-%s", name, feature)

	if featureEnabled {
		// Create or update EtcdInspection resource
		inspection := &etcdv1alpha1.EtcdInspection{
			ObjectMeta: metav1.ObjectMeta{
				Name:      inspectionName,
				Namespace: namespace,
				Labels: map[string]string{
					"cluster":        name,
					"inspectionType": string(feature),
				},
			},
			Spec: etcdv1alpha1.EtcdInspectionSpec{
				ClusterName:    name,
				InspectionType: string(feature),
			},
		}

		// Try to get existing inspection
		existingInspection, err := s.etcdClient.EtcdV1alpha1().EtcdInspections(namespace).Get(
			context.TODO(),
			inspectionName,
			metav1.GetOptions{},
		)

		if err == nil {
			// Update existing inspection
			existingInspection.Spec = inspection.Spec
			_, err = s.etcdClient.EtcdV1alpha1().EtcdInspections(namespace).Update(
				context.TODO(),
				existingInspection,
				metav1.UpdateOptions{},
			)
			if err != nil {
				klog.Errorf("failed to update inspection %s/%s: %v", namespace, inspectionName, err)
				return err
			}
			klog.V(2).Infof("updated inspection %s/%s", namespace, inspectionName)
		} else {
			// Create new inspection
			_, err = s.etcdClient.EtcdV1alpha1().EtcdInspections(namespace).Create(
				context.TODO(),
				inspection,
				metav1.CreateOptions{},
			)
			if err != nil {
				klog.Errorf("failed to create inspection %s/%s: %v", namespace, inspectionName, err)
				return err
			}
			klog.V(2).Infof("created inspection %s/%s", namespace, inspectionName)
		}

		// Update cache
		s.updateInspectionCache(cluster, feature, true)
	} else {
		// Delete EtcdInspection resource
		err := s.etcdClient.EtcdV1alpha1().EtcdInspections(namespace).Delete(
			context.TODO(),
			inspectionName,
			metav1.DeleteOptions{},
		)
		if err != nil {
			klog.Warningf("failed to delete inspection %s/%s: %v", namespace, inspectionName, err)
			// Don't return error if resource doesn't exist
		} else {
			klog.V(2).Infof("deleted inspection %s/%s", namespace, inspectionName)
		}

		// Update cache
		s.updateInspectionCache(cluster, feature, false)
	}

	return nil
}

// GetEtcdClusterInfo retrieves cluster and config information
func (s *Server) GetEtcdClusterInfo(namespace, name string) (*etcdv1alpha1.EtcdCluster, *etcd.ClientConfig, error) {
	// Try to get from cache first
	s.inspectionCacheMux.RLock()
	cacheKey := s.getCacheKey(namespace, name)
	entry, exists := s.inspectionCache[cacheKey]
	s.inspectionCacheMux.RUnlock()

	if exists && entry.cluster != nil && entry.clientConfig != nil {
		klog.V(4).Infof("using cached cluster info for %s/%s", namespace, name)
		return entry.cluster, entry.clientConfig, nil
	}

	// Fetch from API server
	cluster, err := s.etcdClient.EtcdV1alpha1().EtcdClusters(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get etcd cluster %s/%s: %v", namespace, name, err)
	}

	// Get client configuration
	clientConfig, err := s.getClientConfig(cluster)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get client config for cluster %s/%s: %v", namespace, name, err)
	}

	// Update cache
	s.inspectionCacheMux.Lock()
	if entry, exists := s.inspectionCache[cacheKey]; exists {
		entry.cluster = cluster
		entry.clientConfig = clientConfig
	} else {
		s.inspectionCache[cacheKey] = &inspectionCacheEntry{
			cluster:          cluster,
			clientConfig:     clientConfig,
			inspectionStatus: make(map[etcdv1alpha1.KStoneFeature]bool),
		}
	}
	s.inspectionCacheMux.Unlock()

	klog.V(3).Infof("loaded cluster info for %s/%s from API server", namespace, name)
	return cluster, clientConfig, nil
}

// getClientConfig creates etcd client configuration
func (s *Server) getClientConfig(cluster *etcdv1alpha1.EtcdCluster) (*etcd.ClientConfig, error) {
	// Build endpoints from cluster status
	endpoints := []string{cluster.Status.ServiceName}
	if len(cluster.Status.Members) > 0 {
		endpoints = make([]string, 0, len(cluster.Status.Members))
		for _, member := range cluster.Status.Members {
			if member.ExtensionClientUrl != "" {
				endpoints = append(endpoints, member.ExtensionClientUrl)
			} else if member.Endpoint != "" {
				endpoints = append(endpoints, member.Endpoint)
			}
		}
	}

	clientConfig := &etcd.ClientConfig{
		Endpoints: endpoints,
	}

	// Get TLS configuration if needed
	if cluster.Spec.SecureConfig != nil && cluster.Spec.SecureConfig.TLSSecret.Name != "" {
		secretName := cluster.Spec.SecureConfig.TLSSecret.Name
		secret, err := s.kubeClient.CoreV1().Secrets(cluster.Namespace).Get(
			context.TODO(),
			secretName,
			metav1.GetOptions{},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get TLS secret %s: %v", secretName, err)
		}

		secureConfig := &etcd.SecureConfig{}

		if certData, ok := secret.Data["tls.crt"]; ok {
			secureConfig.Cert = certData
		}
		if keyData, ok := secret.Data["tls.key"]; ok {
			secureConfig.Key = keyData
		}
		if caData, ok := secret.Data["ca.crt"]; ok {
			secureConfig.CA = caData
		}

		clientConfig.SecureConfig = secureConfig
	}

	return clientConfig, nil
}

// updateInspectionCache updates the inspection status cache
func (s *Server) updateInspectionCache(cluster *etcdv1alpha1.EtcdCluster, feature etcdv1alpha1.KStoneFeature, enabled bool) {
	s.inspectionCacheMux.Lock()
	defer s.inspectionCacheMux.Unlock()

	cacheKey := s.getCacheKey(cluster.Namespace, cluster.Name)
	entry, exists := s.inspectionCache[cacheKey]

	if !exists {
		entry = &inspectionCacheEntry{
			cluster:          cluster,
			inspectionStatus: make(map[etcdv1alpha1.KStoneFeature]bool),
		}
		s.inspectionCache[cacheKey] = entry
	}

	entry.inspectionStatus[feature] = enabled
	klog.V(3).Infof("updated cache for cluster %s/%s, feature %s, enabled=%v",
		cluster.Namespace, cluster.Name, feature, enabled)
}

// getCacheKey generates cache key for cluster
func (s *Server) getCacheKey(namespace, name string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
}

// InvalidateCache removes cluster from cache
func (s *Server) InvalidateCache(namespace, name string) {
	s.inspectionCacheMux.Lock()
	defer s.inspectionCacheMux.Unlock()

	cacheKey := s.getCacheKey(namespace, name)
	delete(s.inspectionCache, cacheKey)
	klog.V(2).Infof("invalidated cache for cluster %s/%s", namespace, name)
}

// Cleanup performs cleanup operations
func (s *Server) Cleanup() {
	klog.Info("cleaning up inspection server")
	// Stop all watchers
	s.StopAllWatchers()
	// Clear cache
	s.inspectionCacheMux.Lock()
	s.inspectionCache = make(map[string]*inspectionCacheEntry)
	s.inspectionCacheMux.Unlock()
}
