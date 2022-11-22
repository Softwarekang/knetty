package knet

import "github.com/Softwarekang/knet/session"

// PkgCodec codec for session
type PkgCodec interface {
	Write(pkg interface{}) ([]byte, error)
	Read([]byte) (interface{}, error)
}

// EventListener listen for session
type EventListener interface {
	OnMessage(pkg interface{}) error
	OnConnect(session session.Session) error
}
