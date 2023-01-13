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

package net

import (
	"errors"
	"net"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPToSockAddrInet4(t *testing.T) {
	var (
		sa  *syscall.SockaddrInet4
		err error
	)
	sa, err = iPToSockAddrInet4(nil, 8000)
	assert.Nil(t, err)
	assert.Equal(t, sa.Port, 8000)

	sa, err = iPToSockAddrInet4(net.IP{'0', '0', '0', '0'}, 9000)
	assert.Nil(t, err)
	assert.Equal(t, sa.Port, 9000)
}

func TestResolveConnFileDesc(t *testing.T) {
	testCases := []struct {
		name          string
		conn          net.Conn
		expectedFd    int
		expectedError error
	}{
		{
			name:          "nil conn",
			conn:          nil,
			expectedError: errors.New("conn is nil"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fd, err := ResolveConnFileDesc(tc.conn)
			if err != nil {
				if tc.expectedError == nil {
					t.Errorf("unexpected error: %v", err)
				} else if err.Error() != tc.expectedError.Error() {
					t.Errorf("incorrect error message: got %q, want %q", err.Error(), tc.expectedError.Error())
				}
			} else {
				if tc.expectedError != nil {
					t.Errorf("expected error %q but got none", tc.expectedError.Error())
				} else if fd != tc.expectedFd {
					t.Errorf("incorrect fd: got %d, want %d", fd, tc.expectedFd)
				}
			}
		})
	}
}

func TestResolveNetAddrToSocketAddr(t *testing.T) {
	type args struct {
		netAddr net.Addr
	}
	tests := []struct {
		name    string
		args    args
		want    syscall.Sockaddr
		wantErr error
	}{
		{
			name:    "addr is nil",
			args:    args{netAddr: nil},
			want:    nil,
			wantErr: errors.New("netAddr is nil"),
		},
		{
			name: "tcp ipv4",
			args: args{
				netAddr: &net.TCPAddr{
					IP:   net.IP{'0', '0', '0', '0'},
					Port: 9000,
					Zone: "",
				},
			},
			want: &syscall.SockaddrInet4{
				Port: 9000,
				Addr: [4]byte{'0', '0', '0', '0'},
			},
		},
		{
			name: "tcp ipv6",
			args: args{
				netAddr: &net.TCPAddr{
					IP:   net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					Port: 9000,
					Zone: "",
				},
			},
			want: &syscall.SockaddrInet6{
				Port: 9000,
				Addr: [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		{
			name: "udp ipv4",
			args: args{
				netAddr: &net.UDPAddr{
					IP:   net.IP{'0', '0', '0', '0'},
					Port: 9000,
					Zone: "",
				},
			},
			want: &syscall.SockaddrInet4{
				Port: 9000,
				Addr: [4]byte{'0', '0', '0', '0'},
			},
		},
		{
			name: "udp ipv6",
			args: args{
				netAddr: &net.UDPAddr{
					IP:   net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					Port: 9000,
					Zone: "",
				},
			},
			want: &syscall.SockaddrInet6{
				Port: 9000,
				Addr: [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveNetAddrToSocketAddr(tt.args.netAddr)
			if err != nil {
				if tt.wantErr.Error() != err.Error() {
					t.Errorf("incorrect error message: got %q, want %q", err.Error(), tt.wantErr.Error())
				}
			}
			assert.Equalf(t, tt.want, got, "ResolveNetAddrToSocketAddr(%v)", tt.args.netAddr)
		})
	}
}
