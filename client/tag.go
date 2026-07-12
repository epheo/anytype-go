package client

import (
	"context"
	"fmt"
	"iter"
	"net/http"

	"github.com/epheo/anytype-go"
)

type TagClientImpl struct {
	client     *ClientImpl
	spaceID    string
	propertyID string
}

func (tc *TagClientImpl) List(ctx context.Context, opts ...anytype.ListOption) (*anytype.Page[anytype.Tag], error) {
	endpoint := withListParams(fmt.Sprintf("/spaces/%s/properties/%s/tags", tc.spaceID, tc.propertyID), opts)

	req, err := tc.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response anytype.Page[anytype.Tag]
	if err := tc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (tc *TagClientImpl) All(ctx context.Context, opts ...anytype.ListOption) iter.Seq2[anytype.Tag, error] {
	return paginate(ctx, tc.List, opts...)
}

func (tc *TagClientImpl) Create(ctx context.Context, request anytype.CreateTagRequest) (*anytype.TagResponse, error) {
	endpoint := fmt.Sprintf("/spaces/%s/properties/%s/tags", tc.spaceID, tc.propertyID)

	req, err := tc.client.newRequest(ctx, http.MethodPost, endpoint, request)
	if err != nil {
		return nil, err
	}

	var response anytype.TagResponse
	if err := tc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

type TagContextImpl struct {
	client     *ClientImpl
	spaceID    string
	propertyID string
	tagID      string
}

func (tc *TagContextImpl) Get(ctx context.Context) (*anytype.TagResponse, error) {
	endpoint := fmt.Sprintf("/spaces/%s/properties/%s/tags/%s", tc.spaceID, tc.propertyID, tc.tagID)

	req, err := tc.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response anytype.TagResponse
	if err := tc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (tc *TagContextImpl) Update(ctx context.Context, request anytype.UpdateTagRequest) (*anytype.TagResponse, error) {
	endpoint := fmt.Sprintf("/spaces/%s/properties/%s/tags/%s", tc.spaceID, tc.propertyID, tc.tagID)

	req, err := tc.client.newRequest(ctx, http.MethodPatch, endpoint, request)
	if err != nil {
		return nil, err
	}

	var response anytype.TagResponse
	if err := tc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (tc *TagContextImpl) Delete(ctx context.Context) (*anytype.TagResponse, error) {
	endpoint := fmt.Sprintf("/spaces/%s/properties/%s/tags/%s", tc.spaceID, tc.propertyID, tc.tagID)

	req, err := tc.client.newRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response anytype.TagResponse
	if err := tc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
