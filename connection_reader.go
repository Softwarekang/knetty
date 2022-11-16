package knet

// Reader interface for conn
type Reader interface {
	// ReadString will return str length n
	ReadString(n int) (string, error)
	// ReadBytes will return bytes length n
	ReadBytes(n int) ([]byte, error)
	// Len will return readable data size
	Len() int
}

// ReadString .
func (t *tcpConn) ReadString(n int) (string, error) {
	bytes, err := t.ReadBytes(n)
	if err != nil {
		return "", err
	}

	return string(bytes), err
}

// ReadBytes .
func (t *tcpConn) ReadBytes(n int) ([]byte, error) {
	if t.inputBuffer.Len() > n {
		return t.read(n)
	}

	t.waitBufferSize = n
	<-t.waitBufferChan
	t.waitBufferSize = 0
	return t.read(n)
}

func (t *tcpConn) read(n int) ([]byte, error) {
	data := make([]byte, 0, n)
	_, err := t.inputBuffer.Read(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Len .
func (t *tcpConn) Len() int {
	return t.inputBuffer.Len()
}
