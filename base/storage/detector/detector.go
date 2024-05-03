// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package storage_detector

import (
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/openela/mothership/base/storage"
	storage_memory "github.com/openela/mothership/base/storage/memory"
	storage_s3 "github.com/openela/mothership/base/storage/s3"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"net/url"
)

func FromFlags(ctx *cli.Context) (storage.Storage, error) {
	parsedURI, err := url.Parse(ctx.String("storage-connection-string"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse storage connection string")
	}

	switch parsedURI.Scheme {
	case "s3":
		return storage_s3.FromFlags(ctx)
	case "memory":
		return storage_memory.New(osfs.New("/")), nil
	default:
		return nil, errors.Errorf("unknown storage scheme: %s", parsedURI.Scheme)
	}
}
