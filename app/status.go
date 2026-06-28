package app

type statusKind int

const (
	StatusNone statusKind = iota
	StatusInfo
	StatusError
)

type statusMessage struct {
	text string
	kind statusKind
}
