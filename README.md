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
work, _ := client.Spaces().GetByName(ctx, "Work") // exact match, ErrSpaceNotFound otherwise

// Objects
objects, _ := client.Space(spaceID).Objects().List(ctx) // objects.Data, objects.Pagination
object, _ := client.Space(spaceID).Object(objectID).Get(ctx)
markdown, _ := client.Space(spaceID).Object(objectID).Get(ctx, anytype.WithFormat("md"))

// Types and templates
types, _ := client.Space(spaceID).Types().List(ctx)
task, _ := client.Space(spaceID).Types().Get(ctx, "task") // by key; Type(id) takes an ID
templates, _ := client.Space(spaceID).Type(task.ID).Templates().List(ctx)

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

### Pagination

Every `List` and `Search` returns a `*Page[T]` holding one page of `Data` plus a
`Pagination` cursor, and accepts `WithLimit`/`WithOffset`:

```go
page, _ := client.Space(spaceID).Objects().List(ctx, anytype.WithLimit(50), anytype.WithOffset(100))
for _, obj := range page.Data { /* ... */ }
if page.Pagination.HasMore { /* fetch the next page */ }
```

To walk every page without managing offsets yourself, use `All`, which returns a
Go 1.23 iterator (`iter.Seq2[T, error]`):

```go
for obj, err := range client.Space(spaceID).Objects().All(ctx) {
    if err != nil {
        return err
    }
    // ... use obj
}
```

`All` is available on every list, and on `Search` (`client.Search().All(ctx, req)`,
`client.Space(spaceID).SearchAll(ctx, req)`).

### Errors

Non-2xx responses come back as `*anytype.APIError` carrying the server's
`Status`, `Code` and `Message`. Any 404, and the by-key/by-name lookups that
scan a list, match `anytype.ErrNotFound`:

```go
_, err := client.Space(spaceID).Types().Get(ctx, "nope")
if errors.Is(err, anytype.ErrNotFound) { /* no such type */ }

var apiErr *anytype.APIError
if errors.As(err, &apiErr) && apiErr.Status == 429 { /* back off */ }
```

### Property values

Property formats, tag colors and creatable type layouts are typed constants
(`anytype.PropertyFormatText`, `anytype.ColorRed`, `anytype.TypeLayoutNote`).
When the format is only known at runtime, `NewPropertyLinkValue` picks the
right typed value and rejects a mismatched Go type before the request is sent:

```go
v, err := anytype.NewPropertyLinkValue(prop.Key, prop.Format, 42)
obj, _ := client.Space(spaceID).Objects().Create(ctx, anytype.CreateObjectRequest{
    TypeKey:    "task",
    Name:       "Pay invoice",
    Icon:       anytype.EmojiIcon("💸"),
    Properties: []anytype.PropertyLinkValue{v},
})
```

See [examples/usage_examples.go](./examples/usage_examples.go) for complete working examples and [GoDoc](https://godoc.org/github.com/epheo/anytype-go) for the full API reference.

## Testing SDK coverage

```bash
go test -v ./tests_api_coverage/... # API coverage tests
```

## License

Apache License 2.0 - see [LICENSE](LICENSE).
