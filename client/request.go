// Package client implements the client interfaces defined in the anytype package
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/epheo/anytype-go/middleware"
)

// newRequest creates a new HTTP request with the appropriate headers
func (c *ClientImpl) newRequest(ctx context.Context, method, urlPath string, body interface{}) (*http.Request, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}
	// Ensure the /v1 prefix is included in the path
	u.Path = path.Join(u.Path, "/v1", urlPath)

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	// Set standard headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Anytype-Version", "2025-11-08")

	// Set authentication with app key as the Bearer token
	if c.appKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.appKey))
	}

	return req, nil
}

// doRequest executes the HTTP request and unmarshals the response into the result
func (c *ClientImpl) doRequest(req *http.Request, result interface{}) error {
	// Create a middleware chain
	chain := middleware.NewChain(c.httpClient)
	// Add middleware (e.g., retry, validation)
	chain.Use(middleware.WithRetry())
	// chain.Use(middleware.WithValidation())

	// Build the client with middleware chain
	client := chain.Build()

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle non-2xx responses
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response if a result is expected
	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return err
		}
	}

	return nil
}
