// everything related to printing errors to the terminal
// sometimes stop the program too
package xerr

import (
	"os"
	"runtime/debug"

	tl "github.com/tuumbleweed/tintlog/logger"
	"github.com/tuumbleweed/tintlog/palette"
)

// error (red), warning (yellow) or skip (purple)
type ErrorType string

const (
	ErrorTypeError   ErrorType = "error"
	ErrorTypeWarning ErrorType = "warning"
	ErrorTypeSkip    ErrorType = "skip"
)

/*
Check if error is nil. If not—print it with a themed color set (error/warning/skip).
If stopCode > 0, exit the program after printing.

Notes:
- Uses go-tintl colorizers (truecolor + per-line tinting).
- Message "banner" uses bold foreground on background; details use bold/normal tints.
*/
func (e *Error) PrintErrorWithOptions(
	logLevel tl.LogLevel, errorType ErrorType, stopCode int,
	printContext, printDebugStack bool,
) {
	if e == nil {
		return
	}

	// Build colorizers (local, no registry dependency for BG variants)
	errColors := map[string]palette.Colorizer{
		"msg":     palette.FgBgColorizer("ErrMsg", palette.BlackColor, palette.RedColor, true), // bold on red bg
		"err":     palette.FgColorizer("ErrErr", palette.BrightRedColor, true),                 // bold bright red
		"where":   palette.FgColorizer("ErrWhere", palette.RedColor, true),                     // bold red
		"debug":   palette.FgColorizer("ErrDebug", palette.DimRedColor, false),                 // dim red
		"context": palette.FgColorizer("ErrCtx", palette.DimRedColor, false),
	}

	warnColors := map[string]palette.Colorizer{
		"msg":     palette.FgBgColorizer("WarnMsg", palette.BlackColor, palette.YellowColor, true),
		"err":     palette.FgColorizer("WarnErr", palette.BrightYellowColor, true),
		"where":   palette.FgColorizer("WarnWhere", palette.YellowColor, true),
		"debug":   palette.FgColorizer("WarnDebug", palette.DimYellowColor, false),
		"context": palette.FgColorizer("WarnCtx", palette.DimYellowColor, false),
	}

	skipColors := map[string]palette.Colorizer{
		"msg":     palette.FgBgColorizer("SkipMsg", palette.BlackColor, palette.PurpleColor, true),
		"err":     palette.FgColorizer("SkipErr", palette.BrightPurpleColor, true),
		"where":   palette.FgColorizer("SkipWhere", palette.PurpleColor, true),
		"debug":   palette.FgColorizer("SkipDebug", palette.DimPurpleColor, false),
		"context": palette.FgColorizer("SkipCtx", palette.DimPurpleColor, false),
	}

	defColors := errColors

	var colors map[string]palette.Colorizer
	switch errorType {
	case ErrorTypeError:
		colors = errColors
	case ErrorTypeWarning:
		colors = warnColors
	case ErrorTypeSkip:
		colors = skipColors
	default:
		colors = defColors
	}

	tl.Log(logLevel, colors["msg"], "Msg: '%s'", e.Msg)
	tl.Log(logLevel+1, colors["err"], "Err: '%s'", e.ErrStr)
	tl.Log(logLevel+2, colors["where"], "Where: '%s'", e.Where)

	if printDebugStack {
		tl.Log(logLevel+3, colors["debug"], "Debug stack:\n```\n%s\n```", string(debug.Stack()))
	}
	if printContext {
		tl.Log(logLevel+4, colors["context"], "Context:\n```\n%s\n```", e.Context)
	}

	if stopCode > 0 {
		os.Exit(stopCode)
	}
}

// Wrappers around PrintErrorWithOptions

// Print prints the error without context/debug stack.
func (e *Error) Print(errorType ErrorType, logLevel tl.LogLevel, stopCode int) {
	e.PrintErrorWithOptions(logLevel, errorType, stopCode, false, false)
}

// PrintWithContext prints the error with context and a debug stack.
func (e *Error) PrintWithContext(errorType ErrorType, logLevel tl.LogLevel, stopCode int) {
	e.PrintErrorWithOptions(logLevel, errorType, stopCode, true, true)
}

// QuitIf prints the error with context/stack at Critical1 and exits(1).
func (e *Error) QuitIf(errorType ErrorType) {
	e.PrintErrorWithOptions(tl.Critical1, errorType, 1, true, true)
}
