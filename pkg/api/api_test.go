package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/etcd-monitor/taskmaster/pkg/monitor"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestApi_Comprehensive(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("NewServer", func(t *testing.T) {
		config := &Config{
			Port:    8080,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		
		// Create a mock monitor service
		monitorService := &monitor.MonitorService{}
		
		server := NewServer(config, monitorService, logger)
		
		assert.NotNil(t, server)
		assert.Equal(t, monitorService, server.monitorService)
		assert.Equal(t, logger, server.logger)
	})

	t.Run("NewServer with nil logger", func(t *testing.T) {
		config := &Config{
			Port:    8080,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		
		monitorService := &monitor.MonitorService{}
		server := NewServer(config, monitorService, nil)
		
		assert.NotNil(t, server)
		assert.NotNil(t, server.logger) // Should create a production logger
	})

	t.Run("Server structure", func(t *testing.T) {
		config := &Config{
			Port:    8080,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		
		monitorService := &monitor.MonitorService{}
		server := NewServer(config, monitorService, logger)
		
		assert.NotNil(t, server.router)
		assert.NotNil(t, server.server)
		assert.Equal(t, monitorService, server.monitorService)
		assert.Equal(t, logger, server.logger)
	})

	t.Run("Start server", func(t *testing.T) {
		config := &Config{
			Port:    0, // Use port 0 for testing
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		
		monitorService := &monitor.MonitorService{}
		server := NewServer(config, monitorService, logger)
		
		// Start server in a goroutine
		go func() {
			err := server.Start()
			assert.NoError(t, err)
		}()
		
		// Give server time to start
		time.Sleep(100 * time.Millisecond)
		
		// Stop server
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := server.Stop(ctx)
		assert.NoError(t, err)
	})

	t.Run("Stop server", func(t *testing.T) {
		config := &Config{
			Port:    0, // Use port 0 for testing
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		
		monitorService := &monitor.MonitorService{}
		server := NewServer(config, monitorService, logger)
		
		// Start server
		go func() {
			server.Start()
		}()
		
		time.Sleep(100 * time.Millisecond)
		
		// Stop server
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := server.Stop(ctx)
		assert.NoError(t, err)
	})

	t.Run("Health endpoint", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		
		monitorService := &monitor.MonitorService{}
		server := NewServer(config, monitorService, logger)
		
		// Create a test request
		req, err := http.NewRequest("GET", "/health", nil)
		assert.NoError(t, err)
		
		// Create a response recorder
		rr := httptest.NewRecorder()
		
		// Use the router to handle the request
		server.router.ServeHTTP(rr, req)
		
		// Check the response
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "healthy")
	})

	t.Run("Metrics endpoint", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		
		monitorService := &monitor.MonitorService{}
		server := NewServer(config, monitorService, logger)
		
		// Create a test request
		req, err := http.NewRequest("GET", "/api/v1/metrics/current", nil)
		assert.NoError(t, err)
		
		// Create a response recorder
		rr := httptest.NewRecorder()
		
		// Use the router to handle the request
		server.router.ServeHTTP(rr, req)
		
		// Check the response
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Alerts endpoint", func(t *testing.T) {
		config := &Config{
			Port:    0,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		
		monitorService := &monitor.MonitorService{}
		server := NewServer(config, monitorService, logger)
		
		// Create a test request
		req, err := http.NewRequest("GET", "/api/v1/alerts", nil)
		assert.NoError(t, err)
		
		// Create a response recorder
		rr := httptest.NewRecorder()
		
		// Use the router to handle the request
		server.router.ServeHTTP(rr, req)
		
		// Check the response
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Config validation", func(t *testing.T) {
		config := &Config{
			Port:    8080,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		
		assert.Equal(t, 8080, config.Port)
		assert.Equal(t, "localhost", config.Host)
		assert.Equal(t, 30*time.Second, config.Timeout)
	})

	t.Run("Router setup", func(t *testing.T) {
		config := &Config{
			Port:    8080,
			Host:    "localhost",
			Timeout: 30 * time.Second,
		}
		
		monitorService := &monitor.MonitorService{}
		server := NewServer(config, monitorService, logger)
		
		// Check that router is properly configured
		assert.NotNil(t, server.router)
		
		// Test that routes are registered
		req, err := http.NewRequest("GET", "/health", nil)
		assert.NoError(t, err)
		
		// This should not panic
		server.router.ServeHTTP(httptest.NewRecorder(), req)
	})
}