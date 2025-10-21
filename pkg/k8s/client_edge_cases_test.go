package k8s

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
)

// TestGetClientConfig_EdgeCases tests edge cases for GetClientConfig
func TestGetClientConfig_EdgeCases(t *testing.T) {
	t.Run("Empty kubeconfig triggers in-cluster config", func(t *testing.T) {
		// When kubeconfig is empty, should try in-cluster config
		_, err := GetClientConfig("")
		// Will error in test environment (no in-cluster config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to load in-cluster configuration")
	})

	t.Run("Non-empty kubeconfig with invalid path", func(t *testing.T) {
		// Should error when file doesn't exist
		_, err := GetClientConfig("/path/that/does/not/exist.yaml")
		assert.Error(t, err)
	})

	t.Run("Non-empty kubeconfig with invalid YAML", func(t *testing.T) {
		// Create temp file with invalid YAML
		tmpDir := t.TempDir()
		kubeconfigPath := filepath.Join(tmpDir, "invalid-kubeconfig.yaml")
		err := os.WriteFile(kubeconfigPath, []byte("not: valid: yaml: content"), 0644)
		require.NoError(t, err)

		_, err = GetClientConfig(kubeconfigPath)
		assert.Error(t, err)
	})

	t.Run("Non-empty kubeconfig with empty file", func(t *testing.T) {
		// Create temp empty file
		tmpDir := t.TempDir()
		kubeconfigPath := filepath.Join(tmpDir, "empty-kubeconfig.yaml")
		err := os.WriteFile(kubeconfigPath, []byte(""), 0644)
		require.NoError(t, err)

		_, err = GetClientConfig(kubeconfigPath)
		assert.Error(t, err)
	})

	t.Run("Non-empty kubeconfig with valid structure but no clusters", func(t *testing.T) {
		// Create temp file with minimal valid YAML but no clusters
		tmpDir := t.TempDir()
		kubeconfigPath := filepath.Join(tmpDir, "minimal-kubeconfig.yaml")
		kubeconfig := `apiVersion: v1
kind: Config
clusters: []
contexts: []
current-context: ""
preferences: {}
users: []
`
		err := os.WriteFile(kubeconfigPath, []byte(kubeconfig), 0644)
		require.NoError(t, err)

		_, err = GetClientConfig(kubeconfigPath)
		// Should error because no current context
		assert.Error(t, err)
	})

	t.Run("Non-empty kubeconfig with valid minimal config", func(t *testing.T) {
		// Create temp file with minimal valid YAML
		tmpDir := t.TempDir()
		kubeconfigPath := filepath.Join(tmpDir, "valid-kubeconfig.yaml")
		kubeconfig := `apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://localhost:6443
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test-context
current-context: test-context
preferences: {}
users:
- name: test-user
  user:
    token: test-token
`
		err := os.WriteFile(kubeconfigPath, []byte(kubeconfig), 0644)
		require.NoError(t, err)

		config, err := GetClientConfig(kubeconfigPath)
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, "https://localhost:6443", config.Host)
		assert.Equal(t, "test-token", config.BearerToken)
	})
}

// TestGenerateInformer_EdgeCases tests edge cases for GenerateInformer
func TestGenerateInformer_EdgeCases(t *testing.T) {
	t.Run("Nil config causes panic", func(t *testing.T) {
		// GenerateInformer with nil config should panic
		defer func() {
			if r := recover(); r != nil {
				// Expected panic
				assert.NotNil(t, r)
			}
		}()
		_, _, _, _, err := GenerateInformer(nil, "")
		// If we get here without panic, check for error
		if err == nil {
			t.Error("Expected error or panic with nil config")
		}
	})

	t.Run("Valid config with empty host", func(t *testing.T) {
		config := &rest.Config{
			Host: "",
		}
		_, _, _, _, err := GenerateInformer(config, "")
		// May or may not error depending on defaults
		_ = err // Just test it doesn't panic
	})

	t.Run("Valid config with localhost", func(t *testing.T) {
		config := &rest.Config{
			Host: "https://localhost:6443",
		}
		kubeClient, clustetClient, kubeInformer, etcdInformer, err := GenerateInformer(config, "")
		assert.NoError(t, err)
		assert.NotNil(t, kubeClient)
		assert.NotNil(t, clustetClient)
		assert.NotNil(t, kubeInformer)
		assert.NotNil(t, etcdInformer)
	})

	t.Run("Valid config with TLS", func(t *testing.T) {
		config := &rest.Config{
			Host: "https://localhost:6443",
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true,
			},
		}
		kubeClient, clustetClient, kubeInformer, etcdInformer, err := GenerateInformer(config, "")
		assert.NoError(t, err)
		assert.NotNil(t, kubeClient)
		assert.NotNil(t, clustetClient)
		assert.NotNil(t, kubeInformer)
		assert.NotNil(t, etcdInformer)
	})

	t.Run("Valid config with bearer token", func(t *testing.T) {
		config := &rest.Config{
			Host:        "https://localhost:6443",
			BearerToken: "test-token",
		}
		kubeClient, clustetClient, kubeInformer, etcdInformer, err := GenerateInformer(config, "")
		assert.NoError(t, err)
		assert.NotNil(t, kubeClient)
		assert.NotNil(t, clustetClient)
		assert.NotNil(t, kubeInformer)
		assert.NotNil(t, etcdInformer)
	})

	t.Run("Valid config with multiple label selectors", func(t *testing.T) {
		config := &rest.Config{
			Host: "https://localhost:6443",
		}
		labelSelectors := []string{
			"",
			"app=test",
			"app=test,version=v1",
			"app=test,version=v1,env=prod",
			"app!=test",
			"app in (test,prod)",
			"app notin (test,prod)",
		}

		for _, selector := range labelSelectors {
			kubeClient, clustetClient, kubeInformer, etcdInformer, err := GenerateInformer(config, selector)
			assert.NoError(t, err, "Failed with selector: %s", selector)
			assert.NotNil(t, kubeClient)
			assert.NotNil(t, clustetClient)
			assert.NotNil(t, kubeInformer)
			assert.NotNil(t, etcdInformer)
		}
	})

	t.Run("Config with timeout", func(t *testing.T) {
		config := &rest.Config{
			Host:    "https://localhost:6443",
			Timeout: 30,
		}
		kubeClient, clustetClient, kubeInformer, etcdInformer, err := GenerateInformer(config, "")
		assert.NoError(t, err)
		assert.NotNil(t, kubeClient)
		assert.NotNil(t, clustetClient)
		assert.NotNil(t, kubeInformer)
		assert.NotNil(t, etcdInformer)
	})

	t.Run("Config with QPS and Burst", func(t *testing.T) {
		config := &rest.Config{
			Host:  "https://localhost:6443",
			QPS:   100,
			Burst: 200,
		}
		kubeClient, clustetClient, kubeInformer, etcdInformer, err := GenerateInformer(config, "")
		assert.NoError(t, err)
		assert.NotNil(t, kubeClient)
		assert.NotNil(t, clustetClient)
		assert.NotNil(t, kubeInformer)
		assert.NotNil(t, etcdInformer)
	})
}

// TestGetClientConfig_BothBranches ensures both branches of GetClientConfig are tested
func TestGetClientConfig_BothBranches(t *testing.T) {
	t.Run("Empty kubeconfig branch", func(t *testing.T) {
		// Empty string should trigger in-cluster config path
		_, err := GetClientConfig("")
		assert.Error(t, err) // Will fail in test environment
		assert.Contains(t, err.Error(), "unable to load in-cluster configuration")
	})

	t.Run("Non-empty kubeconfig branch", func(t *testing.T) {
		// Create a valid kubeconfig file
		tmpDir := t.TempDir()
		kubeconfigPath := filepath.Join(tmpDir, "valid-kubeconfig.yaml")
		kubeconfig := `apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://kubernetes.default.svc
  name: default-cluster
contexts:
- context:
    cluster: default-cluster
    user: default-user
  name: default-context
current-context: default-context
preferences: {}
users:
- name: default-user
  user:
    token: default-token
`
		err := os.WriteFile(kubeconfigPath, []byte(kubeconfig), 0644)
		require.NoError(t, err)

		// Non-empty string should trigger kubeconfig file path
		config, err := GetClientConfig(kubeconfigPath)
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.NotNil(t, config.Host)
		assert.Equal(t, "https://kubernetes.default.svc", config.Host)
	})
}

// TestGenerateInformer_AllPaths tests all execution paths
func TestGenerateInformer_AllPaths(t *testing.T) {
	t.Run("Success path with all components", func(t *testing.T) {
		config := &rest.Config{
			Host: "https://test-cluster:6443",
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true,
			},
		}

		kubeClient, clustetClient, kubeInformer, etcdInformer, err := GenerateInformer(config, "app=test")

		// Verify all return values
		assert.NoError(t, err)
		assert.NotNil(t, kubeClient, "kubeClient should not be nil")
		assert.NotNil(t, clustetClient, "clustetClient should not be nil")
		assert.NotNil(t, kubeInformer, "kubeInformer should not be nil")
		assert.NotNil(t, etcdInformer, "etcdInformer should not be nil")
	})

	t.Run("Verify informer factory resync period", func(t *testing.T) {
		config := &rest.Config{
			Host: "https://test-cluster:6443",
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true,
			},
		}

		_, _, kubeInformer, etcdInformer, err := GenerateInformer(config, "")

		assert.NoError(t, err)
		assert.NotNil(t, kubeInformer)
		assert.NotNil(t, etcdInformer)

		// The function uses time.Second*30 for resync period
		// We can't directly test the resync period, but we verified the factories are created
	})
}
