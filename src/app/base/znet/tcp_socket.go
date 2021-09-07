package znet

import "net"

type TcpSocket struct {
	listener net.Listener
}

func NewTcpSocket() *TcpSocket {
	return &TcpSocket{}
}

// 创建服务
func (this *TcpSocket) PassiveOpen(addr *SocketAddress) error {
	var err error
	this.listener, err = net.Listen("tcp", addr.GetListenAddr())
	return err
}

func (this *TcpSocket) CloseSocket() error {
	return this.listener.Close()
}

func (this *TcpSocket) Accept() (net.Conn, error) {
	return this.listener.Accept()
}


