package anytype

import (
	"context"
)

// SpaceClient provides operations on spaces
type SpaceClient interface {
	// List returns all spaces accessible to the user
	List(ctx context.Context) (*SpaceListResponse, error)

	// Create creates a new space
	Create(ctx context.Context, request CreateSpaceRequest) (*CreateSpaceResponse, error)
}

// SpaceContext provides operations within a specific space
type SpaceContext interface {
	// Get retrieves information about this space
	Get(ctx context.Context) (*SpaceResponse, error)

	// Objects returns an ObjectClient for this space
	Objects() ObjectClient

	// Object returns an ObjectContext for a specific object in this space
	Object(objectID string) ObjectContext

	// Types returns a TypeClient for this space
	Types() TypeClient

	// Type returns a TypeContext for a specific type in this space
	Type(typeID string) TypeContext

	// Search searches for objects within this space
	Search(ctx context.Context, request SearchRequest) (*SearchResponse, error)

	// Lists returns a ListClient for this space
	Lists() ListClient

	// List returns a ListContext for a specific list in this space
	List(listID string) ListContext

	// Members returns a MemberClient for this space
	Members() MemberClient

	// Member returns a MemberContext for a specific member in this space
	Member(memberID string) MemberContext
}

// Space represents an Anytype workspace/space
type Space struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Icon        *Icon  `json:"icon,omitempty"`
	NetworkID   string `json:"network_id,omitempty"`
	GatewayURL  string `json:"gateway_url,omitempty"`
	Object      string `json:"object"` // Data model identifier (enum: "space", "chat")

	// Legacy fields - may be removed in future versions
	HomeID       string `json:"home_id,omitempty"`
	ArchiveID    string `json:"archive_id,omitempty"`
	ProfileID    string `json:"profile_id,omitempty"`
	CreatedAt    int64  `json:"created_at,omitempty"`
	LastOpenedAt int64  `json:"last_opened_at,omitempty"`
}

// SpaceListResponse represents the response from List spaces
type SpaceListResponse struct {
	Data []Space `json:"data"`
}

// CreateSpaceResponse represents the response from Create space
type CreateSpaceResponse struct {
	Space Space `json:"space"`
}

// SpaceResponse represents the response from Get space
type SpaceResponse struct {
	Space Space `json:"space"`
}

// CreateSpaceRequest represents the request body for creating a new space
type CreateSpaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Icon        *Icon  `json:"icon,omitempty"`
}

// UpdateSpaceRequest represents the request body for updating a space
type UpdateSpaceRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Icon        *Icon   `json:"icon,omitempty"`
}
