package client

import (
	"context"
	"iter"
	"net/http"

	"github.com/epheo/anytype-go"
)

type SearchClientImpl struct {
	client *ClientImpl
}

// Search searches for objects across all spaces
func (sc *SearchClientImpl) Search(ctx context.Context, request anytype.SearchRequest, opts ...anytype.ListOption) (*anytype.Page[anytype.Object], error) {
	endpoint := withListParams("/search", opts)

	req, err := sc.client.newRequest(ctx, http.MethodPost, endpoint, request)
	if err != nil {
		return nil, err
	}

	response := &anytype.Page[anytype.Object]{}
	if err := sc.client.doRequest(req, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (sc *SearchClientImpl) All(ctx context.Context, request anytype.SearchRequest, opts ...anytype.ListOption) iter.Seq2[anytype.Object, error] {
	return paginate(ctx, func(ctx context.Context, o ...anytype.ListOption) (*anytype.Page[anytype.Object], error) {
		return sc.Search(ctx, request, o...)
	}, opts...)
}
