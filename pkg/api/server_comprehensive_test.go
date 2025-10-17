package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/etcd-monitor/taskmaster/pkg/monitor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// MockMonitorService is a mock implementation of MonitorService for testing
type MockMonitorService struct {
	clusterStatus    *monitor.ClusterStatus
	clusterStatusErr error
	metrics          *monitor.MetricsSnapshot
	metricsErr       error
	healthChecker    *MockHealthChecker
	alertManager     *MockAlertManager
}

func (m *MockMonitorService) GetClusterStatus() (*monitor.ClusterStatus, error) {
	return m.clusterStatus, m.clusterStatusErr
}

func (m *MockMonitorService) GetCurrentMetrics() (*monitor.MetricsSnapshot, error) {
	return m.metrics, m.metricsErr
}

func (m *MockMonitorService) GetHealthChecker() *MockHealthChecker {
	return m.healthChecker
}

func (m *MockMonitorService) GetAlertManager() *MockAlertManager {
	return m.alertManager
}

// MockHealthChecker is a mock implementation of HealthChecker for testing
type MockHealthChecker struct {
	memberList    []*monitor.MemberInfo
	memberListErr error
}

func (m *MockHealthChecker) GetMemberList(ctx context.Context) ([]*monitor.MemberInfo, error) {
	return m.memberList, m.memberListErr
}

// MockAlertManager is a mock implementation of AlertManager for testing
type MockAlertManager struct {
	activeAlerts []monitor.Alert
	alertHistory []monitor.Alert
}

func (m *MockAlertManager) GetActiveAlerts() []monitor.Alert {
	return m.activeAlerts
}

func (m *MockAlertManager) GetAlertHistory() []monitor.Alert {
	return m.alertHistory
}

func TestServer_Comprehensive(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Server Creation", func(t *testing.T) {
		config := &Config{
			Port:    8080,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}

		server := NewServer(config, nil, logger)
		assert.NotNil(t, server)
		assert.NotNil(t, server.router)
		assert.NotNil(t, server.server)
		assert.Equal(t, logger, server.logger)
	})

	t.Run("Health Endpoint", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		server := NewServer(config, nil, logger)

		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "healthy", response["status"])
		assert.Equal(t, "etcd-monitor-api", response["service"])
		assert.NotEmpty(t, response["timestamp"])
	})

	t.Run("Cluster Status Endpoint - Success", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		mockService := &MockMonitorService{
			clusterStatus: &monitor.ClusterStatus{
				Healthy:     true,
				LeaderID:    1,
				MemberCount: 3,
				QuorumSize:  2,
				HasLeader:   true,
				LastCheck:   time.Now(),
			},
		}
		server := NewServer(config, mockService, logger)

		req, err := http.NewRequest("GET", "/api/v1/cluster/status", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response monitor.ClusterStatus
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response.Healthy)
		assert.Equal(t, uint64(1), response.LeaderID)
		assert.Equal(t, 3, response.MemberCount)
	})

	t.Run("Cluster Status Endpoint - Error", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		mockService := &MockMonitorService{
			clusterStatusErr: fmt.Errorf("cluster unavailable"),
		}
		server := NewServer(config, mockService, logger)

		req, err := http.NewRequest("GET", "/api/v1/cluster/status", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Failed to get cluster status", response["error"])
	})

	t.Run("Cluster Members Endpoint - Success", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		mockHealthChecker := &MockHealthChecker{
			memberList: []*monitor.MemberInfo{
				{ID: 1, Name: "member1", PeerURLs: []string{"http://localhost:2380"}},
				{ID: 2, Name: "member2", PeerURLs: []string{"http://localhost:2381"}},
			},
		}
		mockService := &MockMonitorService{
			clusterStatus: &monitor.ClusterStatus{
				MemberCount: 2,
				QuorumSize:  2,
			},
			healthChecker: mockHealthChecker,
		}
		server := NewServer(config, mockService, logger)

		req, err := http.NewRequest("GET", "/api/v1/cluster/members", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 2, response["member_count"])
		assert.Equal(t, 2, response["quorum_size"])
		assert.NotNil(t, response["members"])
	})

	t.Run("Cluster Members Endpoint - Health Checker Error", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		mockService := &MockMonitorService{
			clusterStatus: &monitor.ClusterStatus{
				MemberCount: 2,
				QuorumSize:  2,
			},
			healthChecker: nil, // No health checker
		}
		server := NewServer(config, mockService, logger)

		req, err := http.NewRequest("GET", "/api/v1/cluster/members", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("Cluster Leader Endpoint", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		mockService := &MockMonitorService{
			clusterStatus: &monitor.ClusterStatus{
				LeaderID:        1,
				HasLeader:       true,
				LastLeaderChange: time.Now(),
			},
		}
		server := NewServer(config, mockService, logger)

		req, err := http.NewRequest("GET", "/api/v1/cluster/leader", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), response["leader_id"])
		assert.True(t, response["has_leader"].(bool))
	})

	t.Run("Current Metrics Endpoint - Success", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		mockService := &MockMonitorService{
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
		server := NewServer(config, mockService, logger)

		req, err := http.NewRequest("GET", "/api/v1/metrics/current", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response monitor.MetricsSnapshot
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 100.0, response.RequestRate)
		assert.Equal(t, 1.0, response.ReadLatencyP50)
	})

	t.Run("Current Metrics Endpoint - Error", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		mockService := &MockMonitorService{
			metricsErr: fmt.Errorf("metrics unavailable"),
		}
		server := NewServer(config, mockService, logger)

		req, err := http.NewRequest("GET", "/api/v1/metrics/current", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("Latency Metrics Endpoint", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		mockService := &MockMonitorService{
			metrics: &monitor.MetricsSnapshot{
				Timestamp:       time.Now(),
				ReadLatencyP50:  1.0,
				ReadLatencyP95:  5.0,
				ReadLatencyP99:  10.0,
				WriteLatencyP50: 2.0,
				WriteLatencyP95: 8.0,
				WriteLatencyP99: 15.0,
			},
		}
		server := NewServer(config, mockService, logger)

		req, err := http.NewRequest("GET", "/api/v1/metrics/latency", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.NotNil(t, response["read_latency"])
		assert.NotNil(t, response["write_latency"])
	})

	t.Run("Metrics History Endpoint", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		server := NewServer(config, nil, logger)

		req, err := http.NewRequest("GET", "/api/v1/metrics/history", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotImplemented, rr.Code)
	})

	t.Run("Alerts Endpoint - Success", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		mockAlertManager := &MockAlertManager{
			activeAlerts: []monitor.Alert{
				{
					Level:     monitor.AlertLevelWarning,
					Type:      monitor.AlertTypeHighLatency,
					Message:   "High latency detected",
					Timestamp: time.Now(),
				},
			},
		}
		mockService := &MockMonitorService{
			alertManager: mockAlertManager,
		}
		server := NewServer(config, mockService, logger)

		req, err := http.NewRequest("GET", "/api/v1/alerts", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 1, response["count"])
		assert.NotNil(t, response["active_alerts"])
	})

	t.Run("Alerts Endpoint - No Alert Manager", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		mockService := &MockMonitorService{
			alertManager: nil,
		}
		server := NewServer(config, mockService, logger)

		req, err := http.NewRequest("GET", "/api/v1/alerts", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("Alert History Endpoint", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		mockAlertManager := &MockAlertManager{
			alertHistory: []monitor.Alert{
				{
					Level:     monitor.AlertLevelInfo,
					Type:      monitor.AlertTypeClusterHealth,
					Message:   "Cluster healthy",
					Timestamp: time.Now(),
				},
			},
		}
		mockService := &MockMonitorService{
			alertManager: mockAlertManager,
		}
		server := NewServer(config, mockService, logger)

		req, err := http.NewRequest("GET", "/api/v1/alerts/history", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 1, response["count"])
		assert.NotNil(t, response["history"])
	})

	t.Run("Benchmark Endpoint - Valid Request", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		server := NewServer(config, nil, logger)

		requestBody := map[string]interface{}{
			"operations": 1000,
		}
		jsonBody, _ := json.Marshal(requestBody)

		req, err := http.NewRequest("POST", "/api/v1/performance/benchmark", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotImplemented, rr.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Benchmark - not yet implemented", response["message"])
		assert.Equal(t, 1000, response["operations"])
	})

	t.Run("Benchmark Endpoint - Invalid Request", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		server := NewServer(config, nil, logger)

		req, err := http.NewRequest("POST", "/api/v1/performance/benchmark", bytes.NewBufferString("invalid json"))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Benchmark Endpoint - Default Operations", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		server := NewServer(config, nil, logger)

		requestBody := map[string]interface{}{}
		jsonBody, _ := json.Marshal(requestBody)

		req, err := http.NewRequest("POST", "/api/v1/performance/benchmark", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotImplemented, rr.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 100, response["operations"]) // Default value
	})

	t.Run("CORS Middleware", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		server := NewServer(config, nil, logger)

		req, err := http.NewRequest("OPTIONS", "/health", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
	})

	t.Run("Logging Middleware", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		server := NewServer(config, nil, logger)

		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		// Logging middleware should not affect the response
	})

	t.Run("Server Start and Stop", func(t *testing.T) {
		config := &Config{Port: 0, Host: "localhost", Timeout: 30 * time.Second} // Port 0 for random port
		server := NewServer(config, nil, logger)

		// Test start in goroutine
		go func() {
			server.Start()
		}()

		// Give server time to start
		time.Sleep(100 * time.Millisecond)

		// Test stop
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := server.Stop(ctx)
		assert.NoError(t, err)
	})

	t.Run("Response Writer", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rw := &responseWriter{ResponseWriter: rr, statusCode: http.StatusOK}

		// Test WriteHeader
		rw.WriteHeader(http.StatusNotFound)
		assert.Equal(t, http.StatusNotFound, rw.statusCode)

		// Test Write
		_, err := rw.Write([]byte("test"))
		assert.NoError(t, err)
		assert.Equal(t, "test", rr.Body.String())
	})

	t.Run("Error Handling", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		server := NewServer(config, nil, logger)

		// Test with invalid JSON in writeJSON
		w := httptest.NewRecorder()
		invalidData := make(chan int) // Channels cannot be JSON marshaled
		server.writeJSON(w, http.StatusOK, invalidData)

		// Should handle JSON encoding error gracefully
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Concurrent Requests", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		mockService := &MockMonitorService{
			clusterStatus: &monitor.ClusterStatus{
				Healthy:     true,
				MemberCount: 3,
				QuorumSize:  2,
			},
		}
		server := NewServer(config, mockService, logger)

		// Test concurrent requests
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				defer func() { done <- true }()
				
				req, err := http.NewRequest("GET", "/health", nil)
				require.NoError(t, err)

				rr := httptest.NewRecorder()
				server.router.ServeHTTP(rr, req)

				assert.Equal(t, http.StatusOK, rr.Code)
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
	})

	t.Run("Invalid Routes", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		server := NewServer(config, nil, logger)

		req, err := http.NewRequest("GET", "/invalid/route", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Method Not Allowed", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		server := NewServer(config, nil, logger)

		req, err := http.NewRequest("POST", "/health", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("Large Request Body", func(t *testing.T) {
		config := &Config{Port: 8080, Host: "localhost", Timeout: 30 * time.Second}
		server := NewServer(config, nil, logger)

		// Create a large request body
		largeBody := strings.Repeat("a", 10000)
		requestBody := map[string]interface{}{
			"operations": 1000,
			"data":       largeBody,
		}
		jsonBody, _ := json.Marshal(requestBody)

		req, err := http.NewRequest("POST", "/api/v1/performance/benchmark", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		// Should handle large request body gracefully
		assert.Equal(t, http.StatusNotImplemented, rr.Code)
	})

	t.Run("Edge Cases", func(t *testing.T) {
		// Test with nil config
		server := NewServer(nil, nil, logger)
		assert.NotNil(t, server)

		// Test with zero timeout
		config := &Config{Port: 8080, Host: "localhost", Timeout: 0}
		server = NewServer(config, nil, logger)
		assert.NotNil(t, server)

		// Test with empty host
		config = &Config{Port: 8080, Host: "", Timeout: 30 * time.Second}
		server = NewServer(config, nil, logger)
		assert.NotNil(t, server)
	})
}
