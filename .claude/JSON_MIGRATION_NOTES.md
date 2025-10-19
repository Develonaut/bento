# JSON Migration Notes for Bentobox

**Date:** 2025-10-18
**Decision:** Use JSON instead of YAML for flow file format
**Rationale:** Simpler parsing, better Go stdlib support, no external dependencies

---

## Changes from Original Plan

### YAML → JSON

**Original (TypeScript):**
```yaml
# flow.yaml
id: my-flow
type: group
version: 1.0.0
name: My Workflow
nodes:
  - id: node-1
    type: http-request
    parameters:
      method: GET
      url: https://api.example.com/data
```

**New (Go with JSON):**
```json
{
  "id": "my-flow",
  "type": "group",
  "version": "1.0.0",
  "name": "My Workflow",
  "nodes": [
    {
      "id": "node-1",
      "type": "http-request",
      "parameters": {
        "method": "GET",
        "url": "https://api.example.com/data"
      }
    }
  ]
}
```

---

## Package Impact

### Storage Package

**Before (YAML):**
```go
import "gopkg.in/yaml.v3"

func LoadNodeFile(path string) (*nodes.Definition, error) {
    data, _ := os.ReadFile(path)
    var node nodes.Definition
    yaml.Unmarshal(data, &node)
    return &node, nil
}
```

**After (JSON):**
```go
import "encoding/json"

func LoadNodeFile(path string) (*nodes.Definition, error) {
    data, _ := os.ReadFile(path)
    var node nodes.Definition
    json.Unmarshal(data, &node)
    return &node, nil
}
```

**Advantages:**
- ✅ Standard library only (no external dependencies)
- ✅ Faster parsing (JSON is simpler than YAML)
- ✅ Better error messages
- ✅ Struct tags for validation

---

## Struct Tags

**Go structs with JSON tags:**

```go
type NodeDefinition struct {
    ID          string                 `json:"id"`
    Type        string                 `json:"type"`
    Version     string                 `json:"version"`
    ParentID    *string                `json:"parentId,omitempty"`
    Name        string                 `json:"name"`
    Position    Position               `json:"position"`
    Metadata    Metadata               `json:"metadata"`
    Parameters  map[string]interface{} `json:"parameters"`
    Fields      *FieldsConfig          `json:"fields,omitempty"`
    InputPorts  []Port                 `json:"inputPorts"`
    OutputPorts []Port                 `json:"outputPorts"`
    Nodes       []NodeDefinition       `json:"nodes,omitempty"`
    Edges       []Edge                 `json:"edges,omitempty"`
}

type Position struct {
    X float64 `json:"x"`
    Y float64 `json:"y"`
}

type Port struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    Handle string `json:"handle,omitempty"`
}

type Edge struct {
    ID           string `json:"id"`
    Source       string `json:"source"`
    Target       string `json:"target"`
    SourceHandle string `json:"sourceHandle,omitempty"`
    TargetHandle string `json:"targetHandle,omitempty"`
}
```

---

## File Extension

**Convention:** `.flow.json` or just `.json`

**Examples:**
```bash
my-workflow.flow.json
data-pipeline.flow.json
image-processor.flow.json
```

**CLI Usage:**
```bash
bentobox run workflow.flow.json
bentobox validate pipeline.flow.json
bentobox list ~/flows/*.flow.json
```

---

## Updated Dependencies

### Remove:

- ❌ `gopkg.in/yaml.v3` - Not needed
- ❌ `@atomiton/yaml` package - Skip entirely

### Use Instead:

- ✅ `encoding/json` (stdlib) - JSON parsing
- ✅ `encoding/json` (stdlib) - JSON generation

---

## Pretty Printing

**Save flow with indentation:**

```go
func SaveNodeFile(path string, node *nodes.Definition) error {
    data, err := json.MarshalIndent(node, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(path, data, 0644)
}
```

**Output:**
```json
{
  "id": "flow-1",
  "type": "group",
  "version": "1.0.0",
  "name": "Example Flow",
  "nodes": [
    {
      "id": "node-1",
      "type": "http-request"
    }
  ]
}
```

---

## Validation with JSON Schema (Optional)

**If you want JSON Schema validation:**

```go
import "github.com/xeipuuv/gojsonschema"

func ValidateFlowFile(path string) error {
    schemaLoader := gojsonschema.NewReferenceLoader("file://./schema.json")
    documentLoader := gojsonschema.NewReferenceLoader("file://" + path)

    result, err := gojsonschema.Validate(schemaLoader, documentLoader)
    if err != nil {
        return err
    }

    if !result.Valid() {
        return errors.New("Invalid flow file")
    }

    return nil
}
```

**But honestly:** Struct tags + validation package is simpler and faster

---

## Updated Package Structure

```go
pkg/storage/
├── engine.go          # IStorageEngine interface
├── filesystem.go      # Filesystem implementation
├── memory.go          # In-memory implementation
└── json.go            # JSON serialization (NOT yaml.go)
    ├── LoadNodeFile()
    ├── SaveNodeFile()
    ├── ParseJSON()
    └── ToJSON()
```

---

## Migration Effort Impact

**Time Savings:**
- YAML parsing: 1-2 days
- JSON parsing: 1 day (stdlib only)

**Net benefit:** Save 1 day, simpler code, no external deps

---

## Compatibility

**TypeScript to Go conversion:**

If users have existing YAML flows, provide a one-time conversion tool:

```go
// tools/yaml-to-json/main.go
func ConvertYAMLToJSON(yamlPath, jsonPath string) error {
    // Read YAML
    yamlData, _ := os.ReadFile(yamlPath)

    // Parse YAML
    var data interface{}
    yaml.Unmarshal(yamlData, &data)

    // Write JSON
    jsonData, _ := json.MarshalIndent(data, "", "  ")
    return os.WriteFile(jsonPath, jsonData, 0644)
}
```

**Usage:**
```bash
bentobox convert workflow.yaml workflow.flow.json
```

---

## Final Recommendation

✅ **Proceed with JSON format**

**Reasons:**
1. Simpler implementation
2. No external dependencies
3. Faster parsing
4. Better error messages
5. Industry standard for APIs/configs
6. Easier to manipulate programmatically

---

**Status:** JSON format approved
**Next Step:** Implement storage package with JSON support
