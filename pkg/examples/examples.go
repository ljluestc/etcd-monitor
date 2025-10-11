// Package examples provides practical examples of etcd patterns
package examples

import (
	"context"
	"fmt"
	"time"

	"github.com/etcd-monitor/taskmaster/pkg/patterns"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// Example1_DistributedLock demonstrates distributed locking
func Example1_DistributedLock(client *clientv3.Client, logger *zap.Logger) error {
	logger.Info("=== Example 1: Distributed Lock ===")

	// Create a distributed lock
	lock, err := patterns.NewDistributedLock(client, patterns.LockConfig{
		LockKey: "/locks/resource1",
		TTL:     30,
	}, logger)
	if err != nil {
		return err
	}
	defer lock.Close()

	ctx := context.Background()

	// Acquire the lock
	if err := lock.Lock(ctx); err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}

	logger.Info("Lock acquired, performing critical section")

	// Simulate work in critical section
	time.Sleep(2 * time.Second)

	// Release the lock
	if err := lock.Unlock(ctx); err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	logger.Info("Lock released")

	// Example with WithLock helper
	err = lock.WithLock(ctx, func() error {
		logger.Info("Inside WithLock, doing critical work")
		time.Sleep(1 * time.Second)
		return nil
	})

	return err
}

// Example2_ServiceDiscovery demonstrates service registration and discovery
func Example2_ServiceDiscovery(client *clientv3.Client, logger *zap.Logger) error {
	logger.Info("=== Example 2: Service Discovery ===")

	// Create service registry
	registry, err := patterns.NewServiceRegistry(client, "/services", 30, logger)
	if err != nil {
		return err
	}
	defer registry.Close()

	ctx := context.Background()

	// Register a service
	service := patterns.ServiceInfo{
		Name:    "api-server",
		ID:      "api-1",
		Address: "192.168.1.100",
		Port:    8080,
		Metadata: map[string]string{
			"version": "1.0.0",
			"region":  "us-east-1",
		},
	}

	if err := registry.Register(ctx, service); err != nil {
		return err
	}

	logger.Info("Service registered", zap.String("id", service.ID))

	// Discover services
	services, err := registry.Discover(ctx, "api-server")
	if err != nil {
		return err
	}

	logger.Info("Services discovered", zap.Int("count", len(services)))
	for _, svc := range services {
		logger.Info("Found service",
			zap.String("id", svc.ID),
			zap.String("address", fmt.Sprintf("%s:%d", svc.Address, svc.Port)))
	}

	// Watch for service changes
	err = registry.Watch(ctx, "api-server", func(event patterns.ServiceEvent) {
		logger.Info("Service event",
			zap.String("type", string(event.Type)),
			zap.String("service", event.Service.ID))
	})

	return err
}

// Example3_LoadBalancing demonstrates load balancing with service discovery
func Example3_LoadBalancing(client *clientv3.Client, logger *zap.Logger) error {
	logger.Info("=== Example 3: Load Balancing ===")

	registry, err := patterns.NewServiceRegistry(client, "/services", 30, logger)
	if err != nil {
		return err
	}
	defer registry.Close()

	ctx := context.Background()

	// Register multiple service instances
	for i := 1; i <= 3; i++ {
		service := patterns.ServiceInfo{
			Name:    "web-server",
			ID:      fmt.Sprintf("web-%d", i),
			Address: fmt.Sprintf("192.168.1.%d", 100+i),
			Port:    8080,
		}
		if err := registry.Register(ctx, service); err != nil {
			return err
		}
	}

	// Create load balancer
	lb := patterns.NewLoadBalancer(registry, patterns.RoundRobin)

	// Get service instances using round-robin
	for i := 0; i < 5; i++ {
		service, err := lb.GetService(ctx, "web-server")
		if err != nil {
			return err
		}
		logger.Info("Selected service",
			zap.Int("request", i+1),
			zap.String("instance", service.ID),
			zap.String("address", fmt.Sprintf("%s:%d", service.Address, service.Port)))
	}

	return nil
}

// Example4_ConfigurationManagement demonstrates configuration management
func Example4_ConfigurationManagement(client *clientv3.Client, logger *zap.Logger) error {
	logger.Info("=== Example 4: Configuration Management ===")

	// Create config manager
	configMgr := patterns.NewConfigManager(client, "/config/app", logger)

	ctx := context.Background()

	// Set configuration values
	if err := configMgr.Set(ctx, "database/host", "localhost"); err != nil {
		return err
	}
	if err := configMgr.Set(ctx, "database/port", 5432); err != nil {
		return err
	}
	if err := configMgr.Set(ctx, "api/timeout", 30); err != nil {
		return err
	}

	logger.Info("Configuration values set")

	// Get configuration values
	dbHost, err := configMgr.GetString(ctx, "database/host")
	if err != nil {
		return err
	}
	logger.Info("Retrieved config", zap.String("database/host", dbHost))

	dbPort, err := configMgr.GetInt(ctx, "database/port")
	if err != nil {
		return err
	}
	logger.Info("Retrieved config", zap.Int("database/port", dbPort))

	// Watch for configuration changes
	configMgr.Watch(ctx, "api/timeout", func(key string, oldValue, newValue interface{}) {
		logger.Info("Configuration changed",
			zap.String("key", key),
			zap.Any("old", oldValue),
			zap.Any("new", newValue))
	})

	// Get all configurations under a prefix
	allConfigs, err := configMgr.GetAll(ctx, "database")
	if err != nil {
		return err
	}

	logger.Info("All database configs", zap.Int("count", len(allConfigs)))
	for key, value := range allConfigs {
		logger.Info("Config", zap.String("key", key), zap.Any("value", value))
	}

	return nil
}

// Example5_FeatureFlags demonstrates feature flag management
func Example5_FeatureFlags(client *clientv3.Client, logger *zap.Logger) error {
	logger.Info("=== Example 5: Feature Flags ===")

	configMgr := patterns.NewConfigManager(client, "/config/app", logger)
	featureFlags := patterns.NewFeatureFlag(configMgr)

	ctx := context.Background()

	// Enable feature flags
	if err := featureFlags.Enable(ctx, "new-ui"); err != nil {
		return err
	}
	if err := featureFlags.Enable(ctx, "beta-api"); err != nil {
		return err
	}
	if err := featureFlags.Disable(ctx, "old-search"); err != nil {
		return err
	}

	logger.Info("Feature flags configured")

	// Check if features are enabled
	newUIEnabled, err := featureFlags.IsEnabled(ctx, "new-ui")
	if err != nil {
		return err
	}

	if newUIEnabled {
		logger.Info("new-ui feature is enabled")
	} else {
		logger.Info("new-ui feature is disabled")
	}

	// Get all feature flags
	allFlags, err := featureFlags.GetAllFlags(ctx)
	if err != nil {
		return err
	}

	logger.Info("All feature flags", zap.Int("count", len(allFlags)))
	for flag, enabled := range allFlags {
		logger.Info("Feature flag", zap.String("name", flag), zap.Bool("enabled", enabled))
	}

	return nil
}

// Example6_LeaderElection demonstrates leader election
func Example6_LeaderElection(client *clientv3.Client, logger *zap.Logger, nodeID string) error {
	logger.Info("=== Example 6: Leader Election ===", zap.String("node_id", nodeID))

	// Create leader election
	election, err := patterns.NewLeaderElection(client, "/election/primary", nodeID, 30, logger)
	if err != nil {
		return err
	}
	defer election.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Set callbacks
	election.OnElected(func() {
		logger.Info("I am now the leader!")
	})

	election.OnDefeat(func() {
		logger.Info("I lost leadership")
	})

	// Campaign for leadership
	go func() {
		if err := election.Campaign(ctx); err != nil {
			logger.Error("Campaign failed", zap.Error(err))
		}
	}()

	// Observe leadership changes
	observeChan := election.Observe(ctx)

	select {
	case leader := <-observeChan:
		logger.Info("Current leader", zap.String("leader", leader))
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// Example7_LeaderTask demonstrates leader-only tasks
func Example7_LeaderTask(client *clientv3.Client, logger *zap.Logger, nodeID string) error {
	logger.Info("=== Example 7: Leader Task ===", zap.String("node_id", nodeID))

	election, err := patterns.NewLeaderElection(client, "/election/worker", nodeID, 30, logger)
	if err != nil {
		return err
	}
	defer election.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Define a task that should only run on the leader
	leaderTask := patterns.NewLeaderTask(
		election,
		func(ctx context.Context) error {
			logger.Info("Executing leader-only task")
			// Perform periodic cleanup, aggregation, etc.
			return nil
		},
		10*time.Second,
		logger,
	)

	// Campaign for leadership
	go func() {
		if err := election.Campaign(ctx); err != nil {
			logger.Error("Campaign failed", zap.Error(err))
		}
	}()

	// Start the leader task
	leaderTask.Start(ctx)

	// Wait
	<-ctx.Done()
	leaderTask.Stop()

	return nil
}

// Example8_Transaction demonstrates transaction usage
func Example8_Transaction(client *clientv3.Client, logger *zap.Logger) error {
	logger.Info("=== Example 8: Transactions ===")

	ctx := context.Background()

	// Example 1: Compare-and-swap
	key := "/counter"
	value := "0"

	// Set initial value
	_, err := client.Put(ctx, key, value)
	if err != nil {
		return err
	}

	// Transaction: increment counter if current value is 0
	txn := client.Txn(ctx)
	txn.If(
		clientv3.Compare(clientv3.Value(key), "=", "0"),
	).Then(
		clientv3.OpPut(key, "1"),
	).Else(
		clientv3.OpGet(key),
	)

	resp, err := txn.Commit()
	if err != nil {
		return err
	}

	if resp.Succeeded {
		logger.Info("Transaction succeeded, counter incremented")
	} else {
		logger.Info("Transaction failed, counter was modified")
	}

	// Example 2: Multi-key transaction
	txn2 := client.Txn(ctx)
	txn2.If(
		clientv3.Compare(clientv3.Value("/account/a/balance"), ">=", "100"),
	).Then(
		clientv3.OpPut("/account/a/balance", "50"),
		clientv3.OpPut("/account/b/balance", "150"),
		clientv3.OpPut("/transfer/log", "A->B: 50"),
	).Else(
		clientv3.OpGet("/account/a/balance"),
	)

	resp2, err := txn2.Commit()
	if err != nil {
		return err
	}

	if resp2.Succeeded {
		logger.Info("Transfer succeeded")
	} else {
		logger.Info("Transfer failed, insufficient balance")
	}

	return nil
}

// Example9_Leases demonstrates lease management
func Example9_Leases(client *clientv3.Client, logger *zap.Logger) error {
	logger.Info("=== Example 9: Leases ===")

	ctx := context.Background()

	// Create a lease with 30 second TTL
	leaseResp, err := client.Grant(ctx, 30)
	if err != nil {
		return err
	}

	leaseID := leaseResp.ID
	logger.Info("Lease granted", zap.Int64("lease_id", int64(leaseID)), zap.Int64("ttl", leaseResp.TTL))

	// Put a key with the lease
	_, err = client.Put(ctx, "/session/user123", "active", clientv3.WithLease(leaseID))
	if err != nil {
		return err
	}

	logger.Info("Key stored with lease")

	// Keep the lease alive
	keepAliveChan, err := client.KeepAlive(ctx, leaseID)
	if err != nil {
		return err
	}

	// Read keep-alive responses
	go func() {
		for ka := range keepAliveChan {
			logger.Debug("Keep-alive received", zap.Int64("ttl", ka.TTL))
		}
	}()

	// Check lease TTL
	ttlResp, err := client.TimeToLive(ctx, leaseID)
	if err != nil {
		return err
	}

	logger.Info("Lease TTL", zap.Int64("ttl", ttlResp.TTL))

	// Simulate work
	time.Sleep(5 * time.Second)

	// Revoke the lease
	_, err = client.Revoke(ctx, leaseID)
	if err != nil {
		return err
	}

	logger.Info("Lease revoked")

	// Try to get the key (should be gone)
	getResp, err := client.Get(ctx, "/session/user123")
	if err != nil {
		return err
	}

	if len(getResp.Kvs) == 0 {
		logger.Info("Key was deleted with lease")
	}

	return nil
}

// Example10_Watch demonstrates watching for changes
func Example10_Watch(client *clientv3.Client, logger *zap.Logger) error {
	logger.Info("=== Example 10: Watch ===")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Watch a single key
	watchChan := client.Watch(ctx, "/config/setting")

	go func() {
		logger.Info("Starting to watch /config/setting")
		for watchResp := range watchChan {
			if watchResp.Err() != nil {
				logger.Error("Watch error", zap.Error(watchResp.Err()))
				continue
			}

			for _, event := range watchResp.Events {
				logger.Info("Watch event",
					zap.String("type", event.Type.String()),
					zap.String("key", string(event.Kv.Key)),
					zap.String("value", string(event.Kv.Value)))
			}
		}
	}()

	// Make some changes
	time.Sleep(1 * time.Second)
	client.Put(ctx, "/config/setting", "value1")
	time.Sleep(1 * time.Second)
	client.Put(ctx, "/config/setting", "value2")
	time.Sleep(1 * time.Second)
	client.Delete(ctx, "/config/setting")

	// Watch with prefix
	watchChan2 := client.Watch(ctx, "/config/", clientv3.WithPrefix())

	go func() {
		logger.Info("Starting to watch /config/ (prefix)")
		for watchResp := range watchChan2 {
			for _, event := range watchResp.Events {
				logger.Info("Prefix watch event",
					zap.String("key", string(event.Kv.Key)))
			}
		}
	}()

	time.Sleep(2 * time.Second)

	return nil
}

// Example11_HighAvailability demonstrates HA setup
func Example11_HighAvailability(client *clientv3.Client, logger *zap.Logger, nodeID string) error {
	logger.Info("=== Example 11: High Availability ===", zap.String("node_id", nodeID))

	election, err := patterns.NewLeaderElection(client, "/election/ha", nodeID, 30, logger)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Define primary and secondary functions
	primary := func(ctx context.Context) error {
		logger.Info("Running as PRIMARY")
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				logger.Info("PRIMARY: Processing requests")
			}
		}
	}

	secondary := func(ctx context.Context) error {
		logger.Info("Running as SECONDARY (standby)")
		<-ctx.Done()
		return nil
	}

	// Create HA coordinator
	haCoord := patterns.NewHighAvailabilityCoordinator(election, primary, secondary, logger)

	if err := haCoord.Start(ctx); err != nil {
		return err
	}

	// Wait
	<-ctx.Done()
	return haCoord.Stop()
}

// RunAllExamples runs all examples sequentially
func RunAllExamples(client *clientv3.Client, logger *zap.Logger) {
	examples := []struct {
		name string
		fn   func(*clientv3.Client, *zap.Logger) error
	}{
		{"Distributed Lock", Example1_DistributedLock},
		{"Service Discovery", Example2_ServiceDiscovery},
		{"Load Balancing", Example3_LoadBalancing},
		{"Configuration Management", Example4_ConfigurationManagement},
		{"Feature Flags", Example5_FeatureFlags},
		{"Transactions", Example8_Transaction},
		{"Leases", Example9_Leases},
		{"Watch", Example10_Watch},
	}

	for _, example := range examples {
		logger.Info("\n" + example.name + "\n")
		if err := example.fn(client, logger); err != nil {
			logger.Error("Example failed", zap.String("name", example.name), zap.Error(err))
		}
		time.Sleep(2 * time.Second)
	}
}
