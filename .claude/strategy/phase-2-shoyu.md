# Phase 2: Shoyu Package (醤油 - "Soy Sauce")

**Duration:** 3-4 days
**Package:** `pkg/shoyu/`
**Dependencies:** None

---

## TDD Philosophy

> **Write tests FIRST to define contracts**

The logger is a critical piece of infrastructure. Tests should verify:
1. Log levels work correctly (debug, info, warn, error)
2. Structured logging produces correct JSON output
3. Context-aware logging passes through request IDs, trace IDs
4. CLI-friendly output is human-readable (not just JSON)
5. Streaming output works for long-running processes

---

## Phase Overview

The shoyu (soy sauce) package adds "flavor" and transparency to the bento system through structured logging. It wraps `zerolog` to provide:

- **Structured logging:** JSON output for machine parsing
- **CLI-friendly output:** Human-readable for terminal tailing
- **Context-aware:** Automatic trace IDs, request IDs
- **Performance:** Zero-allocation logging with zerolog
- **Streaming:** Real-time output for long Blender renders

### Why "Shoyu"?

Soy sauce (shoyu) adds flavor and depth—our logger adds transparency and observability. Just as you can't have good sushi without good shoyu, you can't have a debuggable CLI without good logging.

---

## Success Criteria

**Phase 2 Complete When:**
- [ ] Logger initialization with configurable levels
- [ ] Structured JSON output mode
- [ ] Human-readable console output mode
- [ ] Context-aware logging (trace IDs, neta IDs)
- [ ] Integration tests for all log levels
- [ ] Streaming output support
- [ ] Files < 250 lines each
- [ ] File-level documentation complete
- [ ] `/code-review` run with Karen + Colossus approval

---

## Test-First Approach

### Step 1: Define logger interface via tests

Create `pkg/shoyu/shoyu_test.go`:

```go
package shoyu_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/bento/pkg/shoyu"
)

// Test: Logger should output JSON in structured mode
func TestLogger_StructuredJSON(t *testing.T) {
	var buf bytes.Buffer

	logger := shoyu.New(shoyu.Config{
		Level:  shoyu.LevelInfo,
		Format: shoyu.FormatJSON,
		Output: &buf,
	})

	logger.Info().
		Str("neta_type", "http-request").
		Str("neta_id", "node-1").
		Msg("Executing HTTP request")

	// Parse JSON output
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if logEntry["level"] != "info" {
		t.Errorf("level = %v, want info", logEntry["level"])
	}

	if logEntry["neta_type"] != "http-request" {
		t.Errorf("neta_type = %v, want http-request", logEntry["neta_type"])
	}

	if logEntry["message"] != "Executing HTTP request" {
		t.Errorf("message = %v, want 'Executing HTTP request'", logEntry["message"])
	}
}

// Test: Logger should output human-readable text in console mode
func TestLogger_ConsoleOutput(t *testing.T) {
	var buf bytes.Buffer

	logger := shoyu.New(shoyu.Config{
		Level:  shoyu.LevelInfo,
		Format: shoyu.FormatConsole,
		Output: &buf,
	})

	logger.Info().
		Str("url", "https://api.example.com").
		Msg("Fetching data")

	output := buf.String()

	// Should be human-readable, not JSON
	if strings.Contains(output, "{") {
		t.Error("Console output should not contain JSON")
	}

	if !strings.Contains(output, "Fetching data") {
		t.Error("Console output should contain message")
	}

	if !strings.Contains(output, "https://api.example.com") {
		t.Error("Console output should contain URL")
	}
}

// Test: Log levels should filter correctly
func TestLogger_Levels(t *testing.T) {
	var buf bytes.Buffer

	// Set level to WARN - should not see INFO
	logger := shoyu.New(shoyu.Config{
		Level:  shoyu.LevelWarn,
		Format: shoyu.FormatJSON,
		Output: &buf,
	})

	logger.Info().Msg("This should not appear")
	logger.Warn().Msg("This should appear")

	output := buf.String()

	if strings.Contains(output, "This should not appear") {
		t.Error("Info message should be filtered out")
	}

	if !strings.Contains(output, "This should appear") {
		t.Error("Warn message should be present")
	}
}

// Test: Context should propagate through logger
func TestLogger_WithContext(t *testing.T) {
	var buf bytes.Buffer

	logger := shoyu.New(shoyu.Config{
		Level:  shoyu.LevelInfo,
		Format: shoyu.FormatJSON,
		Output: &buf,
	})

	// Create logger with context (like itamae would do)
	contextLogger := logger.With().
		Str("trace_id", "trace-123").
		Str("bento_id", "my-workflow").
		Logger()

	contextLogger.Info().Msg("Executing neta")

	// Parse JSON
	var logEntry map[string]interface{}
	json.Unmarshal(buf.Bytes(), &logEntry)

	if logEntry["trace_id"] != "trace-123" {
		t.Errorf("trace_id = %v, want trace-123", logEntry["trace_id"])
	}

	if logEntry["bento_id"] != "my-workflow" {
		t.Errorf("bento_id = %v, want my-workflow", logEntry["bento_id"])
	}
}
```

### Step 2: Test streaming output support

```go
// Test: Streaming callback should work for long-running processes
func TestLogger_StreamingCallback(t *testing.T) {
	var buf bytes.Buffer
	var streamLines []string

	logger := shoyu.New(shoyu.Config{
		Level:  shoyu.LevelInfo,
		Format: shoyu.FormatConsole,
		Output: &buf,
		OnStream: func(line string) {
			streamLines = append(streamLines, line)
		},
	})

	// Simulate Blender output streaming
	logger.Stream("Fra:1 Mem:12.00M (Peak 12.00M) | Rendering 1/100")
	logger.Stream("Fra:2 Mem:12.00M (Peak 12.00M) | Rendering 2/100")
	logger.Stream("Fra:3 Mem:12.00M (Peak 12.00M) | Rendering 3/100")

	if len(streamLines) != 3 {
		t.Errorf("Expected 3 stream lines, got %d", len(streamLines))
	}

	if !strings.Contains(streamLines[2], "Rendering 3/100") {
		t.Error("Stream lines should contain Blender output")
	}
}
```

---

## File Structure

```
pkg/shoyu/
├── shoyu.go           # Main logger implementation (~200 lines)
├── config.go          # Configuration types (~80 lines)
├── levels.go          # Log level constants (~50 lines)
├── context.go         # Context-aware logging (~100 lines)
├── streaming.go       # Streaming output support (~100 lines)
└── shoyu_test.go      # Integration tests (~300 lines)
```

---

## Implementation Guidance

**File: `pkg/shoyu/shoyu.go`**

```go
// Package shoyu provides structured logging for the bento system.
//
// "Shoyu" (醤油 - soy sauce) adds flavor and transparency—our logger provides
// observability into bento workflow execution.
//
// The logger wraps zerolog for high-performance, zero-allocation logging with
// both structured (JSON) and human-readable (console) output modes.
//
// Usage:
//
//	logger := shoyu.New(shoyu.Config{
//	    Level:  shoyu.LevelInfo,
//	    Format: shoyu.FormatConsole,
//	})
//
//	logger.Info().
//	    Str("neta_type", "http-request").
//	    Str("url", "https://api.example.com").
//	    Msg("Executing HTTP request")
//
// Learn more about zerolog: https://github.com/rs/zerolog
package shoyu

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

// Logger wraps zerolog.Logger with bento-specific functionality.
type Logger struct {
	zl       zerolog.Logger
	config   Config
	onStream func(string) // Callback for streaming output
}

// New creates a new Logger with the given configuration.
func New(cfg Config) *Logger {
	// Set default output to stdout if not specified
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}

	// Set default level to Info if not specified
	if cfg.Level == "" {
		cfg.Level = LevelInfo
	}

	// Create zerolog logger
	var zl zerolog.Logger

	switch cfg.Format {
	case FormatJSON:
		zl = zerolog.New(cfg.Output).With().Timestamp().Logger()
	case FormatConsole:
		zl = zerolog.New(zerolog.ConsoleWriter{Out: cfg.Output}).
			With().Timestamp().Logger()
	default:
		zl = zerolog.New(zerolog.ConsoleWriter{Out: cfg.Output}).
			With().Timestamp().Logger()
	}

	// Set log level
	switch cfg.Level {
	case LevelDebug:
		zl = zl.Level(zerolog.DebugLevel)
	case LevelInfo:
		zl = zl.Level(zerolog.InfoLevel)
	case LevelWarn:
		zl = zl.Level(zerolog.WarnLevel)
	case LevelError:
		zl = zl.Level(zerolog.ErrorLevel)
	}

	return &Logger{
		zl:       zl,
		config:   cfg,
		onStream: cfg.OnStream,
	}
}

// Info starts a new info-level log entry.
func (l *Logger) Info() *zerolog.Event {
	return l.zl.Info()
}

// Debug starts a new debug-level log entry.
func (l *Logger) Debug() *zerolog.Event {
	return l.zl.Debug()
}

// Warn starts a new warn-level log entry.
func (l *Logger) Warn() *zerolog.Event {
	return l.zl.Warn()
}

// Error starts a new error-level log entry.
func (l *Logger) Error() *zerolog.Event {
	return l.zl.Error()
}

// With creates a child logger with additional context fields.
// This is used by itamae to add trace IDs, bento IDs, etc.
func (l *Logger) With() zerolog.Context {
	return l.zl.With()
}

// Stream outputs a line for streaming processes (like Blender renders).
// This bypasses normal log levels for real-time output.
func (l *Logger) Stream(line string) {
	if l.onStream != nil {
		l.onStream(line)
	}

	// Also log at debug level
	l.zl.Debug().Str("stream", line).Msg("")
}
```

**File: `pkg/shoyu/config.go`**

```go
package shoyu

import "io"

// Config contains logger configuration.
type Config struct {
	// Level is the minimum log level to output.
	Level Level

	// Format is the output format (JSON or Console).
	Format Format

	// Output is where logs are written (default: os.Stdout).
	Output io.Writer

	// OnStream is called for streaming output (optional).
	// Used for real-time output from long-running processes.
	OnStream func(string)
}

// Level represents a log level.
type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Format represents the output format.
type Format string

const (
	FormatJSON    Format = "json"    // Structured JSON output
	FormatConsole Format = "console" // Human-readable console output
)
```

**File: `pkg/shoyu/context.go`**

```go
package shoyu

import "github.com/rs/zerolog"

// WithBentoID adds bento_id to the logger context.
// The itamae will use this to tag all logs from a specific workflow.
func WithBentoID(logger zerolog.Logger, bentoID string) zerolog.Logger {
	return logger.With().Str("bento_id", bentoID).Logger()
}

// WithNetaID adds neta_id to the logger context.
// Used to track which neta is currently executing.
func WithNetaID(logger zerolog.Logger, netaID string) zerolog.Logger {
	return logger.With().Str("neta_id", netaID).Logger()
}

// WithTraceID adds trace_id to the logger context.
// Used for distributed tracing (future feature).
func WithTraceID(logger zerolog.Logger, traceID string) zerolog.Logger {
	return logger.With().Str("trace_id", traceID).Logger()
}

// WithNetaType adds neta_type to the logger context.
func WithNetaType(logger zerolog.Logger, netaType string) zerolog.Logger {
	return logger.With().Str("neta_type", netaType).Logger()
}
```

**File: `pkg/shoyu/streaming.go`**

```go
package shoyu

import (
	"bufio"
	"io"
)

// StreamReader wraps an io.Reader and calls a callback for each line.
// Used by shell-command neta to stream Blender output in real-time.
//
// Example usage:
//
//	cmd := exec.Command("blender", args...)
//	stdout, _ := cmd.StdoutPipe()
//
//	go shoyu.StreamReader(stdout, logger, func(line string) {
//	    logger.Stream(line)
//	})
func StreamReader(r io.Reader, logger *Logger, callback func(string)) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		if callback != nil {
			callback(line)
		}
	}
}
```

---

## Common Go Pitfalls to Avoid

1. **Zerolog event must end with .Msg()**: Always call Msg() to emit the log
   ```go
   // ❌ BAD - log is silently dropped
   logger.Info().Str("key", "value")

   // ✅ GOOD
   logger.Info().Str("key", "value").Msg("Something happened")
   ```

2. **Context is immutable**: With() returns a NEW logger, doesn't modify the original
   ```go
   // ❌ BAD - doesn't actually add the field
   logger.With().Str("field", "value")
   logger.Info().Msg("Test")  // field is NOT present

   // ✅ GOOD
   contextLogger := logger.With().Str("field", "value").Logger()
   contextLogger.Info().Msg("Test")  // field IS present
   ```

---

## Critical for Phase 8

**Streaming Output:**
- The shell-command neta will use `StreamReader` to stream Blender output
- Must handle multi-minute processes without buffering everything in memory
- Callback should be called for EACH line as it's received

**Context Propagation:**
- Itamae will create context loggers with `WithBentoID` and `WithNetaID`
- Every log from a neta should include its ID for debugging

---

## Bento Box Principle Checklist

- [ ] File < 250 lines (shoyu.go ~200, config.go ~80, etc.)
- [ ] Functions < 20 lines
- [ ] Single responsibility (logging only, no business logic)
- [ ] No utility grab bags
- [ ] Clear interface (Logger wraps zerolog)
- [ ] File-level documentation

---

## Phase Completion

**Phase 2 MUST end with:**

1. All tests passing (`go test ./pkg/shoyu/...`)
2. Run `/code-review` slash command
3. Address feedback from Karen and Colossus
4. Get explicit approval from both agents
5. Document any decisions in `.claude/strategy/`

**Do not proceed to Phase 3 until code review is approved.**

---

## Claude Prompt Template

```
I need to implement Phase 2: shoyu (logger package) following TDD principles.

Please read:
- .claude/strategy/phase-2-shoyu.md (this file)
- .claude/BENTO_BOX_PRINCIPLE.md

Then:

1. Create `pkg/shoyu/shoyu_test.go` with integration tests for:
   - JSON structured output
   - Console human-readable output
   - Log level filtering
   - Context propagation (trace IDs, bento IDs, neta IDs)
   - Streaming output support

2. Watch the tests fail

3. Implement the following files to make tests pass:
   - pkg/shoyu/shoyu.go (~200 lines)
   - pkg/shoyu/config.go (~80 lines)
   - pkg/shoyu/context.go (~100 lines)
   - pkg/shoyu/streaming.go (~100 lines)

4. Add file-level documentation explaining:
   - What zerolog is and why we use it
   - How to use the logger
   - Common pitfalls (event.Msg() required, context is immutable)

Remember:
- Write tests FIRST
- Keep files < 250 lines
- Keep functions < 20 lines
- Integration tests, not unit tests
- CRITICAL: Streaming support for Phase 8 (Blender renders)

When complete, run `/code-review` and get Karen + Colossus approval.
```

---

## Dependencies to Add

```bash
go get github.com/rs/zerolog
```

---

## Notes

- Shoyu is infrastructure - get it right, it's used everywhere
- Streaming output is CRITICAL for Phase 8 (Blender takes minutes to render)
- Context propagation helps with debugging complex workflows
- Zero-allocation logging means minimal performance impact
- Both JSON (for machines) and console (for humans) modes are needed

---

**Status:** Ready for implementation
**Next Phase:** Phase 3 (omakase validation) - depends on completion of Phase 2
