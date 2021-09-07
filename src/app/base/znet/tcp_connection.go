package znet

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"
	"github.com/golang/protobuf/proto"
)
//<property name="MaxReadMsgSize">65536</property>  <!-- 接收包的最大长度 -->
//30     <property name="ReadMsgQueueSize">10240</property>  <!-- 接收包的缓存队列长度 -->
//31     <property name="ReadTimeOut">600</property>  <!-- 接收包的超时(second) -->
//32     <property name="MaxWriteMsgSize">65536</property> <!-- 发送包的最大长度 -->
//33     <property name="WriteMsgQueueSize">10240</property> <!-- 发送包的缓存队列长度  -->
//34     <property name="WriteTimeOut">600</property>  <!-- 发送包的超时(second) -->

const MAX_READ_MESSAGE_SIZE = 65535

const MESSAGE_HEAD_LENGTH = 4 // msg_id(2) + message_length(2)

const DEFAULT_READ_TIMEOUT_SEC = 60 // 60秒没有数据到来， 超时处理
const DEFAULT_WRITE_TIMEOUT_SEC = 10

var ErrorMessageNotEnough = errors.New("message not enough")
var ErrorMessageOverflow = errors.New("message over flow")
var ErrorInvalidMessageType = errors.New("invalid message type")
var ErrorMessaegWriteOverFlow = errors.New("message write overflow")


type TcpConnection struct {
	fd     int
	conn   net.Conn
	remote_addr string
	local_addr string

	messageHead []byte

	reactor *TcpReactor

	writeTimeout time.Duration
	readTimeout time.Duration
}

func NewConnection(conn net.Conn, reactor *TcpReactor, write_timeout_sec int32, read_timeout_sec int32) *TcpConnection {

	s, _ := conn.(*net.TCPConn);
	f, _ := s.File()
	var fd int = int(f.Fd())

	c := &TcpConnection{
		fd: fd,
		conn: conn,
		remote_addr: conn.RemoteAddr().String(),
		local_addr: conn.LocalAddr().String(),
		messageHead: make([]byte, MESSAGE_HEAD_LENGTH, MESSAGE_HEAD_LENGTH),
		reactor: nil,
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
	this.reactor = nil
	return err
}

func (this *TcpConnection) WriteLoop() {
	max_size := 65535 // 发送包的最大长度

	write_buf := make([]byte, max_size)
	head_buf := make([]byte, MESSAGE_HEAD_LEN)
	data_buf := make([]byte, max_size - MESSAGE_HEAD_LEN)

	for {
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
	this.reactor.wg.Done()
}

func marshalWithBytes(pb proto.Message, data []byte) ([]byte, error) {
	p := proto.NewBuffer(data)
	err := p.Marshal(pb)

	return p.Bytes(), err
}

func Marshal(message *NetMessage, max_size int, head_buff, data_buff []byte) (lengRet int, dataRet []byte, errRet error) {

	var data []byte
	switch v := message.Data.(type) {
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

	length := len(data) + MESSAGE_HEAD_LEN

	if length > max_size {
		return 0, nil, ErrorMessaegWriteOverFlow
	}

	EncodeHead(head_buff, message.MessageId, uint16(length))

	return length, data, nil
}

func (this *TcpConnection) ReadLoop() {
	reader := bufio.NewReader(this.conn)

	for {
		err := this.conn.SetReadDeadline(time.Now().Add(this.readTimeout))
		if err != nil {
			continue
		}

		if err := this.read(reader); err != nil {
			goto exit
		}
	}
exit:
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
	_, err := io.ReadFull(r, this.messageHead)
	if err != nil {
		return err
	}

	messageId := DecodeUint16(this.messageHead[0:2])
	messageSize := DecodeUint16(this.messageHead[2:])

	if messageSize < MESSAGE_HEAD_LENGTH {
		return ErrorMessageNotEnough
	}
	if messageSize > MAX_READ_MESSAGE_SIZE {
		return ErrorMessageOverflow
	}

	dataBuf := make([]byte, messageSize - MESSAGE_HEAD_LEN, messageSize - MESSAGE_HEAD_LEN)
	_, err = io.ReadFull(r, dataBuf)
	if err != nil {
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
