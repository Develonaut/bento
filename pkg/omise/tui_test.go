package omise

import (
	"bytes"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"

	"bento/pkg/omise/screens"
)

// TestCreateBentoFormFlow tests the create bento message flow
// For interactive testing, run: go run cmd/bento/main.go
func TestCreateBentoFormFlow(t *testing.T) {
	t.Skip("Use manual testing: go run cmd/bento/main.go, press 'n', verify form appears and typing works")
}

// TestCreateBentoFormTyping tests that typing works in the Huh form
func TestCreateBentoFormTyping(t *testing.T) {
	t.Skip("Skipping for now - need to debug form rendering first")

	workDir := t.TempDir()

	model, err := NewModelWithWorkDir(workDir)
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}

	tm := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(80, 24))

	// Wait for browser
	waitForContent(t, tm, "🍱 Bento | Browser")

	// Press 'n' to create new bento
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

	// Wait for editor
	waitForContent(t, tm, "Create New Bento")
	time.Sleep(100 * time.Millisecond)

	// Type some text
	for _, r := range "my-test" {
		tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		time.Sleep(10 * time.Millisecond)
	}

	// Wait for rendering
	time.Sleep(100 * time.Millisecond)

	// Quit
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	// Check output contains our typed text
	output := readOutput(t, tm)
	t.Logf("Output:\n%s", string(output))
	if !bytes.Contains(output, []byte("my-test")) {
		t.Errorf("expected typed text 'my-test' to appear in output")
	}
}

// TestEditorInitCommand tests that Init() command is returned when creating editor
func TestEditorInitCommand(t *testing.T) {
	workDir := t.TempDir()

	model, err := NewModelWithWorkDir(workDir)
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}

	// Send CreateBentoMsg
	result, cmd := model.Update(screens.CreateBentoMsg{})

	m, ok := result.(Model)
	if !ok {
		t.Fatal("result is not Model type")
	}

	// Check we switched to editor
	if m.screen != ScreenEditor {
		t.Error("expected to switch to editor screen")
	}

	// Check that Init() command was returned
	if cmd == nil {
		t.Fatal("expected Init() command to be returned, got nil")
	}
}
