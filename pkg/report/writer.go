package report

import (
	"sort"

	"github.com/Portshift/dockle/config"
	"github.com/Portshift/dockle/pkg/types"
)

var AlertLabels = map[int]string{
	types.InfoLevel:   "INFO",
	types.WarnLevel:   "WARN",
	types.FatalLevel:  "FATAL",
	types.PassLevel:   "PASS",
	types.SkipLevel:   "SKIP",
	types.IgnoreLevel: "IGNORE",
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
