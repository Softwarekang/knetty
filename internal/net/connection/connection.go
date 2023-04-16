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

// Package connection  implements tcp, udp and other protocols for network connection.
package connection

import (
	"github.com/Softwarekang/knetty/internal/net/poll"

	"go.uber.org/atomic"
)

type ConnType int

var (
	idBuilder atomic.Uint64
)

// define tcp、upd、webSocket connType
const (
	// TCPCONNECTION tcp conn
	TCPCONNECTION ConnType = iota
)

// EventTrigger define connection event notification behavior.
type EventTrigger interface {
	// OnConnReadable triggered when the connection gets data from the network.
	OnConnReadable([]byte) int
	// OnConnHup triggered when the connection gets an error / close event from the network.
	OnConnHup()
}

// Connection is a network connection oriented towards byte streams, based on an event-driven mechanism.
type Connection interface {
	// ID return a uin-type value that uniquely identifies each stream connection。
	ID() uint64
	// FD return socket fd.
	FD() int
	// LocalAddr return the actual local connection address.
	LocalAddr() string
	// RemoteAddr return the actual remote connection address.
	RemoteAddr() string
	// WriteBuffer does not immediately write the byte stream to the network,
	// but rather writes it to a local output buffer.
	WriteBuffer(bytes []byte) (int, error)
	// FlushBuffer writes all data in the local output buffer to the network.
	// It is a blocking operation and waits until all the data has been written to the network.
	FlushBuffer() error
	// SetEventTrigger set the EventTrigger, when the network connection is readable or the connection is abnormal,
	// the trigger will be driven
	SetEventTrigger(trigger EventTrigger)
	// Len returns the maximum readable bytes of the current connection.
	Len() int
	// Type Return the current connection network type tcp/udp/ws.
	Type() ConnType
	// Register register conn in poller with event.
	Register(eventType poll.EventType) error
	// Close the network connection, regardless of the ongoing blocking non-blocking read and write will return an error.
	Close() error
}
