package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBufferWrite(t *testing.T) {
	buffer := NewDoubleBuf()

	l, error := buffer.Write([]byte("Hello World!!"))
	require.Nil(t, error, "Failed to write to buffer")

	assert.Equal(t, "Hello World!!", string(buffer.secondary.byteArr[:l]), "Incorrect data written to secondary buffer")
	assert.Equal(t, "", string(buffer.primary.byteArr[:buffer.primary.length]), "Data written into primary buffer")
}

func TestBufferSync(t *testing.T) {
	buffer := NewDoubleBuf()

	l, error := buffer.Write([]byte("Hello World!!"))
	require.Nil(t, error, "Failed to write to buffer")

	assert.Equal(t, "Hello World!!", string(buffer.secondary.byteArr[:l]), "Incorrect data written to secondary buffer")
	assert.Equal(t, "", string(buffer.primary.byteArr[:buffer.primary.length]), "Data written into primary buffer")

	error = buffer.Sync()
	require.Nil(t, error, "Failed to swap buffers")

	assert.Equal(t, "Hello World!!", string(buffer.primary.byteArr[:buffer.primary.length]), "Primary buffer has incorrect data")
	assert.Equal(t, "Hello World!!", string(buffer.secondary.byteArr[:buffer.secondary.length]), "Secondary buffer data corrupted in sync")
}

func TestBufferWriteExtra(t *testing.T) {
	buffer := NewDoubleBuf()

	l, error := buffer.Write([]byte("Hello World!!"))
	require.Nil(t, error, "Failed to write to buffer")

	assert.Equal(t, "Hello World!!", string(buffer.secondary.byteArr[:l]), "Incorrect data written to secondary buffer")
	assert.Equal(t, "", string(buffer.primary.byteArr[:buffer.primary.length]), "Data written into primary buffer")

	error = buffer.Sync()
	require.Nil(t, error, "Failed to swap buffers")

	assert.Equal(t, "Hello World!!", string(buffer.primary.byteArr[:buffer.primary.length]), "Data not swap between buffers")
	assert.Equal(t, "Hello World!!", string(buffer.secondary.byteArr[:buffer.secondary.length]), "Data not swapped between buffers")

	l, error = buffer.Write([]byte("Hello World!!               \n\n\t\n\n\n\n\n          something"))
	require.Nil(t, error, "Failed to write to buffer")

	assert.Equal(t, "Hello World!!               \n\n\t\n\n\n\n\n          something", string(buffer.secondary.byteArr[:l]), "Incorrect data written to secondary buffer")
	assert.Equal(t, "Hello World!!", string(buffer.primary.byteArr[:buffer.primary.length]), "Data written into primary buffer")

	error = buffer.Sync()
	require.Nil(t, error, "Failed to swap buffers")

	assert.Equal(t, "Hello World!!               \n\n\t\n\n\n\n\n          something", string(buffer.primary.byteArr[:buffer.primary.length]), "Data not swap between buffers")
	assert.Equal(t, "Hello World!!               \n\n\t\n\n\n\n\n          something", string(buffer.secondary.byteArr[:buffer.secondary.length]), "Data not swapped between buffers")

}
