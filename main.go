package main

import (
	"github.com/cool-pants/modus/src/editor"
	"github.com/cool-pants/modus/src/ui"
)

func main() {
	termState := ui.NewTermState()
	editor := editor.InitEditor(termState)
	editor.Start()
}
