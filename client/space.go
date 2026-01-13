// Package client implements the client interfaces defined in the anytype package
package client

import (
	"context"
	"net/http"

	"github.com/epheo/anytype-go"
)

// SpaceClientImpl implements the SpaceClient interface
type SpaceClientImpl struct {
	client *ClientImpl
}

// Create creates a new space
func (sc *SpaceClientImpl) Create(ctx context.Context, request anytype.CreateSpaceRequest) (*anytype.CreateSpaceResponse, error) {
	// Create HTTP request
	req, err := sc.client.newRequest(ctx, http.MethodPost, "/spaces", request)
	if err != nil {
		return nil, err
	}

	// Make the request and parse the response
	response := &anytype.CreateSpaceResponse{}
	err = sc.client.doRequest(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// List returns all spaces accessible to the user
func (sc *SpaceClientImpl) List(ctx context.Context) (*anytype.SpaceListResponse, error) {
	// Create HTTP request
	req, err := sc.client.newRequest(ctx, http.MethodGet, "/spaces", nil)
	if err != nil {
		return nil, err
	}

	// Make the request and parse the response
	response := &anytype.SpaceListResponse{}
	err = sc.client.doRequest(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SpaceContextImpl implements the SpaceContext interface
type SpaceContextImpl struct {
	client  *ClientImpl
	spaceID string
}

// Lists returns a ListClient for this space
func (sc *SpaceContextImpl) Lists() anytype.ListClient {
	return &ListClientImpl{
		client:  sc.client,
		spaceID: sc.spaceID,
	}
}

// Get retrieves information about this space
func (sc *SpaceContextImpl) Get(ctx context.Context) (*anytype.SpaceResponse, error) {
	// Create HTTP request
	endpoint := "/spaces/" + sc.spaceID
	req, err := sc.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Make the request and parse the response
	response := &anytype.SpaceResponse{}
	err = sc.client.doRequest(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Objects returns an ObjectClient for this space
func (sc *SpaceContextImpl) Objects() anytype.ObjectClient {
	return &ObjectClientImpl{
		client:  sc.client,
		spaceID: sc.spaceID,
	}
}

// Object returns an ObjectContext for a specific object in this space
func (sc *SpaceContextImpl) Object(objectID string) anytype.ObjectContext {
	return &ObjectContextImpl{
		client:   sc.client,
		spaceID:  sc.spaceID,
		objectID: objectID,
	}
}

// List returns a ListContext for a specific list in this space
func (sc *SpaceContextImpl) List(listID string) anytype.ListContext {
	return &ListContextImpl{
		client:  sc.client,
		spaceID: sc.spaceID,
		listID:  listID,
	}
}

// Types returns a TypeClient for this space
func (sc *SpaceContextImpl) Types() anytype.TypeClient {
	return &TypeClientImpl{
		client:  sc.client,
		spaceID: sc.spaceID,
	}
}

// Type returns a TypeContext for a specific type in this space
func (sc *SpaceContextImpl) Type(typeID string) anytype.TypeContext {
	return &TypeContextImpl{
		client:  sc.client,
		spaceID: sc.spaceID,
		typeID:  typeID,
	}
}

// Search searches for objects within this space
func (sc *SpaceContextImpl) Search(ctx context.Context, request anytype.SearchRequest) (*anytype.SearchResponse, error) {
	endpoint := "/spaces/" + sc.spaceID + "/search"

	// Create HTTP request with the request struct directly
	// The newRequest method will handle JSON marshaling
	req, err := sc.client.newRequest(ctx, http.MethodPost, endpoint, request)
	if err != nil {
		return nil, err
	}

	// Make the request and parse the response
	response := &anytype.SearchResponse{}
	err = sc.client.doRequest(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Members returns a MemberClient for this space
func (sc *SpaceContextImpl) Members() anytype.MemberClient {
	return &MemberClientImpl{
		client:  sc.client,
		spaceID: sc.spaceID,
	}
}

// Member returns a MemberContext for a specific member in this space
func (sc *SpaceContextImpl) Member(memberID string) anytype.MemberContext {
	return &MemberContextImpl{
		client:   sc.client,
		spaceID:  sc.spaceID,
		memberID: memberID,
	}
}

// Properties returns a SpacePropertyClient for this space
func (sc *SpaceContextImpl) Properties() anytype.SpacePropertyClient {
	return &SpacePropertyClientImpl{
		client:  sc.client,
		spaceID: sc.spaceID,
	}
}

// SpacePropertyClientImpl implements the SpacePropertyClient interface
type SpacePropertyClientImpl struct {
	client  *ClientImpl
	spaceID string
}

// Create creates a new property in the space
func (pc *SpacePropertyClientImpl) Create(ctx context.Context, request anytype.CreatePropertyRequest) (*anytype.PropertyResponse, error) {
	endpoint := "/spaces/" + pc.spaceID + "/properties"

	req, err := pc.client.newRequest(ctx, http.MethodPost, endpoint, request)
	if err != nil {
		return nil, err
	}

	var response anytype.PropertyResponse
	if err := pc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// List returns all properties in the space
func (pc *SpacePropertyClientImpl) List(ctx context.Context) ([]anytype.Property, error) {
	endpoint := "/spaces/" + pc.spaceID + "/properties"

	req, err := pc.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []anytype.Property `json:"data"`
	}
	if err := pc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}
