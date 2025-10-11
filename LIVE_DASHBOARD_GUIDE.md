# Live Dashboard Setup Guide

## Quick Start - Get Live Dashboard Running in 5 Minutes!

This guide will help you get the complete etcd monitoring stack running with live Grafana dashboards.

## Prerequisites

- Docker and Docker Compose installed
- At least 4GB RAM available
- Ports 2379, 3000, 8080, 9090 available

## Step-by-Step Setup

### Step 1: Build and Start Everything

```bash
# Navigate to project directory
cd etcd-monitor

# Start the full stack (etcd cluster + monitor + Prometheus + Grafana)
docker-compose -f docker-compose-full.yml up -d

# Check that all containers are running
docker-compose -f docker-compose-full.yml ps
```

You should see all services running:
```
NAME            STATUS    PORTS
etcd1           Up        0.0.0.0:2379->2379/tcp, 0.0.0.0:2380->2380/tcp
etcd2           Up        0.0.0.0:22379->2379/tcp, 0.0.0.0:22380->2380/tcp
etcd3           Up        0.0.0.0:32379->2379/tcp, 0.0.0.0:32380->2380/tcp
etcd-monitor    Up        0.0.0.0:8080->8080/tcp
prometheus      Up        0.0.0.0:9090->9090/tcp
grafana         Up        0.0.0.0:3000->3000/tcp
```

### Step 2: Verify etcd Cluster

```bash
# Check cluster health
docker exec etcd1 etcdctl endpoint health --cluster

# Should output:
# http://etcd1:2379 is healthy
# http://etcd2:2379 is healthy
# http://etcd3:2379 is healthy
```

### Step 3: Verify etcd-monitor API

```bash
# Check health
curl http://localhost:8080/health

# Get cluster status
curl http://localhost:8080/api/v1/cluster/status | jq

# Get current metrics
curl http://localhost:8080/api/v1/metrics/current | jq
```

### Step 4: Verify Prometheus

Open browser to: http://localhost:9090

1. Go to **Status > Targets**
2. Verify all targets are **UP**:
   - `etcd-monitor` (1 target)
   - `etcd-cluster` (3 targets)
   - `prometheus` (1 target)

3. Try a query in the **Graph** tab:
   ```promql
   etcd_cluster_healthy
   ```

### Step 5: Access Grafana Dashboard

Open browser to: http://localhost:3000

**Login Credentials:**
- Username: `admin`
- Password: `admin`
- (Change password when prompted, or skip)

### Step 6: Import etcd Dashboard

1. Click **☰ Menu** (top left) → **Dashboards** → **Import**

2. **Option A: Import from JSON file**
   - Upload the dashboard JSON from `grafana/dashboards/etcd-dash.json`
   - Select **Prometheus** as the data source
   - Click **Import**

3. **Option B: Import from ID** (if available)
   - Enter dashboard ID: `3070` (etcd by Prometheus)
   - Click **Load**
   - Select **Prometheus** as the data source
   - Click **Import**

4. **Option C: Create Custom Dashboard**
   - Click **+ Create** → **Dashboard**
   - Add a new panel
   - Use the queries below to create visualizations

### Step 7: View Live Metrics!

You should now see live metrics updating in real-time:

**Key Panels to Check:**
- **Cluster Health**: Should show green/healthy
- **Leader Info**: Should show current leader
- **Member Count**: Should show 3
- **Request Latency**: Should show low latency (< 10ms typically)
- **Database Size**: Should show current DB size
- **Proposal Metrics**: Should show Raft consensus activity

## Useful Prometheus Queries for Dashboards

### Cluster Health
```promql
# Cluster healthy (1 = yes, 0 = no)
etcd_cluster_healthy

# Has leader (1 = yes, 0 = no)
etcd_cluster_has_leader

# Member count
etcd_cluster_member_count

# Leader changes (rate)
rate(etcd_cluster_leader_changes_total[5m])
```

### Performance Metrics
```promql
# Read latency P99
etcd_request_read_latency_p99_milliseconds

# Write latency P99
etcd_request_write_latency_p99_milliseconds

# Request rate
etcd_request_rate_per_second

# Database size
etcd_mvcc_db_total_size_bytes / 1024 / 1024  # Convert to MB
```

### Raft Metrics
```promql
# Proposals committed
etcd_server_proposals_committed_total

# Proposals pending
etcd_server_proposals_pending

# Proposals failed
rate(etcd_server_proposals_failed_total[5m])
```

### Resource Usage
```promql
# Memory usage
etcd_server_memory_usage_bytes / 1024 / 1024  # Convert to MB

# Disk usage
etcd_server_disk_usage_bytes / 1024 / 1024  # Convert to MB

# Active connections
etcd_server_active_connections

# Watchers
etcd_server_watchers
```

## Create a Custom Dashboard

Here's a quick dashboard configuration with essential panels:

### 1. Cluster Overview Row

**Panel 1: Cluster Health**
- Visualization: Stat
- Query: `etcd_cluster_healthy`
- Thresholds: 1 = Green, 0 = Red
- Title: "Cluster Healthy"

**Panel 2: Has Leader**
- Visualization: Stat
- Query: `etcd_cluster_has_leader`
- Thresholds: 1 = Green, 0 = Red
- Title: "Has Leader"

**Panel 3: Member Count**
- Visualization: Stat
- Query: `etcd_cluster_member_count`
- Title: "Members"

**Panel 4: Leader Changes**
- Visualization: Stat
- Query: `increase(etcd_cluster_leader_changes_total[1h])`
- Thresholds: 0-3 = Green, 4-10 = Yellow, >10 = Red
- Title: "Leader Changes (1h)"

### 2. Performance Row

**Panel 5: Request Latency**
- Visualization: Graph
- Queries:
  - `etcd_request_read_latency_p99_milliseconds` (Label: "Read P99")
  - `etcd_request_write_latency_p99_milliseconds` (Label: "Write P99")
- Y-axis: milliseconds
- Title: "Request Latency"

**Panel 6: Request Rate**
- Visualization: Graph
- Query: `etcd_request_rate_per_second`
- Y-axis: ops/sec
- Title: "Request Rate"

### 3. Database Row

**Panel 7: Database Size**
- Visualization: Graph
- Query: `etcd_mvcc_db_total_size_bytes / 1024 / 1024`
- Y-axis: MB
- Title: "Database Size"

**Panel 8: Proposals**
- Visualization: Graph
- Queries:
  - `rate(etcd_server_proposals_committed_total[1m])` (Label: "Committed/sec")
  - `etcd_server_proposals_pending` (Label: "Pending")
- Title: "Raft Proposals"

## Generate Test Load

To see the dashboard in action with real metrics:

```bash
# Write test data
for i in {1..1000}; do
  docker exec etcd1 etcdctl put /test/key-$i "value-$i"
done

# Read test data
for i in {1..1000}; do
  docker exec etcd1 etcdctl get /test/key-$i > /dev/null
done

# Watch the dashboard update with latency and throughput metrics!
```

## Run Benchmark Tests

Generate realistic load:

```bash
# Build etcd-monitor binary (if not already built)
go build -o etcd-monitor cmd/etcd-monitor/main.go

# Run write benchmark
./etcd-monitor \
  --endpoints=localhost:2379,localhost:22379,localhost:32379 \
  --run-benchmark \
  --benchmark-type=write \
  --benchmark-ops=10000

# Run mixed benchmark
./etcd-monitor \
  --endpoints=localhost:2379,localhost:22379,localhost:32379 \
  --run-benchmark \
  --benchmark-type=mixed \
  --benchmark-ops=20000

# Watch the metrics spike in Grafana!
```

## Troubleshooting

### Services won't start
```bash
# Check logs
docker-compose -f docker-compose-full.yml logs

# Restart specific service
docker-compose -f docker-compose-full.yml restart etcd-monitor

# Rebuild if needed
docker-compose -f docker-compose-full.yml build --no-cache
docker-compose -f docker-compose-full.yml up -d
```

### No metrics in Grafana
1. Check Prometheus targets: http://localhost:9090/targets
2. Verify etcd-monitor is scraping: `curl http://localhost:8080/metrics`
3. Check Grafana data source configuration
4. Verify time range in Grafana (default: Last 6 hours)

### etcd cluster unhealthy
```bash
# Check each member
docker exec etcd1 etcdctl endpoint health
docker exec etcd2 etcdctl endpoint health
docker exec etcd3 etcdctl endpoint health

# Check logs
docker logs etcd1
docker logs etcd2
docker logs etcd3
```

## Cleanup

```bash
# Stop all services
docker-compose -f docker-compose-full.yml down

# Remove volumes (caution: deletes all data)
docker-compose -f docker-compose-full.yml down -v

# Remove images
docker-compose -f docker-compose-full.yml down --rmi all
```

## Advanced: Persistent Grafana Dashboards

To save your custom dashboards permanently:

1. Export dashboard as JSON:
   - Dashboard → Settings → JSON Model
   - Copy the JSON

2. Save to file:
   ```bash
   # Save to dashboards directory
   cat > grafana/dashboards/my-custom-dashboard.json << 'EOF'
   {your-dashboard-json-here}
   EOF
   ```

3. Update `grafana/dashboards/dashboard.yml` to include it

4. Restart Grafana:
   ```bash
   docker-compose -f docker-compose-full.yml restart grafana
   ```

## Production Considerations

For production deployments:

1. **Use persistent volumes** for etcd data
2. **Enable TLS** for etcd client and peer communication
3. **Set strong passwords** for Grafana
4. **Configure alerting** in Grafana or Prometheus Alertmanager
5. **Set up backup automation** for etcd snapshots
6. **Monitor resource usage** of the monitoring stack itself
7. **Use dedicated hardware** for etcd (especially SSD storage)
8. **Tune etcd parameters** based on workload

## Next Steps

1. ✅ Explore the live dashboard
2. ✅ Run benchmark tests to generate load
3. ✅ Create custom panels for your specific needs
4. ✅ Set up alerting rules in Grafana
5. ✅ Configure Slack/PagerDuty notifications
6. ✅ Set up automated backups
7. ✅ Plan capacity based on metrics trends

## Resources

- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000
- **etcd-monitor API**: http://localhost:8080
- **etcd endpoints**:
  - http://localhost:2379 (etcd1)
  - http://localhost:22379 (etcd2)
  - http://localhost:32379 (etcd3)

Enjoy your live etcd monitoring dashboard! 🎉📊
