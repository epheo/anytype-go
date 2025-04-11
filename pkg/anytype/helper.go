package anytype

import (
	"net/url"
	"strconv"
)

// EncodeQueryParams encodes a map of parameters into a query string.
// It handles string, int, and []string values in the params map.
func EncodeQueryParams(params map[string]interface{}) string {
	if len(params) == 0 {
		return ""
	}

	values := url.Values{}

	for key, value := range params {
		switch v := value.(type) {
		case string:
			values.Add(key, v)
		case int:
			values.Add(key, strconv.Itoa(v))
		case []string:
			// Handle string arrays by adding multiple entries with the same key
			for _, item := range v {
				values.Add(key, item)
			}
		}
	}

	queryString := values.Encode()
	if queryString != "" {
		return "?" + queryString
	}
	return ""
}
