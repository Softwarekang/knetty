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

// Package connection  implements tcp, udp and other protocols for network connection
package connection

type ConnType int

// define tcp、upd、webSocket conn
const (
	TCPCONNECTION ConnType = iota
	UDPCONNECTION
	WEBSOCKETCONNECTION
)

// CloseCallBackFunc will runs at conn on Interrupt
type CloseCallBackFunc func()

// Connection some connection  operations
type Connection interface {
	// ID for conn
	ID() uint32
	// LocalAddr local address for conn
	LocalAddr() string
	// RemoteAddr remote address for conn
	RemoteAddr() string
	// WriteBuffer will write bytes to conn buffer
	WriteBuffer(bytes []byte) (int, error)
	// FlushBuffer will send conn buffer data to net
	FlushBuffer() error
	// SetEventTrigger set eventTrigger for conn
	SetEventTrigger(trigger EventTrigger)
	// Len will return conn readable data size
	Len() int
	// Type  will return conn type
	Type() ConnType
	// Close will interrupt conn
	Close() error
}

// EventTrigger trigger
type EventTrigger interface {
	OnConnBufferReadable([]byte) int
	OnConnHup()
}
