package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClientConfig contains configuration for etcd client
type ClientConfig struct {
	Endpoints     []string
	DialTimeout   time.Duration
	DialOptions   []interface{}
	Context       context.Context
	TLS           *tls.Config
	Username      string
	Password      string
	SecureConfig  *SecureConfig
}

// SecureConfig contains TLS configuration
type SecureConfig struct {
	Cert       []byte
	Key        []byte
	CA         []byte
	CertFile   string
	KeyFile    string
	TrustedCAFile string
	InsecureSkipVerify bool
	Username   string
	Password   string
}

// ClientConfigGetter is an interface to get etcd client config
type ClientConfigGetter interface {
	GetClientConfig(namespace, name string) (*ClientConfig, error)
}

// ClientConfigSecret implements ClientConfigGetter using Kubernetes secrets
type ClientConfigSecret struct {
	kubeClient kubernetes.Interface
}

// NewClientConfigSecret creates a new ClientConfigSecret
func NewClientConfigSecret(kubeClient kubernetes.Interface) ClientConfigGetter {
	return &ClientConfigSecret{
		kubeClient: kubeClient,
	}
}

// GetClientConfig gets etcd client config from Kubernetes secret
func (c *ClientConfigSecret) GetClientConfig(namespace, name string) (*ClientConfig, error) {
	secret, err := c.kubeClient.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return clientConfigFromSecret(secret)
}

// clientConfigFromSecret creates ClientConfig from Kubernetes secret
func clientConfigFromSecret(secret *corev1.Secret) (*ClientConfig, error) {
	config := &ClientConfig{
		DialTimeout: 5 * time.Second,
		Context:     context.Background(),
	}

	if endpoints, ok := secret.Data["endpoints"]; ok {
		config.Endpoints = []string{string(endpoints)}
	}

	certData := secret.Data["tls.crt"]
	keyData := secret.Data["tls.key"]
	caData := secret.Data["ca.crt"]

	if len(certData) > 0 && len(keyData) > 0 {
		cert, err := tls.X509KeyPair(certData, keyData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate: %v", err)
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		if len(caData) > 0 {
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caData)
			tlsConfig.RootCAs = caCertPool
		}

		config.TLS = tlsConfig
	}

	if username, ok := secret.Data["username"]; ok {
		config.Username = string(username)
	}

	if password, ok := secret.Data["password"]; ok {
		config.Password = string(password)
	}

	return config, nil
}

// NewClientv3 creates a new etcd v3 client
func NewClientv3(config *ClientConfig) (*clientv3.Client, error) {
	cfg := clientv3.Config{
		Endpoints:   config.Endpoints,
		DialTimeout: config.DialTimeout,
		TLS:         config.TLS,
		Username:    config.Username,
		Password:    config.Password,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		klog.Errorf("failed to create etcd v3 client: %v", err)
		return nil, err
	}

	return client, nil
}
