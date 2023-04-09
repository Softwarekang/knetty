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

// Package ring_buffer  implement buffer cache pool to improve memory reuse.
package ring_buffer

import (
	"sync"

	"github.com/Softwarekang/knetty/pkg/utils"
)

const (
	// The minimum buffer is 1 byte.
	minAllocBit  = 0
	minAllocSize = 1 << minAllocBit
	// The maximum buffer size is 512 mb.
	maxAllocBit  = 29
	maxAllocSize = 1 << maxAllocBit
)

var (
	ringBufferPool = map[int]*sync.Pool{}
)

func init() {
	for i := minAllocBit; i <= maxAllocBit; i++ {
		size := 1 << i
		ringBufferPool[size] = &sync.Pool{
			New: func() any {
				buf := make([]byte, size)
				return &buf
			},
		}
	}
}

// Get the ringBuffer pool will return []byte of size is the first value greater than or equal its 2^n.
// newCap  between  1 byte. and 512 mb。
func Get(size int) []byte {
	if size < minAllocSize {
		size = minAllocSize
	}

	if size >= maxAllocSize {
		size = maxAllocSize
	}

	size = utils.AdjustNToPowerOfTwo(size)
	bufPointer := ringBufferPool[size].Get().(*[]byte)
	return (*bufPointer)[:size]
}

// Put if buf is legal then it will be put into the buffer pool.
func Put(buf []byte) {
	c := cap(buf)
	if !utils.IsPowerOfTwo(c) {
		return
	}

	if c > maxAllocSize || c < minAllocSize {
		return
	}
	releasedBuf := buf[:0]
	ringBufferPool[c].Put(&releasedBuf)
}
