package scanner

import (
	"archive/tar"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goodwithtech/deckoder/utils"

	deckodertypes "github.com/goodwithtech/deckoder/types"

	"github.com/goodwithtech/dockle/pkg/types"

	"github.com/goodwithtech/dockle/pkg/assessor"

	"github.com/goodwithtech/deckoder/analyzer"
	"github.com/goodwithtech/deckoder/extractor"
	"golang.org/x/crypto/ssh/terminal"
)

func ScanImage(imageName, filePath string, dockerOption deckodertypes.DockerOption) (assessments []*types.Assessment, err error) {
	ctx := context.Background()
	var files extractor.FileMap
	filterFunc := createPathPermissionFilterFunc(assessor.LoadRequiredFiles(), assessor.LoadRequiredPermissions())
	if imageName != "" {
		files, err = analyzer.Analyze(ctx, imageName, filterFunc, dockerOption)
		if err != nil {
			return nil, fmt.Errorf("failed to analyze image: %w", err)
		}
	} else if filePath != "" {
		rc, err := openStream(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open stream: %w", err)
		}

		files, err = analyzer.AnalyzeFromFile(ctx, rc, filterFunc)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, types.ErrSetImageOrFile
	}

	assessments = assessor.GetAssessments(files)
	return assessments, nil
}

func createPathPermissionFilterFunc(filenames []string, permissions []os.FileMode) deckodertypes.FilterFunc {
	requiredDirNames := []string{}
	requiredFileNames := []string{}
	for _, filename := range filenames {
		if filename[len(filename)-1] == '/' {
			// if filename end "/", it is directory and requiredDirNames removes last "/"
			requiredDirNames = append(requiredDirNames, filepath.Clean(filename))
		} else {
			requiredFileNames = append(requiredFileNames, filename)
		}
	}

	return func(h *tar.Header) (bool, error) {
		filePath := filepath.Clean(h.Name)
		fileName := filepath.Base(filePath)
		if utils.StringInSlice(filePath, requiredFileNames) || utils.StringInSlice(fileName, requiredFileNames) {
			return true, nil
		}

		fileDir := filepath.Dir(filePath)
		fileDirBase := filepath.Base(fileDir)
		if utils.StringInSlice(fileDir, requiredDirNames) || utils.StringInSlice(fileDirBase, requiredDirNames) {
			return true, nil
		}

		fi := h.FileInfo()
		fileMode := fi.Mode()
		for _, p := range permissions {
			if fileMode&p != 0 {
				return true, nil
			}
		}
		return false, nil
	}
}

func openStream(path string) (*os.File, error) {
	if path == "-" {
		if terminal.IsTerminal(0) {
			flag.Usage()
			os.Exit(64)
		} else {
			return os.Stdin, nil
		}
	}
	return os.Open(path)
}
