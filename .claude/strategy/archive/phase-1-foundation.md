# Phase 1: Foundation

**Status**: Pending
**Duration**: 2-3 hours
**Prerequisites**: Phase 0 complete

## Overview

Establish the Go workspace structure and core packages that form the foundation of Bento. This phase creates the fundamental types and orchestration patterns that all future phases will build upon.

## Pre-Work Checklist

Before starting, you MUST:

1. ✅ Read [BENTO_BOX_PRINCIPLE.md](../BENTO_BOX_PRINCIPLE.md)
2. ✅ Confirm: "I understand the Bento Box Principle and will follow it"
3. ✅ Use TodoWrite to track all tasks

## Goals

1. Initialize Go workspace with 6 modules
2. Create `pkg/neta/` - Core node definition types
3. Create `pkg/itamae/` - Chef/orchestrator foundation
4. Create `pkg/pantry/` - Node registry structure
5. Bootstrap CLI entry point with Cobra
6. Establish testing patterns
7. Validate Bento Box compliance

## Module Structure

```
bento/
├── go.work              # Workspace definition (gitignored)
├── go.work.sum
├── cmd/
│   └── bento/
│       ├── go.mod       # Main application module
│       ├── main.go
│       └── root.go      # Cobra root command
├── pkg/
│   ├── neta/
│   │   ├── go.mod
│   │   ├── definition.go       # Neta definition types
│   │   ├── definition_test.go
│   │   ├── executable.go       # Executable interface
│   │   └── executable_test.go
│   ├── itamae/
│   │   ├── go.mod
│   │   ├── itamae.go          # Chef/orchestrator
│   │   ├── itamae_test.go
│   │   ├── context.go         # Execution context
│   │   └── context_test.go
│   └── pantry/
│       ├── go.mod
│       ├── registry.go        # Node registry
│       ├── registry_test.go
│       └── lookup.go          # Node lookup
├── internal/
│   └── version/
│       └── version.go         # Version info
├── Makefile
└── README.md
```

## Deliverables

### 1. Go Workspace Initialization

**File**: `go.work`
```go
go 1.23

use (
    ./cmd/bento
    ./pkg/neta
    ./pkg/itamae
    ./pkg/pantry
    ./pkg/jubako
    ./pkg/omise
)
```

**Note**: `go.work` is gitignored (like .vscode settings)

### 2. Neta Package (Node Definitions)

**Purpose**: Foundation types for all nodes
**Dependencies**: None (foundation package)
**File Size Target**: < 150 lines total

#### `pkg/neta/definition.go`
```go
// Package neta defines the core node types for Bento.
// Neta (ネタ) means "ingredients" or "toppings" in sushi terminology.
package neta

import (
    "context"
)

// Definition describes a node that can be executed by Itamae.
// It may be a single executable node or a group containing other nodes.
type Definition struct {
    // Type identifies what kind of node this is (http, transform, group, etc)
    Type string

    // Name is the human-readable identifier for this node
    Name string

    // Parameters contains type-specific configuration
    Parameters map[string]interface{}

    // Nodes contains child nodes (for group types)
    // If empty/nil, this is a leaf node
    Nodes []Definition
}

// IsGroup returns true if this definition contains child nodes.
func (d Definition) IsGroup() bool {
    return len(d.Nodes) > 0
}
```

#### `pkg/neta/executable.go`
```go
package neta

import "context"

// Result represents the outcome of executing a node.
type Result struct {
    // Output contains the result data
    Output interface{}

    // Error contains any execution error
    Error error

    // Metadata contains execution details (duration, etc)
    Metadata map[string]interface{}
}

// Executable is implemented by all node types that can be executed.
// Accept interfaces, return structs.
type Executable interface {
    // Execute runs the node and returns its result
    Execute(ctx context.Context, params map[string]interface{}) (Result, error)
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Type definitions only
- ✅ No execution logic (that's itamae's job)
- ✅ Clear boundaries: Just data structures
- ✅ Files < 100 lines each
- ✅ Functions < 20 lines

### 3. Itamae Package (Orchestrator)

**Purpose**: Execution coordination and orchestration
**Dependencies**: neta
**File Size Target**: < 200 lines total

#### `pkg/itamae/itamae.go`
```go
// Package itamae provides the orchestration engine for executing neta definitions.
// Itamae (板前) means "sushi chef" - the one who prepares each piece.
package itamae

import (
    "context"
    "fmt"

    "bento/pkg/neta"
)

// Itamae orchestrates the execution of neta definitions.
type Itamae struct {
    pantry Registry // Injected dependency
}

// Registry provides node type lookup.
type Registry interface {
    Get(nodeType string) (neta.Executable, error)
}

// New creates a new Itamae with the provided registry.
func New(registry Registry) *Itamae {
    return &Itamae{
        pantry: registry,
    }
}

// Execute runs a neta definition and returns the result.
func (i *Itamae) Execute(ctx context.Context, def neta.Definition) (neta.Result, error) {
    if def.IsGroup() {
        return i.executeGroup(ctx, def)
    }
    return i.executeSingle(ctx, def)
}

// executeSingle runs a single node.
func (i *Itamae) executeSingle(ctx context.Context, def neta.Definition) (neta.Result, error) {
    exec, err := i.pantry.Get(def.Type)
    if err != nil {
        return neta.Result{}, fmt.Errorf("node type not found: %s: %w", def.Type, err)
    }
    return exec.Execute(ctx, def.Parameters)
}

// executeGroup runs a group of nodes in sequence.
func (i *Itamae) executeGroup(ctx context.Context, def neta.Definition) (neta.Result, error) {
    results := make([]neta.Result, 0, len(def.Nodes))
    for _, child := range def.Nodes {
        result, err := i.Execute(ctx, child)
        if err != nil {
            return neta.Result{}, err
        }
        results = append(results, result)
    }
    return neta.Result{Output: results}, nil
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Orchestration only
- ✅ Small functions (< 20 lines each)
- ✅ Clear dependency injection (Registry interface)
- ✅ Composable execution pattern
- ✅ File < 100 lines

### 4. Pantry Package (Registry)

**Purpose**: Node type registry and lookup
**Dependencies**: neta
**File Size Target**: < 150 lines total

#### `pkg/pantry/registry.go`
```go
// Package pantry provides the node type registry.
// Pantry stores all available neta types and provides lookup.
package pantry

import (
    "fmt"
    "sync"

    "bento/pkg/neta"
)

// Pantry is a thread-safe registry of node types.
type Pantry struct {
    mu    sync.RWMutex
    nodes map[string]neta.Executable
}

// New creates a new empty Pantry.
func New() *Pantry {
    return &Pantry{
        nodes: make(map[string]neta.Executable),
    }
}

// Register adds a node type to the pantry.
func (p *Pantry) Register(nodeType string, exec neta.Executable) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    if _, exists := p.nodes[nodeType]; exists {
        return fmt.Errorf("node type already registered: %s", nodeType)
    }

    p.nodes[nodeType] = exec
    return nil
}

// Get retrieves a node type from the pantry.
func (p *Pantry) Get(nodeType string) (neta.Executable, error) {
    p.mu.RLock()
    defer p.mu.RUnlock()

    exec, exists := p.nodes[nodeType]
    if !exists {
        return nil, fmt.Errorf("node type not found: %s", nodeType)
    }

    return exec, nil
}

// List returns all registered node types.
func (p *Pantry) List() []string {
    p.mu.RLock()
    defer p.mu.RUnlock()

    types := make([]string, 0, len(p.nodes))
    for t := range p.nodes {
        types = append(types, t)
    }
    return types
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Registry only
- ✅ Thread-safe implementation
- ✅ Functions < 20 lines
- ✅ Clear API (Register, Get, List)
- ✅ File < 100 lines

### 5. CLI Bootstrap

**Purpose**: Entry point with Cobra framework
**Dependencies**: cobra, neta, itamae, pantry
**File Size Target**: < 100 lines total

#### `cmd/bento/main.go`
```go
package main

import (
    "fmt"
    "os"
)

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

#### `cmd/bento/root.go`
```go
package main

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "bento",
    Short: "🍱 Bento - Organized workflow orchestration",
    Long: `Bento is a Go-based CLI orchestration tool.

Run 'bento' without arguments to launch the interactive TUI.
Or use commands directly: prepare, pack, pantry, taste.

Also available as 'b3o' alias.`,
    Run: func(cmd *cobra.Command, args []string) {
        // Phase 4 will launch TUI here
        // For now, show help
        cmd.Help()
    },
}

func init() {
    // Phase 3 will add subcommands here
}
```

**Bento Box Compliance**:
- ✅ Minimal main.go (< 10 lines)
- ✅ Clear separation of concerns
- ✅ Ready for Phase 3 expansion

### 6. Makefile

```makefile
.PHONY: build test clean fmt lint install

# Build the bento binary
build:
	go build -o bin/bento ./cmd/bento

# Run all tests with race detector
test:
	go test -v -race ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Format all Go files
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Install to GOPATH
install:
	go install ./cmd/bento

# Run all quality checks (Karen's requirements)
check: fmt lint test build
	@echo "✅ All quality checks passed!"
```

### 7. README.md

```markdown
# 🍱 Bento

Organized workflow orchestration in Go.

## Installation

\`\`\`bash
# Build from source
make build

# Or install to GOPATH
make install
\`\`\`

## Usage

\`\`\`bash
# Launch interactive TUI (Phase 4)
bento

# Or use the alias
b3o

# Direct commands (Phase 3)
bento prepare flow.bento.yaml
bento pack flow.bento.yaml
bento pantry list
bento taste flow.bento.yaml
\`\`\`

## Development

\`\`\`bash
# Run tests
make test

# Format code
make fmt

# Lint
make lint

# All checks (required before commit)
make check
\`\`\`

## The Bento Box Principle

This project follows the [Bento Box Principle](./.claude/BENTO_BOX_PRINCIPLE.md):

- 🍙 Single Responsibility
- 🚫 No Utility Grab Bags
- 🔲 Clear Boundaries
- 🧩 Composable
- ✂️ YAGNI

Files < 250 lines. Functions < 20 lines. Zero utils packages.

## License

MIT
```

## Testing Strategy

Each package must have comprehensive tests:

### Neta Tests
- `definition_test.go`: IsGroup() logic
- `executable_test.go`: Interface compliance

### Itamae Tests
- `itamae_test.go`: Single and group execution
- `context_test.go`: Context handling
- Mock registry for testing

### Pantry Tests
- `registry_test.go`: Register/Get/List operations
- Thread safety tests
- Error cases (duplicate registration, not found)

## Validation Commands

Before marking phase complete, run:

```bash
# Format
go fmt ./...

# Lint
golangci-lint run

# Test with race detector
go test -v -race ./...

# Build
go build ./cmd/bento

# Module tidy
go mod tidy

# File size check
find pkg -name "*.go" -exec wc -l {} + | sort -rn | head -10

# Check for utils packages (should be empty)
find pkg -type d -name "*util*"
```

## Success Criteria

Phase 1 is complete when:

1. ✅ Go workspace initialized with 6 modules
2. ✅ Neta package: definition.go + executable.go (< 150 lines total)
3. ✅ Itamae package: orchestrator implementation (< 200 lines total)
4. ✅ Pantry package: registry implementation (< 150 lines total)
5. ✅ CLI: main.go + root.go with Cobra (< 100 lines total)
6. ✅ All tests passing (go test -race ./...)
7. ✅ All files < 250 lines
8. ✅ All functions < 20 lines
9. ✅ Zero utils packages
10. ✅ golangci-lint clean
11. ✅ Builds successfully
12. ✅ **Karen's approval granted**

## Common Pitfalls to Avoid

1. ❌ **Creating utils/ package** - Organize by domain instead
2. ❌ **God objects** - Keep Itamae focused on orchestration
3. ❌ **Premature abstraction** - Simple concrete types first
4. ❌ **Large files** - Extract to focused files if approaching 200 lines
5. ❌ **Complex functions** - Break down if approaching 20 lines
6. ❌ **Mixing concerns** - Execution in neta, validation in pantry, etc.

## Next Phase

After Karen approval, proceed to **[Phase 2: Neta Library](./phase-2-neta-library.md)** to:
- Implement concrete node types (HTTP, transform, conditional, loop)
- Create .bento.yaml examples
- Build out the neta library

## Execution Prompt

```
I'm ready to begin Phase 1: Foundation.

I have read the Bento Box Principle and will follow it.

Please initialize the Go workspace and create the core packages:
- neta (node definitions)
- itamae (orchestrator)
- pantry (registry)
- CLI bootstrap with Cobra

I will use TodoWrite to track progress and get Karen's approval before completing.
```

---

**Phase 1 Foundation**: Core packages and architecture 🍱
