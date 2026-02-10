package anytype

import (
	"context"
)

type SpaceClient interface {
	List(ctx context.Context) (*SpaceListResponse, error)
	Create(ctx context.Context, request CreateSpaceRequest) (*CreateSpaceResponse, error)
}

type SpaceContext interface {
	Get(ctx context.Context) (*SpaceResponse, error)
	Objects() ObjectClient
	Object(objectID string) ObjectContext
	Types() TypeClient
	Type(typeID string) TypeContext
	Search(ctx context.Context, request SearchRequest) (*SearchResponse, error)
	Lists() ListClient
	List(listID string) ListContext
	Members() MemberClient
	Member(memberID string) MemberContext
}

type Space struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Icon        *Icon  `json:"icon,omitempty"`
	NetworkID   string `json:"network_id,omitempty"`
	GatewayURL  string `json:"gateway_url,omitempty"`
	Object      string `json:"object"`

	// Legacy fields - may be removed in future versions
	HomeID       string `json:"home_id,omitempty"`
	ArchiveID    string `json:"archive_id,omitempty"`
	ProfileID    string `json:"profile_id,omitempty"`
	CreatedAt    int64  `json:"created_at,omitempty"`
	LastOpenedAt int64  `json:"last_opened_at,omitempty"`
}

type SpaceListResponse struct {
	Data []Space `json:"data"`
}

type CreateSpaceResponse struct {
	Space Space `json:"space"`
}

type SpaceResponse struct {
	Space Space `json:"space"`
}

type CreateSpaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Icon        *Icon  `json:"icon,omitempty"`
}

type UpdateSpaceRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Icon        *Icon   `json:"icon,omitempty"`
}
