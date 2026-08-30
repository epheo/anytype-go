package anytype

const (
	// Follows AnyType Versioning
	Version    = "0.55.6"
	APIVersion = "2025-11-08"
)

type VersionInfo struct {
	Version    string `json:"version"`
	APIVersion string `json:"api_version"`
}

func GetVersionInfo() VersionInfo {
	return VersionInfo{
		Version:    Version,
		APIVersion: APIVersion,
	}
}
