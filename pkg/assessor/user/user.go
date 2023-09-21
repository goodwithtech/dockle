package user

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/Portshift/dockle/pkg/log"
	"github.com/Portshift/dockle/pkg/types"
)

type UserAssessor struct{}

func (a UserAssessor) Assess(imageData *types.ImageData) ([]*types.Assessment, error) {
	log.Logger.Debug("Start scan : /etc/passwd")

	var existFile bool
	assesses := []*types.Assessment{}
	for _, filename := range a.RequiredFiles() {
		file, ok := imageData.FileMap[filename]
		if !ok {
			continue
		}
		existFile = true

		content, err := file.ReadContent(imageData.Image)
		if err != nil {
			return nil, fmt.Errorf("failed to read content: %w", err)
		}

		scanner := bufio.NewScanner(bytes.NewBuffer(content))
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
					&types.Assessment{
						Code:     types.AvoidDuplicateUserGroup,
						Filename: filename,
						Desc:     fmt.Sprintf("duplicate UID %s : username %s", uid, uname),
					})
			}
			uidMap[uid] = struct{}{}
		}
		if scanner.Err() != nil {
			return nil, fmt.Errorf("failed to create scanner: %w", err)
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
	return []string{"/etc/passwd"}
}

func (a UserAssessor) RequiredExtensions() []string {
	return []string{}
}

func (a UserAssessor) RequiredPermissions() []os.FileMode {
	return []os.FileMode{}
}
