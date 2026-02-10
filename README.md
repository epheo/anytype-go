# Anytype-Go SDK

[![Go Report Card](https://goreportcard.com/badge/github.com/epheo/anytype-go)](https://goreportcard.com/report/github.com/epheo/anytype-go)
[![GoDoc](https://godoc.org/github.com/epheo/anytype-go?status.svg)](https://godoc.org/github.com/epheo/anytype-go)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A Go SDK for interacting with the [Anytype](https://anytype.io) API to manage spaces, objects, and perform searches. This library provides a feature-rich, fluent interface to integrate Anytype functionality into your Go applications.

## Table of Contents

- [Overview](#overview)
- [API Version Support](#api-version-support)
- [Installation](#installation)
- [Quick Start](#quick-start)
  - [Authentication](#authentication)
  - [Working with Spaces](#working-with-spaces)
  - [Working with Objects](#working-with-objects)
  - [Searching](#searching)
- [Advanced Examples](#advanced-examples)
  - [Working with Object Types and Templates](#working-with-object-types-and-templates)
  - [Managing Object Properties](#managing-object-properties)
  - [Working with Lists and Views](#working-with-lists-and-views)
- [Design Philosophy](#design-philosophy)
  - [1. Fluent Interface Pattern](#1-fluent-interface-pattern)
  - [2. Domain-Driven Design](#2-domain-driven-design)
  - [3. Naming Convention](#3-naming-convention)
  - [4. Middleware Architecture](#4-middleware-architecture)
- [API Reference](#api-reference)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)
- [Testing](#testing)
  - [Unit Tests with Mocks](#unit-tests-with-mocks)
  - [API Coverage Tests](#api-coverage-tests)
- [Contributing](#contributing)
- [License](#license)

## Overview

Anytype-Go provides a Go SDK for interacting with Anytype's local API. This SDK offers a clean, fluent interface for:

- Managing spaces and their members
- Creating, reading, updating, and deleting objects
- Searching for objects using various filters
- Exporting objects to different formats
- Working with object types, templates, and properties

## API Version Support

This SDK is compatible with Anytype API version `2025-11-08`. The SDK follows Anytype's API versioning scheme, which uses date-based versioning for stability and compatibility:

- All API requests include the `Anytype-Version` header set to `2025-11-08`
- Type keys such as `page` and `collection` follow the latest API specification
- Authentication uses the app key Bearer token method

If you encounter any compatibility issues when Anytype updates its API, please check for an updated version of this SDK that supports the new API version.

## Installation

```bash
go get github.com/epheo/anytype-go
```

Import both the main package and client implementation:

```go
import (
    "github.com/epheo/anytype-go"
    _ "github.com/epheo/anytype-go/client" // Register client implementation
)
```

## Quick Start

> **More complete examples** can be found in the [examples](./examples) directory, including full implementations of authentication, working with spaces, objects, and more.

### Authentication

To use the Anytype API, you need to obtain an AppKey. Here's how to get it:

```go
// Step 1: Initiate authentication and get challenge ID
authResponse, err := client.Auth().DisplayCode(ctx, "MyAnytypeApp")
if err != nil {
    log.Fatalf("Failed to initiate authentication: %v", err)
}
challengeID := authResponse.ChallengeID

// Step 2: User needs to enter the code shown in Anytype app
fmt.Println("Please enter the authentication code shown in Anytype:")
var code string
fmt.Scanln(&code)

// Step 3: Complete authentication and get tokens
tokenResponse, err := client.Auth().GetToken(ctx, challengeID, code)
if err != nil {
    log.Fatalf("Authentication failed: %v", err)
}

// Now you have your authentication token
appKey := tokenResponse.AppKey

// Create authenticated client
client := anytype.NewClient(
    anytype.WithBaseURL("http://localhost:31009"), // Default Anytype local API URL
    anytype.WithAppKey(appKey),
)
```

> **Note**: The authentication flow requires user interaction. When you call `DisplayCode`, Anytype will show a verification code that must be entered in your application.

### Working with Spaces

```go
// List all spaces
spaces, err := client.Spaces().List(ctx)

// Get a specific space
space, err := client.Space(spaceID).Get(ctx)

// Create a new space
newSpace, err := client.Spaces().Create(ctx, anytype.CreateSpaceRequest{
    Name:        "My New Workspace",
    Description: "Created via the Go SDK",
})
```

### Working with Objects

```go
// Get an object
object, err := client.Space(spaceID).Object(objectID).Get(ctx)

// Delete an object
err = client.Space(spaceID).Object(objectID).Delete(ctx)

// Export an object to markdown
exportResult, err := client.Space(spaceID).Object(objectID).Export(ctx, "markdown")

// Create a new object
newObject, err := client.Space(spaceID).Objects().Create(ctx, anytype.CreateObjectRequest{
    TypeKey:     "page",
    Name:        "My New Page",
    Description: "Created via the Go SDK",
    Body:        "# This is a new page\n\nWith some content in markdown format.",
    Icon: &anytype.Icon{
        Format: anytype.IconFormatEmoji,
        Emoji: "📄",
    },
})
```

### Searching

```go
// Search within a specific space
results, err := client.Space(spaceID).Search(ctx, anytype.SearchRequest{
    Query: "important notes",
    Sort: &anytype.SortOptions{
        Property:  anytype.SortPropertyLastModifiedDate,
        Direction: anytype.SortDirectionDesc,
    },
    Types: []string{"note", "page"}, // Filter by specific types
})
```

## Advanced Examples

### Working with Object Types and Templates

```go
// List available object types in a space
objectTypes, err := client.Space(spaceID).Types().List(ctx)

// Get details of a specific object type
typeDetails, err := client.Space(spaceID).Type(typeKey).Get(ctx)

// List templates for a specific object type
templates, err := client.Space(spaceID).Type(typeKey).Templates().List(ctx)

// Get details of a specific template
template, err := client.Space(spaceID).Type(typeKey).Template(templateID).Get(ctx)
```

### Managing Object Properties

```go
// Update object properties
err := client.Space(spaceID).Object(objectID).UpdateProperties(ctx, anytype.UpdatePropertiesRequest{
    Properties: map[string]interface{}{
        "name":        "Updated Title",
        "description": "Updated description",
        "status":      "In Progress",
        "priority":    "High",
        "deadline":    time.Now().AddDate(0, 0, 14).Format(time.RFC3339),
    },
})

// Add a relation to another object
err := client.Space(spaceID).Object(objectID).AddRelation(ctx, relatedObjectID, "related-to")
```

### Working with Lists and Views

```go
// Create a new list to organize objects
newList, err := client.Space(spaceID).Lists().Create(ctx, anytype.CreateListRequest{
    Name:        "Project Tasks",
    Description: "All tasks for our current project",
    Icon: &anytype.Icon{
        Type:  "emoji",
        Value: "📝",
    },
})

// Add objects to a list
err := client.Space(spaceID).List(listID).AddObjects(ctx, []string{objectID1, objectID2})

// Create a custom view for a list
newView, err := client.Space(spaceID).List(listID).Views().Create(ctx, anytype.CreateViewRequest{
    Name: "Priority View",
    Type: "board",
    GroupBy: []string{"priority"},
    SortBy: []anytype.SortOptions{
        {
            Property:  "deadline",
            Direction: anytype.SortDirectionAsc,
        },
    },
    Filters: []anytype.FilterCondition{
        {
            Property: "status",
            Operator: anytype.OperatorNotEquals,
            Value:    "Completed",
        },
    },
})
```

## Design Philosophy

The Anytype-Go SDK is built around three core design principles:

### 1. Fluent Interface Pattern

```go
exportedMarkdown, err := client.
    Space(spaceID).
    Object(objectID).
    Export(ctx, "markdown")
```

**Benefits:**

- Readable code that mirrors natural language
- IDE autocomplete reveals available operations
- Compile-time type safety
- Reduced boilerplate code

### 2. Domain-Driven Design

The SDK is organized around Anytype's core concepts (spaces, objects, types) with interfaces that map directly to these concepts:

- `SpaceClient` and `SpaceContext` for spaces
- `ObjectClient` and `ObjectContext` for objects
- `TypeClient` for object types
- `ListClient` for lists and views

### 3. Naming Convention

The naming of interfaces and types in this library follows a clear and consistent pattern to improve code readability and API fluency:

- **`<Entity>Client`**: Represents a client that operates on collections of entities (e.g., `SpaceClient` for working with multiple spaces, `TypeClient` for working with multiple object types). These clients handle operations like listing, searching, and creating new entities.
- **`<Entity>Context`**: Represents a client that operates on a single, specific entity instance (e.g., `SpaceContext` for a single space, `ObjectContext` for a single object). These handle operations like getting details, updating, or deleting a specific entity.

This naming convention creates a fluent, chainable API where:

1. Collection operations use the `<Entity>Client` pattern (e.g., `client.Spaces().List()`)
2. Instance operations use the `<Entity>Context` pattern (e.g., `client.Space(spaceID).Get()`)
3. Nested resources follow a natural hierarchy (e.g., `client.Space(spaceID).Object(objectID).Export()`)

This design enables intuitive navigation through the API that mirrors natural language and domain concepts.

### 4. Middleware Architecture

```text
HTTP Request → ValidationMiddleware → RetryMiddleware → DisconnectMiddleware → HTTP Client → API
```

Each middleware handles a specific concern:

- **Validation**: Validates requests before sending
- **Retry**: Handles transient errors with configurable policies
- **Disconnect**: Manages network interruptions

## API Reference

For detailed API documentation, see [GoDoc](https://godoc.org/github.com/epheo/anytype-go).

## Best Practices

1. **Reuse the client instance** across your application
2. **Use context for cancellation** to control timeouts
3. **Handle rate limiting** with appropriate backoff strategies
4. **Validate inputs** before making API calls
5. **Check for errors** and handle them appropriately
6. **Use the fluent interface** for cleaner, more readable code

## Troubleshooting

- **Authentication Failures**: Verify your app key
- **Connection Issues**: Ensure Anytype is running locally
- **Rate Limiting**: Implement backoff if making many requests
- **API Version Mismatch**: If you get errors about unknown fields or unexpected responses, check that your Anytype app version is compatible with the API version this SDK supports (2025-11-08)

## Testing

The SDK testing approach focuses on behavior verification using mock implementations:

### Unit Tests with Mocks

`tests`: Tests ensure that client interfaces behave according to specifications using mock implementations to simulate API responses.

```bash
go test -v ./tests/...
```

### API Coverage Tests

`tests_api_coverage`: Tests verify that all API endpoints are properly defined and can be called with appropriate parameters.

```bash
go test -v ./tests_api_coverage/...
```

The test infrastructure uses mock implementations (in `tests/mocks`) to simulate the Anytype API, allowing thorough testing without requiring a running Anytype instance.

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

Apache License 2.0 - see [LICENSE](LICENSE) file for details.
