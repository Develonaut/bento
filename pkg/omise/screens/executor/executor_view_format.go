package executor

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"bento/pkg/neta"
)

// formatResult formats the execution result with syntax highlighting
func formatResult(result interface{}) string {
	jsonStr, ok := extractJSON(result)
	if !ok {
		return jsonStr // Already a fallback message
	}
	return jsonStr
}

// extractJSON extracts and formats JSON from result
func extractJSON(result interface{}) (string, bool) {
	if result == nil {
		return "No output", false
	}

	// Type assert to neta.Result
	netaResult, ok := result.(neta.Result)
	if !ok {
		return fmt.Sprintf("%v", result), false
	}

	return marshalNetaResult(netaResult)
}

// marshalNetaResult converts neta.Result to a simple string representation
func marshalNetaResult(result neta.Result) (string, bool) {
	if result.Output == nil {
		return "No output", false
	}

	// Try to extract simple values instead of showing full JSON
	switch v := result.Output.(type) {
	case string:
		return v, true
	case int, int64, float64, bool:
		return fmt.Sprintf("%v", v), true
	case map[string]interface{}:
		return formatMapOutput(v)
	default:
		// For other types, try compact JSON
		jsonBytes, err := json.Marshal(result.Output)
		if err != nil {
			return fmt.Sprintf("%v", result.Output), false
		}
		return string(jsonBytes), true
	}
}

// formatMapOutput formats map output as string
func formatMapOutput(v map[string]interface{}) (string, bool) {
	// For objects, try to extract a simple "value" or "message" field
	if val, ok := v["value"]; ok {
		return fmt.Sprintf("%v", val), true
	}
	if msg, ok := v["message"]; ok {
		return fmt.Sprintf("%v", msg), true
	}
	// Fall back to compact JSON for small objects
	jsonBytes, err := json.Marshal(v)
	if err == nil && len(jsonBytes) < 100 {
		return string(jsonBytes), true
	}
	// For larger objects, show summary
	return fmt.Sprintf("{%d fields}", len(v)), true
}

// truncateToOneLine truncates a string to fit on one line
func truncateToOneLine(s string, maxLen int) string {
	// Remove all newlines and extra whitespace
	s = lipgloss.NewStyle().Inline(true).Render(s)

	// Find first newline if any
	for i, r := range s {
		if r == '\n' || r == '\r' {
			s = s[:i]
			break
		}
	}

	// Truncate to max length
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}
