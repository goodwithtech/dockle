package pkg

import (
	"bufio"
	l "log"
	"os"
	"strings"

	"github.com/tomoyamachi/lyon/pkg/log"
	"github.com/tomoyamachi/lyon/pkg/types"
	"github.com/urfave/cli"
)

const (
	lyonIgnore = ".lyonignore"
)

func Run(c *cli.Context) (err error) {
	// cliVersion := c.App.Version
	result := types.ScanResult{}
	debug := c.Bool("debug")
	if err = log.InitLogger(debug); err != nil {
		l.Fatal(err)
	}

	exitCode := c.Int("exit-code")
	if exitCode != 0 {
		os.Exit(handleResult(result))
	}

	return nil
}

func handleResult(r types.ScanResult) (exitCode int) {
	optMap := getIgnoredOptMap()
	for key, targetErr := range r {
		// skip if ignore opt
		if _, ok := optMap[key]; ok {
			continue
		}

		if targetErr != nil {
			exitCode = 1
		}
	}

	return exitCode
}

func getIgnoredOptMap() map[string]struct{} {
	f, err := os.Open(lyonIgnore)
	if err != nil {
		// lyon must work even if there isn't ignore file
		return nil
	}

	ignoredMap := map[string]struct{}{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		ignoredMap[line] = struct{}{}
	}
	return ignoredMap
}
