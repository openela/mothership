package main

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"log/slog"
	"os"

	"github.com/openela/mothership/base"
	"github.com/openela/mothership/base/bugtracker"
	github_bugtracker "github.com/openela/mothership/base/bugtracker/github"
	"github.com/openela/mothership/base/forge"
	github_forge "github.com/openela/mothership/base/forge/github"
	storage_detector "github.com/openela/mothership/base/storage/detector"
	mothership_worker_server "github.com/openela/mothership/worker_server"
	"github.com/urfave/cli/v2"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"golang.org/x/crypto/openpgp"
)

//go:embed rh_public_key.asc
var defaultGpgKey []byte

func run(ctx *cli.Context) error {
	temporalClient, err := base.GetTemporalClientFromFlags(ctx, client.Options{})
	if err != nil {
		return err
	}

	db := base.GetDBFromFlags(ctx)
	storage, err := storage_detector.FromFlags(ctx)
	if err != nil {
		return err
	}

	// Create pgp keys
	var gpgKeys openpgp.EntityList
	for _, key := range ctx.StringSlice("allowed-gpg-keys") {
		decoded, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			return err
		}
		keyRing, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(decoded))
		if err != nil {
			return err
		}

		gpgKeys = append(gpgKeys, keyRing...)
	}

	// Create forge based on git provider
	var remoteForge forge.Forge
	switch ctx.String("git-provider") {
	case "github":
		var appPrivateKey []byte
		if ctx.Bool("github-app-private-key-base64") {
			appPrivateKey, err = base64.StdEncoding.DecodeString(ctx.String("github-app-private-key"))
			if err != nil {
				return err
			}
		} else {
			appPrivateKey = []byte(ctx.String("github-app-private-key"))
		}

		remoteForge, err = github_forge.New(
			ctx.String("github-org"),
			ctx.String("github-app-id"),
			appPrivateKey,
			ctx.Bool("github-make-repo-public"),
		)
		if err != nil {
			return err
		}
	default:
		return cli.Exit("git-provider must be github", 1)
	}

	// Create bugtracker
	var remoteTracker bugtracker.Bugtracker
	switch ctx.String("bugtracker-provider") {
	case "github":
		if ctx.Bool("bugtracker-github-use-forge-auth") {
			var appPrivateKey []byte
			if ctx.Bool("github-app-private-key-base64") {
				appPrivateKey, err = base64.StdEncoding.DecodeString(ctx.String("github-app-private-key"))
				if err != nil {
					return err
				}
			} else {
				appPrivateKey = []byte(ctx.String("github-app-private-key"))
			}

			remoteTracker, err = github_bugtracker.New(
				ctx.String("bugtracker-github-repo"),
				ctx.String("github-app-id"),
				appPrivateKey,
			)
		} else {
			var appPrivateKey []byte
			if ctx.Bool("bugtracker-github-app-private-key-base64") {
				appPrivateKey, err = base64.StdEncoding.DecodeString(ctx.String("bugtracker-github-app-private-key"))
				if err != nil {
					return err
				}
			} else {
				appPrivateKey = []byte(ctx.String("bugtracker-github-app-private-key"))
			}

			remoteTracker, err = github_bugtracker.New(
				ctx.String("bugtracker-github-repo"),
				ctx.String("bugtracker-github-app-id"),
				appPrivateKey,
			)
		}
		if err != nil {
			return err
		}
	}

	remoteForge = forge.NewCacher(remoteForge)

	publicURI := ctx.String("public-uri")
	if remoteTracker != nil && publicURI == "" {
		return cli.Exit("public-uri is required if bugtracker is used", 1)
	}

	w := worker.New(temporalClient, ctx.String("temporal-task-queue"), worker.Options{})
	workerServer := mothership_worker_server.New(
		db,
		storage,
		gpgKeys,
		remoteForge,
		remoteTracker,
		ctx.Bool("import-rolling-release"),
		publicURI,
	)

	// Register workflows
	w.RegisterWorkflow(mothership_worker_server.ProcessRPMWorkflow)
	w.RegisterWorkflow(mothership_worker_server.RetractEntryWorkflow)
	w.RegisterWorkflow(mothership_worker_server.SealBatchWorkflow)

	// Register activities
	w.RegisterActivity(workerServer)

	// Start worker
	return w.Run(worker.InterruptCh())
}

func main() {
	flags := base.WithFlags(
		base.WithDatabaseFlags("mothership"),
		base.WithTemporalFlags("", "mship_worker_server"),
		base.WithStorageFlags(),
		[]cli.Flag{
			&cli.StringSliceFlag{
				Name:    "allowed-gpg-keys",
				Usage:   "Armored GPG keys that we verify SRPMs with. Must be base64 encoded",
				EnvVars: []string{"ALLOWED_GPG_KEYS"},
			},
			&cli.BoolFlag{
				Name:    "import-rolling-release",
				Usage:   "Whether to import packages in rolling release mode",
				EnvVars: []string{"IMPORT_ROLLING_RELEASE"},
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "public-uri",
				Usage:   "Public URI for Mothership UI. Required if bugtracker is used.",
				EnvVars: []string{"PUBLIC_URI"},
			},
			&cli.StringFlag{
				Name: "git-provider",
				Action: func(ctx *cli.Context, s string) error {
					// Can only be github for now
					if s != "github" {
						return cli.Exit("git-provider must be github", 1)
					}

					return nil
				},
				Usage:   "Git provider to use. Currently only github is supported",
				EnvVars: []string{"GIT_PROVIDER"},
			},
			// Github only
			&cli.StringFlag{
				Name:    "github-org",
				Usage:   "Github organization to use",
				EnvVars: []string{"GITHUB_ORG"},
				Action: func(ctx *cli.Context, s string) error {
					// Required for github
					if ctx.String("git-provider") == "github" && s == "" {
						return cli.Exit("github-org is required for github", 1)
					}

					return nil
				},
			},
			&cli.StringFlag{
				Name:    "github-app-id",
				Usage:   "Github app ID",
				EnvVars: []string{"GITHUB_APP_ID"},
				Action: func(ctx *cli.Context, s string) error {
					// Required for github
					if ctx.String("git-provider") == "github" && s == "" {
						return cli.Exit("github-org is required for github", 1)
					}

					return nil
				},
			},
			&cli.StringFlag{
				Name:    "github-app-private-key",
				Usage:   "Github app private key",
				EnvVars: []string{"GITHUB_APP_PRIVATE_KEY"},
				Action: func(ctx *cli.Context, s string) error {
					// Required for github
					if ctx.String("git-provider") == "github" && s == "" {
						return cli.Exit("github-org is required for github", 1)
					}

					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "github-app-private-key-base64",
				Usage:   "Whether the Github app private key is base64 encoded",
				EnvVars: []string{"GITHUB_APP_PRIVATE_KEY_BASE64"},
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "github-make-repo-public",
				Usage:   "Whether to make the Github repository public",
				EnvVars: []string{"GITHUB_MAKE_REPO_PUBLIC"},
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "bugtracker-provider",
				Usage:   "Bugtracker provider to use. Currently only github is supported",
				EnvVars: []string{"BUGTRACKER_PROVIDER"},
				Value:   "github",
			},
			&cli.StringFlag{
				Name:    "bugtracker-github-repo",
				Usage:   "Github repository to use for bugtracker",
				EnvVars: []string{"BUGTRACKER_GITHUB_REPO"},
			},
			&cli.StringFlag{
				Name:    "bugtracker-github-app-id",
				Usage:   "Github app ID for bugtracker",
				EnvVars: []string{"BUGTRACKER_GITHUB_APP_ID"},
			},
			&cli.StringFlag{
				Name:    "bugtracker-github-app-private-key",
				Usage:   "Github app private key for bugtracker",
				EnvVars: []string{"BUGTRACKER_GITHUB_APP_PRIVATE_KEY"},
			},
			&cli.BoolFlag{
				Name:    "bugtracker-github-app-private-key-base64",
				Usage:   "Whether the Github app private key for bugtracker is base64 encoded",
				EnvVars: []string{"BUGTRACKER_GITHUB_APP_PRIVATE_KEY_BASE64"},
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "bugtracker-github-use-forge-auth",
				Usage:   "Whether to use forge authentication for bugtracker",
				EnvVars: []string{"BUGTRACKER_GITHUB_USE_FORGE_AUTH"},
				Value:   false,
			},
		},
	)

	base64EncodedDefaultGpgKey := base64.StdEncoding.EncodeToString(defaultGpgKey)
	base.RareUseChangeDefault("ALLOWED_GPG_KEYS", base64EncodedDefaultGpgKey)

	app := &cli.App{
		Name:   "mship_worker_server",
		Action: run,
		Flags:  flags,
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("failed to run mship_worker_server", "error", err)
		os.Exit(1)
	}
}
