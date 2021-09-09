package sudp

type LogicServer interface {
	OnConnect(sessionId uint32, addr string)
	OnMessage(sessionId uint32, messageId uint16, data []byte)
}

type UdpServer struct {
	service *udpService
	logicServer LogicServer
	closeChan chan struct{}
}

func NewUdpServer(logicServ LogicServer) *UdpServer {
	return &UdpServer{
		service: NewUdpService(),
		logicServer: logicServ,
		closeChan: make(chan struct{}),
	}
}

func (this *UdpServer) CreateServer(host string, port int) error {
	return this.service.CreateServer(host, port, this)
}

func (this *UdpServer) Loop() {
	this.service.Loop()
}

func (this *UdpServer) Shutdown() {
	this.service.Shutdown()
	close(this.closeChan)
}

func (this *UdpServer) onConnect(sessionId uint32, addr string) {
	this.logicServer.OnConnect(sessionId, addr)
}

// 这个接口是并发的
func (this *UdpServer) onMessage(sessionId uint32, messageId uint16, data []byte) {
	this.logicServer.OnMessage(sessionId, messageId, data)
}

func (this *UdpServer) SendMessage(sessionId uint32, message_id uint16, data []byte) {
	this.service.SendMessageToSession(sessionId, message_id, data)
}