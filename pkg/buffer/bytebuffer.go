/*
	Copyright 2022 ankangan

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package buffer

import (
	"errors"

	"github.com/Softwarekang/knetty/pkg/math"
)

const (
	kb            = 1024
	mb            = 1024 * kb
	maxBufferSize = 512 * mb
)

// ByteBuffer buffer for byte
type ByteBuffer struct {
	buf []byte
	len int
}

// NewByteBuffer .
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

/*
tryGrowSlice expand the cache capacity when it is insufficient
When the buffer capacity is less than 1mb, each buffer amplification is doubled, that is, newCap = 2*oldCap.
When the buffer capacity exceeds 1mb, the buffer is amplified by 1mb each time, that is, newCap = oldCap + 1mb.
*/
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

// Read .
func (b *ByteBuffer) Read(data []byte) (int, error) {
	if len(data) == 0 || b.len == 0 {
		return 0, nil
	}

	n := copy(data, b.buf[:b.len])
	b.buf = b.buf[n:]
	b.len = math.Max(b.len-n, 0)
	return n, nil
}

// Bytes .
func (b *ByteBuffer) Bytes() []byte {
	return b.buf[:b.len]
}

// Len .
func (b *ByteBuffer) Len() int {
	return b.len
}

// WriteString .
func (b *ByteBuffer) WriteString(s string) error {
	return b.Write([]byte(s))
}

// IsEmpty .
func (b *ByteBuffer) IsEmpty() bool {
	return b.len == 0
}

// Release .
func (b *ByteBuffer) Release(n int) {
	if cap(b.buf) < n {
		b.Clear()
		return
	}

	b.len = math.Max(b.len-n, 0)
	b.buf = b.buf[n:]
}

// Clear .
func (b *ByteBuffer) Clear() {
	b.buf = nil
	b.len = 0
}
