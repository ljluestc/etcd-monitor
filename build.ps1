# PowerShell build script for etcd-monitor on Windows

param(
    [string]$Command = "build",
    [string]$Platform = "windows",
    [string]$Version = "1.0.0"
)

$BinaryName = "etcd-monitor"
$BuildDir = "bin"
$GoFlags = "-v"

function Build-Application {
    Write-Host "Building $BinaryName..." -ForegroundColor Green

    if (!(Test-Path $BuildDir)) {
        New-Item -ItemType Directory -Path $BuildDir | Out-Null
    }

    $output = Join-Path $BuildDir "$BinaryName.exe"
    $ldflags = "-X main.version=$Version"

    go build $GoFlags -ldflags $ldflags -o $output cmd/etcd-monitor/main.go

    if ($LASTEXITCODE -eq 0) {
        Write-Host "Build complete: $output" -ForegroundColor Green
    } else {
        Write-Host "Build failed!" -ForegroundColor Red
        exit 1
    }
}

function Build-TaskMaster {
    Write-Host "Building TaskMaster CLI..." -ForegroundColor Green

    if (!(Test-Path $BuildDir)) {
        New-Item -ItemType Directory -Path $BuildDir | Out-Null
    }

    $output = Join-Path $BuildDir "taskmaster.exe"

    go build $GoFlags -o $output taskmaster.go

    if ($LASTEXITCODE -eq 0) {
        Write-Host "TaskMaster built: $output" -ForegroundColor Green
    } else {
        Write-Host "Build failed!" -ForegroundColor Red
        exit 1
    }
}

function Build-Examples {
    Write-Host "Building examples..." -ForegroundColor Green

    if (!(Test-Path $BuildDir)) {
        New-Item -ItemType Directory -Path $BuildDir | Out-Null
    }

    $output = Join-Path $BuildDir "examples.exe"
    $ldflags = "-X main.version=$Version"

    go build $GoFlags -ldflags $ldflags -o $output cmd/examples/main.go

    if ($LASTEXITCODE -eq 0) {
        Write-Host "Examples built: $output" -ForegroundColor Green
    } else {
        Write-Host "Build failed!" -ForegroundColor Red
        exit 1
    }
}

function Build-All {
    Build-Application
    Build-TaskMaster
    Build-Examples
}

function Clean-Build {
    Write-Host "Cleaning build artifacts..." -ForegroundColor Yellow

    if (Test-Path $BuildDir) {
        Remove-Item -Recurse -Force $BuildDir
    }

    if (Test-Path "coverage.out") {
        Remove-Item "coverage.out"
    }

    go clean
    Write-Host "Clean complete" -ForegroundColor Green
}

function Get-Dependencies {
    Write-Host "Downloading dependencies..." -ForegroundColor Green
    go mod download
    go mod verify
    Write-Host "Dependencies downloaded" -ForegroundColor Green
}

function Update-Dependencies {
    Write-Host "Updating dependencies..." -ForegroundColor Green
    go get -u ./...
    go mod tidy
    Write-Host "Dependencies updated" -ForegroundColor Green
}

function Run-Tests {
    Write-Host "Running tests..." -ForegroundColor Green
    go test -v ./...
}

function Run-Coverage {
    Write-Host "Running tests with coverage..." -ForegroundColor Green
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    Write-Host "Coverage report generated: coverage.html" -ForegroundColor Green
}

function Format-Code {
    Write-Host "Formatting code..." -ForegroundColor Green
    go fmt ./...
    Write-Host "Code formatted" -ForegroundColor Green
}

function Run-Application {
    Build-Application
    Write-Host "Running $BinaryName..." -ForegroundColor Green
    & ".\$BuildDir\$BinaryName.exe" --endpoints=localhost:2379 --api-port=8080
}

function Run-Benchmark {
    Build-Application
    Write-Host "Running benchmark..." -ForegroundColor Green
    & ".\$BuildDir\$BinaryName.exe" `
        --run-benchmark `
        --benchmark-type=mixed `
        --benchmark-ops=10000 `
        --endpoints=localhost:2379
}

function Show-Help {
    Write-Host "etcd-monitor Build Script" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Usage: .\build.ps1 -Command <command>" -ForegroundColor White
    Write-Host ""
    Write-Host "Commands:" -ForegroundColor Yellow
    Write-Host "  build           - Build the monitoring application"
    Write-Host "  build-taskmaster - Build TaskMaster CLI"
    Write-Host "  build-examples  - Build pattern examples"
    Write-Host "  build-all       - Build all applications"
    Write-Host "  clean           - Remove build artifacts"
    Write-Host "  deps            - Download dependencies"
    Write-Host "  deps-update     - Update dependencies"
    Write-Host "  test            - Run tests"
    Write-Host "  coverage        - Run tests with coverage"
    Write-Host "  fmt             - Format code"
    Write-Host "  run             - Build and run the monitoring application"
    Write-Host "  benchmark       - Run performance benchmark"
    Write-Host "  help            - Show this help message"
    Write-Host ""
}

# Main execution
switch ($Command.ToLower()) {
    "build" { Build-Application }
    "build-taskmaster" { Build-TaskMaster }
    "build-examples" { Build-Examples }
    "build-all" { Build-All }
    "clean" { Clean-Build }
    "deps" { Get-Dependencies }
    "deps-update" { Update-Dependencies }
    "test" { Run-Tests }
    "coverage" { Run-Coverage }
    "fmt" { Format-Code }
    "run" { Run-Application }
    "benchmark" { Run-Benchmark }
    "help" { Show-Help }
    default {
        Write-Host "Unknown command: $Command" -ForegroundColor Red
        Show-Help
        exit 1
    }
}
