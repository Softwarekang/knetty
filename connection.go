package knet

import (
	"bytes"
	"net"
	"syscall"
	"time"

	"go.uber.org/atomic"
)

var (
	connID atomic.Uint32
)

const (
	netIOTimeout = time.Second // 1s
)

// Connection some connection  operations
type Connection interface {
	// ID for conn
	ID() uint32
	// LocalAddr local address for conn
	LocalAddr() string
	// RemoteAddr remote address for conn
	RemoteAddr() string
	// ReadTimeout timeout for read
	ReadTimeout() time.Duration
	// SetReadTimeout setup read timeout
	SetReadTimeout(time.Duration)
	// WriteTimeout timeout for write
	WriteTimeout() time.Duration
	// SetWriteTimeout setup write timeout
	SetWriteTimeout(time.Duration)
	// Close will interrupt conn
	Close()
}

// Conn net.conn with fd
type Conn interface {
	net.Conn

	// FD will return conn fd
	FD() int
}

type wrappedConn struct {
	net.Conn
	fd int
}

// NewWrappedConn .
func NewWrappedConn(conn net.Conn) (*wrappedConn, error) {
	file, err := conn.(*net.TCPConn).File()
	if err != nil {
		return nil, err
	}
	return &wrappedConn{
		Conn: conn,
		fd:   int(file.Fd()),
	}, nil
}

// FD .
func (w *wrappedConn) FD() int {
	return w.fd
}

type kNetConn struct {
	id            uint32
	readTimeOut   *atomic.Duration
	writeTimeOut  *atomic.Duration
	localAddress  string
	remoteAddress string
	netFD         *NetFileDesc
	poller        Poll
	inputBuffer   bytes.Buffer
}

// RegisterPoller register in poller
func (c *kNetConn) RegisterPoller() error {
	c.netFD.OnRead = c.OnRead
	if err := c.poller.Register(c.netFD, Read); err != nil {
		return err
	}
	return nil
}

// OnRead refactor for conn
func (c *kNetConn) OnRead() error {
	// 0.25m bytes
	bytes := make([]byte, 256)
	n, err := syscall.Read(c.netFD.FD, bytes)
	if err != nil {
		if err != syscall.EAGAIN {
			return err
		}
	}

	c.inputBuffer.Write(bytes[:n])
	return nil
}

type tcpConn struct {
	kNetConn
	conn net.Conn
}

func NewTcpConn(conn Conn) *tcpConn {
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

	return &tcpConn{
		kNetConn: kNetConn{
			id:            connID.Inc(),
			readTimeOut:   atomic.NewDuration(netIOTimeout),
			writeTimeOut:  atomic.NewDuration(netIOTimeout),
			localAddress:  localAddress,
			remoteAddress: remoteAddress,
			poller:        PollerManager.Pick(),
			netFD: &NetFileDesc{
				FD: conn.FD(),
			},
		},
		conn: conn,
	}
}

func (t tcpConn) ID() uint32 {
	return t.id
}

func (t tcpConn) LocalAddr() string {
	return t.localAddress
}

func (t tcpConn) RemoteAddr() string {
	return t.remoteAddress
}

func (t tcpConn) ReadTimeout() time.Duration {
	return t.readTimeOut.Load()
}

func (t tcpConn) SetReadTimeout(rTimeout time.Duration) {
	if rTimeout < 1 {
		panic("SetReadTimeout(rTimeout time.Duration):@rTimeout < 0")
	}
	t.readTimeOut = atomic.NewDuration(rTimeout)
}

func (t tcpConn) WriteTimeout() time.Duration {
	return t.writeTimeOut.Load()
}

func (t tcpConn) SetWriteTimeout(wTimeout time.Duration) {
	if wTimeout < 1 {
		panic("SetWriteTimeout(wTimeout time.Duration):@wTimeout < 0")
	}

	t.writeTimeOut = atomic.NewDuration(wTimeout)
}

func (t tcpConn) WriteString(str string) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (t tcpConn) Close() {
	//TODO implement me
	panic("implement me")
}
