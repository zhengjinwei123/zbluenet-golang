package net

import "sync"

// socket read and write processor
type TcpReactor struct {
	tcp_conn *TcpConnection
	tcp_session TcpSession

	ReadMessageQueue chan []byte
	writeMessageQueue chan *NetMessage

	wg sync.WaitGroup
}

func NewTcpReactor(session TcpSession) *TcpReactor {
	return &TcpReactor{
		tcp_session:       session,
	}
}