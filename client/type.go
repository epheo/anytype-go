package client

import (
	"context"
	"fmt"
	"iter"
	"net/http"

	"github.com/epheo/anytype-go"
)

type TypeClientImpl struct {
	client  *ClientImpl
	spaceID string
}

func (tc *TypeClientImpl) List(ctx context.Context, opts ...anytype.ListOption) (*anytype.Page[anytype.Type], error) {
	endpoint := withListParams(fmt.Sprintf("/spaces/%s/types", tc.spaceID), opts)

	req, err := tc.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response anytype.Page[anytype.Type]
	if err := tc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (tc *TypeClientImpl) All(ctx context.Context, opts ...anytype.ListOption) iter.Seq2[anytype.Type, error] {
	return paginate(ctx, tc.List, opts...)
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

// GetKeyByName looks up a type key by its name.
// The API has no name-based lookup, so we scan every page and filter.
func (tc *TypeClientImpl) GetKeyByName(ctx context.Context, name string) (string, error) {
	for t, err := range tc.All(ctx) {
		if err != nil {
			return "", err
		}
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
