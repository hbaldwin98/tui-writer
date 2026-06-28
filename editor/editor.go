package editor

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/hbaldwin98/tui-writer/input"
)

type E interface {
	Load(path string) error
	Save() error
	Modified() bool
	FilePath() string
	SetKeymap(keymap input.Keymap)
	GetAction(key string) (input.Action, bool)
}

type Editor struct {
	mode        input.InputMode
	textArea    textarea.Model
	previewView viewport.Model
	keymap      input.Keymap
	filePath    string
	modified    bool
	preview     bool
}

func New() Editor {
	ta := textarea.New()
	ta.Focus()
	ta.Placeholder = "Start typing..."
	ta.Prompt = ""
	ta.ShowLineNumbers = false
	ta.SetWidth(0)

	return Editor{
		textArea:    ta,
		previewView: viewport.New(0, 0),
		mode:        input.ModeInsert,
		keymap:      input.DefaultKeymap,
	}
}

// TODO: Show the error to the end user. Store it somewhere/pop a message.
// TODO: Allow the user to name the newly created file BEFORE they save it. Otherwise save to the opened file.
func (e *Editor) Save() error {
	filePath := e.filePath
	if e.filePath == "" {
		return errors.New("cannot save untitled file")
	}

	if err := os.WriteFile(filePath, []byte(e.textArea.Value()), 0o644); err != nil {
		return err
	}

	e.filePath = filePath
	e.modified = false

	return nil
}

func (e *Editor) Load(filePath string) error {
	contents, err := os.ReadFile(e.filePath)
	if err != nil {
		return err
	}

	e.filePath = filePath
	e.textArea.SetValue(string(contents))

	return nil
}

func (e Editor) Title() string {
	title := "Untitled"
	if e.filePath != "" {
		title = filepath.Base(e.filePath)
	}

	if e.modified {
		title += " *"
	}
	if e.preview {
		title += " [preview]"
	}

	return title
}

// Modified returns whether the editor has been modified
func (e Editor) Modified() bool {
	return e.modified
}

// FilePath returns the filepath for the working file in the editor
func (e Editor) FilePath() string {
	return e.filePath
}

// TODO: implement this feature
func (e *Editor) SetInputMode(mode input.InputMode) {
	e.mode = mode
}

func (e *Editor) SetKeymap(km input.Keymap) {
	e.keymap = km
}

func (e Editor) GetAction(key string) (input.Action, bool) {
	if action, ok := e.keymap[key]; ok {
		return action, ok
	}

	return input.ActionNone, false
}

func (e Editor) View() string {
	if e.preview {
		return e.previewView.View()
	}

	return e.textArea.View()
}

func (e *Editor) RenderMarkdown(width int) {
	if width <= 0 {
		width = 80
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		e.previewView.SetContent("Unable to render markdown: " + err.Error())
		return
	}

	rendered, err := renderer.Render(e.textArea.Value())
	if err != nil {
		e.previewView.SetContent("Unable to render markdown: " + err.Error())
		return
	}

	e.previewView.SetContent(rendered)
}

func (e *Editor) TogglePreview() {
	e.preview = !e.preview
	if e.preview {
		e.previewView.GotoTop()
		e.RenderMarkdown(e.previewView.Width)
	}
}

func (e *Editor) Resize(width, height int) {
	e.textArea.SetWidth(width)
	e.textArea.SetHeight(height)
	e.previewView.Width = width
	e.previewView.Height = height
}

func (e *Editor) Update(msg tea.Msg) tea.Cmd {
	if e.preview {
		var cmd tea.Cmd
		e.previewView, cmd = e.previewView.Update(msg)
		return cmd
	}

	previous := e.textArea.Value()

	var cmd tea.Cmd
	e.textArea, cmd = e.textArea.Update(msg)

	if e.textArea.Value() != previous {
		e.modified = true
	}

	return cmd
}
