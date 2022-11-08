package knet

// NetFileDesc file-desc for net-fd
type NetFileDesc struct {
	FD int
	// listener for poller
	NetPollListener
}

// NetPollListener listener for net poller
type NetPollListener struct {
	OnRead
	OnWrite
	OnInterrupt
}
