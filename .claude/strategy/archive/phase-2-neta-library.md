# Phase 2: Neta Library

**Status**: Pending
**Duration**: 3-4 hours
**Prerequisites**: Phase 1 complete, Karen approved

## Overview

Build out the concrete node type implementations (neta library). Each node type gets its own focused package following the Bento Box Principle. This phase creates the "ingredients" that users can compose into workflows.

## Pre-Work Checklist

Before starting, you MUST:

1. ✅ Read [BENTO_BOX_PRINCIPLE.md](../BENTO_BOX_PRINCIPLE.md)
2. ✅ Confirm: "I understand the Bento Box Principle and will follow it"
3. ✅ Use TodoWrite to track all tasks
4. ✅ Phase 1 approved by Karen

## Goals

1. Implement 5+ concrete node types
2. Each type in focused package (no mixing!)
3. Register all types with pantry
4. Comprehensive tests for each type
5. Create example .bento.yaml files
6. Validate Bento Box compliance

## Node Type Implementations

### Package Structure

```
pkg/neta/
├── go.mod
├── definition.go       # Core types (from Phase 1)
├── executable.go       # Interface (from Phase 1)
├── http/
│   ├── client.go       # HTTP node implementation
│   ├── client_test.go
│   └── types.go        # HTTP-specific types
├── transform/
│   ├── jq.go          # JQ transformation
│   ├── jq_test.go
│   ├── template.go    # Template transformation
│   └── template_test.go
├── conditional/
│   ├── if.go          # If/else logic
│   ├── if_test.go
│   ├── switch.go      # Switch/case logic
│   └── switch_test.go
├── loop/
│   ├── for.go         # For loop
│   ├── for_test.go
│   ├── while.go       # While loop
│   └── while_test.go
└── group/
    ├── sequence.go    # Sequential execution
    ├── sequence_test.go
    ├── parallel.go    # Parallel execution
    └── parallel_test.go
```

## Deliverables

### 1. HTTP Node (`pkg/neta/http/`)

**Purpose**: HTTP request execution
**File Size Target**: < 150 lines per file

#### `pkg/neta/http/client.go`
```go
// Package http provides HTTP request execution nodes.
package http

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "strings"

    "bento/pkg/neta"
)

// Client executes HTTP requests.
type Client struct {
    client *http.Client
}

// New creates a new HTTP client node.
func New() *Client {
    return &Client{
        client: &http.Client{},
    }
}

// Execute performs an HTTP request.
func (c *Client) Execute(ctx context.Context, params map[string]interface{}) (neta.Result, error) {
    req, err := buildRequest(ctx, params)
    if err != nil {
        return neta.Result{}, err
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return neta.Result{}, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return neta.Result{}, err
    }

    return neta.Result{
        Output: string(body),
        Metadata: map[string]interface{}{
            "status_code": resp.StatusCode,
            "headers":     resp.Header,
        },
    }, nil
}

// buildRequest creates an HTTP request from parameters.
func buildRequest(ctx context.Context, params map[string]interface{}) (*http.Request, error) {
    method := getStringParam(params, "method", "GET")
    url := getStringParam(params, "url", "")
    if url == "" {
        return nil, fmt.Errorf("url parameter required")
    }

    body := getStringParam(params, "body", "")
    req, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(body))
    if err != nil {
        return nil, err
    }

    // Add headers if provided
    if headers, ok := params["headers"].(map[string]string); ok {
        for k, v := range headers {
            req.Header.Set(k, v)
        }
    }

    return req, nil
}

// getStringParam extracts a string parameter with default.
func getStringParam(params map[string]interface{}, key, defaultVal string) string {
    if val, ok := params[key].(string); ok {
        return val
    }
    return defaultVal
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: HTTP execution only
- ✅ Functions < 20 lines
- ✅ Clear helper functions (buildRequest, getStringParam)
- ✅ File < 150 lines

### 2. Transform Node (`pkg/neta/transform/`)

**Purpose**: Data transformation (JQ, templates)
**File Size Target**: < 100 lines per file

#### `pkg/neta/transform/jq.go`
```go
// Package transform provides data transformation nodes.
package transform

import (
    "context"
    "fmt"

    "github.com/itchyny/gojq"
    "bento/pkg/neta"
)

// JQ applies jq transformations to data.
type JQ struct{}

// New creates a new JQ transformer.
func New() *JQ {
    return &JQ{}
}

// Execute applies a jq query to input data.
func (j *JQ) Execute(ctx context.Context, params map[string]interface{}) (neta.Result, error) {
    query := getStringParam(params, "query", ".")
    input := params["input"]

    result, err := applyQuery(query, input)
    if err != nil {
        return neta.Result{}, fmt.Errorf("jq transform failed: %w", err)
    }

    return neta.Result{Output: result}, nil
}

// applyQuery executes a jq query on data.
func applyQuery(queryStr string, data interface{}) (interface{}, error) {
    query, err := gojq.Parse(queryStr)
    if err != nil {
        return nil, err
    }

    iter := query.Run(data)
    v, ok := iter.Next()
    if !ok {
        return nil, fmt.Errorf("no result")
    }
    if err, ok := v.(error); ok {
        return nil, err
    }

    return v, nil
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: JQ transformation
- ✅ Small functions (< 20 lines)
- ✅ Clear separation (Execute vs applyQuery)
- ✅ File < 100 lines

### 3. Conditional Node (`pkg/neta/conditional/`)

**Purpose**: If/else and switch logic
**File Size Target**: < 100 lines per file

#### `pkg/neta/conditional/if.go`
```go
// Package conditional provides conditional execution nodes.
package conditional

import (
    "context"
    "fmt"

    "bento/pkg/neta"
)

// If executes nodes based on conditions.
type If struct {
    itamae Executor
}

// Executor can execute neta definitions.
type Executor interface {
    Execute(ctx context.Context, def neta.Definition) (neta.Result, error)
}

// New creates a new If node.
func New(executor Executor) *If {
    return &If{itamae: executor}
}

// Execute evaluates condition and runs appropriate branch.
func (i *If) Execute(ctx context.Context, params map[string]interface{}) (neta.Result, error) {
    condition := getBoolParam(params, "condition", false)

    if condition {
        return i.executeBranch(ctx, params, "then")
    }
    return i.executeBranch(ctx, params, "else")
}

// executeBranch runs the specified branch.
func (i *If) executeBranch(ctx context.Context, params map[string]interface{}, branch string) (neta.Result, error) {
    def, ok := params[branch].(neta.Definition)
    if !ok {
        return neta.Result{}, nil
    }
    return i.itamae.Execute(ctx, def)
}

// getBoolParam extracts a bool parameter.
func getBoolParam(params map[string]interface{}, key string, defaultVal bool) bool {
    if val, ok := params[key].(bool); ok {
        return val
    }
    return defaultVal
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Conditional logic
- ✅ Dependency injection (Executor)
- ✅ Functions < 15 lines
- ✅ File < 100 lines

### 4. Loop Node (`pkg/neta/loop/`)

**Purpose**: Iteration over collections
**File Size Target**: < 100 lines per file

#### `pkg/neta/loop/for.go`
```go
// Package loop provides iteration nodes.
package loop

import (
    "context"
    "fmt"

    "bento/pkg/neta"
)

// For iterates over a collection.
type For struct {
    itamae Executor
}

// Executor can execute neta definitions.
type Executor interface {
    Execute(ctx context.Context, def neta.Definition) (neta.Result, error)
}

// New creates a new For loop node.
func New(executor Executor) *For {
    return &For{itamae: executor}
}

// Execute iterates over items and executes body for each.
func (f *For) Execute(ctx context.Context, params map[string]interface{}) (neta.Result, error) {
    items, err := getItems(params)
    if err != nil {
        return neta.Result{}, err
    }

    body, ok := params["body"].(neta.Definition)
    if !ok {
        return neta.Result{}, fmt.Errorf("body required")
    }

    return f.iterate(ctx, items, body)
}

// iterate executes body for each item.
func (f *For) iterate(ctx context.Context, items []interface{}, body neta.Definition) (neta.Result, error) {
    results := make([]neta.Result, 0, len(items))
    for _, item := range items {
        // Inject current item into body params
        bodyParams := make(map[string]interface{})
        bodyParams["item"] = item
        body.Parameters = bodyParams

        result, err := f.itamae.Execute(ctx, body)
        if err != nil {
            return neta.Result{}, err
        }
        results = append(results, result)
    }
    return neta.Result{Output: results}, nil
}

// getItems extracts the items array from params.
func getItems(params map[string]interface{}) ([]interface{}, error) {
    items, ok := params["items"].([]interface{})
    if !ok {
        return nil, fmt.Errorf("items parameter required")
    }
    return items, nil
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Iteration
- ✅ Composable with other nodes
- ✅ Functions < 20 lines
- ✅ File < 100 lines

### 5. Group Node (`pkg/neta/group/`)

**Purpose**: Sequential and parallel execution
**File Size Target**: < 100 lines per file

#### `pkg/neta/group/sequence.go`
```go
// Package group provides group execution nodes.
package group

import (
    "context"

    "bento/pkg/neta"
)

// Sequence executes nodes in order.
type Sequence struct {
    itamae Executor
}

// Executor can execute neta definitions.
type Executor interface {
    Execute(ctx context.Context, def neta.Definition) (neta.Result, error)
}

// NewSequence creates a sequential group executor.
func NewSequence(executor Executor) *Sequence {
    return &Sequence{itamae: executor}
}

// Execute runs nodes one after another.
func (s *Sequence) Execute(ctx context.Context, params map[string]interface{}) (neta.Result, error) {
    nodes := getNodes(params)
    results := make([]neta.Result, 0, len(nodes))

    for _, node := range nodes {
        result, err := s.itamae.Execute(ctx, node)
        if err != nil {
            return neta.Result{}, err
        }
        results = append(results, result)
    }

    return neta.Result{Output: results}, nil
}

// getNodes extracts child nodes from params.
func getNodes(params map[string]interface{}) []neta.Definition {
    nodes, ok := params["nodes"].([]neta.Definition)
    if !ok {
        return []neta.Definition{}
    }
    return nodes
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Sequential execution
- ✅ Simple, focused implementation
- ✅ Functions < 15 lines
- ✅ File < 100 lines

## Example .bento.yaml Files

### Example 1: HTTP Request
```yaml
# examples/http-get.bento.yaml
type: http
name: Fetch User Data
parameters:
  method: GET
  url: https://api.example.com/users/1
  headers:
    Accept: application/json
```

### Example 2: Transform Pipeline
```yaml
# examples/transform-pipeline.bento.yaml
type: group
name: Transform Pipeline
nodes:
  - type: http
    name: Fetch Data
    parameters:
      url: https://api.example.com/data

  - type: transform
    name: Extract Names
    parameters:
      query: .users[].name
```

### Example 3: Conditional Logic
```yaml
# examples/conditional.bento.yaml
type: conditional
name: Check and Process
parameters:
  condition: true
  then:
    type: http
    name: Success Path
    parameters:
      url: https://api.example.com/success
  else:
    type: http
    name: Failure Path
    parameters:
      url: https://api.example.com/failure
```

## Testing Strategy

Each node type must have:

1. **Happy path tests** - Normal execution
2. **Error cases** - Invalid params, network errors, etc
3. **Edge cases** - Empty data, large data, timeouts
4. **Integration tests** - With real itamae orchestrator

### Example Test Structure
```go
func TestHTTPClient_Execute(t *testing.T) {
    tests := []struct {
        name    string
        params  map[string]interface{}
        want    neta.Result
        wantErr bool
    }{
        {
            name: "successful GET request",
            params: map[string]interface{}{
                "method": "GET",
                "url":    "https://httpbin.org/get",
            },
            wantErr: false,
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Pantry Registration

Update pantry to register all node types:

```go
// In cmd/bento or initialization code
func initializePantry() *pantry.Pantry {
    p := pantry.New()

    p.Register("http", http.New())
    p.Register("transform.jq", transform.New())
    p.Register("conditional.if", conditional.New(itamae))
    p.Register("loop.for", loop.New(itamae))
    p.Register("group.sequence", group.NewSequence(itamae))
    p.Register("group.parallel", group.NewParallel(itamae))

    return p
}
```

## Validation Commands

Before marking phase complete:

```bash
# Format
go fmt ./...

# Lint
golangci-lint run

# Test with race detector
go test -v -race ./pkg/neta/...

# Build
go build ./cmd/bento

# File size check
find pkg/neta -name "*.go" -exec wc -l {} + | sort -rn | head -10

# Check for utils packages (should be empty)
find pkg/neta -type d -name "*util*"
```

## Success Criteria

Phase 2 is complete when:

1. ✅ 5+ node types implemented (http, transform, conditional, loop, group)
2. ✅ Each type in focused package (no mixing!)
3. ✅ All types registered with pantry
4. ✅ Comprehensive tests (>80% coverage)
5. ✅ 3+ example .bento.yaml files
6. ✅ All tests passing (go test -race ./...)
7. ✅ All files < 250 lines
8. ✅ All functions < 20 lines
9. ✅ Zero utils packages
10. ✅ golangci-lint clean
11. ✅ **Karen's approval granted**

## Common Pitfalls to Avoid

1. ❌ **Creating neta/utils/** - Keep helpers in same package
2. ❌ **Mixing concerns** - HTTP client shouldn't do transformations
3. ❌ **Large files** - Split http/client.go if it grows too large
4. ❌ **God functions** - Extract helpers for readability
5. ❌ **No tests** - Every public function needs tests
6. ❌ **Circular dependencies** - Nodes shouldn't depend on each other

## Next Phase

After Karen approval, proceed to **[Phase 3: CLI Commands](./phase-3-cli-commands.md)** to:
- Implement Cobra commands (prepare, pack, pantry, taste)
- Add Viper configuration
- Create command-line interface

## Execution Prompt

```
I'm ready to begin Phase 2: Neta Library.

I have read the Bento Box Principle and will follow it.

Please implement concrete node types:
- HTTP client
- JQ transform
- Conditional (if/else)
- Loop (for)
- Group (sequence/parallel)

Each type gets its own focused package. I will use TodoWrite to track progress and get Karen's approval before completing.
```

---

**Phase 2 Neta Library**: Node type implementations 🍱
