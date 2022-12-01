package session

import (
	"fmt"

	"github.com/Softwarekang/knet/session"
)

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
	fmt.Printf("session close")
}

func (e *helloWorldListener) OnError(s session.Session, err error) {
	fmt.Printf("err :%v", err)
}
