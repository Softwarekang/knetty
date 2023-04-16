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

package listener

import (
	"github.com/Softwarekang/knetty/internal/net/connection"
	errors "github.com/Softwarekang/knetty/pkg/err"
	netutil "github.com/Softwarekang/knetty/pkg/net"
	"net"

	"golang.org/x/sys/unix"
)

// Listener  A Listener is a generic network listener for stream-oriented protocols.
type Listener interface {
	// Accept waits for and returns the next connection to the listener.
	Accept() (connection.Connection, error)

	// Close closes the listener.
	// Any blocked Accept operations will be unblocked and return errors.
	Close() error

	// Addr returns the listener's network address.
	Addr() net.Addr

	// FD returns the listener's fd
	FD() int
}

// TcpListener tcp network listener.
type TcpListener struct {
	Fd      int
	TcpAddr *net.TCPAddr
}

// Accept implements Listener.
func (t *TcpListener) Accept() (connection.Connection, error) {
	if !t.ok() {
		return nil, errors.IllegalListenerErr("tcp")
	}

	cfd, sa, err := unix.Accept(t.Fd)
	if err != nil {
		if err == unix.EAGAIN {
			return nil, nil
		}
		return nil, err
	}

	rsa := netutil.SocketAddrToAddr(sa)
	return connection.NewTcpConn(cfd, t.TcpAddr, rsa), unix.SetNonblock(cfd, true)
}

// Close implements Listener.
func (t *TcpListener) Close() error {
	if t.Fd != 0 {
		return unix.Close(t.Fd)
	}

	return nil
}

// Addr implements Listener.
func (t *TcpListener) Addr() net.Addr {
	return t.TcpAddr
}

// FD implements Listener.
func (t *TcpListener) FD() int {
	return t.Fd
}

func (t *TcpListener) ok() bool {
	if t.Fd != 0 && t.TcpAddr != nil {
		return true
	}

	return false
}
