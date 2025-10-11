package upstream

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/etcd-monitor/pingap-go/pkg/config"
)

// Manager manages upstream servers
type Manager struct {
	upstreams map[string]*Upstream
	mu        sync.RWMutex
}

// Upstream represents a backend service with multiple servers
type Upstream struct {
	Name    string
	Config  *config.ParsedUpstream
	Servers []*Server
	lb      LoadBalancer
	hc      *HealthChecker
}

// Server represents a single backend server
type Server struct {
	Addr      string
	Weight    int
	Healthy   bool
	mu        sync.RWMutex
	failures  int
	lastCheck time.Time
}

// NewManager creates a new upstream manager
func NewManager() *Manager {
	return &Manager{
		upstreams: make(map[string]*Upstream),
	}
}

// AddUpstream adds an upstream configuration
func (m *Manager) AddUpstream(name string, cfg *config.ParsedUpstream) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	servers := make([]*Server, len(cfg.Addrs))
	for i, addr := range cfg.Addrs {
		weight := 1
		if len(cfg.Weights) > i {
			weight = cfg.Weights[i]
		}
		servers[i] = &Server{
			Addr:    addr,
			Weight:  weight,
			Healthy: true,
		}
	}

	// Create load balancer based on algorithm
	var lb LoadBalancer
	switch cfg.Algorithm {
	case "least_conn":
		lb = NewLeastConnLB(servers)
	case "ip_hash":
		lb = NewIPHashLB(servers)
	case "weighted_round_robin":
		lb = NewWeightedRoundRobinLB(servers)
	default: // round_robin
		lb = NewRoundRobinLB(servers)
	}

	upstream := &Upstream{
		Name:    name,
		Config:  cfg,
		Servers: servers,
		lb:      lb,
	}

	// Start health checker if configured
	if cfg.HealthCheck != "" {
		upstream.hc = NewHealthChecker(upstream, cfg)
		go upstream.hc.Start(context.Background())
	}

	m.upstreams[name] = upstream
	return nil
}

// Get retrieves an upstream by name
func (m *Manager) Get(name string) (*Upstream, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	upstream, ok := m.upstreams[name]
	if !ok {
		return nil, fmt.Errorf("upstream %s not found", name)
	}
	return upstream, nil
}

// NextServer selects the next server using load balancing
func (u *Upstream) NextServer(r *http.Request) (*Server, error) {
	return u.lb.Next(r)
}

// URL returns the full URL for the server
func (s *Server) URL() *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   s.Addr,
	}
}

// MarkHealthy marks server as healthy
func (s *Server) MarkHealthy() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Healthy = true
	s.failures = 0
}

// MarkUnhealthy marks server as unhealthy
func (s *Server) MarkUnhealthy() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Healthy = false
	s.failures++
}

// IsHealthy returns if server is healthy
func (s *Server) IsHealthy() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Healthy
}

// CreateTransport creates an HTTP transport for the upstream
func (u *Upstream) CreateTransport() *http.Transport {
	cfg := u.Config

	// Create TLS config if needed
	var tlsConfig *tls.Config
	if cfg.TLSCert != "" && cfg.TLSKey != "" {
		// Load certificates (implementation needed)
		tlsConfig = &tls.Config{
			InsecureSkipVerify: !cfg.TLSVerify,
		}
	}

	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   cfg.ConnectionTimeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          cfg.MaxIdleConns,
		MaxIdleConnsPerHost:   cfg.MaxIdleConnsPerHost,
		IdleConnTimeout:       cfg.IdleTimeout,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       tlsConfig,
	}
}
