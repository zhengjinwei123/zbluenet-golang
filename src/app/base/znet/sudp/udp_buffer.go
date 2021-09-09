package sudp

import (
	"app/base/znet"
)

// 安全UDP 消息包头

const UDP_MESSAGE_TYPE_ACK = 1
const UDP_MESSAGE_TYPE_MESSAGE = 0

const UDP_PACKET_HEAD_SIZE = 21

type udpBuffer struct {
	DataSize 			  uint16 // 报文大小 字节 2
	SessionId 			  uint32 // 会话ID 用来标识身份 4
	SN        			  uint32 // 包体的序号， 用来保证顺序 4
	Time                  uint64 // 发送时间戳, 用来控制超时重传 8
	MessageType           uint8  // 协议类型 ACK ， data 1
	MessageId             uint16 // 消息ID 2
	Data                  []byte // 业务报文数据

	reSendCount           int // 重发次数， 不是业务数据
	buffer                []byte // 最终受到的数据 或者 发送的数据
	isFull                bool
}


// 构建请求报文
func NewRequestUdpBuffer(session_id uint32, sn uint32, message_type uint8, message_id uint16, data []byte) *udpBuffer {
	return &udpBuffer{
		DataSize: uint16(len(data)),
		Data: data,
		SessionId: session_id,
		SN: sn,
		MessageType: message_type,
		MessageId: message_id,
	}
}

// 构建接收到的报文
func NewRecvUpdBuffer(buffer []byte) *udpBuffer {
	tmp := &udpBuffer{
		buffer: buffer,
		isFull: false,
	}

	tmp.isFull = tmp.DeCode()

	return tmp
}

// 构建ack 报文
func NewAckUpdBuffer(buffer *udpBuffer) *udpBuffer {
	tmp := &udpBuffer{
		DataSize: 0,
		SessionId: buffer.SessionId,
		SN: buffer.SN,
		Time: buffer.Time,
		MessageType: UDP_MESSAGE_TYPE_ACK,
		MessageId: buffer.MessageId,
	}
	tmp.buffer = tmp.Encode(true)

	return tmp
}


func (this *udpBuffer) Encode(is_ack bool) []byte {
	if is_ack == true {
		this.DataSize = 0
	}
	data := make([]byte, UDP_PACKET_HEAD_SIZE + this.DataSize)

	dataSizeBuf := make([]byte, 2, 2)
	znet.EncodeUint16(this.DataSize, dataSizeBuf)

	sessionIdBuf := make([]byte, 4, 4)
	znet.EncodeUint32(this.SessionId, sessionIdBuf)

	snBuf := make([]byte, 4, 4)
	znet.EncodeUint32(this.SN, snBuf)

	timeBuf := make([]byte, 8, 8)
	znet.EncodeUint64(this.Time, timeBuf)

	messageTypeBuf := make([]byte, 1, 1)
	znet.EncodeUint8(this.MessageType, messageTypeBuf)

	messageIdBuf := make([]byte, 2, 2)
	znet.EncodeUint16(this.MessageId, messageIdBuf)


	index := 0
	copy(data[index:], dataSizeBuf)

	index += 2
	copy(data[index:], sessionIdBuf)

	index += 4
	copy(data[index:], snBuf)

	index += 4
	copy(data[index:], timeBuf)

	index += 8
	copy(data[index:], messageTypeBuf)

	index += 1
	copy(data[index:], messageIdBuf)

	if is_ack == false {
		index += 2
		copy(data[index:], data)
	}

	this.buffer = data

	return data
}

func (this *udpBuffer) DeCode() bool {
	index := 0

	bufferSize := len(this.buffer)

	if bufferSize >= 2 {
		this.DataSize = znet.DecodeUint16(this.buffer[index:])
		index += 2

		if bufferSize != (int(this.DataSize) + UDP_PACKET_HEAD_SIZE) {
			return false
		}
	} else {
		return false
	}

	this.SessionId = znet.DecodeUint32(this.buffer[index:])
	index += 4

	this.SN = znet.DecodeUint32(this.buffer[index:])
	index += 4

	this.Time = znet.DecodeUint64(this.buffer[index:])
	index += 8

	this.MessageType = znet.DecodeUint8(this.buffer[index:])
	index += 1

	this.MessageId = znet.DecodeUint16(this.buffer[index:])


	if this.MessageType != UDP_MESSAGE_TYPE_ACK {
		index += 2

		data := make([]byte, this.DataSize, this.DataSize)
		// 将 buffer 里 剩余的数据复制到data, 得到业务数据
		copy(data, this.buffer[index:])
	}

	return true
}


