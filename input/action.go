package input

type Action string

const (
	ActionNone     Action = "none"
	ActionSave     Action = "save"
	ActionLoad     Action = "load"
	ActionQuit     Action = "quit"
	ActionPreview  Action = "preview"
	ActionHelp     Action = "help"
	ActionCancel   Action = "cancel"
	ActionConfirm  Action = "confirm"
	ActionMoveUp   Action = "move_up"
	ActionMoveDown Action = "move_down"
)
