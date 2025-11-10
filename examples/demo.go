package main

import (
	"errors"

	tl "github.com/meeeraaakiii/tintlog/logger"
	xerr "github.com/meeeraaakiii/xerr"
)

func main() {
	// Simulate a cause
	cause := errors.New("failed to connect to database")

	// Wrap with message + context
	e := xerr.NewError(cause, "initialization failed", map[string]any{
		"dsn":     "postgres://user@host/db",
		"retries": 3,
	})

	// Print (no exit): choose error type, log level, and stop code
	e.PrintWithContext(xerr.ErrorTypeError, tl.Critical, 1)
}
