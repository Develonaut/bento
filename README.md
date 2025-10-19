# 🍱 Bento

**High-performance workflow automation CLI written in Go**

Bento is a complete rewrite of Atomiton's core execution engine, designed to be fast, portable, and powerful. Build automated workflows using composable "neta" (ingredients) that can be connected together like a carefully crafted bento box.

## Why Bento?

- **🚀 Fast:** 15-30x faster startup than Node.js, 6x less memory
- **📦 Small:** 20MB single binary vs 202MB node_modules
- **🎯 Portable:** Cross-platform (Linux, macOS, Windows)
- **⚡ Powerful:** True parallelism with goroutines
- **🎨 Simple:** JSON-based workflow files (`.bento.json`)

## Commands

```bash
bento savor workflow.bento.json    # Execute a bento workflow
bento sample workflow.bento.json   # Validate a bento file
bento menu                         # List available bentos
bento box my-workflow              # Create a new bento
bento recipe readme                # View docs with glow
```

## Architecture

Bento uses a sushi-themed package architecture:

- **neta** (ネタ) - "Ingredients" - Core workflow node types
- **itamae** (板前) - "Sushi Chef" - Orchestration engine
- **pantry** - Registry of available neta types
- **hangiri** (はんぎり) - "Wooden Rice Tub" - Storage layer
- **shoyu** - "Soy Sauce" - Structured logging
- **omakase** - "Chef's Choice" - Validation

## Workflow Nodes (Neta)

All 10 neta types from Atomiton, ported to Go:

1. **edit-fields** - Field editor with templates
2. **http-request** - HTTP client
3. **file-system** - File operations
4. **shell-command** - Execute shell commands
5. **group** - Sequential/parallel execution
6. **loop** - Iteration (forEach, times, while)
7. **parallel** - Advanced parallelism
8. **spreadsheet** - Excel/CSV processing
9. **image** - Image processing (govips)
10. **transform** - Data transformation (expr)

## Example Workflow

```json
{
  "id": "product-automation",
  "type": "group",
  "version": "1.0.0",
  "name": "Product Photo Automation",
  "nodes": [
    {
      "id": "read-csv",
      "type": "spreadsheet",
      "parameters": {
        "operation": "read",
        "format": "csv",
        "path": "products.csv"
      }
    },
    {
      "id": "process-products",
      "type": "loop",
      "parameters": {
        "mode": "forEach",
        "items": "{{.read-csv.rows}}"
      }
    }
  ]
}
```

## Development Roadmap

See [.claude/README.md](./.claude/README.md) for detailed project overview and roadmap.

Phase-by-phase implementation strategy in [.claude/strategy/](./.claude/strategy/).

## Bento Box Principle

Bento follows the **Bento Box Principle**:

- 🍙 Single Responsibility per file/package
- 🚫 No utility grab bags
- 🔲 Clear boundaries (interfaces)
- 🧩 Composable (small, focused components)
- ✂️ YAGNI (no future-proofing)

See [.claude/BENTO_BOX_PRINCIPLE.md](./.claude/BENTO_BOX_PRINCIPLE.md) for details.

## Documentation

- [Project Overview](./.claude/README.md)
- [Package Naming](./.claude/PACKAGE_NAMING.md)
- [Bento Box Principle](./.claude/BENTO_BOX_PRINCIPLE.md)
- [Approved Emojis](./.claude/EMOJIS.md) 🍱
- [Phase Documents](./.claude/strategy/)

## License

MIT

---

**Status:** 🏗️ In Development
**Goal:** Rock-solid CLI for workflow automation with true parallelism and blazing speed
