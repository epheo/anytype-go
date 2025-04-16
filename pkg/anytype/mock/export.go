package mock

import (
	"context"
	"fmt"

	"github.com/epheo/anyblog/pkg/anytype"
)

// ExportObject mocks the ExportObject method
func (c *Client) ExportObject(ctx context.Context, spaceID, objectID, exportPath, format string) (string, error) {
	// Record the call
	c.LastSpaceID = spaceID
	c.LastObjectID = objectID

	// Mock successful export
	if c.GetObjectError != nil {
		return "", c.GetObjectError
	}

	// Return a mock file path
	return fmt.Sprintf("%s/mock_export_%s.%s", exportPath, objectID, format), nil
}

// ExportObjects mocks the ExportObjects method
func (c *Client) ExportObjects(ctx context.Context, spaceID string, objects []anytype.Object, exportPath, format string) ([]string, error) {
	// Record the call
	c.LastSpaceID = spaceID

	// Mock successful export
	if c.GetObjectError != nil {
		return nil, c.GetObjectError
	}

	// Return mock file paths
	paths := make([]string, len(objects))
	for i, obj := range objects {
		paths[i] = fmt.Sprintf("%s/mock_export_%s.%s", exportPath, obj.ID, format)
	}

	return paths, nil
}
