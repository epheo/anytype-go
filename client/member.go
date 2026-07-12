package client

import (
	"context"
	"iter"
	"net/http"

	"github.com/epheo/anytype-go"
)

type MemberClientImpl struct {
	client  *ClientImpl
	spaceID string
}

func (mc *MemberClientImpl) List(ctx context.Context, opts ...anytype.ListOption) (*anytype.Page[anytype.Member], error) {
	endpoint := withListParams("/spaces/"+mc.spaceID+"/members", opts)

	req, err := mc.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response anytype.Page[anytype.Member]
	if err := mc.client.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (mc *MemberClientImpl) All(ctx context.Context, opts ...anytype.ListOption) iter.Seq2[anytype.Member, error] {
	return paginate(ctx, mc.List, opts...)
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
