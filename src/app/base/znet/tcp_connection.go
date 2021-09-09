package znet

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
	"github.com/golang/protobuf/proto"
)

const MAX_READ_MESSAGE_SIZE = 65535

const MESSAGE_HEAD_LENGTH = 4 // msg_id(2) + message_length(2)

const DEFAULT_READ_TIMEOUT_SEC = 5 // 60秒没有数据到来， 超时处理
const DEFAULT_WRITE_TIMEOUT_SEC = 10

var ErrorMessageOverflow = errors.New("message over flow")
var ErrorInvalidMessageType = errors.New("invalid message type")
var ErrorMessageWriteOverFlow = errors.New("message write overflow")


type TcpConnection struct {
	id     uint32
	conn   net.Conn
	remote_addr string
	local_addr string

	messageHead []byte

	reactor *TcpReactor

	writeTimeout time.Duration
	readTimeout time.Duration
}

func NewConnection(id uint32, conn net.Conn, reactor *TcpReactor, write_timeout_sec int32, read_timeout_sec int32) *TcpConnection {

	// only work on linux
	if s, ok := conn.(*net.TCPConn);ok {
		f, err := s.File()
		if err == nil {
			id = uint32(f.Fd())
		}
	}

	c := &TcpConnection{
		id: id,
		conn: conn,
		remote_addr: conn.RemoteAddr().String(),
		local_addr: conn.LocalAddr().String(),
		messageHead: make([]byte, MESSAGE_HEAD_LENGTH, MESSAGE_HEAD_LENGTH),
		reactor: reactor,
	}

	if write_timeout_sec > 0 {
		c.writeTimeout = time.Duration(write_timeout_sec) * time.Second
	} else {
		c.writeTimeout = time.Duration(DEFAULT_WRITE_TIMEOUT_SEC) * time.Second
	}

	if read_timeout_sec > 0 {
		c.readTimeout = time.Duration(read_timeout_sec) * time.Second
	} else {
		c.readTimeout = time.Duration(DEFAULT_READ_TIMEOUT_SEC) * time.Second
	}


	return c
}

func (this *TcpConnection) Close() error {
	err := this.conn.Close()
	return err
}

func (this *TcpConnection) WriteLoop() {
	max_size := 65535 // 发送包的最大长度

	write_buf := make([]byte, max_size)
	head_buf := make([]byte, MESSAGE_HEAD_LEN)
	data_buf := make([]byte, max_size - MESSAGE_HEAD_LEN)

	for this.reactor != nil {
		select {
		case msg := <- this.reactor.writeMessageQueue:
			length, data, err := Marshal(msg, max_size, head_buf, data_buf)
			if err == nil {
				index := 0
				copy(write_buf, head_buf)
				copy(write_buf[MESSAGE_HEAD_LEN:], data)
				index += length

				for more := true; more; {
					select {
					case msg := <- this.reactor.writeMessageQueue:
						length, data, err := Marshal(msg, max_size, head_buf, data_buf)
						if err == nil {
							if index + int(length) <= max_size {
								copy(write_buf[index:], head_buf)
								copy(write_buf[index+MESSAGE_HEAD_LEN:], data)
								index += length
							} else {
								if !this.write(write_buf[:index]) {
									goto exit;
								}
								index = 0
								copy(write_buf, head_buf)
								copy(write_buf[MESSAGE_HEAD_LEN:], data)
								index += length
							}
						}
					case <-this.reactor.CloseChan:
						goto exit
					default:
						more= false
					}
				}

				if !this.write(write_buf[:index]) {
					goto exit
				}
			}
		case <-this.reactor.CloseChan:
			goto exit
		}
	}
exit:
	fmt.Printf("WriteLoop exit \n")
	this.reactor.wg.Done()
}

func marshalWithBytes(pb proto.Message, data []byte) ([]byte, error) {
	p := proto.NewBuffer(data)
	err := p.Marshal(pb)

	return p.Bytes(), err
}

func Marshal(message *NetMessage, max_size int, head_buff, data_buff []byte) (lengRet int, dataRet []byte, errRet error) {

	var data []byte
	switch v := message.DataSend.(type) {
	case []byte:
		data = v;
	case proto.Message:
		if mdata, err := marshalWithBytes(v, data_buff[0:0]); err == nil {
			data = mdata
		} else {
			return 0, nil, err
		}
	default:
		return 0, nil, ErrorInvalidMessageType
	}

	msgSize := len(data)

	if (msgSize + MESSAGE_HEAD_LEN) > max_size {
		return 0, nil, ErrorMessageWriteOverFlow
	}

	EncodeHead(head_buff, message.MessageId, uint16(msgSize))

	return msgSize + MESSAGE_HEAD_LEN, data, nil
}

func (this *TcpConnection) ReadLoop() {
	reader := bufio.NewReader(this.conn)
	for {
		err := this.conn.SetReadDeadline(time.Now().Add(this.readTimeout))
		if err != nil {
			continue
		}

		if err := this.read(reader); err != nil {
			fmt.Printf("read loop err: %d | %s \n", this.id, err.Error())
			goto exit
		}
	}
exit:
	fmt.Printf("ReadLoop exit \n")
	this.reactor.Stop()
	this.reactor.wg.Done()
}

func (this *TcpConnection) write(msg []byte) bool {
	_, err := this.conn.Write(msg)
	msg = nil
	if err != nil {
		return false
	}

	return true
}

func (this *TcpConnection) read(r io.Reader) error {
	readSize, err := io.ReadFull(r, this.messageHead)
	if err != nil {
		fmt.Printf("read error: %s (%d) \n", err.Error(), readSize)
		return err
	}

	messageId := DecodeUint16(this.messageHead[:2])
	messageSize := DecodeUint16(this.messageHead[2:])

	if messageSize + MESSAGE_HEAD_LENGTH > MAX_READ_MESSAGE_SIZE {
		return ErrorMessageOverflow
	}

	dataBuf := make([]byte, messageSize, messageSize)
	_, err = io.ReadFull(r, dataBuf)
	if err != nil {
		fmt.Printf("read error222: %s \n", err.Error())
		return err
	}

	this.reactor.recvPacket(messageId, messageSize, dataBuf)

	return nil
}


//func (this *TcpConnection) read1(reader *bufio.Reader) error {
//	readableBytes := reader.Size()
//	if readableBytes < MESSAGE_HEAD_LENGTH {
//		return ErrorMessageNotEnough
//	}
//
//	buf, err := reader.Peek(MESSAGE_HEAD_LENGTH)
//	if err != nil {
//		return err
//	}
//
//	messageLength := DecodeUint16(buf[2:])
//
//	packLength := int(messageLength) + MESSAGE_HEAD_LENGTH
//	if readableBytes < packLength {
//		return ErrorMessageNotEnough
//	}
//
//	r := io.Reader(reader)
//
//	packBuf := make([]byte, packLength, packLength)
//
//	_, err = io.ReadFull(r, packBuf)
//	if err != nil {
//		return err
//	}
//
//	messageId := DecodeUint16(buf[0:2])
//
//	this.reactor.recvPacket(messageId, messageLength, packBuf)
//
//	return nil
//}
