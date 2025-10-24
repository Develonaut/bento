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
	"io"
	"os/exec"
	"strconv"
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

// commandParams holds extracted and validated command parameters.
type commandParams struct {
	command  string
	args     []string
	timeout  int
	stream   bool
	onOutput func(string)
}

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
	cmdParams, err := s.extractCommandParams(params)
	if err != nil {
		return nil, err
	}

	cmdCtx, cancel := context.WithTimeout(ctx, time.Duration(cmdParams.timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, cmdParams.command, cmdParams.args...)

	if cmdParams.stream && cmdParams.onOutput != nil {
		return s.executeStreaming(cmdCtx, cmd, cmdParams)
	}
	return s.executeBuffered(cmdCtx, cmd, cmdParams.timeout)
}

// extractCommandParams extracts and validates command parameters from the params map.
func (s *ShellCommandNeta) extractCommandParams(params map[string]interface{}) (*commandParams, error) {
	command, ok := params["command"].(string)
	if !ok {
		return nil, fmt.Errorf("command parameter is required and must be a string")
	}

	args, err := s.extractArgs(params)
	if err != nil {
		return nil, err
	}

	timeout := s.extractTimeout(params)
	stream, _ := params["stream"].(bool)

	var onOutput func(string)
	if callback, ok := params["_onOutput"].(func(string)); ok {
		onOutput = callback
	}

	return &commandParams{
		command:  command,
		args:     args,
		timeout:  timeout,
		stream:   stream,
		onOutput: onOutput,
	}, nil
}

// extractArgs extracts and validates command arguments.
func (s *ShellCommandNeta) extractArgs(params map[string]interface{}) ([]string, error) {
	argsRaw, ok := params["args"].([]interface{})
	if !ok {
		return nil, nil
	}

	args := make([]string, len(argsRaw))
	for i, arg := range argsRaw {
		strArg, ok := arg.(string)
		if !ok {
			return nil, fmt.Errorf("all args must be strings, got %T at index %d", arg, i)
		}
		args[i] = strArg
	}
	return args, nil
}

// extractTimeout extracts timeout value, handling int, float64, and string from templates.
func (s *ShellCommandNeta) extractTimeout(params map[string]interface{}) int {
	if t, ok := params["timeout"].(int); ok {
		return t
	}
	if t, ok := params["timeout"].(float64); ok {
		return int(t)
	}
	// Handle string values from template resolution
	if t, ok := params["timeout"].(string); ok {
		if timeout, err := strconv.Atoi(t); err == nil {
			return timeout
		}
	}
	return DefaultTimeout
}

// executeStreaming runs a command with streaming output.
func (s *ShellCommandNeta) executeStreaming(cmdCtx context.Context, cmd *exec.Cmd, params *commandParams) (interface{}, error) {
	var stdoutBuilder, stderrBuilder strings.Builder

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	s.streamOutput(stdoutPipe, &stdoutBuilder, params.onOutput)
	s.streamOutput(stderrPipe, &stderrBuilder, nil)

	err = cmd.Wait()
	return s.handleCommandResult(cmdCtx, err, &stdoutBuilder, &stderrBuilder, params.timeout)
}

// streamOutput reads from a pipe line-by-line and optionally calls a callback.
func (s *ShellCommandNeta) streamOutput(pipe io.ReadCloser, builder *strings.Builder, callback func(string)) {
	scanner := bufio.NewScanner(pipe)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			builder.WriteString(line)
			builder.WriteString("\n")
			if callback != nil {
				callback(line)
			}
		}
	}()
}

// executeBuffered runs a command with buffered output.
func (s *ShellCommandNeta) executeBuffered(cmdCtx context.Context, cmd *exec.Cmd, timeout int) (interface{}, error) {
	var stdoutBuilder, stderrBuilder strings.Builder

	cmd.Stdout = &stdoutBuilder
	cmd.Stderr = &stderrBuilder

	err := cmd.Run()
	return s.handleCommandResult(cmdCtx, err, &stdoutBuilder, &stderrBuilder, timeout)
}

// handleCommandResult processes command execution results and errors.
func (s *ShellCommandNeta) handleCommandResult(cmdCtx context.Context, err error, stdout, stderr *strings.Builder, timeout int) (interface{}, error) {
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
		"stdout":   stdout.String(),
		"stderr":   stderr.String(),
		"exitCode": exitCode,
	}, nil
}
