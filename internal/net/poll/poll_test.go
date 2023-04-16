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

package poll

import (
	"fmt"
	"log"
	"net"
	"syscall"
	"testing"
	"time"
)

func TestPoll(t *testing.T) {
	poller := NewDefaultPoller()
	// start server
	ln, err := net.Listen("tcp", "127.0.0.1:8000")
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

			tcpConn := conn.(*net.TCPConn)
			file, err := tcpConn.File()
			if err != nil {
				log.Fatalln(err)
			}
			if err = poller.Register(&NetFileDesc{
				FD: int(file.Fd()),
				NetPollListener: NetPollListener{
					OnRead: func() error {
						buf := make([]byte, 5)
						n, err := syscall.Read(int(file.Fd()), buf)
						if err != nil {
							return err
						}

						if n != 4 && string(buf) != "hello" {
							log.Fatalln("read pkg illegal")
						}
						fmt.Printf("sever got data:%s\n", "hello")
						return nil
					}, OnInterrupt: func() error {
						defer func() {
							fmt.Printf("sever got fd:%d closed\n", int(file.Fd()))
						}()
						return poller.Register(&NetFileDesc{
							FD: int(file.Fd()),
						}, DeleteRead)
					},
				},
			}, Read); err != nil {
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
		network, address := "tcp", "127.0.0.1:8000"
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
		if err := conn.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	time.Sleep(3 * time.Second)
	if err := poller.Close(); err != nil {
		log.Fatalln(err)
	}
}
