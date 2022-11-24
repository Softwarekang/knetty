package net

import (
	"github.com/stretchr/testify/assert"
	"net"
	"syscall"
	"testing"
)

func TestIPToSockAddrInet4(t *testing.T) {
	var (
		sa  *syscall.SockaddrInet4
		err error
	)
	sa, err = IPToSockAddrInet4(nil, 8000)
	assert.Nil(t, err)
	assert.Equal(t, sa.Port, 8000)

	sa, err = IPToSockAddrInet4(net.IP{'0', '0', '0', '0'}, 9000)
	assert.Nil(t, err)
	assert.Equal(t, sa.Port, 9000)
}
