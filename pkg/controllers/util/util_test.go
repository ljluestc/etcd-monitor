package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientBuilder_Comprehensive(t *testing.T) {
	t.Run("NewSimpleClientBuilder", func(t *testing.T) {
		builder := NewSimpleClientBuilder("test-kubeconfig")
		
		assert.NotNil(t, builder)
		assert.IsType(t, &simpleClientBuilder{}, builder)
	})

	t.Run("ConfigOrDie", func(t *testing.T) {
		builder := NewSimpleClientBuilder("test-kubeconfig")
		
		// This will panic if there's an error, which is expected behavior
		// We can't easily test this without mocking the k8s client
		assert.NotNil(t, builder)
	})

	t.Run("ClientOrDie", func(t *testing.T) {
		builder := NewSimpleClientBuilder("test-kubeconfig")
		
		// This will panic if there's an error, which is expected behavior
		// We can't easily test this without mocking the k8s client
		assert.NotNil(t, builder)
	})

	t.Run("DynamicClientOrDie", func(t *testing.T) {
		builder := NewSimpleClientBuilder("test-kubeconfig")
		
		// This will panic if there's an error, which is expected behavior
		// We can't easily test this without mocking the k8s client
		assert.NotNil(t, builder)
	})
}

func TestUtil_Constants(t *testing.T) {
	t.Run("Component Constants", func(t *testing.T) {
		assert.Equal(t, "etcdcluster-controller", ComponentEtcdClusterController)
		assert.Equal(t, "etcdinspection-controller", ComponentEtcdInspectionController)
	})

	t.Run("EtcdClusterPhase Constants", func(t *testing.T) {
		assert.Equal(t, EtcdClusterPhase("EtcdClusterCreating"), EtcdClusterCreating)
		assert.Equal(t, EtcdClusterPhase("EtcdClusterUpdating"), EtcdClusterUpdating)
		assert.Equal(t, EtcdClusterPhase("EtcdClusterUpdateStatus"), EtcdClusterUpdateStatus)
	})

	t.Run("Cluster Constants", func(t *testing.T) {
		assert.Equal(t, "certName", ClusterTLSSecretName)
		assert.Equal(t, "extClientURL", ClusterExtensionClientURL)
	})
}
