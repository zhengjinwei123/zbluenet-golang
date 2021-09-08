package main

import (
	"app/base/znet"
	"app/protocol"
	"fmt"
	"github.com/golang/protobuf/proto"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ZoneServer struct {

}

func (this *ZoneServer) OnConnect(fd int64, remoteAddr string) {
	fmt.Printf("OnConnect [%d] [%s] \n", fd, remoteAddr)
}

func (this *ZoneServer) OnDisConnect(fd int64, remoteAddr string) {
	fmt.Printf("OnDisConnect [%d] [%s] \n", fd, remoteAddr)
}


// 这个接口是并发的， 需要把消息转到channel ， 然后再处理比较好，
func (this *ZoneServer) OnMessage(fd int64,  messageId uint16, data []byte) {
	fmt.Printf("OnMessage [%d] [%d]\n", fd, messageId)

	recv := &protocol.S2CLoginResp{}
	err := proto.Unmarshal(data, recv)
	if err != nil {
		fmt.Printf("OnMessage error: %s\n", err.Error())
		return
	}
	fmt.Printf("OnMessage success: %s\n", recv.String())
}

var g_signal  = make(chan os.Signal, 1)


func main() {

	zonserv := &ZoneServer{}

	server := znet.NewTcpServer(zonserv)
	if err := server.CreateServer("127.0.0.1", 8888, -1); err != nil {
		fmt.Printf("CreateServer failed err:%s \n", err.Error())
		return
	}

	go server.Loop()

	fmt.Printf("server start success, addr: %s\n", server.GetListenAddress())

	listenSignal(server)
}

func listenSignal(server *znet.TcpServer) {
	signal.Notify(g_signal, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGINT)

	select {
	case sig := <-g_signal:
		fmt.Printf("catch signal: %s \n", sig.String())
		server.Shutdown()
	}

	time.Sleep(time.Duration(1) * time.Second)
}
