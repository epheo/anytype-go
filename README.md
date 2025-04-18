# Anytype-Go

[![Go Report Card](https://goreportcard.com/badge/github.com/epheo/anytype-go)](https://goreportcard.com/report/github.com/epheo/anytype-go)
[![GoDoc](https://godoc.org/github.com/epheo/anytype-go?status.svg)](https://godoc.org/github.com/epheo/anytype-go)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A Go SDK and command-line tool for interacting with the [Anytype](https://anytype.io) API to manage spaces, objects, and perform searches.

## 📑 Table of Contents

- [Overview](#-overview)
- [Installation](#-installation)
  - [Library Installation](#library-installation)
  - [CLI Installation](#cli-installation)
- [Getting Started](#-getting-started)
  - [Authentication](#authentication)
  - [Working with Spaces](#working-with-spaces)
  - [Searching for Objects](#searching-for-objects)
  - [Working with Objects](#working-with-objects)
  - [Exporting Objects](#exporting-objects)
- [Advanced Usage](#-advanced-usage)
  - [Query Builder](#query-builder)
  - [Error Handling](#error-handling)
  - [Best Practices](#best-practices)
  - [Troubleshooting](#troubleshooting)
- [CLI Tool](#-cli-tool)
  - [Configuration](#configuration)
  - [Command Line Usage](#command-line-usage)
  - [Command Line Options](#command-line-options)
- [API Reference](#-api-reference)
  - [Available Packages](#available-packages)
  - [Core API Components](#core-api-components)
- [Project Info](#-project-info)
  - [Features](#features)
  - [Project Structure](#-project-structure)
  - [Versioning Policy](#-versioning-policy)
  - [Building from Source](#building-from-source)
  - [Contributing](#contributing)
  - [License](#license)

## 📖 Overview

Anytype-Go provides both a Go library and a command-line interface for interacting with Anytype's local API. It enables you to:

- Search for objects within spaces using various filtering criteria
- Create, read, update, and delete objects
- Explore available spaces and their members
- Export objects to different formats (Markdown, HTML)
- Work with objects using custom types and tags
- Automate Anytype operations via scripts or applications

## 📥 Installation

### Library Installation

To use anytype-go as a package in your Go project:

```bash
go get github.com/epheo/anytype-go
```

### CLI Installation

To install the command-line tool:

```bash
# Install directly from GitHub
go install github.com/epheo/anytype-go/cmd/anytype-go@latest

# Or build from source
git clone https://github.com/epheo/anytype-go.git
cd anytype-go
go install ./cmd/anytype-go
```

## 🚀 Getting Started

This section provides a step-by-step guide to using the Anytype-Go library in your projects.

### Authentication

The Anytype API requires authentication. Here's how to set it up in your code:

```go
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
    authManager := auth.NewManager(auth.ManagerOptions{
        ConfigPath: "~/.config/anytype-go/anytype_auth.json",
    })

    // Get or create auth config (this will handle the authentication flow if needed)
    authConfig, err := authManager.GetOrCreateConfig()
    if err != nil {
        log.Fatalf("Failed to get or create auth config: %v", err)
    }

    // Create an API client with the auth config
    client, err := anytype.FromAuthConfig(authConfig, 
        anytype.WithDebug(false),
        anytype.WithTimeout(30 * time.Second),
    )
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Now you can use the client to interact with the Anytype API
    fmt.Println("Successfully authenticated with Anytype!")
}
```

### Working with Spaces

Spaces are the primary containers for objects in Anytype. Here's how to work with them:

```go
// Get all available spaces
ctx := context.Background()
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
} else {
    log.Fatal("No spaces available")
}
```

### Searching for Objects

You can search for objects using various criteria:

```go
// Basic search with a text query
searchParams := &anytype.SearchParams{
    Query: "project meeting",
    Limit: 50,
}

results, err := client.Search(ctx, targetSpace.ID, searchParams)
if err != nil {
    log.Fatalf("Search failed: %v", err)
}

fmt.Printf("Found %d objects matching the search query\n", len(results.Data))

// Search for objects of specific types
typeSearchParams := &anytype.SearchParams{
    Types: []string{"ot-note", "ot-task"},
    Limit: 25,
}

typeResults, err := client.Search(ctx, targetSpace.ID, typeSearchParams)
if err != nil {
    log.Fatalf("Type search failed: %v", err)
}

fmt.Printf("Found %d notes and tasks\n", len(typeResults.Data))

// Search for objects with specific tags
tagSearchParams := &anytype.SearchParams{
    Tags: []string{"important", "work"},
    Limit: 25,
}

tagResults, err := client.Search(ctx, targetSpace.ID, tagSearchParams)
if err != nil {
    log.Fatalf("Tag search failed: %v", err)
}

fmt.Printf("Found %d objects with tags 'important' or 'work'\n", len(tagResults.Data))
```

### Working with Objects

#### Retrieving Objects

```go
// Get a specific object by ID
objectParams := &anytype.GetObjectParams{
    SpaceID:  targetSpace.ID,
    ObjectID: "obj123456", // Replace with a real object ID
}

object, err := client.GetObject(ctx, objectParams)
if err != nil {
    log.Fatalf("Failed to get object: %v", err)
}

fmt.Printf("Retrieved object: %s (Type: %s)\n", object.Name, object.Type.Name)
```

#### Creating Objects

```go
// Create a new note
newNote := &anytype.Object{
    Name:    "Meeting Notes",
    TypeKey: "ot-note", // Use the appropriate type key
    Icon: &anytype.Icon{
        Format: "emoji",
        Emoji:  "📝",
    },
    Tags: []string{"meeting", "work"},
}

createdNote, err := client.CreateObject(ctx, targetSpace.ID, newNote)
if err != nil {
    log.Fatalf("Failed to create note: %v", err)
}

fmt.Printf("Created note with ID: %s\n", createdNote.ID)

// Create a task with additional properties
newTask := &anytype.Object{
    Name:    "Complete documentation",
    TypeKey: "ot-task",
    Icon: &anytype.Icon{
        Format: "emoji",
        Emoji:  "✅",
    },
    Tags: []string{"work", "urgent"},
    // Add task-specific properties if needed
}

createdTask, err := client.CreateObject(ctx, targetSpace.ID, newTask)
if err != nil {
    log.Fatalf("Failed to create task: %v", err)
}

fmt.Printf("Created task with ID: %s\n", createdTask.ID)
```

#### Updating Objects

```go
// Update an existing object
updateObj := &anytype.Object{
    Name: "Updated Meeting Notes",
    Icon: &anytype.Icon{
        Format: "emoji",
        Emoji:  "📌",
    },
    Tags: []string{"meeting", "work", "important"},
}

updatedObject, err := client.UpdateObject(ctx, targetSpace.ID, createdNote.ID, updateObj)
if err != nil {
    log.Fatalf("Failed to update object: %v", err)
}

fmt.Printf("Updated object: %s\n", updatedObject.Name)
```

#### Deleting Objects

```go
// Delete an object
err = client.DeleteObject(ctx, targetSpace.ID, createdTask.ID)
if err != nil {
    log.Fatalf("Failed to delete object: %v", err)
}

fmt.Println("Object successfully deleted")
```

### Exporting Objects

You can export objects to markdown or HTML files:

```go
// Export a single object as markdown
filePath, err := client.ExportObject(ctx, targetSpace.ID, createdNote.ID, "./exports", "md")
if err != nil {
    log.Fatalf("Failed to export object: %v", err)
}

fmt.Printf("Object exported to: %s\n", filePath)

// Export multiple objects (from search results)
exportedFiles, err := client.ExportObjects(ctx, targetSpace.ID, results.Data, "./exports", "md")
if err != nil {
    log.Fatalf("Failed to export objects: %v", err)
}

fmt.Printf("Exported %d objects:\n", len(exportedFiles))
for i, file := range exportedFiles {
    fmt.Printf("%d. %s\n", i+1, file)
}
```

## 🧩 Advanced Usage

This section covers advanced features and techniques for using Anytype-Go more effectively.

### Query Builder

For more complex queries, you can use the query builder:

```go
// Create a query builder
qb := client.NewQueryBuilder(targetSpace.ID)

// Build a complex query
qb.WithType("ot-note")
  .WithTag("important")
  .WithTextSearch("project")
  .WithLimit(25)
  .OrderBy("name", "asc")

// Execute the query
results, err := qb.Execute(ctx)
if err != nil {
    log.Fatalf("Query failed: %v", err)
}

fmt.Printf("Found %d results with query builder\n", len(results.Data))
```

### Error Handling

The API functions return specific error types that you can handle:

```go
// Example of error handling
_, err = client.GetObject(ctx, &anytype.GetObjectParams{
    SpaceID:  targetSpace.ID,
    ObjectID: "nonexistent-id",
})

if err != nil {
    // Check for specific error types
    if errors.Is(err, anytype.ErrInvalidObjectID) {
        fmt.Println("The object ID is invalid")
    } else if errors.Is(err, anytype.ErrNotFound) {
        fmt.Println("The object was not found")
    } else {
        fmt.Printf("An error occurred: %v\n", err)
    }
}
```

### Best Practices

#### Connection Reuse

Reuse the same client for multiple operations to benefit from connection pooling:

```go
// Create a client once
client, err := anytype.FromAuthConfig(authConfig)
if err != nil {
    log.Fatal(err)
}

// Use the same client for multiple operations
spaces, _ := client.GetSpaces(ctx)
results, _ := client.Search(ctx, spaceID, searchParams)
object, _ := client.GetObject(ctx, objectParams)
```

#### Context Management

Use contexts to manage timeouts and cancellation:

```go
// Create a context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel() // Always defer cancellation

// Use the context in API calls
result, err := client.Search(ctx, spaceID, params)
```

#### Batch Operations

When working with multiple objects, consider using batch operations:

```go
// Export multiple objects in one call
exportedFiles, err := client.ExportObjects(ctx, spaceID, objects, "./exports", "md")
```

### Troubleshooting

#### Common Issues

1. **Authentication Failures**: Ensure your app key and session token are valid and not expired.

2. **Connection Issues**: Check that Anytype is running locally if you're using the default API URL.

3. **Rate Limiting**: If you're making many requests in a short time, you might experience rate limiting.

#### Debug Mode

Enable debug mode to see detailed logs of API requests and responses:

```go
client, err := anytype.FromAuthConfig(authConfig, anytype.WithDebug(true))
```

#### Viewing Request Details

Use the curl option to see the equivalent curl command for each request:

```go
client, err := anytype.FromAuthConfig(authConfig, anytype.WithCurl(true))
```

## 🖥️ CLI Tool

Anytype-Go includes a command-line interface for interacting with the Anytype API without writing code.

### Configuration

The tool stores authentication configuration in:
`~/.config/anytype-go/anytype_auth.json`

On first run, it will guide you through the authentication process.

### Command Line Usage

```bash
# Basic usage with text output
anytype-go

# Display version information
anytype-go -version

# Use JSON output format
anytype-go -format json

# Enable debug mode to see API requests
anytype-go -debug

# Set log level (available levels: error, info, debug)
anytype-go -loglevel debug

# Search in a specific space
anytype-go -space "My Space" -query "search term"

# Search for objects of a specific type
anytype-go -type "Note" -query "search term"

# Search for objects of multiple types
anytype-go -types "Note,Task,Person" -query "search term"

# Search in a specific space for objects of a specific type
anytype-go -space "My Space" -type "Note" -query "search term"

# Search in a specific space for objects of multiple types
anytype-go -space "My Space" -types "Note,Task" -query "search term"

# Search for objects with specific tags
anytype-go -tags "important,work" -query "search term"

# Search for objects with specific tags in a specific space
anytype-go -space "My Space" -tags "important,work" -query "search term"

# Search for objects with specific tags and types
anytype-go -types "Note,Task" -tags "important,work" -query "search term"

# Combined search with space, types, and tags
anytype-go -space "My Space" -types "Note,Task" -tags "important,work" -query "search term"

# Print curl equivalent of all API requests
anytype-go -curl

# Set custom timeout
anytype-go -timeout 60s

# Disable colored output
anytype-go -no-color
```

### Command Line Options

- `-format`: Output format (text or json) [default: text]
- `-no-color`: Disable colored output
- `-debug`: Enable debug mode to see API requests (equivalent to -loglevel debug)
- `-loglevel`: Set logging level (error, info, debug) [default: error]
- `-timeout`: Operation timeout [default: 30s]
- `-space`: Space name to use
- `-type`: Type name to filter search results (deprecated, use -types instead)
- `-types`: Comma-separated list of type names to filter search results (e.g., 'Note,Task,Person')
- `-query`: Search query
- `-tags`: Comma-separated list of tags to filter by (e.g., 'important,work')
- `-curl`: Print curl equivalent of API requests
- `-export`: Export objects as files
- `-export-path`: Path to export files to [default: ./exports]
- `-export-format`: Format to export objects as (md, html) [default: md]
- `-version`: Display version information

## 📚 API Reference

This section details the available packages and API components for developers.

### Available Packages

#### anytype

The `anytype` package provides the main client for interacting with the Anytype API:

**Client Creation:**
- `NewClient(...ClientOption)`: Create a new API client with custom options
- `FromAuthConfig(*AuthConfig, ...ClientOption)`: Create a client from an auth configuration

**Client Options:**
- `WithTimeout(time.Duration)`: Set a custom timeout for the HTTP client
- `WithDebug(bool)`: Enable or disable debug logging
- `WithLogger(log.Logger)`: Set a custom logger for the client
- `WithCurl(bool)`: Enable printing curl equivalent of API requests
- `WithURL(string)`: Set a custom API URL
- `WithToken(string)`: Set the session token for authentication
- `WithAppKey(string)`: Set the application key for authentication

**Space Operations:**
- `GetSpaces(ctx)`: Retrieve all available spaces
- `GetSpaceByID(ctx, spaceID)`: Retrieve a specific space by ID
- `GetMembers(ctx, spaceID)`: Retrieve members of a space

**Object Operations:**
- `GetObject(ctx, params)`: Get a specific object by ID
- `CreateObject(ctx, spaceID, object)`: Create a new object
- `UpdateObject(ctx, spaceID, objectID, object)`: Update an existing object
- `DeleteObject(ctx, spaceID, objectID)`: Delete an object by ID

**Search Operations:**
- `Search(ctx, spaceID, params)`: Search for objects with filters
- `NewQueryBuilder(spaceID)`: Create a fluent query builder for complex searches

**Export Operations:**
- `ExportObject(ctx, spaceID, objectID, path, format)`: Export a single object to a file
- `ExportObjects(ctx, spaceID, objects, path, format)`: Export multiple objects to files
- `DownloadImage(ctx, imageURL, outputDir)`: Download an image from a URL

**Type Operations:**
- `GetTypes(ctx, params)`: Get all types in a space
- `GetTypeByName(ctx, spaceID, typeName)`: Find a type key by name
- `GetTypeName(ctx, spaceID, typeKey)`: Find a type name by key

**Utility Functions:**
- `Version()`: Get version information for the client library
- `GetVersionInfo()`: Get detailed version information including API version

#### auth

The `auth` package handles authentication with the Anytype API:

**Authentication Management:**
- `NewAuthManager(...ManagerOption)`: Create a new authentication manager
- `WithAPIURL(string)`: Set a custom API URL for authentication
- `WithConfigPath(string)`: Set a custom path for the auth configuration file
- `WithNonInteractive(bool)`: Enable or disable interactive authentication
- `WithSilent(bool)`: Enable or disable informational messages

**Configuration Operations:**
- `GetConfiguration()`: Get the saved authentication configuration
- `SaveConfiguration(*AuthConfig)`: Save an authentication configuration
- `GetOrCreateConfig()`: Get existing config or create a new one via authentication
- `ClearConfiguration()`: Delete the saved authentication configuration

**Authentication Methods:**
- `AuthenticateInteractive()`: Perform interactive authentication using challenge-response
- `Authenticate()`: Authenticate non-interactively using saved credentials

### Core API Components

For detailed usage examples of each component, refer to the [GoDoc documentation](https://godoc.org/github.com/epheo/anytype-go).

## 📋 Project Info

### Features

#### Core Functionality
- 🔐 Authentication management with automatic token refresh
- 📋 List and manage spaces and their members
- 🔍 Powerful search capabilities with multiple filtering options
- 📝 Complete CRUD operations for objects (Create, Read, Update, Delete)
- 📤 Export objects to multiple formats (Markdown, HTML)

#### Developer Experience
- 🧩 Fluent query builder for complex search operations
- 🔍 Type resolution from human-readable names to internal keys
- 🏷️ Tag-based filtering for better content organization
- 🐞 Debug mode with detailed logging of API interactions

#### CLI Features
- 🎨 Colored terminal output with configurable formatting
- 🔄 JSON output option for programmatic consumption
- ⏱️ Configurable operation timeout
- 🔍 Intuitive search syntax with multiple filter types

## 📁 Project Structure

The project is organized into several key packages:

```
anytype-go/
├── cmd/
│   └── anytype-go/     # Command line interface implementation
│       └── main.go     # CLI entry point and command handlers
│
├── pkg/                # Public API packages
│   ├── anytype/        # Core Anytype API client
│   │   ├── api.go      # API methods (Search, GetObject, etc.)
│   │   ├── client.go   # HTTP client and configuration
│   │   ├── errors.go   # Error types and handling
│   │   ├── export.go   # Object export functionality
│   │   ├── models.go   # Data structures for API objects
│   │   └── ...
│   │
│   └── auth/           # Authentication package
│       ├── auth.go     # Authentication management
│       └── config.go   # Config file handling
│
└── internal/           # Internal implementation details
    ├── display/        # Output formatting for CLI
    │   ├── display.go  # Display utility functions
    │   ├── printer.go  # Text/JSON output formatting
    │   └── ...
    │
    └── log/            # Logging functionality
        └── log.go      # Logging interface and implementation
```

### Key Components

- **Client (`pkg/anytype/client.go`)**: Core API client that handles HTTP communication and authentication
- **API Methods (`pkg/anytype/api.go`)**: Implementation of Anytype API operations (search, object CRUD, etc.)
- **Models (`pkg/anytype/models.go`)**: Data structures for representing Anytype entities (spaces, objects, types)
- **Query Builder (`pkg/anytype/query_builder.go`)**: Fluent interface for constructing complex search queries
- **Auth Manager (`pkg/auth/auth.go`)**: Handles authentication flow and token management
- **CLI Implementation (`cmd/anytype-go/main.go`)**: Command-line interface that utilizes the API client
- **Display (`internal/display`)**: Handles formatting output (text, JSON, tables) for the CLI

## 📊 Versioning Policy

Anytype-Go follows [Semantic Versioning 2.0.0](https://semver.org/) 

### Semantic Version Format: MAJOR.MINOR.PATCH

- **MAJOR** version increments indicate incompatible API changes
- **MINOR** version increments indicate new functionality added in a backward-compatible manner
- **PATCH** version increments indicate backward-compatible bug fixes

### API Stability

- **Public API:** All exported functions, types, and constants are considered part of the public API
- **Breaking Changes:** Will only occur with a MAJOR version increment and until we reach first major release
- **Deprecation Policy:** Features will be marked as deprecated for at least one MINOR release before removal in a MAJOR release
- **Experimental Features:** May be marked with an `Experimental` prefix and do not follow the same stability

### Version Information

You can access the current SDK version programmatically:

```go
import "github.com/epheo/anytype-go/pkg/anytype"

// Get version info
versionInfo := anytype.GetVersionInfo()
fmt.Printf("SDK Version: %s, API Version: %s\n", versionInfo.Version, versionInfo.APIVersion)
```

### Building from Source

```bash
# Build the CLI tool
go build ./cmd/anytype-go

# Install directly to your GOPATH
go install ./cmd/anytype-go
```

## Contributing

[Anytype API Reference](https://raw.githubusercontent.com/anyproto/anytype-heart/refs/heads/main/core/api/docs/swagger.yaml)

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

Apache License 2.0 - see [LICENSE](LICENSE) file for details.
