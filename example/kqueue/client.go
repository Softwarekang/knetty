package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	network, address := "tcp", "127.0.0.1:8000"
	conn, err := net.Dial(network, address)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(conn.RemoteAddr())
}
