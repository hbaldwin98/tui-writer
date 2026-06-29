package app

import (
	"time"

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
		if action, ok := m.editor.ResolveAction(msg.String()); ok {
			switch action {
			case input.ActionQuit:
				return m, tea.Quit
			case input.ActionSave:
				err := m.editor.Save()
				if err != nil {
					m.setStatus(StatusError, err.Error())
					return m, clearStatusAfter(time.Second)
				}

				m.setStatus(StatusInfo, "Saved...")
				return m, clearStatusAfter(time.Second)
			case input.ActionPreview:
				m.editor.TogglePreview()
				return m, nil
			default:
				return m, m.editor.HandleAction(action)
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.editor.Resize(m.currentEditorWidth(), m.contentHeight())
	case clearStatusMsg:
		m.clearStatus()
	}

	cmd := m.editor.Update(msg)
	return m, cmd
}
