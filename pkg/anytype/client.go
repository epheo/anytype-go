package anytype

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	// HTTP client timeout
	httpTimeout = 10 * time.Second
	// Current API version
	apiVersion = "2025-03-17"
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
	}
}

// NewClient creates a new API client with options
func NewClient(apiURL, sessionToken, appKey string, opts ...ClientOption) *Client {
	client := &Client{
		apiURL:       apiURL,
		sessionToken: sessionToken,
		appKey:       appKey,
		httpClient:   &http.Client{Timeout: httpTimeout},
		debug:        false,
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	return client
}

// makeRequest is a helper function to make HTTP requests
func (c *Client) makeRequest(ctx context.Context, method, path string, body io.Reader) ([]byte, error) {
	url := c.apiURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set standard headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.appKey))
	req.Header.Set("Anytype-Version", apiVersion)

	if c.debug {
		c.printCurlRequest(method, url, req.Header, bodyToBytes(body))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, extractErrorFromResponse(path, resp.StatusCode, responseData)
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
func extractErrorFromResponse(path string, statusCode int, responseData []byte) error {
	var apiError struct {
		Message string `json:"message,omitempty"`
		Error   string `json:"error,omitempty"`
	}

	if err := json.Unmarshal(responseData, &apiError); err != nil {
		return fmt.Errorf("API error: %s returned status %d", path, statusCode)
	}

	if msg := apiError.Message; msg != "" {
		return fmt.Errorf("API error: %s returned status %d - %s", path, statusCode, msg)
	}
	if msg := apiError.Error; msg != "" {
		return fmt.Errorf("API error: %s returned status %d - %s", path, statusCode, msg)
	}

	return fmt.Errorf("API error: %s returned status %d", path, statusCode)
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

	fmt.Printf("\nDebug curl command:\n%s\n", sb.String())
}
