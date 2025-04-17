package anytype

import (
	"errors"
	"fmt"
	"net/http"
)

// Common error types
var (
	// Base errors
	ErrInvalidSpaceID   = errors.New("invalid space ID")
	ErrInvalidObjectID  = errors.New("invalid object ID")
	ErrInvalidTypeID    = errors.New("invalid type ID")
	ErrInvalidTemplate  = errors.New("invalid template")
	ErrInvalidParameter = errors.New("invalid parameter")

	// API-related errors
	ErrEmptyResponse      = errors.New("empty response from API")
	ErrInvalidResponse    = errors.New("invalid response format")
	ErrMissingRequired    = errors.New("missing required parameter")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrNotFound           = errors.New("resource not found")
	ErrServerError        = errors.New("server error")
	ErrNetworkError       = errors.New("network error")
	ErrOperationTimeout   = errors.New("operation timed out")
	ErrSpaceNotFound      = errors.New("space not found")
	ErrObjectNotFound     = errors.New("object not found")
	ErrTypeNotFound       = errors.New("type not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Error wraps API errors with additional context
type Error struct {
	StatusCode int    // HTTP status code if applicable
	Message    string // Human-readable error message
	Path       string // API path that caused the error
	Err        error  // Underlying error for unwrapping
	Details    string // Additional error details
}

// Implements the error interface
func (e *Error) Error() string {
	if e.StatusCode != 0 {
		if e.Details != "" {
			return fmt.Sprintf("API error: %s (status %d) - %s - %s", e.Path, e.StatusCode, e.Message, e.Details)
		}
		return fmt.Sprintf("API error: %s (status %d) - %s", e.Path, e.StatusCode, e.Message)
	}
	if e.Details != "" {
		return fmt.Sprintf("API error: %s - %s - %s", e.Path, e.Message, e.Details)
	}
	return fmt.Sprintf("API error: %s - %s", e.Path, e.Message)
}

// Unwrap returns the underlying error
func (e *Error) Unwrap() error {
	return e.Err
}

// WrapError creates a new Error with context
func WrapError(path string, statusCode int, message string, err error) *Error {
	return &Error{
		StatusCode: statusCode,
		Message:    message,
		Path:       path,
		Err:        err,
	}
}

// WrapErrorWithDetails creates a new Error with context and additional details
func WrapErrorWithDetails(path string, statusCode int, message string, details string, err error) *Error {
	return &Error{
		StatusCode: statusCode,
		Message:    message,
		Path:       path,
		Err:        err,
		Details:    details,
	}
}

// StatusCodeToError maps HTTP status codes to appropriate error types
func StatusCodeToError(statusCode int) error {
	switch statusCode {
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrUnauthorized
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		return ErrServerError
	case http.StatusRequestTimeout, http.StatusGatewayTimeout:
		return ErrOperationTimeout
	default:
		if statusCode >= 400 && statusCode < 500 {
			return fmt.Errorf("client error: status code %d", statusCode)
		}
		if statusCode >= 500 {
			return fmt.Errorf("server error: status code %d", statusCode)
		}
	}
	return nil
}

// IsNotFoundError checks if an error is a not found error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound) ||
		errors.Is(err, ErrSpaceNotFound) ||
		errors.Is(err, ErrObjectNotFound) ||
		errors.Is(err, ErrTypeNotFound)
}

// IsAuthenticationError checks if an error is authentication-related
func IsAuthenticationError(err error) bool {
	return errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrInvalidCredentials)
}

// IsServerError checks if an error is server-related
func IsServerError(err error) bool {
	return errors.Is(err, ErrServerError) || errors.Is(err, ErrOperationTimeout)
}

// IsClientError checks if an error is client-related
func IsClientError(err error) bool {
	var apiErr *Error
	if errors.As(err, &apiErr) && apiErr.StatusCode >= 400 && apiErr.StatusCode < 500 {
		return true
	}
	return errors.Is(err, ErrInvalidParameter) ||
		errors.Is(err, ErrInvalidSpaceID) ||
		errors.Is(err, ErrInvalidObjectID) ||
		errors.Is(err, ErrInvalidTypeID) ||
		errors.Is(err, ErrInvalidTemplate) ||
		errors.Is(err, ErrMissingRequired)
}
