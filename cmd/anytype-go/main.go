package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/epheo/anytype-go/internal/display"
	"github.com/epheo/anytype-go/pkg/anytype"
	"github.com/epheo/anytype-go/pkg/auth"
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
	version      bool   // Display version information
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

// setupClient creates and configures the API client and display printer.
//
// This function initializes both the Anytype API client and the display printer
// based on the provided command line flags. It handles:
//
// 1. Setting up the display printer with the appropriate format and color options
// 2. Configuring the logging level based on flags
// 3. Initializing the authentication manager with appropriate options
// 4. Setting up the client options (debug mode, timeout, curl output)
// 5. Creating the client using the auth manager's helper function
//
// Parameters:
//   - f: A pointer to the parsed command line flags
//
// Returns:
//   - An initialized Anytype API client
//   - A configured display printer for output formatting
//   - Any error encountered during setup
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

	// Initialize auth manager with options
	authManager := auth.NewAuthManager(
		auth.WithAPIURL(""),            // Use default
		auth.WithNonInteractive(false), // Allow interactive authentication
		auth.WithSilent(false),         // Show informational messages
	)

	// Set up the client options
	clientOpts := []anytype.ClientOption{
		anytype.WithDebug(f.debug),
		anytype.WithCurl(f.curl),
	}

	// Create client using the auth manager's helper function
	client, err := authManager.GetClient(clientOpts...)
	if err != nil {
		return nil, printer, fmt.Errorf("failed to create API client: %w", err)
	}

	return client, printer, nil
}

// setupSpaces gets and displays spaces, and finds the target space.
//
// This function retrieves all available spaces from the Anytype API, displays them
// to the user, and then determines which space to use for subsequent operations.
// If a specific space name is provided, it will look for a space with that name.
// Otherwise, it will use the first available space.
//
// Parameters:
//   - ctx: Context for the API request
//   - client: The initialized Anytype API client
//   - spaceName: Optional name of the space to use (empty string means use first available)
//   - printer: Display printer for output formatting
//
// Returns:
//   - A pointer to the selected Space
//   - Any error encountered during the process
//
// The function delegates to findTargetSpace to select the appropriate space based
// on the provided spaceName.
func setupSpaces(ctx context.Context, client *anytype.Client, spaceName string, printer display.Printer) (*anytype.Space, error) {
	// Get spaces
	spaces, err := client.GetSpaces(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get spaces: %w", err)
	}

	if err := printer.PrintSpaces(spaces.Data); err != nil {
		return nil, fmt.Errorf("failed to display spaces: %w", err)
	}

	// Find target space
	return findTargetSpace(spaces, spaceName, printer)
}

// findTargetSpace finds the target space based on name or returns the first available.
//
// This function searches through the available spaces for one matching the provided
// spaceName. If a matching space is found, it returns that space. If no match is found
// and spaces are available, it returns the first space. If no spaces are available,
// it returns an error.
//
// Parameters:
//   - spaces: Response containing available spaces
//   - spaceName: Name of the space to look for (can be empty)
//   - printer: Display printer for output formatting
//
// Returns:
//   - A pointer to the selected Space
//   - An error if no spaces are available or if a specified space name cannot be found
//
// This function is typically called by setupSpaces after retrieving the list of spaces
// from the Anytype API.
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

// processTypeFilters resolves type names to type keys for search filtering
func processTypeFilters(ctx context.Context, client *anytype.Client, spaceID string, typeNames []string, printer display.Printer) ([]string, []string) {
	typeKeys := []string{}
	typeNamesFound := []string{}

	for _, typeName := range typeNames {
		typeName = strings.TrimSpace(typeName)
		if typeName == "" {
			continue
		}

		typeKey, err := client.GetTypeByName(ctx, spaceID, typeName)
		if err != nil {
			printer.PrintError("Could not find type '%s': %v", typeName, err)
		} else if typeKey != "" {
			typeKeys = append(typeKeys, typeKey)
			typeNamesFound = append(typeNamesFound, typeName)
		} else {
			printer.PrintError("Type key for '%s' resolved to an empty string, skipping", typeName)
		}
	}

	return typeKeys, typeNamesFound
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

// handleDefaultExport performs a default export of all objects when no search parameters are provided
func handleDefaultExport(ctx context.Context, client *anytype.Client, targetSpace *anytype.Space, exportOpts *exportOptions, printer display.Printer) error {
	printer.PrintInfo("No search parameters provided, exporting all objects from space %s (%s)", targetSpace.Name, targetSpace.ID)

	searchParams := &anytype.SearchParams{
		Types: []string{"ot-page"}, // Default to ot-page type
		Limit: 100,
	}

	return handleSearch(ctx, client, targetSpace, searchParams, printer, exportOpts)
}

// prepareSearchParams creates and populates a SearchParams object based on command line flags.
//
// This function constructs a SearchParams object using the search criteria provided
// in the command line flags. It handles:
//
// 1. Setting the text query from the -query flag
// 2. Processing type filters from either -types (comma-separated) or -type (single)
// 3. Converting type names to type keys by querying the Anytype API
// 4. Adding tag filters from the -tags flag
//
// The function uses processTypeFilters to resolve type names to internal type keys
// that the Anytype API requires for filtering.
//
// Parameters:
//   - ctx: Context for the API request
//   - client: The initialized Anytype API client
//   - spaceID: ID of the space to search within
//   - f: Parsed command line flags containing search criteria
//   - printer: Display printer for output formatting
//
// Returns:
//   - A populated SearchParams object ready for use with the Search API
//   - Any error encountered during parameter preparation
func prepareSearchParams(ctx context.Context, client *anytype.Client, spaceID string, f *flags, printer display.Printer) (*anytype.SearchParams, error) {
	searchParams := &anytype.SearchParams{
		Query: strings.TrimSpace(f.query),
		Limit: 100,
	}

	// Process type filters (priority given to -types over -type for backwards compatibility)
	var typeKeys []string
	var typeNamesFound []string

	if f.types != "" {
		// Handle multiple types
		typeNames := strings.Split(f.types, ",")
		typeKeys, typeNamesFound = processTypeFilters(ctx, client, spaceID, typeNames, printer)
	} else if f.typeName != "" {
		// For backward compatibility: handle single type
		typeKeys, typeNamesFound = processTypeFilters(ctx, client, spaceID, []string{f.typeName}, printer)
	}

	if len(typeKeys) > 0 {
		searchParams.Types = typeKeys
		printer.PrintInfo("Filtering search results by types: %s", strings.Join(typeNamesFound, ", "))
	} else {
		printer.PrintInfo("No valid types found, proceeding with search without type filtering")
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

	return searchParams, nil
}

// setupExportOptions creates export options if export is enabled
func setupExportOptions(f *flags, printer display.Printer) *exportOptions {
	if !f.export {
		return nil
	}

	exportOpts := &exportOptions{
		enabled: true,
		path:    f.exportPath,
		format:  f.exportFormat,
	}
	printer.PrintInfo("Export enabled. Objects will be exported to %s in %s format", f.exportPath, f.exportFormat)
	return exportOpts
}

// handleSearchCase prepares search parameters and executes the search.
//
// This function handles the search case when user provides search parameters via flags.
// It first prepares the search parameters based on the command line flags, then delegates
// to handleSearch to execute the search and process the results, including any export
// operations if requested.
//
// Parameters:
//   - ctx: Context for the API request
//   - client: The initialized Anytype API client
//   - targetSpace: The space to search within
//   - f: Parsed command line flags containing search criteria
//   - printer: Display printer for output formatting
//   - exportOpts: Export options if export is enabled, or nil
//
// Returns:
//   - Any error encountered during search parameter preparation or execution
//
// This function is called from the main run() function when search parameters
// are detected in the command line flags.
func handleSearchCase(ctx context.Context, client *anytype.Client, targetSpace *anytype.Space, f *flags, printer display.Printer, exportOpts *exportOptions) error {
	// Prepare search parameters
	searchParams, err := prepareSearchParams(ctx, client, targetSpace.ID, f, printer)
	if err != nil {
		return err
	}

	// Execute search and handle results
	return handleSearch(ctx, client, targetSpace, searchParams, printer, exportOpts)
}

func run() error {
	// Parse command line flags
	f := parseFlags()

	// Check if version flag is set
	if f.version {
		versionInfo := anytype.GetVersionInfo()
		fmt.Printf("Anytype-Go v%s (API version: %s)\n", versionInfo.Version, versionInfo.APIVersion)
		return nil
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel()

	// Setup client and printer
	client, printer, err := setupClient(f)
	if err != nil {
		return err
	}

	// Setup spaces and find target space
	targetSpace, err := setupSpaces(ctx, client, f.spaceName, printer)
	if err != nil {
		return err
	}

	// Set up export options if export is enabled
	exportOpts := setupExportOptions(f, printer)

	// Determine if we need to perform a search
	hasSearchParams := f.query != "" || f.tags != "" || f.typeName != "" || f.types != ""

	if hasSearchParams {
		// Handle search case
		if err := handleSearchCase(ctx, client, targetSpace, f, printer, exportOpts); err != nil {
			return err
		}
	} else if f.export {
		// If export is enabled but no search parameters are provided,
		// perform default export
		if err := handleDefaultExport(ctx, client, targetSpace, exportOpts, printer); err != nil {
			return err
		}
	} else {
		// If no search parameters or export options are provided, perform a default search
		searchParams := &anytype.SearchParams{
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

	// Version information
	flag.BoolVar(&f.version, "version", false, "Display version information")

	flag.Parse()

	return f
}
