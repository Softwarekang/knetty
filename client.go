package knetty

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Softwarekang/knetty/net/connection"
	"github.com/Softwarekang/knetty/net/poll"
	merr "github.com/Softwarekang/knetty/pkg/err"
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
	for _, opt := range opts {
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
