package input

type Keymap map[string]Action

var DefaultKeymap = Keymap{
	"ctrl+s": ActionSave,
	"ctrl+q": ActionQuit,
	"ctrl+h": ActionHelp,
	"esc":    ActionCancel,
	"enter":  ActionConfirm,
}
