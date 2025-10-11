package config

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/BurntSushi/toml"
)

// Loader handles configuration loading and parsing
type Loader struct {
	filePath string
}

// NewLoader creates a new configuration loader
func NewLoader(filePath string) *Loader {
	return &Loader{
		filePath: filePath,
	}
}

// Load reads and parses the configuration file
func (l *Loader) Load() (*ParsedConfig, error) {
	data, err := os.ReadFile(l.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse TOML: %w", err)
	}

	// Parse and validate
	parsed, err := l.parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := l.validate(parsed); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return parsed, nil
}

// parse converts string durations to time.Duration and compiles regex patterns
func (l *Loader) parse(cfg *Config) (*ParsedConfig, error) {
	parsed := &ParsedConfig{
		Config:          cfg,
		ParsedUpstreams: make(map[string]*ParsedUpstream),
		ParsedServers:   make(map[string]*ParsedServer),
	}

	// Parse grace period
	if cfg.Basic.GracePeriod != "" {
		d, err := time.ParseDuration(cfg.Basic.GracePeriod)
		if err != nil {
			return nil, fmt.Errorf("invalid grace_period: %w", err)
		}
		parsed.GracePeriod = d
	} else {
		parsed.GracePeriod = 30 * time.Second
	}

	// Parse upstreams
	for name, upstream := range cfg.Upstreams {
		pu, err := l.parseUpstream(&upstream)
		if err != nil {
			return nil, fmt.Errorf("upstream %s: %w", name, err)
		}
		parsed.ParsedUpstreams[name] = pu
	}

	// Parse locations
	for name, location := range cfg.Locations {
		pl, err := l.parseLocation(name, &location)
		if err != nil {
			return nil, fmt.Errorf("location %s: %w", name, err)
		}
		parsed.ParsedLocations = append(parsed.ParsedLocations, pl)
	}

	// Parse servers
	for name, server := range cfg.Servers {
		ps, err := l.parseServer(name, &server)
		if err != nil {
			return nil, fmt.Errorf("server %s: %w", name, err)
		}
		parsed.ParsedServers[name] = ps
	}

	return parsed, nil
}

func (l *Loader) parseUpstream(u *Upstream) (*ParsedUpstream, error) {
	pu := &ParsedUpstream{Upstream: u}

	var err error
	if u.HealthCheckInterval != "" {
		pu.HealthCheckInterval, err = time.ParseDuration(u.HealthCheckInterval)
		if err != nil {
			return nil, fmt.Errorf("invalid health_check_interval: %w", err)
		}
	} else {
		pu.HealthCheckInterval = 10 * time.Second
	}

	if u.ConnectionTimeout != "" {
		pu.ConnectionTimeout, err = time.ParseDuration(u.ConnectionTimeout)
		if err != nil {
			return nil, fmt.Errorf("invalid connection_timeout: %w", err)
		}
	} else {
		pu.ConnectionTimeout = 10 * time.Second
	}

	if u.ReadTimeout != "" {
		pu.ReadTimeout, err = time.ParseDuration(u.ReadTimeout)
		if err != nil {
			return nil, fmt.Errorf("invalid read_timeout: %w", err)
		}
	} else {
		pu.ReadTimeout = 30 * time.Second
	}

	if u.WriteTimeout != "" {
		pu.WriteTimeout, err = time.ParseDuration(u.WriteTimeout)
		if err != nil {
			return nil, fmt.Errorf("invalid write_timeout: %w", err)
		}
	} else {
		pu.WriteTimeout = 30 * time.Second
	}

	if u.IdleTimeout != "" {
		pu.IdleTimeout, err = time.ParseDuration(u.IdleTimeout)
		if err != nil {
			return nil, fmt.Errorf("invalid idle_timeout: %w", err)
		}
	} else {
		pu.IdleTimeout = 90 * time.Second
	}

	return pu, nil
}

func (l *Loader) parseLocation(name string, loc *Location) (ParsedLocation, error) {
	pl := ParsedLocation{
		Location: loc,
		Name:     name,
	}

	// Compile path pattern if regex
	if loc.PathType == "regex" && loc.Path != "" {
		re, err := regexp.Compile(loc.Path)
		if err != nil {
			return pl, fmt.Errorf("invalid path regex: %w", err)
		}
		pl.PathPattern = re
	}

	return pl, nil
}

func (l *Loader) parseServer(name string, srv *Server) (*ParsedServer, error) {
	ps := &ParsedServer{
		Server: srv,
		Name:   name,
	}

	var err error
	if srv.ReadTimeout != "" {
		ps.ReadTimeout, err = time.ParseDuration(srv.ReadTimeout)
		if err != nil {
			return nil, fmt.Errorf("invalid read_timeout: %w", err)
		}
	} else {
		ps.ReadTimeout = 30 * time.Second
	}

	if srv.WriteTimeout != "" {
		ps.WriteTimeout, err = time.ParseDuration(srv.WriteTimeout)
		if err != nil {
			return nil, fmt.Errorf("invalid write_timeout: %w", err)
		}
	} else {
		ps.WriteTimeout = 30 * time.Second
	}

	if srv.IdleTimeout != "" {
		ps.IdleTimeout, err = time.ParseDuration(srv.IdleTimeout)
		if err != nil {
			return nil, fmt.Errorf("invalid idle_timeout: %w", err)
		}
	} else {
		ps.IdleTimeout = 90 * time.Second
	}

	return ps, nil
}

// validate checks configuration for errors
func (l *Loader) validate(cfg *ParsedConfig) error {
	// Validate upstreams
	for name, upstream := range cfg.ParsedUpstreams {
		if len(upstream.Addrs) == 0 {
			return fmt.Errorf("upstream %s has no addresses", name)
		}

		algo := upstream.Algorithm
		if algo == "" {
			algo = "round_robin"
		}
		if algo != "round_robin" && algo != "least_conn" && algo != "ip_hash" && algo != "weighted_round_robin" {
			return fmt.Errorf("upstream %s: invalid algorithm %s", name, algo)
		}

		if algo == "weighted_round_robin" {
			if len(upstream.Weights) != len(upstream.Addrs) {
				return fmt.Errorf("upstream %s: weights count must match addrs count", name)
			}
		}
	}

	// Validate locations
	for _, loc := range cfg.ParsedLocations {
		if loc.Upstream == "" {
			return fmt.Errorf("location %s: upstream is required", loc.Name)
		}
		if _, ok := cfg.ParsedUpstreams[loc.Upstream]; !ok {
			return fmt.Errorf("location %s: upstream %s not found", loc.Name, loc.Upstream)
		}

		pathType := loc.PathType
		if pathType == "" {
			pathType = "prefix"
		}
		if pathType != "exact" && pathType != "prefix" && pathType != "regex" {
			return fmt.Errorf("location %s: invalid path_type %s", loc.Name, pathType)
		}
	}

	// Validate servers
	for name, server := range cfg.ParsedServers {
		if server.Addr == "" {
			return fmt.Errorf("server %s: addr is required", name)
		}
		if len(server.Locations) == 0 {
			return fmt.Errorf("server %s: at least one location is required", name)
		}

		// Check all referenced locations exist
		for _, locName := range server.Locations {
			found := false
			for _, loc := range cfg.ParsedLocations {
				if loc.Name == locName {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("server %s: location %s not found", name, locName)
			}
		}

		// Validate TLS configuration
		if server.TLSCert != "" || server.TLSKey != "" {
			if server.TLSCert == "" || server.TLSKey == "" {
				return fmt.Errorf("server %s: both tls_cert and tls_key must be specified", name)
			}
		}
	}

	return nil
}

// LoadFromBytes loads configuration from byte slice (for testing/etcd)
func LoadFromBytes(data []byte) (*ParsedConfig, error) {
	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse TOML: %w", err)
	}

	loader := &Loader{}
	parsed, err := loader.parse(&cfg)
	if err != nil {
		return nil, err
	}

	if err := loader.validate(parsed); err != nil {
		return nil, err
	}

	return parsed, nil
}
