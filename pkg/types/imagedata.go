package types

import (
	"fmt"
	"io"
	"os"

	"github.com/anchore/stereoscope/pkg/file"
	"github.com/anchore/stereoscope/pkg/image"

	"github.com/Portshift/dockle/pkg/log"
)

type ImageData struct {
	*image.Image
	FileMap map[string]FileData
}

type FileData struct {
	RealPath file.Path
	FileMode os.FileMode
}

type FilterFunc func(filePath string, fileMode os.FileMode) (bool, error)

func (f *FileData) ReadContent(img *image.Image) ([]byte, error) {
	contentReader, err := img.OpenPathFromSquash(f.RealPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", f.RealPath, err)
	}
	content, err := io.ReadAll(contentReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read content of %s: %w", f.RealPath, err)
	}
	if err := contentReader.Close(); err != nil {
		log.Logger.Errorf("Failed to close content reader fo %s: %w", f.RealPath, err)
	}

	return content, nil
}
