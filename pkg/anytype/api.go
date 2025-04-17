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

// extractTags is a helper function to extract tags from an object's Relations and Properties
func extractTags(obj *Object) {
	// Extract tags from the Relations field if present
	if tags, ok := obj.Relations["tags"]; ok {
		// Extract tags based on different possible formats
		switch v := tags.(type) {
		case []interface{}:
			// Handle array of tags
			obj.Tags = make([]string, 0, len(v))
			for _, tag := range v {
				switch t := tag.(type) {
				case string:
					obj.Tags = append(obj.Tags, t)
				case map[string]interface{}:
					// Handle case where tag might be an object with a name field
					if name, ok := t["name"].(string); ok {
						obj.Tags = append(obj.Tags, name)
					}
				}
			}
		case []string:
			// Handle string array directly
			obj.Tags = v
		case string:
			// Handle single tag as string
			obj.Tags = []string{v}
		case map[string]interface{}:
			// Handle case where the whole tags field is a single object
			if name, ok := v["name"].(string); ok {
				obj.Tags = []string{name}
			}
		}
	}

	// Make sure Tags is at least an empty slice, not nil
	if obj.Tags == nil {
		obj.Tags = []string{}
	}

	// Extract tags from properties array if present
	if len(obj.Properties) > 0 {
		for _, prop := range obj.Properties {
			if prop.Name == "Tag" && prop.Format == "multi_select" && len(prop.MultiSelect) > 0 {
				for _, tag := range prop.MultiSelect {
					if tag.Name != "" {
						obj.Tags = append(obj.Tags, tag.Name)
					}
				}
			}
		}
	}
}

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
	Query string       `json:"query,omitempty"`
	Types []string     `json:"types,omitempty"`
	Sort  *SortOptions `json:"sort,omitempty"`
	// These fields are for internal use and not part of the official API
	SpaceID string      `json:"spaceId,omitempty"`
	Tags    []string    `json:"tags,omitempty"`
	Filter  string      `json:"filter,omitempty"`
	Limit   int         `json:"limit,omitempty"`
	Offset  int         `json:"offset,omitempty"`
	Custom  interface{} `json:"custom,omitempty"`
}

// TypeResponse represents the structure of a type response
type TypeResponse struct {
	Data       []TypeInfo `json:"data"`
	Pagination struct {
		Total   int  `json:"total"`
		Offset  int  `json:"offset"`
		Limit   int  `json:"limit"`
		HasMore bool `json:"has_more"`
	} `json:"pagination"`
}

// GetSpaces retrieves spaces from the API
func (c *Client) GetSpaces(ctx context.Context) (*SpacesResponse, error) {
	data, err := c.makeRequest(ctx, http.MethodGet, "/v1/spaces", nil)
	if err != nil {
		return nil, wrapError("/v1/spaces", 0, "failed to get spaces", err)
	}

	if c.debug && c.logger != nil {
		c.logger.Debug("Raw spaces response: %s", string(data))
	}

	var response SpacesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, wrapError("/v1/spaces", 0, "failed to parse spaces response", err)
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
		return nil, wrapError("/v1/spaces/{id}", 0, "space ID is required", ErrInvalidSpaceID)
	}

	path := fmt.Sprintf("/v1/spaces/%s", spaceID)
	data, err := c.makeRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, wrapError(path, 0, fmt.Sprintf("failed to get space %s", spaceID), err)
	}

	if c.debug && c.logger != nil {
		c.logger.Debug("Raw space response: %s", string(data))
	}

	var response struct {
		Space Space `json:"space"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, wrapError(path, 0, "failed to parse space response", err)
	}

	return &response.Space, nil
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
		c.typeCache[spaceID][t.Key] = t.Name
	}

	return &response, nil
}

// GetTypeByName retrieves key for a specific type by its name
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
		c.typeCache[spaceID][t.Key] = t.Name
		reverseCache[t.Name] = t.Key

		// Also handle special case: "Page" -> "ot-page", "Note" -> "ot-note" etc.
		// This helps with common type lookups when the name is a common known type
		if strings.HasPrefix(t.Key, "ot-") && strings.EqualFold(strings.TrimPrefix(t.Key, "ot-"), typeName) {
			if c.debug && c.logger != nil {
				c.logger.Debug("Found matching type by key prefix: '%s' -> '%s'", typeName, t.Key)
			}
			reverseCache[typeName] = t.Key
		}
	}

	// Now check if we have the type after refreshing the cache
	if typeKey, found := reverseCache[typeName]; found {
		return typeKey, nil
	}

	// Try case-insensitive matching as fallback if exact match not found
	for name, key := range reverseCache {
		if strings.EqualFold(name, typeName) {
			if c.debug && c.logger != nil {
				c.logger.Debug("Found type using case-insensitive match: '%s' -> '%s'", typeName, name)
			}
			return key, nil
		}
	}

	// As a final fallback, try to construct a standard key format (e.g., "Page" -> "ot-page")
	standardKey := "ot-" + strings.ToLower(typeName)
	for _, t := range types.Data {
		if t.Key == standardKey {
			if c.debug && c.logger != nil {
				c.logger.Debug("Found type using standard key construction: '%s' -> '%s'", typeName, standardKey)
			}
			return standardKey, nil
		}
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

	// Create search request body according to API spec
	requestBody := SearchRequestBody{
		Query:   params.Query,
		Limit:   params.Limit,
		Offset:  params.Offset,
		SpaceID: spaceID,
	}

	// Filter out empty type strings before adding to request
	if len(params.Types) > 0 {
		typeKeys := make([]string, 0, len(params.Types))
		for _, t := range params.Types {
			if t != "" {
				typeKeys = append(typeKeys, t)
			}
		}
		if len(typeKeys) > 0 {
			requestBody.Types = typeKeys
		}
	}

	// Include tags if present - the API may or may not handle them natively
	// otherwise we'll filter results afterward
	if len(requestedTags) > 0 {
		requestBody.Tags = requestedTags
	}

	// Set sort options if provided
	if params.Sort != nil {
		requestBody.Sort = params.Sort
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

	// For empty responses, return an empty result
	if len(data) == 0 || string(data) == "{}" || string(data) == "[]" {
		if c.debug && c.logger != nil {
			c.logger.Debug("Empty search response, returning empty result")
		}
		return &SearchResponse{
			Data:       []Object{},
			Pagination: Pagination{Total: 0, Limit: params.Limit, Offset: params.Offset},
		}, nil
	}

	// First try to parse the response as a map to understand its structure
	var rawResponse map[string]interface{}
	if err := json.Unmarshal(data, &rawResponse); err == nil && c.debug && c.logger != nil {
		keys := make([]string, 0, len(rawResponse))
		for k := range rawResponse {
			keys = append(keys, k)
		}
		c.logger.Debug("Response structure keys: %v", keys)
	}

	// Adjust the response format - first try to unmarshal into a structure with items array
	// which is how the API returns data in newer versions
	if _, hasItems := rawResponse["items"]; hasItems {
		// New format with "items" array
		var responseItems struct {
			Items      []Object   `json:"items"`
			Total      int        `json:"total"`
			Limit      int        `json:"limit"`
			Offset     int        `json:"offset"`
			Pagination Pagination `json:"pagination,omitempty"`
		}

		if err := json.Unmarshal(data, &responseItems); err != nil {
			if c.debug && c.logger != nil {
				c.logger.Debug("Failed to unmarshal items format: %v", err)
			}
			return nil, wrapError(path, 0, "failed to unmarshal search response (items format)", err)
		}

		// Successfully parsed in the items format
		if c.debug && c.logger != nil {
			c.logger.Debug("Successfully parsed search response in items format")
		}

		// Extract tags from objects
		for i := range responseItems.Items {
			extractTags(&responseItems.Items[i])
		}

		return &SearchResponse{
			Data: responseItems.Items,
			Pagination: Pagination{
				Total:   responseItems.Total,
				Limit:   responseItems.Limit,
				Offset:  responseItems.Offset,
				HasMore: responseItems.Total > (responseItems.Offset + responseItems.Limit),
			},
		}, nil
	} else if rawObjects, hasObjects := rawResponse["objects"]; hasObjects {
		// Format with "objects" array which appears in the search by tags response
		if c.debug && c.logger != nil {
			c.logger.Debug("Found objects array in response")
		}

		var objects []Object
		rawObjectsData, err := json.Marshal(rawObjects)
		if err != nil {
			if c.debug && c.logger != nil {
				c.logger.Debug("Failed to marshal objects data: %v", err)
			}
			return nil, wrapError(path, 0, "failed to marshal objects data", err)
		}

		if err := json.Unmarshal(rawObjectsData, &objects); err != nil {
			if c.debug && c.logger != nil {
				c.logger.Debug("Failed to unmarshal objects array: %v", err)
			}
			return nil, wrapError(path, 0, "failed to unmarshal objects array", err)
		}

		// Extract pagination info
		var pagination struct {
			Total   int  `json:"total"`
			Offset  int  `json:"offset"`
			Limit   int  `json:"limit"`
			HasMore bool `json:"has_more"`
		}

		if raw, ok := rawResponse["pagination"]; ok {
			rawData, err := json.Marshal(raw)
			if err == nil {
				_ = json.Unmarshal(rawData, &pagination)
			}
		} else {
			// Default pagination if not provided
			pagination.Total = len(objects)
			pagination.Limit = params.Limit
			pagination.Offset = params.Offset
		}

		if c.debug && c.logger != nil {
			c.logger.Debug("Successfully parsed search response in objects format")
		}

		// Extract tags from objects
		for i := range objects {
			extractTags(&objects[i])
		}

		return &SearchResponse{
			Data: objects,
			Pagination: Pagination{
				Total:   pagination.Total,
				Limit:   pagination.Limit,
				Offset:  pagination.Offset,
				HasMore: pagination.HasMore,
			},
		}, nil
	} else if _, hasData := rawResponse["data"]; hasData {
		// Original format with "data" array
		var response SearchResponse
		if err := json.Unmarshal(data, &response); err != nil {
			if c.debug && c.logger != nil {
				c.logger.Debug("Failed to unmarshal data format: %v", err)
			}
			return nil, wrapError(path, 0, "failed to unmarshal search response (data format)", err)
		}

		if c.debug && c.logger != nil {
			c.logger.Debug("Successfully parsed search response in data format")
		}

		// Extract tags from the Relations field for each object
		for i := range response.Data {
			extractTags(&response.Data[i])
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

	// If we get here, the response didn't match any expected format
	if c.debug && c.logger != nil {
		c.logger.Debug("Response didn't match any expected format")
	}

	// Return empty response
	return &SearchResponse{
		Data:       []Object{},
		Pagination: Pagination{Total: 0, Limit: params.Limit, Offset: params.Offset},
	}, nil
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

	// First try to unmarshal with a wrapper structure to handle nested response format
	var objectResponse struct {
		Object Object `json:"object"`
	}
	if err := json.Unmarshal(data, &objectResponse); err == nil {
		// Successfully parsed with wrapper structure
		return &objectResponse.Object, nil
	}

	// If that didn't work, try direct unmarshaling
	var object Object
	if err := json.Unmarshal(data, &object); err != nil {
		return nil, fmt.Errorf("failed to parse object response: %w", err)
	}

	// Extract tags from the Relations field if present
	if tags, ok := object.Relations["tags"]; ok {
		// Extract tags based on different possible formats
		switch v := tags.(type) {
		case []interface{}:
			// Handle array of tags
			object.Tags = make([]string, 0, len(v))
			for _, tag := range v {
				switch t := tag.(type) {
				case string:
					object.Tags = append(object.Tags, t)
				case map[string]interface{}:
					// Handle case where tag might be an object with a name field
					if name, ok := t["name"].(string); ok {
						object.Tags = append(object.Tags, name)
					}
				}
			}
		case []string:
			// Handle string array directly
			object.Tags = v
		case string:
			// Handle single tag as string
			object.Tags = []string{v}
		case map[string]interface{}:
			// Handle case where the whole tags field is a single object
			if name, ok := v["name"].(string); ok {
				object.Tags = []string{name}
			}
		}
	}

	// Make sure Tags is at least an empty slice, not nil
	if object.Tags == nil {
		object.Tags = []string{}
	}

	// Extract tags from properties array if present
	if len(object.Properties) > 0 {
		for _, prop := range object.Properties {
			if prop.Name == "Tag" && prop.Format == "multi_select" && len(prop.MultiSelect) > 0 {
				for _, tag := range prop.MultiSelect {
					if tag.Name != "" {
						object.Tags = append(object.Tags, tag.Name)
					}
				}
			}
		}
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

	// Ensure we add tags to Relations if they're specified in the Tags field
	if len(object.Tags) > 0 {
		if object.Relations == nil {
			object.Relations = make(map[string]interface{})
		}
		object.Relations["tags"] = object.Tags
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

	// Extract tags from the Relations field and properties to ensure consistency
	if tags, ok := created.Relations["tags"]; ok {
		// Extract tags based on different possible formats
		switch v := tags.(type) {
		case []interface{}:
			// Handle array of tags
			created.Tags = make([]string, 0, len(v))
			for _, tag := range v {
				switch t := tag.(type) {
				case string:
					created.Tags = append(created.Tags, t)
				case map[string]interface{}:
					// Handle case where tag might be an object with a name field
					if name, ok := t["name"].(string); ok {
						created.Tags = append(created.Tags, name)
					}
				}
			}
		case []string:
			// Handle string array directly
			created.Tags = v
		case string:
			// Handle single tag as string
			created.Tags = []string{v}
		case map[string]interface{}:
			// Handle case where the whole tags field is a single object
			if name, ok := v["name"].(string); ok {
				created.Tags = []string{name}
			}
		}
	}

	// Make sure Tags is at least an empty slice, not nil
	if created.Tags == nil {
		created.Tags = []string{}
	}

	// Extract tags from properties array if present
	if len(created.Properties) > 0 {
		for _, prop := range created.Properties {
			if prop.Name == "Tag" && prop.Format == "multi_select" && len(prop.MultiSelect) > 0 {
				for _, tag := range prop.MultiSelect {
					if tag.Name != "" {
						created.Tags = append(created.Tags, tag.Name)
					}
				}
			}
		}
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
		return nil, wrapError("/v1/spaces/{id}/members", 0, "space ID is required", ErrInvalidSpaceID)
	}

	path := fmt.Sprintf("/v1/spaces/%s/members", spaceID)
	data, err := c.makeRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, wrapError(path, 0, fmt.Sprintf("failed to get members for space %s", spaceID), err)
	}

	if c.debug && c.logger != nil {
		c.logger.Debug("Raw members response: %s", string(data))
	}

	var response MembersResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, wrapError(path, 0, "failed to parse members response", err)
	}

	return &response, nil
}
