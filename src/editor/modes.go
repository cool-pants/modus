package editor

type NavHandler interface {
	moveUp(y int)
	moveDown(y int)
	moveRight(x int)
	moveLeft(x int)
}

type Closable interface {
	close()
}

type ModeHandler interface {
	switchMode(rune)
}

type Writer interface {
	Write(string)
}

type EditorMode interface {
	processKeyPress(key []byte)
	getModeName() string
}

const (
	NORMAL_MODE = 'n'
	INSERT_MODE = 'i'
)

type ModeManager struct {
	editor         *Editor
	ActiveMode     EditorMode
	SupportedModes map[rune]EditorMode
}

func (m *ModeManager) close() {
	m.ActiveMode = m.SupportedModes[NORMAL_MODE]
	m.editor.mode = m.ActiveMode
	m.editor.renderEssentials()
}
func (m *ModeManager) switchMode(c rune) {
	m.ActiveMode = m.SupportedModes[c]
	m.editor.mode = m.ActiveMode
	m.editor.renderEssentials()
}

func NewModeManager(editor *Editor) *ModeManager {
	modeManager := &ModeManager{
		editor:         editor,
		SupportedModes: make(map[rune]EditorMode),
	}
	normalMode := &NormalMode{
		Name:        "NORMAL",
		Identifier:  NORMAL_MODE,
		Write:       editor.buf,
		termHandler: editor,
		modeHandler: modeManager,
	}
	insertMode := &InsertMode{
		Name:        "INSERT",
		Identifier:  INSERT_MODE,
		Writer:      editor.buf,
		termHandler: modeManager,
	}
	modeManager.ActiveMode = normalMode
	modeManager.SupportedModes[INSERT_MODE] = insertMode
	modeManager.SupportedModes[NORMAL_MODE] = normalMode

	return modeManager
}
