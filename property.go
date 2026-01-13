package anytype

import (
	"context"
)

// PropertyClient provides operations on properties of an object
type PropertyClient interface {
	// Get retrieves a specific property value
	Get(ctx context.Context, key string) (*Property, error)

	// Set sets a property value
	Set(ctx context.Context, key string, value interface{}) error

	// Delete removes a property
	Delete(ctx context.Context, key string) error

	// List returns all properties of the object
	List(ctx context.Context) ([]Property, error)
}

// Property represents an object property
type Property struct {
	ID          string   `json:"id,omitempty"`
	Format      string   `json:"format,omitempty"`
	Key         string   `json:"key,omitempty"`
	Name        string   `json:"name,omitempty"`
	Object      string   `json:"object,omitempty"` // Data model of the object, e.g., "property"
	Text        string   `json:"text,omitempty"`
	Number      float64  `json:"number,omitempty"`
	Checkbox    bool     `json:"checkbox,omitempty"`
	Date        string   `json:"date,omitempty"`
	URL         string   `json:"url,omitempty"`
	Email       string   `json:"email,omitempty"`
	Phone       string   `json:"phone,omitempty"`
	Files       []string `json:"files,omitempty"` // Changed from File to Files to match API
	Select      *Tag     `json:"select,omitempty"`
	MultiSelect []Tag    `json:"multi_select,omitempty"`
	Objects     []string `json:"objects,omitempty"` // Changed from ObjectLinks to Objects
	Required    bool     `json:"required,omitempty"`
}

// Tag represents a select or multi-select value option
type Tag struct {
	ID     string `json:"id,omitempty"`
	Key    string `json:"key,omitempty"` // Added key field from API definition
	Name   string `json:"name,omitempty"`
	Color  string `json:"color,omitempty"`
	Object string `json:"object,omitempty"` // Data model identifier
}

// Relation represents a relation within a property
type Relation struct {
	ID       string
	Type     string
	Format   string
	ObjectID string `json:"object_id"`
}

// CreatePropertyRequest contains parameters for creating a new property
type CreatePropertyRequest struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Format string `json:"format"`
}

// PropertyResponse represents the response from creating or getting a property
type PropertyResponse struct {
	Property Property `json:"property"`
}
