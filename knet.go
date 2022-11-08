package knet

// OnRead the callback function when the net fd state is readable
type OnRead func() error

// OnWrite The callback function when the net fd state is writable
type OnWrite func() error

// OnInterrupt The callback function when the net fd state is interrupt
type OnInterrupt func() error
