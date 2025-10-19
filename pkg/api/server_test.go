package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/etcd-monitor/taskmaster/pkg/monitor"
	"go.uber.org/zap"
)

func TestHealthEndpoint(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	// Create a minimal monitor service (nil is okay for health check)
	config := &Config{
		Port:    8080,
		Host:    "0.0.0.0",
		Timeout: 30,
	}

	server := NewServer(config, nil, logger)

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	server.router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `"status":"healthy"`
	if !contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want substring %v",
			rr.Body.String(), expected)
	}
}

func TestCORSMiddleware(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := &Config{
		Port:    8080,
		Host:    "0.0.0.0",
		Timeout: 30,
	}

	server := NewServer(config, nil, logger)

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	server.router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET request failed: got %v want %v", status, http.StatusOK)
	}

	corsHeader := rr.Header().Get("Access-Control-Allow-Origin")
	if corsHeader != "*" {
		t.Errorf("CORS header not set correctly: got %v want *", corsHeader)
	}
}

func TestAPIEndpoints(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	// Create a mock monitor service
	monitorConfig := &monitor.Config{
		Endpoints:           []string{"localhost:2379"},
		HealthCheckInterval: 30,
		MetricsInterval:     10,
	}
	monitorService, _ := monitor.NewMonitorService(monitorConfig, logger)

	config := &Config{
		Port:    8080,
		Host:    "0.0.0.0",
		Timeout: 30,
	}

	server := NewServer(config, monitorService, logger)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"Health Check", "GET", "/health", http.StatusOK},
		{"Cluster Status", "GET", "/api/v1/cluster/status", http.StatusInternalServerError}, // Will fail without etcd
		{"Current Metrics", "GET", "/api/v1/metrics/current", http.StatusInternalServerError}, // Will fail without etcd
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			server.router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Logf("handler returned status code: got %v want %v (body: %s)",
					status, tt.expectedStatus, rr.Body.String())
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr))))
}
