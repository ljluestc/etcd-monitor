package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/etcd-monitor/pingap-go/pkg/config"
	"github.com/etcd-monitor/pingap-go/pkg/plugin"
	"github.com/etcd-monitor/pingap-go/pkg/router"
	"github.com/etcd-monitor/pingap-go/pkg/upstream"
)

// Server is the main proxy server
type Server struct {
	config          *config.ParsedConfig
	router          *router.Router
	upstreamManager *upstream.Manager
	pluginChain     *plugin.Chain
	httpServer      *http.Server
}

// NewServer creates a new proxy server
func NewServer(cfg *config.ParsedConfig) (*Server, error) {
	// Create upstream manager
	upstreamMgr := upstream.NewManager()
	for name, upstream := range cfg.ParsedUpstreams {
		if err := upstreamMgr.AddUpstream(name, upstream); err != nil {
			return nil, fmt.Errorf("failed to add upstream %s: %w", name, err)
		}
	}

	// Create router
	r := router.NewRouter(cfg.ParsedLocations)

	// Create plugin chain
	pluginChain := plugin.NewChain()
	// TODO: Register plugins based on configuration

	srv := &Server{
		config:          cfg,
		router:          r,
		upstreamManager: upstreamMgr,
		pluginChain:     pluginChain,
	}

	return srv, nil
}

// Start starts all configured servers
func (s *Server) Start() error {
	for name, serverCfg := range s.config.ParsedServers {
		go func(n string, cfg *config.ParsedServer) {
			handler := s.createHandler(cfg)

			server := &http.Server{
				Addr:           cfg.Addr,
				Handler:        handler,
				ReadTimeout:    cfg.ReadTimeout,
				WriteTimeout:   cfg.WriteTimeout,
				IdleTimeout:    cfg.IdleTimeout,
				MaxHeaderBytes: cfg.MaxHeaderBytes,
			}

			// Setup TLS if configured
			if cfg.TLSCert != "" && cfg.TLSKey != "" {
				tlsConfig := &tls.Config{
					MinVersion: tls.VersionTLS12,
				}
				server.TLSConfig = tlsConfig

				fmt.Printf("Starting HTTPS server %s on %s\n", n, cfg.Addr)
				if err := server.ListenAndServeTLS(cfg.TLSCert, cfg.TLSKey); err != nil {
					fmt.Printf("Server %s error: %v\n", n, err)
				}
			} else {
				fmt.Printf("Starting HTTP server %s on %s\n", n, cfg.Addr)
				if err := server.ListenAndServe(); err != nil {
					fmt.Printf("Server %s error: %v\n", n, err)
				}
			}
		}(name, serverCfg)
	}

	// Block forever
	select {}
}

// createHandler creates HTTP handler for a server
func (s *Server) createHandler(serverCfg *config.ParsedServer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Match route
		location, matched := s.router.Match(r)
		if !matched {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		// Get upstream
		upstream, err := s.upstreamManager.Get(location.Upstream)
		if err != nil {
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}

		// Create plugin context
		ctx := plugin.NewContext(w, r)

		// Execute request plugins
		if err := s.pluginChain.ExecuteRequest(ctx); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// If response was written by plugin, skip proxying
		if ctx.ResponseWritten {
			return
		}

		// Proxy request
		s.proxyRequest(ctx, upstream, location)
	})
}

// proxyRequest forwards the request to upstream
func (s *Server) proxyRequest(ctx *plugin.Context, upstream *upstream.Upstream, location *config.ParsedLocation) {
	// Select backend server
	server, err := upstream.NextServer(ctx.Request)
	if err != nil {
		ctx.WriteError(http.StatusServiceUnavailable, "No available backend servers")
		return
	}

	// Create reverse proxy
	proxy := &ReverseProxy{
		upstream: upstream,
		server:   server,
		location: location,
	}

	// Forward request
	proxy.ServeHTTP(ctx.Writer, ctx.Request)

	// Execute response plugins
	s.pluginChain.ExecuteResponse(ctx)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}
