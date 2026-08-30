package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	anytype "github.com/epheo/anytype-go"
)

// pagedServer answers path with pages of 2 items built by mk, so any walk
// below must cross a page boundary to find the last item.
func pagedServer[T any](t *testing.T, path string, n int, mk func(i int) T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			t.Errorf("unexpected request %s", r.URL.Path)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		var data []T
		for i := offset; i < offset+2 && i < n; i++ {
			data = append(data, mk(i))
		}
		page := anytype.Page[T]{Data: data}
		page.Pagination = anytype.Pagination{Total: n, Limit: 2, Offset: offset, HasMore: offset+len(data) < n}
		_ = json.NewEncoder(w).Encode(page)
	}))
}

func TestSpacesGetByName(t *testing.T) {
	names := []string{"Work", "Home", "Archive"}
	srv := pagedServer(t, "/v1/spaces", len(names), func(i int) anytype.Space {
		return anytype.Space{ID: "id-" + names[i], Name: names[i]}
	})
	defer srv.Close()

	c := anytype.NewClient(anytype.WithBaseURL(srv.URL))
	got, err := c.Spaces().GetByName(context.Background(), "Archive")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "id-Archive" {
		t.Fatalf("got %+v", got)
	}

	_, err = c.Spaces().GetByName(context.Background(), "nope")
	if !errors.Is(err, anytype.ErrSpaceNotFound) || !errors.Is(err, anytype.ErrNotFound) {
		t.Fatalf("err = %v, want ErrSpaceNotFound", err)
	}
}

func TestSpaceSearchAll(t *testing.T) {
	srv := pagedServer(t, "/v1/spaces/s/search", 5, func(i int) anytype.Object {
		return anytype.Object{ID: strconv.Itoa(i)}
	})
	defer srv.Close()

	c := anytype.NewClient(anytype.WithBaseURL(srv.URL))
	var ids []string
	for obj, err := range c.Space("s").SearchAll(context.Background(), anytype.SearchRequest{Query: "x"}) {
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, obj.ID)
	}
	if len(ids) != 5 || ids[4] != "4" {
		t.Fatalf("got %v", ids)
	}
}
