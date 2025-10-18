// Package shell provides shell command execution nodes.
package shell

import (
	"context"
	"fmt"
	"os/exec"

	"bento/pkg/neta"
)

// Command executes shell commands.
type Command struct{}

// NewCommand creates a new shell command node.
func NewCommand() *Command {
	return &Command{}
}

// Execute runs a shell command and returns output.
func (c *Command) Execute(ctx context.Context, params map[string]interface{}) (neta.Result, error) {
	command := neta.GetStringParam(params, "command", "")
	if command == "" {
		return neta.Result{}, fmt.Errorf("command parameter required")
	}

	// Get args (default: empty slice)
	args := make([]string, 0)
	if argsInterface, ok := params["args"]; ok {
		if argsList, ok := argsInterface.([]interface{}); ok {
			for _, arg := range argsList {
				if strArg, ok := arg.(string); ok {
					args = append(args, strArg)
				} else {
					args = append(args, fmt.Sprintf("%v", arg))
				}
			}
		}
	}

	// Get working directory (default: "")
	workDir := neta.GetStringParam(params, "working_dir", "")

	// Execute command
	cmd := exec.CommandContext(ctx, command, args...)

	if workDir != "" {
		cmd.Dir = workDir
	}

	// Capture combined output (stdout + stderr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return neta.Result{}, fmt.Errorf("command failed: %w\nOutput: %s", err, string(output))
	}

	result := map[string]interface{}{
		"stdout":    string(output),
		"exit_code": 0,
	}

	return neta.Result{Output: result}, nil
}
