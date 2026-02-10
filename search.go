package anytype

import (
	"context"
)

// SearchClient provides global search operations across all spaces
type SearchClient interface {
	// Search searches for objects across all spaces
	Search(ctx context.Context, request SearchRequest) (*SearchResponse, error)
}

// SearchResponse represents the response from a search operation
type SearchResponse struct {
	Data []Object `json:"data"`
}

// SearchRequest represents a search query
type SearchRequest struct {
	Query   string            `json:"query"`
	Types   []string          `json:"types,omitempty"`   // Object type keys or IDs to filter by
	Sort    *SortOptions      `json:"sort,omitempty"`    // Use pointer so it's omitted when nil
	Filters *FilterExpression `json:"filters,omitempty"` // Advanced filtering
}

// SortOptions represents sorting options for search results
type SortOptions struct {
	Property  SortProperty  `json:"property"`
	Direction SortDirection `json:"direction"`
}

// SortProperty represents the property to sort search results by
type SortProperty string

// SortDirection represents the direction to sort search results
type SortDirection string

const (
	// Sort properties
	SortPropertyCreatedDate      SortProperty = "created_date"
	SortPropertyLastModifiedDate SortProperty = "last_modified_date"
	SortPropertyLastOpenedDate   SortProperty = "last_opened_date"
	SortPropertyName             SortProperty = "name"

	// Sort directions
	SortDirectionAsc  SortDirection = "asc"
	SortDirectionDesc SortDirection = "desc"
)

// FilterExpression represents a recursive filter structure for advanced search
type FilterExpression struct {
	Operator   FilterOperator     `json:"operator,omitempty"`   // "and" or "or"
	Filters    []FilterExpression `json:"filters,omitempty"`    // Nested filter expressions
	Conditions []FilterItem       `json:"conditions,omitempty"` // Actual filter conditions
}

// FilterOperator represents logical operators for combining filters
type FilterOperator string

const (
	FilterOperatorAnd FilterOperator = "and"
	FilterOperatorOr  FilterOperator = "or"
)

// FilterItem represents a single filter condition
type FilterItem struct {
	Key       string           `json:"key"`                 // Property key to filter by
	Condition FilterCondition  `json:"condition"`           // Filter condition
	Value     interface{}      `json:"value,omitempty"`     // Value to compare against
	Values    []interface{}    `json:"values,omitempty"`    // Multiple values for "in" conditions
}

// FilterCondition represents the comparison operator for filters
type FilterCondition string

const (
	// General conditions
	FilterConditionEqual        FilterCondition = "equal"
	FilterConditionNotEqual     FilterCondition = "not_equal"
	FilterConditionIn           FilterCondition = "in"
	FilterConditionNotIn        FilterCondition = "not_in"
	FilterConditionEmpty        FilterCondition = "empty"
	FilterConditionNotEmpty     FilterCondition = "not_empty"

	// Text conditions
	FilterConditionContains     FilterCondition = "contains"
	FilterConditionNotContains  FilterCondition = "not_contains"

	// Numeric conditions
	FilterConditionGreater      FilterCondition = "greater"
	FilterConditionLess         FilterCondition = "less"
	FilterConditionGreaterEqual FilterCondition = "greater_or_equal"
	FilterConditionLessEqual    FilterCondition = "less_or_equal"

	// Date conditions
	FilterConditionDaysAgo      FilterCondition = "days_ago"
	FilterConditionDaysAhead    FilterCondition = "days_ahead"

	// Checkbox conditions
	FilterConditionChecked      FilterCondition = "checked"
	FilterConditionNotChecked   FilterCondition = "not_checked"
)
