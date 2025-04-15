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

	"github.com/epheo/anyblog/pkg/anytype"
	"github.com/olekukonko/tablewriter"
)

// printer implements the Printer interface
type printer struct {
	writer    io.Writer
	format    string
	useColors bool
	debug     bool
}

// NewPrinter creates a new Printer instance
func NewPrinter(format string, useColors bool, debug bool) Printer {
	return &printer{
		writer:    os.Stdout,
		format:    format,
		useColors: useColors,
		debug:     debug,
	}
}

// SetWriter changes the output writer
func (p *printer) SetWriter(w io.Writer) {
	p.writer = w
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

	fmt.Fprintf(p.writer, "\n%s:\n%s\n", label, prettyJSON.String())
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
		if obj.Icon != "" {
			name = obj.Icon + " " + name
		}

		// Truncate name if too long
		if len(name) > maxNameLength {
			name = name[:maxNameLength-3] + "..."
		}

		// Get friendly type name
		typeName := obj.Type
		if client != nil {
			typeName = client.GetTypeName(ctx, obj.SpaceID, obj.Type)
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

		table.Append([]string{name, typeName, layout, tags})
	}

	fmt.Fprintf(p.writer, "\n%s:\n", label)
	table.Render()

	if p.debug {
		fmt.Fprintf(p.writer, "\nDebug: Raw objects: %+v\n", objects)
	}

	return nil
}

// PrintError prints an error message
func (p *printer) PrintError(format string, args ...interface{}) {
	prefix := "Error: "
	if p.useColors {
		prefix = colorRed + "Error:" + colorReset + " "
	}
	fmt.Fprintf(p.writer, prefix+format+"\n", args...)
}

// PrintSuccess prints a success message
func (p *printer) PrintSuccess(format string, args ...interface{}) {
	prefix := "Success: "
	if p.useColors {
		prefix = colorGreen + "Success:" + colorReset + " "
	}
	fmt.Fprintf(p.writer, prefix+format+"\n", args...)
}

// PrintInfo prints an informational message
func (p *printer) PrintInfo(format string, args ...interface{}) {
	prefix := "Info: "
	if p.useColors {
		prefix = colorBlue + "Info:" + colorReset + " "
	}
	fmt.Fprintf(p.writer, prefix+format+"\n", args...)
}
