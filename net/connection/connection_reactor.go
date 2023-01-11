package connection

import (
	"syscall"

	"github.com/Softwarekang/knetty/net/poll"
)

// OnRead refactor for conn
func (c *knettyConn) OnRead() error {
	if _, err := c.inputBuffer.CopyFromFd(c.fd); err != nil {
		return err
	}

	waitBufferSize := c.waitBufferSize.Load()
	if waitBufferSize > 0 && int64(c.inputBuffer.Len()) > waitBufferSize {
		c.waitBufferChan <- struct{}{}
	}
	return nil
}

// OnWrite refactor for conn
func (c *knettyConn) OnWrite() error {
	if _, err := c.outputBuffer.WriteToFd(c.fd); err != nil && err != syscall.EAGAIN {
		return err
	}

	if c.outputBuffer.IsEmpty() {
		if err := c.Register(poll.RwToRead); err != nil {
			return err
		}

		c.writeNetBufferChan <- struct{}{}
	}
	return nil
}

// OnInterrupt refactor for conn
func (c *knettyConn) OnInterrupt() error {
	c.close.Store(1)
	c.closeWaitBufferCh()
	if err := c.poller.Register(&poll.NetFileDesc{
		FD: c.fd,
	}, poll.DeleteRead); err != nil {
		return err
	}

	if err := c.closeCallBackFn; err != nil {
		c.closeCallBackFn()
	}
	return nil
}
