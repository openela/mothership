// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"fmt"
	"log/slog"
)

// LogErrorf logs an error message
// todo(mustafa): remove this and use slog.Error properly
// deprecated
func LogErrorf(format string, args ...interface{}) {
	slog.Error(fmt.Sprintf(format, args...))
}
