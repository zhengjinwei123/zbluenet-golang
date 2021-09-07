package znet

import "fmt"

type LogicServer interface {
	OnConnect(fd int)
}

type TcpServer struct {
	service *TcpService
}

func NewTcpServer() *TcpServer {
	return &TcpServer{
		service: NewTcpService(),
	}
}

func (this *TcpServer) CreateServer(host string, port int, max_connection int) error {
	return this.service.CreateServer(host, port, max_connection)
}

func (this *TcpServer) Loop() {
	this.service.Loop(this)
}

func (this *TcpServer) OnConnect(reactor *TcpReactor) {
	fmt.Printf("OnConnect (%d)\n", reactor.GetFd())
}

func (this *TcpServer) OnClose(reactor *TcpReactor) {
	fmt.Printf("OnClose (%d)\n", reactor.GetFd())
}

func (this *TcpServer) OnRecvMessage(reactor *TcpReactor, message *NetMessage) {
	fmt.Printf("OnRecvMessage (%d)\n", reactor.GetFd())
}

func (this *TcpServer) SendMessage(reactor *TcpReactor, messageId uint16, data interface{}) {
	reactor.SendMessage(messageId, data)
}


