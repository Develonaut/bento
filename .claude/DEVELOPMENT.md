# Bento Development Guide

## Development Setup

### Prerequisites
- Go 1.21 or later
- golangci-lint for code quality checks

### Building
```bash
go build -o bento ./cmd/bento
```

### Running Tests
```bash
# Run all tests with race detector
go test ./... -race

# Run specific package tests
go test ./pkg/jubako/...

# Run with verbose output
go test -v ./...
```

### Code Quality
```bash
# Run linter
golangci-lint run

# Format code
go fmt ./...
```

## Debugging the TUI

### Enable Debug Mode

Bento includes a debug mode that logs all Bubble Tea messages to a file for debugging:

```bash
DEBUG=1 ./bento
```

Or when running from source:

```bash
DEBUG=1 go run ./cmd/bento
```

### Viewing Debug Output

Debug messages are written to `/tmp/bento-debug.log`:

```bash
# Tail the debug log in a separate terminal
tail -f /tmp/bento-debug.log
```

The debug log includes:
- All Bubble Tea messages received by the Update() method
- Message types and full content (using spew for deep inspection)
- Context labels for each message group

### Debug Functions

Available debug functions in `pkg/omise/screens/shared/debug.go`:

```go
// Initialize debug mode (automatically called by Launch if DEBUG env var is set)
shared.InitDebug()

// Close debug file (automatically called on exit)
shared.CloseDebug()

// Log a message with context
shared.DebugMsg(msg, "Context Label")

// Log formatted message
shared.DebugPrintf("Format string: %v", value)

// Check if debug mode is enabled
if shared.IsDebugMode() {
    // Debug-only code
}
```

## Development Workflow

### Live Development

Currently, there is no automated live reload configured. To test changes:

1. Build the binary:
   ```bash
   go build -o bento ./cmd/bento
   ```

2. Run with debug mode enabled:
   ```bash
   DEBUG=1 ./bento
   ```

3. In another terminal, watch debug output:
   ```bash
   tail -f /tmp/bento-debug.log
   ```

4. Make changes and rebuild to test

### Testing TUI Components

For testing Bubble Tea components, consider using:
- [teatest](https://github.com/charmbracelet/teatest) - E2E testing framework for Bubbletea
- Manual testing with debug mode enabled

## Code Standards

### The Bento Box Principle

All code should follow the **Bento Box Principle** (see `.claude/BENTO_BOX_PRINCIPLE.md`):

1. **Files ≤ 250 lines** - Split larger files into focused components
2. **Functions ≤ 30 lines** - Extract complex logic into helper functions
3. **No panic() in production** - Use error returns and panic recovery
4. **100% passing tests** - Run with `-race` detector

### Panic Recovery

Critical commands should use panic recovery to ensure the terminal is always restored:

```go
import "bento/pkg/omise/screens/shared"

func myCommand() tea.Cmd {
    return shared.RecoverFromPanic(func() tea.Msg {
        // Command logic here
        return MyMsg{}
    }, "myCommand")
}
```

Or use defer/recover in background goroutines:

```go
func backgroundTask() {
    defer func() {
        if r := recover(); r != nil {
            err := fmt.Errorf("panic: %v", r)
            // Send error message to TUI
        }
    }()
    // Task logic here
}
```

### Layout Calculations

Use the centralized `LayoutHelper` to avoid error-prone manual arithmetic:

```go
import "bento/pkg/omise/screens/shared"

layout := shared.NewLayoutHelper(width, height)

// Calculate content dimensions
contentHeight := layout.ContentHeight(headerHeight, footerHeight)
contentWidth := layout.ContentWidth(leftMargin, rightMargin)

// Split screen
leftWidth, rightWidth := layout.SplitHorizontal(0.5) // 50/50 split
topHeight, bottomHeight := layout.SplitVertical(0.3) // 30/70 split

// Calculate remaining space after rendering components
remaining := layout.RemainingHeight(header, footer, otherComponent)

// Utility functions
clamped := shared.Clamp(value, min, max)
larger := shared.MaxInt(a, b)
smaller := shared.MinInt(a, b)
```

## Architecture

### Package Structure

```
bento/
├── cmd/bento/          # CLI entry point
├── pkg/
│   ├── itamae/        # Execution engine (chef)
│   ├── jubako/        # Bento storage (box)
│   ├── neta/          # Plugin system (ingredients)
│   ├── omise/         # TUI application (shop)
│   │   ├── screens/   # Screen implementations
│   │   │   ├── browser/        # Bento browser
│   │   │   ├── executor/       # Bento executor
│   │   │   ├── guided_creation/# Guided creation modal
│   │   │   ├── help/          # Help screen
│   │   │   ├── pantry/        # Plugin registry
│   │   │   ├── settings/      # Settings screen
│   │   │   └── shared/        # Shared utilities
│   │   ├── config/    # Configuration management
│   │   └── styles/    # Theme and styling
│   └── pantry/        # Plugin registry
```

### Japanese Package Names

Bento uses Japanese culinary terms for package names:

- **itamae** (板前) - "chef" - Executes bentos
- **jubako** (重箱) - "stacked boxes" - Stores bentos
- **neta** (ネタ) - "ingredient" - Individual plugins
- **omise** (お店) - "shop" - Customer-facing TUI
- **pantry** - Registry of available netas

## Making Changes

### Before Committing

1. **Run tests**:
   ```bash
   go test ./... -race
   ```

2. **Run linter**:
   ```bash
   golangci-lint run
   ```

3. **Format code**:
   ```bash
   go fmt ./...
   ```

4. **Check file sizes** - Keep files under 250 lines
5. **Check function sizes** - Keep functions under 30 lines
6. **Review checklist** - See `.claude/workflow/MANDATORY_CHECKLIST.md`

### Commit Messages

Follow conventional commit format:

```
feat: add new feature
fix: fix bug in component
refactor: improve code structure
test: add or update tests
docs: update documentation
```

## Getting Help

- Review `.claude/` documentation for architecture and patterns
- Check phase documents in `.claude/strategy/` for implementation details
- Review agent guidelines in `.claude/agents/` for code review criteria
