// Package session for knetty
package session

import (
	"errors"
	"fmt"
	"log"

	"github.com/Softwarekang/knetty/net/connection"
	"github.com/Softwarekang/knetty/pkg/buffer"

	"go.uber.org/atomic"
)

const (
	onceReadBufferSize = 1024
)

type CloseCallBackFunc func() error

// Session client、server session
type Session interface {
	connection.Connection
	LocalAddr() string
	SetCodec(codec Codec)
	RemoteAddr() string
	SetEventListener(eventListener EventListener)
	WritePkg(pkg interface{}) error
	Run() error
	SetSessionCloseCallBack(fn CloseCallBackFunc)
	Shutdown() error
}

type session struct {
	connection.Connection

	closeCallBackFn CloseCallBackFunc
	pkgCodec        Codec
	eventListener   EventListener
	close           atomic.Int32
}

func NewSession(c connection.Connection) *session {
	s := &session{
		Connection: c,
	}
	s.Connection.SetCloseCallBack(s.onClose)
	return s
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

	s.eventListener.OnConnect(s)
	return s.handlePkg()
}

func (s *session) handlePkg() error {
	var err error
	defer func() {
		if s.closeCallBackFn != nil {
			if err := s.closeCallBackFn(); err != nil {
				log.Println(err)
			}
		}

		if err != nil {
			s.eventListener.OnError(s, err)
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
			fmt.Println("session closed")
			return nil
		}

		p := make([]byte, onceReadBufferSize)
		pkgLen, err := s.Connection.Read(p)
		if err != nil {
			return err
		}

		if err := buf.Write(p[:pkgLen]); err != nil {
			return err
		}

		for buf.Len() > 0 {
			pkg, pkgLen, err := s.pkgCodec.Decode(buf.Bytes())
			if err != nil {
				return err
			}

			if pkg == nil {
				break
			}

			s.eventListener.OnMessage(s, pkg)
			buf.Release(pkgLen)
		}
	}
}

func (s *session) isActive() bool {
	return s.close.Load() == 0
}

func (s *session) LocalAddr() string {
	return s.Connection.LocalAddr()
}

func (s *session) RemoteAddr() string {
	return s.Connection.RemoteAddr()
}

func (s *session) SetSessionCloseCallBack(fn CloseCallBackFunc) {
	s.closeCallBackFn = fn
}

func (s *session) Shutdown() error {
	s.Connection.Close()
	return s.onClose()
}

func (s *session) onClose() error {
	if !s.isActive() {
		return nil
	}

	s.eventListener.OnClose(s)
	s.close.Store(1)
	return nil
}
