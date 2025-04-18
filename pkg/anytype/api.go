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

	extractTagsFromRelations(obj)
	extractTagsFromProperties(obj)
}

// extractTagsFromRelations extracts tags from object's Relations field
func extractTagsFromRelations(obj *Object) {
	if obj.Relations == nil || obj.Relations.Items == nil {
		return
	}

	tagRelations, ok := obj.Relations.Items["tags"]
	if !ok || len(tagRelations) == 0 {
		return
	}

	// Extract name from each relation
	for _, relation := range tagRelations {
		if relation.Name != "" {
			obj.Tags = append(obj.Tags, relation.Name)
		}
	}
}

// extractTagsFromProperties extracts tags from object's Properties array
func extractTagsFromProperties(obj *Object) {
	if len(obj.Properties) == 0 {
		return
	}

	for _, prop := range obj.Properties {
		if !isTagProperty(prop) {
			continue
		}

		for _, tag := range prop.MultiSelect {
			if tag.Name != "" {
				obj.Tags = append(obj.Tags, tag.Name)
			}
		}
	}
}

// isTagProperty checks if a property is a tag property
func isTagProperty(prop Property) bool {
	return prop.Name == "Tag" &&
		prop.Format == "multi_select" &&
		len(prop.MultiSelect) > 0
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
	c.initializeTypeCache(spaceID)

	// First try to find from the existing cache
	reverseCache := c.buildReverseCache(spaceID)
	if typeKey, found := reverseCache[typeName]; found {
		return typeKey, nil
	}

	// If not in cache, fetch all types and update cache
	types, err := c.GetTypes(ctx, &GetTypesParams{SpaceID: spaceID})
	if err != nil {
		return "", err
	}

	// Update cache with fresh data
	reverseCache = c.updateTypeCaches(spaceID, types.Data, typeName)

	// Try different matching strategies in order
	return c.findTypeKeyWithStrategies(ctx, spaceID, typeName, types.Data, reverseCache)
}

// initializeTypeCache initializes the type cache for a specific space
func (c *Client) initializeTypeCache(spaceID string) {
	if _, ok := c.typeCache[spaceID]; !ok {
		c.typeCache[spaceID] = make(map[string]string)
	}
}

// buildReverseCache builds a reverse lookup (name -> key) from existing cache
func (c *Client) buildReverseCache(spaceID string) map[string]string {
	reverseCache := make(map[string]string)

	if cache, ok := c.typeCache[spaceID]; ok && len(cache) > 0 {
		for key, name := range cache {
			reverseCache[name] = key
		}
	}

	return reverseCache
}

// updateTypeCaches updates both the regular and reverse caches with type data
func (c *Client) updateTypeCaches(spaceID string, types []TypeInfo, typeName string) map[string]string {
	reverseCache := make(map[string]string)

	for _, t := range types {
		c.typeCache[spaceID][t.Key] = t.Name
		reverseCache[t.Name] = t.Key

		// Handle special case: "Page" -> "ot-page", "Note" -> "ot-note" etc.
		if c.isOtPrefixMatch(t.Key, typeName) {
			reverseCache[typeName] = t.Key
		}
	}

	return reverseCache
}

// isOtPrefixMatch checks if a key with "ot-" prefix matches the typeName
func (c *Client) isOtPrefixMatch(key, typeName string) bool {
	if strings.HasPrefix(key, "ot-") && strings.EqualFold(strings.TrimPrefix(key, "ot-"), typeName) {
		if c.debug && c.logger != nil {
			c.logger.Debug("Found matching type by key prefix: '%s' -> '%s'", typeName, key)
		}
		return true
	}
	return false
}

// findTypeKeyWithStrategies tries different strategies to find a type key
func (c *Client) findTypeKeyWithStrategies(ctx context.Context, spaceID, typeName string,
	types []TypeInfo, reverseCache map[string]string) (string, error) {

	// Strategy 1: Exact match from updated cache
	if typeKey, found := reverseCache[typeName]; found {
		return typeKey, nil
	}

	// Strategy 2: Case-insensitive matching
	typeKey := c.findTypeKeyCaseInsensitive(typeName, reverseCache)
	if typeKey != "" {
		return typeKey, nil
	}

	// Strategy 3: Standard key construction (e.g., "Page" -> "ot-page")
	typeKey = c.findTypeKeyByStandardConstruction(typeName, types)
	if typeKey != "" {
		return typeKey, nil
	}

	return "", fmt.Errorf("type '%s' not found", typeName)
}

// findTypeKeyCaseInsensitive tries to find a type key using case-insensitive matching
func (c *Client) findTypeKeyCaseInsensitive(typeName string, reverseCache map[string]string) string {
	for name, key := range reverseCache {
		if strings.EqualFold(name, typeName) {
			if c.debug && c.logger != nil {
				c.logger.Debug("Found type using case-insensitive match: '%s' -> '%s'", typeName, name)
			}
			return key
		}
	}
	return ""
}

// findTypeKeyByStandardConstruction tries to construct a standard key format
func (c *Client) findTypeKeyByStandardConstruction(typeName string, types []TypeInfo) string {
	standardKey := "ot-" + strings.ToLower(typeName)
	for _, t := range types {
		if t.Key == standardKey {
			if c.debug && c.logger != nil {
				c.logger.Debug("Found type using standard key construction: '%s' -> '%s'", typeName, standardKey)
			}
			return standardKey
		}
	}
	return ""
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
// Search function is implemented in search.go

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
