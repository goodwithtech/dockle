package group

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

type GroupAssessor struct{}

func (a GroupAssessor) Assess(fileMap deckodertypes.FileMap) ([]*types.Assessment, error) {
	log.Logger.Debug("Start scan : /etc/group")

	var existFile bool
	assesses := []*types.Assessment{}
	for _, filename := range a.RequiredFiles() {
		file, ok := fileMap[filename]
		if !ok {
			continue
		}
		existFile = true
		scanner := bufio.NewScanner(bytes.NewBuffer(file.Body))
		gidMap := map[string]struct{}{}

		for scanner.Scan() {
			line := scanner.Text()
			data := strings.Split(line, ":")
			gname := data[0]
			gid := data[2]

			if _, ok := gidMap[gid]; ok {
				assesses = append(
					assesses,
					&types.Assessment{
						Code:     types.AvoidDuplicateUserGroup,
						Filename: filename,
						Desc:     fmt.Sprintf("duplicate GroupID %s : username %s", gid, gname),
					})
			}
			gidMap[gid] = struct{}{}
		}
	}
	if !existFile {
		assesses = []*types.Assessment{
			{
				Code:  types.AvoidDuplicateUserGroup,
				Level: types.SkipLevel,
				Desc:  fmt.Sprintf("failed to detect %s", strings.Join(a.RequiredFiles(), ",")),
			},
		}
	}

	return assesses, nil
}

func (a GroupAssessor) RequiredFiles() []string {
	return []string{"etc/group"}
}

func (a GroupAssessor) RequiredExtensions() []string {
	return []string{}
}

func (a GroupAssessor) RequiredPermissions() []os.FileMode {
	return []os.FileMode{}
}
