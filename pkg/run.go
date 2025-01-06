package pkg

import (
	"context"
	"errors"
	"fmt"
	l "log"
	"os"
	"strings"

	"github.com/goodwithtech/dockle/pkg/assessor/credential"
	"github.com/goodwithtech/dockle/pkg/assessor/manifest"

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
	quiet := c.Bool("quiet")
	if err = log.InitLogger(debug, quiet); err != nil {
		l.Fatal(err)
	}

	config.CreateFromCli(c)

	cliVersion := "v" + c.App.Version
	if c.Bool("version-check") {
		latestVersion, err := utils.FetchLatestVersion(ctx)
		// check latest version
		if err != nil {
			log.Logger.Infof("Failed to check latest version. %s", err)
		} else if cliVersion != latestVersion && c.App.Version != "dev" {
			log.Logger.Warnf("A new version %s is now available! You have %s.", latestVersion, cliVersion)
		}
	} else {
		log.Logger.Debug("Skipped update confirmation")
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
		DockerDaemonHost:      getDockerSockPath(c),
		DockerDaemonCertPath:  c.String("cert-path"),
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
	manifest.AddSensitiveWords(c.StringSlice("sensitive-word"))
	manifest.AddAcceptanceKeys(c.StringSlice("accept-key"))
	credential.AddSensitiveFiles(c.StringSlice("sensitive-file"))
	scanner.AddAcceptanceFiles(c.StringSlice("accept-file"))
	credential.AddSensitiveFileExtensions(c.StringSlice("sensitive-file-extension"))
	scanner.AddAcceptanceExtensions(c.StringSlice("accept-file-extension"))
	log.Logger.Debug("Start assessments...")
	assessments, err := scanner.ScanImage(ctx, imageName, filePath, dockerOption)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("Pull it with \"docker pull %s\" or \"dockle --timeout 600s\" to increase the timeout\n%w", imageName, err)
		}
		return err
	}
	if useLatestTag {
		assessments = append(assessments, &types.Assessment{
			Code:     types.AvoidLatestTag,
			Filename: imageName,
			Desc:     "Avoid 'latest' tag",
		})
	}

	log.Logger.Debug("End assessments...")

	assessmentMap := types.CreateAssessmentMap(assessments, config.Conf.IgnoreMap, debug)
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
		writer = &report.JsonWriter{Output: output, ImageName: imageName}
	case "sarif":
		writer = &report.SarifWriter{Output: output}
	default:
		writer = &report.ListWriter{Output: output, NoColor: c.Bool("no-color")}
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

func getDockerSockPath(c *cli.Context) string {
	if c.String("host") != "" {
		return c.String("host")
	}
	xdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if c.Bool("use-xdg") && xdgRuntimeDir != "" {
		return fmt.Sprintf("unix://%s/docker.sock", xdgRuntimeDir)
	}
	return "unix:///var/run/docker.sock"
}

func useLatest(imageName string) (bool, error) {
	ref, err := alltransports.ParseImageName("docker://" + imageName)
	if err != nil {
		return false, err

	}
	return strings.HasSuffix(ref.DockerReference().String(), ":latest"), nil
}
