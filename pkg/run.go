package pkg

import (
	"bufio"
	l "log"
	"os"
	"strings"

	"github.com/goodwithtech/dockle/pkg/utils"

	"github.com/goodwithtech/dockle/pkg/writer"

	"github.com/genuinetools/reg/registry"
	"github.com/goodwithtech/dockle/pkg/scanner"
	"github.com/knqyf263/fanal/cache"
	"golang.org/x/xerrors"

	"github.com/goodwithtech/dockle/pkg/log"
	"github.com/goodwithtech/dockle/pkg/types"
	"github.com/urfave/cli"
)

var (
	ignoreCheckpointMap map[string]struct{}
)

const (
	dockleIgnore = ".dockleignore"
)

func Run(c *cli.Context) (err error) {
	debug := c.Bool("debug")
	if err = log.InitLogger(debug); err != nil {
		l.Fatal(err)
	}

	cliVersion := "v" + c.App.Version
	latestVersion, err := utils.FetchLatestVersion()

	// check latest version
	if err == nil && cliVersion != latestVersion && c.App.Version != "dev" {
		log.Logger.Warnf("A new version %s is now available! You have %s.", latestVersion, cliVersion)
	}

	clearCache := c.Bool("clear-cache")
	if clearCache {
		log.Logger.Info("Removing image caches...")
		if err = cache.Clear(); err != nil {
			return xerrors.New("failed to remove image layer cache")
		}
	}
	args := c.Args()
	filePath := c.String("input")
	if filePath == "" && len(args) == 0 {
		log.Logger.Info(`"dockle" requires at least 1 argument or --input option.`)
		cli.ShowAppHelpAndExit(c, 1)
		return
	}

	var imageName string
	if filePath == "" {
		imageName = args[0]
	}

	var useLatestTag bool
	// Check whether 'latest' tag is used
	if imageName != "" {
		image, err := registry.ParseImage(imageName)
		if err != nil {
			return xerrors.Errorf("invalid image: %w", err)
		}
		if image.Tag == "latest" && !clearCache {
			useLatestTag = true
			log.Logger.Warn("You should avoid using the :latest tag as it is cached. You need to specify '--clear-cache' option when :latest image is changed")
		}
	}

	log.Logger.Debug("Start assessments...")
	assessments, err := scanner.ScanImage(imageName, filePath)
	if err != nil {
		return err
	}
	if useLatestTag {
		assessments = append(assessments, &types.Assessment{
			Type:     types.AvoidLatestTag,
			Filename: "image tag",
			Desc:     "Avoid 'latest' tag",
		})
	}

	log.Logger.Debug("End assessments...")

	exitCode := c.Int("exit-code")

	// Store ignore checkpoint code
	getIgnoreCheckpointMap()

	var abendAssessments []*types.Assessment

	targetType := types.MinTypeNumber
	for targetType <= types.MaxTypeNumber {
		filtered := filteredAssessments(targetType, assessments)
		writer.ShowTargetResult(targetType, filtered)

		for _, assessment := range filtered {
			abendAssessments = filterAbendAssessments(abendAssessments, assessment)
		}
		targetType++
	}

	if exitCode != 0 && len(abendAssessments) > 0 {
		os.Exit(exitCode)
	}

	return nil
}

func filteredAssessments(target int, assessments []*types.Assessment) (filtered []*types.Assessment) {
	detail := types.AlertDetails[target]
	for _, assessment := range assessments {
		if assessment.Type == target {
			if _, ok := ignoreCheckpointMap[detail.Code]; ok {
				assessment.Level = types.IgnoreLevel
			}
			filtered = append(filtered, assessment)
		}
	}
	return filtered
}

func filterAbendAssessments(abendAssessments []*types.Assessment, assessment *types.Assessment) []*types.Assessment {
	if assessment.Level == types.SkipLevel {
		return abendAssessments
	}

	detail := types.AlertDetails[assessment.Type]
	if _, ok := ignoreCheckpointMap[detail.Code]; ok {
		return abendAssessments
	}
	return append(abendAssessments, assessment)
}

func getIgnoreCheckpointMap() {
	f, err := os.Open(dockleIgnore)
	if err != nil {
		log.Logger.Debug("There is no .dockleignore file")
		// dockle must work even if there isn't ignore file
		return
	}

	ignoreCheckpointMap = map[string]struct{}{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		log.Logger.Debugf("Add new ignore code: %s", line)
		ignoreCheckpointMap[line] = struct{}{}
	}
}
