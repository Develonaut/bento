package components

import (
	"bento/pkg/omise/styles"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

// DirSelectedMsg signals that a directory was selected
type DirSelectedMsg struct {
	Path string
}

// DirPicker wraps bubbles/filepicker for directory selection
type DirPicker struct {
	filepicker.Model
	defaultDir string
}

// Init initializes the directory picker and loads the current directory
func (dp DirPicker) Init() tea.Cmd {
	return dp.Model.Init()
}

// NewDirPicker creates a themed directory picker
// startDir is where the picker opens initially
// defaultDir is the directory to return to when reset is pressed
func NewDirPicker(startDir string, defaultDir string) DirPicker {
	fp := filepicker.New()
	fp.AllowedTypes = nil
	fp.DirAllowed = true
	fp.FileAllowed = false
	fp.ShowHidden = true // Show hidden directories like .bento
	fp.CurrentDirectory = startDir
	fp.SetHeight(15)

	// Apply theme styling
	fp = applyDirPickerStyles(fp)

	return DirPicker{
		Model:      fp,
		defaultDir: defaultDir,
	}
}

// applyDirPickerStyles applies theme colors to filepicker
func applyDirPickerStyles(fp filepicker.Model) filepicker.Model {
	s := filepicker.DefaultStyles()
	s.Cursor = s.Cursor.Foreground(styles.Primary)
	s.Symlink = s.Symlink.Foreground(styles.Secondary)
	s.Directory = s.Directory.Foreground(styles.Primary)
	s.File = s.File.Foreground(styles.Text)
	s.Permission = s.Permission.Foreground(styles.Muted)
	s.Selected = s.Selected.Foreground(styles.Success)
	s.FileSize = s.FileSize.Foreground(styles.Muted)
	fp.Styles = s
	return fp
}

// Update handles directory picker messages
func (dp DirPicker) Update(msg tea.Msg) (DirPicker, tea.Cmd) {
	var cmd tea.Cmd

	// Handle 's' key to select current directory
	// (Enter is used by filepicker for navigation into directories)
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "s" {
			// Select the current directory
			return dp, func() tea.Msg {
				return DirSelectedMsg{Path: dp.Model.CurrentDirectory}
			}
		}
	}

	// Let filepicker handle Enter for navigation
	dp.Model, cmd = dp.Model.Update(msg)

	// Check if a directory was selected from the list
	if didSelect, path := dp.Model.DidSelectFile(msg); didSelect {
		return dp, func() tea.Msg {
			return DirSelectedMsg{Path: path}
		}
	}

	return dp, cmd
}

// RebuildStyles updates the picker styles with current theme colors
func (dp DirPicker) RebuildStyles() DirPicker {
	dp.Model = applyDirPickerStyles(dp.Model)
	return dp
}

// ResetToDefault resets the directory picker to the default directory
func (dp DirPicker) ResetToDefault() DirPicker {
	dp.Model.CurrentDirectory = dp.defaultDir
	return dp
}
