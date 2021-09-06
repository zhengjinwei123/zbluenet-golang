package net

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"
)
//<property name="MaxReadMsgSize">65536</property>  <!-- 接收包的最大长度 -->
//30     <property name="ReadMsgQueueSize">10240</property>  <!-- 接收包的缓存队列长度 -->
//31     <property name="ReadTimeOut">600</property>  <!-- 接收包的超时(second) -->
//32     <property name="MaxWriteMsgSize">65536</property> <!-- 发送包的最大长度 -->
//33     <property name="WriteMsgQueueSize">10240</property> <!-- 发送包的缓存队列长度  -->
//34     <property name="WriteTimeOut">600</property>  <!-- 发送包的超时(second) -->


const MESSAGE_HEAD_LENGTH = 4 // msg_id(2) + message_length(2)

const DEFAULT_READ_TIMEOUT_SEC = 60 // 60秒没有数据到来， 超时处理
const DEFAULT_WRITE_TIMEOUT_SEC = 10

var ErrorMessageNotEnough = errors.New("message not enough")


type TcpConnection struct {
	conn   net.Conn
	remote_addr string
	local_addr string

	messageId []byte
	messageLength []byte

	reactor *TcpReactor

	writeTimeout time.Duration
	readTimeout time.Duration
}

func NewConnection(conn net.Conn, reactor *TcpReactor, write_timeout_sec int32, read_timeout_sec int32) *TcpConnection {
	c := &TcpConnection{
		conn: conn,
		remote_addr: conn.RemoteAddr().String(),
		local_addr: conn.LocalAddr().String(),
		messageId: make([]byte, 2, 2),
		messageLength: make([]byte, 2, 2),
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

func (this *TcpConnection) WriteLoop() {
	//max_size := 65536 // 发送包的最大长度
	//
	//write_buf := make([]byte, max_size)
	//head_buf := make([]byte, MESSAGE_HEAD_LEN)
	//data_buf := make([]byte, max_size - MESSAGE_HEAD_LEN)
	//
	//for {
	//	select {
	//	case msg <- this.reactor.writeMessageQueue:
	//		length, data, err :=
	//	}
	//}
}

func (this *TcpConnection) ReadLoop() {
	reader := bufio.NewReader(this.conn)

	for {

		err := this.conn.SetReadDeadline(time.Now().Add(this.readTimeout))
		if err != nil {
			continue
		}

		if err := this.read(reader); err != nil {

			if err != ErrorMessageNotEnough {
				goto exit
			}

			continue
		}
	}
exit:
	this.reactor.wg.Done()
}



func (this *TcpConnection) read(reader *bufio.Reader) error {
	readableBytes := reader.Size()
	if readableBytes < MESSAGE_HEAD_LENGTH {
		return ErrorMessageNotEnough
	}

	buf, err := reader.Peek(MESSAGE_HEAD_LENGTH)
	if err != nil {
		return err
	}

	messageLength := DecodeUint16(buf[2:])

	packLength := int(messageLength) + MESSAGE_HEAD_LENGTH
	if readableBytes < packLength {
		return ErrorMessageNotEnough
	}

	r := io.Reader(reader)

	packBuf := make([]byte, packLength, packLength)

	_, err = io.ReadFull(r, packBuf)
	if err != nil {
		return err
	}

	messageId := DecodeUint16(buf[0:2])

	this.reactor.recvPacket(messageId, messageLength, packBuf)

	return nil
}
