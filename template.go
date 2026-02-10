package anytype

import (
	"context"
)

type TemplateClient interface {
	List(ctx context.Context) ([]Template, error)
	Get(ctx context.Context, templateID string) (*Template, error)
}

type Template struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Icon     *Icon  `json:"icon,omitempty"`
	Archived bool   `json:"archived"`
}

type TemplateContext interface {
	Get(ctx context.Context) (*TemplateResponse, error)
}

type TemplateResponse struct {
	Template Template `json:"template"`
}
