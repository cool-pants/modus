package editor

import (
	"fmt"
	"os"
)

type EditorMode interface {
	processKeyPress()
	getMode() rune
}

const (
	NORMAL_MODE = 'n'
	INSERT_MODE = 'i'
)

var (
	inp = make([]byte, 1)
)

type NormalMode struct {
	name       string
	identifier rune
	editor     *Editor
}

func (m *NormalMode) getMode() rune {
	return m.identifier
}

func (m *NormalMode) processKeyPress() {
	switch inp[0] {
	case 'j':
		m.editor.cy = min(m.editor.cy+1, m.editor.term.Rows)
		m.editor.buf.WriteToBuf(fmt.Sprintf("\x1b[%d;%dH", m.editor.cy+1, m.editor.cx+1))
		break
	case 'k':
		m.editor.cy = max(m.editor.cy-1, 0)
		m.editor.buf.WriteToBuf(fmt.Sprintf("\x1b[%d;%dH", m.editor.cy+1, m.editor.cx+1))
		break
	default:
		break
	}
}

type InsertMode struct {
	name       string
	identifier rune
	editor     *Editor
}

func (m *InsertMode) getMode() rune {
	return m.identifier
}

func (m *InsertMode) processKeyPress() {
	switch inp[0] {
	case CTRLKey('q'):
		break
	case 127:
		m.editor.buf.WriteToBuf("\x1b[3~")
		break
	default:
		if !isCtrl(inp) {
			m.editor.buf.WriteToBuf(string(inp))
		}
		break
	}
}

type ModeManager struct {
	editor       *Editor
	CurrentMode  EditorMode
	PreviousMode EditorMode
}

func NewModeManager(editor *Editor) *ModeManager {
	return &ModeManager{
		CurrentMode: &NormalMode{
			name:       "NORMAL",
			identifier: NORMAL_MODE,
			editor:     editor,
		},
		PreviousMode: nil,
	}
}

func (m *ModeManager) processKeyPress() {
	os.Stdin.Read(inp)

	if m.CurrentMode.getMode() == NORMAL_MODE {
		if inp[0] == CTRLKey('q') {
			m.editor.term.Close()
			return
		} else if inp[0] == 'i' {
			m.PreviousMode = m.CurrentMode
			m.CurrentMode = &InsertMode{
				name:       "INSERT",
				identifier: INSERT_MODE,
				editor:     m.editor,
			}
			return
		}
	}
	if m.CurrentMode.getMode() == INSERT_MODE && inp[0] == '\x1b' {
		temp := m.PreviousMode
		m.PreviousMode = m.CurrentMode
		m.CurrentMode = temp
		return
	}
	m.CurrentMode.processKeyPress()

}
