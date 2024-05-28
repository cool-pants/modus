package editor

import (
	"fmt"
	"os"

	"github.com/cool-pants/modus/src/ui"
)

type Editor struct {
	term *ui.TermState
}

func InitEditor(term *ui.TermState) *Editor {
	editor := &Editor{
		term: term,
	}
	editor.drawEditor()
	return editor
}

func (e *Editor) Start() {
	for {
		e.drawEditor()
		e.processKeyPress()
	}
}

func CTRLKey(c rune) byte {
	return byte(c) & 0x1f
}

func isCtrl(c []byte) bool {
	// ASCII characters 32-126 are printable chars https://www.asciitable.com/
	if c[0] >= 32 && c[0] <= 126 {
		return false
	}
	return true
}

func (e *Editor) processKeyPress() {
	inp := make([]byte, 1)

	os.Stdin.Read(inp)

	switch inp[0] {
	case CTRLKey('q'):
		e.term.Close()
		break
	}
}

func (e *Editor) drawEditor() {
	// Clear Screen
	fmt.Fprint(os.Stdin, "\x1b[2J")

	// Move Cursor to top left
	fmt.Fprint(os.Stdin, "\x1b[H")

	e.drawCols()

	// Move Cursor to top left
	fmt.Fprint(os.Stdin, "\x1b[H")
}

func (e *Editor) drawCols() {
	for y := 0; y < e.term.Cols; y++ {
		fmt.Fprint(os.Stdout, "~")
		if y < e.term.Cols-1 {
			fmt.Fprint(os.Stdout, "\r\n")
		}
	}
}
