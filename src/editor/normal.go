package editor

type NormalModeHandler interface {
	NavHandler
	Closable
}

type NormalMode struct {
	Name        string
	Identifier  rune
	Write       Writer
	termHandler NormalModeHandler
	modeHandler ModeHandler
}

func (n *NormalMode) getModeName() string {
	return n.Name
}

func (n *NormalMode) processKeyPress(key []byte) {
	switch key[0] {
	case 'j':
		n.termHandler.moveDown(1)
		break
	case 'k':
		n.termHandler.moveUp(1)
		break
	case 'l':
		n.termHandler.moveRight(1)
		break
	case 'h':
		n.termHandler.moveLeft(1)
		break
	case 'i':
		n.modeHandler.switchMode(INSERT_MODE)
		break
	case CTRLKey('q'):
		n.termHandler.close()
		break
	}
}
