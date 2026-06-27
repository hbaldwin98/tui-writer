package editor

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/hbaldwin98/tui-writer/input"
)

type Editor struct {
	mode     input.InputMode
	TextArea textarea.Model
	FilePath string
	Modified bool
}

func New() Editor {
	ta := textarea.New()
	ta.Focus()
	ta.Placeholder = "Start typing..."
	ta.ShowLineNumbers = true

	return Editor{
		TextArea: ta,
		mode:     input.ModeInsert,
	}
}

func (e *Editor) SetInputMode(mode input.InputMode) {
	e.mode = mode
}
