package main

import (
	"app/base/znet"
	"app/protocol"
	"fmt"
	"github.com/golang/protobuf/proto"
	"os"
	"os/signal"
	"syscall"
)

type BattleServer struct {

}

func (this *BattleServer) OnConnect(sessionId uint32, addr string) {
	fmt.Printf("BattleServer OnConnect [%d] [%s] \n", sessionId, addr)
}

func (this *BattleServer) OnMessage(sessionId uint32, messageId uint16, data []byte) {
	fmt.Printf("BattleServer OnMessage [%d] [%d] [%s] \n", sessionId, messageId, string(data))
}

func (this *BattleServer) OnDisconnect(sessionId uint32, remoteAddr string) {
	fmt.Printf("BattleServer OnDisconnect [%d] [%s] \n", sessionId, remoteAddr)
}



type ZoneServer struct {

}

func (this *ZoneServer) OnConnect(fd uint32, remoteAddr string) {
	fmt.Printf("OnConnect [%d] [%s] \n", fd, remoteAddr)
}


func (this *ZoneServer) OnDisconnect(fd uint32, remoteAddr string) {
	fmt.Printf("OnDisconnect [%d] [%s] \n", fd, remoteAddr)
}

func (this *ZoneServer) OnDisConnect(fd uint32, remoteAddr string) {
	fmt.Printf("OnDisConnect [%d] [%s] \n", fd, remoteAddr)
}


// 这个接口是并发的， 需要把消息转到channel ， 然后再处理比较好，
func (this *ZoneServer) OnMessage(fd uint32,  messageId uint16, data []byte) {
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
		fmt.Printf("TCP CreateServer failed err:%s \n", err.Error())
		return
	}

	go server.Loop()

	fmt.Printf("tcp server start success, addr: %s\n", server.GetListenAddress())


	battleServer := &BattleServer{}
	udpServ := znet.NewUdpServer(battleServer)

	if err := udpServ.CreateServer("127.0.0.1", 8889); err != nil {
		fmt.Printf("UDP CreateServer failed err:%s \n", err.Error())
		return
	}
	go udpServ.Loop()

	fmt.Printf("udp server start success, addr: %s\n", udpServ.GetListenAddress())

	listenSignal(server, udpServ)
}

func listenSignal(tcpServer *znet.TcpServer, udpServ *znet.UdpServer) {
	signal.Notify(g_signal, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGINT)

	select {
	case sig := <-g_signal:
		fmt.Printf("catch signal: %s \n", sig.String())

		tcpServer.Shutdown()

		udpServ.Shutdown()
	}

}
