package pkg

import (
	"bufio"
	l "log"
	"os"
	"strings"

	"github.com/goodwithtech/dockle/pkg/utils"

	"github.com/goodwithtech/dockle/pkg/report"

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

	// delete image cache each time
	if err = cache.Clear(); err != nil {
		return xerrors.New("failed to remove image layer cache")
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
		if image.Tag == "latest" {
			useLatestTag = true
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
	ignoreRules := c.StringSlice("ignore")
	getIgnoreCheckpointMap(ignoreRules)
	o := c.String("output")
	output := os.Stdout
	if o != "" {
		if output, err = os.Create(o); err != nil {
			return xerrors.Errorf("failed to create an output file: %w", err)
		}
	}

	var writer report.Writer
	switch format := c.String("format"); format {
	case "json":
		writer = &report.JsonWriter{Output: output, IgnoreMap: ignoreCheckpointMap}
	default:
		writer = &report.ListWriter{Output: output, IgnoreMap: ignoreCheckpointMap}
	}

	abend, err := writer.Write(assessments)
	if err != nil {
		return xerrors.Errorf("failed to write results: %w", err)
	}
	if exitCode != 0 && abend {
		os.Exit(exitCode)
	}

	return nil
}

func getIgnoreCheckpointMap(ignoreRules []string) {
	ignoreCheckpointMap = map[string]struct{}{}
	for _, rule := range ignoreRules {
		ignoreCheckpointMap[rule] = struct{}{}
	}

	f, err := os.Open(dockleIgnore)
	if err != nil {
		log.Logger.Debug("There is no .dockleignore file")
		// dockle must work even if there isn't ignore file
		return
	}

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
