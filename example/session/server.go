package session

import (
	"fmt"
	"log"
	"net"

	"github.com/Softwarekang/knet/net/connection"
	kpoll "github.com/Softwarekang/knet/net/poll"
	"github.com/Softwarekang/knet/session"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatal(err)
	}

	kpoll.PollerManager.SetPollerNums(1)
	poller := kpoll.PollerManager.Pick()
	onRead := func() error {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			return err
		}

		fmt.Printf("server %s get accept new client conn:%v \n", conn.LocalAddr().String(), conn.RemoteAddr().String())
		wrappedConn, err := connection.NewWrappedConn(conn)
		if err != nil {
			log.Fatal(err)
			return err
		}

		tcpConn := connection.NewTcpConn(wrappedConn)
		if err := tcpConn.Register(kpoll.Read); err != nil {
			return err
		}

		if err := session.NewSession(tcpConn).Run(); err != nil {
			log.Fatal(err)
		}
		return nil
	}

	file, err := listener.(*net.TCPListener).File()
	if err != nil {
		log.Fatal(err)
	}
	if err = poller.Register(&kpoll.NetFileDesc{
		FD: int(file.Fd()),
		NetPollListener: kpoll.NetPollListener{
			OnRead: onRead,
		},
	}, kpoll.Read); err != nil {
		log.Fatal(err)
	}

	// block
	poller.Wait()
}
