# Phase 3: CLI Commands

**Status**: Pending
**Duration**: 2-3 hours
**Prerequisites**: Phase 2 complete, Karen approved

## Overview

Implement Cobra-based CLI commands for direct workflow operations. Users can either use the TUI (Phase 4) or these direct commands for scripting and automation.

## Pre-Work Checklist

Before starting, you MUST:

1. ✅ Read [BENTO_BOX_PRINCIPLE.md](../BENTO_BOX_PRINCIPLE.md)
2. ✅ Confirm: "I understand the Bento Box Principle and will follow it"
3. ✅ Use TodoWrite to track all tasks
4. ✅ Phase 2 approved by Karen

## Goals

1. Implement 4 core commands (prepare, pack, pantry, taste)
2. Integrate Viper for configuration
3. Add error handling and user feedback
4. Create help documentation
5. Integration tests for each command
6. Validate Bento Box compliance

## Command Structure

```
cmd/bento/
├── main.go              # Entry point
├── root.go              # Root command + TUI launch
├── prepare.go           # Validate .bento.yaml
├── pack.go              # Execute workflow
├── pantry.go            # List/search neta types
├── taste.go             # Dry run
└── config.go            # Viper configuration
```

## Deliverables

### 1. Root Command (Updated)

**File**: `cmd/bento/root.go`
**File Size Target**: < 100 lines

```go
package main

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var (
    cfgFile string
    verbose bool
)

var rootCmd = &cobra.Command{
    Use:   "bento",
    Short: "🍱 Bento - Organized workflow orchestration",
    Long: `Bento is a Go-based CLI orchestration tool.

Run 'bento' without arguments to launch the interactive TUI (Phase 4).
Or use commands directly for scripting and automation.

Available commands:
  prepare - Validate a .bento.yaml file
  pack    - Execute a workflow
  pantry  - List/search available neta types
  taste   - Dry run a workflow

Also available as 'b3o' alias.`,
    Version: "0.1.0",
    Run: func(cmd *cobra.Command, args []string) {
        // Phase 4 will launch TUI here
        fmt.Println("🍱 Bento TUI coming in Phase 4!")
        fmt.Println("For now, use commands: prepare, pack, pantry, taste")
        cmd.Help()
    },
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    cobra.OnInitialize(initConfig)

    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bento.yaml)")
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

    viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, err := os.UserHomeDir()
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error finding home: %v\n", err)
            os.Exit(1)
        }

        viper.AddConfigPath(home)
        viper.SetConfigType("yaml")
        viper.SetConfigName(".bento")
    }

    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err == nil {
        if verbose {
            fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
        }
    }
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Root command setup
- ✅ Functions < 20 lines
- ✅ Clear configuration handling
- ✅ File < 100 lines

### 2. Prepare Command

**Purpose**: Validate .bento.yaml files
**File**: `cmd/bento/prepare.go`
**File Size Target**: < 150 lines

```go
package main

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "gopkg.in/yaml.v3"

    "bento/pkg/neta"
)

var prepareCmd = &cobra.Command{
    Use:   "prepare [file.bento.yaml]",
    Short: "Validate a bento workflow file",
    Long: `Prepare validates a .bento.yaml file without executing it.

This checks:
- YAML syntax is valid
- Node types are registered
- Required parameters are present
- Structure is well-formed`,
    Args: cobra.ExactArgs(1),
    RunE: runPrepare,
}

func init() {
    rootCmd.AddCommand(prepareCmd)
}

func runPrepare(cmd *cobra.Command, args []string) error {
    filename := args[0]

    def, err := loadDefinition(filename)
    if err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    if err := validateDefinition(def); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    fmt.Printf("✅ %s is valid\n", filename)
    printDefinitionSummary(def)
    return nil
}

// loadDefinition reads and parses a .bento.yaml file.
func loadDefinition(filename string) (neta.Definition, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return neta.Definition{}, err
    }

    var def neta.Definition
    if err := yaml.Unmarshal(data, &def); err != nil {
        return neta.Definition{}, err
    }

    return def, nil
}

// validateDefinition checks if a definition is well-formed.
func validateDefinition(def neta.Definition) error {
    if def.Type == "" {
        return fmt.Errorf("type is required")
    }

    if def.IsGroup() {
        return validateGroup(def)
    }

    return validateLeaf(def)
}

// validateGroup validates a group definition.
func validateGroup(def neta.Definition) error {
    if len(def.Nodes) == 0 {
        return fmt.Errorf("group must have child nodes")
    }

    for i, child := range def.Nodes {
        if err := validateDefinition(child); err != nil {
            return fmt.Errorf("node %d: %w", i, err)
        }
    }

    return nil
}

// validateLeaf validates a leaf node definition.
func validateLeaf(def neta.Definition) error {
    // Check if type is registered (Phase 2 pantry integration)
    // For now, just check type is not empty
    if def.Type == "" {
        return fmt.Errorf("type is required")
    }
    return nil
}

// printDefinitionSummary prints a summary of the definition.
func printDefinitionSummary(def neta.Definition) {
    fmt.Printf("\nSummary:\n")
    fmt.Printf("  Type: %s\n", def.Type)
    if def.Name != "" {
        fmt.Printf("  Name: %s\n", def.Name)
    }
    if def.IsGroup() {
        fmt.Printf("  Nodes: %d\n", len(def.Nodes))
    }
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Validation only
- ✅ Functions < 20 lines
- ✅ Clear error messages
- ✅ File < 150 lines

### 3. Pack Command

**Purpose**: Execute workflows
**File**: `cmd/bento/pack.go`
**File Size Target**: < 150 lines

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/spf13/cobra"

    "bento/pkg/itamae"
    "bento/pkg/pantry"
)

var (
    packTimeout time.Duration
    packDryRun  bool
)

var packCmd = &cobra.Command{
    Use:   "pack [file.bento.yaml]",
    Short: "Execute a bento workflow",
    Long: `Pack executes a .bento.yaml workflow file.

This runs all nodes in the workflow and reports results.`,
    Args: cobra.ExactArgs(1),
    RunE: runPack,
}

func init() {
    rootCmd.AddCommand(packCmd)

    packCmd.Flags().DurationVar(&packTimeout, "timeout", 5*time.Minute, "execution timeout")
    packCmd.Flags().BoolVar(&packDryRun, "dry-run", false, "validate without executing")
}

func runPack(cmd *cobra.Command, args []string) error {
    filename := args[0]

    def, err := loadDefinition(filename)
    if err != nil {
        return fmt.Errorf("failed to load: %w", err)
    }

    if packDryRun {
        return runPrepare(cmd, args)
    }

    result, err := executeBento(def, packTimeout)
    if err != nil {
        return fmt.Errorf("execution failed: %w", err)
    }

    printResult(result)
    return nil
}

// executeBento runs the workflow with the itamae orchestrator.
func executeBento(def neta.Definition, timeout time.Duration) (interface{}, error) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    // Initialize pantry and itamae (from Phase 1 & 2)
    p := pantry.New()
    // Register node types here (from Phase 2)

    chef := itamae.New(p)

    result, err := chef.Execute(ctx, def)
    if err != nil {
        return nil, err
    }

    return result.Output, nil
}

// printResult displays execution results.
func printResult(result interface{}) {
    fmt.Println("\n✅ Execution complete!")
    fmt.Printf("\nResult:\n%v\n", result)
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Execution
- ✅ Functions < 20 lines
- ✅ Timeout handling
- ✅ File < 150 lines

### 4. Pantry Command

**Purpose**: List/search neta types
**File**: `cmd/bento/pantry.go`
**File Size Target**: < 100 lines

```go
package main

import (
    "fmt"
    "sort"
    "strings"

    "github.com/spf13/cobra"

    "bento/pkg/pantry"
)

var pantryCmd = &cobra.Command{
    Use:   "pantry [search]",
    Short: "List or search available neta types",
    Long: `Pantry shows all registered node types.

Optionally provide a search term to filter results.`,
    Args: cobra.MaximumNArgs(1),
    RunE: runPantry,
}

func init() {
    rootCmd.AddCommand(pantryCmd)
}

func runPantry(cmd *cobra.Command, args []string) error {
    p := initializePantry()

    types := p.List()
    sort.Strings(types)

    if len(args) > 0 {
        types = filterTypes(types, args[0])
    }

    printTypes(types)
    return nil
}

// filterTypes returns types matching the search term.
func filterTypes(types []string, search string) []string {
    search = strings.ToLower(search)
    filtered := []string{}

    for _, t := range types {
        if strings.Contains(strings.ToLower(t), search) {
            filtered = append(filtered, t)
        }
    }

    return filtered
}

// printTypes displays the type list.
func printTypes(types []string) {
    fmt.Printf("🍱 Available neta types (%d):\n\n", len(types))

    for _, t := range types {
        fmt.Printf("  • %s\n", t)
    }

    if len(types) == 0 {
        fmt.Println("  (no types found)")
    }
}

// initializePantry creates and populates the pantry (from Phase 2).
func initializePantry() *pantry.Pantry {
    p := pantry.New()
    // Register all node types here
    return p
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Type listing
- ✅ Functions < 15 lines
- ✅ Simple filtering
- ✅ File < 100 lines

### 5. Taste Command

**Purpose**: Dry run workflows
**File**: `cmd/bento/taste.go`
**File Size Target**: < 100 lines

```go
package main

import (
    "fmt"

    "github.com/spf13/cobra"
)

var tasteCmd = &cobra.Command{
    Use:   "taste [file.bento.yaml]",
    Short: "Dry run a workflow (alias for prepare)",
    Long: `Taste validates a workflow without executing it.

This is an alias for 'bento prepare' with more verbose output.`,
    Args: cobra.ExactArgs(1),
    RunE: runTaste,
}

func init() {
    rootCmd.AddCommand(tasteCmd)
}

func runTaste(cmd *cobra.Command, args []string) error {
    fmt.Println("🍱 Tasting your bento...\n")

    if err := runPrepare(cmd, args); err != nil {
        fmt.Println("\n❌ This bento doesn't taste right!")
        return err
    }

    fmt.Println("\n✨ Delicious! Ready to pack.")
    return nil
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Dry run
- ✅ Simple wrapper around prepare
- ✅ File < 50 lines

### 6. Configuration

**File**: `$HOME/.bento.yaml`
**Example**:

```yaml
# Bento configuration file
verbose: false
timeout: 5m
pantry:
  # Custom node type locations (future)
  paths: []
```

## Integration Tests

Create integration tests for each command:

**File**: `cmd/bento/commands_test.go`

```go
package main

import (
    "os"
    "path/filepath"
    "testing"
)

func TestPrepareCommand(t *testing.T) {
    // Create temp .bento.yaml file
    tmpDir := t.TempDir()
    file := filepath.Join(tmpDir, "test.bento.yaml")

    content := []byte(`
type: http
name: Test
parameters:
  url: https://example.com
`)
    if err := os.WriteFile(file, content, 0644); err != nil {
        t.Fatal(err)
    }

    // Test prepare command
    rootCmd.SetArgs([]string{"prepare", file})
    if err := rootCmd.Execute(); err != nil {
        t.Errorf("prepare failed: %v", err)
    }
}

// Similar tests for pack, pantry, taste
```

## Validation Commands

Before marking phase complete:

```bash
# Format
go fmt ./...

# Lint
golangci-lint run

# Test
go test -v ./cmd/bento/...

# Build and test all commands
go build ./cmd/bento
./bento prepare examples/http-get.bento.yaml
./bento pantry
./bento taste examples/http-get.bento.yaml

# File size check
find cmd/bento -name "*.go" -exec wc -l {} + | sort -rn
```

## Success Criteria

Phase 3 is complete when:

1. ✅ 4 commands implemented (prepare, pack, pantry, taste)
2. ✅ Viper configuration working
3. ✅ Error handling comprehensive
4. ✅ Help documentation clear
5. ✅ Integration tests passing
6. ✅ All files < 250 lines
7. ✅ All functions < 20 lines
8. ✅ golangci-lint clean
9. ✅ Commands tested manually
10. ✅ **Karen's approval granted**

## Common Pitfalls to Avoid

1. ❌ **God command file** - Keep each command in separate file
2. ❌ **Business logic in commands** - Commands should orchestrate, not implement
3. ❌ **Poor error messages** - User-friendly errors are critical
4. ❌ **No help text** - Every command needs good documentation
5. ❌ **Skipping integration tests** - Test the actual CLI experience

## Next Phase

After Karen approval, proceed to **[Phase 4: Omise TUI](./phase-4-omise-tui.md)** to:
- Build Bubble Tea TUI
- Create interactive workflow browser
- Implement execution viewer with progress
- Style with Lip Gloss

## Execution Prompt

```
I'm ready to begin Phase 3: CLI Commands.

I have read the Bento Box Principle and will follow it.

Please implement Cobra commands:
- prepare (validate .bento.yaml)
- pack (execute workflow)
- pantry (list neta types)
- taste (dry run)

Each command in its own file. I will use TodoWrite to track progress and get Karen's approval before completing.
```

---

**Phase 3 CLI Commands**: Direct command interface 🍱
