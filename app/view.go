package app

import "github.com/charmbracelet/lipgloss"

func (m model) View() string {
	header := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Render(m.editor.Title())

	footer := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Render("")

	content := lipgloss.NewStyle().
		Width(m.width).
		Height(m.contentHeight()).
		Align(lipgloss.Center, lipgloss.Top).
		Render(m.editor.View())

	return lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
}
