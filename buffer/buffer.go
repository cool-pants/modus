package buffer

import (
	"bytes"
)

type AppendBuffer struct {
	buf        []*bytes.Buffer
	rows, cols uint
}

// AppendBuffer will always append to the buffer, meaning
// In case a character is present in x,y the buf moves the
// char to the right and appends the new char in the position
func (abuf *AppendBuffer) append(row, col int, byteData []byte) {
	bufBytes := abuf.buf[row].Bytes()
	newArr := make([]byte, len(bufBytes)+len(byteData))
	newArr = append(newArr, bufBytes[:col]...)
	newArr = append(newArr, byteData...)
	newArr = append(newArr, bufBytes[col:]...)
	abuf.buf[row].Truncate(0)
	abuf.buf[row].ReadFrom(bytes.NewReader(newArr))
}

type DoubleBuffer struct {
	readBuf   *bytes.Buffer
	writeBuf  *bytes.Buffer
	WriteChan chan string
	Close     chan struct{}
}

func NewDoubleBuffer() *DoubleBuffer {
	buf := &DoubleBuffer{
		WriteChan: make(chan string),
		readBuf:   bytes.NewBuffer(make([]byte, 0)),
		writeBuf:  bytes.NewBuffer(make([]byte, 0)),
		Close:     make(chan struct{}),
	}
	return buf
}

func (buf *DoubleBuffer) Write(data string) {
	buf.writeBuf.WriteString(data)
}

func (buf *DoubleBuffer) Read() string {
	// Sync Bufs
	buf.readBuf.Truncate(0)
	buf.readBuf.ReadFrom(bytes.NewReader(buf.writeBuf.Bytes()))

	return buf.readBuf.String()
}
