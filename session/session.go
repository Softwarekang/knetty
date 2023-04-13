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

// Package session for knetty
package session

import (
	"errors"
	"fmt"

	"github.com/Softwarekang/knetty/net/connection"
	merr "github.com/Softwarekang/knetty/pkg/err"

	"go.uber.org/atomic"
)

type ExecStatus int

const (
	Normal ExecStatus = iota
)

// CloseCallBackFunc exec when session stopping
type CloseCallBackFunc func(Session)

// Session client、server session
type Session interface {
	// LocalAddr return local address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
	LocalAddr() string
	// RemoteAddr return remote address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
	RemoteAddr() string
	// SetCodec setting yourself codec is necessary, otherwise a panic will occur at runtime
	SetCodec(Codec)
	// SetEventListener setting yourself eventListener is necessary, otherwise a panic will occur at runtime
	SetEventListener(EventListener)
	// WritePkg will encode any type of data as a []byte type using the codec and writes it to the conn buffer.
	// If you want the other end of the network to receive it,
	// call the FlushBuffer API to send all the data from the conn buffer out
	WritePkg(pkg interface{}) (int, error)
	// WriteBuffer will write bytes to conn buffer
	WriteBuffer(bytes []byte) (int, error)
	// FlushBuffer will send conn buffer data to net
	FlushBuffer() error
	//	Run will run this session
	Run() error
	// SetCloseCallBackFunc setting closeBackFunc for session
	SetCloseCallBackFunc(fn CloseCallBackFunc)
	// Info return session info
	Info() string
	// Close will stop session
	Close() error
}

type session struct {
	conn            connection.Connection
	closeCallBackFn CloseCallBackFunc
	pkgCodec        Codec
	eventListener   EventListener
	close           atomic.Int32
}

// NewSession create new session.
func NewSession(conn connection.Connection) Session {
	s := &session{
		conn: conn,
	}

	return s
}

// LocalAddr  implements Session.
func (s *session) LocalAddr() string {
	return s.conn.LocalAddr()
}

// RemoteAddr implements Session.
func (s *session) RemoteAddr() string {
	return s.conn.RemoteAddr()
}

// SetCodec implements Session.
func (s *session) SetCodec(codec Codec) {
	if codec == nil {
		panic("codec is nil")
	}
	s.pkgCodec = codec
}

// SetEventListener implements Session.
func (s *session) SetEventListener(eventListener EventListener) {
	if eventListener == nil {
		panic("eventListener is nil")
	}
	s.eventListener = eventListener
}

// WritePkg implements Session.
func (s *session) WritePkg(pkg interface{}) (int, error) {
	data, err := s.pkgCodec.Encode(pkg)
	if err != nil {
		return 0, err
	}

	return s.conn.WriteBuffer(data)
}

// WriteBuffer implements Session.
func (s *session) WriteBuffer(data []byte) (int, error) {
	return s.conn.WriteBuffer(data)
}

// FlushBuffer implements Session.
func (s *session) FlushBuffer() error {
	return s.conn.FlushBuffer()
}

// Run implements Session.
func (s *session) Run() error {
	if s.pkgCodec == nil {
		return errors.New("session pkgCodec is nil")
	}

	if s.eventListener == nil {
		return errors.New("session eventListener is nil")
	}

	if s.conn == nil {
		return errors.New("session connection is nil")
	}

	// notify listen onConnection func
	s.eventListener.OnConnect(s)
	// set conn eventTrigger
	s.conn.SetEventTrigger(NewSessionEventTrigger(s))
	return nil
}

func (s *session) isActive() bool {
	return s.close.Load() == 0
}

// SetCloseCallBackFunc setting CloseCallBackFunc is only useful when set for the first time
func (s *session) SetCloseCallBackFunc(fn CloseCallBackFunc) {
	if s.closeCallBackFn != nil {
		return
	}
	s.closeCallBackFn = fn
}

// Info implements Session.
func (s *session) Info() string {
	return fmt.Sprintf("[localAddr:%s remoteAddr:%s]", s.LocalAddr(), s.RemoteAddr())
}

// Close implements Session.
func (s *session) Close() error {
	s.onClose()
	return s.conn.Close()
}

func (s *session) handlePkg(buf []byte) (usedBufLen int) {
	var err error
	defer func() {
		if err != nil {
			s.eventListener.OnError(s, err)
		}
	}()

	switch s.conn.Type() {
	case connection.TCPCONNECTION:
		if usedBufLen, err = s.handleTcpPkg(buf); err != nil {
			return
		}
	default:
		err = errors.New("session unSupport connection type")
		_ = s.Close()
		return
	}

	return
}

func (s *session) handleTcpPkg(buf []byte) (int, error) {
	var processedBufLen int
	for {
		if !s.isActive() {
			return processedBufLen, merr.ConnClosedErr
		}

		pkg, pkgLen, err := s.pkgCodec.Decode(buf)
		if err != nil {
			return processedBufLen, err
		}

		if pkg == nil {
			return processedBufLen, nil
		}

		processedBufLen += pkgLen
		buf = buf[pkgLen:]
		switch s.eventListener.OnMessage(s, pkg) {
		case Normal:
			continue
		}
	}
}

func (s *session) onClose() {
	if !s.isActive() {
		return
	}

	s.close.Store(1)
	s.eventListener.OnClose(s)
	if s.closeCallBackFn != nil {
		s.closeCallBackFn(s)
	}
}

type WrappedEventTrigger struct {
	session *session
}

func NewSessionEventTrigger(session *session) *WrappedEventTrigger {
	return &WrappedEventTrigger{session: session}
}
func (s WrappedEventTrigger) OnConnReadable(buf []byte) int {
	return s.session.handlePkg(buf)
}

func (s WrappedEventTrigger) OnConnHup() {
	s.session.onClose()
}
