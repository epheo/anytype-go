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
	spaceName := spaces.Data[0].Name
	fmt.Printf("Using space: %s (%s)\n", spaceName, spaceID)

	// Create a new note
	fmt.Println("\n=== Creating a New Note ===")
	newNote := &anytype.Object{
		Name: fmt.Sprintf("Meeting Notes - %s", time.Now().Format("Jan 02, 2006")),
		Type: &anytype.TypeInfo{
			Key:  "ot-note",
			Name: "Note",
		},
		Icon: &anytype.Icon{
			Format: "emoji",
			Emoji:  "📝",
		},
		Tags: []string{"meeting", "example"},
		// You can add more fields as needed
	}

	createdNote, err := client.CreateObject(ctx, spaceID, newNote)
	if err != nil {
		log.Fatalf("Failed to create note: %v", err)
	}

	fmt.Printf("Created note: %s (ID: %s)\n", createdNote.Name, createdNote.ID)
	fmt.Printf("Type: %s, Tags: %v\n", createdNote.Type.Name, createdNote.Tags)

	// Create a new task
	fmt.Println("\n=== Creating a New Task ===")
	newTask := &anytype.Object{
		Name: "Follow up on meeting action items",
		Type: &anytype.TypeInfo{
			Key:  "ot-task",
			Name: "Task",
		},
		Icon: &anytype.Icon{
			Format: "emoji",
			Emoji:  "✅",
		},
		Tags: []string{"work", "priority"},
		// Add task-specific properties if needed
	}

	createdTask, err := client.CreateObject(ctx, spaceID, newTask)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}

	fmt.Printf("Created task: %s (ID: %s)\n", createdTask.Name, createdTask.ID)
	fmt.Printf("Type: %s, Tags: %v\n", createdTask.Type.Name, createdTask.Tags)

	// Get the created note to verify
	fmt.Println("\n=== Getting Created Note ===")
	params := &anytype.GetObjectParams{
		SpaceID:  spaceID,
		ObjectID: createdNote.ID,
	}

	retrievedObject, err := client.GetObject(ctx, params)
	if err != nil {
		log.Fatalf("Failed to get object: %v", err)
	}

	fmt.Printf("Retrieved object: %s (ID: %s)\n", retrievedObject.Name, retrievedObject.ID)
	fmt.Printf("Type: %s, Tags: %v\n", retrievedObject.Type.Name, retrievedObject.Tags)

	// Update the note
	fmt.Println("\n=== Updating the Note ===")
	updateObj := &anytype.Object{
		Name: fmt.Sprintf("Updated Meeting Notes - %s", time.Now().Format("Jan 02, 2006")),
		Icon: &anytype.Icon{
			Format: "emoji",
			Emoji:  "📌",
		},
		Tags: []string{"meeting", "example", "updated"},
	}

	updatedObject, err := client.UpdateObject(ctx, spaceID, createdNote.ID, updateObj)
	if err != nil {
		log.Fatalf("Failed to update object: %v", err)
	}

	fmt.Printf("Updated object: %s (ID: %s)\n", updatedObject.Name, updatedObject.ID)
	fmt.Printf("Type: %s, Tags: %v\n", updatedObject.Type.Name, updatedObject.Tags)

	// Delete the task (uncomment if you want to actually delete)
	/*
		fmt.Println("\n=== Deleting the Task ===")
		err = client.DeleteObject(ctx, spaceID, createdTask.ID)
		if err != nil {
			log.Fatalf("Failed to delete object: %v", err)
		}
		fmt.Println("Task successfully deleted")
	*/

	// Instead, just show how to call the delete method
	fmt.Println("\n=== How to Delete an Object ===")
	fmt.Printf("// To delete an object, use:\n")
	fmt.Printf("err := client.DeleteObject(ctx, \"%s\", \"%s\")\n", spaceID, createdTask.ID)
	fmt.Printf("if err != nil {\n\tlog.Fatalf(\"Failed to delete object: %%v\", err)\n}\n")
}
