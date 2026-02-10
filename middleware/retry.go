package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

// RetryConfig configures the retry middleware
type RetryConfig struct {
	// MaxRetries is the maximum number of retries
	MaxRetries int
	// RetryDelay is the base delay between retries
	RetryDelay time.Duration
	// MaxRetryDelay is the maximum delay between retries
	MaxRetryDelay time.Duration
	// RetryableStatusCodes defines HTTP status codes that should trigger a retry
	RetryableStatusCodes []int
	// ShouldRetry is a custom function to determine if a request should be retried
	ShouldRetry func(*http.Response, error) bool
}

// DefaultRetryConfig provides sensible defaults for retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:    5,
		RetryDelay:    200 * time.Millisecond,
		MaxRetryDelay: 30 * time.Second,
		ShouldRetry:   nil, // Will use default logic with RetryableStatusCodes
		RetryableStatusCodes: []int{
			http.StatusRequestTimeout,
			http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
		},
	}
}

// RetryMiddleware implements a retry strategy for HTTP requests
type RetryMiddleware struct {
	Next   HTTPDoer
	Config RetryConfig
}

// NewRetryMiddleware creates a new retry middleware with the given configuration
func NewRetryMiddleware(next HTTPDoer, config RetryConfig) *RetryMiddleware {
	return &RetryMiddleware{
		Next:   next,
		Config: config,
	}
}

// Do executes an HTTP request with retries according to the retry policy
func (m *RetryMiddleware) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	// Try initial request first to avoid retry overhead for successful requests.
	// Most requests succeed on first attempt, so we optimize for the happy path.
	resp, err = m.Next.Do(req)

	// Check if we need to retry
	if !m.shouldRetry(resp, err) {
		return resp, err
	}

	// Prepare body for retries if needed
	var getBody func() (io.ReadCloser, error)
	if req.Body != nil {
		// Use GetBody if available (avoids reading entire body into memory for large requests)
		if req.GetBody != nil {
			getBody = req.GetBody
		} else {
			// Fall back to manual body cloning when GetBody is unavailable
			// (e.g., custom io.Reader implementations)
			bodyBytes, readErr := io.ReadAll(req.Body)
			if readErr != nil {
				return nil, readErr
			}
			req.Body.Close()
			// Set up getBody function for retries
			getBody = func() (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewReader(bodyBytes)), nil
			}
			// Reset body for first retry
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
	}

	// Retry loop
	retries := 0
	for retries < m.Config.MaxRetries && m.shouldRetry(resp, err) {
		// Delay before retry with exponential backoff
		delay := exponentialBackoff(m.Config.RetryDelay, retries, m.Config.MaxRetryDelay)
		select {
		case <-req.Context().Done():
			// Context was canceled
			return resp, req.Context().Err()
		case <-time.After(delay):
			// Continue with retry
		}

		// Clone the request for retry
		retryReq := req.Clone(req.Context())
		if getBody != nil {
			retryReq.Body, err = getBody()
			if err != nil {
				return resp, err
			}
		}

		// Execute the retry
		resp, err = m.Next.Do(retryReq)
		retries++
	}

	return resp, err
}

// shouldRetry determines if a request should be retried based on the response and error
func (m *RetryMiddleware) shouldRetry(resp *http.Response, err error) bool {
	if m.Config.ShouldRetry != nil {
		return m.Config.ShouldRetry(resp, err)
	}
	return m.defaultShouldRetry(resp, err)
}

// defaultShouldRetry provides default retry logic using configured status codes
func (m *RetryMiddleware) defaultShouldRetry(resp *http.Response, err error) bool {
	// Retry on connection errors
	if err != nil {
		return true
	}

	// Retry on configured status codes
	if resp != nil {
		for _, code := range m.Config.RetryableStatusCodes {
			if resp.StatusCode == code {
				return true
			}
		}
	}

	return false
}

// WithRetry returns a middleware function that applies retry logic with default configuration
func WithRetry() func(HTTPDoer) HTTPDoer {
	return WithCustomRetry(DefaultRetryConfig())
}

// WithCustomRetry returns a middleware function that applies retry logic with custom configuration
func WithCustomRetry(config RetryConfig) func(HTTPDoer) HTTPDoer {
	return func(next HTTPDoer) HTTPDoer {
		return NewRetryMiddleware(next, config)
	}
}

// exponentialBackoff calculates exponential backoff delay
func exponentialBackoff(baseDelay time.Duration, retry int, maxDelay time.Duration) time.Duration {
	// Exponential backoff: delay doubles with each retry (2^retry).
	// Bit shift is used for efficient power-of-2 multiplication.
	delay := baseDelay * (1 << uint(retry))
	if delay > maxDelay {
		delay = maxDelay
	}
	return delay
}
