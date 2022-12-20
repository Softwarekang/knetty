package connection

import (
	"context"
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
		panic("newTcpConn(conn net.Conn):@conn is nil")
	}

	var localAddress, remoteAddress string
	if conn.LocalAddr() != nil {
		localAddress = conn.LocalAddr().String()
	}

	if conn.RemoteAddr() != nil {
		remoteAddress = conn.RemoteAddr().String()
	}

	tcpConn := conn.(*net.TCPConn)
	file, err := tcpConn.File()
	if err != nil {
		return nil, err
	}

	tcpAddr := conn.RemoteAddr().(*net.TCPAddr)
	remoteSocketAdder, err := mnet.IPToSockAddrInet4(tcpAddr.IP, tcpAddr.Port)
	if err != nil {
		return nil, err
	}

	// set conn no block
	_ = msyscall.SetConnectionNoBlock(int(file.Fd()))
	return &TcpConn{
		knettyConn: knettyConn{
			fd:                 int(file.Fd()),
			remoteSocketAddr:   remoteSocketAdder,
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
	for {
		if !t.isActive() {
			return 0, merr.ConnClosedErr
		}
		if t.inputBuffer.Len() == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		return t.inputBuffer.Read(p)
	}
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
func (t *TcpConn) Close() {
	if !t.isActive() {
		return
	}
	if err := t.OnInterrupt(); err != nil {
		return
	}
}

// Type .
func (t *TcpConn) Type() ConnType {
	return TCPCONNECTION
}
