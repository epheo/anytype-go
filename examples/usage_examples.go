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

	// Step 8: Working with properties and tags
	workWithProperties(ctx, client, spaceID)

	// Step 9: Searching for objects
	searchObjects(ctx, client, spaceID)

	// Step 10: Working with members
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
	for _, objType := range types.Data {
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

	if len(templates.Data) == 0 {
		fmt.Println("No templates found for this type")
		return ""
	}

	// Use the first template
	template := templates.Data[0]
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

// workWithTypes demonstrates creating custom types and objects of those types
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
		Properties: []anytype.PropertyLinkValue{
			anytype.NumberPropertyLinkValue{Key: "price", Number: 1299.99},
			anytype.TextPropertyLinkValue{Key: "description", Text: "High-performance gaming laptop with advanced graphics card."},
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
			for _, prop := range objectDetails.Object.Properties {
				fmt.Printf("    %s: %s\n", prop.Name, formatPropertyValue(prop))
			}
		}
	}

	// Clean up the object we created
	fmt.Println("\nCleaning up custom type object...")
	_, err = client.Space(spaceID).Object(objectID).Delete(ctx)
	if err != nil {
		log.Printf("Failed to delete object: %v", err)
	} else {
		fmt.Printf("Deleted object: %s\n", newObject.Object.Name)
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

	fmt.Printf("Found %d objects in space (%d total)\n", len(objects.Data), objects.Pagination.Total)

	// All() walks every page transparently; no manual offset loop needed.
	fmt.Println("Counting all objects across every page...")
	count := 0
	for _, err := range client.Space(spaceID).Objects().All(ctx) {
		if err != nil {
			log.Printf("Failed while iterating objects: %v", err)
			break
		}
		count++
	}
	fmt.Printf("Iterated %d objects in total\n", count)

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

	// Update the object
	fmt.Println("Updating the object...")
	updatedObject, err := client.Space(spaceID).Object(objectID).Update(ctx, anytype.UpdateObjectRequest{
		Name:     "Go SDK Example Page (Updated)",
		Markdown: "# Go SDK Example (Updated)\n\nThis page was updated using the Anytype Go SDK.\n\n## Features\n\n- Easy authentication\n- Space management\n- Object creation, update, and deletion\n- Search capabilities",
	})
	if err != nil {
		log.Printf("Failed to update object: %v", err)
	} else {
		fmt.Printf("Updated object name: %s\n", updatedObject.Object.Name)
	}

	// Get object as markdown
	fmt.Println("Getting object as markdown...")
	exportResp, err := client.Space(spaceID).Object(objectID).Get(ctx, anytype.WithFormat("md"))
	if err != nil {
		log.Printf("Failed to get object as markdown: %v", err)
	} else {
		fmt.Printf("Markdown export preview: %s\n", safeSubstring(exportResp.Object.Markdown, 100))
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
		_, err = client.Space(spaceID).Object(tempObjectID).Delete(ctx)
		if err != nil {
			log.Printf("Failed to delete object: %v", err)
		} else {
			fmt.Printf("Successfully deleted object: %s\n", tempObject.Object.Name)
		}
	}

	// Clean up the example page
	fmt.Println("Cleaning up example page...")
	_, err = client.Space(spaceID).Object(objectID).Delete(ctx)
	if err != nil {
		log.Printf("Failed to delete example page: %v", err)
	} else {
		fmt.Printf("Deleted object: %s\n", newObject.Object.Name)
	}

	return objectID
}

// workWithLists demonstrates operations with lists and views
func workWithLists(ctx context.Context, client anytype.Client, spaceID, objectID string) string {
	fmt.Println("\n=== Working with Lists and Views ===")

	// Create a collection for this demo
	fmt.Println("Creating a collection...")
	collectionResp, err := client.Space(spaceID).Objects().Create(ctx, anytype.CreateObjectRequest{
		TypeKey: "collection",
		Name:    "SDK Demo Collection",
		Icon: &anytype.Icon{
			Format: anytype.IconFormatEmoji,
			Emoji:  "📋",
		},
	})
	if err != nil {
		log.Printf("Failed to create collection: %v", err)
		return ""
	}

	listID := collectionResp.Object.ID
	fmt.Printf("Created collection: %s (ID: %s)\n", collectionResp.Object.Name, listID)

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

	// Clean up the collection
	fmt.Println("Cleaning up collection...")
	_, err = client.Space(spaceID).Object(listID).Delete(ctx)
	if err != nil {
		log.Printf("Failed to delete collection: %v", err)
	} else {
		fmt.Printf("Deleted collection: %s\n", collectionResp.Object.Name)
	}

	return listID
}

// workWithProperties demonstrates space-level property and tag management
func workWithProperties(ctx context.Context, client anytype.Client, spaceID string) {
	fmt.Println("\n=== Working with Properties and Tags ===")

	// List existing properties
	fmt.Println("Listing space properties...")
	properties, err := client.Space(spaceID).Properties().List(ctx)
	if err != nil {
		log.Printf("Failed to list properties: %v", err)
		return
	}
	fmt.Printf("Found %d properties in space\n", len(properties.Data))

	// Create a multi-select property with tags
	fmt.Println("Creating a multi-select property with tags...")
	propResp, err := client.Space(spaceID).Properties().Create(ctx, anytype.CreatePropertyRequest{
		Name:   "Priority",
		Format: "multi_select",
		Tags: []anytype.CreateTagRequest{
			{Name: "High", Color: "red"},
			{Name: "Medium", Color: "orange"},
			{Name: "Low", Color: "blue"},
		},
	})
	if err != nil {
		log.Printf("Failed to create property: %v", err)
		return
	}
	fmt.Printf("Created property: %s (Key: %s, Format: %s)\n",
		propResp.Property.Name, propResp.Property.Key, propResp.Property.Format)

	// Get the property details
	propDetails, err := client.Space(spaceID).Property(propResp.Property.ID).Get(ctx)
	if err != nil {
		log.Printf("Failed to get property details: %v", err)
	} else {
		fmt.Printf("Property details: %s (Format: %s)\n",
			propDetails.Property.Name, propDetails.Property.Format)
	}

	// List tags on this property
	fmt.Println("Listing tags on the property...")
	tags, err := client.Space(spaceID).Property(propResp.Property.ID).Tags().List(ctx)
	if err != nil {
		log.Printf("Failed to list tags: %v", err)
	} else {
		fmt.Printf("Found %d tags\n", len(tags.Data))
		for _, tag := range tags.Data {
			fmt.Printf("  - %s (Color: %s)\n", tag.Name, tag.Color)
		}
	}

	// Add a new tag to the property
	fmt.Println("Adding a new tag...")
	newTag, err := client.Space(spaceID).Property(propResp.Property.ID).Tags().Create(ctx, anytype.CreateTagRequest{
		Name:  "Critical",
		Color: "red",
	})
	if err != nil {
		log.Printf("Failed to create tag: %v", err)
	} else {
		fmt.Printf("Created tag: %s (Color: %s)\n", newTag.Tag.Name, newTag.Tag.Color)

		// Clean up the tag
		_, err = client.Space(spaceID).Property(propResp.Property.ID).Tag(newTag.Tag.ID).Delete(ctx)
		if err != nil {
			log.Printf("Failed to delete tag: %v", err)
		} else {
			fmt.Printf("Deleted tag: %s\n", newTag.Tag.Name)
		}
	}

	// Clean up the property
	fmt.Println("Cleaning up property...")
	_, err = client.Space(spaceID).Property(propResp.Property.ID).Delete(ctx)
	if err != nil {
		log.Printf("Failed to delete property: %v", err)
	} else {
		fmt.Printf("Deleted property: %s\n", propResp.Property.Name)
	}
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

// formatPropertyValue returns a human-readable string for a property based on its format
func formatPropertyValue(prop anytype.Property) string {
	switch prop.Format {
	case "text":
		return prop.Text
	case "number":
		return fmt.Sprintf("%g", prop.Number)
	case "checkbox":
		return fmt.Sprintf("%v", prop.Checkbox)
	case "date":
		return prop.Date
	case "url":
		return prop.URL
	case "email":
		return prop.Email
	case "phone":
		return prop.Phone
	case "objects":
		if len(prop.Objects) == 0 {
			return "[]"
		}
		return fmt.Sprintf("%v", prop.Objects)
	case "multi_select":
		names := make([]string, len(prop.MultiSelect))
		for i, tag := range prop.MultiSelect {
			names[i] = tag.Name
		}
		return fmt.Sprintf("%v", names)
	case "select":
		if prop.Select != nil {
			return prop.Select.Name
		}
		return ""
	default:
		return fmt.Sprintf("(%s)", prop.Format)
	}
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
