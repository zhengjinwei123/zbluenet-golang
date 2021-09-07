package znet

import "net"



type TcpService struct {
	sock *TcpSocket
	addr *SocketAddress
	max_con int
}

func NewTcpService() *TcpService {
	return &TcpService{
		sock: nil,
		addr: nil,
	}
}

func (this *TcpService) CreateServer(host string, port int, max_connection int) error {
	this.addr = &SocketAddress{
		host:            host,
		port:            port,
		socket_protocol: 1,
	}

	this.sock = &TcpSocket{
		listener:nil,
	}
	if err := this.sock.PassiveOpen(this.addr); err != nil {
		return err
	}

	return nil
}

func (this *TcpService) Loop(srv *TcpServer) {

	for {
		conn, err := this.sock.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			}
			return
		}
		NewTcpReactor(srv).Loop(conn)
	}
}