package passwd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	deckodertypes "github.com/goodwithtech/deckoder/types"

	"github.com/goodwithtech/dockle/pkg/log"

	"github.com/goodwithtech/dockle/pkg/types"
)

type PasswdAssessor struct{}

func (a PasswdAssessor) Assess(fileMap deckodertypes.FileMap) ([]*types.Assessment, error) {
	log.Logger.Debug("Start scan : password files")

	var existFile bool
	assesses := []*types.Assessment{}
	for _, filename := range a.RequiredFiles() {
		file, ok := fileMap[filename]
		if !ok {
			continue
		}
		existFile = true
		scanner := bufio.NewScanner(bytes.NewBuffer(file.Body))
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) == 0 || line[0] == '#' {
				continue
			}
			passData := strings.Split(line, ":")
			if len(passData) < 2 {
				log.Logger.Debug("The password format may be invalid.", line)
			} else if passData[1] == "" {
				assesses = append(
					assesses,
					&types.Assessment{
						Code:     types.AvoidEmptyPassword,
						Filename: filename,
						Desc:     fmt.Sprintf("No password user found! username : %s", passData[0]),
					})
			}
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
	return []string{"etc/shadow", "etc/master.passwd"}
}

func (a PasswdAssessor) RequiredExtensions() []string {
	return []string{}
}

func (a PasswdAssessor) RequiredPermissions() []os.FileMode {
	return []os.FileMode{}
}
