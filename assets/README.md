# knetty

## Introduction

knetty is a high-performance, asynchronous IO network library focused on RPC websockets and supporting TCP, UDP, and WebSocket protocols. It is similar to Netty in functionality.

## Contents

- [knetty](#knetty)
  - [Contents](#contents)

### Start Server

```go
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

```

