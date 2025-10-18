# Phase 3: Pre-Migration - User Bento Conversion

**IMPORTANT**: Do this BEFORE touching any parser code!

---

## Step 0: Backup and Convert User Bentos

### User's Current Bentos

Located in `~/.bento/bentos/`:
1. `hello-world-file.bento.yaml` - File writing workflow
2. `hello-world-http.bento.yaml` - HTTP API workflow

### Conversion Process

#### 1. Create Backup (30 seconds)

```bash
cp ~/.bento/bentos/hello-world-file.bento.yaml ~/.bento/bentos/hello-world-file.bento.yaml.backup
cp ~/.bento/bentos/hello-world-http.bento.yaml ~/.bento/bentos/hello-world-http.bento.yaml.backup
```

#### 2. Create Conversion Tool (15 minutes)

**File**: `tools/yaml2json.go`

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"

    "gopkg.in/yaml.v3"
    "bento/pkg/neta"
)

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: go run yaml2json.go <input.yaml> <output.json>")
        os.Exit(1)
    }

    inputPath := os.Args[1]
    outputPath := os.Args[2]

    // Read YAML
    yamlData, err := os.ReadFile(inputPath)
    if err != nil {
        fmt.Printf("Error reading YAML: %v\n", err)
        os.Exit(1)
    }

    // Parse YAML
    var def neta.Definition
    if err := yaml.Unmarshal(yamlData, &def); err != nil {
        fmt.Printf("Error parsing YAML: %v\n", err)
        os.Exit(1)
    }

    // Write JSON
    jsonData, err := json.MarshalIndent(def, "", "  ")
    if err != nil {
        fmt.Printf("Error marshaling JSON: %v\n", err)
        os.Exit(1)
    }

    if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
        fmt.Printf("Error writing JSON: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("✅ Converted %s → %s\n", inputPath, outputPath)
}
```

#### 3. Convert User Bentos (2 minutes)

```bash
cd /Users/Ryan/Code/bento

# Convert hello-world-file
go run tools/yaml2json.go \
    ~/.bento/bentos/hello-world-file.bento.yaml \
    ~/.bento/bentos/hello-world-file.bento.json

# Convert hello-world-http
go run tools/yaml2json.go \
    ~/.bento/bentos/hello-world-http.bento.yaml \
    ~/.bento/bentos/hello-world-http.bento.json
```

#### 4. Verify JSON Files (5 minutes)

**Expected Output for `hello-world-file.bento.json`:**

```json
{
  "version": "1.0",
  "type": "group.sequence",
  "name": "Hello World File",
  "icon": "🍙",
  "description": "Creates a timestamped greeting and writes it to a file",
  "nodes": [
    {
      "id": "node-1",
      "version": "1.0",
      "type": "transform.jq",
      "name": "Create Greeting with Timestamp",
      "parameters": {
        "input": "{\"user\": \"Bento User\"}",
        "query": "\"Hello, \\(.user)! Welcome to Bento!\\n\\nGenerated at: \" + (now | strftime(\"%Y-%m-%d %H:%M:%S\"))"
      }
    },
    {
      "id": "node-2",
      "version": "1.0",
      "type": "file.write",
      "name": "Write Greeting to File",
      "parameters": {
        "path": "/tmp/bento-hello.txt"
      }
    },
    {
      "id": "node-3",
      "version": "1.0",
      "type": "transform.jq",
      "name": "Create Confirmation",
      "parameters": {
        "query": "\"Successfully wrote \\(.bytes) bytes to \\(.path)\""
      }
    }
  ],
  "edges": [
    {
      "id": "edge-1-2",
      "source": "node-1",
      "target": "node-2"
    },
    {
      "id": "edge-2-3",
      "source": "node-2",
      "target": "node-3"
    }
  ]
}
```

#### 5. Test YAML Bentos Still Work (5 minutes)

**BEFORE migration - ensure current code works:**

```bash
# Test with YAML (current format)
bento pack ~/.bento/bentos/hello-world-file.bento.yaml

# Expected: Creates /tmp/bento-hello.txt with greeting
cat /tmp/bento-hello.txt
```

#### 6. Update Migration Plan

Once JSON files are created and verified:

**Phase 3C Step Order:**
1. ✅ Convert user bentos (DONE in pre-migration)
2. Replace YAML parser with JSON parser
3. Test JSON bentos work with new parser
4. Convert example files in repo
5. Update tests
6. Remove YAML files (including backups)

---

## Success Criteria

- [ ] 2 backup YAML files created
- [ ] 2 JSON files created
- [ ] JSON files validated (proper structure)
- [ ] YAML bentos still run (current code works)
- [ ] Ready to switch parser

---

## After Pre-Migration

Once this is done:
1. Start Phase 3C (migration)
2. First thing: replace parser
3. Test that JSON bentos work
4. Then proceed with rest of migration
