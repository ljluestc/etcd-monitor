package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"github.com/etcd-monitor/taskmaster/pkg/monitor"
)

// MockMonitorService is a mock implementation of MonitorService for testing
type MockMonitorService struct {
	clusterStatus *monitor.ClusterStatus
	metrics       *monitor.MetricsSnapshot
	alertManager  *monitor.AlertManager
	shouldError   bool
}

func (m *MockMonitorService) GetClusterStatus() (*monitor.ClusterStatus, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	return m.clusterStatus, nil
}

func (m *MockMonitorService) GetCurrentMetrics() (*monitor.MetricsSnapshot, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	return m.metrics, nil
}

func (m *MockMonitorService) GetAlertManager() *monitor.AlertManager {
	return m.alertManager
}

func (m *MockMonitorService) GetHealthChecker() *monitor.HealthChecker {
	return nil
}

func (m *MockMonitorService) GetMetricsCollector() *monitor.MetricsCollector {
	return nil
}

func (m *MockMonitorService) IsRunning() bool {
	return true
}

func TestServer_Comprehensive(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("NewServer with valid config", func(t *testing.T) {
		config := &Config{
			Port:    8080,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		assert.NotNil(t, server)
		assert.NotNil(t, server.monitorService)
		assert.NotNil(t, server.logger)
		assert.NotNil(t, server.router)
		assert.NotNil(t, server.server)
	})

	t.Run("NewServer with nil config", func(t *testing.T) {
		monitorService := &MockMonitorService{}
		server := NewServer(nil, monitorService, logger)
		
		assert.NotNil(t, server)
		assert.NotNil(t, server.monitorService)
		assert.NotNil(t, server.logger)
		assert.NotNil(t, server.router)
		assert.Nil(t, server.server) // Should be nil when config is nil
	})

	t.Run("NewServer with nil monitor service", func(t *testing.T) {
		config := &Config{
			Port:    8080,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		server := NewServer(config, nil, logger)
		
		assert.NotNil(t, server)
		assert.Nil(t, server.monitorService)
		assert.NotNil(t, server.logger)
		assert.NotNil(t, server.router)
		assert.NotNil(t, server.server)
	})

	t.Run("NewServer with nil logger", func(t *testing.T) {
		config := &Config{
			Port:    8080,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, nil)
		
		assert.NotNil(t, server)
		assert.NotNil(t, server.monitorService)
		assert.NotNil(t, server.logger) // Should create a production logger
		assert.NotNil(t, server.router)
		assert.NotNil(t, server.server)
	})

	t.Run("Start server", func(t *testing.T) {
		config := &Config{
			Port:    0, // Use random port
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		// Start server in a goroutine since it blocks
		go func() {
			server.Start()
		}()
		
		// Give server time to start
		time.Sleep(100 * time.Millisecond)
		
		// Stop server
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := server.Stop(ctx)
		assert.NoError(t, err)
	})

	t.Run("Start server with invalid port", func(t *testing.T) {
		config := &Config{
			Port:    -1,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		// Start server in a goroutine since it blocks
		go func() {
			err := server.Start()
			assert.Error(t, err)
		}()
		
		// Give server time to fail
		time.Sleep(100 * time.Millisecond)
	})

	t.Run("Start server with invalid host", func(t *testing.T) {
		config := &Config{
			Port:    8080,
			Host:    "invalid host",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		// Start server in a goroutine since it blocks
		go func() {
			err := server.Start()
			assert.Error(t, err)
		}()
		
		// Give server time to fail
		time.Sleep(100 * time.Millisecond)
	})

	t.Run("Start server with invalid timeout", func(t *testing.T) {
		config := &Config{
			Port:    8080,
			Host:    "localhost",
			Timeout: -1 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		// Start server in a goroutine since it blocks
		go func() {
			err := server.Start()
			assert.Error(t, err)
		}()
		
		// Give server time to fail
		time.Sleep(100 * time.Millisecond)
	})

	t.Run("Stop server", func(t *testing.T) {
		config := &Config{
			Port:    0, // Use random port
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		// Start server in a goroutine
		go func() {
			server.Start()
		}()
		
		// Give server time to start
		time.Sleep(100 * time.Millisecond)
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := server.Stop(ctx)
		assert.NoError(t, err)
	})

	t.Run("Stop server without starting", func(t *testing.T) {
		config := &Config{
			Port:    8080,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := server.Stop(ctx)
		assert.NoError(t, err) // Should not error when stopping a non-started server
	})
}

func TestServer_Endpoints(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Health endpoint", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)
		
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "healthy")
	})

	t.Run("Metrics endpoint", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{
			metrics: &monitor.MetricsSnapshot{
				Timestamp:       time.Now(),
				RequestRate:     100.0,
				ReadLatencyP50:  1.0,
				ReadLatencyP95:  5.0,
				ReadLatencyP99:  10.0,
				WriteLatencyP50: 2.0,
				WriteLatencyP95: 8.0,
				WriteLatencyP99: 15.0,
				DBSize:          1024 * 1024,
				DBSizeInUse:     512 * 1024,
			},
		}
		server := NewServer(config, monitorService, logger)
		
		req, err := http.NewRequest("GET", "/api/v1/metrics/current", nil)
		require.NoError(t, err)
		
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		
		var metrics monitor.MetricsSnapshot
		err = json.Unmarshal(rr.Body.Bytes(), &metrics)
		assert.NoError(t, err)
		assert.Equal(t, 100.0, metrics.RequestRate)
		assert.Equal(t, 1.0, metrics.ReadLatencyP50)
	})

	t.Run("Metrics endpoint with error", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{
			shouldError: true,
		}
		server := NewServer(config, monitorService, logger)
		
		req, err := http.NewRequest("GET", "/api/v1/metrics/current", nil)
		require.NoError(t, err)
		
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("Alerts endpoint", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		alertManager := monitor.NewAlertManager(monitor.AlertThresholds{}, logger)
		alertManager.TriggerAlert(monitor.Alert{
			Level:     monitor.AlertLevelCritical,
			Type:      monitor.AlertTypeClusterHealth,
			Message:   "Test alert",
			Timestamp: time.Now(),
		})
		
		monitorService := &MockMonitorService{
			alertManager: alertManager,
		}
		server := NewServer(config, monitorService, logger)
		
		req, err := http.NewRequest("GET", "/api/v1/alerts", nil)
		require.NoError(t, err)
		
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		
		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "active_alerts")
		assert.Contains(t, response, "count")
	})

	t.Run("Alerts endpoint with error", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{
			shouldError: true,
		}
		server := NewServer(config, monitorService, logger)
		
		req, err := http.NewRequest("GET", "/api/v1/alerts", nil)
		require.NoError(t, err)
		
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

}

func TestServer_Middleware(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("CORS middleware", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)
		req.Header.Set("Origin", "http://example.com")
		
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
	})

	t.Run("Logging middleware", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)
		
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		// Logging middleware should not affect response
	})
}

func TestServer_EdgeCases(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Concurrent requests", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		done := make(chan bool, 10)
		
		// Send concurrent requests
		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()
				
				req, err := http.NewRequest("GET", "/health", nil)
				require.NoError(t, err)
				
				rr := httptest.NewRecorder()
				server.router.ServeHTTP(rr, req)
				
				assert.Equal(t, http.StatusOK, rr.Code)
			}(i)
		}
		
		// Wait for all requests to complete
		for i := 0; i < 10; i++ {
			<-done
		}
	})

	t.Run("Invalid HTTP methods", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		req, err := http.NewRequest("PATCH", "/health", nil)
		require.NoError(t, err)
		
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("Non-existent endpoint", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		req, err := http.NewRequest("GET", "/nonexistent", nil)
		require.NoError(t, err)
		
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

func TestServer_Performance(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("High throughput requests", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		start := time.Now()
		
		// Send 1000 requests
		for i := 0; i < 1000; i++ {
			req, err := http.NewRequest("GET", "/health", nil)
			require.NoError(t, err)
			
			rr := httptest.NewRecorder()
			server.router.ServeHTTP(rr, req)
			
			assert.Equal(t, http.StatusOK, rr.Code)
		}
		
		duration := time.Since(start)
		
		// Should handle 1000 requests in less than 5 seconds
		assert.True(t, duration < 5*time.Second, "Handling 1000 requests took too long: %v", duration)
	})
}

func TestServer_Integration(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Full API workflow", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		alertManager := monitor.NewAlertManager(monitor.AlertThresholds{}, logger)
		monitorService := &MockMonitorService{
			alertManager: alertManager,
		}
		server := NewServer(config, monitorService, logger)
		
		// 1. Check health
		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)
		
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		
		// 2. Get metrics (should be empty initially)
		req, err = http.NewRequest("GET", "/api/v1/metrics/current", nil)
		require.NoError(t, err)
		
		rr = httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		
		// 3. Get alerts (should be empty initially)
		req, err = http.NewRequest("GET", "/api/v1/alerts", nil)
		require.NoError(t, err)
		
		rr = httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		
		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "active_alerts")
		assert.Contains(t, response, "count")
	})
}
