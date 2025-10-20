package integration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFolderSetup_CreateProductFolders tests the folder setup bento.
// This validates that the loop neta correctly iterates over CSV rows
// and that the file-system neta creates directories for each product.
func TestFolderSetup_CreateProductFolders(t *testing.T) {
	projectRoot := "../../"

	// Cleanup products/ folder before and after test (at project root)
	defer CleanupTestDir(t, projectRoot+"products")
	CleanupTestDir(t, projectRoot+"products")

	// Execute bento (runs from project root, CSV path is relative to project root)
	bentoPath := "examples/phase8/folder-setup.bento.json"
	output, err := RunBento(t, bentoPath, nil)

	// Verify execution succeeded
	require.NoError(t, err, "Folder setup bento should execute successfully\nOutput: %s", string(output))

	// Verify output shows success
	assert.Contains(t, output, "Delicious! Bento executed successfully", "Should show success message")

	// Verify folders were created (relative to project root)
	// RunBento runs from project root, so products/ is created there
	VerifyFileExists(t, projectRoot+"products/MOCK-001")
	VerifyFileExists(t, projectRoot+"products/MOCK-002")
	VerifyFileExists(t, projectRoot+"products/MOCK-003")

	// Verify they are directories
	info, err := os.Stat(projectRoot + "products/MOCK-001")
	require.NoError(t, err, "products/MOCK-001 should exist")
	assert.True(t, info.IsDir(), "products/MOCK-001 should be a directory")

	info, err = os.Stat(projectRoot + "products/MOCK-002")
	require.NoError(t, err, "products/MOCK-002 should exist")
	assert.True(t, info.IsDir(), "products/MOCK-002 should be a directory")

	info, err = os.Stat(projectRoot + "products/MOCK-003")
	require.NoError(t, err, "products/MOCK-003 should exist")
	assert.True(t, info.IsDir(), "products/MOCK-003 should be a directory")

	t.Log("✓ Successfully created folders for 3 products")
	t.Log("✓ Validated loop neta with forEach mode")
	t.Log("✓ Validated context passing from CSV to loop")
}

// TestFolderSetup_AlreadyExists tests idempotency.
// Running the bento twice should succeed without errors.
func TestFolderSetup_AlreadyExists(t *testing.T) {
	projectRoot := "../../"

	// Cleanup products/ folder before and after test (at project root)
	defer CleanupTestDir(t, projectRoot+"products")
	CleanupTestDir(t, projectRoot+"products")

	// Pre-create one folder to test idempotency
	err := os.MkdirAll(projectRoot+"products/MOCK-001", 0755)
	require.NoError(t, err, "Should pre-create folder")

	// Execute bento (runs from project root)
	bentoPath := "examples/phase8/folder-setup.bento.json"
	output, err := RunBento(t, bentoPath, nil)

	// Should still succeed (idempotent)
	require.NoError(t, err, "Should succeed even if folder exists\nOutput: %s", string(output))

	// All folders should exist
	VerifyFileExists(t, projectRoot+"products/MOCK-001")
	VerifyFileExists(t, projectRoot+"products/MOCK-002")
	VerifyFileExists(t, projectRoot+"products/MOCK-003")

	t.Log("✓ Idempotent: Successfully ran with pre-existing folder")
}
