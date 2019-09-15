package assessor

import (
	"os"

	"github.com/goodwithtech/dockle/pkg/assessor/cache"
	"github.com/goodwithtech/dockle/pkg/assessor/privilege"

	"github.com/goodwithtech/dockle/pkg/assessor/contentTrust"
	"github.com/goodwithtech/dockle/pkg/assessor/credential"
	"github.com/goodwithtech/dockle/pkg/assessor/hosts"

	"github.com/goodwithtech/dockle/pkg/assessor/group"
	"github.com/goodwithtech/dockle/pkg/assessor/manifest"
	"github.com/goodwithtech/dockle/pkg/assessor/passwd"
	"github.com/goodwithtech/dockle/pkg/assessor/user"

	"github.com/goodwithtech/deckoder/extractor"
	"github.com/goodwithtech/dockle/pkg/log"
	"github.com/goodwithtech/dockle/pkg/types"
)

var assessors []Assessor

type Assessor interface {
	Assess(extractor.FileMap) ([]*types.Assessment, error)
	RequiredFiles() []string
	RequiredPermissions() []os.FileMode
}

func init() {
	RegisterAssessor(passwd.PasswdAssessor{})
	RegisterAssessor(privilege.PrivilegeAssessor{})
	RegisterAssessor(user.UserAssessor{})
	RegisterAssessor(group.GroupAssessor{})
	RegisterAssessor(hosts.HostsAssessor{})
	RegisterAssessor(credential.CredentialAssessor{})
	RegisterAssessor(manifest.ManifestAssessor{})
	RegisterAssessor(contentTrust.ContentTrustAssessor{})
	RegisterAssessor(cache.CacheAssessor{})
}

func GetAssessments(files extractor.FileMap) (assessments []*types.Assessment) {
	for _, assessor := range assessors {
		results, err := assessor.Assess(files)
		if err != nil {
			log.Logger.Error(err)
		}
		assessments = append(assessments, results...)
	}
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

func LoadRequiredPermissions() (permissions []os.FileMode) {
	for _, assessor := range assessors {
		permissions = append(permissions, assessor.RequiredPermissions()...)
	}
	return permissions
}
