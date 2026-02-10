// Package anytype provides a Go SDK for interacting with the Anytype API.
package anytype

// Version information
const (
	// Version is the current version of the anytype-go SDK.
	// This follows Semantic Versioning (https://semver.org/):
	// MAJOR.MINOR.PATCH where:
	// - MAJOR version changes with incompatible API changes
	// - MINOR version adds functionality in a backwards compatible manner
	// - PATCH version makes backwards compatible bug fixes
	Version = "0.53.1"
)

// VersionInfo holds detailed version information
type VersionInfo struct {
	// Version is the semantic version of the SDK
	Version string `json:"version"`
	// APIVersion is the Anytype API version this client works with
	APIVersion string `json:"api_version"`
}

// GetVersionInfo returns version information for the SDK
func GetVersionInfo() VersionInfo {
	return VersionInfo{
		Version: Version,
	}
}
