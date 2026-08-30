package anytype

import (
	"context"
	"fmt"
	"iter"
)

type ObjectClient interface {
	List(ctx context.Context, opts ...ListOption) (*Page[Object], error)
	All(ctx context.Context, opts ...ListOption) iter.Seq2[Object, error]
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

// NewPropertyLinkValue builds the typed value for a property whose format is only
// known at runtime (from Property.Format). A nil value clears select and date;
// number accepts any Go integer or float.
func NewPropertyLinkValue(key string, format PropertyFormat, value any) (PropertyLinkValue, error) {
	mismatch := func(want string) error {
		return fmt.Errorf("property %q: format %s expects %s, got %T", key, format, want, value)
	}
	str := func() (string, bool) { s, ok := value.(string); return s, ok }
	strs := func() ([]string, bool) { s, ok := value.([]string); return s, ok }
	optStr := func() (*string, bool) {
		switch v := value.(type) {
		case nil:
			return nil, true
		case string:
			return &v, true
		case *string:
			return v, true
		}
		return nil, false
	}

	switch format {
	case PropertyFormatText:
		if v, ok := str(); ok {
			return TextPropertyLinkValue{Key: key, Text: v}, nil
		}
		return nil, mismatch("string")
	case PropertyFormatURL:
		if v, ok := str(); ok {
			return URLPropertyLinkValue{Key: key, URL: v}, nil
		}
		return nil, mismatch("string")
	case PropertyFormatEmail:
		if v, ok := str(); ok {
			return EmailPropertyLinkValue{Key: key, Email: v}, nil
		}
		return nil, mismatch("string")
	case PropertyFormatPhone:
		if v, ok := str(); ok {
			return PhonePropertyLinkValue{Key: key, Phone: v}, nil
		}
		return nil, mismatch("string")
	case PropertyFormatNumber:
		if v, ok := toFloat(value); ok {
			return NumberPropertyLinkValue{Key: key, Number: v}, nil
		}
		return nil, mismatch("number")
	case PropertyFormatCheckbox:
		if v, ok := value.(bool); ok {
			return CheckboxPropertyLinkValue{Key: key, Checkbox: v}, nil
		}
		return nil, mismatch("bool")
	case PropertyFormatSelect:
		if v, ok := optStr(); ok {
			return SelectPropertyLinkValue{Key: key, Select: v}, nil
		}
		return nil, mismatch("string or nil")
	case PropertyFormatDate:
		if v, ok := optStr(); ok {
			return DatePropertyLinkValue{Key: key, Date: v}, nil
		}
		return nil, mismatch("string or nil")
	case PropertyFormatMultiSelect:
		if v, ok := strs(); ok {
			return MultiSelectPropertyLinkValue{Key: key, MultiSelect: v}, nil
		}
		return nil, mismatch("[]string")
	case PropertyFormatFiles:
		if v, ok := strs(); ok {
			return FilesPropertyLinkValue{Key: key, Files: v}, nil
		}
		return nil, mismatch("[]string")
	case PropertyFormatObjects:
		if v, ok := strs(); ok {
			return ObjectsPropertyLinkValue{Key: key, Objects: v}, nil
		}
		return nil, mismatch("[]string")
	}
	return nil, fmt.Errorf("property %q: unknown format %q", key, format)
}

func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int8:
		return float64(n), true
	case int16:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint8:
		return float64(n), true
	case uint16:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	}
	return 0, false
}

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
	Color  Color      `json:"color,omitempty"`
}

func EmojiIcon(emoji string) *Icon {
	return &Icon{Format: IconFormatEmoji, Emoji: emoji}
}

func FileIcon(fileID string) *Icon {
	return &Icon{Format: IconFormatFile, File: fileID}
}

func NamedIcon(name string, color Color) *Icon {
	return &Icon{Format: IconFormatIcon, Name: name, Color: color}
}
