# etcd-monitor Complete Quickstart Guide

## Overview

etcd-monitor is a comprehensive monitoring, benchmarking, and management system for etcd clusters based on the [clay-wangzhi/etcd-metrics](https://github.com/clay-wangzhi/etcd-metrics) reference architecture and industry best practices.

## Features Implemented

### ✅ Phase 1: Foundation (Complete)
- [x] TLS/SSL connection support for etcd
- [x] Cluster health monitoring with leader tracking
- [x] Split-brain detection
- [x] Core metrics collection (USE + RED methods)
- [x] Comprehensive alerting system
- [x] Multiple alert channels (Email, Slack, PagerDuty, Webhook, Console)
- [x] Basic API server with health and metrics endpoints

### ✅ Phase 2: Enhanced Monitoring (Complete)
- [x] Advanced metrics collection (latency percentiles, throughput, Raft metrics)
- [x] Prometheus metrics exporter
- [x] Comprehensive benchmark testing (write, read, mixed workloads)
- [x] etcdctl command wrapper for administrative operations
- [x] Historical metrics tracking
- [x] Performance measurement and reporting

### ✅ Deployment & Visualization (Complete)
- [x] Docker Compose configuration with full stack
- [x] Grafana dashboard templates
- [x] Prometheus configuration
- [x] 3-node etcd cluster setup

## Quick Start

### Prerequisites

- Docker and Docker Compose installed
- Go 1.19+ (for local development)
- 8GB RAM minimum (for running full stack)

### 1. Clone and Build

```bash
cd etcd-monitor

# Build the etcd-monitor binary
go build -o etcd-monitor ./cmd/etcd-monitor

# Or use Docker Compose (recommended)
docker-compose -f docker-compose-full.yml up -d
```

### 2. Access the Services

After starting with Docker Compose, access:

- **etcd-monitor API**: http://localhost:8080
  - Health check: http://localhost:8080/health
  - Cluster status: http://localhost:8080/api/v1/cluster/status
  - Current metrics: http://localhost:8080/api/v1/metrics/current
  - Prometheus metrics: http://localhost:8080/metrics

- **Prometheus**: http://localhost:9090
  - Query interface for metrics
  - Alerts and rules management

- **Grafana**: http://localhost:3000
  - Username: `admin`
  - Password: `admin`
  - Pre-configured etcd dashboard available

- **etcd Cluster**:
  - Node 1: http://localhost:2379
  - Node 2: http://localhost:22379
  - Node 3: http://localhost:32379

### 3. Verify Cluster Health

```bash
# Check cluster status
curl http://localhost:8080/api/v1/cluster/status | jq

# Check current metrics
curl http://localhost:8080/api/v1/metrics/current | jq

# Check Prometheus metrics
curl http://localhost:8080/metrics
```

### 4. Run Benchmarks

```bash
# Run write benchmark
./etcd-monitor --endpoints=localhost:2379 --run-benchmark --benchmark-type=write --benchmark-ops=10000

# Run read benchmark
./etcd-monitor --endpoints=localhost:2379 --run-benchmark --benchmark-type=read --benchmark-ops=10000

# Run mixed benchmark
./etcd-monitor --endpoints=localhost:2379 --run-benchmark --benchmark-type=mixed --benchmark-ops=10000
```

## Configuration

### Environment Variables

```bash
# etcd connection
ENDPOINTS=etcd1:2379,etcd2:2379,etcd3:2379
DIAL_TIMEOUT=5s

# TLS (optional)
CERT_FILE=/path/to/cert.pem
KEY_FILE=/path/to/key.pem
CA_FILE=/path/to/ca.pem
INSECURE_SKIP_TLS=false

# API server
API_HOST=0.0.0.0
API_PORT=8080

# Monitoring intervals
HEALTH_CHECK_INTERVAL=30s
METRICS_INTERVAL=10s

# Alert thresholds
MAX_LATENCY_MS=100
MAX_DB_SIZE_MB=8192
MIN_AVAILABLE_NODES=2
MAX_LEADER_CHANGES=3

# Benchmark
BENCHMARK_ENABLED=false
BENCHMARK_INTERVAL=1h
```

### Command Line Flags

```bash
./etcd-monitor \
  --endpoints=localhost:2379,localhost:22379,localhost:32379 \
  --dial-timeout=5s \
  --api-port=8080 \
  --health-check-interval=30s \
  --metrics-interval=10s \
  --max-latency-ms=100 \
  --max-db-size-mb=8192 \
  --min-available-nodes=2 \
  --max-leader-changes=3
```

## API Endpoints

### Health & Status

- `GET /health` - API health check
- `GET /api/v1/cluster/status` - Cluster health status
- `GET /api/v1/cluster/members` - Cluster member information
- `GET /api/v1/cluster/leader` - Current leader information

### Metrics

- `GET /api/v1/metrics/current` - Current metrics snapshot
- `GET /api/v1/metrics/history` - Historical metrics (planned)
- `GET /api/v1/metrics/latency` - Latency-specific metrics
- `GET /metrics` - Prometheus metrics endpoint

### Alerts

- `GET /api/v1/alerts` - Current active alerts
- `GET /api/v1/alerts/history` - Alert history

### Performance

- `POST /api/v1/performance/benchmark` - Run performance benchmark

## Monitoring Metrics

### Cluster Health Metrics
- `etcd_cluster_healthy` - Cluster health status (0/1)
- `etcd_cluster_has_leader` - Leader presence (0/1)
- `etcd_cluster_member_count` - Number of cluster members
- `etcd_cluster_quorum_size` - Required quorum size
- `etcd_cluster_leader_changes_total` - Total leader changes

### Performance Metrics
- `etcd_request_read_latency_p{50,95,99}_milliseconds` - Read latency percentiles
- `etcd_request_write_latency_p{50,95,99}_milliseconds` - Write latency percentiles
- `etcd_request_rate_per_second` - Request rate (ops/sec)

### Database Metrics
- `etcd_mvcc_db_total_size_bytes` - Total database size
- `etcd_mvcc_db_total_size_in_use_bytes` - Database size in use

### Raft Consensus Metrics
- `etcd_server_proposals_committed_total` - Total proposals committed
- `etcd_server_proposals_applied_total` - Total proposals applied
- `etcd_server_proposals_pending` - Current pending proposals
- `etcd_server_proposals_failed_total` - Total failed proposals

### Disk Performance Metrics
- `etcd_disk_wal_fsync_duration_p95_milliseconds` - WAL fsync latency P95
- `etcd_disk_backend_commit_duration_p95_milliseconds` - Backend commit latency P95

## Alert Thresholds

### Critical Alerts
- No leader elected
- Leader election frequency > threshold
- Disk fsync latency P99 > 100ms
- Database size approaching quota
- Member unavailable

### Warning Alerts
- High proposal failure rate
- Disk usage > 80%
- Memory usage > 80%
- Request latency degradation
- High pending proposals

## Grafana Dashboards

The pre-configured Grafana dashboard includes:

1. **Cluster Overview Panel**
   - Health status
   - Leader status
   - Member count

2. **Performance Panel**
   - Read/write latency (P50, P95, P99)
   - Request rate
   - Throughput

3. **Database Panel**
   - Database size
   - Size in use
   - Growth rate

4. **Raft Consensus Panel**
   - Proposals committed
   - Proposals applied
   - Proposals pending
   - Proposals failed

5. **Disk Performance Panel**
   - WAL fsync duration
   - Backend commit duration

## Benchmark Testing

### Write Performance Test
```bash
./etcd-monitor --endpoints=localhost:2379 \
  --run-benchmark \
  --benchmark-type=write \
  --benchmark-ops=10000
```

Expected results (etcd 3.5, 8vCPU, 16GB RAM, SSD):
- Throughput: ~20,000 ops/sec
- Latency P99: < 100ms

### Read Performance Test
```bash
./etcd-monitor --endpoints=localhost:2379 \
  --run-benchmark \
  --benchmark-type=read \
  --benchmark-ops=10000
```

Expected results:
- Throughput: ~40,000 ops/sec
- Latency P99: < 50ms

### Mixed Workload Test
```bash
./etcd-monitor --endpoints=localhost:2379 \
  --run-benchmark \
  --benchmark-type=mixed \
  --benchmark-ops=10000
```

70% reads, 30% writes workload.

## Troubleshooting

### etcd-monitor cannot connect to etcd

```bash
# Check etcd is running
docker ps | grep etcd

# Check etcd endpoints
curl http://localhost:2379/health

# Check network connectivity
docker network inspect etcd-monitor_etcd-network
```

### High latency detected

```bash
# Check disk performance
docker exec -it etcd1 etcdctl endpoint status --write-out=table

# Check network latency
docker exec -it etcd-monitor ping etcd1

# Review Prometheus metrics
curl http://localhost:8080/metrics | grep latency
```

### No leader elected

```bash
# Check all members
curl http://localhost:8080/api/v1/cluster/status | jq '.MemberCount'

# Check etcd logs
docker logs etcd1
docker logs etcd2
docker logs etcd3

# Verify quorum
# Need at least (N/2)+1 members healthy
```

### Prometheus not scraping metrics

```bash
# Check Prometheus targets
open http://localhost:9090/targets

# Verify metrics endpoint
curl http://localhost:8080/metrics

# Check Prometheus config
docker exec -it prometheus cat /etc/prometheus/prometheus.yml
```

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        etcd-monitor                         │
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────────┐  │
│  │   Health     │  │   Metrics    │  │     Alert       │  │
│  │   Checker    │  │  Collector   │  │    Manager      │  │
│  └──────────────┘  └──────────────┘  └─────────────────┘  │
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────────┐  │
│  │ Benchmark    │  │   etcdctl    │  │   Prometheus    │  │
│  │   Runner     │  │   Wrapper    │  │    Exporter     │  │
│  └──────────────┘  └──────────────┘  └─────────────────┘  │
│                      API Server (REST)                      │
└─────────────────────────────────────────────────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
  ┌─────▼──────┐    ┌──────▼──────┐   ┌──────▼──────┐
  │   etcd1    │    │    etcd2    │   │    etcd3    │
  │  :2379     │◄───┤   :22379    ├───►   :32379    │
  └────────────┘    └─────────────┘   └─────────────┘
        │                  │                  │
        └──────────────────┼──────────────────┘
                           │
                    ┌──────▼──────┐
                    │ Prometheus  │
                    │   :9090     │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │   Grafana   │
                    │    :3000    │
                    └─────────────┘
```

## Development

### Build from Source

```bash
# Install dependencies
go mod download

# Build
go build -o etcd-monitor ./cmd/etcd-monitor

# Run tests
go test ./...

# Run with hot reload (using air)
air
```

### Docker Build

```bash
# Build image
docker build -t etcd-monitor:latest .

# Run container
docker run -d \
  --name etcd-monitor \
  -p 8080:8080 \
  -e ENDPOINTS=etcd1:2379,etcd2:2379,etcd3:2379 \
  etcd-monitor:latest
```

## References

- [etcd Official Documentation](https://etcd.io/docs/)
- [etcd Metrics Reference](https://etcd.io/docs/latest/metrics/)
- [clay-wangzhi/etcd-metrics](https://github.com/clay-wangzhi/etcd-metrics)
- [etcd Performance Tuning](https://etcd.io/docs/latest/tuning/)
- [etcd Operations Guide](https://etcd.io/docs/latest/op-guide/)

## Support

For issues, questions, or contributions:
- GitHub Issues: [etcd-monitor/issues](https://github.com/your-org/etcd-monitor/issues)
- Documentation: [docs/](./docs/)
- Community: Join our Slack channel

## License

[Your License Here]
