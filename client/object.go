package client

import (
	"context"
	"fmt"
	"iter"
	"net/http"
	"net/url"

	"github.com/epheo/anytype-go"
)

type ObjectClientImpl struct {
	client  *ClientImpl
	spaceID string
}

func (oc *ObjectClientImpl) List(ctx context.Context, opts ...anytype.ListOption) (*anytype.Page[anytype.Object], error) {
	endpoint := withListParams(fmt.Sprintf("/spaces/%s/objects", oc.spaceID), opts)

	req, err := oc.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response anytype.Page[anytype.Object]
	if err := oc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (oc *ObjectClientImpl) All(ctx context.Context, opts ...anytype.ListOption) iter.Seq2[anytype.Object, error] {
	return paginate(ctx, oc.List, opts...)
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
