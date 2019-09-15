package scanner

import (
	"archive/tar"
	"context"
	"errors"
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

func ScanImage(imageName, filePath string) (assessments []*types.Assessment, err error) {
	ctx := context.Background()
	var files extractor.FileMap

	filterFunc := createPathPermissionFilterFunc(assessor.LoadRequiredFiles(), assessor.LoadRequiredPermissions())
	if imageName != "" {
		dockerOption, err := types.GetDockerOption()
		if err != nil {
			return nil, fmt.Errorf("failed to get docker option: %w", err)
		}
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
		return nil, errors.New("image name or image file must be specified")
	}

	assessments = assessor.GetAssessments(files)
	return assessments, nil
}

func createPathPermissionFilterFunc(filenames []string, permissions []os.FileMode) deckodertypes.FilterFunc {
	return func(h *tar.Header) (bool, error) {
		filePath := filepath.Clean(h.Name)
		fileName := filepath.Base(filePath)
		fileDirBase := filepath.Base(filepath.Dir(filePath))

		for _, s := range filenames {
			if s[len(s)-1] == '/' {
				if filepath.Clean(s) == fileDirBase {
					return true, nil
				}
			}
		}

		if utils.StringInSlice(filePath, filenames) || utils.StringInSlice(fileName, filenames) {
			return true, nil
		}
		fi := h.FileInfo()
		fileMode := fi.Mode()
		for _, p := range permissions {
			if fileMode&p != 0 {
				return true, nil
				break
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
