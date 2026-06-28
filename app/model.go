package app

import (
	"io"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hbaldwin98/tui-writer/editor"
	"github.com/hbaldwin98/tui-writer/input"
)

const defaultEditorWidth = 88

type model struct {
	editor *editor.Editor

	dump        io.Writer
	width       int
	height      int
	editorWidth int
	saveErr     error
}

func Init() model {
	var dump *os.File
	if _, ok := os.LookupEnv("DEBUG"); ok {
		var err error
		dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}
	}

	editor := editor.New()
	editor.SetInputMode(input.ModeInsert)
	if len(os.Args) > 1 {
		if err := editor.Load(os.Args[1]); err != nil {
			// TODO: display error
		}
	}

	return model{
		editor:      &editor,
		dump:        dump,
		editorWidth: defaultEditorWidth,
	}
}

func (m model) contentHeight() int {
	headerHeight := 1
	footerHeight := 1

	return max(0, m.height-headerHeight-footerHeight)
}

func (m model) currentEditorWidth() int {
	return min(m.width, m.editorWidth)
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}
