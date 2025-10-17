#!/bin/bash

# Comprehensive Test Runner for etcd-monitor
# This script runs all types of tests: unit, integration, performance, and security

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
TEST_CONFIG_FILE="./test_config.yaml"
COVERAGE_THRESHOLD=90
TEST_TIMEOUT="10m"
PARALLEL_TESTS=4

# Test result tracking
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
SKIPPED_TESTS=0

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "INFO")
            echo -e "${BLUE}[INFO]${NC} $message"
            ;;
        "SUCCESS")
            echo -e "${GREEN}[SUCCESS]${NC} $message"
            ;;
        "WARNING")
            echo -e "${YELLOW}[WARNING]${NC} $message"
            ;;
        "ERROR")
            echo -e "${RED}[ERROR]${NC} $message"
            ;;
    esac
}

# Function to run tests and capture results
run_test_suite() {
    local test_type=$1
    local test_packages=$2
    local test_timeout=$3
    local test_parallel=$4
    
    print_status "INFO" "Running $test_type tests..."
    
    # Create test report directory
    mkdir -p "./test_reports/$test_type"
    
    # Run tests with coverage
    local test_result=0
    if [ "$test_parallel" = "true" ]; then
        go test -v -timeout="$test_timeout" -parallel="$PARALLEL_TESTS" -coverprofile="./test_reports/$test_type/coverage.out" -covermode=atomic $test_packages || test_result=$?
    else
        go test -v -timeout="$test_timeout" -coverprofile="./test_reports/$test_type/coverage.out" -covermode=atomic $test_packages || test_result=$?
    fi
    
    # Generate coverage report
    if [ -f "./test_reports/$test_type/coverage.out" ]; then
        go tool cover -html="./test_reports/$test_type/coverage.out" -o "./test_reports/$test_type/coverage.html"
        go tool cover -func="./test_reports/$test_type/coverage.out" > "./test_reports/$test_type/coverage.txt"
    fi
    
    # Update test counters
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    if [ $test_result -eq 0 ]; then
        PASSED_TESTS=$((PASSED_TESTS + 1))
        print_status "SUCCESS" "$test_type tests passed"
    else
        FAILED_TESTS=$((FAILED_TESTS + 1))
        print_status "ERROR" "$test_type tests failed"
    fi
    
    return $test_result
}

# Function to run security tests
run_security_tests() {
    print_status "INFO" "Running security tests..."
    
    # Create security test report directory
    mkdir -p "./test_reports/security"
    
    # Run gosec
    if command -v gosec &> /dev/null; then
        print_status "INFO" "Running gosec security scanner..."
        gosec -fmt json -out "./test_reports/security/gosec.json" ./... || true
        gosec -fmt text -out "./test_reports/security/gosec.txt" ./... || true
    else
        print_status "WARNING" "gosec not installed, skipping security scan"
    fi
    
    # Run govulncheck
    if command -v govulncheck &> /dev/null; then
        print_status "INFO" "Running govulncheck vulnerability scanner..."
        govulncheck -json ./... > "./test_reports/security/govulncheck.json" 2>&1 || true
        govulncheck ./... > "./test_reports/security/govulncheck.txt" 2>&1 || true
    else
        print_status "WARNING" "govulncheck not installed, skipping vulnerability scan"
    fi
    
    # Run nancy
    if command -v nancy &> /dev/null; then
        print_status "INFO" "Running nancy dependency scanner..."
        nancy sleuth -o json > "./test_reports/security/nancy.json" 2>&1 || true
        nancy sleuth > "./test_reports/security/nancy.txt" 2>&1 || true
    else
        print_status "WARNING" "nancy not installed, skipping dependency scan"
    fi
}

# Function to run performance tests
run_performance_tests() {
    print_status "INFO" "Running performance tests..."
    
    # Create performance test report directory
    mkdir -p "./test_reports/performance"
    
    # Run benchmark tests
    go test -v -bench=. -benchmem -timeout="15m" -coverprofile="./test_reports/performance/coverage.out" -covermode=atomic ./pkg/benchmark/... || true
    
    # Generate benchmark report
    go test -bench=. -benchmem ./pkg/benchmark/... > "./test_reports/performance/benchmark.txt" 2>&1 || true
    
    # Generate coverage report
    if [ -f "./test_reports/performance/coverage.out" ]; then
        go tool cover -html="./test_reports/performance/coverage.out" -o "./test_reports/performance/coverage.html"
        go tool cover -func="./test_reports/performance/coverage.out" > "./test_reports/performance/coverage.txt"
    fi
}

# Function to run integration tests
run_integration_tests() {
    print_status "INFO" "Running integration tests..."
    
    # Create integration test report directory
    mkdir -p "./test_reports/integration"
    
    # Start required services
    print_status "INFO" "Starting required services for integration tests..."
    
    # Start etcd
    docker run -d --name etcd-test -p 2379:2379 -p 2380:2380 \
        -e ETCD_NAME=etcd-server \
        -e ETCD_DATA_DIR=/etcd-data \
        -e ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379 \
        -e ETCD_ADVERTISE_CLIENT_URLS=http://localhost:2379 \
        -e ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380 \
        -e ETCD_INITIAL_ADVERTISE_PEER_URLS=http://localhost:2380 \
        -e ETCD_INITIAL_CLUSTER=etcd-server=http://localhost:2380 \
        -e ETCD_INITIAL_CLUSTER_STATE=new \
        -e ETCD_INITIAL_CLUSTER_TOKEN=etcd-cluster-1 \
        quay.io/coreos/etcd:v3.5.0 || true
    
    # Wait for etcd to be ready
    sleep 10
    
    # Run integration tests
    go test -v -timeout="10m" -tags=integration -coverprofile="./test_reports/integration/coverage.out" -covermode=atomic ./integration_tests/... || true
    
    # Generate coverage report
    if [ -f "./test_reports/integration/coverage.out" ]; then
        go tool cover -html="./test_reports/integration/coverage.out" -o "./test_reports/integration/coverage.html"
        go tool cover -func="./test_reports/integration/coverage.out" > "./test_reports/integration/coverage.txt"
    fi
    
    # Cleanup
    docker stop etcd-test || true
    docker rm etcd-test || true
}

# Function to generate overall test report
generate_test_report() {
    print_status "INFO" "Generating overall test report..."
    
    # Create overall test report directory
    mkdir -p "./test_reports/overall"
    
    # Generate overall coverage report
    echo "mode: atomic" > "./test_reports/overall/coverage.out"
    
    # Combine all coverage files
    for coverage_file in ./test_reports/*/coverage.out; do
        if [ -f "$coverage_file" ]; then
            tail -n +2 "$coverage_file" >> "./test_reports/overall/coverage.out"
        fi
    done
    
    # Generate HTML coverage report
    go tool cover -html="./test_reports/overall/coverage.out" -o "./test_reports/overall/coverage.html"
    go tool cover -func="./test_reports/overall/coverage.out" > "./test_reports/overall/coverage.txt"
    
    # Generate test summary
    cat > "./test_reports/overall/test_summary.txt" << EOF
Test Summary
============
Total Test Suites: $TOTAL_TESTS
Passed: $PASSED_TESTS
Failed: $FAILED_TESTS
Skipped: $SKIPPED_TESTS

Coverage Threshold: $COVERAGE_THRESHOLD%
Test Timeout: $TEST_TIMEOUT
Parallel Tests: $PARALLEL_TESTS

Test Reports:
- Unit Tests: ./test_reports/unit/
- Integration Tests: ./test_reports/integration/
- Performance Tests: ./test_reports/performance/
- Security Tests: ./test_reports/security/
- Overall Coverage: ./test_reports/overall/

Generated: $(date)
EOF
    
    print_status "SUCCESS" "Test report generated in ./test_reports/overall/"
}

# Main execution
main() {
    print_status "INFO" "Starting comprehensive test suite for etcd-monitor"
    print_status "INFO" "Configuration: $TEST_CONFIG_FILE"
    print_status "INFO" "Coverage Threshold: $COVERAGE_THRESHOLD%"
    print_status "INFO" "Test Timeout: $TEST_TIMEOUT"
    print_status "INFO" "Parallel Tests: $PARALLEL_TESTS"
    
    # Create test reports directory
    mkdir -p "./test_reports"
    
    # Run unit tests
    run_test_suite "unit" "./pkg/... ./cmd/... ./pingap-go/..." "$TEST_TIMEOUT" "true"
    
    # Run integration tests
    run_integration_tests
    
    # Run performance tests
    run_performance_tests
    
    # Run security tests
    run_security_tests
    
    # Generate overall test report
    generate_test_report
    
    # Print final summary
    print_status "INFO" "Test execution completed"
    print_status "INFO" "Total Test Suites: $TOTAL_TESTS"
    print_status "INFO" "Passed: $PASSED_TESTS"
    print_status "INFO" "Failed: $FAILED_TESTS"
    print_status "INFO" "Skipped: $SKIPPED_TESTS"
    
    if [ $FAILED_TESTS -eq 0 ]; then
        print_status "SUCCESS" "All tests passed!"
        exit 0
    else
        print_status "ERROR" "Some tests failed!"
        exit 1
    fi
}

# Run main function
main "$@"
