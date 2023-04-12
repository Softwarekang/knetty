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
	"github.com/Softwarekang/knetty/net/poll"
	"github.com/Softwarekang/knetty/session"
)

func main() {
	poll.PollerManager.SetPollerNums(8)
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
	log.Printf("server pid:%d", os.Getpid())
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
	return nil
}

type helloWorldListener struct {
}

func (e *helloWorldListener) OnMessage(s session.Session, pkg interface{}) session.ExecStatus {
	s1 := pkg.(string)
	_, err := s.WriteBuffer([]byte(s1))
	if err = s.FlushBuffer(); err != nil {
	}
	return session.Normal
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

	data, _ := pkg.(string)

	return []byte(data), nil
}

func (c codec) Decode(bytes []byte) (interface{}, int, error) {
	if bytes == nil {
		return nil, 0, errors.New("bytes is nil")
	}

	data := string(bytes)

	if len(data) == 0 {
		return nil, 0, nil
	}
	return data, len(data), nil
}
