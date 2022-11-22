package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Softwarekang/knet/connection"
	kpoll "github.com/Softwarekang/knet/poll"
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
		if err := tcpConn.Register(); err != nil {
			return err
		}

		go func() {
			for {
				data, err := tcpConn.Read(10)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(string(data) + "\n")
				n, err := tcpConn.Write(data)
				fmt.Printf("server write length:%v data:%s\n", n, string(data))
				time.Sleep(3 * time.Second)
			}

		}()
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
