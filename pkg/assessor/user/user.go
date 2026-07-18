package user

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/goodwithtech/dockle/pkg/log"
	"github.com/goodwithtech/dockle/pkg/types"
)

type UserAssessor struct{}

func (a UserAssessor) Assess(fileMap types.FileMap) ([]*types.Assessment, error) {
	log.Logger.Debug("Start scan : /etc/passwd")

	var existFile bool
	assesses := []*types.Assessment{}
	for _, filename := range a.RequiredFiles() {
		file, ok := fileMap[filename]
		if !ok {
			continue
		}
		existFile = true
		scanner := bufio.NewScanner(bytes.NewBuffer(file.Body))
		uidMap := map[string]struct{}{}
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) == 0 || line[0] == '#' {
				continue
			}
			data := strings.Split(line, ":")
			if len(data) < 3 {
				log.Logger.Debug("The passwd format may be invalid.", line)
				continue
			}
			uname := data[0]
			uid := data[2]

			// check duplicate UID
			if _, ok := uidMap[uid]; ok {
				assesses = append(
					assesses,
					&types.Assessment{
						Code:     types.AvoidDuplicateUserGroup,
						Filename: filename,
						Desc:     fmt.Sprintf("duplicate UID %s : username %s", uid, uname),
					})
			}
			uidMap[uid] = struct{}{}
		}
	}
	if !existFile {
		assesses = []*types.Assessment{{
			Code:  types.AvoidDuplicateUserGroup,
			Level: types.SkipLevel,
			Desc:  fmt.Sprintf("failed to detect %s", strings.Join(a.RequiredFiles(), ",")),
		}}
	}

	return assesses, nil
}

func (a UserAssessor) RequiredFiles() []string {
	return []string{"etc/passwd"}
}

func (a UserAssessor) RequiredExtensions() []string {
	return []string{}
}

func (a UserAssessor) RequiredPermissions() []os.FileMode {
	return []os.FileMode{}
}
