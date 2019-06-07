package manifest

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goodwithtech/docker-guard/pkg/log"

	"github.com/goodwithtech/docker-guard/pkg/types"
	"github.com/knqyf263/fanal/extractor"
	"golang.org/x/xerrors"
)

type ManifestAssessor struct{}

var sensitiveDirs = map[string]struct{}{"/boot": {}, "/dev": {}, "/etc": {}, "/lib": {}, "/proc": {}, "/sys": {}, "/usr": {}}
var suspiciousEnvKey = []string{"PASSWD", "PASSWORD", "SECRET", "KEY", "ACCESS"}
var acceptanceEnvKey = map[string]struct{}{"GPG_KEY": {}}

func (a ManifestAssessor) Assess(fileMap extractor.FileMap) (assesses []*types.Assessment, err error) {
	log.Logger.Debug("Scan start : config file")
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

func checkAssessments(img types.Image) (assesses []*types.Assessment, err error) {
	if img.Config.User == "" || img.Config.User == "root" {
		assesses = append(assesses, &types.Assessment{
			Type:     types.AvoidRootDefault,
			Filename: "docker config",
			Desc:     "Last user should not be root",
		})
	}

	for _, envVar := range img.Config.Env {
		e := strings.Split(envVar, "=")
		envKey := e[0]
		for _, suspiciousKey := range suspiciousEnvKey {
			if strings.Contains(envKey, suspiciousKey) {
				if _, ok := acceptanceEnvKey[envKey]; ok{
					continue
				}
				assesses = append(assesses, &types.Assessment{
					Type:     types.AvoidEnvKeySecret,
					Filename: "docker config",
					Desc:     fmt.Sprintf("Suspicious ENV key found : %s", envKey),
				})
			}
		}
	}

	if img.Config.Healthcheck == nil {
		assesses = append(assesses, &types.Assessment{
			Type:     types.AddHealthcheck,
			Filename: "docker config",
			Desc:     "not found HEALTHCHECK statement",
		})
	}

	// TODO: use goroutine
	for index, cmd := range img.History {
		if reducableApkAdd(cmd.CreatedBy) {
			assesses = append(assesses, &types.Assessment{
				Type:     types.UseApkAddNoCache,
				Filename: "docker config",
				Desc:     fmt.Sprintf("Use --no-cache option if use 'apk add': %s", cmd.CreatedBy),
			})
		}

		if reducableAptGetInstall(cmd.CreatedBy) {
			assesses = append(assesses, &types.Assessment{
				Type:     types.MinimizeAptGet,
				Filename: "docker config",
				Desc:     fmt.Sprintf("Use 'apt-get clean && rm -rf /var/lib/apt/lists/*' : %s", cmd.CreatedBy),
			})
		}

		if reducableAptGetUpdate(cmd.CreatedBy) {
			assesses = append(assesses, &types.Assessment{
				Type:     types.UseAptGetUpdateNoCache,
				Filename: "docker config",
				Desc:     fmt.Sprintf("Use 'apt-get update --no-cache' : %s", cmd.CreatedBy),
			})
		}

		if strings.Contains(cmd.CreatedBy, "upgrade") {
			assesses = append(assesses, &types.Assessment{
				Type:     types.AvoidDistUpgrade,
				Filename: "docker config",
				Desc:     fmt.Sprintf("Avoid upgrade in container : %s", cmd.CreatedBy),
			})
		}
		if strings.Contains(cmd.CreatedBy, "sudo") {
			assesses = append(assesses, &types.Assessment{
				Type:     types.AvoidSudo,
				Filename: "docker config",
				Desc:     fmt.Sprintf("Avoid sudo in container : %s", cmd.CreatedBy),
			})
		}

		if index != 0 && strings.Contains(cmd.CreatedBy, "ADD") {
			assesses = append(assesses, &types.Assessment{
				Type:     types.UseCOPY,
				Filename: "docker config",
				Desc:     fmt.Sprintf("Use COPY : %s", cmd.CreatedBy),
			})
		}
	}

	for volume := range img.Config.Volumes {
		if _, ok := sensitiveDirs[volume]; ok {
			assesses = append(assesses, &types.Assessment{
				Type:     types.AvoidSensitiveDirectoryMounting,
				Filename: "docker config",
				Desc:     fmt.Sprintf("Avoid mounting sensitive dirs : %s", volume),
			})
		}
	}
	return assesses, nil
}

func reducableAptGetUpdate(cmd string) bool {
	if strings.Contains(cmd, "apt-get update") {
		if !strings.Contains(cmd, "--no-cache") {
			return true
		}
	}
	return false
}

func reducableAptGetInstall(cmd string) bool {
	if strings.Contains(cmd, "apt-get install") {
		if strings.Contains(cmd, "apt-get clean") && strings.Contains(cmd, "rm -rf /var/lib/apt/lists") {
			return false
		}
		return true
	}
	return false
}

func reducableApkAdd(cmd string) bool {
	if strings.Contains(cmd, "apk add") {
		if !strings.Contains(cmd, "--no-cache") {
			return true
		}
	}
	return false
}

// manifest contains /config
func (a ManifestAssessor) RequiredFiles() []string {
	return []string{}
}
