package proxy

import (
	"io"
	"net/http"
	"strings"

	"github.com/etcd-monitor/pingap-go/pkg/config"
	"github.com/etcd-monitor/pingap-go/pkg/upstream"
)

// ReverseProxy handles proxying requests to upstream
type ReverseProxy struct {
	upstream *upstream.Upstream
	server   *upstream.Server
	location *config.ParsedLocation
	transport *http.Transport
}

// ServeHTTP implements http.Handler
func (p *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create transport if needed
	if p.transport == nil {
		p.transport = p.upstream.CreateTransport()
	}

	// Clone request
	outReq := r.Clone(r.Context())

	// Set target URL
	targetURL := p.server.URL()
	outReq.URL.Scheme = targetURL.Scheme
	outReq.URL.Host = targetURL.Host

	// Apply rewrite if configured
	if p.location.Rewrite != "" {
		outReq.URL.Path = strings.Replace(outReq.URL.Path, p.location.Path, p.location.Rewrite, 1)
	}

	// Set headers
	outReq.RequestURI = "" // Must be cleared for client requests
	outReq.Host = targetURL.Host

	// Apply header modifications
	p.modifyHeaders(outReq)

	// Execute request
	resp, err := p.transport.RoundTrip(outReq)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Apply response header modifications
	for _, hdr := range p.location.HeadersRemove {
		w.Header().Del(hdr)
	}

	for key, value := range p.location.HeadersSet {
		w.Header().Set(key, value)
	}

	for _, hdr := range p.location.HeadersAdd {
		parts := strings.SplitN(hdr, ":", 2)
		if len(parts) == 2 {
			w.Header().Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	// Write status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	io.Copy(w, resp.Body)
}

// modifyHeaders modifies request headers
func (p *ReverseProxy) modifyHeaders(r *http.Request) {
	// Set X-Forwarded headers
	if clientIP := r.RemoteAddr; clientIP != "" {
		r.Header.Set("X-Forwarded-For", clientIP)
	}
	r.Header.Set("X-Forwarded-Proto", "http")
	if r.TLS != nil {
		r.Header.Set("X-Forwarded-Proto", "https")
	}

	// Apply custom headers from proxy_set_headers
	for key, value := range p.location.ProxySetHeaders {
		r.Header.Set(key, value)
	}
}
