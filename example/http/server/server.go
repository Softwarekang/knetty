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
	s.SetCodec(newHttpEchoCodec())
	s.SetEventListener(&httpEventListener{})
	s.SetReadTimeout(1 * time.Second)
	s.SetWriteTimeout(1 * time.Second)
	return nil
}

type httpEventListener struct {
}

func (e *httpEventListener) OnMessage(s session.Session, pkg interface{}) {
	data := pkg.(string)
	fmt.Printf("server got data:%s\n", data)
	if err := echoHello(s, data); err != nil {
		log.Printf("server echo err:%v", err)
	}
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
	fmt.Printf("local:%s get a remote:%s connection\n", s.LocalAddr(), s.RemoteAddr())
}

func (e *httpEventListener) OnClose(s session.Session) {
	fmt.Printf("server session: %s closed\n", s.Info())
}

func (e *httpEventListener) OnError(s session.Session, err error) {
	fmt.Printf("session: %s got err :%v\n", s.Info(), err)
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
