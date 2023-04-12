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

package main

import (
	"errors"
	"fmt"
	"log"
	"syscall"

	"github.com/Softwarekang/knetty"
	"github.com/Softwarekang/knetty/session"
)

func main() {
	// setting optional options for the server
	options := []knetty.ClientOption{
		knetty.WithClientNewSessionCallBackFunc(newSessionCallBackFn),
	}
	client := knetty.NewClient("tcp", "127.0.0.1:8000", options...)

	if err := client.Run(); err != nil {
		log.Printf("run client: %s\n", err)
	}
}

// set the necessary parameters for the session to run.
func newSessionCallBackFn(s session.Session) error {
	s.SetCodec(codec{})
	s.SetEventListener(&pkgListener{})
	return nil
}

func sendHello(s session.Session) {
	n, err := s.WritePkg("hello")
	if err != nil && err != syscall.EAGAIN {
		log.Println(err)
	}

	fmt.Printf("client session send %v bytes data to server\n", n)
	if err := s.FlushBuffer(); err != nil {
		log.Println(err)
	}
}

type codec struct{}

func (c codec) Encode(pkg interface{}) ([]byte, error) {
	if pkg == nil {
		return nil, errors.New("pkg is illegal")
	}
	data, ok := pkg.(string)
	if !ok {
		return nil, errors.New("pkg type must be string")
	}

	return []byte(data), nil
}

func (c codec) Decode(bytes []byte) (interface{}, int, error) {
	if bytes == nil {
		return nil, 0, errors.New("bytes is nil")
	}

	if len(bytes) < 5 {
		return nil, 0, nil
	}

	data := string(bytes)
	if len(bytes) > 5 {
		data = data[0:5]
	}
	return data, len(data), nil
}

type pkgListener struct {
}

func (e *pkgListener) OnMessage(s session.Session, pkg interface{}) session.ExecStatus {
	data := pkg.(string)
	fmt.Printf("client got data:%s\n", data)
	return session.Normal
}

func (e *pkgListener) OnConnect(s session.Session) {
	fmt.Printf("local:%s get a remote:%s connection\n", s.LocalAddr(), s.RemoteAddr())
	sendHello(s)
}

func (e *pkgListener) OnClose(s session.Session) {
	fmt.Printf("client session: %s closed\n", s.Info())

}

func (e *pkgListener) OnError(s session.Session, err error) {
	fmt.Printf("client session: %s got err :%v\n", s.Info(), err)
}
