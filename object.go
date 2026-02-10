package anytype

import (
	"context"

	"github.com/epheo/anytype-go/options"
)

// ObjectClient provides operations on objects within a space
type ObjectClient interface {
	// List returns all objects in the space
	List(ctx context.Context, opts ...options.ListOption) ([]Object, error)

	// Create creates a new object in the space
	Create(ctx context.Context, request CreateObjectRequest) (*ObjectResponse, error)
}

// ObjectContext provides operations on a specific object
type ObjectContext interface {
	// Get retrieves the object
	Get(ctx context.Context) (*ObjectResponse, error)

	// Delete deletes the object
	Delete(ctx context.Context) (*ObjectResponse, error)

	// Export exports the object in the specified format
	Export(ctx context.Context, format string) (*ExportResult, error)
}

// Object represents an Anytype object
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
	Markdown   string `json:"markdown,omitempty"` // Content in markdown format when requested with format=md
}

// ObjectResponse wraps an Object in a response according to the API specification
type ObjectResponse struct {
	Object *Object `json:"object"`
}

// CreateObjectRequest contains parameters for creating a new object
type CreateObjectRequest struct {
	TypeKey    string                 `json:"type_key"`
	Name       string                 `json:"name,omitempty"`
	Body       string                 `json:"body,omitempty"`
	Icon       *Icon                  `json:"icon,omitempty"`
	TemplateID string                 `json:"template_id,omitempty"`
	Properties []PropertyLinkWithValue `json:"properties,omitempty"`
}

// UpdateObjectRequest contains parameters for updating an object
type UpdateObjectRequest struct {
	Name       string                 `json:"name,omitempty"`
	Markdown   string                 `json:"markdown,omitempty"`
	Icon       *Icon                  `json:"icon,omitempty"`
	TypeKey    string                 `json:"type_key,omitempty"`
	Properties []PropertyLinkWithValue `json:"properties,omitempty"`
}

// PropertyLinkWithValue represents a property key-value pair for object creation/update
type PropertyLinkWithValue struct {
	Key         string      `json:"key"`
	Text        *string     `json:"text,omitempty"`
	Number      *float64    `json:"number,omitempty"`
	Checkbox    *bool       `json:"checkbox,omitempty"`
	Date        *int64      `json:"date,omitempty"`
	URL         *string     `json:"url,omitempty"`
	Email       *string     `json:"email,omitempty"`
	Phone       *string     `json:"phone,omitempty"`
	Files       []string    `json:"files,omitempty"`
	Select      *string     `json:"select,omitempty"`
	MultiSelect []string    `json:"multi_select,omitempty"`
	Objects     []string    `json:"objects,omitempty"`
}

// IconFormat represents the type of icon
type IconFormat string

const (
	// IconFormatEmoji represents an emoji icon
	IconFormatEmoji IconFormat = "emoji"
	// IconFormatFile represents a file icon
	IconFormatFile IconFormat = "file"
	// IconFormatIcon represents a named icon
	IconFormatIcon IconFormat = "icon"
)

// Icon represents an object icon
type Icon struct {
	Format IconFormat `json:"format,omitempty"` // Type of icon: emoji, file, or icon
	Emoji  string     `json:"emoji,omitempty"`  // The emoji character if format is emoji
	File   string     `json:"file,omitempty"`   // The file URL if format is file
	Name   string     `json:"name,omitempty"`   // The name of the icon if format is icon
	Color  string     `json:"color,omitempty"`  // The color of the icon
}
