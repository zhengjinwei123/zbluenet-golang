package sudp

import (
	"app/base/znet"
	"net"
)

//type MessageCallback func(addr *net.UDPAddr, read_size int, data []byte)

type udpSocket struct {
	listenConn *net.UDPConn
}

func NewUdpSocket() *udpSocket {
	return &udpSocket{}
}

func (this *udpSocket) PassiveOpen(addr *znet.SocketAddress) error {
	addr1, err := net.ResolveUDPAddr("udp", addr.GetListenAddr())
	if err != nil {
		return err
	}

	this.listenConn, err = net.ListenUDP("udp", addr1)
	if err != nil {
		return err
	}

	return nil
}

func (this *udpSocket) CloseSocket() error {
	return this.listenConn.Close()
}

func (this *udpSocket) Accept(data []byte) (int, *net.UDPAddr, error) {
	return this.listenConn.ReadFromUDP(data)
}

func (this *udpSocket) SendMessage(addr *net.UDPAddr, data []byte) error {
	_, err := this.listenConn.WriteToUDP(data, addr)
	return err
}