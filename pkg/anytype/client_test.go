package anytype

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestClient tests the API client and its methods
func TestClient(t *testing.T) {
	// Configure test server that mocks the Anytype API
	var receivedRequest *http.Request
	var receivedBody []byte

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Save the request for inspection in the test
		receivedRequest = r

		// Read and save the request body
		body, _ := io.ReadAll(r.Body)
		receivedBody = body

		// Set appropriate response based on path
		switch r.URL.Path {
		case "/v1/spaces":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
				"data": [
					{"id": "space1", "name": "Test Space 1", "type": "workspace"}, 
					{"id": "space2", "name": "Test Space 2", "type": "workspace"}
				]
			}`))

		case "/v1/spaces/space1":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
				"id": "space1", 
				"name": "Test Space 1", 
				"type": "workspace"
			}`))

		case "/v1/spaces/space1/types":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
				"data": [
					{
						"type": "object-type",
						"id": "type1",
						"unique_key": "blog-post",
						"name": "Blog Post",
						"icon": "📝",
						"recommended_layout": "blog"
					},
					{
						"type": "object-type",
						"id": "type2",
						"unique_key": "note",
						"name": "Note",
						"icon": "📌",
						"recommended_layout": "note"
					}
				],
				"pagination": {}
			}`))

		case "/v1/spaces/space1/objects":
			// Handle search or create based on method
			if r.Method == http.MethodGet {
				// This is a search request
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{
					"data": [
						{"id": "obj1", "name": "Test Object 1", "type": "blog-post"},
						{"id": "obj2", "name": "Test Object 2", "type": "blog-post"}
					],
					"pagination": {"total": 2}
				}`))
			} else if r.Method == http.MethodPost {
				// This is a create request
				// Parse the incoming object and echo it back with an ID
				var obj Object
				json.Unmarshal(body, &obj)
				obj.ID = "new-obj-id" // Assign a mock ID

				respBytes, _ := json.Marshal(obj)
				w.Header().Set("Content-Type", "application/json")
				w.Write(respBytes)
			}

		case "/v1/spaces/space1/objects/obj1":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
				"id": "obj1", 
				"name": "Test Object 1", 
				"type": "blog-post"
			}`))

		default:
			// Unknown path, return 404
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "Not found"}`))
		}
	}))
	defer mockServer.Close()

	// Create a client that points to our test server
	client := NewClient(
		mockServer.URL, // Use test server URL as API URL
		"test-session-token",
		"test-app-key",
		WithTimeout(2*time.Second),
	)

	ctx := context.Background()

	// Test GetSpaces
	t.Run("GetSpaces", func(t *testing.T) {
		spaces, err := client.GetSpaces(ctx)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if len(spaces.Data) != 2 {
			t.Errorf("Expected 2 spaces, got: %d", len(spaces.Data))
		}
		if spaces.Data[0].ID != "space1" {
			t.Errorf("Expected ID space1, got: %s", spaces.Data[0].ID)
		}

		// Verify request headers
		if receivedRequest.Header.Get("Authorization") != "Bearer test-app-key" {
			t.Errorf("Expected Authorization header with app key, got: %s",
				receivedRequest.Header.Get("Authorization"))
		}
	})

	// Test GetSpaceByID
	t.Run("GetSpaceByID", func(t *testing.T) {
		space, err := client.GetSpaceByID(ctx, "space1")

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if space.ID != "space1" {
			t.Errorf("Expected ID space1, got: %s", space.ID)
		}
		if space.Name != "Test Space 1" {
			t.Errorf("Expected name 'Test Space 1', got: %s", space.Name)
		}
	})

	// Test GetTypes
	t.Run("GetTypes", func(t *testing.T) {
		types, err := client.GetTypes(ctx, "space1")

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if len(types.Data) != 2 {
			t.Errorf("Expected 2 types, got: %d", len(types.Data))
		}

		// Verify the cache is populated
		typeName := client.GetTypeName(ctx, "space1", "blog-post")
		if typeName != "Blog Post" {
			t.Errorf("Expected type name 'Blog Post', got: %s", typeName)
		}
	})

	// Test GetTypeByName
	t.Run("GetTypeByName", func(t *testing.T) {
		// First call will populate the cache
		typeKey, err := client.GetTypeByName(ctx, "space1", "Blog Post")

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if typeKey != "blog-post" {
			t.Errorf("Expected key 'blog-post', got: %s", typeKey)
		}

		// Second call should use the cache
		typeKey2, err := client.GetTypeByName(ctx, "space1", "Note")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if typeKey2 != "note" {
			t.Errorf("Expected key 'note', got: %s", typeKey2)
		}
	})

	// Test Search
	t.Run("Search", func(t *testing.T) {
		params := &SearchParams{
			Query: "test",
			Types: []string{"blog-post"},
			Limit: 10,
		}

		results, err := client.Search(ctx, "space1", params)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if len(results.Data) != 2 {
			t.Errorf("Expected 2 objects, got: %d", len(results.Data))
		}

		// Verify search parameters were sent properly
		var sentParams SearchParams
		json.Unmarshal(receivedBody, &sentParams)

		if sentParams.Query != "test" {
			t.Errorf("Expected query 'test', got: %s", sentParams.Query)
		}
		if len(sentParams.Types) != 1 || sentParams.Types[0] != "blog-post" {
			t.Errorf("Expected types ['blog-post'], got: %v", sentParams.Types)
		}
	})

	// Test GetObject
	t.Run("GetObject", func(t *testing.T) {
		obj, err := client.GetObject(ctx, "space1", "obj1")

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if obj.ID != "obj1" {
			t.Errorf("Expected ID obj1, got: %s", obj.ID)
		}
		if obj.Name != "Test Object 1" {
			t.Errorf("Expected name 'Test Object 1', got: %s", obj.Name)
		}
	})

	// Test CreateObject
	t.Run("CreateObject", func(t *testing.T) {
		newObj := &Object{
			Name: "New Test Object",
			Type: "blog-post",
		}

		createdObj, err := client.CreateObject(ctx, "space1", newObj)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if createdObj.ID != "new-obj-id" {
			t.Errorf("Expected ID new-obj-id, got: %s", createdObj.ID)
		}
		if createdObj.Name != "New Test Object" {
			t.Errorf("Expected name 'New Test Object', got: %s", createdObj.Name)
		}
	})

	// Test error handling
	t.Run("ErrorHandling", func(t *testing.T) {
		// Test with invalid space ID
		_, err := client.GetSpaceByID(ctx, "")
		if err == nil {
			t.Error("Expected error for empty space ID, got nil")
		}

		// Test with non-existent endpoint
		_, err = client.GetSpaceByID(ctx, "nonexistent")
		if err == nil {
			t.Error("Expected error for non-existent space, got nil")
		}
	})
}
