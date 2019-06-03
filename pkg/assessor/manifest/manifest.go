package manifest

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goodwithtech/docker-guard/pkg/types"
	"github.com/knqyf263/fanal/extractor"
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

	var d types.Image

	err = json.Unmarshal(file, &d)
	if err != nil {
		return nil, xerrors.New("Fail to parse docker config file.")
	}

	return checkAssessments(d)
}

func checkAssessments(img types.Image) (assesses []types.Assessment, err error) {
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
		if strings.Contains(cmd.CreatedBy, "apk add") {
			if !strings.Contains(cmd.CreatedBy, "--no-cache") {
				assesses = append(assesses, types.Assessment{
					Type:     types.UseNoCacheAPK,
					Filename: "docker config",
					Desc:     fmt.Sprintf("Use --no-cache option if use apk add : %s", cmd.CreatedBy),
				})
			}
		}

		if minimizableAptGet(cmd.CreatedBy) {
			assesses = append(assesses, types.Assessment{
				Type:     types.MinimizeAptGet,
				Filename: "docker config",
				Desc:     fmt.Sprintf("Use 'apt-get clean' and 'rm -rf /var/lib/apt/lists/*' : %s", cmd.CreatedBy),
			})
		}

		if strings.Contains(cmd.CreatedBy, "upgrade") {
			assesses = append(assesses, types.Assessment{
				Type:     types.AvoidUpgrade,
				Filename: "docker config",
				Desc:     fmt.Sprintf("Avoid upgrade in container : %s", cmd.CreatedBy),
			})
		}
		if strings.Contains(cmd.CreatedBy, "sudo") {
			assesses = append(assesses, types.Assessment{
				Type:     types.AvoidSudo,
				Filename: "docker config",
				Desc:     fmt.Sprintf("Avoid sudo in container : %s", cmd.CreatedBy),
			})
		}
	}

	for volume := range img.Config.Volumes {
		if _, ok := sensitiveDirs[volume]; ok {
			assesses = append(assesses, types.Assessment{
				Type:     types.AvoidSensitiveDirectoryMounting,
				Filename: "docker config",
				Desc:     fmt.Sprintf("Avoid mounting sensitive dirs : %s", volume),
			})
		}
	}
	return assesses, nil
}

func minimizableAptGet(cmd string) bool {
	if strings.Contains(cmd, "apt-get install") {
		if strings.Contains(cmd, "apt-get clean") && strings.Contains(cmd, "rm -rf /var/lib/apt/lists") {
			return false
		}
		return true
	}
	return false
}

// manifest contains /config
func (a ManifestAssessor) RequiredFiles() []string {
	return []string{}
}
