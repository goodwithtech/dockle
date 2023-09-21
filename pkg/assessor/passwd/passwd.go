package passwd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/Portshift/dockle/pkg/log"
	"github.com/Portshift/dockle/pkg/types"
)

type PasswdAssessor struct{}

func (a PasswdAssessor) Assess(imageData *types.ImageData) ([]*types.Assessment, error) {
	log.Logger.Debug("Start scan : password files")

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
		for scanner.Scan() {
			line := scanner.Text()
			passData := strings.Split(line, ":")
			// password must given
			if passData[1] == "" {
				assesses = append(
					assesses,
					&types.Assessment{
						Code:     types.AvoidEmptyPassword,
						Filename: filename,
						Desc:     fmt.Sprintf("No password user found! username : %s", passData[0]),
					})
			}
		}
		if scanner.Err() != nil {
			return nil, fmt.Errorf("failed to create scanner: %w", err)
		}
	}
	if !existFile {
		assesses = []*types.Assessment{
			{
				Code:  types.AvoidEmptyPassword,
				Level: types.SkipLevel,
				Desc:  fmt.Sprintf("failed to detect %s", strings.Join(a.RequiredFiles(), ",")),
			},
		}
	}
	return assesses, nil
}

func (a PasswdAssessor) RequiredFiles() []string {
	return []string{"/etc/shadow", "/etc/master.passwd"}
}

func (a PasswdAssessor) RequiredExtensions() []string {
	return []string{}
}

func (a PasswdAssessor) RequiredPermissions() []os.FileMode {
	return []os.FileMode{}
}
