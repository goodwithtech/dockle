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

	"github.com/goodwithtech/dockle/pkg/report"

	"github.com/goodwithtech/dockle/pkg/scanner"

	"github.com/urfave/cli"

	"github.com/goodwithtech/dockle/pkg/log"
	"github.com/goodwithtech/dockle/pkg/types"
)

func RunFromCli(c *cli.Context) (err error) {
	if err = log.InitLogger(c.Bool("debug")); err != nil {
		l.Fatal(err)
	}
	config.CreateFromCli(c)
	_, err = run()

	return err
}

func RunFromConfig(conf *config.Config) (types.AssessmentMap, error) {
	config.Conf = *conf

	return run()
}

func run() (ret types.AssessmentMap, err error) {
	ctx := context.Background()
	if err = log.InitLogger(config.Conf.Debug); err != nil {
		l.Fatal(err)
	}

	// TODO: Check latest version

	// set docker option
	dockerOption := deckodertypes.DockerOption{
		Timeout:  config.Conf.Timeout,
		UserName: config.Conf.Username,
		Password: config.Conf.Password,
		SkipPing: true,
	}

	var useLatestTag bool
	// Check whether 'latest' tag is used
	if config.Conf.ImageName != "" {
		if useLatestTag, err = useLatest(config.Conf.ImageName); err != nil {
			return nil, fmt.Errorf("invalid image: %w", err)
		}
	}
	log.Logger.Debug("Start assessments...")

	assessments, err := scanner.ScanImage(ctx, config.Conf.ImageName, config.Conf.FilePath, dockerOption)
	if err != nil {
		return nil, err
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
	o := config.Conf.Output
	output := os.Stdout
	if o != "" {
		if output, err = os.Create(o); err != nil {
			return nil, fmt.Errorf("failed to create an output file: %w", err)
		}
	}

	var writer report.Writer
	switch format := config.Conf.Format; format {
	case "json":
		writer = &report.JsonWriter{Output: output}
	default:
		writer = &report.ListWriter{Output: output}
	}

	abend, err := writer.Write(assessmentMap)
	if err != nil {
		return nil, fmt.Errorf("failed to write results: %w", err)
	}

	if config.Conf.ExitCode != 0 && abend {
		os.Exit(config.Conf.ExitCode)
	}

	return assessmentMap, nil
}

func useLatest(imageName string) (bool, error) {
	ref, err := alltransports.ParseImageName("docker://" + imageName)
	if err != nil {
		return false, err

	}
	return strings.HasSuffix(ref.DockerReference().String(), ":latest"), nil
}
