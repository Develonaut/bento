// Package shellcommand provides shell command execution for the bento workflow system.
//
// The shellcommand neta allows you to execute shell commands and capture their output.
// This is CRITICAL for Phase 8 (Blender automation) which requires:
//   - Long-running commands (5-30 minute Blender renders)
//   - Streaming output (line-by-line render progress)
//   - Configurable timeouts
//   - Exit code capture
//
// Example usage:
//
//	// Basic command execution
//	params := map[string]interface{}{
//	    "command": "ls",
//	    "args": []string{"-la"},
//	}
//
//	// Long-running Blender render with streaming output
//	params := map[string]interface{}{
//	    "command": "blender",
//	    "args": []string{
//	        "--background", "scene.blend",
//	        "--render-output", "/tmp/frame_####",
//	        "--frame-start", "1",
//	        "--frame-end", "10",
//	        "--render-anim",
//	    },
//	    "timeout": 1800,  // 30 minutes in seconds
//	    "stream": true,   // Enable line-by-line output streaming
//	}
//
// The result contains:
//   - stdout: Command standard output
//   - stderr: Command standard error
//   - exitCode: Process exit code
//
// Learn more about Go's os/exec package: https://pkg.go.dev/os/exec
package shellcommand

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/Develonaut/bento/pkg/neta"
)

const (
	// DefaultTimeout is the default command timeout in seconds.
	// Set to 2 minutes for typical commands, but can be overridden
	// for long-running operations like Blender renders (30+ minutes).
	DefaultTimeout = 120
)

// ShellCommandNeta implements shell command execution.
type ShellCommandNeta struct{}

// New creates a new shellcommand neta instance.
func New() neta.Executable {
	return &ShellCommandNeta{}
}

// Execute runs a shell command based on the provided parameters.
//
// Parameters:
//   - command (string, required): The command to execute
//   - args ([]interface{}, optional): Command arguments
//   - timeout (int, optional): Timeout in seconds (default: 120)
//   - stream (bool, optional): Enable line-by-line output streaming
//   - _onOutput (func(string), optional): Callback for streaming output
//
// Returns a map with:
//   - stdout (string): Standard output
//   - stderr (string): Standard error
//   - exitCode (int): Exit code (0 = success)
func (s *ShellCommandNeta) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// Extract command
	command, ok := params["command"].(string)
	if !ok {
		return nil, fmt.Errorf("command parameter is required and must be a string")
	}

	// Extract args (optional)
	var args []string
	if argsRaw, ok := params["args"].([]interface{}); ok {
		args = make([]string, len(argsRaw))
		for i, arg := range argsRaw {
			if strArg, ok := arg.(string); ok {
				args[i] = strArg
			} else {
				return nil, fmt.Errorf("all args must be strings, got %T at index %d", arg, i)
			}
		}
	}

	// Extract timeout (optional, default 120 seconds)
	timeout := DefaultTimeout
	if t, ok := params["timeout"].(int); ok {
		timeout = t
	}

	// Extract streaming settings
	stream, _ := params["stream"].(bool)
	var onOutput func(string)
	if callback, ok := params["_onOutput"].(func(string)); ok {
		onOutput = callback
	}

	// Create context with timeout
	cmdCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	// Create command
	cmd := exec.CommandContext(cmdCtx, command, args...)

	// Capture stdout and stderr
	var stdoutBuilder strings.Builder
	var stderrBuilder strings.Builder

	if stream && onOutput != nil {
		// Streaming mode: read line-by-line and call callback
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
		}

		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
		}

		// Start the command
		if err := cmd.Start(); err != nil {
			return nil, fmt.Errorf("failed to start command: %w", err)
		}

		// Read stdout line-by-line
		stdoutScanner := bufio.NewScanner(stdoutPipe)
		go func() {
			for stdoutScanner.Scan() {
				line := stdoutScanner.Text()
				stdoutBuilder.WriteString(line)
				stdoutBuilder.WriteString("\n")
				onOutput(line)
			}
		}()

		// Read stderr line-by-line
		stderrScanner := bufio.NewScanner(stderrPipe)
		go func() {
			for stderrScanner.Scan() {
				line := stderrScanner.Text()
				stderrBuilder.WriteString(line)
				stderrBuilder.WriteString("\n")
			}
		}()

		// Wait for command to finish
		err = cmd.Wait()

		// Check context errors first
		if cmdCtx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("command timeout after %d seconds", timeout)
		}
		if cmdCtx.Err() == context.Canceled {
			return nil, fmt.Errorf("command killed due to context cancellation")
		}

		exitCode := 0
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				return nil, fmt.Errorf("command failed: %w", err)
			}
		}

		return map[string]interface{}{
			"stdout":   stdoutBuilder.String(),
			"stderr":   stderrBuilder.String(),
			"exitCode": exitCode,
		}, nil
	} else {
		// Non-streaming mode: capture all output
		cmd.Stdout = &stdoutBuilder
		cmd.Stderr = &stderrBuilder

		err := cmd.Run()

		// Check context errors first
		if cmdCtx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("command timeout after %d seconds", timeout)
		}
		if cmdCtx.Err() == context.Canceled {
			return nil, fmt.Errorf("command killed due to context cancellation")
		}

		exitCode := 0
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				return nil, fmt.Errorf("command failed: %w", err)
			}
		}

		return map[string]interface{}{
			"stdout":   stdoutBuilder.String(),
			"stderr":   stderrBuilder.String(),
			"exitCode": exitCode,
		}, nil
	}
}
