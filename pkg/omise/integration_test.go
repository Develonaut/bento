package omise

import (
	"bytes"
	"io"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"

	"bento/pkg/jubako"
	"bento/pkg/neta"
)

// TestBrowserToExecutorFlow tests selecting and executing a bento
func TestBrowserToExecutorFlow(t *testing.T) {
	// Create test bento in temp directory
	workDir := t.TempDir()
	createTestBento(t, workDir)

	m, err := NewModelWithWorkDir(workDir)
	if err != nil {
		t.Fatalf("NewModelWithWorkDir error: %v", err)
	}

	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	// Wait for initial render - look for app header
	waitForContent(t, tm, "🍱 Bento v0.1.0")

	// Navigate down to skip the "+ Create New Bento" item
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	time.Sleep(50 * time.Millisecond)

	// Press 'r' to quick run first actual bento (bypasses action menu)
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'r'},
	})

	// Give time for screen transition
	time.Sleep(100 * time.Millisecond)

	// Send quit to finish
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'q'},
	})

	// Verify final model state
	tm.WaitFinished(t, teatest.WithFinalTimeout(2*time.Second))
	fm := tm.FinalModel(t)
	model, ok := fm.(Model)
	if !ok {
		t.Fatal("Final model is not Model type")
	}

	// Should be on executor screen after pressing 'r'
	if model.screen != ScreenExecutor {
		t.Errorf("Expected screen to be Executor, got %v", model.screen)
	}
}

// TestQuitBehavior tests that 'q' and ctrl+c quit the app
func TestQuitBehavior(t *testing.T) {
	tests := []struct {
		name string
		key  tea.KeyMsg
	}{
		{
			name: "quit with q",
			key: tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{'q'},
			},
		},
		{
			name: "quit with ctrl+c",
			key: tea.KeyMsg{
				Type: tea.KeyCtrlC,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel()
			tm := teatest.NewTestModel(
				t, m,
				teatest.WithInitialTermSize(80, 24),
			)

			// Wait for initial render
			waitForContent(t, tm, "🍱 Bento v0.1.0")

			// Send quit key
			tm.Send(tt.key)

			// Verify app finishes
			tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

			// Verify final output contains goodbye message
			output := readOutput(t, tm)
			if !bytes.Contains(output, []byte("Thanks for using Bento!")) {
				t.Error("Expected goodbye message in output")
			}
		})
	}
}

// TestHelpScreen tests that '?' opens help screen
func TestHelpScreen(t *testing.T) {
	m := NewModel()
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	// Give initial render time
	time.Sleep(50 * time.Millisecond)

	// Press '?' to open help
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'?'},
	})

	// Give screen time to update
	time.Sleep(50 * time.Millisecond)

	// Quit
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'q'},
	})

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	// Verify final model state
	fm := tm.FinalModel(t)
	model, ok := fm.(Model)
	if !ok {
		t.Fatal("Final model is not Model type")
	}
	if model.screen != ScreenHelp {
		t.Errorf("Expected screen to be Help, got %v", model.screen)
	}
}

// waitForContent waits for the given content to appear in output
func waitForContent(t *testing.T, tm *teatest.TestModel, content string) {
	t.Helper()
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte(content))
		},
		teatest.WithCheckInterval(time.Millisecond*10),
		teatest.WithDuration(time.Second*3),
	)
}

// waitForContentOrTimeout waits for content with custom timeout
func waitForContentOrTimeout(t *testing.T, tm *teatest.TestModel, content string, timeout time.Duration) {
	t.Helper()
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte(content))
		},
		teatest.WithCheckInterval(time.Millisecond*50),
		teatest.WithDuration(timeout),
	)
}

// readOutput reads all output from the test model
func readOutput(t *testing.T, tm *teatest.TestModel) []byte {
	t.Helper()
	output, err := io.ReadAll(tm.FinalOutput(t))
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}
	return output
}

// createTestBento creates a simple test bento file
func createTestBento(t *testing.T, workDir string) {
	t.Helper()

	store, err := jubako.NewStore(workDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	def := neta.Definition{
		Version: "1.0",
		Type:    "http",
		Name:    "test-bento",
		Parameters: map[string]interface{}{
			"method": "GET",
			"url":    "https://httpbin.org/get",
		},
	}

	if err := store.Save("test-bento", def); err != nil {
		t.Fatalf("Failed to save test bento: %v", err)
	}
}

// TestExecutorExecution tests the executor actually runs bentos
func TestExecutorExecution(t *testing.T) {
	t.Skip("Skipping for now - debugging screen transitions")

	workDir := t.TempDir()
	createTestBento(t, workDir)

	m, err := NewModelWithWorkDir(workDir)
	if err != nil {
		t.Fatalf("NewModelWithWorkDir error: %v", err)
	}

	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	// Wait for browser
	waitForContent(t, tm, "Browser")

	// Navigate to bento and run it
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	time.Sleep(100 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

	// Wait for executor screen
	time.Sleep(200 * time.Millisecond)

	// Debug: print what we got
	output := readCurrentOutput(t, tm)
	t.Logf("Current output:\n%s", string(output))

	waitForContent(t, tm, "Executor")

	// Wait for execution to complete
	waitForContentOrTimeout(t, tm, "Success", 10*time.Second)

	// Quit
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

// readCurrentOutput reads current output from test model
func readCurrentOutput(t *testing.T, tm *teatest.TestModel) []byte {
	t.Helper()
	buf := make([]byte, 8192)
	n, _ := tm.Output().Read(buf)
	return buf[:n]
}

// TestBrowserToEditorFlow tests pressing 'e' on a bento to open the editor
// TODO: This test is skipped until the editor feature is implemented
func TestBrowserToEditorFlow(t *testing.T) {
	t.Skip("Editor feature not yet implemented - this test serves as a specification for future implementation")

	// Create test bento in temp directory
	workDir := t.TempDir()
	createTestBento(t, workDir)

	m, err := NewModelWithWorkDir(workDir)
	if err != nil {
		t.Fatalf("NewModelWithWorkDir error: %v", err)
	}

	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	// Wait for initial render - look for app header
	waitForContent(t, tm, "🍱 Bento v0.1.0")

	// Navigate down to skip the "+ Create New Bento" item and select first bento
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	time.Sleep(50 * time.Millisecond)

	// Press 'e' to edit the selected bento
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'e'},
	})

	// Give time for screen transition to editor
	time.Sleep(100 * time.Millisecond)

	// Quit to finish test
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'q'},
	})

	// Verify final model state
	tm.WaitFinished(t, teatest.WithFinalTimeout(2*time.Second))
	fm := tm.FinalModel(t)
	model, ok := fm.(Model)
	if !ok {
		t.Fatal("Final model is not Model type")
	}

	// Should be on editor screen after pressing 'e'
	if model.screen != ScreenEditor {
		t.Errorf("Expected screen to be Editor, got %v", model.screen)
	}

	// When editor is implemented, this test should also verify:
	// - Editor shows the bento name
	// - Editor displays the bento's nodes/configuration
	// - User can navigate and edit the bento
}
