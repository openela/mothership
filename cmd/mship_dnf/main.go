package main

import (
	"github.com/openela/mothership/base"
	"github.com/openela/mothership/worker_client"
	"github.com/urfave/cli/v2"
	"log/slog"
	"os"
)

func run(ctx *cli.Context) error {
	return worker_client.Run(ctx)
}

func main() {
	flags := base.WithFlags(
		base.WithStorageFlags(),
		[]cli.Flag{
			&cli.StringFlag{
				Name:  "sync-location",
				Value: "/opt/srcs",
			},
		},
	)

	app := &cli.App{
		Name:   "mship_dnf",
		Usage:  "DNF worker for Mothership",
		Action: run,
		Flags:  flags,
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("failed to run mship_worker_server", "error", err)
		os.Exit(1)
	}
}
