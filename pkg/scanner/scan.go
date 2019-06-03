package scanner

import (
	"context"
	"flag"
	"os"

	"github.com/goodwithtech/docker-guard/pkg/types"

	"github.com/goodwithtech/docker-guard/pkg/assessor"

	"github.com/knqyf263/fanal/analyzer"
	"github.com/knqyf263/fanal/extractor"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/xerrors"
)

func ScanImage(imageName, filePath string) (assessments []types.Assessment, err error) {
	ctx := context.Background()
	var target string
	var files extractor.FileMap

	// add required files to fanal's analyzer
	analyzer.AddRequiredFilenames(assessor.LoadRequiredFiles())
	if imageName != "" {
		target = imageName
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

	assessments = assessor.GetAssessments(files)
	if len(assessments) == 0 {
		return nil, xerrors.Errorf("failed scan %s: %w", target, err)
	}
	return assessments, nil
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
