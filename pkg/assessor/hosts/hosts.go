package hosts

import (
	"os"

	"github.com/Portshift/dockle/pkg/log"
	"github.com/Portshift/dockle/pkg/types"
)

type HostsAssessor struct{}

func (a HostsAssessor) Assess(_ *types.ImageData) ([]*types.Assessment, error) {
	log.Logger.Debug("Start scan : /etc/hosts")

	assesses := []*types.Assessment{}
	// TODO : check hosts setting
	return assesses, nil
}

func (a HostsAssessor) RequiredFiles() []string {
	return []string{"etc/hosts"}
}

func (a HostsAssessor) RequiredExtensions() []string {
	return []string{}
}

func (a HostsAssessor) RequiredPermissions() []os.FileMode {
	return []os.FileMode{}
}
