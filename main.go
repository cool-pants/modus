package main

import (
	"bytes"
	"fmt"
	"os"

	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

type ModusTermState struct {
	originalState *term.State
	height        int
	width         int
}

var (
	inp       = make([]byte, 1)
	termState = &ModusTermState{}
)

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

func editorReadKey() bool {
	_, err := os.Stdin.Read(inp)
	if err != nil {
		return false
	}
	return true
}

func getCursorPos() int {
	n, err := fmt.Fprintf(os.Stdout, "\x1b[6n")
	assert(err == nil, fmt.Sprintf("Failed to write to Stdout with %v", err))
	assert(n == 4, "Incorrect length of 4 bytes")
	buf := make([]byte, 32)
	var i int

	fmt.Printf("\r\n")

	for i < len(buf)-1 {
		if !editorReadKey() {
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
	_, err = fmt.Fscanf(bytes.NewReader(buf[2:]), "%d;%d", &termState.height, &termState.width)
	assert(err == nil, fmt.Sprintf("Failed to parse buffer with %v", err))
	fmt.Fprintf(os.Stdout, "\r\n%d: %d\r\n", termState.height, termState.width)
	editorReadKey()

	return 0
}

func getWindowSize() int {
	ws, err := unix.IoctlGetWinsize(int(os.Stdin.Fd()), 1)
	if err == nil {
		termState.height = int(ws.Row)
		termState.width = int(ws.Col)
	}

	n, err := fmt.Fprint(os.Stdout, "\x1b[999C\x1b[999B")
	if err != nil || n != 12 {
		return -1
	}
	return getCursorPos()
}

func editorProcessKeyPress() {
	editorReadKey()

	switch inp[0] {
	case CTRLKey('q'):
		term.Restore(int(os.Stdin.Fd()), termState.originalState)
		// Clear Screen
		fmt.Fprint(os.Stdin, "\x1b[2J")

		// Move Cursor to top left
		fmt.Fprint(os.Stdin, "\x1b[H")
		os.Exit(0)
		break
	}
}

/*** output ***/
func editorDrawRows() {
	for y := 0; y < termState.height; y++ {
		fmt.Fprint(os.Stdout, "~")
		if y < termState.height-1 {
			fmt.Fprintf(os.Stdout, "\r\n")
		}
	}
}

func editorRefreshScreen() {
	// Clear Screen
	fmt.Fprint(os.Stdin, "\x1b[2J")

	// Move Cursor to top left
	fmt.Fprint(os.Stdin, "\x1b[H")

	editorDrawRows()

	// Move Cursor to top left
	fmt.Fprint(os.Stdin, "\x1b[H")
}

func initEditor() {
	if getWindowSize() == -1 {
		die("Failed to load Window Size with ioctl")
	}
}

// If Assertion is false, kills term with msg
func assert(assertion bool, msg string) {
	if !assertion {
		die(msg)
	}
}

func die(msg string) {
	term.Restore(int(os.Stdin.Fd()), termState.originalState)
	fmt.Fprint(os.Stdin, "\x1b[2J")
	fmt.Fprint(os.Stdin, "\x1b[H")
	fmt.Printf("Load Failed with: %s", msg)
	os.Exit(1)
}

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	termState.originalState = oldState
	initEditor()

	for {
		editorRefreshScreen()
		editorProcessKeyPress()
	}
}
