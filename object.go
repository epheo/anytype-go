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

	// Update updates the object
	Update(ctx context.Context, request UpdateObjectRequest) (*ObjectResponse, error)

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
	TypeKey    string `json:"type_key"`
	Name       string
	Body       string
	Icon       *Icon
	TemplateID string           `json:"template_id,omitempty"`
	Properties []map[string]any `json:"properties"`
}

// UpdateObjectRequest contains parameters for updating an object
type UpdateObjectRequest struct {
	Name        *string
	Description *string
	Icon        *Icon

	// This items of this slice should be any of the *PropertyLinkValue types.
	Properties []any
}

// TextPropertyLinkValue contains parameters for updating a Text property
type TextPropertyLinkValue struct {
	Key  string `json:"key"`
	Text string `json:"text"`
}

// NumberPropertyLinkValue contains parameters for updating a Number property
type NumberPropertyLinkValue struct {
	Key    string  `json:"key"`
	Number float64 `json:"number"`
}

// SelectPropertyLinkValue contains parameters for updating a Select property
type SelectPropertyLinkValue struct {
	Key    string  `json:"key"`
	Select *string `json:"select"`
}

// MultiSelectPropertyLinkValue contains parameters for updating a MultiSelect property
type MultiSelectPropertyLinkValue struct {
	Key         string   `json:"key"`
	MultiSelect []string `json:"multi_select"`
}

// DatePropertyLinkValue contains parameters for updating a Date property
type DatePropertyLinkValue struct {
	Key  string  `json:"key"`
	Date *string `json:"date"`
}

// FilesPropertyLinkValue contains parameters for updating a Files property
type FilesPropertyLinkValue struct {
	Key   string   `json:"key"`
	Files []string `json:"files"`
}

// CheckboxPropertyLinkValue contains parameters for updating a Checkbox property
type CheckboxPropertyLinkValue struct {
	Key      string `json:"key"`
	Checkbox bool   `json:"checkbox"`
}

// UrlPropertyLinkValue contains parameters for updating a Url property
type UrlPropertyLinkValue struct {
	Key string `json:"key"`
	Url string `json:"url"`
}

// EmailPropertyLinkValue contains parameters for updating a Email property
type EmailPropertyLinkValue struct {
	Key   string `json:"key"`
	Email string `json:"email"`
}

// PhonePropertyLinkValue contains parameters for updating a Phone property
type PhonePropertyLinkValue struct {
	Key   string `json:"key"`
	Phone string `json:"phone"`
}

// ObjectsPropertyLinkValue contains parameters for updating a Objects property
type ObjectsPropertyLinkValue struct {
	Key     string   `json:"key"`
	Objects []string `json:"objects"`
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
