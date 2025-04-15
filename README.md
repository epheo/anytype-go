# Anyblog

A command-line tool for interacting with the Anytype API to manage spaces, objects, and perform searches.

## Installation

```bash
go install github.com/epheo/anyblog@latest
```

## Configuration

The tool stores authentication configuration in:
`~/.config/anyblog/anytype_auth.json`

On first run, it will guide you through the authentication process.

## Usage

```bash
# Basic usage with text output
anyblog

# Use JSON output format
anyblog -format json

# Enable debug mode to see API requests
anyblog -debug

# Search in a specific space
anyblog -space "My Space" -query "search term"

# Search for objects of a specific type
anyblog -type "Note" -query "search term"

# Search in a specific space for objects of a specific type
anyblog -space "My Space" -type "Note" -query "search term"

# Set custom timeout
anyblog -timeout 60s

# Disable colored output
anyblog -no-color
```

### Command Line Options

- `-format`: Output format (text or json) [default: text]
- `-no-color`: Disable colored output
- `-debug`: Enable debug mode to see API requests
- `-timeout`: Operation timeout [default: 30s]
- `-space`: Space name to use
- `-type`: Type name to filter search results
- `-query`: Search query

## Features

- Authentication management with automatic token refresh
- List available spaces
- Search for objects within spaces
- Filter searches by object type (uses type's unique key internally)
- Query object types
- Pretty-printed output in text or JSON format
- Colored terminal output
- Debug mode for API requests
- Configurable operation timeout

## Development

### Project Structure

```
cmd/
  anyblog/        # Command line interface
    main.go       # Entry point
pkg/
  anytype/        # Anytype API client
    api.go        # API methods
    client.go     # HTTP client
    helper.go     # Utility functions
    models.go     # Data models
  auth/           # Authentication
    auth.go       # Auth management
    config.go     # Config handling
  display/        # Output formatting
    display.go    # Display utilities
```

### Building

```bash
go build ./cmd/anyblog
```

### Running Tests

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

MIT License - see LICENSE file for details

