package pkg

import (
	"bufio"
	l "log"
	"os"
	"strings"

	"github.com/goodwithtech/docker-guard/pkg/utils"

	"github.com/goodwithtech/docker-guard/pkg/writer"

	"github.com/genuinetools/reg/registry"
	"github.com/goodwithtech/docker-guard/pkg/scanner"
	"github.com/knqyf263/fanal/cache"
	"golang.org/x/xerrors"

	"github.com/goodwithtech/docker-guard/pkg/log"
	"github.com/goodwithtech/docker-guard/pkg/types"
	"github.com/urfave/cli"
)

const (
	guardIgnore = ".guardignore"
)

func Run(c *cli.Context) (err error) {
	debug := c.Bool("debug")
	if err = log.InitLogger(debug); err != nil {
		l.Fatal(err)
	}

	cliVersion := c.App.Version
	latestVersion, err := utils.FetchLatestVersion()

	// check latest version
	if err == nil && cliVersion != latestVersion {
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
		log.Logger.Info(`"docker-guard" requires at least 1 argument or --input option.`)
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
		assessments = append(assessments, types.Assessment{
			Type:     types.AvoidLatestTag,
			Filename: "image tag",
			Desc:     "Avoid 'latest' tag",
		})
	}

	log.Logger.Debug("End assessments...")

	targetType := types.MinTypeNumber
	for targetType <= types.MaxTypeNumber {
		filtered := filteredAssessments(targetType, assessments)
		writer.ShowTargetResult(targetType, filtered)
		targetType++
	}

	exitCode := c.Int("exit-code")
	if exitCode != 0 {
		os.Exit(handleResult(assessments))
	}

	return nil
}

func filteredAssessments(target int, assessments []types.Assessment) (filtered []types.Assessment) {
	for _, assessment := range assessments {
		if assessment.Type == target {
			filtered = append(filtered, assessment)
		}
	}
	return filtered
}

func handleResult(assessments []types.Assessment) (exitCode int) {
	optMap := getIgnoredOptMap()
	for _, assessment := range assessments {
		// skip if ignore opt
		if assessment.Level == types.SkipLevel {
			continue
		}

		detail := types.AlertDetails[assessment.Type]

		if _, ok := optMap[detail.Code]; ok {
			continue
		}

		writer.ShowWhyABEND(detail.Code, assessment)
		exitCode = 1
	}

	return exitCode
}

func getIgnoredOptMap() map[string]struct{} {
	f, err := os.Open(guardIgnore)
	if err != nil {
		log.Logger.Debug("There is no .guardignore file")
		// docker-guard must work even if there isn't ignore file
		return nil
	}

	ignoredMap := map[string]struct{}{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		log.Logger.Debugf("Add new ignore code: %s", line)
		ignoredMap[line] = struct{}{}
	}
	return ignoredMap
}
