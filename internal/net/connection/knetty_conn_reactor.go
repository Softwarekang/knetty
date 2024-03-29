/*
	Copyright 2022 Phoenix

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
	"github.com/Softwarekang/knetty/internal/net/poll"

	"golang.org/x/sys/unix"
)

// OnRead executed when the network connection FD is readable.
// the network data first enters the connection buffer as much as possible,
// and then drives the EventTrigger OnConnReadable function of the upper layer to process the data in the buffer.
func (c *knettyConn) OnRead() (err error) {
	if _, err = c.inputBuffer.CopyFromFd(c.fd); err != nil {
		return
	}

	// return a copied buf for session
	buf := c.inputBuffer.Bytes()
	usedBufLen := c.eventTrigger.OnConnReadable(buf)
	c.inputBuffer.Release(usedBufLen)
	return
}

// OnWrite executed when the network connection FD is writeable.
// in some cases, there may be an `abnormality (EAGAIN)` in which data is written to the network.
// When the network FD becomes writable, data should be written to the network as much as possible.
func (c *knettyConn) OnWrite() (err error) {
	if _, err = c.outputBuffer.WriteToFd(c.fd); err != nil {
		return err
	}

	if c.outputBuffer.IsEmpty() {
		c.writeable = true
		// unregister the connection FD readable event to avoid too many invalid readable event triggers by poll.
		return c.Register(poll.RwToRead)
	}
	return
}

// OnInterrupt executed when the network connection FD is close/hup.
// when the network connection needs to be closed or the exception needs to close the entire connection.
func (c *knettyConn) OnInterrupt() error {
	// set connection status
	c.close.Store(1)
	// trigger OnConnHup fn
	c.eventTrigger.OnConnHup()
	// clean up the connection FD in poll to avoid resource leaks
	if err := c.poller.Register(&poll.NetFileDesc{
		FD: c.fd,
	}, poll.DeleteRead); err != nil {
		return err
	}

	return unix.Close(c.fd)
}
