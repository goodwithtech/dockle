package main

import (
	l "log"
	"os"
	"time"

	"github.com/goodwithtech/dockle/pkg"
	"github.com/goodwithtech/dockle/pkg/log"
	"github.com/urfave/cli"
)

var (
	version = "dev"
)

/*
NewDockleCli Factory for Dockle CLI creation.
An Enabler for programmatic usage of Dockle
*/
func NewDockleCli() *cli.App {
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
			Name:   "accept-file, af",
			EnvVar: "DOCKLE_ACCEPT_FILES",
			Usage:  "For CIS-DI-0010. You can add acceptable file names. e.g) -af id_rsa -af config.json",
		},
		cli.StringSliceFlag{
			Name:   "accept-file-extension, ae",
			EnvVar: "DOCKLE_ACCEPT_FILE_EXTENSIONS",
			Usage:  "For CIS-DI-0010. You can add acceptable file extensions. e.g) -ae pem -ae log",
		},
		cli.StringFlag{
			Name:  "format, f",
			Value: "",
			Usage: "format (json)",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "output file name",
		},
		cli.IntFlag{
			Name:  "exit-code, c",
			Usage: "exit code when alert were found",
			Value: 0,
		},
		cli.StringFlag{
			Name:  "exit-level, l",
			Usage: "change ABEND level when use exit-code=1",
			Value: "WARN",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "debug mode",
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "suppress log output",
		},
		cli.BoolFlag{
			Name:   "no-color",
			EnvVar: "NO_COLOR",
			Usage:  "suppress log output",
		},

		// Registry flag
		cli.DurationFlag{
			Name:   "timeout, t",
			Value:  time.Second * 90,
			EnvVar: "DOCKLE_TIMEOUT",
			Usage:  "docker timeout. e.g) 5s, 5m...",
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
			Name:  "cache-dir",
			Usage: "cache directory",
		},
	}

	app.Action = pkg.Run
	return app
}

func main() {
	app := NewDockleCli()
	err := app.Run(os.Args)

	if err != nil {
		if log.Logger != nil {
			log.Fatal(err)
		}
		l.Fatal(err)
	}
}
