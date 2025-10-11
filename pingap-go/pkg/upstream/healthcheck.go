package upstream

import (
	"context"
	"net/http"
	"time"

	"github.com/etcd-monitor/pingap-go/pkg/config"
)

// HealthChecker performs health checks on upstream servers
type HealthChecker struct {
	upstream *Upstream
	config   *config.ParsedUpstream
	client   *http.Client
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(upstream *Upstream, cfg *config.ParsedUpstream) *HealthChecker {
	return &HealthChecker{
		upstream: upstream,
		config:   cfg,
		client: &http.Client{
			Timeout: cfg.ConnectionTimeout,
		},
	}
}

// Start begins health checking
func (hc *HealthChecker) Start(ctx context.Context) {
	ticker := time.NewTicker(hc.config.HealthCheckInterval)
	defer ticker.Stop()

	// Initial health check
	hc.checkAll()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			hc.checkAll()
		}
	}
}

// checkAll checks all servers
func (hc *HealthChecker) checkAll() {
	for _, server := range hc.upstream.Servers {
		go hc.checkServer(server)
	}
}

// checkServer checks a single server
func (hc *HealthChecker) checkServer(server *Server) {
	url := hc.config.HealthCheck
	if url == "" {
		url = "http://" + server.Addr + "/health"
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		server.MarkUnhealthy()
		return
	}

	resp, err := hc.client.Do(req)
	if err != nil {
		server.MarkUnhealthy()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		server.MarkHealthy()
	} else {
		server.MarkUnhealthy()
	}
}
