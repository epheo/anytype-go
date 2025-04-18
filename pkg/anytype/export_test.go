package anytype

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestExportObject tests the ExportObject method with a mocked API response
func TestExportObject(t *testing.T) {
	// Create a temporary directory for exports
	tempDir, err := ioutil.TempDir("", "anytype-export-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Mock response for object content
	mockResponse := `# Test Note

This is the content of a test note.

- Item 1
- Item 2
- Item 3

## Heading

More content here.
`

	// Set up test server
	server := http.NewServeMux()
	server.HandleFunc("/v1/spaces/space123/objects/obj123", func(w http.ResponseWriter, r *http.Request) {
		// This endpoint would be called first to get object metadata
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "obj123",
			"name": "Test Note",
			"type": {
				"key": "ot-note",
				"name": "Note"
			}
		}`))
	})
	server.HandleFunc("/v1/spaces/space123/objects/obj123/markdown", func(w http.ResponseWriter, r *http.Request) {
		// This endpoint handles the export
		w.Header().Set("Content-Type", "text/markdown")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	})

	testServer := httptest.NewServer(server)
	defer testServer.Close()

	// Create client with test server URL and required app key
	client, err := NewClient(
		WithURL(testServer.URL),
		WithAppKey("test-app-key"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Call the export method
	ctx := context.Background()
	filePath, err := client.ExportObject(ctx, "space123", "obj123", tempDir, "md")

	// Check for errors
	if err != nil {
		t.Fatalf("ExportObject failed: %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("Export file was not created at %s", filePath)
	}

	// Read the exported file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read exported file: %v", err)
	}

	// Verify the content
	if string(content) != mockResponse {
		t.Fatalf("Exported content doesn't match expected content. Got:\n%s\n\nExpected:\n%s", content, mockResponse)
	}
}

// TestExportObjects tests the ExportObjects method with multiple objects
func TestExportObjects(t *testing.T) {
	// Create a temporary directory for exports
	tempDir, err := ioutil.TempDir("", "anytype-export-test-multi")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set up test server with handlers for multiple objects
	server := http.NewServeMux()

	// First object
	server.HandleFunc("/v1/spaces/space123/objects/obj123", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "obj123",
			"name": "First Note",
			"type": {
				"key": "ot-note",
				"name": "Note"
			}
		}`))
	})
	server.HandleFunc("/v1/spaces/space123/objects/obj123/markdown", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/markdown")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("# First Note content"))
	})

	// Second object
	server.HandleFunc("/v1/spaces/space123/objects/obj456", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "obj456",
			"name": "Second Note",
			"type": {
				"key": "ot-note",
				"name": "Note"
			}
		}`))
	})
	server.HandleFunc("/v1/spaces/space123/objects/obj456/markdown", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/markdown")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("# Second Note content"))
	})

	testServer := httptest.NewServer(server)
	defer testServer.Close()

	// Create client with test server URL and required app key
	client, err := NewClient(
		WithURL(testServer.URL),
		WithAppKey("test-app-key"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Define objects to export
	objects := []Object{
		{ID: "obj123", Name: "First Note", Type: &TypeInfo{Key: "ot-note", Name: "Note"}},
		{ID: "obj456", Name: "Second Note", Type: &TypeInfo{Key: "ot-note", Name: "Note"}},
	}

	// Call the export method
	ctx := context.Background()
	filePaths, err := client.ExportObjects(ctx, "space123", objects, tempDir, "md")

	// Check for errors
	if err != nil {
		t.Fatalf("ExportObjects failed: %v", err)
	}

	// Verify the files were created
	if len(filePaths) != 2 {
		t.Fatalf("Expected 2 exported files, got %d", len(filePaths))
	}

	// Check each file
	for i, filePath := range filePaths {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Fatalf("Export file %d was not created at %s", i, filePath)
		}
	}

	// Read and verify first file
	content1, err := ioutil.ReadFile(filePaths[0])
	if err != nil {
		t.Fatalf("Failed to read first exported file: %v", err)
	}
	if string(content1) != "# First Note content" {
		t.Fatalf("First exported content incorrect. Got:\n%s", string(content1))
	}

	// Read and verify second file
	content2, err := ioutil.ReadFile(filePaths[1])
	if err != nil {
		t.Fatalf("Failed to read second exported file: %v", err)
	}
	if string(content2) != "# Second Note content" {
		t.Fatalf("Second exported content incorrect. Got:\n%s", string(content2))
	}
}

// TestExportValidation tests error cases for export functions
func TestExportValidation(t *testing.T) {
	client, err := NewClient(WithAppKey("test-app-key")) // Add app key for tests
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	ctx := context.Background()

	// Test with empty space ID
	_, err = client.ExportObject(ctx, "", "obj123", "/tmp", "md")
	if err == nil {
		t.Fatal("ExportObject should fail with empty space ID")
	}

	// Test with empty object ID
	_, err = client.ExportObject(ctx, "space123", "", "/tmp", "md")
	if err == nil {
		t.Fatal("ExportObject should fail with empty object ID")
	}

	// Test with empty export path
	_, err = client.ExportObject(ctx, "space123", "obj123", "", "md")
	if err == nil {
		t.Fatal("ExportObject should fail with empty export path")
	}

	// Test ExportObjects with empty object list
	_, err = client.ExportObjects(ctx, "space123", []Object{}, "/tmp", "md")
	if err == nil {
		t.Fatal("ExportObjects should fail with empty object list")
	}
}
