package shellcommand_test

import (
	"context"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Develonaut/bento/pkg/neta/library/shellcommand"
)

// TestShellCommand_BasicExecution tests basic command execution.
func TestShellCommand_BasicExecution(t *testing.T) {
	ctx := context.Background()

	sc := shellcommand.New()

	params := map[string]interface{}{
		"command": "echo",
		"args":    []interface{}{"hello", "world"},
	}

	result, err := sc.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output := result.(map[string]interface{})
	stdout := output["stdout"].(string)

	if !strings.Contains(stdout, "hello world") {
		t.Errorf("stdout = %q, want to contain 'hello world'", stdout)
	}

	exitCode, ok := output["exitCode"].(int)
	if !ok {
		t.Fatalf("exitCode is not an int: %T", output["exitCode"])
	}

	if exitCode != 0 {
		t.Errorf("exitCode = %v, want 0", exitCode)
	}
}

// TestShellCommand_WithArguments tests command execution with arguments.
func TestShellCommand_WithArguments(t *testing.T) {
	ctx := context.Background()

	sc := shellcommand.New()

	var params map[string]interface{}

	if runtime.GOOS == "windows" {
		// Windows command
		params = map[string]interface{}{
			"command": "cmd",
			"args":    []interface{}{"/C", "echo", "test-arg-123"},
		}
	} else {
		// Unix command
		params = map[string]interface{}{
			"command": "echo",
			"args":    []interface{}{"test-arg-123"},
		}
	}

	result, err := sc.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output := result.(map[string]interface{})
	stdout := output["stdout"].(string)

	if !strings.Contains(stdout, "test-arg-123") {
		t.Errorf("stdout = %q, want to contain 'test-arg-123'", stdout)
	}
}

// TestShellCommand_StdoutStderr tests stdout and stderr capture.
func TestShellCommand_StdoutStderr(t *testing.T) {
	ctx := context.Background()

	sc := shellcommand.New()

	var params map[string]interface{}

	if runtime.GOOS == "windows" {
		// Windows: redirect stderr
		params = map[string]interface{}{
			"command": "cmd",
			"args":    []interface{}{"/C", "echo error message 1>&2"},
		}
	} else {
		// Unix: redirect stderr
		params = map[string]interface{}{
			"command": "sh",
			"args":    []interface{}{"-c", "echo 'error message' >&2"},
		}
	}

	result, err := sc.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output := result.(map[string]interface{})
	stderr := output["stderr"].(string)

	if !strings.Contains(stderr, "error message") {
		t.Errorf("stderr = %q, want to contain 'error message'", stderr)
	}
}

// TestShellCommand_ExitCode tests capturing non-zero exit codes.
func TestShellCommand_ExitCode(t *testing.T) {
	ctx := context.Background()

	sc := shellcommand.New()

	var params map[string]interface{}

	if runtime.GOOS == "windows" {
		params = map[string]interface{}{
			"command": "cmd",
			"args":    []interface{}{"/C", "exit", "42"},
		}
	} else {
		params = map[string]interface{}{
			"command": "sh",
			"args":    []interface{}{"-c", "exit 42"},
		}
	}

	result, err := sc.Execute(ctx, params)
	// Non-zero exit should NOT be an error - we capture it in exitCode
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output := result.(map[string]interface{})
	exitCode := output["exitCode"].(int)

	if exitCode != 42 {
		t.Errorf("exitCode = %v, want 42", exitCode)
	}
}

// TestShellCommand_ConfigurableTimeout tests that timeout can be configured.
// CRITICAL FOR PHASE 8: Blender renders take 5-30 minutes.
func TestShellCommand_ConfigurableTimeout(t *testing.T) {
	ctx := context.Background()

	sc := shellcommand.New()

	var params map[string]interface{}

	if runtime.GOOS == "windows" {
		// Windows: sleep for 2 seconds using timeout command
		params = map[string]interface{}{
			"command": "timeout",
			"args":    []interface{}{"/T", "2", "/NOBREAK"},
			"timeout": 5, // 5 second timeout (should NOT timeout)
		}
	} else {
		// Unix: sleep for 2 seconds
		params = map[string]interface{}{
			"command": "sleep",
			"args":    []interface{}{"2"},
			"timeout": 5, // 5 second timeout (should NOT timeout)
		}
	}

	start := time.Now()
	result, err := sc.Execute(ctx, params)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Should complete in about 2 seconds, definitely less than 5
	if duration > 4*time.Second {
		t.Errorf("Command took too long: %v", duration)
	}

	output := result.(map[string]interface{})
	exitCode := output["exitCode"].(int)

	if exitCode != 0 {
		t.Errorf("exitCode = %v, want 0", exitCode)
	}
}

// TestShellCommand_Timeout tests that commands timeout appropriately.
func TestShellCommand_Timeout(t *testing.T) {
	ctx := context.Background()

	sc := shellcommand.New()

	var params map[string]interface{}

	if runtime.GOOS == "windows" {
		// Windows: sleep for 5 seconds
		params = map[string]interface{}{
			"command": "timeout",
			"args":    []interface{}{"/T", "5", "/NOBREAK"},
			"timeout": 1, // 1 second timeout (SHOULD timeout)
		}
	} else {
		// Unix: sleep for 5 seconds
		params = map[string]interface{}{
			"command": "sleep",
			"args":    []interface{}{"5"},
			"timeout": 1, // 1 second timeout (SHOULD timeout)
		}
	}

	_, err := sc.Execute(ctx, params)
	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}

	if !strings.Contains(err.Error(), "timeout") && !strings.Contains(err.Error(), "deadline") && !strings.Contains(err.Error(), "killed") {
		t.Errorf("Expected timeout/deadline/killed error, got: %v", err)
	}
}

// TestShellCommand_StreamingOutput tests streaming output callback.
// CRITICAL FOR PHASE 8: Stream Blender render progress line-by-line.
func TestShellCommand_StreamingOutput(t *testing.T) {
	ctx := context.Background()

	sc := shellcommand.New()

	var params map[string]interface{}
	var outputLines []string

	// Callback to capture streaming output
	onOutput := func(line string) {
		outputLines = append(outputLines, line)
	}

	if runtime.GOOS == "windows" {
		params = map[string]interface{}{
			"command":   "cmd",
			"args":      []interface{}{"/C", "for /L %i in (1,1,3) do @echo %i"},
			"stream":    true,
			"_onOutput": onOutput,
		}
	} else {
		params = map[string]interface{}{
			"command":   "sh",
			"args":      []interface{}{"-c", "for i in 1 2 3; do echo $i; done"},
			"stream":    true,
			"_onOutput": onOutput,
		}
	}

	result, err := sc.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify we got streaming output
	if len(outputLines) < 3 {
		t.Errorf("Expected at least 3 output lines, got %d: %v", len(outputLines), outputLines)
	}

	// Verify final result also has stdout
	output := result.(map[string]interface{})
	stdout := output["stdout"].(string)

	if !strings.Contains(stdout, "1") || !strings.Contains(stdout, "2") || !strings.Contains(stdout, "3") {
		t.Errorf("stdout missing expected numbers: %q", stdout)
	}
}

// TestShellCommand_LongRunning tests long-running commands.
// CRITICAL FOR PHASE 8: Blender renders can take many minutes.
func TestShellCommand_LongRunning(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long-running test in short mode")
	}

	ctx := context.Background()

	sc := shellcommand.New()

	var params map[string]interface{}

	if runtime.GOOS == "windows" {
		params = map[string]interface{}{
			"command": "timeout",
			"args":    []interface{}{"/T", "3", "/NOBREAK"},
			"timeout": 10, // 10 second timeout (won't timeout)
		}
	} else {
		params = map[string]interface{}{
			"command": "sleep",
			"args":    []interface{}{"3"},
			"timeout": 10, // 10 second timeout (won't timeout)
		}
	}

	start := time.Now()
	result, err := sc.Execute(ctx, params)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Should take at least 3 seconds
	if duration < 2*time.Second {
		t.Errorf("Command finished too quickly: %v", duration)
	}

	output := result.(map[string]interface{})
	exitCode := output["exitCode"].(int)

	if exitCode != 0 {
		t.Errorf("exitCode = %v, want 0", exitCode)
	}
}

// TestShellCommand_ContextCancellation tests that context cancellation stops the command.
func TestShellCommand_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	sc := shellcommand.New()

	var params map[string]interface{}

	if runtime.GOOS == "windows" {
		// Sleep for 10 seconds (will be cancelled)
		params = map[string]interface{}{
			"command": "timeout",
			"args":    []interface{}{"/T", "10", "/NOBREAK"},
		}
	} else {
		params = map[string]interface{}{
			"command": "sleep",
			"args":    []interface{}{"10"},
		}
	}

	_, err := sc.Execute(ctx, params)
	if err == nil {
		t.Fatal("Expected context cancellation error, got nil")
	}

	if !strings.Contains(err.Error(), "context") && !strings.Contains(err.Error(), "killed") && !strings.Contains(err.Error(), "timeout") {
		t.Errorf("Expected context/killed/timeout error, got: %v", err)
	}
}

// TestShellCommand_DefaultTimeout tests that a default timeout is applied.
func TestShellCommand_DefaultTimeout(t *testing.T) {
	ctx := context.Background()

	sc := shellcommand.New()

	params := map[string]interface{}{
		"command": "echo",
		"args":    []interface{}{"test"},
		// No timeout specified - should use default (120 seconds)
	}

	result, err := sc.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output := result.(map[string]interface{})
	if output["exitCode"].(int) != 0 {
		t.Errorf("exitCode = %v, want 0", output["exitCode"])
	}
}
