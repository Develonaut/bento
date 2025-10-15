# Phase 0: Bootstrap

**Status**: ✅ Complete
**Duration**: 1-2 hours
**Prerequisites**: None

## Overview

Establish the foundational repository structure, documentation, and agent configurations for the Bento project. This phase creates the .claude/ directory with all necessary agents, workflows, and strategy files to guide development.

## Goals

1. Create bento repository structure
2. Initialize git with appropriate .gitignore
3. Establish Bento Box Principle as core philosophy
4. Adapt 4 agents from Atomiton for Go development
5. Create workflow checklists for quality enforcement
6. Document phased implementation strategy

## Deliverables

### ✅ Repository Structure
```
/Users/Ryan/Code/bento/
├── .git/
├── .gitignore
└── .claude/
    ├── BENTO_BOX_PRINCIPLE.md
    ├── CHARM_STACK_GUIDE.md
    ├── GO_CLI_FEASIBILITY_ANALYSIS.md
    ├── agents/
    │   ├── Guilliman.md (Go Standards Guardian)
    │   ├── Karen.md (Bento Box Enforcer)
    │   ├── Michael.md (Architecture Agent)
    │   └── Voorhees.md (Complexity Slasher)
    ├── workflow/
    │   ├── MANDATORY_CHECKLIST.md
    │   ├── ENFORCEMENT_CHECKLIST.md
    │   └── REVIEW_CHECKLIST.md
    └── strategy/
        ├── README.md
        ├── phase-0-bootstrap.md (this file)
        ├── phase-1-foundation.md
        ├── phase-2-neta-library.md
        ├── phase-3-cli-commands.md
        ├── phase-4-omise-tui.md
        └── phase-5-jubako.md
```

### ✅ Core Philosophy Document

**BENTO_BOX_PRINCIPLE.md** establishes 5 core principles:

1. **Single Responsibility** - One thing per file/function
2. **No Utility Grab Bags** - No utils/ packages
3. **Clear Boundaries** - Clean package interfaces
4. **Composable** - Small functions working together
5. **YAGNI** - No unused code

**Quality Standards**:
- Files < 250 lines (500 max)
- Functions < 20 lines (30 max)
- Zero utils packages

### ✅ Agent Team

#### 🏛️ Michael (Architecture Agent)
- **Role**: Architecture integrity guardian
- **Unchanged from Atomiton**: Language-agnostic principles
- **Responsibilities**: Domain boundaries, package structure, dependency management

#### 📏 Guilliman (Go Standards Guardian)
- **Role**: Idiomatic Go enforcement
- **Complete rewrite from TypeScript version**
- **Responsibilities**: Go idioms, standard library usage, context/error patterns
- **Key additions**: Go Proverbs enforcement, module boundary protection

#### 👮‍♀️ Karen (Bento Box Enforcer)
- **Role**: Task completion validator and quality gate
- **Adapted for Go**: Validation commands changed to Go tools
- **Responsibilities**: Bento Box compliance, test/build verification, zero tolerance enforcement
- **Critical additions**: File/function size limits, utils package detection

#### 🔪 Voorhees (Complexity Slasher)
- **Role**: Simplicity enforcer
- **Adapted for Go**: Anti-patterns updated for Go idioms
- **Responsibilities**: Over-abstraction removal, nesting reduction, YAGNI guardian
- **Key focus**: Unnecessary interfaces, premature generics, wrapper functions

### ✅ Workflow Checklists

**MANDATORY_CHECKLIST.md**: Pre-work requirements
- Must read Bento Box Principle before starting
- Verbal confirmation required
- TodoWrite tracking mandatory

**ENFORCEMENT_CHECKLIST.md**: Karen's enforcement rules
- Auto-reject conditions (utils/, file size, function size)
- Zero tolerance policy
- Enforcement actions

**REVIEW_CHECKLIST.md**: Comprehensive validation checklist
- Go quality checks (fmt, lint, test, race, build)
- Bento Box compliance (CRITICAL section)
- Testing requirements
- Documentation standards
- Security and performance checks
- Module hygiene

### ✅ Strategy Documentation

Seven strategy files documenting phased implementation:
- README.md (this overview)
- phase-0-bootstrap.md (completed)
- phase-1-foundation.md (core packages)
- phase-2-neta-library.md (node implementations)
- phase-3-cli-commands.md (Cobra CLI)
- phase-4-omise-tui.md (Bubble Tea TUI)
- phase-5-jubako.md (storage layer)

## Key Decisions

### Naming & Theming

**Project**: Bento 🍱
- Metaphor: Organized compartments (nodes) in a box (flow)
- Command alias: `b3o` (like a6n from atomiton)
- File extension: `.bento.yaml`

**Package Names** (sushi/bento theme):
- `neta` (ネタ) - Ingredients/toppings - Node definitions
- `itamae` (板前) - Sushi chef - Orchestrator/conductor
- `pantry` - Western metaphor - Node registry
- `jubako` (重箱) - Stacked boxes - Storage layer
- `omise` (お店) - Shop - TUI/customer interaction

### TUI-First Design

- `bento` (no args) launches interactive TUI
- Commands available for direct usage: prepare, pack, pantry, taste
- Inspired by `gh` command behavior

### Technology Stack

**CLI Framework**: Cobra + Viper
**TUI Stack**: Bubble Tea + Lip Gloss + Bubbles + Huh
**Configuration**: gopkg.in/yaml.v3
**Monorepo**: Go workspace mode (native, no Turborepo)

## Validation

This phase is complete when:

1. ✅ Repository created at `/Users/Ryan/Code/bento/`
2. ✅ Git initialized with .gitignore
3. ✅ BENTO_BOX_PRINCIPLE.md created
4. ✅ 4 agents adapted and documented
5. ✅ 3 workflow checklists created
6. ✅ 7 strategy files created
7. ✅ All documentation reviewed and approved

## Success Criteria

- [x] Repository structure in place
- [x] Core philosophy documented
- [x] Agent team configured
- [x] Workflow enforcement ready
- [x] Strategy documented
- [x] Ready for Phase 1

## Next Phase

Proceed to **[Phase 1: Foundation](./phase-1-foundation.md)** to:
- Initialize Go workspace
- Create core packages (neta, itamae, pantry)
- Bootstrap CLI with Cobra
- Implement foundational types

## Execution Prompt for Phase 1

```
I'm ready to begin Phase 1: Foundation.

Please review .claude/strategy/phase-1-foundation.md and begin implementation.

I understand the Bento Box Principle and will follow it.

Let's create the Go workspace and core packages.
```

---

**Phase 0 Complete** ✅

All foundational documentation and agent configurations are in place. The Bento Box Principle is established as the non-negotiable code quality standard. Karen is ready to enforce with zero tolerance. 🍱
