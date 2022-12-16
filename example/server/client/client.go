package main

import (
	"errors"
	"fmt"
	"log"
	"time"

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
	s.SetReadTimeout(1 * time.Second)
	s.SetWriteTimeout(1 * time.Second)
	s.SetCodec(codec{})
	s.SetEventListener(&pkgListener{})
	return nil
}

func sendHello(s session.Session) {
	if err := s.WritePkg("hello"); err != nil {
		log.Fatalln(err)
	}

	if err := s.FlushBuffer(); err != nil {
		log.Fatalln(err)
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

func (e *pkgListener) OnMessage(s session.Session, pkg interface{}) {
	data := pkg.(string)
	fmt.Println(data)
}

func (e *pkgListener) OnConnect(s session.Session) {
	fmt.Printf("local:%s get a remote:%s connection\n", s.LocalAddr(), s.RemoteAddr())
	sendHello(s)
}

func (e *pkgListener) OnClose(s session.Session) {
	fmt.Printf("session close\n")
}

func (e *pkgListener) OnError(s session.Session, err error) {
	fmt.Printf("session got err :%v\n", err)
}
