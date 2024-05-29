package editor

import (
	"fmt"
	"os"
	"strings"

	"github.com/cool-pants/modus/buffer"
	"github.com/cool-pants/modus/src/ui"
)

var (
	inp = make([]byte, 1)
)

type NavHandler interface {
	moveUp(y int)
	moveDown(y int)
	moveRight(x int)
	moveLeft(x int)
	close()
}

type Editor struct {
	term *ui.TermState
	buf  *buffer.DoubleBuffer
	kill chan struct{}
	mode EditorMode

	cx, cy int
}

func (e *Editor) moveUp(y int) {
	e.cy = max(0, e.cy-1)
	e.writeStatusLine(true)
	e.buf.Write(fmt.Sprintf("\x1b[%d;%dH", e.cy+1, e.cx+1))
	e.buf.Write("\x1b[?25h")
}

func (e *Editor) moveDown(y int) {
	e.cy = min(e.term.Cols-3, e.cy+1)
	e.writeStatusLine(true)
	e.buf.Write(fmt.Sprintf("\x1b[%d;%dH", e.cy+1, e.cx+1))
	e.buf.Write("\x1b[?25h")
}
func (e *Editor) moveLeft(x int) {
	e.cx = max(3, e.cx-1)
	e.writeStatusLine(true)
	e.buf.Write(fmt.Sprintf("\x1b[%d;%dH", e.cy+1, e.cx+1))
	e.buf.Write("\x1b[?25h")
}

func (e *Editor) moveRight(x int) {
	e.cx = min(e.term.Rows-1, e.cx+1)
	e.writeStatusLine(true)
	e.buf.Write(fmt.Sprintf("\x1b[%d;%dH", e.cy+1, e.cx+1))
	e.buf.Write("\x1b[?25h")
}

func (e *Editor) close() {
	e.term.Close()
}

func InitEditor(term *ui.TermState) *Editor {
	doubleBuffer := buffer.NewDoubleBuffer()
	editor := &Editor{
		term: term,
		buf:  doubleBuffer,
		kill: make(chan struct{}),
		cx:   3,
	}
	modeManager := NewModeManager(editor)
	editor.mode = modeManager.ActiveMode
	editor.drawEditor()
	return editor
}

func (e *Editor) Start() {
	for {
		fmt.Fprint(os.Stdin, e.buf.Read())
		os.Stdin.Read(inp)
		e.mode.processKeyPress(inp)
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

func (e *Editor) drawEditor() {
	// Hide Cursor
	e.buf.Write("\x1b[?25l")

	// // Clear Screen
	// e.buf.WriteToBuf("\x1b[2J")

	// Move Cursor to top left
	e.buf.Write("\x1b[1;1H")

	e.drawCols()

	// Move Cursor to top left
	e.buf.Write(fmt.Sprintf("\x1b[%d;%dH", e.cy+1, e.cx+1))

	// UnHide Cursor
	e.buf.Write("\x1b[?25h")

	fmt.Fprint(os.Stdout, e.buf.Read())
}

func (e *Editor) drawCols() {
	for y := 0; y < e.term.Cols-1; y++ {
		if y < e.term.Cols-2 {
			e.buf.Write("~")
		}
		e.buf.Write("\x1b[K\r\n")
	}
	e.writeStatusLine(false)
}

func (e *Editor) writeStatusLine(hideCursor bool) {
	if hideCursor {
		e.buf.Write("\x1b[?25l")
	}
	e.buf.Write(fmt.Sprintf("\x1b[%d;0H", e.term.Rows))
	modeName := fmt.Sprintf("----%s----", e.mode.getModeName())
	position := fmt.Sprintf("%d: %d", e.cx+1, e.cy+1)
	e.buf.Write(e.spaceBetween(modeName, position))
	e.buf.Write("\x1b[K")
}

func (e *Editor) spaceBetween(s1 string, s2 string) string {
	buf := strings.Builder{}
	buf.WriteString(s1)
	for i := 0; i < e.term.Rows-(len(s2)+len(s1)); i++ {
		buf.WriteString(" ")
	}
	buf.WriteString(s2)
	return buf.String()
}
