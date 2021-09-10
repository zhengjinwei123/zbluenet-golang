package znet

import (
	"fmt"
	"net"
	"sync"
)


const MIN_NET_ID = 1000

type udpService struct {
	socket *udpSocket
	addr *SocketAddress
	net_id_allocator *NetIdAllocator

	running bool
	connections map[uint32]*udpConnection
	server *UdpServer

	conMutex *sync.Mutex
	closeChan chan struct{}
	wg *sync.WaitGroup
}

func NewUdpService() *udpService {
	return &udpService{
		net_id_allocator: NewNetIdAllocator(MIN_NET_ID),
		socket: nil,
		addr: nil,
		running: false,
		connections: make(map[uint32]*udpConnection),
		server: nil,
		conMutex: &sync.Mutex{},
		closeChan: make(chan struct{}),
		wg: &sync.WaitGroup{},
	}
}

func (this *udpService) CreateServer(host string, port int, server *UdpServer) error {
	this.addr = NewSocketAddress(host, port, SOCKET_PROTOCOL_IPV4)
	this.socket = NewUdpSocket()

	if err := this.socket.PassiveOpen(this.addr); err != nil {
		return err
	}

	this.server = server

	return nil
}

func (this *udpService) GetListenAddr() string {
	return this.addr.GetListenAddr()
}

func (this *udpService) stop(shutdownWg *sync.WaitGroup) {
	if this.socket == nil {
		shutdownWg.Done()
		return
	}

	this.wg.Add(1)
	close(this.closeChan)

	go func() {
		this.wg.Wait()

		fmt.Printf("udpService stop\n")
		_ = this.socket.CloseSocket()
		this.socket = nil

		for _, con := range this.connections {
			con.close()
		}

		shutdownWg.Done()
	}()
}

func (this *udpService) Shutdown(shutdownWg *sync.WaitGroup) {
	this.running = false
	this.stop(shutdownWg)
}

func (this *udpService) Loop() {
	for {
		select {
		case <-this.closeChan:
			goto exit
		default:
			data := make([]byte, 65535)
			n, addr, err := this.socket.Accept(data)
			if err != nil {
				// 超时
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					continue
				}

				continue
			}

			go this.handleConn(addr, data, n)
		}
	}
exit:
	this.wg.Done()
	fmt.Printf("udpservice loop exit \n")
}


func (this *udpService) handleConn(addr *net.UDPAddr, buf []byte, size int) {

	if size < UDP_PACKET_HEAD_SIZE {
		return
	}
	// 解包
	buffer := NewRecvUpdBuffer(buf)
	if buffer.isFull == false {
		return
	}

	// 获取 SessionId
	if buffer.SessionId < MIN_NET_ID {
		// 如果 sessionid < MIN_NET_ID , 说明是连接请求， 分配一个net_id, 然后放到包体中
		id := this.net_id_allocator.NextId()
		buffer.SessionId = id

		this.createConnection(addr, buffer.SessionId)
	}

	if _, exists := this.connections[buffer.SessionId]; exists {
		this.connections[buffer.SessionId].onRecvMessage(buffer)
	}
}


func (this *udpService) createConnection(addr *net.UDPAddr, session_id uint32) {
	if _, exists := this.connections[session_id]; exists {
		return
	}
	this.conMutex.Lock()
	this.connections[session_id] = NewUdpConnection(addr, session_id, nil)
	this.conMutex.Unlock()

	this.server.onConnect(session_id, addr.String())
}

func (this *udpService) dispatchMessage(buffer *udpBuffer) {
	this.server.onMessage(buffer.SessionId, buffer.MessageId, buffer.Data)
}

func (this *udpService) removeClient(sessionId uint32) {
	this.conMutex.Lock()
	defer this.conMutex.Unlock()

	if conn, exists := this.connections[sessionId]; exists {

		this.server.onClose(sessionId, conn.addr.String())
		delete(this.connections, sessionId)
	}
}

func (this *udpService) sendMessage(addr *net.UDPAddr, buffer *udpBuffer) error {
	return this.socket.SendMessage(addr, buffer.Encode(false))
}

func (this *udpService) SendMessageToSession(session_id uint32, message_id uint16, data []byte) {
	if conn, exists := this.connections[session_id]; exists {

		if conn.isConnected == false {
			return
		}

		buffer := NewRequestUdpBuffer(session_id, 0, UDP_MESSAGE_TYPE_MESSAGE, message_id, data)
		conn.sendMessage(buffer)
	}
}
