package main

import (
	"fmt"
	"log"
	"net"
	"syscall"
	"time"

	"github.com/Softwarekang/knet"
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

	poller := knet.NewDefaultPoller()
	listenerFD := int(file.Fd())
	onRead := func() error {
		nfd, stockade, err := syscall.Accept(listenerFD)
		if err != nil {
			log.Fatal(err)
		}

		if err := poller.Register(&knet.NetFileDesc{
			FD: listenerFD,
			NetPollListener: knet.NetPollListener{
				OnRead: func() error {
					buf := make([]byte, 0, 4)
					n, err := syscall.Read(listenerFD, buf)
					if err != nil {
						return err
					}

					fmt.Printf("read %d bytes, data:%s\n", n, string(buf))
					return nil
				},
			},
		}, knet.Read); err != nil {
			return err
		}
		stockadeInt4 := stockade.(*syscall.SockaddrInet4)
		tcpAddr := &net.TCPAddr{
			IP:   stockadeInt4.Addr[0:],
			Port: stockadeInt4.Port,
		}
		fmt.Printf("server  get client conn fd:%d addr:%v", nfd, tcpAddr.String())
		time.Sleep(5 * time.Second)
		// after 5 second close conn
		return syscall.Close(nfd)
	}

	if err = poller.Register(&knet.NetFileDesc{
		FD: listenerFD,
		NetPollListener: knet.NetPollListener{
			OnRead: onRead,
		},
	}, knet.Read); err != nil {
		log.Fatal(err)
	}

	// block
	poller.Wait()
}
