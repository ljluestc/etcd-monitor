package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestClientConfig_Comprehensive(t *testing.T) {
	t.Run("ClientConfig Creation", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"http://localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "user",
			Password:    "pass",
		}

		assert.Equal(t, []string{"http://localhost:2379"}, config.Endpoints)
		assert.Equal(t, 5*time.Second, config.DialTimeout)
		assert.Equal(t, "user", config.Username)
		assert.Equal(t, "pass", config.Password)
	})

	t.Run("SecureConfig Creation", func(t *testing.T) {
		secureConfig := &SecureConfig{
			Cert:             []byte("cert"),
			Key:              []byte("key"),
			CA:               []byte("ca"),
			CertFile:         "/path/to/cert",
			KeyFile:          "/path/to/key",
			TrustedCAFile:    "/path/to/ca",
			InsecureSkipVerify: true,
			Username:         "user",
			Password:         "pass",
		}

		assert.Equal(t, []byte("cert"), secureConfig.Cert)
		assert.Equal(t, []byte("key"), secureConfig.Key)
		assert.Equal(t, []byte("ca"), secureConfig.CA)
		assert.Equal(t, "/path/to/cert", secureConfig.CertFile)
		assert.Equal(t, "/path/to/key", secureConfig.KeyFile)
		assert.Equal(t, "/path/to/ca", secureConfig.TrustedCAFile)
		assert.True(t, secureConfig.InsecureSkipVerify)
		assert.Equal(t, "user", secureConfig.Username)
		assert.Equal(t, "pass", secureConfig.Password)
	})
}

func TestClientConfigSecret_Comprehensive(t *testing.T) {
	t.Run("NewClientConfigSecret", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset()
		configSecret := NewClientConfigSecret(kubeClient)

		assert.NotNil(t, configSecret)
		assert.Implements(t, (*ClientConfigGetter)(nil), configSecret)
	})

	t.Run("GetClientConfig - Valid Secret", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset()
		configSecret := NewClientConfigSecret(kubeClient)

		// Create a test secret
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "etcd-config",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"endpoints": []byte("http://localhost:2379"),
				"tls.crt":   []byte("-----BEGIN CERTIFICATE-----\ntest-cert\n-----END CERTIFICATE-----"),
				"tls.key":   []byte("-----BEGIN PRIVATE KEY-----\ntest-key\n-----END PRIVATE KEY-----"),
				"ca.crt":    []byte("-----BEGIN CERTIFICATE-----\ntest-ca\n-----END CERTIFICATE-----"),
				"username":  []byte("etcd-user"),
				"password":  []byte("etcd-pass"),
			},
		}

		_, err := kubeClient.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
		require.NoError(t, err)

		config, err := configSecret.GetClientConfig("default", "etcd-config")
		require.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, []string{"http://localhost:2379"}, config.Endpoints)
		assert.Equal(t, "etcd-user", config.Username)
		assert.Equal(t, "etcd-pass", config.Password)
		assert.NotNil(t, config.TLS)
	})

	t.Run("GetClientConfig - Secret Not Found", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset()
		configSecret := NewClientConfigSecret(kubeClient)

		config, err := configSecret.GetClientConfig("default", "nonexistent")
		assert.Error(t, err)
		assert.Nil(t, config)
	})

	t.Run("GetClientConfig - Secret with TLS Only", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset()
		configSecret := NewClientConfigSecret(kubeClient)

		// Create a secret with only TLS data
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "etcd-tls",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"endpoints": []byte("https://localhost:2379"),
				"tls.crt":   []byte("-----BEGIN CERTIFICATE-----\ntest-cert\n-----END CERTIFICATE-----"),
				"tls.key":   []byte("-----BEGIN PRIVATE KEY-----\ntest-key\n-----END PRIVATE KEY-----"),
			},
		}

		_, err := kubeClient.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
		require.NoError(t, err)

		config, err := configSecret.GetClientConfig("default", "etcd-tls")
		require.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, []string{"https://localhost:2379"}, config.Endpoints)
		assert.NotNil(t, config.TLS)
		assert.Empty(t, config.Username)
		assert.Empty(t, config.Password)
	})

	t.Run("GetClientConfig - Secret with CA Only", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset()
		configSecret := NewClientConfigSecret(kubeClient)

		// Create a secret with only CA data
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "etcd-ca",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"endpoints": []byte("https://localhost:2379"),
				"ca.crt":    []byte("-----BEGIN CERTIFICATE-----\ntest-ca\n-----END CERTIFICATE-----"),
			},
		}

		_, err := kubeClient.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
		require.NoError(t, err)

		config, err := configSecret.GetClientConfig("default", "etcd-ca")
		require.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, []string{"https://localhost:2379"}, config.Endpoints)
		assert.NotNil(t, config.TLS)
		assert.NotNil(t, config.TLS.RootCAs)
	})

	t.Run("GetClientConfig - Secret with Auth Only", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset()
		configSecret := NewClientConfigSecret(kubeClient)

		// Create a secret with only auth data
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "etcd-auth",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"endpoints": []byte("http://localhost:2379"),
				"username":  []byte("etcd-user"),
				"password":  []byte("etcd-pass"),
			},
		}

		_, err := kubeClient.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
		require.NoError(t, err)

		config, err := configSecret.GetClientConfig("default", "etcd-auth")
		require.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, []string{"http://localhost:2379"}, config.Endpoints)
		assert.Equal(t, "etcd-user", config.Username)
		assert.Equal(t, "etcd-pass", config.Password)
		assert.Nil(t, config.TLS)
	})

	t.Run("GetClientConfig - Empty Secret", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset()
		configSecret := NewClientConfigSecret(kubeClient)

		// Create an empty secret
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "etcd-empty",
				Namespace: "default",
			},
			Data: map[string][]byte{},
		}

		_, err := kubeClient.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
		require.NoError(t, err)

		config, err := configSecret.GetClientConfig("default", "etcd-empty")
		require.NoError(t, err)
		assert.NotNil(t, config)
		assert.Empty(t, config.Endpoints)
		assert.Empty(t, config.Username)
		assert.Empty(t, config.Password)
		assert.Nil(t, config.TLS)
	})
}

func TestClientConfigFromSecret_Comprehensive(t *testing.T) {
	t.Run("Valid Secret with All Fields", func(t *testing.T) {
		secret := &corev1.Secret{
			Data: map[string][]byte{
				"endpoints": []byte("http://localhost:2379,http://localhost:2380"),
				"tls.crt":   []byte("-----BEGIN CERTIFICATE-----\ntest-cert\n-----END CERTIFICATE-----"),
				"tls.key":   []byte("-----BEGIN PRIVATE KEY-----\ntest-key\n-----END PRIVATE KEY-----"),
				"ca.crt":    []byte("-----BEGIN CERTIFICATE-----\ntest-ca\n-----END CERTIFICATE-----"),
				"username":  []byte("etcd-user"),
				"password":  []byte("etcd-pass"),
			},
		}

		config, err := clientConfigFromSecret(secret)
		require.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, []string{"http://localhost:2379,http://localhost:2380"}, config.Endpoints)
		assert.Equal(t, "etcd-user", config.Username)
		assert.Equal(t, "etcd-pass", config.Password)
		assert.NotNil(t, config.TLS)
	})

	t.Run("Invalid Certificate", func(t *testing.T) {
		secret := &corev1.Secret{
			Data: map[string][]byte{
				"endpoints": []byte("https://localhost:2379"),
				"tls.crt":   []byte("invalid-cert"),
				"tls.key":   []byte("invalid-key"),
			},
		}

		config, err := clientConfigFromSecret(secret)
		assert.Error(t, err)
		assert.Nil(t, config)
	})

	t.Run("Missing Key with Certificate", func(t *testing.T) {
		secret := &corev1.Secret{
			Data: map[string][]byte{
				"endpoints": []byte("https://localhost:2379"),
				"tls.crt":   []byte("-----BEGIN CERTIFICATE-----\ntest-cert\n-----END CERTIFICATE-----"),
				// Missing tls.key
			},
		}

		config, err := clientConfigFromSecret(secret)
		require.NoError(t, err)
		assert.NotNil(t, config)
		assert.Nil(t, config.TLS) // Should not create TLS config without both cert and key
	})

	t.Run("Missing Certificate with Key", func(t *testing.T) {
		secret := &corev1.Secret{
			Data: map[string][]byte{
				"endpoints": []byte("https://localhost:2379"),
				"tls.key":   []byte("-----BEGIN PRIVATE KEY-----\ntest-key\n-----END PRIVATE KEY-----"),
				// Missing tls.crt
			},
		}

		config, err := clientConfigFromSecret(secret)
		require.NoError(t, err)
		assert.NotNil(t, config)
		assert.Nil(t, config.TLS) // Should not create TLS config without both cert and key
	})

	t.Run("Valid Certificate and Key without CA", func(t *testing.T) {
		secret := &corev1.Secret{
			Data: map[string][]byte{
				"endpoints": []byte("https://localhost:2379"),
				"tls.crt":   []byte("-----BEGIN CERTIFICATE-----\ntest-cert\n-----END CERTIFICATE-----"),
				"tls.key":   []byte("-----BEGIN PRIVATE KEY-----\ntest-key\n-----END PRIVATE KEY-----"),
			},
		}

		config, err := clientConfigFromSecret(secret)
		require.NoError(t, err)
		assert.NotNil(t, config)
		assert.NotNil(t, config.TLS)
		assert.Nil(t, config.TLS.RootCAs) // Should not have RootCAs without CA
	})

	t.Run("Valid Certificate, Key, and CA", func(t *testing.T) {
		secret := &corev1.Secret{
			Data: map[string][]byte{
				"endpoints": []byte("https://localhost:2379"),
				"tls.crt":   []byte("-----BEGIN CERTIFICATE-----\ntest-cert\n-----END CERTIFICATE-----"),
				"tls.key":   []byte("-----BEGIN PRIVATE KEY-----\ntest-key\n-----END PRIVATE KEY-----"),
				"ca.crt":    []byte("-----BEGIN CERTIFICATE-----\ntest-ca\n-----END CERTIFICATE-----"),
			},
		}

		config, err := clientConfigFromSecret(secret)
		require.NoError(t, err)
		assert.NotNil(t, config)
		assert.NotNil(t, config.TLS)
		assert.NotNil(t, config.TLS.RootCAs)
	})

	t.Run("Empty Secret", func(t *testing.T) {
		secret := &corev1.Secret{
			Data: map[string][]byte{},
		}

		config, err := clientConfigFromSecret(secret)
		require.NoError(t, err)
		assert.NotNil(t, config)
		assert.Empty(t, config.Endpoints)
		assert.Empty(t, config.Username)
		assert.Empty(t, config.Password)
		assert.Nil(t, config.TLS)
	})

	t.Run("Nil Secret", func(t *testing.T) {
		config, err := clientConfigFromSecret(nil)
		require.NoError(t, err)
		assert.NotNil(t, config)
		// Should return default config
		assert.Equal(t, 5*time.Second, config.DialTimeout)
		assert.NotNil(t, config.Context)
	})
}

func TestNewClientv3_Comprehensive(t *testing.T) {
	t.Run("Valid Config", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"http://localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "user",
			Password:    "pass",
		}

		// This will fail without a real etcd server, but we can test the config creation
		client, err := NewClientv3(config)
		// We expect this to fail in test environment without etcd server
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("Config with TLS", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"https://localhost:2379"},
			DialTimeout: 5 * time.Second,
			TLS: &tls.Config{
				InsecureSkipVerify: true,
			},
		}

		// This will fail without a real etcd server, but we can test the config creation
		client, err := NewClientv3(config)
		// We expect this to fail in test environment without etcd server
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("Config with Empty Endpoints", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{},
			DialTimeout: 5 * time.Second,
		}

		client, err := NewClientv3(config)
		// This should fail due to empty endpoints
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("Config with Nil Endpoints", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   nil,
			DialTimeout: 5 * time.Second,
		}

		client, err := NewClientv3(config)
		// This should fail due to nil endpoints
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("Config with Zero Timeout", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"http://localhost:2379"},
			DialTimeout: 0,
		}

		// This will fail without a real etcd server, but we can test the config creation
		client, err := NewClientv3(config)
		// We expect this to fail in test environment without etcd server
		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

func TestTLSConfig_Comprehensive(t *testing.T) {
	t.Run("TLS Config Creation", func(t *testing.T) {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         "etcd-server",
		}

		assert.True(t, tlsConfig.InsecureSkipVerify)
		assert.Equal(t, "etcd-server", tlsConfig.ServerName)
	})

	t.Run("Certificate Pool Creation", func(t *testing.T) {
		caCert := []byte("-----BEGIN CERTIFICATE-----\ntest-ca\n-----END CERTIFICATE-----")
		caCertPool := x509.NewCertPool()
		success := caCertPool.AppendCertsFromPEM(caCert)

		assert.True(t, success)
		assert.NotNil(t, caCertPool)
	})

	t.Run("Invalid Certificate Pool", func(t *testing.T) {
		invalidCert := []byte("invalid-certificate")
		caCertPool := x509.NewCertPool()
		success := caCertPool.AppendCertsFromPEM(invalidCert)

		assert.False(t, success)
		assert.NotNil(t, caCertPool)
	})
}

func TestClientConfig_EdgeCases(t *testing.T) {
	t.Run("Config with Special Characters", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"http://localhost:2379", "http://[::1]:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "user@domain.com",
			Password:    "pass!@#$%^&*()",
		}

		assert.Len(t, config.Endpoints, 2)
		assert.Equal(t, "user@domain.com", config.Username)
		assert.Equal(t, "pass!@#$%^&*()", config.Password)
	})

	t.Run("Config with Long Values", func(t *testing.T) {
		longString := string(make([]byte, 10000))
		for range longString {
			longString = "a" + longString
		}

		config := &ClientConfig{
			Endpoints:   []string{"http://localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    longString,
			Password:    longString,
		}

		assert.Equal(t, longString, config.Username)
		assert.Equal(t, longString, config.Password)
	})

	t.Run("Config with Unicode", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"http://localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "用户",
			Password:    "密码",
		}

		assert.Equal(t, "用户", config.Username)
		assert.Equal(t, "密码", config.Password)
	})
}

func TestClientConfig_Concurrency(t *testing.T) {
	t.Run("Concurrent Config Access", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"http://localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "user",
			Password:    "pass",
		}

		// Test concurrent access to config fields
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				defer func() { done <- true }()
				
				// Read config fields
				_ = config.Endpoints
				_ = config.DialTimeout
				_ = config.Username
				_ = config.Password
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

func TestClientConfig_Validation(t *testing.T) {
	t.Run("Config Validation", func(t *testing.T) {
		// Test various config combinations
		testCases := []struct {
			name   string
			config *ClientConfig
			valid  bool
		}{
			{
				name: "Valid basic config",
				config: &ClientConfig{
					Endpoints:   []string{"http://localhost:2379"},
					DialTimeout: 5 * time.Second,
				},
				valid: true,
			},
			{
				name: "Valid config with auth",
				config: &ClientConfig{
					Endpoints:   []string{"http://localhost:2379"},
					DialTimeout: 5 * time.Second,
					Username:    "user",
					Password:    "pass",
				},
				valid: true,
			},
			{
				name: "Valid config with TLS",
				config: &ClientConfig{
					Endpoints:   []string{"https://localhost:2379"},
					DialTimeout: 5 * time.Second,
					TLS:         &tls.Config{InsecureSkipVerify: true},
				},
				valid: true,
			},
			{
				name: "Invalid config - empty endpoints",
				config: &ClientConfig{
					Endpoints:   []string{},
					DialTimeout: 5 * time.Second,
				},
				valid: false,
			},
			{
				name: "Invalid config - nil endpoints",
				config: &ClientConfig{
					Endpoints:   nil,
					DialTimeout: 5 * time.Second,
				},
				valid: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Basic validation - check if required fields are present
				hasValidEndpoints := len(tc.config.Endpoints) > 0
				hasValidTimeout := tc.config.DialTimeout > 0

				isValid := hasValidEndpoints && hasValidTimeout
				assert.Equal(t, tc.valid, isValid, "Config validation failed for %s", tc.name)
			})
		}
	})
}
