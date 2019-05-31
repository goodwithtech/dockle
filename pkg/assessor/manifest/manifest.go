package manifest

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/docker/docker/image"
	"github.com/goodwithtech/docker-guard/pkg/types"
	"github.com/knqyf263/fanal/extractor"
	"golang.org/x/xerrors"
)

type ManifestAssessor struct{}

var suspitiousEnvKey = []string{
	"PASSWD", "PASSWORD", "SECRET", "ENV",
}

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
			Filename: "config file",
			Desc:     "Avoid default user to root",
		})
	}

	for _, envVar := range img.Config.Env {
		e := strings.Split(envVar, "=")
		envKey := e[0]
		for _, suspitiousKey := range suspitiousEnvKey {
			if strings.Contains(envKey, suspitiousKey) {
				assesses = append(assesses, types.Assessment{
					Type:     types.AvoidEnvKeySecret,
					Filename: "config file",
					Desc:     fmt.Sprintf("Suspitious key found : %s", envKey),
				})
			}
		}
	}

	return assesses, nil
}

// manifest contains /config
func (a ManifestAssessor) RequiredFiles() []string {
	return []string{}
}
