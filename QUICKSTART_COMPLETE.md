# etcd-monitor Quickstart Guide
## Complete Implementation Based on Clay Wang's etcd Guide

This quickstart guide helps you get etcd-monitor up and running quickly.

## Prerequisites

- Go 1.19+ installed
- Running etcd cluster (v3.4+)
- (Optional) Docker for containerized setup

## Quick Start

### 1. Install Dependencies

```bash
# Navigate to project directory
cd etcd-monitor

# Download dependencies
go mod download
```

### 2. Start a Local etcd Cluster (Testing)

If you don't have an etcd cluster running:

```bash
# Using docker-compose
docker-compose up -d etcd

# Or download and run etcd binary
ETCD_VER=v3.5.9
DOWNLOAD_URL=https://github.com/etcd-io/etcd/releases/download
curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
tar xzvf /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz -C /tmp
/tmp/etcd-${ETCD_VER}-linux-amd64/etcd
```

### 3. Build etcd-monitor

```bash
# Build the binary
go build -o etcd-monitor cmd/etcd-monitor/main.go

# Or use the Makefile
make build
```

### 4. Run etcd-monitor

```bash
# Basic usage with local etcd
./etcd-monitor \
  --endpoints=localhost:2379 \
  --api-port=8080

# With all options
./etcd-monitor \
  --endpoints=localhost:2379 \
  --dial-timeout=5s \
  --health-check-interval=30s \
  --metrics-interval=10s \
  --api-port=8080 \
  --max-latency-ms=100 \
  --max-db-size-mb=8192

# With TLS (production)
./etcd-monitor \
  --endpoints=etcd1:2379,etcd2:2379,etcd3:2379 \
  --cert=/path/to/cert.pem \
  --key=/path/to/key.pem \
  --cacert=/path/to/ca.pem \
  --api-port=8080
```

## Features Overview

### 1. Health Monitoring

The health checker continuously monitors:
- Cluster membership
- Leader status and changes
- Network partitions
- Member health and connectivity
- etcd alarms (NOSPACE, CORRUPT, etc.)

**API Endpoints:**
```bash
# Get cluster status
curl http://localhost:8080/api/v1/cluster/status

# Get leader information
curl http://localhost:8080/api/v1/cluster/leader

# Check overall health
curl http://localhost:8080/health
```

### 2. Performance Metrics

Collects comprehensive metrics:
- **Request Latency:** Read/write latency (p50, p95, p99)
- **Database Size:** Total size and in-use size
- **Raft Metrics:** Proposals committed, applied, pending
- **Client Metrics:** Active connections, watchers

**API Endpoints:**
```bash
# Get current metrics snapshot
curl http://localhost:8080/api/v1/metrics/current

# Get latency metrics
curl http://localhost:8080/api/v1/metrics/latency
```

**Example Response:**
```json
{
  "timestamp": "2025-10-11T10:30:00Z",
  "readLatencyP50": 1.5,
  "readLatencyP95": 3.2,
  "readLatencyP99": 5.8,
  "writeLatencyP50": 2.1,
  "writeLatencyP95": 4.5,
  "writeLatencyP99": 8.2,
  "dbSize": 10485760,
  "dbSizeInUse": 8388608,
  "proposalCommitted": 15234,
  "proposalApplied": 15234
}
```

### 3. Benchmarking

Built-in performance benchmarking tool:

```bash
# Run a write benchmark
./etcd-monitor \
  --run-benchmark \
  --benchmark-type=write \
  --benchmark-ops=10000

# Run a read benchmark
./etcd-monitor \
  --run-benchmark \
  --benchmark-type=read \
  --benchmark-ops=10000

# Run a mixed benchmark (70% read, 30% write)
./etcd-monitor \
  --run-benchmark \
  --benchmark-type=mixed \
  --benchmark-ops=10000
```

**Expected Output:**
```
=== Benchmark Results ===
Type:             write
Duration:         5.234s
Total Operations: 10000
Successful:       10000
Failed:           0
Throughput:       1910.45 ops/sec

Latency:
  Average:        5.23 ms
  P50:            4.12 ms
  P95:            8.45 ms
  P99:            12.34 ms
  Min:            1.23 ms
  Max:            45.67 ms

Latency Distribution:
  <1ms         :    125 (1.25%)
  1-5ms        :   7823 (78.23%)
  5-10ms       :   1567 (15.67%)
  10-50ms      :    485 (4.85%)
  50-100ms     :      0 (0.00%)
  100-500ms    :      0 (0.00%)
  >500ms       :      0 (0.00%)
```

### 4. Alerting

Automated alerting based on configurable thresholds:

**Alert Types:**
- Cluster unhealthy
- No leader elected
- High latency (> threshold)
- High disk usage (> threshold)
- Network partition detected
- etcd alarms triggered

**Alert Channels:**
- Console logging (default)
- Email (SMTP)
- Slack webhooks
- PagerDuty Events API
- Generic webhooks

**Configuration Example:**
```go
// In code or via configuration file
alertManager.AddChannel(&SlackChannel{
    WebhookURL: "https://hooks.slack.com/services/YOUR/WEBHOOK/URL",
    Channel:    "#etcd-alerts",
    Username:   "etcd-monitor",
})

alertManager.AddChannel(&PagerDutyChannel{
    IntegrationKey: "your-pagerduty-integration-key",
})
```

### 5. etcdctl Operations

Programmatic etcdctl operations via API:

```bash
# Get a key
curl -X POST http://localhost:8080/api/v1/etcdctl/get \
  -d '{"key": "/config/app"}'

# Put a key-value
curl -X POST http://localhost:8080/api/v1/etcdctl/put \
  -d '{"key": "/config/app", "value": "production"}'

# List members
curl http://localhost:8080/api/v1/cluster/members

# Check endpoint health
curl http://localhost:8080/api/v1/cluster/health
```

## Configuration

### Configuration File

Create `config.yaml`:

```yaml
# etcd connection
endpoints:
  - localhost:2379
  - localhost:22379
  - localhost:32379

dial_timeout: 5s

# TLS configuration (optional)
tls:
  cert_file: /path/to/cert.pem
  key_file: /path/to/key.pem
  ca_file: /path/to/ca.pem
  insecure_skip_verify: false

# Monitoring intervals
health_check_interval: 30s
metrics_interval: 10s

# Alert thresholds
alert_thresholds:
  max_latency_ms: 100
  max_database_size_mb: 8192
  min_available_nodes: 2
  max_leader_changes_per_hour: 3
  max_error_rate: 0.05
  min_disk_space_percent: 10.0

# API server
api:
  host: 0.0.0.0
  port: 8080

# Alerting channels
alerts:
  - type: slack
    webhook_url: https://hooks.slack.com/services/YOUR/WEBHOOK/URL
    channel: "#etcd-alerts"

  - type: email
    smtp_server: smtp.gmail.com
    smtp_port: 587
    from: etcd-monitor@example.com
    to:
      - ops@example.com
      - oncall@example.com
    username: your-email@gmail.com
    password: your-app-password

# Benchmarking
benchmark:
  enabled: true
  interval: 1h
```

Load configuration:
```bash
./etcd-monitor --config=config.yaml
```

## Docker Deployment

### Build Docker Image

```bash
docker build -t etcd-monitor:latest .
```

### Run with Docker

```bash
docker run -d \
  --name etcd-monitor \
  -p 8080:8080 \
  -e ENDPOINTS=etcd1:2379,etcd2:2379,etcd3:2379 \
  etcd-monitor:latest
```

### Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  etcd:
    image: quay.io/coreos/etcd:v3.5.9
    command:
      - /usr/local/bin/etcd
      - --name=etcd1
      - --listen-client-urls=http://0.0.0.0:2379
      - --advertise-client-urls=http://etcd:2379
      - --listen-peer-urls=http://0.0.0.0:2380
    ports:
      - "2379:2379"

  etcd-monitor:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ENDPOINTS=etcd:2379
      - API_PORT=8080
      - HEALTH_CHECK_INTERVAL=30s
      - METRICS_INTERVAL=10s
    depends_on:
      - etcd
```

Start:
```bash
docker-compose up -d
```

## Kubernetes Deployment

### Using Helm

```bash
# Add helm repository (if published)
helm repo add etcd-monitor https://your-repo.com/charts

# Install
helm install etcd-monitor etcd-monitor/etcd-monitor \
  --set etcd.endpoints="etcd-client:2379" \
  --set api.port=8080
```

### Using Manifests

```bash
# Apply Kubernetes manifests
kubectl apply -f k8s/

# Check status
kubectl get pods -l app=etcd-monitor
kubectl logs -f deployment/etcd-monitor
```

## Monitoring Best Practices

Based on Clay Wang's etcd guide and industry best practices:

### 1. Key Metrics to Watch

**Critical Metrics:**
- Leader changes per hour (should be < 3)
- Request latency p99 (should be < 100ms)
- Disk fsync latency p99 (should be < 10ms for production)
- Database size (monitor growth rate)
- Proposal failure rate (should be near 0)

**Resource Metrics:**
- CPU usage
- Memory usage
- Disk usage (keep < 90%)
- Network bandwidth

### 2. Alert Thresholds

**Production Recommendations:**
```yaml
# Critical alerts
no_leader_timeout: 1m
fsync_latency_p99: 10ms (warning), 100ms (critical)
backend_commit_latency_p99: 25ms (warning), 250ms (critical)
database_size_percent: 80% (warning), 90% (critical)

# Warning alerts
leader_changes_per_hour: 3
proposal_failure_rate: 1%
disk_usage: 80%
memory_usage: 80%
```

### 3. Performance Tuning

**Hardware:**
- ✅ Use SSD storage (critical for etcd performance)
- ✅ Dedicated disk for etcd WAL and data
- ✅ 4+ CPU cores
- ✅ 8+ GB RAM
- ✅ Low-latency network (< 10ms between members)

**OS Tuning:**
- Set CPU governor to performance mode
- Use deadline or noop I/O scheduler for SSDs
- Increase file descriptor limits

**etcd Configuration:**
- `--snapshot-count=10000` (adjust based on workload)
- `--heartbeat-interval=100` (default)
- `--election-timeout=1000` (default)
- `--quota-backend-bytes=8GB` (adjust based on needs)
- `--auto-compaction-retention=1h` (enable auto-compaction)

### 4. Backup Strategy

```bash
# Create snapshot
./etcd-monitor --endpoints=localhost:2379 snapshot save backup.db

# Verify snapshot
./etcd-monitor --endpoints=localhost:2379 snapshot status backup.db

# Schedule regular backups
# Add to crontab:
0 */6 * * * /path/to/etcd-monitor snapshot save /backups/etcd-$(date +\%Y\%m\%d-\%H\%M).db
```

### 5. Capacity Planning

Monitor these trends:
- Database size growth rate
- Request rate trends
- Latency trends
- Resource utilization trends

**Example Analysis:**
```bash
# Get metrics every hour for 24 hours
for i in {1..24}; do
  curl http://localhost:8080/api/v1/metrics/current >> metrics-$(date +\%Y\%m\%d).json
  sleep 3600
done

# Analyze growth
# DB size growing at 100MB/day = need to plan for compaction/cleanup
```

## Troubleshooting

### High Latency

**Symptoms:**
- Request latency > 100ms
- Slow application response

**Diagnosis:**
```bash
# Check disk performance
./etcd-monitor disk-test --directory=/var/lib/etcd

# Check current latency
curl http://localhost:8080/api/v1/metrics/latency

# Run benchmark
./etcd-monitor --run-benchmark --benchmark-type=mixed
```

**Solutions:**
1. Check disk I/O (should use SSD)
2. Check network latency between members
3. Check for high write load
4. Consider compaction if DB is large
5. Check CPU/memory usage

### Leader Elections

**Symptoms:**
- Frequent leader changes
- Cluster instability

**Diagnosis:**
```bash
# Check leader history
curl http://localhost:8080/api/v1/cluster/status | jq '.leaderChanges'

# Check network latency
ping etcd-member-2
```

**Solutions:**
1. Check network stability
2. Tune heartbeat-interval and election-timeout
3. Check for resource contention
4. Ensure member clocks are synchronized (NTP)

### Database Size Growing

**Symptoms:**
- DB size continuously growing
- Performance degradation

**Diagnosis:**
```bash
# Check current size
curl http://localhost:8080/api/v1/metrics/current | jq '.dbSize'

# Check fragmentation
etcdctl endpoint status --write-out=table
```

**Solutions:**
1. Enable auto-compaction
2. Manual compaction: `etcdctl compact <revision>`
3. Defragment: `etcdctl defrag`
4. Review data retention policies
5. Check for lease leaks

### No Space Alarm

**Symptoms:**
- etcd stops accepting writes
- NOSPACE alarm triggered

**Diagnosis:**
```bash
# Check alarms
curl http://localhost:8080/api/v1/cluster/status | jq '.alarms'
etcdctl alarm list
```

**Solutions:**
```bash
# 1. Check disk space
df -h /var/lib/etcd

# 2. Increase quota if space available
etcdctl --endpoints=localhost:2379 \
  member update --quota-backend-bytes=8589934592  # 8GB

# 3. Compact and defrag
etcdctl compact $(etcdctl endpoint status --write-out="json" | jq '.[0].Status.header.revision')
etcdctl defrag

# 4. Disarm alarm
etcdctl alarm disarm
```

## Advanced Features

### Custom Metrics Export

Export metrics to Prometheus:

```bash
# Prometheus will scrape this endpoint
curl http://localhost:8080/metrics
```

Add to `prometheus.yml`:
```yaml
scrape_configs:
  - job_name: 'etcd-monitor'
    static_configs:
      - targets: ['localhost:8080']
```

### Grafana Dashboards

Import pre-built dashboards from `grafana/dashboards/`:
1. etcd Cluster Overview
2. Performance Metrics
3. Alerting Dashboard

### Programmatic Access

Use the Go client library:

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/etcd-monitor/pkg/monitor"
    clientv3 "go.etcd.io/etcd/client/v3"
    "go.uber.org/zap"
)

func main() {
    logger, _ := zap.NewProduction()

    config := &monitor.Config{
        Endpoints:   []string{"localhost:2379"},
        DialTimeout: 5 * time.Second,
        HealthCheckInterval: 30 * time.Second,
        MetricsInterval:     10 * time.Second,
    }

    monitorService, err := monitor.NewMonitorService(config, logger)
    if err != nil {
        log.Fatal(err)
    }

    if err := monitorService.Start(); err != nil {
        log.Fatal(err)
    }

    // Get cluster status
    status, _ := monitorService.GetClusterStatus()
    log.Printf("Cluster healthy: %v, Leader: %d", status.Healthy, status.LeaderID)

    // Get metrics
    metrics, _ := monitorService.GetCurrentMetrics()
    log.Printf("Read latency P99: %.2fms", metrics.ReadLatencyP99)
}
```

## Next Steps

1. **Set up production monitoring:**
   - Configure TLS
   - Set up alerting channels
   - Configure backup schedules

2. **Integrate with existing stack:**
   - Export to Prometheus
   - Create Grafana dashboards
   - Set up log aggregation

3. **Tune for your workload:**
   - Run benchmarks
   - Adjust thresholds
   - Optimize etcd configuration

4. **Set up automation:**
   - Automated backups
   - Automated compaction
   - Automated alerting

## Resources

- **etcd Documentation:** https://etcd.io/docs/
- **etcd Metrics Guide:** https://etcd.io/docs/latest/metrics/
- **Clay Wang's Guide:** https://clay-wangzhi.com/cloudnative/tools-ops/etcd.html
- **PRD:** See PRD_COMPLETE.md
- **Implementation Status:** See IMPLEMENTATION_STATUS.md

## Support

For issues and questions:
- GitHub Issues: https://github.com/your-org/etcd-monitor/issues
- Documentation: See docs/ directory
- Examples: See pkg/examples/

## License

[Your License Here]
