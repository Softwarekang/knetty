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

package knetty

import (
	"context"
	"fmt"
	"net"

	"github.com/Softwarekang/knetty/net/connection"
	"github.com/Softwarekang/knetty/net/poll"
	merr "github.com/Softwarekang/knetty/pkg/err"
	"github.com/Softwarekang/knetty/pkg/log"
	"github.com/Softwarekang/knetty/session"
)

type Client struct {
	ClientOptions

	session session.Session
	closeCh chan struct{}
}

// NewClient init the client
// network and address are necessary parameters
// network like tcp、udp、websocket
// address like 127.0.0.1:8000、localhost:8000.
func NewClient(network, address string, opts ...ClientOption) *Client {
	c := &Client{
		closeCh: make(chan struct{}),
	}
	opts = append(opts, withClientNetwork(network), withClientAddress(address))
	for _, opt := range mergeCustomClientOptions(opts...) {
		opt(&c.ClientOptions)
	}

	return c
}

func (c *Client) Run() error {
	if !c.isActive() {
		return merr.ClientClosedErr
	}

	switch c.network {
	case "tcp":
		return c.tcpEventloop()
	default:
		return fmt.Errorf("client not support network:%v", c.network)
	}
}

func (c *Client) tcpEventloop() error {
	conn, err := c.dicTcp()
	if err != nil {
		return err
	}

	newSession := session.NewSession(conn)
	newSession.SetCloseCallBackFunc(c.quit)
	if err := c.newSession(newSession); err != nil {
		return err
	}

	c.session = newSession
	if err := newSession.Run(); err != nil {
		log.Errorf("session run err:%s", err.Error())
		return err
	}
	c.waitQuit()
	return nil
}

func (c *Client) dicTcp() (connection.Connection, error) {
	netConn, err := net.Dial(c.network, c.address)
	if err != nil {
		return nil, err
	}

	tcpConn, err := connection.NewTcpConn(netConn)
	if err != nil {
		return nil, err
	}

	if err := tcpConn.Register(poll.Read); err != nil {
		return nil, err
	}

	return tcpConn, nil
}

func (c *Client) waitQuit() {
	<-c.closeCh
}

func (c *Client) quit(session session.Session) {
	c.closeClientCh()
}

func (c *Client) isActive() bool {
	select {
	case <-c.closeCh:
		return false
	default:
		return true
	}
}

// Shutdown closeCh the client within the maximum allowed time in ctx, otherwise return timeout err.
func (c *Client) Shutdown(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("server shutdown caused by:%s", ctx.Err())
		case <-c.closeCh:
			return merr.ClientClosedErr
		default:
			c.quit(nil)
			if c.session != nil {
				return c.session.Close()
			}
			return nil
		}
	}
}

func (c *Client) closeClientCh() {
	select {
	case <-c.closeCh:
	default:
		close(c.closeCh)
	}
}
