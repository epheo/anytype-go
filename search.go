package anytype

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
)

type SearchClient interface {
	Search(ctx context.Context, request SearchRequest, opts ...ListOption) (*Page[Object], error)
	All(ctx context.Context, request SearchRequest, opts ...ListOption) iter.Seq2[Object, error]
}

type SearchRequest struct {
	Query   string            `json:"query"`
	Types   []string          `json:"types,omitempty"`
	Sort    *SortOptions      `json:"sort,omitempty"`
	Filters *FilterExpression `json:"filters,omitempty"`
}

type SortOptions struct {
	Property  SortProperty  `json:"property_key"`
	Direction SortDirection `json:"direction"`
}

type SortProperty string

type SortDirection string

const (
	SortPropertyCreatedDate      SortProperty = "created_date"
	SortPropertyLastModifiedDate SortProperty = "last_modified_date"
	SortPropertyLastOpenedDate   SortProperty = "last_opened_date"
	SortPropertyName             SortProperty = "name"

	SortDirectionAsc  SortDirection = "asc"
	SortDirectionDesc SortDirection = "desc"
)

type FilterExpression struct {
	Operator   FilterOperator     `json:"operator,omitempty"`
	Filters    []FilterExpression `json:"filters,omitempty"`
	Conditions []FilterItem       `json:"conditions,omitempty"`
}

type FilterOperator string

const (
	FilterOperatorAnd FilterOperator = "and"
	FilterOperatorOr  FilterOperator = "or"
)

// FilterFormat is PropertyFormat under its historical name; the constants
// below are kept so existing callers compile unchanged.
type FilterFormat = PropertyFormat

const (
	FilterFormatText        FilterFormat = "text"
	FilterFormatNumber      FilterFormat = "number"
	FilterFormatSelect      FilterFormat = "select"
	FilterFormatMultiSelect FilterFormat = "multi_select"
	FilterFormatDate        FilterFormat = "date"
	FilterFormatCheckbox    FilterFormat = "checkbox"
	FilterFormatFiles       FilterFormat = "files"
	FilterFormatURL         FilterFormat = "url"
	FilterFormatEmail       FilterFormat = "email"
	FilterFormatPhone       FilterFormat = "phone"
	FilterFormatObjects     FilterFormat = "objects"
)

// FilterItem represents a single filter condition. Set Format to indicate
// which type of filter value is being used (e.g. FilterFormatText, FilterFormatCheckbox).
// For empty/nempty conditions, Format and Value can be omitted.
type FilterItem struct {
	Key       string          `json:"-"`
	Format    FilterFormat    `json:"-"`
	Condition FilterCondition `json:"-"`
	Value     interface{}     `json:"-"`
}

func (f FilterItem) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"property_key": f.Key,
		"condition":    f.Condition,
	}
	if f.Format != "" && f.Value != nil {
		m[string(f.Format)] = f.Value
	}
	return json.Marshal(m)
}

func (f *FilterItem) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if v, ok := raw["property_key"]; ok {
		if err := json.Unmarshal(v, &f.Key); err != nil {
			return err
		}
	}
	if v, ok := raw["condition"]; ok {
		if err := json.Unmarshal(v, &f.Condition); err != nil {
			return err
		}
	}

	formatFields := []FilterFormat{
		FilterFormatText, FilterFormatNumber, FilterFormatSelect,
		FilterFormatMultiSelect, FilterFormatDate, FilterFormatCheckbox,
		FilterFormatFiles, FilterFormatURL, FilterFormatEmail,
		FilterFormatPhone, FilterFormatObjects,
	}

	for _, ff := range formatFields {
		if v, ok := raw[string(ff)]; ok {
			f.Format = ff
			var val interface{}
			if err := json.Unmarshal(v, &val); err != nil {
				return fmt.Errorf("unmarshal %s: %w", ff, err)
			}
			f.Value = val
			return nil
		}
	}

	return nil
}

// Convenience constructors for common filter types

func TextFilter(key string, condition FilterCondition, text string) FilterItem {
	return FilterItem{Key: key, Format: FilterFormatText, Condition: condition, Value: text}
}

func NumberFilter(key string, condition FilterCondition, number float64) FilterItem {
	return FilterItem{Key: key, Format: FilterFormatNumber, Condition: condition, Value: number}
}

func CheckboxFilter(key string, condition FilterCondition, checked bool) FilterItem {
	return FilterItem{Key: key, Format: FilterFormatCheckbox, Condition: condition, Value: checked}
}

func SelectFilter(key string, condition FilterCondition, tagID string) FilterItem {
	return FilterItem{Key: key, Format: FilterFormatSelect, Condition: condition, Value: tagID}
}

func MultiSelectFilter(key string, condition FilterCondition, tagIDs []string) FilterItem {
	return FilterItem{Key: key, Format: FilterFormatMultiSelect, Condition: condition, Value: tagIDs}
}

func DateFilter(key string, condition FilterCondition, date string) FilterItem {
	return FilterItem{Key: key, Format: FilterFormatDate, Condition: condition, Value: date}
}

func URLFilter(key string, condition FilterCondition, url string) FilterItem {
	return FilterItem{Key: key, Format: FilterFormatURL, Condition: condition, Value: url}
}

func EmailFilter(key string, condition FilterCondition, email string) FilterItem {
	return FilterItem{Key: key, Format: FilterFormatEmail, Condition: condition, Value: email}
}

func PhoneFilter(key string, condition FilterCondition, phone string) FilterItem {
	return FilterItem{Key: key, Format: FilterFormatPhone, Condition: condition, Value: phone}
}

func FilesFilter(key string, condition FilterCondition, fileIDs []string) FilterItem {
	return FilterItem{Key: key, Format: FilterFormatFiles, Condition: condition, Value: fileIDs}
}

func ObjectsFilter(key string, condition FilterCondition, objectIDs []string) FilterItem {
	return FilterItem{Key: key, Format: FilterFormatObjects, Condition: condition, Value: objectIDs}
}

func EmptyFilter(key string, condition FilterCondition) FilterItem {
	return FilterItem{Key: key, Condition: condition}
}

type FilterCondition string

const (
	FilterConditionEq     FilterCondition = "eq"
	FilterConditionNe     FilterCondition = "ne"
	FilterConditionIn     FilterCondition = "in"
	FilterConditionNin    FilterCondition = "nin"
	FilterConditionEmpty  FilterCondition = "empty"
	FilterConditionNEmpty FilterCondition = "nempty"

	FilterConditionContains  FilterCondition = "contains"
	FilterConditionNContains FilterCondition = "ncontains"

	FilterConditionGt  FilterCondition = "gt"
	FilterConditionLt  FilterCondition = "lt"
	FilterConditionGte FilterCondition = "gte"
	FilterConditionLte FilterCondition = "lte"

	FilterConditionAll FilterCondition = "all"
)
