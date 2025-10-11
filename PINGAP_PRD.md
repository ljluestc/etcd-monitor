# Pingap-Go: High-Performance Reverse Proxy

## Executive Summary

**Project:** pingap-go
**Version:** 1.0.0
**Based on:** Pingap (Rust) architecture, implemented in Go
**Purpose:** Production-grade reverse proxy with advanced routing, load balancing, and observability

## Core Features (Must-Have - P0)

### 1. Reverse Proxy Engine
- HTTP/1.1, HTTP/2, h2c support
- gRPC and gRPC-web proxying
- WebSocket support
- Efficient connection pooling
- Request/response buffering
- Streaming support

### 2. Upstream Management
- Multiple upstream servers per service
- Health checking (active/passive)
- Load balancing algorithms:
  - Round Robin
  - Least Connections
  - IP Hash
  - Weighted Round Robin
- Dynamic upstream discovery
- Circuit breaker pattern

### 3. Routing System
- Host-based routing
- Path-based routing (exact, prefix, regex)
- Header-based routing
- Query parameter routing
- Priority-based rule matching
- Path rewriting
- Request/response header manipulation

### 4. Configuration Management
- TOML configuration format
- Hot reload (zero downtime)
- File-based configuration
- etcd configuration backend
- Configuration validation
- Version control and rollback
- Configuration history tracking

### 5. Plugin System
- Request/response lifecycle hooks
- Chain-able plugin execution
- Built-in plugins:
  - **Authentication:** JWT, API Key, Basic Auth
  - **Security:** CSRF, IP Restriction, Referer Check, User-Agent Filter
  - **Traffic Control:** Rate Limiting, Request Size Limit
  - **Caching:** Response caching with TTL
  - **Compression:** Gzip, Brotli
  - **Redirect:** HTTP to HTTPS, domain redirect
  - **Request ID:** X-Request-ID injection
  - **CORS:** Cross-Origin Resource Sharing
  - **Rewrite:** URL rewriting
  - **Mock:** Mock responses for testing

### 6. TLS/SSL
- TLS 1.2, TLS 1.3 support
- Automatic certificate management (Let's Encrypt)
- SNI (Server Name Indication)
- Certificate hot reload
- Custom certificate configuration
- OCSP stapling

### 7. Observability
- Prometheus metrics endpoint
- OpenTelemetry tracing
- Structured logging
- Access logs with 30+ variables:
  - Request/response info
  - Timing metrics
  - Upstream details
  - Error tracking
- Custom log format
- Log rotation

### 8. Web UI
- Dashboard with real-time metrics
- Configuration editor
- Service/upstream management
- Plugin configuration
- Log viewer
- Certificate management
- Health status visualization

## Architecture

### Component Structure

```
pingap-go/
├── cmd/
│   └── pingap/
│       └── main.go              # Entry point
├── pkg/
│   ├── config/
│   │   ├── types.go             # Config data structures
│   │   ├── loader.go            # Config loading (file/etcd)
│   │   ├── validator.go         # Config validation
│   │   └── watcher.go           # Hot reload
│   ├── proxy/
│   │   ├── server.go            # HTTP server
│   │   ├── handler.go           # Request handler
│   │   ├── reverseproxy.go      # Core proxy logic
│   │   └── transport.go         # HTTP transport pool
│   ├── upstream/
│   │   ├── manager.go           # Upstream management
│   │   ├── loadbalancer.go      # Load balancing
│   │   ├── healthcheck.go       # Health checking
│   │   └── discovery.go         # Service discovery
│   ├── router/
│   │   ├── router.go            # Route matching
│   │   ├── matcher.go           # Match rules
│   │   └── rewrite.go           # Path rewriting
│   ├── plugin/
│   │   ├── plugin.go            # Plugin interface
│   │   ├── chain.go             # Plugin chain
│   │   ├── auth/                # Auth plugins
│   │   ├── security/            # Security plugins
│   │   ├── traffic/             # Traffic control
│   │   └── transform/           # Content transformation
│   ├── tls/
│   │   ├── manager.go           # Certificate management
│   │   └── acme.go              # Let's Encrypt
│   ├── metrics/
│   │   ├── prometheus.go        # Prometheus exporter
│   │   └── collector.go         # Metrics collection
│   ├── tracing/
│   │   └── otlp.go              # OpenTelemetry
│   ├── logging/
│   │   └── logger.go            # Structured logging
│   └── admin/
│       ├── api.go               # Admin API
│       └── ui/                  # Web UI assets
├── web/                         # React Web UI
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
├── configs/
│   └── example.toml             # Example config
├── docker/
│   ├── Dockerfile
│   └── docker-compose.yml
└── README.md
```

## Configuration Example

```toml
[basic]
name = "pingap-go"
error_template = "error.html"
pid_file = "/var/run/pingap.pid"
upgrade_sock = "/tmp/pingap_upgrade.sock"
user = "nobody"
group = "nobody"
threads = 4

[upstreams.backend_api]
addrs = [
  "127.0.0.1:8001",
  "127.0.0.1:8002",
  "127.0.0.1:8003"
]
sni = "api.example.com"
algorithm = "round_robin"
health_check = "http://127.0.0.1:8001/health"
connection_timeout = "10s"
read_timeout = "30s"
idle_timeout = "90s"

[locations.api_v1]
upstream = "backend_api"
path = "/api/v1"
plugins = ["jwt_auth", "rate_limit", "request_id"]
rewrite = "/"
headers_add = [
  "X-Forwarded-Proto:https"
]
headers_remove = ["Server"]

[servers.main]
addr = "0.0.0.0:80"
locations = ["api_v1"]
access_log = "/var/log/pingap/access.log"

[servers.https]
addr = "0.0.0.0:443"
tls_cert = "/etc/pingap/certs/cert.pem"
tls_key = "/etc/pingap/certs/key.pem"
locations = ["api_v1"]

[plugins.jwt_auth]
type = "jwt"
secret = "your-secret-key"
algorithm = "HS256"
header = "Authorization"

[plugins.rate_limit]
type = "rate_limit"
rate = 100
burst = 200
key = "$remote_addr"

[plugins.request_id]
type = "request_id"
header = "X-Request-ID"

[metrics]
enabled = true
path = "/metrics"
addr = "127.0.0.1:9090"

[tracing]
enabled = true
endpoint = "localhost:4317"
service_name = "pingap-go"

[admin]
enabled = true
addr = "127.0.0.1:3000"
user = "admin"
password = "changeme"
```

## Plugin Interface

```go
type Plugin interface {
    Name() string
    Priority() int
    OnRequest(ctx *Context) error
    OnResponse(ctx *Context) error
}

type Context struct {
    Request     *http.Request
    Response    *http.Response
    Upstream    *Upstream
    StartTime   time.Time
    Variables   map[string]interface{}
    Logger      *Logger
}
```

## Performance Targets

- **Throughput:** >100k requests/sec on commodity hardware
- **Latency:** <1ms proxy overhead (p99)
- **Memory:** <100MB baseline, efficient under load
- **CPU:** <5% at 10k req/sec
- **Connections:** Support 10k+ concurrent connections
- **Config Reload:** <100ms

## Implementation Phases

### Phase 1: Core Proxy (Week 1-2)
- Basic reverse proxy engine
- HTTP/1.1 and HTTP/2 support
- Simple upstream management
- Round-robin load balancing
- TOML configuration loading
- Basic logging

### Phase 2: Advanced Routing (Week 3)
- Host/path/header routing
- Route matching and priority
- Path rewriting
- Header manipulation
- Request/response buffering

### Phase 3: Plugin System (Week 4-5)
- Plugin architecture
- Plugin chain execution
- Core plugins:
  - JWT authentication
  - API key authentication
  - Rate limiting
  - IP restriction
  - Request ID
  - CORS

### Phase 4: Observability (Week 6)
- Prometheus metrics
- Access logging
- OpenTelemetry tracing
- Performance profiling

### Phase 5: Advanced Features (Week 7-8)
- TLS/SSL management
- Auto certificates (Let's Encrypt)
- Health checking
- Circuit breaker
- Caching
- Compression

### Phase 6: Configuration & Management (Week 9)
- Hot reload
- etcd backend
- Configuration validation
- Version control
- Admin API

### Phase 7: Web UI (Week 10)
- React dashboard
- Configuration editor
- Metrics visualization
- Log viewer
- Certificate management

### Phase 8: Production Ready (Week 11-12)
- Comprehensive testing
- Performance optimization
- Documentation
- Docker packaging
- Examples and tutorials

## Success Criteria

1. **Functionality:** All P0 features implemented and working
2. **Performance:** Meets performance targets
3. **Reliability:** 99.99% uptime in testing
4. **Security:** Pass security audit
5. **Usability:** Intuitive configuration and management
6. **Documentation:** Complete API docs and user guide

## Technology Stack

- **Language:** Go 1.21+
- **Web Framework:** net/http (stdlib)
- **Configuration:** BurntSushi/toml
- **Metrics:** prometheus/client_golang
- **Tracing:** go.opentelemetry.io
- **TLS:** crypto/acme (Let's Encrypt)
- **Storage:** etcd.io/client/v3
- **Web UI:** React + TypeScript + Vite
- **Testing:** testify, httptest
- **Containerization:** Docker

## Dependencies

```go
require (
    github.com/BurntSushi/toml v1.3.2
    github.com/prometheus/client_golang v1.18.0
    go.opentelemetry.io/otel v1.22.0
    go.opentelemetry.io/otel/exporters/otlp v0.46.0
    go.etcd.io/etcd/client/v3 v3.5.11
    golang.org/x/crypto v0.18.0
    golang.org/x/net v0.20.0
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/gorilla/mux v1.8.1
    go.uber.org/zap v1.26.0
)
```
