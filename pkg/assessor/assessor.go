package assessor

import (
	"fmt"

	"github.com/tomoyamachi/docker-guard/pkg/assessor/contentTrust"

	"github.com/tomoyamachi/docker-guard/pkg/assessor/group"
	"github.com/tomoyamachi/docker-guard/pkg/assessor/manifest"
	"github.com/tomoyamachi/docker-guard/pkg/assessor/passwd"
	"github.com/tomoyamachi/docker-guard/pkg/assessor/user"

	"github.com/knqyf263/fanal/extractor"
	"github.com/tomoyamachi/docker-guard/pkg/log"
	"github.com/tomoyamachi/docker-guard/pkg/types"
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
