package anytype

import (
	"context"
)

type SearchClient interface {
	Search(ctx context.Context, request SearchRequest) (*SearchResponse, error)
}

type SearchResponse struct {
	Data []Object `json:"data"`
}

type SearchRequest struct {
	Query   string            `json:"query"`
	Types   []string          `json:"types,omitempty"`
	Sort    *SortOptions      `json:"sort,omitempty"`
	Filters *FilterExpression `json:"filters,omitempty"`
}

type SortOptions struct {
	Property  SortProperty  `json:"property"`
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

type FilterItem struct {
	Key       string          `json:"key"`
	Condition FilterCondition `json:"condition"`
	Value     interface{}     `json:"value,omitempty"`
	Values    []interface{}   `json:"values,omitempty"`
}

type FilterCondition string

const (
	FilterConditionEqual    FilterCondition = "equal"
	FilterConditionNotEqual FilterCondition = "not_equal"
	FilterConditionIn       FilterCondition = "in"
	FilterConditionNotIn    FilterCondition = "not_in"
	FilterConditionEmpty    FilterCondition = "empty"
	FilterConditionNotEmpty FilterCondition = "not_empty"

	FilterConditionContains    FilterCondition = "contains"
	FilterConditionNotContains FilterCondition = "not_contains"

	FilterConditionGreater      FilterCondition = "greater"
	FilterConditionLess         FilterCondition = "less"
	FilterConditionGreaterEqual FilterCondition = "greater_or_equal"
	FilterConditionLessEqual    FilterCondition = "less_or_equal"

	FilterConditionDaysAgo   FilterCondition = "days_ago"
	FilterConditionDaysAhead FilterCondition = "days_ahead"

	FilterConditionChecked    FilterCondition = "checked"
	FilterConditionNotChecked FilterCondition = "not_checked"
)
