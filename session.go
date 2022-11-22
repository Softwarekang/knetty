package knet

// Session client、server session
type Session interface {
	EventListener

	PkgCodec
}
