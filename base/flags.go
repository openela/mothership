package base

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/urfave/cli/v2"
	"go.temporal.io/sdk/client"
	"log/slog"
	"os"
)

type EnvVar string

const (
	EnvVarGRPCPort                EnvVar = "GRPC_PORT"
	EnvVarGatewayPort             EnvVar = "GATEWAY_PORT"
	EnvVarDatabaseURI             EnvVar = "DATABASE_URI"
	EnvVarPort                    EnvVar = "PORT"
	EnvVarOIDCIssuer              EnvVar = "OIDC_ISSUER"
	EnvVarOIDCClientID            EnvVar = "OIDC_CLIENT_ID"
	EnvVarOIDCClientSecret        EnvVar = "OIDC_CLIENT_SECRET"
	EnvVarOIDCUserInfoOverride    EnvVar = "OIDC_USERINFO_OVERRIDE"
	EnvVarRequiredOIDCGroup       EnvVar = "REQUIRED_OIDC_GROUP"
	EnvVarTemporalNamespace       EnvVar = "TEMPORAL_NAMESPACE"
	EnvVarTemporalAddress         EnvVar = "TEMPORAL_ADDRESS"
	EnvVarTemporalTaskQueue       EnvVar = "TEMPORAL_TASK_QUEUE"
	EnvVarSelf                    EnvVar = "SELF"
	EnvVarStorageEndpoint         EnvVar = "STORAGE_ENDPOINT"
	EnvVarStorageConnectionString EnvVar = "STORAGE_CONNECTION_STRING"
	EnvVarStorageRegion           EnvVar = "STORAGE_REGION"
	EnvVarStorageSecure           EnvVar = "STORAGE_SECURE"
	EnvVarStoragePathStyle        EnvVar = "STORAGE_PATH_STYLE"
)

func WithDatabaseFlags(appName string) []cli.Flag {
	if appName == "" {
		appName = "root"
	}

	return []cli.Flag{
		&cli.StringFlag{
			Name:    "database-uri",
			Aliases: []string{"d"},
			Usage:   "database uri",
			EnvVars: []string{string(EnvVarDatabaseURI)},
			Value:   "postgres://postgres:postgres@localhost:5432/" + appName + "?sslmode=disable",
		},
	}
}

func WithTemporalFlags(defaultNamespace string, defaultTaskQueue string) []cli.Flag {
	if defaultNamespace == "" {
		defaultNamespace = "default"
	}

	return []cli.Flag{
		&cli.StringFlag{
			Name:    "temporal-namespace",
			Aliases: []string{"n"},
			Usage:   "temporal namespace",
			EnvVars: []string{string(EnvVarTemporalNamespace)},
			Value:   defaultNamespace,
		},
		&cli.StringFlag{
			Name:    "temporal-address",
			Aliases: []string{"a"},
			Usage:   "temporal address",
			EnvVars: []string{string(EnvVarTemporalAddress)},
			Value:   "localhost:7233",
		},
		&cli.StringFlag{
			Name:    "temporal-task-queue",
			Aliases: []string{"q"},
			Usage:   "temporal task queue",
			EnvVars: []string{string(EnvVarTemporalTaskQueue)},
			Value:   defaultTaskQueue,
		},
	}
}

func WithGrpcFlags(defaultPort int) []cli.Flag {
	if defaultPort == 0 {
		defaultPort = 8080
	}

	return []cli.Flag{
		&cli.IntFlag{
			Name:    "grpc-port",
			Usage:   "gRPC port",
			EnvVars: []string{string(EnvVarGRPCPort)},
			Value:   defaultPort,
		},
	}
}

func WithGatewayFlags(defaultPort int) []cli.Flag {
	if defaultPort == 0 {
		defaultPort = 8081
	}

	return []cli.Flag{
		&cli.IntFlag{
			Name:    "gateway-port",
			Usage:   "gRPC gateway port",
			EnvVars: []string{string(EnvVarGatewayPort)},
			Value:   defaultPort,
		},
	}
}

func WithOidcFlags(defaultOidcIssuer string, defaultGroup string) []cli.Flag {
	if defaultOidcIssuer == "" {
		defaultOidcIssuer = "https://accounts.rockylinux.org/auth/realms/rocky"
	}

	return []cli.Flag{
		&cli.StringFlag{
			Name:    "oidc-issuer",
			Usage:   "OIDC issuer",
			EnvVars: []string{string(EnvVarOIDCIssuer)},
			Value:   defaultOidcIssuer,
		},
		&cli.StringFlag{
			Name:    "required-oidc-group",
			Usage:   "OIDC group that is required to access the frontend",
			EnvVars: []string{string(EnvVarRequiredOIDCGroup)},
			Value:   defaultGroup,
		},
	}
}

func WithFrontendFlags(defaultPort int) []cli.Flag {
	if defaultPort == 0 {
		defaultPort = 9111
	}

	return []cli.Flag{
		&cli.IntFlag{
			Name:    "port",
			Usage:   "frontend port",
			EnvVars: []string{string(EnvVarPort)},
			Value:   defaultPort,
		},
	}
}

func WithFrontendAuthFlags(defaultOidcIssuer string) []cli.Flag {
	if defaultOidcIssuer == "" {
		defaultOidcIssuer = "https://accounts.rockylinux.org/auth/realms/rocky"
	}

	return []cli.Flag{
		&cli.StringFlag{
			Name:    "oidc-issuer",
			Usage:   "OIDC issuer",
			EnvVars: []string{string(EnvVarOIDCIssuer)},
			Value:   defaultOidcIssuer,
		},
		&cli.StringFlag{
			Name:    "oidc-client-id",
			Usage:   "OIDC client ID",
			EnvVars: []string{string(EnvVarOIDCClientID)},
		},
		&cli.StringFlag{
			Name:    "oidc-client-secret",
			Usage:   "OIDC client secret",
			EnvVars: []string{string(EnvVarOIDCClientSecret)},
		},
		&cli.StringFlag{
			Name:    "oidc-userinfo-override",
			Usage:   "OIDC userinfo override",
			EnvVars: []string{string(EnvVarOIDCUserInfoOverride)},
		},
		&cli.StringFlag{
			Name:    "required-oidc-group",
			Usage:   "OIDC group that is required to access the frontend",
			EnvVars: []string{string(EnvVarRequiredOIDCGroup)},
		},
		&cli.StringFlag{
			Name:    "self",
			Usage:   "Endpoint pointing to the frontend",
			EnvVars: []string{string(EnvVarSelf)},
		},
	}
}

// WithStorageFlags adds the storage flags to the app.
func WithStorageFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "storage-endpoint",
			Usage:   "storage endpoint",
			EnvVars: []string{string(EnvVarStorageEndpoint)},
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "storage-connection-string",
			Usage:   "storage connection string",
			EnvVars: []string{string(EnvVarStorageConnectionString)},
		},
		&cli.StringFlag{
			Name:    "storage-region",
			Usage:   "storage region",
			EnvVars: []string{string(EnvVarStorageRegion)},
			// RESF default region
			Value: "us-east-2",
		},
		&cli.BoolFlag{
			Name:    "storage-secure",
			Usage:   "storage secure",
			EnvVars: []string{string(EnvVarStorageSecure)},
			Value:   true,
		},
		&cli.BoolFlag{
			Name:    "storage-path-style",
			Usage:   "storage path style",
			EnvVars: []string{string(EnvVarStoragePathStyle)},
			Value:   false,
		},
	}
}

func WithFlags(flags ...[]cli.Flag) []cli.Flag {
	var result []cli.Flag

	for _, f := range flags {
		result = append(result, f...)
	}

	return result
}

// FlagsToGRPCServerOptions converts the cli flags to gRPC server options.
func FlagsToGRPCServerOptions(ctx *cli.Context) []GRPCServerOption {
	return []GRPCServerOption{
		WithGRPCPort(ctx.Int("grpc-port")),
		WithGatewayPort(ctx.Int("gateway-port")),
	}
}

// FlagsToFrontendInfo converts the cli flags to frontend info.
func FlagsToFrontendInfo(ctx *cli.Context) *FrontendInfo {
	return &FrontendInfo{
		Title:                ctx.App.Name,
		Port:                 ctx.Int("port"),
		Self:                 ctx.String("self"),
		OIDCIssuer:           ctx.String("oidc-issuer"),
		OIDCClientID:         ctx.String("oidc-client-id"),
		OIDCClientSecret:     ctx.String("oidc-client-secret"),
		OIDCGroup:            ctx.String("required-oidc-group"),
		OIDCUserInfoOverride: ctx.String("oidc-userinfo-override"),
	}
}

// FlagsToOidcInterceptorDetails converts the cli flags to oidc interceptor details.
func FlagsToOidcInterceptorDetails(ctx *cli.Context) (*OidcInterceptorDetails, error) {
	provider, err := oidc.NewProvider(ctx.Context, ctx.String("oidc-issuer"))
	if err != nil {
		return nil, err
	}

	return &OidcInterceptorDetails{
		Provider: &OidcProviderImpl{provider},
		Group:    ctx.String("required-oidc-group"),
	}, nil
}

// GetDBFromFlags gets the database from the cli flags.
func GetDBFromFlags(ctx *cli.Context) *DB {
	// Create database.
	db, err := NewDB(ctx.String("database-uri"))
	if err != nil {
		slog.Error("failed to create database", "error", err)
		os.Exit(1)
	}

	return db
}

// GetTemporalClientFromFlags gets the temporal client from the cli flags.
func GetTemporalClientFromFlags(ctx *cli.Context, opts client.Options) (client.Client, error) {
	return NewTemporalClient(
		ctx.String("temporal-address"),
		ctx.String("temporal-namespace"),
		ctx.String("temporal-task-queue"),
		opts,
	)
}

// RareUseChangeDefault changes the default value of an arbitrary environment variable.
func RareUseChangeDefault(envVar string, newDefault string) {
	// Check if the environment variable is set.
	if _, ok := os.LookupEnv(envVar); ok {
		return
	}

	// Change the default value.
	if err := os.Setenv(envVar, newDefault); err != nil {
		slog.Error("failed to set environment variable", "error", err, "envVar", envVar)
		os.Exit(1)
	}
}