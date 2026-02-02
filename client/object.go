package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/epheo/anytype-go"
	"github.com/epheo/anytype-go/options"
)

// ObjectClientImpl implements the ObjectClient interface
type ObjectClientImpl struct {
	client  *ClientImpl
	spaceID string
}

// List returns all objects in the space
func (oc *ObjectClientImpl) List(ctx context.Context, opts ...options.ListOption) ([]anytype.Object, error) {
	endpoint := fmt.Sprintf("/spaces/%s/objects", oc.spaceID)

	// Apply all list options to create query parameters
	listOpts := options.ApplyListOptions(opts...)
	if listOpts.Limit > 0 {
		endpoint = fmt.Sprintf("%s?limit=%d", endpoint, listOpts.Limit)
	}
	if listOpts.Offset > 0 {
		if listOpts.Limit > 0 {
			endpoint = fmt.Sprintf("%s&offset=%d", endpoint, listOpts.Offset)
		} else {
			endpoint = fmt.Sprintf("%s?offset=%d", endpoint, listOpts.Offset)
		}
	}

	req, err := oc.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data       []anytype.Object           `json:"data"`
		Pagination options.PaginationMetadata `json:"pagination"`
	}

	err = oc.client.doRequest(req, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// Create creates a new object in the space
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

// ObjectContextImpl implements the ObjectContext interface
type ObjectContextImpl struct {
	client   *ClientImpl
	spaceID  string
	objectID string
}

// Get retrieves the object
func (oc *ObjectContextImpl) Get(ctx context.Context) (*anytype.ObjectResponse, error) {
	endpoint := fmt.Sprintf("/spaces/%s/objects/%s", oc.spaceID, oc.objectID)

	req, err := oc.client.newRequest(ctx, http.MethodGet, endpoint, nil)
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

// Update updates the object
func (oc *ObjectContextImpl) Update(ctx context.Context, request anytype.UpdateObjectRequest) (*anytype.ObjectResponse, error) {
	endpoint := fmt.Sprintf("/spaces/%s/objects/%s", oc.spaceID, oc.objectID)

	req, err := oc.client.newRequest(ctx, http.MethodPatch, endpoint, request)
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

// Delete deletes the object
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

// Blocks returns a BlockClient for this object
func (oc *ObjectContextImpl) Blocks() anytype.BlockClient {
	return &BlockClientImpl{
		client:   oc.client,
		spaceID:  oc.spaceID,
		objectID: oc.objectID,
	}
}

// Properties returns a PropertyClient for this object
func (oc *ObjectContextImpl) Properties() anytype.PropertyClient {
	return &PropertyClientImpl{
		client:   oc.client,
		spaceID:  oc.spaceID,
		objectID: oc.objectID,
	}
}

// UpdateName updates the name of the object by updating the name property
func (oc *ObjectContextImpl) UpdateName(ctx context.Context, name string) error {
	// Since there's no direct UPDATE endpoint, we'll use the property update
	return oc.Properties().Set(ctx, "name", name)
}

// UpdateIcon updates the icon of the object by updating the icon property
func (oc *ObjectContextImpl) UpdateIcon(ctx context.Context, icon *anytype.Icon) error {
	// Since there's no direct UPDATE endpoint, we'll use the property update
	return oc.Properties().Set(ctx, "icon", icon)
}

// Export exports the object in the specified format
func (oc *ObjectContextImpl) Export(ctx context.Context, format string) (*anytype.ExportResult, error) {
	// Export by fetching the object and returning its Markdown field
	resp, err := oc.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &anytype.ExportResult{
		Markdown: resp.Object.Markdown,
	}, nil
}
