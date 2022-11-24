package knet

import (
	"github.com/Softwarekang/knet/connection"
)

// Session client、server session
type Session interface {
	SetPkgCodec(codec PkgCodec)
	SetEventListener(eventListener EventListener)
	WritePkg(pgk interface{}) error
	Close() error
}

type session struct {
	connection.Connection

	pkgCodec PkgCodec

	eventListener EventListener
}

func NewSession(c connection.Connection) *session {
	return &session{
		Connection: c,
	}
}

func (s *session) SetPkgCodec(codec PkgCodec) {
	if codec == nil {
		panic("session.SetPkgCodec codec is nil")
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
	data, err := s.pkgCodec.Write(pkg)
	if err != nil {
		return err
	}

	if err := s.Connection.WriteBuffer(data); err != nil {
		return err
	}
	return nil
}

func (s *session) Close() error {
	s.Connection.Close()
	return nil
}
