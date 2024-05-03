// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package storage_s3

import (
	base "github.com/openela/mothership/base"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"log/slog"
	"net/url"
	"strings"
)

func FromFlags(ctx *cli.Context) (*S3, error) {
	// Parse the connection string
	parsedURI, err := url.Parse(ctx.String("storage-connection-string"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse storage connection string")
	}

	// Retrieve the bucket name
	bucket := parsedURI.Host

	// Remove the leading/trailing slashes
	bucket = strings.TrimSuffix(strings.TrimPrefix(bucket, "/"), "/")

	// Convert certain flags into environment variables so that they can be used by the AWS SDK
	base.RareUseChangeDefault("AWS_REGION", ctx.String("storage-region"))
	base.RareUseChangeDefault("AWS_ENDPOINT", ctx.String("storage-endpoint"))

	if !ctx.Bool("storage-secure") {
		base.RareUseChangeDefault("AWS_DISABLE_SSL", "true")
	}

	if ctx.Bool("storage-path-style") {
		base.RareUseChangeDefault("AWS_S3_FORCE_PATH_STYLE", "true")
	}

	slog.Info("Using S3 bucket", "bucket", bucket)

	return New(bucket)
}
