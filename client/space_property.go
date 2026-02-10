package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/epheo/anytype-go"
)

type SpacePropertyClientImpl struct {
	client  *ClientImpl
	spaceID string
}

func (pc *SpacePropertyClientImpl) List(ctx context.Context) ([]anytype.Property, error) {
	endpoint := fmt.Sprintf("/spaces/%s/properties", pc.spaceID)

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

func (pc *SpacePropertyClientImpl) Create(ctx context.Context, request anytype.CreatePropertyRequest) (*anytype.PropertyResponse, error) {
	endpoint := fmt.Sprintf("/spaces/%s/properties", pc.spaceID)

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

type PropertyContextImpl struct {
	client     *ClientImpl
	spaceID    string
	propertyID string
}

func (pc *PropertyContextImpl) Get(ctx context.Context) (*anytype.PropertyResponse, error) {
	endpoint := fmt.Sprintf("/spaces/%s/properties/%s", pc.spaceID, pc.propertyID)

	req, err := pc.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response anytype.PropertyResponse
	if err := pc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (pc *PropertyContextImpl) Update(ctx context.Context, request anytype.UpdatePropertyRequest) (*anytype.PropertyResponse, error) {
	endpoint := fmt.Sprintf("/spaces/%s/properties/%s", pc.spaceID, pc.propertyID)

	req, err := pc.client.newRequest(ctx, http.MethodPatch, endpoint, request)
	if err != nil {
		return nil, err
	}

	var response anytype.PropertyResponse
	if err := pc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (pc *PropertyContextImpl) Delete(ctx context.Context) (*anytype.PropertyResponse, error) {
	endpoint := fmt.Sprintf("/spaces/%s/properties/%s", pc.spaceID, pc.propertyID)

	req, err := pc.client.newRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response anytype.PropertyResponse
	if err := pc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (pc *PropertyContextImpl) Tags() anytype.TagClient {
	return &TagClientImpl{
		client:     pc.client,
		spaceID:    pc.spaceID,
		propertyID: pc.propertyID,
	}
}

func (pc *PropertyContextImpl) Tag(tagID string) anytype.TagContext {
	return &TagContextImpl{
		client:     pc.client,
		spaceID:    pc.spaceID,
		propertyID: pc.propertyID,
		tagID:      tagID,
	}
}
