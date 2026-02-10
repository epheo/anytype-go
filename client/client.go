package client

import (
	"net/http"

	"github.com/epheo/anytype-go"
	"github.com/epheo/anytype-go/middleware"
)

type ClientImpl struct {
	doer    middleware.HTTPDoer
	baseURL string
	appKey  string
}

func init() {
	// Register constructor to enable the dependency injection pattern
	anytype.RegisterClientConstructor(NewClient)
}

func NewClient(options anytype.ClientOptions) anytype.Client {
	chain := middleware.NewChain(http.DefaultClient)
	chain.Use(middleware.WithRetry())

	return &ClientImpl{
		doer:    chain.Build(),
		baseURL: options.BaseURL,
		appKey:  options.AppKey,
	}
}

func (c *ClientImpl) Spaces() anytype.SpaceClient {
	return &SpaceClientImpl{client: c}
}

func (c *ClientImpl) Space(spaceID string) anytype.SpaceContext {
	return &SpaceContextImpl{
		client:  c,
		spaceID: spaceID,
	}
}

func (c *ClientImpl) Search() anytype.SearchClient {
	return &SearchClientImpl{client: c}
}

func (c *ClientImpl) Auth() anytype.AuthClient {
	return &AuthClientImpl{client: c}
}
