package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/epheo/anyblog/pkg/anytype"
	"github.com/epheo/anyblog/pkg/auth"
	"github.com/epheo/anyblog/pkg/display"
)

const defaultTimeout = 30 * time.Second

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// run encapsulates the main program logic
func run(ctx context.Context) error {
	// Get configuration
	config, err := auth.GetConfiguration()
	if err != nil {
		return err
	}

	// Create API client
	client := anytype.NewClient(config.ApiURL, config.SessionToken, config.AppKey)

	// Get spaces
	spacesData, err := client.GetSpaces(ctx)
	if err != nil {
		return fmt.Errorf("error getting spaces: %w", err)
	}

	// Parse spaces response
	var spacesResp anytype.SpacesResponse
	if err := json.Unmarshal(spacesData, &spacesResp); err != nil {
		return fmt.Errorf("error parsing spaces JSON: %w", err)
	}

	filters := map[string]interface{}{
		"types": []string{"ot-67f7210ccebba02cb2576fb2"},
		"tags":  []string{"published"},
		"query": "frfrefrefre",
	}
	// Retrieve the first available space ID
	if len(spacesResp.Data) == 0 {
		return fmt.Errorf("no spaces available")
	}
	// Look for a space with a specific name (prefer "My private space" if it exists)
	var spaceID string
	var spaceName string
	searchName := "epheo" // Default space to look for

	for _, space := range spacesResp.Data {
		if space.Name == searchName {
			spaceID = space.ID
			spaceName = space.Name
			fmt.Printf("Found exact match for space: %s\n", space.Name)
			break
		}
	}

	// If we didn't find the preferred space, use the first one
	if spaceID == "" {
		spaceID = spacesResp.Data[0].ID
		fmt.Printf("Using default space: %s (%s)\n", spacesResp.Data[0].Name, spaceID)
	}

	// Fetch available types
	typesData, err := client.GetTypes(ctx, spaceID)
	if err != nil {
		return fmt.Errorf("error getting types: %w", err)
	}

	// Parse types response
	var typesResp map[string]interface{}
	if err := json.Unmarshal(typesData, &typesResp); err != nil {
		return fmt.Errorf("error parsing types JSON: %w", err)
	}

	// Display the available types
	display.PrettyPrintJSON("Available Types", typesData)

	// Find a specific type by name (optional)
	typeName := "Page" // Change this to the type you're looking for
	typeDetails, err := client.GetTypeByName(ctx, spaceID, typeName)
	if err != nil {
		fmt.Printf("Could not find type '%s': %v\n", typeName, err)
	} else {
		display.PrettyPrintJSON(fmt.Sprintf("Type details for '%s'", typeName), []byte(typeDetails))
	}

	objects, err := client.SearchObjectsWithFilters(ctx, spaceID, filters)
	if err != nil {
		return err
	}

	// Pretty-print the successful JSON response
	display.PrettyPrintJSON(fmt.Sprintf("Objects from space %s", spaceName), objects)

	return nil
}
