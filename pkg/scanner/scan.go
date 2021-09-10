package scanner

import (
	"archive/tar"
	"context"
	"os"
	"path/filepath"

	"github.com/goodwithtech/deckoder/analyzer"
	"github.com/goodwithtech/deckoder/extractor"
	"github.com/goodwithtech/deckoder/extractor/docker"
	deckodertypes "github.com/goodwithtech/deckoder/types"

	"github.com/goodwithtech/dockle/pkg/types"

	"github.com/goodwithtech/dockle/pkg/assessor"
)

var (
	acceptanceFiles      = map[string]struct{}{}
	acceptanceExtensions = map[string]struct{}{}
)

func AddAcceptanceFiles(keys []string) {
	for _, key := range keys {
		acceptanceFiles[key] = struct{}{}
	}
}

func AddAcceptanceExtensions(keys []string) {
	for _, key := range keys {
		// file extension must start with .
		acceptanceExtensions["."+key] = struct{}{}
	}
}

func ScanImage(ctx context.Context, imageName, filePath string, dockerOption deckodertypes.DockerOption) (assessments []*types.Assessment, err error) {
	var ext extractor.Extractor
	var cleanup func()
	if imageName != "" {
		ext, cleanup, err = docker.NewDockerExtractor(ctx, imageName, dockerOption)
		if err != nil {
			return nil, err
		}
	} else if filePath != "" {
		ext, cleanup, err = docker.NewDockerArchiveExtractor(ctx, filePath, dockerOption)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, types.ErrSetImageOrFile
	}
	defer cleanup()
	ac := analyzer.New(ext)
	var files deckodertypes.FileMap
	filterFunc := createPathPermissionFilterFunc(assessor.LoadRequiredFiles(), assessor.LoadRequiredExtensions(), assessor.LoadRequiredPermissions())
	if files, err = ac.Analyze(ctx, filterFunc); err != nil {
		return nil, err
	}

	assessments = assessor.GetAssessments(files)
	return assessments, nil
}

func createPathPermissionFilterFunc(filenames, extensions []string, permissions []os.FileMode) deckodertypes.FilterFunc {
	requiredDirNames := map[string]struct{}{}
	requiredFileNames := map[string]struct{}{}
	requiredExts := map[string]struct{}{}
	for _, filename := range filenames {
		if filename[len(filename)-1] == '/' {
			// if filename end "/", it is directory and requiredDirNames removes last "/"
			requiredDirNames[filepath.Clean(filename)] = struct{}{}
		} else {
			requiredFileNames[filename] = struct{}{}
		}
	}
	for _, extension := range extensions {
		requiredExts[extension] = struct{}{}
	}

	return func(h *tar.Header) (bool, error) {
		filePath := filepath.Clean(h.Name)
		fileName := filepath.Base(filePath)
		// Skip check if acceptance files
		if _, ok := acceptanceExtensions[filepath.Ext(fileName)]; ok {
			return false, nil
		}
		if _, ok := acceptanceFiles[filePath]; ok {
			return false, nil
		}
		if _, ok := acceptanceFiles[fileName]; ok {
			return false, nil
		}

		// Check with file names
		if _, ok := requiredFileNames[filePath]; ok {
			return true, nil
		}
		if _, ok := requiredFileNames[fileName]; ok {
			return true, nil
		}

		// Check with file extensions
		if _, ok := requiredExts[filepath.Ext(fileName)]; ok {
			return true, nil
		}

		// Check with file directory name
		fileDir := filepath.Dir(filePath)
		if _, ok := requiredDirNames[fileDir]; ok {
			return true, nil
		}
		fileDirBase := filepath.Base(fileDir)
		if _, ok := requiredDirNames[fileDirBase]; ok {
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
