package anytype

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestQueryBuilder tests the query builder functionality
func TestQueryBuilder(t *testing.T) {
	// We'll create the response directly in the handler

	// Set up test server to mock the type resolution and search requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Mock the type resolution endpoint
		if strings.Contains(r.URL.Path, "/types") && r.Method == http.MethodGet {
			w.Write([]byte(`{
				"data": [
					{
						"key": "ot-note",
						"name": "Note"
					}
				]
			}`))
			return
		}

		// Mock the search endpoint
		if strings.Contains(r.URL.Path, "/search") && r.Method == http.MethodPost {
			// Create proper response object that matches the actual API implementation
			response := map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"id":   "obj123",
						"name": "Test Note",
						"type": map[string]interface{}{
							"key":  "ot-note",
							"name": "Note",
						},
						"relations": map[string]interface{}{
							"items": map[string]interface{}{
								"tags": []map[string]interface{}{
									{"name": "important"},
									{"name": "work"},
								},
							},
						},
					},
				},
				"pagination": map[string]interface{}{
					"total":    1,
					"offset":   0,
					"limit":    25,
					"has_more": false,
				},
			}

			responseData, _ := json.Marshal(response)
			w.Write(responseData)
			return
		}

		http.Error(w, "Not found", http.StatusNotFound)
	}))
	defer server.Close()

	// Create client with test server URL and required app key
	client, err := NewClient(
		WithURL(server.URL),
		WithAppKey("test-app-key"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create a query builder
	ctx := context.Background()
	qb := client.NewQueryBuilder("space123")

	// Build a complex query using the actual API methods
	qb.WithQuery("test")
	qb.WithTypeKeys("ot-note")
	qb.WithTag("important")
	qb.WithLimit(25)
	qb.WithSortField("name", true) // Fix: Using WithSortField instead of OrderBy

	// Execute the query
	results, err := qb.Execute(ctx)

	// Check for errors
	if err != nil {
		t.Fatalf("Query builder execution failed: %v", err)
	}

	// Validate the response
	if len(results.Data) != 1 {
		t.Fatalf("Expected 1 object, got %d", len(results.Data))
	}
	if results.Data[0].ID != "obj123" || results.Data[0].Name != "Test Note" {
		t.Fatalf("Object data incorrect: %+v", results.Data[0])
	}
}

// TestQueryBuilderWithCallback tests the ExecuteWithCallback method
func TestQueryBuilderWithCallback(t *testing.T) {
	// Mock response for search
	mockResponse := `{
		"data": [
			{
				"id": "obj123",
				"name": "Test Note",
				"tags": ["important"]
			},
			{
				"id": "obj456",
				"name": "Another Note",
				"tags": ["work"]
			}
		],
		"pagination": {
			"total": 2,
			"offset": 0,
			"limit": 25,
			"has_more": false
		}
	}`

	// Set up test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Create client with test server URL and required app key
	client, err := NewClient(
		WithURL(server.URL),
		WithAppKey("test-app-key"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create a query builder
	ctx := context.Background()
	qb := client.NewQueryBuilder("space123")

	// Build a query with the actual API
	qb.WithTypeKeys("ot-note")
	qb.WithLimit(25)

	// Process results with a callback
	var processedObjects []Object
	err = qb.ExecuteWithCallback(ctx, func(obj Object) error {
		processedObjects = append(processedObjects, obj)
		return nil
	})

	// Check for errors
	if err != nil {
		t.Fatalf("Query builder callback execution failed: %v", err)
	}

	// Validate that callback processed all objects
	if len(processedObjects) != 2 {
		t.Fatalf("Expected callback to process 2 objects, got %d", len(processedObjects))
	}
	if processedObjects[0].ID != "obj123" || processedObjects[0].Name != "Test Note" {
		t.Fatalf("First processed object incorrect: %+v", processedObjects[0])
	}
	if processedObjects[1].ID != "obj456" || processedObjects[1].Name != "Another Note" {
		t.Fatalf("Second processed object incorrect: %+v", processedObjects[1])
	}
}

// TestQueryBuilderWithTimeout tests the WithTimeout method
func TestQueryBuilderWithTimeout(t *testing.T) {
	// Create a client with required app key
	client, err := NewClient(WithAppKey("test-app-key"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create a query builder with timeout
	qb := client.NewQueryBuilder("space123")
	timeout := 5 * time.Second
	qb.WithTimeout(timeout)

	// Verify the timeout was set properly
	// Note: We can't directly access the timeout field as it's private,
	// but we can indirectly test it by checking if an error occurs
	// with a canceled context that would normally trigger the timeout
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately to simulate timeout

	// This should fail due to context cancellation, not due to the query builder's timeout
	_, err = qb.Execute(ctx)
	if err == nil {
		t.Fatal("Expected error when context is canceled")
	}
	// Accept any error that occurs with a canceled context - we just want to make sure it fails
	// and doesn't hang or timeout in an unexpected way
}

// TestQueryBuilderErrors tests error handling in the query builder
func TestQueryBuilderErrors(t *testing.T) {
	// Create a client with required app key
	client, err := NewClient(WithAppKey("test-app-key"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test with invalid limit
	qb := client.NewQueryBuilder("space123")
	qb.WithLimit(-1) // Invalid limit

	ctx := context.Background()
	_, err = qb.Execute(ctx)
	if err == nil {
		t.Fatal("Expected error for negative limit, got nil")
	}

	// Test with invalid offset
	qb = client.NewQueryBuilder("space123")
	qb.WithOffset(-5) // Invalid offset

	_, err = qb.Execute(ctx)
	if err == nil {
		t.Fatal("Expected error for negative offset, got nil")
	}

	// Test with empty space ID
	qb = client.NewQueryBuilder("") // Empty space ID
	_, err = qb.Execute(ctx)
	if err == nil {
		t.Fatal("Expected error for empty space ID, got nil")
	}
}
