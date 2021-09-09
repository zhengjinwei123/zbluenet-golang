package sudp

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const SEND_MESSAGE_TIMEOUT_MS = 150 // 150 超时重传

type udpConnection struct {
	addr        *net.UDPAddr
	sendSN      uint32
	handleSN    uint32
	sessionId   uint32
	sendedPackageMap map[uint32]*udpBuffer // 存储已经发送的报文
	awaitPackageMap map[uint32]*udpBuffer // 存储错序的报文

	sendMutex  *sync.Mutex
	awaitMutex *sync.Mutex

	isConnected bool
	service *udpService

	closeChan chan struct{}
}

func NewUdpConnection(addr *net.UDPAddr, session_id uint32, service *udpService) *udpConnection {
	t := &udpConnection{
		addr: addr,
		sessionId: session_id,
		sendSN: 0,
		handleSN: 0,
		sendedPackageMap: make(map[uint32]*udpBuffer, 0),
		awaitPackageMap: make(map[uint32]*udpBuffer, 0),
		isConnected: true,
		service: service,
		closeChan: make(chan struct{}),
		sendMutex: &sync.Mutex{},
		awaitMutex: &sync.Mutex{},
	}
	go t.checkOutTime()
	return t
}

func (this *udpConnection) checkOutTime() {
	timer := time.NewTicker(time.Duration(SEND_MESSAGE_TIMEOUT_MS) * time.Millisecond)

	running := true
	for running == true {
		select {
		case <-timer.C:
			this.handleTimeoutPackage()
		case <-this.closeChan:
			running = false
		}
	}

	fmt.Printf("checkOutTime exit \n")
}


func (this *udpConnection) handleTimeoutPackage() {

	// copy
	this.sendMutex.Lock()
	tmpMap := make(map[uint32]*udpBuffer, len(this.sendedPackageMap))
	for sessionId, buffer := range this.sendedPackageMap {
		tmpMap[sessionId] = buffer
	}
	this.sendMutex.Unlock()

	now := time.Now().Unix()

	for sessionId, buffer := range tmpMap {
		if buffer.reSendCount >= 10 {
			this.service.removeClient(sessionId)
			return
		}

		if int(now - int64(buffer.Time)) >= (buffer.reSendCount + 1) * SEND_MESSAGE_TIMEOUT_MS {
			buffer.reSendCount += 1

			fmt.Printf("超时重发: (sn:%d) (message_id:%d) \n", buffer.SN, buffer.MessageId)
			_ = this.service.sendMessage(this.addr, buffer)
		}

	}
}

func (this *udpConnection) close() {
	this.isConnected = false
}

func (this *udpConnection) onRecvMessage(buffer *udpBuffer) {
	if buffer.MessageType == UDP_MESSAGE_TYPE_ACK {
		// ack 报文
		this.sendMutex.Lock()
		defer this.sendMutex.Unlock()
		delete(this.sendedPackageMap, buffer.SN)
	} else {
		// 业务报文
		this.handleLogicPackage(buffer)
	}
}

func (this *udpConnection) handleLogicPackage(buffer *udpBuffer) {
	if buffer.SN <= this.handleSN {
		// 已经处理过的消息， 直接过滤
		return
	}

	if buffer.SN - this.handleSN > 1 {
		// 收到错序的报文, save it
		this.awaitMutex.Lock()
		defer this.awaitMutex.Unlock()
		this.awaitPackageMap[buffer.SN] = buffer

		return
	}

	this.handleSN = buffer.SN

	// 分发消息
	this.service.dispatchMessage(buffer)

	// 处理队列里的消息
	if _, exists := this.awaitPackageMap[this.handleSN + 1]; exists {
		this.handleLogicPackage(this.awaitPackageMap[this.handleSN + 1])


		this.awaitMutex.Lock()
		defer this.awaitMutex.Unlock()
		delete(this.awaitPackageMap, this.handleSN + 1)
	}
}

func (this *udpConnection) Send(buffer *udpBuffer) {
	if this.isConnected == false {
		return
	}

	buffer.Time = uint64(time.Now().Unix())
	this.sendSN += 1

	buffer.SN = this.sendSN

	buffer.Encode(false)

	_ = this.service.sendMessage(this.addr, buffer)
	if this.sessionId != 0 {
		this.sendMutex.Lock()
		defer this.sendMutex.Unlock()
		this.sendedPackageMap[buffer.SN] = buffer
	}
}
