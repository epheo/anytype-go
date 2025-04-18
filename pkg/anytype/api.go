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

// extractTags is a helper function to extract tags from an object's Relations and Properties
func extractTags(obj *Object) {
	// Initialize Tags as an empty slice if nil
	if obj.Tags == nil {
		obj.Tags = []string{}
	}

	// Extract tags from the Relations field if present
	if obj.Relations != nil && obj.Relations.Items != nil {
		// Check if there are any relations of type "tags"
		if tagRelations, ok := obj.Relations.Items["tags"]; ok && len(tagRelations) > 0 {
			// Extract name from each relation
			for _, relation := range tagRelations {
				if relation.Name != "" {
					obj.Tags = append(obj.Tags, relation.Name)
				}
			}
		}
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

// GetSpaces retrieves all available spaces from the Anytype API.
//
// This method fetches all spaces that the authenticated user has access to.
// For each space, it also attempts to fetch and populate the space's members.
// If member fetching fails for any space, the error is logged (in debug mode) but
// the space is still included in the results.
//
// Example:
//
//	spaces, err := client.GetSpaces(ctx)
//	if err != nil {
//	    log.Fatalf("Failed to get spaces: %v", err)
//	}
//
//	fmt.Printf("Found %d spaces:\n", len(spaces.Data))
//	for _, space := range spaces.Data {
//	    fmt.Printf("- %s (ID: %s)\n", space.Name, space.ID)
//	}
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
func (c *Client) GetTypes(ctx context.Context, params *GetTypesParams) (*TypeResponse, error) {
	if params == nil {
		return nil, ErrInvalidParameter
	}
	if err := params.Validate(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v1/spaces/%s/types", params.SpaceID)
	data, err := c.makeRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get types for space %s: %w", params.SpaceID, err)
	}

	// This follows the API's pagination response format for types
	var response TypeResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse types response: %w", err)
	}

	// Update the type cache with the retrieved types
	// Initialize cache for this space if needed
	if _, ok := c.typeCache[params.SpaceID]; !ok {
		c.typeCache[params.SpaceID] = make(map[string]string)
	}

	// Update cache with all types
	for _, t := range response.Data {
		c.typeCache[params.SpaceID][t.Key] = t.Name
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
	types, err := c.GetTypes(ctx, &GetTypesParams{SpaceID: spaceID})
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

	// The API returns search responses in a paginated format
	// which has a data array of objects and a pagination object
	var response SearchResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, wrapError(path, 0, "failed to parse search response", err)
	}

	// Extract tags from all objects
	for i := range response.Data {
		extractTags(&response.Data[i])
		if c.debug && c.logger != nil {
			c.logger.Debug("Object '%s' has tags: %v", response.Data[i].Name, response.Data[i].Tags)
		}
	}

	// Apply tag filtering if requested
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
						if c.debug && c.logger != nil {
							c.logger.Debug("Object '%s' matches tag '%s' with '%s'", obj.Name, requestedTag, objTag)
						}
						break
					}
				}
				if hasTag {
					break
				}
			}

			if hasTag {
				filteredObjects = append(filteredObjects, obj)
			} else if c.debug && c.logger != nil {
				c.logger.Debug("Object '%s' was filtered out, tags: %v", obj.Name, obj.Tags)
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

// GetObject retrieves a specific object by ID.
//
// This method fetches an object from Anytype by its unique ID. The object includes
// metadata such as name, type, icon, tags, and other properties defined in the Object struct.
//
// The method requires a valid context and GetObjectParams struct containing the space ID
// and object ID. If successful, it returns a populated Object struct and nil error.
//
// If the object doesn't exist or cannot be fetched due to permissions or network issues,
// an appropriate error will be returned.
//
// Example:
//
//	params := &anytype.GetObjectParams{
//	    SpaceID:  "space123",
//	    ObjectID: "obj456",
//	}
//
//	object, err := client.GetObject(ctx, params)
//	if err != nil {
//	    log.Fatalf("Failed to get object: %v", err)
//	}
//
//	fmt.Printf("Object name: %s\n", object.Name)
func (c *Client) GetObject(ctx context.Context, params *GetObjectParams) (*Object, error) {
	if params == nil {
		return nil, ErrInvalidParameter
	}
	if err := params.Validate(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v1/spaces/%s/objects/%s", params.SpaceID, params.ObjectID)
	data, err := c.makeRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get object %s: %w", params.ObjectID, err)
	}

	// The API response is structured with an "object" field
	// containing the Object data in the standard response format
	var objectResponse struct {
		Object Object `json:"object"`
	}

	if err := json.Unmarshal(data, &objectResponse); err != nil {
		return nil, fmt.Errorf("failed to parse object response: %w", err)
	}

	// Extract tags from the retrieved object
	extractTags(&objectResponse.Object)

	return &objectResponse.Object, nil
}

// CreateObject creates a new object in a space.
//
// This method allows you to create new objects in a specified Anytype space.
// The object parameter must contain all required fields for the object type,
// such as name, type key, and any properties specific to that type.
//
// If the object contains tags in the Tags field, they will automatically be
// added to the object's Relations.
//
// Example:
//
//	// Create a new note
//	newObject := &anytype.Object{
//	    Name:    "Meeting Notes",
//	    TypeKey: "ot-note",
//	    Icon: &anytype.Icon{
//	        Format: "emoji",
//	        Emoji:  "📝",
//	    },
//	    Tags: []string{"work", "meeting"},
//	}
//
//	created, err := client.CreateObject(ctx, "space123", newObject)
//	if err != nil {
//	    log.Fatalf("Failed to create object: %v", err)
//	}
//
//	fmt.Printf("Created object with ID: %s\n", created.ID)
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
			object.Relations = &Relations{
				Items: make(map[string][]Relation),
			}
		}

		// Create tag relations from tag names
		tagRelations := make([]Relation, 0, len(object.Tags))
		for _, tagName := range object.Tags {
			tagRelations = append(tagRelations, Relation{
				Name: tagName,
			})
		}

		// Add to the relations map
		object.Relations.Items["tags"] = tagRelations
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

	// The API response is structured with an "object" field
	// containing the Object data in the standard response format
	var objectResponse struct {
		Object Object `json:"object"`
	}

	if err := json.Unmarshal(data, &objectResponse); err != nil {
		return nil, fmt.Errorf("failed to parse created object response: %w", err)
	}

	// Extract tags using the helper function
	extractTags(&objectResponse.Object)

	return &objectResponse.Object, nil
}

// DeleteObject deletes an object from a space.
//
// This method permanently removes an object identified by its objectID from the
// specified space. Once deleted, the object cannot be recovered through the API.
//
// Example:
//
//	err := client.DeleteObject(ctx, "space123", "obj456")
//	if err != nil {
//	    log.Fatalf("Failed to delete object: %v", err)
//	}
//
//	fmt.Println("Object successfully deleted")
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

// UpdateObject updates an existing object in a space.
//
// This method allows you to update an existing object identified by its objectID
// within the specified space. The provided object parameter should contain the
// fields to be updated. The object ID in the object parameter will be overridden
// with the objectID parameter to ensure consistency.
//
// If the object contains tags in the Tags field, they will automatically be
// added to the object's Relations, replacing any existing tag relations.
//
// Example:
//
//	// Update an existing object
//	updateObj := &anytype.Object{
//	    Name: "Updated Meeting Notes",
//	    Icon: &anytype.Icon{
//	        Format: "emoji",
//	        Emoji:  "📌",
//	    },
//	    Tags: []string{"work", "important", "meeting"},
//	}
//
//	updated, err := client.UpdateObject(ctx, "space123", "obj456", updateObj)
//	if err != nil {
//	    log.Fatalf("Failed to update object: %v", err)
//	}
//
//	fmt.Printf("Updated object: %s\n", updated.Name)
func (c *Client) UpdateObject(ctx context.Context, spaceID, objectID string, object *Object) (*Object, error) {
	if spaceID == "" {
		return nil, ErrInvalidSpaceID
	}
	if objectID == "" {
		return nil, ErrInvalidObjectID
	}
	if object == nil {
		return nil, ErrInvalidParameter
	}

	// Ensure the object ID in the path matches the object ID in the body
	object.ID = objectID

	// Ensure we add tags to Relations if they're specified in the Tags field
	if len(object.Tags) > 0 {
		if object.Relations == nil {
			object.Relations = &Relations{
				Items: make(map[string][]Relation),
			}
		}

		// Create tag relations from tag names
		tagRelations := make([]Relation, 0, len(object.Tags))
		for _, tagName := range object.Tags {
			tagRelations = append(tagRelations, Relation{
				Name: tagName,
			})
		}

		// Add to the relations map
		object.Relations.Items["tags"] = tagRelations
	}

	path := fmt.Sprintf("/v1/spaces/%s/objects/%s", spaceID, objectID)
	body, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal object: %w", err)
	}

	data, err := c.makeRequest(ctx, http.MethodPut, path, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to update object %s: %w", objectID, err)
	}

	// The API response is structured with an "object" field
	// containing the Object data in the standard response format
	var objectResponse struct {
		Object Object `json:"object"`
	}

	if err := json.Unmarshal(data, &objectResponse); err != nil {
		return nil, fmt.Errorf("failed to parse updated object response: %w", err)
	}

	// Extract tags using the helper function
	extractTags(&objectResponse.Object)

	return &objectResponse.Object, nil
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
