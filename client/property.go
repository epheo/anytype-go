package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/epheo/anytype-go"
)

type PropertyClientImpl struct {
	client   *ClientImpl
	spaceID  string
	objectID string
}

func (pc *PropertyClientImpl) Get(ctx context.Context, key string) (*anytype.Property, error) {
	urlPath := fmt.Sprintf("/spaces/%s/objects/%s/properties/%s", pc.spaceID, pc.objectID, key)

	req, err := pc.client.newRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Property anytype.Property `json:"property"`
	}

	err = pc.client.doRequest(req, &response)
	if err != nil {
		return nil, err
	}

	return &response.Property, nil
}

func (pc *PropertyClientImpl) Set(ctx context.Context, key string, value interface{}) error {
	urlPath := fmt.Sprintf("/spaces/%s/objects/%s/properties/%s", pc.spaceID, pc.objectID, key)

	body := map[string]interface{}{
		"value": value,
	}

	req, err := pc.client.newRequest(ctx, http.MethodPut, urlPath, body)
	if err != nil {
		return err
	}

	return pc.client.doRequest(req, nil)
}

func (pc *PropertyClientImpl) Delete(ctx context.Context, key string) error {
	urlPath := fmt.Sprintf("/spaces/%s/objects/%s/properties/%s", pc.spaceID, pc.objectID, key)

	req, err := pc.client.newRequest(ctx, http.MethodDelete, urlPath, nil)
	if err != nil {
		return err
	}

	return pc.client.doRequest(req, nil)
}

func (pc *PropertyClientImpl) List(ctx context.Context) ([]anytype.Property, error) {
	urlPath := fmt.Sprintf("/spaces/%s/objects/%s/properties", pc.spaceID, pc.objectID)

	req, err := pc.client.newRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Properties []anytype.Property `json:"properties"`
	}

	err = pc.client.doRequest(req, &response)
	if err != nil {
		return nil, err
	}

	return response.Properties, nil
}
