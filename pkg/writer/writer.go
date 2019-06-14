package writer

import (
	"fmt"

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

var AlertLabels = []string{
	"INFO",
	"WARN",
	"FATAL",
	"PASS",
	"SKIP",
	"IGNORE",
}

var AlertLevelColors = []color.Color{
	color.Magenta,
	color.Yellow,
	color.Red,
	color.Green,
	color.Blue,
	color.Blue,
}

func ShowTargetResult(assessmentType int, assessments []*types.Assessment) {
	if len(assessments) == 0 {
		showTitleLine(assessmentType, types.PassLevel)
		return
	}

	if assessments[0].Level == types.SkipLevel {
		showTitleLine(assessmentType, types.SkipLevel)
		return
	}
	detail := types.AlertDetails[assessmentType]
	level := detail.DefaultLevel
	if assessments[0].Level == types.IgnoreLevel {
		level = types.IgnoreLevel
	}
	showTitleLine(assessmentType, level)
	for _, assessment := range assessments {
		showDescription(assessment)
	}
}

func showTitleLine(assessmentType int, level int) {
	cyan := color.Cyan
	detail := types.AlertDetails[assessmentType]
	fmt.Print(colorizeAlert(level), TAB, "-", SPACE, cyan.Add(detail.Code), COLON, SPACE, detail.Title, NEWLINE)
}

func showDescription(assessment *types.Assessment) {
	fmt.Print(TAB, LISTMARK, SPACE, assessment.Desc, NEWLINE)
}

func colorizeAlert(alertLevel int) string {
	return AlertLevelColors[alertLevel].Add(AlertLabels[alertLevel])
}
