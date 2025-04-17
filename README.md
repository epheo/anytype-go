# Anytype-Go

[![Go Report Card](https://goreportcard.com/badge/github.com/epheo/anytype-go)](https://goreportcard.com/report/github.com/epheo/anytype-go)
[![GoDoc](https://godoc.org/github.com/epheo/anytype-go?status.svg)](https://godoc.org/github.com/epheo/anytype-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A Go SDK and command-line tool for interacting with the [Anytype](https://anytype.io) API to manage spaces, objects, and perform searches.

## 📚 Usage as a CLI Tool

### 🔧 Installation

```bash
go install github.com/epheo/anytype-go@latest
```

### ⚙️ Configuration

The tool stores authentication configuration in:
`~/.config/anytype-go/anytype_auth.json`

On first run, it will guide you through the authentication process.

### 💻 Command Line Usage

```bash
# Basic usage with text output
anytype-go

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

### 🎛️ Command Line Options

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

## 📦 Usage as a Go Package

The `anytype-go` project can be used as a Go package in your own applications to interact with the Anytype API.

### 📥 Installation

```bash
go get github.com/epheo/anytype-go
```

### Example Usage

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/epheo/anytype-go/pkg/anytype"
	"github.com/epheo/anytype-go/pkg/auth"
)

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize auth manager
	authManager := auth.NewAuthManager("")

	// Get configuration
	config, err := authManager.GetConfiguration()
	if err != nil {
		fmt.Printf("Authentication failed: %v\n", err)
		return
	}

	// Create API client
	client := anytype.NewClient(
		config.ApiURL,
		config.SessionToken,
		config.AppKey,
		anytype.WithDebug(false),
	)

	// Get spaces
	spaces, err := client.GetSpaces(ctx)
	if err != nil {
		fmt.Printf("Failed to get spaces: %v\n", err)
		return
	}

	// Print spaces
	fmt.Println("Available spaces:")
	for _, space := range spaces.Data {
		fmt.Printf("- %s (%s)\n", space.Name, space.ID)
	}

	// Use the first space
	if len(spaces.Data) > 0 {
		space := spaces.Data[0]
		
		// Search for objects
		searchParams := &anytype.SearchParams{
			Query: "example search",
			Limit: 10,
		}
		
		results, err := client.Search(ctx, space.ID, searchParams)
		if err != nil {
			fmt.Printf("Search failed: %v\n", err)
			return
		}
		
		// Print results
		fmt.Printf("Found %d objects:\n", len(results.Data))
		for _, obj := range results.Data {
			fmt.Printf("- %s (%s)\n", obj.Name, obj.ID)
		}
		
		// Export an object
		if len(results.Data) > 0 {
			obj := results.Data[0]
			filePath, err := client.ExportObject(ctx, space.ID, obj.ID, "./exports", "markdown")
			if err != nil {
				fmt.Printf("Export failed: %v\n", err)
				return
			}
			fmt.Printf("Object exported to: %s\n", filePath)
		}
	}
}
```

### Available Packages

#### anytype

The `anytype` package provides the main client for interacting with the Anytype API:

- `NewClient()`: Create a new API client
- `GetSpaces()`: Retrieve available spaces
- `GetObject()`: Get a specific object by ID
- `Search()`: Search for objects with filters
- `ExportObject()`: Export an object to a file
- `ExportObjects()`: Export multiple objects to files

#### auth

The `auth` package handles authentication with the Anytype API:

- `NewAuthManager()`: Create a new authentication manager
- `GetConfiguration()`: Get the saved authentication configuration
- `AuthenticateInteractive()`: Perform interactive authentication

## ✨ Features

- 🔐 Authentication management with automatic token refresh
- 📋 List available spaces
- 🔍 Search for objects within spaces
- 🏷️ Filter searches by object type (uses type's unique key internally)
- 🎨 Colored terminal output
- 🐞 Debug mode for API requests
- ⏱️ Configurable operation timeout
- 📤 Export objects to files in different formats

### 📁 Project Structure

```
cmd/
  anytype-go/        # Command line interface
    main.go          # Entry point
pkg/
  anytype/           # Anytype API client
    api.go           # API methods
    client.go        # HTTP client
    export.go        # Object export functionality
    helper.go        # Utility functions
    models.go        # Data models
  auth/              # Authentication
    auth.go          # Auth management
    config.go        # Config handling
internal/
  display/           # Output formatting (for internal use only)
```

### 🔨 Building

```bash
go build ./cmd/anytype-go
```

## Contributing

[Anytype API Reference](https://github.com/anyproto/anytype-heart/blob/main/core/api/docs/swagger.yaml)

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

Apache License 2.0 - see [LICENSE](LICENSE) file for details.
