package anytype

import (
	"encoding/json"
	"strings"

	"github.com/epheo/anytype-go/internal/log"
)

// ParseSearchResponse attempts to parse the API response data into a SearchResponse
// using multiple formats to handle different API versions or formats
func ParseSearchResponse(data []byte, debug bool, logger log.Logger) (*SearchResponse, error) {
	// Try direct format first
	var response SearchResponse
	err := json.Unmarshal(data, &response)
	if err == nil && len(response.Data) > 0 {
		return &response, nil
	}

	// If that didn't work, try the items format
	var itemsResponse struct {
		Items      []Object   `json:"items"`
		Total      int        `json:"total"`
		Limit      int        `json:"limit"`
		Offset     int        `json:"offset"`
		Pagination Pagination `json:"pagination,omitempty"`
	}

	err2 := json.Unmarshal(data, &itemsResponse)
	if err2 == nil && len(itemsResponse.Items) > 0 {
		// Convert to standard SearchResponse format
		return &SearchResponse{
			Data: itemsResponse.Items,
			Pagination: Pagination{
				Total:   itemsResponse.Total,
				Limit:   itemsResponse.Limit,
				Offset:  itemsResponse.Offset,
				HasMore: itemsResponse.Total > (itemsResponse.Offset + itemsResponse.Limit),
			},
		}, nil
	}

	// Try a more flexible approach with manual JSON structure extraction
	var anyFormat map[string]interface{}
	err3 := json.Unmarshal(data, &anyFormat)
	if err3 == nil {
		// Check for common variations in the response structure
		if debug && logger != nil {
			var keys []string
			for k := range anyFormat {
				keys = append(keys, k)
			}
			logger.Debug("Response keys: %v", strings.Join(keys, ", "))
		}

		// Add more handler variations here based on the actual API response structure

		// If we get this far, check if response is actually an error message
		if errorData, ok := anyFormat["error"]; ok {
			if errObj, ok := errorData.(map[string]interface{}); ok {
				if msg, ok := errObj["message"].(string); ok {
					return nil, NewError(msg)
				}
			}
		}
	}

	// If all attempts fail, return original error
	if debug && logger != nil {
		logger.Debug("Failed to parse search response: %v", err)
	}

	return nil, err
}

// NewError creates a simple error with a message
func NewError(message string) error {
	return &SearchError{Message: message}
}

// SearchError represents an error in search processing
type SearchError struct {
	Message string
}

func (e *SearchError) Error() string {
	return e.Message
}
