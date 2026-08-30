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

// typesByKey serves the list endpoint in pages of 2 and rejects any
// /types/{x} lookup, mirroring the real API which only accepts IDs there.
func typesByKey(t *testing.T, keys ...string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/spaces/s/types" {
			t.Errorf("unexpected request %s", r.URL.Path)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		limit := 2
		var data []anytype.Type
		for i := offset; i < offset+limit && i < len(keys); i++ {
			data = append(data, anytype.Type{ID: "id-" + keys[i], Key: keys[i]})
		}
		page := anytype.Page[anytype.Type]{Data: data}
		page.Pagination = anytype.Pagination{
			Total: len(keys), Limit: limit, Offset: offset,
			HasMore: offset+len(data) < len(keys),
		}
		_ = json.NewEncoder(w).Encode(page)
	}))
}

func TestTypesGetByKeyWalksList(t *testing.T) {
	srv := typesByKey(t, "page", "task", "note")
	defer srv.Close()

	c := anytype.NewClient(anytype.WithBaseURL(srv.URL))
	got, err := c.Space("s").Types().Get(context.Background(), "note")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "id-note" {
		t.Fatalf("got %+v, want id-note", got)
	}
}

func TestTypesGetUnknownKey(t *testing.T) {
	srv := typesByKey(t, "page")
	defer srv.Close()

	c := anytype.NewClient(anytype.WithBaseURL(srv.URL))
	_, err := c.Space("s").Types().Get(context.Background(), "missing")
	if !errors.Is(err, anytype.ErrTypeNotFound) {
		t.Fatalf("err = %v, want ErrTypeNotFound", err)
	}
}
