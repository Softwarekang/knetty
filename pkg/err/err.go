// Package err wrapped err for knetty
package err

// knettyErr wrapped err for net
type knettyErr interface {
	error

	// TimeoutError if an error is caused by a timeout will return true
	TimeoutError() bool
}

var (
	// NetIOTimeoutErr net io err
	NetIOTimeoutErr = &netTimeoutErr{}
	// ConnClosedErr conn closed err
	ConnClosedErr = &connClosedErr{}
	// ClientClosedErr client closed err
	ClientClosedErr = &clientClosedErr{}
	// ServerClosedErr server closed err
	ServerClosedErr = &serverClosedErr{}
)

type netTimeoutErr struct{}

func (n netTimeoutErr) Error() string {
	return "net io timeout"
}

func (n netTimeoutErr) TimeoutError() bool {
	return true
}

type connClosedErr struct{}

func (c *connClosedErr) Error() string {
	return "net conn is closed"
}

func (c *connClosedErr) TimeoutError() bool {
	return false
}

type clientClosedErr struct{}

func (c *clientClosedErr) Error() string {
	return "client has already been closed"
}

func (c *clientClosedErr) TimeoutError() bool {
	return false
}

type serverClosedErr struct{}

func (s *serverClosedErr) Error() string {
	return "server has already been closed"
}

func (s *serverClosedErr) TimeoutError() bool {
	return false
}
