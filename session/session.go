// Package session for knetty
package session

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Softwarekang/knetty/net/connection"
	"github.com/Softwarekang/knetty/pkg/buffer"

	"go.uber.org/atomic"
)

const (
	onceReadBufferSize = 1024
)

// CloseCallBackFunc exec when session stopping
type CloseCallBackFunc func(Session) error

// Session client、server session
type Session interface {
	// LocalAddr return local address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
	LocalAddr() string
	// RemoteAddr return remote address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
	RemoteAddr() string
	// SetReadTimeout setting read timeout
	SetReadTimeout(time.Duration)
	// SetWriteTimeout setting write timeout
	SetWriteTimeout(time.Duration)
	// SetCodec setting yourself codec is necessary, otherwise a panic will occur at runtime
	SetCodec(Codec)
	// SetEventListener setting yourself eventListener is necessary, otherwise a panic will occur at runtime
	SetEventListener(EventListener)
	// WritePkg will encode any type of data as a []byte type using the codec and writes it to the conn buffer.
	// If you want the other end of the network to receive it,
	// call the FlushBuffer API to send all the data from the conn buffer out
	WritePkg(pkg interface{}) error
	// WriteBuffer will write bytes to conn buffer
	WriteBuffer(bytes []byte) error
	// FlushBuffer will send conn buffer data to net
	FlushBuffer() error
	/*
		Run will run this session, so it is blocking.
		You can use to avoid blocking the main program
		go func(){
			if err:=session.Run();err!=nil{
				// handle err
			}
		}
	*/
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

// NewSession .
func NewSession(conn connection.Connection) Session {
	s := &session{
		conn: conn,
	}
	s.conn.SetCloseCallBack(s.onClose)
	return s
}

// LocalAddr .
func (s *session) LocalAddr() string {
	return s.conn.LocalAddr()
}

// RemoteAddr .
func (s *session) RemoteAddr() string {
	return s.conn.RemoteAddr()
}

// SetReadTimeout .
func (s *session) SetReadTimeout(duration time.Duration) {
	s.conn.SetReadTimeout(duration)
}

// SetWriteTimeout .
func (s *session) SetWriteTimeout(duration time.Duration) {
	s.conn.SetWriteTimeout(duration)
}

// SetCodec .
func (s *session) SetCodec(codec Codec) {
	if codec == nil {
		panic("codec is nil")
	}
	s.pkgCodec = codec
}

// SetEventListener .
func (s *session) SetEventListener(eventListener EventListener) {
	if eventListener == nil {
		panic("eventListener is nil")
	}
	s.eventListener = eventListener
}

// WritePkg .
func (s *session) WritePkg(pkg interface{}) error {
	data, err := s.pkgCodec.Encode(pkg)
	if err != nil {
		return err
	}

	if err := s.conn.WriteBuffer(data); err != nil {
		return err
	}
	return nil
}

// WriteBuffer .
func (s *session) WriteBuffer(data []byte) error {
	return s.conn.WriteBuffer(data)
}

// FlushBuffer .
func (s *session) FlushBuffer() error {
	return s.conn.FlushBuffer()
}

// Run .
func (s *session) Run() error {
	if s.pkgCodec == nil {
		return errors.New("session pkgCodec is nil")
	}

	if s.eventListener == nil {
		return errors.New("session eventListener is nil")
	}

	// notify listen onConnection func
	s.eventListener.OnConnect(s)
	return s.handlePkg()
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

// Info .
func (s *session) Info() string {
	return fmt.Sprintf("[localAddr:%s remoteAddr:%s]", s.LocalAddr(), s.RemoteAddr())
}

// Close .
func (s *session) Close() error {
	if err := s.onClose(); err != nil {
		return err
	}

	return s.conn.Close()
}

func (s *session) handlePkg() error {
	var err error
	defer func() {
		if s.closeCallBackFn != nil {
			if err := s.closeCallBackFn(s); err != nil {
				log.Println(err)
			}
		}

		if err != nil {
			s.eventListener.OnError(s, err)
		}
	}()
	if s.conn == nil {
		err = errors.New("session connection is nil")
		return err
	}

	switch s.conn.Type() {
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
		pkgLen, err := s.conn.Read(p)
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

func (s *session) onClose() error {
	if !s.isActive() {
		return nil
	}

	s.close.Store(1)
	s.eventListener.OnClose(s)
	return nil
}
