package ui

import (
	"fmt"
	"io"
)

type Buffer struct {
	byteArr []byte
	length  int
}

type DoubleBuffer struct {
	length     int
	contentLen int
	primary    *Buffer
	secondary  *Buffer
}

func newBuf(cap int) *Buffer {
	return &Buffer{
		byteArr: make([]byte, cap, cap),
		length:  0,
	}
}

func (buf *Buffer) expand(cap int) {
	newByteArr := make([]byte, cap, cap)
	copy(newByteArr, buf.byteArr)
	buf.byteArr = newByteArr
}

func (buf *Buffer) writeBufn(length int, data []byte) {
	for i := 0; i < length; i++ {
		buf.byteArr[buf.length+i] = data[i]
	}
	buf.length += length
}

// Creates a new empty buffer.
// And returns a Pointer object to said buffer
func NewDoubleBuf() *DoubleBuffer {
	return &DoubleBuffer{
		length:     0,
		contentLen: 0,
		primary:    newBuf(0),
		secondary:  newBuf(0),
	}
}

// Double Cap of the primary and secondary buffers
func (b *DoubleBuffer) ExpandCap(n int) {
	// Double len of the buffer
	b.secondary.expand(b.secondary.length + n)
}

// Write into the buffer at some offset.
// The Write will always write to the secondary buffer
// to prevent screen flicker.
func (b *DoubleBuffer) WriteToBuf(data string) (int, error) {
	b.secondary.expand(b.secondary.length + len(data))
	b.secondary.writeBufn(len(data), []byte(data))
	return b.secondary.length, nil
}

// Sync Primary and secondary buffers
func (b *DoubleBuffer) Sync() error {
	b.primary.expand(b.secondary.length)
	copy(b.primary.byteArr, b.secondary.byteArr)
	b.primary.length = b.secondary.length
	b.length = b.primary.length
	return nil
}

func (b *DoubleBuffer) WriteToWriter(writer io.Writer) {
	fmt.Fprint(writer, string(b.primary.byteArr))
}
