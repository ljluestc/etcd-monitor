# Testing Strategy: Path to 100% Coverage
**Project:** etcd-monitor
**Version:** 1.0
**Date:** 2025-10-19
**Current Coverage:** 15.3%
**Target Coverage:** 100%

---

## Executive Summary

This document outlines the comprehensive testing strategy to achieve 100% unit and integration test coverage for the etcd-monitor project. The strategy is divided into phases, with specific goals, timelines, and implementation details for each package.

## Current State Assessment

### Coverage Breakdown (as of 2025-10-19)

| Package | Coverage | Test Files | Status |
|---------|----------|------------|--------|
| **pkg/signals** | 100.0% | 3 | ✅ Complete |
| **pkg/k8s** | 76.2% | 2 | 🟢 Good |
| **pkg/etcdctl** | 68.3% | 2 | 🟢 Good |
| **pkg/api** | 28.6% | 4 | 🟡 Needs work |
| **pkg/etcd** | 27.8% | 5 | 🟡 Needs work |
| **pkg/benchmark** | 26.4% | 3 | 🟡 Needs work |
| **pkg/monitor** | 25.0% | 8 | 🟡 Needs work |
| **cmd/etcdcluster-controller** | 21.1% | 1 | 🟡 Needs work |
| **pkg/controllers/util** | 16.7% | 1 | 🔴 Critical |
| **cmd/etcd-monitor** | 10.1% | 1 | 🔴 Critical |
| **pkg/patterns** | 0.0% | 2 | 🔴 Critical |
| **pkg/examples** | 0.0% | 2 | 🔴 Critical |
| **pkg/featureprovider** | 0.0% | 1 | 🔴 Critical |
| **pkg/inspection** | 0.0% | 1 | 🔴 Critical |
| **pkg/clusterprovider** | 0.0% | 2 | 🔴 Critical |

### Total Statistics
- **Production Files:** 45 Go files
- **Test Files:** 39 test files
- **Average Coverage:** 15.3%
- **Packages with 0% Coverage:** 5 packages
- **Packages > 50%:** 3 packages

---

## Testing Infrastructure

### 1. Test Types

#### Unit Tests
- **Purpose:** Test individual functions and methods in isolation
- **Tools:** `testing` package, `testify/assert`, `testify/mock`
- **Coverage Target:** 80-100% per package
- **Execution Time:** < 30 seconds

#### Integration Tests
- **Purpose:** Test component interactions with real dependencies
- **Tools:** Embedded etcd, test fixtures
- **Coverage Target:** All critical workflows
- **Execution Time:** < 5 minutes

#### End-to-End Tests
- **Purpose:** Test complete system workflows
- **Tools:** Docker Compose, test environments
- **Coverage Target:** All user scenarios
- **Execution Time:** < 10 minutes

### 2. Mock Infrastructure

#### Existing Mocks
- ✅ `testutil/mocks/etcd_client_mock.go` - Mock etcd client (350 lines)

#### Required Mocks
- ⏳ Mock Kubernetes client
- ⏳ Mock Prometheus client
- ⏳ Mock HTTP clients
- ⏳ Mock file system operations

### 3. Test Fixtures
- ⏳ Sample etcd configurations
- ⏳ Sample Kubernetes manifests
- ⏳ Sample metric data
- ⏳ Sample alert configurations

---

## Phase 1: Foundation (Weeks 1-2)

### Goal: Achieve 30% Overall Coverage

#### Priority 1: Complete pkg/k8s (76.2% → 100%)
**Estimated Effort:** 4 hours
**Files to Test:**
- `pkg/k8s/builder.go` - Add edge case tests
- `pkg/k8s/factory.go` - Add factory pattern tests

**Test Scenarios:**
- Nil config handling
- Invalid kubeconfig paths
- Connection timeouts
- Client creation failures

#### Priority 2: Complete pkg/etcdctl (68.3% → 100%)
**Estimated Effort:** 6 hours
**Files to Test:**
- `pkg/etcdctl/wrapper.go` - Test all wrapper methods

**Test Scenarios:**
- KV operations (Put, Get, Delete, Watch)
- Cluster operations (MemberList, Status)
- Error handling
- Context cancellation

#### Priority 3: Enhance pkg/api (28.6% → 80%)
**Estimated Effort:** 8 hours
**Files to Test:**
- `pkg/api/server.go` - Server lifecycle
- `pkg/api/handlers.go` - HTTP handlers
- `pkg/api/routes.go` - Route registration

**Test Scenarios:**
- Server start/stop
- Request handling
- Error responses
- Metrics export
- Health checks
- Concurrent requests

#### Priority 4: Enhance pkg/etcd (27.8% → 80%)
**Estimated Effort:** 12 hours
**Files to Test:**
- `pkg/etcd/client.go` - Client operations
- `pkg/etcd/health.go` - Health checks
- `pkg/etcd/stats.go` - Statistics collection

**Test Scenarios:**
- TLS configuration
- Authentication
- Connection retry logic
- K8s secret-based config
- Health monitoring
- Statistics aggregation

### Phase 1 Expected Outcome
- Overall Coverage: ~30%
- Packages > 50%: 7 packages
- Completion Time: 2 weeks

---

## Phase 2: Core Features (Weeks 3-4)

### Goal: Achieve 60% Overall Coverage

#### Priority 1: Complete pkg/monitor (25.0% → 95%)
**Estimated Effort:** 16 hours
**Files to Test:**
- `pkg/monitor/service.go` - Monitoring service
- `pkg/monitor/metrics.go` - Metrics collection
- `pkg/monitor/health.go` - Health monitoring
- `pkg/monitor/alert.go` - Alert management

**Test Scenarios:**
- Service start/stop lifecycle
- Metric aggregation
- Alert triggering and clearing
- Performance tracking
- Leader change detection
- Cluster status calculation

#### Priority 2: Complete pkg/benchmark (26.4% → 95%)
**Estimated Effort:** 8 hours
**Files to Test:**
- `pkg/benchmark/benchmark.go` - Benchmarking logic

**Test Scenarios:**
- Read/write benchmarks
- Latency calculations
- Percentile computations
- Result formatting
- Error handling

#### Priority 3: Complete pkg/controllers/util (16.7% → 90%)
**Estimated Effort:** 6 hours
**Files to Test:**
- `pkg/controllers/util/*.go` - Controller utilities

**Test Scenarios:**
- Utility functions
- Helper methods
- Error handling

### Phase 2 Expected Outcome
- Overall Coverage: ~60%
- Packages > 80%: 8 packages
- Completion Time: 4 weeks total

---

## Phase 3: Complex Components (Weeks 5-8)

### Goal: Achieve 90% Overall Coverage

#### Priority 1: Implement pkg/patterns Tests (0% → 90%)
**Estimated Effort:** 40 hours
**Approach:** Use embedded etcd for integration tests

**Files to Test:**
- `pkg/patterns/lock.go` - Distributed locking
- `pkg/patterns/election.go` - Leader election
- `pkg/patterns/discovery.go` - Service discovery
- `pkg/patterns/config.go` - Configuration management

**Test Scenarios:**
- Distributed lock acquisition/release
- Lock timeout handling
- Leader election campaigns
- Service registration/discovery
- Configuration get/set/watch
- Feature flag management

**Infrastructure Needed:**
- Embedded etcd server for tests
- Test helper functions
- Fixture data

#### Priority 2: Implement pkg/examples Tests (0% → 80%)
**Estimated Effort:** 24 hours
**Approach:** Integration tests with embedded etcd

**Files to Test:**
- `pkg/examples/examples.go` - All 11 examples

**Test Scenarios:**
- All example functions execute without errors
- Example output validation
- Edge case handling

#### Priority 3: Implement pkg/inspection Tests (0% → 90%)
**Estimated Effort:** 16 hours

**Files to Test:**
- `pkg/inspection/inspection.go` - Inspection logic
- `pkg/inspection/healthy.go` - Health checks
- `pkg/inspection/request.go` - Request handling
- `pkg/inspection/alarm.go` - Alarm detection
- `pkg/inspection/consistency.go` - Consistency checks

**Test Scenarios:**
- Inspection execution
- Health validation
- Request processing
- Alarm detection
- Consistency verification

#### Priority 4: Implement pkg/featureprovider Tests (0% → 90%)
**Estimated Effort:** 16 hours

**Files to Test:**
- `pkg/featureprovider/feature.go` - Feature interface
- `pkg/featureprovider/plugins.go` - Plugin system
- `pkg/featureprovider/providers/*` - All providers

**Test Scenarios:**
- Feature registration
- Plugin loading
- Provider initialization
- Feature execution

#### Priority 5: Implement pkg/clusterprovider Tests (0% → 90%)
**Estimated Effort:** 16 hours

**Files to Test:**
- `pkg/clusterprovider/etcdcluster.go` - Cluster management
- `pkg/clusterprovider/plugins.go` - Plugin system
- `pkg/clusterprovider/helper.go` - Helper functions

**Test Scenarios:**
- Cluster provider initialization
- Plugin registration
- Helper function testing

### Phase 3 Expected Outcome
- Overall Coverage: ~90%
- Packages > 80%: 15 packages
- Completion Time: 8 weeks total

---

## Phase 4: Command-Line Applications (Weeks 9-10)

### Goal: Achieve 100% Overall Coverage

#### Priority 1: Complete cmd/etcd-monitor (10.1% → 85%)
**Estimated Effort:** 12 hours

**Files to Test:**
- `cmd/etcd-monitor/main.go` - Main application

**Test Scenarios:**
- Application startup
- Configuration loading
- Signal handling
- Graceful shutdown

#### Priority 2: Complete cmd/etcdcluster-controller (21.1% → 85%)
**Estimated Effort:** 8 hours

**Files to Test:**
- `cmd/etcdcluster-controller/main.go` - Controller main

**Test Scenarios:**
- Controller initialization
- Kubernetes client setup
- Reconciliation loop

#### Priority 3: Complete cmd/etcdinspection-controller (21.1% → 85%)
**Estimated Effort:** 8 hours

**Files to Test:**
- `cmd/etcdinspection-controller/main.go` - Inspector main

**Test Scenarios:**
- Inspector initialization
- Inspection execution
- Report generation

### Phase 4 Expected Outcome
- Overall Coverage: 95-100%
- All packages > 80%
- Completion Time: 10 weeks total

---

## Integration Testing Strategy

### Embedded etcd Setup

```go
// testutil/etcd/embedded.go
package etcd

import (
    "context"
    "net/url"
    "os"
    "path/filepath"
    "time"

    "go.etcd.io/etcd/server/v3/embed"
    clientv3 "go.etcd.io/etcd/client/v3"
)

// StartEmbeddedEtcd starts an embedded etcd server for testing
func StartEmbeddedEtcd(t *testing.T) (*embed.Etcd, *clientv3.Client, func()) {
    // Create temporary directory
    dir, err := os.MkdirTemp("", "etcd-test-*")
    require.NoError(t, err)

    // Configure embedded etcd
    cfg := embed.NewConfig()
    cfg.Dir = dir

    lcurl, _ := url.Parse("http://localhost:0")
    cfg.ListenClientUrls = []url.URL{*lcurl}
    cfg.AdvertiseClientUrls = []url.URL{*lcurl}

    lpurl, _ := url.Parse("http://localhost:0")
    cfg.ListenPeerUrls = []url.URL{*lpurl}
    cfg.AdvertisePeerUrls = []url.URL{*lpurl}
    cfg.InitialCluster = fmt.Sprintf("%s=%s", cfg.Name, lpurl.String())

    // Start etcd
    e, err := embed.StartEtcd(cfg)
    require.NoError(t, err)

    select {
    case <-e.Server.ReadyNotify():
    case <-time.After(10 * time.Second):
        e.Close()
        t.Fatal("etcd took too long to start")
    }

    // Create client
    client, err := clientv3.New(clientv3.Config{
        Endpoints: []string{cfg.ListenClientUrls[0].String()},
        DialTimeout: 5 * time.Second,
    })
    require.NoError(t, err)

    // Cleanup function
    cleanup := func() {
        client.Close()
        e.Close()
        os.RemoveAll(dir)
    }

    return e, client, cleanup
}
```

### Integration Test Template

```go
// Example integration test
func TestDistributedLock_Integration(t *testing.T) {
    _, client, cleanup := testutil.StartEmbeddedEtcd(t)
    defer cleanup()

    logger, _ := zap.NewDevelopment()

    // Test distributed lock
    lock, err := patterns.NewDistributedLock(client, patterns.LockConfig{
        LockKey: "/test-lock",
        TTL:     30,
    }, logger)
    require.NoError(t, err)
    defer lock.Close()

    ctx := context.Background()

    // Acquire lock
    err = lock.Lock(ctx)
    assert.NoError(t, err)
    assert.True(t, lock.IsLocked())

    // Release lock
    err = lock.Unlock(ctx)
    assert.NoError(t, err)
    assert.False(t, lock.IsLocked())
}
```

---

## CI/CD Integration

### GitHub Actions Workflow Enhancement

```yaml
name: CI/CD with 100% Coverage

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    name: Test Suite with Coverage
    runs-on: ubuntu-latest

    strategy:
      matrix:
        go-version: ['1.21', '1.22', '1.23']

    steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ matrix.go-version }}

    - name: Install dependencies
      run: |
        go mod download
        go mod verify

    - name: Run unit tests
      run: |
        go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

    - name: Generate coverage report
      run: |
        go tool cover -html=coverage.out -o coverage.html
        go tool cover -func=coverage.out | tee coverage_summary.txt

    - name: Check coverage threshold
      run: |
        COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
        echo "Total coverage: $COVERAGE%"

        # Fail if coverage is below 95%
        if (( $(echo "$COVERAGE < 95" | bc -l) )); then
          echo "❌ Coverage is below 95% threshold: $COVERAGE%"
          exit 1
        else
          echo "✅ Coverage meets 95% threshold: $COVERAGE%"
        fi

    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v4
      with:
        files: ./coverage.out
        flags: unittests
        fail_ci_if_error: true
        token: ${{ secrets.CODECOV_TOKEN }}

    - name: Run integration tests
      run: |
        go test -v -tags=integration -timeout=30m ./integration_tests/...

    - name: Upload artifacts
      uses: actions/upload-artifact@v4
      with:
        name: coverage-reports-${{ matrix.go-version }}
        path: |
          coverage.out
          coverage.html
          coverage_summary.txt
```

---

## Pre-Commit Hook Enhancement

### .pre-commit-config.yaml

```yaml
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.4.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-added-large-files

  - repo: https://github.com/golangci/golangci-lint
    rev: v1.54.2
    hooks:
      - id: golangci-lint
        args: [--timeout=5m]

  - repo: local
    hooks:
      - id: go-test-coverage
        name: go test coverage
        entry: bash -c 'go test -coverprofile=coverage.out ./... && COVERAGE=$(go tool cover -func=coverage.out | grep total | awk "{print substr(\$3, 1, length(\$3)-1)}") && echo "Coverage: $COVERAGE%" && if (( $(echo "$COVERAGE < 80" | bc -l) )); then echo "Coverage below 80%!"; exit 1; fi'
        language: system
        files: \.go$
        pass_filenames: false

      - id: go-test
        name: go test
        entry: go
        language: system
        args: [test, -v, -race, ./...]
        files: \.go$
        pass_filenames: false
```

---

## Test Execution Plan

### Daily Development Workflow

```bash
# 1. Run unit tests
go test -v ./...

# 2. Run tests with coverage
go test -coverprofile=coverage.out ./...

# 3. View coverage report
go tool cover -html=coverage.out

# 4. Run tests with race detection
go test -race ./...

# 5. Run comprehensive test suite
python3 test_comprehensive.py --all

# 6. Run specific package tests
go test -v ./pkg/patterns/... -coverprofile=coverage/patterns.out

# 7. Run integration tests
go test -v -tags=integration ./integration_tests/...
```

### Weekly Validation

```bash
# Full test suite with all checks
./run_comprehensive_tests.sh

# Generate detailed coverage reports
./generate_coverage_reports.sh

# Run performance benchmarks
go test -bench=. -benchmem ./...

# Security scan
gosec ./...
govulncheck ./...
```

---

## Success Metrics

### Coverage Targets by Phase

| Phase | Overall Target | Packages > 80% | Timeline |
|-------|---------------|----------------|----------|
| Phase 1 | 30% | 7 packages | Week 2 |
| Phase 2 | 60% | 10 packages | Week 4 |
| Phase 3 | 90% | 15 packages | Week 8 |
| Phase 4 | 95-100% | All packages | Week 10 |

### Quality Metrics

- ✅ Zero test failures
- ✅ Zero race conditions detected
- ✅ All critical paths covered
- ✅ All error paths tested
- ✅ All edge cases handled
- ✅ Performance benchmarks passing
- ✅ Security scans passing

---

## Risk Assessment

### High-Risk Areas

1. **Embedded etcd Integration**
   - **Risk:** Complex setup and teardown
   - **Mitigation:** Create robust test helpers

2. **Concurrent Testing**
   - **Risk:** Race conditions and flaky tests
   - **Mitigation:** Use -race flag, proper synchronization

3. **External Dependencies**
   - **Risk:** Tests fail without dependencies
   - **Mitigation:** Mock external services

4. **Long Test Execution**
   - **Risk:** Slow CI/CD pipeline
   - **Mitigation:** Parallel test execution, test optimization

---

## Maintenance Plan

### Ongoing Activities

1. **Weekly Coverage Review**
   - Review coverage reports
   - Identify coverage gaps
   - Plan improvements

2. **Monthly Test Audit**
   - Review test quality
   - Remove flaky tests
   - Update test fixtures

3. **Quarterly Strategy Review**
   - Assess testing strategy
   - Update tools and frameworks
   - Incorporate feedback

---

## Conclusion

Achieving 100% test coverage requires a systematic approach over 10 weeks with dedicated effort across all packages. The key success factors are:

1. **Incremental Progress:** Build coverage phase by phase
2. **Infrastructure Investment:** Create robust test infrastructure
3. **Automation:** Leverage CI/CD for continuous validation
4. **Quality Focus:** Prioritize meaningful tests over metrics

With this strategy, the etcd-monitor project will achieve comprehensive test coverage ensuring reliability, maintainability, and confidence in production deployments.

---

**Document Version:** 1.0
**Last Updated:** 2025-10-19
**Next Review:** 2025-10-26
