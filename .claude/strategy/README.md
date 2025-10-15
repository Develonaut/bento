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

### Phase 5: Jubako (Storage)
**Duration**: 2-3 hours
**Status**: Pending
**File**: [phase-5-jubako.md](./phase-5-jubako.md)

Storage layer - Stacked Boxes:
- .bento.yaml parsing
- Flow file management
- History/versioning
- Import/export

**Deliverables**:
- [ ] YAML parser
- [ ] File operations
- [ ] History tracking
- [ ] Migration support
- [ ] Storage tests
- [ ] Karen approval

---

## Timeline

**Total Estimated Duration**: 14-20 hours

```
Week 1:
├── Phase 0: Bootstrap          [Complete] ✅
├── Phase 1: Foundation         [2-3 hours]
└── Phase 2: Neta Library       [3-4 hours]

Week 2:
├── Phase 3: CLI Commands       [2-3 hours]
├── Phase 4: Omise TUI          [4-5 hours]
└── Phase 5: Jubako Storage     [2-3 hours]
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
5. **Validate** with Karen's checklist
6. **Get approval** from Karen before proceeding

## Critical Rules

### ✅ DO
- Follow Bento Box Principle religiously
- Keep files under 250 lines
- Keep functions under 20 lines
- Organize by domain (no utils/)
- Use TodoWrite for tracking
- Get Karen's approval before completing phase

### ❌ DON'T
- Create utils/ packages (INSTANT REJECTION)
- Exceed file/function limits
- Skip validation steps
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
9. ✅ **Karen's approval granted**

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
