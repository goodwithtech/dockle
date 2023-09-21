package scanner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/anchore/stereoscope"
	"github.com/anchore/stereoscope/pkg/image"

	"github.com/Portshift/dockle/config"
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

func ScanImage(ctx context.Context, cfg config.Config) ([]*types.Assessment, error) {
	registryOptions := image.RegistryOptions{
		InsecureSkipTLSVerify: config.Conf.Insecure,
		Credentials: []image.RegistryCredentials{
			{
				Username: config.Conf.Username,
				Password: config.Conf.Password,
				Token:    config.Conf.Token,
			},
		},
		InsecureUseHTTP: config.Conf.NonSSL,
	}

	var userInput string
	if cfg.FilePath != "" {
		userInput = cfg.FilePath
	} else if cfg.ImageName != "" {
		userInput = setImageSource(cfg.LocalImage, cfg.ImageName)
	} else {
		return nil, types.ErrSetImageOrFile
	}

	opts := stereoscope.WithRegistryOptions(registryOptions)
	img, err := stereoscope.GetImage(ctx, userInput, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get image from user input=%s: %w", userInput, err)
	}
	defer img.Cleanup()

	filterFunc := createPathPermissionFilterFunc(assessor.LoadRequiredFiles(), assessor.LoadRequiredExtensions(), assessor.LoadRequiredPermissions())
	files, err := createFileMap(img, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to create file map: %w", err)
	}

	imageData := &types.ImageData{
		Image:   img,
		FileMap: files,
	}

	return assessor.GetAssessments(imageData), nil
}

func createFileMap(img *image.Image, filterFunc types.FilterFunc) (map[string]types.FileData, error) {
	files := make(map[string]types.FileData)
	refs := img.SquashedTree().AllFiles()
	for i := range refs {
		entry, err := img.FileCatalog.Get(refs[i])
		if err != nil {
			return nil, fmt.Errorf("failed to get entry from file catalog reference=%+v: %w", refs[i], err)
		}
		fileMode := entry.FileInfo.Mode()
		ok, err := filterFunc(entry.Path, fileMode)
		if err != nil {
			return nil, fmt.Errorf("failed to run filter function on file=%s: %w", entry.RealPath, err)
		}
		if !ok {
			continue
		}

		files[entry.Path] = types.FileData{
			RealPath: entry.RealPath,
			FileMode: fileMode,
		}
	}

	return files, nil
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

func setImageSource(local bool, source string) string {
	if local {
		return "docker:" + source
	}
	return "registry:" + source
}
