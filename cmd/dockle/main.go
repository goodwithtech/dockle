package main

import (
	"github.com/goodwithtech/dockle/pkg"
	"github.com/goodwithtech/dockle/pkg/log"
	l "log"
	"os"
)

func main() {
	app := pkg.NewApp()
	err := app.Run(os.Args)

	if err != nil {
		if log.Logger != nil {
			log.Fatal(err)
		}
		l.Fatal(err)
	}
}
