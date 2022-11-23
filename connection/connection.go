package connection

import (
	"bytes"
	mnet "github.com/Softwarekang/knet/pkg/net"
	"net"
	"syscall"
	"time"

	"github.com/Softwarekang/knet/poll"

	"go.uber.org/atomic"
)

var (
	connID atomic.Uint32
)

const (
	netIOTimeout = time.Second // 1s
)

type CloseCallBackFunc func() error

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
	// Read will return length n bytes
	Read(n int) ([]byte, error)
	// Write will write bytes to conn buffer
	WriteBuffer(bytes []byte) error
	// Flush will send conn buffer data to net
	Flush() error
	// SetCloseCallBack set close callback fun when conn on interrupt
	SetCloseCallBack(fn CloseCallBackFunc)
	// Len will return conn readable data size
	Len() int
	// Close will interrupt conn
	Close()
}

type kNetConn struct {
	id                 uint32
	fd                 int
	readTimeOut        *atomic.Duration
	writeTimeOut       *atomic.Duration
	remoteSocketAddr   syscall.Sockaddr
	localAddress       string
	remoteAddress      string
	poller             poll.Poll
	inputBuffer        bytes.Buffer
	outputBuffer       bytes.Buffer
	closeCallBackFn    CloseCallBackFunc
	waitBufferSize     atomic.Int64
	netFd              *poll.NetFileDesc
	writeNetBufferChan chan struct{}
	waitBufferChan     chan struct{}
	close              atomic.Int32
}

// Register register in poller
func (c *kNetConn) Register(eventType poll.PollEventType) error {
	c.initNetFd()
	if err := c.poller.Register(c.netFd, eventType); err != nil {
		return err
	}
	return nil
}

func (c *kNetConn) initNetFd() {
	if c.netFd != nil {
		return
	}

	c.netFd = &poll.NetFileDesc{
		FD: c.fd,
		NetPollListener: poll.NetPollListener{
			OnRead:      c.OnRead,
			OnInterrupt: c.OnInterrupt,
		},
	}
}

// OnRead refactor for conn
func (c *kNetConn) OnRead() error {
	// 0.25m bytes
	bytes := make([]byte, 256)
	n, err := syscall.Read(c.fd, bytes)
	if err != nil {
		if err != syscall.EAGAIN {
			return err
		}
	}

	c.inputBuffer.Write(bytes[:n])
	waitBufferSize := c.waitBufferSize.Load()
	if waitBufferSize > 0 && int64(c.inputBuffer.Len()) > waitBufferSize {
		c.waitBufferChan <- struct{}{}
	}
	return nil
}

// OnWrite refactor for conn
func (c *kNetConn) OnWrite() error {
	_, err := syscall.SendmsgN(c.fd, c.outputBuffer.Bytes(), nil, c.remoteSocketAddr, 0)
	if err != nil && err != syscall.EAGAIN {
		return err
	}

	return nil
}

// OnInterrupt refactor for conn
func (c *kNetConn) OnInterrupt() error {
	if err := c.poller.Register(&poll.NetFileDesc{
		FD: c.fd,
	}, poll.DeleteRead); err != nil {
		return err
	}

	if c.closeCallBackFn != nil {
		c.closeCallBackFn()
	}
	c.close.Store(1)
	return nil
}

// Conn wrapped net.conn with fd、remote sa
type Conn interface {
	net.Conn

	// FD will return conn fd
	FD() int

	// RemoteSocketAddr will return conn remote sa
	RemoteSocketAddr() syscall.Sockaddr
}

type wrappedConn struct {
	net.Conn
	remoteSocketAddr syscall.Sockaddr
	fd               int
}

// NewWrappedConn .
func NewWrappedConn(conn net.Conn) (*wrappedConn, error) {
	tcpConn := conn.(*net.TCPConn)
	file, err := tcpConn.File()
	if err != nil {
		return nil, err
	}

	tcpAddr := conn.RemoteAddr().(*net.TCPAddr)
	remoteSocketAdder, err := mnet.IPToSockAddrInet4(tcpAddr.IP, tcpAddr.Port)
	if err != nil {
		panic("")
	}

	return &wrappedConn{
		Conn:             conn,
		fd:               int(file.Fd()),
		remoteSocketAddr: remoteSocketAdder,
	}, nil
}

// FD .
func (w *wrappedConn) FD() int {
	return w.fd
}

// RemoteSocketAddr .
func (w *wrappedConn) RemoteSocketAddr() syscall.Sockaddr {
	return w.remoteSocketAddr
}
