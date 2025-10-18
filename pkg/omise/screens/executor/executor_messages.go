package executor

import (
	"time"

	"bento/pkg/neta"
)

// ExecutionStartMsg signals bento execution has started
type ExecutionStartMsg struct{}

// ExecutionInitMsg carries the definition for node initialization
type ExecutionInitMsg struct {
	Definition neta.Definition
}

// ExecutionProgressMsg updates execution progress
type ExecutionProgressMsg struct {
	NodeName string
	Progress float64 // 0.0 to 1.0
	Status   string
}

// ExecutionCompleteMsg signals bento execution finished
type ExecutionCompleteMsg struct {
	Success bool
	Result  neta.Result
	Error   error
}

// ExecutionErrorMsg signals an execution error
type ExecutionErrorMsg struct {
	Error error
}

// CopyResultMsg provides feedback for copy operation
type CopyResultMsg string

// NodeStartedMsg signals a node has started execution
type NodeStartedMsg struct {
	Path     string
	Name     string
	NodeType string
}

// NodeCompletedMsg signals a node has finished execution
type NodeCompletedMsg struct {
	Path     string
	Duration time.Duration
	Error    error
}

// GraphStateUpdateMsg carries ExecutionGraphStore state updates to the UI
type GraphStateUpdateMsg struct {
	State neta.ExecutionGraphState
}
