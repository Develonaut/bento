# Phase 9: Examples & Templates

**Status**: Pending
**Duration**: 2-3 hours
**Prerequisites**: Phase 8 complete, Karen approved

## Overview

Add built-in example bentos embedded in the binary to give users starting points and demonstrate capabilities. Create a separate Examples section in the browser with read-only examples that users can copy to create their own bentos.

## Pre-Work Checklist

Before starting, you MUST:

1. ✅ Read [BENTO_BOX_PRINCIPLE.md](../BENTO_BOX_PRINCIPLE.md)
2. ✅ Confirm: "I understand the Bento Box Principle and will follow it"
3. ✅ Use TodoWrite to track all tasks
4. ✅ Phase 8 approved by Karen

## Goals

1. Embed example .bento.yaml files in binary (go:embed)
2. Create Examples section in browser
3. Visual indicator for read-only examples
4. Copy-from-template functionality
5. Examples for each major node type
6. Clear distinction between user bentos and examples
7. Validate Bento Box compliance

## Examples Structure

### Example Categories

1. **HTTP Requests** - Basic API calls
2. **Data Transformation** - JQ filtering and mapping
3. **Conditionals** - If/else logic
4. **Loops** - Iteration examples
5. **Sequences** - Multi-step workflows
6. **Complete Workflows** - End-to-end examples

## File Structure

```
pkg/omise/examples/
├── examples.go          # Embedded examples (NEW)
├── examples_test.go     # Tests (NEW)
├── templates/           # Template YAML files (NEW)
│   ├── http-get.bento.yaml
│   ├── http-post.bento.yaml
│   ├── transform-jq.bento.yaml
│   ├── conditional-if.bento.yaml
│   ├── loop-for.bento.yaml
│   ├── sequence.bento.yaml
│   └── complete-api-workflow.bento.yaml

pkg/omise/screens/
├── browser.go           # Add examples mode (modify)
└── examples.go          # Examples browser (NEW)
```

## Deliverables

### 1. Example Templates

**Files**: `pkg/omise/examples/templates/*.bento.yaml` (NEW)

**http-get.bento.yaml:**
```yaml
version: "1.0"
type: http
name: Simple GET Request
parameters:
  method: GET
  url: https://api.github.com/users/octocat
```

**http-post.bento.yaml:**
```yaml
version: "1.0"
type: http
name: POST with JSON Body
parameters:
  method: POST
  url: https://httpbin.org/post
  headers:
    Content-Type: application/json
  body: '{"message": "Hello from Bento!"}'
```

**transform-jq.bento.yaml:**
```yaml
version: "1.0"
type: transform.jq
name: Extract User IDs
parameters:
  filter: ".data | map(.id)"
```

**conditional-if.bento.yaml:**
```yaml
version: "1.0"
type: conditional.if
name: Check Success Status
parameters:
  condition: ".status == 200"
  then:
    type: http
    name: Success Handler
    parameters:
      method: GET
      url: https://httpbin.org/status/200
  else:
    type: http
    name: Error Handler
    parameters:
      method: GET
      url: https://httpbin.org/status/500
```

**loop-for.bento.yaml:**
```yaml
version: "1.0"
type: loop.for
name: Process Multiple Users
parameters:
  items: [1, 2, 3, 4, 5]
  body:
    type: http
    name: Fetch User
    parameters:
      method: GET
      url: https://api.github.com/users/{{item}}
```

**sequence.bento.yaml:**
```yaml
version: "1.0"
type: group.sequence
name: Multi-Step Workflow
nodes:
  - type: http
    name: Fetch Data
    parameters:
      method: GET
      url: https://api.github.com/users/octocat

  - type: transform.jq
    name: Extract Name
    parameters:
      filter: ".name"

  - type: http
    name: Log Result
    parameters:
      method: POST
      url: https://httpbin.org/post
      body: "{{previous}}"
```

**complete-api-workflow.bento.yaml:**
```yaml
version: "1.0"
type: group.sequence
name: Complete API Workflow Example
nodes:
  - type: http
    name: Fetch User List
    parameters:
      method: GET
      url: https://api.github.com/users

  - type: transform.jq
    name: Extract Active Users
    parameters:
      filter: ".[] | select(.type == \"User\") | {login, id}"

  - type: conditional.if
    name: Check Results
    parameters:
      condition: ". | length > 0"
      then:
        type: http
        name: Post Success
        parameters:
          method: POST
          url: https://httpbin.org/post
          body: "Found {{length}} users"
      else:
        type: http
        name: Post Empty
        parameters:
          method: POST
          url: https://httpbin.org/post
          body: "No users found"
```

### 2. Examples Loader

**File**: `pkg/omise/examples/examples.go` (NEW)
**Target Size**: < 200 lines

```go
// Package examples provides built-in example bentos
package examples

import (
	_ "embed"
	"fmt"

	"bento/pkg/jubako"
	"bento/pkg/neta"
)

//go:embed templates/http-get.bento.yaml
var httpGetExample string

//go:embed templates/http-post.bento.yaml
var httpPostExample string

//go:embed templates/transform-jq.bento.yaml
var transformJQExample string

//go:embed templates/conditional-if.bento.yaml
var conditionalIfExample string

//go:embed templates/loop-for.bento.yaml
var loopForExample string

//go:embed templates/sequence.bento.yaml
var sequenceExample string

//go:embed templates/complete-api-workflow.bento.yaml
var completeWorkflowExample string

// Example represents an example bento
type Example struct {
	ID          string
	Name        string
	Description string
	Category    string
	Content     string
}

// Category groups examples
type Category struct {
	Name     string
	Examples []Example
}

// GetAll returns all examples organized by category
func GetAll() []Category {
	return []Category{
		{
			Name: "HTTP Requests",
			Examples: []Example{
				{
					ID:          "http-get",
					Name:        "Simple GET Request",
					Description: "Fetch data from an API endpoint",
					Category:    "HTTP Requests",
					Content:     httpGetExample,
				},
				{
					ID:          "http-post",
					Name:        "POST with JSON Body",
					Description: "Send data to an API endpoint",
					Category:    "HTTP Requests",
					Content:     httpPostExample,
				},
			},
		},
		{
			Name: "Data Transformation",
			Examples: []Example{
				{
					ID:          "transform-jq",
					Name:        "Extract User IDs",
					Description: "Use JQ to filter and transform JSON data",
					Category:    "Data Transformation",
					Content:     transformJQExample,
				},
			},
		},
		{
			Name: "Control Flow",
			Examples: []Example{
				{
					ID:          "conditional-if",
					Name:        "Check Success Status",
					Description: "Execute different actions based on conditions",
					Category:    "Control Flow",
					Content:     conditionalIfExample,
				},
				{
					ID:          "loop-for",
					Name:        "Process Multiple Users",
					Description: "Iterate over a list of items",
					Category:    "Control Flow",
					Content:     loopForExample,
				},
			},
		},
		{
			Name: "Complete Workflows",
			Examples: []Example{
				{
					ID:          "sequence",
					Name:        "Multi-Step Workflow",
					Description: "Chain multiple operations together",
					Category:    "Complete Workflows",
					Content:     sequenceExample,
				},
				{
					ID:          "complete-api-workflow",
					Name:        "Complete API Workflow",
					Description: "Full example with HTTP, transform, and conditionals",
					Category:    "Complete Workflows",
					Content:     completeWorkflowExample,
				},
			},
		},
	}
}

// Get returns a specific example by ID
func Get(id string) (*Example, error) {
	categories := GetAll()
	for _, cat := range categories {
		for _, ex := range cat.Examples {
			if ex.ID == id {
				return &ex, nil
			}
		}
	}
	return nil, fmt.Errorf("example not found: %s", id)
}

// Parse parses an example into a Definition
func Parse(ex Example) (neta.Definition, error) {
	parser := jubako.NewParser()
	return parser.ParseBytes([]byte(ex.Content))
}

// List returns all examples as a flat list
func List() []Example {
	categories := GetAll()
	examples := []Example{}

	for _, cat := range categories {
		examples = append(examples, cat.Examples...)
	}

	return examples
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Example management
- ✅ Embedded files (no external deps)
- ✅ Functions < 20 lines
- ✅ File < 200 lines

### 3. Enhanced Browser with Examples

**File**: `pkg/omise/screens/browser.go` (modify)

Add mode switching:

```go
// BrowserMode defines browser mode
type BrowserMode int

const (
	BrowserModeUser BrowserMode = iota
	BrowserModeExamples
)

// Browser add mode field
type Browser struct {
	// ... existing fields
	mode BrowserMode
}

// Update handleKey add mode toggle
func (b Browser) handleKey(msg tea.KeyMsg) (Browser, tea.Cmd) {
	// ... existing code

	switch msg.String() {
	// ... existing cases

	case "x":
		// Toggle examples mode
		if b.mode == BrowserModeUser {
			b.mode = BrowserModeExamples
			return b.loadExamples()
		} else {
			b.mode = BrowserModeUser
			return b.refreshList()
		}

	// ... rest of cases
	}
}

// loadExamples loads example bentos
func (b Browser) loadExamples() (Browser, tea.Cmd) {
	examples := examples.List()
	items := make([]list.Item, len(examples))

	for i, ex := range examples {
		items[i] = exampleItem{
			id:          ex.ID,
			name:        ex.Name,
			description: ex.Description,
			category:    ex.Category,
		}
	}

	b.list = components.NewStyledList(items, "📚 Example Bentos (Read-Only)")
	return b, nil
}
```

Add example item type:

```go
// exampleItem represents an example in the list
type exampleItem struct {
	id          string
	name        string
	description string
	category    string
}

// Title returns the item title
func (i exampleItem) Title() string {
	return fmt.Sprintf("📖 %s", i.name)
}

// Description returns the item description
func (i exampleItem) Description() string {
	return fmt.Sprintf("%s • %s", i.category, i.description)
}

// FilterValue returns the value to filter by
func (i exampleItem) FilterValue() string {
	return i.name
}
```

Handle example selection:

```go
// In handleKey, when Enter/Space/r pressed on example
if b.mode == BrowserModeExamples {
	if item, ok := b.list.SelectedItem().(exampleItem); ok {
		// Copy example to user bentos
		return b, b.copyExample(item)
	}
}

// copyExample copies example to user bento
func (b Browser) copyExample(item exampleItem) tea.Cmd {
	return func() tea.Msg {
		ex, err := examples.Get(item.id)
		if err != nil {
			return BentoOperationCompleteMsg{
				Operation: "copy-example",
				Success:   false,
				Error:     err,
			}
		}

		def, err := examples.Parse(*ex)
		if err != nil {
			return BentoOperationCompleteMsg{
				Operation: "copy-example",
				Success:   false,
				Error:     err,
			}
		}

		// Save with unique name
		newName := fmt.Sprintf("%s-from-example", def.Name)
		if err := b.store.Save(newName, def); err != nil {
			return BentoOperationCompleteMsg{
				Operation: "copy-example",
				Success:   false,
				Error:     err,
			}
		}

		return BentoOperationCompleteMsg{
			Operation: "copy-example",
			Success:   true,
		}
	}
}
```

### 4. Visual Indicators

Update list title and help to show mode:

```go
// renderFooter update
func (b Browser) renderFooter() string {
	if b.mode == BrowserModeExamples {
		return styles.Subtle.Render("enter: Copy to My Bentos • x: Back to My Bentos • ?: Help")
	}

	return styles.Subtle.Render("enter/r: Run • e: Edit • c: Copy • d: Delete • n: New • x: Examples • ?: Help")
}
```

### 5. Tests

**File**: `pkg/omise/examples/examples_test.go` (NEW)

```go
package examples

import (
	"testing"
)

func TestGetAll(t *testing.T) {
	categories := GetAll()

	if len(categories) == 0 {
		t.Error("should have categories")
	}

	totalExamples := 0
	for _, cat := range categories {
		totalExamples += len(cat.Examples)
	}

	if totalExamples < 5 {
		t.Errorf("should have at least 5 examples, got %d", totalExamples)
	}
}

func TestGet(t *testing.T) {
	ex, err := Get("http-get")
	if err != nil {
		t.Fatalf("should find http-get example: %v", err)
	}

	if ex.Name != "Simple GET Request" {
		t.Errorf("wrong example name: %s", ex.Name)
	}
}

func TestParse(t *testing.T) {
	ex, err := Get("http-get")
	if err != nil {
		t.Fatal(err)
	}

	def, err := Parse(*ex)
	if err != nil {
		t.Fatalf("should parse example: %v", err)
	}

	if def.Version != "1.0" {
		t.Error("example should have version 1.0")
	}

	if def.Type != "http" {
		t.Errorf("wrong type: %s", def.Type)
	}
}

func TestList(t *testing.T) {
	examples := List()

	if len(examples) < 5 {
		t.Errorf("should have at least 5 examples, got %d", len(examples))
	}

	// Check each example has required fields
	for _, ex := range examples {
		if ex.ID == "" {
			t.Error("example missing ID")
		}
		if ex.Name == "" {
			t.Error("example missing name")
		}
		if ex.Content == "" {
			t.Error("example missing content")
		}
	}
}
```

## Integration

### Update Help

Show examples shortcut in help screen:

**File**: `pkg/omise/screens/help.go` (modify)

```go
func (h Help) View() string {
	help := `
🍱 Bento Keyboard Shortcuts

Browser:
  enter/space/r  Run bento
  e              Edit bento
  c              Copy bento
  d              Delete bento
  n              Create new bento
  x              Toggle examples          ← NEW
  ?              Toggle help
  tab            Next screen
  q              Quit

Examples:                                  ← NEW
  enter          Copy example to My Bentos
  x              Back to My Bentos

Editor:
  ↑/↓            Navigate nodes
  e              Edit selected node
  m              Move node
  d              Delete node
  r              Run bento
  v              Toggle view mode
  a              Add node
  s/enter        Save
  esc            Cancel

Press ? again to close.
`
	return styles.Subtle.Render(help)
}
```

## Validation Commands

```bash
# Test examples
go test -v ./pkg/omise/examples/

# Verify examples are embedded
go build ./cmd/bento
./bento pantry # Should show all node types used in examples

# Integration test
./bento

# In browser:
# 1. Press 'x' to switch to Examples
# 2. Should see list of example bentos with 📖 icon
# 3. List title should say "Example Bentos (Read-Only)"
# 4. Select an example
# 5. Press Enter to copy to My Bentos
# 6. Press 'x' to switch back
# 7. Should see copied bento in My Bentos
# 8. Can edit/run/delete copied bento normally
```

## Success Criteria

Phase 9 is complete when:

1. ✅ Example templates created (7+ examples)
2. ✅ Examples embedded in binary (go:embed)
3. ✅ Examples loader implemented
4. ✅ Browser mode toggle working
5. ✅ Examples list displays correctly
6. ✅ Visual indicators for read-only examples
7. ✅ Copy-from-template working
8. ✅ Help updated with examples info
9. ✅ All files < 250 lines
10. ✅ All functions < 20 lines
11. ✅ Tests passing
12. ✅ **Karen's approval granted**

## Common Pitfalls to Avoid

1. ❌ **Invalid example YAML** - All examples must parse correctly
2. ❌ **Missing versions** - All examples must have version field
3. ❌ **No copy functionality** - Users should copy, not edit examples
4. ❌ **Unclear mode** - Clear indication of user vs examples mode
5. ❌ **Poor examples** - Examples should demonstrate best practices

## Documentation Updates

Update README.md:

```markdown
## 📚 Examples

Bento includes built-in examples demonstrating various features:

1. Launch the TUI: `bento`
2. Press `x` to view examples
3. Browse example workflows
4. Press Enter to copy an example to your bentos
5. Press `x` again to return to your bentos

Example categories:
- **HTTP Requests** - API calls and data fetching
- **Data Transformation** - JQ filtering and mapping
- **Control Flow** - Conditionals and loops
- **Complete Workflows** - End-to-end examples

All examples are embedded in the binary and available offline.
```

## Next Phase

After Karen approval, proceed to **[Phase 10: Real-World Use Case](./phase-10-proof-of-concept.md)** to:
- Discuss your specific bento requirements
- Create a real-world bento using the editor
- Validate the editor/flow supports your needs
- Proof-of-concept for the entire system

## Execution Prompt

```
I'm ready to begin Phase 9: Examples & Templates.

I have read the Bento Box Principle and will follow it.

Please create examples system with:
- 7+ example templates in YAML
- Embedded examples (go:embed)
- Browser mode toggle
- Copy-from-template functionality
- Clear visual indicators

Each file < 250 lines, functions < 20 lines. I will use TodoWrite to track progress and get Karen's approval before completing.
```

---

**Phase 9 Examples**: Built-in templates and starting points 📚

**After this phase**: Ready for your real-world proof-of-concept! 🎯
