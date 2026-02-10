package anytype

import (
	"context"
)

type PropertyClient interface {
	Get(ctx context.Context, key string) (*Property, error)
	Set(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	List(ctx context.Context) ([]Property, error)
}

type Property struct {
	ID          string   `json:"id,omitempty"`
	Format      string   `json:"format,omitempty"`
	Key         string   `json:"key,omitempty"`
	Name        string   `json:"name,omitempty"`
	Object      string   `json:"object,omitempty"`
	Text        string   `json:"text,omitempty"`
	Number      float64  `json:"number,omitempty"`
	Checkbox    bool     `json:"checkbox,omitempty"`
	Date        string   `json:"date,omitempty"`
	URL         string   `json:"url,omitempty"`
	Email       string   `json:"email,omitempty"`
	Phone       string   `json:"phone,omitempty"`
	Files       []string `json:"files,omitempty"`
	Select      *Tag     `json:"select,omitempty"`
	MultiSelect []Tag    `json:"multi_select,omitempty"`
	Objects     []string `json:"objects,omitempty"`
	Required    bool     `json:"required,omitempty"`
}

type Tag struct {
	ID     string `json:"id,omitempty"`
	Key    string `json:"key,omitempty"`
	Name   string `json:"name,omitempty"`
	Color  string `json:"color,omitempty"`
	Object string `json:"object,omitempty"`
}

type Relation struct {
	ID       string
	Type     string
	Format   string
	ObjectID string `json:"object_id"`
}
