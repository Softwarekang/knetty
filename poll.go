package knet

// Poll net poll interface
type Poll interface {
	// Register netFd in the poller. events is the type of event that the poller focus on
	Register(netFd *NetFileDesc, eventType PollEventType) error

	// Wait
	// poller will focus on all registered netFd, wait for netFd to satisfy the condition and
	// notify the registered listener, so it is blocked
	Wait() error

	// Close the poller
	Close() error
}
