package writer

import (
	"fmt"

	"github.com/goodwithtech/docker-guard/pkg/types"

	"github.com/fatih/color"
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

var AlertLevelColors = []func(a ...interface{}) string{
	color.New(color.FgMagenta).SprintFunc(),
	color.New(color.FgYellow).SprintFunc(),
	color.New(color.FgRed).SprintFunc(),
	color.New(color.FgGreen).SprintFunc(),
	color.New(color.FgBlue).SprintFunc(),
	color.New(color.FgBlue).SprintFunc(),
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
	cyan := color.New(color.FgCyan).SprintFunc()
	detail := types.AlertDetails[assessmentType]
	fmt.Print(colorizeAlert(level), TAB, "-", SPACE, cyan(detail.Code), COLON, SPACE, detail.Title, NEWLINE)
}

func showDescription(assessment *types.Assessment) {
	fmt.Print(TAB, LISTMARK, SPACE, assessment.Desc, NEWLINE)
}

func colorizeAlert(alertLevel int) string {
	return AlertLevelColors[alertLevel](AlertLabels[alertLevel])
}
