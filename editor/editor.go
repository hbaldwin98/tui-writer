package editor

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/hbaldwin98/tui-writer/input"
)

type Editor struct {
	mode        input.InputMode
	textArea    textarea.Model
	previewView viewport.Model
	keymap      input.Keymap
	filePath    string
	modified    bool
	preview     bool
	pendingKey  string
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
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	e.filePath = filePath
	e.textArea.SetValue(string(contents))
	// When we set the contents, the cursor is pushed all the way to the bottom.
	// TextArea supports going back to the beginning using Ctrl+Home, so we can
	// make use of this to push the view to the top before we render anything
	e.textArea, _ = e.textArea.Update(tea.KeyMsg{Type: tea.KeyCtrlHome})
	e.modified = false

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
	e.pendingKey = ""
	e.textArea.Cursor.Blink = mode == input.ModeInsert
}

func (e *Editor) SetKeymap(km input.Keymap) {
	e.keymap = km
}

func (e Editor) GetAction(key string) (input.Action, bool) {
	modeKeymap, ok := e.keymap[e.mode]
	if !ok {
		return input.ActionNone, false
	}

	if action, ok := modeKeymap[key]; ok {
		return action, ok
	}

	return input.ActionNone, false
}

func (e Editor) hasBindingPrefix(prefix string) bool {
	modeKeymap, ok := e.keymap[e.mode]
	if !ok {
		return false
	}

	for key := range modeKeymap {
		if key != prefix && strings.HasPrefix(key, prefix) {
			return true
		}
	}

	return false
}

func (e *Editor) ResolveAction(key string) (input.Action, bool, bool) {
	if e.pendingKey != "" {
		key = e.pendingKey + key
		e.pendingKey = ""
	}

	action, ok := e.GetAction(key)
	if !ok {
		if e.hasBindingPrefix(key) {
			e.pendingKey = key
			return input.ActionNone, false, true
		}

		return input.ActionNone, false, false
	}

	if e.hasBindingPrefix(key) {
		e.pendingKey = key
		return input.ActionNone, false, true
	}

	return action, true, false
}

func (e *Editor) HandleAction(action input.Action) tea.Cmd {
	var cmd tea.Cmd
	switch action {
	case input.ActionInsertMode, input.ActionInsertModeNext, input.ActionInsertModeAbove, input.ActionInsertModeBelow:
		e.enterInsertMode(action)
	case input.ActionNormalMode:
		e.SetInputMode(input.ModeNormal)
	case input.ActionMoveUp:
		e.textArea.CursorUp()
	case input.ActionMoveDown:
		e.textArea.CursorDown()
	case input.ActionMoveLeft:
		e.textArea, cmd = e.textArea.Update(tea.KeyMsg{Type: tea.KeyLeft})
	case input.ActionMoveRight:
		e.textArea, cmd = e.textArea.Update(tea.KeyMsg{Type: tea.KeyRight})
	case input.ActionNextWord:
		e.textArea, cmd = e.textArea.Update(tea.KeyMsg{Type: tea.KeyRight, Alt: true})
	case input.ActionPreviousWord:
		e.textArea, cmd = e.textArea.Update(tea.KeyMsg{Type: tea.KeyLeft, Alt: true})
	case input.ActionEndOfWord:
		e.textArea, _ = e.textArea.Update(tea.KeyMsg{Type: tea.KeyRight})
		e.textArea, cmd = e.textArea.Update(tea.KeyMsg{Type: tea.KeyRight, Alt: true})
		e.textArea, _ = e.textArea.Update(tea.KeyMsg{Type: tea.KeyLeft})
	case input.ActionStartOfLine:
		e.textArea, cmd = e.textArea.Update(tea.KeyMsg{Type: tea.KeyHome})
	case input.ActionEndOfLine:
		e.textArea, cmd = e.textArea.Update(tea.KeyMsg{Type: tea.KeyEnd})
	case input.ActionStartOfFile:
		e.textArea, cmd = e.textArea.Update(tea.KeyMsg{Type: tea.KeyCtrlHome})
	case input.ActionEndOfFile:
		e.textArea, cmd = e.textArea.Update(tea.KeyMsg{Type: tea.KeyCtrlEnd})
	}

	return cmd
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

	switch e.mode {
	case input.ModeInsert:
		return e.handleModeInsert(msg)
	case input.ModeNormal:
		return e.handleModeNormal(msg)
	default:
		return nil
	}
}

func (e *Editor) enterInsertMode(action input.Action) {
	switch action {
	case input.ActionInsertModeNext:
		e.textArea, _ = e.textArea.Update(tea.KeyMsg{Type: tea.KeyRight})
	case input.ActionInsertModeAbove:
		e.textArea, _ = e.textArea.Update(tea.KeyMsg{Type: tea.KeyHome})
		e.textArea, _ = e.textArea.Update(tea.KeyMsg{Type: tea.KeyEnter})
		e.textArea.CursorUp()
		e.modified = true
	case input.ActionInsertModeBelow:
		e.textArea, _ = e.textArea.Update(tea.KeyMsg{Type: tea.KeyEnd})
		e.textArea, _ = e.textArea.Update(tea.KeyMsg{Type: tea.KeyEnter})
		e.modified = true
	}

	e.SetInputMode(input.ModeInsert)
}

func (e *Editor) handleModeNormal(msg tea.Msg) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}

	action, ok, waiting := e.ResolveAction(keyMsg.String())
	if waiting {
		return nil
	}
	if !ok {
		return nil
	}

	return e.HandleAction(action)
}

func (e *Editor) handleModeInsert(msg tea.Msg) tea.Cmd {
	previous := e.textArea.Value()

	var cmd tea.Cmd
	e.textArea, cmd = e.textArea.Update(msg)

	if e.textArea.Value() != previous {
		e.modified = true
	}

	return cmd
}
