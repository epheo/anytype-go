package client

import (
	"context"
	"iter"
	"net/http"

	"github.com/epheo/anytype-go"
)

type SpaceClientImpl struct {
	client *ClientImpl
}

func (sc *SpaceClientImpl) Create(ctx context.Context, request anytype.CreateSpaceRequest) (*anytype.SpaceResponse, error) {
	req, err := sc.client.newRequest(ctx, http.MethodPost, "/spaces", request)
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

func (sc *SpaceClientImpl) List(ctx context.Context, opts ...anytype.ListOption) (*anytype.Page[anytype.Space], error) {
	req, err := sc.client.newRequest(ctx, http.MethodGet, withListParams("/spaces", opts), nil)
	if err != nil {
		return nil, err
	}

	response := &anytype.Page[anytype.Space]{}
	if err := sc.client.doRequest(req, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (sc *SpaceClientImpl) All(ctx context.Context, opts ...anytype.ListOption) iter.Seq2[anytype.Space, error] {
	return paginate(ctx, sc.List, opts...)
}

func (sc *SpaceClientImpl) GetByName(ctx context.Context, name string) (*anytype.Space, error) {
	for s, err := range sc.All(ctx) {
		if err != nil {
			return nil, err
		}
		if s.Name == name {
			return &s, nil
		}
	}

	return nil, anytype.ErrSpaceNotFound
}

type SpaceContextImpl struct {
	client  *ClientImpl
	spaceID string
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

func (sc *SpaceContextImpl) Update(ctx context.Context, request anytype.UpdateSpaceRequest) (*anytype.SpaceResponse, error) {
	endpoint := "/spaces/" + sc.spaceID

	req, err := sc.client.newRequest(ctx, http.MethodPatch, endpoint, request)
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

func (sc *SpaceContextImpl) Search(ctx context.Context, request anytype.SearchRequest, opts ...anytype.ListOption) (*anytype.Page[anytype.Object], error) {
	endpoint := withListParams("/spaces/"+sc.spaceID+"/search", opts)

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

func (sc *SpaceContextImpl) SearchAll(ctx context.Context, request anytype.SearchRequest, opts ...anytype.ListOption) iter.Seq2[anytype.Object, error] {
	return paginate(ctx, func(ctx context.Context, o ...anytype.ListOption) (*anytype.Page[anytype.Object], error) {
		return sc.Search(ctx, request, o...)
	}, opts...)
}

func (sc *SpaceContextImpl) Properties() anytype.SpacePropertyClient {
	return &SpacePropertyClientImpl{
		client:  sc.client,
		spaceID: sc.spaceID,
	}
}

func (sc *SpaceContextImpl) Property(propertyID string) anytype.PropertyContext {
	return &PropertyContextImpl{
		client:     sc.client,
		spaceID:    sc.spaceID,
		propertyID: propertyID,
	}
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
