package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/Softwarekang/knet"
	"github.com/Softwarekang/knet/session"
)

func main() {
	options := []knet.ServerOption{
		knet.WithServiceNewSessionCallBackFunc(newSessionCallBackFn),
	}
	server := knet.NewServer("tcp", "127.0.0.1:8000", options...)
	if err := server.Server(); err != nil {
		log.Fatalln(err)
		return
	}
}

func newSessionCallBackFn(s session.Session) error {
	s.SetCodec(&codec{})
	s.SetEventListener(&helloWorldListener{})
	return nil
}

type helloWorldListener struct {
}

func (e *helloWorldListener) OnMessage(s session.Session, pkg interface{}) {
	data := pkg.(string)
	fmt.Println(data)
}

func (e *helloWorldListener) OnConnect(s session.Session) {
	fmt.Printf("local:%s get a remote:%s connection\n", s.LocalAddr(), s.RemoteAddr())
}

func (e *helloWorldListener) OnClose(s session.Session) {
	fmt.Printf("session close\n")
}

func (e *helloWorldListener) OnError(s session.Session, err error) {
	fmt.Printf("err :%v\n", err)
}

type codec struct {
}

func (c codec) Encode(pkg interface{}) ([]byte, error) {
	if pkg == nil {
		return nil, errors.New("pkg is illegal")
	}
	data, ok := pkg.(string)
	if !ok {
		return nil, errors.New("pkg type must be string")
	}

	if len(data) != 5 || data != "hello" {
		return nil, errors.New("pkg string must be \"hello\"")
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
	if data != "hello" {
		return nil, 0, errors.New("data is not 'hello'")
	}
	return data, len(data), nil
}
