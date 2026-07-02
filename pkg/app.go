package pkg

import (
	"time"

	"github.com/urfave/cli/v3"
)

var (
	version = "dev"
)

// rootHelpTemplate keeps the compact single-command help layout that
// Dockle has always used (no COMMANDS section, two-space indentation).
const rootHelpTemplate = `NAME:
  {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}
USAGE:
  {{if .UsageText}}{{.UsageText}}{{else}}{{.FullName}} {{if .VisibleFlags}}[options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}
VERSION:
  {{.Version}}{{end}}{{end}}{{if .Description}}
DESCRIPTION:
  {{.Description}}{{end}}{{if len .Authors}}
AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
  {{range $index, $author := .Authors}}{{if $index}}
  {{end}}{{$author}}{{end}}{{end}}{{if .VisibleFlags}}
OPTIONS:
  {{range $index, $option := .VisibleFlags}}{{if $index}}
  {{end}}{{$option}}{{end}}{{end}}
`

/*
NewApp Factory for Dockle CLI creation.
An Enabler for programmatic usage of Dockle
*/
func NewApp() *cli.Command {
	cmd := &cli.Command{
		Name:      "dockle",
		Version:   version,
		ArgsUsage: "image_name",
		Usage:     "Container Image Linter for Security, Helping build the Best-Practice Docker Image, Easy to start",

		CustomRootCommandHelpTemplate: rootHelpTemplate,

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "input",
				Usage: "input file path instead of image name",
			},
			&cli.StringSliceFlag{
				Name:    "ignore",
				Aliases: []string{"i"},
				Sources: cli.EnvVars("DOCKLE_IGNORES"),
				Usage:   "checkpoints to ignore. You can use .dockleignore too.",
			},
			&cli.StringSliceFlag{
				Name:    "accept-key",
				Aliases: []string{"ak"},
				Sources: cli.EnvVars("DOCKLE_ACCEPT_KEYS"),
				Usage:   "For CIS-DI-0010. You can add acceptable keywords. e.g) -ak GPG_KEY -ak KEYCLOAK",
			},
			&cli.StringSliceFlag{
				Name:    "sensitive-word",
				Aliases: []string{"sw"},
				Sources: cli.EnvVars("DOCKLE_REJECT_KEYS"),
				Usage:   "For CIS-DI-0010. You can add sensitive keywords to look for. e.g) -ak api_password -sw keys",
			},
			&cli.StringSliceFlag{
				Name:    "accept-file",
				Aliases: []string{"af"},
				Sources: cli.EnvVars("DOCKLE_ACCEPT_FILES"),
				Usage:   "For CIS-DI-0010. You can add acceptable file names. e.g) -af id_rsa -af config.json",
			},
			&cli.StringSliceFlag{
				Name:    "sensitive-file",
				Aliases: []string{"sf"},
				Sources: cli.EnvVars("DOCKLE_REJECT_FILES"),
				Usage:   "For CIS-DI-0010. You can add sensitive files to look for. e.g) -sf .git",
			},
			&cli.StringSliceFlag{
				Name:    "accept-file-extension",
				Aliases: []string{"ae"},
				Sources: cli.EnvVars("DOCKLE_ACCEPT_FILE_EXTENSIONS"),
				Usage:   "For CIS-DI-0010. You can add acceptable file extensions. e.g) -ae pem -ae log",
			},
			&cli.StringSliceFlag{
				Name:    "sensitive-file-extension",
				Aliases: []string{"se"},
				Sources: cli.EnvVars("DOCKLE_REJECT_FILE_EXTENSIONS"),
				Usage:   "For CIS-DI-0010. You can add sensitive files to look for. e.g) -se .pfx",
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   "",
				Sources: cli.EnvVars("DOCKLE_OUTPUT_FORMAT"),
				Usage:   "output format (list, json, sarif)",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Sources: cli.EnvVars("DOCKLE_OUTPUT_FILE"),
				Usage:   "output file name",
			},
			&cli.IntFlag{
				Name:    "exit-code",
				Aliases: []string{"c"},
				Value:   0,
				Sources: cli.EnvVars("DOCKLE_EXIT_CODE"),
				Usage:   "exit code when alert were found",
			},
			&cli.StringFlag{
				Name:    "exit-level",
				Aliases: []string{"l"},
				Value:   "WARN",
				Sources: cli.EnvVars("DOCKLE_EXIT_LEVEL"),
				Usage:   "change ABEND level when use exit-code=1",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Sources: cli.EnvVars("DOCKLE_DEBUG"),
				Usage:   "debug mode",
			},
			&cli.BoolFlag{
				Name:    "quiet",
				Aliases: []string{"q"},
				Sources: cli.EnvVars("DOCKLE_QUIET"),
				Usage:   "suppress log output",
			},
			&cli.BoolFlag{
				Name:    "no-color",
				Sources: cli.EnvVars("NO_COLOR"),
				Usage:   "disabling color output",
			},
			&cli.BoolFlag{
				Name:    "version-check",
				Sources: cli.EnvVars("DOCKLE_VERSION_CHECK"),
				Usage:   "show an update notification",
			},

			// Registry flag
			&cli.DurationFlag{
				Name:    "timeout",
				Aliases: []string{"t"},
				Value:   time.Second * 90,
				Sources: cli.EnvVars("DOCKLE_TIMEOUT"),
				Usage:   "docker timeout. e.g) 5s, 5m...",
			},
			&cli.BoolFlag{
				Name:    "use-xdg",
				Aliases: []string{"x"},
				Sources: cli.EnvVars("USE_XDG"),
				Usage:   "Docker daemon host file XDG_RUNTIME_DIR",
			},
			&cli.StringFlag{
				Name:    "host",
				Sources: cli.EnvVars("DOCKLE_HOST"),
				Usage:   "docker daemon host",
			},
			&cli.StringFlag{
				Name:    "authurl",
				Sources: cli.EnvVars("DOCKLE_AUTH_URL"),
				Usage:   "registry authenticate url",
			},
			&cli.StringFlag{
				Name:    "username",
				Sources: cli.EnvVars("DOCKLE_USERNAME"),
				Usage:   "registry login username",
			},
			&cli.StringFlag{
				Name:    "password",
				Sources: cli.EnvVars("DOCKLE_PASSWORD"),
				Usage:   "registry login password. Using --password via CLI is insecure.",
			},
			&cli.BoolFlag{
				Name:    "insecure",
				Sources: cli.EnvVars("DOCKLE_INSECURE"),
				Usage:   "registry connect insecure",
			},
			&cli.BoolFlag{
				Name: "nonssl",
				// v1 used BoolTFlag, so the default stays true.
				Value:   true,
				Sources: cli.EnvVars("DOCKLE_NON_SSL"),
				Usage:   "registry connect without ssl",
			},
			&cli.StringFlag{
				Name:    "cert-path",
				Sources: cli.EnvVars("DOCKLE_CERT_PATH"),
				Usage:   "docker daemon certificate path",
			},
			&cli.StringFlag{
				Name:  "cache-dir",
				Usage: "cache directory",
			},
		},

		Action: Run,
	}
	return cmd
}
