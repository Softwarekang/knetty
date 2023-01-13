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
	"syscall"

	"github.com/Softwarekang/knetty/pkg/math"
	msycall "github.com/Softwarekang/knetty/pkg/syscall"
)

const (
	// defaultCacheSize default cache size 64 kb
	defaultCacheSize = 64 * K
	// maxCacheSize max cache size 1 GB
	maxCacheSize = 1 * G
)

const (
	B = 1
	K = 1024 * B
	M = 1024 * K
	G = 1024 * M
)

// RingBuffer lock-free cache for a read-write goroutine.
type RingBuffer struct {
	p   []byte
	r   int
	w   int
	cap int
}

// NewRingBuffer .
func NewRingBuffer() *RingBuffer {
	return NewRingBufferWithCap(defaultCacheSize)
}

// NewRingBufferWithCap .
func NewRingBufferWithCap(cap int) *RingBuffer {
	if cap <= 0 {
		cap = defaultCacheSize
	}

	if (cap & (cap - 1)) != 0 {
		cap = math.Min(adjust(cap), maxCacheSize)
	}

	return &RingBuffer{
		p:   make([]byte, cap),
		r:   0,
		w:   0,
		cap: cap,
	}
}

// CopyFromFd .
func (r *RingBuffer) CopyFromFd(fd int) (int, error) {
	rr := r.r
	if r.full(rr) {
		return 0, syscall.EAGAIN
	}

	writeIndex, readIndex := r.index(r.w), r.index(rr)
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
	n, err := msycall.Readv(fd, bs)
	if err != nil {
		return 0, err
	}

	r.w += n
	return n, nil
}

// WriteToFd .
func (r *RingBuffer) WriteToFd(fd int) (int, error) {
	rw := r.w
	if rw == r.r {
		return 0, syscall.EAGAIN
	}

	writeIndex, readIndex := r.index(rw), r.index(r.r)
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
	n, err := msycall.Writev(fd, bs)
	if err != nil {
		return 0, err
	}
	r.r += n
	return n, nil
}

// Write .
func (r *RingBuffer) Write(p []byte) (int, error) {
	rr := r.r
	if r.full(rr) {
		return 0, syscall.EAGAIN
	}

	l := len(p)
	if l <= 0 {
		return 0, nil
	}

	writeIndex, readIndex := r.index(r.w), r.index(rr)
	if writeIndex < readIndex {
		n := copy(r.p[writeIndex:readIndex], p)
		r.w += n
		return n, nil
	}

	writeableSize := math.Min(r.cap+readIndex-writeIndex, l)
	n := copy(r.p[writeIndex:], p)
	if n < writeableSize {
		n += copy(r.p[:readIndex], p[n:])
	}

	r.w += n
	return n, nil
}

// Read .
func (r *RingBuffer) Read(p []byte) (int, error) {
	rw := r.w
	if rw == r.r {
		return 0, syscall.EAGAIN
	}

	l := len(p)
	if l <= 0 {
		return 0, nil
	}

	writeIndex, readIndex := r.index(rw), r.index(r.r)
	if readIndex < writeIndex {
		n := copy(p, r.p[readIndex:writeIndex])
		r.r += n
		return n, nil
	}

	readableSize := math.Min(r.readableSize(rw), l)
	n := copy(p, r.p[readIndex:])
	if n < readableSize {
		n += copy(p[n:], r.p[:writeIndex])
	}

	r.r += n
	return n, nil
}

// Bytes .
func (r *RingBuffer) Bytes() []byte {
	var zeroBytes []byte
	rw := r.w
	if rw == r.r {
		return zeroBytes
	}
	readableSize := r.readableSize(rw)
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

// Len .
func (r *RingBuffer) Len() int {
	return r.readableSize(r.w)
}

// WriteString .
func (r *RingBuffer) WriteString(s string) (int, error) {
	return r.Write([]byte(s))
}

// IsEmpty .
func (r *RingBuffer) IsEmpty() bool {
	return r.r == r.w
}

// Release .
func (r *RingBuffer) Release(n int) {
	rw := r.w
	releasableSize := rw - r.r
	if releasableSize == 0 {
		return
	}

	if n > releasableSize {
		r.r = rw
		return
	}

	r.r += n
}

// Cap .
func (r *RingBuffer) Cap() int {
	return r.cap
}

// Clear .
func (r *RingBuffer) Clear() {
	r.r, r.w, r.cap, r.p = 0, 0, 0, nil
}

func (r *RingBuffer) full(rr int) bool {
	return r.w-rr == r.cap
}

func (r *RingBuffer) readableSize(rw int) int {
	if rw == r.r {
		return 0
	}

	writeIndex, readIndex := r.index(rw), r.index(r.r)
	if writeIndex > readIndex {
		return writeIndex - readIndex
	}

	return r.cap + writeIndex - readIndex
}

func (r *RingBuffer) index(i int) int {
	return i & (r.cap - 1)
}

func adjust(n int) int {
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	return n + 1
}
