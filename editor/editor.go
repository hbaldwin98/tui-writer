package editor

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hbaldwin98/tui-writer/input"
)

type Editor struct {
	mode     input.InputMode
	TextArea textarea.Model
	Keymap   input.Keymap
	FilePath string
	Modified bool
}

func New() Editor {
	ta := textarea.New()
	ta.Focus()
	ta.Placeholder = "Start typing..."
	ta.Prompt = ""
	ta.ShowLineNumbers = false
	ta.SetWidth(0)

	return Editor{
		TextArea: ta,
		mode:     input.ModeInsert,
		Keymap:   input.DefaultKeymap,
	}
}

// TODO: implement this
func (e *Editor) SetInputMode(mode input.InputMode) {
	e.mode = mode
}

// TODO: Show the error to the end user. Store it somewhere/pop a message.
// TODO: Allow the user to name the newly created file BEFORE they save it. Otherwise save to the opened file.
func (e *Editor) Save() error {
	filePath := e.FilePath
	if e.FilePath == "" {
		path, err := os.Getwd()
		if err != nil {
			return err
		}

		filePath = filepath.Join(path, "undefined.md")
	}

	if err := os.WriteFile(filePath, []byte(e.TextArea.Value()), 0o644); err != nil {
		return err
	}

	e.FilePath = filePath
	e.Modified = false

	return nil
}

func (e Editor) Title() string {
	if e.FilePath == "" {
		return "Untitled"
	}

	title := filepath.Base(e.FilePath)
	if e.Modified {
		title += " *"
	}

	return title
}

func (e Editor) View() string {
	return e.TextArea.View()
}

func (e *Editor) Resize(width, height int) {
	e.TextArea.SetWidth(width)
	e.TextArea.SetHeight(height)
}

func (e *Editor) Update(msg tea.Msg) tea.Cmd {
	previous := e.TextArea.Value()

	var cmd tea.Cmd
	e.TextArea, cmd = e.TextArea.Update(msg)

	if e.TextArea.Value() != previous {
		e.Modified = true
	}

	return cmd
}
