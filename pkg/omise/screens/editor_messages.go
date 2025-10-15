package screens

// BentoNameEnteredMsg signals bento name was entered
type BentoNameEnteredMsg struct {
	Name string
}

// NodeTypeSelectedMsg signals node type was selected from Pantry
type NodeTypeSelectedMsg struct {
	Type string
}

// NodeConfiguredMsg signals node parameters were configured
type NodeConfiguredMsg struct {
	Type       string
	Name       string
	Parameters map[string]interface{}
}

// EditorSavedMsg signals bento was saved successfully
type EditorSavedMsg struct {
	Name string
}

// EditorSaveErrorMsg signals save error occurred
type EditorSaveErrorMsg struct {
	Error error
}

// EditorCancelledMsg signals editor was cancelled
type EditorCancelledMsg struct{}

// ReturnToBrowserMsg signals return to browser screen
type ReturnToBrowserMsg struct{}

// EditNodeMsg signals node should be edited
type EditNodeMsg struct {
	Index int
	Type  string
}

// RunBentoFromEditorMsg signals run bento from editor
type RunBentoFromEditorMsg struct {
	Def interface{}
}
