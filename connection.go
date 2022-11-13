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
	ID() uint32
	LocalAddr() string
	RemoteAddr() string
	readTimeout() time.Duration
	SetReadTimeout(time.Duration)
	writeTimeout() time.Duration
	SetWriteTimeout(time.Duration)
	close()
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
	bytes := make([]byte, 0, 256)
	_, err := syscall.Read(c.netFD.FD, bytes)
	if err != nil {
		if err != syscall.EAGAIN {
			return err
		}
	}

	c.inputBuffer.Write(bytes)
	return nil
}

type tcpConn struct {
	kNetConn
	conn net.Conn
}

func newTcpConn(conn net.Conn) *tcpConn {
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

func (t tcpConn) readTimeout() time.Duration {
	return t.readTimeOut.Load()
}

func (t tcpConn) SetReadTimeout(rTimeout time.Duration) {
	if rTimeout < 1 {
		panic("SetReadTimeout(rTimeout time.Duration):@rTimeout < 0")
	}
	t.readTimeOut = atomic.NewDuration(rTimeout)
}

func (t tcpConn) writeTimeout() time.Duration {
	return t.writeTimeOut.Load()
}

func (t tcpConn) SetWriteTimeout(wTimeout time.Duration) {
	if wTimeout < 1 {
		panic("SetWriteTimeout(wTimeout time.Duration):@wTimeout < 0")
	}

	t.writeTimeOut = atomic.NewDuration(wTimeout)
}

func (t tcpConn) writeString(str string) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (t *tcpConn) Read(n int) ([]byte, error) {
	return nil, nil
}

func (t tcpConn) close() {
	//TODO implement me
	panic("implement me")
}
