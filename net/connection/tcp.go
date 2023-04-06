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
	"errors"
	"net"
	"syscall"

	"github.com/Softwarekang/knetty/net/poll"
	"github.com/Softwarekang/knetty/pkg/buffer"
	mnet "github.com/Softwarekang/knetty/pkg/net"
	msyscall "github.com/Softwarekang/knetty/pkg/syscall"
)

// TcpConn tcp conn in knetty, impl Connection
type TcpConn struct {
	knettyConn
	conn net.Conn
}

// NewTcpConn .
func NewTcpConn(conn net.Conn) (*TcpConn, error) {
	if conn == nil {
		return nil, errors.New("conn is nil")
	}

	var localAddress, remoteAddress string
	if conn.LocalAddr() != nil {
		localAddress = conn.LocalAddr().String()
	}

	if conn.RemoteAddr() != nil {
		remoteAddress = conn.RemoteAddr().String()
	}

	fd, err := mnet.ResolveConnFileDesc(conn)
	if err != nil {
		return nil, err
	}

	// set conn no block
	_ = msyscall.SetConnectionNoBlock(fd)
	return &TcpConn{
		knettyConn: knettyConn{
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

// ID .
func (t *TcpConn) ID() uint32 {
	return t.id
}

// LocalAddr .
func (t *TcpConn) LocalAddr() string {
	return t.localAddress
}

// RemoteAddr .
func (t *TcpConn) RemoteAddr() string {
	return t.remoteAddress
}

// WriteBuffer .
func (t *TcpConn) WriteBuffer(bytes []byte) (int, error) {
	return t.outputBuffer.Write(bytes)
}

// FlushBuffer .
func (t *TcpConn) FlushBuffer() error {
	if _, err := t.outputBuffer.WriteToFd(t.fd); err != nil && err != syscall.EAGAIN {
		return err
	}

	if t.outputBuffer.IsEmpty() {
		return nil
	}

	// net buffer is full
	if err := t.Register(poll.ReadToRW); err != nil {
		return err
	}

	<-t.writeNetBufferChan
	return nil
}

// Len .
func (t *TcpConn) Len() int {
	return t.inputBuffer.Len()
}

func (t *TcpConn) isActive() bool {
	return t.close.Load() == 0
}

// SetEventTrigger .
func (t *TcpConn) SetEventTrigger(trigger EventTrigger) {
	t.eventTrigger = trigger
}

// Close .
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

// Type .
func (t *TcpConn) Type() ConnType {
	return TCPCONNECTION
}
