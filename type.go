package anytype

import (
	"context"
)

type TypeClient interface {
	List(ctx context.Context) ([]Type, error)
	Get(ctx context.Context, typeKey string) (*Type, error)
	GetKeyByName(ctx context.Context, name string) (string, error)
	Create(ctx context.Context, request CreateTypeRequest) (*TypeResponse, error)
	Type(typeID string) TypeContext
}

type TypeContext interface {
	Get(ctx context.Context) (*TypeResponse, error)
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
	Key    string `json:"key,omitempty"`
	Name   string `json:"name"`
	Format string `json:"format"`
}

type PropertyOption struct {
	ID    string
	Value string
	Color string
}

type CreateTypeRequest struct {
	Key        string               `json:"key,omitempty"`
	Name       string               `json:"name"`
	Icon       *Icon                `json:"icon,omitempty"`
	Layout     string               `json:"layout"`
	PluralName string               `json:"plural_name,omitempty"`
	Properties []PropertyDefinition `json:"properties,omitempty"`
}

type TypeResponse struct {
	Type Type `json:"type"`
}
