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

// Package net general net func
package net

import (
	"errors"
	"golang.org/x/sys/unix"
	"net"
)

// ResolveConnFileDesc  get the real file descriptor of net.conn.
func ResolveConnFileDesc(conn net.Conn) (int, error) {
	if conn == nil {
		return 0, errors.New("conn is nil")
	}

	var fdCopy uintptr
	switch c := conn.(type) {
	// When a TCP/UDP connection returns a file, it is a copy of the underlying OS file, and closing this file with
	// syscall.close(file.FD()) has no effect on closing the connection.
	// Thus, it is necessary to use the SyscallConn method to obtain the actual connection and invoke
	// Control(f func(fd uintptr)) error to obtain the actual FD for effective connection management.
	case *net.TCPConn:
		rawConn, err := c.SyscallConn()
		if err != nil {
			return 0, nil
		}
		if err := rawConn.Control(func(fd uintptr) {
			fdCopy = fd
		}); err != nil {
			return 0, err
		}
		return int(fdCopy), nil
	case *net.UDPConn:
		rawConn, err := c.SyscallConn()
		if err != nil {
			return 0, nil
		}
		if err := rawConn.Control(func(fd uintptr) {
			fdCopy = fd
		}); err != nil {
			return 0, err
		}
		return int(fdCopy), nil
	default:
		return 0, errors.New("resolveConnFileDesc only support tcp、udp")
	}
}

// ResolveNetAddrToSocketAddr resolve net addr to socket addr
func ResolveNetAddrToSocketAddr(netAddr net.Addr) (unix.Sockaddr, error) {
	if netAddr == nil {
		return nil, errors.New("netAddr is nil")
	}

	switch addr := netAddr.(type) {
	case *net.TCPAddr:
		return convertAddrToSocketAddr(addr.IP, addr.Port)
	case *net.UDPAddr:
		return convertAddrToSocketAddr(addr.IP, addr.Port)
	default:
		return nil, errors.New("ResolveNetAddrToSocketAddr only support tcp、udp addr")
	}
}

func convertAddrToSocketAddr(ip net.IP, port int) (unix.Sockaddr, error) {
	parsedIP := net.ParseIP(ip.String())
	if parsedIP == nil {
		return nil, errors.New("ip is illegal")
	}
	if ip.To4() != nil {
		return iPToSockAddrInet4(ip, port)
	}

	return ipToSockaddrInet6(ip, port)
}

// iPToSockAddrInet4 convert ip port to  unix.SockaddrInet4
func iPToSockAddrInet4(ip net.IP, port int) (*unix.SockaddrInet4, error) {
	if len(ip) == 0 {
		ip = net.IPv4zero
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return nil, &net.AddrError{Err: "non-IPv4 address", Addr: ip.String()}
	}
	sa := &unix.SockaddrInet4{Port: port}
	copy(sa.Addr[:], ip4)
	return sa, nil
}

// ipToSockaddrInet6 convert ip port to unix.SockaddrInet6
func ipToSockaddrInet6(ip net.IP, port int) (*unix.SockaddrInet6, error) {
	// In general, an IP wildcard address, which is either
	// "0.0.0.0" or "::", means the entire IP addressing
	// space. For some historical reason, it is used to
	// specify "any available address" on some operations
	// of IP node.
	//
	// When the IP node supports IPv4-mapped IPv6 address,
	// we allow a listener to listen to the wildcard
	// address of both IP addressing spaces by specifying
	// IPv6 wildcard address.
	if len(ip) == 0 || ip.Equal(net.IPv4zero) {
		ip = net.IPv6zero
	}
	// We accept any IPv6 address including IPv4-mapped
	// IPv6 address.
	ip6 := ip.To16()
	if ip6 == nil {
		return &unix.SockaddrInet6{}, &net.AddrError{Err: "non-IPv6 address", Addr: ip.String()}
	}
	sa := &unix.SockaddrInet6{Port: port}
	copy(sa.Addr[:], ip6)
	return sa, nil
}

// SocketAddrToAddr returns a go/net friendly address
func SocketAddrToAddr(sa unix.Sockaddr) net.Addr {
	var a net.Addr
	switch sa := sa.(type) {
	case *unix.SockaddrInet4:
		a = &net.TCPAddr{
			IP:   sa.Addr[0:],
			Port: sa.Port,
		}
	case *unix.SockaddrInet6:
		var zone string
		if sa.ZoneId != 0 {
			if ifi, err := net.InterfaceByIndex(int(sa.ZoneId)); err == nil {
				zone = ifi.Name
			}
		}
		// if zone == "" && sa.ZoneId != 0 {
		// }
		a = &net.TCPAddr{
			IP:   sa.Addr[0:],
			Port: sa.Port,
			Zone: zone,
		}
	case *unix.SockaddrUnix:
		a = &net.UnixAddr{Net: "unix", Name: sa.Name}
	}
	return a
}
