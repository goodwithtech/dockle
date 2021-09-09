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
	NoColor bool
}

func (lw ListWriter) Write(assessMap types.AssessmentMap) (abend bool, err error) {
	codeOrderLevel := getCodeOrder()
	for _, ass := range codeOrderLevel {
		if _, ok := assessMap[ass.Code]; !ok {
			continue
		}
		assess := assessMap[ass.Code]
		lw.showTargetResult(assess.Code, assess.Level, assess.Assessments)
		if assess.Level >= config.Conf.ExitLevel {
			abend = true
		}
	}
	return abend, nil
}

func (lw ListWriter) showTargetResult(code string, level int, assessments []*types.Assessment) {
	lw.showTitleLine(code, level)
	if level > types.IgnoreLevel {
		for _, assessment := range assessments {
			lw.showDescription(assessment)
		}
	}
}

func (lw ListWriter) showTitleLine(code string, level int) {
	if lw.NoColor {
		fmt.Fprint(lw.Output, AlertLabels[level], TAB, "-", SPACE, code, COLON, SPACE, types.TitleMap[code], NEWLINE)
		return
	}
	cyan := color.Cyan
	fmt.Fprint(lw.Output, colorizeAlert(level), TAB, "-", SPACE, cyan.Add(code), COLON, SPACE, types.TitleMap[code], NEWLINE)
}

func (lw ListWriter) showDescription(assessment *types.Assessment) {
	fmt.Fprint(lw.Output, TAB, LISTMARK, SPACE, assessment.Desc, NEWLINE)
}

func colorizeAlert(alertLevel int) string {
	return AlertLevelColors[alertLevel].Add(AlertLabels[alertLevel])
}
