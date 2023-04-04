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

import (
	"time"
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
type CloseCallBackFunc func()

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
