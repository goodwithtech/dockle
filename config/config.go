package config

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Portshift/dockle/pkg/types"

	"github.com/urfave/cli"
)

const (
	dockleIgnore = ".dockleignore"
)

var ExitLevelMap = map[string]int{
	"info":  types.InfoLevel,
	"INFO":  types.InfoLevel,
	"warn":  types.WarnLevel,
	"WARN":  types.WarnLevel,
	"fatal": types.FatalLevel,
	"FATAL": types.FatalLevel,
}

type Config struct {
	Debug                bool
	Quiet                bool
	Timeout              time.Duration
	AuthURL              string
	Username             string
	Password             string
	Token                string
	Insecure             bool
	NonSSL               bool
	ImageName            string
	LocalImage           bool
	FilePath             string
	Output               string
	Format               string
	IgnoreMap            map[string]struct{}
	ExitCode             int
	ExitLevel            int
	AcceptanceKeys       []string
	AcceptanceFiles      []string
	AcceptanceExtensions []string
	NoColor              bool
}

var Conf Config

func CreateFromCli(c *cli.Context) {
	Conf = Config{}
	args := c.Args()

	Conf.FilePath = c.String("input")
	if Conf.FilePath == "" && len(args) == 0 {
		log.Printf(`"dockle" requires at least 1 argument or --input option.`)
		cli.ShowAppHelpAndExit(c, 1)
		return
	}
	if Conf.FilePath == "" {
		Conf.ImageName = args[0]
	}
	Conf.LocalImage = c.Bool("local")
	Conf.IgnoreMap = GetIgnoreCheckpointMap(c.StringSlice("ignore"))
	Conf.Debug = c.Bool("debug")
	Conf.Quiet = c.Bool("quiet")
	Conf.Timeout = c.Duration("timeout")
	Conf.AuthURL = c.String("authurl")
	Conf.Username = c.String("username")
	Conf.Password = c.String("password")
	Conf.Token = c.String("token")
	Conf.Insecure = c.Bool("insecure")
	Conf.NonSSL = c.Bool("nonssl")
	Conf.Output = c.String("output")
	Conf.NoColor = c.Bool("no-color")
	Conf.Format = c.String("format")
	Conf.ExitCode = c.Int("exit-code")
	Conf.ExitLevel = getExitLevel(c.String("exit-level"))
	Conf.AcceptanceKeys = c.StringSlice("accept-key")
	Conf.AcceptanceFiles = c.StringSlice("accept-file")
	Conf.AcceptanceExtensions = c.StringSlice("accept-file-extension")
}

func getExitLevel(param string) (exitLevel int) {
	exitLevel, ok := ExitLevelMap[param]
	if !ok {
		return types.WarnLevel
	}
	return exitLevel
}

func GetIgnoreCheckpointMap(ignoreRules []string) map[string]struct{} {
	ignoreCheckpointMap := map[string]struct{}{}
	// from cli command
	for _, rule := range ignoreRules {
		ignoreCheckpointMap[rule] = struct{}{}
	}

	// from ignore file
	f, err := os.Open(dockleIgnore)
	if err != nil {
		log.Printf("There is no .dockleignore file")
		// dockle must work even if there isn't ignore file
		return ignoreCheckpointMap
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		log.Printf("Add new ignore code: %s", line)
		ignoreCheckpointMap[line] = struct{}{}
	}
	return ignoreCheckpointMap
}
