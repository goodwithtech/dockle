package report

import (
	"fmt"
	"io"

	"github.com/goodwithtech/dockle/config"

	"github.com/goodwithtech/dockle/pkg/color"
	"github.com/goodwithtech/dockle/pkg/types"
)

const (
	LISTMARK = "*"
	COLON    = ":"
	SPACE    = " "
	TAB      = "	"
	NEWLINE  = "\n"
)

var AlertLevelColors = map[int]color.Color{
	types.InfoLevel:   color.Magenta,
	types.WarnLevel:   color.Yellow,
	types.FatalLevel:  color.Red,
	types.PassLevel:   color.Green,
	types.SkipLevel:   color.Blue,
	types.IgnoreLevel: color.Blue,
}

type ListWriter struct {
	Output io.Writer
}

func (lw ListWriter) Write(assessMap types.AssessmentMap) (bool, error) {
	abend := types.AssessmentSlice{}
	abendAssessments := &abend

	codeOrderLevel := getCodeOrder()
	for _, ass := range codeOrderLevel {
		assesses, ok := assessMap[ass.Code]
		if !ok {
			continue
		}
		showTargetResult(ass.Code, ass.Level, assesses)
		for _, assessment := range assesses {
			abendAssessments.AddAbend(assessment, config.Conf.ExitLevel)
		}
	}
	return len(*abendAssessments) > 0, nil
}

func showTargetResult(code string, level int, assessments []*types.Assessment) {
	showTitleLine(code, level)
	if level != types.IgnoreLevel {
		for _, assessment := range assessments {
			showDescription(assessment)
		}
	}
}

func showTitleLine(code string, level int) {
	cyan := color.Cyan
	fmt.Print(colorizeAlert(level), TAB, "-", SPACE, cyan.Add(code), COLON, SPACE, types.TitleMap[code], NEWLINE)
}

func showDescription(assessment *types.Assessment) {
	fmt.Print(TAB, LISTMARK, SPACE, assessment.Desc, NEWLINE)
}

func colorizeAlert(alertLevel int) string {
	return AlertLevelColors[alertLevel].Add(AlertLabels[alertLevel])
}
