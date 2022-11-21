package connection

import (
	"context"
	"fmt"
	"net"
	"syscall"
	"time"

	"github.com/Softwarekang/knet"
	"go.uber.org/atomic"
)

type tcpConn struct {
	kNetConn
	conn net.Conn
}

func NewTcpConn(conn Conn) *tcpConn {
	if conn == nil {
		panic("newTcpConn(conn net.Conn):@conn is nil")
	}

	var localAddress, remoteAddress string
	if conn.LocalAddr() != nil {
		localAddress = conn.LocalAddr().String()
	}

	if conn.RemoteAddr() != nil {
		remoteAddress = conn.RemoteAddr().String()
	}

	// set conn no block
	syscall.SetNonblock(conn.FD(), true)
	return &tcpConn{
		kNetConn: kNetConn{
			fd:               conn.FD(),
			remoteSocketAddr: conn.RemoteSocketAddr(),
			readTimeOut:      atomic.NewDuration(netIOTimeout),
			writeTimeOut:     atomic.NewDuration(netIOTimeout),
			localAddress:     localAddress,
			remoteAddress:    remoteAddress,
			poller:           knet.PollerManager.Pick(),
			waitBufferChan:   make(chan struct{}, 1),
		},
		conn: conn,
	}
}

func (t tcpConn) ID() uint32 {
	return t.id
}

func (t tcpConn) LocalAddr() string {
	return t.localAddress
}

func (t tcpConn) RemoteAddr() string {
	return t.remoteAddress
}

func (t tcpConn) ReadTimeout() time.Duration {
	return t.readTimeOut.Load()
}

func (t tcpConn) SetReadTimeout(rTimeout time.Duration) {
	if rTimeout < 1 {
		panic("SetReadTimeout(rTimeout time.Duration):@rTimeout < 0")
	}
	t.readTimeOut = atomic.NewDuration(rTimeout)
}

func (t tcpConn) WriteTimeout() time.Duration {
	return t.writeTimeOut.Load()
}

func (t tcpConn) SetWriteTimeout(wTimeout time.Duration) {
	if wTimeout < 1 {
		panic("SetWriteTimeout(wTimeout time.Duration):@wTimeout < 0")
	}

	t.writeTimeOut = atomic.NewDuration(wTimeout)
}

// Read .
func (t *tcpConn) Read(n int) ([]byte, error) {
	if err := t.waitReadBuffer(n); err != nil {
		return nil, err
	}

	return t.read(n)
}

func (t *tcpConn) waitReadBuffer(n int) error {
	if t.inputBuffer.Len() >= n {
		return nil
	}

	t.waitBufferSize.Store(int64(n))
	defer t.waitBufferSize.Store(0)
	if t.inputBuffer.Len() >= n {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.TODO(), t.readTimeOut.Load())
	defer cancel()
	for t.inputBuffer.Len() < n {
		if !t.isActive() {
			return fmt.Errorf("waitReadBufferWithTimeout conn is closed")
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("waitReadBufferWithTimeout ctx timeout")
		case <-t.waitBufferChan:
			continue
		}
	}

	return nil
}

func (t *tcpConn) read(n int) ([]byte, error) {
	data := make([]byte, n)
	n, err := t.inputBuffer.Read(data)
	if err != nil {
		return nil, err
	}

	fmt.Printf("read %d length data from input buffer", n)
	return data, nil
}

// Write .
func (t *tcpConn) Write(bytes []byte) (int, error) {
	return syscall.SendmsgN(t.fd, bytes, nil, t.remoteSocketAddr, 0)
}

func ipToSockaddrInet4(ip net.IP, port int) (*syscall.SockaddrInet4, error) {
	if len(ip) == 0 {
		ip = net.IPv4zero
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return nil, &net.AddrError{Err: "non-IPv4 address", Addr: ip.String()}
	}
	sa := &syscall.SockaddrInet4{Port: port}
	copy(sa.Addr[:], ip4)
	return sa, nil
}

// Len .
func (t *tcpConn) Len() int {
	return t.inputBuffer.Len()
}

func (t *tcpConn) isActive() bool {
	return t.close.Load() == 0
}

// Close .
func (t tcpConn) Close() {
	t.OnInterrupt()
}
