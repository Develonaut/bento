package itamae

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

// executionContext holds data passed between nodes during execution.
type executionContext struct {
	nodeData map[string]interface{} // Data from each executed node
}

// newExecutionContext creates a new execution context.
// Initializes nodeData with environment variables so templates can access them.
func newExecutionContext() *executionContext {
	nodeData := make(map[string]interface{})

	// Load all environment variables into context
	// This allows templates like {{.FIGMA_API_URL}} to work
	for _, env := range os.Environ() {
		// Split on first '=' to handle values that contain '='
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			nodeData[parts[0]] = parts[1]
		}
	}

	return &executionContext{
		nodeData: nodeData,
	}
}

// set stores output from a node.
func (ec *executionContext) set(nodeID string, data interface{}) {
	ec.nodeData[nodeID] = data
}

// resolveValue recursively resolves template strings in a value.
func (ec *executionContext) resolveValue(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		return ec.resolveString(v)
	case map[string]interface{}:
		return ec.resolveMap(v)
	case []interface{}:
		return ec.resolveSlice(v)
	default:
		return value
	}
}

// resolveString resolves template syntax in a string.
// If the string is ONLY a template (no literal text), return the actual value.
// Otherwise, return the string interpolation.
func (ec *executionContext) resolveString(s string) interface{} {
	// Check if string contains template syntax
	if !containsTemplate(s) {
		return s
	}

	// Special case: if the entire string is a single template expression,
	// try to return the actual value instead of string representation.
	// This allows passing arrays/maps through templates.
	if isExactTemplate(s) {
		if val := ec.resolveExactTemplate(s); val != nil {
			return val
		}
	}

	// Parse and execute template (returns string interpolation)
	tmpl, err := template.New("param").Parse(s)
	if err != nil {
		return s // Return original if parse fails
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ec.nodeData); err != nil {
		return s // Return original if execute fails
	}

	return buf.String()
}

// isExactTemplate checks if a string is EXACTLY one template (no literal text).
func isExactTemplate(s string) bool {
	trimmed := strings.TrimSpace(s)
	return strings.HasPrefix(trimmed, "{{") && strings.HasSuffix(trimmed, "}}")
}

// resolveExactTemplate resolves a template that is exactly one expression.
// Returns the actual value from context (array, map, etc.) instead of string.
func (ec *executionContext) resolveExactTemplate(s string) interface{} {
	// Extract the expression between {{ and }}
	trimmed := strings.TrimSpace(s)
	expr := strings.TrimSpace(trimmed[2 : len(trimmed)-2])

	// Handle "index . \"key1\" \"key2\"..." syntax
	if strings.HasPrefix(expr, "index .") {
		return ec.resolveIndexExpression(expr)
	}

	// Handle simple ".key" or ".key.subkey" syntax
	if strings.HasPrefix(expr, ".") {
		return ec.resolveDotExpression(expr[1:]) // Remove leading dot
	}

	return nil
}

// resolveIndexExpression resolves {{index . "key1" "key2"}} expressions.
func (ec *executionContext) resolveIndexExpression(expr string) interface{} {
	// Parse: index . "key1" "key2" ...
	parts := strings.Fields(expr)
	if len(parts) < 3 || parts[0] != "index" || parts[1] != "." {
		return nil
	}

	// Extract keys (remove quotes)
	keys := make([]string, 0, len(parts)-2)
	for i := 2; i < len(parts); i++ {
		key := strings.Trim(parts[i], "\"")
		keys = append(keys, key)
	}

	// Navigate through the context
	var current interface{} = ec.nodeData
	for _, key := range keys {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		current = m[key]
	}

	return current
}

// resolveDotExpression resolves {{.key.subkey}} expressions.
func (ec *executionContext) resolveDotExpression(expr string) interface{} {
	keys := strings.Split(expr, ".")
	var current interface{} = ec.nodeData

	for _, key := range keys {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		current = m[key]
	}

	return current
}

// resolveMap resolves templates in a map.
func (ec *executionContext) resolveMap(m map[string]interface{}) map[string]interface{} {
	resolved := make(map[string]interface{})
	for k, v := range m {
		resolved[k] = ec.resolveValue(v)
	}
	return resolved
}

// resolveSlice resolves templates in a slice.
func (ec *executionContext) resolveSlice(s []interface{}) []interface{} {
	resolved := make([]interface{}, len(s))
	for i, v := range s {
		resolved[i] = ec.resolveValue(v)
	}
	return resolved
}

// containsTemplate checks if a string contains template syntax.
func containsTemplate(s string) bool {
	return len(s) > 4 && strings.Contains(s, "{{") && strings.Contains(s, "}}")
}

// copy creates a shallow copy of the execution context.
// Note: This performs a shallow copy - the nodeData map is copied,
// but the values within the map are not deep-copied. This is intentional
// for performance and works correctly because node outputs are immutable
// after being set.
func (ec *executionContext) copy() *executionContext {
	newCtx := newExecutionContext()
	for k, v := range ec.nodeData {
		newCtx.nodeData[k] = v
	}
	return newCtx
}

// toMap converts the context to a map for external use.
func (ec *executionContext) toMap() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range ec.nodeData {
		result[k] = v
	}
	return result
}

// String returns a string representation for debugging.
func (ec *executionContext) String() string {
	return fmt.Sprintf("executionContext{nodes: %d}", len(ec.nodeData))
}
