package znet

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
)

const (
	CONNECTION_STATE_INIT = 1
	CONNECTION_STATE_CONNECTED = 2
	CONNECTION_STATE_DISCONNECTED = 3
)


// socket read and write processor
type TcpReactor struct {
	tcp_conn *TcpConnection
	tcp_server *TcpServer

	ReadMessageQueue chan *NetMessage
	writeMessageQueue chan *NetMessage

	state int32
	CloseChan chan struct{}

	wg sync.WaitGroup
}

func NewTcpReactor(serv *TcpServer) *TcpReactor {
	return &TcpReactor{
		tcp_server:       serv,
		state: CONNECTION_STATE_INIT,
	}
}


func (this *TcpReactor) GetNetId() uint32 {
	return this.tcp_conn.id
}

func (this *TcpReactor) LocalAddress() string {
	return this.tcp_conn.local_addr
}

func (this *TcpReactor) RemoteAddress() string {
	return this.tcp_conn.remote_addr
}

func (this *TcpReactor) Status() int32 {
	return atomic.LoadInt32(&this.state)
}

func (this *TcpReactor) Stop() {
	if !atomic.CompareAndSwapInt32(&this.state, CONNECTION_STATE_CONNECTED, CONNECTION_STATE_DISCONNECTED) {
		return
	}

	// triggle close signal
	close(this.CloseChan)

	_ = this.tcp_conn.Close()

	go func() {
		this.wg.Wait()

		this.tcp_server.onClose(this)

		close(this.writeMessageQueue)
		close(this.ReadMessageQueue)

		atomic.StoreInt32(&this.state, CONNECTION_STATE_INIT)

		fmt.Printf("reactor closed \n")
	}()
}

func (this *TcpReactor) Loop(id uint32, conn net.Conn) bool {

	if !atomic.CompareAndSwapInt32(&this.state, CONNECTION_STATE_INIT, CONNECTION_STATE_CONNECTED) {
		return false
	}


	this.tcp_conn = NewConnection(id, conn, this,600, 600)
	this.CloseChan = make(chan struct{})

	this.writeMessageQueue = make(chan *NetMessage, 10240)
	this.ReadMessageQueue = make(chan *NetMessage, 10240)

	this.tcp_server.onConnect(this)

	// 写
	this.wg.Add(1)
	go this.tcp_conn.WriteLoop()

	// 读
	this.wg.Add(1)
	go this.tcp_conn.ReadLoop()

	// 派发
	this.wg.Add(1)
	go this.dispatchMessage()

	return true
}

// 得到一个完整的包
func (this *TcpReactor) recvPacket(messageId uint16, messageLength uint16, messageBuf []byte) {
	// 通过messageId 创建对应的消息， 然后

	pack := NewNetMessage(messageId, this.tcp_conn.id, messageBuf)
	select {
	case this.ReadMessageQueue <- pack:
	case <- this.CloseChan:
	}
}

func (this *TcpReactor) dispatchMessage() {
	for more := true; more; {
		select {
		case msg := <- this.ReadMessageQueue:

			// 将字节数组 转化为 消息结构
			this.tcp_server.onRecvMessage(this, msg)
		case <-this.CloseChan:
			more = false
		case <-this.tcp_server.closeChan:
			more = false
		}
	}
	fmt.Printf("dispatchMessage exit \n")
	this.wg.Done()
	this.Stop()
}

func (this *TcpReactor) SendMessage(messageId uint16, data interface{}) bool {
	msg := &NetMessage{
		NetId: 0,
		MessageId: messageId,
		Data: nil,
		DataSend: data,
	}
	select {
	case this.writeMessageQueue <- msg:
		return true
	case <-this.CloseChan:
		return false
	}
}

