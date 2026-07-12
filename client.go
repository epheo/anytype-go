// Package anytype provides a Go SDK for interacting with the Anytype API.
package anytype

type Pagination struct {
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"has_more"`
}

// Page is the shape every list and search endpoint returns: one slice of
// results plus the pagination cursor. All() walks every page for you.
type Page[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type ListOptions struct {
	Limit  int
	Offset int
}

type ListOption func(*ListOptions)

func WithLimit(limit int) ListOption {
	return func(opts *ListOptions) {
		opts.Limit = limit
	}
}

func WithOffset(offset int) ListOption {
	return func(opts *ListOptions) {
		opts.Offset = offset
	}
}

func ApplyListOptions(opts ...ListOption) ListOptions {
	var options ListOptions
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

type GetOptions struct {
	Format string
}

type GetOption func(*GetOptions)

func WithFormat(format string) GetOption {
	return func(opts *GetOptions) {
		opts.Format = format
	}
}

func ApplyGetOptions(opts ...GetOption) GetOptions {
	var options GetOptions
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

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
