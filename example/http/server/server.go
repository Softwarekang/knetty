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
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Softwarekang/knetty"
	"github.com/Softwarekang/knetty/session"

	"github.com/evanphx/wildcat"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

func init() {
	initLogger()
}

func main() {
	// setting optional options for the server
	options := []knetty.ServerOption{
		knetty.WithServiceNewSessionCallBackFunc(newSessionCallBackFn),
	}

	knetty.SetLogger(logger)
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
	logger.Infof("shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("server starting shutdown:", err)
	}

	logger.Info("server exiting")
}

func initLogger() {
	logger = logrus.StandardLogger()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
}

// set the necessary parameters for the session to run.
func newSessionCallBackFn(s session.Session) error {
	s.SetCodec(newHttpEchoCodec())
	s.SetEventListener(&httpEventListener{})
	return nil
}

type httpEventListener struct {
}

func (e *httpEventListener) OnMessage(s session.Session, pkg interface{}) session.ExecStatus {
	data := pkg.(string)
	logger.Infof("server got data:%s", data)
	if err := echoHello(s, data); err != nil {
		logger.Infof("server echo err:%v", err)
	}
	return session.Normal
}

func echoHello(s session.Session, data string) error {
	rsp := []byte("HTTP/1.1 200 OK\r\nServer: knetty\r\nContent-Type: text/plain\r\nDate: ")
	rsp = time.Now().AppendFormat(rsp, "Mon, 02 Jan 2006 15:04:05 GMT")
	rsp = append(rsp, fmt.Sprintf("\r\nContent-Length: %d\r\n\r\n%s!", len(data)+1, data)...)

	if _, err := s.WriteBuffer(rsp); err != nil {
		return err
	}

	return s.FlushBuffer()
}

func (e *httpEventListener) OnConnect(s session.Session) {
	logger.Infof("local:%s get a remote:%s connection\n", s.LocalAddr(), s.RemoteAddr())
}

func (e *httpEventListener) OnClose(s session.Session) {
	logger.Infof("server client:%s  closed", s.RemoteAddr())
}

func (e *httpEventListener) OnError(s session.Session, err error) {
	logger.Errorf("session: %s got err :%v", s.Info(), err)
}

type httpEchoCodec struct {
	httpParser *wildcat.HTTPParser
}

func newHttpEchoCodec() *httpEchoCodec {
	return &httpEchoCodec{
		httpParser: wildcat.NewHTTPParser(),
	}
}
func (h *httpEchoCodec) Encode(pkg interface{}) ([]byte, error) {
	return []byte(nil), nil
}

func (h *httpEchoCodec) Decode(bytes []byte) (interface{}, int, error) {
	headerOffset, err := h.httpParser.Parse(bytes)
	if err != nil {
		if errors.Is(err, wildcat.ErrMissingData) {
			return nil, 0, nil
		}

		return nil, 0, err
	}

	if !h.httpParser.Get() {
		return nil, 0, fmt.Errorf("http method need %value", http.MethodGet)
	}

	url, err := url.Parse(string(h.httpParser.Path))
	if err != nil {
		return nil, 0, err
	}

	value := url.Query().Get("echo")
	if value == "" {
		value = "default"
	}
	return value, headerOffset, nil
}
