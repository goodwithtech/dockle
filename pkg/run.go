package pkg

import (
	"context"
	"fmt"
	l "log"
	"os"
	"strings"

	"github.com/containers/image/v5/transports/alltransports"
	deckodertypes "github.com/goodwithtech/deckoder/types"

	"github.com/goodwithtech/dockle/config"
	"github.com/goodwithtech/dockle/pkg/utils"

	"github.com/goodwithtech/dockle/pkg/report"

	"github.com/goodwithtech/dockle/pkg/scanner"

	"github.com/urfave/cli"

	"github.com/goodwithtech/dockle/pkg/log"
	"github.com/goodwithtech/dockle/pkg/types"
)

func Run(c *cli.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Duration("timeout"))
	defer cancel()
	debug := c.Bool("debug")
	if err = log.InitLogger(debug); err != nil {
		l.Fatal(err)
	}

	config.CreateFromCli(c)

	cliVersion := "v" + c.App.Version
	latestVersion, err := utils.FetchLatestVersion(ctx)

	// check latest version
	if err != nil {
		log.Logger.Infof("Failed to check latest version. %s", err)
	} else if cliVersion != latestVersion && c.App.Version != "dev" {
		log.Logger.Warnf("A new version %s is now available! You have %s.", latestVersion, cliVersion)
	}

	args := c.Args()
	filePath := c.String("input")
	if filePath == "" && len(args) == 0 {
		log.Logger.Info(`"dockle" requires at least 1 argument or --input option.`)
		cli.ShowAppHelpAndExit(c, 1)
		return
	}
	// set docker option
	dockerOption := deckodertypes.DockerOption{
		Timeout:               c.Duration("timeout"),
		UserName:              c.String("username"),
		Password:              c.String("password"),
		InsecureSkipTLSVerify: c.Bool("insecure"),
		SkipPing:              true,
	}
	var imageName string
	if filePath == "" {
		imageName = args[0]
	}

	var useLatestTag bool
	// Check whether 'latest' tag is used
	if imageName != "" {
		if useLatestTag, err = useLatest(imageName); err != nil {
			return fmt.Errorf("invalid image: %w", err)
		}
	}
	log.Logger.Debug("Start assessments...")

	assessments, err := scanner.ScanImage(ctx, imageName, filePath, dockerOption)
	if err != nil {
		return err
	}
	if useLatestTag {
		assessments = append(assessments, &types.Assessment{
			Code:     types.AvoidLatestTag,
			Filename: "image tag",
			Desc:     "Avoid 'latest' tag",
		})
	}

	log.Logger.Debug("End assessments...")

	assessmentMap := types.CreateAssessmentMap(assessments, config.Conf.IgnoreMap)
	// Store ignore checkpoint code
	o := c.String("output")
	output := os.Stdout
	if o != "" {
		if output, err = os.Create(o); err != nil {
			return fmt.Errorf("failed to create an output file: %w", err)
		}
	}

	var writer report.Writer
	switch format := c.String("format"); format {
	case "json":
		writer = &report.JsonWriter{Output: output}
	default:
		writer = &report.ListWriter{Output: output}
	}

	abend, err := writer.Write(assessmentMap)
	if err != nil {
		return fmt.Errorf("failed to write results: %w", err)
	}

	if config.Conf.ExitCode != 0 && abend {
		os.Exit(config.Conf.ExitCode)
	}

	return nil
}

func useLatest(imageName string) (bool, error) {
	ref, err := alltransports.ParseImageName("docker://" + imageName)
	if err != nil {
		return false, err

	}
	return strings.HasSuffix(ref.DockerReference().String(), ":latest"), nil
}
