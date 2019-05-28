package scanner

import (
	"context"

	"github.com/tomoyamachi/lyon/pkg/types"

	"github.com/tomoyamachi/lyon/pkg/assessor"

	"github.com/knqyf263/fanal/analyzer"
	"github.com/knqyf263/fanal/extractor"
	"golang.org/x/xerrors"
)

func ScanImage(imageName, filePath string) (map[string][]types.Assessment, error) {
	var err error
	results := map[string][]types.Assessment{}
	ctx := context.Background()

	var target string
	var files extractor.FileMap

	// register init assessors
	assessor.InitAssessors()
	// add required files to fanal's analyzer
	analyzer.AddRequiredFilenames(assessor.LoadRequiredFiles())
	if imageName != "" {
		target = imageName
		if err != nil {
			return nil, xerrors.Errorf("failed to get docker option: %w", err)
		}
		files, err = analyzer.Analyze(ctx, imageName)
		if err != nil {
			return nil, xerrors.Errorf("failed to analyze image: %w", err)
		}
	} else if filePath != "" {
		target = filePath
		rc, err := openStream(filePath)
		if err != nil {
			return nil, xerrors.Errorf("failed to open stream: %w", err)
		}

		files, err = analyzer.AnalyzeFromFile(ctx, rc)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, xerrors.New("image name or image file must be specified")
	}

	results = assessor.GetAssessments(files)

	if len(results) == 0 {
		return nil, xerrors.Errorf("failed scan %s: %w", target, err)
	}
	return results, nil
}
