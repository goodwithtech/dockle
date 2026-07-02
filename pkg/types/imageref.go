package types

import (
	"time"

	digest "github.com/opencontainers/go-digest"
)

type FilePath string

type ImageReference struct {
	Name     string // image name or tar file name
	ID       digest.Digest
	LayerIDs []string
}

type ImageDetail struct {
	Files FileMap
}

// ImageInfo is stored in cache
type ImageInfo struct {
	SchemaVersion int
	Architecture  string
	Created       time.Time
	DockerVersion string
	OS            string
}

// LayerInfo is stored in cache
type LayerInfo struct {
	ID            digest.Digest `json:",omitempty"`
	SchemaVersion int
	TargetFiles   FileMap
	OpaqueDirs    []string `json:",omitempty"`
	WhiteoutFiles []string `json:",omitempty"`
}
