package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	network, address := "tcp", "127.0.0.1:8000"
	conn, err := net.Dial(network, address)
	if err != nil {
		log.Fatal(err)
	}

	n, err := conn.Write([]byte("hello"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("client send %d bytes data to server:%s\n", n, conn.RemoteAddr().String())
	time.Sleep(5 * time.Second)
}
