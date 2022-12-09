package knet

// ServerOption option for server
type ServerOption func(*ServerOptions)

// ServerOptions options for server
type ServerOptions struct {
	network string
	address string
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
