package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/epheo/anytype-go/pkg/anytype"
	"github.com/epheo/anytype-go/pkg/auth"
)

func main() {
	// Create auth manager and get client
	authManager := auth.NewAuthManager()
	client, err := authManager.GetClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Get spaces
	spaces, err := client.GetSpaces(ctx)
	if err != nil {
		log.Fatalf("Failed to get spaces: %v", err)
	}

	// Select the first space
	if len(spaces.Data) == 0 {
		log.Fatal("No spaces available")
	}
	spaceID := spaces.Data[0].ID
	spaceName := spaces.Data[0].Name
	fmt.Printf("Using space: %s (%s)\n", spaceName, spaceID)

	// Create export directory
	exportDir := filepath.Join(".", "exports", "example_exports")
	err = os.MkdirAll(exportDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create export directory: %v", err)
	}
	fmt.Printf("Created export directory: %s\n", exportDir)

	// Search for objects to export
	fmt.Println("\n=== Finding Objects to Export ===")
	searchParams := &anytype.SearchParams{
		Types: []string{"ot-note", "ot-page"},
		Limit: 10,
	}

	results, err := client.Search(ctx, spaceID, searchParams)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}
	fmt.Printf("Found %d objects to export\n", len(results.Data))

	// Export a single object (if any found)
	if len(results.Data) > 0 {
		fmt.Println("\n=== Exporting a Single Object ===")
		objectToExport := results.Data[0]
		fmt.Printf("Exporting object: %s (ID: %s)\n", objectToExport.Name, objectToExport.ID)

		filePath, err := client.ExportObject(ctx, spaceID, objectToExport.ID, exportDir, "md")
		if err != nil {
			log.Fatalf("Failed to export object: %v", err)
		}

		fmt.Printf("Object exported to: %s\n", filePath)
		fmt.Printf("File content preview:\n")
		printFilePreview(filePath)
	}

	// Export multiple objects
	if len(results.Data) > 1 {
		fmt.Println("\n=== Exporting Multiple Objects ===")
		fmt.Printf("Exporting %d objects to %s in markdown format\n", len(results.Data), exportDir)

		exportedFiles, err := client.ExportObjects(ctx, spaceID, results.Data, exportDir, "md")
		if err != nil {
			log.Fatalf("Failed to export objects: %v", err)
		}

		fmt.Printf("Successfully exported %d objects:\n", len(exportedFiles))
		for i, file := range exportedFiles {
			if i < 3 {
				fmt.Printf("  %d. %s\n", i+1, file)
			} else if i == 3 {
				fmt.Printf("  ... and %d more files\n", len(exportedFiles)-3)
				break
			}
		}
	}

	// Export as HTML (if supported)
	if len(results.Data) > 0 {
		fmt.Println("\n=== Exporting as HTML ===")
		htmlExportDir := filepath.Join(exportDir, "html")
		err = os.MkdirAll(htmlExportDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create HTML export directory: %v", err)
		}

		filePath, err := client.ExportObject(ctx, spaceID, results.Data[0].ID, htmlExportDir, "html")
		if err != nil {
			fmt.Printf("Note: HTML export failed: %v (this format may not be supported)\n", err)
		} else {
			fmt.Printf("HTML export successful: %s\n", filePath)
		}
	}
}

// Helper function to print a preview of an exported file
func printFilePreview(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("  Error reading file: %v\n", err)
		return
	}

	// Print first few lines or characters with a limit
	const previewLength = 200
	preview := string(content)
	if len(preview) > previewLength {
		preview = preview[:previewLength] + "..."
	}

	fmt.Printf("----- File Preview -----\n")
	fmt.Println(preview)
	fmt.Printf("------------------------\n")
}
