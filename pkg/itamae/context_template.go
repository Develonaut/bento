package itamae

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

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

// resolveSecretsInString resolves {{SECRETS.X}} placeholders from keychain.
// Returns the string with secrets resolved, or original string if error occurs.
func (ec *executionContext) resolveSecretsInString(s string) string {
	if ec.secretsManager == nil || !strings.Contains(s, "{{SECRETS.") {
		return s
	}

	resolved, err := ec.secretsManager.ResolveTemplate(s)
	if err != nil {
		// SECRET RESOLUTION FAILED - This is a CRITICAL error
		fmt.Fprintf(os.Stderr, "\n❌ ERROR: Failed to resolve secrets in template: %v\n", err)
		fmt.Fprintf(os.Stderr, "   Template: %s\n", s)
		fmt.Fprintf(os.Stderr, "   This will likely cause authentication failures!\n\n")
		return s
	}
	return resolved
}

// executeGoTemplate parses and executes a Go template string.
// Returns the interpolated string, or the input if parsing/execution fails.
func (ec *executionContext) executeGoTemplate(s string) string {
	tmpl, err := template.New("param").Funcs(templateFuncs()).Parse(s)
	if err != nil {
		return s
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ec.nodeData); err != nil {
		return s
	}

	return buf.String()
}

// resolveString resolves template syntax in a string.
// Secrets ({{SECRETS.X}}) are resolved FIRST, then regular templates ({{.X}}).
// If the string is ONLY a template (no literal text), return the actual value.
// Otherwise, return the string interpolation.
func (ec *executionContext) resolveString(s string) interface{} {
	// Step 1: Resolve {{SECRETS.X}} placeholders from keychain
	resolvedSecrets := ec.resolveSecretsInString(s)

	// Step 2: Check if string contains Go template syntax ({{.X}})
	if !containsTemplate(resolvedSecrets) {
		return resolvedSecrets
	}

	// Step 3: Special case - if entire string is single template, return actual value
	if isExactTemplate(resolvedSecrets) {
		if val := ec.resolveExactTemplate(resolvedSecrets); val != nil {
			return val
		}
	}

	// Step 4: Parse and execute Go template (returns string interpolation)
	return ec.executeGoTemplate(resolvedSecrets)
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

// templateFuncs returns custom template functions.
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"basename": filepath.Base,
	}
}
