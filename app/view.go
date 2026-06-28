package app

import "github.com/charmbracelet/lipgloss"

func (m model) View() string {
	header := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Render(m.editor.Title())

	statusStyle := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center)

	switch m.status.kind {
	case StatusError:
		statusStyle = statusStyle.Foreground(lipgloss.Color("9"))
	case StatusInfo:
		statusStyle = statusStyle.Foreground(lipgloss.Color("10"))
	default:
		break
	}

	footer := statusStyle.Render(m.status.text)

	content := lipgloss.NewStyle().
		Width(m.width).
		Height(m.contentHeight()).
		Align(lipgloss.Center, lipgloss.Top).
		Render(m.editor.View())

	return lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
}
