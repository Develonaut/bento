// Package shoyu provides structured logging for the bento system.
//
// Tests for charm/log integration.
package shoyu

import (
	"bytes"
	"strings"
	"testing"
)

// TestLogger_CharmEnabled verifies charm/log is used when UseCharm is true.
func TestLogger_CharmEnabled(t *testing.T) {
	var buf bytes.Buffer
	logger := New(Config{
		Level:    LevelInfo,
		Output:   &buf,
		UseCharm: true,
	})

	logger.Info("Test message", "key", "value")

	output := buf.String()
	if output == "" {
		t.Error("Charm logger produced no output")
	}

	// Charm/log output should be more colorful/formatted than slog
	// We can't test exact format due to ANSI codes, but we can check for content
	if !strings.Contains(output, "Test message") {
		t.Errorf("Output missing message: %s", output)
	}
}

// TestLogger_CharmDisabled verifies slog is used when UseCharm is false.
func TestLogger_CharmDisabled(t *testing.T) {
	var buf bytes.Buffer
	logger := New(Config{
		Level:    LevelInfo,
		Output:   &buf,
		UseCharm: false,
	})

	logger.Info("Test message", "key", "value")

	output := buf.String()
	if output == "" {
		t.Error("Slog logger produced no output")
	}

	if !strings.Contains(output, "Test message") {
		t.Errorf("Output missing message: %s", output)
	}
}

// TestLogger_CharmDefaultIsFalse verifies UseCharm defaults to false.
func TestLogger_CharmDefaultIsFalse(t *testing.T) {
	var buf bytes.Buffer
	logger := New(Config{
		Level:  LevelInfo,
		Output: &buf,
		// UseCharm not specified, should default to false
	})

	logger.Info("Test message")

	output := buf.String()
	if output == "" {
		t.Error("Logger produced no output")
	}

	// Should use slog format (text handler)
	if !strings.Contains(output, "level=INFO") {
		t.Errorf("Expected slog format, got: %s", output)
	}
}

// TestLogger_CharmLevels verifies charm/log respects log levels.
func TestLogger_CharmLevels(t *testing.T) {
	var buf bytes.Buffer
	logger := New(Config{
		Level:    LevelWarn,
		Output:   &buf,
		UseCharm: true,
	})

	// Info should be filtered out
	logger.Info("Should not appear")
	if buf.String() != "" {
		t.Errorf("Info message should be filtered: %s", buf.String())
	}

	// Warn should appear
	buf.Reset()
	logger.Warn("Should appear")
	if buf.String() == "" {
		t.Error("Warn message should appear")
	}
}

// TestLogger_CharmWith verifies With() creates child loggers.
func TestLogger_CharmWith(t *testing.T) {
	var buf bytes.Buffer
	logger := New(Config{
		Level:    LevelInfo,
		Output:   &buf,
		UseCharm: true,
	})

	childLogger := logger.With("parent_id", "123")
	childLogger.Info("Child message")

	output := buf.String()
	if !strings.Contains(output, "Child message") {
		t.Errorf("Output missing child message: %s", output)
	}

	// Charm log should include the context field
	if !strings.Contains(output, "parent_id") || !strings.Contains(output, "123") {
		t.Errorf("Output missing context fields: %s", output)
	}
}

// TestLogger_CharmStream verifies Stream() works with charm/log.
func TestLogger_CharmStream(t *testing.T) {
	var buf bytes.Buffer
	streamCalled := false
	streamLine := ""

	logger := New(Config{
		Level:    LevelDebug,
		Output:   &buf,
		UseCharm: true,
		OnStream: func(line string) {
			streamCalled = true
			streamLine = line
		},
	})

	logger.Stream("streaming output")

	if !streamCalled {
		t.Error("OnStream callback was not called")
	}

	if streamLine != "streaming output" {
		t.Errorf("Stream line = %q, want %q", streamLine, "streaming output")
	}

	// Should also be logged at debug level
	output := buf.String()
	if !strings.Contains(output, "stream") {
		t.Errorf("Stream not logged: %s", output)
	}
}
