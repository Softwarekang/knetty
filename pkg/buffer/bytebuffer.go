package buffer

import (
	"errors"
)

const (
	kb            = 1024
	mb            = 1024 * kb
	maxBufferSize = 512 * mb
)

type ByteBuffer struct {
	buf []byte
	len int
}

func NewByteBuffer() *ByteBuffer {
	return newByteBuffer(0)
}

func newByteBuffer(cap int) *ByteBuffer {
	var bufferSize int
	if cap < 0 {
		bufferSize = 0
	}

	if cap > maxBufferSize {
		bufferSize = maxBufferSize
	}

	return &ByteBuffer{
		buf: make([]byte, bufferSize),
	}
}

// Write .
func (b *ByteBuffer) Write(data []byte) error {
	n := len(data)
	if n == 0 {
		return nil
	}

	if err := b.tryGrowSlice(n); err != nil {
		return err
	}

	l := copy(b.buf[b.len:], data)
	b.len += l
	return nil
}

func (b *ByteBuffer) tryGrowSlice(n int) error {
	newCap := n + b.len
	if newCap < cap(b.buf) {
		return nil
	}

	if newCap > maxBufferSize {
		return errors.New("buffer too large")
	}

	if newCap < mb {
		newCap *= 2
	} else {
		newCap += mb
	}

	if newCap > maxBufferSize {
		newCap = maxBufferSize
	}

	newBuffer := make([]byte, newCap)
	copy(newBuffer, b.buf)
	b.buf = newBuffer
	return nil
}

func (b *ByteBuffer) Read(data []byte) (int, error) {
	if len(data) == 0 || b.len == 0 {
		return 0, nil
	}

	n := copy(data, b.buf[:b.len])
	b.buf = b.buf[n:]
	b.len = max(b.len-n, 0)
	return n, nil
}

func (b *ByteBuffer) Bytes() []byte {
	return b.buf[:b.len]
}

func (b *ByteBuffer) Len() int {
	return b.len
}

func (b *ByteBuffer) WriteString(s string) error {
	return b.Write([]byte(s))
}

func (b *ByteBuffer) IsEmpty() bool {
	return b.len == 0
}

func (b *ByteBuffer) Release(n int) {
	if cap(b.buf) < n {
		b.Clear()
		return
	}

	b.len = max(b.len-n, 0)
	b.buf = b.buf[n:]
	return
}

func (b *ByteBuffer) Clear() {
	b.buf = nil
	b.len = 0
}

func max(a, b int) int {
	if a < b {
		return b
	}

	return a
}
