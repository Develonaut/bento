# Cross-Platform Compatibility Research for Bento CLI

## Executive Summary

The Bento CLI codebase is currently **100% cross-platform compatible**. No platform-specific code or problematic patterns were found. The codebase follows Go best practices by using standard library functions appropriately and avoiding platform-specific assumptions. The only minor recommendations are for future-proofing as the project grows, particularly around file operations and build configurations.

## Current State Analysis

### Platform-Specific Code Found
**NONE** - The current codebase contains no platform-specific code patterns or hard-coded assumptions.

### Clean Cross-Platform Code
✅ **All current code is platform-agnostic:**
- No hard-coded path separators (`/` or `\`)
- No Unix-specific file permissions
- No platform-specific environment variables
- No shell-specific commands
- Uses standard `os.Stderr` for error output
- Clean use of `os.Exit()` for process termination
- Cobra dependency is fully cross-platform

## Go Best Practices for Cross-Platform CLIs

### File Path Handling

**Standard Library Packages:**
- `path/filepath` - For OS-specific file paths (ALWAYS use for filesystem operations)
- `path` - For URL paths and slash-separated paths (never for filesystem)

**Key Functions to Use:**
```go
filepath.Join()          // Constructs paths with correct separator
filepath.Clean()         // Cleans paths for the OS
filepath.Abs()           // Gets absolute path
filepath.FromSlash()     // Converts "/" to OS separator
filepath.ToSlash()       // Converts OS separator to "/"
os.PathSeparator         // The OS-specific path separator
os.PathListSeparator     // For PATH-like lists (: on Unix, ; on Windows)
```

### Environment & Configuration

**Home Directory Detection:**
```go
// Go 1.12+ provides os.UserHomeDir()
home, err := os.UserHomeDir()

// Config locations per platform:
// Linux:   $XDG_CONFIG_HOME or ~/.config/bento/
// macOS:   ~/Library/Application Support/bento/
// Windows: %APPDATA%\bento\
```

**Standard Config Paths:**
```go
func ConfigDir() (string, error) {
    switch runtime.GOOS {
    case "linux":
        if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
            return filepath.Join(xdg, "bento"), nil
        }
        home, err := os.UserHomeDir()
        if err != nil {
            return "", err
        }
        return filepath.Join(home, ".config", "bento"), nil
    case "darwin":
        home, err := os.UserHomeDir()
        if err != nil {
            return "", err
        }
        return filepath.Join(home, "Library", "Application Support", "bento"), nil
    case "windows":
        return filepath.Join(os.Getenv("APPDATA"), "bento"), nil
    default:
        home, err := os.UserHomeDir()
        if err != nil {
            return "", err
        }
        return filepath.Join(home, ".bento"), nil
    }
}
```

### Build & Distribution

**Cross-Compilation Commands:**
```bash
# Build for all major platforms
GOOS=linux GOARCH=amd64 go build -o bin/bento-linux-amd64 ./cmd/bento
GOOS=darwin GOARCH=amd64 go build -o bin/bento-darwin-amd64 ./cmd/bento
GOOS=darwin GOARCH=arm64 go build -o bin/bento-darwin-arm64 ./cmd/bento
GOOS=windows GOARCH=amd64 go build -o bin/bento.exe ./cmd/bento
```

**Binary Naming Conventions:**
- Linux: `bento` or `bento-linux-amd64`
- macOS: `bento` or `bento-darwin-amd64` / `bento-darwin-arm64`
- Windows: `bento.exe` or `bento-windows-amd64.exe`

### Testing Strategies

**Platform-Specific Test Files:**
```go
// file_unix_test.go
//go:build !windows

// file_windows_test.go
//go:build windows
```

**CI/CD Matrix Testing (GitHub Actions):**
```yaml
strategy:
  matrix:
    os: [ubuntu-latest, macos-latest, windows-latest]
    go: [1.22, 1.23]
runs-on: ${{ matrix.os }}
```

## Recommendations for Bento

### Immediate Actions
**NO CRITICAL ACTIONS REQUIRED** - Current code is clean.

### Future Considerations (As Bento Grows)

1. **File Operations (Phase 3/4)**
   - Always use `filepath.Join()` for paths
   - Use `0755` for directories, `0644` for files (Go handles Windows conversion)
   - Test file operations on all platforms

2. **Shell/Command Execution (When Added)**
   ```go
   // Use exec.Command correctly
   cmd := exec.Command("sh", "-c", script)  // Unix
   cmd := exec.Command("cmd", "/c", script) // Windows

   // Or better, use exec.LookPath for tools:
   path, err := exec.LookPath("git")
   ```

3. **Configuration Files**
   - Implement proper config directory detection
   - Support both YAML and JSON (both are cross-platform)
   - Use `os.UserConfigDir()` (Go 1.13+) or custom function above

4. **Terminal/TTY Detection (Phase 4 TUI)**
   ```go
   // Check if stdout is a terminal
   if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
       // Terminal attached
   }
   ```

5. **Signal Handling**
   ```go
   // Cross-platform signal handling
   sigChan := make(chan os.Signal, 1)
   signal.Notify(sigChan, os.Interrupt) // Works on all platforms
   // SIGTERM is Unix-only, use build tags if needed
   ```

### Build Configuration

**Enhanced Makefile for Cross-Platform Builds:**
```makefile
# Detect OS
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
    BINARY_NAME = bento
endif
ifeq ($(UNAME_S),Darwin)
    BINARY_NAME = bento
endif
ifeq ($(OS),Windows_NT)
    BINARY_NAME = bento.exe
endif

.PHONY: build-all
build-all: build-linux build-darwin build-windows

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/bento-linux-amd64 ./cmd/bento

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o bin/bento-darwin-amd64 ./cmd/bento
	GOOS=darwin GOARCH=arm64 go build -o bin/bento-darwin-arm64 ./cmd/bento

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/bento-windows-amd64.exe ./cmd/bento

# Clean with proper Windows support
clean:
	$(RM) -r bin/
```

### Documentation Needed

Create `.claude/docs/CROSS_PLATFORM.md`:
```markdown
# Cross-Platform Development Guidelines

## Always Use
- `filepath.Join()` for paths, never string concatenation
- `os.PathSeparator` when needed, never hard-code `/` or `\`
- `os.UserHomeDir()` or `os.UserConfigDir()` for user directories
- Build tags for platform-specific code

## Never Use
- Hard-coded paths like `/tmp` (use `os.TempDir()`)
- Unix-only signals without build tags
- Shell-specific syntax without checking `runtime.GOOS`
- Assumptions about file permissions on Windows

## Testing
- Run tests on all platforms via CI
- Use GitHub Actions matrix builds
- Test file operations explicitly on Windows
```

## Code Examples

### Good Patterns to Use
```go
// Correct path construction
configPath := filepath.Join(homeDir, ".config", "bento", "config.yaml")

// Proper temp file creation
tmpFile, err := os.CreateTemp("", "bento-*.json")

// Cross-platform home directory
home, err := os.UserHomeDir()

// Platform detection when necessary
switch runtime.GOOS {
case "windows":
    // Windows-specific code
case "darwin":
    // macOS-specific code
default:
    // Linux and others
}

// Safe file permissions (Go handles Windows)
err := os.MkdirAll(configDir, 0755)
err := os.WriteFile(configFile, data, 0644)
```

### Patterns to Avoid
```go
// BAD: Hard-coded separators
configPath := homeDir + "/.config/bento/config.yaml"  // WRONG!

// BAD: Unix-only paths
tmpFile := "/tmp/bento.json"  // WRONG!

// BAD: Assuming shell
cmd := exec.Command("sh", "-c", script)  // Check GOOS first!

// BAD: Platform-specific without build tags
signal.Notify(sigChan, syscall.SIGTERM)  // SIGTERM doesn't exist on Windows!

// BAD: String manipulation for paths
path := strings.Replace(unixPath, "/", "\\", -1)  // Use filepath!
```

## Reference Resources

- [Go Documentation - File Paths](https://pkg.go.dev/path/filepath)
- [Go Documentation - Runtime](https://pkg.go.dev/runtime)
- [Go Documentation - OS Package](https://pkg.go.dev/os)
- [Building Go Programs](https://go.dev/doc/tutorial/compile-install)
- [Go Cross Compilation](https://go.dev/doc/install/source#environment)
- [Cobra CLI Framework](https://github.com/spf13/cobra) - Fully cross-platform
- [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html)
- [Dave Cheney - Cross Compilation](https://dave.cheney.net/2015/08/22/cross-compilation-with-go-1-5)

## Conclusion

The Bento CLI is currently in excellent shape for cross-platform compatibility. The codebase follows Go idioms and uses no platform-specific code. The use of Cobra for CLI handling ensures consistent behavior across platforms.

**Approval Status: ✅ APPROVED**

The main recommendations are preventive measures for future development phases:
1. Maintain current clean practices
2. Use `filepath` package when file operations are added
3. Implement proper config directory detection when needed
4. Add cross-platform CI/CD testing before Phase 3
5. Document guidelines for contributors

The standard library has already solved cross-platform compatibility. Continue using it properly, and Bento will run seamlessly on Linux, macOS, and Windows.