package knetty

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/netip"

	"github.com/Softwarekang/knetty/net/connection"
	"github.com/Softwarekang/knetty/net/poll"
	"github.com/Softwarekang/knetty/session"
)

// Server for knetty
type Server struct {
	ServerOptions

	tcpListener net.Listener
	poller      poll.Poll
	close       chan struct{}
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
		close:  make(chan struct{}),
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
	file, err := streamListener.(*net.TCPListener).File()
	if err != nil {
		return err
	}

	if err := s.poller.Register(&poll.NetFileDesc{
		FD: int(file.Fd()),
		NetPollListener: poll.NetPollListener{
			OnRead: s.onRead,
		},
	}, poll.Read); err != nil {
		return err
	}

	s.waitQuit()
	return nil
}

func (s *Server) onRead() error {
	netConn, err := s.tcpListener.Accept()
	if err != nil {
		return err
	}

	tcpConn, err := connection.NewTcpConn(netConn)
	if err != nil {
		return err
	}

	if err := tcpConn.Register(poll.Read); err != nil {
		return err
	}

	newSession := session.NewSession(tcpConn)
	if err := s.newSession(newSession); err != nil {
		return err
	}

	go func() {
		if err := newSession.Run(); err != nil {
			log.Println(err)
		}
	}()

	return nil
}

func (s *Server) waitQuit() {
	<-s.close
}

// Shutdown stop server
func (s *Server) Shutdown(ctx context.Context) error {
	// todo:pref(shutdown)
	if err := s.tcpListener.Close(); err != nil {
		return err
	}

	return s.poller.Close()
}
