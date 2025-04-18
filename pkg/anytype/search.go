package anytype

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// validateSearchParams validates search parameters
func (c *Client) validateSearchParams(spaceID string, params *SearchParams) error {
	if spaceID == "" {
		return wrapError("/search", 0, "space ID is required", ErrMissingRequired)
	}
	if params == nil {
		params = NewSearchParams()
	}
	if err := params.Validate(); err != nil {
		return wrapError("/search", 0, "invalid search parameters", err)
	}
	return nil
}

// prepareSearchRequest prepares the search request body
func (c *Client) prepareSearchRequest(spaceID string, params *SearchParams) (*SearchRequestBody, []string) {
	// Create search request body according to API spec
	requestBody := &SearchRequestBody{
		Query:   params.Query,
		Limit:   params.Limit,
		Offset:  params.Offset,
		SpaceID: spaceID,
	}

	// Save the tags for post-filtering
	requestedTags := params.Tags

	// Add non-empty type strings
	c.addTypeFilters(requestBody, params.Types)

	// Include tags if present - the API may or may not handle them natively
	if len(requestedTags) > 0 {
		requestBody.Tags = requestedTags

		// Increase limit for tag filtering to ensure we get enough matches
		c.adjustLimitForTagFiltering(requestBody, requestedTags)
	}

	// Set sort options if provided
	if params.Sort != nil {
		requestBody.Sort = params.Sort
	}

	return requestBody, requestedTags
}

// addTypeFilters adds type filters to the search request
func (c *Client) addTypeFilters(requestBody *SearchRequestBody, types []string) {
	if len(types) == 0 {
		return
	}

	typeKeys := make([]string, 0, len(types))
	for _, t := range types {
		if t != "" {
			typeKeys = append(typeKeys, t)
		}
	}

	if len(typeKeys) > 0 {
		requestBody.Types = typeKeys
	}
}

// adjustLimitForTagFiltering increases the limit for tag filtering
func (c *Client) adjustLimitForTagFiltering(requestBody *SearchRequestBody, tags []string) {
	if len(tags) > 0 && requestBody.Limit < 1000 {
		if c.debug && c.logger != nil {
			c.logger.Debug("Increasing search limit for tag filtering: %d -> 1000", requestBody.Limit)
		}
		requestBody.Limit = 1000
	}
}

// executeSearch executes the search request and parses the response
func (c *Client) executeSearch(ctx context.Context, spaceID string, requestBody *SearchRequestBody, params *SearchParams) (*SearchResponse, error) {
	path := fmt.Sprintf("/v1/spaces/%s/search", spaceID)

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, wrapError("/search", 0, "failed to marshal search params", err)
	}

	if c.debug && c.logger != nil {
		c.logger.Debug("Search request body: %s", string(body))
	}

	data, err := c.makeRequest(ctx, http.MethodPost, path, bytes.NewBuffer(body))
	if err != nil {
		return nil, wrapError(path, 0, "failed to perform search", err)
	}

	if c.debug && c.logger != nil {
		c.logger.Debug("Raw search response: %s", string(data))
	}

	// Handle empty responses
	if c.isEmptyResponse(data) {
		if c.debug && c.logger != nil {
			c.logger.Debug("Empty search response, returning empty result")
		}
		return &SearchResponse{
			Data:       []Object{},
			Pagination: Pagination{Total: 0, Limit: params.Limit, Offset: params.Offset},
		}, nil
	}

	// Parse response
	var response SearchResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, wrapError(path, 0, "failed to parse search response", err)
	}

	return &response, nil
}

// isEmptyResponse checks if the response data is empty
func (c *Client) isEmptyResponse(data []byte) bool {
	return len(data) == 0 || string(data) == "{}" || string(data) == "[]"
}

// postProcessSearchResults processes search results (extracts tags and filters by tags)
func (c *Client) postProcessSearchResults(response *SearchResponse, requestedTags []string) {
	// Extract tags from all objects
	for i := range response.Data {
		extractTags(&response.Data[i])
		if c.debug && c.logger != nil {
			c.logger.Debug("Object '%s' has tags: %v", response.Data[i].Name, response.Data[i].Tags)
		}
	}

	// Apply tag filtering if requested
	if len(requestedTags) > 0 {
		c.filterObjectsByTags(response, requestedTags)
	}
}

// filterObjectsByTags filters objects by the requested tags
func (c *Client) filterObjectsByTags(response *SearchResponse, requestedTags []string) {
	if c.debug && c.logger != nil {
		c.logger.Debug("Filtering %d objects by tags: %v", len(response.Data), requestedTags)
	}

	filteredObjects := make([]Object, 0)

	// Filter objects that contain ANY of the requested tags
	for _, obj := range response.Data {
		if c.objectMatchesAnyTag(obj, requestedTags) {
			filteredObjects = append(filteredObjects, obj)
		}
	}

	// Update response with filtered objects and fix pagination
	if c.debug && c.logger != nil {
		c.logger.Debug("Tag filtering reduced results: %d -> %d", len(response.Data), len(filteredObjects))
	}
	response.Data = filteredObjects
	response.Pagination.Total = len(filteredObjects)
}

// objectMatchesAnyTag checks if an object matches any of the requested tags
func (c *Client) objectMatchesAnyTag(obj Object, requestedTags []string) bool {
	for _, requestedTag := range requestedTags {
		for _, objTag := range obj.Tags {
			if strings.EqualFold(requestedTag, objTag) {
				if c.debug && c.logger != nil {
					c.logger.Debug("Object '%s' matches tag '%s' with '%s'", obj.Name, requestedTag, objTag)
				}
				return true
			}
		}
	}
	return false
}

// ensureSearchPagination ensures pagination is properly set
func (c *Client) ensureSearchPagination(response *SearchResponse, params *SearchParams) {
	if response.Pagination.Limit == 0 {
		response.Pagination.Limit = params.Limit
	}
	if response.Pagination.Offset == 0 && params.Offset > 0 {
		response.Pagination.Offset = params.Offset
	}
}

// Search performs a search in a space with the given parameters.
//
// This method allows you to search for objects within a specific space based on various criteria.
// You can search by text query, filter by object types, and limit the number of results.
//
// The spaceID parameter specifies which space to search in. If params is nil, default search
// parameters will be used. The search results include objects matching the criteria and any
// related metadata.
//
// Tag filtering is performed client-side after retrieving the results from the API.
//
// Example:
//
//	// Search for notes containing "meeting"
//	params := &anytype.SearchParams{
//	    Query: "meeting",
//	    Types: []string{"ot-note"},
//	    Limit: 50,
//	}
//
//	results, err := client.Search(ctx, "space123", params)
//	if err != nil {
//	    log.Fatalf("Search failed: %v", err)
//	}
//
//	fmt.Printf("Found %d objects matching the search criteria\n", len(results.Data))
//
//	// Search with tag filtering
//	params := &anytype.SearchParams{
//	    Tags: []string{"important", "work"},
//	    Limit: 25,
//	}
//
//	results, err := client.Search(ctx, "space123", params)
func (c *Client) Search(ctx context.Context, spaceID string, params *SearchParams) (*SearchResponse, error) {
	// Validate input parameters
	if err := c.validateSearchParams(spaceID, params); err != nil {
		return nil, err
	}

	// Prepare search request
	requestBody, requestedTags := c.prepareSearchRequest(spaceID, params)

	// Execute search request and parse response
	response, err := c.executeSearch(ctx, spaceID, requestBody, params)
	if err != nil {
		return nil, err
	}

	// Post-process results (extract and filter tags)
	c.postProcessSearchResults(response, requestedTags)

	// Ensure pagination is properly set
	c.ensureSearchPagination(response, params)

	return response, nil
}
