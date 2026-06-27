package input

type InputMode int

const (
	ModeInsert InputMode = iota
	ModeNormal
	ModeCommand
	ModePrompt
	ModeConfirm
)
