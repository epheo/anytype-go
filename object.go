package anytype

import (
	"context"
)

type ObjectClient interface {
	List(ctx context.Context, opts ...ListOption) ([]Object, error)
	Create(ctx context.Context, request CreateObjectRequest) (*ObjectResponse, error)
}

type ObjectContext interface {
	Get(ctx context.Context, opts ...GetOption) (*ObjectResponse, error)
	Update(ctx context.Context, request UpdateObjectRequest) (*ObjectResponse, error)
	Delete(ctx context.Context) (*ObjectResponse, error)
}

type Object struct {
	ID         string     `json:"id"`
	Object     string     `json:"object"`
	Name       string     `json:"name"`
	SpaceID    string     `json:"space_id"`
	Layout     string     `json:"layout"`
	Archived   bool       `json:"archived"`
	Icon       *Icon      `json:"icon,omitempty"`
	Snippet    string     `json:"snippet"`
	Properties []Property `json:"properties"`
	Type       *Type      `json:"type,omitempty"`
	Markdown   string     `json:"markdown,omitempty"`
}

type ObjectResponse struct {
	Object *Object `json:"object"`
}

type CreateObjectRequest struct {
	TypeKey    string              `json:"type_key"`
	Name       string              `json:"name,omitempty"`
	Body       string              `json:"body,omitempty"`
	Icon       *Icon               `json:"icon,omitempty"`
	TemplateID string              `json:"template_id,omitempty"`
	Properties []PropertyLinkValue `json:"properties,omitempty"`
}

type UpdateObjectRequest struct {
	Name       string              `json:"name,omitempty"`
	Markdown   string              `json:"markdown,omitempty"`
	Icon       *Icon               `json:"icon,omitempty"`
	TypeKey    string              `json:"type_key,omitempty"`
	Properties []PropertyLinkValue `json:"properties,omitempty"`
}

// PropertyLinkValue is implemented by the typed property value structs below.
// The API requires each property to be sent as a distinct type (oneOf),
// so that null values can be sent explicitly to clear a property.
type PropertyLinkValue interface {
	propertyLinkValue()
}

type TextPropertyLinkValue struct {
	Key  string `json:"key"`
	Text string `json:"text"`
}

type NumberPropertyLinkValue struct {
	Key    string  `json:"key"`
	Number float64 `json:"number"`
}

type SelectPropertyLinkValue struct {
	Key    string  `json:"key"`
	Select *string `json:"select"`
}

type MultiSelectPropertyLinkValue struct {
	Key         string   `json:"key"`
	MultiSelect []string `json:"multi_select"`
}

type DatePropertyLinkValue struct {
	Key  string  `json:"key"`
	Date *string `json:"date"`
}

type FilesPropertyLinkValue struct {
	Key   string   `json:"key"`
	Files []string `json:"files"`
}

type CheckboxPropertyLinkValue struct {
	Key      string `json:"key"`
	Checkbox bool   `json:"checkbox"`
}

type URLPropertyLinkValue struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

type EmailPropertyLinkValue struct {
	Key   string `json:"key"`
	Email string `json:"email"`
}

type PhonePropertyLinkValue struct {
	Key   string `json:"key"`
	Phone string `json:"phone"`
}

type ObjectsPropertyLinkValue struct {
	Key     string   `json:"key"`
	Objects []string `json:"objects"`
}

func (TextPropertyLinkValue) propertyLinkValue()        {}
func (NumberPropertyLinkValue) propertyLinkValue()      {}
func (SelectPropertyLinkValue) propertyLinkValue()      {}
func (MultiSelectPropertyLinkValue) propertyLinkValue() {}
func (DatePropertyLinkValue) propertyLinkValue()        {}
func (FilesPropertyLinkValue) propertyLinkValue()       {}
func (CheckboxPropertyLinkValue) propertyLinkValue()    {}
func (URLPropertyLinkValue) propertyLinkValue()         {}
func (EmailPropertyLinkValue) propertyLinkValue()       {}
func (PhonePropertyLinkValue) propertyLinkValue()       {}
func (ObjectsPropertyLinkValue) propertyLinkValue()     {}

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
