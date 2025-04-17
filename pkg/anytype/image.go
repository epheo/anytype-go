package anytype

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// imageURLPattern is the regex pattern to find image references in markdown
var imageURLPattern = regexp.MustCompile(`!\[([^\]]*)\]\((http://127\.0\.0\.1:[0-9]+/image/[^)]+)\)`)

// DownloadImage downloads an image from a URL and saves it to the specified path
func (c *Client) DownloadImage(ctx context.Context, imageURL, outputDir string) (string, error) {
	// Extract the image hash from the URL
	urlParts := strings.Split(imageURL, "/")
	if len(urlParts) < 1 {
		return "", fmt.Errorf("invalid image URL format: %s", imageURL)
	}

	imageHash := urlParts[len(urlParts)-1]
	if imageHash == "" {
		return "", fmt.Errorf("couldn't extract image hash from URL: %s", imageURL)
	}

	// Create a filename for the image
	imgDir := filepath.Join(outputDir, "static")
	if err := os.MkdirAll(imgDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create image directory: %w", err)
	}

	filename := filepath.Join(imgDir, imageHash+".png")

	// Check if the file already exists (to avoid redownloading)
	if _, err := os.Stat(filename); err == nil {
		if c.logger != nil {
			c.logger.Debug("Image already exists: %s", filename)
		}
		return "static/" + imageHash + ".png", nil
	}

	// Create HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", imageURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image, status: %s", resp.Status)
	}

	// Create the output file
	out, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Copy the response body to the output file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	if c.logger != nil {
		c.logger.Debug("Downloaded image to: %s", filename)
	}

	// Return the relative path for use in markdown
	return "static/" + imageHash + ".png", nil
}

// ProcessMarkdownImages processes a markdown string, downloads all images, and updates image references
func (c *Client) ProcessMarkdownImages(ctx context.Context, markdown, outputDir string) (string, error) {
	// Find all image references
	matches := imageURLPattern.FindAllStringSubmatch(markdown, -1)

	if len(matches) == 0 {
		// No images to process
		return markdown, nil
	}

	processedMarkdown := markdown

	// Process each image
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		altText := match[1]
		imageURL := match[2]
		originalReference := match[0]

		// Download the image
		localPath, err := c.DownloadImage(ctx, imageURL, outputDir)
		if err != nil {
			if c.logger != nil {
				c.logger.Error("Failed to download image %s: %v", imageURL, err)
			}
			continue
		}

		// Create the new markdown reference
		newReference := fmt.Sprintf("![%s](%s)", altText, localPath)

		// Replace the reference in the markdown
		processedMarkdown = strings.Replace(processedMarkdown, originalReference, newReference, 1)
	}

	return processedMarkdown, nil
}
