package executor

import (
	"fmt"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

// CopyResultCmd copies result to clipboard and returns feedback message
func CopyResultCmd(result, bentoName, errorMsg string, success bool) tea.Msg {
	content := buildClipboardContent(result, bentoName, errorMsg, success)
	if content == "" {
		return CopyResultMsg("No output or error to copy")
	}

	if err := clipboard.WriteAll(content); err != nil {
		return CopyResultMsg(fmt.Sprintf("Failed to copy: %s", err.Error()))
	}

	return CopyResultMsg("✓ Copied to clipboard!")
}

// CopyEntireViewCmd copies the entire view content to clipboard
func CopyEntireViewCmd(viewContent string) tea.Msg {
	if viewContent == "" {
		return CopyResultMsg("No view content to copy")
	}

	if err := clipboard.WriteAll(viewContent); err != nil {
		return CopyResultMsg(fmt.Sprintf("Failed to copy view: %s", err.Error()))
	}

	return CopyResultMsg("✓ Entire view copied to clipboard!")
}

// buildClipboardContent formats content for clipboard
func buildClipboardContent(result, bentoName, errorMsg string, success bool) string {
	if success && result != "" {
		return fmt.Sprintf("Bento: %s\n\nStatus: Success\n\nOutput:\n%s", bentoName, result)
	}
	if !success && errorMsg != "" {
		return fmt.Sprintf("Bento: %s\n\nStatus: Failed\n\nError:\n%s", bentoName, errorMsg)
	}
	if result != "" {
		return fmt.Sprintf("Bento: %s\n\nOutput:\n%s", bentoName, result)
	}
	if errorMsg != "" {
		return fmt.Sprintf("Bento: %s\n\nError:\n%s", bentoName, errorMsg)
	}
	return ""
}
