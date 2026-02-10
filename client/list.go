package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/epheo/anytype-go"
)

type ListClientImpl struct {
	client  *ClientImpl
	spaceID string
	listID  string
}

type AddObjectsToListRequest struct {
	Objects []string `json:"objects"`
}

func (lc *ListClientImpl) Add(ctx context.Context, objectIDs []string) error {
	endpoint := "/spaces/" + lc.spaceID + "/lists/" + lc.listID + "/objects"

	requestBody := AddObjectsToListRequest{
		Objects: objectIDs,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := lc.client.newRequest(ctx, http.MethodPost, endpoint, jsonData)
	if err != nil {
		return err
	}

	err = lc.client.doRequest(req, nil)
	if err != nil {
		return err
	}

	return nil
}

func (lc *ListClientImpl) GetViews(ctx context.Context, listID string) ([]anytype.ListView, error) {
	endpoint := "/spaces/" + lc.spaceID + "/lists/" + listID + "/views"
	req, err := lc.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	response := &anytype.ViewListResponse{}
	err = lc.client.doRequest(req, response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (lc *ListClientImpl) GetObjects(ctx context.Context, listID string, viewID string) ([]anytype.Object, error) {
	endpoint := "/spaces/" + lc.spaceID + "/lists/" + listID + "/views/" + viewID + "/objects"
	req, err := lc.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	response := &anytype.ObjectListResponse{}
	err = lc.client.doRequest(req, response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (lc *ListClientImpl) AddObjects(ctx context.Context, listID string, objectIDs []string) error {
	endpoint := "/spaces/" + lc.spaceID + "/lists/" + listID + "/objects"

	requestBody := AddObjectsToListRequest{
		Objects: objectIDs,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := lc.client.newRequest(ctx, http.MethodPost, endpoint, jsonData)
	if err != nil {
		return err
	}

	return lc.client.doRequest(req, nil)
}

func (lc *ListClientImpl) RemoveObject(ctx context.Context, listID string, objectID string) error {
	endpoint := "/spaces/" + lc.spaceID + "/lists/" + listID + "/objects/" + objectID

	req, err := lc.client.newRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	return lc.client.doRequest(req, nil)
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

func (vc *ViewClientImpl) List(ctx context.Context) (*anytype.ViewListResponse, error) {
	endpoint := "/spaces/" + vc.spaceID + "/lists/" + vc.listID + "/views"
	req, err := vc.client.newRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	response := &anytype.ViewListResponse{}
	err = vc.client.doRequest(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
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

func (olc *ObjectListClientImpl) List(ctx context.Context) (*anytype.ObjectListResponse, error) {
	endpoint := "/spaces/" + olc.spaceID + "/lists/" + olc.listID + "/objects"
	req, err := olc.client.newRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	response := &anytype.ObjectListResponse{}
	err = olc.client.doRequest(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (olc *ObjectListClientImpl) Add(ctx context.Context, objectIDs []string) error {
	endpoint := "/spaces/" + olc.spaceID + "/lists/" + olc.listID + "/objects"

	requestBody := AddObjectsToListRequest{
		Objects: objectIDs,
	}

	req, err := olc.client.newRequest(ctx, "POST", endpoint, requestBody)
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
	req, err := olc.client.newRequest(ctx, "DELETE", endpoint, nil)
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

func (ovc *ObjectViewClientImpl) List(ctx context.Context) (*anytype.ObjectListResponse, error) {
	endpoint := "/spaces/" + ovc.spaceID + "/lists/" + ovc.listID + "/views/" + ovc.viewID + "/objects"
	req, err := ovc.client.newRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	response := &anytype.ObjectListResponse{}
	err = ovc.client.doRequest(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
