package main

import (
	"log/slog"
	"os"

	"github.com/openela/mothership/base"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "pkgpatch",
		Commands: []*cli.Command{
			{
				Name:   "open",
				Usage:  "open an entry for patching",
				Action: open,
			},
			// todo(mustafa): finish the "generate" command
		},
		Flags: base.WithFlags(),
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("Could not run pkgpatch", "err", err)
	}
}
