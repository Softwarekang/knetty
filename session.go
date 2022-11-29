package knet

import (
	"errors"
	"go.uber.org/atomic"

	"github.com/Softwarekang/knet/connection"
	"github.com/Softwarekang/knet/pkg/buffer"
)

const (
	onceReadBufferSize = 1024
)

// Session client、server session
type Session interface {
	SetCodec(codec Codec)
	SetEventListener(eventListener EventListener)
	WritePkg(pkg interface{}) error
	Close() error
	Run() error
}

type session struct {
	connection.Connection

	pkgCodec      Codec
	eventListener EventListener
	close         atomic.Int32
}

func NewSession(c connection.Connection) *session {
	return &session{
		Connection: c,
	}
}

func (s *session) SetCodec(codec Codec) {
	if codec == nil {
		panic("session.SetCodec codec is nil")
	}
	s.pkgCodec = codec
}

func (s *session) SetEventListener(eventListener EventListener) {
	if eventListener == nil {
		panic("session.SetEventListener eventListener is nil")
	}
	s.eventListener = eventListener
}

func (s *session) WritePkg(pkg interface{}) error {
	data, err := s.pkgCodec.Encode(pkg)
	if err != nil {
		return err
	}

	if err := s.Connection.WriteBuffer(data); err != nil {
		return err
	}
	return nil
}

func (s *session) Run() error {
	if s.pkgCodec == nil {
		return errors.New("session pkgCodec is nil")
	}

	if s.eventListener == nil {
		return errors.New("session eventListener is nil")
	}

	if err := s.eventListener.OnConnect(s); err != nil {
		return err
	}

	return s.handlePkg()
}

func (s *session) handlePkg() error {
	var err error
	defer func() {
		if err != nil {
			s.eventListener.OnError(s)
		}
	}()
	if s.Connection == nil {
		err = errors.New("session connection is nil")
		return err
	}

	switch s.Connection.Type() {
	case connection.TCPCONNECTION:
		if err = s.handleTcpPkg(); err != nil {
			return nil
		}
	default:
		err = errors.New("session unSupport connection type")
		return err
	}

	return nil
}

func (s *session) handleTcpPkg() error {
	buf := buffer.NewByteBuffer()
	for {
		if !s.isActive() {
			return nil
		}

		p := make([]byte, onceReadBufferSize)
		if _, err := s.Connection.Read(p); err != nil {
			return err
		}

		if err := buf.Write(p); err != nil {
			return err
		}

		pkg, pkgLen, err := s.pkgCodec.Decode(buf.Bytes())
		if err != nil {
			return err
		}

		if pkg == nil {
			continue
		}

		buf.Release(pkgLen)

		s.eventListener.OnMessage(s, pkg)
	}
}

func (s *session) isActive() bool {
	return s.close.Load() == 0
}
func (s *session) Close() error {
	s.eventListener.OnClose(s)
	s.close.Store(1)
	s.Connection.Close()
	return nil
}
