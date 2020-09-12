package privilege

import (
	"fmt"
	"os"

	"github.com/goodwithtech/deckoder/extractor"
	"github.com/goodwithtech/dockle/pkg/types"
)

type PrivilegeAssessor struct{}

func (a PrivilegeAssessor) Assess(fileMap extractor.FileMap) ([]*types.Assessment, error) {
	var assesses []*types.Assessment

	for filename, filedata := range fileMap {
		if filedata.FileMode&os.ModeSetuid != 0 {
			assesses = append(
				assesses,
				&types.Assessment{
					Code:     types.CheckSuidGuid,
					Filename: filename,
					Desc:     fmt.Sprintf("setuid file: %s %s", filename, filedata.FileMode),
				})
		}
		if filedata.FileMode&os.ModeSetgid != 0 {
			assesses = append(
				assesses,
				&types.Assessment{
					Code:     types.CheckSuidGuid,
					Filename: filename,
					Desc:     fmt.Sprintf("setgid file: %s %s", filename, filedata.FileMode),
				})
		}

	}
	return assesses, nil
}

func (a PrivilegeAssessor) RequiredFiles() []string {
	return []string{}
}

//const GidMode os.FileMode = 4000
func (a PrivilegeAssessor) RequiredPermissions() []os.FileMode {
	return []os.FileMode{os.ModeSetgid, os.ModeSetuid}
}
