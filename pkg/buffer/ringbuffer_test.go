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
	"net"
	"reflect"
	"sync"
	"testing"

	errors "github.com/Softwarekang/knetty/pkg/err"

	"github.com/stretchr/testify/assert"
)

func TestNewRingBuffer(t *testing.T) {
	type args struct {
		cap int
	}
	tests := []struct {
		name    string
		args    args
		wantCap int
	}{
		{
			name:    "normal",
			args:    args{cap: 1},
			wantCap: 1,
		},
		{
			name:    "normal cap != 2^n",
			args:    args{cap: 3},
			wantCap: 4,
		},
		{
			name:    "normal cap = maxCacheSize",
			args:    args{cap: maxCacheSize},
			wantCap: maxCacheSize,
		},
		{
			name:    "cap <= 0",
			args:    args{cap: -1},
			wantCap: defaultCacheSize,
		},
		{
			name:    "cap > maxCacheSize",
			args:    args{cap: maxCacheSize + 1},
			wantCap: maxCacheSize,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRingBufferWithCap(tt.args.cap); !reflect.DeepEqual(got.Cap(), tt.wantCap) {
				t.Errorf("NewRingBuffer() cap= %v, want %v", got, tt.wantCap)
			}
		})
	}
}

func TestRingBuffer_Write(t *testing.T) {
	var (
		n   int
		err error
	)

	defaultRingBuffer := NewRingBuffer()
	assert.Equal(t, defaultCacheSize, defaultRingBuffer.Cap())

	ringBuffer := NewRingBufferWithCap(10)
	assert.Equal(t, true, ringBuffer.IsEmpty())

	n, err = ringBuffer.WriteString("helloworld")
	assert.Nil(t, err)
	assert.Equal(t, 10, n)

	assert.Equal(t, "helloworld", string(ringBuffer.Bytes()))
	assert.Equal(t, 10, ringBuffer.Len())
	assert.Equal(t, 16, ringBuffer.Cap())

	n, err = ringBuffer.WriteString("helloworld")
	assert.Nil(t, err)
	assert.Equal(t, 10, n)
	assert.Equal(t, 32, ringBuffer.Cap())

	assert.Equal(t, "helloworldhelloworld", string(ringBuffer.Bytes()))
	assert.Equal(t, 20, ringBuffer.Len())

	ringBuffer.Release(5)

	assert.Equal(t, "worldhelloworld", string(ringBuffer.Bytes()))
	assert.Equal(t, 15, ringBuffer.Len())

	n, err = ringBuffer.WriteString("123456")
	assert.Equal(t, 6, n)
	assert.Nil(t, err)

	assert.Equal(t, "worldhelloworld123456", string(ringBuffer.Bytes()))
	assert.Equal(t, 21, ringBuffer.Len())

	ringBuffer.Release(10)

	assert.Equal(t, "world123456", string(ringBuffer.Bytes()))
	assert.Equal(t, 11, ringBuffer.Len())

	ringBuffer.Release(ringBuffer.Len() + 1)
	var emptyValue []byte
	assert.Equal(t, emptyValue, ringBuffer.Bytes())
	assert.Equal(t, 0, ringBuffer.Len())

	ringBuffer.Release(1)

	assert.Equal(t, 0, ringBuffer.Len())

	n, err = ringBuffer.Write(nil)
	assert.Equal(t, 0, n)
	assert.Nil(t, err)
}

func TestRingBuffer_Read(t *testing.T) {
	var (
		n    int
		err  error
		rBuf []byte
	)

	ringBuffer := NewRingBufferWithCap(10)
	assert.Equal(t, true, ringBuffer.IsEmpty())

	n, err = ringBuffer.WriteString("helloworld")
	assert.Nil(t, err)
	assert.Equal(t, 10, n)

	assert.Equal(t, "helloworld", string(ringBuffer.Bytes()))
	assert.Equal(t, 10, ringBuffer.Len())
	assert.Equal(t, 16, ringBuffer.Cap())

	n, err = ringBuffer.WriteString("helloworld")
	assert.Nil(t, err)
	assert.Equal(t, 10, n)

	assert.Equal(t, "helloworldhelloworld", string(ringBuffer.Bytes()))
	assert.Equal(t, 20, ringBuffer.Len())

	rBuf = make([]byte, 5)
	n, err = ringBuffer.Read(rBuf)
	assert.Nil(t, err)
	assert.Equal(t, 5, n)

	assert.Equal(t, "worldhelloworld", string(ringBuffer.Bytes()))
	assert.Equal(t, 15, ringBuffer.Len())

	n, err = ringBuffer.WriteString("123456")
	assert.Equal(t, 6, n)
	assert.Nil(t, err)

	assert.Equal(t, "worldhelloworld123456", string(ringBuffer.Bytes()))
	assert.Equal(t, 21, ringBuffer.Len())

	rBuf = make([]byte, 10)
	n, err = ringBuffer.Read(rBuf)
	assert.Nil(t, err)
	assert.Equal(t, 10, n)

	assert.Equal(t, "world123456", string(ringBuffer.Bytes()))
	assert.Equal(t, 11, ringBuffer.Len())

	n, err = ringBuffer.WriteString("789")
	assert.Equal(t, 3, n)
	assert.Nil(t, err)

	assert.Equal(t, "world123456789", string(ringBuffer.Bytes()))
	assert.Equal(t, 14, ringBuffer.Len())

	rBuf = make([]byte, ringBuffer.Len()+1)
	n, err = ringBuffer.Read(rBuf)
	assert.Nil(t, err)
	assert.Equal(t, 14, n)

	var emptyValue []byte
	assert.Equal(t, emptyValue, ringBuffer.Bytes())
	assert.Equal(t, 0, ringBuffer.Len())

	rBuf = make([]byte, 1)
	n, err = ringBuffer.Read(rBuf)
	assert.Equal(t, errors.BufferEmptyErr, err)
	assert.Equal(t, 0, n)

	assert.Equal(t, 0, ringBuffer.Len())

	n, err = ringBuffer.Write(nil)
	assert.Equal(t, 0, n)
	assert.Nil(t, err)

	n, err = ringBuffer.WriteString("hello")
	assert.Nil(t, err)
	assert.Equal(t, 5, n)

	rBuf = make([]byte, 2)
	n, err = ringBuffer.Read(rBuf)
	assert.Nil(t, err)
	assert.Equal(t, 2, n)

	assert.Equal(t, "llo", string(ringBuffer.Bytes()))
	assert.Equal(t, 3, ringBuffer.Len())

	n, err = ringBuffer.Read(nil)
	assert.Nil(t, err)
	assert.Equal(t, 0, n)

	ringBuffer.Clear()
}

func TestRingBuffer_CopyFromFd(t *testing.T) {
	g := sync.WaitGroup{}
	w := make(chan struct{}, 1)

	var (
		wconn net.Conn
		fd    int
	)
	g.Add(2)
	go func() {
		l, err := net.Listen("tcp", "127.0.0.1:10000")
		if err != nil {
			t.Errorf("net listen err:%v", err)
		}

		w <- struct{}{}
		rconn, err := l.Accept()
		if err != nil {
			t.Errorf("net accept err:%v", err)
		}

		f, _ := rconn.(*net.TCPConn).File()
		fd = int(f.Fd())
		g.Done()
	}()

	go func() {
		<-w
		var err error
		wconn, err = net.Dial("tcp", "127.0.0.1:10000")
		if err != nil {
			t.Errorf("net dial err:%v", err)
		}

		g.Done()
	}()

	g.Wait()

	n, err := wconn.Write([]byte("hello"))
	assert.Equal(t, 5, n)
	assert.Nil(t, err)

	ringBuffer := NewRingBufferWithCap(5)
	n, err = ringBuffer.CopyFromFd(fd)
	assert.Equal(t, 5, n)
	assert.Nil(t, err)
	assert.Equal(t, "hello", string(ringBuffer.Bytes()))

	ringBuffer.Release(2)
	assert.Equal(t, "llo", string(ringBuffer.Bytes()))

	n, err = wconn.Write([]byte("12345"))
	assert.Equal(t, 5, n)
	assert.Nil(t, err)

	n, err = ringBuffer.CopyFromFd(fd)
	assert.Equal(t, 5, n)
	assert.Nil(t, err)
	assert.Equal(t, "llo12345", string(ringBuffer.Bytes()))

	ringBuffer.Release(2)
	assert.Equal(t, "o12345", string(ringBuffer.Bytes()))

	n, err = wconn.Write([]byte("67891"))
	assert.Equal(t, 5, n)
	assert.Nil(t, err)

	n, err = ringBuffer.CopyFromFd(fd)
	assert.Equal(t, 2, n)
	assert.Nil(t, err)
	assert.Equal(t, "o1234567", string(ringBuffer.Bytes()))

	n, err = ringBuffer.CopyFromFd(fd)
	assert.Equal(t, 3, n)
	assert.Nil(t, err)
	assert.Equal(t, "o1234567891", string(ringBuffer.Bytes()))
	assert.Equal(t, 16, ringBuffer.Cap())
}

func TestRingBuffer_WriteToFd(t *testing.T) {
	g := sync.WaitGroup{}
	w := make(chan struct{}, 1)

	var (
		rfd int
		wfd int
	)
	g.Add(2)
	go func() {
		l, err := net.Listen("tcp", "127.0.0.1:10001")
		if err != nil {
			t.Errorf("net listen err:%v", err)
		}

		w <- struct{}{}
		rconn, err := l.Accept()
		if err != nil {
			t.Errorf("net accept err:%v", err)
		}

		f, _ := rconn.(*net.TCPConn).File()
		rfd = int(f.Fd())
		g.Done()
	}()

	go func() {
		<-w
		wconn, err := net.Dial("tcp", "127.0.0.1:10001")
		if err != nil {
			t.Errorf("net dial err:%v", err)
		}

		f, _ := wconn.(*net.TCPConn).File()
		wfd = int(f.Fd())
		g.Done()
	}()

	g.Wait()

	readRingBuffer := NewRingBufferWithCap(5)
	writeRingBuffer := NewRingBufferWithCap(10)

	n, err := writeRingBuffer.WriteString("hello")
	assert.Equal(t, 5, n)
	assert.Nil(t, err)

	n, err = writeRingBuffer.WriteToFd(wfd)
	assert.Equal(t, 5, n)
	assert.Nil(t, err)

	n, err = readRingBuffer.CopyFromFd(rfd)
	assert.Equal(t, 5, n)
	assert.Nil(t, err)
	assert.Equal(t, "hello", string(readRingBuffer.Bytes()))

	readRingBuffer.Release(2)
	assert.Equal(t, "llo", string(readRingBuffer.Bytes()))

	n, err = writeRingBuffer.WriteString("12345")
	assert.Equal(t, 5, n)
	assert.Nil(t, err)

	n, err = writeRingBuffer.WriteToFd(wfd)
	assert.Equal(t, 5, n)
	assert.Nil(t, err)

	n, err = readRingBuffer.CopyFromFd(rfd)
	assert.Equal(t, 5, n)
	assert.Nil(t, err)
	assert.Equal(t, "llo12345", string(readRingBuffer.Bytes()))

	readRingBuffer.Release(2)
	assert.Equal(t, "o12345", string(readRingBuffer.Bytes()))

	n, err = writeRingBuffer.WriteString("67891")
	assert.Equal(t, 5, n)
	assert.Nil(t, err)

	n, err = writeRingBuffer.WriteToFd(wfd)
	assert.Equal(t, 5, n)
	assert.Nil(t, err)

	n, err = readRingBuffer.CopyFromFd(rfd)
	assert.Equal(t, 2, n)
	assert.Nil(t, err)
	assert.Equal(t, "o1234567", string(readRingBuffer.Bytes()))

	n, err = readRingBuffer.CopyFromFd(rfd)
	assert.Equal(t, 3, n)
	assert.Nil(t, err)
	assert.Equal(t, 16, readRingBuffer.Cap())
	assert.Equal(t, "o1234567891", string(readRingBuffer.Bytes()))

	n, err = writeRingBuffer.WriteString("23456")
	assert.Equal(t, 5, n)
	assert.Nil(t, err)

	n, err = writeRingBuffer.WriteToFd(wfd)
	assert.Equal(t, 5, n)
	assert.Nil(t, err)

	n, err = writeRingBuffer.WriteToFd(wfd)
	assert.Equal(t, 0, n)
	assert.Nil(t, err)

}
