package anytype

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SupportedExportFormats defines the available export formats
var SupportedExportFormats = []string{"md", "markdown", "html"}

// ExportObject exports an object's content to a file in the specified format
func (c *Client) ExportObject(ctx context.Context, spaceID, objectID, exportPath, format string) (string, error) {
	if spaceID == "" {
		return "", ErrInvalidSpaceID
	}
	if objectID == "" {
		return "", ErrInvalidObjectID
	}
	if exportPath == "" {
		return "", fmt.Errorf("export path cannot be empty")
	}

	// Normalize format
	format = strings.ToLower(format)
	if format == "markdown" {
		format = "md" // Normalize markdown to md
	}

	// Validate format
	validFormat := false
	for _, f := range SupportedExportFormats {
		if format == f {
			validFormat = true
			break
		}
	}
	if !validFormat {
		return "", fmt.Errorf("unsupported export format: %s", format)
	}

	// Get the object to get its metadata
	object, err := c.GetObject(ctx, spaceID, objectID)
	if err != nil {
		return "", fmt.Errorf("failed to get object %s: %w", objectID, err)
	}

	// Construct the export path
	// Create directory if it doesn't exist
	if err := os.MkdirAll(exportPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create export directory: %w", err)
	}

	// Sanitize object name for file system use
	sanitizedName := sanitizeFilename(object.Name)
	if sanitizedName == "" {
		sanitizedName = fmt.Sprintf("object_%s", objectID)
	}

	// Construct file path with timestamp to avoid overwriting
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.%s", sanitizedName, timestamp, format)
	filePath := filepath.Join(exportPath, filename)

	// Get object content in the requested format
	content, err := c.getObjectContent(ctx, spaceID, objectID, format)
	if err != nil {
		return "", fmt.Errorf("failed to get object content: %w", err)
	}

	// Write content to file
	if err := ioutil.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write to file: %w", err)
	}

	return filePath, nil
}

// sanitizeFilename removes characters that are invalid in filenames
func sanitizeFilename(name string) string {
	// Replace common invalid filename characters
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return replacer.Replace(name)
}

// getObjectContent retrieves the content of an object in the specified format
func (c *Client) getObjectContent(ctx context.Context, spaceID, objectID, format string) (string, error) {
	// Construct API path for content export
	path := fmt.Sprintf("/v1/spaces/%s/objects/%s/export?format=%s", spaceID, objectID, format)

	// Make API request
	data, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return "", fmt.Errorf("failed to export object %s: %w", objectID, err)
	}

	// Check if the response is empty
	if len(data) == 0 {
		return "", fmt.Errorf("received empty response for object %s", objectID)
	}

	// Parse the response - content should be directly in the response body
	return string(data), nil
}

// ExportObjects exports multiple objects to files
func (c *Client) ExportObjects(ctx context.Context, spaceID string, objects []Object, exportPath, format string) ([]string, error) {
	if len(objects) == 0 {
		return nil, fmt.Errorf("no objects to export")
	}

	exportedFiles := make([]string, 0, len(objects))
	errors := make([]string, 0)

	for _, obj := range objects {
		filePath, err := c.ExportObject(ctx, spaceID, obj.ID, exportPath, format)
		if err != nil {
			// Log error but continue with other objects
			errMsg := fmt.Sprintf("Failed to export object %s (%s): %v", obj.ID, obj.Name, err)
			errors = append(errors, errMsg)
			if c.logger != nil {
				c.logger.Error(errMsg)
			}
			continue
		}
		exportedFiles = append(exportedFiles, filePath)
	}

	if len(exportedFiles) == 0 {
		if len(errors) > 0 {
			// Return the first few errors to help diagnose the problem
			maxErrors := 3
			if len(errors) < maxErrors {
				maxErrors = len(errors)
			}
			return nil, fmt.Errorf("failed to export any objects. First %d errors: %s",
				maxErrors, strings.Join(errors[:maxErrors], "; "))
		}
		return nil, fmt.Errorf("failed to export any objects")
	}

	// If some objects were exported successfully but others failed, log the count
	if len(errors) > 0 && c.logger != nil {
		c.logger.Info("Exported %d objects successfully, %d objects failed",
			len(exportedFiles), len(errors))
	}

	return exportedFiles, nil
}
