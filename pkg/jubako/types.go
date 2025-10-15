// Package jubako provides storage and file management for Bento workflows.
// Jubako (重箱) means "stacked boxes" - a traditional Japanese food container.
package jubako

import "time"

// WorkflowInfo contains metadata about a workflow file.
type WorkflowInfo struct {
	Name     string
	Path     string
	Type     string
	Modified time.Time
}

// ExecutionRecord tracks a workflow execution.
type ExecutionRecord struct {
	ID        string
	Workflow  string
	StartTime time.Time
	EndTime   time.Time
	Success   bool
	Error     string
	// Result holds the execution output. Uses interface{} to support
	// heterogeneous workflow result types (HTTP responses, file contents, etc.)
	Result interface{}
}

// HistoryFilter filters execution history.
type HistoryFilter struct {
	Workflow    string
	SuccessOnly bool
	Limit       int
}
