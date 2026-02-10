package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/epheo/anytype-go"
)

type SpaceClientImpl struct {
	client *ClientImpl
}

func (sc *SpaceClientImpl) Create(ctx context.Context, request anytype.CreateSpaceRequest) (*anytype.CreateSpaceResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := sc.client.newRequest(ctx, http.MethodPost, "/spaces", jsonData)
	if err != nil {
		return nil, err
	}

	response := &anytype.CreateSpaceResponse{}
	err = sc.client.doRequest(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (sc *SpaceClientImpl) List(ctx context.Context) (*anytype.SpaceListResponse, error) {
	req, err := sc.client.newRequest(ctx, http.MethodGet, "/spaces", nil)
	if err != nil {
		return nil, err
	}

	response := &anytype.SpaceListResponse{}
	err = sc.client.doRequest(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type SpaceContextImpl struct {
	client  *ClientImpl
	spaceID string
}

func (sc *SpaceContextImpl) Lists() anytype.ListClient {
	return &ListClientImpl{
		client:  sc.client,
		spaceID: sc.spaceID,
	}
}

func (sc *SpaceContextImpl) Get(ctx context.Context) (*anytype.SpaceResponse, error) {
	endpoint := "/spaces/" + sc.spaceID
	req, err := sc.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	response := &anytype.SpaceResponse{}
	err = sc.client.doRequest(req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (sc *SpaceContextImpl) Objects() anytype.ObjectClient {
	return &ObjectClientImpl{
		client:  sc.client,
		spaceID: sc.spaceID,
	}
}

func (sc *SpaceContextImpl) Object(objectID string) anytype.ObjectContext {
	return &ObjectContextImpl{
		client:   sc.client,
		spaceID:  sc.spaceID,
		objectID: objectID,
	}
}

func (sc *SpaceContextImpl) List(listID string) anytype.ListContext {
	return &ListContextImpl{
		client:  sc.client,
		spaceID: sc.spaceID,
		listID:  listID,
	}
}

func (sc *SpaceContextImpl) Types() anytype.TypeClient {
	return &TypeClientImpl{
		client:  sc.client,
		spaceID: sc.spaceID,
	}
}

func (sc *SpaceContextImpl) Type(typeID string) anytype.TypeContext {
	return &TypeContextImpl{
		client:  sc.client,
		spaceID: sc.spaceID,
		typeID:  typeID,
	}
}

func (sc *SpaceContextImpl) Search(ctx context.Context, request anytype.SearchRequest) (*anytype.SearchResponse, error) {
	endpoint := "/spaces/" + sc.spaceID + "/search"

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

func (sc *SpaceContextImpl) Members() anytype.MemberClient {
	return &MemberClientImpl{
		client:  sc.client,
		spaceID: sc.spaceID,
	}
}

func (sc *SpaceContextImpl) Member(memberID string) anytype.MemberContext {
	return &MemberContextImpl{
		client:   sc.client,
		spaceID:  sc.spaceID,
		memberID: memberID,
	}
}
