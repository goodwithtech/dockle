package config

import (
	"bufio"
	"os"
	"strings"

	"github.com/goodwithtech/dockle/pkg/types"

	"github.com/goodwithtech/dockle/pkg/log"
	"github.com/urfave/cli/v3"
)

const (
	dockleIgnore = ".dockleignore"
)

var exitLevelMap = map[string]int{
	"info":  types.InfoLevel,
	"INFO":  types.InfoLevel,
	"warn":  types.WarnLevel,
	"WARN":  types.WarnLevel,
	"fatal": types.FatalLevel,
	"FATAL": types.FatalLevel,
}

type Config struct {
	IgnoreMap map[string]struct{}
	ExitCode  int
	ExitLevel int
}

var Conf Config

func CreateFromCli(cmd *cli.Command) {
	ignoreRules := cmd.StringSlice("ignore")
	Conf = Config{
		IgnoreMap: getIgnoreCheckpointMap(ignoreRules),
		ExitCode:  cmd.Int("exit-code"),
		ExitLevel: getExitLevel(cmd.String("exit-level")),
	}
}

func getExitLevel(param string) (exitLevel int) {
	exitLevel, ok := exitLevelMap[param]
	if !ok {
		return types.WarnLevel
	}
	return exitLevel
}

func getIgnoreCheckpointMap(ignoreRules []string) map[string]struct{} {
	ignoreCheckpointMap := map[string]struct{}{}
	// from cli command
	for _, rule := range ignoreRules {
		ignoreCheckpointMap[rule] = struct{}{}
	}

	// from ignore file
	f, err := os.Open(dockleIgnore)
	if err != nil {
		log.Logger.Debug("There is no .dockleignore file")
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
		log.Logger.Debugf("Add new ignore code: %s", line)
		ignoreCheckpointMap[line] = struct{}{}
	}
	return ignoreCheckpointMap
}
