package manifest

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goodwithtech/docker-guard/pkg/types"
	"github.com/knqyf263/fanal/extractor"
	"github.com/moby/moby/image"
	"golang.org/x/xerrors"
)

type ManifestAssessor struct{}

var sensitiveDirs = map[string]struct{}{"/boot": {}, "/dev": {}, "/etc": {}, "/lib": {}, "/proc": {}, "/sys": {}, "/usr": {}}
var suspitiousEnvKey = []string{"PASSWD", "PASSWORD", "SECRET", "ENV", "ACCESS"}

func (a ManifestAssessor) Assess(fileMap extractor.FileMap) (assesses []types.Assessment, err error) {
	file, ok := fileMap["/config"]
	if !ok {
		return nil, xerrors.New("config json file doesn't exist")
	}

	var d image.Image
	json.Unmarshal(file, &d)
	return checkAssessments(d)
}

func checkAssessments(img image.Image) (assesses []types.Assessment, err error) {
	if img.Config.User == "" || img.Config.User == "root" {
		assesses = append(assesses, types.Assessment{
			Type:     types.AvoidRootDefault,
			Filename: "docker config",
			Desc:     "Avoid default user set root",
		})
	}

	for _, envVar := range img.Config.Env {
		e := strings.Split(envVar, "=")
		envKey := e[0]
		for _, suspitiousKey := range suspitiousEnvKey {
			if strings.Contains(envKey, suspitiousKey) {
				assesses = append(assesses, types.Assessment{
					Type:     types.AvoidEnvKeySecret,
					Filename: "docker config",
					Desc:     fmt.Sprintf("Suspitious keyname found : %s", envKey),
				})
			}
		}
	}

	for _, cmd := range img.History {
		if strings.Contains("update", cmd.CreatedBy) {
			assesses = append(assesses, types.Assessment{
				Type:     types.AvoidUpdate,
				Filename: "docker config",
				Desc:     fmt.Sprintf("Avoid update in container : %s", cmd.CreatedBy),
			})
		}
		if strings.Contains("upgrade", cmd.CreatedBy) {
			assesses = append(assesses, types.Assessment{
				Type:     types.AvoidUpgrade,
				Filename: "docker config",
				Desc:     fmt.Sprintf("Avoid upgrade in container : %s", cmd.CreatedBy),
			})
		}
		if strings.Contains("sudo", cmd.CreatedBy) {
			assesses = append(assesses, types.Assessment{
				Type:     types.AvoidSudo,
				Filename: "docker config",
				Desc:     fmt.Sprintf("Avoid sudo in container : %s", cmd.CreatedBy),
			})
		}
	}

	for volume, _ := range img.Config.Volumes {
		if _, ok := sensitiveDirs[volume]; ok {
			assesses = append(assesses, types.Assessment{
				Type:     types.AvoidMountSensitiveDir,
				Filename: "docker config",
				Desc:     fmt.Sprintf("Avoid mounting danger point : %s", volume),
			})
		}

	}
	return assesses, nil
}

// manifest contains /config
func (a ManifestAssessor) RequiredFiles() []string {
	return []string{}
}
