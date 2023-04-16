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

package net

import (
	"fmt"
	"net"

	"github.com/Softwarekang/knetty/internal/net/connection"
	"github.com/Softwarekang/knetty/internal/net/listener"
	errors "github.com/Softwarekang/knetty/pkg/err"
	netutil "github.com/Softwarekang/knetty/pkg/net"

	"golang.org/x/sys/unix"
)

func Listen(network, address string) (listener.Listener, error) {
	switch network {
	case "tcp":
		return listenTcp(network, address)
	default:
		return nil, errors.UnKnowNetworkErr(network)
	}

}

func listenTcp(network, address string) (*listener.TcpListener, error) {
	tcpAddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return nil, err
	}

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}

	if err := unix.Bind(fd, &unix.SockaddrInet4{
		Port: tcpAddr.Port,
		Addr: tcpAddr.AddrPort().Addr().As4(),
	}); err != nil {
		return nil, err
	}

	if err := unix.Listen(fd, unix.SOMAXCONN); err != nil {
		return nil, err
	}
	return &listener.TcpListener{
		Fd:      fd,
		TcpAddr: tcpAddr,
	}, unix.SetNonblock(fd, true)
}

func Dial(network, address string) (connection.Connection, error) {
	switch network {
	case "tcp":
		return dialTcp(network, address)
	default:
		return nil, errors.UnKnowNetworkErr(network)
	}

}

func dialTcp(network string, address string) (*connection.TcpConn, error) {
	tcpAddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return nil, err
	}

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}

	rsa, err := netutil.ResolveNetAddrToSocketAddr(tcpAddr)
	if err != nil {
		return nil, err
	}
	fmt.Println(rsa)
	if err = unix.Connect(fd, rsa); err != nil {
		return nil, err
	}

	lsa, err := unix.Getsockname(fd)
	if err != nil {
		return nil, err
	}
	return connection.NewTcpConn(fd, netutil.SocketAddrToAddr(lsa), netutil.SocketAddrToAddr(rsa)), unix.SetNonblock(fd, true)
}
