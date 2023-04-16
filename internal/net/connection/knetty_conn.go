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
	"github.com/Softwarekang/knetty/pkg/buffer"

	"go.uber.org/atomic"
)

// knettyConn defines common information and methods for various network implementations.
type knettyConn struct {
	id            uint64
	fd            int
	localAddress  string
	remoteAddress string
	poller        poll.Poll
	inputBuffer   *buffer.RingBuffer
	outputBuffer  *buffer.RingBuffer
	netFd         *poll.NetFileDesc
	writeable     bool
	eventTrigger  EventTrigger
	close         atomic.Int32
}

// Register the network connection to poll.
func (c *knettyConn) Register(eventType poll.EventType) error {
	// check the connected network fd is initialized.
	c.initNetFd()
	return c.poller.Register(c.netFd, eventType)
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
			OnWrite:     c.OnWrite,
		},
	}
}
