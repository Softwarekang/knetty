// Package connection  implements tcp, udp and other protocols for network connection
package connection

import (
	"time"

	"github.com/Softwarekang/knetty/net/poll"
	"github.com/Softwarekang/knetty/pkg/buffer"
	"go.uber.org/atomic"
)

const (
	// default timeout for net io
	netIOTimeout = time.Second // 1s
)

type ConnType int

// define tcp、upd、webSocket conn
const (
	TCPCONNECTION ConnType = iota
	UDPCONNECTION
	WEBSOCKETCONNECTION
)

// CloseCallBackFunc will runs at conn on Interrupt
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
	// Next will return length n bytes
	Next(n int) ([]byte, error)
	// Read will return max len(p) data
	Read(p []byte) (int, error)
	// WriteBuffer will write bytes to conn buffer
	WriteBuffer(bytes []byte) (int, error)
	// FlushBuffer will send conn buffer data to net
	FlushBuffer() error
	// SetCloseCallBack set close callback fun when conn on interrupt
	SetCloseCallBack(fn CloseCallBackFunc)
	// Len will return conn readable data size
	Len() int
	// Type  will return conn type
	Type() ConnType
	// Close will interrupt conn
	Close() error
}

type knettyConn struct {
	id                 uint32
	fd                 int
	readTimeOut        *atomic.Duration
	writeTimeOut       *atomic.Duration
	localAddress       string
	remoteAddress      string
	poller             poll.Poll
	inputBuffer        *buffer.RingBuffer
	outputBuffer       *buffer.RingBuffer
	closeCallBackFn    CloseCallBackFunc
	waitBufferSize     atomic.Int64
	netFd              *poll.NetFileDesc
	writeNetBufferChan chan struct{}
	waitBufferChan     chan struct{}
	close              atomic.Int32
}

// Register conn in poller
func (c *knettyConn) Register(eventType poll.EventType) error {
	c.initNetFd()
	if err := c.poller.Register(c.netFd, eventType); err != nil {
		return err
	}
	return nil
}

func (c *knettyConn) initNetFd() {
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

func (c *knettyConn) closeWaitBufferCh() {
	select {
	case <-c.waitBufferChan:
	default:
		close(c.waitBufferChan)
	}
}
