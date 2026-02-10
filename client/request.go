// Package client implements the client interfaces defined in the anytype package
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	anytype "github.com/epheo/anytype-go"
)

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
