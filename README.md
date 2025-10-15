# 🍱 Bento

**Organized workflow orchestration in Go**

Bento is a lightweight CLI tool for defining and executing workflows using a simple YAML-based syntax. Like a traditional bento box with its organized compartments, Bento keeps your automation workflows clean, focused, and composable.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-1.23+-blue.svg)](https://golang.org/dl/)

## ✨ Features

- 🍙 **Simple YAML Syntax** - Define workflows in `.bento.yaml` files
- 👨‍🍳 **Smart Orchestration** - Sequential and parallel execution
- 🏮 **Interactive TUI** - Beautiful terminal UI powered by Bubble Tea
- 📦 **Extensible** - Plugin system for custom node types
- 🧪 **Well-Tested** - Comprehensive test coverage
- 🔍 **Type-Safe** - Built with Go's strong typing

## 🚀 Quick Start

### Installation

```bash
# Install from source
git clone https://github.com/Develonaut/bento.git
cd bento
make install

# Or use the alias
b3o --version
```

### Your First Workflow

Create a `hello.bento.yaml` file:

```yaml
type: http
name: Fetch User Data
parameters:
  method: GET
  url: https://api.github.com/users/octocat
```

Run it:

```bash
# Interactive TUI (recommended)
bento

# Or use direct commands
bento pack hello.bento.yaml

# Dry run / validation
bento prepare hello.bento.yaml
bento taste hello.bento.yaml
```

## 📦 Available Commands

### Interactive Mode

```bash
bento        # Launch TUI (default)
b3o          # Alias for bento
```

The TUI provides:
- 📂 Workflow browser
- ▶️  Execution viewer with progress
- 🏪 Node type explorer
- ⚙️  Settings
- ❓ Help

### Direct Commands

```bash
bento prepare <file>    # Validate a workflow
bento pack <file>       # Execute a workflow
bento pantry [search]   # List/search available node types
bento taste <file>      # Dry run (alias for prepare)
```

## 🍙 Node Types

Bento comes with several built-in node types:

### Network
- **http** - HTTP requests (GET, POST, PUT, DELETE)

### Data Transformation
- **transform.jq** - JQ transformations
- **transform.template** - Template rendering

### Control Flow
- **conditional.if** - If/else logic
- **conditional.switch** - Switch/case logic
- **loop.for** - For loop iteration
- **loop.while** - While loop iteration

### Grouping
- **group.sequence** - Sequential execution
- **group.parallel** - Parallel execution

See [`pkg/README.md`](./pkg/README.md) for package details.

## 🏗️ Architecture

Bento follows the **Bento Box Principle** - organized compartments with clear boundaries:

```
┌─────────────────────────────────────┐
│  Omise (お店) - TUI Shop            │  Interactive interface
├─────────────────────────────────────┤
│  Jubako (重箱) - Storage            │  File management
├─────────────────────────────────────┤
│  Itamae (板前) - Chef               │  Orchestration
├─────────────────────────────────────┤
│  Pantry - Registry                  │  Node type lookup
├─────────────────────────────────────┤
│  Neta (ネタ) - Ingredients          │  Node definitions
└─────────────────────────────────────┘
```

Each layer has a single, well-defined responsibility. See [BENTO_BOX_PRINCIPLE.md](./.claude/BENTO_BOX_PRINCIPLE.md) for details.

## 🛠️ Development

### Prerequisites

- Go 1.23 or higher
- golangci-lint (for linting)

### Building

```bash
# Build
make build

# Run tests
make test

# Format & lint
make fmt
make lint

# All quality checks
make check
```

### Code Quality

This project follows strict quality standards:

- **Files < 250 lines** (500 max)
- **Functions < 20 lines** (30 max)
- **Zero utils packages**
- **100% test coverage on critical paths**

Run `/code-review` (requires Claude Code) for comprehensive code review.

## 🍱 The Name

**Bento** (弁当) - A traditional Japanese meal in a box, with organized compartments.

**b3o** - The alias (like "a6n" for "atomiton"), representing:
- **b** = bento
- **3** = three key principles (organized, composable, simple)
- **o** = orchestration

The package names use sushi/bento terminology:
- **Neta** (ネタ) - Ingredients/toppings
- **Itamae** (板前) - Sushi chef
- **Jubako** (重箱) - Stacked boxes
- **Omise** (お店) - Shop/restaurant

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

Built with love using:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration

## 🤝 Contributing

Contributions welcome! Please ensure:

1. All code follows the Bento Box Principle
2. Files stay under 250 lines
3. Functions stay under 20 lines
4. Tests pass: `make check`
5. No utils packages

See [`.claude/workflow/`](./.claude/workflow/) for detailed guidelines.

---

**Made with 🍱 by the Bento team**
