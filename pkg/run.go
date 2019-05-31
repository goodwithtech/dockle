package pkg

import (
	"bufio"
	l "log"
	"os"
	"strings"

	"github.com/tomoyamachi/docker-guard/pkg/writer"

	"github.com/genuinetools/reg/registry"
	"github.com/knqyf263/fanal/cache"
	"github.com/tomoyamachi/docker-guard/pkg/scanner"
	"golang.org/x/xerrors"

	"github.com/tomoyamachi/docker-guard/pkg/log"
	"github.com/tomoyamachi/docker-guard/pkg/types"
	"github.com/urfave/cli"
)

const (
	guardIgnore = ".guardignore"
)

func Run(c *cli.Context) (err error) {
	// cliVersion := c.App.Version
	result := types.ScanResult{}
	debug := c.Bool("debug")
	if err = log.InitLogger(debug); err != nil {
		l.Fatal(err)
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

	// Check whether 'latest' tag is used
	if imageName != "" {
		image, err := registry.ParseImage(imageName)
		if err != nil {
			return xerrors.Errorf("invalid image: %w", err)
		}
		if image.Tag == "latest" && !clearCache {
			log.Logger.Warn("You should avoid using the :latest tag as it is cached. You need to specify '--clear-cache' option when :latest image is changed")
		}
	}

	assessments, err := scanner.ScanImage(imageName, filePath)
	if err != nil {
		return err
	}

	targetType := types.MinTypeNumber
	for targetType < types.MaxTypeNumber {
		filtered := filteredAssessments(targetType, assessments)
		writer.ShowTitleLine(targetType, len(filtered) == 0)
		targetType++
	}

	exitCode := c.Int("exit-code")
	if exitCode != 0 {
		os.Exit(handleResult(result))
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

func handleResult(r types.ScanResult) (exitCode int) {
	optMap := getIgnoredOptMap()
	for key, targetErr := range r {
		// skip if ignore opt
		if _, ok := optMap[key]; ok {
			continue
		}

		if targetErr != nil {
			exitCode = 1
		}
	}

	return exitCode
}

func getIgnoredOptMap() map[string]struct{} {
	f, err := os.Open(guardIgnore)
	if err != nil {
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
		ignoredMap[line] = struct{}{}
	}
	return ignoredMap
}
