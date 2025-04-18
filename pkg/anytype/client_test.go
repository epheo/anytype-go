package anytype

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestNewClient tests the creation of a new client with various options
func TestNewClient(t *testing.T) {
	// For all tests, we need to provide an app key
	appKeyOption := WithAppKey("test-app-key")

	// Test with default options
	client, err := NewClient(appKeyOption)
	if err != nil {
		t.Fatalf("Failed to create client with default options: %v", err)
	}
	if client == nil {
		t.Fatal("Client should not be nil")
	}

	// Test with custom timeout
	customTimeout := 60 * time.Second
	client, err = NewClient(appKeyOption, WithTimeout(customTimeout))
	if err != nil {
		t.Fatalf("Failed to create client with custom timeout: %v", err)
	}
	if client == nil {
		t.Fatal("Client should not be nil")
	}

	// Test with debug enabled
	client, err = NewClient(appKeyOption, WithDebug(true))
	if err != nil {
		t.Fatalf("Failed to create client with debug enabled: %v", err)
	}
	if client == nil {
		t.Fatal("Client should not be nil")
	}

	// Test with custom URL
	customURL := "http://localhost:8000"
	client, err = NewClient(appKeyOption, WithURL(customURL))
	if err != nil {
		t.Fatalf("Failed to create client with custom URL: %v", err)
	}
	if client == nil {
		t.Fatal("Client should not be nil")
	}
}

// TestFromAuthConfig tests creating a client from auth config
func TestFromAuthConfig(t *testing.T) {
	// Create a mock auth config
	authConfig := &AuthConfig{
		SessionToken: "test-token",
		AppKey:       "test-app-key",
		ApiURL:       "http://localhost:8000",
	}

	// Test with auth config
	client, err := FromAuthConfig(authConfig)
	if err != nil {
		t.Fatalf("Failed to create client from auth config: %v", err)
	}
	if client == nil {
		t.Fatal("Client should not be nil")
	}

	// We can't directly test private fields, so we'll just verify the client was created
}

// createMockServerForClient creates a test server that returns the given response
func createMockServerForClient(t *testing.T, statusCode int, response string) (*httptest.Server, *Client) {
	// Create a test server that returns the response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write([]byte(response))
	}))

	// Create a client that points to the test server
	client, err := NewClient(
		WithURL(server.URL),
		WithAppKey("test-app-key"), // Add app key for tests
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	return server, client
}
