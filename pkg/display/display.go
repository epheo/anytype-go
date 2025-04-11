package display

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/epheo/anyblog/pkg/anytype"
)

// Output formats
const (
	formatJSON = "json"
	formatText = "text"
)

// Printer handles formatted output
type Printer struct {
	writer    io.Writer
	format    string
	useColors bool
}

// NewPrinter creates a new Printer instance
func NewPrinter(format string, useColors bool) *Printer {
	return &Printer{
		writer:    os.Stdout,
		format:    format,
		useColors: useColors,
	}
}

// SetWriter changes the output writer
func (p *Printer) SetWriter(w io.Writer) {
	p.writer = w
}

// PrintJSON formats and prints JSON data
func (p *Printer) PrintJSON(label string, data interface{}) error {
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

// PrintSpaces displays available spaces
func (p *Printer) PrintSpaces(spaces []anytype.Space) error {
	if p.format == formatJSON {
		return p.PrintJSON("Available spaces", spaces)
	}

	w := tabwriter.NewWriter(p.writer, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "\nAvailable spaces:")
	fmt.Fprintln(w, "ID\tName\tIcon")
	fmt.Fprintln(w, "----\t----\t----")

	for _, space := range spaces {
		icon := space.Icon
		if icon == "" {
			icon = "📄"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", space.ID, space.Name, icon)
	}

	return w.Flush()
}

// PrintObjects displays objects with their properties
func (p *Printer) PrintObjects(label string, objects []anytype.Object) error {
	if p.format == formatJSON {
		return p.PrintJSON(label, objects)
	}

	w := tabwriter.NewWriter(p.writer, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "\n%s:\n", label)
	fmt.Fprintln(w, "ID\tName\tType\tTags")
	fmt.Fprintln(w, "----\t----\t----\t----")

	for _, obj := range objects {
		tags := "none"
		if len(obj.Tags) > 0 {
			tags = fmt.Sprintf("%v", obj.Tags)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", obj.ID, obj.Name, obj.Type, tags)
	}

	return w.Flush()
}

// PrintError prints an error message
func (p *Printer) PrintError(format string, args ...interface{}) {
	prefix := "Error: "
	if p.useColors {
		prefix = "\033[31mError:\033[0m "
	}
	fmt.Fprintf(p.writer, prefix+format+"\n", args...)
}

// PrintSuccess prints a success message
func (p *Printer) PrintSuccess(format string, args ...interface{}) {
	prefix := "Success: "
	if p.useColors {
		prefix = "\033[32mSuccess:\033[0m "
	}
	fmt.Fprintf(p.writer, prefix+format+"\n", args...)
}

// PrintInfo prints an informational message
func (p *Printer) PrintInfo(format string, args ...interface{}) {
	prefix := "Info: "
	if p.useColors {
		prefix = "\033[34mInfo:\033[0m "
	}
	fmt.Fprintf(p.writer, prefix+format+"\n", args...)
}
