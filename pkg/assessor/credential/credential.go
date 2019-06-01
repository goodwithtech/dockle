package credential

import (
	"fmt"

	"github.com/goodwithtech/docker-guard/pkg/types"

	"github.com/knqyf263/fanal/extractor"
)

type CredentialAssessor struct{}

func (a CredentialAssessor) Assess(fileMap extractor.FileMap) ([]types.Assessment, error) {
	assesses := []types.Assessment{}
	for _, filename := range a.RequiredFiles() {
		_, ok := fileMap[filename]
		if !ok {
			continue
		}
		assesses = append(
			assesses,
			types.Assessment{
				Type:     types.AvoidCredential,
				Filename: filename,
				Desc:     fmt.Sprintf("Suspicious file found : %s ", filename),
			})

	}
	return assesses, nil
}

func (a CredentialAssessor) RequiredFiles() []string {
	return []string{"credentials.json", "credential.json", "credentials", "credential"}
}
