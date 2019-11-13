package types

type Assessment struct {
	Code     string
	Level    int
	Filename string
	Desc     string
}

type AssessmentMap map[string][]*Assessment

func CreateAssessmentMap(as AssessmentSlice) AssessmentMap {
	asMap := AssessmentMap{}
	for _, a := range as {
		if _, ok := asMap[a.Code]; !ok {
			asMap[a.Code] = []*Assessment{
				a,
			}
		} else {
			asMap[a.Code] = append(asMap[a.Code], a)
		}
	}
	return asMap
}

type AssessmentSlice []*Assessment
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
