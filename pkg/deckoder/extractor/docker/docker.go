package docker

import (
	"archive/tar"
	"context"
	"crypto/sha256"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	digest "github.com/opencontainers/go-digest"
	"golang.org/x/xerrors"

	"github.com/goodwithtech/dockle/pkg/deckoder/extractor/image/token/ecr"
	"github.com/goodwithtech/dockle/pkg/deckoder/extractor/image/token/gcr"

	"github.com/knqyf263/nested"

	"github.com/goodwithtech/dockle/pkg/deckoder/extractor/image"
	"github.com/goodwithtech/dockle/pkg/deckoder/types"
)

const (
	opq string = ".wh..wh..opq"
	wh  string = ".wh."
)

type Config struct {
	ContainerConfig containerConfig `json:"container_config"`
	History         []History
}

type containerConfig struct {
	Env []string
}

type History struct {
	Created   time.Time
	CreatedBy string `json:"created_by"`
}

type Extractor struct {
	option types.DockerOption
	image  image.Image
}

func init() {
	image.RegisterRegistry(&gcr.GCR{})
	image.RegisterRegistry(&ecr.ECR{})
}

func NewDockerExtractor(ctx context.Context, imageName string, option types.DockerOption) (Extractor, func(), error) {
	ref := image.Reference{Name: imageName, IsFile: false}
	transports := []string{"docker-daemon:", "docker://"}
	return newDockerExtractor(ctx, ref, transports, option)
}

func NewDockerArchiveExtractor(ctx context.Context, fileName string, option types.DockerOption) (Extractor, func(), error) {
	ref := image.Reference{Name: fileName, IsFile: true}
	transports := []string{"docker-archive:"}
	return newDockerExtractor(ctx, ref, transports, option)
}

func newDockerExtractor(ctx context.Context, imgRef image.Reference, transports []string,
	option types.DockerOption) (Extractor, func(), error) {
	ctx, cancel := context.WithTimeout(ctx, option.Timeout)
	defer cancel()

	img, err := image.NewImage(ctx, imgRef, transports, option)
	if err != nil {
		return Extractor{}, nil, xerrors.Errorf("unable to initialize a image struct: %w", err)
	}

	cleanup := func() {
		_ = img.Close()
	}

	return Extractor{
		option: option,
		image:  img,
	}, cleanup, nil
}

func ApplyLayers(layers []types.LayerInfo) types.FileMap {
	sep := "/"
	nestedMap := nested.Nested{}

	for _, layer := range layers {
		for _, opqDir := range layer.OpaqueDirs {
			_ = nestedMap.DeleteByString(opqDir, sep)
		}
		for _, whFile := range layer.WhiteoutFiles {
			_ = nestedMap.DeleteByString(whFile, sep)
		}

		for filePath, content := range layer.TargetFiles {
			// fileName := filepath.Base(filePath)
			// fileDir := filepath.Dir(filePath)
			nestedMap.SetByString(filePath, sep, content)
		}

	}

	fileMap := types.FileMap{}
	_ = nestedMap.Walk(func(keys []string, value interface{}) error {
		content, ok := value.(types.FileData)
		if !ok {
			return nil
		}
		path := strings.Join(keys, sep)
		fileMap[path] = content
		return nil
	})

	return fileMap
}

func (d Extractor) ImageName() string {
	return d.image.Name()
}

func (d Extractor) ImageID() digest.Digest {
	return d.image.ConfigInfo().Digest
}

func (d Extractor) ConfigBlob(ctx context.Context) ([]byte, error) {
	return d.image.ConfigBlob(ctx)
}

func (d Extractor) LayerIDs() []string {
	return d.image.LayerIDs()
}

func (d Extractor) ExtractLayerFiles(ctx context.Context, dig digest.Digest, filterFunc types.FilterFunc) (types.FileMap, []string, []string, error) {
	img, err := d.image.GetLayer(ctx, dig)
	if err != nil {
		return nil, nil, nil, xerrors.Errorf("failed to get a blob: %w", err)
	}
	defer img.Close()

	// calculate decompressed layer ID
	sha256hash := sha256.New()
	r := io.TeeReader(img, sha256hash)

	files, opqDirs, whFiles, err := d.extractFiles(r, filterFunc)
	if err != nil {
		return nil, nil, nil, xerrors.Errorf("failed to extract files: %w", err)
	}

	return files, opqDirs, whFiles, nil
}

// trace another layers if once checked file
var tracingFilepath = map[string]struct{}{}

func (d Extractor) extractFiles(layer io.Reader, filterFunc types.FilterFunc) (types.FileMap, []string, []string, error) {
	data := make(map[string]types.FileData)
	var opqDirs, whFiles []string

	tr := tar.NewReader(layer)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return data, nil, nil, xerrors.Errorf("failed to extract the archive: %w", err)
		}

		filePath := hdr.Name
		filePath = filepath.Clean(filePath)
		fi := hdr.FileInfo()
		fileMode := fi.Mode()

		fileDir, fileName := filepath.Split(filePath)

		// e.g. etc/.wh..wh..opq
		if opq == fileName {
			opqDirs = append(opqDirs, fileDir)
			continue
		}

		// Determine if we should extract the element
		extract := false
		if _, ok := tracingFilepath[filePath]; ok {
			extract = true
		}

		// etc/.wh.hostname
		if strings.HasPrefix(fileName, wh) {
			name := strings.TrimPrefix(fileName, wh)
			fpath := filepath.Join(fileDir, name)
			whFiles = append(whFiles, fpath)
			continue
		}

		if !extract {
			// Determine if we should extract the element
			extract, err = filterFunc(hdr)
			if err != nil {
				return data, nil, nil, xerrors.Errorf("failed to filtering file: %w", err)
			}
			if !extract {
				continue
			}
			tracingFilepath[filePath] = struct{}{}
		}

		if !extract {
			continue
		}

		if hdr.Typeflag == tar.TypeSymlink || hdr.Typeflag == tar.TypeLink || hdr.Typeflag == tar.TypeReg {
			d, err := ioutil.ReadAll(tr)
			if err != nil {
				return nil, nil, nil, xerrors.Errorf("failed to read file: %w", err)
			}
			data[filePath] = types.FileData{
				Body:     d,
				FileMode: fileMode,
			}
		}
	}

	return data, opqDirs, whFiles, nil
}
