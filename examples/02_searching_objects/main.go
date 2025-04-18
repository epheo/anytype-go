package main

import (
	"context"
	"fmt"
	"log"
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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
	fmt.Printf("Using space: %s (%s)\n", spaces.Data[0].Name, spaceID)

	// Demonstrate different search options
	fmt.Println("\n=== Basic Search ===")
	basicSearchParams := &anytype.SearchParams{
		Query: "project meeting", // Text search
		Limit: 10,
	}
	basicResults, err := client.Search(ctx, spaceID, basicSearchParams)
	if err != nil {
		log.Fatalf("Basic search failed: %v", err)
	}
	fmt.Printf("Found %d objects matching 'project meeting'\n", len(basicResults.Data))
	printObjectSummary(basicResults.Data)

	fmt.Println("\n=== Search by Type ===")
	typeSearchParams := &anytype.SearchParams{
		Types: []string{"ot-note", "ot-task"},
		Limit: 10,
	}
	typeResults, err := client.Search(ctx, spaceID, typeSearchParams)
	if err != nil {
		log.Fatalf("Type search failed: %v", err)
	}
	fmt.Printf("Found %d notes and tasks\n", len(typeResults.Data))
	printObjectSummary(typeResults.Data)

	fmt.Println("\n=== Search by Tags ===")
	tagSearchParams := &anytype.SearchParams{
		Tags:  []string{"important", "work"},
		Limit: 10,
	}
	tagResults, err := client.Search(ctx, spaceID, tagSearchParams)
	if err != nil {
		log.Fatalf("Tag search failed: %v", err)
	}
	fmt.Printf("Found %d objects with tags 'important' or 'work'\n", len(tagResults.Data))
	printObjectSummary(tagResults.Data)

	fmt.Println("\n=== Combined Search ===")
	combinedSearchParams := &anytype.SearchParams{
		Query: "project",
		Types: []string{"ot-note"},
		Tags:  []string{"important"},
		Limit: 10,
	}
	combinedResults, err := client.Search(ctx, spaceID, combinedSearchParams)
	if err != nil {
		log.Fatalf("Combined search failed: %v", err)
	}
	fmt.Printf("Found %d important note objects related to projects\n", len(combinedResults.Data))
	printObjectSummary(combinedResults.Data)

	fmt.Println("\n=== Using Query Builder ===")
	qb := client.NewQueryBuilder(spaceID)
	qb.WithTypeKeys("ot-note")
	qb.WithTag("important")
	qb.WithQuery("project")
	qb.WithLimit(10)
	qb.WithSortField("name", true)

	qbResults, err := qb.Execute(ctx)
	if err != nil {
		log.Fatalf("Query builder search failed: %v", err)
	}
	fmt.Printf("Found %d results with query builder\n", len(qbResults.Data))
	printObjectSummary(qbResults.Data)

	fmt.Println("\n=== Using ExecuteWithCallback for Processing ===")
	processedCount := 0
	err = qb.ExecuteWithCallback(ctx, func(obj anytype.Object) error {
		processedCount++
		fmt.Printf("Processing object: %s (%s)\n", obj.Name, obj.ID)
		// Do something with each object here
		return nil
	})
	if err != nil {
		log.Fatalf("Callback processing failed: %v", err)
	}
	fmt.Printf("Processed %d objects with callback\n", processedCount)
}

// Helper function to print a summary of objects
func printObjectSummary(objects []anytype.Object) {
	if len(objects) == 0 {
		fmt.Println("  No objects found")
		return
	}

	for i, obj := range objects {
		if i >= 5 {
			fmt.Printf("  ... and %d more\n", len(objects)-5)
			break
		}
		typeName := "Unknown"
		if obj.Type != nil {
			typeName = obj.Type.Name
		}
		fmt.Printf("  %d. %s (Type: %s, ID: %s)\n", i+1, obj.Name, typeName, obj.ID)
		if len(obj.Tags) > 0 {
			fmt.Printf("     Tags: %v\n", obj.Tags)
		}
	}
}
