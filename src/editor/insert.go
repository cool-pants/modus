package editor

type InsertMode struct {
	Name        string
	Identifier  rune
	Writer      Writer
	termHandler Closable
}

func (n *InsertMode) getModeName() string {
	return n.Name
}

func (n *InsertMode) processKeyPress(key []byte) {
	if !isCtrl(key) {
		n.Writer.Write(string(key))
	}
	switch key[0] {
	case CTRLKey('q'):
		n.termHandler.close()
		break
	}
}
