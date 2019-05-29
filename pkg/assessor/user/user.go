package user

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/tomoyamachi/lyon/pkg/types"

	"github.com/knqyf263/fanal/extractor"
)

type UserAssessor struct{}

func (a UserAssessor) Assess(fileMap extractor.FileMap) ([]types.Assessment, error) {
	assesses := []types.Assessment{}
	for _, filename := range a.RequiredFiles() {
		file, ok := fileMap[filename]
		if !ok {
			continue
		}
		scanner := bufio.NewScanner(bytes.NewBuffer(file))
		uidMap := map[string]struct{}{}
		for scanner.Scan() {
			line := scanner.Text()
			data := strings.Split(line, ":")
			uname := data[0]
			uid := data[2]

			// check duplicate UID
			if _, ok := uidMap[uid]; ok {
				assesses = append(
					assesses,
					types.Assessment{
						Type:     types.AvoidDuplicateUser,
						Filename: filename,
						Desc:     fmt.Sprintf("duplicate UID %s : username %s", uid, uname),
					})
			}
			uidMap[uid] = struct{}{}
		}
	}
	return assesses, nil
}

func (a UserAssessor) RequiredFiles() []string {
	return []string{"etc/passwd"}
}
