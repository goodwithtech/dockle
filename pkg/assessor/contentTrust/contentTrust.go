package contentTrust

import (
	"os"

	"github.com/goodwithtech/deckoder/extractor"
	"github.com/goodwithtech/dockle/pkg/log"
	"github.com/goodwithtech/dockle/pkg/types"
)

type ContentTrustAssessor struct{}

func (a ContentTrustAssessor) Assess(fileMap extractor.FileMap) ([]*types.Assessment, error) {
	log.Logger.Debug("Scan start : DOCKER_CONTENT_TRUST")

	if os.Getenv("DOCKER_CONTENT_TRUST") != "1" {
		return []*types.Assessment{
			{
				Type:     types.UseContentTrust,
				Filename: "ENVIRONMENT variable",
				Desc:     "export DOCKER_CONTENT_TRUST=1 before docker pull/build",
			},
		}, nil
	}
	return nil, nil
}

func (a ContentTrustAssessor) RequiredFiles() []string {
	return []string{}
}

func (a ContentTrustAssessor) RequiredPermissions() []os.FileMode {
	return []os.FileMode{}
}
