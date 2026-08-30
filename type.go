package anytype

import (
	"context"
	"errors"
	"iter"
)

var ErrTypeNotFound = errors.New("type not found")

type TypeClient interface {
	List(ctx context.Context, opts ...ListOption) (*Page[Type], error)
	All(ctx context.Context, opts ...ListOption) iter.Seq2[Type, error]
	// Get finds a type by key; returns ErrTypeNotFound when no type matches.
	Get(ctx context.Context, typeKey string) (*Type, error)
	GetKeyByName(ctx context.Context, name string) (string, error)
	Create(ctx context.Context, request CreateTypeRequest) (*TypeResponse, error)
}

type TypeContext interface {
	Get(ctx context.Context) (*TypeResponse, error)
	Update(ctx context.Context, request UpdateTypeRequest) (*TypeResponse, error)
	Delete(ctx context.Context) (*TypeResponse, error)
	Templates() TemplateClient
	Template(templateID string) TemplateContext
}

type Type struct {
	ID         string     `json:"id"`
	Key        string     `json:"key"`
	Name       string     `json:"name"`
	Object     string     `json:"object"`
	Icon       *Icon      `json:"icon,omitempty"`
	Layout     string     `json:"layout"`
	Archived   bool       `json:"archived"`
	PluralName string     `json:"plural_name,omitempty"`
	Properties []Property `json:"properties,omitempty"`
}

type PropertyDefinition struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Format string `json:"format"`
}

type CreateTypeRequest struct {
	Key        string               `json:"key,omitempty"`
	Name       string               `json:"name"`
	Icon       *Icon                `json:"icon,omitempty"`
	Layout     string               `json:"layout"`
	PluralName string               `json:"plural_name"`
	Properties []PropertyDefinition `json:"properties,omitempty"`
}

type UpdateTypeRequest struct {
	Key        string               `json:"key,omitempty"`
	Name       string               `json:"name,omitempty"`
	Icon       *Icon                `json:"icon,omitempty"`
	Layout     string               `json:"layout,omitempty"`
	PluralName string               `json:"plural_name,omitempty"`
	Properties []PropertyDefinition `json:"properties,omitempty"`
}

type TypeResponse struct {
	Type Type `json:"type"`
}
