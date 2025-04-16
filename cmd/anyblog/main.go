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
	format       string
	noColor      bool
	debug        bool
	logLevel     string
	timeout      time.Duration
	spaceName    string
	typeName     string // Single type name (deprecated)
	types        string // Comma-separated list of type names
	query        string
	tags         string // Comma-separated list of tags to filter by
	curl         bool   // Print curl equivalent of API requests
	export       bool   // Export objects as files
	exportPath   string // Path to export files to
	exportFormat string // Format to export objects as (md, html, etc.)
}

// exportOptions defines options for exporting objects
type exportOptions struct {
	enabled bool
	path    string
	format  string
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

	// Set log level (debug flag overrides loglevel flag)
	if f.debug {
		printer.SetLogLevel(display.LogLevelDebug)
	} else {
		level := display.ParseLogLevel(f.logLevel)
		printer.SetLogLevel(level)
	}

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
		anytype.WithCurl(f.curl),
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
func handleSearch(ctx context.Context, client *anytype.Client, targetSpace *anytype.Space, params *anytype.SearchParams, printer display.Printer, exportOptions *exportOptions) error {
	results, err := client.Search(ctx, targetSpace.ID, params)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if err := printer.PrintObjects("Search Results", results.Data, client, ctx); err != nil {
		return fmt.Errorf("failed to display search results: %w", err)
	}

	// Handle export if enabled
	if exportOptions != nil && exportOptions.enabled {
		printer.PrintInfo("Exporting %d objects to %s in %s format", len(results.Data), exportOptions.path, exportOptions.format)

		// Create export directory if it doesn't exist
		if err := os.MkdirAll(exportOptions.path, 0755); err != nil {
			return fmt.Errorf("failed to create export directory: %w", err)
		}

		exportedFiles, err := client.ExportObjects(ctx, targetSpace.ID, results.Data, exportOptions.path, exportOptions.format)
		if err != nil {
			return fmt.Errorf("export failed: %w", err)
		}

		printer.PrintSuccess("Successfully exported %d objects:", len(exportedFiles))
		for i, file := range exportedFiles {
			printer.PrintInfo("  %d. %s", i+1, file)
		}
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

	// Identify the type ID if a type name is specified
	var typeID string
	if f.typeName != "" {
		var err error
		typeID, err = client.GetTypeByName(ctx, targetSpace.ID, f.typeName)
		if err != nil {
			printer.PrintError("Could not find type '%s': %v", f.typeName, err)
		} else {
			printer.PrintSuccess("Found type '%s' with ID: %s", f.typeName, typeID)
		}
	}

	// Set up export options if export is enabled
	var exportOpts *exportOptions
	if f.export {
		exportOpts = &exportOptions{
			enabled: true,
			path:    f.exportPath,
			format:  f.exportFormat,
		}
		printer.PrintInfo("Export enabled. Objects will be exported to %s in %s format", f.exportPath, f.exportFormat)
	}

	// Perform search if query, tags, or types are provided
	if f.query != "" || f.tags != "" || f.typeName != "" || f.types != "" {
		searchParams := &anytype.SearchParams{
			Query: strings.TrimSpace(f.query),
			Types: []string{"ot-page"}, // Default to ot-page type
			Limit: 100,
		}

		// Process type filters (priority given to -types over -type for backwards compatibility)
		if f.types != "" {
			// Handle multiple types
			typeNames := strings.Split(f.types, ",")
			typeKeys := []string{}
			typeNamesFound := []string{}

			for _, typeName := range typeNames {
				typeName = strings.TrimSpace(typeName)
				if typeName == "" {
					continue
				}

				typeKey, err := client.GetTypeByName(ctx, targetSpace.ID, typeName)
				if err != nil {
					printer.PrintError("Could not find type '%s': %v", typeName, err)
				} else {
					typeKeys = append(typeKeys, typeKey)
					typeNamesFound = append(typeNamesFound, typeName)
				}
			}

			if len(typeKeys) > 0 {
				searchParams.Types = typeKeys
				printer.PrintInfo("Filtering search results by types: %s", strings.Join(typeNamesFound, ", "))
			}
		} else if f.typeName != "" {
			// For backward compatibility: handle single type
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

		if err := handleSearch(ctx, client, targetSpace, searchParams, printer, exportOpts); err != nil {
			return err
		}
	} else if f.export {
		// If export is enabled but no search parameters are provided,
		// fetch all objects from the space
		printer.PrintInfo("No search parameters provided, exporting all objects from space %s (%s)", targetSpace.Name, targetSpace.ID)

		searchParams := &anytype.SearchParams{
			Types: []string{"ot-page"}, // Default to ot-page type
			Limit: 100,
		}

		if err := handleSearch(ctx, client, targetSpace, searchParams, printer, exportOpts); err != nil {
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
	flag.StringVar(&f.logLevel, "loglevel", "error", "Log level (error, info, debug)")
	flag.DurationVar(&f.timeout, "timeout", defaultTimeout, "Operation timeout")
	flag.StringVar(&f.spaceName, "space", "", "Space name to use")
	flag.StringVar(&f.typeName, "type", "", "Type name to look for (deprecated, use -types instead)")
	flag.StringVar(&f.types, "types", "", "Comma-separated list of type names to filter by (e.g., 'Note,Task')")
	flag.StringVar(&f.query, "query", "", "Search query")
	flag.StringVar(&f.tags, "tags", "", "Comma-separated list of tags to filter by (e.g., 'important,work')")
	flag.BoolVar(&f.curl, "curl", false, "Print curl equivalent of API requests")

	// Export options
	flag.BoolVar(&f.export, "export", false, "Export objects as files")
	flag.StringVar(&f.exportPath, "export-path", "./exports", "Path to export files to")
	flag.StringVar(&f.exportFormat, "export-format", "md", "Format to export objects as (md, html)")

	flag.Parse()

	return f
}
