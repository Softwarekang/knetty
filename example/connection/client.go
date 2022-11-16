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

	n, err := conn.Write([]byte("J9eMSxIxdVr2b74m8333zZGm6er5iaBnKc7wClWX3WMdI72Wqkxsx6eEHXCRgzrqxust8zeIc5nSXGcTEMXvQ5VR089oxQBFSEt0hGE4MtV2dCdxeMyzFsbRxtkylGtmXqidpuheUDH7CHLzidMF9X4E2MXZivB7Ubn9tV4WGT8Pbt6UeuEjIm3LtuImf63S0gwP0McRfafzUIGyCe2BIMueIICgvgjes4o6xuFKibxCWhp0aHOy7mqmoTsNc7XbnZ9"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("client send %d bytes data to server:%s\n", n, conn.RemoteAddr().String())
	time.Sleep(50 * time.Second)
}
