package input

type Action string

const (
	ActionSave     Action = "save"
	ActionQuit     Action = "quit"
	ActionHelp     Action = "help"
	ActionCancel   Action = "cancel"
	ActionConfirm  Action = "confirm"
	ActionMoveUp   Action = "move_up"
	ActionMoveDOwn Action = "move_down"
)
