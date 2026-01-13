package anytype

import (
	"context"
)

// TypeClient provides operations on object types within a space
type TypeClient interface {
	// List returns all available object types in the space
	List(ctx context.Context) ([]Type, error)

	// Get retrieves details of a specific type by key
	Get(ctx context.Context, typeKey string) (*Type, error)

	// GetKeyByName looks up a type key by its name
	GetKeyByName(ctx context.Context, name string) (string, error)

	// Create creates a new type in the space
	Create(ctx context.Context, request CreateTypeRequest) (*TypeResponse, error)

	// Type returns a TypeContext for a specific type in this space
	Type(typeID string) TypeContext
}

// TypeContext provides operations on a specific object type
type TypeContext interface {
	// Get retrieves details of this specific type
	Get(ctx context.Context) (*TypeResponse, error)

	// Templates returns a TemplateClient for this type
	Templates() TemplateClient

	// Template returns a TemplateContext for a specific template of this type
	Template(templateID string) TemplateContext
}

// Type represents an object type in Anytype
type Type struct {
	ID                string `json:"id"`
	Key               string
	Name              string
	Description       string
	Icon              *Icon
	Layout            string
	RecommendedLayout string `json:"recommended_layout"`
	IsArchived        bool   `json:"is_archived"`
	IsHidden          bool   `json:"is_hidden"`

	// Available property definitions for this type
	PropertyDefinitions []PropertyDefinition `json:"property_definitions"`
}

// PropertyDefinition defines a property that can be used with a type
type PropertyDefinition struct {
	Key    string `json:"key,omitempty"`
	Name   string `json:"name"`
	Format string `json:"format"`
}

// PropertyOption represents an option for select/multi-select properties
type PropertyOption struct {
	ID    string
	Value string
	Color string
}

// CreateTypeRequest represents the request payload for creating a new type
type CreateTypeRequest struct {
	Key        string               `json:"key,omitempty"`
	Name       string               `json:"name"`
	Icon       *Icon                `json:"icon,omitempty"`
	Layout     string               `json:"layout"`
	PluralName string               `json:"plural_name,omitempty"`
	Properties []PropertyDefinition `json:"properties,omitempty"`
}

// TypeResponse represents the response from a Get call on a type
type TypeResponse struct {
	Type Type `json:"type"`
}
