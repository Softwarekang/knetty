package connection

import (
	"context"
	"errors"
	"net"
	"syscall"
	"time"

	"github.com/Softwarekang/knetty/net/poll"
	"github.com/Softwarekang/knetty/pkg/buffer"
	merr "github.com/Softwarekang/knetty/pkg/err"
	mnet "github.com/Softwarekang/knetty/pkg/net"
	msyscall "github.com/Softwarekang/knetty/pkg/syscall"

	"go.uber.org/atomic"
)

// TcpConn tcp conn in knetty, impl Connection
type TcpConn struct {
	knettyConn
	conn net.Conn
}

// NewTcpConn .
func NewTcpConn(conn net.Conn) (*TcpConn, error) {
	if conn == nil {
		return nil, errors.New("conn is nil")
	}

	var localAddress, remoteAddress string
	if conn.LocalAddr() != nil {
		localAddress = conn.LocalAddr().String()
	}

	if conn.RemoteAddr() != nil {
		remoteAddress = conn.RemoteAddr().String()
	}

	fd, err := mnet.ResolveConnFileDesc(conn)
	if err != nil {
		return nil, err
	}

	// set conn no block
	_ = msyscall.SetConnectionNoBlock(fd)
	return &TcpConn{
		knettyConn: knettyConn{
			fd:            fd,
			readTimeOut:   atomic.NewDuration(netIOTimeout),
			writeTimeOut:  atomic.NewDuration(netIOTimeout),
			localAddress:  localAddress,
			remoteAddress: remoteAddress,
			poller:        poll.PollerManager.Pick(),
			// todo:fix use options set buffer size
			inputBuffer:        buffer.NewRingBuffer(),
			outputBuffer:       buffer.NewRingBuffer(),
			waitBufferChan:     make(chan struct{}, 1),
			writeNetBufferChan: make(chan struct{}, 1),
		},
		conn: conn,
	}, nil
}

// ID .
func (t *TcpConn) ID() uint32 {
	return t.id
}

// LocalAddr .
func (t *TcpConn) LocalAddr() string {
	return t.localAddress
}

// RemoteAddr .
func (t *TcpConn) RemoteAddr() string {
	return t.remoteAddress
}

// ReadTimeout .
func (t *TcpConn) ReadTimeout() time.Duration {
	return t.readTimeOut.Load()
}

// SetReadTimeout .
func (t *TcpConn) SetReadTimeout(rTimeout time.Duration) {
	if rTimeout < 1 {
		panic("SetReadTimeout(rTimeout time.Duration):@rTimeout < 0")
	}
	t.readTimeOut = atomic.NewDuration(rTimeout)
}

// WriteTimeout .
func (t *TcpConn) WriteTimeout() time.Duration {
	return t.writeTimeOut.Load()
}

// SetWriteTimeout .
func (t *TcpConn) SetWriteTimeout(wTimeout time.Duration) {
	if wTimeout < 1 {
		panic("SetWriteTimeout(wTimeout time.Duration):@wTimeout < 0")
	}

	t.writeTimeOut = atomic.NewDuration(wTimeout)
}

// Next .
func (t *TcpConn) Next(n int) ([]byte, error) {
	if err := t.waitReadBuffer(n, true); err != nil {
		return nil, err
	}

	p := make([]byte, n)
	if _, err := t.inputBuffer.Read(p); err != nil {
		return nil, err
	}

	return p, nil
}

// Read .
func (t *TcpConn) Read(p []byte) (int, error) {
	if err := t.waitReadBuffer(1, false); err != nil {
		return 0, err
	}

	return t.inputBuffer.Read(p)
}

func (t *TcpConn) waitReadBuffer(n int, timeout bool) error {
	if t.inputBuffer.Len() >= n {
		return nil
	}

	t.waitBufferSize.Store(int64(n))
	defer t.waitBufferSize.Store(0)
	if timeout {
		return t.waitWithTimeout(n)
	}

	for t.inputBuffer.Len() < n {
		<-t.waitBufferChan
		if !t.isActive() {
			return merr.ConnClosedErr
		}
	}

	return nil
}

func (t *TcpConn) waitWithTimeout(n int) error {
	ctx, cancel := context.WithTimeout(context.TODO(), t.readTimeOut.Load())
	defer cancel()
	for t.inputBuffer.Len() < n {
		select {
		case <-ctx.Done():
			return merr.NetIOTimeoutErr
		case <-t.waitBufferChan:
		}

		if !t.isActive() {
			return merr.ConnClosedErr
		}
	}

	return nil
}

// WriteBuffer .
func (t *TcpConn) WriteBuffer(bytes []byte) (int, error) {
	return t.outputBuffer.Write(bytes)
}

// FlushBuffer .
func (t *TcpConn) FlushBuffer() error {
	if _, err := t.outputBuffer.WriteToFd(t.fd); err != nil && err != syscall.EAGAIN {
		return err
	}

	if t.outputBuffer.IsEmpty() {
		return nil
	}

	// net buffer is full
	if err := t.Register(poll.ReadToRW); err != nil {
		return err
	}

	<-t.writeNetBufferChan
	return nil
}

// Len .
func (t *TcpConn) Len() int {
	return t.inputBuffer.Len()
}

func (t *TcpConn) isActive() bool {
	return t.close.Load() == 0
}

// SetCloseCallBack .
func (t *TcpConn) SetCloseCallBack(fn CloseCallBackFunc) {
	t.closeCallBackFn = fn
}

// Close .
func (t *TcpConn) Close() error {
	if !t.isActive() {
		return nil
	}

	return t.OnInterrupt()
}

// Type .
func (t *TcpConn) Type() ConnType {
	return TCPCONNECTION
}
