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
	"errors"
	"net"
	"syscall"

	"github.com/Softwarekang/knetty/net/poll"
	"github.com/Softwarekang/knetty/pkg/buffer"
	netutil "github.com/Softwarekang/knetty/pkg/net"
	syscallutil "github.com/Softwarekang/knetty/pkg/syscall"
)

// TcpConn tcp connection implements the Connection interface.
type TcpConn struct {
	knettyConn

	conn net.Conn
}

// NewTcpConn create a new tcp connection, conn implements net.Conn in the standard library.
func NewTcpConn(conn net.Conn) (*TcpConn, error) {
	if conn == nil {
		return nil, errors.New("NewTcpConn(conn net.Conn): conn is nil")
	}

	var localAddress, remoteAddress string
	if conn.LocalAddr() != nil {
		localAddress = conn.LocalAddr().String()
	}

	if conn.RemoteAddr() != nil {
		remoteAddress = conn.RemoteAddr().String()
	}

	// get real fd from conn
	fd, err := netutil.ResolveConnFileDesc(conn)
	if err != nil {
		return nil, err
	}

	// set conn no block
	_ = syscallutil.SetConnectionNoBlock(fd)
	return &TcpConn{
		knettyConn: knettyConn{
			id:                 idBuilder.Inc(),
			fd:                 fd,
			localAddress:       localAddress,
			remoteAddress:      remoteAddress,
			poller:             poll.PollerManager.Pick(),
			inputBuffer:        buffer.NewRingBuffer(),
			outputBuffer:       buffer.NewRingBuffer(),
			writeNetBufferChan: make(chan struct{}, 1),
		},
		conn: conn,
	}, nil
}

// ID implements Connection.
func (t *TcpConn) ID() uint64 {
	return t.id
}

// LocalAddr implements Connection.
func (t *TcpConn) LocalAddr() string {
	return t.localAddress
}

// RemoteAddr implements Connection.
func (t *TcpConn) RemoteAddr() string {
	return t.remoteAddress
}

// WriteBuffer implements Connection.
func (t *TcpConn) WriteBuffer(bytes []byte) (int, error) {
	return t.outputBuffer.Write(bytes)
}

// FlushBuffer implements Connection.
func (t *TcpConn) FlushBuffer() error {
	if _, err := t.outputBuffer.WriteToFd(t.fd); err != nil && err != syscall.EAGAIN {
		return err
	}

	if t.outputBuffer.IsEmpty() {
		return nil
	}

	// When the network data cannot be written, register the write event to poll,
	// and write the buffer data to the network when it is writable again.
	if err := t.Register(poll.ReadToRW); err != nil {
		return err
	}

	// block the goroutine.
	<-t.writeNetBufferChan
	return nil
}

// Len implements Connection.
func (t *TcpConn) Len() int {
	return t.inputBuffer.Len()
}

func (t *TcpConn) isActive() bool {
	return t.close.Load() == 0
}

// SetEventTrigger implements Connection.
func (t *TcpConn) SetEventTrigger(trigger EventTrigger) {
	t.eventTrigger = trigger
}

// Close implements Connection.
func (t *TcpConn) Close() error {
	if !t.isActive() {
		return nil
	}
	t.close.Store(1)
	if et := t.eventTrigger; et != nil {
		et.OnConnHup()
	}
	return syscall.Close(t.fd)
}

// Type implements Connection.
func (t *TcpConn) Type() ConnType {
	return TCPCONNECTION
}
