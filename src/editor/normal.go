package editor

type NormalMode struct {
	Name        string
	Identifier  rune
	bufChan     chan string
	termHandler NavHandler
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
	case CTRLKey('q'):
		n.termHandler.close()
		break
	}
}
