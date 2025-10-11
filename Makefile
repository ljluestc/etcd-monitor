# Makefile for etcd-monitor

# Variables
BINARY_NAME=etcd-monitor
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe
VERSION=1.0.0
BUILD_DIR=bin
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

.PHONY: all build clean test coverage lint run help deps install docker-build docker-run benchmark

# Default target
all: clean deps build test

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) cmd/etcd-monitor/main.go
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_UNIX) cmd/etcd-monitor/main.go
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_UNIX)"

# Build for Windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_WINDOWS) cmd/etcd-monitor/main.go
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_WINDOWS)"

# Build for all platforms
build-all: build-linux build-windows build

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out
	@$(GO) clean
	@echo "Clean complete"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod verify
	@echo "Dependencies downloaded"

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo "Dependencies updated"

# Run tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Installing..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...
	@echo "Code formatted"

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME) --endpoints=localhost:2379 --api-port=8080

# Run with custom endpoints
run-custom:
	@echo "Running $(BINARY_NAME) with custom configuration..."
	./$(BUILD_DIR)/$(BINARY_NAME) \
		--endpoints=$(ENDPOINTS) \
		--api-port=$(API_PORT) \
		--health-check-interval=$(HEALTH_INTERVAL) \
		--metrics-interval=$(METRICS_INTERVAL)

# Run benchmark
benchmark: build
	@echo "Running benchmark..."
	./$(BUILD_DIR)/$(BINARY_NAME) \
		--run-benchmark \
		--benchmark-type=mixed \
		--benchmark-ops=10000 \
		--endpoints=localhost:2379

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete"

# Uninstall the binary
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstall complete"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t etcd-monitor:$(VERSION) .
	docker tag etcd-monitor:$(VERSION) etcd-monitor:latest
	@echo "Docker image built: etcd-monitor:$(VERSION)"

# Run in Docker
docker-run:
	@echo "Running Docker container..."
	docker run -d \
		--name etcd-monitor \
		-p 8080:8080 \
		-v $(PWD)/config.yaml:/root/config.yaml \
		etcd-monitor:latest

# Stop Docker container
docker-stop:
	@echo "Stopping Docker container..."
	docker stop etcd-monitor
	docker rm etcd-monitor

# Run development environment with docker-compose
dev:
	@echo "Starting development environment..."
	docker-compose up -d

# Stop development environment
dev-stop:
	@echo "Stopping development environment..."
	docker-compose down

# Show help
help:
	@echo "etcd-monitor Makefile commands:"
	@echo ""
	@echo "  make build          - Build the application"
	@echo "  make build-all      - Build for all platforms"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make deps           - Download dependencies"
	@echo "  make deps-update    - Update dependencies"
	@echo "  make test           - Run tests"
	@echo "  make coverage       - Run tests with coverage"
	@echo "  make lint           - Run linter"
	@echo "  make fmt            - Format code"
	@echo "  make run            - Build and run the application"
	@echo "  make benchmark      - Run performance benchmark"
	@echo "  make install        - Install binary to /usr/local/bin"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-run     - Run in Docker container"
	@echo "  make dev            - Start development environment"
	@echo "  make help           - Show this help message"
	@echo ""

# Build TaskMaster CLI
build-taskmaster:
	@echo "Building TaskMaster CLI..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/taskmaster taskmaster.go
	@echo "TaskMaster built: $(BUILD_DIR)/taskmaster"

# Build examples
build-examples:
	@echo "Building examples..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/examples cmd/examples/main.go
	@echo "Examples built: $(BUILD_DIR)/examples"

# Build all applications
build-all-apps: build build-taskmaster build-examples
