package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hbaldwin98/tui-writer/editor"
	"github.com/hbaldwin98/tui-writer/input"
)

func TestUpdateDoesNotBufferInsertModePrefixes(t *testing.T) {
	ed := editor.New()
	ed.SetKeymap(input.Keymap{
		input.ModeInsert: {
			"aa": input.ActionSave,
		},
	})
	m := model{editor: &ed}

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = updated.(model)

	if got := m.editor.Value(); got != "a" {
		t.Fatalf("expected insert key to reach textarea, got %q", got)
	}
}

func TestUpdateAllowsEnterInInsertMode(t *testing.T) {
	ed := editor.New()
	m := model{editor: &ed}

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)

	if got := m.editor.Value(); got != "\n" {
		t.Fatalf("expected enter to insert newline, got %q", got)
	}
}
