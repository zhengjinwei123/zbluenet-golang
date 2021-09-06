package net

import (
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
	tcp_session TcpSession

	ReadMessageQueue chan *NetMessage
	writeMessageQueue chan *NetMessage

	state int32
	CloseChan chan struct{}

	wg sync.WaitGroup
}

func NewTcpReactor(session TcpSession) *TcpReactor {
	return &TcpReactor{
		tcp_session:       session,
	}
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

func (this *TcpReactor) Loop(conn net.Conn) bool {
	if !atomic.CompareAndSwapInt32(&this.state, CONNECTION_STATE_INIT, CONNECTION_STATE_CONNECTED) {
		return false
	}

	this.tcp_conn = NewConnection(conn, this,20, 60)
	this.CloseChan = make(chan struct{})

	this.writeMessageQueue = make(chan *NetMessage, 819200)
	this.ReadMessageQueue = make(chan *NetMessage, 819200)

	this.wg.Add(1)
	go this.tcp_conn.WriteLoop()

	this.wg.Add(1)
	go this.tcp_conn.ReadLoop()

	this.tcp_session.OnInit(this)

	return true
}

// 得到一个完整的包
func (this *TcpReactor) recvPacket(messageId uint16, messageLength uint16, messageBuf []byte) {
	// 通过messageId 创建对应的消息， 然后

	pack := NewNetMessage(messageId, messageLength, messageBuf)
	select {
	case this.ReadMessageQueue <- pack:
	case <- this.CloseChan:
	}
}