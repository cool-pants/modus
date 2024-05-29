package editor

type EditorMode interface {
	processKeyPress(key []byte)
	getModeName() string
}

const (
	NORMAL_MODE = 'n'
	INSERT_MODE = 'i'
)

type ModeManager struct {
	ActiveMode EditorMode
}

func NewModeManager(editor *Editor) *ModeManager {
	modeManager := &ModeManager{
		ActiveMode: &NormalMode{
			Name:        "NORMAL",
			Identifier:  NORMAL_MODE,
			bufChan:     editor.buf.WriteChan,
			termHandler: editor,
		},
	}

	return modeManager
}
