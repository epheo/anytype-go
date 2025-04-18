package anytype

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// QueryBuilder provides a fluent interface for building and executing search queries
type QueryBuilder struct {
	client  *Client
	spaceID string
	params  *SearchParams
	timeout time.Duration
	err     error
}

// NewQueryBuilder creates a new query builder for the given space
func (c *Client) NewQueryBuilder(spaceID string) *QueryBuilder {
	return &QueryBuilder{
		client:  c,
		spaceID: spaceID,
		params:  NewSearchParams(),
	}
}

// WithQuery adds a text search term to the query
func (qb *QueryBuilder) WithQuery(query string) *QueryBuilder {
	if qb.err != nil {
		return qb
	}
	qb.params.Query = strings.TrimSpace(query)
	return qb
}

// WithType adds a single type filter to the query.
// This resolves the type name to its key automatically
func (qb *QueryBuilder) WithType(ctx context.Context, typeName string) *QueryBuilder {
	if qb.err != nil {
		return qb
	}

	if typeName == "" {
		return qb
	}

	typeKey, err := qb.client.GetTypeByName(ctx, qb.spaceID, typeName)
	if err != nil {
		qb.err = fmt.Errorf("failed to resolve type name '%s': %w", typeName, err)
		return qb
	}

	qb.params.Types = append(qb.params.Types, typeKey)
	return qb
}

// WithTypes adds multiple type filters to the query.
// This resolves all type names to their keys automatically
func (qb *QueryBuilder) WithTypes(ctx context.Context, typeNames ...string) *QueryBuilder {
	if qb.err != nil {
		return qb
	}

	for _, typeName := range typeNames {
		if typeName != "" {
			qb = qb.WithType(ctx, typeName)
			if qb.err != nil {
				return qb
			}
		}
	}

	return qb
}

// WithTypeKeys adds type filters directly using internal type keys.
// Use this when you already know the type keys (e.g., "ot-page", "ot-note")
func (qb *QueryBuilder) WithTypeKeys(typeKeys ...string) *QueryBuilder {
	if qb.err != nil {
		return qb
	}

	for _, typeKey := range typeKeys {
		if typeKey != "" {
			qb.params.Types = append(qb.params.Types, typeKey)
		}
	}

	return qb
}

// WithTag adds a single tag filter
func (qb *QueryBuilder) WithTag(tag string) *QueryBuilder {
	if qb.err != nil {
		return qb
	}

	tag = strings.TrimSpace(tag)
	if tag != "" {
		qb.params.Tags = append(qb.params.Tags, tag)
	}

	return qb
}

// WithTags adds multiple tag filters
func (qb *QueryBuilder) WithTags(tags ...string) *QueryBuilder {
	if qb.err != nil {
		return qb
	}

	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			qb.params.Tags = append(qb.params.Tags, tag)
		}
	}

	return qb
}

// WithLimit sets the maximum number of results to return
func (qb *QueryBuilder) WithLimit(limit int) *QueryBuilder {
	if qb.err != nil {
		return qb
	}

	if limit <= 0 {
		qb.err = fmt.Errorf("limit must be greater than 0")
		return qb
	}

	qb.params.Limit = limit
	return qb
}

// WithOffset sets the pagination offset
func (qb *QueryBuilder) WithOffset(offset int) *QueryBuilder {
	if qb.err != nil {
		return qb
	}

	if offset < 0 {
		qb.err = fmt.Errorf("offset cannot be negative")
		return qb
	}

	qb.params.Offset = offset
	return qb
}

// WithSortField adds sorting by a specific field
func (qb *QueryBuilder) WithSortField(field string, ascending bool) *QueryBuilder {
	if qb.err != nil {
		return qb
	}

	field = strings.TrimSpace(field)
	if field == "" {
		qb.err = fmt.Errorf("sort field cannot be empty")
		return qb
	}

	direction := "desc"
	if ascending {
		direction = "asc"
	}

	qb.params.Sort = &SortOptions{
		Property:  field,
		Direction: direction,
	}

	return qb
}

// WithSortByName sets sorting by object name
func (qb *QueryBuilder) WithSortByName(ascending bool) *QueryBuilder {
	return qb.WithSortField("name", ascending)
}

// WithSortByCreatedAt sets sorting by creation date
func (qb *QueryBuilder) WithSortByCreatedAt(ascending bool) *QueryBuilder {
	return qb.WithSortField("createdAt", ascending)
}

// WithSortByUpdatedAt sets sorting by last update date
func (qb *QueryBuilder) WithSortByUpdatedAt(ascending bool) *QueryBuilder {
	return qb.WithSortField("updatedAt", ascending)
}

// WithTimeout sets a custom timeout for this specific search operation
func (qb *QueryBuilder) WithTimeout(timeout time.Duration) *QueryBuilder {
	if qb.err != nil {
		return qb
	}

	if timeout <= 0 {
		qb.err = fmt.Errorf("timeout must be greater than 0")
		return qb
	}

	// Store the timeout in the QueryBuilder struct
	qb.timeout = timeout
	return qb
}

// GetParams returns a copy of the current search parameters
// This can be useful for saving a query for later use
func (qb *QueryBuilder) GetParams() (*SearchParams, error) {
	if qb.err != nil {
		return nil, qb.err
	}

	// Return a copy of the search parameters
	paramsCopy := *qb.params
	return &paramsCopy, nil
}

// Error returns any error that occurred during query building
func (qb *QueryBuilder) Error() error {
	return qb.err
}

// Execute runs the search with the built parameters
func (qb *QueryBuilder) Execute(ctx context.Context) (*SearchResponse, error) {
	if qb.err != nil {
		return nil, qb.err
	}

	// Apply custom timeout if specified
	if qb.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, qb.timeout)
		defer cancel()
	}

	// Call the search method with the built parameters
	return qb.client.Search(ctx, qb.spaceID, qb.params)
}

// ExecuteWithCallback runs the search and processes results with a callback function
func (qb *QueryBuilder) ExecuteWithCallback(ctx context.Context, callback func(obj Object) error) error {
	if qb.err != nil {
		return qb.err
	}

	// Apply custom timeout if specified
	if qb.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, qb.timeout)
		defer cancel()
	}

	// Execute the search
	results, err := qb.client.Search(ctx, qb.spaceID, qb.params)
	if err != nil {
		return err
	}

	// Process each result with the callback
	for _, obj := range results.Data {
		if err := callback(obj); err != nil {
			return err
		}
	}

	return nil
}
