package types

type Assessment struct {
	Code     string
	Level    int
	Filename string
	Desc     string
}
type AssessmentSlice []*Assessment
type CodeInfo struct {
	Code        string
	Level       int
	Assessments AssessmentSlice
}
type AssessmentMap map[string]CodeInfo

func CreateAssessmentMap(as AssessmentSlice, ignoreMap map[string]struct{}) AssessmentMap {
	asMap := AssessmentMap{}
	for _, a := range as {
		level := DefaultLevelMap[a.Code]
		if _, ok := ignoreMap[a.Code]; ok {
			level = IgnoreLevel
		}
		if _, ok := asMap[a.Code]; !ok {
			asMap[a.Code] = CodeInfo{
				Code:        a.Code,
				Level:       level,
				Assessments: []*Assessment{a},
			}
		} else {
			asMap[a.Code] = CodeInfo{
				Code:        a.Code,
				Level:       level,
				Assessments: append(asMap[a.Code].Assessments, a),
			}
		}
	}
	return asMap
}

type ByLevel []Assessment

func (a ByLevel) Len() int { return len(a) }
func (a ByLevel) Less(i, j int) bool {
	if a[i].Level == a[j].Level {
		return a[i].Code < a[j].Code
	}
	return a[i].Level > a[j].Level
}
func (a ByLevel) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// AddAbend add assessment to AssessmentSlice pointer if abend level
func (as AssessmentSlice) AddAbend(assessment *Assessment, exitLevel int) {
	level := assessment.Level
	if level == 0 {
		level = DefaultLevelMap[assessment.Code]
	}
	if level < exitLevel {
		return
	}
	as = append(as, assessment)
}
