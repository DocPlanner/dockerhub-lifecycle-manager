package dockerhub

import "time"

type Image struct {
	Repository string    `json:"repository"`
	LastPushed time.Time `json:"last_pushed"`
	LastPulled time.Time `json:"last_pulled"`
	Digest     string    `json:"digest"`
	Tags       []*Tags   `json:"tags"`
}

type DeleteImagesRequest struct {
	DryRun         bool              `json:"dry_run"`
	ActiveFrom     time.Time         `json:"active_from,omitempty"`
	Manifests      []*Manifest       `json:"manifests"`
	IgnoreWarnings []*IgnoreWarnings `json:"ignore_warnings"`
}

type Manifest struct {
	Repository string `json:"repository"`
	Digest     string `json:"digest"`
}
type IgnoreWarnings struct {
	Repository string   `json:"repository"`
	Digest     string   `json:"digest"`
	Warning    string   `json:"warning"`
	Tags       []string `json:"tags"`
}

type Tags struct {
	Tag       string `json:"tag"`
	IsCurrent bool   `json:"is_current"`
}

type DeletedImagesResponse struct {
	Metrics *Metrics `json:"metrics"`
}

type Metrics struct {
	ManifestDeletes int `json:"manifest_deletes"`
	ManifestErrors  int `json:"manifest_errors"`
	TagDeletes      int `json:"tag_deletes"`
	TagErrors       int `json:"tag_errors"`
}
