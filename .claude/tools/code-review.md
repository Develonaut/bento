---
description: Comprehensive code review coordinating Guilliman (Go standards), Karen (Bento Box compliance), and Voorhees (complexity slashing) using all workflow checklists
---

I need you to perform a comprehensive code review of our Go code.

## Review Team

This review coordinates three specialized agents:

1. **Guilliman** - Go Standards Guardian (`.claude/agents/Guilliman.md`)
2. **Karen** - Bento Box Enforcer & Quality Gate (`.claude/agents/Karen.md`)
3. **Voorhees** - Complexity Slasher (`.claude/agents/Voorhees.md`)

## Required Reading (ALL Reviewers)

Before starting, EVERY reviewer MUST read:

### Core Principles
- **`.claude/BENTO_BOX_PRINCIPLE.md`** - The foundational philosophy

### Workflow Checklists (CRITICAL)
- **`.claude/workflow/MANDATORY_CHECKLIST.md`** - Pre-work requirements
- **`.claude/workflow/ENFORCEMENT_CHECKLIST.md`** - Karen's enforcement rules
- **`.claude/workflow/REVIEW_CHECKLIST.md`** - Comprehensive review criteria

### Agent Responsibilities
- Review your own agent file in `.claude/agents/`
- Understand your specific role and triggers

## Review Process

### Phase 1: Guilliman's Go Standards Review

**Focus**: Idiomatic Go and best practices

**Reference**: [REVIEW_CHECKLIST.md](../.claude/workflow/REVIEW_CHECKLIST.md) - Go Quality & Code Quality sections

#### Check:

✅ **Go Idioms**:
- `context.Context` as first parameter where applicable
- `error` as last return value
- No empty `interface{}` without justification
- Accept interfaces, return structs
- Standard library preferred over dependencies

✅ **Error Handling**:
- All errors checked (no ignored errors)
- Proper error wrapping with `fmt.Errorf` and `%w`
- No panics in production code

✅ **Code Quality**:
- Proper godoc comments (not jsdoc style)
- No redundant comments (remove obvious ones)
- No commented-out code (use git history)
- Clear, descriptive naming

✅ **Module Hygiene**:
- `go mod tidy` run
- Dependencies justified
- No circular imports
- Clean dependency graph

**Commands to Run**:
```bash
go fmt ./...
golangci-lint run
go test -v -race ./...
go build ./cmd/bento
go mod tidy
```

**Report Format**:
```
## 📏 Guilliman's Go Standards Review

### Build & Quality Status
✅/❌ Format: [go fmt output]
✅/❌ Lint: [golangci-lint results with warning count]
✅/❌ Tests: [test results - count passing/failing]
✅/❌ Race: [race detector results]
✅/❌ Build: [build status]

### Go Idiom Compliance
- [ ] Context patterns correct
- [ ] Error handling complete
- [ ] Interface usage appropriate
- [ ] Standard library preferred
- [ ] Naming conventions followed

### Code Quality Issues
[List specific violations with file:line references]

### Module Health
- [ ] No circular dependencies
- [ ] Dependencies justified
- [ ] go mod tidy clean

### Recommendations
[Specific improvements needed]

### Approval
[APPROVED / CHANGES REQUIRED]
```

---

### Phase 2: Voorhees' Complexity Review

**Focus**: Slash complexity and over-engineering

**Reference**: [Voorhees.md](../.claude/agents/Voorhees.md) - Code Smells & Anti-Patterns sections

#### Hunt For:

🔪 **Over-Abstraction**:
- Unnecessary interfaces (single implementation)
- Abstract factories for simple types
- Builders for simple structs
- Premature generics

🔪 **Complexity**:
- Functions > 20 lines (target) / > 30 lines (VIOLATION)
- Deep nesting (> 3 levels)
- God objects/structs
- Massive functions doing everything

🔪 **Wrapper Hell**:
- Pointless wrapper functions
- Functions that just call one other function
- Unnecessary indirection

🔪 **Bad Patterns**:
- Manager/Handler/Service suffixes without reason
- Configuration for things that never change
- Premature optimization
- Enterprise patterns in simple apps

**Commands to Run**:
```bash
# Find large functions (manual inspection needed)
grep -n "^func " pkg/**/*.go cmd/**/*.go | while read line; do
  # Check function length
done

# Check for deeply nested code
# Check for wrapper functions
# Check for unnecessary abstractions
```

**Report Format**:
```
## 🔪 Voorhees' Complexity Review

### Complexity Issues

#### Over-Abstraction Found
[List unnecessary interfaces, factories, builders with file:line]

#### Large Functions (> 20 lines)
[List functions needing breakdown with line counts]

#### God Functions (> 30 lines) - VIOLATIONS
[List critical violations with line counts]

#### Deep Nesting (> 3 levels)
[List files with excessive nesting]

#### Wrapper Functions
[List pointless wrappers that should be inlined or removed]

#### Anti-Patterns
[List Manager/Handler misuse, premature optimization, etc.]

### Recommendations
[Specific refactoring suggestions]

### Slash Count
[Number of items that should be deleted or simplified]

### Approval
[APPROVED / CHANGES REQUIRED]
```

---

### Phase 3: Karen's Bento Box Compliance Review

**Focus**: Bento Box Principle enforcement with ZERO TOLERANCE

**Reference**:
- [REVIEW_CHECKLIST.md](../.claude/workflow/REVIEW_CHECKLIST.md) - Bento Box Compliance section
- [ENFORCEMENT_CHECKLIST.md](../.claude/workflow/ENFORCEMENT_CHECKLIST.md) - Red Flags section
- [Karen.md](../.claude/agents/Karen.md) - Validation Criteria

#### Critical Checks:

🍱 **Bento Box Principles** (from REVIEW_CHECKLIST.md):
- [ ] **Single Responsibility** - Each file/function does ONE thing
- [ ] **No Utility Grab Bags** - No `utils/` or `helpers/` dumping grounds
- [ ] **Clear Boundaries** - Well-defined package interfaces
- [ ] **Composable** - Small functions working together
- [ ] **YAGNI** - No unused code, no premature features
- [ ] **No circular dependencies** - Clean import graph

🍱 **Size Limits** (from MANDATORY_CHECKLIST.md):
- [ ] **Files < 250 lines** - Target, 500 absolute max
- [ ] **Functions < 20 lines** - Target, 30 absolute max
- [ ] **No utils packages** - INSTANT REJECTION

🍱 **Package Organization**:
- [ ] Single purpose packages
- [ ] Public API minimal
- [ ] Internal/ for private code
- [ ] No mixed responsibilities

**Commands to Run**:
```bash
# All quality checks
go fmt ./...
golangci-lint run
go test -v -race ./...
go build ./cmd/bento
go mod tidy

# File size check
find pkg cmd -name "*.go" -exec wc -l {} + | sort -rn | head -20

# Utils package check (MUST BE EMPTY)
find . -type d -name "*util*" -o -name "*helper*"

# Function size check (manual)
# Scan all .go files for functions > 20 lines
```

**Report Format**:
```
## 👮‍♀️ Karen's Bento Box Compliance Review

### Critical Quality Gates (from REVIEW_CHECKLIST.md)
✅/❌ Format: [go fmt status]
✅/❌ Lint: [golangci-lint - ZERO warnings required]
✅/❌ Tests: [All X tests passing]
✅/❌ Race: [No race conditions]
✅/❌ Build: [Build successful]
✅/❌ Module: [go mod tidy clean]

### 🍱 Bento Box Principle Compliance

#### Single Responsibility
✅/❌ [Assessment with violations if any]

#### No Utility Grab Bags
✅/❌ Utils packages found: [NONE or list with INSTANT REJECTION]

#### Clear Boundaries
✅/❌ [Package interface assessment]

#### Composable
✅/❌ [Function composition assessment]

#### YAGNI
✅/❌ [Unused code check]

### Size Limits Analysis

#### Files > 250 lines (target) or > 500 lines (VIOLATION)
[List files with line counts - flag violations]

#### Functions > 20 lines (target) or > 30 lines (VIOLATION)
[List functions with line counts - flag violations]

### Red Flags (from ENFORCEMENT_CHECKLIST.md)

#### Auto-Reject Conditions
- [ ] Utils package detected - [YES/NO]
- [ ] Files over 500 lines - [YES/NO]
- [ ] Functions over 30 lines - [YES/NO]
- [ ] Mixed responsibilities - [YES/NO]
- [ ] Unused code present - [YES/NO]

### Package Organization
- [ ] Single purpose packages
- [ ] No utils packages
- [ ] Public API minimal
- [ ] Internal/ usage correct

### Overall Compliance
[COMPLIANT / VIOLATIONS FOUND]

### Approval Status
✅ **APPROVED** - All Bento Box principles followed
❌ **REJECTED** - [Specific violations requiring fixes]

### Enforcement Action
[If rejected: List required changes before approval]
```

---

## Phase 4: Combined Final Report

After all three reviews complete, provide consolidated summary:

```
# 🍱 Bento Code Review - Final Report

## Review Summary

### 📏 Guilliman (Go Standards)
Status: [APPROVED / CHANGES REQUIRED]
Key Issues: [Count and summary]

### 🔪 Voorhees (Complexity)
Status: [APPROVED / CHANGES REQUIRED]
Items to Slash: [Count]

### 👮‍♀️ Karen (Bento Box Compliance)
Status: [APPROVED / REJECTED]
Violations: [Count]

## Overall Verdict

[✅ APPROVED / ❌ CHANGES REQUIRED / ❌ REJECTED]

### Must Fix Before Approval
1. [Numbered list of critical issues from all reviewers]
2. ...

### Recommended Improvements
1. [Numbered list of suggestions]
2. ...

### Compliance Score
- Go Standards: [X/Y checks passed]
- Complexity: [X items need slashing]
- Bento Box: [COMPLIANT / X violations]

## Commit Readiness

[YES / NO]

**All three reviewers must approve before code can be committed.**

## Next Steps

[If approved:]
✅ Code is ready for commit
✅ All quality gates passed
✅ Bento Box compliance verified
⚠️  **DO NOT commit without explicit user permission**

[If changes required:]
❌ Address issues listed above
❌ Re-run /code-review after fixes
❌ All three approvals required
```

---

## Critical Rules

### From MANDATORY_CHECKLIST.md
1. All work must follow Bento Box Principle
2. No utils packages - INSTANT REJECTION
3. Files < 250 lines (500 max)
4. Functions < 20 lines (30 max)

### From ENFORCEMENT_CHECKLIST.md
1. **NO EXCEPTIONS** - Bento Box Principle is mandatory
2. Auto-reject conditions (utils/, file size, function size, mixed responsibilities)
3. Karen enforces with ZERO TOLERANCE

### From REVIEW_CHECKLIST.md
1. All critical items must pass
2. All Bento Box principles must be followed
3. Tests must pass (including race detector)
4. Build must succeed
5. No security issues

## Remember

- **Three reviewers, three perspectives**: Go quality, Complexity, Bento Box
- **Use the checklists**: They define the standards
- **Zero tolerance**: Violations = rejection
- **All must approve**: Need unanimous approval
- **NO COMMIT**: Without explicit user permission

---

**Karen enforces these rules with ZERO TOLERANCE.** 🍱
