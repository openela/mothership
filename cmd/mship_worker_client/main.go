// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"crypto/tls"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/openela/mothership/base"
	storage_detector "github.com/openela/mothership/base/storage/detector"
	mothershippb "github.com/openela/mothership/proto/v1"
	"github.com/openela/mothership/worker_client"
	"github.com/openela/mothership/worker_client/state/system_state"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"gopkg.in/yaml.v3"
)

func run(ctx *cli.Context) error {
	configPath := ctx.String("config")
	slog.Info("starting mship_worker_client", "config", configPath)

	var args system_state.Args
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return errors.Wrap(err, "failed to read config file")
	}
	err = yaml.Unmarshal(configBytes, &args)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal config file")
	}

	storage, err := storage_detector.FromFlags(ctx)
	if err != nil {
		return err
	}
	args.Storage = storage

	systemState, err := system_state.New(&args)
	if err != nil {
		return errors.Wrap(err, "failed to create system state")
	}

	pollMinutes := ctx.Int("poll-minutes")
	slog.Info("polling for changes", "interval", pollMinutes)

	var creds credentials.TransportCredentials
	if ctx.Bool("insecure") {
		creds = insecure.NewCredentials()
	} else {
		creds = credentials.NewTLS(&tls.Config{})
	}
	grpcDial, err := grpc.Dial(ctx.String("api-endpoint"), grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	blobURI := ctx.String("blob-uri")

	srpmArchiverClient := mothershippb.NewSrpmArchiverClient(grpcDial)

	outgoingCtx := metadata.NewOutgoingContext(context.TODO(), metadata.Pairs("x-mship-worker-secret", args.WorkerSecret))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			err := worker_client.Run(outgoingCtx, blobURI, ctx.String("force-release"), systemState, srpmArchiverClient)
			if err != nil {
				slog.Error("failed polling for changes", "error", err)
			}

			// Sleep for the polling interval
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(pollMinutes) * time.Minute):
			}
		}
	}()

	wg.Wait()

	return nil
}

func main() {
	flags := base.WithFlags(
		base.WithStorageFlags(),
		[]cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Usage:   "Path to config file",
				EnvVars: []string{"CONFIG"},
				Value:   "/etc/mship/config.yaml",
			},
			&cli.IntFlag{
				Name:    "poll-minutes",
				Usage:   "Polling interval in minutes",
				EnvVars: []string{"POLL_MINUTES"},
				Value:   15,
			},
			&cli.StringFlag{
				Name:    "api-endpoint",
				Usage:   "API endpoint",
				EnvVars: []string{"API_ENDPOINT"},
				Value:   "imports.openela.org:443",
			},
			&cli.StringFlag{
				Name:    "blob-uri",
				Usage:   "Blob URI",
				EnvVars: []string{"BLOB_URI"},
				Value:   "s3://mship-srpm1",
			},
			&cli.StringFlag{
				Name:    "force-release",
				Usage:   "Release value to be used always instead of dynamically fetching from dnf",
				EnvVars: []string{"FORCE_RELEASE"},
			},
			&cli.BoolFlag{
				Name:    "insecure",
				Usage:   "Use insecure connection",
				EnvVars: []string{"INSECURE"},
			},
		},
	)

	app := &cli.App{
		Name:   "mship_worker_client",
		Action: run,
		Flags:  flags,
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("failed to run mship_worker_client", "error", err)
		os.Exit(1)
	}
}
