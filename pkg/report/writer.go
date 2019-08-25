package report

import (
	"github.com/goodwithtech/dockle/config"
	"github.com/goodwithtech/dockle/pkg/types"
)

var AlertLabels = map[int]string{
	types.InfoLevel:   "INFO",
	types.WarnLevel:   "WARN",
	types.FatalLevel:  "FATAL",
	types.PassLevel:   "PASS",
	types.SkipLevel:   "SKIP",
	types.IgnoreLevel: "IGNORE",
}

type AssessmentSlice []*types.Assessment
type Writer interface {
	Write(assessments AssessmentSlice) (bool, error)
}

// FilteredByTargetCode returns only target type assessments from all assessments slice
func (as *AssessmentSlice) FilteredByTargetCode(target int) (filtered AssessmentSlice) {
	detail := types.AlertDetails[target]
	for _, assessment := range *as {
		if assessment.Type == target {
			if _, ok := config.Conf.IgnoreMap[detail.Code]; ok {
				assessment.Level = types.IgnoreLevel
			}
			filtered = append(filtered, assessment)
		}
	}
	return filtered
}

// AddAbend add assessment to AssessmentSlice pointer if abend level
func (as *AssessmentSlice) AddAbend(assessment *types.Assessment) {
	level := assessment.Level
	detail := types.AlertDetails[assessment.Type]
	if level == 0 {
		level = detail.DefaultLevel
	}
	if level < config.Conf.ExitLevel {
		return
	}
	if _, ok := config.Conf.IgnoreMap[detail.Code]; ok {
		return
	}
	*as = append(*as, assessment)
}
