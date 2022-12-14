package knetty

import "github.com/Softwarekang/knetty/session"

/*
NewSessionCallBackFunc It is executed when a new session is established,
so some necessary parameters for drawing need to be set to ensure that the session starts properly.
*/
type NewSessionCallBackFunc func(s session.Session) error

// ServerOption option for server
type ServerOption func(*ServerOptions)

// ServerOptions options for server
type ServerOptions struct {
	network    string
	address    string
	newSession NewSessionCallBackFunc
}

// withServerNetwork set network
func withServerNetwork(network string) ServerOption {
	return func(opt *ServerOptions) {
		opt.network = network
	}
}

// withServerAddress set address
func withServerAddress(address string) ServerOption {
	return func(opt *ServerOptions) {
		opt.address = address
	}
}

// WithServiceNewSessionCallBackFunc set newSessionCallBackFunc
func WithServiceNewSessionCallBackFunc(f NewSessionCallBackFunc) ServerOption {
	return func(opt *ServerOptions) {
		opt.newSession = f
	}
}

// ClientOption option for client
type ClientOption func(options *ClientOptions)

// ClientOptions options for client
type ClientOptions struct {
	network    string
	address    string
	newSession NewSessionCallBackFunc
}

// withClientNetwork set network
func withClientNetwork(network string) ClientOption {
	return func(opt *ClientOptions) {
		opt.network = network
	}
}

// withClientAddress set address
func withClientAddress(address string) ClientOption {
	return func(opt *ClientOptions) {
		opt.address = address
	}
}

// WithClientNewSessionCallBackFunc set newSessionCallBackFunc
func WithClientNewSessionCallBackFunc(f NewSessionCallBackFunc) ClientOption {
	return func(opt *ClientOptions) {
		opt.newSession = f
	}
}
