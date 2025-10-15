package screens

// BentoSelectedMsg signals that a bento was selected for execution (legacy)
// Deprecated: Use WorkflowSelectedMsg instead
type BentoSelectedMsg struct {
	Name string
	Path string
}

// WorkflowSelectedMsg signals that a bento was selected for execution
type WorkflowSelectedMsg struct {
	Name string
	Path string
}

// EditBentoMsg signals user wants to edit a bento
type EditBentoMsg struct {
	Name string
	Path string
}

// CreateBentoMsg signals user wants to create new bento
type CreateBentoMsg struct{}

// CopyBentoMsg signals user wants to copy a bento
type CopyBentoMsg struct {
	Name string
	Path string
}

// DeleteBentoMsg signals user confirmed deletion
type DeleteBentoMsg struct {
	Name string
	Path string
}

// BentoListRefreshMsg signals bento list should reload
type BentoListRefreshMsg struct{}

// BentoOperationCompleteMsg signals operation completed
type BentoOperationCompleteMsg struct {
	Operation string // "copy", "delete", "create"
	Success   bool
	Error     error
}
