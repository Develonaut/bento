# Phase 5.5: Definition Versioning

**Status**: Pending
**Duration**: 1-2 hours
**Prerequisites**: Phase 5 complete, Karen approved

## Overview

Add version field to bento definitions to enable future schema evolution and breaking changes. This ensures we can safely migrate workflows when the definition format needs to change.

## Pre-Work Checklist

Before starting, you MUST:

1. ✅ Read [BENTO_BOX_PRINCIPLE.md](../BENTO_BOX_PRINCIPLE.md)
2. ✅ Confirm: "I understand the Bento Box Principle and will follow it"
3. ✅ Use TodoWrite to track all tasks
4. ✅ Phase 5 approved by Karen

## Goals

1. Add `Version` field to `neta.Definition`
2. Implement version validation in parser
3. Update all existing example bentos with version
4. Add version compatibility checking
5. Lay foundation for future migration tooling
6. Validate Bento Box compliance

## Version Format

**Semantic Versioning:** `MAJOR.MINOR`

- **Major** (1.x) = Breaking changes to definition format
- **Minor** (x.1) = Backward compatible additions

**Initial Version:** `1.0`

**Compatibility Rules:**
- Parser accepts same major version (1.0 accepts 1.0, 1.1, 1.2...)
- Parser rejects different major versions (1.0 rejects 2.0)
- Clear error messages for incompatible versions

## Deliverables

### 1. Update Definition Type

**File**: `pkg/neta/definition.go`
**Target Size**: < 100 lines (currently ~50)

```go
// Package neta defines the core node types for Bento.
// Neta (ネタ) means "ingredients" or "toppings" in sushi terminology.
package neta

// Definition describes a node that can be executed by Itamae.
// It may be a single executable node or a group containing other nodes.
type Definition struct {
	// Version specifies the definition schema version (e.g., "1.0")
	// REQUIRED: Must be present in all .bento.yaml files
	// Format: MAJOR.MINOR (semantic versioning)
	Version string `yaml:"version" json:"version"`

	// Type identifies what kind of node this is (http, transform, group, etc)
	Type string `yaml:"type" json:"type"`

	// Name is the human-readable identifier for this node
	Name string `yaml:"name" json:"name"`

	// Parameters contains type-specific configuration.
	Parameters map[string]interface{} `yaml:"parameters,omitempty" json:"parameters,omitempty"`

	// Nodes contains child nodes (for group types)
	// If empty/nil, this is a leaf node
	Nodes []Definition `yaml:"nodes,omitempty" json:"nodes,omitempty"`
}

// CurrentVersion is the version of definitions this build supports
const CurrentVersion = "1.0"

// IsGroup returns true if this definition contains child nodes
func (d Definition) IsGroup() bool {
	return len(d.Nodes) > 0
}

// IsVersionCompatible checks if the definition version is compatible
func (d Definition) IsVersionCompatible() bool {
	return isCompatibleVersion(d.Version)
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Core definition type
- ✅ Clear constant for current version
- ✅ File < 100 lines

### 2. Add Version Validation

**File**: `pkg/neta/version.go` (NEW)
**Target Size**: < 100 lines

```go
package neta

import (
	"fmt"
	"strconv"
	"strings"
)

// isCompatibleVersion checks if a version string is compatible
func isCompatibleVersion(v string) bool {
	if v == "" {
		return false
	}

	major, err := parseMajorVersion(v)
	if err != nil {
		return false
	}

	currentMajor, err := parseMajorVersion(CurrentVersion)
	if err != nil {
		return false
	}

	return major == currentMajor
}

// parseMajorVersion extracts the major version number
func parseMajorVersion(v string) (int, error) {
	parts := strings.Split(v, ".")
	if len(parts) < 1 {
		return 0, fmt.Errorf("invalid version format")
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid major version: %w", err)
	}

	return major, nil
}

// ValidateVersion checks version and returns descriptive error
func ValidateVersion(v string) error {
	if v == "" {
		return fmt.Errorf("version is required (current version: %s)", CurrentVersion)
	}

	if !isCompatibleVersion(v) {
		return fmt.Errorf("incompatible version %s (current version: %s)", v, CurrentVersion)
	}

	return nil
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Version checking
- ✅ Separate file for version logic
- ✅ Functions < 20 lines
- ✅ Clear error messages

### 3. Update Parser Validation

**File**: `pkg/jubako/parser.go` (modify existing validateDefinition)

Update the `validateDefinition` function:

```go
// validateDefinition ensures a definition is well-formed.
func validateDefinition(def neta.Definition) error {
	// Validate version first
	if err := neta.ValidateVersion(def.Version); err != nil {
		return fmt.Errorf("version error: %w", err)
	}

	if def.Type == "" {
		return fmt.Errorf("type is required")
	}

	if def.IsGroup() {
		for i, child := range def.Nodes {
			if err := validateDefinition(child); err != nil {
				return fmt.Errorf("node %d: %w", i, err)
			}
		}
	}

	return nil
}
```

**Changes**:
- Add version validation as first check
- Recursive validation ensures all nodes have versions

### 4. Add Version Tests

**File**: `pkg/neta/version_test.go` (NEW)
**Target Size**: < 150 lines

```go
package neta

import (
	"testing"
)

func TestIsCompatibleVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    bool
	}{
		{
			name:    "same version",
			version: "1.0",
			want:    true,
		},
		{
			name:    "same major, different minor",
			version: "1.1",
			want:    true,
		},
		{
			name:    "different major",
			version: "2.0",
			want:    false,
		},
		{
			name:    "empty version",
			version: "",
			want:    false,
		},
		{
			name:    "invalid format",
			version: "abc",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isCompatibleVersion(tt.version)
			if got != tt.want {
				t.Errorf("isCompatibleVersion(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{
			name:    "valid version",
			version: "1.0",
			wantErr: false,
		},
		{
			name:    "missing version",
			version: "",
			wantErr: true,
		},
		{
			name:    "incompatible version",
			version: "2.0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVersion(%q) error = %v, wantErr %v", tt.version, err, tt.wantErr)
			}
		})
	}
}

func TestDefinition_IsVersionCompatible(t *testing.T) {
	tests := []struct {
		name string
		def  Definition
		want bool
	}{
		{
			name: "compatible version",
			def:  Definition{Version: "1.0", Type: "http"},
			want: true,
		},
		{
			name: "incompatible version",
			def:  Definition{Version: "2.0", Type: "http"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.def.IsVersionCompatible()
			if got != tt.want {
				t.Errorf("Definition.IsVersionCompatible() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

### 5. Update Parser Tests

**File**: `pkg/jubako/parser_test.go` (add test cases)

```go
func TestParser_ValidateVersion(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name    string
		yaml    string
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid version",
			yaml: `
version: "1.0"
type: http
name: Test
parameters:
  url: https://example.com
`,
			wantErr: false,
		},
		{
			name: "missing version",
			yaml: `
type: http
name: Test
parameters:
  url: https://example.com
`,
			wantErr: true,
			errMsg:  "version is required",
		},
		{
			name: "incompatible version",
			yaml: `
version: "2.0"
type: http
name: Test
parameters:
  url: https://example.com
`,
			wantErr: true,
			errMsg:  "incompatible version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser.ParseBytes([]byte(tt.yaml))

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("error message = %v, want to contain %v", err, tt.errMsg)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
```

## Example YAML Format

**Before (Phase 5)**:
```yaml
type: http
name: Fetch User Data
parameters:
  method: GET
  url: https://api.github.com/users/octocat
```

**After (Phase 5.5)**:
```yaml
version: "1.0"
type: http
name: Fetch User Data
parameters:
  method: GET
  url: https://api.github.com/users/octocat
```

## Migration Strategy

### For Existing Bentos

Create a migration tool (optional, for convenience):

**File**: `cmd/bento/migrate.go` (NEW)

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"bento/pkg/neta"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [file or directory]",
	Short: "Add version field to existing bento files",
	Args:  cobra.ExactArgs(1),
	RunE:  runMigrate,
}

func runMigrate(cmd *cobra.Command, args []string) error {
	path := args[0]

	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return migrateDir(path)
	}
	return migrateFile(path)
}

func migrateFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var def neta.Definition
	if err := yaml.Unmarshal(data, &def); err != nil {
		return err
	}

	// Already has version
	if def.Version != "" {
		fmt.Printf("✓ %s (already versioned: %s)\n", path, def.Version)
		return nil
	}

	// Add version
	def.Version = neta.CurrentVersion

	newData, err := yaml.Marshal(def)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, newData, 0644); err != nil {
		return err
	}

	fmt.Printf("✓ %s (added version: %s)\n", path, neta.CurrentVersion)
	return nil
}

func migrateDir(dir string) error {
	pattern := filepath.Join(dir, "*.bento.yaml")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for _, path := range matches {
		if err := migrateFile(path); err != nil {
			fmt.Printf("✗ %s: %v\n", path, err)
		}
	}

	return nil
}
```

## Future: Migration Framework

When version 2.0 arrives, we'll add:

**File**: `pkg/jubako/migrator.go` (FUTURE)

```go
package jubako

import (
	"fmt"
	"bento/pkg/neta"
)

// Migrator handles definition version upgrades
type Migrator struct {
	migrations map[string]MigrationFunc
}

// MigrationFunc transforms a definition to a newer version
type MigrationFunc func(neta.Definition) (neta.Definition, error)

// NewMigrator creates a migrator with registered migrations
func NewMigrator() *Migrator {
	m := &Migrator{
		migrations: make(map[string]MigrationFunc),
	}

	// Register migrations
	m.Register("1.0", "2.0", migrate1to2)

	return m
}

// Register adds a migration path
func (m *Migrator) Register(from, to string, fn MigrationFunc) {
	key := fmt.Sprintf("%s->%s", from, to)
	m.migrations[key] = fn
}

// Upgrade migrates a definition to target version
func (m *Migrator) Upgrade(def neta.Definition, target string) (neta.Definition, error) {
	// Implementation for chaining migrations
	return def, nil
}

// migrate1to2 is an example migration
func migrate1to2(def neta.Definition) (neta.Definition, error) {
	// Example: Rename field, transform structure, etc.
	def.Version = "2.0"
	return def, nil
}
```

## Integration Points

### CLI Commands

All commands now require version:

```bash
# prepare validates version
bento prepare workflow.bento.yaml

# pack requires valid version
bento pack workflow.bento.yaml
```

### Jubako Store

Store operations automatically validate version:

```go
// Load validates version
def, err := store.Load("my-workflow")
if err != nil {
    // "incompatible version 2.0 (current version: 1.0)"
}

// Save ensures version is set
def.Version = neta.CurrentVersion
store.Save("my-workflow", def)
```

### Omise TUI

Display version in UI:

```go
// Browser shows version
fmt.Sprintf("%s (v%s)", workflow.Name, workflow.Version)

// Editor sets version on new bentos
def := neta.Definition{
    Version: neta.CurrentVersion,
    Type:    selectedType,
    Name:    enteredName,
}
```

## Validation Commands

```bash
# Format
go fmt ./pkg/neta/... ./pkg/jubako/...

# Test version logic
go test -v ./pkg/neta/

# Test parser validation
go test -v ./pkg/jubako/

# Integration test
cat > test-version.bento.yaml <<EOF
version: "1.0"
type: http
name: Test
parameters:
  url: https://httpbin.org/get
EOF

./bento prepare test-version.bento.yaml  # Should pass

cat > test-bad-version.bento.yaml <<EOF
version: "2.0"
type: http
name: Test
parameters:
  url: https://httpbin.org/get
EOF

./bento prepare test-bad-version.bento.yaml  # Should fail with clear error
```

## Success Criteria

Phase 5.5 is complete when:

1. ✅ `neta.Definition` has `Version` field
2. ✅ `neta.CurrentVersion` constant defined
3. ✅ Version validation implemented
4. ✅ Parser rejects missing/incompatible versions
5. ✅ All tests passing
6. ✅ Clear error messages
7. ✅ Example bentos updated
8. ✅ Migration tool created (optional)
9. ✅ Files < 150 lines
10. ✅ Functions < 20 lines
11. ✅ **Karen's approval granted**

## Common Pitfalls to Avoid

1. ❌ **Forgetting version in tests** - All test YAML must include version
2. ❌ **Complex version logic** - Keep it simple: major version only
3. ❌ **No validation** - Must validate version before other fields
4. ❌ **Poor error messages** - Tell user current version and how to fix
5. ❌ **Breaking existing bentos** - Migration tool should help transition

## Documentation Updates

Update README.md with versioned examples:

```markdown
### Your First Workflow

Create a `hello.bento.yaml` file:

```yaml
version: "1.0"
type: http
name: Fetch User Data
parameters:
  method: GET
  url: https://api.github.com/users/octocat
```
```

## Execution Prompt

```
I'm ready to begin Phase 5.5: Definition Versioning.

I have read the Bento Box Principle and will follow it.

Please add versioning to bento definitions:
- Add Version field to neta.Definition
- Implement version validation
- Update parser to check versions
- Add comprehensive tests
- Create migration tool

Each file < 150 lines, functions < 20 lines. I will use TodoWrite to track progress and get Karen's approval before completing.
```

---

**Phase 5.5 Definition Versioning**: Future-proof schema evolution 🔢

**After this phase**: All bentos have versions, ready for CRUD operations!
