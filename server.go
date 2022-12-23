package knetty

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/netip"
	"sync"

	"github.com/Softwarekang/knetty/net/connection"
	"github.com/Softwarekang/knetty/net/poll"
	merr "github.com/Softwarekang/knetty/pkg/err"
	"github.com/Softwarekang/knetty/session"
)

// Server for knetty
type Server struct {
	ServerOptions

	mu          sync.Mutex
	sessions    map[session.Session]struct{}
	tcpListener net.Listener
	poller      poll.Poll
	closeCh     chan struct{}
}

// NewServer init the server
// network and address are necessary parameters
// network like tcp、udp、websocket
// address like 127.0.0.1:8000、localhost:8000.
func NewServer(network, address string, opts ...ServerOption) *Server {
	s := &Server{
		poller:   poll.PollerManager.Pick(),
		sessions: make(map[session.Session]struct{}),
		closeCh:  make(chan struct{}),
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

	fmt.Printf("sever started listen on: [%s]....", s.address)
	s.waitQuit()
	return nil
}

func (s *Server) onRead() error {
	if !s.isActive() {
		return merr.ServerClosedErr
	}

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

	newSession.SetCloseCallBackFunc(s.onSessionClose)
	s.mu.Lock()
	s.sessions[newSession] = struct{}{}
	s.mu.Unlock()
	go func() {
		if err := newSession.Run(); err != nil {
			log.Println(err)
		}
	}()

	return nil
}

func (s *Server) waitQuit() {
	<-s.closeCh
}

func (s *Server) onSessionClose(session session.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, session)
}

func (s *Server) isActive() bool {
	select {
	case <-s.closeCh:
		return false
	default:
		return true
	}
}

// Shutdown stop server
func (s *Server) Shutdown(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("server shutdown caused by:%s", ctx.Err())
		case <-s.closeCh:
			return merr.ServerClosedErr
		default:
			s.closeServerCloseCh()
			if s.tcpListener != nil {
				if err := s.tcpListener.Close(); err != nil {
					log.Printf("tcpListener closeCh err caused by:%s", err.Error())
				}
			}

			s.mu.Lock()
			for ss := range s.sessions {
				if err := ss.Close(); err != nil {
					log.Printf("session closeCh err caused by:%s", err.Error())
				}
			}
			s.mu.Unlock()
			return s.poller.Close()
		}
	}
}

func (s *Server) closeServerCloseCh() {
	select {
	case <-s.closeCh:
	default:
		close(s.closeCh)
	}
}
