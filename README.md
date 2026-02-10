# Anytype-Go SDK

[![Go Report Card](https://goreportcard.com/badge/github.com/epheo/anytype-go)](https://goreportcard.com/report/github.com/epheo/anytype-go)
[![GoDoc](https://godoc.org/github.com/epheo/anytype-go?status.svg)](https://godoc.org/github.com/epheo/anytype-go)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Go SDK for the [Anytype](https://anytype.io) local API. Compatible with API version `2025-11-08`.

## Installation

```bash
go get github.com/epheo/anytype-go
```

```go
import (
    "github.com/epheo/anytype-go"
    _ "github.com/epheo/anytype-go/client"
)
```

## Authentication

```go
client := anytype.NewClient(
    anytype.WithBaseURL("http://localhost:31009"),
)

// Initiate challenge — Anytype app will display a verification code
auth, _ := client.Auth().CreateChallenge(ctx, "MyApp")

// User enters the code
var code string
fmt.Scanln(&code)

// Exchange for API key
token, _ := client.Auth().CreateApiKey(ctx, auth.ChallengeID, code)

// Create authenticated client
client = anytype.NewClient(
    anytype.WithBaseURL("http://localhost:31009"),
    anytype.WithAppKey(token.ApiKey),
)
```

## Usage

The SDK uses a fluent interface that mirrors the resource hierarchy:

```go
// Spaces
spaces, _ := client.Spaces().List(ctx)
space, _ := client.Space(spaceID).Get(ctx)

// Objects
objects, _ := client.Space(spaceID).Objects().List(ctx)
object, _ := client.Space(spaceID).Object(objectID).Get(ctx)
markdown, _ := client.Space(spaceID).Object(objectID).Get(ctx, anytype.WithFormat("md"))

// Types and templates
types, _ := client.Space(spaceID).Types().List(ctx)
templates, _ := client.Space(spaceID).Type(typeKey).Templates().List(ctx)

// Lists and views
views, _ := client.Space(spaceID).List(listID).Views().List(ctx)
_ = client.Space(spaceID).List(listID).Objects().Add(ctx, []string{objectID})

// Properties and tags
props, _ := client.Space(spaceID).Properties().List(ctx)
tags, _ := client.Space(spaceID).Property(propID).Tags().List(ctx)

// Members
members, _ := client.Space(spaceID).Members().List(ctx)

// Search (within a space or globally)
results, _ := client.Space(spaceID).Search(ctx, anytype.SearchRequest{
    Query: "meeting notes",
    Types: []string{"page"},
    Sort:  &anytype.SortOptions{
        Property:  anytype.SortPropertyLastModifiedDate,
        Direction: anytype.SortDirectionDesc,
    },
})
globalResults, _ := client.Search().Search(ctx, anytype.SearchRequest{Query: "notes"})
```

See [examples/usage_examples.go](./examples/usage_examples.go) for complete working examples and [GoDoc](https://godoc.org/github.com/epheo/anytype-go) for the full API reference.

## Testing SDK coverage

```bash
go test -v ./tests_api_coverage/... # API coverage tests
```

## License

Apache License 2.0 - see [LICENSE](LICENSE).
