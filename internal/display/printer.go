package display

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/epheo/anytype-go/pkg/anytype"
	"github.com/olekukonko/tablewriter"
)

// printer implements the Printer interface
type printer struct {
	writer    io.Writer
	format    string
	useColors bool
	debug     bool
	logLevel  LogLevel
}

// NewPrinter creates a new Printer instance
func NewPrinter(format string, useColors bool, debug bool) *printer {
	logLevel := LogLevelInfo
	if debug {
		logLevel = LogLevelDebug
	}
	return &printer{
		writer:    os.Stdout,
		format:    format,
		useColors: useColors,
		debug:     debug,
		logLevel:  logLevel,
	}
}

// SetWriter changes the output writer
func (p *printer) SetWriter(w io.Writer) {
	p.writer = w
}

// SetLogLevel sets the current log level
func (p *printer) SetLogLevel(level LogLevel) {
	p.logLevel = level
}

// GetLogLevel returns the current log level
func (p *printer) GetLogLevel() LogLevel {
	return p.logLevel
}

// PrintJSON formats and prints JSON data
func (p *printer) PrintJSON(label string, data interface{}) error {
	var rawJSON []byte
	var err error

	switch v := data.(type) {
	case []byte:
		rawJSON = v
	default:
		rawJSON, err = json.Marshal(v)
		if err != nil {
			return fmt.Errorf("error marshaling JSON: %w", err)
		}
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, rawJSON, "", "  "); err != nil {
		return fmt.Errorf("error formatting JSON: %w", err)
	}

	// In JSON format mode, only output the raw JSON without any labels
	if p.format == formatJSON {
		fmt.Fprintln(p.writer, prettyJSON.String())
	} else {
		fmt.Fprintf(p.writer, "\n%s:\n%s\n", label, prettyJSON.String())
	}
	return nil
}

// handleError formats and prints an error message appropriately based on error type
func (p *printer) handleError(err error) {
	if err == nil {
		return
	}

	var apiErr *anytype.Error
	if errors.As(err, &apiErr) {
		// Handle API-specific errors with more context
		if apiErr.StatusCode != 0 {
			p.PrintError("%s (Status: %d)", apiErr.Message, apiErr.StatusCode)
		} else {
			p.PrintError("%s", apiErr.Message)
		}
		return
	}

	// Handle standard errors
	p.PrintError("%v", err)
}

// PrintSpaces displays available spaces with better error handling
func (p *printer) PrintSpaces(spaces []anytype.Space) error {
	if spaces == nil {
		return anytype.ErrEmptyResponse
	}

	if p.format == formatJSON {
		return p.PrintJSON("Available spaces", spaces)
	}

	table := tablewriter.NewWriter(p.writer)
	table.SetHeader([]string{"Name", "Members"})
	setupTable(table)

	for _, space := range spaces {
		if err := space.Validate(); err != nil {
			p.handleError(err)
			continue
		}

		// Format members list
		members := "-"
		if len(space.Members) > 0 {
			memberStrs := make([]string, 0, len(space.Members))
			for _, member := range space.Members {
				memberStr := member.Name
				if member.Role != "" {
					memberStr += fmt.Sprintf(" (%s)", member.Role)
				}
				memberStrs = append(memberStrs, memberStr)
			}
			members = strings.Join(memberStrs, ", ")
			if len(members) > maxMembersLength {
				members = members[:maxMembersLength-3] + "..."
			}
		}

		table.Append([]string{space.Name, members})
	}

	fmt.Fprintln(p.writer, "\nAvailable spaces:")
	table.Render()
	return nil
}

// PrintObjects displays objects with improved error handling
func (p *printer) PrintObjects(label string, objects []anytype.Object, client *anytype.Client, ctx context.Context) error {
	if objects == nil {
		return anytype.ErrEmptyResponse
	}

	// Pre-fetch type information for all objects in one go if client is available
	if client != nil && len(objects) > 0 {
		// We only need to pre-fetch once per space
		spacesSeen := make(map[string]bool)

		for _, obj := range objects {
			if obj.SpaceID != "" && !spacesSeen[obj.SpaceID] {
				// This will populate the type cache for this space with a single API call
				client.GetTypeName(ctx, obj.SpaceID, "dummy-key")
				spacesSeen[obj.SpaceID] = true
			}
		}
	}

	if p.format == formatJSON {
		return p.PrintJSON(label, objects)
	}

	table := tablewriter.NewWriter(p.writer)
	table.SetHeader([]string{"Name", "Type", "Layout", "Tags"})
	setupTable(table)

	for _, obj := range objects {
		if err := obj.Validate(); err != nil {
			p.handleError(err)
			continue
		}

		name := obj.Name
		if name == "" {
			name = "<no name>"
		}

		// Format the display name with fixed-width icon space
		var displayName string
		// Normalize the icon (ensure it's not nil and handle special cases)
		var iconStr string
		if obj.Icon != nil {
			if obj.Icon.Emoji != "" {
				iconStr = obj.Icon.Emoji
			} else if obj.Icon.Name != "" {
				iconStr = obj.Icon.Name
			}
		}
		// Use GetPaddedIcon to ensure consistent spacing regardless of icon presence or type
		paddedIcon := GetPaddedIcon(iconStr, iconFixedWidth)
		displayName = fmt.Sprintf("%s%s", paddedIcon, name)

		// Truncate name if too long
		if len(displayName) > maxNameLength {
			displayName = displayName[:maxNameLength-3] + "..."
		}

		// Get friendly type name
		var typeNameStr string
		if obj.Type != nil {
			typeNameStr = obj.Type.Name
			if client != nil && obj.Type.Key != "" {
				typeNameStr = client.GetTypeName(ctx, obj.SpaceID, obj.Type.Key)
			}
		} else {
			typeNameStr = "Unknown"
		}

		layout := obj.Layout
		if layout == "" {
			layout = "-"
		}

		// Format tags
		tags := "-"
		if len(obj.Tags) > 0 {
			tags = strings.Join(obj.Tags, ", ")
			if len(tags) > maxTagsLength {
				tags = tags[:maxTagsLength-3] + "..."
			}
		}

		table.Append([]string{displayName, typeNameStr, layout, tags})
	}

	fmt.Fprintf(p.writer, "\n%s:\n", label)
	table.Render()

	if p.debug {
		fmt.Fprintf(p.writer, "\nDebug: Raw objects: %+v\n", objects)
	}

	return nil
}

// PrintError prints an error message (always enabled unless in JSON mode)
func (p *printer) PrintError(format string, args ...interface{}) {
	// Skip logging in JSON mode
	if p.format == formatJSON {
		return
	}

	prefix := "Error: "
	if p.useColors {
		prefix = colorRed + "Error:" + colorReset + " "
	}
	fmt.Fprintf(p.writer, prefix+format+"\n", args...)
}

// PrintSuccess prints a success message if info logging is enabled
func (p *printer) PrintSuccess(format string, args ...interface{}) {
	// Skip logging in JSON mode
	if p.format == formatJSON {
		return
	}

	if p.logLevel >= LogLevelInfo {
		prefix := "Success: "
		if p.useColors {
			prefix = colorGreen + "Success:" + colorReset + " "
		}
		fmt.Fprintf(p.writer, prefix+format+"\n", args...)
	}
}

// PrintInfo prints an informational message if info logging is enabled
func (p *printer) PrintInfo(format string, args ...interface{}) {
	// Skip logging in JSON mode
	if p.format == formatJSON {
		return
	}

	if p.logLevel >= LogLevelInfo {
		prefix := "Info: "
		if p.useColors {
			prefix = colorBlue + "Info:" + colorReset + " "
		}
		fmt.Fprintf(p.writer, prefix+format+"\n", args...)
	}
}

// PrintDebug prints a debug message if debug logging is enabled
func (p *printer) PrintDebug(format string, args ...interface{}) {
	// Skip logging in JSON mode
	if p.format == formatJSON {
		return
	}

	if p.logLevel >= LogLevelDebug {
		prefix := "Debug: "
		if p.useColors {
			prefix = colorCyan + "Debug:" + colorReset + " "
		}
		fmt.Fprintf(p.writer, prefix+format+"\n", args...)
	}
}
