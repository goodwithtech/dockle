package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	deckodertypes "github.com/goodwithtech/deckoder/types"
	"github.com/goodwithtech/deckoder/utils"

	"github.com/goodwithtech/dockle/pkg/log"
	"github.com/goodwithtech/dockle/pkg/types"
)

var (
	reqFiles = []string{"Dockerfile", "docker-compose.yml", ".vimrc", ".DS_Store"}
	// Directory ends "/" separator
	reqDirs            = []string{".cache/", ".aws/", ".azure/", ".gcp/", ".git/", ".vscode/", ".idea/", ".npm/"}
	uncontrollableDirs = []string{"node_modules/", "vendor/"}
	detectedDir        = map[string]struct{}{}
)

type CacheAssessor struct{}

func (a CacheAssessor) Assess(fileMap deckodertypes.FileMap) ([]*types.Assessment, error) {
	log.Logger.Debug("Start scan : cache files")
	assesses := []*types.Assessment{}
	for filename := range fileMap {
		fileBase := filepath.Base(filename)
		dirName := filepath.Dir(filename)
		dirBase := filepath.Base(dirName)

		// match Directory
		if utils.StringInSlice(dirBase+"/", reqDirs) || utils.StringInSlice(dirName+"/", reqDirs) {
			if _, ok := detectedDir[dirName]; ok {
				continue
			}
			detectedDir[dirName] = struct{}{}

			// Skip uncontrollable dependency directory e.g) npm : node_modules, php: composer
			if inIgnoreDir(filename) {
				continue
			}

			assesses = append(
				assesses,
				&types.Assessment{
					Code:     types.InfoDeletableFiles,
					Filename: dirName,
					Desc:     fmt.Sprintf("Suspicious directory : %s ", dirName),
				})

		}

		// match File
		if utils.StringInSlice(filename, reqFiles) || utils.StringInSlice(fileBase, reqFiles) {
			assesses = append(
				assesses,
				&types.Assessment{
					Code:     types.InfoDeletableFiles,
					Filename: filename,
					Desc:     fmt.Sprintf("unnecessary file : %s ", filename),
				})
		}
	}
	return assesses, nil
}

// check and register uncontrollable directory e.g) npm : node_modules, php: composer
func inIgnoreDir(filename string) bool {
	for _, ignoreDir := range uncontrollableDirs {
		if strings.Contains(filename, ignoreDir) {
			return true
		}
	}
	return false
}

func (a CacheAssessor) RequiredFiles() []string {
	return append(reqFiles, reqDirs...)
}

func (a CacheAssessor) RequiredExtensions() []string {
	return []string{}
}

func (a CacheAssessor) RequiredPermissions() []os.FileMode {
	return []os.FileMode{}
}
