// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/openela/mothership/base"
	mothership_migrations "github.com/openela/mothership/migrations"
	mothership_rpc "github.com/openela/mothership/rpc"
	"github.com/urfave/cli/v2"
	"go.temporal.io/sdk/client"
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
		Commands: []*cli.Command{
			{
				Name:  "migrate",
				Usage: "run database migrations",
				Flags: base.WithFlags(
					base.WithDatabaseFlags("mothership"),
				),
				Action: func(ctx *cli.Context) error {
					db := base.GetDBFromFlags(ctx)
					c := &postgres.Config{}
					instance, err := postgres.WithInstance(db.DB().DB, c)
					if err != nil {
						return err
					}

					// Write all SQL files to temp directory
					tempDir, err := os.MkdirTemp("", "migrations")
					if err != nil {
						return err
					}

					ls, err := mothership_migrations.UpSQLs.ReadDir(".")
					if err != nil {
						return err
					}

					for _, f := range ls {
						b, err := mothership_migrations.UpSQLs.ReadFile(f.Name())
						if err != nil {
							return err
						}

						name := strings.TrimPrefix(f.Name(), "./")
						if err := os.WriteFile(tempDir+"/"+name, b, 0644); err != nil {
							return err
						}
					}

					m, err := migrate.NewWithDatabaseInstance("file:///"+tempDir, c.DatabaseName, instance)
					if err != nil {
						return err
					}

					if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
						return err
					}

					slog.Info("migrations ran successfully")

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("failed to start mship_api", "error", err)
		os.Exit(1)
	}
}
