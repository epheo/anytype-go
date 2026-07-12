package anytype

import (
	"context"
	"iter"
)

type SpacePropertyClient interface {
	List(ctx context.Context, opts ...ListOption) (*Page[Property], error)
	All(ctx context.Context, opts ...ListOption) iter.Seq2[Property, error]
	Create(ctx context.Context, request CreatePropertyRequest) (*PropertyResponse, error)
}

type TagClient interface {
	List(ctx context.Context, opts ...ListOption) (*Page[Tag], error)
	All(ctx context.Context, opts ...ListOption) iter.Seq2[Tag, error]
	Create(ctx context.Context, request CreateTagRequest) (*TagResponse, error)
}

type CreatePropertyRequest struct {
	Name   string             `json:"name"`
	Format string             `json:"format"`
	Key    string             `json:"key,omitempty"`
	Tags   []CreateTagRequest `json:"tags,omitempty"`
}

type UpdatePropertyRequest struct {
	Name string `json:"name"`
	Key  string `json:"key,omitempty"`
}

type PropertyResponse struct {
	Property Property `json:"property"`
}

type CreateTagRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Key   string `json:"key,omitempty"`
}

type UpdateTagRequest struct {
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
	Key   string `json:"key,omitempty"`
}

type TagResponse struct {
	Tag Tag `json:"tag"`
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
}

type Tag struct {
	ID     string `json:"id,omitempty"`
	Key    string `json:"key,omitempty"`
	Name   string `json:"name,omitempty"`
	Color  string `json:"color,omitempty"`
	Object string `json:"object,omitempty"`
}
