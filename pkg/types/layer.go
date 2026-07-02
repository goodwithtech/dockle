package types

import (
	digest "github.com/opencontainers/go-digest"
)

const LayerJSONSchemaVersion = 1

// LayerInfo is stored in cache
type LayerInfo struct {
	ID            digest.Digest `json:",omitempty"`
	SchemaVersion int
	TargetFiles   FileMap
	OpaqueDirs    []string `json:",omitempty"`
	WhiteoutFiles []string `json:",omitempty"`
}
