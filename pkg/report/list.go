package report

import (
	"fmt"
	"io"

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

func (lw ListWriter) Write(assessments AssessmentSlice) (bool, error) {
	abend := AssessmentSlice{}
	abendAssessments := &abend
	targetType := types.MinTypeNumber
	for targetType <= types.MaxTypeNumber {
		filtered := assessments.FilteredByTargetCode(targetType)
		showTargetResult(targetType, filtered)

		for _, assessment := range filtered {
			abendAssessments.AddAbend(assessment)
		}
		targetType++
	}
	return len(*abendAssessments) > 0, nil
}

func showTargetResult(assessmentType int, assessments []*types.Assessment) {
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
	if level != types.IgnoreLevel {
		for _, assessment := range assessments {
			showDescription(assessment)
		}
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
