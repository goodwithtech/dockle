package report

import (
	"github.com/goodwithtech/dockle/pkg/types"
)

var AlertLabels = []string{
	"INFO",
	"WARN",
	"FATAL",
	"PASS",
	"SKIP",
	"IGNORE",
}

type Writer interface {
	Write(assessments []*types.Assessment) (bool, error)
}

func filteredAssessments(ignoreCheckpointMap map[string]struct{}, target int, assessments []*types.Assessment) (filtered []*types.Assessment) {
	detail := types.AlertDetails[target]
	for _, assessment := range assessments {
		if assessment.Type == target {
			if _, ok := ignoreCheckpointMap[detail.Code]; ok {
				assessment.Level = types.IgnoreLevel
			}
			filtered = append(filtered, assessment)
		}
	}
	return filtered
}

func filterAbendAssessments(ignoreCheckpointMap map[string]struct{}, abendAssessments []*types.Assessment, assessment *types.Assessment) []*types.Assessment {
	if assessment.Level == types.SkipLevel {
		return abendAssessments
	}

	detail := types.AlertDetails[assessment.Type]
	if _, ok := ignoreCheckpointMap[detail.Code]; ok {
		return abendAssessments
	}
	return append(abendAssessments, assessment)
}
