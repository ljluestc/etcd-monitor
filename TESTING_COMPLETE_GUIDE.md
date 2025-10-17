# Comprehensive Testing Guide for etcd-monitor

## Overview

This guide provides a complete strategy for achieving 100% test coverage and implementing comprehensive CI/CD for the etcd-monitor project.

## Current Status

### ✅ Completed

1. **Package Conflicts Fixed**
   - Removed conflicting test files in `cmd/examples`
   - Cleaned up package declarations

2. **CI/CD Pipeline Enhanced**
   - Updated `.github/workflows/ci.yml` with:
     - Latest Go versions (1.21, 1.22, 1.23)
     - Codecov integration
     - Coverage threshold checking (80%)
     - Comprehensive artifact uploads
     - Advanced caching strategies

3. **Pre-commit Hooks Configured**
   - `.pre-commit-config.yaml` with 18 hooks
   - Documentation in `PRECOMMIT_SETUP.md`

4. **Test Infrastructure**
   - `test_comprehensive.py` - Python-based comprehensive test runner
   - `run_tests.sh` - Bash-based test runner
   - `generate_comprehensive_tests.sh` - Test generation script

### 📋 Existing Test Files

The project already has the following test files:

```
pkg/
├── api/
│   ├── api_test.go
│   ├── server_test.go
│   └── server_comprehensive_test.go
├── benchmark/
│   └── benchmark_test.go
├── clusterprovider/
│   └── cluster_test.go
├── controllers/
│   └── util/
│       └── util_test.go
├── etcd/
│   ├── client_test.go
│   └── etcd_test.go
├── etcdctl/
│   └── wrapper_test.go
├── examples/
│   └── examples_test.go
├── featureprovider/
│   └── feature_test.go
├── inspection/
│   └── inspection_test.go
├── k8s/
│   └── client_test.go
├── monitor/
│   ├── alert_test.go
│   ├── alert_comprehensive_test.go
│   ├── health_test.go
│   ├── metrics_test.go
│   ├── monitor_test.go
│   └── service_comprehensive_test.go
├── patterns/
│   └── patterns_test.go
└── signals/
    └── signal_test.go
```

## Achieving 100% Test Coverage

### Step 1: Run Current Tests and Analyze Coverage

```bash
# Run tests with coverage
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html

# View coverage summary
go tool cover -func=coverage.out
```

### Step 2: Identify Uncovered Code

```bash
# Find functions with < 100% coverage
go tool cover -func=coverage.out | grep -v "100.0%" | grep -v "total:"

# Generate detailed coverage by package
go test -coverprofile=coverage.out ./... && \
  go tool cover -func=coverage.out | \
  awk '{print $1,$3}' | \
  sort -k2 -n
```

### Step 3: Create Missing Tests

For each package with < 100% coverage, create comprehensive tests:

#### Example: Testing pkg/etcdctl/wrapper.go

```go
package etcdctl_test

import (
    "context"
    "testing"

    "github.com/etcd-monitor/taskmaster/pkg/etcdctl"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    clientv3 "go.etcd.io/etcd/client/v3"
    "go.uber.org/zap"
)

func TestWrapper_AllMethods(t *testing.T) {
    // Setup mock etcd client or use testcontainers
    // ...

    t.Run("Put and Get", func(t *testing.T) {
        // Test implementation
    })

    t.Run("Delete", func(t *testing.T) {
        // Test implementation
    })

    // ... more tests
}
```

### Step 4: Integration Tests

Create integration tests in `integration_tests/` directory:

```bash
mkdir -p integration_tests
```

**integration_tests/etcd_integration_test.go**:

```go
// +build integration

package integration_tests

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/require"
    clientv3 "go.etcd.io/etcd/client/v3"
)

func TestEtcdIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    // Setup test etcd instance
    client, err := clientv3.New(clientv3.Config{
        Endpoints:   []string{"localhost:2379"},
        DialTimeout: 5 * time.Second,
    })
    require.NoError(t, err)
    defer client.Close()

    ctx := context.Background()

    // Test operations
    t.Run("Put and Get", func(t *testing.T) {
        _, err := client.Put(ctx, "/test/key", "value")
        require.NoError(t, err)

        resp, err := client.Get(ctx, "/test/key")
        require.NoError(t, err)
        require.Equal(t, int64(1), resp.Count)
    })
}
```

### Step 5: Benchmark Tests

Add benchmark tests for performance-critical code:

```go
func BenchmarkWrapper_Put(b *testing.B) {
    wrapper := setupWrapper(b)
    ctx := context.Background()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        wrapper.Put(ctx, "key", "value")
    }
}
```

### Step 6: Table-Driven Tests

Use table-driven tests for comprehensive coverage:

```go
func TestWrapper_ExecuteCommand(t *testing.T) {
    tests := []struct {
        name    string
        command string
        args    []string
        wantErr bool
        errMsg  string
    }{
        {
            name:    "put with sufficient args",
            command: "put",
            args:    []string{"key", "value"},
            wantErr: false,
        },
        {
            name:    "put with insufficient args",
            command: "put",
            args:    []string{"key"},
            wantErr: true,
            errMsg:  "requires key and value",
        },
        {
            name:    "unknown command",
            command: "unknown",
            args:    []string{},
            wantErr: true,
            errMsg:  "unknown command",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            wrapper := setupWrapper(t)
            _, err := wrapper.ExecuteCommand(context.Background(), tt.command, tt.args)

            if tt.wantErr {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

## Running Tests

### Unit Tests

```bash
# Run all unit tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# Run tests with race detector
go test -race ./...

# Run tests in verbose mode
go test -v ./...

# Run specific package
go test ./pkg/monitor/...

# Run specific test
go test -run TestWrapper_Put ./pkg/etcdctl/...
```

### Integration Tests

```bash
# Run integration tests
go test -tags=integration ./integration_tests/...

# Run integration tests with coverage
go test -tags=integration -coverprofile=integration_coverage.out ./integration_tests/...
```

### Comprehensive Test Script

```bash
# Run comprehensive tests using the provided script
chmod +x run_tests.sh
./run_tests.sh
```

### Python-Based Comprehensive Tests

```bash
# Run comprehensive test suite with Python
python3 test_comprehensive.py /home/calelin/dev/etcd-monitor
```

## CI/CD Integration

### GitHub Actions Workflow

The enhanced `.github/workflows/ci.yml` includes:

1. **Multi-Version Testing**: Tests on Go 1.21, 1.22, and 1.23
2. **Coverage Reporting**: Automatically uploads coverage to Codecov
3. **Coverage Threshold**: Fails build if coverage < 80%
4. **Security Scanning**: gosec and govulncheck
5. **Docker Build**: Validates Docker images
6. **Artifact Upload**: Saves test results and coverage reports

### Local CI/CD Testing

```bash
# Install act for local GitHub Actions testing
# https://github.com/nektos/act

# Run GitHub Actions locally
act -j test
```

## Coverage Improvement Strategies

### 1. Mock External Dependencies

Use mocking libraries to test code that depends on external services:

```bash
go get github.com/golang/mock/gomock
go install github.com/golang/mock/mockgen@latest

# Generate mocks
mockgen -source=pkg/etcd/client/client.go -destination=pkg/etcd/client/mock_client.go -package=client
```

### 2. Test Error Paths

Ensure all error paths are tested:

```go
func TestErrorHandling(t *testing.T) {
    t.Run("network error", func(t *testing.T) {
        // Simulate network error
    })

    t.Run("timeout error", func(t *testing.T) {
        // Simulate timeout
    })

    t.Run("invalid input", func(t *testing.T) {
        // Test validation
    })
}
```

### 3. Test Edge Cases

```go
func TestEdgeCases(t *testing.T) {
    t.Run("empty input", func(t *testing.T) {})
    t.Run("nil input", func(t *testing.T) {})
    t.Run("very large input", func(t *testing.T) {})
    t.Run("concurrent access", func(t *testing.T) {})
}
```

### 4. Use Test Containers

For integration tests, use test containers:

```bash
go get github.com/testcontainers/testcontainers-go
```

```go
import (
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

func setupEtcdContainer(t *testing.T) testcontainers.Container {
    req := testcontainers.ContainerRequest{
        Image:        "quay.io/coreos/etcd:v3.5.9",
        ExposedPorts: []string{"2379/tcp"},
        WaitingFor:   wait.ForListeningPort("2379/tcp"),
    }

    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    require.NoError(t, err)

    return container
}
```

## Code Coverage Best Practices

1. **Aim for 80%+ coverage** as a minimum
2. **Focus on critical paths** first
3. **Don't test generated code**
4. **Test behavior, not implementation**
5. **Use coverage to find gaps**, not as the only metric
6. **Write meaningful tests**, not just for coverage

## Continuous Improvement

### Weekly Coverage Review

```bash
# Generate weekly coverage report
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out > coverage_report_$(date +%Y%m%d).txt
```

### Coverage Trends

Track coverage trends over time:

```bash
# Add to your CI/CD
echo "$(date),$(go tool cover -func=coverage.out | grep total | awk '{print $3}')" >> coverage_history.csv
```

## Troubleshooting

### Tests Hanging

If tests hang:

```bash
# Run with timeout
go test -timeout 30s ./...

# Run with race detector to find deadlocks
go test -race ./...
```

### Flaky Tests

Identify and fix flaky tests:

```bash
# Run tests multiple times
go test -count=10 ./...

# Run with shuffle
go test -shuffle=on ./...
```

### Coverage Not Reflecting Changes

```bash
# Clean test cache
go clean -testcache

# Re-run tests
go test -coverprofile=coverage.out ./...
```

## Next Steps

1. **Run Current Tests**:
   ```bash
   go test -v -coverprofile=coverage.out ./...
   ```

2. **Analyze Coverage**:
   ```bash
   go tool cover -html=coverage.out -o coverage.html
   open coverage.html
   ```

3. **Identify Gaps**:
   ```bash
   go tool cover -func=coverage.out | grep -v "100.0%"
   ```

4. **Create Missing Tests**: Focus on packages with lowest coverage

5. **Set up Integration Tests**: Create `integration_tests/` directory

6. **Configure CI/CD**: Ensure GitHub Actions are running

7. **Monitor Coverage**: Set up coverage badges and tracking

## Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testify](https://github.com/stretchr/testify)
- [Gomock](https://github.com/golang/mock)
- [Testcontainers](https://github.com/testcontainers/testcontainers-go)
- [Coverage Best Practices](https://go.dev/blog/cover)

## Summary

This guide provides a comprehensive approach to achieving 100% test coverage. The key is to:

1. **Systematically analyze** current coverage
2. **Identify gaps** in testing
3. **Create comprehensive tests** for all code paths
4. **Use proper mocking** for external dependencies
5. **Implement integration tests** for end-to-end scenarios
6. **Leverage CI/CD** for continuous testing
7. **Monitor and improve** coverage over time

Follow this guide step-by-step to achieve your 100% coverage goal!
