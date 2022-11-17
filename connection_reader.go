package knet

import (
	"context"
	"fmt"
)

// Reader interface for conn
type Reader interface {
	// ReadString will return str length n
	ReadString(n int) (string, error)
	// ReadBytes will return bytes length n
	ReadBytes(n int) ([]byte, error)
	// Len will return readable data size
	Len() int
}

// ReadString .
func (t *tcpConn) ReadString(n int) (string, error) {
	bytes, err := t.ReadBytes(n)
	if err != nil {
		return "", err
	}

	return string(bytes), err
}

// ReadBytes .
func (t *tcpConn) ReadBytes(n int) ([]byte, error) {
	if err := t.waitReadBuffer(n); err != nil {
		return nil, err
	}

	return t.read(n)
}

func (t *tcpConn) waitReadBuffer(n int) error {
	if t.inputBuffer.Len() >= n {
		return nil
	}

	t.waitBufferSize.Store(int64(n))
	defer t.waitBufferSize.Store(0)
	if t.inputBuffer.Len() >= n {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.TODO(), t.readTimeOut.Load())
	defer cancel()
	for t.inputBuffer.Len() < n {
		if !t.isActive() {
			return fmt.Errorf("waitReadBufferWithTimeout conn is closed")
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("waitReadBufferWithTimeout ctx timeout")
		case <-t.waitBufferChan:
			continue
		}
	}

	return nil
}

func (t *tcpConn) read(n int) ([]byte, error) {
	data := make([]byte, n)
	n, err := t.inputBuffer.Read(data)
	if err != nil {
		return nil, err
	}

	fmt.Printf("read %d length data from input buffer", n)
	return data, nil
}

func (t *tcpConn) isActive() bool {
	return t.close.Load() == 0
}

// Len .
func (t *tcpConn) Len() int {
	return t.inputBuffer.Len()
}
