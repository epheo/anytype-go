package anytype

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Search defaults
const (
	defaultSearchLimit  = 100
	defaultSearchOffset = 0
)

// SearchRequestBody represents the structure of a search request
type SearchRequestBody struct {
	SpaceID string      `json:"spaceId"`
	Query   string      `json:"query"`
	Types   []string    `json:"types,omitempty"`
	Tags    []string    `json:"tags,omitempty"`
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

	types, err := c.GetTypes(ctx, spaceID)
	if err != nil {
		return "", err
	}

	for _, t := range types.Data {
		if t.Name == typeName {
			return t.UniqueKey, nil
		}
	}

	return "", fmt.Errorf("type '%s' not found", typeName)
}

// Search performs a search in a space with the given parameters
func (c *Client) Search(ctx context.Context, spaceID string, params *SearchParams) (*SearchResponse, error) {
	if spaceID == "" {
		return nil, ErrInvalidSpaceID
	}
	if params == nil {
		params = NewSearchParams()
	}
	if err := params.Validate(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/v1/spaces/%s/search", spaceID)
	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search params: %w", err)
	}

	data, err := c.makeRequest(ctx, http.MethodPost, path, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to perform search: %w", err)
	}

	var response SearchResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
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

// PrintCurlRequest prints a curl command equivalent to the HTTP request for debugging
func (c *Client) PrintCurlRequest(method, url string, headers map[string]string, body []byte) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("curl -X %s \\\n  '%s'", method, url))

	// Add headers
	for key, value := range headers {
		sb.WriteString(fmt.Sprintf(" \\\n  -H '%s: %s'", key, value))
	}

	// Add body with proper JSON formatting if possible
	if len(body) > 0 {
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, body, "  ", "  "); err == nil {
			sb.WriteString(fmt.Sprintf(" \\\n  -d '%s'", prettyJSON.String()))
		} else {
			sb.WriteString(fmt.Sprintf(" \\\n  -d '%s'", string(body)))
		}
	}

	log.Printf("Debug curl command:\n%s", sb.String())
}
