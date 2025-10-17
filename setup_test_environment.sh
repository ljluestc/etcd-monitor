#!/bin/bash

# Test Environment Setup Script for etcd-monitor
# This script sets up the complete test environment including dependencies

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
ETCD_VERSION="v3.5.0"
PROMETHEUS_VERSION="v2.40.0"
GRAFANA_VERSION="9.3.0"
DOCKER_COMPOSE_FILE="docker-compose-test.yml"

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

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to install Go tools
install_go_tools() {
    print_status "INFO" "Installing Go tools..."
    
    # Install golangci-lint
    if ! command_exists golangci-lint; then
        print_status "INFO" "Installing golangci-lint..."
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
    else
        print_status "INFO" "golangci-lint already installed"
    fi
    
    # Install gosec
    if ! command_exists gosec; then
        print_status "INFO" "Installing gosec..."
        go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
    else
        print_status "INFO" "gosec already installed"
    fi
    
    # Install govulncheck
    if ! command_exists govulncheck; then
        print_status "INFO" "Installing govulncheck..."
        go install golang.org/x/vuln/cmd/govulncheck@latest
    else
        print_status "INFO" "govulncheck already installed"
    fi
    
    # Install nancy
    if ! command_exists nancy; then
        print_status "INFO" "Installing nancy..."
        go install github.com/sonatypecommunity/nancy@latest
    else
        print_status "INFO" "nancy already installed"
    fi
    
    # Install pre-commit
    if ! command_exists pre-commit; then
        print_status "INFO" "Installing pre-commit..."
        pip install pre-commit
    else
        print_status "INFO" "pre-commit already installed"
    fi
}

# Function to create test docker-compose file
create_test_docker_compose() {
    print_status "INFO" "Creating test docker-compose file..."
    
    cat > "$DOCKER_COMPOSE_FILE" << EOF
version: '3.8'

services:
  etcd:
    image: quay.io/coreos/etcd:${ETCD_VERSION}
    container_name: etcd-test
    ports:
      - "2379:2379"
      - "2380:2380"
    environment:
      - ETCD_NAME=etcd-server
      - ETCD_DATA_DIR=/etcd-data
      - ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379
      - ETCD_ADVERTISE_CLIENT_URLS=http://localhost:2379
      - ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380
      - ETCD_INITIAL_ADVERTISE_PEER_URLS=http://localhost:2380
      - ETCD_INITIAL_CLUSTER=etcd-server=http://localhost:2380
      - ETCD_INITIAL_CLUSTER_STATE=new
      - ETCD_INITIAL_CLUSTER_TOKEN=etcd-cluster-1
    volumes:
      - etcd-data:/etcd-data
    healthcheck:
      test: ["CMD", "etcdctl", "endpoint", "health"]
      interval: 10s
      timeout: 5s
      retries: 3
      start_period: 10s

  prometheus:
    image: prom/prometheus:${PROMETHEUS_VERSION}
    container_name: prometheus-test
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus-data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--storage.tsdb.retention.time=200h'
      - '--web.enable-lifecycle'
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:9090/-/healthy"]
      interval: 10s
      timeout: 5s
      retries: 3
      start_period: 10s

  grafana:
    image: grafana/grafana:${GRAFANA_VERSION}
    container_name: grafana-test
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_USERS_ALLOW_SIGN_UP=false
    volumes:
      - grafana-data:/var/lib/grafana
      - ./grafana/datasources:/etc/grafana/provisioning/datasources
      - ./grafana/dashboards:/etc/grafana/provisioning/dashboards
    healthcheck:
      test: ["CMD-SHELL", "curl -f http://localhost:3000/api/health || exit 1"]
      interval: 10s
      timeout: 5s
      retries: 3
      start_period: 10s

volumes:
  etcd-data:
  prometheus-data:
  grafana-data:
EOF

    print_status "SUCCESS" "Test docker-compose file created: $DOCKER_COMPOSE_FILE"
}

# Function to start test services
start_test_services() {
    print_status "INFO" "Starting test services..."
    
    # Start services
    docker-compose -f "$DOCKER_COMPOSE_FILE" up -d
    
    # Wait for services to be ready
    print_status "INFO" "Waiting for services to be ready..."
    
    # Wait for etcd
    print_status "INFO" "Waiting for etcd..."
    timeout 60 bash -c 'until docker exec etcd-test etcdctl endpoint health; do sleep 2; done'
    
    # Wait for prometheus
    print_status "INFO" "Waiting for prometheus..."
    timeout 60 bash -c 'until curl -f http://localhost:9090/-/healthy; do sleep 2; done'
    
    # Wait for grafana
    print_status "INFO" "Waiting for grafana..."
    timeout 60 bash -c 'until curl -f http://localhost:3000/api/health; do sleep 2; done'
    
    print_status "SUCCESS" "All test services are ready"
}

# Function to stop test services
stop_test_services() {
    print_status "INFO" "Stopping test services..."
    docker-compose -f "$DOCKER_COMPOSE_FILE" down -v
    print_status "SUCCESS" "Test services stopped"
}

# Function to setup pre-commit hooks
setup_pre_commit_hooks() {
    print_status "INFO" "Setting up pre-commit hooks..."
    
    if [ -f ".pre-commit-config.yaml" ]; then
        pre-commit install
        print_status "SUCCESS" "Pre-commit hooks installed"
    else
        print_status "WARNING" "Pre-commit config file not found"
    fi
}

# Function to run tests
run_tests() {
    print_status "INFO" "Running tests..."
    
    # Run unit tests
    print_status "INFO" "Running unit tests..."
    go test -v -race -coverprofile=coverage.out ./...
    
    # Run integration tests
    print_status "INFO" "Running integration tests..."
    go test -v -race -tags=integration ./integration_tests/...
    
    # Generate coverage report
    print_status "INFO" "Generating coverage report..."
    go tool cover -html=coverage.out -o coverage.html
    go tool cover -func=coverage.out > coverage.txt
    
    print_status "SUCCESS" "Tests completed"
}

# Function to run security scans
run_security_scans() {
    print_status "INFO" "Running security scans..."
    
    # Run gosec
    if command_exists gosec; then
        print_status "INFO" "Running gosec..."
        gosec ./...
    fi
    
    # Run govulncheck
    if command_exists govulncheck; then
        print_status "INFO" "Running govulncheck..."
        govulncheck ./...
    fi
    
    # Run nancy
    if command_exists nancy; then
        print_status "INFO" "Running nancy..."
        nancy sleuth
    fi
    
    print_status "SUCCESS" "Security scans completed"
}

# Function to run linting
run_linting() {
    print_status "INFO" "Running linting..."
    
    # Run golangci-lint
    if command_exists golangci-lint; then
        print_status "INFO" "Running golangci-lint..."
        golangci-lint run
    fi
    
    # Run go vet
    print_status "INFO" "Running go vet..."
    go vet ./...
    
    # Run go fmt
    print_status "INFO" "Running go fmt..."
    go fmt ./...
    
    print_status "SUCCESS" "Linting completed"
}

# Function to cleanup
cleanup() {
    print_status "INFO" "Cleaning up..."
    
    # Stop test services
    stop_test_services
    
    # Remove test docker-compose file
    if [ -f "$DOCKER_COMPOSE_FILE" ]; then
        rm "$DOCKER_COMPOSE_FILE"
    fi
    
    print_status "SUCCESS" "Cleanup completed"
}

# Main execution
main() {
    print_status "INFO" "Setting up test environment for etcd-monitor"
    
    # Check if Docker is available
    if ! command_exists docker; then
        print_status "ERROR" "Docker is required but not installed"
        exit 1
    fi
    
    # Check if docker-compose is available
    if ! command_exists docker-compose; then
        print_status "ERROR" "docker-compose is required but not installed"
        exit 1
    fi
    
    # Install Go tools
    install_go_tools
    
    # Create test docker-compose file
    create_test_docker_compose
    
    # Setup pre-commit hooks
    setup_pre_commit_hooks
    
    # Start test services
    start_test_services
    
    # Run tests
    run_tests
    
    # Run security scans
    run_security_scans
    
    # Run linting
    run_linting
    
    print_status "SUCCESS" "Test environment setup completed"
    print_status "INFO" "Test services are running:"
    print_status "INFO" "  - etcd: http://localhost:2379"
    print_status "INFO" "  - prometheus: http://localhost:9090"
    print_status "INFO" "  - grafana: http://localhost:3000 (admin/admin)"
    print_status "INFO" "Coverage report: coverage.html"
    print_status "INFO" "To stop test services, run: docker-compose -f $DOCKER_COMPOSE_FILE down -v"
}

# Handle script arguments
case "${1:-}" in
    "start")
        start_test_services
        ;;
    "stop")
        stop_test_services
        ;;
    "test")
        run_tests
        ;;
    "security")
        run_security_scans
        ;;
    "lint")
        run_linting
        ;;
    "cleanup")
        cleanup
        ;;
    *)
        main
        ;;
esac
