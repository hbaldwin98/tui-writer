package editor

import (
	"testing"

	"github.com/hbaldwin98/tui-writer/input"
)

func TestResolveActionDoesNotBufferInsertModePrefixes(t *testing.T) {
	ed := New()
	ed.SetKeymap(input.Keymap{
		input.ModeInsert: {
			"aa": input.ActionSave,
		},
	})

	action, ok := ed.ResolveAction("a")
	if action != input.ActionNone || ok {
		t.Fatalf("expected insert prefix not to buffer, got action=%q ok=%v", action, ok)
	}
}

func TestResolveActionBuffersNormalModePrefixes(t *testing.T) {
	ed := New()
	ed.SetInputMode(input.ModeNormal)
	ed.SetKeymap(input.Keymap{
		input.ModeNormal: {
			"aa": input.ActionSave,
		},
	})

	action, ok := ed.ResolveAction("a")
	if action != input.ActionNone || !ok {
		t.Fatalf("expected normal prefix to be consumed, got action=%q ok=%v", action, ok)
	}
}

func TestResolveActionReturnsExactInsertModeAction(t *testing.T) {
	ed := New()
	ed.SetKeymap(input.VimKeymap)

	action, ok := ed.ResolveAction("esc")
	if !ok || action != input.ActionNormalMode {
		t.Fatalf("expected esc to resolve normal mode, got action=%q ok=%v", action, ok)
	}
}
