# knet

## Introduction

Knet is a high-performance, asynchronous IO network library focused on RPC websockets and supporting TCP, UDP, and WebSocket protocols. It is similar to Netty in functionality.

## Contents

- [Knet](#knet)
  - [Contents](#contents)

### Start Server

```go
// setting optional options for the server
options := []knet.ServerOption{
    knet.WithServiceNewSessionCallBackFunc(newSessionCallBackFn),
}

// creating a new server with network settings such as tcp/upd, address such as 127.0.0.1:8000, and optional options
server := knet.NewServer("tcp", "127.0.0.1:8000", options...)
// starting the server
if err := server.Server(); err != nil {
    log.Fatalln(err)
    return
}

```

