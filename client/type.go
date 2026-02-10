package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/epheo/anytype-go"
)

type TypeClientImpl struct {
	client  *ClientImpl
	spaceID string
}

func (tc *TypeClientImpl) List(ctx context.Context) ([]anytype.Type, error) {
	req, err := tc.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/spaces/%s/types", tc.spaceID), nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []anytype.Type `json:"data"`
	}
	if err := tc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (tc *TypeClientImpl) Get(ctx context.Context, typeKey string) (*anytype.Type, error) {
	req, err := tc.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/spaces/%s/types/%s", tc.spaceID, typeKey), nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Type anytype.Type `json:"type"`
	}
	if err := tc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response.Type, nil
}

// GetKeyByName looks up a type key by its name
func (tc *TypeClientImpl) GetKeyByName(ctx context.Context, name string) (string, error) {
	types, err := tc.List(ctx)
	if err != nil {
		return "", err
	}

	// Linear search is acceptable because type lists are typically small (<100 items).
	// The API doesn't provide name-based lookup, so we must fetch all and filter.
	for _, t := range types {
		if t.Name == name {
			return t.Key, nil
		}
	}

	return "", anytype.ErrTypeNotFound
}

func (tc *TypeClientImpl) Create(ctx context.Context, request anytype.CreateTypeRequest) (*anytype.TypeResponse, error) {
	endpoint := fmt.Sprintf("/spaces/%s/types", tc.spaceID)

	req, err := tc.client.newRequest(ctx, http.MethodPost, endpoint, request)
	if err != nil {
		return nil, err
	}

	var response anytype.TypeResponse
	err = tc.client.doRequest(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type TypeContextImpl struct {
	client  *ClientImpl
	spaceID string
	typeID  string
}

func (tc *TypeContextImpl) Get(ctx context.Context) (*anytype.TypeResponse, error) {
	req, err := tc.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/spaces/%s/types/%s", tc.spaceID, tc.typeID), nil)
	if err != nil {
		return nil, err
	}

	var response anytype.TypeResponse
	if err := tc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (tc *TypeContextImpl) Update(ctx context.Context, request anytype.UpdateTypeRequest) (*anytype.TypeResponse, error) {
	endpoint := fmt.Sprintf("/spaces/%s/types/%s", tc.spaceID, tc.typeID)

	req, err := tc.client.newRequest(ctx, http.MethodPatch, endpoint, request)
	if err != nil {
		return nil, err
	}

	var response anytype.TypeResponse
	if err := tc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (tc *TypeContextImpl) Delete(ctx context.Context) (*anytype.TypeResponse, error) {
	endpoint := fmt.Sprintf("/spaces/%s/types/%s", tc.spaceID, tc.typeID)

	req, err := tc.client.newRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response anytype.TypeResponse
	if err := tc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (tc *TypeContextImpl) Templates() anytype.TemplateClient {
	return &TemplateClientImpl{
		client:  tc.client,
		spaceID: tc.spaceID,
		typeID:  tc.typeID,
	}
}

func (tc *TypeContextImpl) Template(templateID string) anytype.TemplateContext {
	return &TemplateContextImpl{
		client:     tc.client,
		spaceID:    tc.spaceID,
		typeID:     tc.typeID,
		templateID: templateID,
	}
}
