package hosts

import (
	"github.com/goodwithtech/docker-guard/pkg/log"
	"github.com/goodwithtech/docker-guard/pkg/types"
	"github.com/knqyf263/fanal/extractor"
)

type HostsAssessor struct{}

func (a HostsAssessor) Assess(fileMap extractor.FileMap) ([]types.Assessment, error) {
	log.Logger.Debug("Start scan : /etc/hosts")

	assesses := []types.Assessment{}
	// TODO : check hosts setting
	return assesses, nil
}

func (a HostsAssessor) RequiredFiles() []string {
	return []string{"etc/hosts"}
}
