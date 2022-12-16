package knetty

import (
	"context"
	"fmt"
	"github.com/Softwarekang/knetty/net/poll"
	"log"
	"net"

	"github.com/Softwarekang/knetty/net/connection"
	"github.com/Softwarekang/knetty/session"
)

type Client struct {
	ClientOptions

	close chan struct{}
}

// NewClient init the client
//
//	network and address are necessary parameters
//	network like tcp、udp、websocket
//	address like 127.0.0.1:8000、localhost:8000.
func NewClient(network, address string, opts ...ClientOption) *Client {
	c := &Client{
		close: make(chan struct{}),
	}
	opts = append(opts, withClientNetwork(network), withClientAddress(address))
	for _, opt := range opts {
		opt(&c.ClientOptions)
	}

	return c
}

func (c *Client) Run() error {
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
	if err := c.newSession(newSession); err != nil {
		return err
	}

	newSession.SetSessionCloseCallBack(c.quit)
	go func() {
		if err := newSession.Run(); err != nil {
			log.Println(err)
		}
	}()

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
	<-c.close
}

func (c *Client) quit() error {
	select {
	case c.close <- struct{}{}:
	default:
	}
	return nil
}

// Shutdown close the client within the maximum allowed time in ctx, otherwise return timeout err.
func (c *Client) Shutdown(ctx context.Context) error {
	// todo:fix shutdown
	_ = c.quit()
	return nil
}
