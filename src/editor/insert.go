package editor

type TermHandler interface {
	Closable
	Sizable
}

type InsertMode struct {
	Name        string
	Identifier  rune
	Writer      Writer
	termHandler TermHandler
}

func (n *InsertMode) getModeName() string {
	return n.Name
}

func (n *InsertMode) processKeyPress(key []byte) {
	if !isCtrl(key) {
		n.Writer.Write(n.termHandler.getCurY(), n.termHandler.getCurX(), string(key))
	}
	switch key[0] {
	case CTRLKey('q'):
		n.termHandler.close()
		break
	}
}
