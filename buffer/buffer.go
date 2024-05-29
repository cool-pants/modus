package buffer

import (
	"bytes"
)

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
