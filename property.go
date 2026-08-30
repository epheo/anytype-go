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

// PropertyFormat names the value shape of a property; it doubles as the JSON
// field name the API uses for that value (see NewPropertyLinkValue and FilterItem).
type PropertyFormat string

const (
	PropertyFormatText        PropertyFormat = "text"
	PropertyFormatNumber      PropertyFormat = "number"
	PropertyFormatSelect      PropertyFormat = "select"
	PropertyFormatMultiSelect PropertyFormat = "multi_select"
	PropertyFormatDate        PropertyFormat = "date"
	PropertyFormatFiles       PropertyFormat = "files"
	PropertyFormatCheckbox    PropertyFormat = "checkbox"
	PropertyFormatURL         PropertyFormat = "url"
	PropertyFormatEmail       PropertyFormat = "email"
	PropertyFormatPhone       PropertyFormat = "phone"
	PropertyFormatObjects     PropertyFormat = "objects"
)

// Color is the palette accepted for tags and named icons.
type Color string

const (
	ColorGrey   Color = "grey"
	ColorYellow Color = "yellow"
	ColorOrange Color = "orange"
	ColorRed    Color = "red"
	ColorPink   Color = "pink"
	ColorPurple Color = "purple"
	ColorBlue   Color = "blue"
	ColorIce    Color = "ice"
	ColorTeal   Color = "teal"
	ColorLime   Color = "lime"
)

type CreatePropertyRequest struct {
	Name   string             `json:"name"`
	Format PropertyFormat     `json:"format"`
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
	Color Color  `json:"color"`
	Key   string `json:"key,omitempty"`
}

type UpdateTagRequest struct {
	Name  string `json:"name,omitempty"`
	Color Color  `json:"color,omitempty"`
	Key   string `json:"key,omitempty"`
}

type TagResponse struct {
	Tag Tag `json:"tag"`
}

type Property struct {
	ID          string         `json:"id,omitempty"`
	Format      PropertyFormat `json:"format,omitempty"`
	Key         string         `json:"key,omitempty"`
	Name        string         `json:"name,omitempty"`
	Object      string         `json:"object,omitempty"`
	Text        string         `json:"text,omitempty"`
	Number      float64        `json:"number,omitempty"`
	Checkbox    bool           `json:"checkbox,omitempty"`
	Date        string         `json:"date,omitempty"`
	URL         string         `json:"url,omitempty"`
	Email       string         `json:"email,omitempty"`
	Phone       string         `json:"phone,omitempty"`
	Files       []string       `json:"files,omitempty"`
	Select      *Tag           `json:"select,omitempty"`
	MultiSelect []Tag          `json:"multi_select,omitempty"`
	Objects     []string       `json:"objects,omitempty"`
}

type Tag struct {
	ID     string `json:"id,omitempty"`
	Key    string `json:"key,omitempty"`
	Name   string `json:"name,omitempty"`
	Color  Color  `json:"color,omitempty"`
	Object string `json:"object,omitempty"`
}
