package config

import (
	"time"
)

// Config represents the complete configuration
type Config struct {
	Basic     Basic                `toml:"basic"`
	Upstreams map[string]Upstream  `toml:"upstreams"`
	Locations map[string]Location  `toml:"locations"`
	Servers   map[string]Server    `toml:"servers"`
	Plugins   map[string]Plugin    `toml:"plugins"`
	Metrics   Metrics              `toml:"metrics"`
	Tracing   Tracing              `toml:"tracing"`
	Admin     Admin                `toml:"admin"`
}

// Basic contains basic server configuration
type Basic struct {
	Name          string `toml:"name"`
	ErrorTemplate string `toml:"error_template"`
	PidFile       string `toml:"pid_file"`
	UpgradeSock   string `toml:"upgrade_sock"`
	User          string `toml:"user"`
	Group         string `toml:"group"`
	Threads       int    `toml:"threads"`
	GracePeriod   string `toml:"grace_period"`
}

// Upstream represents a backend service
type Upstream struct {
	Addrs              []string      `toml:"addrs"`
	SNI                string        `toml:"sni"`
	Algorithm          string        `toml:"algorithm"` // round_robin, least_conn, ip_hash, weighted_round_robin
	HealthCheck        string        `toml:"health_check"`
	HealthCheckInterval string       `toml:"health_check_interval"`
	ConnectionTimeout  string        `toml:"connection_timeout"`
	ReadTimeout        string        `toml:"read_timeout"`
	WriteTimeout       string        `toml:"write_timeout"`
	IdleTimeout        string        `toml:"idle_timeout"`
	MaxIdleConns       int           `toml:"max_idle_conns"`
	MaxIdleConnsPerHost int          `toml:"max_idle_conns_per_host"`
	TLSVerify          bool          `toml:"tls_verify"`
	TLSCert            string        `toml:"tls_cert"`
	TLSKey             string        `toml:"tls_key"`
	Weights            []int         `toml:"weights"`
}

// Location represents a routing location
type Location struct {
	Upstream      string            `toml:"upstream"`
	Path          string            `toml:"path"`
	PathType      string            `toml:"path_type"` // exact, prefix, regex
	Host          string            `toml:"host"`
	Plugins       []string          `toml:"plugins"`
	Rewrite       string            `toml:"rewrite"`
	ProxySetHeaders map[string]string `toml:"proxy_set_headers"`
	HeadersAdd    []string          `toml:"headers_add"`
	HeadersRemove []string          `toml:"headers_remove"`
	HeadersSet    map[string]string `toml:"headers_set"`
	Priority      int               `toml:"priority"`
	Methods       []string          `toml:"methods"`
}

// Server represents an HTTP server
type Server struct {
	Addr           string   `toml:"addr"`
	Locations      []string `toml:"locations"`
	TLSCert        string   `toml:"tls_cert"`
	TLSKey         string   `toml:"tls_key"`
	TLSAuto        bool     `toml:"tls_auto"`
	TLSDomains     []string `toml:"tls_domains"`
	AccessLog      string   `toml:"access_log"`
	AccessLogFormat string  `toml:"access_log_format"`
	ErrorLog       string   `toml:"error_log"`
	EnableHTTP2    bool     `toml:"enable_http2"`
	EnableWebSocket bool    `toml:"enable_websocket"`
	ReadTimeout    string   `toml:"read_timeout"`
	WriteTimeout   string   `toml:"write_timeout"`
	IdleTimeout    string   `toml:"idle_timeout"`
	MaxHeaderBytes int      `toml:"max_header_bytes"`
}

// Plugin represents a plugin configuration
type Plugin struct {
	Type     string                 `toml:"type"`
	Enabled  bool                   `toml:"enabled"`
	Priority int                    `toml:"priority"`
	Config   map[string]interface{} `toml:"config"`
}

// Metrics configuration
type Metrics struct {
	Enabled bool   `toml:"enabled"`
	Path    string `toml:"path"`
	Addr    string `toml:"addr"`
}

// Tracing configuration
type Tracing struct {
	Enabled     bool   `toml:"enabled"`
	Endpoint    string `toml:"endpoint"`
	ServiceName string `toml:"service_name"`
	SampleRate  float64 `toml:"sample_rate"`
}

// Admin configuration
type Admin struct {
	Enabled  bool   `toml:"enabled"`
	Addr     string `toml:"addr"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	TLSCert  string `toml:"tls_cert"`
	TLSKey   string `toml:"tls_key"`
}

// ParsedConfig contains parsed durations and validated config
type ParsedConfig struct {
	*Config
	ParsedUpstreams  map[string]*ParsedUpstream
	ParsedLocations  []ParsedLocation
	ParsedServers    map[string]*ParsedServer
	GracePeriod      time.Duration
}

// ParsedUpstream with parsed durations
type ParsedUpstream struct {
	*Upstream
	HealthCheckInterval time.Duration
	ConnectionTimeout   time.Duration
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	IdleTimeout         time.Duration
}

// ParsedLocation with compiled patterns
type ParsedLocation struct {
	*Location
	Name         string
	PathPattern  interface{} // *regexp.Regexp for regex type
	HostPattern  interface{} // *regexp.Regexp for regex matching
}

// ParsedServer with parsed durations
type ParsedServer struct {
	*Server
	Name         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}
