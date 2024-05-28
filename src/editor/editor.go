package editor

import (
	"os"

	"github.com/cool-pants/modus/src/ui"
)

type Editor struct {
	term *ui.TermState
	buf  *ui.DoubleBuffer
}

func InitEditor(term *ui.TermState) *Editor {
	editor := &Editor{
		term: term,
		buf:  ui.NewDoubleBuf(),
	}
	editor.drawEditor()
	return editor
}

func (e *Editor) Start() {
	for {
		e.buf.Sync()
		e.buf.WriteToWriter(os.Stdout)
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
	default:
		if !isCtrl(inp) {
			e.buf.WriteToBuf(string(inp))
		}
		break
	}
}

func (e *Editor) drawEditor() {
	// Hide Cursor
	e.buf.WriteToBuf("\x1b[?25l")

	// // Clear Screen
	// e.buf.WriteToBuf("\x1b[2J")

	// Move Cursor to top left
	e.buf.WriteToBuf("\x1b[1;1H")

	e.drawCols()

	// Move Cursor to top left
	e.buf.WriteToBuf("\x1b[1;3H")

	// UnHide Cursor
	e.buf.WriteToBuf("\x1b[?25h")

	e.buf.Sync()
	e.buf.WriteToWriter(os.Stdout)
}

func (e *Editor) drawCols() {
	for y := 0; y < e.term.Cols; y++ {
		e.buf.WriteToBuf("~")
		e.buf.WriteToBuf("\x1b[K")
		if y < e.term.Cols-1 {
			e.buf.WriteToBuf("\r\n")
		}
	}
}
