package anytype

import (
	"context"
)

type TemplateClient interface {
	List(ctx context.Context) ([]Object, error)
}

type TemplateContext interface {
	Get(ctx context.Context) (*TemplateResponse, error)
}

type TemplateResponse struct {
	Template Object `json:"template"`
}
