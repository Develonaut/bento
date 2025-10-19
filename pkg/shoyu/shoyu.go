package shoyu

import (
	"context"
	"log/slog"
	"os"
)

// Logger wraps slog.Logger with bento-specific functionality.
// It provides structured logging with support for both JSON and
// human-readable console output, as well as streaming for long-running
// processes like Blender renders.
type Logger struct {
	sl       *slog.Logger
	config   Config
	onStream StreamCallback
}

// New creates a new Logger with the given configuration.
// Defaults are applied for missing configuration values:
//   - Output: os.Stdout
//   - Level: LevelInfo
//   - Format: FormatConsole
func New(cfg Config) *Logger {
	cfg = applyDefaults(cfg)
	handler := createHandler(cfg)
	sl := slog.New(handler)

	return &Logger{
		sl:       sl,
		config:   cfg,
		onStream: cfg.OnStream,
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
// Args must be provided as alternating keys and values.
func (l *Logger) Info(msg string, args ...any) {
	l.sl.Info(msg, args...)
}

// Debug logs a debug message with optional key-value pairs.
// Args must be provided as alternating keys and values.
func (l *Logger) Debug(msg string, args ...any) {
	l.sl.Debug(msg, args...)
}

// Warn logs a warning message with optional key-value pairs.
// Args must be provided as alternating keys and values.
func (l *Logger) Warn(msg string, args ...any) {
	l.sl.Warn(msg, args...)
}

// Error logs an error message with optional key-value pairs.
// Args must be provided as alternating keys and values.
func (l *Logger) Error(msg string, args ...any) {
	l.sl.Error(msg, args...)
}

// InfoContext logs an informational message with context.
// This is the preferred method when context is available.
func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.sl.InfoContext(ctx, msg, args...)
}

// DebugContext logs a debug message with context.
func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.sl.DebugContext(ctx, msg, args...)
}

// WarnContext logs a warning message with context.
func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.sl.WarnContext(ctx, msg, args...)
}

// ErrorContext logs an error message with context.
func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.sl.ErrorContext(ctx, msg, args...)
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
		sl:       l.sl.With(args...),
		config:   l.config,
		onStream: l.onStream,
	}
}

// Stream outputs a line for streaming processes (like Blender renders).
// This bypasses normal log levels and calls the OnStream callback if set.
// Critical for Phase 8: real-time output from shell-command neta.
//
// The line is also logged at debug level for record-keeping.
func (l *Logger) Stream(line string) {
	if l.onStream != nil {
		l.onStream(line)
	}

	// Also log at debug level for record-keeping
	l.sl.Debug("stream", "output", line)
}
