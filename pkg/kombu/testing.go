package kombu

// SetConfigDirForTesting allows tests to override the config directory.
// Returns a cleanup function that should be called with t.Cleanup().
//
// This is exported for use by other packages' tests (e.g., miso tests).
// Example usage:
//
//	cleanup := kombu.SetConfigDirForTesting(t.TempDir())
//	t.Cleanup(cleanup)
func SetConfigDirForTesting(dir string) func() {
	original := configDir
	configDir = func() (string, error) {
		return dir, nil
	}
	return func() {
		configDir = original
	}
}
