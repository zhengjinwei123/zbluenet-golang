package znet

import (
	"fmt"
	"net"
)



type TcpService struct {
	sock *TcpSocket
	addr *SocketAddress
	max_con int
	net_id_allocator *NetIdAllocator

	con_num int

	running bool
}

func NewTcpService() *TcpService {
	return &TcpService{
		sock: nil,
		addr: nil,
		running: false,
		max_con: -1,
		con_num: 0,
		net_id_allocator: NewNetIdAllocator(0),
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

	this.max_con = max_connection

	if err := this.sock.PassiveOpen(this.addr); err != nil {
		return err
	}

	return nil
}

func (this *TcpService) stop() {
	if this.sock == nil {
		return
	}
	fmt.Printf("tcp service listen socket close\n")
	_ = this.sock.CloseSocket()
	this.sock = nil
}

func (this *TcpService) Shutdown() {
	fmt.Printf("tcp service ready shutdown \n")
	this.running = false
	this.stop()
}

func (this *TcpService) GetListenAddr() string {
	return this.addr.GetListenAddr()
}

func (this *TcpService) Loop(srv *TcpServer) {
	this.running = true

	for this.running == true {

		conn, err := this.sock.Accept()
		if err != nil {

			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			}
			fmt.Printf("tcp service loop err:%s \n", err.Error())
			return
		}

		if this.max_con > 0 && this.con_num >= this.max_con {
			_ = conn.Close()
			continue
		}

		this.con_num ++
		NewTcpReactor(srv).Loop(this.net_id_allocator.NextId(), conn)
	}

	fmt.Printf("tcp service loop exit \n")
}