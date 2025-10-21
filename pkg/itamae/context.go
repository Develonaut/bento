package itamae

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/Develonaut/bento/pkg/wasabi"
)

// executionContext holds data passed between nodes during execution.
type executionContext struct {
	nodeData       map[string]interface{} // Data from each executed node
	secretsManager *wasabi.Manager        // Secrets manager for {{SECRETS.X}} resolution
	depth          int                    // Nesting depth for logging indentation
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

	// Initialize secrets manager for {{SECRETS.X}} resolution
	// If initialization fails, proceed without secrets (logged as warning)
	secretsMgr, err := wasabi.NewManager()
	if err != nil {
		// Note: We don't fail here because secrets might not be needed
		// The error will surface when trying to resolve {{SECRETS.X}} if used
		secretsMgr = nil
	}

	return &executionContext{
		nodeData:       nodeData,
		secretsManager: secretsMgr,
		depth:          0,
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
// Secrets ({{SECRETS.X}}) are resolved FIRST, then regular templates ({{.X}}).
// If the string is ONLY a template (no literal text), return the actual value.
// Otherwise, return the string interpolation.
func (ec *executionContext) resolveString(s string) interface{} {
	// Step 1: Resolve {{SECRETS.X}} placeholders from keychain
	// This happens BEFORE Go template resolution to prevent template errors
	// and maintain strict separation of concerns
	resolvedSecrets := s
	if ec.secretsManager != nil && strings.Contains(s, "{{SECRETS.") {
		var err error
		resolvedSecrets, err = ec.secretsManager.ResolveTemplate(s)
		if err != nil {
			// SECRET RESOLUTION FAILED - This is a CRITICAL error
			// Print to stderr so user sees it immediately (not buried in logs)
			// Returning unresolved template will likely cause downstream errors
			fmt.Fprintf(os.Stderr, "\n❌ ERROR: Failed to resolve secrets in template: %v\n", err)
			fmt.Fprintf(os.Stderr, "   Template: %s\n", s)
			fmt.Fprintf(os.Stderr, "   This will likely cause authentication failures!\n\n")
			return s
		}
	}

	// Step 2: Check if string contains Go template syntax ({{.X}})
	if !containsTemplate(resolvedSecrets) {
		return resolvedSecrets
	}

	// Step 3: Special case - if entire string is single template, return actual value
	// This allows passing arrays/maps through templates
	if isExactTemplate(resolvedSecrets) {
		if val := ec.resolveExactTemplate(resolvedSecrets); val != nil {
			return val
		}
	}

	// Step 4: Parse and execute Go template (returns string interpolation)
	tmpl, err := template.New("param").Parse(resolvedSecrets)
	if err != nil {
		return resolvedSecrets // Return secrets-resolved string if parse fails
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ec.nodeData); err != nil {
		return resolvedSecrets // Return secrets-resolved string if execute fails
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
// after being set. The secrets manager is shared across copies.
func (ec *executionContext) copy() *executionContext {
	newCtx := newExecutionContext()
	for k, v := range ec.nodeData {
		newCtx.nodeData[k] = v
	}
	// Share the same secrets manager (thread-safe)
	newCtx.secretsManager = ec.secretsManager
	// Preserve depth
	newCtx.depth = ec.depth
	return newCtx
}

// withDepth returns a copy of the context with incremented depth.
func (ec *executionContext) withDepth(increment int) *executionContext {
	newCtx := ec.copy()
	newCtx.depth = ec.depth + increment
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
