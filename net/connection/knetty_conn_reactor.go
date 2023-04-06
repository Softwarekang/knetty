/*
	Copyright 2022 ankangan

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

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

	buf := c.inputBuffer.Bytes()
	usedBufLen := c.eventTrigger.OnConnBufferReadable(buf)
	c.inputBuffer.Release(usedBufLen)
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
	if err := c.poller.Register(&poll.NetFileDesc{
		FD: c.fd,
	}, poll.DeleteRead); err != nil {
		return err
	}

	c.eventTrigger.OnConnHup()
	return nil
}
