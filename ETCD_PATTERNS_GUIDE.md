# Complete Guide to etcd Patterns - Full Implementation

This guide provides comprehensive implementations of all key etcd patterns based on the complete guide to etcd as a distributed key-value store powering cloud infrastructure.

## Table of Contents

1. [Distributed Locking](#1-distributed-locking)
2. [Service Discovery](#2-service-discovery)
3. [Configuration Management](#3-configuration-management)
4. [Leader Election](#4-leader-election)
5. [Transactions](#5-transactions)
6. [Leases & TTL](#6-leases--ttl)
7. [Watches](#7-watches)
8. [High Availability](#8-high-availability)

---

## 1. Distributed Locking

**Location:** `pkg/patterns/lock.go`

### Purpose
Provide mutual exclusion across distributed systems to ensure only one process accesses a critical resource at a time.

### Implementation

#### Basic Distributed Lock

```go
lock, err := patterns.NewDistributedLock(client, patterns.LockConfig{
    LockKey: "/locks/resource1",
    TTL:     30,  // seconds
}, logger)
defer lock.Close()

// Acquire lock
if err := lock.Lock(ctx); err != nil {
    return err
}

// Critical section
// ... do work ...

// Release lock
if err := lock.Unlock(ctx); err != nil {
    return err
}
```

#### Using WithLock Helper

```go
err := lock.WithLock(ctx, func() error {
    // Critical section
    fmt.Println("Inside critical section")
    return nil
})
```

#### Try Lock (Non-blocking)

```go
acquired, err := lock.TryLock(ctx)
if acquired {
    defer lock.Unlock(ctx)
    // Got the lock
} else {
    // Lock not available
}
```

### Features

- **Automatic TTL**: Lock expires after TTL to prevent deadlocks
- **Session-based**: Uses etcd sessions for automatic cleanup
- **Non-blocking option**: TryLock for immediate return
- **Lock Manager**: Manage multiple locks
- **Read-Write Locks**: Support for reader-writer patterns

### Read-Write Lock

```go
rwLock, err := patterns.NewReadWriteLock(client, "/locks/data", 30, logger)

// Multiple readers
rwLock.RLock(ctx)
// ... read data ...
rwLock.RUnlock(ctx)

// Single writer
rwLock.Lock(ctx)
// ... write data ...
rwLock.Unlock(ctx)
```

---

## 2. Service Discovery

**Location:** `pkg/patterns/discovery.go`

### Purpose
Enable dynamic service registration and discovery for microservices architecture.

### Implementation

#### Register a Service

```go
registry, err := patterns.NewServiceRegistry(client, "/services", 30, logger)
defer registry.Close()

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

err = registry.Register(ctx, service)
```

#### Discover Services

```go
services, err := registry.Discover(ctx, "api-server")
for _, svc := range services {
    fmt.Printf("Found: %s at %s:%d\n", svc.ID, svc.Address, svc.Port)
}
```

#### Watch for Service Changes

```go
registry.Watch(ctx, "api-server", func(event patterns.ServiceEvent) {
    switch event.Type {
    case patterns.ServiceAdded:
        fmt.Printf("Service added: %s\n", event.Service.ID)
    case patterns.ServiceRemoved:
        fmt.Printf("Service removed: %s\n", event.Service.ID)
    case patterns.ServiceUpdated:
        fmt.Printf("Service updated: %s\n", event.Service.ID)
    }
})
```

### Features

- **Automatic heartbeat**: Services send periodic heartbeats
- **TTL-based**: Services expire if heartbeat stops
- **Event notifications**: Watch for service changes
- **Metadata support**: Store arbitrary service metadata
- **Health checking**: Optional health check integration

### Load Balancing

```go
lb := patterns.NewLoadBalancer(registry, patterns.RoundRobin)

// Get service instance
service, err := lb.GetService(ctx, "api-server")
fmt.Printf("Route to: %s:%d\n", service.Address, service.Port)
```

**Strategies:**
- `RoundRobin`: Distribute requests evenly
- `Random`: Random selection
- `LeastConn`: Least connections (planned)

---

## 3. Configuration Management

**Location:** `pkg/patterns/config.go`

### Purpose
Centralized, dynamic configuration management with change notifications.

### Implementation

#### Store Configuration

```go
configMgr := patterns.NewConfigManager(client, "/config/app", logger)

// Set values
configMgr.Set(ctx, "database/host", "localhost")
configMgr.Set(ctx, "database/port", 5432)
configMgr.Set(ctx, "api/timeout", 30)
```

#### Retrieve Configuration

```go
// Get string
dbHost, err := configMgr.GetString(ctx, "database/host")

// Get integer
dbPort, err := configMgr.GetInt(ctx, "database/port")

// Get boolean
enabled, err := configMgr.GetBool(ctx, "feature/enabled")

// Get all configs under prefix
allConfigs, err := configMgr.GetAll(ctx, "database")
```

#### Watch for Changes

```go
configMgr.Watch(ctx, "api/timeout", func(key string, oldValue, newValue interface{}) {
    fmt.Printf("Config changed: %s from %v to %v\n", key, oldValue, newValue)
})
```

### Features

- **Type-safe getters**: GetString, GetInt, GetBool
- **Hierarchical keys**: Organize configs in tree structure
- **Change notifications**: React to configuration updates
- **Caching**: In-memory cache for performance
- **Batch operations**: Get all configs under prefix

### Feature Flags

```go
featureFlags := patterns.NewFeatureFlag(configMgr)

// Enable feature
featureFlags.Enable(ctx, "new-ui")

// Check if enabled
if enabled, _ := featureFlags.IsEnabled(ctx, "new-ui"); enabled {
    // Use new UI
}

// Get all flags
allFlags, err := featureFlags.GetAllFlags(ctx)
```

### Configuration Snapshots

```go
snapshotMgr := patterns.NewSnapshotManager(configMgr)

// Create snapshot
snapshot, err := snapshotMgr.CreateSnapshot(ctx, 1)

// Restore snapshot
err = snapshotMgr.RestoreSnapshot(ctx, 1)
```

---

## 4. Leader Election

**Location:** `pkg/patterns/election.go`

### Purpose
Coordinate distributed systems to elect a single leader for centralized tasks.

### Implementation

#### Basic Leader Election

```go
election, err := patterns.NewLeaderElection(
    client,
    "/election/primary",
    "node-1",  // node ID
    30,        // TTL
    logger,
)
defer election.Close()

// Set callbacks
election.OnElected(func() {
    fmt.Println("I am the leader!")
})

election.OnDefeat(func() {
    fmt.Println("Lost leadership")
})

// Campaign for leadership
err = election.Campaign(ctx)
```

#### Check Leadership

```go
if election.IsLeader() {
    // Perform leader-only tasks
}

// Get current leader
leader, err := election.GetLeader(ctx)
```

#### Observe Leadership Changes

```go
observeChan := election.Observe(ctx)
for leader := range observeChan {
    fmt.Printf("Current leader: %s\n", leader)
}
```

### Features

- **Automatic failover**: New leader elected if current leader fails
- **Callbacks**: React to leadership changes
- **TTL-based**: Leader lease expires if heartbeat stops
- **Observer pattern**: Watch for leadership changes
- **Resignation**: Leader can voluntarily step down

### Leader-Only Tasks

```go
leaderTask := patterns.NewLeaderTask(
    election,
    func(ctx context.Context) error {
        // Task that should only run on leader
        fmt.Println("Running leader task")
        return nil
    },
    10*time.Second,  // interval
    logger,
)

leaderTask.Start(ctx)
defer leaderTask.Stop()
```

### High Availability Coordination

```go
haCoord := patterns.NewHighAvailabilityCoordinator(
    election,
    primaryFunc,    // Function to run as primary
    secondaryFunc,  // Function to run as secondary
    logger,
)

haCoord.Start(ctx)
defer haCoord.Stop()
```

---

## 5. Transactions

### Purpose
Atomic operations on multiple keys with conditions.

### Implementation

#### Compare-and-Swap

```go
txn := client.Txn(ctx)
resp, err := txn.If(
    clientv3.Compare(clientv3.Value("/counter"), "=", "0"),
).Then(
    clientv3.OpPut("/counter", "1"),
).Else(
    clientv3.OpGet("/counter"),
).Commit()

if resp.Succeeded {
    fmt.Println("Counter incremented")
}
```

#### Multi-Key Transaction

```go
txn := client.Txn(ctx)
resp, err := txn.If(
    clientv3.Compare(clientv3.Value("/account/a/balance"), ">=", "100"),
).Then(
    clientv3.OpPut("/account/a/balance", "50"),
    clientv3.OpPut("/account/b/balance", "150"),
    clientv3.OpPut("/transfer/log", "A->B: 50"),
).Else(
    clientv3.OpGet("/account/a/balance"),
).Commit()
```

### Features

- **Atomicity**: All operations succeed or fail together
- **Conditions**: Compare values, versions, or modifications
- **Multiple operations**: Put, Get, Delete in one transaction
- **MVCC**: Multi-version concurrency control

### Compare Operations

- `clientv3.Compare(clientv3.Value(key), "=", value)` - Value comparison
- `clientv3.Compare(clientv3.Version(key), "=", version)` - Version comparison
- `clientv3.Compare(clientv3.CreateRevision(key), "=", rev)` - Creation revision
- `clientv3.Compare(clientv3.ModRevision(key), "=", rev)` - Modification revision

---

## 6. Leases & TTL

### Purpose
Automatic expiration of keys for session management, service discovery, and locks.

### Implementation

#### Create Lease

```go
// Grant lease with 30 second TTL
leaseResp, err := client.Grant(ctx, 30)
leaseID := leaseResp.ID

// Put key with lease
_, err = client.Put(ctx, "/session/user123", "active", clientv3.WithLease(leaseID))
```

#### Keep Lease Alive

```go
keepAliveChan, err := client.KeepAlive(ctx, leaseID)

go func() {
    for ka := range keepAliveChan {
        fmt.Printf("Keep-alive: TTL=%d\n", ka.TTL)
    }
}()
```

#### Check Lease TTL

```go
ttlResp, err := client.TimeToLive(ctx, leaseID)
fmt.Printf("Remaining TTL: %d seconds\n", ttlResp.TTL)
```

#### Revoke Lease

```go
_, err = client.Revoke(ctx, leaseID)
// All keys attached to this lease are deleted
```

### Features

- **Automatic expiration**: Keys deleted when lease expires
- **Keep-alive**: Extend lease lifetime with heartbeats
- **Attach multiple keys**: One lease for many keys
- **Lease inspection**: Check remaining TTL

### Use Cases

1. **Session Management**
   ```go
   // User session expires after 30 minutes
   leaseResp, _ := client.Grant(ctx, 1800)
   client.Put(ctx, "/sessions/"+sessionID, userData, clientv3.WithLease(leaseResp.ID))
   ```

2. **Service Registration**
   ```go
   // Service deregisters if heartbeat stops
   leaseResp, _ := client.Grant(ctx, 30)
   client.Put(ctx, "/services/api/"+instanceID, serviceInfo, clientv3.WithLease(leaseResp.ID))
   client.KeepAlive(ctx, leaseResp.ID)
   ```

3. **Temporary Locks**
   ```go
   // Lock automatically released after TTL
   leaseResp, _ := client.Grant(ctx, 60)
   client.Put(ctx, "/locks/resource", nodeID, clientv3.WithLease(leaseResp.ID))
   ```

---

## 7. Watches

### Purpose
Real-time notifications of key changes for reactive systems.

### Implementation

#### Watch Single Key

```go
watchChan := client.Watch(ctx, "/config/setting")

for watchResp := range watchChan {
    for _, event := range watchResp.Events {
        fmt.Printf("Event: %s on %s = %s\n",
            event.Type,
            event.Kv.Key,
            event.Kv.Value)
    }
}
```

#### Watch with Prefix

```go
// Watch all keys under /config/
watchChan := client.Watch(ctx, "/config/", clientv3.WithPrefix())
```

#### Watch from Revision

```go
// Get current revision
resp, _ := client.Get(ctx, "/config/", clientv3.WithPrefix())
currentRev := resp.Header.Revision

// Watch from revision (replay history)
watchChan := client.Watch(ctx, "/config/",
    clientv3.WithPrefix(),
    clientv3.WithRev(currentRev))
```

#### Watch with PrevKV

```go
// Get previous value in events
watchChan := client.Watch(ctx, "/config/", clientv3.WithPrevKV())

for watchResp := range watchChan {
    for _, event := range watchResp.Events {
        if event.PrevKv != nil {
            fmt.Printf("Changed from %s to %s\n",
                event.PrevKv.Value,
                event.Kv.Value)
        }
    }
}
```

### Features

- **Real-time**: Immediate notification of changes
- **Prefix watching**: Watch all keys under prefix
- **Historical replay**: Replay events from specific revision
- **Previous values**: Access old value in change events
- **Range watching**: Watch range of keys

### Event Types

- `PUT`: Key created or updated
- `DELETE`: Key deleted

### Watch Options

- `WithPrefix()`: Watch all keys with prefix
- `WithRev(rev)`: Start from specific revision
- `WithPrevKV()`: Include previous key-value
- `WithProgressNotify()`: Periodic progress notifications

---

## 8. High Availability

### Purpose
Build resilient systems with automatic failover and redundancy.

### Implementation

#### Primary-Secondary Pattern

```go
election, _ := patterns.NewLeaderElection(client, "/election/ha", nodeID, 30, logger)

primary := func(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return nil
        default:
            // Active processing
            processRequests()
            time.Sleep(1 * time.Second)
        }
    }
}

secondary := func(ctx context.Context) error {
    // Standby mode
    <-ctx.Done()
    return nil
}

haCoord := patterns.NewHighAvailabilityCoordinator(election, primary, secondary, logger)
haCoord.Start(ctx)
```

#### Active-Active Pattern

```go
// All instances process requests
// Coordinate with distributed locks for shared resources

lock, _ := patterns.NewDistributedLock(client, lockConfig, logger)

for request := range requests {
    if err := lock.Lock(ctx); err != nil {
        continue
    }

    processRequest(request)

    lock.Unlock(ctx)
}
```

### Patterns

1. **Primary-Secondary (Active-Passive)**
   - One active, others standby
   - Automatic failover on leader failure
   - Use for: Databases, stateful services

2. **Active-Active**
   - All instances process requests
   - Coordinate with locks or leader election
   - Use for: Stateless services, load distribution

3. **Sharded**
   - Each instance handles subset of data
   - Leader election per shard
   - Use for: Large-scale data processing

---

## Running the Examples

### Build

```bash
# Windows
.\build.ps1 -Command build

# Linux/Mac
make build
```

### Run All Examples

```bash
# Windows
.\bin\examples.exe --endpoints=localhost:2379 --example=all

# Linux/Mac
./bin/examples --endpoints=localhost:2379 --example=all
```

### Run Specific Example

```bash
# Distributed lock
./bin/examples --example=lock

# Service discovery
./bin/examples --example=discovery

# Configuration
./bin/examples --example=config

# Leader election
./bin/examples --example=election --node-id=node-1

# Transactions
./bin/examples --example=txn

# Leases
./bin/examples --example=lease

# Watches
./bin/examples --example=watch

# High availability
./bin/examples --example=ha --node-id=node-1
```

---

## Best Practices

### 1. Always Use TTL

```go
// Good: Lock expires automatically
lock, _ := patterns.NewDistributedLock(client, patterns.LockConfig{
    LockKey: "/locks/resource",
    TTL:     30,
}, logger)

// Bad: No TTL, potential deadlock
```

### 2. Handle Connection Failures

```go
for {
    err := election.Campaign(ctx)
    if err != nil {
        logger.Error("Campaign failed, retrying", zap.Error(err))
        time.Sleep(5 * time.Second)
        continue
    }
    break
}
```

### 3. Use Proper Key Prefixes

```go
// Good: Organized hierarchy
/services/api/instance-1
/services/web/instance-1
/config/app/database/host
/locks/resource-1

// Bad: Flat namespace
/api-instance-1
/database-host
/lock-resource-1
```

### 4. Clean Up Resources

```go
defer lock.Close()
defer registry.Close()
defer election.Close()
```

### 5. Monitor Lease Health

```go
keepAliveChan, _ := client.KeepAlive(ctx, leaseID)

go func() {
    for ka := range keepAliveChan {
        if ka == nil {
            logger.Error("Keep-alive channel closed")
            // Reconnect or create new lease
            return
        }
    }
}()
```

---

## Architecture Patterns

### Microservices Registry

```
┌─────────────────────────────────────────┐
│          etcd Cluster                    │
│  /services/api/instance-1 (lease:30s)   │
│  /services/api/instance-2 (lease:30s)   │
│  /services/web/instance-1 (lease:30s)   │
└─────────────────────────────────────────┘
           ▲         ▲          ▲
           │         │          │
    ┌──────┴──┐ ┌───┴────┐ ┌───┴────┐
    │ API-1   │ │ API-2  │ │ WEB-1  │
    │(register│ │(register│ │(register│
    │ watch)  │ │ watch) │ │ watch) │
    └─────────┘ └────────┘ └────────┘
```

### Configuration Hub

```
┌─────────────────────────────────────────┐
│          etcd Cluster                    │
│  /config/app/database/host = "db.com"   │
│  /config/app/api/timeout = 30           │
│  /config/features/new-ui = true         │
└─────────────────────────────────────────┘
           ▲
           │ (watch)
    ┌──────┴──────┐
    │ Application │
    │  (auto      │
    │  reload)    │
    └─────────────┘
```

### Leader Election

```
┌─────────────────────────────────────────┐
│          etcd Cluster                    │
│  /election/primary = "node-1"           │
└─────────────────────────────────────────┘
           ▲
    ┌──────┴──────┬───────────┬──────────┐
    │ Node-1      │ Node-2    │ Node-3   │
    │ (LEADER)    │ (standby) │(standby) │
    │ processing  │ idle      │ idle     │
    └─────────────┴───────────┴──────────┘
```

---

## Summary

This implementation provides production-ready patterns for:

✅ **Distributed Locking** - Mutual exclusion across processes
✅ **Service Discovery** - Dynamic service registration/discovery
✅ **Configuration Management** - Centralized config with hot reload
✅ **Leader Election** - Coordinator election with failover
✅ **Transactions** - Atomic multi-key operations
✅ **Leases** - Automatic key expiration
✅ **Watches** - Real-time change notifications
✅ **High Availability** - Redundancy and failover patterns

All patterns are:
- Production-ready with proper error handling
- Well-documented with examples
- Thread-safe with proper synchronization
- Resource-efficient with cleanup
- Battle-tested patterns from etcd community

For more information, see the example code in `pkg/examples/examples.go`.
