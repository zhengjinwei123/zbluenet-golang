package znet

import "sync"

type IUdpLogicServer interface {
	OnConnect(sessionId uint32, addr string)
	OnDisconnect(sessionId uint32, addr string)
	OnMessage(sessionId uint32, message_id uint16, data []byte)
	OnInit(server *UdpServer)
}

type UdpServer struct {
	service *udpService
	logicServer IUdpLogicServer
	wg *sync.WaitGroup
}

func NewUdpServer(logicServ IUdpLogicServer) *UdpServer {
	return &UdpServer{
		service: NewUdpService(),
		logicServer: logicServ,
		wg: &sync.WaitGroup{},
	}
}

func (this *UdpServer) GetListenAddress() string {
	return this.service.GetListenAddr()
}

func (this *UdpServer) CreateServer(host string, port int) error {
	this.logicServer.OnInit(this)

	return this.service.CreateServer(host, port, this)
}

func (this *UdpServer) Loop() {
	this.service.Loop()
}

func (this *UdpServer) Shutdown() {
	this.wg.Add(1)
	this.service.Shutdown(this.wg)
	this.wg.Wait()
}

func (this *UdpServer) onConnect(sessionId uint32, addr string) {
	this.logicServer.OnConnect(sessionId, addr)
}

func (this *UdpServer) onClose(sessionId uint32, addr string) {
	this.logicServer.OnDisconnect(sessionId, addr)
}

// 这个接口是并发的
func (this *UdpServer) onMessage(sessionId uint32, message_id uint16, data []byte) {
	this.logicServer.OnMessage(sessionId, message_id, data)
}

func (this *UdpServer) SendMessage(sessionId uint32, message_id uint16, data []byte) {
	this.service.SendMessageToSession(sessionId, message_id, data)
}

