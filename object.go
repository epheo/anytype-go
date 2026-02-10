package anytype

import (
	"context"

	"github.com/epheo/anytype-go/options"
)

type ObjectClient interface {
	List(ctx context.Context, opts ...options.ListOption) ([]Object, error)
	Create(ctx context.Context, request CreateObjectRequest) (*ObjectResponse, error)
}

type ObjectContext interface {
	Get(ctx context.Context) (*ObjectResponse, error)
	Delete(ctx context.Context) (*ObjectResponse, error)
	Export(ctx context.Context, format string) (*ExportResult, error)
}

type Object struct {
	ID         string
	Name       string
	SpaceID    string `json:"space_id"`
	TypeKey    string
	Layout     string
	Archived   bool
	Icon       *Icon
	Snippet    string
	Properties []Property
	Type       *Type  `json:"type,omitempty"`
	Markdown   string `json:"markdown,omitempty"`
}

type ObjectResponse struct {
	Object *Object `json:"object"`
}

type CreateObjectRequest struct {
	TypeKey    string                  `json:"type_key"`
	Name       string                  `json:"name,omitempty"`
	Body       string                  `json:"body,omitempty"`
	Icon       *Icon                   `json:"icon,omitempty"`
	TemplateID string                  `json:"template_id,omitempty"`
	Properties []PropertyLinkWithValue `json:"properties,omitempty"`
}

type UpdateObjectRequest struct {
	Name       string                  `json:"name,omitempty"`
	Markdown   string                  `json:"markdown,omitempty"`
	Icon       *Icon                   `json:"icon,omitempty"`
	TypeKey    string                  `json:"type_key,omitempty"`
	Properties []PropertyLinkWithValue `json:"properties,omitempty"`
}

type PropertyLinkWithValue struct {
	Key         string   `json:"key"`
	Text        *string  `json:"text,omitempty"`
	Number      *float64 `json:"number,omitempty"`
	Checkbox    *bool    `json:"checkbox,omitempty"`
	Date        *int64   `json:"date,omitempty"`
	URL         *string  `json:"url,omitempty"`
	Email       *string  `json:"email,omitempty"`
	Phone       *string  `json:"phone,omitempty"`
	Files       []string `json:"files,omitempty"`
	Select      *string  `json:"select,omitempty"`
	MultiSelect []string `json:"multi_select,omitempty"`
	Objects     []string `json:"objects,omitempty"`
}

type IconFormat string

const (
	IconFormatEmoji IconFormat = "emoji"
	IconFormatFile  IconFormat = "file"
	IconFormatIcon  IconFormat = "icon"
)

type Icon struct {
	Format IconFormat `json:"format,omitempty"`
	Emoji  string     `json:"emoji,omitempty"`
	File   string     `json:"file,omitempty"`
	Name   string     `json:"name,omitempty"`
	Color  string     `json:"color,omitempty"`
}
