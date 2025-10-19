// Package shoyu provides structured logging for the bento system.
//
// Charm/log integration for beautiful terminal output.
package shoyu

import (
	"context"
	"io"

	"github.com/charmbracelet/log"
)

// charmLogger wraps charmbracelet/log to implement our logging interface.
type charmLogger struct {
	cl       *log.Logger
	onStream StreamCallback
}

// newCharmLogger creates a charm-based logger.
func newCharmLogger(cfg Config) *charmLogger {
	cl := log.NewWithOptions(cfg.Output, log.Options{
		Level:           convertLevelToCharm(cfg.Level),
		ReportTimestamp: true,
		ReportCaller:    false,
	})

	return &charmLogger{
		cl:       cl,
		onStream: cfg.OnStream,
	}
}

// convertLevelToCharm converts our Level type to charm/log Level.
func convertLevelToCharm(level Level) log.Level {
	switch level {
	case LevelDebug:
		return log.DebugLevel
	case LevelInfo:
		return log.InfoLevel
	case LevelWarn:
		return log.WarnLevel
	case LevelError:
		return log.ErrorLevel
	default:
		return log.InfoLevel
	}
}

// Info logs an informational message with optional key-value pairs.
func (cl *charmLogger) Info(msg string, args ...any) {
	cl.cl.Info(msg, argsToKeyvals(args)...)
}

// Debug logs a debug message with optional key-value pairs.
func (cl *charmLogger) Debug(msg string, args ...any) {
	cl.cl.Debug(msg, argsToKeyvals(args)...)
}

// Warn logs a warning message with optional key-value pairs.
func (cl *charmLogger) Warn(msg string, args ...any) {
	cl.cl.Warn(msg, argsToKeyvals(args)...)
}

// Error logs an error message with optional key-value pairs.
func (cl *charmLogger) Error(msg string, args ...any) {
	cl.cl.Error(msg, argsToKeyvals(args)...)
}

// InfoContext logs an informational message with context.
// Note: charm/log doesn't have native context support, so we ignore context.
func (cl *charmLogger) InfoContext(ctx context.Context, msg string, args ...any) {
	cl.Info(msg, args...)
}

// DebugContext logs a debug message with context.
func (cl *charmLogger) DebugContext(ctx context.Context, msg string, args ...any) {
	cl.Debug(msg, args...)
}

// WarnContext logs a warning message with context.
func (cl *charmLogger) WarnContext(ctx context.Context, msg string, args ...any) {
	cl.Warn(msg, args...)
}

// ErrorContext logs an error message with context.
func (cl *charmLogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	cl.Error(msg, args...)
}

// With creates a child logger with additional context fields.
func (cl *charmLogger) With(args ...any) logger {
	return &charmLogger{
		cl:       cl.cl.With(argsToKeyvals(args)...),
		onStream: cl.onStream,
	}
}

// Stream outputs a line for streaming processes.
func (cl *charmLogger) Stream(line string) {
	if cl.onStream != nil {
		cl.onStream(line)
	}

	// Also log at debug level for record-keeping
	cl.cl.Debug("stream", "output", line)
}

// SetOutput changes the output destination.
func (cl *charmLogger) SetOutput(w io.Writer) {
	cl.cl.SetOutput(w)
}

// argsToKeyvals converts variadic args to key-value pairs.
// Charm/log expects the same format as slog (alternating keys and values).
func argsToKeyvals(args []any) []any {
	return args
}
