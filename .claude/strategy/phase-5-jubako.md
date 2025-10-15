# Phase 5: Jubako (Storage)

**Status**: Pending
**Duration**: 2-3 hours
**Prerequisites**: Phase 4 complete, Karen approved

## Overview

Implement the storage layer - "Jubako" (重箱 - stacked boxes). This package handles .bento.yaml file management, workflow discovery, history tracking, and import/export functionality.

## Pre-Work Checklist

Before starting, you MUST:

1. ✅ Read [BENTO_BOX_PRINCIPLE.md](../BENTO_BOX_PRINCIPLE.md)
2. ✅ Confirm: "I understand the Bento Box Principle and will follow it"
3. ✅ Use TodoWrite to track all tasks
4. ✅ Phase 4 approved by Karen

## Goals

1. Implement .bento.yaml file parsing and validation
2. Create workflow discovery and indexing
3. Add execution history tracking
4. Implement import/export functionality
5. File watching for auto-reload (optional)
6. Validate Bento Box compliance

## Storage Structure

```
pkg/jubako/
├── go.mod
├── parser.go           # YAML parsing
├── parser_test.go
├── store.go            # File storage management
├── store_test.go
├── history.go          # Execution history
├── history_test.go
├── discovery.go        # Workflow discovery
├── discovery_test.go
└── types.go            # Storage types
```

## Deliverables

### 1. Parser

**Purpose**: Parse and validate .bento.yaml files
**File**: `pkg/jubako/parser.go`
**File Size Target**: < 150 lines

```go
// Package jubako provides storage and file management for Bento workflows.
// Jubako (重箱) means "stacked boxes" - a traditional Japanese food container.
package jubako

import (
    "fmt"
    "os"

    "gopkg.in/yaml.v3"

    "bento/pkg/neta"
)

// Parser handles .bento.yaml file parsing.
type Parser struct{}

// NewParser creates a new parser.
func NewParser() *Parser {
    return &Parser{}
}

// Parse reads and parses a .bento.yaml file.
func (p *Parser) Parse(path string) (neta.Definition, error) {
    data, err := readFile(path)
    if err != nil {
        return neta.Definition{}, fmt.Errorf("read failed: %w", err)
    }

    return p.ParseBytes(data)
}

// ParseBytes parses .bento.yaml from bytes.
func (p *Parser) ParseBytes(data []byte) (neta.Definition, error) {
    var def neta.Definition
    if err := yaml.Unmarshal(data, &def); err != nil {
        return neta.Definition{}, fmt.Errorf("invalid YAML: %w", err)
    }

    if err := validateDefinition(def); err != nil {
        return neta.Definition{}, fmt.Errorf("validation failed: %w", err)
    }

    return def, nil
}

// Format converts a definition to YAML.
func (p *Parser) Format(def neta.Definition) ([]byte, error) {
    data, err := yaml.Marshal(def)
    if err != nil {
        return nil, fmt.Errorf("marshal failed: %w", err)
    }
    return data, nil
}

// readFile reads a file from disk.
func readFile(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    return data, nil
}

// validateDefinition ensures a definition is well-formed.
func validateDefinition(def neta.Definition) error {
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

**Bento Box Compliance**:
- ✅ Single responsibility: YAML parsing
- ✅ Functions < 20 lines
- ✅ Clear error messages
- ✅ File < 150 lines

### 2. Store

**Purpose**: Workflow file storage management
**File**: `pkg/jubako/store.go`
**File Size Target**: < 200 lines

```go
package jubako

import (
    "fmt"
    "os"
    "path/filepath"

    "bento/pkg/neta"
)

// Store manages workflow file storage.
type Store struct {
    workDir string
    parser  *Parser
}

// NewStore creates a new store.
func NewStore(workDir string) (*Store, error) {
    if err := ensureDir(workDir); err != nil {
        return nil, err
    }

    return &Store{
        workDir: workDir,
        parser:  NewParser(),
    }, nil
}

// Load reads a workflow by name.
func (s *Store) Load(name string) (neta.Definition, error) {
    path := s.pathFor(name)
    return s.parser.Parse(path)
}

// Save writes a workflow to disk.
func (s *Store) Save(name string, def neta.Definition) error {
    path := s.pathFor(name)

    data, err := s.parser.Format(def)
    if err != nil {
        return err
    }

    return writeFile(path, data)
}

// Delete removes a workflow.
func (s *Store) Delete(name string) error {
    path := s.pathFor(name)
    return os.Remove(path)
}

// List returns all workflows in the store.
func (s *Store) List() ([]WorkflowInfo, error) {
    pattern := filepath.Join(s.workDir, "*.bento.yaml")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return nil, err
    }

    infos := make([]WorkflowInfo, 0, len(matches))
    for _, path := range matches {
        info, err := s.getInfo(path)
        if err != nil {
            continue // Skip invalid files
        }
        infos = append(infos, info)
    }

    return infos, nil
}

// pathFor returns the file path for a workflow name.
func (s *Store) pathFor(name string) string {
    if !filepath.Ext(name) == ".bento.yaml" {
        name += ".bento.yaml"
    }
    return filepath.Join(s.workDir, name)
}

// getInfo extracts workflow info from a file.
func (s *Store) getInfo(path string) (WorkflowInfo, error) {
    def, err := s.parser.Parse(path)
    if err != nil {
        return WorkflowInfo{}, err
    }

    stat, err := os.Stat(path)
    if err != nil {
        return WorkflowInfo{}, err
    }

    return WorkflowInfo{
        Name:     filepath.Base(path),
        Path:     path,
        Type:     def.Type,
        Modified: stat.ModTime(),
    }, nil
}

// ensureDir creates directory if it doesn't exist.
func ensureDir(path string) error {
    return os.MkdirAll(path, 0755)
}

// writeFile writes data to a file.
func writeFile(path string, data []byte) error {
    return os.WriteFile(path, data, 0644)
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: File management
- ✅ Functions < 20 lines
- ✅ Clear API (Load, Save, Delete, List)
- ✅ File < 200 lines

### 3. Types

**File**: `pkg/jubako/types.go`
**File Size Target**: < 100 lines

```go
package jubako

import "time"

// WorkflowInfo contains metadata about a workflow file.
type WorkflowInfo struct {
    Name     string
    Path     string
    Type     string
    Modified time.Time
}

// ExecutionRecord tracks a workflow execution.
type ExecutionRecord struct {
    ID        string
    Workflow  string
    StartTime time.Time
    EndTime   time.Time
    Success   bool
    Error     string
    Result    interface{}
}

// HistoryFilter filters execution history.
type HistoryFilter struct {
    Workflow   string
    SuccessOnly bool
    Limit      int
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Type definitions
- ✅ Clear data structures
- ✅ File < 100 lines

### 4. Discovery

**Purpose**: Find .bento.yaml files in directories
**File**: `pkg/jubako/discovery.go`
**File Size Target**: < 150 lines

```go
package jubako

import (
    "os"
    "path/filepath"
)

// Discovery finds workflow files.
type Discovery struct {
    searchPaths []string
}

// NewDiscovery creates a new discovery instance.
func NewDiscovery(paths ...string) *Discovery {
    if len(paths) == 0 {
        paths = defaultPaths()
    }

    return &Discovery{
        searchPaths: paths,
    }
}

// Find searches for .bento.yaml files.
func (d *Discovery) Find() ([]string, error) {
    found := []string{}

    for _, path := range d.searchPaths {
        files, err := findInPath(path)
        if err != nil {
            continue // Skip inaccessible paths
        }
        found = append(found, files...)
    }

    return found, nil
}

// Watch monitors directories for changes (optional, Phase 5+).
func (d *Discovery) Watch() (<-chan string, error) {
    // Future: fsnotify integration
    return nil, nil
}

// findInPath searches a single path for .bento.yaml files.
func findInPath(root string) ([]string, error) {
    found := []string{}

    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return nil // Skip errors
        }

        if info.IsDir() {
            return nil
        }

        if filepath.Ext(path) == ".yaml" && isBentoFile(path) {
            found = append(found, path)
        }

        return nil
    })

    return found, err
}

// isBentoFile checks if a file is a .bento.yaml file.
func isBentoFile(path string) bool {
    base := filepath.Base(path)
    return filepath.Ext(base) == ".yaml" &&
           len(base) > 11 &&
           base[len(base)-11:] == ".bento.yaml"
}

// defaultPaths returns default search paths.
func defaultPaths() []string {
    home, err := os.UserHomeDir()
    if err != nil {
        return []string{"."}
    }

    return []string{
        ".",
        filepath.Join(home, ".bento"),
        filepath.Join(home, "bento"),
    }
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: File discovery
- ✅ Functions < 20 lines
- ✅ Graceful error handling
- ✅ File < 150 lines

### 5. History

**Purpose**: Track workflow execution history
**File**: `pkg/jubako/history.go`
**File Size Target**: < 200 lines

```go
package jubako

import (
    "encoding/json"
    "os"
    "path/filepath"
    "time"

    "github.com/google/uuid"
)

// History manages execution history.
type History struct {
    historyDir string
}

// NewHistory creates a new history manager.
func NewHistory(dir string) (*History, error) {
    if err := ensureDir(dir); err != nil {
        return nil, err
    }

    return &History{historyDir: dir}, nil
}

// Record saves an execution record.
func (h *History) Record(rec ExecutionRecord) error {
    if rec.ID == "" {
        rec.ID = uuid.New().String()
    }

    path := h.recordPath(rec.ID)
    data, err := json.MarshalIndent(rec, "", "  ")
    if err != nil {
        return err
    }

    return writeFile(path, data)
}

// Get retrieves an execution record by ID.
func (h *History) Get(id string) (ExecutionRecord, error) {
    path := h.recordPath(id)
    data, err := os.ReadFile(path)
    if err != nil {
        return ExecutionRecord{}, err
    }

    var rec ExecutionRecord
    if err := json.Unmarshal(data, &rec); err != nil {
        return ExecutionRecord{}, err
    }

    return rec, nil
}

// List returns execution history with optional filtering.
func (h *History) List(filter HistoryFilter) ([]ExecutionRecord, error) {
    files, err := h.listFiles()
    if err != nil {
        return nil, err
    }

    records := []ExecutionRecord{}
    for _, file := range files {
        rec, err := h.loadRecord(file)
        if err != nil {
            continue
        }

        if matchesFilter(rec, filter) {
            records = append(records, rec)
        }

        if filter.Limit > 0 && len(records) >= filter.Limit {
            break
        }
    }

    return records, nil
}

// Clear removes all history records.
func (h *History) Clear() error {
    return os.RemoveAll(h.historyDir)
}

// recordPath returns the file path for a record.
func (h *History) recordPath(id string) string {
    return filepath.Join(h.historyDir, id+".json")
}

// listFiles returns all history files sorted by modification time.
func (h *History) listFiles() ([]string, error) {
    pattern := filepath.Join(h.historyDir, "*.json")
    return filepath.Glob(pattern)
}

// loadRecord loads a record from a file.
func (h *History) loadRecord(path string) (ExecutionRecord, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return ExecutionRecord{}, err
    }

    var rec ExecutionRecord
    if err := json.Unmarshal(data, &rec); err != nil {
        return ExecutionRecord{}, err
    }

    return rec, nil
}

// matchesFilter checks if a record matches the filter.
func matchesFilter(rec ExecutionRecord, filter HistoryFilter) bool {
    if filter.Workflow != "" && rec.Workflow != filter.Workflow {
        return false
    }

    if filter.SuccessOnly && !rec.Success {
        return false
    }

    return true
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: History tracking
- ✅ Functions < 20 lines
- ✅ JSON for persistence
- ✅ File < 200 lines

## Integration Points

### Store Integration with CLI

Update `cmd/bento/prepare.go`:

```go
import "bento/pkg/jubako"

func runPrepare(cmd *cobra.Command, args []string) error {
    parser := jubako.NewParser()
    def, err := parser.Parse(args[0])
    // ...
}
```

### History Integration with Executor

```go
func executeBento(def neta.Definition, timeout time.Duration) (interface{}, error) {
    hist, _ := jubako.NewHistory(historyDir)

    rec := jubako.ExecutionRecord{
        Workflow:  def.Name,
        StartTime: time.Now(),
    }

    result, err := chef.Execute(ctx, def)

    rec.EndTime = time.Now()
    rec.Success = (err == nil)
    rec.Result = result.Output

    hist.Record(rec)

    return result.Output, err
}
```

### Discovery Integration with TUI

Update `pkg/omise/screens/browser.go`:

```go
func loadWorkflows() []list.Item {
    disc := jubako.NewDiscovery()
    files, _ := disc.Find()

    items := make([]list.Item, len(files))
    for i, file := range files {
        items[i] = workflowItem{
            name: filepath.Base(file),
            path: file,
        }
    }
    return items
}
```

## Default Directory Structure

```
$HOME/.bento/
├── workflows/           # User workflows
│   ├── example.bento.yaml
│   └── my-flow.bento.yaml
├── history/             # Execution history
│   ├── abc123.json
│   └── def456.json
└── config.yaml          # Bento configuration
```

## Testing Strategy

### Parser Tests
```go
func TestParser_Parse(t *testing.T) {
    tests := []struct {
        name    string
        yaml    string
        want    neta.Definition
        wantErr bool
    }{
        {
            name: "valid http node",
            yaml: `
type: http
name: Test
parameters:
  url: https://example.com
`,
            wantErr: false,
        },
        // ... more test cases
    }
}
```

### Store Tests
- Test Load/Save/Delete operations
- Test List with multiple files
- Test error cases (missing files, invalid YAML)

### History Tests
- Test Record/Get operations
- Test List with filtering
- Test Clear operation

### Discovery Tests
- Test Find in various directories
- Test .bento.yaml file detection
- Test default paths

## Validation Commands

Before marking phase complete:

```bash
# Format
go fmt ./...

# Lint
golangci-lint run

# Test
go test -v -race ./pkg/jubako/...

# Integration test
mkdir -p ~/.bento/workflows
cat > ~/.bento/workflows/test.bento.yaml <<EOF
type: http
name: Test
parameters:
  url: https://httpbin.org/get
EOF

./bento prepare ~/.bento/workflows/test.bento.yaml
./bento pack ~/.bento/workflows/test.bento.yaml

# File size check
find pkg/jubako -name "*.go" -exec wc -l {} + | sort -rn
```

## Success Criteria

Phase 5 is complete when:

1. ✅ Parser implemented and tested
2. ✅ Store with Load/Save/Delete/List
3. ✅ History tracking functional
4. ✅ Discovery finding .bento.yaml files
5. ✅ Integration with CLI commands
6. ✅ Integration with TUI
7. ✅ All tests passing (>80% coverage)
8. ✅ All files < 250 lines
9. ✅ All functions < 20 lines
10. ✅ golangci-lint clean
11. ✅ **Karen's approval granted**

## Common Pitfalls to Avoid

1. ❌ **God store** - Keep parser, history, discovery separate
2. ❌ **Ignoring errors** - File I/O needs comprehensive error handling
3. ❌ **No validation** - Always validate YAML before parsing
4. ❌ **Hardcoded paths** - Use configurable directories
5. ❌ **Large files** - Split if approaching limits

## Final Integration

After Phase 5, all components work together:

```
User runs: bento

┌─────────────────────────────────────┐
│  Omise (TUI) - pkg/omise            │
│  • Bubble Tea interface             │
│  • Lists workflows from Jubako      │
│  • Shows execution progress         │
└─────────────────────────────────────┘
          ↓
┌─────────────────────────────────────┐
│  Jubako (Storage) - pkg/jubako      │
│  • Discovers .bento.yaml files      │
│  • Parses and validates             │
│  • Tracks execution history         │
└─────────────────────────────────────┘
          ↓
┌─────────────────────────────────────┐
│  Itamae (Orchestrator) - pkg/itamae │
│  • Executes workflows               │
│  • Coordinates node execution       │
└─────────────────────────────────────┘
          ↓
┌─────────────────────────────────────┐
│  Neta (Nodes) - pkg/neta            │
│  • HTTP, Transform, Conditional...  │
│  • Registered in Pantry             │
└─────────────────────────────────────┘
```

## Project Complete!

After Phase 5 approval from Karen, the Bento project is **feature complete**:

✅ Core packages (neta, itamae, pantry)
✅ Node library (http, transform, conditional, loop, group)
✅ CLI commands (prepare, pack, pantry, taste)
✅ Interactive TUI (omise)
✅ Storage layer (jubako)
✅ **100% Bento Box compliant**

## Execution Prompt

```
I'm ready to begin Phase 5: Jubako Storage.

I have read the Bento Box Principle and will follow it.

Please implement the storage layer:
- Parser for .bento.yaml files
- Store for file management
- History for execution tracking
- Discovery for finding workflows

Each component in focused file < 200 lines. I will use TodoWrite to track progress and get Karen's approval before completing.
```

---

**Phase 5 Jubako**: Stacked boxes storage layer 🍱

**After this phase**: Bento is complete and ready for users! 🎉
