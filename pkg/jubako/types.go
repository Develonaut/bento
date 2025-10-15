// Package jubako provides storage and file management for Bentos.
// Jubako (重箱) means "stacked boxes" - a traditional Japanese food container.
package jubako

import "time"

// BentoInfo contains metadata about a bento file.
type BentoInfo struct {
	Name     string
	Path     string
	Type     string
	Modified time.Time
}

// ExecutionRecord tracks a bento execution.
type ExecutionRecord struct {
	ID        string
	Bento     string
	StartTime time.Time
	EndTime   time.Time
	Success   bool
	Error     string
	// Result holds the execution output. Uses interface{} to support
	// heterogeneous bento result types (HTTP responses, file contents, etc.)
	Result interface{}
}

// HistoryFilter filters execution history.
type HistoryFilter struct {
	Bento       string
	SuccessOnly bool
	Limit       int
}
