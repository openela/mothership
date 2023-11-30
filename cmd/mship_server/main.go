package main

import (
	"github.com/openela/mothership/base"
	mothership_rpc "github.com/openela/mothership/rpc"
	"github.com/urfave/cli/v2"
	"go.temporal.io/sdk/client"
	"log/slog"
	"os"
)

func run(ctx *cli.Context) error {
	temporalClient, err := base.GetTemporalClientFromFlags(ctx, client.Options{})
	if err != nil {
		return err
	}

	s, err := mothership_rpc.NewServer(
		base.GetDBFromFlags(ctx),
		temporalClient,
		base.FlagsToGRPCServerOptions(ctx)...,
	)
	if err != nil {
		return err
	}
	return s.Start()
}

func main() {
	app := &cli.App{
		Name:   "mship_server",
		Action: run,
		Flags: base.WithFlags(
			base.WithDatabaseFlags("mothership"),
			base.WithTemporalFlags("", "mship_worker_server"),
			base.WithGrpcFlags(6677),
			base.WithGatewayFlags(6678),
		),
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("failed to start mship_api", "error", err)
		os.Exit(1)
	}
}
