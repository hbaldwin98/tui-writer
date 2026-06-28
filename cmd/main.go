package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hbaldwin98/tui-writer/app"
)

func main() {
	p := tea.NewProgram(app.Init(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
