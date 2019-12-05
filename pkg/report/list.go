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

func (lw ListWriter) Write(assessMap types.AssessmentMap) (abend bool, err error) {
	codeOrderLevel := getCodeOrder()
	for _, ass := range codeOrderLevel {
		if _, ok := assessMap[ass.Code]; !ok {
			continue
		}
		assess := assessMap[ass.Code]
		showTargetResult(assess.Code, assess.Level, assess.Assessments)
		if assess.Level >= config.Conf.ExitLevel {
			abend = true
		}
	}
	return abend, nil
}

func showTargetResult(code string, level int, assessments []*types.Assessment) {
	showTitleLine(code, level)
	if level > types.IgnoreLevel {
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
