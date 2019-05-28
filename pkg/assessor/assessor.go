package assessor

import (
	"github.com/knqyf263/fanal/extractor"
	"github.com/tomoyamachi/lyon/pkg/assessor/passwd"
	"github.com/tomoyamachi/lyon/pkg/log"
	"github.com/tomoyamachi/lyon/pkg/types"
)

var assessors []Assessor

type Assessor interface {
	GetType() string
	Assess(extractor.FileMap) ([]types.Assessment, error)
	RequiredFiles() []string
}

func GetAssessments(files extractor.FileMap) map[string][]types.Assessment {
	assessments := map[string][]types.Assessment{}
	for _, assessor := range assessors {
		results, err := assessor.Assess(files)
		if err != nil {
			log.Logger.Error(err)
		}
		assessments[assessor.GetType()] = results
	}
	return assessments
}

func InitAssessors() {
	RegisterAssessor(passwd.PasswdAssessor{})
}

func RegisterAssessor(a Assessor) {
	assessors = append(assessors, a)
}

func LoadRequiredFiles() (filenames []string) {
	for _, assessor := range assessors {
		filenames = append(filenames, assessor.RequiredFiles()...)
	}
	return filenames
}
