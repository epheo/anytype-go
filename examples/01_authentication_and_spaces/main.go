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
	// Create an auth manager
	authManager := auth.NewAuthManager(
		auth.WithAPIURL(""),
		auth.WithNonInteractive(false),
		auth.WithSilent(false),
	)

	// Get client directly from auth manager (handles authentication internally)
	client, err := authManager.GetClient(
		anytype.WithDebug(false),
		anytype.WithTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("Successfully authenticated with Anytype!")

	// Create context with timeout for operations
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get all available spaces
	spaces, err := client.GetSpaces(ctx)
	if err != nil {
		log.Fatalf("Failed to get spaces: %v", err)
	}

	// Display available spaces
	fmt.Printf("Found %d spaces:\n", len(spaces.Data))
	for i, space := range spaces.Data {
		fmt.Printf("%d. %s (ID: %s)\n", i+1, space.Name, space.ID)
	}

	// Select a space to work with
	var targetSpace *anytype.Space
	if len(spaces.Data) > 0 {
		targetSpace = &spaces.Data[0]
		fmt.Printf("Using space: %s\n", targetSpace.Name)

		// Display space members if available
		if len(targetSpace.Members) > 0 {
			fmt.Printf("Space has %d members\n", len(targetSpace.Members))
			for i, member := range targetSpace.Members {
				fmt.Printf("  Member %d: %s (ID: %s)\n", i+1, member.Name, member.ID)
			}
		}
	} else {
		log.Fatal("No spaces available")
	}
}
