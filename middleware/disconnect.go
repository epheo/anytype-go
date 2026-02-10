package middleware

import (
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// DisconnectConfig configures the disconnect handling middleware
type DisconnectConfig struct {
	// ReconnectDelay is the delay before attempting to reconnect
	ReconnectDelay time.Duration
	// MaxReconnectAttempts is the maximum number of reconnect attempts
	MaxReconnectAttempts int
	// OnDisconnect is called when a disconnection is detected
	OnDisconnect func()
	// OnReconnect is called when a connection is re-established
	OnReconnect func()
}

// DefaultDisconnectConfig provides sensible defaults for disconnect configuration
func DefaultDisconnectConfig() DisconnectConfig {
	return DisconnectConfig{
		ReconnectDelay:       5 * time.Second,
		MaxReconnectAttempts: 5,
	}
}

// DisconnectMiddleware handles temporary network disconnections
type DisconnectMiddleware struct {
	Next   HTTPDoer
	Config DisconnectConfig

	// Connection state is protected by mutex because it's accessed from
	// concurrent HTTP requests and callback functions.
	connected bool
	mu        sync.RWMutex
}

// NewDisconnectMiddleware creates a new disconnect handling middleware
func NewDisconnectMiddleware(next HTTPDoer, config DisconnectConfig) *DisconnectMiddleware {
	return &DisconnectMiddleware{
		Next:      next,
		Config:    config,
		connected: true,
	}
}

// Do executes an HTTP request with disconnect handling
func (m *DisconnectMiddleware) Do(req *http.Request) (*http.Response, error) {
	resp, err := m.Next.Do(req)

	// Check if this is a network-related error
	if m.isNetworkError(err) {
		m.handleDisconnect()
		return resp, err
	}

	// If we were previously disconnected, but now succeeded
	// mark as reconnected
	m.mu.RLock()
	wasDisconnected := !m.connected
	m.mu.RUnlock()

	if wasDisconnected && err == nil {
		m.handleReconnect()
	}

	return resp, err
}

// isNetworkError checks if an error is related to network connectivity
func (m *DisconnectMiddleware) isNetworkError(err error) bool {
	if err == nil {
		return false
	}

	// Handle both direct net.Error and wrapped url.Error types because
	// Go's http.Client wraps network errors in url.Error
	if _, ok := err.(net.Error); ok {
		return true
	}

	// Check if it's a wrapped network error
	if err, ok := err.(*url.Error); ok {
		if _, ok := err.Err.(net.Error); ok {
			return true
		}
	}

	// Fallback to string matching when type assertions fail - necessary because
	// some network errors don't implement net.Error interface
	errStr := err.Error()
	return strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "no such host") ||
		strings.Contains(errStr, "host unreachable") ||
		strings.Contains(errStr, "i/o timeout") ||
		strings.Contains(errStr, "network is unreachable") ||
		strings.Contains(errStr, "connection reset by peer")
}

// handleDisconnect handles a disconnection event
func (m *DisconnectMiddleware) handleDisconnect() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.connected {
		m.connected = false
		if m.Config.OnDisconnect != nil {
			m.Config.OnDisconnect()
		}
	}
}

// handleReconnect handles a reconnection event
func (m *DisconnectMiddleware) handleReconnect() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.connected {
		m.connected = true
		if m.Config.OnReconnect != nil {
			m.Config.OnReconnect()
		}
	}
}

// IsConnected returns the current connection status
func (m *DisconnectMiddleware) IsConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.connected
}
