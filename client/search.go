package client

import (
	"context"
	"net/http"

	"github.com/epheo/anytype-go"
)

type SearchClientImpl struct {
	client *ClientImpl
}

// Search searches for objects across all spaces
func (sc *SearchClientImpl) Search(ctx context.Context, request anytype.SearchRequest) (*anytype.SearchResponse, error) {
	endpoint := "/search"

	req, err := sc.client.newRequest(ctx, http.MethodPost, endpoint, request)
	if err != nil {
		return nil, err
	}

	response := &anytype.SearchResponse{}

	err = sc.client.doRequest(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
