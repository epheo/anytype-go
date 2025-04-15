package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/epheo/anyblog/pkg/anytype"
	"github.com/epheo/anyblog/pkg/auth"
	"github.com/epheo/anyblog/pkg/display"
)

// Command line flags
type flags struct {
	format    string
	noColor   bool
	debug     bool
	timeout   time.Duration
	spaceName string
	typeName  string
	query     string
	tags      string // New: comma-separated list of tags to filter by
}

const defaultTimeout = 30 * time.Second

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// setupClient creates and configures the API client
func setupClient(f *flags) (*anytype.Client, display.Printer, error) {
	// Initialize display
	printer := display.NewPrinter(f.format, !f.noColor, f.debug)

	// Initialize auth manager
	authManager := auth.NewAuthManager("")

	// Get configuration
	config, err := authManager.GetConfiguration()
	if err != nil {
		return nil, printer, fmt.Errorf("authentication failed: %w", err)
	}

	// Create API client with options
	client := anytype.NewClient(
		config.ApiURL,
		config.SessionToken,
		config.AppKey,
		anytype.WithDebug(f.debug),
	)

	return client, printer, nil
}

// findTargetSpace finds the target space based on name or returns the first available
func findTargetSpace(spaces *anytype.SpacesResponse, spaceName string, printer display.Printer) (*anytype.Space, error) {
	for _, space := range spaces.Data {
		if space.Name == spaceName {
			spacePtr := space
			printer.PrintInfo("Found space: %s (%s)", space.Name, space.ID)
			return &spacePtr, nil
		}
	}

	if len(spaces.Data) > 0 {
		spacePtr := spaces.Data[0]
		printer.PrintInfo("Using default space: %s (%s)", spacePtr.Name, spacePtr.ID)
		return &spacePtr, nil
	}

	return nil, fmt.Errorf("no spaces available")
}

// handleSearch performs the search operation with the given parameters
func handleSearch(ctx context.Context, client *anytype.Client, targetSpace *anytype.Space, params *anytype.SearchParams, printer display.Printer) error {
	results, err := client.Search(ctx, targetSpace.ID, params)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if err := printer.PrintObjects("Search Results", results.Data, client, ctx); err != nil {
		return fmt.Errorf("failed to display search results: %w", err)
	}

	return nil
}

func run() error {
	// Parse command line flags
	f := parseFlags()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel()

	// Setup client and printer
	client, printer, err := setupClient(f)
	if err != nil {
		return err
	}

	// Get spaces
	spaces, err := client.GetSpaces(ctx)
	if err != nil {
		return fmt.Errorf("failed to get spaces: %w", err)
	}

	if err := printer.PrintSpaces(spaces.Data); err != nil {
		return fmt.Errorf("failed to display spaces: %w", err)
	}

	// Find target space
	targetSpace, err := findTargetSpace(spaces, f.spaceName, printer)
	if err != nil {
		return err
	}

	// Find specific type if requested
	if f.typeName != "" {
		typeID, err := client.GetTypeByName(ctx, targetSpace.ID, f.typeName)
		if err != nil {
			printer.PrintError("Could not find type '%s': %v", f.typeName, err)
		} else {
			printer.PrintSuccess("Found type '%s' with ID: %s", f.typeName, typeID)
		}
	}

	// Perform search if query or tags provided
	if f.query != "" || f.tags != "" {
		searchParams := &anytype.SearchParams{
			Query: strings.TrimSpace(f.query),
			Types: []string{"ot-page"}, // Default to ot-page type
			Limit: 100,
		}

		// Add type filter if type name is specified
		if f.typeName != "" {
			typeKey, err := client.GetTypeByName(ctx, targetSpace.ID, f.typeName)
			if err != nil {
				printer.PrintError("Could not find type '%s': %v", f.typeName, err)
			} else {
				searchParams.Types = []string{typeKey}
				printer.PrintInfo("Filtering search results by type: %s", f.typeName)
			}
		}

		// Add tags filter if tags are specified
		if f.tags != "" {
			tags := strings.Split(f.tags, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
			searchParams.Tags = tags
			printer.PrintInfo("Filtering search results by tags: %s", strings.Join(tags, ", "))
		}

		if err := handleSearch(ctx, client, targetSpace, searchParams, printer); err != nil {
			return err
		}
	}

	return nil
}

func parseFlags() *flags {
	f := &flags{}

	flag.StringVar(&f.format, "format", "text", "Output format (text or json)")
	flag.BoolVar(&f.noColor, "no-color", false, "Disable colored output")
	flag.BoolVar(&f.debug, "debug", false, "Enable debug mode")
	flag.DurationVar(&f.timeout, "timeout", defaultTimeout, "Operation timeout")
	flag.StringVar(&f.spaceName, "space", "", "Space name to use")
	flag.StringVar(&f.typeName, "type", "", "Type name to look for")
	flag.StringVar(&f.query, "query", "", "Search query")
	flag.StringVar(&f.tags, "tags", "", "Comma-separated list of tags to filter by (e.g., 'important,work')")

	flag.Parse()

	return f
}
