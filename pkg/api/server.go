// Package api provides REST API for etcd-monitor
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/etcd-monitor/taskmaster/pkg/monitor"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// MonitorServiceInterface defines the interface for MonitorService
type MonitorServiceInterface interface {
	GetClusterStatus() (*monitor.ClusterStatus, error)
	GetCurrentMetrics() (*monitor.MetricsSnapshot, error)
	GetAlertManager() *monitor.AlertManager
	GetHealthChecker() *monitor.HealthChecker
	GetMetricsCollector() *monitor.MetricsCollector
	IsRunning() bool
}

// Server is the API server
type Server struct {
	router         *mux.Router
	server         *http.Server
	monitorService MonitorServiceInterface
	logger         *zap.Logger
}

// Config holds the API server configuration
type Config struct {
	Port    int
	Host    string
	Timeout time.Duration
}

// NewServer creates a new API server
func NewServer(config *Config, monitorService MonitorServiceInterface, logger *zap.Logger) *Server {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	s := &Server{
		router:         mux.NewRouter(),
		monitorService: monitorService,
		logger:         logger,
	}

	s.setupRoutes()

	if config != nil {
		s.server = &http.Server{
			Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
			Handler:      s.router,
			ReadTimeout:  config.Timeout,
			WriteTimeout: config.Timeout,
		}
	}

	return s
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Health check
	s.router.HandleFunc("/health", s.handleHealth).Methods("GET")

	// Cluster endpoints
	s.router.HandleFunc("/api/v1/cluster/status", s.handleClusterStatus).Methods("GET")
	s.router.HandleFunc("/api/v1/cluster/members", s.handleClusterMembers).Methods("GET")
	s.router.HandleFunc("/api/v1/cluster/leader", s.handleClusterLeader).Methods("GET")

	// Metrics endpoints
	s.router.HandleFunc("/api/v1/metrics/current", s.handleCurrentMetrics).Methods("GET")
	s.router.HandleFunc("/api/v1/metrics/history", s.handleMetricsHistory).Methods("GET")
	s.router.HandleFunc("/api/v1/metrics/latency", s.handleLatencyMetrics).Methods("GET")

	// Alert endpoints
	s.router.HandleFunc("/api/v1/alerts", s.handleAlerts).Methods("GET")
	s.router.HandleFunc("/api/v1/alerts/history", s.handleAlertHistory).Methods("GET")

	// Performance endpoints
	s.router.HandleFunc("/api/v1/performance/benchmark", s.handleBenchmark).Methods("POST")

	// Add middleware
	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.corsMiddleware)
}

// Start starts the API server
func (s *Server) Start() error {
	s.logger.Info("Starting API server", zap.String("address", s.server.Addr))
	return s.server.ListenAndServe()
}

// Stop stops the API server
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping API server")
	return s.server.Shutdown(ctx)
}

// handleHealth returns the health status of the API server
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "etcd-monitor-api",
	}

	s.writeJSON(w, http.StatusOK, response)
}

// handleClusterStatus returns the current cluster status
func (s *Server) handleClusterStatus(w http.ResponseWriter, r *http.Request) {
	status, err := s.monitorService.GetClusterStatus()
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to get cluster status", err)
		return
	}

	s.writeJSON(w, http.StatusOK, status)
}

// handleClusterMembers returns cluster member information
func (s *Server) handleClusterMembers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get cluster status which includes member info
	status, err := s.monitorService.GetClusterStatus()
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to get cluster status", err)
		return
	}

	// Get detailed member information from health checker
	healthChecker := s.monitorService.GetHealthChecker()
	if healthChecker == nil {
		s.writeError(w, http.StatusInternalServerError, "Health checker not available", nil)
		return
	}

	// Get member list
	members, err := healthChecker.GetMemberList(ctx)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to get member list", err)
		return
	}

	response := map[string]interface{}{
		"member_count": status.MemberCount,
		"quorum_size":  status.QuorumSize,
		"members":      members,
		"timestamp":    time.Now().Format(time.RFC3339),
	}

	s.writeJSON(w, http.StatusOK, response)
}

// handleClusterLeader returns the current leader information
func (s *Server) handleClusterLeader(w http.ResponseWriter, r *http.Request) {
	status, err := s.monitorService.GetClusterStatus()
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to get cluster status", err)
		return
	}

	response := map[string]interface{}{
		"leader_id":   status.LeaderID,
		"has_leader":  status.HasLeader,
		"last_change": status.LastLeaderChange,
	}

	s.writeJSON(w, http.StatusOK, response)
}

// handleCurrentMetrics returns current metrics snapshot
func (s *Server) handleCurrentMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := s.monitorService.GetCurrentMetrics()
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to get metrics", err)
		return
	}

	s.writeJSON(w, http.StatusOK, metrics)
}

// handleMetricsHistory returns historical metrics
func (s *Server) handleMetricsHistory(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	// duration := r.URL.Query().Get("duration")
	// interval := r.URL.Query().Get("interval")

	// This would fetch historical metrics from the database
	response := map[string]interface{}{
		"message": "Historical metrics - not yet implemented",
	}
	s.writeJSON(w, http.StatusNotImplemented, response)
}

// handleLatencyMetrics returns latency metrics
func (s *Server) handleLatencyMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := s.monitorService.GetCurrentMetrics()
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to get metrics", err)
		return
	}

	response := map[string]interface{}{
		"read_latency": map[string]interface{}{
			"p50": metrics.ReadLatencyP50,
			"p95": metrics.ReadLatencyP95,
			"p99": metrics.ReadLatencyP99,
		},
		"write_latency": map[string]interface{}{
			"p50": metrics.WriteLatencyP50,
			"p95": metrics.WriteLatencyP95,
			"p99": metrics.WriteLatencyP99,
		},
		"timestamp": metrics.Timestamp,
	}

	s.writeJSON(w, http.StatusOK, response)
}

// handleAlerts returns current active alerts
func (s *Server) handleAlerts(w http.ResponseWriter, r *http.Request) {
	alertManager := s.monitorService.GetAlertManager()
	if alertManager == nil {
		s.writeError(w, http.StatusInternalServerError, "Alert manager not available", nil)
		return
	}

	// Get active alerts from the alert manager
	activeAlerts := alertManager.GetActiveAlerts()

	response := map[string]interface{}{
		"active_alerts": activeAlerts,
		"count":         len(activeAlerts),
		"timestamp":     time.Now().Format(time.RFC3339),
	}

	s.writeJSON(w, http.StatusOK, response)
}

// handleAlertHistory returns alert history
func (s *Server) handleAlertHistory(w http.ResponseWriter, r *http.Request) {
	alertManager := s.monitorService.GetAlertManager()
	if alertManager == nil {
		s.writeError(w, http.StatusInternalServerError, "Alert manager not available", nil)
		return
	}

	// Parse query parameters for filtering
	// limit := r.URL.Query().Get("limit") // Future enhancement
	// since := r.URL.Query().Get("since") // Future enhancement

	// Get alert history
	history := alertManager.GetAlertHistory()

	response := map[string]interface{}{
		"history":   history,
		"count":     len(history),
		"timestamp": time.Now().Format(time.RFC3339),
	}

	s.writeJSON(w, http.StatusOK, response)
}

// handleBenchmark runs a performance benchmark
func (s *Server) handleBenchmark(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Operations int `json:"operations"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if request.Operations == 0 {
		request.Operations = 100
	}

	response := map[string]interface{}{
		"message":    "Benchmark - not yet implemented",
		"operations": request.Operations,
	}
	s.writeJSON(w, http.StatusNotImplemented, response)
}

// writeJSON writes a JSON response
func (s *Server) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.Error("Failed to encode JSON response", zap.Error(err))
	}
}

// writeError writes an error response
func (s *Server) writeError(w http.ResponseWriter, status int, message string, err error) {
	s.logger.Error(message, zap.Error(err))

	response := map[string]interface{}{
		"error":     message,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	if err != nil {
		response["details"] = err.Error()
	}

	s.writeJSON(w, status, response)
}

// loggingMiddleware logs all HTTP requests
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer that captures the status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		s.logger.Info("HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", rw.statusCode),
			zap.Duration("duration", duration),
			zap.String("remote_addr", r.RemoteAddr))
	})
}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
