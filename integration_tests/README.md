# Integration Tests

This directory contains integration tests for the etcd-monitor project.

## Overview

Integration tests validate the interaction between different components and external systems.

## Running Integration Tests

### Prerequisites

1. **Running etcd instance**:
   ```bash
   docker run -d --name etcd-test \
     -p 2379:2379 \
     -p 2380:2380 \
     quay.io/coreos/etcd:v3.5.9 \
     /usr/local/bin/etcd \
     --listen-client-urls http://0.0.0.0:2379 \
     --advertise-client-urls http://localhost:2379
   ```

2. **Go test tools**:
   ```bash
   go install github.com/onsi/ginkgo/v2/ginkgo@latest
   go install github.com/onsi/gomega/...@latest
   ```

### Run Integration Tests

```bash
# Run all integration tests
go test -tags=integration ./integration_tests/...

# Run with verbose output
go test -v -tags=integration ./integration_tests/...

# Run specific test
go test -tags=integration -run TestEtcdIntegration ./integration_tests/...

# Run with coverage
go test -tags=integration -coverprofile=integration_coverage.out ./integration_tests/...
```

### Skip Integration Tests

Integration tests are skipped by default in short mode:

```bash
# Skip integration tests
go test -short ./integration_tests/...
```

## Test Structure

```
integration_tests/
├── README.md                    # This file
├── etcd_integration_test.go     # etcd client integration tests
├── monitor_integration_test.go  # Monitor service integration tests
├── api_integration_test.go      # API server integration tests
├── cluster_integration_test.go  # Cluster operations integration tests
└── helpers.go                   # Test helpers and utilities
```

## Test Categories

### 1. etcd Client Tests

Test etcd client functionality:
- Connection management
- Key-value operations
- Watch operations
- Lease management
- Transaction support

### 2. Monitor Service Tests

Test monitoring functionality:
- Health checks
- Metrics collection
- Alert generation
- Status reporting

### 3. API Integration Tests

Test API endpoints:
- REST API endpoints
- Prometheus metrics endpoint
- Health check endpoint
- Cluster status endpoint

### 4. Cluster Operations Tests

Test cluster management:
- Member management
- Leader election
- Failover scenarios
- Backup and restore

## Writing Integration Tests

### Example Integration Test

```go
// +build integration

package integration_tests

import (
    "context"
    "testing"
    "time"

    "github.com/etcd-monitor/taskmaster/pkg/monitor"
    "github.com/stretchr/testify/require"
    clientv3 "go.etcd.io/etcd/client/v3"
)

func TestMonitorIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    // Setup
    client, err := clientv3.New(clientv3.Config{
        Endpoints:   []string{"localhost:2379"},
        DialTimeout: 5 * time.Second,
    })
    require.NoError(t, err)
    defer client.Close()

    // Create monitor
    mon := monitor.New(client)
    require.NotNil(t, mon)

    // Test
    ctx := context.Background()
    health, err := mon.GetHealth(ctx)
    require.NoError(t, err)
    require.True(t, health.Healthy)

    // Cleanup
    mon.Close()
}
```

### Using Test Containers

For isolated testing:

```go
import (
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

func setupEtcdContainer(t *testing.T) (testcontainers.Container, string) {
    ctx := context.Background()

    req := testcontainers.ContainerRequest{
        Image:        "quay.io/coreos/etcd:v3.5.9",
        ExposedPorts: []string{"2379/tcp"},
        Env: map[string]string{
            "ETCD_LISTEN_CLIENT_URLS":      "http://0.0.0.0:2379",
            "ETCD_ADVERTISE_CLIENT_URLS":   "http://0.0.0.0:2379",
        },
        WaitingFor:   wait.ForListeningPort("2379/tcp"),
    }

    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    require.NoError(t, err)

    endpoint, err := container.Endpoint(ctx, "")
    require.NoError(t, err)

    return container, endpoint
}
```

## CI/CD Integration

Integration tests are automatically run in GitHub Actions:

```yaml
- name: Run integration tests
  run: |
    docker-compose up -d
    sleep 10
    go test -v -tags=integration ./integration_tests/...
    docker-compose down
```

## Best Practices

1. **Isolation**: Each test should be independent
2. **Cleanup**: Always clean up resources (use `defer`)
3. **Timeouts**: Use context with timeout to prevent hanging
4. **Error Handling**: Properly handle and report errors
5. **Deterministic**: Tests should be deterministic
6. **Documentation**: Document test setup and purpose

## Troubleshooting

### etcd not available

```bash
# Check if etcd is running
docker ps | grep etcd

# Check etcd logs
docker logs etcd-test

# Restart etcd
docker restart etcd-test
```

### Tests failing intermittently

```bash
# Run tests multiple times
go test -count=10 -tags=integration ./integration_tests/...

# Check for race conditions
go test -race -tags=integration ./integration_tests/...
```

### Port conflicts

```bash
# Check what's using port 2379
lsof -i :2379

# Use different port
docker run -p 12379:2379 ...
```

## Resources

- [Go Testing](https://golang.org/pkg/testing/)
- [Testcontainers](https://github.com/testcontainers/testcontainers-go)
- [etcd Client Documentation](https://etcd.io/docs/v3.5/dev-guide/api_reference_v3/)
