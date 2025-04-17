package anytype

import (
	"encoding/json"
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

// Icon represents an icon in Anytype
// Matches the util.Icon schema in the API documentation
type Icon struct {
	Format string `json:"format,omitempty"` // Format of the icon: emoji, file, or icon
	Emoji  string `json:"emoji,omitempty"`  // Emoji character if format is emoji
	Name   string `json:"name,omitempty"`   // Name of the icon if format is icon
	File   string `json:"file,omitempty"`   // File URL if format is file
	Color  string `json:"color,omitempty"`  // Color of the icon
}

// Block represents a content block in an object
// Matches the object.Block schema in the API documentation
type Block struct {
	ID              string      `json:"id,omitempty"`               // Block ID
	ChildrenIDs     []string    `json:"children_ids,omitempty"`     // Child block IDs
	Align           string      `json:"align,omitempty"`            // Alignment: AlignLeft, AlignCenter, etc.
	VerticalAlign   string      `json:"vertical_align,omitempty"`   // Vertical alignment
	BackgroundColor string      `json:"background_color,omitempty"` // Background color
	Text            *TextBlock  `json:"text,omitempty"`             // Text content if applicable
	File            *FileBlock  `json:"file,omitempty"`             // File content if applicable
	Property        interface{} `json:"property,omitempty"`         // Property information if applicable
}

// TextBlock represents text content in a block
// Matches the object.Text schema in the API documentation
type TextBlock struct {
	Text    string `json:"text,omitempty"`    // Text content
	Style   string `json:"style,omitempty"`   // Style: Paragraph, Header1, etc.
	Checked bool   `json:"checked,omitempty"` // Whether the text is checked (for checkboxes)
	Color   string `json:"color,omitempty"`   // Text color
	Icon    *Icon  `json:"icon,omitempty"`    // Icon for the text block
}

// FileBlock represents file content in a block
// Matches the object.File schema in the API documentation
type FileBlock struct {
	Hash           string `json:"hash,omitempty"`             // File hash
	Name           string `json:"name,omitempty"`             // File name
	Mime           string `json:"mime,omitempty"`             // MIME type
	Size           int    `json:"size,omitempty"`             // File size in bytes
	Type           string `json:"type,omitempty"`             // File type
	State          string `json:"state,omitempty"`            // File state
	Style          string `json:"style,omitempty"`            // Display style
	TargetObjectID string `json:"target_object_id,omitempty"` // Target object ID
	AddedAt        int64  `json:"added_at,omitempty"`         // Timestamp when added
}

// TypeInfo represents a type in Anytype
// Matches the object.Type schema in the API documentation
type TypeInfo struct {
	Object            string `json:"object,omitempty"`             // Data model, always "type"
	ID                string `json:"id,omitempty"`                 // Unique ID of the type
	Key               string `json:"key,omitempty"`                // Consistent key across spaces (e.g., "ot-page")
	Name              string `json:"name,omitempty"`               // Display name of the type
	Icon              *Icon  `json:"icon,omitempty"`               // Type icon
	Archived          bool   `json:"archived,omitempty"`           // Whether the type is archived
	RecommendedLayout string `json:"recommended_layout,omitempty"` // Recommended layout for this type
}

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
	// Matches the space.Space schema in the API documentation
	Space struct {
		Object            string   `json:"object,omitempty"`              // Data model, e.g. "space"
		ID                string   `json:"id,omitempty"`                  // Unique ID of the space
		Name              string   `json:"name,omitempty"`                // Display name of the space
		Icon              *Icon    `json:"icon,omitempty"`                // Space icon
		Description       string   `json:"description,omitempty"`         // Description of the space
		GatewayURL        string   `json:"gateway_url,omitempty"`         // Gateway URL for files and media
		NetworkID         string   `json:"network_id,omitempty"`          // Network ID of the space
		HomeObjectID      string   `json:"home_object_id,omitempty"`      // Home object ID
		ArchiveObjectID   string   `json:"archive_object_id,omitempty"`   // Archive object ID
		ProfileObjectID   string   `json:"profile_object_id,omitempty"`   // Profile object ID
		WorkspaceObjectID string   `json:"workspace_object_id,omitempty"` // Workspace object ID
		DeviceID          string   `json:"device_id,omitempty"`           // Device ID
		AccountSpaceID    string   `json:"account_space_id,omitempty"`    // Account space ID
		SpaceViewID       string   `json:"space_view_id,omitempty"`       // Space view ID
		LocalPath         string   `json:"local_path,omitempty"`          // Local path
		Timezone          string   `json:"timezone,omitempty"`            // Timezone
		IsReadOnly        bool     `json:"is_read_only,omitempty"`        // Whether the space is read-only
		CanDelete         bool     `json:"can_delete,omitempty"`          // Whether the space can be deleted
		CanLeave          bool     `json:"can_leave,omitempty"`           // Whether the space can be left
		Role              string   `json:"role,omitempty"`                // User's role in the space
		Members           []Member `json:"-"`                             // Members populated separately
	}

	// Member represents a member of a space
	// Matches the space.Member schema in the API documentation
	Member struct {
		Object     string `json:"object,omitempty"`      // Data model, e.g. "member"
		ID         string `json:"id,omitempty"`          // Member ID
		Name       string `json:"name,omitempty"`        // Member name
		Icon       *Icon  `json:"icon,omitempty"`        // Member icon
		Identity   string `json:"identity,omitempty"`    // Network identity
		GlobalName string `json:"global_name,omitempty"` // Global name in network
		Role       string `json:"role,omitempty"`        // Role: viewer, editor, owner, no_permission
		Status     string `json:"status,omitempty"`      // Status: joining, active, removed, declined, removing, canceled
	}

	// MembersResponse represents the response from the members endpoint
	// Matches the pagination.PaginatedResponse-space_Member schema
	MembersResponse struct {
		Data       []Member   `json:"data"`
		Pagination Pagination `json:"pagination"`
	}

	// SpacesResponse represents the response from the spaces endpoint
	// Matches the pagination.PaginatedResponse-space_Space schema
	SpacesResponse struct {
		Data       []Space    `json:"data"`
		Pagination Pagination `json:"pagination"`
	}

	// Pagination represents common pagination information
	// Matches the pagination.PaginationMeta schema
	Pagination struct {
		Total   int  `json:"total"`    // Total available items
		Offset  int  `json:"offset"`   // Items skipped
		Limit   int  `json:"limit"`    // Max items in response
		HasMore bool `json:"has_more"` // More items available beyond current result
	}

	// PropertyTag represents a tag in a multi_select property
	// Matches the object.Tag schema
	PropertyTag struct {
		ID    string `json:"id,omitempty"`    // Tag ID
		Name  string `json:"name,omitempty"`  // Tag name
		Color string `json:"color,omitempty"` // Tag color
	}

	// Property represents a property of an object
	// Matches the object.Property schema
	Property struct {
		ID          string        `json:"id,omitempty"`           // Property ID
		Name        string        `json:"name,omitempty"`         // Property name
		Format      string        `json:"format,omitempty"`       // Property format type
		MultiSelect []PropertyTag `json:"multi_select,omitempty"` // Multi-select values
		Date        string        `json:"date,omitempty"`         // Date value
		Object      []string      `json:"object,omitempty"`       // Object references
		Number      float64       `json:"number,omitempty"`       // Number value
		Text        string        `json:"text,omitempty"`         // Text value
		URL         string        `json:"url,omitempty"`          // URL value
		Email       string        `json:"email,omitempty"`        // Email value
		Phone       string        `json:"phone,omitempty"`        // Phone value
		Checkbox    bool          `json:"checkbox,omitempty"`     // Checkbox value
		File        []string      `json:"file,omitempty"`         // File references
		Select      *PropertyTag  `json:"select,omitempty"`       // Single select value
	}

	// Object represents an object in a space
	// Matches the object.Object schema in the API documentation
	Object struct {
		Object     string                 `json:"object,omitempty"`     // Data model, e.g. "object"
		ID         string                 `json:"id,omitempty"`         // Unique ID of the object
		Name       string                 `json:"name,omitempty"`       // Display name of the object
		Type       *TypeInfo              `json:"type,omitempty"`       // Type information
		Icon       *Icon                  `json:"icon,omitempty"`       // Object icon
		Archived   bool                   `json:"archived,omitempty"`   // Whether the object is archived
		SpaceID    string                 `json:"space_id,omitempty"`   // ID of the space the object belongs to
		Snippet    string                 `json:"snippet,omitempty"`    // Preview/snippet of the object content
		Layout     string                 `json:"layout,omitempty"`     // Layout of the object e.g. "basic"
		Blocks     []Block                `json:"blocks,omitempty"`     // Content blocks of the object
		Relations  map[string]interface{} `json:"relations,omitempty"`  // Relations/links to other objects
		Properties []Property             `json:"properties,omitempty"` // Properties/metadata of the object
		Tags       []string               `json:"-"`                    // Tags is a client-side representation for convenience
	}

	// SearchParams represents search parameters
	// Matches the search.SearchRequest schema
	SearchParams struct {
		SpaceID string       `json:"space_id,omitempty"` // Space ID to search in
		Query   string       `json:"query,omitempty"`    // Search term
		Types   []string     `json:"types,omitempty"`    // Object types to include
		Tags    []string     `json:"tags,omitempty"`     // Tags to filter by (client-side)
		Sort    *SortOptions `json:"sort,omitempty"`     // Sorting options
		Limit   int          `json:"limit,omitempty"`    // Result limit
		Offset  int          `json:"offset,omitempty"`   // Result offset
	}

	// SortOptions represents sorting criteria for search results
	// Matches the search.SortOptions schema
	SortOptions struct {
		Property  string `json:"property,omitempty"`  // Property to sort by
		Direction string `json:"direction,omitempty"` // Sort direction (asc or desc)
	}

	// SearchResponse represents the response from search endpoints
	// Matches the pagination.PaginatedResponse-object_Object schema
	SearchResponse struct {
		Data       []Object   `json:"data"`
		Pagination Pagination `json:"pagination"`
	}
)

// GetObjectParams represents parameters for retrieving an object
type GetObjectParams struct {
	SpaceID  string `json:"space_id"`  // Space ID the object belongs to
	ObjectID string `json:"object_id"` // Object ID to retrieve
}

// Validate validates GetObjectParams fields
func (p *GetObjectParams) Validate() error {
	if p.SpaceID == "" {
		return ErrInvalidSpaceID
	}
	if p.ObjectID == "" {
		return ErrInvalidObjectID
	}
	return nil
}

// CreateObjectParams represents parameters for creating an object
type CreateObjectParams struct {
	SpaceID string  `json:"space_id"`         // Space ID where object should be created
	Object  *Object `json:"object,omitempty"` // Object to create
}

// Validate validates CreateObjectParams fields
func (p *CreateObjectParams) Validate() error {
	if p.SpaceID == "" {
		return ErrInvalidSpaceID
	}
	if p.Object == nil {
		return ErrInvalidParameter
	}
	return p.Object.Validate()
}

// UpdateObjectParams represents parameters for updating an object
type UpdateObjectParams struct {
	SpaceID  string  `json:"space_id"`         // Space ID the object belongs to
	ObjectID string  `json:"object_id"`        // Object ID to update
	Object   *Object `json:"object,omitempty"` // Object data for update
}

// Validate validates UpdateObjectParams fields
func (p *UpdateObjectParams) Validate() error {
	if p.SpaceID == "" {
		return ErrInvalidSpaceID
	}
	if p.ObjectID == "" {
		return ErrInvalidObjectID
	}
	if p.Object == nil {
		return ErrInvalidParameter
	}
	return nil
}

// DeleteObjectParams represents parameters for deleting an object
type DeleteObjectParams struct {
	SpaceID  string `json:"space_id"`  // Space ID the object belongs to
	ObjectID string `json:"object_id"` // Object ID to delete
}

// Validate validates DeleteObjectParams fields
func (p *DeleteObjectParams) Validate() error {
	if p.SpaceID == "" {
		return ErrInvalidSpaceID
	}
	if p.ObjectID == "" {
		return ErrInvalidObjectID
	}
	return nil
}

// GetSpacesParams represents parameters for retrieving spaces
type GetSpacesParams struct {
	IncludeMembers bool `json:"include_members,omitempty"` // Whether to include member information for each space
}

// NewGetSpacesParams creates a new GetSpacesParams with default values
func NewGetSpacesParams() *GetSpacesParams {
	return &GetSpacesParams{
		IncludeMembers: true,
	}
}

// GetSpaceByIDParams represents parameters for retrieving a specific space
type GetSpaceByIDParams struct {
	SpaceID string `json:"space_id"` // Space ID to retrieve
}

// Validate validates GetSpaceByIDParams fields
func (p *GetSpaceByIDParams) Validate() error {
	if p.SpaceID == "" {
		return ErrInvalidSpaceID
	}
	return nil
}

// GetTypesParams represents parameters for retrieving types from a space
type GetTypesParams struct {
	SpaceID string `json:"space_id"` // Space ID to get types from
}

// Validate validates GetTypesParams fields
func (p *GetTypesParams) Validate() error {
	if p.SpaceID == "" {
		return ErrInvalidSpaceID
	}
	return nil
}

// GetTypeByNameParams represents parameters for retrieving a type by name
type GetTypeByNameParams struct {
	SpaceID  string `json:"space_id"`  // Space ID to get type from
	TypeName string `json:"type_name"` // Type name to lookup
}

// Validate validates GetTypeByNameParams fields
func (p *GetTypeByNameParams) Validate() error {
	if p.SpaceID == "" {
		return ErrInvalidSpaceID
	}
	if p.TypeName == "" {
		return ErrInvalidTypeID
	}
	return nil
}

// GetMembersParams represents parameters for retrieving space members
type GetMembersParams struct {
	SpaceID string `json:"space_id"` // Space ID to get members from
}

// Validate validates GetMembersParams fields
func (p *GetMembersParams) Validate() error {
	if p.SpaceID == "" {
		return ErrInvalidSpaceID
	}
	return nil
}

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
	if o.Type == nil || o.Type.Key == "" {
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

// UnmarshalJSON implements custom JSON unmarshaling for Object
func (o *Object) UnmarshalJSON(data []byte) error {
	type Alias Object
	aux := &struct {
		*Alias
		Details []struct {
			ID      string `json:"id"`
			Details struct {
				Tags []struct {
					ID    string `json:"id"`
					Name  string `json:"name"`
					Color string `json:"color"`
				} `json:"tags"`
			} `json:"details"`
		} `json:"details"`
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Extract tags from the nested details structure
	o.Tags = []string{}
	for _, detail := range aux.Details {
		if detail.ID == "tags" {
			for _, tag := range detail.Details.Tags {
				o.Tags = append(o.Tags, tag.Name)
			}
			break
		}
	}

	return nil
}
