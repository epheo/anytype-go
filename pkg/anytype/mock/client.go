package mock

import (
	"context"

	"github.com/epheo/anyblog/pkg/anytype"
)

// Client is a mock implementation of the anytype.Client for testing
type Client struct {
	// Mock data that will be returned by each method
	Spaces        *anytype.SpacesResponse
	Space         *anytype.Space
	Types         *anytype.TypeResponse
	TypeName      string
	SearchResults *anytype.SearchResponse
	Object        *anytype.Object

	// Keep track of calls to verify in tests
	GetSpacesCalled     bool
	GetSpaceByIDCalled  bool
	GetTypesCalled      bool
	GetTypeByNameCalled bool
	GetTypeNameCalled   bool
	SearchCalled        bool
	GetObjectCalled     bool
	CreateObjectCalled  bool

	// Store parameters passed to methods
	LastSpaceID  string
	LastTypeName string
	LastTypeKey  string
	LastQuery    *anytype.SearchParams
	LastObjectID string
	LastObject   *anytype.Object

	// Configured errors to return
	SpacesError       error
	SpaceByIDError    error
	TypesError        error
	TypeByNameError   error
	SearchError       error
	GetObjectError    error
	CreateObjectError error
}

// NewClient creates a new mock client with default values
func NewClient() *Client {
	return &Client{
		// Initialize with empty responses
		Spaces: &anytype.SpacesResponse{
			Data: []anytype.Space{},
		},
		Types: &anytype.TypeResponse{
			Data: []struct {
				Type              string `json:"type"`
				ID                string `json:"id"`
				UniqueKey         string `json:"unique_key"`
				Name              string `json:"name"`
				Icon              string `json:"icon"`
				RecommendedLayout string `json:"recommended_layout"`
			}{},
		},
		SearchResults: &anytype.SearchResponse{
			Data: []anytype.Object{},
		},
	}
}

// GetSpaces mocks the GetSpaces method
func (c *Client) GetSpaces(ctx context.Context) (*anytype.SpacesResponse, error) {
	c.GetSpacesCalled = true
	return c.Spaces, c.SpacesError
}

// GetSpaceByID mocks the GetSpaceByID method
func (c *Client) GetSpaceByID(ctx context.Context, spaceID string) (*anytype.Space, error) {
	c.GetSpaceByIDCalled = true
	c.LastSpaceID = spaceID
	return c.Space, c.SpaceByIDError
}

// GetTypes mocks the GetTypes method
func (c *Client) GetTypes(ctx context.Context, spaceID string) (*anytype.TypeResponse, error) {
	c.GetTypesCalled = true
	c.LastSpaceID = spaceID
	return c.Types, c.TypesError
}

// GetTypeByName mocks the GetTypeByName method
func (c *Client) GetTypeByName(ctx context.Context, spaceID, typeName string) (string, error) {
	c.GetTypeByNameCalled = true
	c.LastSpaceID = spaceID
	c.LastTypeName = typeName

	// If an error is configured, return it
	if c.TypeByNameError != nil {
		return "", c.TypeByNameError
	}

	// Otherwise return a mock type key
	if c.TypeName == "" {
		return "mock-type-key-for-" + typeName, nil
	}
	return c.TypeName, nil
}

// GetTypeName mocks the GetTypeName method
func (c *Client) GetTypeName(ctx context.Context, spaceID, typeKey string) string {
	c.GetTypeNameCalled = true
	c.LastSpaceID = spaceID
	c.LastTypeKey = typeKey

	// Return the configured type name or a default one
	if c.TypeName == "" {
		return "Mock Type for " + typeKey
	}
	return c.TypeName
}

// Search mocks the Search method
func (c *Client) Search(ctx context.Context, spaceID string, params *anytype.SearchParams) (*anytype.SearchResponse, error) {
	c.SearchCalled = true
	c.LastSpaceID = spaceID
	c.LastQuery = params
	return c.SearchResults, c.SearchError
}

// GetObject mocks the GetObject method
func (c *Client) GetObject(ctx context.Context, spaceID, objectID string) (*anytype.Object, error) {
	c.GetObjectCalled = true
	c.LastSpaceID = spaceID
	c.LastObjectID = objectID
	return c.Object, c.GetObjectError
}

// CreateObject mocks the CreateObject method
func (c *Client) CreateObject(ctx context.Context, spaceID string, object *anytype.Object) (*anytype.Object, error) {
	c.CreateObjectCalled = true
	c.LastSpaceID = spaceID
	c.LastObject = object

	// If an object is set to be returned, return it, otherwise return the input object
	if c.Object != nil {
		return c.Object, c.CreateObjectError
	}
	return object, c.CreateObjectError
}
