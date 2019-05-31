package writer

import (
	"fmt"

	"github.com/tomoyamachi/docker-guard/pkg/types"

	"github.com/fatih/color"
)

const (
	SPACE   = " "
	TAB     = "	"
	NEWLINE = "\n"
)

var AlertLabels = []string{
	"INFO",
	"WARN",
	"FATAL",
	"PASS",
}

var AlertLevelColors = []func(a ...interface{}) string{
	color.New(color.FgBlue).SprintFunc(),
	color.New(color.FgYellow).SprintFunc(),
	color.New(color.FgRed).SprintFunc(),
	color.New(color.FgCyan).SprintFunc(),
}

func ShowTitleLine(assessmentType int, passed bool) {
	detail := types.AlertDetails[assessmentType]
	level := detail.DefaultLevel
	if passed {
		level = types.Pass
	}
	fmt.Print(ColorizeAlert(level), TAB, detail.Title, NEWLINE)
}

func ColorizeAlert(alertLevel int) string {
	return AlertLevelColors[alertLevel](AlertLabels[alertLevel])
	// return color.New(color.FgCyan).SprintFunc()(alertLevel)
}
