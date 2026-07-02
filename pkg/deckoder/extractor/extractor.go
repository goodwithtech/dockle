package extractor

import (
	"context"

	digest "github.com/opencontainers/go-digest"

	"github.com/goodwithtech/dockle/pkg/deckoder/types"
)

type Extractor interface {
	ImageName() (imageName string)
	ImageID() (imageID digest.Digest)
	ConfigBlob(ctx context.Context) (configBlob []byte, err error)
	LayerIDs() (layerIDs []string)
	ExtractLayerFiles(ctx context.Context, dig digest.Digest, filterFunc types.FilterFunc) (files types.FileMap, opqDirs, whFiles []string, err error)
}
