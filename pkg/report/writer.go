package report

import (
	"sort"

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
var sarifAlertLabels = map[int]string{
	types.InfoLevel:   "note",
	types.WarnLevel:   "warning",
	types.FatalLevel:  "error",
	types.PassLevel:   "none",
	types.SkipLevel:   "none",
	types.IgnoreLevel: "none",
}

type Writer interface {
	Write(assessments types.AssessmentMap) (bool, error)
}

func getCodeOrder() []types.Assessment {
	ass := types.ByLevel{}
	for code, level := range types.DefaultLevelMap {
		if _, ok := config.Conf.IgnoreMap[code]; ok {
			ass = append(ass, types.Assessment{
				Code:  code,
				Level: types.IgnoreLevel,
			})
			continue
		}
		ass = append(ass, types.Assessment{
			Code:  code,
			Level: level,
		})
	}
	sort.Sort(ass)
	return ass
}
