package main

import (
	"fmt"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("listen err =", err)
		return
	}

	defer listen.Close()

	fmt.Println("server start at:", listen.Addr())

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		fmt.Printf("client [%s] connected \n", conn.RemoteAddr().String())

		go ProcessNewConnection(conn)
	}
}

func ProcessNewConnection(conn net.Conn) {
	addr := conn.RemoteAddr().String()

	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("client [%s] has exited", addr)
			return
		}

		fmt.Printf("client[%s] recv message:[%s]",addr, string(buf[:n]))
	}
}