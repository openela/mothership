package main

import (
	"github.com/openela/mothership/base"
	mshipadminpb "github.com/openela/mothership/proto/admin/v1"
	mothershippb "github.com/openela/mothership/proto/v1"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"os"
)

func run(ctx *cli.Context) error {
	frontendInfo := base.FlagsToFrontendInfo(ctx)

	// Set ctx
	if instanceNameCli := ctx.String("instance-name"); instanceNameCli != "" {
		instanceName = instanceNameCli
	}

	srpmArchiverDial, err := grpc.Dial(
		ctx.String("mothership-public-api"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	srpmArchiver := mothershippb.NewSrpmArchiverClient(srpmArchiverDial)

	mshipAdminDial, err := grpc.Dial(
		ctx.String("mothership-admin-api"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	mshipAdmin := mshipadminpb.NewMshipAdminClient(mshipAdminDial)

	s := &server{
		srpmArchiver: srpmArchiver,
		mshipAdmin:   mshipAdmin,
		frontendInfo: frontendInfo,
	}

	return s.run()
}

func main() {
	app := &cli.App{
		Name:   "mship_ui",
		Action: run,
		Flags: base.WithFlags(
			base.WithFrontendFlags(9111),
			base.WithFrontendAuthFlags(""),
			[]cli.Flag{
				&cli.StringFlag{
					Name:    "mothership-public-api",
					Usage:   "mothership public api",
					EnvVars: []string{"MOTHERSHIP_PUBLIC_API"},
					Value:   "localhost:6677",
				},
				&cli.StringFlag{
					Name:    "mothership-admin-api",
					Usage:   "mothership admin api",
					EnvVars: []string{"MOTHERSHIP_ADMIN_API"},
					Value:   "localhost:6687",
				},
				&cli.StringFlag{
					Name:    "instance-name",
					Usage:   "instance name",
					EnvVars: []string{"INSTANCE_NAME"},
				},
			},
		),
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("failed to start mship_ui", "error", err)
		os.Exit(1)
	}
}
