package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHandleClusterStatus(t *testing.T) {
	t.Skip("Skipping - requires mock monitor service to avoid panics")
}

func TestHandleClusterMembers(t *testing.T) {
	t.Skip("Skipping - requires mock monitor service")
}

func TestHandleCurrentMetrics(t *testing.T) {
	t.Skip("Skipping - requires mock monitor service")
}

func TestHandleMetricsHistory(t *testing.T) {
	t.Skip("Skipping - requires mock monitor service")
}

func TestHandleAlerts(t *testing.T) {
	t.Skip("Skipping - requires mock monitor service")
}

func TestHandleBenchmark(t *testing.T) {
	t.Skip("Skipping - requires mock monitor service")
}

func TestServerLifecycle(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Server creation and basic setup", func(t *testing.T) {
		config := &Config{
			Host:    "localhost",
			Port:    0, // Use port 0 for testing (random available port)
			Timeout: 30 * time.Second,
		}

		server := NewServer(config, nil, logger)

		assert.NotNil(t, server)
		assert.NotNil(t, server.router)
		assert.NotNil(t, server.server)
		assert.NotNil(t, server.logger)
	})

	t.Run("Server with nil logger", func(t *testing.T) {
		config := &Config{
			Host:    "localhost",
			Port:    0,
			Timeout: 30 * time.Second,
		}

		server := NewServer(config, nil, nil)

		assert.NotNil(t, server)
		assert.NotNil(t, server.logger) // Should have default logger
	})

	t.Run("Server stop before start", func(t *testing.T) {
		config := &Config{
			Host:    "localhost",
			Port:    0,
			Timeout: 30 * time.Second,
		}

		server := NewServer(config, nil, logger)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := server.Stop(ctx)
		// Should handle gracefully even if not started
		assert.NoError(t, err)
	})
}

func TestWriteJSON(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := &Config{
		Host:    "localhost",
		Port:    0,
		Timeout: 30 * time.Second,
	}
	server := NewServer(config, nil, logger)

	t.Run("Write valid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := map[string]string{"message": "test"}

		server.writeJSON(w, http.StatusOK, data)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "test", response["message"])
	})

	t.Run("Write empty JSON", func(t *testing.T) {
		w := httptest.NewRecorder()

		server.writeJSON(w, http.StatusOK, map[string]interface{}{})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	})
}

func TestWriteError(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := &Config{
		Host:    "localhost",
		Port:    0,
		Timeout: 30 * time.Second,
	}
	server := NewServer(config, nil, logger)

	t.Run("Write error response", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := assert.AnError

		server.writeError(w, http.StatusBadRequest, "bad request", err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]interface{}
		jsonErr := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, jsonErr)
		assert.Contains(t, response, "error")
		assert.Equal(t, "bad request", response["error"])
	})

	t.Run("Write 500 error", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := assert.AnError

		server.writeError(w, http.StatusInternalServerError, "internal server error", err)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		jsonErr := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, jsonErr)
		assert.Equal(t, "internal server error", response["error"])
	})
}

func TestMiddlewares(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	config := &Config{
		Host:    "localhost",
		Port:    0,
		Timeout: 30 * time.Second,
	}
	server := NewServer(config, nil, logger)

	t.Run("Logging middleware logs requests", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := server.loggingMiddleware(handler)

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("CORS middleware sets headers", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := server.corsMiddleware(handler)

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
	})

	t.Run("CORS middleware handles OPTIONS", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Handler should not be called for OPTIONS")
		})

		middleware := server.corsMiddleware(handler)

		req := httptest.NewRequest("OPTIONS", "/test", nil)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	})
}
