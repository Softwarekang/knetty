package knet

import (
	"fmt"
	"github.com/Softwarekang/knet/net/poll"
	"github.com/Softwarekang/knet/session"
	"go.uber.org/atomic"
	"net"
	"net/netip"
)

// Server for kNet
type Server struct {
	ServerOptions

	tcpListener net.Listener
	session     session.Session
	poller      poll.Poll
	close       atomic.Bool
}

/*
	NewServer init the server
	network and address are necessary parameters
	network like tcp、udp、websocket
	address like 127.0.0.1:8000、localhost:8000.
*/
func NewServer(network, address string, opts ...ServerOption) *Server {
	s := &Server{
		poller: poll.PollerManager.Pick(),
	}
	opts = append(opts, withServerNetwork(network), withServerAddress(address))
	for _, opt := range opts {
		opt(&s.ServerOptions)
	}

	return s
}

// Server listen and run event-loop
func (s *Server) Server() error {
	switch s.network {
	case "tcp":
		return s.tcpServer()
	default:
		return fmt.Errorf("server not support network:%v", s.network)
	}
}

func (s *Server) tcpServer() error {
	if err := s.listenTcp(); err != nil {
		return err
	}

	return nil
}

func (s *Server) listenTcp() error {
	// validate ipv4,ipv4
	address, err := netip.ParseAddrPort(s.address)
	if err != nil {
		return err
	}

	streamListener, err := net.Listen(s.network, address.String())
	if err != nil {
		return err
	}

	s.tcpListener, s.address = streamListener, streamListener.Addr().String()
	return nil
}

func (s *Server) accept() error {
	return nil
}
