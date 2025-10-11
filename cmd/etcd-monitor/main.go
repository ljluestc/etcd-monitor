// Package main is the entry point for etcd-monitor
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/etcd-monitor/taskmaster/pkg/api"
	"github.com/etcd-monitor/taskmaster/pkg/benchmark"
	"github.com/etcd-monitor/taskmaster/pkg/etcdctl"
	"github.com/etcd-monitor/taskmaster/pkg/monitor"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	"go.uber.org/zap"
)

var (
	// etcd connection flags
	endpoints       = flag.String("endpoints", "localhost:2379", "Comma-separated list of etcd endpoints")
	dialTimeout     = flag.Duration("dial-timeout", 5*time.Second, "Timeout for connecting to etcd")
	certFile        = flag.String("cert", "", "TLS certificate file")
	keyFile         = flag.String("key", "", "TLS key file")
	caFile          = flag.String("cacert", "", "TLS CA certificate file")
	insecureSkipTLS = flag.Bool("insecure-skip-tls-verify", false, "Skip TLS verification")

	// API server flags
	apiHost = flag.String("api-host", "0.0.0.0", "API server host")
	apiPort = flag.Int("api-port", 8080, "API server port")

	// Monitoring flags
	healthCheckInterval = flag.Duration("health-check-interval", 30*time.Second, "Health check interval")
	metricsInterval     = flag.Duration("metrics-interval", 10*time.Second, "Metrics collection interval")

	// Alert thresholds
	maxLatencyMs           = flag.Int("max-latency-ms", 100, "Maximum acceptable latency in milliseconds")
	maxDatabaseSizeMB      = flag.Int("max-db-size-mb", 8192, "Maximum database size in MB")
	minAvailableNodes      = flag.Int("min-available-nodes", 2, "Minimum available nodes")
	maxLeaderChangesPerHour = flag.Int("max-leader-changes", 3, "Maximum leader changes per hour")

	// Benchmark flags
	benchmarkEnabled  = flag.Bool("benchmark-enabled", false, "Enable benchmark testing")
	benchmarkInterval = flag.Duration("benchmark-interval", 1*time.Hour, "Benchmark test interval")

	// Mode flags
	runBenchmark = flag.Bool("run-benchmark", false, "Run a single benchmark and exit")
	benchmarkType = flag.String("benchmark-type", "mixed", "Benchmark type: write, read, mixed")
	benchmarkOps  = flag.Int("benchmark-ops", 10000, "Number of operations for benchmark")

	// Version
	version = flag.Bool("version", false, "Print version and exit")
)

const (
	appVersion = "1.0.0"
	appName    = "etcd-monitor"
)

func main() {
	flag.Parse()

	// Print version and exit
	if *version {
		fmt.Printf("%s version %s\n", appName, appVersion)
		os.Exit(0)
	}

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting etcd-monitor",
		zap.String("version", appVersion),
		zap.String("endpoints", *endpoints))

	// Parse endpoints
	endpointList := parseEndpoints(*endpoints)

	// Create etcd client
	client, err := createEtcdClient(endpointList, logger)
	if err != nil {
		logger.Fatal("Failed to create etcd client", zap.Error(err))
	}
	defer client.Close()

	// If benchmark mode, run benchmark and exit
	if *runBenchmark {
		runBenchmarkMode(client, logger)
		return
	}

	// Create monitor configuration
	monitorConfig := &monitor.Config{
		Endpoints:   endpointList,
		DialTimeout: *dialTimeout,
		HealthCheckInterval: *healthCheckInterval,
		MetricsInterval:     *metricsInterval,
		AlertThresholds: monitor.AlertThresholds{
			MaxLatencyMs:           *maxLatencyMs,
			MaxDatabaseSizeMB:      *maxDatabaseSizeMB,
			MinAvailableNodes:      *minAvailableNodes,
			MaxLeaderChangesPerHour: *maxLeaderChangesPerHour,
			MaxErrorRate:           0.05,
			MinDiskSpacePercent:    10.0,
		},
		BenchmarkEnabled:  *benchmarkEnabled,
		BenchmarkInterval: *benchmarkInterval,
	}

	// Create monitor service
	monitorService, err := monitor.NewMonitorService(monitorConfig, logger)
	if err != nil {
		logger.Fatal("Failed to create monitor service", zap.Error(err))
	}

	// Start monitor service
	if err := monitorService.Start(); err != nil {
		logger.Fatal("Failed to start monitor service", zap.Error(err))
	}
	defer monitorService.Stop()

	// Create API server
	apiConfig := &api.Config{
		Host:    *apiHost,
		Port:    *apiPort,
		Timeout: 30 * time.Second,
	}

	apiServer := api.NewServer(apiConfig, monitorService, logger)

	// Create and register Prometheus exporter
	prometheusExporter := api.NewPrometheusExporter(monitorService, logger)
	prometheusExporter.RegisterWithServer(apiServer)

	// Start API server in a goroutine
	go func() {
		if err := apiServer.Start(); err != nil {
			logger.Error("API server error", zap.Error(err))
		}
	}()

	logger.Info("etcd-monitor started successfully",
		zap.String("api_address", fmt.Sprintf("%s:%d", *apiHost, *apiPort)))

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := apiServer.Stop(ctx); err != nil {
		logger.Error("Error stopping API server", zap.Error(err))
	}

	logger.Info("Shutdown complete")
}

// createEtcdClient creates an etcd client
func createEtcdClient(endpoints []string, logger *zap.Logger) (*clientv3.Client, error) {
	config := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: *dialTimeout,
	}

	// Add TLS configuration if certificates are provided
	if *certFile != "" && *keyFile != "" && *caFile != "" {
		tlsInfo := transport.TLSInfo{
			CertFile:      *certFile,
			KeyFile:       *keyFile,
			TrustedCAFile: *caFile,
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to create TLS config: %w", err)
		}
		if *insecureSkipTLS {
			tlsConfig.InsecureSkipVerify = true
		}
		config.TLS = tlsConfig
		logger.Info("TLS enabled for etcd connection")
	}

	client, err := clientv3.New(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to etcd: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Status(ctx, endpoints[0])
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to get etcd status: %w", err)
	}

	logger.Info("Successfully connected to etcd", zap.Strings("endpoints", endpoints))

	return client, nil
}

// parseEndpoints parses comma-separated endpoints
func parseEndpoints(endpointsStr string) []string {
	endpoints := make([]string, 0)
	for i := 0; i < len(endpointsStr); {
		start := i
		for i < len(endpointsStr) && endpointsStr[i] != ',' {
			i++
		}
		endpoint := endpointsStr[start:i]
		if endpoint != "" {
			endpoints = append(endpoints, endpoint)
		}
		i++
	}
	return endpoints
}

// runBenchmarkMode runs a benchmark and exits
func runBenchmarkMode(client *clientv3.Client, logger *zap.Logger) {
	logger.Info("Running benchmark",
		zap.String("type", *benchmarkType),
		zap.Int("operations", *benchmarkOps))

	// Create benchmark configuration
	benchConfig := &benchmark.Config{
		Type:            benchmark.BenchmarkType(*benchmarkType),
		Connections:     10,
		Clients:         10,
		KeySize:         32,
		ValueSize:       256,
		TotalOperations: *benchmarkOps,
		KeyPrefix:       "/benchmark-test",
		TargetLeader:    false,
		RateLimit:       0,
	}

	// Create benchmark runner
	runner := benchmark.NewRunner(client, benchConfig, logger)

	// Run benchmark
	ctx := context.Background()
	result, err := runner.Run(ctx)
	if err != nil {
		logger.Fatal("Benchmark failed", zap.Error(err))
	}

	// Print results
	benchmark.PrintResult(result, logger)

	fmt.Printf("\n=== Benchmark Results ===\n")
	fmt.Printf("Type:             %s\n", result.Type)
	fmt.Printf("Duration:         %s\n", result.Duration)
	fmt.Printf("Total Operations: %d\n", result.TotalOperations)
	fmt.Printf("Successful:       %d\n", result.SuccessfulOps)
	fmt.Printf("Failed:           %d\n", result.FailedOps)
	fmt.Printf("Throughput:       %.2f ops/sec\n", result.Throughput)
	fmt.Printf("\nLatency:\n")
	fmt.Printf("  Average:        %.2f ms\n", result.AvgLatency)
	fmt.Printf("  P50:            %.2f ms\n", result.P50Latency)
	fmt.Printf("  P95:            %.2f ms\n", result.P95Latency)
	fmt.Printf("  P99:            %.2f ms\n", result.P99Latency)
	fmt.Printf("  Min:            %.2f ms\n", result.MinLatency)
	fmt.Printf("  Max:            %.2f ms\n", result.MaxLatency)
	fmt.Printf("\nLatency Distribution:\n")
	for bucket, count := range result.LatencyHistogram {
		percentage := float64(count) / float64(result.TotalOperations) * 100
		fmt.Printf("  %-12s: %6d (%.2f%%)\n", bucket, count, percentage)
	}
}
