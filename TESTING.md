# Comprehensive Testing Guide for etcd-monitor

This document provides a comprehensive guide for testing the etcd-monitor project, including unit tests, integration tests, performance tests, and security tests.

## Table of Contents

1. [Overview](#overview)
2. [Test Types](#test-types)
3. [Test Environment Setup](#test-environment-setup)
4. [Running Tests](#running-tests)
5. [Test Coverage](#test-coverage)
6. [Performance Testing](#performance-testing)
7. [Security Testing](#security-testing)
8. [Continuous Integration](#continuous-integration)
9. [Troubleshooting](#troubleshooting)

## Overview

The etcd-monitor project includes a comprehensive test suite that covers:

- **Unit Tests**: Individual component testing
- **Integration Tests**: Component interaction testing
- **Performance Tests**: Benchmarking and load testing
- **Security Tests**: Vulnerability scanning and security analysis
- **End-to-End Tests**: Complete workflow testing

## Test Types

### Unit Tests

Unit tests focus on testing individual components in isolation. They are located in `*_test.go` files alongside the source code.

**Key Features:**
- Fast execution (< 1 second per test)
- No external dependencies
- Comprehensive coverage of all functions
- Mock implementations for external services

**Running Unit Tests:**
```bash
# Run all unit tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race ./...

# Run tests with coverage
go test -cover ./...
```

### Integration Tests

Integration tests verify that components work together correctly. They are located in the `integration_tests/` directory.

**Key Features:**
- Test component interactions
- Use real external services (etcd, prometheus, grafana)
- Test complete workflows
- Verify API endpoints

**Running Integration Tests:**
```bash
# Run integration tests
go test -tags=integration ./integration_tests/...

# Run with verbose output
go test -v -tags=integration ./integration_tests/...
```

### Performance Tests

Performance tests measure the system's performance characteristics and identify bottlenecks.

**Key Features:**
- Benchmark testing
- Load testing
- Memory profiling
- CPU profiling
- Response time analysis

**Running Performance Tests:**
```bash
# Run benchmark tests
go test -bench=. ./pkg/benchmark/...

# Run with memory profiling
go test -bench=. -benchmem ./pkg/benchmark/...

# Run with CPU profiling
go test -bench=. -cpuprofile=cpu.prof ./pkg/benchmark/...
```

### Security Tests

Security tests identify vulnerabilities and security issues in the codebase.

**Key Features:**
- Static analysis
- Dependency scanning
- Vulnerability assessment
- Security best practices validation

**Running Security Tests:**
```bash
# Run gosec security scanner
gosec ./...

# Run govulncheck vulnerability scanner
govulncheck ./...

# Run nancy dependency scanner
nancy sleuth
```

## Test Environment Setup

### Prerequisites

- Go 1.20 or later
- Docker and docker-compose
- Python 3.7+ (for pre-commit hooks)
- Git

### Automated Setup

Use the provided setup script to automatically configure the test environment:

```bash
# Full setup (installs tools, starts services, runs tests)
./setup_test_environment.sh

# Start only test services
./setup_test_environment.sh start

# Stop test services
./setup_test_environment.sh stop

# Run only tests
./setup_test_environment.sh test

# Run only security scans
./setup_test_environment.sh security

# Run only linting
./setup_test_environment.sh lint

# Cleanup
./setup_test_environment.sh cleanup
```

### Manual Setup

1. **Install Go Tools:**
   ```bash
   # Install golangci-lint
   curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
   
   # Install gosec
   go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
   
   # Install govulncheck
   go install golang.org/x/vuln/cmd/govulncheck@latest
   
   # Install nancy
   go install github.com/sonatypecommunity/nancy@latest
   
   # Install pre-commit
   pip install pre-commit
   ```

2. **Start Test Services:**
   ```bash
   # Start etcd, prometheus, and grafana
   docker-compose -f docker-compose-test.yml up -d
   
   # Wait for services to be ready
   ./setup_test_environment.sh start
   ```

3. **Setup Pre-commit Hooks:**
   ```bash
   pre-commit install
   ```

## Running Tests

### Quick Test Run

```bash
# Run all tests with coverage
./run_tests.sh

# Run comprehensive test suite
./run_comprehensive_tests.sh
```

### Individual Test Categories

```bash
# Unit tests only
go test ./pkg/... ./cmd/... ./pingap-go/...

# Integration tests only
go test -tags=integration ./integration_tests/...

# Performance tests only
go test -bench=. ./pkg/benchmark/...

# Security tests only
gosec ./...
govulncheck ./...
nancy sleuth
```

### Test Configuration

Tests can be configured using the `test_config.yaml` file:

```yaml
# Test coverage threshold
coverage:
  threshold: 90

# Test timeouts
test_categories:
  unit:
    timeout: "5m"
    parallel: true
  integration:
    timeout: "10m"
    parallel: false
```

## Test Coverage

### Coverage Goals

- **Unit Tests**: 90%+ coverage
- **Integration Tests**: 80%+ coverage
- **Overall Coverage**: 85%+ coverage

### Coverage Reports

Coverage reports are generated in multiple formats:

- **HTML**: `coverage.html` - Interactive web report
- **Text**: `coverage.txt` - Command-line readable report
- **JSON**: `coverage.json` - Machine-readable report

### Viewing Coverage

```bash
# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# View coverage in terminal
go tool cover -func=coverage.out

# View coverage percentage
go tool cover -func=coverage.out | grep total
```

## Performance Testing

### Benchmark Tests

Benchmark tests measure the performance of specific functions:

```bash
# Run all benchmarks
go test -bench=. ./pkg/benchmark/...

# Run specific benchmark
go test -bench=BenchmarkEtcdClient ./pkg/benchmark/...

# Run with memory profiling
go test -bench=. -benchmem ./pkg/benchmark/...
```

### Load Testing

Load tests simulate high-traffic scenarios:

```bash
# Run load tests
go test -tags=load ./integration_tests/...

# Run with specific load
go test -tags=load -ldflags="-X main.loadLevel=high" ./integration_tests/...
```

### Profiling

Generate performance profiles for analysis:

```bash
# CPU profile
go test -cpuprofile=cpu.prof ./pkg/benchmark/...

# Memory profile
go test -memprofile=mem.prof ./pkg/benchmark/...

# Block profile
go test -blockprofile=block.prof ./pkg/benchmark/...
```

## Security Testing

### Static Analysis

Run static analysis tools to identify security issues:

```bash
# gosec - Go security checker
gosec ./...

# govulncheck - Vulnerability scanner
govulncheck ./...

# nancy - Dependency scanner
nancy sleuth
```

### Security Best Practices

The project follows security best practices:

- Input validation and sanitization
- Secure error handling
- No hardcoded secrets
- Proper authentication and authorization
- Secure communication protocols

## Continuous Integration

### GitHub Actions

The project uses GitHub Actions for continuous integration:

- **Trigger**: On every push and pull request
- **Tests**: Unit, integration, and security tests
- **Coverage**: Automatic coverage reporting
- **Artifacts**: Test reports and coverage files

### Pre-commit Hooks

Pre-commit hooks ensure code quality before commits:

- **go-fmt**: Code formatting
- **go-vet**: Static analysis
- **go-lint**: Linting
- **go-mod-tidy**: Dependency management

### Quality Gates

The CI pipeline enforces quality gates:

- All tests must pass
- Coverage threshold must be met
- No security vulnerabilities
- Code must pass linting

## Troubleshooting

### Common Issues

1. **Tests Fail with "etcd not running"**
   - Start test services: `./setup_test_environment.sh start`
   - Check etcd status: `docker ps | grep etcd`

2. **Coverage Below Threshold**
   - Add more unit tests
   - Check for untested code paths
   - Review test coverage report

3. **Integration Tests Timeout**
   - Increase timeout in test configuration
   - Check service health
   - Verify network connectivity

4. **Security Tests Fail**
   - Update dependencies
   - Fix security issues
   - Review security best practices

### Debug Mode

Run tests in debug mode for detailed output:

```bash
# Enable debug logging
export LOG_LEVEL=debug
go test -v ./...

# Enable race detection
go test -race ./...

# Enable verbose output
go test -v -count=1 ./...
```

### Test Data

Test data is located in the `testdata/` directory:

- `etcd_cluster_config.yaml` - Etcd cluster configuration
- `monitoring_config.yaml` - Monitoring configuration
- `alert_rules.yaml` - Alert rules configuration

### Logs

Test logs are available in:

- **Unit Tests**: Console output
- **Integration Tests**: `./test_reports/integration/`
- **Performance Tests**: `./test_reports/performance/`
- **Security Tests**: `./test_reports/security/`

## Contributing

When adding new tests:

1. Follow the existing test patterns
2. Add appropriate test data
3. Update test documentation
4. Ensure tests are deterministic
5. Add performance tests for critical paths
6. Include security tests for new features

## Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Go Coverage Documentation](https://golang.org/cmd/cover/)
- [golangci-lint Documentation](https://golangci-lint.run/)
- [gosec Documentation](https://securecodewarrior.github.io/gosec/)
- [govulncheck Documentation](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)
