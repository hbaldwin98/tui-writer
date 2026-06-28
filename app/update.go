package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
	"github.com/hbaldwin98/tui-writer/input"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.dump != nil {
		spew.Fdump(m.dump, msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if action, ok := m.editor.GetAction(msg.String()); ok {
			switch action {
			case input.ActionQuit:
				return m, tea.Quit
			case input.ActionSave:
				err := m.editor.Save()
				if err != nil {
					m.setStatus(StatusError, err.Error())
					return m, nil
				}

				m.setStatus(StatusInfo, "Saved...")
				return m, nil
			case input.ActionPreview:
				m.editor.TogglePreview()
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.editor.Resize(m.currentEditorWidth(), m.contentHeight())
	}

	cmd := m.editor.Update(msg)
	return m, cmd
}
