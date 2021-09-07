package main

import (
	"app/base/znet"
	"fmt"
)


func main() {
	server := znet.NewTcpServer()
	if err := server.CreateServer("127.0.0.1", 8888, 100000); err != nil {
		fmt.Printf("CreateServer failed err:%s \n", err.Error())
		return
	}
	server.Loop()
}
