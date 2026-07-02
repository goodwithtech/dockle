package main

import (
	"context"
	l "log"
	"os"

	"github.com/goodwithtech/dockle/pkg"
	"github.com/goodwithtech/dockle/pkg/log"
)

func main() {
	app := pkg.NewApp()
	err := app.Run(context.Background(), os.Args)

	if err != nil {
		if log.Logger != nil {
			log.Fatal(err)
		}
		l.Fatal(err)
	}
}
