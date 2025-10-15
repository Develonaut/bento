# 🍱 Bento Implementation Strategy

## Overview

This directory contains the phased implementation strategy for Bento, a Go-based CLI orchestration tool built on the Bento Box Principle.

## Core Philosophy

All implementation phases MUST follow the **[Bento Box Principle](../BENTO_BOX_PRINCIPLE.md)**:

1. **Single Responsibility** - One thing per file/function
2. **No Utility Grab Bags** - No `utils/` packages
3. **Clear Boundaries** - Clean package interfaces
4. **Composable** - Small functions working together
5. **YAGNI** - No unused code

## Quality Standards

- **Files < 250 lines** (500 max)
- **Functions < 20 lines** (30 max)
- **Zero utils packages**
- All Go quality checks pass (fmt, lint, test -race, build)

## Implementation Phases

### Phase 0: Bootstrap ✅
**Duration**: 1-2 hours
**Status**: Complete
**File**: [phase-0-bootstrap.md](./phase-0-bootstrap.md)

Foundation setup:
- Repository structure
- .claude/ folder with agents and workflow
- BENTO_BOX_PRINCIPLE.md
- Strategy files (this phase)

**Deliverables**:
- ✅ Git repository initialized
- ✅ .gitignore configured
- ✅ 4 agents adapted for Go (Karen, Michael, Guilliman, Voorhees)
- ✅ Workflow checklists created
- ✅ Strategy files documented

---

### Phase 1: Foundation
**Duration**: 2-3 hours
**Status**: Pending
**File**: [phase-1-foundation.md](./phase-1-foundation.md)

Core packages and architecture:
- Go workspace setup (6 modules)
- `pkg/neta/` - Node definition types
- `pkg/itamae/` - Chef/orchestrator core
- `pkg/pantry/` - Node registry
- Bootstrap CLI with Cobra

**Deliverables**:
- [ ] Go workspace initialized
- [ ] Core types defined (Neta, Executable)
- [ ] Itamae orchestrator structure
- [ ] Pantry registry foundation
- [ ] Cobra CLI skeleton
- [ ] All tests passing
- [ ] Karen approval

---

### Phase 2: Neta Library
**Duration**: 3-4 hours
**Status**: Pending
**File**: [phase-2-neta-library.md](./phase-2-neta-library.md)

Node type implementations:
- HTTP neta (GET/POST/PUT/DELETE)
- Transform neta (jq, template)
- Conditional neta (if/switch)
- Loop neta (for/while)
- Group neta (sequence/parallel)

**Deliverables**:
- [ ] 5+ neta types implemented
- [ ] Each neta in focused package
- [ ] Comprehensive tests
- [ ] Example .bento.yaml files
- [ ] Bento Box compliant
- [ ] Karen approval

---

### Phase 3: CLI Commands
**Duration**: 2-3 hours
**Status**: Pending
**File**: [phase-3-cli-commands.md](./phase-3-cli-commands.md)

Cobra command implementation:
- `bento prepare` - Validate .bento.yaml
- `bento pack` - Execute bento file
- `bento pantry` - List/search neta
- `bento taste` - Dry run

**Deliverables**:
- [ ] 4 commands implemented
- [ ] Viper config integration
- [ ] Error handling
- [ ] Help documentation
- [ ] Integration tests
- [ ] Karen approval

---

### Phase 4: Omise (TUI)
**Duration**: 4-5 hours
**Status**: Pending
**File**: [phase-4-omise-tui.md](./phase-4-omise-tui.md)

Bubble Tea TUI - The Shop:
- `bento` (no args) launches TUI
- Flow browser and selector
- Execution viewer with progress
- Pantry explorer
- Settings screen

**Deliverables**:
- [ ] Bubble Tea app structure
- [ ] 5 screen implementations
- [ ] Lip Gloss styling
- [ ] Bubbles components integrated
- [ ] Huh forms for wizards
- [ ] TUI tests
- [ ] Karen approval

---

### Phase 5: Jubako (Storage) ✅
**Duration**: 2-3 hours
**Status**: Complete
**File**: [phase-5-jubako.md](./archive/phase-5-jubako.md)

Storage layer - Stacked Boxes:
- .bento.yaml parsing
- Flow file management
- History/versioning
- Import/export

**Deliverables**:
- ✅ YAML parser
- ✅ File operations
- ✅ History tracking
- ✅ Migration support
- ✅ Storage tests
- ✅ Karen approval

---

### Phase 5.5: Definition Versioning ✅
**Duration**: 1-2 hours
**Status**: Complete
**File**: [phase-5.5-versioning.md](./phase-5.5-versioning.md)

Add versioning to bento definitions:
- Version field in neta.Definition
- Version validation in parser
- Future-proof for breaking changes
- Migration framework foundation

**Deliverables**:
- ✅ Version field added
- ✅ Validation implemented
- ✅ Tests passing
- ✅ Migration tool created
- ✅ CLI commands enforce version validation
- ✅ Karen approval

---

### Phase 5.7: Node and Bento Validation
**Duration**: 3-4 hours
**Status**: Pending
**File**: [phase-5.7-validation.md](./phase-5.7-validation.md)

Add structured validation framework:
- Schema definitions for all node types
- Parameter validation with clear errors
- Integration with Huh forms in editor
- Prevent invalid bento compartments
- Provide field metadata for UI

**Deliverables**:
- [ ] Validation framework implemented
- [ ] Schemas for all node types
- [ ] Parser integration
- [ ] Clear error messages
- [ ] Schema metadata for editor
- [ ] Karen approval

---

### Phase 6: Enhanced Browser & CRUD
**Duration**: 3-4 hours
**Status**: Pending
**File**: [phase-6-enhanced-browser.md](./phase-6-enhanced-browser.md)

Full-featured bento management:
- Keyboard shortcuts (r, e, c, d, n)
- Jubako integration for dynamic lists
- Delete confirmation dialogs
- Copy and delete operations
- Real-time bento discovery

**Deliverables**:
- [ ] Keyboard shortcuts working
- [ ] Jubako integration complete
- [ ] Confirmation dialogs
- [ ] CRUD operations
- [ ] Karen approval

---

### Phase 7: Bento Editor - Node Builder
**Duration**: 5-6 hours
**Status**: Pending
**File**: [phase-7-bento-editor-builder.md](./phase-7-bento-editor-builder.md)

Guided bento creation:
- New Editor screen
- Pantry integration for node selection
- Huh wizards for parameter configuration
- Definition building
- Create and edit modes

**Deliverables**:
- [ ] Editor screen created
- [ ] Create/edit modes
- [ ] Pantry integration
- [ ] Huh wizards
- [ ] Save to Jubako
- [ ] Karen approval

---

### Phase 8: Bento Editor - Visualization
**Duration**: 4-5 hours
**Status**: Pending
**File**: [phase-8-bento-visualization.md](./phase-8-bento-visualization.md)

Visual bento box with navigation:
- Text-based bento box rendering
- Arrow key navigation
- Edit node in-place
- Move/delete nodes
- Run from editor

**Deliverables**:
- [ ] Visual bento box
- [ ] Navigation working
- [ ] Node operations (edit/move/delete)
- [ ] Run from editor
- [ ] Karen approval

---

### Phase 9: Examples & Templates
**Duration**: 2-3 hours
**Status**: Pending
**File**: [phase-9-examples.md](./phase-9-examples.md)

Built-in examples and templates:
- Embedded example .bento.yaml files
- Examples section in browser
- Copy-from-template functionality
- Examples for each node type

**Deliverables**:
- [ ] 7+ example templates
- [ ] go:embed integration
- [ ] Examples browser mode
- [ ] Copy functionality
- [ ] Karen approval

---

### Phase 10: Real-World Proof-of-Concept
**Duration**: 1-2 hours
**Status**: Pending
**File**: [phase-10-proof-of-concept.md](./phase-10-proof-of-concept.md)

Validate system with real use case:
- Create actual bento for user's workflow
- Test editor/flow with real requirements
- Validate all phases work together
- Final system validation

**Deliverables**:
- [ ] User's bento created
- [ ] Editor validated
- [ ] System integration verified
- [ ] Documentation complete

---

## Timeline

**Total Estimated Duration**: 31-42 hours

```
Week 1-2 (Foundation):
├── Phase 0: Bootstrap          [Complete] ✅
├── Phase 1: Foundation         [Complete] ✅
├── Phase 2: Neta Library       [Complete] ✅
├── Phase 3: CLI Commands       [Complete] ✅
├── Phase 4: Omise TUI          [Complete] ✅
└── Phase 5: Jubako Storage     [Complete] ✅

Week 3 (Enhancement):
├── Phase 5.5: Versioning       [Complete] ✅
├── Phase 5.7: Validation       [3-4 hours]
├── Phase 6: Enhanced Browser   [3-4 hours]
├── Phase 7: Editor Builder     [5-6 hours]
├── Phase 8: Visualization      [4-5 hours]
└── Phase 9: Examples           [2-3 hours]

Week 4 (Validation):
└── Phase 10: Proof-of-Concept  [1-2 hours]
```

## Agent Responsibilities

### 🏛️ Michael (Architecture)
- Phase structure validation
- Package boundary enforcement
- Dependency graph reviews
- Domain model integrity

### 📏 Guilliman (Go Standards)
- Idiomatic Go enforcement
- Standard library first
- Context and error patterns
- Module boundary guardian

### 👮‍♀️ Karen (Validation)
- **Bento Box Principle enforcer** (PRIMARY ROLE)
- Quality gate validation
- Test and build verification
- Phase completion approval
- **ZERO TOLERANCE** for violations

### 🔪 Voorhees (Simplicity)
- Complexity slashing
- Over-abstraction removal
- File/function size enforcement
- YAGNI principle guardian

## Workflow

Each phase follows this workflow:

1. **Read** phase strategy file
2. **Confirm** understanding of Bento Box Principle
3. **Plan** with TodoWrite tool
4. **Implement** following Go standards
5. **Clean up** old/deprecated code:
   - Remove unused views or components
   - Delete deprecated functions
   - Clean up commented-out code
   - Remove obsolete tests
6. **Validate** with quality checks:
   - Run `go fmt ./...`
   - Run `golangci-lint run`
   - Run `go test -race ./...`
   - Check file sizes (< 250 lines)
   - Check function sizes (< 20 lines)
7. **Code Review** - Run `/code-review` command
8. **Get approval** from Karen before proceeding

## Critical Rules

### ✅ DO
- Follow Bento Box Principle religiously
- Keep files under 250 lines
- Keep functions under 20 lines
- Organize by domain (no utils/)
- Use TodoWrite for tracking
- Clean up old/deprecated code as you go
- Remove unused views and components
- Run `/code-review` after implementation
- Get Karen's approval before completing phase

### ❌ DON'T
- Create utils/ packages (INSTANT REJECTION)
- Exceed file/function limits
- Leave unused/deprecated code behind
- Skip validation steps
- Skip `/code-review` command
- Proceed without Karen approval
- Ignore Bento Box violations

## Success Criteria

Each phase is complete when:

1. ✅ All deliverables implemented
2. ✅ `go fmt ./...` - Clean
3. ✅ `golangci-lint run` - Zero warnings
4. ✅ `go test -race ./...` - All pass
5. ✅ `go build ./cmd/bento` - Success
6. ✅ Files < 250 lines
7. ✅ Functions < 20 lines
8. ✅ Zero utils packages
9. ✅ **`/code-review` command run** - Peer review complete
10. ✅ **Karen's approval granted**

## Getting Started

To begin Phase 1:

```bash
# Read the phase strategy
cat .claude/strategy/phase-1-foundation.md

# Confirm understanding
echo "I understand the Bento Box Principle and will follow it"

# Begin implementation
# (Follow phase-specific prompts)
```

---

**Remember**: The Bento Box Principle is not optional. Karen enforces it with zero tolerance. 🍱
