#!/bin/bash

# Comprehensive Test Runner for etcd-monitor
# This script runs all tests and generates coverage reports

set -e

echo "🚀 Starting Comprehensive Test Suite for etcd-monitor"
echo "=================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test results tracking
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Function to run tests for a package
run_package_tests() {
    local package_path=$1
    local package_name=$2
    
    echo -e "\n${BLUE}Testing package: $package_name${NC}"
    echo "----------------------------------------"
    
    if [ -d "$package_path" ]; then
        cd "$package_path"
        
        # Run tests with coverage
        if go test -v -coverprofile=coverage.out -covermode=atomic ./... 2>&1 | tee test_output.log; then
            echo -e "${GREEN}✅ $package_name tests PASSED${NC}"
            ((PASSED_TESTS++))
            
            # Generate coverage report
            if [ -f "coverage.out" ]; then
                go tool cover -html=coverage.out -o coverage.html
                echo "📊 Coverage report generated: $package_path/coverage.html"
            fi
        else
            echo -e "${RED}❌ $package_name tests FAILED${NC}"
            ((FAILED_TESTS++))
        fi
        
        cd - > /dev/null
    else
        echo -e "${YELLOW}⚠️  Package directory not found: $package_path${NC}"
    fi
    
    ((TOTAL_TESTS++))
}

# Function to run integration tests
run_integration_tests() {
    echo -e "\n${BLUE}Running Integration Tests${NC}"
    echo "=============================="
    
    # Create integration test directory if it doesn't exist
    mkdir -p integration_tests
    
    # Run integration tests
    if [ -d "integration_tests" ]; then
        cd integration_tests
        if go test -v -tags=integration ./... 2>&1 | tee integration_test_output.log; then
            echo -e "${GREEN}✅ Integration tests PASSED${NC}"
            ((PASSED_TESTS++))
        else
            echo -e "${RED}❌ Integration tests FAILED${NC}"
            ((FAILED_TESTS++))
        fi
        cd - > /dev/null
    else
        echo -e "${YELLOW}⚠️  No integration tests found${NC}"
    fi
    
    ((TOTAL_TESTS++))
}

# Function to run performance tests
run_performance_tests() {
    echo -e "\n${BLUE}Running Performance Tests${NC}"
    echo "=============================="
    
    # Run performance tests with benchmarks
    if go test -v -bench=. -benchmem -run=^$ ./... 2>&1 | tee performance_test_output.log; then
        echo -e "${GREEN}✅ Performance tests PASSED${NC}"
        ((PASSED_TESTS++))
    else
        echo -e "${RED}❌ Performance tests FAILED${NC}"
        ((FAILED_TESTS++))
    fi
    
    ((TOTAL_TESTS++))
}

# Function to generate overall coverage report
generate_overall_coverage() {
    echo -e "\n${BLUE}Generating Overall Coverage Report${NC}"
    echo "====================================="
    
    # Create coverage directory
    mkdir -p coverage_reports
    
    # Find all coverage.out files
    find . -name "coverage.out" -type f | while read -r coverage_file; do
        echo "Found coverage file: $coverage_file"
    done
    
    # Generate overall coverage
    echo "mode: atomic" > coverage_reports/overall_coverage.out
    find . -name "coverage.out" -type f -exec grep -h "^[^m]" {} \; >> coverage_reports/overall_coverage.out
    
    # Generate HTML report
    go tool cover -html=coverage_reports/overall_coverage.out -o coverage_reports/overall_coverage.html
    
    echo "📊 Overall coverage report generated: coverage_reports/overall_coverage.html"
}

# Function to run security tests
run_security_tests() {
    echo -e "\n${BLUE}Running Security Tests${NC}"
    echo "========================"
    
    # Check for security vulnerabilities
    if command -v gosec &> /dev/null; then
        echo "Running gosec security scanner..."
        gosec ./... 2>&1 | tee security_test_output.log
    else
        echo -e "${YELLOW}⚠️  gosec not installed, skipping security tests${NC}"
    fi
    
    # Check for dependency vulnerabilities
    if command -v govulncheck &> /dev/null; then
        echo "Running govulncheck vulnerability scanner..."
        govulncheck ./... 2>&1 | tee vulnerability_test_output.log
    else
        echo -e "${YELLOW}⚠️  govulncheck not installed, skipping vulnerability tests${NC}"
    fi
}

# Function to run linting
run_linting() {
    echo -e "\n${BLUE}Running Linting${NC}"
    echo "==============="
    
    # Run go vet
    if go vet ./... 2>&1 | tee linting_output.log; then
        echo -e "${GREEN}✅ go vet PASSED${NC}"
    else
        echo -e "${RED}❌ go vet FAILED${NC}"
    fi
    
    # Run go fmt check
    if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
        echo -e "${RED}❌ go fmt check FAILED - files need formatting${NC}"
        gofmt -s -l .
    else
        echo -e "${GREEN}✅ go fmt check PASSED${NC}"
    fi
    
    # Run golangci-lint if available
    if command -v golangci-lint &> /dev/null; then
        echo "Running golangci-lint..."
        golangci-lint run ./... 2>&1 | tee golangci_lint_output.log
    else
        echo -e "${YELLOW}⚠️  golangci-lint not installed, skipping advanced linting${NC}"
    fi
}

# Main execution
main() {
    echo "Starting test execution at $(date)"
    echo "=================================="
    
    # Run package tests
    run_package_tests "pkg/controllers/util" "controllers/util"
    run_package_tests "pkg/etcd" "etcd"
    run_package_tests "pkg/etcdctl" "etcdctl"
    run_package_tests "pkg/k8s" "k8s"
    run_package_tests "pkg/signals" "signals"
    run_package_tests "pkg/benchmark" "benchmark"
    run_package_tests "pkg/examples" "examples"
    run_package_tests "pkg/patterns" "patterns"
    run_package_tests "pkg/inspection" "inspection"
    run_package_tests "pkg/featureprovider" "featureprovider"
    run_package_tests "pkg/clusterprovider" "clusterprovider"
    run_package_tests "pkg/monitor" "monitor"
    run_package_tests "pkg/api" "api"
    
    # Run cmd tests
    run_package_tests "cmd/main" "main"
    run_package_tests "cmd/etcd-monitor" "etcd-monitor"
    run_package_tests "cmd/etcdcluster-controller" "etcdcluster-controller"
    run_package_tests "cmd/etcdinspection-controller" "etcdinspection-controller"
    run_package_tests "cmd/examples" "examples"
    
    # Run pingap-go tests
    run_package_tests "pingap-go" "pingap-go"
    run_package_tests "pingap-go/pkg/config" "pingap-go/config"
    run_package_tests "pingap-go/pkg/plugin" "pingap-go/plugin"
    run_package_tests "pingap-go/pkg/proxy" "pingap-go/proxy"
    run_package_tests "pingap-go/pkg/router" "pingap-go/router"
    run_package_tests "pingap-go/pkg/upstream" "pingap-go/upstream"
    run_package_tests "pingap-go/cmd/pingap" "pingap-go/cmd/pingap"
    
    # Run integration tests
    run_integration_tests
    
    # Run performance tests
    run_performance_tests
    
    # Run security tests
    run_security_tests
    
    # Run linting
    run_linting
    
    # Generate overall coverage report
    generate_overall_coverage
    
    # Print summary
    echo -e "\n${BLUE}Test Summary${NC}"
    echo "============"
    echo "Total test suites: $TOTAL_TESTS"
    echo -e "Passed: ${GREEN}$PASSED_TESTS${NC}"
    echo -e "Failed: ${RED}$FAILED_TESTS${NC}"
    
    if [ $FAILED_TESTS -eq 0 ]; then
        echo -e "\n${GREEN}🎉 All tests PASSED!${NC}"
        exit 0
    else
        echo -e "\n${RED}💥 Some tests FAILED!${NC}"
        exit 1
    fi
}

# Run main function
main "$@"
