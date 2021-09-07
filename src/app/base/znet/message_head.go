package znet

type MessageHead struct {
	MessageId uint16
	MessageLength uint16
}

type NetMessage struct {
	Fd        int
	MessageId uint16
	Data interface{}
}

const MESSAGE_HEAD_LEN = 4


func NewNetMessage(messageId uint16, fd int, data []byte) *NetMessage {
	return &NetMessage{
		Fd: fd,
		MessageId: messageId,
		Data: data,
	}
}

func DecodeHead(buf []byte, msgHead *MessageHead) bool {
	if (len(buf) < 4) {
		return false
	}

	msgHead.MessageId = DecodeUint16(buf[0:])
	msgHead.MessageLength = DecodeUint16(buf[2:])

	return true
}

func EncodeHead(buf []byte, message_id uint16, message_size uint16) bool {
	if (message_size < 4) {
		return false
	}

	EncodeUint16(message_id, buf[0:])
	EncodeUint16(message_size, buf[2:])

	return true
}

// big endian
func DecodeUint16(data []byte) uint16 {
	return (uint16(data[0]) << 8) | (uint16(data[1]))
}

// big endian
func EncodeUint16(n uint16, buf []byte) {
	buf[1] = byte(n & 0xff)
	buf[0] = byte((n >> 8) & 0xff)
}

// big endian
func DecodeUint32(data []byte) uint32 {
	return (uint32(data[0]) << 24) |
		(uint32(data[1]) << 16) |
		(uint32(data[2]) << 8) |
		(uint32(data[3]))
}

// big endian
func EncodeUint32(n uint32, buf []byte)  {
	buf[3] = byte(n & 0xff)
	buf[2] = byte((n >> 8) & 0xff)
	buf[1] = byte((n >> 16) & 0xff)
	buf[0] = byte((n >> 24) & 0xff)
}

// big endian
func DecodeUint64(data []byte) uint32 {
	return (uint32(data[0]) << 56) |
		(uint32(data[1]) << 48) |
		(uint32(data[2]) << 40) |
		(uint32(data[3]) << 32) |
		(uint32(data[4]) << 24) |
		(uint32(data[5]) << 16) |
		(uint32(data[6]) << 8) |
		(uint32(data[7]))
}

//big endian
func EncodeUint64(n uint64, buf []byte) {
	buf[7] = byte(n & 0xff)
	buf[6] = byte((n >> 8) & 0xFF)
	buf[5] = byte((n >> 16) & 0xFF)
	buf[4] = byte((n >> 24) & 0xFF)
	buf[3] = byte((n >> 32) & 0xFF)
	buf[2] = byte((n >> 40) & 0xFF)
	buf[1] = byte((n >> 48) & 0xFF)
	buf[0] = byte((n >> 56) & 0xFF)
}