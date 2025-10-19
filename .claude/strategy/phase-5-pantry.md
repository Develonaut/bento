# Phase 5: Pantry Package - "Ingredient Storage"

**Duration:** 2-3 days
**Package:** `pkg/pantry/`
**Dependencies:** `pkg/neta`

---

## TDD Philosophy

> **Write tests FIRST to define contracts**

Registry tests should verify:
1. Register neta types and retrieve them by name
2. Thread-safe concurrent access (multiple goroutines)
3. List all registered neta types
4. Clear error when requesting unregistered type
5. Prevent duplicate registration

---

## Phase Overview

The pantry package provides a thread-safe registry for all neta types. Like a well-organized pantry where ingredients are stored and easily retrieved, this registry keeps track of all available neta implementations.

Key responsibilities:
- **Register:** Add neta types to the registry
- **Get:** Retrieve neta implementation by type name
- **List:** Get all registered neta types
- **Thread-safety:** Safe concurrent access using sync.RWMutex

---

## Success Criteria

**Phase 5 Complete When:**
- [ ] Register() adds neta types to registry
- [ ] Get() retrieves neta by type name
- [ ] List() returns all registered types
- [ ] Thread-safe (concurrent goroutines can access safely)
- [ ] Clear errors for missing types
- [ ] Integration tests for concurrent access
- [ ] Files < 250 lines
- [ ] File-level documentation complete
- [ ] `/code-review` run with Karen + Colossus approval

---

## Test-First Approach

### Step 1: Define registry interface via tests

Create `pkg/pantry/pantry_test.go`:

```go
package pantry_test

import (
	"context"
	"sync"
	"testing"

	"github.com/Develonaut/bento/pkg/neta"
	"github.com/Develonaut/bento/pkg/pantry"
)

// Test: Register and retrieve a neta type
func TestPantry_RegisterAndGet(t *testing.T) {
	p := pantry.New()

	// Register a mock neta
	mockNeta := &MockNeta{name: "test-neta"}
	p.Register("test-neta", mockNeta)

	// Retrieve it
	retrieved, err := p.Get("test-neta")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved != mockNeta {
		t.Error("Retrieved neta is not the same as registered")
	}
}

// Test: Get unregistered type should return error
func TestPantry_GetUnregistered(t *testing.T) {
	p := pantry.New()

	_, err := p.Get("nonexistent-type")
	if err == nil {
		t.Fatal("Expected error for unregistered type")
	}

	if !strings.Contains(err.Error(), "not registered") {
		t.Errorf("Error should mention 'not registered': %v", err)
	}

	if !strings.Contains(err.Error(), "nonexistent-type") {
		t.Errorf("Error should mention the type name: %v", err)
	}
}

// Test: List all registered types
func TestPantry_List(t *testing.T) {
	p := pantry.New()

	// Register multiple types
	p.Register("type-a", &MockNeta{name: "a"})
	p.Register("type-b", &MockNeta{name: "b"})
	p.Register("type-c", &MockNeta{name: "c"})

	// List them
	types := p.List()

	if len(types) != 3 {
		t.Errorf("Expected 3 types, got %d", len(types))
	}

	// Verify all types present
	typeMap := make(map[string]bool)
	for _, typeName := range types {
		typeMap[typeName] = true
	}

	if !typeMap["type-a"] || !typeMap["type-b"] || !typeMap["type-c"] {
		t.Error("Not all types were listed")
	}
}

// Test: Thread-safe concurrent access
func TestPantry_ConcurrentAccess(t *testing.T) {
	p := pantry.New()

	// Register initial types
	for i := 0; i < 10; i++ {
		p.Register(fmt.Sprintf("type-%d", i), &MockNeta{name: fmt.Sprintf("neta-%d", i)})
	}

	// Concurrent reads
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			typeNum := i % 10
			typeName := fmt.Sprintf("type-%d", typeNum)

			_, err := p.Get(typeName)
			if err != nil {
				t.Errorf("Concurrent Get failed: %v", err)
			}

			// Also list
			_ = p.List()
		}(i)
	}

	wg.Wait()
}

// Test: Duplicate registration should overwrite (or error - your choice)
func TestPantry_DuplicateRegistration(t *testing.T) {
	p := pantry.New()

	neta1 := &MockNeta{name: "first"}
	neta2 := &MockNeta{name: "second"}

	p.Register("test-type", neta1)
	p.Register("test-type", neta2) // Duplicate

	retrieved, _ := p.Get("test-type")

	// Verify which one is stored (last one should win)
	if retrieved != neta2 {
		t.Error("Duplicate registration should overwrite previous")
	}
}

// Mock Neta for testing
type MockNeta struct {
	name string
}

func (m *MockNeta) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	return map[string]interface{}{"mock": m.name}, nil
}
```

### Step 2: Test factory pattern for creating neta instances

```go
// Test: GetNew should return a NEW instance (not shared)
func TestPantry_GetNew(t *testing.T) {
	p := pantry.New()

	// Register a factory function instead of instance
	p.RegisterFactory("http-request", func() neta.Executable {
		return &MockNeta{name: "http-request"}
	})

	// Get two instances
	instance1, err := p.GetNew("http-request")
	if err != nil {
		t.Fatalf("GetNew failed: %v", err)
	}

	instance2, err := p.GetNew("http-request")
	if err != nil {
		t.Fatalf("GetNew failed: %v", err)
	}

	// Should be different instances
	if instance1 == instance2 {
		t.Error("GetNew should return new instances, not the same one")
	}
}
```

---

## File Structure

```
pkg/pantry/
├── pantry.go           # Main registry implementation (~200 lines)
├── factory.go          # Factory pattern for neta creation (~100 lines)
└── pantry_test.go      # Integration tests (~250 lines)
```

---

## Implementation Guidance

**File: `pkg/pantry/pantry.go`**

```go
// Package pantry provides a thread-safe registry for neta types.
//
// Like a well-organized pantry stores ingredients, this package stores and
// provides access to all available neta implementations.
//
// Usage:
//
//	p := pantry.New()
//
//	// Register neta types
//	p.RegisterFactory("http-request", func() neta.Executable {
//	    return httpneta.New()
//	})
//
//	// Retrieve a new instance
//	netaInstance, err := p.GetNew("http-request")
//
//	// List all types
//	types := p.List()
//
// Thread Safety:
// The pantry uses sync.RWMutex to ensure safe concurrent access from multiple
// goroutines. Multiple readers can access simultaneously, but writes are exclusive.
//
// Learn more about sync.RWMutex:
// https://pkg.go.dev/sync#RWMutex
package pantry

import (
	"fmt"
	"sync"

	"github.com/Develonaut/bento/pkg/neta"
)

// Factory is a function that creates a new neta instance.
// Each call should return a NEW instance, not a shared one.
type Factory func() neta.Executable

// Pantry is a thread-safe registry for neta types.
type Pantry struct {
	mu        sync.RWMutex          // Protects the maps
	factories map[string]Factory    // Type name -> factory function
}

// New creates a new Pantry with all built-in neta types registered.
func New() *Pantry {
	p := &Pantry{
		factories: make(map[string]Factory),
	}

	// Register all built-in types
	// (This will be filled in during Phase 1 implementation)
	// p.RegisterFactory("edit-fields", func() neta.Executable { return editfields.New() })
	// p.RegisterFactory("http-request", func() neta.Executable { return httpneta.New() })
	// ... etc for all 10 types

	return p
}

// RegisterFactory registers a neta type with a factory function.
//
// The factory function should return a NEW instance each time it's called.
// If a type is already registered, this will overwrite it.
func (p *Pantry) RegisterFactory(typeName string, factory Factory) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.factories[typeName] = factory
}

// Register registers a neta type with a shared instance.
//
// DEPRECATED: Use RegisterFactory instead for better isolation.
// This is kept for backwards compatibility.
func (p *Pantry) Register(typeName string, instance neta.Executable) {
	p.RegisterFactory(typeName, func() neta.Executable {
		return instance
	})
}

// GetNew creates and returns a new instance of the specified neta type.
//
// Returns an error if the type is not registered.
func (p *Pantry) GetNew(typeName string) (neta.Executable, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	factory, exists := p.factories[typeName]
	if !exists {
		return nil, fmt.Errorf("neta type '%s' is not registered in pantry. Available types: %v",
			typeName, p.listUnsafe())
	}

	return factory(), nil
}

// Get returns a neta instance (calls GetNew for compatibility).
//
// DEPRECATED: Use GetNew for clarity that a new instance is created.
func (p *Pantry) Get(typeName string) (neta.Executable, error) {
	return p.GetNew(typeName)
}

// List returns all registered neta type names.
func (p *Pantry) List() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.listUnsafe()
}

// listUnsafe returns all types without locking (caller must hold lock).
func (p *Pantry) listUnsafe() []string {
	types := make([]string, 0, len(p.factories))

	for typeName := range p.factories {
		types = append(types, typeName)
	}

	// Sort for consistent output
	sort.Strings(types)

	return types
}

// Has checks if a neta type is registered.
func (p *Pantry) Has(typeName string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	_, exists := p.factories[typeName]
	return exists
}
```

---

## Common Go Pitfalls to Avoid

1. **RWMutex usage**: Use RLock for reads, Lock for writes
   ```go
   // ❌ BAD - uses write lock for read operation
   func (p *Pantry) Get(name string) {
       p.mu.Lock()
       defer p.mu.Unlock()
       return p.factories[name]
   }

   // ✅ GOOD - uses read lock for read operation
   func (p *Pantry) Get(name string) {
       p.mu.RLock()
       defer p.mu.RUnlock()
       return p.factories[name]
   }
   ```

2. **Defer unlock**: ALWAYS defer unlock immediately after lock
   ```go
   // ❌ BAD - if function panics, mutex stays locked forever
   p.mu.Lock()
   // ... do stuff ...
   p.mu.Unlock()

   // ✅ GOOD - defer ensures unlock even on panic
   p.mu.Lock()
   defer p.mu.Unlock()
   ```

3. **Factory returns new instances**: Don't return the same instance
   ```go
   // ❌ BAD - all calls share same instance
   var sharedInstance = &HttpNeta{}
   factory := func() neta.Executable {
       return sharedInstance
   }

   // ✅ GOOD - each call gets new instance
   factory := func() neta.Executable {
       return httpneta.New()
   }
   ```

---

## Critical for Phase 8

**Factory Pattern:**
- Itamae will call GetNew() for each neta in the bento
- Each neta instance must be independent (no shared state)
- Product automation runs 50+ iterations - must not leak state between iterations

**Thread Safety:**
- If we implement parallel neta execution, multiple goroutines will access pantry
- RWMutex allows concurrent reads (fast) but exclusive writes

---

## Bento Box Principle Checklist

- [ ] Files < 250 lines (pantry.go ~200)
- [ ] Functions < 20 lines
- [ ] Single responsibility (registry only)
- [ ] Thread-safe (RWMutex)
- [ ] Clear error messages (list available types)
- [ ] File-level documentation

---

## Phase Completion

**Phase 5 MUST end with:**

1. All tests passing (`go test ./pkg/pantry/...`)
2. Run `/code-review` slash command
3. Address feedback from Karen and Colossus
4. Get explicit approval from both agents
5. Document any decisions in `.claude/strategy/`

**Do not proceed to Phase 6 until code review is approved.**

---

## Claude Prompt Template

```
I need to implement Phase 5: pantry (registry package) following TDD principles.

Please read:
- .claude/strategy/phase-5-pantry.md (this file)
- .claude/BENTO_BOX_PRINCIPLE.md

Then:

1. Create `pkg/pantry/pantry_test.go` with integration tests for:
   - RegisterFactory and GetNew
   - Get unregistered type (clear error listing available types)
   - List all types
   - Concurrent access (100 goroutines reading simultaneously)
   - Factory pattern (each call returns NEW instance)

2. Watch the tests fail

3. Implement to make tests pass:
   - pkg/pantry/pantry.go (~200 lines)
   - Use sync.RWMutex for thread safety
   - Factory pattern for creating instances

4. Add file-level documentation explaining:
   - What a pantry is (ingredient storage)
   - How to use RegisterFactory, GetNew, List
   - Thread safety with RWMutex
   - Why factory pattern (new instance each time)

Remember:
- Write tests FIRST
- Files < 250 lines
- Functions < 20 lines
- Thread-safe with RWMutex
- Factory returns NEW instances

When complete, run `/code-review` and get Karen + Colossus approval.
```

---

## Dependencies

No additional dependencies needed - uses stdlib only:
- `sync` for RWMutex
- `sort` for sorted type list

---

## Notes

- Factory pattern prevents shared state between neta instances
- RWMutex allows concurrent reads (fast for List() and GetNew())
- Error messages should list available types (helps with typos)
- New() should pre-register all 10 built-in neta types
- This is a simple package - don't over-engineer it

---

**Status:** Ready for implementation
**Next Phase:** Phase 6 (itamae orchestration) - depends on completion of Phase 5
