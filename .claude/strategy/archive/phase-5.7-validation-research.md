# Guilliman's Go Validation Research

## Executive Summary

After researching Go's standard library, major projects (Kubernetes, Docker, Terraform), and popular validation libraries, I recommend a **programmatic validation approach** using simple Go code rather than struct tags or external libraries. This aligns with Go's philosophy of "clear is better than clever" and keeps us dependency-free.

## Standard Library Approach

The Go standard library provides **NO built-in validation**. The `encoding/json` package only handles marshaling/unmarshaling with struct tags for field mapping:

```go
type Config struct {
    URL    string            `json:"url"`
    Method string            `json:"method,omitempty"`
    Headers map[string]string `json:"headers,omitempty"`
}
```

**Verdict**: Standard library doesn't solve validation. We need our own solution.

## Industry Patterns

### Kubernetes Approach
- Uses **OpenAPI schemas** generated from Go structs
- **Kubebuilder markers** (comments) for validation rules:
  ```go
  // +kubebuilder:validation:Required
  // +kubebuilder:validation:MinLength=1
  Name string `json:"name"`
  ```
- Generates CRD validation schemas from these markers
- **Complex toolchain** (controller-gen, kube-openapi)

**Pros**: Powerful, declarative, generates OpenAPI docs
**Cons**: Heavy dependencies, complex tooling, overkill for our needs

### Docker Compose Approach
- Uses **gojsonschema** library for JSON Schema validation
- Defines schemas in JSON/YAML files
- Validates against external schema files

**Pros**: Standard JSON Schema format
**Cons**: External dependency, schemas in separate files, not Go-idiomatic

### Terraform Approach
- **Plugin SDK v2**: Uses `Schema` structs with `ValidateFunc`:
  ```go
  &schema.Schema{
      Type:     schema.TypeString,
      Required: true,
      ValidateFunc: validation.StringInSlice([]string{"GET", "POST"}, false),
  }
  ```
- **go-cty**: Dynamic type system for configuration values
- Built specifically for HCL/configuration languages

**Pros**: Powerful, type-safe, good for config languages
**Cons**: Complex, external dependency (go-cty), designed for HCL

## Popular Go Libraries

### go-playground/validator (25.5k stars)
```go
type Config struct {
    URL     string `validate:"required,url"`
    Method  string `validate:"omitempty,oneof=GET POST PUT DELETE"`
    Headers map[string]string `validate:"dive,keys,required,endkeys,required"`
}
```

**Pros**:
- Most popular validation library
- Declarative struct tags
- Rich set of built-in validators

**Cons**:
- String-based tags (no compile-time checking)
- External dependency
- Complex tag syntax for nested structures
- Error messages need customization

### ozzo-validation (3.7k stars)
```go
func (c Config) Validate() error {
    return validation.ValidateStruct(&c,
        validation.Field(&c.URL, validation.Required, is.URL),
        validation.Field(&c.Method, validation.In("GET", "POST", "PUT", "DELETE")),
    )
}
```

**Pros**:
- Programmatic (compile-time safety)
- Cleaner for complex validations
- Good error messages

**Cons**:
- External dependency
- More verbose than struct tags
- Less popular (smaller community)

## Recommendation for Bento

### Recommended Approach: **Simple Programmatic Validation**

Write our own minimal validation using standard Go patterns. No external dependencies.

#### Why This Is Most Idiomatic for Go:
1. **No magic** - validation logic is explicit and debuggable
2. **No dependencies** - aligns with "a little copying is better than a little dependency"
3. **Compile-time safety** - no string-based tags to mistype
4. **Simple and clear** - anyone can understand the validation logic
5. **Tailored to our needs** - validates `map[string]interface{}` from YAML

#### Example Implementation:

```go
// pkg/validation/schema.go
package validation

import (
    "fmt"
    "net/url"
)

// Schema defines validation rules for a node type
type Schema struct {
    Required []string
    Optional []string
    Validators map[string]FieldValidator
    Defaults map[string]interface{}
}

// FieldValidator validates a single field value
type FieldValidator interface {
    Validate(value interface{}) error
    Metadata() FieldMetadata // For UI generation
}

// FieldMetadata provides information for UI generation
type FieldMetadata struct {
    Type        string   // "string", "number", "boolean", "select"
    Description string
    Options     []string // For select/enum fields
    Default     interface{}
}

// StringField validates string values
type StringField struct {
    Description string
    Required    bool
    Options     []string // Empty for free text, populated for enum
    Default     string
}

func (f StringField) Validate(value interface{}) error {
    str, ok := value.(string)
    if !ok {
        return fmt.Errorf("expected string, got %T", value)
    }

    if f.Required && str == "" {
        return fmt.Errorf("field is required")
    }

    if len(f.Options) > 0 {
        valid := false
        for _, opt := range f.Options {
            if str == opt {
                valid = true
                break
            }
        }
        if !valid {
            return fmt.Errorf("must be one of: %v", f.Options)
        }
    }

    return nil
}

func (f StringField) Metadata() FieldMetadata {
    fieldType := "string"
    if len(f.Options) > 0 {
        fieldType = "select"
    }
    return FieldMetadata{
        Type:        fieldType,
        Description: f.Description,
        Options:     f.Options,
        Default:     f.Default,
    }
}

// URLField validates URL strings
type URLField struct {
    StringField
}

func (f URLField) Validate(value interface{}) error {
    if err := f.StringField.Validate(value); err != nil {
        return err
    }

    str := value.(string)
    if str == "" && !f.Required {
        return nil
    }

    _, err := url.Parse(str)
    if err != nil {
        return fmt.Errorf("invalid URL: %w", err)
    }

    return nil
}

// HTTPNodeSchema defines the schema for HTTP nodes
var HTTPNodeSchema = Schema{
    Required: []string{"url"},
    Optional: []string{"method", "headers", "body"},
    Validators: map[string]FieldValidator{
        "url": URLField{
            StringField: StringField{
                Description: "The URL to fetch",
                Required:    true,
            },
        },
        "method": StringField{
            Description: "HTTP method",
            Options:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
            Default:     "GET",
        },
        "headers": MapField{
            Description: "HTTP headers",
            KeyValidator: StringField{Description: "Header name"},
            ValueValidator: StringField{Description: "Header value"},
        },
    },
    Defaults: map[string]interface{}{
        "method": "GET",
    },
}

// Validate validates parameters against a schema
func (s Schema) Validate(params map[string]interface{}) error {
    // Check required fields
    for _, field := range s.Required {
        if _, exists := params[field]; !exists {
            return fmt.Errorf("required field '%s' is missing", field)
        }
    }

    // Validate each field
    for key, value := range params {
        validator, exists := s.Validators[key]
        if !exists {
            // Check if it's in optional fields
            isOptional := false
            for _, opt := range s.Optional {
                if opt == key {
                    isOptional = true
                    break
                }
            }
            if !isOptional && !isRequired(key, s.Required) {
                return fmt.Errorf("unknown field '%s'", key)
            }
            continue
        }

        if err := validator.Validate(value); err != nil {
            return fmt.Errorf("field '%s': %w", key, err)
        }
    }

    return nil
}

// Usage example:
params := map[string]interface{}{
    "url": "https://api.example.com",
    "method": "POST",
    "headers": map[string]interface{}{
        "Authorization": "Bearer token",
    },
}

if err := HTTPNodeSchema.Validate(params); err != nil {
    return fmt.Errorf("validation failed: %w", err)
}
```

### Pros:
- **Zero dependencies** - pure Go standard library
- **Clear and explicit** - validation logic is readable
- **Compile-time safe** - no string parsing or reflection
- **Tailored for our use case** - works with `map[string]interface{}`
- **Metadata for UI** - same validators provide form generation info
- **Extensible** - easy to add new field types
- **Good error messages** - we control the format

### Cons:
- More code to write initially (but it's simple code)
- No fancy features (but we don't need them)

### Alternative Considered: go-playground/validator

If we were to use an external library, `go-playground/validator` would be the choice due to its popularity and community support. However, it's overkill for our simple needs and introduces:
- External dependency
- String-based configuration (error-prone)
- Complex syntax for map validation
- Need to customize error messages

## Verdict

**Use simple programmatic validation** - no external libraries.

This approach:
1. **Follows Go idioms** - explicit over implicit, clear over clever
2. **Respects the Bento Box principle** - small, focused validation package
3. **Minimizes dependencies** - "a little copying is better than a little dependency"
4. **Provides exactly what we need** - nothing more, nothing less
5. **Is maintainable** - any Go developer can understand and modify it

The validation package will be under 300 lines, provide clear error messages, support UI metadata generation, and integrate perfectly with our YAML-based configuration system. This is the Go way.

## Implementation Plan

1. Create `pkg/validation/` package (< 250 lines)
2. Define `Schema` and `FieldValidator` interfaces
3. Implement common validators: `StringField`, `NumberField`, `BoolField`, `MapField`, `ArrayField`
4. Create schemas for each neta type (HTTP, transform, etc.)
5. Integrate with bento node parameter validation
6. Use metadata for future UI form generation

This keeps us dependency-free, maintains simplicity, and provides exactly the validation we need for our use case.