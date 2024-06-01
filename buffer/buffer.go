package buffer

import (
	"bytes"
	"strings"
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

func (abuf *AppendBuffer) replace(row, col int, byteData []byte) {
	bufBytes := abuf.buf[row].Bytes()
	newArr := make([]byte, len(bufBytes)+len(byteData))
	newArr = append(newArr, bufBytes[:col-1]...)
	newArr = append(newArr, byteData...)
	newArr = append(newArr, bufBytes[col+1:]...)
	abuf.buf[row].Truncate(0)
	abuf.buf[row].ReadFrom(bytes.NewReader(newArr))
}

func (abuf *AppendBuffer) getString() string {
	arr := make([]string, abuf.rows)
	for _, buf := range abuf.buf {
		arr = append(arr, buf.String())
	}
	return strings.Join(arr, "\r\n")
}

type DoubleBuffer struct {
	readBuf   *AppendBuffer
	writeBuf  *AppendBuffer
	WriteChan chan string
	Close     chan struct{}
}

func createBufArr(size int) []*bytes.Buffer {
	buf := make([]*bytes.Buffer, size)

	for i := 0; i < size; i++ {
		buf[i] = bytes.NewBuffer(make([]byte, 0))
	}

	return buf
}

func NewDoubleBuffer(rows, cols uint) *DoubleBuffer {
	buf := &DoubleBuffer{
		WriteChan: make(chan string),
		readBuf: &AppendBuffer{
			buf:  createBufArr(int(rows)),
			rows: rows,
			cols: cols,
		},
		writeBuf: &AppendBuffer{
			buf:  createBufArr(int(rows)),
			rows: rows,
			cols: cols,
		},
		Close: make(chan struct{}),
	}
	return buf
}

func (buf *DoubleBuffer) Write(x, y int, data string) {
	buf.writeBuf.append(x, y, []byte(data))
}

func (buf *DoubleBuffer) Read() string {
	// Sync Bufs
	for i := range buf.writeBuf.buf {
		buf.readBuf.buf[i].Truncate(0)
		buf.readBuf.buf[i].ReadFrom(bytes.NewReader(buf.writeBuf.buf[i].Bytes()))
	}

	return buf.readBuf.getString()
}
