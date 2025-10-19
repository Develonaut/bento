package itamae

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// executionContext holds data passed between nodes during execution.
type executionContext struct {
	nodeData map[string]interface{} // Data from each executed node
}

// newExecutionContext creates a new execution context.
func newExecutionContext() *executionContext {
	return &executionContext{
		nodeData: make(map[string]interface{}),
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
func (ec *executionContext) resolveString(s string) interface{} {
	// Check if string contains template syntax
	if !containsTemplate(s) {
		return s
	}

	// Parse and execute template
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

// copy creates a deep copy of the execution context.
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
