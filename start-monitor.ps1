# etcd-monitor Startup Script for Windows
# PowerShell script to build and start etcd-monitor

param(
    [string]$Endpoints = "localhost:2379",
    [int]$ApiPort = 8080,
    [string]$HealthInterval = "30s",
    [string]$MetricsInterval = "10s",
    [switch]$RunBenchmark,
    [string]$BenchmarkType = "mixed",
    [int]$BenchmarkOps = 10000,
    [switch]$Help
)

function Show-Help {
    Write-Host "etcd-monitor Startup Script" -ForegroundColor Green
    Write-Host ""
    Write-Host "Usage:" -ForegroundColor Yellow
    Write-Host "  .\start-monitor.ps1 [options]"
    Write-Host ""
    Write-Host "Options:" -ForegroundColor Yellow
    Write-Host "  -Endpoints <string>        Comma-separated etcd endpoints (default: localhost:2379)"
    Write-Host "  -ApiPort <int>             API server port (default: 8080)"
    Write-Host "  -HealthInterval <string>   Health check interval (default: 30s)"
    Write-Host "  -MetricsInterval <string>  Metrics collection interval (default: 10s)"
    Write-Host "  -RunBenchmark              Run benchmark mode and exit"
    Write-Host "  -BenchmarkType <string>    Benchmark type: write, read, mixed (default: mixed)"
    Write-Host "  -BenchmarkOps <int>        Number of operations (default: 10000)"
    Write-Host "  -Help                      Show this help message"
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Yellow
    Write-Host "  # Start monitoring local etcd"
    Write-Host "  .\start-monitor.ps1"
    Write-Host ""
    Write-Host "  # Start with custom settings"
    Write-Host "  .\start-monitor.ps1 -Endpoints 'etcd1:2379,etcd2:2379,etcd3:2379' -ApiPort 8080"
    Write-Host ""
    Write-Host "  # Run benchmark"
    Write-Host "  .\start-monitor.ps1 -RunBenchmark -BenchmarkType write -BenchmarkOps 10000"
    Write-Host ""
}

if ($Help) {
    Show-Help
    exit 0
}

Write-Host "=== etcd-monitor Startup ===" -ForegroundColor Green
Write-Host ""

# Check if Go is installed
Write-Host "Checking Go installation..." -ForegroundColor Yellow
try {
    $goVersion = go version
    Write-Host "✓ Go found: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "✗ Go not found. Please install Go 1.19+ from https://golang.org/dl/" -ForegroundColor Red
    exit 1
}

# Check if etcd is accessible
Write-Host "Checking etcd connectivity..." -ForegroundColor Yellow
$firstEndpoint = $Endpoints.Split(",")[0]
try {
    $testConnection = Test-NetConnection -ComputerName $firstEndpoint.Split(":")[0] -Port ([int]$firstEndpoint.Split(":")[1]) -WarningAction SilentlyContinue -ErrorAction Stop
    if ($testConnection.TcpTestSucceeded) {
        Write-Host "✓ etcd endpoint reachable: $firstEndpoint" -ForegroundColor Green
    } else {
        Write-Host "⚠ Warning: Cannot reach etcd at $firstEndpoint" -ForegroundColor Yellow
        Write-Host "  Make sure etcd is running or update the -Endpoints parameter" -ForegroundColor Yellow
    }
} catch {
    Write-Host "⚠ Warning: Cannot verify etcd connectivity" -ForegroundColor Yellow
}

# Build the project
Write-Host ""
Write-Host "Building etcd-monitor..." -ForegroundColor Yellow
try {
    go build -o etcd-monitor.exe cmd/etcd-monitor/main.go
    Write-Host "✓ Build successful" -ForegroundColor Green
} catch {
    Write-Host "✗ Build failed" -ForegroundColor Red
    exit 1
}

# Run etcd-monitor
Write-Host ""
if ($RunBenchmark) {
    Write-Host "Running benchmark (type: $BenchmarkType, ops: $BenchmarkOps)..." -ForegroundColor Yellow
    Write-Host ""

    .\etcd-monitor.exe `
        --endpoints=$Endpoints `
        --run-benchmark `
        --benchmark-type=$BenchmarkType `
        --benchmark-ops=$BenchmarkOps
} else {
    Write-Host "Starting etcd-monitor..." -ForegroundColor Yellow
    Write-Host "  Endpoints: $Endpoints" -ForegroundColor Cyan
    Write-Host "  API Port: $ApiPort" -ForegroundColor Cyan
    Write-Host "  Health Check Interval: $HealthInterval" -ForegroundColor Cyan
    Write-Host "  Metrics Interval: $MetricsInterval" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "API will be available at: http://localhost:$ApiPort" -ForegroundColor Green
    Write-Host "Press Ctrl+C to stop" -ForegroundColor Yellow
    Write-Host ""

    .\etcd-monitor.exe `
        --endpoints=$Endpoints `
        --api-port=$ApiPort `
        --health-check-interval=$HealthInterval `
        --metrics-interval=$MetricsInterval
}
