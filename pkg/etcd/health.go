package etcd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"k8s.io/klog/v2"
)

// HealthClient is an HTTP client for etcd health checks
type HealthClient struct {
	client *http.Client
}

// NewHealthClient creates a new health check client
func NewHealthClient(tlsConfig *tls.Config) *HealthClient {
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return &HealthClient{
		client: &http.Client{
			Transport: transport,
			Timeout:   5 * time.Second,
		},
	}
}

// HealthResponse represents the response from /health endpoint
type HealthResponse struct {
	Health string `json:"health"`
}

// VersionResponse represents the response from /version endpoint
type VersionResponse struct {
	EtcdServer  string `json:"etcdserver"`
	EtcdCluster string `json:"etcdcluster"`
}

// StatsResponse represents the response from /v2/stats/self endpoint
type StatsResponse struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	State     string `json:"state"`
	StartTime string `json:"startTime"`
	LeaderInfo struct {
		Leader    string `json:"leader"`
		Uptime    string `json:"uptime"`
		StartTime string `json:"startTime"`
	} `json:"leaderInfo"`
}

// MemberHealthy checks if an etcd member is healthy
func MemberHealthy(endpoint string, config *ClientConfig) (bool, error) {
	healthClient := NewHealthClient(config.TLS)

	url := fmt.Sprintf("%s/health", endpoint)
	resp, err := healthClient.client.Get(url)
	if err != nil {
		klog.Errorf("failed to check health for %s: %v", endpoint, err)
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var healthResp HealthResponse
	if err := json.Unmarshal(body, &healthResp); err != nil {
		return false, err
	}

	return healthResp.Health == "true", nil
}

// GetVersion gets the version of an etcd member
func GetVersion(endpoint string, config *ClientConfig) (*VersionResponse, error) {
	healthClient := NewHealthClient(config.TLS)

	url := fmt.Sprintf("%s/version", endpoint)
	resp, err := healthClient.client.Get(url)
	if err != nil {
		klog.Errorf("failed to get version for %s: %v", endpoint, err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var versionResp VersionResponse
	if err := json.Unmarshal(body, &versionResp); err != nil {
		return nil, err
	}

	return &versionResp, nil
}

// GetStats gets the statistics of an etcd member
func GetStats(endpoint string, config *ClientConfig) (*StatsResponse, error) {
	healthClient := NewHealthClient(config.TLS)

	url := fmt.Sprintf("%s/v2/stats/self", endpoint)
	resp, err := healthClient.client.Get(url)
	if err != nil {
		klog.Errorf("failed to get stats for %s: %v", endpoint, err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var statsResp StatsResponse
	if err := json.Unmarshal(body, &statsResp); err != nil {
		return nil, err
	}

	return &statsResp, nil
}
