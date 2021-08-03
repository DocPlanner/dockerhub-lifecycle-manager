package dockerhub

import "time"

type Image struct {
	Repository string    `json:"repository"`
	LastPushed time.Time `json:"last_pushed"`
	LastPulled time.Time `json:"last_pulled"`
	Digest     string    `json:"digest"`
}

type DeleteImagesRequest struct {
	DryRun     bool        `json:"dry_run"`
	ActiveFrom time.Time   `json:"active_from,omitempty"`
	Manifests  []*Manifest `json:"manifests"`
}

type Manifest struct {
	Repository string `json:"repository"`
	Digest     string `json:"digest"`
}
