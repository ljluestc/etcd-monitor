package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestClientConfig_Comprehensive(t *testing.T) {
	t.Run("ClientConfig Creation", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "test-user",
			Password:    "test-pass",
		}

		assert.NotNil(t, config)
		assert.Equal(t, []string{"localhost:2379"}, config.Endpoints)
		assert.Equal(t, 5*time.Second, config.DialTimeout)
		assert.Equal(t, "test-user", config.Username)
		assert.Equal(t, "test-pass", config.Password)
	})

	t.Run("ClientConfig with multiple endpoints", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"localhost:2379", "localhost:2380", "localhost:2381"},
			DialTimeout: 10 * time.Second,
		}

		assert.NotNil(t, config)
		assert.Len(t, config.Endpoints, 3)
		assert.Equal(t, "localhost:2379", config.Endpoints[0])
		assert.Equal(t, "localhost:2380", config.Endpoints[1])
		assert.Equal(t, "localhost:2381", config.Endpoints[2])
	})

	t.Run("ClientConfig with empty endpoints", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{},
			DialTimeout: 5 * time.Second,
		}

		assert.NotNil(t, config)
		assert.Len(t, config.Endpoints, 0)
	})

	t.Run("ClientConfig with nil endpoints", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   nil,
			DialTimeout: 5 * time.Second,
		}

		assert.NotNil(t, config)
		assert.Nil(t, config.Endpoints)
	})

	t.Run("ClientConfig with zero timeout", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 0,
		}

		assert.NotNil(t, config)
		assert.Equal(t, time.Duration(0), config.DialTimeout)
	})

	t.Run("ClientConfig with negative timeout", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: -1 * time.Second,
		}

		assert.NotNil(t, config)
		assert.Equal(t, -1*time.Second, config.DialTimeout)
	})

	t.Run("ClientConfig with empty username", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "",
			Password:    "test-pass",
		}

		assert.NotNil(t, config)
		assert.Equal(t, "", config.Username)
		assert.Equal(t, "test-pass", config.Password)
	})

	t.Run("ClientConfig with empty password", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "test-user",
			Password:    "",
		}

		assert.NotNil(t, config)
		assert.Equal(t, "test-user", config.Username)
		assert.Equal(t, "", config.Password)
	})
}

func TestSecureConfig_Comprehensive(t *testing.T) {
	t.Run("SecureConfig Creation", func(t *testing.T) {
		secureConfig := &SecureConfig{
			CertFile:           "/path/to/cert.pem",
			KeyFile:            "/path/to/key.pem",
			TrustedCAFile:      "/path/to/ca.pem",
			InsecureSkipVerify: true,
			Username:           "user",
			Password:           "pass",
		}

		assert.NotNil(t, secureConfig)
		assert.Equal(t, "/path/to/cert.pem", secureConfig.CertFile)
		assert.Equal(t, "/path/to/key.pem", secureConfig.KeyFile)
		assert.Equal(t, "/path/to/ca.pem", secureConfig.TrustedCAFile)
		assert.True(t, secureConfig.InsecureSkipVerify)
		assert.Equal(t, "user", secureConfig.Username)
		assert.Equal(t, "pass", secureConfig.Password)
	})

	t.Run("SecureConfig with empty cert file", func(t *testing.T) {
		secureConfig := &SecureConfig{
			CertFile:           "",
			KeyFile:            "/path/to/key.pem",
			TrustedCAFile:      "/path/to/ca.pem",
			InsecureSkipVerify: false,
		}

		assert.NotNil(t, secureConfig)
		assert.Equal(t, "", secureConfig.CertFile)
		assert.Equal(t, "/path/to/key.pem", secureConfig.KeyFile)
		assert.Equal(t, "/path/to/ca.pem", secureConfig.TrustedCAFile)
		assert.False(t, secureConfig.InsecureSkipVerify)
	})

	t.Run("SecureConfig with empty key file", func(t *testing.T) {
		secureConfig := &SecureConfig{
			CertFile:           "/path/to/cert.pem",
			KeyFile:            "",
			TrustedCAFile:      "/path/to/ca.pem",
			InsecureSkipVerify: false,
		}

		assert.NotNil(t, secureConfig)
		assert.Equal(t, "/path/to/cert.pem", secureConfig.CertFile)
		assert.Equal(t, "", secureConfig.KeyFile)
		assert.Equal(t, "/path/to/ca.pem", secureConfig.TrustedCAFile)
		assert.False(t, secureConfig.InsecureSkipVerify)
	})

	t.Run("SecureConfig with empty CA file", func(t *testing.T) {
		secureConfig := &SecureConfig{
			CertFile:           "/path/to/cert.pem",
			KeyFile:            "/path/to/key.pem",
			TrustedCAFile:      "",
			InsecureSkipVerify: false,
		}

		assert.NotNil(t, secureConfig)
		assert.Equal(t, "/path/to/cert.pem", secureConfig.CertFile)
		assert.Equal(t, "/path/to/key.pem", secureConfig.KeyFile)
		assert.Equal(t, "", secureConfig.TrustedCAFile)
		assert.False(t, secureConfig.InsecureSkipVerify)
	})

	t.Run("SecureConfig with all empty files", func(t *testing.T) {
		secureConfig := &SecureConfig{
			CertFile:           "",
			KeyFile:            "",
			TrustedCAFile:      "",
			InsecureSkipVerify: true,
		}

		assert.NotNil(t, secureConfig)
		assert.Equal(t, "", secureConfig.CertFile)
		assert.Equal(t, "", secureConfig.KeyFile)
		assert.Equal(t, "", secureConfig.TrustedCAFile)
		assert.True(t, secureConfig.InsecureSkipVerify)
	})

	t.Run("SecureConfig with nil files", func(t *testing.T) {
		secureConfig := &SecureConfig{
			CertFile:           "",
			KeyFile:            "",
			TrustedCAFile:      "",
			InsecureSkipVerify: false,
		}

		assert.NotNil(t, secureConfig)
		assert.Equal(t, "", secureConfig.CertFile)
		assert.Equal(t, "", secureConfig.KeyFile)
		assert.Equal(t, "", secureConfig.TrustedCAFile)
		assert.False(t, secureConfig.InsecureSkipVerify)
	})
}

func TestClientConfigSecret_Comprehensive(t *testing.T) {
	t.Run("NewClientConfigSecret", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset()
		configSecret := NewClientConfigSecret(kubeClient)

		assert.NotNil(t, configSecret)
		assert.Implements(t, (*ClientConfigGetter)(nil), configSecret)
	})

	t.Run("NewClientConfigSecret with nil client", func(t *testing.T) {
		configSecret := NewClientConfigSecret(nil)

		assert.NotNil(t, configSecret)
		assert.Implements(t, (*ClientConfigGetter)(nil), configSecret)
	})

	t.Run("GetClientConfig - Secret Not Found", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset()
		configSecret := NewClientConfigSecret(kubeClient)

		config, err := configSecret.GetClientConfig("default", "non-existent-secret")
		assert.NoError(t, err) // Should not error, just return default config
		assert.NotNil(t, config)
		assert.Empty(t, config.Endpoints)
	})

	t.Run("GetClientConfig - Secret with Endpoints Only", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset()
		configSecret := NewClientConfigSecret(kubeClient)

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "endpoints-secret",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"endpoints": []byte("http://etcd-0:2379,http://etcd-1:2379"),
			},
		}
		_, err := kubeClient.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
		require.NoError(t, err)

		config, err := configSecret.GetClientConfig("default", "endpoints-secret")
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, []string{"http://etcd-0:2379", "http://etcd-1:2379"}, config.Endpoints)
	})

	t.Run("GetClientConfig - Secret with Username/Password Only", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset()
		configSecret := NewClientConfigSecret(kubeClient)

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "auth-secret",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"username": []byte("testuser"),
				"password": []byte("testpass"),
			},
		}
		_, err := kubeClient.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
		require.NoError(t, err)

		config, err := configSecret.GetClientConfig("default", "auth-secret")
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, "testuser", config.Username)
		assert.Equal(t, "testpass", config.Password)
	})

	t.Run("GetClientConfig - Secret with DialTimeout", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset()
		configSecret := NewClientConfigSecret(kubeClient)

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "timeout-secret",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"dialTimeout": []byte("10s"),
			},
		}
		_, err := kubeClient.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
		require.NoError(t, err)

		config, err := configSecret.GetClientConfig("default", "timeout-secret")
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, 10*time.Second, config.DialTimeout)
	})

	t.Run("GetClientConfig - Secret with Invalid DialTimeout", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset()
		configSecret := NewClientConfigSecret(kubeClient)

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "invalid-timeout-secret",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"dialTimeout": []byte("invalid"),
			},
		}
		_, err := kubeClient.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
		require.NoError(t, err)

		config, err := configSecret.GetClientConfig("default", "invalid-timeout-secret")
		assert.Error(t, err) // Should error due to invalid duration
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "failed to parse dialTimeout")
	})

	t.Run("GetClientConfig - Secret with All Fields", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset()
		configSecret := NewClientConfigSecret(kubeClient)

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "complete-secret",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"endpoints":   []byte("http://etcd-0:2379,http://etcd-1:2379"),
				"username":    []byte("testuser"),
				"password":    []byte("testpass"),
				"dialTimeout": []byte("15s"),
			},
		}
		_, err := kubeClient.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
		require.NoError(t, err)

		config, err := configSecret.GetClientConfig("default", "complete-secret")
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, []string{"http://etcd-0:2379", "http://etcd-1:2379"}, config.Endpoints)
		assert.Equal(t, "testuser", config.Username)
		assert.Equal(t, "testpass", config.Password)
		assert.Equal(t, 15*time.Second, config.DialTimeout)
	})

	t.Run("GetClientConfig - Secret with Empty Values", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset()
		configSecret := NewClientConfigSecret(kubeClient)

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "empty-secret",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"endpoints":   []byte(""),
				"username":    []byte(""),
				"password":    []byte(""),
				"dialTimeout": []byte(""),
			},
		}
		_, err := kubeClient.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
		require.NoError(t, err)

		config, err := configSecret.GetClientConfig("default", "empty-secret")
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Empty(t, config.Endpoints)
		assert.Empty(t, config.Username)
		assert.Empty(t, config.Password)
		assert.Equal(t, time.Duration(0), config.DialTimeout)
	})

	t.Run("GetClientConfig - Secret with Whitespace Values", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset()
		configSecret := NewClientConfigSecret(kubeClient)

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "whitespace-secret",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"endpoints":   []byte("  http://etcd-0:2379  ,  http://etcd-1:2379  "),
				"username":    []byte("  testuser  "),
				"password":    []byte("  testpass  "),
				"dialTimeout": []byte("  10s  "),
			},
		}
		_, err := kubeClient.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
		require.NoError(t, err)

		config, err := configSecret.GetClientConfig("default", "whitespace-secret")
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, []string{"http://etcd-0:2379", "http://etcd-1:2379"}, config.Endpoints)
		assert.Equal(t, "testuser", config.Username)
		assert.Equal(t, "testpass", config.Password)
		assert.Equal(t, 10*time.Second, config.DialTimeout)
	})

	t.Run("GetClientConfig - Secret with Special Characters", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset()
		configSecret := NewClientConfigSecret(kubeClient)

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "special-secret",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"endpoints":   []byte("http://etcd-0:2379,http://etcd-1:2379,http://etcd-2:2379"),
				"username":    []byte("user@domain.com"),
				"password":    []byte("pass!@#$%^&*()"),
				"dialTimeout": []byte("30s"),
			},
		}
		_, err := kubeClient.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
		require.NoError(t, err)

		config, err := configSecret.GetClientConfig("default", "special-secret")
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, []string{"http://etcd-0:2379", "http://etcd-1:2379", "http://etcd-2:2379"}, config.Endpoints)
		assert.Equal(t, "user@domain.com", config.Username)
		assert.Equal(t, "pass!@#$%^&*()", config.Password)
		assert.Equal(t, 30*time.Second, config.DialTimeout)
	})

	t.Run("GetClientConfig - Secret with Large Values", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset()
		configSecret := NewClientConfigSecret(kubeClient)

		// Create a large password
		largePassword := make([]byte, 1000)
		for i := range largePassword {
			largePassword[i] = byte('a' + i%26)
		}

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "large-secret",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"endpoints":   []byte("http://etcd-0:2379,http://etcd-1:2379"),
				"username":    []byte("testuser"),
				"password":    largePassword,
				"dialTimeout": []byte("10s"),
			},
		}
		_, err := kubeClient.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
		require.NoError(t, err)

		config, err := configSecret.GetClientConfig("default", "large-secret")
		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, []string{"http://etcd-0:2379", "http://etcd-1:2379"}, config.Endpoints)
		assert.Equal(t, "testuser", config.Username)
		assert.Equal(t, string(largePassword), config.Password)
		assert.Equal(t, 10*time.Second, config.DialTimeout)
	})
}

func TestGetTLSConfig_Comprehensive(t *testing.T) {
	t.Run("GetTLSConfig - No TLS Files", func(t *testing.T) {
		tlsConfig, err := GetTLSConfig("", "", "", false)
		assert.NoError(t, err)
		assert.Nil(t, tlsConfig)
	})

	t.Run("GetTLSConfig - InsecureSkipVerify", func(t *testing.T) {
		tlsConfig, err := GetTLSConfig("", "", "", true)
		assert.NoError(t, err)
		assert.NotNil(t, tlsConfig)
		assert.True(t, tlsConfig.InsecureSkipVerify)
	})

	t.Run("GetTLSConfig - Invalid Cert File", func(t *testing.T) {
		tlsConfig, err := GetTLSConfig("/nonexistent/cert.pem", "/nonexistent/key.pem", "", false)
		assert.Error(t, err)
		assert.Nil(t, tlsConfig)
		assert.Contains(t, err.Error(), "failed to load client certs")
	})

	t.Run("GetTLSConfig - Invalid CA File", func(t *testing.T) {
		// Create dummy cert and key files
		certFile := createTempFile(t, "cert.pem", []byte("-----BEGIN CERTIFICATE-----\nCERT\n-----END CERTIFICATE-----"))
		keyFile := createTempFile(t, "key.pem", []byte("-----BEGIN PRIVATE KEY-----\nKEY\n-----END PRIVATE KEY-----"))
		defer os.Remove(certFile)
		defer os.Remove(keyFile)

		tlsConfig, err := GetTLSConfig(certFile, keyFile, "/nonexistent/ca.pem", false)
		assert.Error(t, err)
		assert.Nil(t, tlsConfig)
		assert.Contains(t, err.Error(), "failed to parse root certificate")
	})

	t.Run("GetTLSConfig - Valid TLS Files", func(t *testing.T) {
		// Create dummy cert, key, and ca files
		certFile := createTempFile(t, "cert.pem", []byte("-----BEGIN CERTIFICATE-----\nCERT\n-----END CERTIFICATE-----"))
		keyFile := createTempFile(t, "key.pem", []byte("-----BEGIN PRIVATE KEY-----\nKEY\n-----END PRIVATE KEY-----"))
		caFile := createTempFile(t, "ca.pem", []byte("-----BEGIN CERTIFICATE-----\nCA\n-----END CERTIFICATE-----"))
		defer os.Remove(certFile)
		defer os.Remove(keyFile)
		defer os.Remove(caFile)

		// Mock x509.SystemCertPool to return a non-nil pool
		oldSystemCertPool := x509SystemCertPool
		x509SystemCertPool = func() (*x509.CertPool, error) {
			return x509.NewCertPool(), nil
		}
		defer func() { x509SystemCertPool = oldSystemCertPool }()

		tlsConfig, err := GetTLSConfig(certFile, keyFile, caFile, false)
		assert.NoError(t, err)
		assert.NotNil(t, tlsConfig)
		assert.NotNil(t, tlsConfig.Certificates)
		assert.Len(t, tlsConfig.Certificates, 1)
		assert.NotNil(t, tlsConfig.RootCAs)
	})

	t.Run("GetTLSConfig - Cert File Only", func(t *testing.T) {
		certFile := createTempFile(t, "cert.pem", []byte("-----BEGIN CERTIFICATE-----\nCERT\n-----END CERTIFICATE-----"))
		defer os.Remove(certFile)

		tlsConfig, err := GetTLSConfig(certFile, "", "", false)
		assert.Error(t, err)
		assert.Nil(t, tlsConfig)
		assert.Contains(t, err.Error(), "failed to load client certs")
	})

	t.Run("GetTLSConfig - Key File Only", func(t *testing.T) {
		keyFile := createTempFile(t, "key.pem", []byte("-----BEGIN PRIVATE KEY-----\nKEY\n-----END PRIVATE KEY-----"))
		defer os.Remove(keyFile)

		tlsConfig, err := GetTLSConfig("", keyFile, "", false)
		assert.Error(t, err)
		assert.Nil(t, tlsConfig)
		assert.Contains(t, err.Error(), "failed to load client certs")
	})

	t.Run("GetTLSConfig - CA File Only", func(t *testing.T) {
		caFile := createTempFile(t, "ca.pem", []byte("-----BEGIN CERTIFICATE-----\nCA\n-----END CERTIFICATE-----"))
		defer os.Remove(caFile)

		// Mock x509.SystemCertPool to return a non-nil pool
		oldSystemCertPool := x509SystemCertPool
		x509SystemCertPool = func() (*x509.CertPool, error) {
			return x509.NewCertPool(), nil
		}
		defer func() { x509SystemCertPool = oldSystemCertPool }()

		tlsConfig, err := GetTLSConfig("", "", caFile, false)
		assert.NoError(t, err)
		assert.NotNil(t, tlsConfig)
		assert.NotNil(t, tlsConfig.RootCAs)
		assert.Len(t, tlsConfig.Certificates, 0)
	})

	t.Run("GetTLSConfig - Empty Files", func(t *testing.T) {
		certFile := createTempFile(t, "cert.pem", []byte(""))
		keyFile := createTempFile(t, "key.pem", []byte(""))
		caFile := createTempFile(t, "ca.pem", []byte(""))
		defer os.Remove(certFile)
		defer os.Remove(keyFile)
		defer os.Remove(caFile)

		tlsConfig, err := GetTLSConfig(certFile, keyFile, caFile, false)
		assert.Error(t, err)
		assert.Nil(t, tlsConfig)
	})

	t.Run("GetTLSConfig - Invalid Cert Format", func(t *testing.T) {
		certFile := createTempFile(t, "cert.pem", []byte("invalid cert format"))
		keyFile := createTempFile(t, "key.pem", []byte("-----BEGIN PRIVATE KEY-----\nKEY\n-----END PRIVATE KEY-----"))
		defer os.Remove(certFile)
		defer os.Remove(keyFile)

		tlsConfig, err := GetTLSConfig(certFile, keyFile, "", false)
		assert.Error(t, err)
		assert.Nil(t, tlsConfig)
		assert.Contains(t, err.Error(), "failed to load client certs")
	})

	t.Run("GetTLSConfig - Invalid Key Format", func(t *testing.T) {
		certFile := createTempFile(t, "cert.pem", []byte("-----BEGIN CERTIFICATE-----\nCERT\n-----END CERTIFICATE-----"))
		keyFile := createTempFile(t, "key.pem", []byte("invalid key format"))
		defer os.Remove(certFile)
		defer os.Remove(keyFile)

		tlsConfig, err := GetTLSConfig(certFile, keyFile, "", false)
		assert.Error(t, err)
		assert.Nil(t, tlsConfig)
		assert.Contains(t, err.Error(), "failed to load client certs")
	})

	t.Run("GetTLSConfig - Invalid CA Format", func(t *testing.T) {
		certFile := createTempFile(t, "cert.pem", []byte("-----BEGIN CERTIFICATE-----\nCERT\n-----END CERTIFICATE-----"))
		keyFile := createTempFile(t, "key.pem", []byte("-----BEGIN PRIVATE KEY-----\nKEY\n-----END PRIVATE KEY-----"))
		caFile := createTempFile(t, "ca.pem", []byte("invalid ca format"))
		defer os.Remove(certFile)
		defer os.Remove(keyFile)
		defer os.Remove(caFile)

		tlsConfig, err := GetTLSConfig(certFile, keyFile, caFile, false)
		assert.Error(t, err)
		assert.Nil(t, tlsConfig)
		assert.Contains(t, err.Error(), "failed to parse root certificate")
	})

	t.Run("GetTLSConfig - SystemCertPool Error", func(t *testing.T) {
		certFile := createTempFile(t, "cert.pem", []byte("-----BEGIN CERTIFICATE-----\nCERT\n-----END CERTIFICATE-----"))
		keyFile := createTempFile(t, "key.pem", []byte("-----BEGIN PRIVATE KEY-----\nKEY\n-----END PRIVATE KEY-----"))
		caFile := createTempFile(t, "ca.pem", []byte("-----BEGIN CERTIFICATE-----\nCA\n-----END CERTIFICATE-----"))
		defer os.Remove(certFile)
		defer os.Remove(keyFile)
		defer os.Remove(caFile)

		// Mock x509.SystemCertPool to return an error
		oldSystemCertPool := x509SystemCertPool
		x509SystemCertPool = func() (*x509.CertPool, error) {
			return nil, assert.AnError
		}
		defer func() { x509SystemCertPool = oldSystemCertPool }()

		tlsConfig, err := GetTLSConfig(certFile, keyFile, caFile, false)
		assert.Error(t, err)
		assert.Nil(t, tlsConfig)
		assert.Contains(t, err.Error(), "failed to get system cert pool")
	})
}

func TestClientConfig_EdgeCases(t *testing.T) {
	t.Run("ClientConfig with very long endpoint list", func(t *testing.T) {
		endpoints := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			endpoints[i] = fmt.Sprintf("localhost:%d", 2379+i)
		}

		config := &ClientConfig{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
		}

		assert.NotNil(t, config)
		assert.Len(t, config.Endpoints, 1000)
		assert.Equal(t, "localhost:2379", config.Endpoints[0])
		assert.Equal(t, "localhost:3378", config.Endpoints[999])
	})

	t.Run("ClientConfig with very long username", func(t *testing.T) {
		longUsername := strings.Repeat("a", 10000)
		config := &ClientConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    longUsername,
		}

		assert.NotNil(t, config)
		assert.Equal(t, longUsername, config.Username)
		assert.Len(t, config.Username, 10000)
	})

	t.Run("ClientConfig with very long password", func(t *testing.T) {
		longPassword := strings.Repeat("b", 10000)
		config := &ClientConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Password:    longPassword,
		}

		assert.NotNil(t, config)
		assert.Equal(t, longPassword, config.Password)
		assert.Len(t, config.Password, 10000)
	})

	t.Run("ClientConfig with very long timeout", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 24 * time.Hour,
		}

		assert.NotNil(t, config)
		assert.Equal(t, 24*time.Hour, config.DialTimeout)
	})

	t.Run("ClientConfig with very short timeout", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 1 * time.Nanosecond,
		}

		assert.NotNil(t, config)
		assert.Equal(t, 1*time.Nanosecond, config.DialTimeout)
	})

	t.Run("ClientConfig with special characters in endpoints", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"http://etcd-0:2379", "https://etcd-1:2379", "http://etcd-2:2379"},
			DialTimeout: 5 * time.Second,
		}

		assert.NotNil(t, config)
		assert.Len(t, config.Endpoints, 3)
		assert.Equal(t, "http://etcd-0:2379", config.Endpoints[0])
		assert.Equal(t, "https://etcd-1:2379", config.Endpoints[1])
		assert.Equal(t, "http://etcd-2:2379", config.Endpoints[2])
	})

	t.Run("ClientConfig with special characters in username", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "user@domain.com",
		}

		assert.NotNil(t, config)
		assert.Equal(t, "user@domain.com", config.Username)
	})

	t.Run("ClientConfig with special characters in password", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Password:    "pass!@#$%^&*()",
		}

		assert.NotNil(t, config)
		assert.Equal(t, "pass!@#$%^&*()", config.Password)
	})

	t.Run("ClientConfig with unicode characters", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "用户",
			Password:    "密码",
		}

		assert.NotNil(t, config)
		assert.Equal(t, "用户", config.Username)
		assert.Equal(t, "密码", config.Password)
	})

	t.Run("ClientConfig with control characters", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "user\n\r\t",
			Password:    "pass\n\r\t",
		}

		assert.NotNil(t, config)
		assert.Equal(t, "user\n\r\t", config.Username)
		assert.Equal(t, "pass\n\r\t", config.Password)
	})
}

func TestClientConfig_Performance(t *testing.T) {
	t.Run("ClientConfig creation performance", func(t *testing.T) {
		start := time.Now()
		
		// Create 1000 ClientConfig instances
		for i := 0; i < 1000; i++ {
			config := &ClientConfig{
				Endpoints:   []string{"localhost:2379"},
				DialTimeout: 5 * time.Second,
				Username:    "testuser",
				Password:    "testpass",
			}
			assert.NotNil(t, config)
		}
		
		duration := time.Since(start)
		
		// Should create 1000 instances in less than 100ms
		assert.True(t, duration < 100*time.Millisecond, "Creating 1000 ClientConfig instances took too long: %v", duration)
	})

	t.Run("ClientConfig with large endpoint list performance", func(t *testing.T) {
		endpoints := make([]string, 10000)
		for i := 0; i < 10000; i++ {
			endpoints[i] = fmt.Sprintf("localhost:%d", 2379+i)
		}
		
		start := time.Now()
		
		config := &ClientConfig{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
		}
		
		duration := time.Since(start)
		
		// Should create config with 10000 endpoints in less than 10ms
		assert.True(t, duration < 10*time.Millisecond, "Creating ClientConfig with 10000 endpoints took too long: %v", duration)
		assert.Len(t, config.Endpoints, 10000)
	})
}

func TestClientConfig_Integration(t *testing.T) {
	t.Run("ClientConfig with all fields set", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"http://etcd-0:2379", "http://etcd-1:2379", "http://etcd-2:2379"},
			DialTimeout: 10 * time.Second,
			Username:    "testuser",
			Password:    "testpass",
		}

		assert.NotNil(t, config)
		assert.Len(t, config.Endpoints, 3)
		assert.Equal(t, "http://etcd-0:2379", config.Endpoints[0])
		assert.Equal(t, "http://etcd-1:2379", config.Endpoints[1])
		assert.Equal(t, "http://etcd-2:2379", config.Endpoints[2])
		assert.Equal(t, 10*time.Second, config.DialTimeout)
		assert.Equal(t, "testuser", config.Username)
		assert.Equal(t, "testpass", config.Password)
	})

	t.Run("ClientConfig with minimal fields", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		}

		assert.NotNil(t, config)
		assert.Len(t, config.Endpoints, 1)
		assert.Equal(t, "localhost:2379", config.Endpoints[0])
		assert.Equal(t, 5*time.Second, config.DialTimeout)
		assert.Empty(t, config.Username)
		assert.Empty(t, config.Password)
	})

	t.Run("ClientConfig with empty fields", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   []string{},
			DialTimeout: 0,
			Username:    "",
			Password:    "",
		}

		assert.NotNil(t, config)
		assert.Len(t, config.Endpoints, 0)
		assert.Equal(t, time.Duration(0), config.DialTimeout)
		assert.Empty(t, config.Username)
		assert.Empty(t, config.Password)
	})

	t.Run("ClientConfig with nil fields", func(t *testing.T) {
		config := &ClientConfig{
			Endpoints:   nil,
			DialTimeout: 0,
			Username:    "",
			Password:    "",
		}

		assert.NotNil(t, config)
		assert.Nil(t, config.Endpoints)
		assert.Equal(t, time.Duration(0), config.DialTimeout)
		assert.Empty(t, config.Username)
		assert.Empty(t, config.Password)
	})
}

// Helper function to create a temporary file with content
func createTempFile(t *testing.T, name string, content []byte) string {
	file, err := os.CreateTemp("", name)
	require.NoError(t, err)
	defer file.Close()
	_, err = file.Write(content)
	require.NoError(t, err)
	return file.Name()
}
