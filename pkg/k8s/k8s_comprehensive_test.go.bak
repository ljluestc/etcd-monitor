package k8s

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	corev1 "k8s.io/api/core/v1"
)

func TestK8sClient_Comprehensive(t *testing.T) {
	t.Run("NewK8sClient with valid config", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.NotNil(t, client.config)
		assert.NotNil(t, client.logger)
	})

	t.Run("NewK8sClient with nil config", func(t *testing.T) {
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(nil, logger)
		assert.NotNil(t, client)
		assert.Nil(t, client.config)
		assert.NotNil(t, client.logger)
	})

	t.Run("NewK8sClient with nil logger", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}

		client := NewK8sClient(config, nil)
		assert.NotNil(t, client)
		assert.NotNil(t, client.config)
		assert.NotNil(t, client.logger) // Should create a production logger
	})

	t.Run("NewK8sClient with empty config", func(t *testing.T) {
		config := &Config{}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.NotNil(t, client.config)
		assert.NotNil(t, client.logger)
	})

	t.Run("Start client", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		err := client.Start()
		assert.NoError(t, err)
	})

	t.Run("Stop client", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		err := client.Start()
		require.NoError(t, err)

		client.Stop()
		// Should not panic or error
	})

	t.Run("Stop client without starting", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		client.Stop()
		// Should not panic or error
	})

	t.Run("GetPods", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		pods := client.GetPods()
		assert.NotNil(t, pods)
	})

	t.Run("GetPods with namespace", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		pods := client.GetPodsWithNamespace("test-namespace")
		assert.NotNil(t, pods)
	})

	t.Run("GetPods with empty namespace", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		pods := client.GetPodsWithNamespace("")
		assert.NotNil(t, pods)
	})

	t.Run("GetPods with nil namespace", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		pods := client.GetPodsWithNamespace("")
		assert.NotNil(t, pods)
	})

	t.Run("GetPod", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		pod := client.GetPod("test-pod")
		assert.NotNil(t, pod)
	})

	t.Run("GetPod with empty name", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		pod := client.GetPod("")
		assert.NotNil(t, pod)
	})

	t.Run("GetPod with nil name", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		pod := client.GetPod("")
		assert.NotNil(t, pod)
	})

	t.Run("GetPodInfo", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		info := client.GetPodInfo("test-pod")
		assert.NotNil(t, info)
	})

	t.Run("GetPodInfo with empty name", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		info := client.GetPodInfo("")
		assert.NotNil(t, info)
	})

	t.Run("GetPodMetrics", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		metrics := client.GetPodMetrics("test-pod")
		assert.NotNil(t, metrics)
	})

	t.Run("GetPodMetrics with empty name", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		metrics := client.GetPodMetrics("")
		assert.NotNil(t, metrics)
	})

	t.Run("GetServices", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		services := client.GetServices()
		assert.NotNil(t, services)
	})

	t.Run("GetServices with namespace", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		services := client.GetServicesWithNamespace("test-namespace")
		assert.NotNil(t, services)
	})

	t.Run("GetServices with empty namespace", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		services := client.GetServicesWithNamespace("")
		assert.NotNil(t, services)
	})

	t.Run("GetService", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		service := client.GetService("test-service")
		assert.NotNil(t, service)
	})

	t.Run("GetService with empty name", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		service := client.GetService("")
		assert.NotNil(t, service)
	})

	t.Run("GetServiceInfo", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		info := client.GetServiceInfo("test-service")
		assert.NotNil(t, info)
	})

	t.Run("GetServiceInfo with empty name", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		info := client.GetServiceInfo("")
		assert.NotNil(t, info)
	})

	t.Run("GetServiceMetrics", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		metrics := client.GetServiceMetrics("test-service")
		assert.NotNil(t, metrics)
	})

	t.Run("GetServiceMetrics with empty name", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		metrics := client.GetServiceMetrics("")
		assert.NotNil(t, metrics)
	})

	t.Run("GetNodes", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		nodes := client.GetNodes()
		assert.NotNil(t, nodes)
	})

	t.Run("GetNode", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		node := client.GetNode("test-node")
		assert.NotNil(t, node)
	})

	t.Run("GetNode with empty name", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		node := client.GetNode("")
		assert.NotNil(t, node)
	})

	t.Run("GetNodeInfo", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		info := client.GetNodeInfo("test-node")
		assert.NotNil(t, info)
	})

	t.Run("GetNodeInfo with empty name", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		info := client.GetNodeInfo("")
		assert.NotNil(t, info)
	})

	t.Run("GetNodeMetrics", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		metrics := client.GetNodeMetrics("test-node")
		assert.NotNil(t, metrics)
	})

	t.Run("GetNodeMetrics with empty name", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		metrics := client.GetNodeMetrics("")
		assert.NotNil(t, metrics)
	})

	t.Run("GetSecrets", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		secrets := client.GetSecrets()
		assert.NotNil(t, secrets)
	})

	t.Run("GetSecrets with namespace", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		secrets := client.GetSecretsWithNamespace("test-namespace")
		assert.NotNil(t, secrets)
	})

	t.Run("GetSecrets with empty namespace", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		secrets := client.GetSecretsWithNamespace("")
		assert.NotNil(t, secrets)
	})

	t.Run("GetSecret", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		secret := client.GetSecret("test-secret")
		assert.NotNil(t, secret)
	})

	t.Run("GetSecret with empty name", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		secret := client.GetSecret("")
		assert.NotNil(t, secret)
	})

	t.Run("GetSecretInfo", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		info := client.GetSecretInfo("test-secret")
		assert.NotNil(t, info)
	})

	t.Run("GetSecretInfo with empty name", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		info := client.GetSecretInfo("")
		assert.NotNil(t, info)
	})

	t.Run("GetSecretMetrics", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		metrics := client.GetSecretMetrics("test-secret")
		assert.NotNil(t, metrics)
	})

	t.Run("GetSecretMetrics with empty name", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		metrics := client.GetSecretMetrics("")
		assert.NotNil(t, metrics)
	})
}

func TestK8sClient_EdgeCases(t *testing.T) {
	t.Run("K8sClient with very long kubeconfig path", func(t *testing.T) {
		longPath := strings.Repeat("a", 10000)
		config := &Config{
			KubeconfigPath: longPath,
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.Len(t, client.config.KubeconfigPath, 10000)
	})

	t.Run("K8sClient with very long context", func(t *testing.T) {
		longContext := strings.Repeat("b", 10000)
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        longContext,
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.Len(t, client.config.Context, 10000)
	})

	t.Run("K8sClient with very long namespace", func(t *testing.T) {
		longNamespace := strings.Repeat("c", 10000)
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      longNamespace,
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.Len(t, client.config.Namespace, 10000)
	})

	t.Run("K8sClient with empty kubeconfig path", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.Empty(t, client.config.KubeconfigPath)
	})

	t.Run("K8sClient with empty context", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.Empty(t, client.config.Context)
	})

	t.Run("K8sClient with empty namespace", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.Empty(t, client.config.Namespace)
	})

	t.Run("K8sClient with special characters in kubeconfig path", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig@domain.com",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.Equal(t, "/path/to/kubeconfig@domain.com", client.config.KubeconfigPath)
	})

	t.Run("K8sClient with special characters in context", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context@domain.com",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.Equal(t, "test-context@domain.com", client.config.Context)
	})

	t.Run("K8sClient with special characters in namespace", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default-namespace@domain.com",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.Equal(t, "default-namespace@domain.com", client.config.Namespace)
	})

	t.Run("K8sClient with unicode characters in kubeconfig path", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig/测试",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.Equal(t, "/path/to/kubeconfig/测试", client.config.KubeconfigPath)
	})

	t.Run("K8sClient with unicode characters in context", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "测试上下文",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.Equal(t, "测试上下文", client.config.Context)
	})

	t.Run("K8sClient with unicode characters in namespace", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "默认命名空间",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.Equal(t, "默认命名空间", client.config.Namespace)
	})

	t.Run("K8sClient with control characters in kubeconfig path", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig\n\r\t",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.Equal(t, "/path/to/kubeconfig\n\r\t", client.config.KubeconfigPath)
	})

	t.Run("K8sClient with control characters in context", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context\n\r\t",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.Equal(t, "test-context\n\r\t", client.config.Context)
	})

	t.Run("K8sClient with control characters in namespace", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default-namespace\n\r\t",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.Equal(t, "default-namespace\n\r\t", client.config.Namespace)
	})

	t.Run("K8sClient with very long resource names", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		longName := strings.Repeat("a", 10000)
		
		pod := client.GetPod(longName)
		assert.NotNil(t, pod)
		
		service := client.GetService(longName)
		assert.NotNil(t, service)
		
		node := client.GetNode(longName)
		assert.NotNil(t, node)
		
		secret := client.GetSecret(longName)
		assert.NotNil(t, secret)
	})

	t.Run("K8sClient with special characters in resource names", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		specialName := "test-resource@domain.com"
		
		pod := client.GetPod(specialName)
		assert.NotNil(t, pod)
		
		service := client.GetService(specialName)
		assert.NotNil(t, service)
		
		node := client.GetNode(specialName)
		assert.NotNil(t, node)
		
		secret := client.GetSecret(specialName)
		assert.NotNil(t, secret)
	})

	t.Run("K8sClient with unicode characters in resource names", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		unicodeName := "测试资源"
		
		pod := client.GetPod(unicodeName)
		assert.NotNil(t, pod)
		
		service := client.GetService(unicodeName)
		assert.NotNil(t, service)
		
		node := client.GetNode(unicodeName)
		assert.NotNil(t, node)
		
		secret := client.GetSecret(unicodeName)
		assert.NotNil(t, secret)
	})

	t.Run("K8sClient with control characters in resource names", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		controlName := "test-resource\n\r\t"
		
		pod := client.GetPod(controlName)
		assert.NotNil(t, pod)
		
		service := client.GetService(controlName)
		assert.NotNil(t, service)
		
		node := client.GetNode(controlName)
		assert.NotNil(t, node)
		
		secret := client.GetSecret(controlName)
		assert.NotNil(t, secret)
	})
}

func TestK8sClient_Performance(t *testing.T) {
	t.Run("K8sClient creation performance", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Create 1000 K8sClient instances
		for i := 0; i < 1000; i++ {
			client := NewK8sClient(config, logger)
			assert.NotNil(t, client)
		}
		
		duration := time.Since(start)
		
		// Should create 1000 instances in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Creating 1000 K8sClient instances took too long: %v", duration)
	})

	t.Run("K8sClient with large config performance", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: strings.Repeat("a", 10000),
			Context:        strings.Repeat("b", 10000),
			Namespace:      strings.Repeat("c", 10000),
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		client := NewK8sClient(config, logger)
		
		duration := time.Since(start)
		
		// Should create client with large config in less than 10ms
		assert.True(t, duration < 10*time.Millisecond, "Creating K8sClient with large config took too long: %v", duration)
		assert.Len(t, client.config.KubeconfigPath, 10000)
		assert.Len(t, client.config.Context, 10000)
		assert.Len(t, client.config.Namespace, 10000)
	})

	t.Run("K8sClient start/stop performance", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()
		
		start := time.Now()
		
		// Start and stop 100 clients
		for i := 0; i < 100; i++ {
			client := NewK8sClient(config, logger)
			err := client.Start()
			require.NoError(t, err)
			client.Stop()
		}
		
		duration := time.Since(start)
		
		// Should start and stop 100 clients in less than 1 second
		assert.True(t, duration < 1*time.Second, "Starting and stopping 100 clients took too long: %v", duration)
	})

	t.Run("K8sClient get pods performance", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()
		
		client := NewK8sClient(config, logger)
		
		start := time.Now()
		
		// Get pods 1000 times
		for i := 0; i < 1000; i++ {
			pods := client.GetPods()
			assert.NotNil(t, pods)
		}
		
		duration := time.Since(start)
		
		// Should get pods 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting pods 1000 times took too long: %v", duration)
	})

	t.Run("K8sClient get services performance", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()
		
		client := NewK8sClient(config, logger)
		
		start := time.Now()
		
		// Get services 1000 times
		for i := 0; i < 1000; i++ {
			services := client.GetServices()
			assert.NotNil(t, services)
		}
		
		duration := time.Since(start)
		
		// Should get services 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting services 1000 times took too long: %v", duration)
	})

	t.Run("K8sClient get nodes performance", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()
		
		client := NewK8sClient(config, logger)
		
		start := time.Now()
		
		// Get nodes 1000 times
		for i := 0; i < 1000; i++ {
			nodes := client.GetNodes()
			assert.NotNil(t, nodes)
		}
		
		duration := time.Since(start)
		
		// Should get nodes 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting nodes 1000 times took too long: %v", duration)
	})

	t.Run("K8sClient get secrets performance", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()
		
		client := NewK8sClient(config, logger)
		
		start := time.Now()
		
		// Get secrets 1000 times
		for i := 0; i < 1000; i++ {
			secrets := client.GetSecrets()
			assert.NotNil(t, secrets)
		}
		
		duration := time.Since(start)
		
		// Should get secrets 1000 times in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Getting secrets 1000 times took too long: %v", duration)
	})
}

func TestK8sClient_Integration(t *testing.T) {
	t.Run("K8sClient with all fields set", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.NotNil(t, client.config)
		assert.NotNil(t, client.logger)
		assert.Equal(t, "/path/to/kubeconfig", client.config.KubeconfigPath)
		assert.Equal(t, "test-context", client.config.Context)
		assert.Equal(t, "default", client.config.Namespace)
	})

	t.Run("K8sClient with minimal fields", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.NotNil(t, client.config)
		assert.NotNil(t, client.logger)
		assert.Equal(t, "/path/to/kubeconfig", client.config.KubeconfigPath)
		assert.Equal(t, "test-context", client.config.Context)
		assert.Equal(t, "default", client.config.Namespace)
	})

	t.Run("K8sClient with empty fields", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "",
			Context:        "",
			Namespace:      "",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.NotNil(t, client.config)
		assert.NotNil(t, client.logger)
		assert.Empty(t, client.config.KubeconfigPath)
		assert.Empty(t, client.config.Context)
		assert.Empty(t, client.config.Namespace)
	})

	t.Run("K8sClient with nil fields", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "",
			Context:        "",
			Namespace:      "",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)
		assert.NotNil(t, client.config)
		assert.NotNil(t, client.logger)
		assert.Empty(t, client.config.KubeconfigPath)
		assert.Empty(t, client.config.Context)
		assert.Empty(t, client.config.Namespace)
	})

	t.Run("K8sClient full lifecycle", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)

		// Start client
		err := client.Start()
		assert.NoError(t, err)

		// Get pods
		pods := client.GetPods()
		assert.NotNil(t, pods)

		// Get specific pod
		pod := client.GetPod("test-pod")
		assert.NotNil(t, pod)

		// Get pod info
		podInfo := client.GetPodInfo("test-pod")
		assert.NotNil(t, podInfo)

		// Get pod metrics
		podMetrics := client.GetPodMetrics("test-pod")
		assert.NotNil(t, podMetrics)

		// Get services
		services := client.GetServices()
		assert.NotNil(t, services)

		// Get specific service
		service := client.GetService("test-service")
		assert.NotNil(t, service)

		// Get service info
		serviceInfo := client.GetServiceInfo("test-service")
		assert.NotNil(t, serviceInfo)

		// Get service metrics
		serviceMetrics := client.GetServiceMetrics("test-service")
		assert.NotNil(t, serviceMetrics)

		// Get nodes
		nodes := client.GetNodes()
		assert.NotNil(t, nodes)

		// Get specific node
		node := client.GetNode("test-node")
		assert.NotNil(t, node)

		// Get node info
		nodeInfo := client.GetNodeInfo("test-node")
		assert.NotNil(t, nodeInfo)

		// Get node metrics
		nodeMetrics := client.GetNodeMetrics("test-node")
		assert.NotNil(t, nodeMetrics)

		// Get secrets
		secrets := client.GetSecrets()
		assert.NotNil(t, secrets)

		// Get specific secret
		secret := client.GetSecret("test-secret")
		assert.NotNil(t, secret)

		// Get secret info
		secretInfo := client.GetSecretInfo("test-secret")
		assert.NotNil(t, secretInfo)

		// Get secret metrics
		secretMetrics := client.GetSecretMetrics("test-secret")
		assert.NotNil(t, secretMetrics)

		// Stop client
		client.Stop()
	})

	t.Run("K8sClient with multiple concurrent operations", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)

		// Start client
		err := client.Start()
		require.NoError(t, err)

		// Perform multiple concurrent operations
		done := make(chan bool, 10)
		
		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				// Get pods
				pods := client.GetPods()
				assert.NotNil(t, pods)
				
				// Get specific pod
				pod := client.GetPod(fmt.Sprintf("pod-%d", id))
				assert.NotNil(t, pod)
				
				// Get pod info
				podInfo := client.GetPodInfo(fmt.Sprintf("pod-%d", id))
				assert.NotNil(t, podInfo)
				
				// Get pod metrics
				podMetrics := client.GetPodMetrics(fmt.Sprintf("pod-%d", id))
				assert.NotNil(t, podMetrics)
				
				// Get services
				services := client.GetServices()
				assert.NotNil(t, services)
				
				// Get specific service
				service := client.GetService(fmt.Sprintf("service-%d", id))
				assert.NotNil(t, service)
				
				// Get service info
				serviceInfo := client.GetServiceInfo(fmt.Sprintf("service-%d", id))
				assert.NotNil(t, serviceInfo)
				
				// Get service metrics
				serviceMetrics := client.GetServiceMetrics(fmt.Sprintf("service-%d", id))
				assert.NotNil(t, serviceMetrics)
				
				// Get nodes
				nodes := client.GetNodes()
				assert.NotNil(t, nodes)
				
				// Get specific node
				node := client.GetNode(fmt.Sprintf("node-%d", id))
				assert.NotNil(t, node)
				
				// Get node info
				nodeInfo := client.GetNodeInfo(fmt.Sprintf("node-%d", id))
				assert.NotNil(t, nodeInfo)
				
				// Get node metrics
				nodeMetrics := client.GetNodeMetrics(fmt.Sprintf("node-%d", id))
				assert.NotNil(t, nodeMetrics)
				
				// Get secrets
				secrets := client.GetSecrets()
				assert.NotNil(t, secrets)
				
				// Get specific secret
				secret := client.GetSecret(fmt.Sprintf("secret-%d", id))
				assert.NotNil(t, secret)
				
				// Get secret info
				secretInfo := client.GetSecretInfo(fmt.Sprintf("secret-%d", id))
				assert.NotNil(t, secretInfo)
				
				// Get secret metrics
				secretMetrics := client.GetSecretMetrics(fmt.Sprintf("secret-%d", id))
				assert.NotNil(t, secretMetrics)
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Stop client
		client.Stop()
	})

	t.Run("K8sClient with different resource names", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)

		// Start client
		err := client.Start()
		require.NoError(t, err)

		// Test with different resource names
		resourceNames := []string{
			"test-pod",
			"production-pod",
			"staging-pod",
			"development-pod",
			"pod-1",
			"pod-2",
			"pod-3",
		}

		for _, name := range resourceNames {
			// Get pod
			pod := client.GetPod(name)
			assert.NotNil(t, pod)
			
			// Get pod info
			podInfo := client.GetPodInfo(name)
			assert.NotNil(t, podInfo)
			
			// Get pod metrics
			podMetrics := client.GetPodMetrics(name)
			assert.NotNil(t, podMetrics)
			
			// Get service
			service := client.GetService(name)
			assert.NotNil(t, service)
			
			// Get service info
			serviceInfo := client.GetServiceInfo(name)
			assert.NotNil(t, serviceInfo)
			
			// Get service metrics
			serviceMetrics := client.GetServiceMetrics(name)
			assert.NotNil(t, serviceMetrics)
			
			// Get node
			node := client.GetNode(name)
			assert.NotNil(t, node)
			
			// Get node info
			nodeInfo := client.GetNodeInfo(name)
			assert.NotNil(t, nodeInfo)
			
			// Get node metrics
			nodeMetrics := client.GetNodeMetrics(name)
			assert.NotNil(t, nodeMetrics)
			
			// Get secret
			secret := client.GetSecret(name)
			assert.NotNil(t, secret)
			
			// Get secret info
			secretInfo := client.GetSecretInfo(name)
			assert.NotNil(t, secretInfo)
			
			// Get secret metrics
			secretMetrics := client.GetSecretMetrics(name)
			assert.NotNil(t, secretMetrics)
		}

		// Stop client
		client.Stop()
	})

	t.Run("K8sClient with empty resource names", func(t *testing.T) {
		config := &Config{
			KubeconfigPath: "/path/to/kubeconfig",
			Context:        "test-context",
			Namespace:      "default",
		}
		logger, _ := zap.NewDevelopment()

		client := NewK8sClient(config, logger)
		assert.NotNil(t, client)

		// Start client
		err := client.Start()
		require.NoError(t, err)

		// Test with empty resource name
		pod := client.GetPod("")
		assert.NotNil(t, pod)
		
		podInfo := client.GetPodInfo("")
		assert.NotNil(t, podInfo)
		
		podMetrics := client.GetPodMetrics("")
		assert.NotNil(t, podMetrics)
		
		service := client.GetService("")
		assert.NotNil(t, service)
		
		serviceInfo := client.GetServiceInfo("")
		assert.NotNil(t, serviceInfo)
		
		serviceMetrics := client.GetServiceMetrics("")
		assert.NotNil(t, serviceMetrics)
		
		node := client.GetNode("")
		assert.NotNil(t, node)
		
		nodeInfo := client.GetNodeInfo("")
		assert.NotNil(t, nodeInfo)
		
		nodeMetrics := client.GetNodeMetrics("")
		assert.NotNil(t, nodeMetrics)
		
		secret := client.GetSecret("")
		assert.NotNil(t, secret)
		
		secretInfo := client.GetSecretInfo("")
		assert.NotNil(t, secretInfo)
		
		secretMetrics := client.GetSecretMetrics("")
		assert.NotNil(t, secretMetrics)

		// Stop client
		client.Stop()
	})
}
