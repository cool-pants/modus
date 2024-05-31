package ui

import (
	"bytes"
	"fmt"
	"os"

	"github.com/cool-pants/modus/src/utils"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

type TermState struct {
	Cols int
	Rows int

	origState *term.State
}

var (
	inp                  = make([]byte, 1)
	inputFileDescriptor  = os.Stdin.Fd()
	outputFileDescriptor = os.Stdout.Fd()
)

func read() bool {
	_, err := os.Stdin.Read(inp)
	return err == nil
}

func (t *TermState) getCursorPos() int {
	n, err := fmt.Fprintf(os.Stdout, "\x1b[6n")
	utils.Assert(err == nil, fmt.Sprintf("Failed to write to Stdout with %v", err), t)
	utils.Assert(n == 4, "Incorrect length of 4 bytes", t)
	buf := make([]byte, 32)
	var i int

	fmt.Printf("\r\n")

	for i < len(buf)-1 {
		if !read() {
			return -1
		}
		buf[i] = inp[0]
		if buf[i] == 'R' {
			break
		}
		i++
	}
	if buf[0] != byte('\x1b') || buf[1] != '[' {
		return -1
	}
	_, err = fmt.Fscanf(bytes.NewReader(buf[2:]), "%d;%d", &t.Cols, &t.Rows)
	utils.Assert(err == nil, fmt.Sprintf("Failed to parse buffer with %v", err), t)
	fmt.Fprintf(os.Stdout, "\r\nRows %d: Cols %d\r\n", t.Rows, t.Cols)
	read()

	return 0
}

func (t *TermState) initTerm() int {
	ws, err := unix.IoctlGetWinsize(int(inputFileDescriptor), unix.TIOCGWINSZ)
	if err != nil {
		t.Cols = int(ws.Col)
		t.Rows = int(ws.Row)
	}
	n, err := fmt.Fprint(os.Stdout, "\x1b[999C\x1b[999B")
	if err != nil || n != 12 {
		return -1
	}
	return t.getCursorPos()
}

func NewTermState() *TermState {
	state, err := term.MakeRaw(int(outputFileDescriptor))
	utils.Assert(err == nil, "Couldn't make raw term for STDOUT", nil)
	term := &TermState{
		origState: state,
	}
	if term.initTerm() == -1 {
		term.Kill("Failed to init terminal and sizes")
	}

	return term
}

func (t *TermState) Close() {
	term.Restore(int(os.Stdin.Fd()), t.origState)
	// Clear Screen
	fmt.Fprint(os.Stdin, "\x1b[2J")

	// Move Cursor to top left
	fmt.Fprint(os.Stdin, "\x1b[H")
	os.Exit(0)
}

func (t *TermState) Kill(msg string) {
	term.Restore(int(outputFileDescriptor), t.origState)
	fmt.Fprint(os.Stdin, "\x1b[2J")
	fmt.Fprint(os.Stdin, "\x1b[H")
	fmt.Printf("Load Failed with: %s", msg)
	os.Exit(1)
}
