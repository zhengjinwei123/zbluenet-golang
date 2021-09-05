package net

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

func (this *TcpService) Loop() {

	for {
		conn, err := this.sock.Accept()
		if err != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				continue
			}
			return
		}
		go this.onNewConnection(conn)
	}
}

func (this *TcpService) onNewConnection(conn net.Conn) {
	defer conn.Close()

	for {
		buf
	}
}