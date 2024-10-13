package pkg

import (
	"time"

	"github.com/urfave/cli"
)

var (
	version = "dev"
)

/*
NewApp Factory for Dockle CLI creation.
An Enabler for programmatic usage of Dockle
*/
func NewApp() *cli.App {
	cli.AppHelpTemplate = `NAME:
  {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}
USAGE:
  {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}
VERSION:
  {{.Version}}{{end}}{{end}}{{if .Description}}
DESCRIPTION:
  {{.Description}}{{end}}{{if len .Authors}}
AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
  {{range $index, $author := .Authors}}{{if $index}}
  {{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}
OPTIONS:
  {{range $index, $option := .VisibleFlags}}{{if $index}}
  {{end}}{{$option}}{{end}}{{end}}
`
	app := cli.NewApp()

	var dockerSockPath string
	app.Name = "dockle"
	app.Version = version
	app.ArgsUsage = "image_name"

	app.Usage = "Container Image Linter for Security, Helping build the Best-Practice Docker Image, Easy to start"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input",
			Usage: "input file path instead of image name",
		},
		cli.StringSliceFlag{
			Name:   "ignore, i",
			EnvVar: "DOCKLE_IGNORES",
			Usage:  "checkpoints to ignore. You can use .dockleignore too.",
		},
		cli.StringSliceFlag{
			Name:   "accept-key, ak",
			EnvVar: "DOCKLE_ACCEPT_KEYS",
			Usage:  "For CIS-DI-0010. You can add acceptable keywords. e.g) -ak GPG_KEY -ak KEYCLOAK",
		},
		cli.StringSliceFlag{
			Name:   "sensitive-word, sw",
			EnvVar: "DOCKLE_REJECT_KEYS",
			Usage:  "For CIS-DI-0010. You can add sensitive keywords to look for. e.g) -ak api_password -sw keys",
		},
		cli.StringSliceFlag{
			Name:   "accept-file, af",
			EnvVar: "DOCKLE_ACCEPT_FILES",
			Usage:  "For CIS-DI-0010. You can add acceptable file names. e.g) -af id_rsa -af config.json",
		},
		cli.StringSliceFlag{
			Name:   "sensitive-file, sf",
			EnvVar: "DOCKLE_REJECT_FILES",
			Usage:  "For CIS-DI-0010. You can add sensitive files to look for. e.g) -sf .git",
		},
		cli.StringSliceFlag{
			Name:   "accept-file-extension, ae",
			EnvVar: "DOCKLE_ACCEPT_FILE_EXTENSIONS",
			Usage:  "For CIS-DI-0010. You can add acceptable file extensions. e.g) -ae pem -ae log",
		},
		cli.StringSliceFlag{
			Name:   "sensitive-file-extension, se",
			EnvVar: "DOCKLE_REJECT_FILE_EXTENSIONS",
			Usage:  "For CIS-DI-0010. You can add sensitive files to look for. e.g) -se .pfx",
		},
		cli.StringFlag{
			Name:   "format, f",
			Value:  "",
			EnvVar: "DOCKLE_OUTPUT_FORMAT",
			Usage:  "output format (list, json, sarif)",
		},
		cli.StringFlag{
			Name:   "output, o",
			EnvVar: "DOCKLE_OUTPUT_FILE",
			Usage:  "output file name",
		},
		cli.IntFlag{
			Name:   "exit-code, c",
			Value:  0,
			EnvVar: "DOCKLE_EXIT_CODE",
			Usage:  "exit code when alert were found",
		},
		cli.StringFlag{
			Name:   "exit-level, l",
			Value:  "WARN",
			EnvVar: "DOCKLE_EXIT_LEVEL",
			Usage:  "change ABEND level when use exit-code=1",
		},
		cli.BoolFlag{
			Name:   "debug, d",
			EnvVar: "DOCKLE_DEBUG",
			Usage:  "debug mode",
		},
		cli.BoolFlag{
			Name:   "quiet, q",
			EnvVar: "DOCKLE_QUIET",
			Usage:  "suppress log output",
		},
		cli.BoolFlag{
			Name:   "no-color",
			EnvVar: "NO_COLOR",
			Usage:  "disabling color output",
		},
		cli.BoolFlag{
			Name:   "version-check",
			EnvVar: "DOCKLE_VERSION_CHECK",
			Usage:  "show an update notification",
		},

		// Registry flag
		cli.DurationFlag{
			Name:   "timeout, t",
			Value:  time.Second * 90,
			EnvVar: "DOCKLE_TIMEOUT",
			Usage:  "docker timeout. e.g) 5s, 5m...",
		},
		cli.BoolFlag{
			Name:   "use-xdg, x",
			EnvVar: "USE_XDG",
			Usage:  "Docker daemon host file XDG_RUNTIME_DIR",
		},
		cli.StringFlag{
			Name:   "host",
			EnvVar: "DOCKLE_HOST",
			Usage:  "docker daemon host",
			Value:  dockerSockPath,
		},
		cli.StringFlag{
			Name:   "authurl",
			EnvVar: "DOCKLE_AUTH_URL",
			Usage:  "registry authenticate url",
		},
		cli.StringFlag{
			Name:   "username",
			EnvVar: "DOCKLE_USERNAME",
			Usage:  "registry login username",
		},
		cli.StringFlag{
			Name:   "password",
			EnvVar: "DOCKLE_PASSWORD",
			Usage:  "registry login password. Using --password via CLI is insecure.",
		},
		cli.BoolFlag{
			Name:   "insecure",
			EnvVar: "DOCKLE_INSECURE",
			Usage:  "registry connect insecure",
		},
		cli.BoolTFlag{
			Name:   "nonssl",
			EnvVar: "DOCKLE_NON_SSL",
			Usage:  "registry connect without ssl",
		},
		cli.StringFlag{
			Name:   "cert-path",
			EnvVar: "DOCKLE_CERT_PATH",
			Usage:  "docker daemon certificate path",
			Value:  dockerSockPath,
		},
		cli.StringFlag{
			Name:  "cache-dir",
			Usage: "cache directory",
		},
	}

	app.Action = Run
	return app
}
