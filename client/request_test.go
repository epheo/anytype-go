package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	anytype "github.com/epheo/anytype-go"
)

func TestAPIErrorParsed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":"object_not_found","message":"Resource not found","object":"error","status":404}`))
	}))
	defer srv.Close()

	c := anytype.NewClient(anytype.WithBaseURL(srv.URL))
	_, err := c.Space("s").Type("nope").Get(context.Background())

	var apiErr *anytype.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("err = %T, want *APIError", err)
	}
	if apiErr.Status != 404 || apiErr.Code != "object_not_found" || apiErr.Message != "Resource not found" {
		t.Fatalf("got %+v", apiErr)
	}
	if !errors.Is(err, anytype.ErrNotFound) {
		t.Fatalf("404 should match ErrNotFound")
	}
}

func TestAPIErrorRawBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("<html>bad request</html>"))
	}))
	defer srv.Close()

	c := anytype.NewClient(anytype.WithBaseURL(srv.URL))
	_, err := c.Spaces().List(context.Background())

	var apiErr *anytype.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("err = %T, want *APIError", err)
	}
	if apiErr.Status != 400 || apiErr.Code != "" || apiErr.Message != "<html>bad request</html>" {
		t.Fatalf("got %+v", apiErr)
	}
	if errors.Is(err, anytype.ErrNotFound) {
		t.Fatalf("400 must not match ErrNotFound")
	}
}

func TestListLookupsMatchErrNotFound(t *testing.T) {
	srv := typesByKey(t, "page")
	defer srv.Close()

	c := anytype.NewClient(anytype.WithBaseURL(srv.URL))
	_, err := c.Space("s").Types().Get(context.Background(), "missing")
	if !errors.Is(err, anytype.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}
