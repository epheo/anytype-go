package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/epheo/anytype-go"
	_ "github.com/epheo/anytype-go/client" // Register client implementation
)

func main() {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Step 1: Authentication
	client := authenticate(ctx)
	if client == nil {
		return
	}

	// Step 2: Working with spaces
	spaceID := os.Getenv("ANYTYPE_SPACE_ID")
	// if Space ID isn't set, use one of the spaces available
	if spaceID == "" {
		fmt.Println("Space ID not set, using first space available")
		spaceID = workWithSpaces(ctx, client)
		if spaceID == "" {
			return
		}
	}

	// Step 3: Working with types (creating new types)
	err := workWithTypes(ctx, client, spaceID)
	if err != nil {
		log.Printf("Failed to work with types: %v", err)
		return
	}

	// Step 4: Working with object types
	typeID := workWithObjectTypes(ctx, client, spaceID)
	if typeID == "" {
		return
	}

	// Step 5: Working with templates
	templateID := workWithTemplates(ctx, client, spaceID, typeID)

	// Step 6: Working with objects
	objectID := workWithObjects(ctx, client, spaceID, templateID)
	if objectID == "" {
		return
	}

	// Step 7: Working with lists and views
	listID := workWithLists(ctx, client, spaceID, objectID)
	if listID == "" {
		return
	}

	// Step 8: Searching for objects
	searchObjects(ctx, client, spaceID)

	// Step 9: Working with members
	workWithMembers(ctx, client, spaceID)

	fmt.Println("\nSuccessfully completed all example operations!")
}

// authenticate demonstrates the authentication flow
func authenticate(ctx context.Context) anytype.Client {
	fmt.Println("=== Authentication ===")

	// Create an unauthenticated client for authentication
	client := anytype.NewClient(
		anytype.WithBaseURL("http://localhost:31009"), // Default Anytype local API URL (without /v1)
	)

	// Perform interactive authentication
	fmt.Println("Starting authentication flow...")

	// Step 1: Initiate authentication and get challenge ID
	fmt.Println("Initiating authentication flow...")
	authResponse, err := client.Auth().CreateChallenge(ctx, "GoSDKExample")
	if err != nil {
		log.Printf("Failed to initiate authentication: %v", err)
		return nil
	}
	challengeID := authResponse.ChallengeID

	// Step 2: User needs to enter the code shown in Anytype app
	fmt.Println("Please check your Anytype app and enter the displayed verification code:")
	var code string
	fmt.Scanln(&code)

	// Step 3: Complete authentication and get tokens
	tokenResponse, err := client.Auth().CreateApiKey(ctx, challengeID, code)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		return nil
	}

	fmt.Printf("Authentication successful! ApiKey: %s\n",
		safeSubstring(tokenResponse.ApiKey, 10))

	// Create a new authenticated client with the tokens
	client = anytype.NewClient(
		anytype.WithBaseURL("http://localhost:31009"),
		anytype.WithAppKey(tokenResponse.ApiKey),
	)

	return client
}

// workWithSpaces demonstrates operations with spaces
func workWithSpaces(ctx context.Context, client anytype.Client) string {
	fmt.Println("\n=== Working with Spaces ===")

	// List all spaces
	fmt.Println("Listing all spaces...")
	spacesResp, err := client.Spaces().List(ctx)
	if err != nil {
		log.Printf("Failed to list spaces: %v", err)
		return ""
	}

	if len(spacesResp.Data) == 0 {
		fmt.Println("No spaces found. Creating a new space...")

		// Create a new space
		newSpace, err := client.Spaces().Create(ctx, anytype.CreateSpaceRequest{
			Name:        "My Go SDK Space",
			Description: "A space created via the Go SDK example",
		})
		if err != nil {
			log.Printf("Failed to create space: %v", err)
			return ""
		}

		fmt.Printf("Created new space: %s (ID: %s)\n", newSpace.Space.Name, newSpace.Space.ID)
		return newSpace.Space.ID
	}

	// Use the first available space
	space := spacesResp.Data[0]
	fmt.Printf("Using existing space: %s (ID: %s)\n", space.Name, space.ID)

	// Get details for a specific space
	spaceDetails, err := client.Space(space.ID).Get(ctx)
	if err != nil {
		log.Printf("Failed to get space details: %v", err)
		return ""
	}

	fmt.Printf("Space details: Name=%s, Description=%s\n",
		spaceDetails.Space.Name,
		spaceDetails.Space.Description)

	return space.ID
}

// workWithObjectTypes demonstrates operations with object types
func workWithObjectTypes(ctx context.Context, client anytype.Client, spaceID string) string {
	fmt.Println("\n=== Working with Object Types ===")

	// List all object types
	fmt.Println("Listing object types...")
	types, err := client.Space(spaceID).Types().List(ctx)
	if err != nil {
		log.Printf("Failed to list object types: %v", err)
		return ""
	}

	// Find a page type
	var pageType *anytype.Type
	for _, objType := range types {
		fmt.Printf("Found type: %s (Key: %s)\n", objType.Name, objType.Key)
		if objType.Key == "page" {
			pageType = &objType
			break
		}
	}

	if pageType == nil {
		log.Printf("Could not find 'page' type")
		return ""
	}

	// Display the type information we already have from the List() call
	fmt.Printf("Page type details: Name=%s, Key=%s, Layout=%s\n",
		pageType.Name,
		pageType.Key,
		pageType.Layout)

	return pageType.Key
}

// workWithTemplates demonstrates operations with templates
func workWithTemplates(ctx context.Context, client anytype.Client, spaceID, typeKey string) string {
	fmt.Println("\n=== Working with Templates ===")

	// List templates for the type
	fmt.Println("Listing templates for the page type...")
	templates, err := client.Space(spaceID).Type(typeKey).Templates().List(ctx)
	if err != nil {
		log.Printf("Failed to list templates: %v", err)
		return ""
	}

	if len(templates) == 0 {
		fmt.Println("No templates found for this type")
		return ""
	}

	// Use the first template
	template := templates[0]
	fmt.Printf("Found template: %s (ID: %s)\n", template.Name, template.ID)

	// Get detailed template information
	templateDetails, err := client.Space(spaceID).Type(typeKey).Template(template.ID).Get(ctx)
	if err != nil {
		log.Printf("Failed to get template details: %v", err)
		return ""
	}

	fmt.Printf("Template details: Name=%s\n", templateDetails.Template.Name)

	return template.ID
}

// workWithObjects demonstrates operations with objects
func workWithTypes(ctx context.Context, client anytype.Client, spaceID string) error {
	fmt.Println("\n=== Working with Types ===")

	// Create a new type
	fmt.Println("Creating a new type...")
	createTypeReq := anytype.CreateTypeRequest{
		Name:   "Product",
		Layout: "basic",
		Icon: &anytype.Icon{
			Format: anytype.IconFormatEmoji,
			Emoji:  "📦",
		},
		PluralName: "Products",
		Properties: []anytype.PropertyDefinition{
			{
				Key:    "price",
				Name:   "Price",
				Format: "number",
			},
			{
				Key:    "category",
				Name:   "Category",
				Format: "text",
			},
			{
				Key:    "description",
				Name:   "Description",
				Format: "text",
			},
		},
	}

	newType, err := client.Space(spaceID).Types().Create(ctx, createTypeReq)
	if err != nil {
		return fmt.Errorf("failed to create type: %w", err)
	}

	fmt.Printf("Created type: %s (Key: %s)\n", newType.Type.Name, newType.Type.Key)
	fmt.Printf("Type layout: %s\n", newType.Type.Layout)
	fmt.Printf("Type has %d properties\n", len(newType.Type.Properties))

	for _, prop := range newType.Type.Properties {
		fmt.Printf("  - Property: %s (%s) - Format: %s\n", prop.Name, prop.Key, prop.Format)
	}

	// Create an object of the custom type
	fmt.Println("\nCreating object of custom type...")
	createReq := anytype.CreateObjectRequest{
		TypeKey: newType.Type.Key,
		Name:    "Gaming Laptop",
		Icon: &anytype.Icon{
			Format: anytype.IconFormatEmoji,
			Emoji:  "💻",
		},
		Properties: []anytype.PropertyLinkWithValue{
			{
				Key:    "price",
				Number: func(v float64) *float64 { return &v }(1299.99),
			},
			// Note: For multi-select, we'd need to create tags first or use existing tag IDs
			// Skipping the tag field for now since it requires pre-existing tag options
			{
				Key:  "description",
				Text: func(v string) *string { return &v }("High-performance gaming laptop with advanced graphics card."),
			},
		},
	}

	newObject, err := client.Space(spaceID).Objects().Create(ctx, createReq)
	if err != nil {
		log.Printf("Failed to create object: %v", err)
		return nil
	}

	objectID := newObject.Object.ID
	fmt.Printf("Created object: %s (ID: %s)\n", newObject.Object.Name, objectID)
	fmt.Printf("Object type: %s\n", newObject.Object.Type.Name)

	objectDetails, err := client.Space(spaceID).Object(objectID).Get(ctx)
	if err != nil {
		log.Printf("Failed to get object details: %v", err)
	} else {
		fmt.Printf("Object details:\n")
		fmt.Printf("  Name: %s\n", objectDetails.Object.Name)
		fmt.Printf("  Type: %s\n", objectDetails.Object.Type.Name)
		if objectDetails.Object.Properties != nil {
			fmt.Printf("  Properties:\n")
			for _, value := range objectDetails.Object.Properties {
				fmt.Printf("    %s: %+v\n", value.Key, value)
			}
		}
	}

	return nil
}

func workWithObjects(ctx context.Context, client anytype.Client, spaceID, templateID string) string {
	fmt.Println("\n=== Working with Objects ===")

	// List existing objects
	fmt.Println("Listing existing objects...")
	objects, err := client.Space(spaceID).Objects().List(ctx)
	if err != nil {
		log.Printf("Failed to list objects: %v", err)
		return ""
	}

	fmt.Printf("Found %d objects in space\n", len(objects))

	// Create a new object
	fmt.Println("Creating a new page object...")

	createReq := anytype.CreateObjectRequest{
		TypeKey: "page", // Using the known type key for pages
		Name:    "Go SDK Example Page",
		Body:    "# Go SDK Example\n\nThis page was created using the Anytype Go SDK.\n\n## Features\n\n- Easy authentication\n- Space management\n- Object creation and manipulation\n- Search capabilities",
	}

	// Use template if available
	if templateID != "" {
		createReq.TemplateID = templateID
	}

	// Set an emoji icon
	createReq.Icon = &anytype.Icon{
		Format: anytype.IconFormatEmoji,
		Emoji:  "🚀",
	}

	newObject, err := client.Space(spaceID).Objects().Create(ctx, createReq)
	if err != nil {
		log.Printf("Failed to create object: %v", err)
		return ""
	}

	objectID := newObject.Object.ID
	fmt.Printf("Created new object: %s (ID: %s)\n", newObject.Object.Name, objectID)

	// Get object details
	objectDetails, err := client.Space(spaceID).Object(objectID).Get(ctx)
	if err != nil {
		log.Printf("Failed to get object details: %v", err)
		return ""
	}

	fmt.Printf("Object details: Name=%s, Type=%s\n",
		objectDetails.Object.Name,
		objectDetails.Object.Type.Name)

	// Export the object to markdown
	fmt.Println("Exporting object to markdown...")
	exportResp, err := client.Space(spaceID).Object(objectID).Export(ctx, "markdown")
	if err != nil {
		log.Printf("Failed to export object: %v", err)
	} else {
		fmt.Printf("Markdown export preview: %s\n", previewText(exportResp.Markdown, 100))
	}

	// Create a temporary object to demonstrate deletion
	fmt.Println("Creating a temporary object to demonstrate deletion...")
	tempObject, err := client.Space(spaceID).Objects().Create(ctx, anytype.CreateObjectRequest{
		TypeKey: "page",
		Name:    "Temporary Object for Deletion Demo",
		Icon: &anytype.Icon{
			Format: anytype.IconFormatEmoji,
			Emoji:  "📝",
		},
	})
	if err != nil {
		log.Printf("Failed to create temporary object: %v", err)
	} else {
		tempObjectID := tempObject.Object.ID
		fmt.Printf("Created temporary object: %s (ID: %s)\n", tempObject.Object.Name, tempObjectID)

		// Delete the temporary object
		fmt.Printf("Deleting temporary object...\n")
		deletedObject, err := client.Space(spaceID).Object(tempObjectID).Delete(ctx)
		if err != nil {
			log.Printf("Failed to delete object: %v", err)
		} else {
			fmt.Printf("Successfully deleted object: %s (archived status: %v)\n",
				deletedObject.Object.Name,
				deletedObject.Object.Archived)
		}
	}

	return objectID
}

// workWithLists demonstrates operations with lists and views
func workWithLists(ctx context.Context, client anytype.Client, spaceID, objectID string) string {
	fmt.Println("\n=== Working with Lists and Views ===")

	// For this example, we'll just show how to list views for a list object
	// In a real application, you'd likely search for a collection or create one

	fmt.Println("Finding a list in the current space...")
	searchResp, err := client.Space(spaceID).Search(ctx, anytype.SearchRequest{
		Types: []string{"collection"}, // Search for objects of type: Collection
	})
	if err != nil || len(searchResp.Data) == 0 {
		log.Printf("No collection objects found for demonstration")
		return ""
	}

	// Find a collection object
	var listID string
	for _, obj := range searchResp.Data {
		// Could add more sophisticated detection here
		if obj.Layout == "collection" || obj.Layout == "list" {
			listID = obj.ID
			break
		}
	}

	if listID == "" {
		log.Printf("No suitable list object found")
		return ""
	}

	fmt.Printf("Found list object with ID: %s\n", listID)

	// Get views for the list
	viewsResp, err := client.Space(spaceID).List(listID).Views().List(ctx)
	if err != nil {
		log.Printf("Failed to get list views: %v", err)
		return ""
	}

	fmt.Printf("Found %d views for list\n", len(viewsResp.Data))
	for i, view := range viewsResp.Data {
		fmt.Printf("  View %d: %s (Layout: %s)\n", i+1, view.Name, view.Layout)
	}

	// If we have a view, let's get objects using that view
	if len(viewsResp.Data) > 0 {
		viewID := viewsResp.Data[0].ID
		fmt.Printf("Getting objects using view: %s\n", viewsResp.Data[0].Name)

		objectsResp, err := client.Space(spaceID).List(listID).View(viewID).Objects().List(ctx)
		if err != nil {
			log.Printf("Failed to get objects in list view: %v", err)
		} else {
			fmt.Printf("View contains %d objects\n", len(objectsResp.Data))
		}

		// Add the example object to this list if possible
		fmt.Printf("Attempting to add our example object to the list...\n")
		err = client.Space(spaceID).List(listID).Objects().Add(ctx, []string{objectID})
		if err != nil {
			log.Printf("Could not add object to list: %v", err)
		} else {
			fmt.Printf("Successfully added object to list\n")

			// Now demonstrate removing the object from the list
			fmt.Printf("Removing object from list to demonstrate full API coverage...\n")
			err = client.Space(spaceID).List(listID).Object(objectID).Remove(ctx)
			if err != nil {
				log.Printf("Could not remove object from list: %v", err)
			} else {
				fmt.Printf("Successfully removed object from list\n")
			}
		}
	}

	return listID
}

// searchObjects demonstrates the search functionality
func searchObjects(ctx context.Context, client anytype.Client, spaceID string) {
	fmt.Println("\n=== Searching for Objects ===")

	// Basic search in current space
	fmt.Println("Performing a basic search...")
	basicSearchResp, err := client.Space(spaceID).Search(ctx, anytype.SearchRequest{
		Query: "example", // Search for objects containing "example"
	})
	if err != nil {
		log.Printf("Failed to perform basic search: %v", err)
		return
	}

	fmt.Printf("Found %d objects containing 'example'\n", len(basicSearchResp.Data))

	// Advanced search with sorting and type filtering
	fmt.Println("Performing an advanced search...")
	advancedSearchResp, err := client.Space(spaceID).Search(ctx, anytype.SearchRequest{
		Query: "SDK",
		Sort: &anytype.SortOptions{
			Property:  anytype.SortPropertyLastModifiedDate,
			Direction: anytype.SortDirectionDesc,
		},
		Types: []string{"page"}, // Only search for pages
	})
	if err != nil {
		log.Printf("Failed to perform advanced search: %v", err)
		return
	}

	fmt.Printf("Found %d page objects containing 'SDK', sorted by last modified date\n",
		len(advancedSearchResp.Data))

	// Perform a global search across all spaces
	fmt.Println("Performing a global search across all spaces...")
	searchClient := client.Search() // Get the SearchClient first
	globalSearchResp, err := searchClient.Search(ctx, anytype.SearchRequest{
		Query: "important",
	})
	if err != nil {
		log.Printf("Failed to perform global search: %v", err)
		return
	}

	fmt.Printf("Found %d objects containing 'important' across all spaces\n",
		len(globalSearchResp.Data))
}

// workWithMembers demonstrates operations with space members
func workWithMembers(ctx context.Context, client anytype.Client, spaceID string) {
	fmt.Println("\n=== Working with Members ===")

	// List members in the space
	fmt.Println("Listing space members...")
	membersResp, err := client.Space(spaceID).Members().List(ctx)
	if err != nil {
		log.Printf("Failed to list space members: %v", err)
		return
	}

	fmt.Printf("Found %d members in space\n", len(membersResp.Data))
	for i, member := range membersResp.Data {
		fmt.Printf("  Member %d: %s (Role: %s, Status: %s)\n",
			i+1,
			member.Name,
			member.Role,
			member.Status)
	}

	// If there are members, get details for one
	if len(membersResp.Data) > 0 {
		memberID := membersResp.Data[0].ID
		fmt.Printf("Getting details for member: %s\n", memberID)

		memberDetails, err := client.Space(spaceID).Member(memberID).Get(ctx)
		if err != nil {
			log.Printf("Failed to get member details: %v", err)
			return
		}

		fmt.Printf("Member details:\n")
		fmt.Printf("  Name: %s\n", memberDetails.Member.Name)
		fmt.Printf("  Global Name: %s\n", memberDetails.Member.GlobalName)
		fmt.Printf("  Identity: %s\n", memberDetails.Member.Identity)
		fmt.Printf("  Role: %s\n", memberDetails.Member.Role)
		fmt.Printf("  Status: %s\n", memberDetails.Member.Status)
		if memberDetails.Member.Icon != nil {
			fmt.Printf("  Icon Format: %s\n", memberDetails.Member.Icon.Format)
		}
	}
}

// previewText returns a preview of a longer text string
func previewText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

// safeSubstring returns a substring safely handling empty strings
func safeSubstring(s string, maxLen int) string {
	if len(s) == 0 {
		return "<empty>"
	}
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
