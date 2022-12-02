package main

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/Softwarekang/knet/net/connection"
	kpoll "github.com/Softwarekang/knet/net/poll"
	"github.com/Softwarekang/knet/session"
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

		wrappedConn, err := connection.NewWrappedConn(conn)
		if err != nil {
			log.Fatal(err)
			return err
		}

		tcpConn := connection.NewTcpConn(wrappedConn)
		if err := tcpConn.Register(kpoll.Read); err != nil {
			return err
		}
		newSession := session.NewSession(tcpConn)
		newSession.SetEventListener(&helloWorldListener{})
		newSession.SetCodec(&codec{})
		go func() {
			if err := newSession.Run(); err != nil {
				log.Fatal(err)
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
