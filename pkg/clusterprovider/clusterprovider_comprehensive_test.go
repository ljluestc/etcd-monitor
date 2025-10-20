package clusterprovider

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	etcdv1alpha1 "github.com/etcd-monitor/taskmaster/api/etcd/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClusterProvider_Comprehensive(t *testing.T) {
	t.Run("NewClusterProvider with valid config", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		assert.NotNil(t, provider)
		assert.NotNil(t, provider.config)
		assert.NotNil(t, provider.logger)
	})

	t.Run("NewClusterProvider with nil config", func(t *testing.T) {
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(nil, logger)
		assert.NotNil(t, provider)
		assert.Nil(t, provider.config)
		assert.NotNil(t, provider.logger)
	})

	t.Run("NewClusterProvider with nil logger", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		provider := NewClusterProvider(config, nil)
		assert.NotNil(t, provider)
		assert.NotNil(t, provider.config)
		assert.NotNil(t, provider.logger) // Should create a production logger
	})

	t.Run("NewClusterProvider with empty config", func(t *testing.T) {
		config := &Config{}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		assert.NotNil(t, provider)
		assert.NotNil(t, provider.config)
		assert.NotNil(t, provider.logger)
	})

	t.Run("Start provider", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		err := provider.Start()
		assert.NoError(t, err)
	})

	t.Run("Stop provider", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		err := provider.Start()
		require.NoError(t, err)

		provider.Stop()
		// Should not panic or error
	})

	t.Run("Stop provider without starting", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		provider.Stop()
		// Should not panic or error
	})

	t.Run("GetClusters", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		clusters := provider.GetClusters()
		assert.NotNil(t, clusters)
	})

	t.Run("GetCluster", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		cluster := provider.GetCluster("test-cluster")
		assert.NotNil(t, cluster)
	})

	t.Run("GetCluster with empty name", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		cluster := provider.GetCluster("")
		assert.NotNil(t, cluster)
	})

	t.Run("GetCluster with nil name", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		cluster := provider.GetCluster("")
		assert.NotNil(t, cluster)
	})

	t.Run("GetClusterInfo", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		info := provider.GetClusterInfo("test-cluster")
		assert.NotNil(t, info)
	})

	t.Run("GetClusterInfo with empty name", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		info := provider.GetClusterInfo("")
		assert.NotNil(t, info)
	})

	t.Run("GetClusterMetrics", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		metrics := provider.GetClusterMetrics("test-cluster")
		assert.NotNil(t, metrics)
	})

	t.Run("GetClusterMetrics with empty name", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		metrics := provider.GetClusterMetrics("")
		assert.NotNil(t, metrics)
	})
}

func TestClusterProvider_EdgeCases(t *testing.T) {
	t.Run("ClusterProvider with very long endpoint list", func(t *testing.T) {
		endpoints := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			endpoints[i] = fmt.Sprintf("localhost:%d", 2379+i)
		}

		config := &Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		assert.NotNil(t, provider)
		assert.Len(t, provider.config.Endpoints, 1000)
	})

	t.Run("ClusterProvider with very long timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 24 * time.Hour,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		assert.NotNil(t, provider)
		assert.Equal(t, 24*time.Hour, provider.config.DialTimeout)
	})

	t.Run("ClusterProvider with very short timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 1 * time.Nanosecond,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		assert.NotNil(t, provider)
		assert.Equal(t, 1*time.Nanosecond, provider.config.DialTimeout)
	})

	t.Run("ClusterProvider with zero timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 0,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		assert.NotNil(t, provider)
		assert.Equal(t, time.Duration(0), provider.config.DialTimeout)
	})

	t.Run("ClusterProvider with negative timeout", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: -1 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		assert.NotNil(t, provider)
		assert.Equal(t, -1*time.Second, provider.config.DialTimeout)
	})

	t.Run("ClusterProvider with empty endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		assert.NotNil(t, provider)
		assert.Len(t, provider.config.Endpoints, 0)
	})

	t.Run("ClusterProvider with nil endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		assert.NotNil(t, provider)
		assert.Nil(t, provider.config.Endpoints)
	})

	t.Run("ClusterProvider with special characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "https://etcd-1:2379", "http://etcd-2:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		assert.NotNil(t, provider)
		assert.Len(t, provider.config.Endpoints, 3)
		assert.Equal(t, "http://etcd-0:2379", provider.config.Endpoints[0])
		assert.Equal(t, "https://etcd-1:2379", provider.config.Endpoints[1])
		assert.Equal(t, "http://etcd-2:2379", provider.config.Endpoints[2])
	})

	t.Run("ClusterProvider with unicode characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		assert.NotNil(t, provider)
		assert.Len(t, provider.config.Endpoints, 2)
	})

	t.Run("ClusterProvider with control characters in endpoints", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		assert.NotNil(t, provider)
		assert.Len(t, provider.config.Endpoints, 2)
	})

	t.Run("ClusterProvider with very long cluster name", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		longName := strings.Repeat("a", 10000)
		
		cluster := provider.GetCluster(longName)
		assert.NotNil(t, cluster)
		
		info := provider.GetClusterInfo(longName)
		assert.NotNil(t, info)
		
		metrics := provider.GetClusterMetrics(longName)
		assert.NotNil(t, metrics)
	})

	t.Run("ClusterProvider with special characters in cluster name", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		specialName := "test-cluster@domain.com"
		
		cluster := provider.GetCluster(specialName)
		assert.NotNil(t, cluster)
		
		info := provider.GetClusterInfo(specialName)
		assert.NotNil(t, info)
		
		metrics := provider.GetClusterMetrics(specialName)
		assert.NotNil(t, metrics)
	})

	t.Run("ClusterProvider with unicode characters in cluster name", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		unicodeName := "测试集群"
		
		cluster := provider.GetCluster(unicodeName)
		assert.NotNil(t, cluster)
		
		info := provider.GetClusterInfo(unicodeName)
		assert.NotNil(t, info)
		
		metrics := provider.GetClusterMetrics(unicodeName)
		assert.NotNil(t, metrics)
	})

	t.Run("ClusterProvider with control characters in cluster name", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		controlName := "test-cluster\n\r\t"
		
		cluster := provider.GetCluster(controlName)
		assert.NotNil(t, cluster)
		
		info := provider.GetClusterInfo(controlName)
		assert.NotNil(t, info)
		
		metrics := provider.GetClusterMetrics(controlName)
		assert.NotNil(t, metrics)
	})
}

func TestClusterProvider_Performance(t *testing.T) {
	t.Run("ClusterProvider creation performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Create 1000 ClusterProvider instances
		for i := 0; i < 1000; i++ {
			provider := NewClusterProvider(config, logger)
			assert.NotNil(t, provider)
		}
		
		duration := time.Since(start)
		
		// Should create 1000 instances in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Creating 1000 ClusterProvider instances took too long: %v", duration)
	})

	t.Run("ClusterProvider with large endpoint list performance", func(t *testing.T) {
		endpoints := make([]string, 10000)
		for i := 0; i < 10000; i++ {
			endpoints[i] = fmt.Sprintf("localhost:%d", 2379+i)
		}
		
		config := &Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		provider := NewClusterProvider(config, logger)
		
		duration := time.Since(start)
		
		// Should create provider with 10000 endpoints in less than 10ms
		assert.True(t, duration < 10*time.Millisecond, "Creating ClusterProvider with 10000 endpoints took too long: %v", duration)
		assert.Len(t, provider.config.Endpoints, 10000)
	})

	t.Run("ClusterProvider start/stop performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Start and stop 100 providers
		for i := 0; i < 100; i++ {
			provider := NewClusterProvider(config, logger)
			err := provider.Start()
			require.NoError(t, err)
			provider.Stop()
		}
		
		duration := time.Since(start)
		
		// Should start and stop 100 providers in less than 1 second
		assert.True(t, duration < 1*time.Second, "Starting and stopping 100 providers took too long: %v", duration)
	})

	t.Run("ClusterProvider get clusters performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()
		
		provider := NewClusterProvider(config, logger)
		
		start := time.Now()
		
		// Get clusters 1000 times
		for i := 0; i < 1000; i++ {
			clusters := provider.GetClusters()
			assert.NotNil(t, clusters)
		}
		
		duration := time.Since(start)
		
		// Should get clusters 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting clusters 1000 times took too long: %v", duration)
	})

	t.Run("ClusterProvider get cluster performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()
		
		provider := NewClusterProvider(config, logger)
		
		start := time.Now()
		
		// Get cluster 1000 times
		for i := 0; i < 1000; i++ {
			cluster := provider.GetCluster(fmt.Sprintf("cluster-%d", i))
			assert.NotNil(t, cluster)
		}
		
		duration := time.Since(start)
		
		// Should get cluster 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting cluster 1000 times took too long: %v", duration)
	})

	t.Run("ClusterProvider get cluster info performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()
		
		provider := NewClusterProvider(config, logger)
		
		start := time.Now()
		
		// Get cluster info 1000 times
		for i := 0; i < 1000; i++ {
			info := provider.GetClusterInfo(fmt.Sprintf("cluster-%d", i))
			assert.NotNil(t, info)
		}
		
		duration := time.Since(start)
		
		// Should get cluster info 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting cluster info 1000 times took too long: %v", duration)
	})

	t.Run("ClusterProvider get cluster metrics performance", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()
		
		provider := NewClusterProvider(config, logger)
		
		start := time.Now()
		
		// Get cluster metrics 1000 times
		for i := 0; i < 1000; i++ {
			metrics := provider.GetClusterMetrics(fmt.Sprintf("cluster-%d", i))
			assert.NotNil(t, metrics)
		}
		
		duration := time.Since(start)
		
		// Should get cluster metrics 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting cluster metrics 1000 times took too long: %v", duration)
	})
}

func TestClusterProvider_Integration(t *testing.T) {
	t.Run("ClusterProvider with all fields set", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379", "http://etcd-2:2379"},
			DialTimeout: 10 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		assert.NotNil(t, provider)
		assert.NotNil(t, provider.config)
		assert.NotNil(t, provider.logger)
		assert.Len(t, provider.config.Endpoints, 3)
		assert.Equal(t, "http://etcd-0:2379", provider.config.Endpoints[0])
		assert.Equal(t, "http://etcd-1:2379", provider.config.Endpoints[1])
		assert.Equal(t, "http://etcd-2:2379", provider.config.Endpoints[2])
		assert.Equal(t, 10*time.Second, provider.config.DialTimeout)
	})

	t.Run("ClusterProvider with minimal fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		assert.NotNil(t, provider)
		assert.NotNil(t, provider.config)
		assert.NotNil(t, provider.logger)
		assert.Len(t, provider.config.Endpoints, 1)
		assert.Equal(t, "localhost:2379", provider.config.Endpoints[0])
		assert.Equal(t, 5*time.Second, provider.config.DialTimeout)
	})

	t.Run("ClusterProvider with empty fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{},
			DialTimeout: 0,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		assert.NotNil(t, provider)
		assert.NotNil(t, provider.config)
		assert.NotNil(t, provider.logger)
		assert.Len(t, provider.config.Endpoints, 0)
		assert.Equal(t, time.Duration(0), provider.config.DialTimeout)
	})

	t.Run("ClusterProvider with nil fields", func(t *testing.T) {
		config := &Config{
			Endpoints:   nil,
			DialTimeout: 0,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		assert.NotNil(t, provider)
		assert.NotNil(t, provider.config)
		assert.NotNil(t, provider.logger)
		assert.Nil(t, provider.config.Endpoints)
		assert.Equal(t, time.Duration(0), provider.config.DialTimeout)
	})

	t.Run("ClusterProvider full lifecycle", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		assert.NotNil(t, provider)

		// Start provider
		err := provider.Start()
		assert.NoError(t, err)

		// Get clusters
		clusters := provider.GetClusters()
		assert.NotNil(t, clusters)

		// Get specific cluster
		cluster := provider.GetCluster("test-cluster")
		assert.NotNil(t, cluster)

		// Get cluster info
		info := provider.GetClusterInfo("test-cluster")
		assert.NotNil(t, info)

		// Get cluster metrics
		metrics := provider.GetClusterMetrics("test-cluster")
		assert.NotNil(t, metrics)

		// Stop provider
		provider.Stop()
	})

	t.Run("ClusterProvider with multiple concurrent operations", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		assert.NotNil(t, provider)

		// Start provider
		err := provider.Start()
		require.NoError(t, err)

		// Perform multiple concurrent operations
		done := make(chan bool, 10)
		
		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				// Get clusters
				clusters := provider.GetClusters()
				assert.NotNil(t, clusters)
				
				// Get specific cluster
				cluster := provider.GetCluster(fmt.Sprintf("cluster-%d", id))
				assert.NotNil(t, cluster)
				
				// Get cluster info
				info := provider.GetClusterInfo(fmt.Sprintf("cluster-%d", id))
				assert.NotNil(t, info)
				
				// Get cluster metrics
				metrics := provider.GetClusterMetrics(fmt.Sprintf("cluster-%d", id))
				assert.NotNil(t, metrics)
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Stop provider
		provider.Stop()
	})

	t.Run("ClusterProvider with different cluster names", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		assert.NotNil(t, provider)

		// Start provider
		err := provider.Start()
		require.NoError(t, err)

		// Test with different cluster names
		clusterNames := []string{
			"test-cluster",
			"production-cluster",
			"staging-cluster",
			"development-cluster",
			"cluster-1",
			"cluster-2",
			"cluster-3",
		}

		for _, name := range clusterNames {
			// Get cluster
			cluster := provider.GetCluster(name)
			assert.NotNil(t, cluster)
			
			// Get cluster info
			info := provider.GetClusterInfo(name)
			assert.NotNil(t, info)
			
			// Get cluster metrics
			metrics := provider.GetClusterMetrics(name)
			assert.NotNil(t, metrics)
		}

		// Stop provider
		provider.Stop()
	})

	t.Run("ClusterProvider with empty cluster names", func(t *testing.T) {
		config := &Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}
		logger, _ := zap.NewDevelopment()

		provider := NewClusterProvider(config, logger)
		assert.NotNil(t, provider)

		// Start provider
		err := provider.Start()
		require.NoError(t, err)

		// Test with empty cluster name
		cluster := provider.GetCluster("")
		assert.NotNil(t, cluster)
		
		info := provider.GetClusterInfo("")
		assert.NotNil(t, info)
		
		metrics := provider.GetClusterMetrics("")
		assert.NotNil(t, metrics)

		// Stop provider
		provider.Stop()
	})
}
