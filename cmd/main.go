package main

import (
	"io"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
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

func initialModel() model {
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
		editor.FilePath = os.Args[1]
		if contents, err := os.ReadFile(editor.FilePath); err == nil {
			editor.TextArea.SetValue(string(contents))
		}
	}

	return model{
		editor:      &editor,
		dump:        dump,
		editorWidth: defaultEditorWidth,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) View() string {
	header := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Render(m.editor.Title())

	footer := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Render("")

	contentHeight := max(0, m.height-lipgloss.Height(header)-lipgloss.Height(footer))
	editorWidth := min(m.width, m.editorWidth)
	m.editor.Resize(editorWidth, contentHeight)

	content := lipgloss.NewStyle().
		Width(m.width).
		Height(contentHeight).
		Align(lipgloss.Center, lipgloss.Top).
		Render(m.editor.View())

	return lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.dump != nil {
		spew.Fdump(m.dump, msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q":
			return m, tea.Quit
		case "ctrl+s":
			m.saveErr = m.editor.Save()
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	cmd := m.editor.Update(msg)
	return m, cmd
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
