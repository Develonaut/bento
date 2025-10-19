package shoyu

import (
	"context"
	"io"
	"os"
)

// logger is the internal interface that both slog and charm implementations satisfy.
type logger interface {
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	DebugContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
	With(args ...any) logger
	Stream(line string)
	SetOutput(w io.Writer)
}

// Logger wraps either slog.Logger or charm/log with bento-specific functionality.
// It provides structured logging with support for both JSON and
// human-readable console output, as well as streaming for long-running
// processes like Blender renders.
type Logger struct {
	impl logger
}

// New creates a new Logger with the given configuration.
// Defaults are applied for missing configuration values:
//   - Output: os.Stdout
//   - Level: LevelInfo
//   - Format: FormatConsole
//   - UseCharm: false
func New(cfg Config) *Logger {
	cfg = applyDefaults(cfg)

	var impl logger
	if cfg.UseCharm {
		impl = newCharmLogger(cfg)
	} else {
		impl = newSlogLogger(cfg)
	}

	return &Logger{
		impl: impl,
	}
}

// applyDefaults sets default values for missing config fields.
func applyDefaults(cfg Config) Config {
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}

	if cfg.Level == "" {
		cfg.Level = LevelInfo
	}

	if cfg.Format == "" {
		cfg.Format = FormatConsole
	}

	return cfg
}

// Info logs an informational message with optional key-value pairs.
// Args must be provided as alternating keys and values.
func (l *Logger) Info(msg string, args ...any) {
	l.impl.Info(msg, args...)
}

// Debug logs a debug message with optional key-value pairs.
// Args must be provided as alternating keys and values.
func (l *Logger) Debug(msg string, args ...any) {
	l.impl.Debug(msg, args...)
}

// Warn logs a warning message with optional key-value pairs.
// Args must be provided as alternating keys and values.
func (l *Logger) Warn(msg string, args ...any) {
	l.impl.Warn(msg, args...)
}

// Error logs an error message with optional key-value pairs.
// Args must be provided as alternating keys and values.
func (l *Logger) Error(msg string, args ...any) {
	l.impl.Error(msg, args...)
}

// InfoContext logs an informational message with context.
// This is the preferred method when context is available.
func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.impl.InfoContext(ctx, msg, args...)
}

// DebugContext logs a debug message with context.
func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.impl.DebugContext(ctx, msg, args...)
}

// WarnContext logs a warning message with context.
func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.impl.WarnContext(ctx, msg, args...)
}

// ErrorContext logs an error message with context.
func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.impl.ErrorContext(ctx, msg, args...)
}

// With creates a child logger with additional context fields.
// This is used by itamae to add trace IDs, bento IDs, neta IDs, etc.
//
// Example:
//
//	contextLogger := logger.With(
//	    "bento_id", "my-workflow",
//	    "neta_id", "node-1")
func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		impl: l.impl.With(args...),
	}
}

// Stream outputs a line for streaming processes (like Blender renders).
// This bypasses normal log levels and calls the OnStream callback if set.
// Critical for Phase 8: real-time output from shell-command neta.
//
// The line is also logged at debug level for record-keeping.
func (l *Logger) Stream(line string) {
	l.impl.Stream(line)
}

// SetOutput changes the output destination.
// Note: This may not work for all backend implementations.
func (l *Logger) SetOutput(w io.Writer) {
	l.impl.SetOutput(w)
}
