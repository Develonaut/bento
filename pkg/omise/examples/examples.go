// Package examples provides built-in example bentos
package examples

import (
	_ "embed"
	"fmt"

	"bento/pkg/jubako"
	"bento/pkg/neta"
)

//go:embed templates/http-get.bento.yaml
var httpGetExample string

//go:embed templates/http-post.bento.yaml
var httpPostExample string

//go:embed templates/transform-jq.bento.yaml
var transformJQExample string

//go:embed templates/conditional-if.bento.yaml
var conditionalIfExample string

//go:embed templates/loop-for.bento.yaml
var loopForExample string

//go:embed templates/sequence.bento.yaml
var sequenceExample string

//go:embed templates/complete-api-workflow.bento.yaml
var completeWorkflowExample string

// Example represents an example bento
type Example struct {
	ID          string
	Name        string
	Description string
	Category    string
	Content     string
}

// Category groups examples
type Category struct {
	Name     string
	Examples []Example
}

// GetAll returns all examples organized by category
func GetAll() []Category {
	return []Category{
		createHTTPRequestExamples(),
		createDataTransformationExamples(),
		createControlFlowExamples(),
		createCompleteWorkflowExamples(),
	}
}

// createHTTPRequestExamples creates HTTP request examples
func createHTTPRequestExamples() Category {
	return Category{
		Name: "HTTP Requests",
		Examples: []Example{
			{
				ID:          "http-get",
				Name:        "Simple GET Request",
				Description: "Fetch data from an API endpoint",
				Category:    "HTTP Requests",
				Content:     httpGetExample,
			},
			{
				ID:          "http-post",
				Name:        "POST with JSON Body",
				Description: "Send data to an API endpoint",
				Category:    "HTTP Requests",
				Content:     httpPostExample,
			},
		},
	}
}

// createDataTransformationExamples creates data transformation examples
func createDataTransformationExamples() Category {
	return Category{
		Name: "Data Transformation",
		Examples: []Example{
			{
				ID:          "transform-jq",
				Name:        "Extract User IDs",
				Description: "Use JQ to filter and transform JSON data",
				Category:    "Data Transformation",
				Content:     transformJQExample,
			},
		},
	}
}

// createControlFlowExamples creates control flow examples
func createControlFlowExamples() Category {
	return Category{
		Name: "Control Flow",
		Examples: []Example{
			{
				ID:          "conditional-if",
				Name:        "Check Success Status",
				Description: "Execute different actions based on conditions",
				Category:    "Control Flow",
				Content:     conditionalIfExample,
			},
			{
				ID:          "loop-for",
				Name:        "Process Multiple Users",
				Description: "Iterate over a list of items",
				Category:    "Control Flow",
				Content:     loopForExample,
			},
		},
	}
}

// createCompleteWorkflowExamples creates complete workflow examples
func createCompleteWorkflowExamples() Category {
	return Category{
		Name: "Complete Workflows",
		Examples: []Example{
			{
				ID:          "sequence",
				Name:        "Multi-Step Workflow",
				Description: "Chain multiple operations together",
				Category:    "Complete Workflows",
				Content:     sequenceExample,
			},
			{
				ID:          "complete-api-workflow",
				Name:        "Complete API Workflow",
				Description: "Full example with HTTP, transform, and conditionals",
				Category:    "Complete Workflows",
				Content:     completeWorkflowExample,
			},
		},
	}
}

// Get returns a specific example by ID
func Get(id string) (*Example, error) {
	categories := GetAll()
	for _, cat := range categories {
		for _, ex := range cat.Examples {
			if ex.ID == id {
				return &ex, nil
			}
		}
	}
	return nil, fmt.Errorf("example not found: %s", id)
}

// Parse parses an example into a Definition
func Parse(ex Example) (neta.Definition, error) {
	parser := jubako.NewParser()
	return parser.ParseBytes([]byte(ex.Content))
}

// List returns all examples as a flat list
func List() []Example {
	categories := GetAll()
	examples := []Example{}

	for _, cat := range categories {
		examples = append(examples, cat.Examples...)
	}

	return examples
}
