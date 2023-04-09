/*
	Copyright 2022 Phoenix

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
	"syscall"

	"github.com/Softwarekang/knetty/pkg/math"
	pool "github.com/Softwarekang/knetty/pkg/pool/ringbuffer"
	syscallutil "github.com/Softwarekang/knetty/pkg/syscall"
	"github.com/Softwarekang/knetty/pkg/utils"
)

const (
	// defaultCacheSize default cache size 64 kb.
	defaultCacheSize = 64 * KiByte
	// maxCacheSize max cache size 512 mb.
	maxCacheSize = 512 * MiByte
)

const (
	Byte = 1 << (iota * 10)
	KiByte
	MiByte
	GiByte
)

// RingBuffer an efficient, automatically resizable, and memory-reusable circular buffer implementation.
type RingBuffer struct {
	p   []byte
	r   int
	w   int
	cap int
}

// NewRingBuffer returns a  default 64 kb size circular buffer.
func NewRingBuffer() *RingBuffer {
	return NewRingBufferWithCap(defaultCacheSize)
}

// NewRingBufferWithCap returns a cap byte size circular buffer.
func NewRingBufferWithCap(cap int) *RingBuffer {
	if cap <= 0 {
		cap = defaultCacheSize
	}

	if !utils.IsPowerOfTwo(cap) {
		cap = math.Min(utils.AdjustNToPowerOfTwo(cap), maxCacheSize)
	}

	return &RingBuffer{
		p:   pool.Get(cap),
		r:   0,
		w:   0,
		cap: cap,
	}
}

// CopyFromFd read data from fd to ringBuffer, if ringBuffer is full, return retry error (EAGAIN).
// todo:(Phoenix) ringBuffer does not have an automatic expansion mechanism,
// which will cause many invalid status callbacks in poll.
// It is necessary to add an expansion strategy and provide memory multiplexing capabilities.
func (r *RingBuffer) CopyFromFd(fd int) (int, error) {
	// if ringBuffer is full, increase the double capacity  each time.
	if r.full() {
		// if the ringBuffer is already at its maximum allocatable capacity(512mb),
		// it will return the retryable error EAGAIN.
		if !r.grow(r.cap * 2) {
			return 0, syscall.EAGAIN
		}
	}

	writeIndex, readIndex := r.index(r.w), r.index(r.r)
	if writeIndex < readIndex {
		n, err := syscall.Read(fd, r.p[writeIndex:readIndex])
		if err != nil {
			return 0, err
		}

		r.w += n
		return n, nil
	}

	bs := [][]byte{
		r.p[writeIndex:],
		r.p[:readIndex],
	}
	n, err := syscallutil.Readv(fd, bs)
	if err != nil {
		return 0, err
	}

	r.w += n
	return n, nil
}

// WriteToFd write the ring Buffer data to the network, if the buffer is empty, an EAGAIN error will be returned.
// One will return other more types to error
func (r *RingBuffer) WriteToFd(fd int) (int, error) {
	if r.IsEmpty() {
		return 0, syscall.EAGAIN
	}

	writeIndex, readIndex := r.index(r.w), r.index(r.r)
	if readIndex < writeIndex {
		n, err := syscall.Write(fd, r.p[readIndex:writeIndex])
		if err != nil {
			return 0, err
		}
		r.r += n
		return n, nil
	}

	bs := [][]byte{
		r.p[readIndex:],
		r.p[:writeIndex],
	}
	n, err := syscallutil.Writev(fd, bs)
	if err != nil {
		return 0, err
	}
	r.r += n
	return n, nil
}

// Write max len(p) of data to ringBuffer, returning a retryable EAGAIN error if the buffer is full.
func (r *RingBuffer) Write(p []byte) (int, error) {
	l := len(p)
	if l <= 0 {
		return 0, nil
	}

	writeableSize := r.writeableSize()
	// if the writableSize of the ringBuffer is less than len(p), it needs to be resized.
	// if the ringBuffer has reached the maximum allocatable capacity, a retryable EAGAIN error will be returned.
	if l > writeableSize {
		if !r.grow(r.cap+l) && r.full() {
			return 0, syscall.EAGAIN
		}
	}

	writeIndex, readIndex := r.index(r.w), r.index(r.r)
	if writeIndex < readIndex {
		n := copy(r.p[writeIndex:readIndex], p)
		r.w += n
		return n, nil
	}

	writeableSize = math.Min(writeableSize, l)
	n := copy(r.p[writeIndex:], p)
	if n < writeableSize {
		n += copy(r.p[:readIndex], p[n:])
	}

	r.w += n
	return n, nil
}

// Read max len(p) of data to p, returning a retryable EAGAIN error if the buffer is empty.
func (r *RingBuffer) Read(p []byte) (int, error) {
	l := len(p)
	if l <= 0 {
		return 0, nil
	}

	if r.IsEmpty() {
		return 0, syscall.EAGAIN
	}

	writeIndex, readIndex := r.index(r.w), r.index(r.r)
	if readIndex < writeIndex {
		n := copy(p, r.p[readIndex:writeIndex])
		r.r += n
		return n, nil
	}

	readableSize := math.Min(r.readableSize(), l)
	n := copy(p, r.p[readIndex:])
	if n < readableSize {
		n += copy(p[n:], r.p[:writeIndex])
	}

	r.r += n
	return n, nil
}

// Bytes return all data in ringBuffer.
func (r *RingBuffer) Bytes() []byte {
	var zeroBytes []byte
	rw := r.w
	if rw == r.r {
		return zeroBytes
	}
	readableSize := r.readableSize()
	p := make([]byte, readableSize)
	writeIndex, readIndex := r.index(rw), r.index(r.r)
	if readIndex < writeIndex {
		copy(p, r.p[readIndex:writeIndex])
		return p
	}

	n := copy(p, r.p[readIndex:])
	if n < readableSize {
		copy(p[n:], r.p[:readIndex])
	}

	return p
}

// Len returns the number of readable bytes of the ring Buffer.
func (r *RingBuffer) Len() int {
	return r.readableSize()
}

// WriteString write string into ringBuffer.
func (r *RingBuffer) WriteString(s string) (int, error) {
	return r.Write([]byte(s))
}

// IsEmpty return true if ringBuffer is empty.
func (r *RingBuffer) IsEmpty() bool {
	return r.readableSize() == 0
}

// Release  maximum n bytes of data in ringBuffer.
func (r *RingBuffer) Release(n int) {
	if n <= 0 {
		return
	}

	if n > r.readableSize() {
		r.r = r.w
		return
	}

	r.r += n
}

// Cap return the ringBuffer capacity.
func (r *RingBuffer) Cap() int {
	return r.cap
}

// Clear up the ringBuffer.
func (r *RingBuffer) Clear() {
	r.r, r.w, r.cap, r.p = 0, 0, 0, nil
}

func (r *RingBuffer) grow(needCap int) bool {
	if needCap > maxCacheSize {
		return false
	}

	newCap := utils.AdjustNToPowerOfTwo(needCap)
	buf := pool.Get(newCap)
	n, _ := r.Read(buf)
	pool.Put(r.p)
	r.r, r.w, r.p, r.cap = 0, n, buf, newCap
	return true
}

func (r *RingBuffer) full() bool {
	return r.writeableSize() == 0
}

func (r *RingBuffer) readableSize() int {
	if r.w == r.r {
		return 0
	}

	writeIndex, readIndex := r.index(r.w), r.index(r.r)
	if writeIndex > readIndex {
		return writeIndex - readIndex
	}

	return r.cap + writeIndex - readIndex
}

func (r *RingBuffer) writeableSize() int {
	return r.cap - r.readableSize()
}

func (r *RingBuffer) index(i int) int {
	return i & (r.cap - 1)
}
