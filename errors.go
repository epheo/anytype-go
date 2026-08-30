package anytype

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrNotFound matches any 404 from the API and the by-key/by-name lookups
// that scan a list, so callers test one sentinel regardless of the path taken.
var ErrNotFound = errors.New("not found")

// APIError is the body every non-2xx response carries. A body that is not the
// documented shape still yields an APIError with Status set and the raw text in Message.
type APIError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	if e.Code == "" {
		return fmt.Sprintf("request failed with status %d: %s", e.Status, e.Message)
	}
	return fmt.Sprintf("request failed with status %d (%s): %s", e.Status, e.Code, e.Message)
}

func (e *APIError) Is(target error) bool {
	return target == ErrNotFound && e.Status == http.StatusNotFound
}
