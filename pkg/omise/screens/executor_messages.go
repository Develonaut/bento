package screens

import (
	"bento/pkg/neta"
)

// ExecutionStartMsg signals bento execution has started
type ExecutionStartMsg struct{}

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
