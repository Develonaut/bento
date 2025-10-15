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

	// Wait for initial render - look for Browser header
	waitForContent(t, tm, "🍱 Bento | Browser")

	// Navigate down to skip the "+ Create New Bento" item
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	time.Sleep(50 * time.Millisecond)

	// Press Enter to select first actual bento
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyEnter,
		Runes: []rune{'\r'},
	})

	// Wait for Executor screen - look for both header and content
	waitForContent(t, tm, "Bento Executor")

	// Send quit to finish
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'q'},
	})

	// Verify final model state
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
	fm := tm.FinalModel(t)
	model, ok := fm.(Model)
	if !ok {
		t.Fatal("Final model is not Model type")
	}

	if model.screen != ScreenExecutor {
		t.Errorf("Expected screen to be Executor, got %v", model.screen)
	}

	if !model.executor.IsRunning() {
		t.Error("Expected executor to be running after bento selection")
	}
}

// TestScreenNavigation tests tab cycling through all screens
func TestScreenNavigation(t *testing.T) {
	m := NewModel()
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	// Send multiple tabs to cycle through screens
	for i := 0; i < 5; i++ {
		tm.Send(tea.KeyMsg{Type: tea.KeyTab})
		time.Sleep(50 * time.Millisecond) // Give time for screen to update
	}

	// Quit
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'q'},
	})

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	// Verify final model state - should have cycled back to Browser
	fm := tm.FinalModel(t)
	model, ok := fm.(Model)
	if !ok {
		t.Fatal("Final model is not Model type")
	}
	if model.screen != ScreenBrowser {
		t.Errorf("Expected screen to be Browser after full cycle, got %v", model.screen)
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
			waitForContent(t, tm, "🍱 Bento | Browser")

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
			"url": "https://httpbin.org/get",
		},
	}

	if err := store.Save("test-bento", def); err != nil {
		t.Fatalf("Failed to save test bento: %v", err)
	}
}
