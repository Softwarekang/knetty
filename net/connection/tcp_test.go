/*
	Copyright 2022 Phoenix

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package connection

import (
	"fmt"
	"log"
	"net"
	"syscall"
	"testing"
	"time"

	"github.com/Softwarekang/knetty/net/poll"
)

func TestTcpConnection(t *testing.T) {
	poller := poll.NewDefaultPoller()
	// start server
	ln, err := net.Listen("tcp", "127.0.0.1:8001")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}

			tcpConn, err := NewTcpConn(conn)
			if err != nil {
				log.Fatalln(err)
			}

			if err = poller.Register(&poll.NetFileDesc{
				FD: tcpConn.fd,
				NetPollListener: poll.NetPollListener{
					OnRead: func() error {
						buf := make([]byte, 5)
						n, err := syscall.Read(tcpConn.fd, buf)
						if err != nil {
							return err
						}

						if n != 4 && string(buf) != "hello" {
							log.Fatalln("read pkg illegal")
						}
						fmt.Printf("server got data:%s\n", "hello")
						return nil
					}, OnInterrupt: func() error {
						defer func() {
							fmt.Printf("server got fd:%d closed\n", tcpConn.fd)
						}()

						return poller.Register(&poll.NetFileDesc{
							FD: tcpConn.fd,
						}, poll.DeleteRead)
					},
				},
			}, poll.Read); err != nil {
				log.Fatalln(err)
			}

			if err := poller.Wait(); err != nil {
				log.Fatalln(err)
			}
		}
	}()

	go func() {
		defer func() {
			fmt.Printf("client exit\n")
		}()
		network, address := "tcp", "127.0.0.1:8001"
		conn, err := net.Dial(network, address)
		if err != nil {
			log.Fatal(err)
		}

		n, err := conn.Write([]byte("hello"))
		if err != nil && n != 5 {
			log.Fatal("write pkg illegal")
		}
		fmt.Printf("client write data:%s\n", "hello")

		time.Sleep(1 * time.Second)

		tcpConn, err := NewTcpConn(conn)
		if err != nil {
			log.Fatalln(err)
		}
		if err := tcpConn.Close(); err != nil {
			log.Fatalln(err)
		}
		time.Sleep(10 * time.Second)
	}()

	time.Sleep(3 * time.Second)
}
