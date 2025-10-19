package omakase

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Develonaut/bento/pkg/neta"
)

// preflightShellCommand checks if the command exists in PATH.
func preflightShellCommand(def *neta.Definition) error {
	command, ok := def.Parameters["command"].(string)
	if !ok {
		return nil // Already validated by validateShellCommand
	}

	// Check if command exists in PATH
	_, err := exec.LookPath(command)
	if err != nil {
		return fmt.Errorf("shell-command neta '%s': command '%s' not found in PATH. Please install it first",
			def.ID, command)
	}

	return nil
}

// preflightHTTPRequest checks for required environment variables in URL/headers.
func preflightHTTPRequest(def *neta.Definition) error {
	if err := checkURLEnvVars(def); err != nil {
		return err
	}
	return checkHeaderEnvVars(def)
}

// checkURLEnvVars validates environment variables in URL.
func checkURLEnvVars(def *neta.Definition) error {
	url, _ := def.Parameters["url"].(string)
	envVars := extractEnvVars(url)

	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			return fmt.Errorf("http-request neta '%s': environment variable '%s' not set (required in URL)",
				def.ID, envVar)
		}
	}
	return nil
}

// checkHeaderEnvVars validates environment variables in headers.
func checkHeaderEnvVars(def *neta.Definition) error {
	headers, ok := def.Parameters["headers"].(map[string]string)
	if !ok {
		return nil
	}

	for key, value := range headers {
		if err := checkHeaderValue(def, key, value); err != nil {
			return err
		}
	}
	return nil
}

// checkHeaderValue validates environment variables in a single header value.
func checkHeaderValue(def *neta.Definition, key, value string) error {
	envVars := extractEnvVars(value)
	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			return fmt.Errorf("http-request neta '%s': environment variable '%s' not set (required in header '%s')",
				def.ID, envVar, key)
		}
	}
	return nil
}

// preflightFileSystem checks if file paths exist for read operations.
func preflightFileSystem(def *neta.Definition) error {
	operation, _ := def.Parameters["operation"].(string)

	// For read operations, check file exists
	if operation == "read" {
		if path, ok := def.Parameters["path"].(string); ok {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return fmt.Errorf("file-system neta '%s': file not found: %s", def.ID, path)
			}
		}
	}

	return nil
}

// extractEnvVars finds {{.VAR_NAME}} patterns in a string.
//
// This is a simple string-based approach that doesn't require regex.
// It extracts all environment variable names from Go template syntax.
func extractEnvVars(s string) []string {
	var vars []string

	for {
		varName, remaining, found := findNextVar(s)
		if !found {
			break
		}
		vars = append(vars, varName)
		s = remaining
	}

	return vars
}

// findNextVar finds the next {{.VAR}} pattern and returns (varName, remaining, found).
func findNextVar(s string) (string, string, bool) {
	start := strings.Index(s, "{{.")
	if start == -1 {
		return "", "", false
	}

	end := strings.Index(s[start:], "}}")
	if end == -1 {
		return "", "", false
	}

	varName := s[start+3 : start+end]
	remaining := s[start+end+2:]
	return varName, remaining, true
}
