// Package main provides integration tests for the bento CLI.
//
// These tests verify the CLI commands work correctly end-to-end by:
//   - Creating test bento files on disk
//   - Executing the compiled bento binary
//   - Verifying output and exit codes
//
// Tests use the compiled binary to ensure we're testing the actual user experience.
package main_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Test helpers

// createTestBento creates a test bento file with the given content.
func createTestBento(t *testing.T, name, content string) string {
	t.Helper()

	tmpfile, err := os.CreateTemp("", name)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tmpfile.Name()
}

// createTestBentoInDir creates a test bento in a specific directory.
func createTestBentoInDir(t *testing.T, dir, name, title string) {
	t.Helper()

	content := simpleValidBento(title)
	path := filepath.Join(dir, name)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write bento file: %v", err)
	}
}

// bentoParts returns the parts of a bento definition.
func bentoParts() (header, node, footer string) {
	header = `{
		"id": "test-bento",
		"type": "group",
		"version": "1.0.0",
		"name": "`

	node = `",
		"position": {"x": 0, "y": 0},
		"metadata": {},
		"parameters": {},
		"inputPorts": [],
		"outputPorts": [],
		"nodes": [
			{
				"id": "node-1",
				"type": "edit-fields",
				"version": "1.0.0",
				"name": "Set Field",
				"position": {"x": 0, "y": 0},
				"metadata": {},
				"parameters": {
					"values": {"foo": "bar"}
				},
				"inputPorts": [],
				"outputPorts": []
			}
		],
		"edges": []
	}`

	return
}

// simpleValidBento returns a valid minimal bento definition.
func simpleValidBento(title string) string {
	if title == "" {
		title = "Test Bento"
	}
	h, n, _ := bentoParts()
	return h + title + n
}

// Test: bento eat command

// verifyCommandSuccess checks if command succeeded with expected output.
func verifyCommandSuccess(t *testing.T, cmd *exec.Cmd, expectedText string) string {
	t.Helper()
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput: %s", err, string(output))
	}
	outputStr := string(output)
	if !strings.Contains(outputStr, expectedText) {
		t.Errorf("Output should contain '%s': %s", expectedText, outputStr)
	}
	return outputStr
}

// TestEatCommand_ValidBento verifies that eat executes a valid bento successfully.
func TestEatCommand_ValidBento(t *testing.T) {
	bentoFile := createTestBento(t, "test.bento.json", simpleValidBento(""))
	defer os.Remove(bentoFile)

	cmd := exec.Command("bento", "eat", bentoFile)
	output := verifyCommandSuccess(t, cmd, "Delicious")

	if !strings.Contains(output, "🍱") {
		t.Error("Output should contain bento emoji 🍱")
	}
}

// verifyCommandError checks if command failed with expected error.
func verifyCommandError(t *testing.T, cmd *exec.Cmd) {
	t.Helper()
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("command should fail")
	}
	if exitError, ok := err.(*exec.ExitError); ok {
		if exitError.ExitCode() != 1 {
			t.Errorf("Exit code = %d, want 1", exitError.ExitCode())
		}
	}
	outputStr := strings.ToLower(string(output))
	if !strings.Contains(outputStr, "error") &&
		!strings.Contains(outputStr, "failed") &&
		!strings.Contains(outputStr, "missing") {
		t.Errorf("Output should mention error: %s", string(output))
	}
}

// TestEatCommand_InvalidBento verifies proper error handling for invalid bentos.
func TestEatCommand_InvalidBento(t *testing.T) {
	bentoFile := createTestBento(t, "invalid.bento.json", `{
		"id": "invalid",
		"type": "group",
		"name": "Invalid Bento"
	}`)
	defer os.Remove(bentoFile)

	verifyCommandError(t, exec.Command("bento", "eat", bentoFile))
}

// TestEatCommand_VerboseFlag verifies verbose output includes details.
func TestEatCommand_VerboseFlag(t *testing.T) {
	bentoFile := createTestBento(t, "test.bento.json", simpleValidBento(""))
	defer os.Remove(bentoFile)

	output := verifyCommandSuccess(t, exec.Command("bento", "eat", bentoFile, "--verbose"), "Delicious")
	if !strings.Contains(output, "node-1") {
		t.Error("Verbose output should mention node IDs")
	}
}

// Test: bento peek command

// TestPeekCommand_ValidBento verifies peek validates without executing.
func TestPeekCommand_ValidBento(t *testing.T) {
	bentoFile := createTestBento(t, "test.bento.json", simpleValidBento(""))
	defer os.Remove(bentoFile)

	output := verifyCommandSuccess(t, exec.Command("bento", "peek", bentoFile), "Looks delicious")
	if strings.Contains(output, "Delicious! Bento devoured") {
		t.Error("peek should NOT execute the bento")
	}
}

// invalidHTTPBento returns a bento with missing URL parameter.
func invalidHTTPBento() string {
	return `{
		"id": "invalid",
		"type": "group",
		"version": "1.0.0",
		"name": "Invalid",
		"position": {"x": 0, "y": 0},
		"metadata": {},
		"parameters": {},
		"inputPorts": [],
		"outputPorts": [],
		"nodes": [{
			"id": "http-1",
			"type": "http-request",
			"version": "1.0.0",
			"name": "HTTP",
			"position": {"x": 0, "y": 0},
			"metadata": {},
			"parameters": {"method": "GET"},
			"inputPorts": [],
			"outputPorts": []
		}],
		"edges": []
	}`
}

// TestPeekCommand_InvalidBento verifies peek reports validation errors clearly.
func TestPeekCommand_InvalidBento(t *testing.T) {
	bentoFile := createTestBento(t, "invalid.bento.json", invalidHTTPBento())
	defer os.Remove(bentoFile)

	verifyCommandError(t, exec.Command("bento", "peek", bentoFile))
}

// Test: bento menu command

// verifyBentosListed checks if bentos are listed in output.
func verifyBentosListed(t *testing.T, output string, files ...string) {
	t.Helper()
	for _, file := range files {
		if !strings.Contains(output, file) {
			t.Errorf("Output should list %s", file)
		}
	}
}

// TestMenuCommand_ListBentos verifies menu lists all bentos in directory.
func TestMenuCommand_ListBentos(t *testing.T) {
	tmpDir := t.TempDir()
	createTestBentoInDir(t, tmpDir, "workflow1.bento.json", "Workflow 1")
	createTestBentoInDir(t, tmpDir, "workflow2.bento.json", "Workflow 2")

	output := verifyCommandSuccess(t, exec.Command("bento", "menu", tmpDir), "2 bentos")
	verifyBentosListed(t, output, "workflow1.bento.json", "workflow2.bento.json")
}

// TestMenuCommand_EmptyDirectory verifies menu handles empty directory gracefully.
func TestMenuCommand_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	output := verifyCommandSuccess(t, exec.Command("bento", "menu", tmpDir), "No bentos")
	if !strings.Contains(output, "0 bentos") && !strings.Contains(output, "No bentos") {
		t.Errorf("Output should indicate no bentos found: %s", output)
	}
}

// Test: bento box command

// runBoxInDir runs box command in specified directory.
func runBoxInDir(t *testing.T, dir, name string) {
	t.Helper()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Logf("failed to restore directory: %v", err)
		}
	}()

	verifyCommandSuccess(t, exec.Command("bento", "box", name), "Created")
}

// verifyBentoJSONValid checks if bento file is valid JSON with correct ID.
func verifyBentoJSONValid(t *testing.T, path, expectedID string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("Bento file was not created")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var def map[string]interface{}
	if err := json.Unmarshal(content, &def); err != nil {
		t.Errorf("Created bento is not valid JSON: %v", err)
	}

	if def["id"] != expectedID {
		t.Errorf("Bento ID = %v, want %s", def["id"], expectedID)
	}

	if def["type"] != "group" {
		t.Errorf("Bento type = %v, want group", def["type"])
	}
}

// TestBoxCommand_CreateTemplate verifies box creates a template bento.
func TestBoxCommand_CreateTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	bentoPath := filepath.Join(tmpDir, "my-workflow.bento.json")

	runBoxInDir(t, tmpDir, "my-workflow")
	verifyBentoJSONValid(t, bentoPath, "my-workflow")
}

// createExistingFile creates a file in the directory.
func createExistingFile(t *testing.T, dir, name string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("existing content"), 0644); err != nil {
		t.Fatal(err)
	}
}

// changeToDir changes to directory and sets up cleanup.
func changeToDir(t *testing.T, dir string) {
	t.Helper()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Logf("failed to restore directory: %v", err)
		}
	})
}

// TestBoxCommand_OverwriteProtection verifies box doesn't overwrite existing files.
func TestBoxCommand_OverwriteProtection(t *testing.T) {
	tmpDir := t.TempDir()
	createExistingFile(t, tmpDir, "existing.bento.json")
	changeToDir(t, tmpDir)

	cmd := exec.Command("bento", "box", "existing")
	output, err := cmd.CombinedOutput()

	if err == nil {
		outputStr := strings.ToLower(string(output))
		if !strings.Contains(outputStr, "exists") && !strings.Contains(outputStr, "already") {
			t.Error("Should warn or error about existing file")
		}
	}
}
