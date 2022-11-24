package buffer

// Buffer net in out data buffer
type Buffer interface {
	// Write bytes to buffer
	Write(bytes []byte) error
	// Read buffer data to  bytes
	Read(bytes []byte) (int, error)
	// Bytes will return buffer bytes
	Bytes() []byte
	// Len will return buffer readable length
	Len() int
	// WriteString string to buffer
	WriteString(s string) error
	// IsEmpty will return true if buffer len is zero
	IsEmpty() bool
	// Release will release length n buffer data
	Release(n int)
	// Clear will clear buffer
	Clear()
}
