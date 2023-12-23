package v1

type VersionEntity struct {
	Version       string            `json:"version"`
	LatestVersion string            `json:"latest_version,omitempty"`
	Name          string            `json:"name"`
	Type          string            `json:"type"`
	Tags          map[string]string `json:"tags"`
}
