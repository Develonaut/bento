// Package shoyu provides structured logging for the bento system.
//
// Standard library slog integration.
package shoyu

import (
	"context"
	"io"
	"log/slog"
)

// slogLogger wraps slog.Logger to implement our logging interface.
type slogLogger struct {
	sl       *slog.Logger
	output   io.Writer
	onStream StreamCallback
}

// newSlogLogger creates an slog-based logger.
func newSlogLogger(cfg Config) *slogLogger {
	handler := createHandler(cfg)
	sl := slog.New(handler)

	return &slogLogger{
		sl:       sl,
		output:   cfg.Output,
		onStream: cfg.OnStream,
	}
}

// createHandler creates an slog.Handler based on format and level.
func createHandler(cfg Config) slog.Handler {
	level := convertLevel(cfg.Level)
	opts := &slog.HandlerOptions{
		Level: level,
	}

	switch cfg.Format {
	case FormatJSON:
		return slog.NewJSONHandler(cfg.Output, opts)
	case FormatConsole:
		return slog.NewTextHandler(cfg.Output, opts)
	default:
		return slog.NewTextHandler(cfg.Output, opts)
	}
}

// convertLevel converts our Level type to slog.Level.
func convertLevel(level Level) slog.Level {
	switch level {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Info logs an informational message with optional key-value pairs.
func (sl *slogLogger) Info(msg string, args ...any) {
	sl.sl.Info(msg, args...)
}

// Debug logs a debug message with optional key-value pairs.
func (sl *slogLogger) Debug(msg string, args ...any) {
	sl.sl.Debug(msg, args...)
}

// Warn logs a warning message with optional key-value pairs.
func (sl *slogLogger) Warn(msg string, args ...any) {
	sl.sl.Warn(msg, args...)
}

// Error logs an error message with optional key-value pairs.
func (sl *slogLogger) Error(msg string, args ...any) {
	sl.sl.Error(msg, args...)
}

// InfoContext logs an informational message with context.
func (sl *slogLogger) InfoContext(ctx context.Context, msg string, args ...any) {
	sl.sl.InfoContext(ctx, msg, args...)
}

// DebugContext logs a debug message with context.
func (sl *slogLogger) DebugContext(ctx context.Context, msg string, args ...any) {
	sl.sl.DebugContext(ctx, msg, args...)
}

// WarnContext logs a warning message with context.
func (sl *slogLogger) WarnContext(ctx context.Context, msg string, args ...any) {
	sl.sl.WarnContext(ctx, msg, args...)
}

// ErrorContext logs an error message with context.
func (sl *slogLogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	sl.sl.ErrorContext(ctx, msg, args...)
}

// With creates a child logger with additional context fields.
func (sl *slogLogger) With(args ...any) logger {
	return &slogLogger{
		sl:       sl.sl.With(args...),
		output:   sl.output,
		onStream: sl.onStream,
	}
}

// Stream outputs a line for streaming processes.
func (sl *slogLogger) Stream(line string) {
	if sl.onStream != nil {
		sl.onStream(line)
	}

	// Also log at debug level for record-keeping
	sl.sl.Debug("stream", "output", line)
}

// SetOutput changes the output destination.
func (sl *slogLogger) SetOutput(w io.Writer) {
	sl.output = w
	// Note: slog doesn't support changing output after creation
	// This is a limitation of the adapter pattern
}
