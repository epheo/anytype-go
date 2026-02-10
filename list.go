package anytype

import (
	"context"
)

type ListClient interface {
	Add(ctx context.Context, objectIDs []string) error
}

type ListContext interface {
	Views() ViewClient
	View(viewID string) ViewContext
	Objects() ObjectListClient
	Object(objectID string) ObjectListContext
}

type ViewClient interface {
	List(ctx context.Context) (*ViewListResponse, error)
}

type ViewContext interface {
	Objects() ObjectViewClient
}

type ObjectListClient interface {
	List(ctx context.Context) (*ObjectListResponse, error)
	Add(ctx context.Context, objectIDs []string) error
}

type ObjectListContext interface {
	Remove(ctx context.Context) error
}

type ObjectViewClient interface {
	List(ctx context.Context) (*ObjectListResponse, error)
}

type ListView struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	Layout  string       `json:"layout"`
	Filters []ListFilter `json:"filters,omitempty"`
	Sorts   []ListSort   `json:"sorts,omitempty"`
}

type ListFilter struct {
	ID          string `json:"id"`
	PropertyKey string `json:"property_key"`
	Format      string `json:"format"`
	Condition   string `json:"condition"`
	Value       string `json:"value"`
}

type ListSort struct {
	ID          string `json:"id"`
	PropertyKey string `json:"property_key"`
	Format      string `json:"format"`
	SortType    string `json:"sort_type"`
}

type ViewListResponse struct {
	Data       []ListView `json:"data"`
	Pagination struct {
		Limit   int  `json:"limit"`
		Offset  int  `json:"offset"`
		Total   int  `json:"total"`
		HasMore bool `json:"has_more"`
	} `json:"pagination"`
}

type ObjectListResponse struct {
	Data       []Object `json:"data"`
	Pagination struct {
		Limit   int  `json:"limit"`
		Offset  int  `json:"offset"`
		Total   int  `json:"total"`
		HasMore bool `json:"has_more"`
	} `json:"pagination"`
}
