package anytype

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/epheo/anytype-go/internal/log"
)

// Constants for client configuration
const (
	// HTTP client timeout
	httpTimeout = 10 * time.Second
	// Current API version
	apiVersion = "2025-03-17"
	// Default API URL
	defaultAPIURL = "http://localhost:31009"
)

// ClientOption defines a function type for client configuration
type ClientOption func(*Client)

// Client manages API communication
type Client struct {
	apiURL       string
	sessionToken string
	appKey       string
	httpClient   *http.Client
	debug        bool
	printCurl    bool                         // whether to print curl commands
	typeCache    map[string]map[string]string // spaceID -> typeKey -> typeName
	logger       log.Logger                   // for logging output
}

// WithTimeout sets a custom timeout for the HTTP client
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithDebug enables debug mode for the client
func WithDebug(debug bool) ClientOption {
	return func(c *Client) {
		c.debug = debug
		if c.logger != nil {
			if debug {
				c.logger.SetLevel(log.LevelDebug)
			} else {
				c.logger.SetLevel(log.LevelInfo)
			}
		}
	}
}

// WithLogger sets a logger for the client
func WithLogger(logger log.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithCurl enables printing curl equivalent of API requests
func WithCurl(printCurl bool) ClientOption {
	return func(c *Client) {
		c.printCurl = printCurl
	}
}

// WithURL sets the API URL
func WithURL(url string) ClientOption {
	return func(c *Client) {
		c.apiURL = url
	}
}

// WithToken sets the session token
func WithToken(token string) ClientOption {
	return func(c *Client) {
		c.sessionToken = token
	}
}

// WithAppKey sets the app key
func WithAppKey(appKey string) ClientOption {
	return func(c *Client) {
		c.appKey = appKey
	}
}

// NewClient creates a new API client with options
func NewClient(opts ...ClientOption) (*Client, error) {
	client := &Client{
		apiURL:     defaultAPIURL, // Default API URL
		httpClient: &http.Client{Timeout: httpTimeout},
		debug:      false,
		typeCache:  make(map[string]map[string]string),
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	// Validate required fields
	if client.apiURL == "" {
		return nil, fmt.Errorf("API URL is required")
	}

	if client.appKey == "" {
		return nil, fmt.Errorf("app key is required")
	}

	return client, nil
}

// FromAuthConfig creates a new client from an AuthConfig
func FromAuthConfig(config *AuthConfig, additionalOpts ...ClientOption) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("auth config cannot be nil")
	}

	// Start with basic options from the config
	opts := []ClientOption{
		WithURL(config.ApiURL),
		WithToken(config.SessionToken),
		WithAppKey(config.AppKey),
	}

	// Add any additional options
	opts = append(opts, additionalOpts...)

	return NewClient(opts...)
}

// FromEnvironment creates a new client from environment variables
func FromEnvironment(additionalOpts ...ClientOption) (*Client, error) {
	apiURL := getEnvOrDefault("ANYTYPE_API_URL", defaultAPIURL)
	appKey := os.Getenv("ANYTYPE_APP_KEY")
	sessionToken := os.Getenv("ANYTYPE_SESSION_TOKEN")

	if appKey == "" {
		return nil, fmt.Errorf("ANYTYPE_APP_KEY environment variable is not set")
	}

	// Start with basic options from environment variables
	opts := []ClientOption{
		WithURL(apiURL),
		WithToken(sessionToken),
		WithAppKey(appKey),
	}

	// Add any additional options
	opts = append(opts, additionalOpts...)

	return NewClient(opts...)
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// makeRequest is a helper function to make HTTP requests
func (c *Client) makeRequest(ctx context.Context, method, path string, body io.Reader) ([]byte, error) {
	url := c.apiURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, WrapError(path, 0, "failed to create HTTP request", err)
	}

	// Set standard headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.appKey))
	req.Header.Set("Anytype-Version", apiVersion)

	// Print curl command if debug mode or curl mode is enabled
	if c.debug || c.printCurl {
		c.printCurlRequest(method, url, req.Header, bodyToBytes(body))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Check if context was canceled
		if errors.Is(err, context.Canceled) {
			return nil, WrapError(path, 0, "request canceled", err)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, WrapError(path, 0, "request timed out", ErrOperationTimeout)
		}
		return nil, WrapError(path, 0, "failed to execute HTTP request", fmt.Errorf("%w: %s", ErrNetworkError, err.Error()))
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, WrapError(path, resp.StatusCode, "failed to read response body", err)
	}

	if c.debug && c.logger != nil {
		c.logger.Debug("Response: %s", string(responseData))
	}

	if resp.StatusCode != http.StatusOK {
		baseError := StatusCodeToError(resp.StatusCode)
		return nil, extractErrorFromResponse(path, resp.StatusCode, responseData, baseError)
	}

	return responseData, nil
}

// bodyToBytes converts an io.Reader to bytes for debug printing
func bodyToBytes(body io.Reader) []byte {
	if body == nil {
		return nil
	}
	if bodyBytes, ok := body.(*bytes.Buffer); ok {
		return bodyBytes.Bytes()
	}
	return nil
}

// extractErrorFromResponse tries to extract a meaningful error message from API response
func extractErrorFromResponse(path string, statusCode int, responseData []byte, baseErr error) error {
	var apiError struct {
		Message string `json:"message,omitempty"`
		Error   string `json:"error,omitempty"`
		Details string `json:"details,omitempty"`
		Code    string `json:"code,omitempty"`
	}

	if err := json.Unmarshal(responseData, &apiError); err != nil {
		return WrapError(path, statusCode, "unknown error", baseErr)
	}

	// Extract error message
	message := "unknown error"
	if apiError.Message != "" {
		message = apiError.Message
	} else if apiError.Error != "" {
		message = apiError.Error
	}

	// Extract additional details
	details := ""
	if apiError.Details != "" {
		details = apiError.Details
	} else if apiError.Code != "" {
		details = "code: " + apiError.Code
	}

	return WrapErrorWithDetails(path, statusCode, message, details, baseErr)
}

// printCurlRequest prints a curl command equivalent to the current request
func (c *Client) printCurlRequest(method, url string, headers http.Header, body []byte) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("curl -X %s '%s'", method, url))

	// Add headers
	for key, values := range headers {
		for _, value := range values {
			sb.WriteString(fmt.Sprintf(" \\\n  -H '%s: %s'", key, value))
		}
	}

	// Add body with proper JSON formatting if possible
	if len(body) > 0 {
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, body, "  ", "  "); err == nil {
			sb.WriteString(fmt.Sprintf(" \\\n  -d '%s'", prettyJSON.String()))
		} else {
			sb.WriteString(fmt.Sprintf(" \\\n  -d '%s'", string(body)))
		}
	}

	// If logger is available, use it; otherwise print to stdout
	if c.logger != nil {
		c.logger.Debug("CURL command:\n%s", sb.String())
	} else {
		fmt.Printf("CURL command:\n%s\n", sb.String())
	}
}

// GetTypeName returns the friendly name for a type key, using cache if available
func (c *Client) GetTypeName(ctx context.Context, spaceID, typeKey string) string {
	// Check cache first
	if cache, ok := c.typeCache[spaceID]; ok {
		if name, ok := cache[typeKey]; ok {
			return name
		}
	}

	// Initialize cache for this space if needed
	if _, ok := c.typeCache[spaceID]; !ok {
		c.typeCache[spaceID] = make(map[string]string)
	}

	// If cache is empty for this space, fetch all types at once
	// instead of doing it for each type key separately
	if len(c.typeCache[spaceID]) == 0 {
		// Fetch all types and update cache
		types, err := c.GetTypes(ctx, &GetTypesParams{SpaceID: spaceID})
		if err != nil {
			return typeKey // Return original key if error
		}

		// Update cache with all types
		for _, t := range types.Data {
			c.typeCache[spaceID][t.Key] = t.Name
		}
	}

	// Return cached value or original key if not found
	if name, ok := c.typeCache[spaceID][typeKey]; ok {
		return name
	}
	return typeKey
}
