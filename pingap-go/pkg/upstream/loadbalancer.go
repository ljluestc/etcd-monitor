package upstream

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"net/http"
	"sync/atomic"
)

// LoadBalancer interface for different load balancing strategies
type LoadBalancer interface {
	Next(r *http.Request) (*Server, error)
}

// RoundRobinLB implements round-robin load balancing
type RoundRobinLB struct {
	servers []*Server
	current uint32
}

func NewRoundRobinLB(servers []*Server) *RoundRobinLB {
	return &RoundRobinLB{servers: servers}
}

func (lb *RoundRobinLB) Next(r *http.Request) (*Server, error) {
	n := atomic.AddUint32(&lb.current, 1)

	// Find next healthy server
	for i := 0; i < len(lb.servers); i++ {
		idx := int(n+uint32(i)) % len(lb.servers)
		if lb.servers[idx].IsHealthy() {
			return lb.servers[idx], nil
		}
	}

	// No healthy servers
	if len(lb.servers) > 0 {
		return lb.servers[0], nil
	}
	return nil, fmt.Errorf("no servers available")
}

// WeightedRoundRobinLB implements weighted round-robin
type WeightedRoundRobinLB struct {
	servers []*Server
	current uint32
}

func NewWeightedRoundRobinLB(servers []*Server) *WeightedRoundRobinLB {
	return &WeightedRoundRobinLB{servers: servers}
}

func (lb *WeightedRoundRobinLB) Next(r *http.Request) (*Server, error) {
	n := atomic.AddUint32(&lb.current, 1)

	// Calculate total weight of healthy servers
	totalWeight := 0
	for _, srv := range lb.servers {
		if srv.IsHealthy() {
			totalWeight += srv.Weight
		}
	}

	if totalWeight == 0 {
		if len(lb.servers) > 0 {
			return lb.servers[0], nil
		}
		return nil, fmt.Errorf("no servers available")
	}

	// Select server based on weight
	pos := int(n) % totalWeight
	for _, srv := range lb.servers {
		if !srv.IsHealthy() {
			continue
		}
		if pos < srv.Weight {
			return srv, nil
		}
		pos -= srv.Weight
	}

	return lb.servers[0], nil
}

// IPHashLB implements IP hash load balancing
type IPHashLB struct {
	servers []*Server
}

func NewIPHashLB(servers []*Server) *IPHashLB {
	return &IPHashLB{servers: servers}
}

func (lb *IPHashLB) Next(r *http.Request) (*Server, error) {
	// Get client IP
	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = forwarded
	}

	// Hash the IP
	hash := md5.Sum([]byte(ip))
	idx := binary.BigEndian.Uint32(hash[:4]) % uint32(len(lb.servers))

	// Find healthy server starting from hashed index
	for i := 0; i < len(lb.servers); i++ {
		serverIdx := (int(idx) + i) % len(lb.servers)
		if lb.servers[serverIdx].IsHealthy() {
			return lb.servers[serverIdx], nil
		}
	}

	// No healthy servers
	if len(lb.servers) > 0 {
		return lb.servers[idx], nil
	}
	return nil, fmt.Errorf("no servers available")
}

// LeastConnLB implements least connections load balancing
type LeastConnLB struct {
	servers     []*Server
	connections []uint32
}

func NewLeastConnLB(servers []*Server) *LeastConnLB {
	return &LeastConnLB{
		servers:     servers,
		connections: make([]uint32, len(servers)),
	}
}

func (lb *LeastConnLB) Next(r *http.Request) (*Server, error) {
	minConns := ^uint32(0)
	minIdx := 0

	for i, srv := range lb.servers {
		if !srv.IsHealthy() {
			continue
		}
		conns := atomic.LoadUint32(&lb.connections[i])
		if conns < minConns {
			minConns = conns
			minIdx = i
		}
	}

	atomic.AddUint32(&lb.connections[minIdx], 1)

	// Decrease connection count when done (caller should call Release)
	return lb.servers[minIdx], nil
}

func (lb *LeastConnLB) Release(server *Server) {
	for i, srv := range lb.servers {
		if srv == server {
			atomic.AddUint32(&lb.connections[i], ^uint32(0)) // decrement
			break
		}
	}
}
