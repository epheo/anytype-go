package display

import (
	"context"
	"io"

	"github.com/epheo/anyblog/pkg/anytype"
)

// Printer defines the interface for output formatting
type Printer interface {
	PrintJSON(label string, data interface{}) error
	PrintSpaces(spaces []anytype.Space) error
	PrintObjects(label string, objects []anytype.Object, client *anytype.Client, ctx context.Context) error
	PrintError(format string, args ...interface{})
	PrintSuccess(format string, args ...interface{})
	PrintInfo(format string, args ...interface{})
	SetWriter(w io.Writer)
}
