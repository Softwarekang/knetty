package main

import (
	"github.com/Softwarekang/knet"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatal(err)
	}

	file, err := listener.(*net.TCPListener).File()
	if err != nil {
		log.Fatal(err)
	}

	listenerFD := int(file.Fd())
	onRead := func() error {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
		conn.Close()
		return nil
	}

	poller := knet.NewDefaultPoller()
	if err = poller.Register(&knet.NetFileDesc{
		FD: listenerFD,
		NetPollListener: knet.NetPollListener{
			OnRead: onRead,
		},
	}, knet.Read); err != nil {
		log.Fatal(err)
	}

	poller.Wait()
}
