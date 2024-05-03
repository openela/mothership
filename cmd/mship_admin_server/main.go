// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log/slog"
	"os"

	mothershipadmin_rpc "github.com/openela/mothership/admin/rpc"
	"github.com/openela/mothership/base"
	"github.com/urfave/cli/v2"
	"go.temporal.io/sdk/client"
)

func run(ctx *cli.Context) error {
	if ctx.String("github-team") == "" {
		return cli.ShowAppHelp(ctx)
	}

	temporalClient, err := base.GetTemporalClientFromFlags(ctx, client.Options{})
	if err != nil {
		return err
	}

	s, err := mothershipadmin_rpc.NewServer(
		base.GetDBFromFlags(ctx),
		temporalClient,
		ctx.String("github-team"),
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
			[]cli.Flag{
				&cli.StringFlag{
					Name:     "github-team",
					Usage:    "The user should be in this github team to access the admin server",
					EnvVars:  []string{"GITHUB_TEAM"},
					Required: true,
				},
			},
		),
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("failed to start mship_admin_server", "error", err)
		os.Exit(1)
	}
}
