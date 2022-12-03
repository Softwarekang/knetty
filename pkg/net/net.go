// Package net general net func
package net

import (
	"net"
	"syscall"
)

// IPToSockAddrInet4 convert ip port to tcp socketAddr
func IPToSockAddrInet4(ip net.IP, port int) (*syscall.SockaddrInet4, error) {
	if len(ip) == 0 {
		ip = net.IPv4zero
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return nil, &net.AddrError{Err: "non-IPv4 address", Addr: ip.String()}
	}
	sa := &syscall.SockaddrInet4{Port: port}
	copy(sa.Addr[:], ip4)
	return sa, nil
}
