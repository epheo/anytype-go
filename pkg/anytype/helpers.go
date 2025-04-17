package anytype

import "fmt"

// getMapKeys returns the keys of a map as a string slice
// useful for debugging API responses
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// debugJSON logs the structure of a JSON response
// useful for understanding API response formats
func debugJSON(v interface{}) string {
	return fmt.Sprintf("%+v", v)
}
