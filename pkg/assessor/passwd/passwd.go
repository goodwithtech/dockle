package passwd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/goodwithtech/dockle/pkg/log"

	"github.com/goodwithtech/dockle/pkg/types"

	"github.com/goodwithtech/deckoder/extractor"
)

type PasswdAssessor struct{}

func (a PasswdAssessor) Assess(fileMap extractor.FileMap) ([]*types.Assessment, error) {
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
			passData := strings.Split(line, ":")
			// password must given
			if passData[1] == "" {
				assesses = append(
					assesses,
					&types.Assessment{
						Type:     types.AvoidEmptyPassword,
						Filename: filename,
						Desc:     fmt.Sprintf("No password user found! username : %s", passData[0]),
					})
			}
		}
	}
	if !existFile {
		assesses = []*types.Assessment{{Type: types.AvoidEmptyPassword, Level: types.SkipLevel}}
	}
	return assesses, nil
}

func (a PasswdAssessor) RequiredFiles() []string {
	return []string{"etc/shadow", "etc/master.passwd"}
}

func (a PasswdAssessor) RequiredPermissions() []os.FileMode {
	return []os.FileMode{}
}
