package anytype

import (
	"errors"
	"time"
)

// Common errors
var (
	ErrInvalidSpaceID   = errors.New("invalid space ID")
	ErrInvalidObjectID  = errors.New("invalid object ID")
	ErrInvalidTypeID    = errors.New("invalid type ID")
	ErrInvalidTemplate  = errors.New("invalid template")
	ErrInvalidParameter = errors.New("invalid parameter")
)

// Response types
type (
	// ChallengeResponse represents the response from the challenge endpoint
	ChallengeResponse struct {
		ChallengeID string `json:"challenge_id"`
	}

	// AuthResponse represents the authentication token response
	AuthResponse struct {
		SessionToken string `json:"session_token"`
		AppKey       string `json:"app_key"`
	}

	// AuthConfig stores authentication configuration
	AuthConfig struct {
		ApiURL       string    `json:"api_url"`
		SessionToken string    `json:"session_token"`
		AppKey       string    `json:"app_key"`
		Timestamp    time.Time `json:"timestamp"`
	}

	// Space represents a space in Anytype
	Space struct {
		Type                   string `json:"type"`
		ID                     string `json:"id"`
		Name                   string `json:"name"`
		Icon                   string `json:"icon,omitempty"`
		HomeObjectID           string `json:"home_object_id,omitempty"`
		ArchiveObjectID        string `json:"archive_object_id,omitempty"`
		ProfileObjectID        string `json:"profile_object_id,omitempty"`
		MarketplaceWorkspaceID string `json:"marketplace_workspace_id,omitempty"`
		WorkspaceObjectID      string `json:"workspace_object_id,omitempty"`
		DeviceID               string `json:"device_id,omitempty"`
		AccountSpaceID         string `json:"account_space_id,omitempty"`
		WidgetsID              string `json:"widgets_id,omitempty"`
		SpaceViewID            string `json:"space_view_id,omitempty"`
		TechSpaceID            string `json:"tech_space_id,omitempty"`
		GatewayURL             string `json:"gateway_url,omitempty"`
		LocalStoragePath       string `json:"local_storage_path,omitempty"`
		Timezone               string `json:"timezone,omitempty"`
		AnalyticsID            string `json:"analytics_id,omitempty"`
		NetworkID              string `json:"network_id,omitempty"`
	}

	// SpacesResponse represents the response from the spaces endpoint
	SpacesResponse struct {
		Data       []Space    `json:"data"`
		Pagination Pagination `json:"pagination"`
	}

	// Pagination represents common pagination information
	Pagination struct {
		Total   int  `json:"total"`
		Offset  int  `json:"offset"`
		Limit   int  `json:"limit"`
		HasMore bool `json:"has_more"`
	}

	// Object represents an object in a space
	Object struct {
		ID        string                 `json:"id"`
		Type      string                 `json:"type"`
		Name      string                 `json:"name"`
		Props     map[string]interface{} `json:"props,omitempty"`
		CreatedAt time.Time              `json:"created_at"`
		UpdatedAt time.Time              `json:"updated_at"`
		Tags      []string               `json:"tags,omitempty"`
	}

	// ObjectsResponse represents the response from the space objects endpoint
	ObjectsResponse struct {
		Objects    []Object   `json:"objects"`
		Pagination Pagination `json:"pagination"`
	}

	// SearchParams represents search parameters
	SearchParams struct {
		Query  string   `json:"query,omitempty"`
		Types  []string `json:"types,omitempty"`
		Tags   []string `json:"tags,omitempty"`
		Limit  int      `json:"limit,omitempty"`
		Offset int      `json:"offset,omitempty"`
	}

	// SearchResponse represents the response from search endpoints
	SearchResponse struct {
		Items      []Object   `json:"items"`
		Pagination Pagination `json:"pagination"`
	}
)

// Validate validates Space fields
func (s *Space) Validate() error {
	if s.ID == "" {
		return ErrInvalidSpaceID
	}
	return nil
}

// Validate validates Object fields
func (o *Object) Validate() error {
	if o.ID == "" {
		return ErrInvalidObjectID
	}
	if o.Type == "" {
		return ErrInvalidTypeID
	}
	return nil
}

// Validate validates SearchParams fields
func (p *SearchParams) Validate() error {
	if p.Limit < 0 || p.Offset < 0 {
		return ErrInvalidParameter
	}
	return nil
}

// NewSearchParams creates a new SearchParams with default values
func NewSearchParams() *SearchParams {
	return &SearchParams{
		Limit:  defaultSearchLimit,
		Offset: defaultSearchOffset,
	}
}
