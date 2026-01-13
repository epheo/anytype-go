package mocks

import (
	"context"

	"github.com/epheo/anytype-go"
)

// MockSpacePropertyClient implements the anytype.SpacePropertyClient interface for testing
type MockSpacePropertyClient struct {
	CreateFunc func(ctx context.Context, req anytype.CreatePropertyRequest) (*anytype.PropertyResponse, error)
	ListFunc   func(ctx context.Context) ([]anytype.Property, error)
}

// NewMockSpacePropertyClient creates a new instance of MockSpacePropertyClient with default implementations
func NewMockSpacePropertyClient() *MockSpacePropertyClient {
	return &MockSpacePropertyClient{
		CreateFunc: func(ctx context.Context, req anytype.CreatePropertyRequest) (*anytype.PropertyResponse, error) {
			return &anytype.PropertyResponse{
				Property: anytype.Property{
					Key:    req.Key,
					Name:   req.Name,
					Format: req.Format,
				},
			}, nil
		},
		ListFunc: func(ctx context.Context) ([]anytype.Property, error) {
			return []anytype.Property{
				{
					Key:    "mock-property-key",
					Name:   "Mock Property",
					Format: "text",
				},
				{
					Key:    "mock-number-property",
					Name:   "Mock Number Property",
					Format: "number",
				},
			}, nil
		},
	}
}

// Create calls the mock implementation
func (c *MockSpacePropertyClient) Create(ctx context.Context, req anytype.CreatePropertyRequest) (*anytype.PropertyResponse, error) {
	return c.CreateFunc(ctx, req)
}

// List calls the mock implementation
func (c *MockSpacePropertyClient) List(ctx context.Context) ([]anytype.Property, error) {
	return c.ListFunc(ctx)
}
