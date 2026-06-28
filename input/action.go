package input

type Action string

const (
	ActionNone    Action = "none"
	ActionSave    Action = "save"
	ActionLoad    Action = "load"
	ActionQuit    Action = "quit"
	ActionPreview Action = "preview"
	ActionHelp    Action = "help"
	ActionCancel  Action = "cancel"
	ActionConfirm Action = "confirm"

	// Modes
	ActionInsertMode Action = "insert"

	// Insert Mode Actions
	ActionInsertModeAbove = "insert_above"
	ActionInsertModeBelow = "insert_below"
	ActionInsertModeNext  = "insert_next"

	ActionNormalMode Action = "normal"
	ActionVisualMode Action = "visual"

	// Movement
	ActionMoveUp       Action = "move_up"
	ActionMoveDown     Action = "move_down"
	ActionMoveLeft     Action = "move_left"
	ActionMoveRight    Action = "move_right"
	ActionNextWord     Action = "next_word"
	ActionPreviousWord Action = "previous_word"
	ActionEndOfWord    Action = "end_of_word"
	ActionStartOfLine  Action = "start_of_line"
	ActionEndOfLine    Action = "end_of_line"
	ActionStartOfFile  Action = "start_of_file"
	ActionEndOfFile    Action = "end_of_file"
)
