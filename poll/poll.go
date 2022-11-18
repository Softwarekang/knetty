package poll

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

// NetFileDesc file-desc for net-fd
type NetFileDesc struct {
	// FD system fd
	FD int
	// listener for poller
	NetPollListener
}

// NetPollListener listener for net poller
type NetPollListener struct {
	// OnRead will run where fd is readable
	OnRead
	// OnWrite will run where fd is writeable
	OnWrite
	// OnInterrupt will run where fd is interrupted
	OnInterrupt
}

// OnRead the callback function when the net fd state is readable
type OnRead func() error

// OnWrite The callback function when the net fd state is writable
type OnWrite func() error

// OnInterrupt The callback function when the net fd state is interrupt
type OnInterrupt func() error

// PollEventType event type for poller
type PollEventType int

const (
	Read PollEventType = iota + 1
	DeleteRead
	Write
	DeleteWrite
	OnceWrite
)
