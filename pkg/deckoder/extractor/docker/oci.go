package docker

import (
	"context"

	"github.com/goodwithtech/dockle/pkg/deckoder/extractor/image"
	"github.com/goodwithtech/dockle/pkg/types"
)

// NewOCIDirExtractor builds an Extractor for an OCI image layout directory
// (the "oci-path" format produced by tools such as Bazel's rules_oci and
// `skopeo copy ... oci:...`). It relies on the containers/image "oci:"
// transport, which reads the layout's index.json / blobs directory.
func NewOCIDirExtractor(ctx context.Context, dirPath string, option types.DockerOption) (Extractor, func(), error) {
	ref := image.Reference{Name: dirPath, IsFile: true}
	transports := []string{"oci:"}
	return newDockerExtractor(ctx, ref, transports, option)
}
