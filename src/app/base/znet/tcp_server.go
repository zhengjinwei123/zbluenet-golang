package znet

type TcpLogicServer interface {
	OnConnect(net_id uint32, remoteAddr string)
	OnDisConnect(net_id uint32, remoteAddr string)
	OnMessage(net_id uint32, messageId uint16, data []byte)
}

type TcpServer struct {
	service *TcpService
	logicServer TcpLogicServer
	closeChan chan struct{}
}

func NewTcpServer(logicServ TcpLogicServer) *TcpServer {
	return &TcpServer{
		service: NewTcpService(),
		logicServer: logicServ,
		closeChan: make(chan struct{}),
	}
}

func (this *TcpServer) CreateServer(host string, port int, max_connection int) error {
	return this.service.CreateServer(host, port, max_connection)
}

func (this *TcpServer) Loop() {
	this.service.Loop(this)
}

func (this *TcpServer) GetListenAddress() string {
	return this.service.GetListenAddr()
}

func (this *TcpServer) Shutdown() {
	// 通知监听程序关闭
	this.service.Shutdown()
	// 通知 reactor 关闭
	close(this.closeChan)
}

func (this *TcpServer) onConnect(reactor *TcpReactor) {
	this.logicServer.OnConnect(reactor.GetNetId(), reactor.tcp_conn.remote_addr)
}

func (this *TcpServer) onClose(reactor *TcpReactor) {
	this.logicServer.OnDisConnect(reactor.GetNetId(), reactor.tcp_conn.remote_addr)
}

func (this *TcpServer) onRecvMessage(reactor *TcpReactor, message *NetMessage) {
	this.logicServer.OnMessage(reactor.GetNetId(), message.MessageId, message.Data)
}

func (this *TcpServer) SendMessage(reactor *TcpReactor, messageId uint16, data interface{}) {
	reactor.SendMessage(messageId, data)
}


