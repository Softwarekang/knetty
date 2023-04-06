package connection

import (
	"github.com/Softwarekang/knetty/net/poll"
	"github.com/Softwarekang/knetty/pkg/buffer"

	"go.uber.org/atomic"
)

type knettyConn struct {
	id                 uint32
	fd                 int
	localAddress       string
	remoteAddress      string
	poller             poll.Poll
	inputBuffer        *buffer.RingBuffer
	outputBuffer       *buffer.RingBuffer
	netFd              *poll.NetFileDesc
	eventTrigger       EventTrigger
	writeNetBufferChan chan struct{}
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
