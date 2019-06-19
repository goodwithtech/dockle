package main

import (
	l "log"
	"os"

	"github.com/goodwithtech/dockle/pkg"
	"github.com/goodwithtech/dockle/pkg/log"
	"github.com/urfave/cli"
)

var (
	version = "dev"
)

func main() {
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

	app.Usage = "A Simple Security Checker for Container Image, Suitable for CI"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input",
			Value: "",
			Usage: "input file path instead of image name",
		},
		cli.StringSliceFlag{
			Name:  "ignore, i",
			Usage: "A checkpoint to ignore. You can use .dockleignore too.",
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
			Name:  "exit-code",
			Usage: "Exit code when alert were found",
			Value: 0,
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "debug mode",
		},
		cli.StringFlag{
			Name:  "cache-dir",
			Usage: "cache directory",
		},
	}

	app.Action = pkg.Run
	err := app.Run(os.Args)

	if err != nil {
		if log.Logger != nil {
			log.Fatal(err)
		}
		l.Fatal(err)
	}
}
