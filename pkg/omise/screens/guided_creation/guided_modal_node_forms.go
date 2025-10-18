package guided_creation

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"

	"bento/pkg/neta"
)

// createHTTPNodeForm creates a form for HTTP nodes
func (m *GuidedModal) createHTTPNodeForm() *huh.Form {
	var nodeName, url, method, headers, body string
	method = "GET" // Default

	// Calculate form width: total width - preview width (40) - margins
	formWidth := m.width - 40 - 8

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("node_name").
				Title("Name").
				Description("Name for this HTTP node").
				Value(&nodeName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("name is required")
					}
					return nil
				}),

			huh.NewInput().
				Key("url").
				Title("URL").
				Description("HTTP(S) URL to request").
				Value(&url).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("url is required")
					}
					return nil
				}),

			huh.NewSelect[string]().
				Key("method").
				Title("Method").
				Description("HTTP method").
				Options(
					huh.NewOption("GET", "GET"),
					huh.NewOption("POST", "POST"),
					huh.NewOption("PUT", "PUT"),
					huh.NewOption("DELETE", "DELETE"),
					huh.NewOption("PATCH", "PATCH"),
					huh.NewOption("HEAD", "HEAD"),
					huh.NewOption("OPTIONS", "OPTIONS"),
				).
				Value(&method),

			huh.NewInput().
				Key("headers").
				Title("Headers").
				Description("HTTP headers (optional, e.g., Content-Type: application/json)").
				Value(&headers),

			huh.NewInput().
				Key("body").
				Title("Body").
				Description("Request body for POST/PUT/PATCH (optional)").
				Value(&body),
		).Title("HTTP Node:"),
	).
		WithWidth(formWidth).
		WithShowHelp(false).
		WithShowErrors(false)
}

// createJQNodeForm creates a form for transform.jq nodes
func (m *GuidedModal) createJQNodeForm() *huh.Form {
	var nodeName, query, input string

	// Calculate form width: total width - preview width (40) - margins
	formWidth := m.width - 40 - 8

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("node_name").
				Title("Name").
				Description("Name for this transform node").
				Value(&nodeName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("name is required")
					}
					return nil
				}),

			huh.NewInput().
				Key("query").
				Title("Query").
				Description("jq query to transform data (e.g., '.data | .[] | select(.active)')").
				Value(&query).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("query is required")
					}
					return nil
				}),

			huh.NewInput().
				Key("input").
				Title("Input").
				Description("Static input data (optional, uses previous node output if not specified)").
				Value(&input),
		).Title("Transform (jq) Node:"),
	).
		WithWidth(formWidth).
		WithShowHelp(false).
		WithShowErrors(false)
}

// createFileWriteNodeForm creates a form for file.write nodes
func (m *GuidedModal) createFileWriteNodeForm() *huh.Form {
	var nodeName, path, content string

	// Calculate form width: total width - preview width (40) - margins
	formWidth := m.width - 40 - 8

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("node_name").
				Title("Name").
				Description("Name for this file write node").
				Value(&nodeName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("name is required")
					}
					return nil
				}),

			huh.NewInput().
				Key("path").
				Title("Path").
				Description("File path to write to (e.g., '/tmp/output.txt')").
				Value(&path).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("path is required")
					}
					return nil
				}),

			huh.NewInput().
				Key("content").
				Title("Content").
				Description("Content to write (optional, uses previous node output if not specified)").
				Value(&content),
		).Title("File Write Node:"),
	).
		WithWidth(formWidth).
		WithShowHelp(false).
		WithShowErrors(false)
}

// createSequenceNodeForm creates a form for group.sequence nodes
func (m *GuidedModal) createSequenceNodeForm() *huh.Form {
	var nodeName string

	// Calculate form width: total width - preview width (40) - margins
	formWidth := m.width - 40 - 8

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("node_name").
				Title("Name").
				Description("Name for this sequence group").
				Value(&nodeName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("name is required")
					}
					return nil
				}),
		).Title("Sequence Group Node:"),
	).
		WithWidth(formWidth).
		WithShowHelp(false).
		WithShowErrors(false)
}

// createParallelNodeForm creates a form for group.parallel nodes
func (m *GuidedModal) createParallelNodeForm() *huh.Form {
	var nodeName string

	// Calculate form width: total width - preview width (40) - margins
	formWidth := m.width - 40 - 8

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("node_name").
				Title("Name").
				Description("Name for this parallel group").
				Value(&nodeName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("name is required")
					}
					return nil
				}),
		).Title("Parallel Group Node:"),
	).
		WithWidth(formWidth).
		WithShowHelp(false).
		WithShowErrors(false)
}

// updateCurrentNodeFromNodeForm updates the current node with values from a node-type-specific form
func (m *GuidedModal) updateCurrentNodeFromNodeForm(nodeType string) {
	if m.currentNode == nil {
		m.currentNode = &neta.Definition{
			Version:    "1.0",
			Parameters: make(map[string]interface{}),
		}
	}

	// Set node type
	m.currentNode.Type = nodeType

	// Set node name
	if name := m.form.GetString("node_name"); name != "" {
		m.currentNode.Name = name
	}

	// Set parameters based on node type
	switch nodeType {
	case "http":
		if url := m.form.GetString("url"); url != "" {
			m.currentNode.Parameters["url"] = url
		}
		if method := m.form.GetString("method"); method != "" {
			m.currentNode.Parameters["method"] = method
		}
		if headers := m.form.GetString("headers"); headers != "" {
			m.currentNode.Parameters["headers"] = parseHeaders(headers)
		}
		if body := m.form.GetString("body"); body != "" {
			m.currentNode.Parameters["body"] = body
		}

	case "transform.jq", "jq":
		if query := m.form.GetString("query"); query != "" {
			m.currentNode.Parameters["query"] = query
		}
		if input := m.form.GetString("input"); input != "" {
			m.currentNode.Parameters["input"] = input
		}

	case "file.write":
		if path := m.form.GetString("path"); path != "" {
			m.currentNode.Parameters["path"] = path
		}
		if content := m.form.GetString("content"); content != "" {
			m.currentNode.Parameters["content"] = content
		}

	case "group.sequence", "sequence":
		// Sequence nodes don't have parameters beyond child nodes
		// which will be handled separately

	case "group.parallel", "parallel":
		// Parallel nodes don't have parameters beyond child nodes
		// which will be handled separately
	}
}

// parseHeaders parses a header string into a map
// Supports formats:
// - "Key: Value"
// - "Key1: Value1, Key2: Value2"
// - "Key1: Value1\nKey2: Value2"
func parseHeaders(headerStr string) map[string]string {
	headers := make(map[string]string)
	if headerStr == "" {
		return headers
	}

	// Split by newlines or commas
	var lines []string
	if strings.Contains(headerStr, "\n") {
		lines = strings.Split(headerStr, "\n")
	} else if strings.Contains(headerStr, ",") {
		lines = strings.Split(headerStr, ",")
	} else {
		lines = []string{headerStr}
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Split by first colon
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" {
				headers[key] = value
			}
		}
	}

	return headers
}
