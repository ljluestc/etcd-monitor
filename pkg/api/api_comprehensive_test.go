package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"github.com/etcd-monitor/taskmaster/pkg/monitor"
)

// MockMonitorService is a mock implementation of MonitorService for testing
type MockMonitorService struct {
	metrics     *monitor.Metrics
	alerts      []monitor.Alert
	shouldError bool
}

func (m *MockMonitorService) GetMetrics() *monitor.Metrics {
	if m.shouldError {
		return nil
	}
	return m.metrics
}

func (m *MockMonitorService) GetAlerts() []monitor.Alert {
	if m.shouldError {
		return nil
	}
	return m.alerts
}

func (m *MockMonitorService) TriggerAlert(alert monitor.Alert) {
	if !m.shouldError {
		m.alerts = append(m.alerts, alert)
	}
}

func (m *MockMonitorService) ClearAlert(alertType monitor.AlertType, message string) {
	if !m.shouldError {
		for i, alert := range m.alerts {
			if alert.Type == alertType && alert.Message == message {
				m.alerts = append(m.alerts[:i], m.alerts[i+1:]...)
				break
			}
		}
	}
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
		assert.NotNil(t, server.config)
		assert.NotNil(t, server.monitorService)
		assert.NotNil(t, server.logger)
		assert.NotNil(t, server.router)
	})

	t.Run("NewServer with nil config", func(t *testing.T) {
		monitorService := &MockMonitorService{}
		server := NewServer(nil, monitorService, logger)
		
		assert.NotNil(t, server)
		assert.Nil(t, server.config)
		assert.NotNil(t, server.monitorService)
		assert.NotNil(t, server.logger)
		assert.NotNil(t, server.router)
	})

	t.Run("NewServer with nil monitor service", func(t *testing.T) {
		config := &Config{
			Port:    8080,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		server := NewServer(config, nil, logger)
		
		assert.NotNil(t, server)
		assert.NotNil(t, server.config)
		assert.Nil(t, server.monitorService)
		assert.NotNil(t, server.logger)
		assert.NotNil(t, server.router)
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
		assert.NotNil(t, server.config)
		assert.NotNil(t, server.monitorService)
		assert.NotNil(t, server.logger) // Should create a production logger
		assert.NotNil(t, server.router)
	})

	t.Run("Start server", func(t *testing.T) {
		config := &Config{
			Port:    0, // Use random port
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		err := server.Start()
		assert.NoError(t, err)
		
		// Stop server
		server.Stop()
	})

	t.Run("Start server with invalid port", func(t *testing.T) {
		config := &Config{
			Port:    -1,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		err := server.Start()
		assert.Error(t, err)
	})

	t.Run("Start server with invalid host", func(t *testing.T) {
		config := &Config{
			Port:    8080,
			Host:    "invalid host",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		err := server.Start()
		assert.Error(t, err)
	})

	t.Run("Start server with invalid timeout", func(t *testing.T) {
		config := &Config{
			Port:    8080,
			Host:    "localhost",
			Timeout: -1 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		err := server.Start()
		assert.Error(t, err)
	})

	t.Run("Stop server", func(t *testing.T) {
		config := &Config{
			Port:    0, // Use random port
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		err := server.Start()
		require.NoError(t, err)
		
		server.Stop()
		// Should not panic or error
	})

	t.Run("Stop server without starting", func(t *testing.T) {
		config := &Config{
			Port:    8080,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		server.Stop()
		// Should not panic or error
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
			metrics: &monitor.Metrics{
				ClusterSize: 3,
				LeaderID:    1,
				RaftIndex:   100,
				RaftTerm:    1,
			},
		}
		server := NewServer(config, monitorService, logger)
		
		req, err := http.NewRequest("GET", "/metrics", nil)
		require.NoError(t, err)
		
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		
		var metrics monitor.Metrics
		err = json.Unmarshal(rr.Body.Bytes(), &metrics)
		assert.NoError(t, err)
		assert.Equal(t, 3, metrics.ClusterSize)
		assert.Equal(t, 1, metrics.LeaderID)
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
		
		req, err := http.NewRequest("GET", "/metrics", nil)
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
		monitorService := &MockMonitorService{
			alerts: []monitor.Alert{
				{
					Level:     monitor.AlertLevelCritical,
					Type:      monitor.AlertTypeClusterHealth,
					Message:   "Test alert",
					Timestamp: time.Now(),
				},
			},
		}
		server := NewServer(config, monitorService, logger)
		
		req, err := http.NewRequest("GET", "/alerts", nil)
		require.NoError(t, err)
		
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		
		var alerts []monitor.Alert
		err = json.Unmarshal(rr.Body.Bytes(), &alerts)
		assert.NoError(t, err)
		assert.Len(t, alerts, 1)
		assert.Equal(t, monitor.AlertLevelCritical, alerts[0].Level)
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
		
		req, err := http.NewRequest("GET", "/alerts", nil)
		require.NoError(t, err)
		
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("Trigger alert endpoint", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		alert := monitor.Alert{
			Level:     monitor.AlertLevelWarning,
			Type:      monitor.AlertTypeHighLatency,
			Message:   "Test alert",
			Timestamp: time.Now(),
		}
		
		alertJSON, err := json.Marshal(alert)
		require.NoError(t, err)
		
		req, err := http.NewRequest("POST", "/alerts", bytes.NewBuffer(alertJSON))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Len(t, monitorService.alerts, 1)
	})

	t.Run("Trigger alert endpoint with invalid JSON", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		req, err := http.NewRequest("POST", "/alerts", bytes.NewBuffer([]byte("invalid json")))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Clear alert endpoint", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{
			alerts: []monitor.Alert{
				{
					Level:     monitor.AlertLevelWarning,
					Type:      monitor.AlertTypeHighLatency,
					Message:   "Test alert",
					Timestamp: time.Now(),
				},
			},
		}
		server := NewServer(config, monitorService, logger)
		
		req, err := http.NewRequest("DELETE", "/alerts/high_latency/Test%20alert", nil)
		require.NoError(t, err)
		
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Len(t, monitorService.alerts, 0)
	})

	t.Run("Clear alert endpoint with non-existent alert", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		req, err := http.NewRequest("DELETE", "/alerts/high_latency/Non-existent", nil)
		require.NoError(t, err)
		
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Len(t, monitorService.alerts, 0)
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

	t.Run("CORS preflight request", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		req, err := http.NewRequest("OPTIONS", "/health", nil)
		require.NoError(t, err)
		req.Header.Set("Origin", "http://example.com")
		req.Header.Set("Access-Control-Request-Method", "GET")
		
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

	t.Run("Recovery middleware", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		// Add a route that panics
		server.router.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		}).Methods("GET")
		
		req, err := http.NewRequest("GET", "/panic", nil)
		require.NoError(t, err)
		
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Internal Server Error")
	})
}

func TestServer_EdgeCases(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Large request body", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		// Create a large alert with many details
		details := make(map[string]interface{})
		for i := 0; i < 1000; i++ {
			details[fmt.Sprintf("key%d", i)] = fmt.Sprintf("value%d", i)
		}
		
		alert := monitor.Alert{
			Level:     monitor.AlertLevelInfo,
			Type:      monitor.AlertTypeEtcdAlarm,
			Message:   "Large alert",
			Details:   details,
			Timestamp: time.Now(),
		}
		
		alertJSON, err := json.Marshal(alert)
		require.NoError(t, err)
		
		req, err := http.NewRequest("POST", "/alerts", bytes.NewBuffer(alertJSON))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Len(t, monitorService.alerts, 1)
	})

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

	t.Run("Memory usage under load", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		// Send many requests with large payloads
		for i := 0; i < 100; i++ {
			details := make(map[string]interface{})
			for j := 0; j < 100; j++ {
				details[fmt.Sprintf("key%d", j)] = fmt.Sprintf("value%d", j)
			}
			
			alert := monitor.Alert{
				Level:     monitor.AlertLevelInfo,
				Type:      monitor.AlertTypeEtcdAlarm,
				Message:   fmt.Sprintf("Alert %d", i),
				Details:   details,
				Timestamp: time.Now(),
			}
			
			alertJSON, err := json.Marshal(alert)
			require.NoError(t, err)
			
			req, err := http.NewRequest("POST", "/alerts", bytes.NewBuffer(alertJSON))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			
			rr := httptest.NewRecorder()
			server.router.ServeHTTP(rr, req)
			
			assert.Equal(t, http.StatusOK, rr.Code)
		}
		
		assert.Len(t, monitorService.alerts, 100)
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
		monitorService := &MockMonitorService{}
		server := NewServer(config, monitorService, logger)
		
		// 1. Check health
		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)
		
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		
		// 2. Get metrics (should be empty initially)
		req, err = http.NewRequest("GET", "/metrics", nil)
		require.NoError(t, err)
		
		rr = httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		
		// 3. Get alerts (should be empty initially)
		req, err = http.NewRequest("GET", "/alerts", nil)
		require.NoError(t, err)
		
		rr = httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		
		var alerts []monitor.Alert
		err = json.Unmarshal(rr.Body.Bytes(), &alerts)
		assert.NoError(t, err)
		assert.Len(t, alerts, 0)
		
		// 4. Trigger an alert
		alert := monitor.Alert{
			Level:     monitor.AlertLevelCritical,
			Type:      monitor.AlertTypeClusterHealth,
			Message:   "Cluster health issue",
			Details:   map[string]interface{}{"node": "node1", "status": "down"},
			Timestamp: time.Now(),
		}
		
		alertJSON, err := json.Marshal(alert)
		require.NoError(t, err)
		
		req, err = http.NewRequest("POST", "/alerts", bytes.NewBuffer(alertJSON))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		
		rr = httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		
		// 5. Verify alert is active
		req, err = http.NewRequest("GET", "/alerts", nil)
		require.NoError(t, err)
		
		rr = httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		
		err = json.Unmarshal(rr.Body.Bytes(), &alerts)
		assert.NoError(t, err)
		assert.Len(t, alerts, 1)
		assert.Equal(t, alert.Message, alerts[0].Message)
		
		// 6. Clear the alert
		req, err = http.NewRequest("DELETE", "/alerts/cluster_health/Cluster%20health%20issue", nil)
		require.NoError(t, err)
		
		rr = httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		
		// 7. Verify alert is cleared
		req, err = http.NewRequest("GET", "/alerts", nil)
		require.NoError(t, err)
		
		rr = httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusOK, rr.Code)
		
		err = json.Unmarshal(rr.Body.Bytes(), &alerts)
		assert.NoError(t, err)
		assert.Len(t, alerts, 0)
	})
}
