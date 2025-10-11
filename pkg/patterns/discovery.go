package patterns

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// ServiceInfo contains information about a registered service
type ServiceInfo struct {
	Name      string            `json:"name"`
	ID        string            `json:"id"`
	Address   string            `json:"address"`
	Port      int               `json:"port"`
	Metadata  map[string]string `json:"metadata"`
	Timestamp time.Time         `json:"timestamp"`
}

// ServiceRegistry provides service registration and discovery
type ServiceRegistry struct {
	client      *clientv3.Client
	keyPrefix   string
	leaseID     clientv3.LeaseID
	ttl         int64
	logger      *zap.Logger
	registered  map[string]*ServiceInfo
	watchers    map[string][]ServiceWatcher
	mu          sync.RWMutex
	cancelWatch context.CancelFunc
}

// ServiceWatcher is called when services change
type ServiceWatcher func(event ServiceEvent)

// ServiceEvent represents a service change event
type ServiceEvent struct {
	Type    ServiceEventType
	Service ServiceInfo
}

// ServiceEventType defines the type of service event
type ServiceEventType string

const (
	ServiceAdded   ServiceEventType = "added"
	ServiceRemoved ServiceEventType = "removed"
	ServiceUpdated ServiceEventType = "updated"
)

// NewServiceRegistry creates a new service registry
func NewServiceRegistry(client *clientv3.Client, keyPrefix string, ttl int64, logger *zap.Logger) (*ServiceRegistry, error) {
	if ttl == 0 {
		ttl = 30 // Default 30 seconds
	}

	return &ServiceRegistry{
		client:     client,
		keyPrefix:  keyPrefix,
		ttl:        ttl,
		logger:     logger,
		registered: make(map[string]*ServiceInfo),
		watchers:   make(map[string][]ServiceWatcher),
	}, nil
}

// Register registers a service with automatic lease renewal
func (sr *ServiceRegistry) Register(ctx context.Context, service ServiceInfo) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	// Create a lease
	leaseResp, err := sr.client.Grant(ctx, sr.ttl)
	if err != nil {
		return fmt.Errorf("failed to create lease: %w", err)
	}

	sr.leaseID = leaseResp.ID

	// Marshal service info
	service.Timestamp = time.Now()
	data, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("failed to marshal service info: %w", err)
	}

	// Put service info with lease
	key := sr.getServiceKey(service.Name, service.ID)
	_, err = sr.client.Put(ctx, key, string(data), clientv3.WithLease(sr.leaseID))
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	// Keep lease alive
	keepAliveChan, err := sr.client.KeepAlive(ctx, sr.leaseID)
	if err != nil {
		return fmt.Errorf("failed to keep lease alive: %w", err)
	}

	// Start goroutine to handle keep-alive responses
	go sr.handleKeepAlive(ctx, keepAliveChan, service)

	sr.registered[service.ID] = &service
	sr.logger.Info("Service registered",
		zap.String("name", service.Name),
		zap.String("id", service.ID),
		zap.String("address", service.Address),
		zap.Int("port", service.Port))

	return nil
}

// Deregister removes a service from the registry
func (sr *ServiceRegistry) Deregister(ctx context.Context, serviceName, serviceID string) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	key := sr.getServiceKey(serviceName, serviceID)
	_, err := sr.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}

	delete(sr.registered, serviceID)
	sr.logger.Info("Service deregistered",
		zap.String("name", serviceName),
		zap.String("id", serviceID))

	return nil
}

// Discover finds all instances of a service
func (sr *ServiceRegistry) Discover(ctx context.Context, serviceName string) ([]ServiceInfo, error) {
	keyPrefix := sr.getServicePrefix(serviceName)

	resp, err := sr.client.Get(ctx, keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to discover services: %w", err)
	}

	services := make([]ServiceInfo, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var service ServiceInfo
		if err := json.Unmarshal(kv.Value, &service); err != nil {
			sr.logger.Warn("Failed to unmarshal service info", zap.Error(err))
			continue
		}
		services = append(services, service)
	}

	sr.logger.Debug("Services discovered",
		zap.String("name", serviceName),
		zap.Int("count", len(services)))

	return services, nil
}

// Watch watches for service changes
func (sr *ServiceRegistry) Watch(ctx context.Context, serviceName string, watcher ServiceWatcher) error {
	sr.mu.Lock()
	sr.watchers[serviceName] = append(sr.watchers[serviceName], watcher)
	sr.mu.Unlock()

	// If this is the first watcher for this service, start watching
	if len(sr.watchers[serviceName]) == 1 {
		go sr.watchService(ctx, serviceName)
	}

	return nil
}

// watchService watches for changes to a specific service
func (sr *ServiceRegistry) watchService(ctx context.Context, serviceName string) {
	keyPrefix := sr.getServicePrefix(serviceName)

	watchCtx, cancel := context.WithCancel(ctx)
	sr.cancelWatch = cancel

	watchChan := sr.client.Watch(watchCtx, keyPrefix, clientv3.WithPrefix())

	for {
		select {
		case <-ctx.Done():
			return
		case watchResp := <-watchChan:
			if watchResp.Err() != nil {
				sr.logger.Error("Watch error", zap.Error(watchResp.Err()))
				continue
			}

			for _, event := range watchResp.Events {
				sr.handleWatchEvent(serviceName, event)
			}
		}
	}
}

// handleWatchEvent processes a watch event
func (sr *ServiceRegistry) handleWatchEvent(serviceName string, event *clientv3.Event) {
	var service ServiceInfo
	if err := json.Unmarshal(event.Kv.Value, &service); err != nil {
		sr.logger.Warn("Failed to unmarshal service in watch event", zap.Error(err))
		return
	}

	var eventType ServiceEventType
	switch event.Type {
	case clientv3.EventTypePut:
		if event.IsCreate() {
			eventType = ServiceAdded
		} else {
			eventType = ServiceUpdated
		}
	case clientv3.EventTypeDelete:
		eventType = ServiceRemoved
	}

	serviceEvent := ServiceEvent{
		Type:    eventType,
		Service: service,
	}

	// Notify all watchers
	sr.mu.RLock()
	watchers := sr.watchers[serviceName]
	sr.mu.RUnlock()

	for _, watcher := range watchers {
		go watcher(serviceEvent)
	}

	sr.logger.Debug("Service event",
		zap.String("type", string(eventType)),
		zap.String("service", service.Name),
		zap.String("id", service.ID))
}

// handleKeepAlive handles lease keep-alive responses
func (sr *ServiceRegistry) handleKeepAlive(ctx context.Context, keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse, service ServiceInfo) {
	for {
		select {
		case <-ctx.Done():
			return
		case resp, ok := <-keepAliveChan:
			if !ok {
				sr.logger.Warn("Keep-alive channel closed",
					zap.String("service", service.Name),
					zap.String("id", service.ID))
				return
			}
			if resp == nil {
				sr.logger.Warn("Keep-alive response is nil")
				return
			}
			sr.logger.Debug("Keep-alive received",
				zap.String("service", service.Name),
				zap.Int64("ttl", resp.TTL))
		}
	}
}

// getServiceKey returns the etcd key for a service
func (sr *ServiceRegistry) getServiceKey(serviceName, serviceID string) string {
	return path.Join(sr.keyPrefix, serviceName, serviceID)
}

// getServicePrefix returns the prefix for all instances of a service
func (sr *ServiceRegistry) getServicePrefix(serviceName string) string {
	return path.Join(sr.keyPrefix, serviceName) + "/"
}

// Close stops the service registry
func (sr *ServiceRegistry) Close() error {
	if sr.cancelWatch != nil {
		sr.cancelWatch()
	}

	// Revoke lease
	if sr.leaseID != 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := sr.client.Revoke(ctx, sr.leaseID)
		if err != nil {
			sr.logger.Error("Failed to revoke lease", zap.Error(err))
			return err
		}
	}

	return nil
}

// LoadBalancer provides load balancing across service instances
type LoadBalancer struct {
	registry *ServiceRegistry
	strategy LoadBalanceStrategy
	mu       sync.RWMutex
	index    int
}

// LoadBalanceStrategy defines load balancing strategy
type LoadBalanceStrategy string

const (
	RoundRobin LoadBalanceStrategy = "round_robin"
	Random     LoadBalanceStrategy = "random"
	LeastConn  LoadBalanceStrategy = "least_conn"
)

// NewLoadBalancer creates a new load balancer
func NewLoadBalancer(registry *ServiceRegistry, strategy LoadBalanceStrategy) *LoadBalancer {
	return &LoadBalancer{
		registry: registry,
		strategy: strategy,
		index:    0,
	}
}

// GetService returns a service instance based on load balancing strategy
func (lb *LoadBalancer) GetService(ctx context.Context, serviceName string) (*ServiceInfo, error) {
	services, err := lb.registry.Discover(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("no instances available for service: %s", serviceName)
	}

	switch lb.strategy {
	case RoundRobin:
		return lb.roundRobin(services), nil
	case Random:
		return lb.random(services), nil
	default:
		return &services[0], nil
	}
}

// roundRobin selects service using round-robin
func (lb *LoadBalancer) roundRobin(services []ServiceInfo) *ServiceInfo {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	index := lb.index % len(services)
	lb.index++

	return &services[index]
}

// random selects service randomly
func (lb *LoadBalancer) random(services []ServiceInfo) *ServiceInfo {
	index := time.Now().UnixNano() % int64(len(services))
	return &services[index]
}

// HealthChecker performs health checks on registered services
type HealthChecker struct {
	registry *ServiceRegistry
	interval time.Duration
	logger   *zap.Logger
	cancel   context.CancelFunc
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(registry *ServiceRegistry, interval time.Duration, logger *zap.Logger) *HealthChecker {
	return &HealthChecker{
		registry: registry,
		interval: interval,
		logger:   logger,
	}
}

// Start begins health checking
func (hc *HealthChecker) Start(ctx context.Context) {
	checkCtx, cancel := context.WithCancel(ctx)
	hc.cancel = cancel

	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-checkCtx.Done():
			return
		case <-ticker.C:
			hc.checkServices(checkCtx)
		}
	}
}

// Stop stops health checking
func (hc *HealthChecker) Stop() {
	if hc.cancel != nil {
		hc.cancel()
	}
}

// checkServices checks all registered services
func (hc *HealthChecker) checkServices(ctx context.Context) {
	// This would perform actual health checks
	// For now, it's a placeholder
	hc.logger.Debug("Performing health checks")
}
