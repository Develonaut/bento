package components

// KeyMap aggregates all keymaps for the application.
// Note: This is kept for organizational purposes, but screens should
// use their specific keymaps directly rather than going through wrappers.
type KeyMap struct {
	Navigation NavigationKeyMap
	Browser    BrowserKeyMap
	Editor     EditorKeyMap
	Settings   SettingsKeyMap
	Global     GlobalKeyMap
}

// NewKeyMap creates a new centralized keymap
func NewKeyMap() KeyMap {
	return KeyMap{
		Navigation: NewNavigationKeyMap(),
		Browser:    NewBrowserKeyMap(),
		Editor:     NewEditorKeyMap(),
		Settings:   NewSettingsKeyMap(),
		Global:     NewGlobalKeyMap(),
	}
}
