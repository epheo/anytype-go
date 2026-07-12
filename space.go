package anytype

import (
	"context"
	"iter"
)

type SpaceClient interface {
	List(ctx context.Context, opts ...ListOption) (*Page[Space], error)
	All(ctx context.Context, opts ...ListOption) iter.Seq2[Space, error]
	Create(ctx context.Context, request CreateSpaceRequest) (*SpaceResponse, error)
}

type SpaceContext interface {
	Get(ctx context.Context) (*SpaceResponse, error)
	Update(ctx context.Context, request UpdateSpaceRequest) (*SpaceResponse, error)
	Objects() ObjectClient
	Object(objectID string) ObjectContext
	Types() TypeClient
	Type(typeID string) TypeContext
	Search(ctx context.Context, request SearchRequest, opts ...ListOption) (*Page[Object], error)
	List(listID string) ListContext
	Members() MemberClient
	Member(memberID string) MemberContext
	Properties() SpacePropertyClient
	Property(propertyID string) PropertyContext
}

type PropertyContext interface {
	Get(ctx context.Context) (*PropertyResponse, error)
	Update(ctx context.Context, request UpdatePropertyRequest) (*PropertyResponse, error)
	Delete(ctx context.Context) (*PropertyResponse, error)
	Tags() TagClient
	Tag(tagID string) TagContext
}

type TagContext interface {
	Get(ctx context.Context) (*TagResponse, error)
	Update(ctx context.Context, request UpdateTagRequest) (*TagResponse, error)
	Delete(ctx context.Context) (*TagResponse, error)
}

type Space struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Icon        *Icon  `json:"icon,omitempty"`
	NetworkID   string `json:"network_id,omitempty"`
	GatewayURL  string `json:"gateway_url,omitempty"`
	Object      string `json:"object"`
}

type SpaceResponse struct {
	Space Space `json:"space"`
}

type CreateSpaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type UpdateSpaceRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}
