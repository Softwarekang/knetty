package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Softwarekang/knet"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatal(err)
	}

	knet.PollerManager.SetNumLoops(1)
	poller := knet.PollerManager.Pick()
	onRead := func() error {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			return err
		}

		fmt.Printf("server %s get accept new client conn:%v \n", conn.LocalAddr().String(), conn.RemoteAddr().String())
		wrappedConn, err := knet.NewWrappedConn(conn)
		if err != nil {
			log.Fatal(err)
			return err
		}

		tcpConn := knet.NewTcpConn(wrappedConn)
		if err := tcpConn.Register(); err != nil {
			return err
		}

		go func() {
			data, err := tcpConn.ReadString(5)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(data + "\n")
		}()
		return nil
	}

	file, err := listener.(*net.TCPListener).File()
	if err != nil {
		log.Fatal(err)
	}
	if err = poller.Register(&knet.NetFileDesc{
		FD: int(file.Fd()),
		NetPollListener: knet.NetPollListener{
			OnRead: onRead,
		},
	}, knet.Read); err != nil {
		log.Fatal(err)
	}

	// block
	poller.Wait()
}
