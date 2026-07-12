package client

import (
	"context"
	"iter"
	"net/http"

	"github.com/epheo/anytype-go"
)

type AddObjectsToListRequest struct {
	Objects []string `json:"objects"`
}

type ListContextImpl struct {
	client  *ClientImpl
	spaceID string
	listID  string
}

func (lc *ListContextImpl) Views() anytype.ViewClient {
	return &ViewClientImpl{
		client:  lc.client,
		spaceID: lc.spaceID,
		listID:  lc.listID,
	}
}

func (lc *ListContextImpl) View(viewID string) anytype.ViewContext {
	return &ViewContextImpl{
		client:  lc.client,
		spaceID: lc.spaceID,
		listID:  lc.listID,
		viewID:  viewID,
	}
}

func (lc *ListContextImpl) Objects() anytype.ObjectListClient {
	return &ObjectListClientImpl{
		client:  lc.client,
		spaceID: lc.spaceID,
		listID:  lc.listID,
	}
}

func (lc *ListContextImpl) Object(objectID string) anytype.ObjectListContext {
	return &ObjectListContextImpl{
		client:   lc.client,
		spaceID:  lc.spaceID,
		listID:   lc.listID,
		objectID: objectID,
	}
}

type ViewClientImpl struct {
	client  *ClientImpl
	spaceID string
	listID  string
}

func (vc *ViewClientImpl) List(ctx context.Context, opts ...anytype.ListOption) (*anytype.Page[anytype.ListView], error) {
	endpoint := withListParams("/spaces/"+vc.spaceID+"/lists/"+vc.listID+"/views", opts)
	req, err := vc.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	response := &anytype.Page[anytype.ListView]{}
	if err := vc.client.doRequest(req, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (vc *ViewClientImpl) All(ctx context.Context, opts ...anytype.ListOption) iter.Seq2[anytype.ListView, error] {
	return paginate(ctx, vc.List, opts...)
}

type ViewContextImpl struct {
	client  *ClientImpl
	spaceID string
	listID  string
	viewID  string
}

func (vc *ViewContextImpl) Objects() anytype.ObjectViewClient {
	return &ObjectViewClientImpl{
		client:  vc.client,
		spaceID: vc.spaceID,
		listID:  vc.listID,
		viewID:  vc.viewID,
	}
}

type ObjectListClientImpl struct {
	client  *ClientImpl
	spaceID string
	listID  string
}

func (olc *ObjectListClientImpl) Add(ctx context.Context, objectIDs []string) error {
	endpoint := "/spaces/" + olc.spaceID + "/lists/" + olc.listID + "/objects"

	requestBody := AddObjectsToListRequest{
		Objects: objectIDs,
	}

	req, err := olc.client.newRequest(ctx, http.MethodPost, endpoint, requestBody)
	if err != nil {
		return err
	}

	return olc.client.doRequest(req, nil)
}

type ObjectListContextImpl struct {
	client   *ClientImpl
	spaceID  string
	listID   string
	objectID string
}

func (olc *ObjectListContextImpl) Remove(ctx context.Context) error {
	endpoint := "/spaces/" + olc.spaceID + "/lists/" + olc.listID + "/objects/" + olc.objectID
	req, err := olc.client.newRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	return olc.client.doRequest(req, nil)
}

type ObjectViewClientImpl struct {
	client  *ClientImpl
	spaceID string
	listID  string
	viewID  string
}

func (ovc *ObjectViewClientImpl) List(ctx context.Context, opts ...anytype.ListOption) (*anytype.Page[anytype.Object], error) {
	endpoint := withListParams("/spaces/"+ovc.spaceID+"/lists/"+ovc.listID+"/views/"+ovc.viewID+"/objects", opts)
	req, err := ovc.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	response := &anytype.Page[anytype.Object]{}
	if err := ovc.client.doRequest(req, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (ovc *ObjectViewClientImpl) All(ctx context.Context, opts ...anytype.ListOption) iter.Seq2[anytype.Object, error] {
	return paginate(ctx, ovc.List, opts...)
}
