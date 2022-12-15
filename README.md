# knetty

## Introduction

knetty is a network communication framework written in Go based on the event-driven architecture. It supports TCP, UDP, and websocket protocols and is easy to use like Netty written in Java."

## Contents
- [knetty](#knetty)
  - [Introduction](#introduction)
  - [Contents](#contents)
  - [Installation](#installation)
  - [Quick Start](#quick-start)
  - [More Detail](#more-detail)
    - [Using NewSessionCallBackFunc](#using-newsessioncallbackfunc)
    - [Using Codec](#using-codec)
    - [Using EventListener](#using-eventlistener)
    - [Graceful shutdown](#graceful-shutdown)
  - [Benchmarks](#benchmarks)

## Installation
To install knetty package,you nedd to install Go and set your Go workspace first.

- You first need [Go](https://golang.org/) installed (**version 1.18+ is required**), then you can use the below Go command to install knetty.
```shell
go get -u  github.com/Softwarekang/knetty
```
- import in your code
```go
import "github.com/Softwarekang/knetty"
```
    
## Quick Start

```sh
# View knetty code examples
# work dir in knetty
cd /example/sever
```

```sh
# view server start up code examples
cat server.go
```
```go
package main

import (
	"errors"
	"fmt"
	"log"

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
	// starting the server
	if err := server.Server(); err != nil {
		log.Fatalln(err)
		return
	}
}

// set the necessary parameters for the session to run.
func newSessionCallBackFn(s session.Session) error {
	s.SetCodec(&codec{})
	s.SetEventListener(&helloWorldListener{})
	return nil
}

type helloWorldListener struct {}

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
```
```sh
# start up server 
go run ./example/sever/server/server.go
```
# view server start up code examples
cat server.go
```sh
# view client start up code examples
cat client.go
```
```go
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

```
``` sh
# start up client
go run ./example/sever/client/client.go
```

## More Detail

### Using NewSessionCallBackFunc

definition
```go
/*
	NewSessionCallBackFunc It is executed when a new session is established,
	so some necessary parameters for drawing need to be set to ensure that the session starts properly.
*/
type NewSessionCallBackFunc func(s session.Session) error
```
You can set parameters such as codec, event listener, read/write timeouts, and more for the session via the provided API.
```go
// set the necessary parameters for the session to run.
func newSessionCallBackFn(s session.Session) error {
	s.SetCodec(&codec{})
	s.SetEventListener(&helloWorldListener{})
	s.SetReadTimeout(1 * time.Second)
	s.SetWriteTimeout(1 * time.Second)
	return nil
}
```
### Using Codec
definition
```go
// Codec for session
type Codec interface {
	// Encode will convert object to binary network data
	Encode(pkg interface{}) ([]byte, error)

	// Decode will convert binary network data into upper-layer protocol objects.
	// The following three conditions are used to distinguish abnormal, half - wrapped, normal and sticky packets.
	// Exceptions: nil,0,err
	// Half-pack: nil,0,nil
	// Normal & Sticky package: pkg,pkgLen,nil
	Decode([]byte) (interface{}, int, error)
}
```
Here is an implementation of a hello string boundary encoder that handles semi-packet, sticky packet, and exceptional network data processing logic.
```go
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
```

### Using EventListener
definition
```go
// EventListener listen for session
type EventListener interface {
	// OnMessage runs when the session gets a pkg
	OnMessage(s Session, pkg interface{})
	// OnConnect runs when the connection initialized
	OnConnect(s Session)
	// OnClose runs before the session closed
	OnClose(s Session)
	// OnError runs when the session err
	OnError(s Session, e error)
}
```
Below is a typical event listener.
```go
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
	fmt.Printf("session got err :%v\n", err)
}
```

### Graceful shutdown 
```go
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
```
## Benchmarks