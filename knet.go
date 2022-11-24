package knet

// PkgCodec codec for session
type PkgCodec interface {
	Write(pkg interface{}) ([]byte, error)
	Read([]byte) (interface{}, error)
}

// EventListener listen for session
type EventListener interface {
	OnMessage(pkg interface{}) error
	OnConnect(s Session) error
	OnClose(s Session) error
}
