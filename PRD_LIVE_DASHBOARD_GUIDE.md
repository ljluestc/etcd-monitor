# Product Requirements Document: ETCD-MONITOR: Live Dashboard Guide

---

## Document Information
**Project:** etcd-monitor
**Document:** LIVE_DASHBOARD_GUIDE
**Version:** 1.0.0
**Date:** 2025-10-13
**Status:** READY FOR TASK-MASTER PARSING

---

## 1. EXECUTIVE SUMMARY

### 1.1 Overview
This PRD captures the requirements and implementation details for ETCD-MONITOR: Live Dashboard Guide.

### 1.2 Purpose
This document provides a structured specification that can be parsed by task-master to generate actionable tasks.

### 1.3 Scope
The scope includes all requirements, features, and implementation details from the original documentation.

---

## 2. REQUIREMENTS

### 2.1 Functional Requirements
**Priority:** HIGH

**REQ-001:** help you get the complete etcd monitoring stack running with live Grafana dashboards.

**REQ-002:** see all services running:

**REQ-003:** now see live metrics updating in real-time:

**REQ-004:** show green/healthy

**REQ-005:** show current leader

**REQ-006:** show low latency (< 10ms typically)

**REQ-007:** show current DB size

**REQ-008:** show Raft consensus activity


## 3. TASKS

The following tasks have been identified for implementation:

**TASK_001** [MEDIUM]: Docker and Docker Compose installed

**TASK_002** [MEDIUM]: At least 4GB RAM available

**TASK_003** [MEDIUM]: Ports 2379, 3000, 8080, 9090 available

**TASK_004** [HIGH]: Go to **Status > Targets**

**TASK_005** [MEDIUM]: `etcd-monitor` (1 target)

**TASK_006** [MEDIUM]: `etcd-cluster` (3 targets)

**TASK_007** [MEDIUM]: `prometheus` (1 target)

**TASK_008** [MEDIUM]: Username: `admin`

**TASK_009** [MEDIUM]: Password: `admin`

**TASK_010** [MEDIUM]: (Change password when prompted, or skip)

**TASK_011** [HIGH]: Click **☰ Menu** (top left) → **Dashboards** → **Import**

**TASK_012** [HIGH]: **Option A: Import from JSON file**

**TASK_013** [MEDIUM]: Upload the dashboard JSON from `grafana/dashboards/etcd-dash.json`

**TASK_014** [MEDIUM]: Select **Prometheus** as the data source

**TASK_015** [MEDIUM]: Click **Import**

**TASK_016** [HIGH]: **Option B: Import from ID** (if available)

**TASK_017** [MEDIUM]: Enter dashboard ID: `3070` (etcd by Prometheus)

**TASK_018** [MEDIUM]: Click **Load**

**TASK_019** [MEDIUM]: Select **Prometheus** as the data source

**TASK_020** [MEDIUM]: Click **Import**

**TASK_021** [HIGH]: **Option C: Create Custom Dashboard**

**TASK_022** [MEDIUM]: Click **+ Create** → **Dashboard**

**TASK_023** [MEDIUM]: Add a new panel

**TASK_024** [MEDIUM]: Use the queries below to create visualizations

**TASK_025** [MEDIUM]: **Cluster Health**: Should show green/healthy

**TASK_026** [MEDIUM]: **Leader Info**: Should show current leader

**TASK_027** [MEDIUM]: **Member Count**: Should show 3

**TASK_028** [MEDIUM]: **Request Latency**: Should show low latency (< 10ms typically)

**TASK_029** [MEDIUM]: **Database Size**: Should show current DB size

**TASK_030** [MEDIUM]: **Proposal Metrics**: Should show Raft consensus activity

**TASK_031** [MEDIUM]: Visualization: Stat

**TASK_032** [MEDIUM]: Query: `etcd_cluster_healthy`

**TASK_033** [MEDIUM]: Thresholds: 1 = Green, 0 = Red

**TASK_034** [MEDIUM]: Title: "Cluster Healthy"

**TASK_035** [MEDIUM]: Visualization: Stat

**TASK_036** [MEDIUM]: Query: `etcd_cluster_has_leader`

**TASK_037** [MEDIUM]: Thresholds: 1 = Green, 0 = Red

**TASK_038** [MEDIUM]: Title: "Has Leader"

**TASK_039** [MEDIUM]: Visualization: Stat

**TASK_040** [MEDIUM]: Query: `etcd_cluster_member_count`

**TASK_041** [MEDIUM]: Title: "Members"

**TASK_042** [MEDIUM]: Visualization: Stat

**TASK_043** [MEDIUM]: Query: `increase(etcd_cluster_leader_changes_total[1h])`

**TASK_044** [MEDIUM]: Thresholds: 0-3 = Green, 4-10 = Yellow, >10 = Red

**TASK_045** [MEDIUM]: Title: "Leader Changes (1h)"

**TASK_046** [MEDIUM]: Visualization: Graph

**TASK_047** [MEDIUM]: `etcd_request_read_latency_p99_milliseconds` (Label: "Read P99")

**TASK_048** [MEDIUM]: `etcd_request_write_latency_p99_milliseconds` (Label: "Write P99")

**TASK_049** [MEDIUM]: Y-axis: milliseconds

**TASK_050** [MEDIUM]: Title: "Request Latency"

**TASK_051** [MEDIUM]: Visualization: Graph

**TASK_052** [MEDIUM]: Query: `etcd_request_rate_per_second`

**TASK_053** [MEDIUM]: Y-axis: ops/sec

**TASK_054** [MEDIUM]: Title: "Request Rate"

**TASK_055** [MEDIUM]: Visualization: Graph

**TASK_056** [MEDIUM]: Query: `etcd_mvcc_db_total_size_bytes / 1024 / 1024`

**TASK_057** [MEDIUM]: Title: "Database Size"

**TASK_058** [MEDIUM]: Visualization: Graph

**TASK_059** [MEDIUM]: `rate(etcd_server_proposals_committed_total[1m])` (Label: "Committed/sec")

**TASK_060** [MEDIUM]: `etcd_server_proposals_pending` (Label: "Pending")

**TASK_061** [MEDIUM]: Title: "Raft Proposals"

**TASK_062** [HIGH]: Check Prometheus targets: http://localhost:9090/targets

**TASK_063** [HIGH]: Verify etcd-monitor is scraping: `curl http://localhost:8080/metrics`

**TASK_064** [HIGH]: Check Grafana data source configuration

**TASK_065** [HIGH]: Verify time range in Grafana (default: Last 6 hours)

**TASK_066** [MEDIUM]: Dashboard → Settings → JSON Model

**TASK_067** [MEDIUM]: Copy the JSON

**TASK_068** [HIGH]: Update `grafana/dashboards/dashboard.yml` to include it

**TASK_069** [HIGH]: **Use persistent volumes** for etcd data

**TASK_070** [HIGH]: **Enable TLS** for etcd client and peer communication

**TASK_071** [HIGH]: **Set strong passwords** for Grafana

**TASK_072** [HIGH]: **Configure alerting** in Grafana or Prometheus Alertmanager

**TASK_073** [HIGH]: **Set up backup automation** for etcd snapshots

**TASK_074** [HIGH]: **Monitor resource usage** of the monitoring stack itself

**TASK_075** [HIGH]: **Use dedicated hardware** for etcd (especially SSD storage)

**TASK_076** [HIGH]: **Tune etcd parameters** based on workload

**TASK_077** [HIGH]: ✅ Explore the live dashboard

**TASK_078** [HIGH]: ✅ Run benchmark tests to generate load

**TASK_079** [HIGH]: ✅ Create custom panels for your specific needs

**TASK_080** [HIGH]: ✅ Set up alerting rules in Grafana

**TASK_081** [HIGH]: ✅ Configure Slack/PagerDuty notifications

**TASK_082** [HIGH]: ✅ Set up automated backups

**TASK_083** [HIGH]: ✅ Plan capacity based on metrics trends

**TASK_084** [MEDIUM]: **Prometheus**: http://localhost:9090

**TASK_085** [MEDIUM]: **Grafana**: http://localhost:3000

**TASK_086** [MEDIUM]: **etcd-monitor API**: http://localhost:8080

**TASK_087** [MEDIUM]: http://localhost:2379 (etcd1)

**TASK_088** [MEDIUM]: http://localhost:22379 (etcd2)

**TASK_089** [MEDIUM]: http://localhost:32379 (etcd3)


## 4. DETAILED SPECIFICATIONS

### 4.1 Original Content

The following sections contain the original documentation:


#### Live Dashboard Setup Guide

# Live Dashboard Setup Guide


#### Quick Start Get Live Dashboard Running In 5 Minutes 

## Quick Start - Get Live Dashboard Running in 5 Minutes!

This guide will help you get the complete etcd monitoring stack running with live Grafana dashboards.


#### Prerequisites

## Prerequisites

- Docker and Docker Compose installed
- At least 4GB RAM available
- Ports 2379, 3000, 8080, 9090 available


#### Step By Step Setup

## Step-by-Step Setup


#### Step 1 Build And Start Everything

### Step 1: Build and Start Everything

```bash

#### Navigate To Project Directory

# Navigate to project directory
cd etcd-monitor


#### Start The Full Stack Etcd Cluster Monitor Prometheus Grafana 

# Start the full stack (etcd cluster + monitor + Prometheus + Grafana)
docker-compose -f docker-compose-full.yml up -d


#### Check That All Containers Are Running

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


#### Step 2 Verify Etcd Cluster

### Step 2: Verify etcd Cluster

```bash

#### Check Cluster Health

# Check cluster health
docker exec etcd1 etcdctl endpoint health --cluster


#### Should Output 

# Should output:

#### Http Etcd1 2379 Is Healthy

# http://etcd1:2379 is healthy

#### Http Etcd2 2379 Is Healthy

# http://etcd2:2379 is healthy

#### Http Etcd3 2379 Is Healthy

# http://etcd3:2379 is healthy
```


#### Step 3 Verify Etcd Monitor Api

### Step 3: Verify etcd-monitor API

```bash

#### Check Health

# Check health
curl http://localhost:8080/health


#### Get Cluster Status

# Get cluster status
curl http://localhost:8080/api/v1/cluster/status | jq


#### Get Current Metrics

# Get current metrics
curl http://localhost:8080/api/v1/metrics/current | jq
```


#### Step 4 Verify Prometheus

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


#### Step 5 Access Grafana Dashboard

### Step 5: Access Grafana Dashboard

Open browser to: http://localhost:3000

**Login Credentials:**
- Username: `admin`
- Password: `admin`
- (Change password when prompted, or skip)


#### Step 6 Import Etcd Dashboard

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


#### Step 7 View Live Metrics 

### Step 7: View Live Metrics!

You should now see live metrics updating in real-time:

**Key Panels to Check:**
- **Cluster Health**: Should show green/healthy
- **Leader Info**: Should show current leader
- **Member Count**: Should show 3
- **Request Latency**: Should show low latency (< 10ms typically)
- **Database Size**: Should show current DB size
- **Proposal Metrics**: Should show Raft consensus activity


#### Useful Prometheus Queries For Dashboards

## Useful Prometheus Queries for Dashboards


#### Cluster Health

### Cluster Health
```promql

#### Cluster Healthy 1 Yes 0 No 

# Cluster healthy (1 = yes, 0 = no)
etcd_cluster_healthy


#### Has Leader 1 Yes 0 No 

# Has leader (1 = yes, 0 = no)
etcd_cluster_has_leader


#### Member Count

# Member count
etcd_cluster_member_count


#### Leader Changes Rate 

# Leader changes (rate)
rate(etcd_cluster_leader_changes_total[5m])
```


#### Performance Metrics

### Performance Metrics
```promql

#### Read Latency P99

# Read latency P99
etcd_request_read_latency_p99_milliseconds


#### Write Latency P99

# Write latency P99
etcd_request_write_latency_p99_milliseconds


#### Request Rate

# Request rate
etcd_request_rate_per_second


#### Database Size

# Database size
etcd_mvcc_db_total_size_bytes / 1024 / 1024  # Convert to MB
```


#### Raft Metrics

### Raft Metrics
```promql

#### Proposals Committed

# Proposals committed
etcd_server_proposals_committed_total


#### Proposals Pending

# Proposals pending
etcd_server_proposals_pending


#### Proposals Failed

# Proposals failed
rate(etcd_server_proposals_failed_total[5m])
```


#### Resource Usage

### Resource Usage
```promql

#### Memory Usage

# Memory usage
etcd_server_memory_usage_bytes / 1024 / 1024  # Convert to MB


#### Disk Usage

# Disk usage
etcd_server_disk_usage_bytes / 1024 / 1024  # Convert to MB


#### Active Connections

# Active connections
etcd_server_active_connections


#### Watchers

# Watchers
etcd_server_watchers
```


#### Create A Custom Dashboard

## Create a Custom Dashboard

Here's a quick dashboard configuration with essential panels:


#### 1 Cluster Overview Row

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


#### 2 Performance Row

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


#### 3 Database Row

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


#### Generate Test Load

## Generate Test Load

To see the dashboard in action with real metrics:

```bash

#### Write Test Data

# Write test data
for i in {1..1000}; do
  docker exec etcd1 etcdctl put /test/key-$i "value-$i"
done


#### Read Test Data

# Read test data
for i in {1..1000}; do
  docker exec etcd1 etcdctl get /test/key-$i > /dev/null
done


#### Watch The Dashboard Update With Latency And Throughput Metrics 

# Watch the dashboard update with latency and throughput metrics!
```


#### Run Benchmark Tests

## Run Benchmark Tests

Generate realistic load:

```bash

#### Build Etcd Monitor Binary If Not Already Built 

# Build etcd-monitor binary (if not already built)
go build -o etcd-monitor cmd/etcd-monitor/main.go


#### Run Write Benchmark

# Run write benchmark
./etcd-monitor \
  --endpoints=localhost:2379,localhost:22379,localhost:32379 \
  --run-benchmark \
  --benchmark-type=write \
  --benchmark-ops=10000


#### Run Mixed Benchmark

# Run mixed benchmark
./etcd-monitor \
  --endpoints=localhost:2379,localhost:22379,localhost:32379 \
  --run-benchmark \
  --benchmark-type=mixed \
  --benchmark-ops=20000


#### Watch The Metrics Spike In Grafana 

# Watch the metrics spike in Grafana!
```


#### Troubleshooting

## Troubleshooting


#### Services Won T Start

### Services won't start
```bash

#### Check Logs

# Check logs
docker logs etcd1
docker logs etcd2
docker logs etcd3
```


#### Restart Specific Service

# Restart specific service
docker-compose -f docker-compose-full.yml restart etcd-monitor


#### Rebuild If Needed

# Rebuild if needed
docker-compose -f docker-compose-full.yml build --no-cache
docker-compose -f docker-compose-full.yml up -d
```


#### No Metrics In Grafana

### No metrics in Grafana
1. Check Prometheus targets: http://localhost:9090/targets
2. Verify etcd-monitor is scraping: `curl http://localhost:8080/metrics`
3. Check Grafana data source configuration
4. Verify time range in Grafana (default: Last 6 hours)


#### Etcd Cluster Unhealthy

### etcd cluster unhealthy
```bash

#### Check Each Member

# Check each member
docker exec etcd1 etcdctl endpoint health
docker exec etcd2 etcdctl endpoint health
docker exec etcd3 etcdctl endpoint health


#### Cleanup

## Cleanup

```bash

#### Stop All Services

# Stop all services
docker-compose -f docker-compose-full.yml down


#### Remove Volumes Caution Deletes All Data 

# Remove volumes (caution: deletes all data)
docker-compose -f docker-compose-full.yml down -v


#### Remove Images

# Remove images
docker-compose -f docker-compose-full.yml down --rmi all
```


#### Advanced Persistent Grafana Dashboards

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


#### Production Considerations

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


#### Next Steps

## Next Steps

1. ✅ Explore the live dashboard
2. ✅ Run benchmark tests to generate load
3. ✅ Create custom panels for your specific needs
4. ✅ Set up alerting rules in Grafana
5. ✅ Configure Slack/PagerDuty notifications
6. ✅ Set up automated backups
7. ✅ Plan capacity based on metrics trends


#### Resources

## Resources

- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000
- **etcd-monitor API**: http://localhost:8080
- **etcd endpoints**:
  - http://localhost:2379 (etcd1)
  - http://localhost:22379 (etcd2)
  - http://localhost:32379 (etcd3)

Enjoy your live etcd monitoring dashboard! 🎉📊


---

## 5. TECHNICAL REQUIREMENTS

### 5.1 Dependencies
- All dependencies from original documentation apply
- Standard development environment
- Required tools and libraries as specified

### 5.2 Compatibility
- Compatible with existing infrastructure
- Follows project standards and conventions

---

## 6. SUCCESS CRITERIA

### 6.1 Functional Success Criteria
- All identified tasks completed successfully
- All requirements implemented as specified
- All tests passing

### 6.2 Quality Success Criteria
- Code meets quality standards
- Documentation is complete and accurate
- No critical issues remaining

---

## 7. IMPLEMENTATION PLAN

### Phase 1: Preparation
- Review all requirements and tasks
- Set up development environment
- Gather necessary resources

### Phase 2: Implementation
- Execute tasks in priority order
- Follow best practices
- Test incrementally

### Phase 3: Validation
- Run comprehensive tests
- Validate against requirements
- Document completion

---

## 8. TASK-MASTER INTEGRATION

### How to Parse This PRD

```bash
# Parse this PRD with task-master
task-master parse-prd --input="{doc_name}_PRD.md"

# List generated tasks
task-master list

# Start execution
task-master next
```

### Expected Task Generation
Task-master should generate approximately {len(tasks)} tasks from this PRD.

---

## 9. APPENDIX

### 9.1 References
- Original document: {doc_name}.md
- Project: {project_name}

### 9.2 Change History
| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | {datetime.now().strftime('%Y-%m-%d')} | Initial PRD conversion |

---

*End of PRD*
*Generated by MD-to-PRD Converter*
