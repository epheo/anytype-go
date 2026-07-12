package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	anytype "github.com/epheo/anytype-go"
)

// pagedObjects serves /v1/spaces/{id}/objects as `total` objects split into
// pages of the requested limit, recording the offsets it was asked for.
func pagedObjects(total int, offsets *[]int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit == 0 {
			limit = 2
		}
		*offsets = append(*offsets, offset)

		var data []anytype.Object
		for i := offset; i < offset+limit && i < total; i++ {
			data = append(data, anytype.Object{ID: strconv.Itoa(i)})
		}
		page := anytype.Page[anytype.Object]{Data: data}
		page.Pagination = anytype.Pagination{
			Total:   total,
			Limit:   limit,
			Offset:  offset,
			HasMore: offset+len(data) < total,
		}
		_ = json.NewEncoder(w).Encode(page)
	}))
}

func TestListReturnsPage(t *testing.T) {
	var offsets []int
	srv := pagedObjects(5, &offsets)
	defer srv.Close()

	c := anytype.NewClient(anytype.WithBaseURL(srv.URL))
	page, err := c.Space("s").Objects().List(context.Background(), anytype.WithLimit(2))
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Data) != 2 {
		t.Fatalf("page size = %d, want 2", len(page.Data))
	}
	if page.Pagination.Total != 5 || !page.Pagination.HasMore {
		t.Fatalf("pagination = %+v, want total 5 with more", page.Pagination)
	}
}

func TestAllWalksEveryPage(t *testing.T) {
	var offsets []int
	srv := pagedObjects(5, &offsets)
	defer srv.Close()

	c := anytype.NewClient(anytype.WithBaseURL(srv.URL))
	var ids []string
	for obj, err := range c.Space("s").Objects().All(context.Background(), anytype.WithLimit(2)) {
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, obj.ID)
	}

	if len(ids) != 5 {
		t.Fatalf("iterated %d objects, want 5 (%v)", len(ids), ids)
	}
	if want := []int{0, 2, 4}; !equalInts(offsets, want) {
		t.Fatalf("requested offsets = %v, want %v", offsets, want)
	}
}

func TestAllStopsOnBreak(t *testing.T) {
	var offsets []int
	srv := pagedObjects(100, &offsets)
	defer srv.Close()

	c := anytype.NewClient(anytype.WithBaseURL(srv.URL))
	seen := 0
	for _, err := range c.Space("s").Objects().All(context.Background(), anytype.WithLimit(2)) {
		if err != nil {
			t.Fatal(err)
		}
		seen++
		if seen == 3 {
			break
		}
	}

	if seen != 3 {
		t.Fatalf("saw %d objects before break, want 3", seen)
	}
	// Breaking mid-second-page must not fetch a third page.
	if len(offsets) > 2 {
		t.Fatalf("fetched %d pages after early break, want <= 2", len(offsets))
	}
}

func TestAllYieldsError(t *testing.T) {
	// 400 is not retryable, so the error surfaces without backoff delay.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	c := anytype.NewClient(anytype.WithBaseURL(srv.URL))
	var gotErr error
	iterations := 0
	for _, err := range c.Space("s").Objects().All(context.Background()) {
		iterations++
		gotErr = err
	}
	if gotErr == nil {
		t.Fatal("expected an error from the iterator, got nil")
	}
	// The error is yielded exactly once, then iteration stops.
	if iterations != 1 {
		t.Fatalf("iterated %d times on error, want 1", iterations)
	}
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
