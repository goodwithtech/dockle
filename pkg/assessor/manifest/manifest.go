package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/goodwithtech/dockle/pkg/log"

	"github.com/goodwithtech/dockle/pkg/types"
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

	err = json.Unmarshal(file.Body, &d)
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
				if _, ok := acceptanceEnvKey[envKey]; ok {
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

	assessesCh := make(chan []*types.Assessment)
	for index, cmd := range img.History {
		go func(index int, cmd types.History) {
			assessesCh <- assessHistory(index, cmd)
		}(index, cmd)
	}

	timeout := time.After(10 * time.Second)
	for i := 0; i < len(img.History); i++ {
		select {
		case results := <-assessesCh:
			assesses = append(assesses, results...)
		case <-timeout:
			return nil, xerrors.New("timeout: manifest assessor")
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

func splitByCommands(line string) map[int][]string {
	commands := strings.Split(line, "&&")

	tokens := map[int][]string{}
	for index, command := range commands {
		splitted := strings.Split(command, " ")
		cmds := []string{}
		for _, cmd := range splitted {
			trimmed := strings.TrimSpace(cmd)
			if trimmed != "" {
				cmds = append(cmds, cmd)
			}

		}
		tokens[index] = cmds
	}
	return tokens
}

func assessHistory(index int, cmd types.History) []*types.Assessment {
	var assesses []*types.Assessment
	cmdSlices := splitByCommands(cmd.CreatedBy)
	if reducableApkAdd(cmdSlices) {
		assesses = append(assesses, &types.Assessment{
			Type:     types.UseApkAddNoCache,
			Filename: "docker config",
			Desc:     fmt.Sprintf("Use --no-cache option if use 'apk add': %s", cmd.CreatedBy),
		})
	}

	if reducableAptGetInstall(cmdSlices) {
		assesses = append(assesses, &types.Assessment{
			Type:     types.MinimizeAptGet,
			Filename: "docker config",
			Desc:     fmt.Sprintf("Use 'apt-get clean && rm -rf /var/lib/apt/lists/*' : %s", cmd.CreatedBy),
		})
	}

	if reducableAptGetUpdate(cmdSlices) {
		assesses = append(assesses, &types.Assessment{
			Type:     types.UseAptGetUpdateNoCache,
			Filename: "docker config",
			Desc:     fmt.Sprintf("Always combine 'apt-get update' with 'apt-get install' : %s", cmd.CreatedBy),
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
	return assesses
}

func reducableAptGetUpdate(cmdSlices map[int][]string) bool {
	var useAptUpdate bool
	var useAptInstall bool
	for _, cmdSlice := range cmdSlices {
		if !useAptUpdate && ContainAll(cmdSlice, []string{"apt-get", "update"}) {
			useAptUpdate = true
		}

		if !useAptInstall && ContainAll(cmdSlice, []string{"apt-get", "install"}) {
			useAptInstall = true
		}
		if useAptUpdate && useAptInstall {
			return false
		}
	}
	if useAptUpdate && !useAptInstall {
		return true
	}
	return false
}

func reducableAptGetInstall(cmdSlices map[int][]string) bool {
	var useAptInstall bool
	var useRmCache bool
	for _, cmdSlice := range cmdSlices {
		if useAptInstall == false && ContainAll(cmdSlice, []string{"apt-get", "install"}) {
			useAptInstall = true
		}
		if !useRmCache && ContainAll(cmdSlice, []string{"rm", "-rf", "/var/lib/apt/lists"}) {
			useRmCache = true
		}
		if !useRmCache && ContainAll(cmdSlice, []string{"rm", "-fr", "/var/lib/apt/lists"}) {
			useRmCache = true
		}
		if !useRmCache && ContainAll(cmdSlice, []string{"rm", "-fR", "/var/lib/apt/lists"}) {
			useRmCache = true
		}
		if !useRmCache && ContainAll(cmdSlice, []string{"rm", "-rf", "/var/lib/apt/lists/*"}) {
			useRmCache = true
		}
		if !useRmCache && ContainAll(cmdSlice, []string{"rm", "-fr", "/var/lib/apt/lists/*"}) {
			useRmCache = true
		}
		if !useRmCache && ContainAll(cmdSlice, []string{"rm", "-fR", "/var/lib/apt/lists/*"}) {
			useRmCache = true
		}

		if useAptInstall && useRmCache {
			return false
		}
	}
	if useAptInstall && !useRmCache {
		return true
	}
	return false
}

func reducableApkAdd(cmdSlices map[int][]string) bool {
	for _, cmdSlice := range cmdSlices {
		if ContainAll(cmdSlice, []string{"apk", "add"}) {
			if !ContainAll(cmdSlice, []string{"--no-cache"}) {
				return true
			}
		}
	}
	return false
}

// manifest contains /config
func (a ManifestAssessor) RequiredFiles() []string {
	return []string{}
}

func (a ManifestAssessor) RequiredPermissions() []os.FileMode {
	return []os.FileMode{}
}

func ContainAll(heystack []string, needles []string) bool {
	needleMap := map[string]struct{}{}
	for _, n := range needles {
		needleMap[n] = struct{}{}
	}

	for _, v := range heystack {
		if _, ok := needleMap[v]; ok {
			delete(needleMap, v)
			if len(needleMap) == 0 {
				return true
			}
		}
	}
	return false
}
