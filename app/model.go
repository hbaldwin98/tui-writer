package app

import (
	"io"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hbaldwin98/tui-writer/editor"
	"github.com/hbaldwin98/tui-writer/input"
)

const (
	defaultEditorWidth = 88
	maxStatusHistory   = 100
)

type clearStatusMsg struct{}

// Naive way to clear status. Should have it be a little more aware
// in the model so that we don't override a status we actually want
// if one comes in AFTER we call this
func clearStatusAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}

type model struct {
	editor *editor.Editor

	dump        io.Writer
	width       int
	height      int
	editorWidth int
	status      statusMessage
	history     []statusMessage
}

func Init() tea.Model {
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
	model := model{
		editor:      &editor,
		dump:        dump,
		editorWidth: defaultEditorWidth,
	}

	if len(os.Args) > 1 {
		if err := editor.Load(os.Args[1]); err != nil {
			model.setStatus(StatusError, err.Error())
		}
	}

	return model
}

func (m model) contentHeight() int {
	headerHeight := 1
	footerHeight := 1

	return max(0, m.height-headerHeight-footerHeight)
}

func (m model) currentEditorWidth() int {
	return min(m.width, m.editorWidth)
}

func (m *model) setStatus(kind statusKind, text string) {
	m.status = statusMessage{
		text,
		kind,
	}
	m.history = append(m.history, m.status)

	if len(m.history) > maxStatusHistory {
		m.history = m.history[len(m.history)-maxStatusHistory:]
	}
}

func (m *model) clearStatus() {
	m.status = statusMessage{}
}

func (m model) Init() tea.Cmd {
	if m.status.text != "" {
		return tea.Batch(textarea.Blink, clearStatusAfter(3*time.Second))
	}

	return textarea.Blink
}
