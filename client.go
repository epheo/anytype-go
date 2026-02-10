// Package anytype provides a Go SDK for interacting with the Anytype API.
package anytype

type ClientOptions struct {
	BaseURL string
	AppKey  string
}

type Client interface {
	Auth() AuthClient
	Spaces() SpaceClient
	Space(spaceID string) SpaceContext
	Search() SearchClient
}

type clientConstructor func(ClientOptions) Client

// Global constructor enables dependency injection pattern where the client implementation
// package registers itself via init(), allowing the interface package to remain
// implementation-agnostic while avoiding circular dependencies.
var defaultClientConstructor clientConstructor

func RegisterClientConstructor(constructor clientConstructor) {
	defaultClientConstructor = constructor
}

type ClientOption func(*ClientOptions)

func WithBaseURL(url string) ClientOption {
	return func(o *ClientOptions) {
		o.BaseURL = url
	}
}

func WithAppKey(appKey string) ClientOption {
	return func(o *ClientOptions) {
		o.AppKey = appKey
	}
}

func NewClient(opts ...ClientOption) Client {
	if defaultClientConstructor == nil {
		// Panic is intentional - this is a programming error, not a runtime error.
		// The client package must be imported to register its constructor via init().
		panic("No client constructor registered. Import the client implementation package.")
	}

	clientOpts := ClientOptions{}

	for _, opt := range opts {
		opt(&clientOpts)
	}

	client := defaultClientConstructor(clientOpts)
	return client
}
