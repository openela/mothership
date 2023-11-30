package main

import (
	mothershipadmin_rpc "github.com/openela/mothership/admin/rpc"
	"github.com/openela/mothership/base"
	"github.com/urfave/cli/v2"
	"go.temporal.io/sdk/client"
	"log/slog"
	"os"
)

func run(ctx *cli.Context) error {
	oidcInterceptorDetails, err := base.FlagsToOidcInterceptorDetails(ctx)
	if err != nil {
		return err
	}

	temporalClient, err := base.GetTemporalClientFromFlags(ctx, client.Options{})
	if err != nil {
		return err
	}

	s, err := mothershipadmin_rpc.NewServer(
		base.GetDBFromFlags(ctx),
		temporalClient,
		oidcInterceptorDetails,
		base.FlagsToGRPCServerOptions(ctx)...,
	)
	if err != nil {
		return err
	}
	return s.Start()
}

func main() {
	app := &cli.App{
		Name:   "mship_admin_server",
		Action: run,
		Flags: base.WithFlags(
			base.WithDatabaseFlags("mothership"),
			base.WithTemporalFlags("", "mship_worker_server"),
			base.WithGrpcFlags(6687),
			base.WithGatewayFlags(6688),
			base.WithOidcFlags("", "releng"),
		),
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("failed to start mship_admin_server", "error", err)
		os.Exit(1)
	}
}
