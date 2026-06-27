package main

import (
	"fmt"

	"github.com/hbaldwin98/tui-writer/editor"
	"github.com/hbaldwin98/tui-writer/input"
)

func main() {
	editor := editor.New()
	editor.SetInputMode(input.ModeInsert)
	fmt.Println("Starting Editor...")
}
