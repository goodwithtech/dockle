package scanner

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/anchore/stereoscope"
	"github.com/anchore/stereoscope/pkg/image"

	"github.com/Portshift/dockle/pkg/assessor"
	"github.com/Portshift/dockle/pkg/types"
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

func ScanImage(ctx context.Context, imageName, filePath string, reigstryOptions image.RegistryOptions) (assessments []*types.Assessment, err error) {
	//var ext extractor.Extractor
	//var cleanup func()
	//if imageName != "" {
	//	ext, cleanup, err = docker.NewDockerExtractor(ctx, imageName, dockerOption)
	//	if err != nil {
	//		return nil, fmt.Errorf("%v. %w", err, types.ErrorCreateDockerExtractor)
	//	}
	//} else if filePath != "" {
	//	ext, cleanup, err = docker.NewDockerArchiveExtractor(ctx, filePath, dockerOption)
	//	if err != nil {
	//		return nil, fmt.Errorf("%v. %w", err, types.ErrorCreateDockerExtractor)
	//	}
	//} else {
	//	return nil, types.ErrSetImageOrFile
	//}
	//defer cleanup()
	//ac := analyzer.New(ext)
	//var files deckodertypes.FileMap
	//filterFunc := createPathPermissionFilterFunc(assessor.LoadRequiredFiles(), assessor.LoadRequiredExtensions(), assessor.LoadRequiredPermissions())
	//if files, err = ac.Analyze(ctx, filterFunc); err != nil {
	//	return nil, fmt.Errorf("%v. %w", err, types.ErrorAnalyze)
	//}

	var userInput string
	if filePath != "" {
		userInput = filePath
	} else if imageName != "" {
		userInput = "registry:" + imageName
	} else {
		return nil, types.ErrSetImageOrFile
	}

	opts := stereoscope.WithRegistryOptions(reigstryOptions)
	image, err := stereoscope.GetImage(ctx, userInput, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get image from source : %w", err)
	}

	files := make(map[string]types.FileData)
	filterFunc := createPathPermissionFilterFunc(assessor.LoadRequiredFiles(), assessor.LoadRequiredExtensions(), assessor.LoadRequiredPermissions())

	refs := image.SquashedTree().AllFiles()
	for i := range refs {
		entry, err := image.FileCatalog.Get(refs[i])
		if err != nil {
			panic(err)
		}
		fileMode := entry.FileInfo.Mode()
		ok, err := filterFunc(entry.Path, fileMode)
		if err != nil {
			panic(err)
		}
		if !ok {
			continue
		}

		contentReader, err := image.OpenPathFromSquash(entry.RealPath)
		if err != nil {
			panic(err)
		}
		content, err := io.ReadAll(contentReader)
		if err != nil {
			panic(err)
		}
		files[entry.Path] = types.FileData{
			Body:     content,
			FileMode: entry.FileInfo.Mode(),
		}
	}

	assessments = assessor.GetAssessments(files)
	return assessments, nil
}

func createPathPermissionFilterFunc(filenames, extensions []string, permissions []os.FileMode) types.FilterFunc {
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

	return func(filePath string, fileMode os.FileMode) (bool, error) {
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

		for _, p := range permissions {
			if fileMode&p != 0 {
				return true, nil
			}
		}
		return false, nil
	}
}
