package passwd

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/tomoyamachi/docker-guard/pkg/types"

	"github.com/knqyf263/fanal/extractor"
)

type PasswdAssessor struct{}

func (a PasswdAssessor) Assess(fileMap extractor.FileMap) ([]types.Assessment, error) {
	assesses := []types.Assessment{}
	for _, filename := range a.RequiredFiles() {
		file, ok := fileMap[filename]
		if !ok {
			continue
		}
		scanner := bufio.NewScanner(bytes.NewBuffer(file))
		for scanner.Scan() {
			line := scanner.Text()
			passData := strings.Split(line, ":")
			if passData[1] == "" {
				assesses = append(
					assesses,
					types.Assessment{
						Type:     types.SetPassword,
						Filename: filename,
						Desc:     fmt.Sprintf("no passwd setting user %s ", passData[0]),
					})
			}
		}
	}
	return assesses, nil
}

func (a PasswdAssessor) RequiredFiles() []string {
	return []string{"etc/shadow", "etc/master.passwd"}
}
