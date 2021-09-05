package net

type TcpSession interface {
	OnInit(*TcpReactor)
	OnRecvMessage(*MessageHead, []byte)
	OnClose()
}