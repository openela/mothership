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
