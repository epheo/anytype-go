// Package client implements the client interfaces defined in the anytype package
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"
	"net/url"
	"path"
	"strconv"

	anytype "github.com/epheo/anytype-go"
)

// withListParams appends limit/offset query params to an endpoint. Zero values
// are omitted so the server applies its own defaults.
func withListParams(endpoint string, opts []anytype.ListOption) string {
	o := anytype.ApplyListOptions(opts...)
	params := url.Values{}
	if o.Limit > 0 {
		params.Set("limit", strconv.Itoa(o.Limit))
	}
	if o.Offset > 0 {
		params.Set("offset", strconv.Itoa(o.Offset))
	}
	if encoded := params.Encode(); encoded != "" {
		return endpoint + "?" + encoded
	}
	return endpoint
}

// paginate turns any List-shaped fetch into an iterator that walks every page,
// advancing the offset until the server reports no more results. A fetch error
// is yielded once with the zero value, then iteration stops.
func paginate[T any](ctx context.Context, fetch func(context.Context, ...anytype.ListOption) (*anytype.Page[T], error), base ...anytype.ListOption) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		var zero T
		offset := 0
		for {
			opts := append(append([]anytype.ListOption(nil), base...), anytype.WithOffset(offset))
			page, err := fetch(ctx, opts...)
			if err != nil {
				yield(zero, err)
				return
			}
			for _, item := range page.Data {
				if !yield(item, nil) {
					return
				}
			}
			// len==0 guards against a server that reports HasMore with no rows.
			if !page.Pagination.HasMore || len(page.Data) == 0 {
				return
			}
			offset += len(page.Data)
		}
	}
}

func (c *ClientImpl) newRequest(ctx context.Context, method, urlPath string, body interface{}) (*http.Request, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}

	// Separate query string from path before joining
	parsedPath, err := url.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "/v1", parsedPath.Path)
	u.RawQuery = parsedPath.RawQuery

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Anytype-Version", anytype.APIVersion)

	if c.appKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.appKey))
	}

	return req, nil
}

func (c *ClientImpl) doRequest(req *http.Request, result interface{}) error {
	resp, err := c.doer.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Read body for error details to provide actionable error messages
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return err
		}
	}

	return nil
}
