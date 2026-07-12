package client

import (
	"context"
	"fmt"
	"iter"
	"net/http"

	"github.com/epheo/anytype-go"
)

type TemplateClientImpl struct {
	client  *ClientImpl
	spaceID string
	typeID  string
}

func (tc *TemplateClientImpl) List(ctx context.Context, opts ...anytype.ListOption) (*anytype.Page[anytype.Object], error) {
	endpoint := withListParams(fmt.Sprintf("/spaces/%s/types/%s/templates", tc.spaceID, tc.typeID), opts)

	req, err := tc.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response anytype.Page[anytype.Object]
	if err := tc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (tc *TemplateClientImpl) All(ctx context.Context, opts ...anytype.ListOption) iter.Seq2[anytype.Object, error] {
	return paginate(ctx, tc.List, opts...)
}

type TemplateContextImpl struct {
	client     *ClientImpl
	spaceID    string
	typeID     string
	templateID string
}

func (tc *TemplateContextImpl) Get(ctx context.Context) (*anytype.TemplateResponse, error) {
	endpoint := fmt.Sprintf("/spaces/%s/types/%s/templates/%s", tc.spaceID, tc.typeID, tc.templateID)

	req, err := tc.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response anytype.TemplateResponse
	if err := tc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
