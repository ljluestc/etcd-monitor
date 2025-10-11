package plugin

import (
	"context"
	"net/http"
	"time"
)

// Phase represents the request/response lifecycle phase
type Phase int

const (
	PhaseRequest Phase = iota
	PhaseResponse
)

// Plugin interface that all plugins must implement
type Plugin interface {
	// Name returns the plugin name
	Name() string

	// Priority returns execution priority (lower number = earlier execution)
	Priority() int

	// OnRequest is called before forwarding to upstream
	OnRequest(ctx *Context) error

	// OnResponse is called after receiving upstream response
	OnResponse(ctx *Context) error
}

// Context holds request/response context for plugins
type Context struct {
	// HTTP request and response
	Request  *http.Request
	Response *http.Response
	Writer   http.ResponseWriter

	// Upstream information
	UpstreamAddr string
	UpstreamName string

	// Timing
	StartTime time.Time
	EndTime   time.Time

	// Variables for inter-plugin communication
	Variables map[string]interface{}

	// Context for cancellation
	Ctx context.Context

	// Error tracking
	Error error

	// Skip further plugin execution
	Skip bool

	// Response was written by plugin
	ResponseWritten bool
}

// NewContext creates a new plugin context
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Request:   r,
		Writer:    w,
		StartTime: time.Now(),
		Variables: make(map[string]interface{}),
		Ctx:       r.Context(),
	}
}

// Set stores a variable in the context
func (c *Context) Set(key string, value interface{}) {
	c.Variables[key] = value
}

// Get retrieves a variable from the context
func (c *Context) Get(key string) (interface{}, bool) {
	val, ok := c.Variables[key]
	return val, ok
}

// GetString retrieves a string variable
func (c *Context) GetString(key string) string {
	if val, ok := c.Variables[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// WriteError writes an error response and marks response as written
func (c *Context) WriteError(statusCode int, message string) {
	c.Writer.WriteHeader(statusCode)
	c.Writer.Write([]byte(message))
	c.ResponseWritten = true
	c.Skip = true
}

// WriteJSON writes a JSON response
func (c *Context) WriteJSON(statusCode int, data []byte) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(statusCode)
	c.Writer.Write(data)
	c.ResponseWritten = true
	c.Skip = true
}
