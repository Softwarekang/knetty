package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Softwarekang/knetty"
	"github.com/Softwarekang/knetty/session"
)

func main() {
	// setting optional options for the server
	options := []knetty.ServerOption{
		knetty.WithServiceNewSessionCallBackFunc(newSessionCallBackFn),
	}

	// creating a new server with network settings such as tcp/upd, address such as 127.0.0.1:8000, and optional options
	server := knetty.NewServer("tcp", "127.0.0.1:8000", options...)
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := server.Server(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("run server: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server starting shutdown:", err)
	}

	log.Println("server exiting")
}

// set the necessary parameters for the session to run.
func newSessionCallBackFn(s session.Session) error {
	s.SetCodec(&codec{})
	s.SetEventListener(&helloWorldListener{})
	s.SetReadTimeout(1 * time.Second)
	s.SetWriteTimeout(1 * time.Second)
	return nil
}

type helloWorldListener struct {
}

func (e *helloWorldListener) OnMessage(s session.Session, pkg interface{}) {
	data := pkg.(string)
	fmt.Printf("server got data:%s\n", data)
}

func (e *helloWorldListener) OnConnect(s session.Session) {
	fmt.Printf("local:%s get a remote:%s connection\n", s.LocalAddr(), s.RemoteAddr())
}

func (e *helloWorldListener) OnClose(s session.Session) {
	fmt.Printf("server session: %s closed\n", s.Info())
}

func (e *helloWorldListener) OnError(s session.Session, err error) {
	fmt.Printf("session: %s got err :%v\n", s.Info(), err)
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
