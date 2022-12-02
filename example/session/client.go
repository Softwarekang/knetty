package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	var err error
	network, address := "tcp", "127.0.0.1:8000"
	conn, err := net.Dial(network, address)
	if err != nil {
		log.Fatal(err)
	}

	n, err := conn.Write([]byte("hellohello"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("send length:%d", n)
	time.Sleep(5 * time.Second)
}
