package knet

import "github.com/Softwarekang/knet/session"

// ServerOption option for server
type ServerOption func(*ServerOptions)

// NewSessionCallBackFunc exec when client/sever got a new session
type NewSessionCallBackFunc func(s session.Session) error

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
