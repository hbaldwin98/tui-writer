package editor

import "github.com/hbaldwin98/tui-writer/input"

type Editor struct {
	mode input.InputMode
}

func New() Editor {
	return Editor{
		mode: input.ModeInsert,
	}
}

func (e *Editor) SetInputMode(mode input.InputMode) {
	e.mode = mode
}
