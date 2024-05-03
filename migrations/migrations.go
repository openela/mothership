// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package mothership_migrations

import "embed"

//go:embed *.up.sql
var UpSQLs embed.FS
