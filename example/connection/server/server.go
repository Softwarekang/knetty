package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Softwarekang/knetty/net/connection"
	kpoll "github.com/Softwarekang/knetty/net/poll"
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

		tcpConn, err := connection.NewTcpConn(conn)
		if err != nil {
			return err
		}
		if err := tcpConn.Register(kpoll.Read); err != nil {
			return err
		}

		go func() {
			for {
				data, err := tcpConn.Next(10)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(string(data) + "\n")
				if err := tcpConn.WriteBuffer(data); err != nil {
					log.Fatal(err)
				}
				tcpConn.FlushBuffer()
				fmt.Printf("server write data:%s\n", string(data))
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
