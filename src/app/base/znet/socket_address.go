package znet

import "fmt"

const (
	SOCKET_PROTOCOL_IPV4 = 1
	SOCKET_PROTOCOL_IPV6 = 2
)


type SocketAddress struct {
	host string
	port int
	socket_protocol int
}

func NewSocketAddress(host string, port, socket_protocol int) *SocketAddress {
	return &SocketAddress{
		host: host,
		port: port,
		socket_protocol: socket_protocol,
	}
}

func (this *SocketAddress) GetListenAddr() string {
	return fmt.Sprintf("%s:%d", this.GetIp(), this.GetPort());
}

func (this *SocketAddress) GetIp() string {
	return this.host
}

func (this *SocketAddress) GetPort() int {
	return this.port
}



