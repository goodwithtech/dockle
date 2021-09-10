package contentTrust

import (
	"os"

	deckodertypes "github.com/goodwithtech/deckoder/types"

	"github.com/goodwithtech/dockle/pkg/log"
	"github.com/goodwithtech/dockle/pkg/types"
)

var HostEnvironmentFileName = "ENVIRONMENT variable on HOST OS"

type ContentTrustAssessor struct{}

func (a ContentTrustAssessor) Assess(_ deckodertypes.FileMap) ([]*types.Assessment, error) {
	log.Logger.Debug("Scan start : DOCKER_CONTENT_TRUST")

	if os.Getenv("DOCKER_CONTENT_TRUST") != "1" {
		return []*types.Assessment{
			{
				Code:     types.UseContentTrust,
				Filename: HostEnvironmentFileName,
				Desc:     "export DOCKER_CONTENT_TRUST=1 before docker pull/build",
			},
		}, nil
	}
	return nil, nil
}

func (a ContentTrustAssessor) RequiredFiles() []string {
	return []string{}
}

func (a ContentTrustAssessor) RequiredExtensions() []string {
	return []string{}
}

func (a ContentTrustAssessor) RequiredPermissions() []os.FileMode {
	return []os.FileMode{}
}
