package anytype

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// validateExportParams validates the export parameters
func validateExportParams(spaceID, objectID, exportPath string) error {
	if spaceID == "" {
		return ErrInvalidSpaceID
	}
	if objectID == "" {
		return ErrInvalidObjectID
	}
	if exportPath == "" {
		return fmt.Errorf("export path cannot be empty")
	}
	return nil
}

// normalizeExportFormat normalizes and validates the export format
func (c *Client) normalizeExportFormat(format string) string {
	// Normalize format
	format = strings.ToLower(format)

	// Convert "md" to "markdown" for API calls
	if format == "md" {
		format = "markdown"
	}

	return format
}

// validateExportFormat checks if the format is supported and logs a warning if not
func (c *Client) validateExportFormat(format string) bool {
	validFormat := false
	for _, f := range SupportedExportFormats {
		if format == f {
			validFormat = true
			break
		}
	}

	// Only "markdown" is officially supported by the API,
	// but we'll allow other formats with a warning
	if !validFormat && c.logger != nil {
		c.logger.Info("Format '%s' is not officially supported by the Anytype API (only 'markdown' is guaranteed). Attempting export anyway.", format)
	}

	return validFormat
}

// createExportDirectory creates necessary directories for export
func (c *Client) createExportDirectory(exportPath, typeName string) (string, error) {
	// Create type-specific subdirectory
	typeSubdir := filepath.Join(exportPath, typeName)
	if err := os.MkdirAll(typeSubdir, 0755); err != nil {
		return "", fmt.Errorf("failed to create export directory: %w", err)
	}
	return typeSubdir, nil
}

// getTypeNameForExport retrieves a sanitized type name for the export directory
func getTypeNameForExport(object *Object) string {
	typeName := "Unknown"
	if object.Type != nil && object.Type.Name != "" {
		typeName = sanitizeFilename(object.Type.Name)
	}
	return typeName
}

// getExportFilename generates a filename for the exported file
func getExportFilename(object *Object, objectID string, format string) string {
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
	return fmt.Sprintf("%s.%s", sanitizedName, fileExtension)
}

// writeExportFile writes content to the export file
func (c *Client) writeExportFile(filePath string, content string) error {
	// Write the content to the file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	if c.logger != nil {
		c.logger.Info("Object exported successfully to: %s", filePath)
	}

	return nil
}

// processExportContent processes the content based on format
func (c *Client) processExportContent(
	ctx context.Context,
	content string,
	format string,
	exportPath string,
) (string, error) {
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
				c.logger.Info("Continuing with original content without image processing")
			}
			return content, nil
		}
		return processedContent, nil
	}

	return content, nil
}

// ExportObjectImpl is the implementation of ExportObject with reduced cyclomatic complexity
func (c *Client) ExportObjectImpl(ctx context.Context, spaceID, objectID, exportPath, format string) (string, error) {
	// Validate input parameters
	if err := validateExportParams(spaceID, objectID, exportPath); err != nil {
		return "", err
	}

	// Normalize and validate format
	format = c.normalizeExportFormat(format)
	c.validateExportFormat(format)

	// Get the object to retrieve its metadata
	object, err := c.GetObject(ctx, &GetObjectParams{
		SpaceID:  spaceID,
		ObjectID: objectID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get object %s: %w", objectID, err)
	}

	// Log debug information about the object
	if c.logger != nil {
		c.logger.Debug("Exporting object - ID: %s, Name: %s, Type: %s",
			object.ID, object.Name, object.Type)
	}

	// Get type name and create export directory
	typeName := getTypeNameForExport(object)
	typeSubdir, err := c.createExportDirectory(exportPath, typeName)
	if err != nil {
		return "", err
	}

	// Generate the export filename
	filename := getExportFilename(object, objectID, format)
	filePath := filepath.Join(typeSubdir, filename)

	// Get object content in the requested format
	content, err := c.getObjectContent(ctx, spaceID, objectID, format)
	if err != nil {
		return "", fmt.Errorf("failed to get object content: %w", err)
	}

	// Process content based on format (handle images, etc.)
	processedContent, err := c.processExportContent(ctx, content, format, exportPath)
	if err == nil && processedContent != "" {
		content = processedContent
	}

	// Write content to the file
	if err := c.writeExportFile(filePath, content); err != nil {
		return "", err
	}

	return filePath, nil
}
