package assessor

import (
	"fmt"

	"github.com/tomoyamachi/lyon/pkg/assessor/contentTrust"

	"github.com/tomoyamachi/lyon/pkg/assessor/group"
	"github.com/tomoyamachi/lyon/pkg/assessor/manifest"
	"github.com/tomoyamachi/lyon/pkg/assessor/passwd"
	"github.com/tomoyamachi/lyon/pkg/assessor/user"

	"github.com/knqyf263/fanal/extractor"
	"github.com/tomoyamachi/lyon/pkg/log"
	"github.com/tomoyamachi/lyon/pkg/types"
)

var assessors []Assessor

type Assessor interface {
	Assess(extractor.FileMap) ([]types.Assessment, error)
	RequiredFiles() []string
}

func init() {
	RegisterAssessor(passwd.PasswdAssessor{})
	RegisterAssessor(user.UserAssessor{})
	RegisterAssessor(group.GroupAssessor{})
	RegisterAssessor(manifest.ManifestAssessor{})
	RegisterAssessor(contentTrust.ContentTrustAssessor{})
}

func GetAssessments(files extractor.FileMap) (assessments []types.Assessment) {
	for _, assessor := range assessors {
		results, err := assessor.Assess(files)
		if err != nil {
			log.Logger.Error(err)
		}
		assessments = append(assessments, results...)
	}
	fmt.Println(assessments)
	return assessments
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
