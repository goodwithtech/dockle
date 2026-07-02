package analyzer

import (
	"context"

	digest "github.com/opencontainers/go-digest"
	"golang.org/x/xerrors"

	"github.com/goodwithtech/dockle/pkg/extractor"
	"github.com/goodwithtech/dockle/pkg/extractor/docker"
	"github.com/goodwithtech/dockle/pkg/types"
)

var (
	additionalFiles []string

	// ErrUnknownOS occurs when unknown OS is analyzed.
	ErrUnknownOS = xerrors.New("unknown OS")
	// ErrPkgAnalysis occurs when the analysis of packages is failed.
	ErrPkgAnalysis = xerrors.New("failed to analyze packages")
	// ErrNoPkgsDetected occurs when the required files for an OS package manager are not detected
	ErrNoPkgsDetected = xerrors.New("no packages detected")
)

type Config struct {
	Extractor extractor.Extractor
}

func New(ext extractor.Extractor) Config {
	return Config{Extractor: ext}
}

func (ac Config) Analyze(ctx context.Context, filterFunc types.FilterFunc) (types.FileMap, error) {
	// always delete cache
	layerInfos := []types.LayerInfo{}
	layerIDs := ac.Extractor.LayerIDs()
	for _, layerID := range layerIDs {
		dig := digest.Digest(layerID)
		layerInfo, err := ac.analyzeLayer(ctx, dig, filterFunc)
		if err != nil {
			return nil, xerrors.Errorf("failed to analyze layer: %s : %w", dig, err)
		}
		layerInfos = append(layerInfos, layerInfo)
	}

	fileMap := docker.ApplyLayers(layerInfos)
	config, err := ac.analyzeConfig(ctx)
	if err != nil {
		return nil, xerrors.Errorf("unable to analyze config: %w", err)
	}
	fileMap["/config"] = config

	return fileMap, nil
}

func (ac Config) analyzeLayer(ctx context.Context, dig digest.Digest, filterFunc types.FilterFunc) (types.LayerInfo, error) {
	files, opqDirs, whFiles, err := ac.Extractor.ExtractLayerFiles(ctx, dig, filterFunc)
	if err != nil {
		return types.LayerInfo{}, xerrors.Errorf("unable to extract files from layer %s: %w", dig, err)
	}

	layerInfo := types.LayerInfo{
		SchemaVersion: types.LayerJSONSchemaVersion,
		TargetFiles:   files,
		OpaqueDirs:    opqDirs,
		WhiteoutFiles: whFiles,
	}
	return layerInfo, nil
}

func (ac Config) analyzeConfig(ctx context.Context) (types.FileData, error) {
	configBlob, err := ac.Extractor.ConfigBlob(ctx)
	if err != nil {
		return types.FileData{}, xerrors.Errorf("unable to get config blob: %w", err)
	}

	// special file for config
	return types.FileData{
		Body: configBlob,
	}, nil
}
