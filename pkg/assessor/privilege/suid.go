package privilege

import (
	"fmt"
	"os"
	"strings"

	"github.com/goodwithtech/dockle/pkg/types"
	"github.com/knqyf263/fanal/extractor"
)

type PrivilegeAssessor struct{}

var ignorePaths = []string{"bin/", "usr/lib/"}

func (a PrivilegeAssessor) Assess(fileMap extractor.FileMap) ([]*types.Assessment, error) {
	var assesses []*types.Assessment

	for filename, filedata := range fileMap {
		if containIgnorePath(filename) {
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

func containIgnorePath(filename string) bool {
	for _, ignoreDir := range ignorePaths {
		if strings.Contains(filename, ignoreDir) {
			return true
		}
	}
	return false
}

func (a PrivilegeAssessor) RequiredFiles() []string {
	return []string{}
}

//const GidMode os.FileMode = 4000
func (a PrivilegeAssessor) RequiredPermissions() []os.FileMode {
	return []os.FileMode{os.ModeSocket, os.ModeSetuid}
}
