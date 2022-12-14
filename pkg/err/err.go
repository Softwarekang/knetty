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
	// ConnClosedErr conn interrupted err
	ConnClosedErr = &connCloseErr{}
)

type netTimeoutErr struct {
}

func (n netTimeoutErr) Error() string {
	return "net io timeout"
}

func (n netTimeoutErr) TimeoutError() bool {
	return true
}

type connCloseErr struct {
}

func (c *connCloseErr) Error() string {
	return "net conn interrupted"
}

func (c *connCloseErr) TimeoutError() bool {
	return false
}
