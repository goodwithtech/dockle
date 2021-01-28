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

type ImageAssessment struct {
	Assessment AssessmentMap `json:"assessment"`
	Image      string        `json:"image"`
	Success    bool          `json:"success"`
	ScanUUID   string        `json:"scanuuid"`
	ScanErr    *ScanError    `json:"scanError"`
}

func CreateAssessmentMap(as AssessmentSlice, ignoreMap map[string]struct{}) AssessmentMap {
	asMap := AssessmentMap{}
	for _, a := range as {
		level := a.Level
		if level == 0 {
			level = DefaultLevelMap[a.Code]
		}
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
