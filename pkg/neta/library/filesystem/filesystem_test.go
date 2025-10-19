package filesystem_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Develonaut/bento/pkg/neta/library/filesystem"
)

// TestFileSystem_ReadFile tests reading a file.
func TestFileSystem_ReadFile(t *testing.T) {
	ctx := context.Background()

	// Create temp file
	tmpfile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	content := "Hello, bento!"
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	fs := filesystem.New()

	params := map[string]interface{}{
		"operation": "read",
		"path":      tmpfile.Name(),
	}

	result, err := fs.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output := result.(map[string]interface{})
	if output["content"] != content {
		t.Errorf("content = %v, want %v", output["content"], content)
	}
}

// TestFileSystem_WriteFile tests writing to a file.
func TestFileSystem_WriteFile(t *testing.T) {
	ctx := context.Background()

	tmpfile := filepath.Join(os.TempDir(), "test-write.txt")
	defer os.Remove(tmpfile)

	fs := filesystem.New()

	content := "Test content for writing"
	params := map[string]interface{}{
		"operation": "write",
		"path":      tmpfile,
		"content":   content,
	}

	result, err := fs.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output := result.(map[string]interface{})
	if output["path"] != tmpfile {
		t.Errorf("path = %v, want %v", output["path"], tmpfile)
	}

	// Verify file was written
	data, err := os.ReadFile(tmpfile)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	if string(data) != content {
		t.Errorf("file content = %v, want %v", string(data), content)
	}
}

// TestFileSystem_CopyFile tests copying a file.
func TestFileSystem_CopyFile(t *testing.T) {
	ctx := context.Background()

	// Create source file
	srcfile, err := os.CreateTemp("", "test-src-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(srcfile.Name())

	content := "Content to copy"
	if _, err := srcfile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	srcfile.Close()

	// Destination path
	dstfile := filepath.Join(os.TempDir(), "test-dst.txt")
	defer os.Remove(dstfile)

	fs := filesystem.New()

	params := map[string]interface{}{
		"operation": "copy",
		"source":    srcfile.Name(),
		"dest":      dstfile,
	}

	result, err := fs.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output := result.(map[string]interface{})
	if output["dest"] != dstfile {
		t.Errorf("dest = %v, want %v", output["dest"], dstfile)
	}

	// Verify destination file exists with same content
	data, err := os.ReadFile(dstfile)
	if err != nil {
		t.Fatalf("Failed to read copied file: %v", err)
	}

	if string(data) != content {
		t.Errorf("copied content = %v, want %v", string(data), content)
	}
}

// TestFileSystem_MoveFile tests moving a file.
func TestFileSystem_MoveFile(t *testing.T) {
	ctx := context.Background()

	// Create source file
	srcfile, err := os.CreateTemp("", "test-move-src-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(srcfile.Name())

	content := "Content to move"
	if _, err := srcfile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	srcfile.Close()

	// Destination path
	dstfile := filepath.Join(os.TempDir(), "test-move-dst.txt")
	defer os.Remove(dstfile)

	fs := filesystem.New()

	params := map[string]interface{}{
		"operation": "move",
		"source":    srcfile.Name(),
		"dest":      dstfile,
	}

	result, err := fs.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output := result.(map[string]interface{})
	if output["dest"] != dstfile {
		t.Errorf("dest = %v, want %v", output["dest"], dstfile)
	}

	// Verify source no longer exists
	if _, err := os.Stat(srcfile.Name()); !os.IsNotExist(err) {
		t.Error("Source file still exists after move")
	}

	// Verify destination exists with correct content
	data, err := os.ReadFile(dstfile)
	if err != nil {
		t.Fatalf("Failed to read moved file: %v", err)
	}

	if string(data) != content {
		t.Errorf("moved content = %v, want %v", string(data), content)
	}
}

// TestFileSystem_DeleteFile tests deleting a file.
func TestFileSystem_DeleteFile(t *testing.T) {
	ctx := context.Background()

	// Create temp file
	tmpfile, err := os.CreateTemp("", "test-delete-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpfile.Close()

	fs := filesystem.New()

	params := map[string]interface{}{
		"operation": "delete",
		"path":      tmpfile.Name(),
	}

	result, err := fs.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output := result.(map[string]interface{})
	if output["deleted"] != true {
		t.Errorf("deleted = %v, want true", output["deleted"])
	}

	// Verify file no longer exists
	if _, err := os.Stat(tmpfile.Name()); !os.IsNotExist(err) {
		t.Error("File still exists after delete")
	}
}

// TestFileSystem_CreateDirectory tests creating a directory.
func TestFileSystem_CreateDirectory(t *testing.T) {
	ctx := context.Background()

	tmpdir := filepath.Join(os.TempDir(), "test-mkdir")
	defer os.RemoveAll(tmpdir)

	fs := filesystem.New()

	params := map[string]interface{}{
		"operation": "mkdir",
		"path":      tmpdir,
	}

	result, err := fs.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output := result.(map[string]interface{})
	if output["path"] != tmpdir {
		t.Errorf("path = %v, want %v", output["path"], tmpdir)
	}

	// Verify directory exists
	info, err := os.Stat(tmpdir)
	if err != nil {
		t.Fatalf("Directory was not created: %v", err)
	}

	if !info.IsDir() {
		t.Error("Path exists but is not a directory")
	}
}

// TestFileSystem_Exists tests checking if a file exists.
func TestFileSystem_Exists(t *testing.T) {
	ctx := context.Background()

	// Create temp file
	tmpfile, err := os.CreateTemp("", "test-exists-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	fs := filesystem.New()

	// Test existing file
	params := map[string]interface{}{
		"operation": "exists",
		"path":      tmpfile.Name(),
	}

	result, err := fs.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output := result.(map[string]interface{})
	if output["exists"] != true {
		t.Errorf("exists = %v, want true", output["exists"])
	}

	// Test non-existing file
	params["path"] = "/nonexistent/path/file.txt"

	result, err = fs.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output = result.(map[string]interface{})
	if output["exists"] != false {
		t.Errorf("exists = %v, want false", output["exists"])
	}
}

// TestFileSystem_InvalidOperation tests error handling for invalid operations.
func TestFileSystem_InvalidOperation(t *testing.T) {
	ctx := context.Background()

	fs := filesystem.New()

	params := map[string]interface{}{
		"operation": "invalid",
		"path":      "/some/path",
	}

	_, err := fs.Execute(ctx, params)
	if err == nil {
		t.Fatal("Expected error for invalid operation, got nil")
	}

	if !strings.Contains(err.Error(), "unsupported operation") {
		t.Errorf("Expected 'unsupported operation' error, got: %v", err)
	}
}

// TestFileSystem_MissingPath tests error handling when path is missing.
func TestFileSystem_MissingPath(t *testing.T) {
	ctx := context.Background()

	fs := filesystem.New()

	params := map[string]interface{}{
		"operation": "read",
		// Missing "path" parameter
	}

	_, err := fs.Execute(ctx, params)
	if err == nil {
		t.Fatal("Expected error for missing path, got nil")
	}
}
