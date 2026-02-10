package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/epheo/anytype-go"
)

type ObjectClientImpl struct {
	client  *ClientImpl
	spaceID string
}

func (oc *ObjectClientImpl) List(ctx context.Context, opts ...anytype.ListOption) ([]anytype.Object, error) {
	endpoint := fmt.Sprintf("/spaces/%s/objects", oc.spaceID)

	listOpts := anytype.ApplyListOptions(opts...)
	params := url.Values{}
	if listOpts.Limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", listOpts.Limit))
	}
	if listOpts.Offset > 0 {
		params.Set("offset", fmt.Sprintf("%d", listOpts.Offset))
	}
	if encoded := params.Encode(); encoded != "" {
		endpoint += "?" + encoded
	}

	req, err := oc.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data       []anytype.Object   `json:"data"`
		Pagination anytype.Pagination `json:"pagination"`
	}

	err = oc.client.doRequest(req, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (oc *ObjectClientImpl) Create(ctx context.Context, request anytype.CreateObjectRequest) (*anytype.ObjectResponse, error) {
	endpoint := fmt.Sprintf("/spaces/%s/objects", oc.spaceID)

	req, err := oc.client.newRequest(ctx, http.MethodPost, endpoint, request)
	if err != nil {
		return nil, err
	}

	var response anytype.ObjectResponse
	err = oc.client.doRequest(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type ObjectContextImpl struct {
	client   *ClientImpl
	spaceID  string
	objectID string
}

func (oc *ObjectContextImpl) Get(ctx context.Context, opts ...anytype.GetOption) (*anytype.ObjectResponse, error) {
	endpoint := fmt.Sprintf("/spaces/%s/objects/%s", oc.spaceID, oc.objectID)

	getOpts := anytype.ApplyGetOptions(opts...)
	params := url.Values{}
	if getOpts.Format != "" {
		params.Set("format", getOpts.Format)
	}
	if encoded := params.Encode(); encoded != "" {
		endpoint += "?" + encoded
	}

	req, err := oc.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response anytype.ObjectResponse
	if err := oc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (oc *ObjectContextImpl) Update(ctx context.Context, request anytype.UpdateObjectRequest) (*anytype.ObjectResponse, error) {
	endpoint := fmt.Sprintf("/spaces/%s/objects/%s", oc.spaceID, oc.objectID)

	req, err := oc.client.newRequest(ctx, http.MethodPatch, endpoint, request)
	if err != nil {
		return nil, err
	}

	var response anytype.ObjectResponse
	if err := oc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (oc *ObjectContextImpl) Delete(ctx context.Context) (*anytype.ObjectResponse, error) {
	endpoint := fmt.Sprintf("/spaces/%s/objects/%s", oc.spaceID, oc.objectID)

	req, err := oc.client.newRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response anytype.ObjectResponse
	err = oc.client.doRequest(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
