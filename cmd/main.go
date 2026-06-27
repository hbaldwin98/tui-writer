package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hbaldwin98/tui-writer/editor"
	"github.com/hbaldwin98/tui-writer/input"
)

type model struct {
	editor *editor.Editor

	width  int
	height int
}

func initialModel() model {
	editor := editor.New()
	editor.SetInputMode(input.ModeInsert)
	return model{
		editor: &editor,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	return "Hello from the terminal\n\nctrl+q to quit"
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
