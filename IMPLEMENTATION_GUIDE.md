# etcd-monitor Full Implementation Guide

This guide provides a complete overview of the etcd monitoring system implementation, including architecture, key concepts, and usage examples based on etcd best practices.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Core Components](#core-components)
3. [Key Features](#key-features)
4. [Installation & Setup](#installation--setup)
5. [Usage Examples](#usage-examples)
6. [Monitoring Best Practices](#monitoring-best-practices)
7. [API Reference](#api-reference)
8. [Troubleshooting](#troubleshooting)

## Architecture Overview

The etcd-monitor system follows a microservices architecture with the following components:

```
┌─────────────────────────────────────────────────────────┐
│                     etcd-monitor                         │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  ┌────────────────┐      ┌─────────────────┐          │
│  │  Monitor       │      │  REST API       │          │
│  │  Service       │◄────►│  Server         │          │
│  └────────┬───────┘      └─────────────────┘          │
│           │                                             │
│  ┌────────▼───────┐  ┌──────────────┐  ┌───────────┐ │
│  │ Health         │  │  Metrics     │  │  Alert    │ │
│  │ Checker        │  │  Collector   │  │  Manager  │ │
│  └────────────────┘  └──────────────┘  └───────────┘ │
│                                                          │
│  ┌────────────────┐  ┌──────────────┐                 │
│  │  Benchmark     │  │  etcdctl     │                 │
│  │  Runner        │  │  Wrapper     │                 │
│  └────────────────┘  └──────────────┘                 │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
                   ┌───────────────┐
                   │  etcd Cluster │
                   └───────────────┘
```

## Core Components

### 1. Monitor Service (`pkg/monitor/service.go`)

The central orchestrator that:
- Manages the monitoring lifecycle
- Coordinates health checks and metrics collection
- Triggers alerts based on thresholds
- Watches for cluster events

**Key methods:**
- `Start()` - Begins monitoring
- `Stop()` - Gracefully shuts down
- `GetClusterStatus()` - Returns current cluster health
- `GetCurrentMetrics()` - Returns latest metrics snapshot

### 2. Health Checker (`pkg/monitor/health.go`)

Monitors cluster health including:
- Member availability and quorum status
- Leader election and stability
- Network partition detection
- Split-brain scenario detection
- Alarm monitoring

**Key features:**
- Tracks leader change history
- Detects when quorum is lost
- Measures network latency between members
- Identifies alarms (NOSPACE, CORRUPT, etc.)

### 3. Metrics Collector (`pkg/monitor/metrics.go`)

Collects comprehensive performance metrics:
- Request latency (read/write, P50/P95/P99)
- Throughput (operations per second)
- Database size and growth
- Raft consensus metrics
- Resource utilization

**Collection strategy:**
- Parallel collection for efficiency
- Statistical analysis with percentiles
- Historical data retention
- Performance baseline tracking

### 4. Alert Manager (`pkg/monitor/alert.go`)

Multi-channel alerting system supporting:
- **Email** - SMTP-based notifications
- **Slack** - Webhook integration
- **PagerDuty** - Incident management
- **Webhook** - Generic HTTP endpoints
- **Console** - Logging for development

**Features:**
- Alert deduplication (5-minute window)
- Severity levels (Info, Warning, Critical)
- Alert history and audit trail
- Configurable thresholds

### 5. REST API Server (`pkg/api/server.go`)

RESTful API for accessing monitoring data:

**Endpoints:**
- `GET /health` - API health check
- `GET /api/v1/cluster/status` - Cluster status
- `GET /api/v1/cluster/leader` - Leader information
- `GET /api/v1/metrics/current` - Current metrics
- `GET /api/v1/metrics/latency` - Latency metrics
- `POST /api/v1/performance/benchmark` - Run benchmark

### 6. Benchmark Runner (`pkg/benchmark/benchmark.go`)

Performance testing and validation:

**Benchmark types:**
- **Write** - Sequential write operations
- **Read** - Random read operations
- **Mixed** - 70% reads, 30% writes (realistic workload)

**Metrics collected:**
- Throughput (ops/sec)
- Latency distribution (P50/P95/P99)
- Error rates
- Latency histogram

### 7. etcdctl Wrapper (`pkg/etcdctl/wrapper.go`)

Programmatic interface for etcd operations:
- Key-value operations (Put, Get, Delete)
- Cluster management (MemberList, MemberAdd, MemberRemove)
- Maintenance operations (Compact, Defragment, Snapshot)
- Lease management
- Transactions
- Watches

## Key Features

### Health Monitoring

The system continuously monitors:

1. **Cluster Quorum**
   - Tracks available members
   - Calculates quorum size (N/2 + 1)
   - Detects when quorum is lost

2. **Leader Stability**
   - Monitors current leader
   - Tracks leader changes over time
   - Alerts on excessive leader changes (>3/hour)

3. **Network Health**
   - Measures inter-member latency
   - Detects network partitions
   - Identifies split-brain scenarios

4. **System Alarms**
   - NOSPACE - Out of disk space
   - CORRUPT - Data corruption detected
   - Custom alarms

### Performance Metrics

#### USE Method (System Resources)
- CPU utilization
- Memory usage
- Open file descriptors
- Storage space usage
- Network bandwidth

#### RED Method (Application Performance)
- Request Rate (ops/sec)
- Request Latency (P50/P95/P99)
- Error Rate (%)

#### etcd-Specific Metrics
- Proposal commit duration
- Disk WAL fsync duration
- Database backend commit duration
- Raft proposals (committed/applied/pending/failed)
- Database size and growth rate
- Number of watchers
- Connected clients

### Alert Thresholds

Default thresholds (configurable):
- Max latency: 100ms (P99)
- Max database size: 8GB
- Min available nodes: 2
- Max leader changes: 3/hour
- Max error rate: 5%
- Min disk space: 10%

## Installation & Setup

### Prerequisites

- Go 1.21 or later
- Access to an etcd cluster (v3.4+)
- PostgreSQL + TimescaleDB (optional, for historical data)
- Redis (optional, for alerting queue)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/your-org/etcd-monitor.git
cd etcd-monitor

# Download dependencies
go mod download

# Build the main application
go build -o etcd-monitor cmd/etcd-monitor/main.go

# Or use Make (if Makefile provided)
make build
```

### Configuration

1. Copy the example configuration:
```bash
cp config.example.yaml config.yaml
```

2. Edit `config.yaml` with your etcd endpoints and preferences

3. Set up alert channels (Slack, PagerDuty, etc.)

### Running

#### Standard Mode (Continuous Monitoring)

```bash
./etcd-monitor \
  --endpoints="localhost:2379,localhost:22379,localhost:32379" \
  --api-port=8080 \
  --health-check-interval=30s \
  --metrics-interval=10s
```

#### Benchmark Mode (One-time Test)

```bash
./etcd-monitor \
  --run-benchmark \
  --benchmark-type=mixed \
  --benchmark-ops=10000 \
  --endpoints="localhost:2379"
```

### Docker Deployment

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o etcd-monitor cmd/etcd-monitor/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/etcd-monitor .
COPY config.yaml .
CMD ["./etcd-monitor"]
```

```bash
docker build -t etcd-monitor:latest .
docker run -p 8080:8080 -v ./config.yaml:/root/config.yaml etcd-monitor:latest
```

## Usage Examples

### Monitoring a Cluster

```go
package main

import (
    "context"
    "time"

    "github.com/etcd-monitor/taskmaster/pkg/monitor"
    "go.uber.org/zap"
)

func main() {
    logger, _ := zap.NewProduction()

    config := &monitor.Config{
        Endpoints: []string{"localhost:2379"},
        DialTimeout: 5 * time.Second,
        HealthCheckInterval: 30 * time.Second,
        MetricsInterval: 10 * time.Second,
        AlertThresholds: monitor.AlertThresholds{
            MaxLatencyMs: 100,
            MaxDatabaseSizeMB: 8192,
        },
    }

    service, _ := monitor.NewMonitorService(config, logger)
    service.Start()
    defer service.Stop()

    // Monitor runs in background
    time.Sleep(1 * time.Hour)
}
```

### Running a Benchmark

```go
package main

import (
    "context"

    "github.com/etcd-monitor/taskmaster/pkg/benchmark"
    clientv3 "go.etcd.io/etcd/client/v3"
    "go.uber.org/zap"
)

func main() {
    logger, _ := zap.NewProduction()
    client, _ := clientv3.New(clientv3.Config{
        Endpoints: []string{"localhost:2379"},
    })

    config := &benchmark.Config{
        Type: benchmark.BenchmarkTypeMixed,
        Connections: 10,
        Clients: 10,
        TotalOperations: 10000,
        KeyPrefix: "/bench",
    }

    runner := benchmark.NewRunner(client, config, logger)
    result, _ := runner.Run(context.Background())

    benchmark.PrintResult(result, logger)
}
```

### Using etcdctl Wrapper

```go
package main

import (
    "context"

    "github.com/etcd-monitor/taskmaster/pkg/etcdctl"
    clientv3 "go.etcd.io/etcd/client/v3"
    "go.uber.org/zap"
)

func main() {
    logger, _ := zap.NewProduction()
    client, _ := clientv3.New(clientv3.Config{
        Endpoints: []string{"localhost:2379"},
    })

    wrapper := etcdctl.NewWrapper(client, logger)

    // Put a key
    wrapper.Put(context.Background(), "/mykey", "myvalue")

    // Get a key
    resp, _ := wrapper.Get(context.Background(), "/mykey")

    // Check endpoint health
    health, _ := wrapper.EndpointHealth(context.Background())
    for _, h := range health {
        logger.Info("Endpoint",
            zap.String("endpoint", h.Endpoint),
            zap.Bool("healthy", h.Healthy))
    }
}
```

## Monitoring Best Practices

### 1. Health Check Frequency

- **Production**: 30-60 seconds
- **Development**: 10-30 seconds
- **Critical systems**: 10-15 seconds

### 2. Metrics Collection

- **Standard**: Every 10 seconds
- **High-frequency**: Every 5 seconds (increased overhead)
- **Low-frequency**: Every 30 seconds (may miss spikes)

### 3. Alert Threshold Guidelines

#### Latency
- **Warning**: P99 > 50ms
- **Critical**: P99 > 100ms

#### Database Size
- **Warning**: > 6GB
- **Critical**: > 8GB (default limit)

#### Leader Changes
- **Warning**: > 2/hour
- **Critical**: > 3/hour

### 4. Benchmark Recommendations

Run benchmarks:
- Before production deployment
- After major configuration changes
- Monthly for capacity planning
- During maintenance windows

Expected performance (3-node cluster, SSD):
- **Write throughput**: 10,000-20,000 ops/sec
- **Read throughput**: 30,000-50,000 ops/sec
- **P99 latency**: < 100ms

### 5. Disk Performance Requirements

etcd requires fast disks:
- **Minimum**: 99th percentile fsync < 10ms
- **Recommended**: 99th percentile fsync < 5ms
- Use SSD or NVMe storage
- Avoid network-attached storage

## API Reference

### GET /api/v1/cluster/status

Returns current cluster status.

**Response:**
```json
{
  "healthy": true,
  "leader_id": 123456789,
  "member_count": 3,
  "quorum_size": 2,
  "has_leader": true,
  "leader_changes": 1,
  "last_leader_change": "2025-10-11T10:30:00Z",
  "network_partition": false,
  "alarms": [],
  "last_check": "2025-10-11T12:45:00Z"
}
```

### GET /api/v1/metrics/current

Returns current metrics snapshot.

**Response:**
```json
{
  "timestamp": "2025-10-11T12:45:00Z",
  "request_rate": 1250.5,
  "read_latency_p50": 2.5,
  "read_latency_p95": 8.2,
  "read_latency_p99": 15.3,
  "write_latency_p50": 5.1,
  "write_latency_p95": 12.4,
  "write_latency_p99": 25.7,
  "db_size": 104857600,
  "db_size_in_use": 95234816
}
```

## Troubleshooting

### High Latency

**Symptoms:**
- P99 latency > 100ms
- Slow API responses

**Causes:**
1. Disk I/O bottleneck
2. Network congestion
3. Resource exhaustion
4. Large database size

**Solutions:**
- Check disk performance with `fio`
- Compact and defragment database
- Increase resources (CPU, RAM, disk)
- Review network configuration

### Leader Changes

**Symptoms:**
- Frequent leader elections
- Service disruptions

**Causes:**
1. Network instability
2. Resource pressure on leader
3. Clock skew
4. Heartbeat timeouts

**Solutions:**
- Check network latency between members
- Verify system clocks are synchronized (NTP)
- Increase election timeout
- Add resources to leader node

### Database Size Growth

**Symptoms:**
- Rapid database growth
- NOSPACE alarms

**Causes:**
1. No compaction configured
2. Long revision history
3. Many watchers keeping old revisions

**Solutions:**
- Enable auto-compaction
- Run manual compaction regularly
- Defragment after compaction
- Review data retention policies

## Additional Resources

- [etcd Documentation](https://etcd.io/docs/)
- [etcd Metrics Guide](https://etcd.io/docs/latest/metrics/)
- [etcd Performance Benchmarking](https://etcd.io/docs/latest/op-guide/performance/)
- [Kubernetes etcd Best Practices](https://kubernetes.io/docs/tasks/administer-cluster/configure-upgrade-etcd/)

## Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

## License

[License Type] - See LICENSE file for details
