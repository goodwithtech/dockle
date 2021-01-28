package scanner

import (
	"archive/tar"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goodwithtech/deckoder/analyzer"
	"github.com/goodwithtech/deckoder/extractor"
	"github.com/goodwithtech/deckoder/extractor/docker"
	"github.com/goodwithtech/deckoder/utils"

	deckodertypes "github.com/goodwithtech/deckoder/types"

	"github.com/Portshift/dockle/pkg/types"

	"github.com/Portshift/dockle/pkg/assessor"

	"golang.org/x/crypto/ssh/terminal"
)

func ScanImage(ctx context.Context, imageName, filePath string, dockerOption deckodertypes.DockerOption) (assessments []*types.Assessment, err error) {
	var files deckodertypes.FileMap
	filterFunc := createPathPermissionFilterFunc(assessor.LoadRequiredFiles(), assessor.LoadRequiredPermissions())
	var ext extractor.Extractor
	var cleanup func()
	if imageName != "" {
		ext, cleanup, err = docker.NewDockerExtractor(ctx, imageName, dockerOption)
		if err != nil {
			return nil, fmt.Errorf("%v. %w", err, types.ErrorCreateDockerExtractor)
		}
	} else if filePath != "" {
		ext, cleanup, err = docker.NewDockerArchiveExtractor(ctx, filePath, dockerOption)
		if err != nil {
			return nil, fmt.Errorf("%v. %w", err, types.ErrorCreateDockerExtractor)
		}
	} else {
		return nil, types.ErrSetImageOrFile
	}
	defer cleanup()
	ac := analyzer.New(ext)
	if files, err = ac.Analyze(ctx, filterFunc); err != nil {
		return nil, fmt.Errorf("%v. %w", err, types.ErrorAnalyze)
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
