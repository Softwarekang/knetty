package connection

import (
	"context"
	"net"
	"syscall"
	"time"

	"github.com/Softwarekang/knet/pkg/buffer"
	merr "github.com/Softwarekang/knet/pkg/err"
	"github.com/Softwarekang/knet/poll"
	msyscall "github.com/Softwarekang/knet/syscall"

	"go.uber.org/atomic"
)

type TcpConn struct {
	kNetConn
	conn net.Conn
}

func NewTcpConn(conn Conn) *TcpConn {
	if conn == nil {
		panic("newTcpConn(conn net.Conn):@conn is nil")
	}

	var localAddress, remoteAddress string
	if conn.LocalAddr() != nil {
		localAddress = conn.LocalAddr().String()
	}

	if conn.RemoteAddr() != nil {
		remoteAddress = conn.RemoteAddr().String()
	}

	// set conn no block
	msyscall.SetConnectionNoBlock(conn.FD())
	return &TcpConn{
		kNetConn: kNetConn{
			fd:                 conn.FD(),
			remoteSocketAddr:   conn.RemoteSocketAddr(),
			readTimeOut:        atomic.NewDuration(netIOTimeout),
			writeTimeOut:       atomic.NewDuration(netIOTimeout),
			localAddress:       localAddress,
			remoteAddress:      remoteAddress,
			poller:             poll.PollerManager.Pick(),
			inputBuffer:        buffer.NewByteBuffer(),
			outputBuffer:       buffer.NewByteBuffer(),
			waitBufferChan:     make(chan struct{}, 1),
			writeNetBufferChan: make(chan struct{}, 1),
		},
		conn: conn,
	}
}

func (t TcpConn) ID() uint32 {
	return t.id
}

func (t TcpConn) LocalAddr() string {
	return t.localAddress
}

func (t TcpConn) RemoteAddr() string {
	return t.remoteAddress
}

func (t TcpConn) ReadTimeout() time.Duration {
	return t.readTimeOut.Load()
}

func (t TcpConn) SetReadTimeout(rTimeout time.Duration) {
	if rTimeout < 1 {
		panic("SetReadTimeout(rTimeout time.Duration):@rTimeout < 0")
	}
	t.readTimeOut = atomic.NewDuration(rTimeout)
}

func (t TcpConn) WriteTimeout() time.Duration {
	return t.writeTimeOut.Load()
}

func (t TcpConn) SetWriteTimeout(wTimeout time.Duration) {
	if wTimeout < 1 {
		panic("SetWriteTimeout(wTimeout time.Duration):@wTimeout < 0")
	}

	t.writeTimeOut = atomic.NewDuration(wTimeout)
}

// Next .
func (t *TcpConn) Next(n int) ([]byte, error) {
	if err := t.waitReadBuffer(n); err != nil {
		return nil, err
	}

	p := make([]byte, n)
	if _, err := t.read(p); err != nil {
		return nil, err
	}

	return p, nil
}

// Read .
func (t *TcpConn) Read(p []byte) (int, error) {
	return t.read(p)
}

func (t *TcpConn) waitReadBuffer(n int) error {
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
			return merr.ConnClosedErr
		}

		select {
		case <-ctx.Done():
			return merr.NetIOTimeoutErr
		case <-t.waitBufferChan:
			continue
		}
	}

	return nil
}

func (t *TcpConn) read(p []byte) (int, error) {
	return t.inputBuffer.Read(p)
}

// WriteBuffer .
func (t *TcpConn) WriteBuffer(bytes []byte) error {
	return t.outputBuffer.Write(bytes)
}

// FlushBuffer .
func (t *TcpConn) FlushBuffer() error {
	n, err := syscall.SendmsgN(t.fd, t.outputBuffer.Bytes(), nil, t.remoteSocketAddr, 0)
	if err != nil && err != syscall.EAGAIN {
		return err
	}

	t.outputBuffer.Release(n)
	if t.outputBuffer.Len() == 0 {
		return nil
	}

	// net buffer is full
	if err := t.Register(poll.Write); err != nil {
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
func (t *TcpConn) Close() {
	if !t.isActive() {
		return
	}
	t.OnInterrupt()
}

// Type .
func (t *TcpConn) Type() ConnType {
	return TCPCONNECTION
}
