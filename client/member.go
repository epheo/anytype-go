package client

import (
	"context"
	"net/http"

	"github.com/epheo/anytype-go"
	"github.com/epheo/anytype-go/options"
)

type MemberClientImpl struct {
	client  *ClientImpl
	spaceID string
}

func (mc *MemberClientImpl) List(ctx context.Context) (*anytype.MemberListResponse, error) {
	path := "/spaces/" + mc.spaceID + "/members"

	var response struct {
		Data       []anytype.Member           `json:"data"`
		Pagination options.PaginationMetadata `json:"pagination"`
	}

	req, err := mc.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	if err := mc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &anytype.MemberListResponse{
		Data: response.Data,
	}, nil
}

type MemberContextImpl struct {
	client   *ClientImpl
	spaceID  string
	memberID string
}

func (mc *MemberContextImpl) Get(ctx context.Context) (*anytype.MemberResponse, error) {
	path := "/spaces/" + mc.spaceID + "/members/" + mc.memberID

	var response anytype.MemberResponse

	req, err := mc.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	if err := mc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
