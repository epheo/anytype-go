package anytype

import (
	"context"
	"iter"
)

type TemplateClient interface {
	List(ctx context.Context, opts ...ListOption) (*Page[Object], error)
	All(ctx context.Context, opts ...ListOption) iter.Seq2[Object, error]
}

type TemplateContext interface {
	Get(ctx context.Context) (*TemplateResponse, error)
}

type TemplateResponse struct {
	Template Object `json:"template"`
}
