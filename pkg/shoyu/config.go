// Package shoyu provides structured logging for the bento system.
//
// "Shoyu" (醤油 - soy sauce) adds flavor and transparency—our logger provides
// observability into bento workflow execution.
//
// The logger wraps Go's standard library log/slog for high-performance,
// zero-allocation logging with both structured (JSON) and human-readable
// (console) output modes.
//
// # Why log/slog?
//
// log/slog is Go's standard library structured logger (introduced in Go 1.21):
//   - Zero allocation: Minimal memory overhead
//   - Structured logging: JSON output for machine parsing
//   - Handler-based: Flexible output formatting
//   - Context support: Built-in context.Context integration
//   - Standard library: No external dependencies
//
// Learn more: https://pkg.go.dev/log/slog
//
// # Usage
//
//	logger := shoyu.New(shoyu.Config{
//	    Level:  shoyu.LevelInfo,
//	    Format: shoyu.FormatConsole,
//	})
//
//	logger.Info("Executing HTTP request",
//	    "neta_type", "http-request",
//	    "url", "https://api.example.com")
//
// # Common Pitfalls
//
// 1. Args must be key-value pairs
//
//	// ❌ BAD - odd number of arguments
//	logger.Info("message", "key")
//
//	// ✅ GOOD - even number for key-value pairs
//	logger.Info("message", "key", "value")
//
// 2. Context methods require context.Context as first parameter
//
//	// ✅ Using context-aware methods
//	ctx := context.Background()
//	logger.InfoContext(ctx, "message", "key", "value")
package shoyu

import "io"

// Config contains logger configuration.
type Config struct {
	// Level is the minimum log level to output.
	// Defaults to LevelInfo if not specified.
	Level Level

	// Format is the output format (JSON or Console).
	// Defaults to FormatConsole if not specified.
	Format Format

	// Output is where logs are written.
	// Defaults to os.Stdout if not specified.
	Output io.Writer

	// UseCharm enables beautiful charm/log formatting.
	// When true, uses charmbracelet/log instead of slog for pretty output.
	// Default: false (uses slog)
	UseCharm bool

	// OnStream is called for streaming output (optional).
	// Used for real-time output from long-running processes like Blender.
	// Critical for Phase 8: shell-command neta streaming.
	OnStream StreamCallback
}

// StreamCallback is called for each line of streaming output.
// Used by shell-command neta to provide real-time feedback
// from long-running processes like Blender renders.
type StreamCallback func(line string)

// Level represents a log level for filtering messages.
type Level string

const (
	// LevelDebug shows all messages including debug information.
	LevelDebug Level = "debug"

	// LevelInfo shows informational messages and above (default).
	LevelInfo Level = "info"

	// LevelWarn shows warnings and errors only.
	LevelWarn Level = "warn"

	// LevelError shows only error messages.
	LevelError Level = "error"
)

// Format represents the output format for log messages.
type Format string

const (
	// FormatJSON produces structured JSON output for machine parsing.
	// Use this for production environments or when logs need to be parsed.
	FormatJSON Format = "json"

	// FormatConsole produces human-readable console output.
	// Use this for development or when tailing logs in a terminal.
	FormatConsole Format = "console"
)
