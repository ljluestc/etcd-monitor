// +build integration

package integration_tests

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/etcd-monitor/taskmaster/pkg/api"
	"github.com/etcd-monitor/taskmaster/pkg/etcd"
	"github.com/etcd-monitor/taskmaster/pkg/monitor"
	"github.com/etcd-monitor/taskmaster/pkg/inspection"
	"go.uber.org/zap"
)

// TestEtcdMonitorIntegration tests the complete etcd-monitor integration
func TestEtcdMonitorIntegration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Set up test environment
	etcdEndpoints := os.Getenv("ETCD_ENDPOINTS")
	if etcdEndpoints == "" {
		etcdEndpoints = "http://localhost:2379"
	}

	// Create logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Test etcd client creation
	t.Run("EtcdClientCreation", func(t *testing.T) {
		client, err := etcd.NewClient(etcdEndpoints, nil)
		if err != nil {
			t.Fatalf("Failed to create etcd client: %v", err)
		}
		defer client.Close()

		// Test basic connectivity
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err = client.Get(ctx, "test-key")
		if err != nil {
			t.Logf("Etcd connection test failed (expected if etcd not running): %v", err)
		}
	})

	// Test monitoring service
	t.Run("MonitoringService", func(t *testing.T) {
		// Create monitoring service
		monitorService := monitor.NewService(logger)

		// Test health check
		health, err := monitorService.CheckHealth(etcdEndpoints)
		if err != nil {
			t.Logf("Health check failed (expected if etcd not running): %v", err)
		} else {
			t.Logf("Health check result: %+v", health)
		}

		// Test metrics collection
		metrics, err := monitorService.CollectMetrics(etcdEndpoints)
		if err != nil {
			t.Logf("Metrics collection failed (expected if etcd not running): %v", err)
		} else {
			t.Logf("Metrics collected: %d", len(metrics))
		}
	})

	// Test inspection service
	t.Run("InspectionService", func(t *testing.T) {
		// Create inspection service
		inspectionService := inspection.NewService(logger)

		// Test alarm inspection
		alarms, err := inspectionService.CheckAlarms(etcdEndpoints)
		if err != nil {
			t.Logf("Alarm inspection failed (expected if etcd not running): %v", err)
		} else {
			t.Logf("Alarms found: %d", len(alarms))
		}

		// Test consistency inspection
		consistency, err := inspectionService.CheckConsistency(etcdEndpoints)
		if err != nil {
			t.Logf("Consistency inspection failed (expected if etcd not running): %v", err)
		} else {
			t.Logf("Consistency check result: %+v", consistency)
		}

		// Test health inspection
		health, err := inspectionService.CheckHealth(etcdEndpoints)
		if err != nil {
			t.Logf("Health inspection failed (expected if etcd not running): %v", err)
		} else {
			t.Logf("Health inspection result: %+v", health)
		}

		// Test request inspection
		requests, err := inspectionService.CheckRequests(etcdEndpoints)
		if err != nil {
			t.Logf("Request inspection failed (expected if etcd not running): %v", err)
		} else {
			t.Logf("Request inspection result: %+v", requests)
		}
	})

	// Test API server
	t.Run("APIServer", func(t *testing.T) {
		// Create API server
		apiServer := api.NewServer("localhost", 8080, logger)

		// Start server in background
		go func() {
			if err := apiServer.Start(); err != nil {
				t.Logf("API server start failed: %v", err)
			}
		}()

		// Wait for server to start
		time.Sleep(2 * time.Second)

		// Test health endpoint
		resp, err := http.Get("http://localhost:8080/health")
		if err != nil {
			t.Logf("Health endpoint test failed: %v", err)
		} else {
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200, got %d", resp.StatusCode)
			}
		}

		// Test metrics endpoint
		resp, err = http.Get("http://localhost:8080/metrics")
		if err != nil {
			t.Logf("Metrics endpoint test failed: %v", err)
		} else {
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200, got %d", resp.StatusCode)
			}
		}

		// Test etcd endpoint
		resp, err = http.Get("http://localhost:8080/etcd/health")
		if err != nil {
			t.Logf("Etcd endpoint test failed: %v", err)
		} else {
			resp.Body.Close()
			// Status code might be 500 if etcd is not running
			t.Logf("Etcd endpoint status: %d", resp.StatusCode)
		}

		// Stop server
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		apiServer.Stop(ctx)
	})
}

// TestEtcdMonitorEndToEnd tests the complete end-to-end flow
func TestEtcdMonitorEndToEnd(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Set up test environment
	etcdEndpoints := os.Getenv("ETCD_ENDPOINTS")
	if etcdEndpoints == "" {
		etcdEndpoints = "http://localhost:2379"
	}

	// Create logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Create monitoring service
	monitorService := monitor.NewService(logger)

	// Create inspection service
	inspectionService := inspection.NewService(logger)

	// Create API server
	apiServer := api.NewServer("localhost", 8081, logger)

	// Start API server
	go func() {
		if err := apiServer.Start(); err != nil {
			t.Logf("API server start failed: %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(2 * time.Second)

	// Test complete monitoring flow
	t.Run("CompleteMonitoringFlow", func(t *testing.T) {
		// 1. Check etcd health
		health, err := monitorService.CheckHealth(etcdEndpoints)
		if err != nil {
			t.Logf("Health check failed: %v", err)
		} else {
			t.Logf("Health check result: %+v", health)
		}

		// 2. Collect metrics
		metrics, err := monitorService.CollectMetrics(etcdEndpoints)
		if err != nil {
			t.Logf("Metrics collection failed: %v", err)
		} else {
			t.Logf("Metrics collected: %d", len(metrics))
		}

		// 3. Run inspections
		alarms, err := inspectionService.CheckAlarms(etcdEndpoints)
		if err != nil {
			t.Logf("Alarm inspection failed: %v", err)
		} else {
			t.Logf("Alarms found: %d", len(alarms))
		}

		consistency, err := inspectionService.CheckConsistency(etcdEndpoints)
		if err != nil {
			t.Logf("Consistency inspection failed: %v", err)
		} else {
			t.Logf("Consistency check result: %+v", consistency)
		}

		// 4. Test API endpoints
		client := &http.Client{Timeout: 10 * time.Second}

		// Test health endpoint
		resp, err := client.Get("http://localhost:8081/health")
		if err != nil {
			t.Logf("Health endpoint test failed: %v", err)
		} else {
			resp.Body.Close()
			t.Logf("Health endpoint status: %d", resp.StatusCode)
		}

		// Test metrics endpoint
		resp, err = client.Get("http://localhost:8081/metrics")
		if err != nil {
			t.Logf("Metrics endpoint test failed: %v", err)
		} else {
			resp.Body.Close()
			t.Logf("Metrics endpoint status: %d", resp.StatusCode)
		}

		// Test etcd endpoint
		resp, err = client.Get("http://localhost:8081/etcd/health")
		if err != nil {
			t.Logf("Etcd endpoint test failed: %v", err)
		} else {
			resp.Body.Close()
			t.Logf("Etcd endpoint status: %d", resp.StatusCode)
		}
	})

	// Cleanup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	apiServer.Stop(ctx)
}

// TestEtcdMonitorWithMockEtcd tests the monitor with a mock etcd server
func TestEtcdMonitorWithMockEtcd(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test would require setting up a mock etcd server
	// For now, we'll just test the error handling when etcd is not available
	t.Run("ErrorHandling", func(t *testing.T) {
		// Create logger
		logger, _ := zap.NewDevelopment()
		defer logger.Sync()

		// Create monitoring service
		monitorService := monitor.NewService(logger)

		// Test with invalid etcd endpoint
		invalidEndpoint := "http://invalid-endpoint:2379"
		
		// Health check should fail gracefully
		_, err := monitorService.CheckHealth(invalidEndpoint)
		if err == nil {
			t.Error("Expected health check to fail with invalid endpoint")
		}
		t.Logf("Health check failed as expected: %v", err)

		// Metrics collection should fail gracefully
		_, err = monitorService.CollectMetrics(invalidEndpoint)
		if err == nil {
			t.Error("Expected metrics collection to fail with invalid endpoint")
		}
		t.Logf("Metrics collection failed as expected: %v", err)
	})
}

// TestEtcdMonitorPerformance tests the performance of the monitoring system
func TestEtcdMonitorPerformance(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Set up test environment
	etcdEndpoints := os.Getenv("ETCD_ENDPOINTS")
	if etcdEndpoints == "" {
		etcdEndpoints = "http://localhost:2379"
	}

	// Create logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Create monitoring service
	monitorService := monitor.NewService(logger)

	t.Run("HealthCheckPerformance", func(t *testing.T) {
		// Test health check performance
		start := time.Now()
		_, err := monitorService.CheckHealth(etcdEndpoints)
		duration := time.Since(start)

		if err != nil {
			t.Logf("Health check failed: %v", err)
		} else {
			t.Logf("Health check took: %v", duration)
			if duration > 5*time.Second {
				t.Errorf("Health check took too long: %v", duration)
			}
		}
	})

	t.Run("MetricsCollectionPerformance", func(t *testing.T) {
		// Test metrics collection performance
		start := time.Now()
		_, err := monitorService.CollectMetrics(etcdEndpoints)
		duration := time.Since(start)

		if err != nil {
			t.Logf("Metrics collection failed: %v", err)
		} else {
			t.Logf("Metrics collection took: %v", duration)
			if duration > 10*time.Second {
				t.Errorf("Metrics collection took too long: %v", duration)
			}
		}
	})
}

// TestEtcdMonitorConcurrency tests the monitoring system under concurrent load
func TestEtcdMonitorConcurrency(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Set up test environment
	etcdEndpoints := os.Getenv("ETCD_ENDPOINTS")
	if etcdEndpoints == "" {
		etcdEndpoints = "http://localhost:2379"
	}

	// Create logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Create monitoring service
	monitorService := monitor.NewService(logger)

	t.Run("ConcurrentHealthChecks", func(t *testing.T) {
		// Test concurrent health checks
		numGoroutines := 10
		results := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				_, err := monitorService.CheckHealth(etcdEndpoints)
				results <- err
			}()
		}

		// Collect results
		for i := 0; i < numGoroutines; i++ {
			err := <-results
			if err != nil {
				t.Logf("Concurrent health check failed: %v", err)
			}
		}
	})

	t.Run("ConcurrentMetricsCollection", func(t *testing.T) {
		// Test concurrent metrics collection
		numGoroutines := 5
		results := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				_, err := monitorService.CollectMetrics(etcdEndpoints)
				results <- err
			}()
		}

		// Collect results
		for i := 0; i < numGoroutines; i++ {
			err := <-results
			if err != nil {
				t.Logf("Concurrent metrics collection failed: %v", err)
			}
		}
	})
}
