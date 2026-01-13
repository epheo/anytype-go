package mocks

import (
	"context"

	"github.com/epheo/anytype-go"
)

// MockSpacesService implements the anytype.SpaceClient interface for testing
type MockSpacesService struct {
	ListFunc   func(ctx context.Context) (*anytype.SpaceListResponse, error)
	CreateFunc func(ctx context.Context, req anytype.CreateSpaceRequest) (*anytype.CreateSpaceResponse, error)
}

// NewMockSpacesService creates a new instance of MockSpacesService with default implementations
func NewMockSpacesService() *MockSpacesService {
	return &MockSpacesService{
		ListFunc: func(ctx context.Context) (*anytype.SpaceListResponse, error) {
			return &anytype.SpaceListResponse{
				Data: []anytype.Space{
					{
						ID:          "mock-space-id",
						Name:        "Mock Space",
						Description: "A mock space for testing",
					},
				},
			}, nil
		},
		CreateFunc: func(ctx context.Context, req anytype.CreateSpaceRequest) (*anytype.CreateSpaceResponse, error) {
			return &anytype.CreateSpaceResponse{
				Space: anytype.Space{
					ID:          "new-mock-space-id",
					Name:        req.Name,
					Description: req.Description,
				},
			}, nil
		},
	}
}

// List calls the mock implementation
func (s *MockSpacesService) List(ctx context.Context) (*anytype.SpaceListResponse, error) {
	return s.ListFunc(ctx)
}

// Create calls the mock implementation
func (s *MockSpacesService) Create(ctx context.Context, req anytype.CreateSpaceRequest) (*anytype.CreateSpaceResponse, error) {
	return s.CreateFunc(ctx, req)
}

// MockSpaceService implements the anytype.SpaceContext interface for testing
type MockSpaceService struct {
	CurrentSpaceID     string
	GetFunc            func(ctx context.Context) (*anytype.SpaceResponse, error)
	MockTypeService    *MockTypeService
	MockObjectsService *MockObjectsService
	MockMembersService *MockMembersService
	MockPropertyClient *MockSpacePropertyClient
	MockSearchFunc     func(ctx context.Context, req anytype.SearchRequest) (*anytype.SearchResponse, error)
}

// NewMockSpaceService creates a new instance of MockSpaceService with default implementations
func NewMockSpaceService() *MockSpaceService {
	typeService := NewMockTypeService()
	objectsService := NewMockObjectsService()
	membersService := NewMockMembersService()
	propertyClient := NewMockSpacePropertyClient()

	return &MockSpaceService{
		CurrentSpaceID: "mock-space-id",
		GetFunc: func(ctx context.Context) (*anytype.SpaceResponse, error) {
			return &anytype.SpaceResponse{
				Space: anytype.Space{
					ID:          "mock-space-id",
					Name:        "Mock Space",
					Description: "A mock space for testing",
				},
			}, nil
		},
		MockTypeService:    typeService,
		MockObjectsService: objectsService,
		MockMembersService: membersService,
		MockPropertyClient: propertyClient,
		MockSearchFunc: func(ctx context.Context, req anytype.SearchRequest) (*anytype.SearchResponse, error) {
			// Check if this is the specific search we're testing for
			if req.Query == "UniqueTestSearchTerm2025" {
				// Return a matching object with the expected ID format
				return &anytype.SearchResponse{
					Data: []anytype.Object{
						{
							ID:      "new-mock-object-id",
							Name:    "Test Object with UniqueTestSearchTerm2025",
							SpaceID: "mock-space-id",
							TypeKey: "page",
							Layout:  "basic",
						},
					},
				}, nil
			}

			// Default response for other searches
			return &anytype.SearchResponse{
				Data: []anytype.Object{
					{
						ID:      "mock-object-id",
						Name:    "Mock Object",
						SpaceID: "mock-space-id",
						TypeKey: "page",
						Layout:  "basic",
					},
				},
			}, nil
		},
	}
}

// Get calls the mock implementation
func (s *MockSpaceService) Get(ctx context.Context) (*anytype.SpaceResponse, error) {
	return s.GetFunc(ctx)
}

// Types returns the mock type service
func (s *MockSpaceService) Types() anytype.TypeClient {
	return s.MockTypeService
}

// Type returns the mock type context for a specific type
func (s *MockSpaceService) Type(typeKey string) anytype.TypeContext {
	return NewMockTypeContextService(typeKey)
}

// Objects returns the mock objects service
func (s *MockSpaceService) Objects() anytype.ObjectClient {
	return s.MockObjectsService
}

// Object returns the mock object context for a specific object
func (s *MockSpaceService) Object(objectID string) anytype.ObjectContext {
	s.MockObjectsService.SetCurrentObjectID(objectID)
	return s.MockObjectsService
}

// Members returns the mock members service
func (s *MockSpaceService) Members() anytype.MemberClient {
	return s.MockMembersService
}

// Member returns the mock member service for a specific member
func (s *MockSpaceService) Member(memberID string) anytype.MemberContext {
	return s.MockMembersService.Member(memberID)
}

// Lists returns the mock lists service
func (s *MockSpaceService) Lists() anytype.ListClient {
	return NewMockListClient()
}

// List returns the mock list service
func (s *MockSpaceService) List(listID string) anytype.ListContext {
	return NewMockListService(listID)
}

// Search calls the mock implementation
func (s *MockSpaceService) Search(ctx context.Context, req anytype.SearchRequest) (*anytype.SearchResponse, error) {
	return s.MockSearchFunc(ctx, req)
}

// Properties returns the mock property client
func (s *MockSpaceService) Properties() anytype.SpacePropertyClient {
	return s.MockPropertyClient
}

// MockSearchClient implements the anytype.SearchClient interface for testing
type MockSearchClient struct {
	SearchFunc func(ctx context.Context, req anytype.SearchRequest) (*anytype.SearchResponse, error)
}

// NewMockSearchClient creates a new instance of MockSearchClient
func NewMockSearchClient() *MockSearchClient {
	return &MockSearchClient{
		SearchFunc: func(ctx context.Context, req anytype.SearchRequest) (*anytype.SearchResponse, error) {
			// Check if the query contains the test search term
			if req.Query == "UniqueTestSearchTerm2025" {
				// Return a match for the specific test search term used in the test
				return &anytype.SearchResponse{
					Data: []anytype.Object{
						{
							ID:      "new-mock-object-id", // This ID matches what the test expects
							Name:    "Test Object with UniqueTestSearchTerm2025",
							SpaceID: "mock-space-id",
							TypeKey: "page",
						},
					},
				}, nil
			}

			// For any other search term, return default results
			return &anytype.SearchResponse{
				Data: []anytype.Object{
					{
						ID:      "mock-search-result-id",
						Name:    "Mock Search Result",
						SpaceID: "mock-space-id",
						TypeKey: "page",
					},
				},
			}, nil
		},
	}
}

// Search calls the mock implementation
func (c *MockSearchClient) Search(ctx context.Context, req anytype.SearchRequest) (*anytype.SearchResponse, error) {
	return c.SearchFunc(ctx, req)
}
