package net

import (
	"net"
	"time"
)

const (
	CONNECTION_STATE_INIT = 1
	CONNECTION_STATE_CONNECTED = 2
	CONNECTION_STATE_DISCONNECTED = 3
)

type TcpConnection struct {
	conn   net.Conn
	remote_addr string
	local_addr string

	msgLength []byte

	reactor *TcpReactor

	writeTimeout time.Duration
	readTimeout time.Duration
}

func NewConnection(conn net.Conn, reactor *TcpReactor) *TcpConnection {

}
