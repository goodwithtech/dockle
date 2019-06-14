package privilege

import (
	"fmt"
	"os"
	"strings"

	"github.com/goodwithtech/dockle/pkg/types"
	"github.com/knqyf263/fanal/extractor"
)

type privilegeAssessor struct{}

func (a privilegeAssessor) Assess(fileMap extractor.FileMap) ([]*types.Assessment, error) {
	var assesses []*types.Assessment

	for filename, filedata := range fileMap {
		if strings.Contains(filename, "bin/") {
			continue
		}
		if filedata.FileMode&os.ModeSetuid != 0 {
			assesses = append(
				assesses,
				&types.Assessment{
					Type:     types.RemoveSetuidSetgid,
					Filename: filename,
					Desc:     fmt.Sprintf("Found setuid file: %s %s", filename, filedata.FileMode),
				})
		}
		if filedata.FileMode&os.ModeSetgid != 0 {
			assesses = append(
				assesses,
				&types.Assessment{
					Type:     types.RemoveSetuidSetgid,
					Filename: filename,
					Desc:     fmt.Sprintf("Found setuid file: %s %s", filename, filedata.FileMode),
				})
		}

	}
	return assesses, nil
}

func (a privilegeAssessor) RequiredFiles() []string {
	return []string{}
}

//const GidMode os.FileMode = 4000

func (a privilegeAssessor) RequiredPermissions() []os.FileMode {
	return []os.FileMode{os.ModeSocket, os.ModeSetuid}
}
