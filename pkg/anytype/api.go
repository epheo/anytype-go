package anytype

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Search defaults
const (
	defaultSearchLimit  = 100
	defaultSearchOffset = 0
)

// Custom error types for better error handling
var (
	ErrEmptyResponse   = fmt.Errorf("empty response from API")
	ErrInvalidResponse = fmt.Errorf("invalid response format")
	ErrMissingRequired = fmt.Errorf("missing required parameter")
)

// Error wraps API errors with additional context
type Error struct {
	StatusCode int
	Message    string
	Path       string
	Err        error
}

func (e *Error) Error() string {
	if e.StatusCode != 0 {
		return fmt.Sprintf("API error: %s (status %d) - %s", e.Path, e.StatusCode, e.Message)
	}
	return fmt.Sprintf("API error: %s - %s", e.Path, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

// wrapError creates a new Error with context
func wrapError(path string, statusCode int, message string, err error) *Error {
	return &Error{
		StatusCode: statusCode,
		Message:    message,
		Path:       path,
		Err:        err,
	}
}

// SearchRequestBody represents the structure of a search request
type SearchRequestBody struct {
	SpaceID string      `json:"spaceId"`
	Query   string      `json:"query"`
	Types   []string    `json:"types,omitempty"`
	Tags    []string    `json:"tags,omitempty"`
	Filter  string      `json:"filter,omitempty"`
	Sort    string      `json:"sort,omitempty"`
	Limit   int         `json:"limit"`
	Offset  int         `json:"offset"`
	Custom  interface{} `json:"custom,omitempty"`
}

// TypeResponse represents the structure of a type response
type TypeResponse struct {
	Data []struct {
		Type              string `json:"type"`
		ID                string `json:"id"`
		UniqueKey         string `json:"unique_key"`
		Name              string `json:"name"`
		Icon              string `json:"icon"`
		RecommendedLayout string `json:"recommended_layout"`
	} `json:"data"`
	Pagination map[string]interface{} `json:"pagination"`
}

// GetSpaces retrieves spaces from the API
func (c *Client) GetSpaces(ctx context.Context) (*SpacesResponse, error) {
	data, err := c.makeRequest(ctx, http.MethodGet, "/v1/spaces", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get spaces: %w", err)
	}

	var response SpacesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse spaces response: %w", err)
	}

	// Fetch members for each space
	for i := range response.Data {
		if c.debug && c.logger != nil {
			c.logger.Debug("Fetching members for space %s (%s)", response.Data[i].Name, response.Data[i].ID)
		}

		members, err := c.GetMembers(ctx, response.Data[i].ID)
		if err != nil {
			if c.debug && c.logger != nil {
				c.logger.Debug("Warning: failed to get members for space %s: %v", response.Data[i].ID, err)
			}
			continue
		}

		response.Data[i].Members = members.Data
	}

	return &response, nil
}

// GetSpaceByID retrieves a specific space by ID
func (c *Client) GetSpaceByID(ctx context.Context, spaceID string) (*Space, error) {
	if spaceID == "" {
		return nil, ErrInvalidSpaceID
	}

	path := fmt.Sprintf("/v1/spaces/%s", spaceID)
	data, err := c.makeRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get space %s: %w", spaceID, err)
	}

	var space Space
	if err := json.Unmarshal(data, &space); err != nil {
		return nil, fmt.Errorf("failed to parse space response: %w", err)
	}

	return &space, nil
}

// GetTypes retrieves types from a space
func (c *Client) GetTypes(ctx context.Context, spaceID string) (*TypeResponse, error) {
	if spaceID == "" {
		return nil, ErrInvalidSpaceID
	}

	path := fmt.Sprintf("/v1/spaces/%s/types", spaceID)
	data, err := c.makeRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get types for space %s: %w", spaceID, err)
	}

	var response TypeResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse types response: %w", err)
	}

	// Update the type cache with the retrieved types
	// Initialize cache for this space if needed
	if _, ok := c.typeCache[spaceID]; !ok {
		c.typeCache[spaceID] = make(map[string]string)
	}

	// Update cache with all types
	for _, t := range response.Data {
		c.typeCache[spaceID][t.UniqueKey] = t.Name
	}

	return &response, nil
}

// GetTypeByName retrieves unique_key for a specific type by its name
func (c *Client) GetTypeByName(ctx context.Context, spaceID, typeName string) (string, error) {
	if spaceID == "" {
		return "", ErrInvalidSpaceID
	}
	if typeName == "" {
		return "", ErrInvalidTypeID
	}

	// Initialize cache for this space if needed
	if _, ok := c.typeCache[spaceID]; !ok {
		c.typeCache[spaceID] = make(map[string]string)
	}

	// Check if we have a reverse lookup cache for name -> key
	reverseCache := make(map[string]string)

	// First check if we can build a reverse lookup from existing cache
	if cache, ok := c.typeCache[spaceID]; ok && len(cache) > 0 {
		// Build reverse lookup (name -> key) from existing cache (key -> name)
		for key, name := range cache {
			reverseCache[name] = key
		}

		// Check if we already have this type name in our reverse lookup
		if typeKey, found := reverseCache[typeName]; found {
			return typeKey, nil
		}
	}

	// If not in cache, fetch all types and update cache
	types, err := c.GetTypes(ctx, spaceID)
	if err != nil {
		return "", err
	}

	// Update both caches with all types
	for _, t := range types.Data {
		c.typeCache[spaceID][t.UniqueKey] = t.Name
		reverseCache[t.Name] = t.UniqueKey
	}

	// Now check if we have the type after refreshing the cache
	if typeKey, found := reverseCache[typeName]; found {
		return typeKey, nil
	}

	return "", fmt.Errorf("type '%s' not found", typeName)
}

// Search performs a search in a space with the given parameters
func (c *Client) Search(ctx context.Context, spaceID string, params *SearchParams) (*SearchResponse, error) {
	if spaceID == "" {
		return nil, wrapError("/search", 0, "space ID is required", ErrMissingRequired)
	}
	if params == nil {
		params = NewSearchParams()
	}
	if err := params.Validate(); err != nil {
		return nil, wrapError("/search", 0, "invalid search parameters", err)
	}

	path := fmt.Sprintf("/v1/spaces/%s/search", spaceID)

	// Save the tags for post-filtering
	requestedTags := params.Tags

	// Create search request body without tags filtering
	requestBody := SearchRequestBody{
		SpaceID: spaceID,
		Query:   params.Query,
		Types:   params.Types,
		Limit:   params.Limit,
		Offset:  params.Offset,
	}

	// Increase limit if we're filtering by tags to ensure we get enough matches
	if len(requestedTags) > 0 && requestBody.Limit < 1000 {
		if c.debug && c.logger != nil {
			c.logger.Debug("Increasing search limit for tag filtering: %d -> 1000", requestBody.Limit)
		}
		requestBody.Limit = 1000
	}

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

	var response SearchResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, wrapError(path, 0, "failed to unmarshal search response", err)
	}

	// Extract tags from the Relations field for each object
	for i := range response.Data {
		if tags, ok := response.Data[i].Relations["tags"]; ok {
			// Extract tags based on different possible formats
			switch v := tags.(type) {
			case []interface{}:
				// Handle array of tags
				response.Data[i].Tags = make([]string, 0, len(v))
				for _, tag := range v {
					switch t := tag.(type) {
					case string:
						response.Data[i].Tags = append(response.Data[i].Tags, t)
					case map[string]interface{}:
						// Handle case where tag might be an object with a name field
						if name, ok := t["name"].(string); ok {
							response.Data[i].Tags = append(response.Data[i].Tags, name)
						}
					}
				}
			case []string:
				// Handle string array directly
				response.Data[i].Tags = v
			case string:
				// Handle single tag as string
				response.Data[i].Tags = []string{v}
			case map[string]interface{}:
				// Handle case where the whole tags field is a single object
				if name, ok := v["name"].(string); ok {
					response.Data[i].Tags = []string{name}
				}
			}
		}

		// Make sure Tags is at least an empty slice, not nil
		if response.Data[i].Tags == nil {
			response.Data[i].Tags = []string{}
		}
	}

	// Post-process to filter by tags if needed
	if len(requestedTags) > 0 {
		if c.debug && c.logger != nil {
			c.logger.Debug("Filtering %d objects by tags: %v", len(response.Data), requestedTags)
		}

		filteredObjects := make([]Object, 0)

		// Filter objects that contain ANY of the requested tags
		for _, obj := range response.Data {
			// Check if object has any of the requested tags
			hasTag := false
			for _, requestedTag := range requestedTags {
				for _, objTag := range obj.Tags {
					if strings.EqualFold(requestedTag, objTag) {
						hasTag = true
						break
					}
				}
				if hasTag {
					break
				}
			}

			if hasTag {
				filteredObjects = append(filteredObjects, obj)
			}
		}

		// Update the response with filtered objects
		if c.debug && c.logger != nil {
			c.logger.Debug("Filtered to %d objects that match tags", len(filteredObjects))
		}

		// Update pagination info
		response.Data = filteredObjects
		response.Pagination.Total = len(filteredObjects)
	}

	return &response, nil
}

// GetObject retrieves a specific object by ID
func (c *Client) GetObject(ctx context.Context, spaceID, objectID string) (*Object, error) {
	if spaceID == "" {
		return nil, ErrInvalidSpaceID
	}
	if objectID == "" {
		return nil, ErrInvalidObjectID
	}

	path := fmt.Sprintf("/v1/spaces/%s/objects/%s", spaceID, objectID)
	data, err := c.makeRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get object %s: %w", objectID, err)
	}

	var object Object
	if err := json.Unmarshal(data, &object); err != nil {
		return nil, fmt.Errorf("failed to parse object response: %w", err)
	}

	return &object, nil
}

// CreateObject creates a new object in a space
func (c *Client) CreateObject(ctx context.Context, spaceID string, object *Object) (*Object, error) {
	if spaceID == "" {
		return nil, ErrInvalidSpaceID
	}
	if object == nil {
		return nil, ErrInvalidParameter
	}
	if err := object.Validate(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v1/spaces/%s/objects", spaceID)
	body, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal object: %w", err)
	}

	data, err := c.makeRequest(ctx, http.MethodPost, path, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create object: %w", err)
	}

	var created Object
	if err := json.Unmarshal(data, &created); err != nil {
		return nil, fmt.Errorf("failed to parse created object response: %w", err)
	}

	return &created, nil
}

// DeleteObject deletes an object from a space
func (c *Client) DeleteObject(ctx context.Context, spaceID, objectID string) error {
	if spaceID == "" {
		return ErrInvalidSpaceID
	}
	if objectID == "" {
		return ErrInvalidObjectID
	}

	path := fmt.Sprintf("/v1/spaces/%s/objects/%s", spaceID, objectID)
	_, err := c.makeRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("failed to delete object %s: %w", objectID, err)
	}

	return nil
}

// GetMembers retrieves members of a space
func (c *Client) GetMembers(ctx context.Context, spaceID string) (*MembersResponse, error) {
	if spaceID == "" {
		return nil, ErrInvalidSpaceID
	}

	path := fmt.Sprintf("/v1/spaces/%s/members", spaceID)
	data, err := c.makeRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get members for space %s: %w", spaceID, err)
	}

	if c.debug && c.logger != nil {
		c.logger.Debug("Raw members response: %s", string(data))
	}

	var response MembersResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse members response: %w", err)
	}

	return &response, nil
}
