package tests

import (
	"testing"

	"github.com/epheo/anytype-go"
)

// TestProperties tests property-related operations
func TestProperties(t *testing.T) {
	tc := setupTestClient(t)
	defer cleanupTestClient(tc)

	spaceID := findOrCreateTestSpace(t, tc)

	// List properties
	properties, err := tc.Client.Space(spaceID).Properties().List(tc.Ctx)
	if err != nil {
		t.Fatalf("Failed to list properties: %v", err)
	}

	t.Logf("Found %d properties in space", len(properties))

	// Verify we have some properties from the mock
	if len(properties) == 0 {
		t.Error("Expected to find at least one property")
	}
}

// TestCreateProperty tests creating a new property
func TestCreateProperty(t *testing.T) {
	tc := setupTestClient(t)
	defer cleanupTestClient(tc)

	spaceID := findOrCreateTestSpace(t, tc)

	// Create a new property
	createReq := anytype.CreatePropertyRequest{
		Key:    "test-property-key",
		Name:   "Test Property",
		Format: "text",
	}

	propResp, err := tc.Client.Space(spaceID).Properties().Create(tc.Ctx, createReq)
	if err != nil {
		t.Fatalf("Failed to create property: %v", err)
	}

	// Verify the property was created with the correct values
	if propResp.Property.Key != createReq.Key {
		t.Errorf("Property key mismatch: got %s, want %s", propResp.Property.Key, createReq.Key)
	}

	if propResp.Property.Name != createReq.Name {
		t.Errorf("Property name mismatch: got %s, want %s", propResp.Property.Name, createReq.Name)
	}

	if propResp.Property.Format != createReq.Format {
		t.Errorf("Property format mismatch: got %s, want %s", propResp.Property.Format, createReq.Format)
	}

	t.Logf("Successfully created property: %s (key: %s, format: %s)", propResp.Property.Name, propResp.Property.Key, propResp.Property.Format)
}
