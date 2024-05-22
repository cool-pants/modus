package ui

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
}

func (buf *Buffer) overwriteBuffer(length int, data []byte) {
	buf.length = length
	for i := 0; i < buf.length; i++ {
		buf.byteArr[i] = data[i]
	}
}

// Creates a new empty buffer.
// And returns a Pointer object to said buffer
func NewDoubleBuf() *DoubleBuffer {
	return &DoubleBuffer{
		length:     1024,
		contentLen: 0,
		primary:    newBuf(1024),
		secondary:  newBuf(1024),
	}
}

// Double Cap of the primary and secondary buffers
func (b *DoubleBuffer) ExpandCap() {
	// Double len of the buffer
	b.length *= 2
	b.primary.expand(b.length)
	b.secondary.expand(b.length)
}

// Write into the buffer at some offset.
// The Write will always write to the secondary buffer
// to prevent screen flicker.
func (b *DoubleBuffer) Write(data []byte) (int, error) {
	if len(data) > b.secondary.length {
		b.ExpandCap()
	}
	b.secondary.overwriteBuffer(len(data), data)
	return b.secondary.length, nil
}

// Sync Primary and secondary buffers, and then sync secondary with primary
func (b *DoubleBuffer) Sync() error {
	b.primary.overwriteBuffer(b.secondary.length, b.secondary.byteArr)
	return nil
}
