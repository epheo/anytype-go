package anytype

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SupportedExportFormats defines the available export formats
// The API officially supports only "markdown" format
var SupportedExportFormats = []string{"markdown"}

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

	// Convert "md" to "markdown" for API calls
	if format == "md" {
		format = "markdown"
	}

	// Validate format against officially supported formats
	validFormat := false
	for _, f := range SupportedExportFormats {
		if format == f {
			validFormat = true
			break
		}
	}

	// Only "markdown" is officially supported by the API,
	// but we'll allow other formats with a warning
	if !validFormat {
		if c.logger != nil {
			c.logger.Info("Format '%s' is not officially supported by the Anytype API (only 'markdown' is guaranteed). Attempting export anyway.", format)
		}
	}

	// Get the object to get its metadata
	object, err := c.GetObject(ctx, spaceID, objectID)
	if err != nil {
		return "", fmt.Errorf("failed to get object %s: %w", objectID, err)
	}

	// Get type name for the subdirectory
	typeName := "Unknown"
	if object.Type != nil && object.Type.Name != "" {
		typeName = sanitizeFilename(object.Type.Name)
	}

	// Create type-specific subdirectory
	typeSubdir := filepath.Join(exportPath, typeName)
	if err := os.MkdirAll(typeSubdir, 0755); err != nil {
		return "", fmt.Errorf("failed to create export directory: %w", err)
	}

	// Log debug information about the object
	if c.logger != nil {
		c.logger.Debug("Exporting object - ID: %s, Name: %s, Type: %s",
			object.ID, object.Name, object.Type)
	}

	// Sanitize object name for file system use
	var sanitizedName string
	if object.Name != "" {
		sanitizedName = sanitizeFilename(object.Name)
		// Convert spaces to hyphens and ensure filename is clean
		sanitizedName = strings.ReplaceAll(sanitizedName, " ", "-")
		// Remove consecutive hyphens
		for strings.Contains(sanitizedName, "--") {
			sanitizedName = strings.ReplaceAll(sanitizedName, "--", "-")
		}
		// Remove any leading or trailing hyphens
		sanitizedName = strings.Trim(sanitizedName, "-")
	}

	// Use object ID as fallback only if name is empty or becomes empty after sanitization
	if sanitizedName == "" {
		sanitizedName = fmt.Sprintf("object-%s", objectID)
	}

	// Determine proper file extension based on the format
	var fileExtension string
	switch format {
	case "markdown":
		fileExtension = "md"
	default:
		fileExtension = format
	}

	// Create filename without timestamp to allow overwriting
	filename := fmt.Sprintf("%s.%s", sanitizedName, fileExtension)
	filePath := filepath.Join(typeSubdir, filename)

	// Get object content in the requested format
	content, err := c.getObjectContent(ctx, spaceID, objectID, format)
	if err != nil {
		return "", fmt.Errorf("failed to get object content: %w", err)
	}

	// Process images in the content (download them and update references)
	if format == "markdown" {
		if c.logger != nil {
			c.logger.Debug("Processing images in markdown content")
		}
		processedContent, err := c.ProcessMarkdownImages(ctx, content, exportPath)
		if err != nil {
			// Log the error but continue with the original content
			if c.logger != nil {
				c.logger.Error("Failed to process images: %v", err)
			}
		} else {
			content = processedContent
		}
	}

	// Write content to file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
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
	// Construct API path for content export based on the API documentation
	// The API endpoint is /v1/spaces/{space_id}/objects/{object_id}/{format}
	path := fmt.Sprintf("/v1/spaces/%s/objects/%s/%s", spaceID, objectID, format)

	// Make API request
	data, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		// If the export endpoint returned a 404, try to extract content from the regular object GET endpoint
		if strings.Contains(err.Error(), "returned status 404") {
			if c.logger != nil {
				c.logger.Debug("Export endpoint returned 404, trying to extract content from regular object endpoint")
			}
			return c.extractObjectContentFromRegularEndpoint(ctx, spaceID, objectID)
		}
		return "", fmt.Errorf("failed to export object %s: %w", objectID, err)
	}

	// Check if the response is empty
	if len(data) == 0 {
		return "", fmt.Errorf("received empty response for object %s", objectID)
	}

	// Try to parse the response as JSON
	var response struct {
		Markdown string `json:"markdown"`
		Content  string `json:"content"`
	}

	if err := json.Unmarshal(data, &response); err == nil {
		// Successfully parsed as JSON
		// Check for markdown field first
		if response.Markdown != "" {
			return response.Markdown, nil
		}

		// Fall back to generic content field
		if response.Content != "" {
			return response.Content, nil
		}
	}

	// If JSON parsing fails or no appropriate content field found,
	// return the raw data as it might be raw markdown
	return string(data), nil
}

// extractObjectContentFromRegularEndpoint tries to get the content from the regular object endpoint
// and format it as markdown as a fallback when the export endpoint doesn't work
func (c *Client) extractObjectContentFromRegularEndpoint(ctx context.Context, spaceID, objectID string) (string, error) {
	// Get the object's full details from the regular endpoint
	obj, err := c.GetObject(ctx, spaceID, objectID)
	if err != nil {
		return "", fmt.Errorf("failed to get object details: %w", err)
	}

	// Build a proper markdown representation from the object's data
	var sb strings.Builder

	// Add title with icon if available
	if obj.Name != "" {
		if obj.Icon != nil {
			var iconDisplay string
			if obj.Icon.Emoji != "" {
				iconDisplay = obj.Icon.Emoji
			} else if obj.Icon.Name != "" {
				iconDisplay = obj.Icon.Name
			}
			if iconDisplay != "" {
				sb.WriteString(fmt.Sprintf("# %s %s\n\n", iconDisplay, obj.Name))
			} else {
				sb.WriteString(fmt.Sprintf("# %s\n\n", obj.Name))
			}
		} else {
			sb.WriteString(fmt.Sprintf("# %s\n\n", obj.Name))
		}
	}

	// Add tags if available
	if len(obj.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("**Tags:** %s\n\n", strings.Join(obj.Tags, ", ")))
	}

	// Add snippet as content
	if obj.Snippet != "" {
		sb.WriteString(obj.Snippet)
		sb.WriteString("\n\n")
	}

	// Add metadata in a discreet way at the bottom
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("Type: %s  \n", obj.Type.Name))
	if obj.Layout != "" {
		sb.WriteString(fmt.Sprintf("Layout: %s  \n", obj.Layout))
	}

	return sb.String(), nil
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
