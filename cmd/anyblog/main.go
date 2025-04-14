package main

import (
	"context"
	"flag"
	"fmt"
	"os"
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
}

const defaultTimeout = 30 * time.Second

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Parse command line flags
	f := parseFlags()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel()

	// Initialize display
	printer := display.NewPrinter(f.format, !f.noColor)

	// Initialize auth manager
	authManager := auth.NewAuthManager("")

	// Get configuration
	config, err := authManager.GetConfiguration()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create API client with options
	client := anytype.NewClient(
		config.ApiURL,
		config.SessionToken,
		config.AppKey,
		anytype.WithDebug(f.debug),
	)

	// Get spaces
	spaces, err := client.GetSpaces(ctx)
	if err != nil {
		return fmt.Errorf("failed to get spaces: %w", err)
	}

	if err := printer.PrintSpaces(spaces.Data); err != nil {
		return fmt.Errorf("failed to display spaces: %w", err)
	}

	// Find target space
	var targetSpace *anytype.Space
	for _, space := range spaces.Data {
		if space.Name == f.spaceName {
			targetSpace = &space
			printer.PrintInfo("Found space: %s (%s)", space.Name, space.ID)
			break
		}
	}

	// If specified space not found, use first available
	if targetSpace == nil && len(spaces.Data) > 0 {
		targetSpace = &spaces.Data[0]
		printer.PrintInfo("Using default space: %s (%s)", targetSpace.Name, targetSpace.ID)
	}

	if targetSpace == nil {
		return fmt.Errorf("no spaces available")
	}

	// Get types for the space
	types, err := client.GetTypes(ctx, targetSpace.ID)
	if err != nil {
		return fmt.Errorf("failed to get types: %w", err)
	}

	if err := printer.PrintJSON("Available Types", types); err != nil {
		return fmt.Errorf("failed to display types: %w", err)
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

	// Perform search if query provided
	if f.query != "" {
		searchParams := &anytype.SearchParams{
			Query: f.query,
			Limit: 10,
		}

		results, err := client.Search(ctx, targetSpace.ID, searchParams)
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}

		if err := printer.PrintObjects("Search Results", results.Items); err != nil {
			return fmt.Errorf("failed to display search results: %w", err)
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

	flag.Parse()

	return f
}
